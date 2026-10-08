package main

import (
	"context"
	"crypto/sha256"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatdb/internal/api"
	"chatdb/internal/auth"
	"chatdb/internal/config"
	"chatdb/internal/engine"
	"chatdb/internal/mcp"
	"chatdb/internal/migrate"
	"chatdb/internal/security"
	"chatdb/internal/store"
	"chatdb/web"

	"github.com/joho/godotenv"
)

func main() {
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		log.Fatalf("config path: %v", err)
	}

	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := store.OpenMetadataDB(cfg.Metadata.Path)
	if err != nil {
		log.Fatalf("open metadata db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := migrate.Bootstrap(ctx, db); err != nil {
		cancel()
		log.Fatalf("bootstrap metadata: %v", err)
	}
	if err := migrate.Upgrade(ctx, db); err != nil {
		cancel()
		log.Fatalf("upgrade metadata: %v", err)
	}
	cancel()
	log.Printf("metadata sqlite ready: %s", cfg.Metadata.Path)

	_ = godotenv.Load("../.env", ".env")

	// Respect PORT and HOST from environment or .env
	if envPort := os.Getenv("PORT"); envPort != "" {
		host := "127.0.0.1"
		if envHost := os.Getenv("HOST"); envHost != "" {
			host = envHost
		} else if h, _, err := net.SplitHostPort(cfg.Listen); err == nil && h != "" {
			host = h
		}
		cfg.Listen = net.JoinHostPort(host, envPort)
	}

	authApiKey := os.Getenv("AUTH_API_KEY")
	var cryptoKey []byte
	if authApiKey != "" {
		hash := sha256.Sum256([]byte(authApiKey))
		cryptoKey = hash[:]
		log.Println("Using AUTH_API_KEY for encryption")
	} else {
		cryptoKey = []byte(cfg.AppKey)
	}

	crypter, err := security.NewCrypter(cryptoKey)
	if err != nil {
		log.Fatalf("crypter: %v", err)
	}

	st := store.New(db)

	// CLI subcommand: ./chatdb mcp
	// Runs Model Context Protocol (MCP) server over standard input/output.
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		conn, err := st.FirstConnection(context.Background())
		if err != nil {
			log.Fatalf("mcp: no database connection found: %v", err)
		}
		password, err := crypter.Decrypt(conn.ReadPassword)
		if err != nil {
			log.Fatalf("mcp decrypt password: %v", err)
		}
		var sshPass, sshKey string
		if conn.UseSSH {
			if conn.SshPassword != "" {
				sshPass, _ = crypter.Decrypt(conn.SshPassword)
			}
			if conn.SshKey != "" {
				sshKey, _ = crypter.Decrypt(conn.SshKey)
			}
		}
		var tunnel *engine.SSHTunnel
		if conn.UseSSH {
			tunnel, err = engine.NewSSHTunnel(conn.SshHost, conn.SshPort, conn.SshUser, sshPass, sshKey, conn.Host, conn.Port)
			if err != nil {
				log.Fatalf("mcp ssh tunnel: %v", err)
			}
		}
		var eng engine.Engine
		switch config.Driver(conn.Driver) {
		case config.DriverMySQL:
			eng, err = engine.OpenMySQL(tunnel, conn.Host, conn.Port, conn.ReadUsername, password, conn.Database)
		default:
			eng, err = engine.OpenPostgres(tunnel, conn.Host, conn.Port, conn.ReadUsername, password, conn.Database, conn.SslMode)
		}
		if err != nil {
			log.Fatalf("mcp open engine: %v", err)
		}
		defer eng.Close()

		log.Printf("VueNiceDB MCP server ready on stdio for database '%s'", conn.Database)
		if err := mcp.RunStdio(context.Background(), eng, conn.Database); err != nil {
			log.Fatalf("mcp run: %v", err)
		}
		return
	}

	srv := &api.Server{
		Cfg:     cfg,
		Store:   st,
		Crypter: crypter,
		JWT:     auth.NewIssuer(cfg.JWTSecret),
		Pools:   engine.NewManager(),
		Static:  web.SPA(),
	}

	httpSrv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("chatdb listening on http://%s", cfg.Listen)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	runDesktopApp(httpSrv.Handler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	srv.Pools.CloseAll()
}

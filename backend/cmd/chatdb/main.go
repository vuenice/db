package main

import (
	"context"
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatdb/internal/api"
	"chatdb/internal/auth"
	"chatdb/internal/config"
	"chatdb/internal/engine"
	"chatdb/internal/migrate"
	"chatdb/internal/security"
	"chatdb/internal/store"
	"chatdb/web"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
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

	wailsApp := application.New(application.Options{
		Name: "chatdb",
		Description: "chatdb Wails application",
		Assets: application.AssetOptions{
			Handler: httpSrv.Handler,
		},
	})
	
	go func() {
		err := wailsApp.Run()
		if err != nil {
			log.Printf("wails app run: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	srv.Pools.CloseAll()
}

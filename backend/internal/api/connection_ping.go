package api

import (
	"context"

	"chatdb/internal/config"
	"chatdb/internal/engine"
)

// pingTargetDB verifies read credentials against the given database (physical DB name).
func pingTargetDB(ctx context.Context, driver, host string, port int, database, sslMode, readUser, readPass string,
	useSSH bool, sshHost string, sshPort int, sshUser, sshPass, sshKey string) error {
	var tunnel *engine.SSHTunnel
	if useSSH {
		var err error
		tunnel, err = engine.NewSSHTunnel(sshHost, sshPort, sshUser, sshPass, sshKey, host, port)
		if err != nil {
			return err
		}
	}

	switch config.Driver(driver) {
	case config.DriverMySQL:
		e, err := engine.OpenMySQL(tunnel, host, port, readUser, readPass, database)
		if err != nil {
			return err
		}
		defer e.Close()
		return e.Ping(ctx)
	default:
		e, err := engine.OpenPostgres(tunnel, host, port, readUser, readPass, database, sslMode)
		if err != nil {
			return err
		}
		defer e.Close()
		return e.Ping(ctx)
	}
}

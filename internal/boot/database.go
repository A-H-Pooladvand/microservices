package boot

import (
	"context"
	"po/pkg/logger"
	"po/pkg/mysql"
)

type DB struct {
	client *mysql.Client
}

func (db *DB) Boot(ctx context.Context) error {
	client := mysql.New(
		mysql.NewEnvConfig(),
	)

	db.client = client

	select {
	case <-ctx.Done():
		logger.Error()

		return db.Shutdown()
	}
}

func (db *DB) Shutdown() error {
	logger.Info("Shutting down the database...")
	return db.client.Close()
}

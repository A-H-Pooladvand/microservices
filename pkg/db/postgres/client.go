package postgres

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// New creates a new postgres client and returns a pointer to the gorm.DB instance.
func New(config Config) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	dialect := postgres.Open(config.DSN())

	db, err := gorm.Open(dialect, &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	connection, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	for {
		if err = connection.Ping(); err == nil {
			break
		}

		select {
		case <-time.After(500 * time.Millisecond):
			continue
		case <-ctx.Done():
			defer connection.Close()
			return nil, errors.New("unable to connect to postgres client, context deadline exceeded")
		}
	}

	return db, nil
}

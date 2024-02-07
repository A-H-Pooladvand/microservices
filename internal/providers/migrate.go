package providers

import (
	"context"
	"po/internal/db"
	"po/migrations"
)

func Migrations(ctx context.Context) error {
	return db.New().AutoMigrate(migrations.Get()...)
}

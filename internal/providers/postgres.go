package providers

import (
	"context"
	"po/internal/app"
	"po/internal/db"
)

func Postgres(_ context.Context) error {
	client := db.New()

	a := app.Get()
	a.Postgres = client

	return nil
}

package providers

import (
	"po/internal/app"
	"po/internal/db"
)

func Postgres(ctx app.Context) error {
	client := db.New()

	a := app.Get()
	a.Postgres = client

	return nil
}

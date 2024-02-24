package configs

import (
	"context"
	"github.com/caarlos0/env/v10"
	"os"
	"po/internal/vault"
	"strings"
)

func parse(path string, v any) error {
	environment := strings.ToLower(os.Getenv("APP_ENV"))

	prod := environment == "prod" || environment == "production"

	if !prod {
		return env.Parse(v)
	}

	client, err := vault.New()

	if err != nil {
		return err
	}

	return client.Parse(context.Background(), path, v)
}

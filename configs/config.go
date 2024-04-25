package configs

import (
	"context"
	"github.com/caarlos0/env/v10"
	"po/internal/app"
	"po/pkg/vault"
	"time"
)

func Parse(path string, v any, client *vault.Client) error {
	if app.Local() {
		return env.Parse(v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return client.Parse(ctx, path, v)
}

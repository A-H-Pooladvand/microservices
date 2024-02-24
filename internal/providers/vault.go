package providers

import (
	"context"
	"errors"
	"po/internal/app"
	"po/internal/vault"
)

func Vault(_ context.Context) error {
	if app.Production() {
		config, err := vault.NewConfig()

		if err != nil {
			return err
		}

		if config.Empty() {
			return errors.New(`Due to production constraints, environment variables are unavailable.
Please specify the Config configuration.`)
		}
	}

	return nil
}

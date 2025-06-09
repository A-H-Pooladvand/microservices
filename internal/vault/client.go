package vault

import (
	"context"
	"errors"
	"github.com/a-h-pooladvand/microservices/internal/app"
	"github.com/a-h-pooladvand/microservices/pkg/vault"
	"go.uber.org/fx"
)

// New vault implementation
func New(config Config) (*vault.Client, error) {
	c := vault.NewConfig(
		config.Address,
		config.RoleID,
		config.SecretId,
	)

	return vault.New(c)
}

func Provide(lc fx.Lifecycle, config Config) (*vault.Client, error) {
	if app.Local() {
		return nil, nil
	}

	if config.Empty() {
		return nil, errors.New(`Due to production constraints, vault environment variables are unavailable.
Please specify the Config configuration.`)
	}

	client, err := New(Config{
		Address:   config.Address,
		RoleID:    config.RoleID,
		SecretId:  config.SecretId,
		MountPath: config.MountPath,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
	})

	return client, nil
}

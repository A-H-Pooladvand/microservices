package vault

import (
	"context"
	"errors"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/vault"
	"go.uber.org/fx"
)

// New vault implementation
func New(cfg Config) (*vault.Client, error) {
	c := vault.NewConfig(
		cfg.Address,
		cfg.RoleID,
		cfg.SecretId,
	)

	return vault.New(c)
}

func Provide(lc fx.Lifecycle, cfg Config, appCfg *config.Config) (*vault.Client, error) {
	// Skip vault in local/dev environment
	if appCfg.App.Dev() {
		return nil, nil
	}

	if cfg.Empty() {
		return nil, errors.New(`Due to production constraints, vault environment variables are unavailable.
Please specify the Config configuration.`)
	}

	client, err := New(Config{
		Address:   cfg.Address,
		RoleID:    cfg.RoleID,
		SecretId:  cfg.SecretId,
		MountPath: cfg.MountPath,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
	})

	return client, nil
}

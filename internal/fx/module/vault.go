package module

import (
	"context"

	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/vault"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Vault = fx.Module("vault", fx.Provide(
	NewVault,
))

// NewVault creates a new Vault client with observability.
// Returns nil if Vault is not configured (local development).
func NewVault(lc fx.Lifecycle, cfg *config.Config) (*vault.Client, error) {
	// Skip if vault is not configured
	if cfg.Vault.Address == "" {
		return nil, nil
	}

	vaultCfg := vault.NewConfig(
		cfg.Vault.Address,
		cfg.Vault.RoleID,
		cfg.Vault.SecretID,
	)

	if vaultCfg.Empty() {
		return nil, nil
	}

	client, err := vault.New(vaultCfg, zap.L())
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Health(ctx)
		},
	})

	return client, nil
}

package module

import (
	"github.com/a-h-pooladvand/microservices/pkg/vault"
	"go.uber.org/fx"
)

var Vault = fx.Module("vault", fx.Provide(
	vault.NewConfig,
	NewVault,
))

func NewVault(lc fx.Lifecycle, config vault.Config) (*vault.Client, error) {
	// Todo: implement if env == "local" return nil, nil
	return nil, nil

	//	if config.Empty() {
	//		return nil, errors.New(`Due to production constraints, vault environment variables are unavailable.
	//Please specify the Config configuration.`)
	//	}
	//
	//	client, err := vault.New(vault.Config{
	//		Address:   config.Address,
	//		RoleID:    config.RoleID,
	//		SecretId:  config.SecretId,
	//		MountPath: config.MountPath,
	//	})
	//
	//	lc.Append(fx.Hook{
	//		OnStart: func(ctx context.Context) error {
	//			return err
	//		},
	//	})
	//
	//	return client, nil
}

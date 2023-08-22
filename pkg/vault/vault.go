package vault

import (
	"context"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"po/pkg/env"
)

type Vault struct{}

func (v Vault) Boot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
		config := vault.DefaultConfig()
		config.Address = env.Get("VAULT_ADDRESS", "http://127.0.0.1:8200")

		client, err := vault.NewClient(config)

		if err != nil {
			return errors.New(
				fmt.Sprintf("unable to initialize Vault client: %v", err),
			)
		}

		client.SetToken(env.Get("VAULT_TOKEN", "hvs.6ja7wvLDcakpK2EMh6WTUMeo"))

		return nil
	}
}

package vault

import (
	"context"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"os"
)

type Vault struct{}

func (v Vault) Boot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
		config := vault.DefaultConfig()
		config.Address = os.Getenv("VAULT_ADDRESS")

		client, err := vault.NewClient(config)

		if err != nil {
			return errors.New(
				fmt.Sprintf("unable to initialize Vault client: %v", err),
			)
		}

		client.SetToken("hvs.6ja7wvLDcakpK2EMh6WTUMeo")

		return nil
	}
}

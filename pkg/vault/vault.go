package vault

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"log"
)

type Vault struct {
	Client *vault.Client
	Config Config
}

func New(config Config) (*Vault, error) {
	ctx, cancelContextFunc := context.WithCancel(context.Background())
	defer cancelContextFunc()

	cfg := vault.DefaultConfig() // modify for more granular configuration
	cfg.Address = config.Address

	client, err := vault.NewClient(cfg)

	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault Client: %w", err)
	}

	v := &Vault{
		Client: client,
		Config: config,
	}

	token, err := v.login(ctx)

	fmt.Println("---------------")
	fmt.Println(token.TokenID())

	if err != nil {
		return nil, fmt.Errorf("vault login error: %w", err)
	}

	log.Println("connecting to vault: success!")

	return v, nil
}

func (v *Vault) login(ctx context.Context) (*vault.Secret, error) {
	log.Printf("logging in to vault with approle auth; role id: %s", v.Config.ApproleRoleID)

	approleSecretID := &approle.SecretID{
		FromFile: v.Config.ApproleSecretIDFile,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		v.Config.ApproleRoleID,
		approleSecretID,
		approle.WithWrappingToken(), // only required if the SecretID is response-wrapped
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize approle authentication method: %w", err)
	}

	authInfo, err := v.Client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using approle auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no approle info was returned after login")
	}

	log.Println("logging in to vault with approle auth: success!")

	return authInfo, nil
}

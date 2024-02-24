package vault

import (
	"context"
	"encoding/json"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

type Client struct {
	config Config
	conn   *vault.Client
}

func New(config Config) (*Client, error) {
	c := vault.DefaultConfig()
	c.Address = config.Address

	vc, err := vault.NewClient(c)

	if err != nil {
		return nil, err
	}

	client := &Client{
		config: config,
		conn:   vc,
	}

	client.login(context.Background())

	return client, nil
}

func (c *Client) login(ctx context.Context) (*vault.Secret, error) {
	secretID := &auth.SecretID{
		FromString: c.config.SecretID,
	}

	appRoleAuth, err := auth.NewAppRoleAuth(
		c.config.RoleID,
		secretID,
		//auth.WithWrappingToken(),
	)

	if err != nil {
		return nil, err
	}

	return c.conn.Auth().Login(ctx, appRoleAuth)
}

func (c *Client) Get(ctx context.Context, path string) (*vault.KVSecret, error) {
	return c.conn.KVv2("secret").Get(ctx, path)
}

func (c *Client) Parse(ctx context.Context, path string, v any) error {
	secret, err := c.Get(ctx, path)

	if err != nil {
		return err
	}

	b, err := json.Marshal(secret.Data)

	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"po/internal/app"
	"po/pkg/vault"
)

func Vault(ctx app.Context) error {
	config := vault.NewConfig(
		"http://127.0.0.1:8200",
		"demo-web-webserver",
		"/tmp/secret",
		"api-key",
		"kv-v2",
		"api-key-field",
		"database/creds/dev-readonly",
	)

	v, err := vault.New(config)

	//
	log.Println("getting temporary database credentials from vault")

	lease, err := v.Client.Logical().ReadWithContext(ctx, v.Config.DatabaseCredentialsPath)
	if err != nil {
		return fmt.Errorf("unable to read secret: %w", err)
	}

	b, err := json.Marshal(lease.Data)

	if err != nil {
		return fmt.Errorf("malformed credentials returned: %w", err)
	}

	var credentials DatabaseCredentials

	if err := json.Unmarshal(b, &credentials); err != nil {
		return fmt.Errorf("unable to unmarshal credentials: %w", err)
	}

	log.Println("getting temporary database credentials from vault: success!")

	return err
}

// DatabaseCredentials is a set of dynamic credentials retrieved from Vault
type DatabaseCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

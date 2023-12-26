package cfg

import (
	"github.com/caarlos0/env/v10"
)

type Vault struct {
	Host              string `env:"VAULT_HOST" envDefault:"127.0.0.1"`
	Port              string `env:"VAULT_PORT" envDefault:"8200"`
	RoleId            string `env:"VAULT_APPROLE_ROLE_ID,notEmpty" envDefault:"1"`
	SecretIdFile      string `env:"VAULT_APPROLE_SECRET_ID_FILE"`
	ApiKeyPath        string `env:"VAULT_API_KEY_PATH"`
	KeyMountPath      string `env:"VAULT_API_KEY_MOUNT_PATH"`
	ApiKeyField       string `env:"VAULT_API_KEY_FIELD"`
	DatabaseCredsPath string `env:"VAULT_DATABASE_CREDS_PATH"`
}

func NewVault() Vault {
	cfg := Vault{}

	if err := env.Parse(&cfg); err != nil {
		// Todo:: Handle Panic!
		panic(err)
	}

	return cfg
}

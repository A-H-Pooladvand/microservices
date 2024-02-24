package vault

import (
	"github.com/caarlos0/env/v10"
)

// Config Must not be parsed via godotenv.Load()
type Config struct {
	Address   string `env:"VAULT_ADDRESS"`
	RoleID    string `env:"VAULT_APPROLE_ROLE_ID"`
	SecretId  string `env:"VAULT_APPROLE_SECRET_ID"`
	MountPath string `env:"VAULT_APPROLE_MOUNT_PATH"`
}

func NewConfig() (Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)

	return cfg, err
}

func (v Config) NotEmpty() bool {
	return !v.Empty()
}

func (v Config) Empty() bool {
	return v.Address == "" && v.RoleID == "" && v.SecretId == "" && v.MountPath == ""
}

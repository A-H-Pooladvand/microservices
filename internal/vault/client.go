package vault

import (
	"po/pkg/vault"
)

// New vault implementation
func New() (*vault.Client, error) {
	config, err := NewConfig()

	if err != nil {
		return nil, err
	}

	c := vault.NewConfig(
		config.Address,
		config.RoleID,
		config.SecretId,
	)

	return vault.New(c)
}

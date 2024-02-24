package vault

type Config struct {
	Address  string
	RoleID   string
	SecretID string
}

func NewConfig(
	address string,
	roleID string,
	secretID string,
) Config {
	return Config{
		Address:  address,
		RoleID:   roleID,
		SecretID: secretID,
	}
}

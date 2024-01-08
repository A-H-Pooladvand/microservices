package vault

type Config struct {
	Address                 string
	ApproleRoleID           string
	ApproleSecretIDFile     string
	ApiKeyPath              string
	ApiKeyMountPath         string
	ApiKeyField             string
	DatabaseCredentialsPath string
}

func NewConfig(
	address string,
	approleRoleID string,
	approleSecretIDFile string,
	apiKeyPath string,
	apiKeyMountPath string,
	apiKeyField string,
	databaseCredentialsPath string,
) Config {
	return Config{
		Address:                 address,
		ApproleRoleID:           approleRoleID,
		ApproleSecretIDFile:     approleSecretIDFile,
		ApiKeyPath:              apiKeyPath,
		ApiKeyMountPath:         apiKeyMountPath,
		ApiKeyField:             apiKeyField,
		DatabaseCredentialsPath: databaseCredentialsPath,
	}
}

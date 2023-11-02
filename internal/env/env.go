package env

import (
	"context"
	"github.com/jessevdk/go-flags"
	"os"
	"sync"
)

var (
	environment Environment
	once        = sync.Once{}
)

type Environment struct {
	AppAddress   string `env:"APP_ADDRESS"                   default:":8000"                        description:"App port"                                                              long:"port"`
	VaultAddress string `env:"VAULT_ADDRESS"                 default:"127.0.0.1:8200"               description:"Vault address"                                                         long:"vault-address"`
	//VaultAppRoleRoleID       string        `env:"VAULT_APPROLE_ROLE_ID"         required:"true"                        description:"AppRole RoleID to logger in to Vault"                                     long:"vault-approle-role-id"`
	//VaultAppRoleSecretIDFile string        `env:"VAULT_APPROLE_SECRET_ID_FILE"  default:"/tmp/secret"                  description:"AppRole SecretID file path to logger in to Vault"                         long:"vault-approle-secret-id-file"`
	//VaultAPIKeyPath          string        `env:"VAULT_API_KEY_PATH"            default:"api-key"                      description:"Path to the API key used by 'secure-service'"                          long:"vault-api-key-path"`
	//VaultAPIKeyMountPath     string        `env:"VAULT_API_KEY_MOUNT_PATH"      default:"kv-v2"                        description:"The location where the KV v2 secrets engine has been mounted in Vault" long:"vault-api-key-mount-path"`
	//VaultAPIKeyField         string        `env:"VAULT_API_KEY_FIELD"           default:"api-key-field"                description:"The secret field name for the API key"                                 long:"vault-api-key-descriptor"`
	//VaultDatabaseCredsPath   string        `env:"VAULT_DATABASE_CREDS_PATH"     default:"database/creds/dev-readonly"  description:"Temporary database credentials will be generated here"                 long:"vault-database-creds-path"`
	//DatabaseHostname         string        `env:"DATABASE_HOSTNAME"             required:"true"                        description:"PostgreSQL database hostname"                                          long:"database-hostname"`
	//DatabasePort             string        `env:"DATABASE_PORT"                 default:"5432"                         description:"PostgreSQL database port"                                              long:"database-port"`
	//DatabaseName             string        `env:"DATABASE_NAME"                 default:"postgres"                     description:"PostgreSQL database name"                                              long:"database-name"`
	//DatabaseTimeout          time.Duration `env:"DATABASE_TIMEOUT"              default:"10s"                          description:"PostgreSQL database connection timeout"                                long:"database-timeout"`
	//SecureServiceAddress     string        `env:"SECURE_SERVICE_ADDRESS"        required:"true"                        description:"3rd party service that requires secure credentials"                    long:"secure-service-address"`
}

func (e Environment) Get() Environment {
	return environment
}

func (e Environment) Boot(ctx context.Context) error {
	once.Do(func() {
		_, err := flags.Parse(&environment)

		if err != nil {
			if flags.WroteHelp(err) {
				os.Exit(0)
			}

			panic(err)
		}
	})

	for {
		select {
		case <-ctx.Done():
			return e.Shutdown()
		}
	}
}

func (e Environment) Shutdown() error {
	return nil
}

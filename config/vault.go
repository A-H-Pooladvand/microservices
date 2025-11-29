package config

// Vault holds Vault configuration.
type Vault struct {
	Address   string `mapstructure:"address"`
	RoleID    string `mapstructure:"role_id"`
	SecretID  string `mapstructure:"secret_id"`
	MountPath string `mapstructure:"mount_path"`
	Namespace string `mapstructure:"namespace"`
}

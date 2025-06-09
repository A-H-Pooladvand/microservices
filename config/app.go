package config

import "strings"

type App struct {
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"Port"`
	Debug    string `mapstructure:"debug"`
	Env      string `mapstructure:"env"`
	Timezone string `mapstructure:"timezone"`
}

func (a App) Prod() bool {
	env := strings.ToLower(a.Env)

	return env == "production" || env == "prod"
}

func (a App) Dev() bool {
	return !a.Prod()
}

func (a App) Debuggable() bool {
	return strings.ToLower(a.Debug) == "true" || a.Debug == "1"
}

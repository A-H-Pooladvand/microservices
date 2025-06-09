package config

type Redis struct {
	Addr string `mapstructure:"addr"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"password"`
}

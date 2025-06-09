package config

type RabbitMQ struct {
	Addr string `mapstructure:"addr"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

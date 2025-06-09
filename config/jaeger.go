package config

type Jaeger struct {
	Addr          string  `mapstructure:"addr"`
	SamplingRatio float64 `mapstructure:"SamplingRatio"`
	UseTLS        bool    `mapstructure:"UseTLS"`
	Username      string  `mapstructure:"username"`
	Password      string  `mapstructure:"password"`
}

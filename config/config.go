package config

import (
	"fmt"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify(config *Config)
}

type Observer func(config *Config)

type Config struct {
	observers []Observer
	App       App      `mapstructure:"app"`
	Inquiry   Inquiry  `mapstructure:"inquiry"`
	Jaeger    Jaeger   `mapstructure:"jaeger"`
	Logstash  Logstash `mapstructure:"logstash"`
	Postgres  Postgres `mapstructure:"postgres"`
	RabbitMQ  RabbitMQ `mapstructure:"rabbitmq"`
	Redis     Redis    `mapstructure:"redis"`
	GRPC      GRPC     `mapstructure:"grpc"`
}

func New() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetDefault("app.name", "app")
	viper.SetDefault("app.port", 8000)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.WatchConfig()

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Panic("failed to unmarshal config file", zap.Error(err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.Unmarshal(&config)

		if err != nil {
			log.Panic("failed to unmarshal config file", zap.Error(err))
			return
		}

		config.Notify()
	})

	return &config
}

func (c *Config) Attach(observer Observer) {
	c.observers = append(c.observers, observer)
}

func (c *Config) Detach(observer Observer) {
	for i, o := range c.observers {
		// Can't compare function values directly in Go, so skip unless it's the same ref
		if fmt.Sprintf("%p", o) == fmt.Sprintf("%p", observer) {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

func (c *Config) Notify() {
	for _, observer := range c.observers {
		observer(c)
	}
}

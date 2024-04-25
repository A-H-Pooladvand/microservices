package configs

import "po/pkg/vault"

type Logstash struct {
	Address string `env:"LOGSTASH_ADDRESS" json:"address"`
}

func NewLogstash(client *vault.Client) (*Logstash, error) {
	c := &Logstash{}

	return c, Parse("logstash", c, client)
}

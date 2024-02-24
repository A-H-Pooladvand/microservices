package configs

type Logstash struct {
	Address string `env:"LOGSTASH_ADDRESS" json:"address"`
}

func NewLogstash() (Logstash, error) {
	c := Logstash{}

	err := parse("logstash", &c)

	return c, err
}

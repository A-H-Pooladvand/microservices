package configs

type Redis struct {
	Addr string `env:"REDIS_ADDRESS" envDefault:"127.0.0.1:5672" json:"address"`
	User string `env:"REDIS_USER,notEmpty" json:"user"`
	Pass string `env:"REDIS_PASS" json:"pass"`
}

func NewRedis() (Redis, error) {
	c := Redis{}

	err := parse("redis", &c)

	return c, err
}

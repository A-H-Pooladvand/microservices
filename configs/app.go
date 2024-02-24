package configs

type App struct {
	Name  string `env:"APP_NAME,notEmpty" envDefault:"App" json:"name"`
	Port  string `env:"APP_PORT" envDefault:"8000" json:"port"`
	Debug string `env:"APP_DEBUG" envDefault:"true" json:"debug"`
}

func NewApp() (App, error) {
	c := App{}

	err := parse("app", &c)

	return c, err
}

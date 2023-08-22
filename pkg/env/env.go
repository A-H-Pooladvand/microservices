package env

import (
	"os"
	"po/pkg/log"
)

func Get(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Sugar.Debugf("Missing %s env variable, using the default one: %s", key, defaultValue)
		value = defaultValue
	}

	return value
}

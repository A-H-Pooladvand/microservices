package app

import (
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"strings"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func Production() bool {
	env := strings.ToLower(os.Getenv("APP_ENV"))

	return env == "prod" || env == "production"
}

func Local() bool {
	return !Production()
}

func LocalMessage() {
	if Production() {
		return
	}

	warning := color.New(color.FgHiYellow).Add(color.Bold)
	_, _ = warning.Println("----------------------------------------------------------------------------------------------------")
	_, _ = warning.Println("| ⚠️ Warning: Application is running in LOCAL environment. ⚠️                                      |")
	_, _ = warning.Println("| ⚠️ If this is unintended, please switch to the PRODUCTION environment for accurate results. ⚠️   |")
	_, _ = warning.Println("----------------------------------------------------------------------------------------------------")
}

func GetContext(c echo.Context) *Context {
	ctx, ok := c.(*Context)

	if !ok {
		zap.L().Panic("unable to get context")
	}

	return ctx
}

func LoadEnvironmentVariablesInLocalEnv() {
	if Local() {
		if err := godotenv.Load(); err != nil {
			zap.L().Panic("unable to load .env file", zap.Error(err))
		}
	}
}

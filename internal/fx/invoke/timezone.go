package invoke

import (
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"go.uber.org/zap"
	"time"
)

const defaultTimezone = "UTC" // Default timezone if the configured one is invalid.

func Timezone(config *config.Config) error {
	location, err := time.LoadLocation(config.App.Timezone)

	if err == nil {
		time.Local = location
		return nil
	}

	// Log a structured warning if the configured timezone is invalid.
	log.Warn("invalid configured timezone, falling back to default.",
		zap.String("configured_timezone", config.App.Timezone),
		zap.Error(err),
	)

	// If the configured timezone is invalid, try a sensible default.
	location, err = time.LoadLocation(defaultTimezone)
	if err != nil {
		// This error indicates a fundamental issue if "Asia/Tehran" also fails.
		log.Error("failed to load default timezone. Critical error",
			zap.String("default_timezone", defaultTimezone),
			zap.Error(err),
		)
		return err
	}

	// Set the loaded location as the application's local timezone.
	time.Local = location

	return nil
}

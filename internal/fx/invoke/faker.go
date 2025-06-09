package invoke

import (
	"github.com/a-h-pooladvand/microservices/internal/log"
	fake "github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/interfaces"
	"go.uber.org/zap"
)

func Faker() error {
	for name, fn := range taggedFunctions() {
		if err := fake.AddProvider(name, fn); err != nil {
			log.Warn("Failed to add fake function", zap.String("name", name), zap.Error(err))
			return err
		}
	}

	return nil
}

func taggedFunctions() map[string]interfaces.TaggedFunction {
	return map[string]interfaces.TaggedFunction{}
}

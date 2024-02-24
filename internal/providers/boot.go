package providers

import (
	"context"
)

var providers = []Booter{
	Vault,
	Postgres,
	APM,
}

type Booter func(ctx context.Context) error

func Boot(ctx context.Context) error {
	for _, fn := range providers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

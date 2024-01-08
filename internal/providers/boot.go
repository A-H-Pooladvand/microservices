package providers

import "po/internal/app"

var providers = []Booter{
	//Vault,
}

type Booter func(ctx app.Context) error

func Boot(ctx app.Context) error {
	for _, fn := range providers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

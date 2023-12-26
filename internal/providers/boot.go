package providers

import "po/internal/app"

var providers = []Booter{
	&Vault{},
}

type Booter interface {
	Boot(ctx app.Context) error
}

func Boot(ctx app.Context) error {
	for _, provider := range providers {
		if err := provider.Boot(ctx); err != nil {
			return err
		}
	}

	return nil
}

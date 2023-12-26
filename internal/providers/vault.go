package providers

import (
	"po/internal/app"
)

type Vault struct {
}

func (v *Vault) Boot(ctx app.Context) error {
	return nil
}

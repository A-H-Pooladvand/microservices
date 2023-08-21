package boot

import (
	"context"
	"golang.org/x/sync/errgroup"
	"po/pkg/vault"
)

var bootstraps = []Booter{
	vault.Vault{},
}

type Booter interface {
	Boot(ctx context.Context) error
}

func Boot(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, booter := range bootstraps {
		b := booter
		group.Go(func() error {
			return b.Boot(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

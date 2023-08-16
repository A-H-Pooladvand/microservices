package boot

import (
	"context"
	"golang.org/x/sync/errgroup"
)

var bootstraps = []Booter{}

type Booter interface {
	Boot(ctx context.Context) error
}

func Boot(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, booter := range bootstraps {
		group.Go(func() error {
			return booter.Boot(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

package boot

import (
	"context"
	"golang.org/x/sync/errgroup"
)

var bootstraps = []Booter{}

type Booter interface {
	Boot(ctx context.Context) error
	Shutdown() error
}

func Boot(ctx context.Context) (*errgroup.Group, error) {
	group, ctx := errgroup.WithContext(ctx)

	for _, booter := range bootstraps {
		b := booter
		group.Go(func() error {
			return b.Boot(ctx)
		})
	}

	bootstraps = nil

	return group, nil
}

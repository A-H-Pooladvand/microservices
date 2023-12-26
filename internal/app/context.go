package app

import "context"

type Context struct {
	context.Context
}

func WithCancel() (Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	return Context{
		Context: ctx,
	}, cancel
}

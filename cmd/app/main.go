package main

import (
	"context"
	"po/internal/app"
	"po/internal/boot"
	"po/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer panicRecover()

	logger.Boot()

	group, err := boot.Boot(ctx)

	if err != nil {
		panic(err)
	}

	if err = app.Serve(ctx); err != nil {
		panic(err)
	}

	if err = group.Wait(); err != nil {
		panic(err)
	}
}

func panicRecover() {
	if r := recover(); r != nil {
		logger.Error(r)

		panic(r)
	}
}

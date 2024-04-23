package cmd

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	golog "github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"po/internal/app"
	"po/internal/grpc"
	"po/internal/providers"
	"po/internal/webserver"
	"po/pkg/log"
)

type Message string

func NewMessage() Message {
	return "Hello World"
}

type Greeter struct {
	Message Message
}

func NewGreeter(message Message) Greeter {
	return Greeter{
		Message: message,
	}
}

func (g Greeter) Hi() {
	fmt.Println(g.Message)
}

var cmd = &cobra.Command{
	Use:   "app",
	Short: "webserver",
	Long:  `Initializing...`,
}

func Execute() {
	// Load environment variables if APP_ENV set to local
	if app.Local() {
		if err := godotenv.Load(); err != nil {
			zap.L().Panic("unable to load .env file", zap.Error(err))
		}
	}

	// First it's mandatory to load our logger service
	// before any other thing
	log.Boot()

	defer func(l *zap.Logger) {
		if err := l.Sync(); err != nil {
			golog.Error(err)
		}
	}(zap.L())

	a := app.NewSingleton()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.Ctx = &ctx
	defer panicRecover(cancel)

	// Boot third party services
	if err := providers.Boot(ctx); err != nil {
		zap.L().Panic("unable to boot service", zap.Error(err))
	}

	cmd.AddCommand(serveCmd)
	cmd.AddCommand(migrateCmd)
	cmd.AddCommand(seedCmd)

	if err := cmd.Execute(); err != nil {
		zap.L().Panic("err", zap.Error(err))
	}
}

func serve(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.Serve(ctx)
	})

	g.Go(func() error {
		return webserver.New().Serve(ctx)
	})

	return g.Wait()
}

func panicRecover(cancel context.CancelFunc) {
	if r := recover(); r != nil {
		cancel()
		zap.L().Error("panic recovery", zap.Any("panic message", fmt.Sprintf("%v", r)))
		_ = zap.L().Sync()
		app.Get().GracefulShutdown()
	}
}

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"po/internal/app"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the application",
	// Serve necessary protocols such as gRPC, HTTP etc...
	Run: func(cmd *cobra.Command, args []string) {
		app.LocalMessage()

		if err := serve(app.Get().Context()); err != nil {
			zap.L().Panic("err to serve necessary protocols such as gRPC, HTTP etc", zap.Error(err))
		}
	},
}

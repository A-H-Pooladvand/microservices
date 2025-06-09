package cmd

import (
	"github.com/a-h-pooladvand/microservices/internal/fx"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	Long:  "Serve the application and start the server",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	fx.App().Run()
}

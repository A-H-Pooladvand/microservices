package cmd

import (
	"github.com/spf13/cobra"
	"po/internal/app"
	"po/pkg/logger"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Serves the application",
	Long:  `Serving the application`,
	Run:   serve,
}

func init() {
	cmd.AddCommand(serveCommand)
}

func serve(command *cobra.Command, args []string) {
	ctx, cancel := app.WithCancel()
	defer cancel()

	a := app.New()
	defer a.Recover(cancel)

	if err := a.Serve(ctx); err != nil {
		logger.Panic(err)
	}
}

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cmd = &cobra.Command{
	Use:   "app",
	Short: "webserver",
	Long:  `Initializing...`,
}

func Execute() {
	cmd.AddCommand(serveCmd)
	cmd.AddCommand(migrateCmd)
	cmd.AddCommand(seedCmd)

	if err := cmd.Execute(); err != nil {
		zap.L().Panic("err", zap.Error(err))
	}
}

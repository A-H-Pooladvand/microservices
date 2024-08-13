package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cmd = &cobra.Command{
	Use:   "app",
	Short: "app is a simple application",
	Long:  "app is a simple application",
}

func Execute() {
	cmd.AddCommand(serveCmd)
	cmd.AddCommand(migrateCmd)
	cmd.AddCommand(seedCmd)

	if err := cmd.Execute(); err != nil {
		zap.L().Panic("err", zap.Error(err))
	}
}

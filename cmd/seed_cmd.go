package cmd

import (
	"github.com/spf13/cobra"
	"po/cmd/seed"
)

var seeders = []seed.Seeder{
	seed.UserSeeder{},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeds the database",
	// Serve necessary protocols such as gRPC, HTTP etc...
	Run: func(cmd *cobra.Command, args []string) {
		for _, seeder := range seeders {
			seeder.Run()
		}
	},
}

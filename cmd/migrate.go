package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"po/internal/db"
	"po/internal/models"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs the migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Begin running the migrations")

		err := db.New().AutoMigrate(
			migrations(),
		)

		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("All migrations completed successfully")
	},
}

func migrations() []any {
	return []any{
		models.User{},
	}
}

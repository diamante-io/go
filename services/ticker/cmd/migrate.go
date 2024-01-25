package cmd

import (
	"go/services/ticker/internal/tickerdb"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdMigrate)
}

var cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Updates the database to the latest schema version.",
	Run: func(cmd *cobra.Command, args []string) {
		Logger.Info("Refreshing the asset database")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		Logger.Info("Upgrading the database")
		n, err := tickerdb.MigrateDB(&session)
		if err != nil {
			Logger.Fatal("could not upgrade the database:", err)
		}
		Logger.Infof("Database Successfully Upgraded. Applied %d migrations.\n", n)
	},
}

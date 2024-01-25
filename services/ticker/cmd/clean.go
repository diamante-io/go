package cmd

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
)

var DaysToKeep int

func init() {
	rootCmd.AddCommand(cmdClean)
	cmdClean.AddCommand(cmdCleanTrades)

	cmdCleanTrades.Flags().IntVarP(
		&DaysToKeep,
		"keep-days",
		"k",
		7,
		"Trade entries older than keep-days will be deleted",
	)
}

var cmdClean = &cobra.Command{
	Use:   "clean [data type]",
	Short: "Cleans up the database for a given data type",
}

var cmdCleanTrades = &cobra.Command{
	Use:   "trades",
	Short: "Cleans up old trades from the database",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}

		now := time.Now()
		minDate := now.AddDate(0, 0, -DaysToKeep)
		Logger.Infof("Deleting trade entries older than %d days", DaysToKeep)
		err = session.DeleteOldTrades(context.Background(), minDate)
		if err != nil {
			Logger.Fatal("could not delete trade entries:", err)
		}
	},
}

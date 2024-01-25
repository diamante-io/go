package cmd

import (
	"context"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
	ticker "github.com/diamcircle/go/services/ticker/internal"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
)

var ShouldStream bool
var BackfillHours int

func init() {
	rootCmd.AddCommand(cmdIngest)
	cmdIngest.AddCommand(cmdIngestAssets)
	cmdIngest.AddCommand(cmdIngestTrades)
	cmdIngest.AddCommand(cmdIngestOrderbooks)

	cmdIngestTrades.Flags().BoolVar(
		&ShouldStream,
		"stream",
		false,
		"Continuously stream new trades from the Aurora Stream API as a daemon",
	)

	cmdIngestTrades.Flags().IntVar(
		&BackfillHours,
		"num-hours",
		7*24,
		"Number of past hours to backfill trade data",
	)
}

var cmdIngest = &cobra.Command{
	Use:   "ingest [data type]",
	Short: "Ingests new data from data type into the database.",
}

var cmdIngestAssets = &cobra.Command{
	Use:   "assets",
	Short: "Refreshes the asset database with new data retrieved from Aurora.",
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

		ctx := context.Background()
		err = ticker.RefreshAssets(ctx, &session, Client, Logger)
		if err != nil {
			Logger.Fatal("could not refresh asset database:", err)
		}
	},
}

var cmdIngestTrades = &cobra.Command{
	Use:   "trades",
	Short: "Fills the trade database with data retrieved form Aurora.",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ctx := context.Background()
		numDays := float32(BackfillHours) / 24.0
		Logger.Infof(
			"Backfilling Trade data for the past %d hour(s) [%.2f days]\n",
			BackfillHours,
			numDays,
		)
		err = ticker.BackfillTrades(ctx, &session, Client, Logger, BackfillHours, 0)
		if err != nil {
			Logger.Fatal("could not refresh trade database:", err)
		}

		if ShouldStream {
			Logger.Info("Streaming new data (this is a continuous process)")
			err = ticker.StreamTrades(ctx, &session, Client, Logger)
			if err != nil {
				Logger.Fatal("could not refresh trade database:", err)
			}
		}
	},
}

var cmdIngestOrderbooks = &cobra.Command{
	Use:   "orderbooks",
	Short: "Refreshes the orderbook stats database with new data retrieved from Aurora.",
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

		err = ticker.RefreshOrderbookEntries(&session, Client, Logger)
		if err != nil {
			Logger.Fatal("could not refresh orderbook database:", err)
		}
	},
}

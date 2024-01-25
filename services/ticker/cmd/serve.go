package cmd

import (
	"github.com/lib/pq"
	"github.com/spf13/cobra"
	ticker "github.com/diamcircle/go/services/ticker/internal"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
)

var ServerAddr string

func init() {
	rootCmd.AddCommand(cmdServe)

	cmdServe.Flags().StringVar(
		&ServerAddr,
		"address",
		"0.0.0.0:3000",
		"Server address and port",
	)
}

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Runs a GraphQL interface to get Ticker data",
	Run: func(cmd *cobra.Command, args []string) {
		Logger.Info("Starting GraphQL Server")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ticker.StartGraphQLServer(&session, Logger, ServerAddr)
	},
}

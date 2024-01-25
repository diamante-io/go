package cmd

import (
	"context"

	ticker "go/services/ticker/internal"
	"go/services/ticker/internal/tickerdb"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
)

var MarketsOutFile string
var AssetsOutFile string

func init() {
	rootCmd.AddCommand(cmdGenerate)
	cmdGenerate.AddCommand(cmdGenerateMarketData)
	cmdGenerate.AddCommand(cmdGenerateAssetData)

	cmdGenerateMarketData.Flags().StringVarP(
		&MarketsOutFile,
		"out-file",
		"o",
		"markets.json",
		"Set the name of the output file",
	)

	cmdGenerateAssetData.Flags().StringVarP(
		&AssetsOutFile,
		"out-file",
		"o",
		"assets.json",
		"Set the name of the output file",
	)
}

var cmdGenerate = &cobra.Command{
	Use:   "generate [data type]",
	Short: "Generates reports about assets and markets",
}

var cmdGenerateMarketData = &cobra.Command{
	Use:   "market-data",
	Short: "Generate the aggregated market data (for 24h and 7d) and outputs to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}

		Logger.Infof("Starting market data generation, outputting to: %s\n", MarketsOutFile)
		err = ticker.GenerateMarketSummaryFile(&session, Logger, MarketsOutFile)
		if err != nil {
			Logger.Fatal("could not generate market data:", err)
		}
	},
}

var cmdGenerateAssetData = &cobra.Command{
	Use:   "asset-data",
	Short: "Generate the aggregated asset data and outputs to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}

		Logger.Infof("Starting asset data generation, outputting to: %s\n", AssetsOutFile)
		err = ticker.GenerateAssetsFile(context.Background(), &session, Logger, AssetsOutFile)
		if err != nil {
			Logger.Fatal("could not generate asset data:", err)
		}
	},
}

package cmd

import (
	aurora "go/services/aurora/internal"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run aurora server",
	Long:  "serve initializes then starts the aurora HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := aurora.NewAppFromFlags(config, flags)
		if err != nil {
			return err
		}
		return app.Serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

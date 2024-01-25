package cmd

import (
	"fmt"

	"go/clients/auroraclient"
	demo "go/txnbuild/cmd/demo/operations"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create and fund some demo accounts on the TestNet",
	Long:  `This command creates four test accounts for use with further operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Initialising TestNet accounts...")
		keys := demo.InitKeys(4)
		client := auroraclient.DefaultTestNetClient

		demo.Initialise(client, keys)
		fmt.Println("Initialisation complete.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

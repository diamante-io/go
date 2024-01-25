package main

import (
	"log"

	"go/tools/diamcircle-hd-wallet/commands"

	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "diamcircle-hd-wallet",
	Short: "Simple HD wallet for diamcircle Lumens. THIS PROGRAM IS STILL EXPERIMENTAL. USE AT YOUR OWN RISK.",
}

func init() {
	mainCmd.AddCommand(commands.NewCmd)
	mainCmd.AddCommand(commands.AccountsCmd)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

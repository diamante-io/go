package cmd

import (
	"fmt"
	"runtime"

	apkg "go/support/app"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print aurora and Golang runtime version",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(apkg.Version())
		fmt.Println(runtime.Version())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

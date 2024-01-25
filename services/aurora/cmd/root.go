package cmd

import (
	"fmt"
	stdLog "log"

	aurora "go/services/aurora/internal"

	"github.com/spf13/cobra"
)

var (
	config, flags = aurora.Flags()

	RootCmd = &cobra.Command{
		Use:           "aurora",
		Short:         "client-facing api server for the diamcircle network",
		SilenceErrors: true,
		SilenceUsage:  true,
		Long:          "Client-facing API server for the diamcircle network. It acts as the interface between diamcircle Core and applications that want to access the diamcircle network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams and more.",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := aurora.NewAppFromFlags(config, flags)
			if err != nil {
				return err
			}
			return app.Serve()
		},
	}
)

// ErrUsage indicates we should print the usage string and exit with code 1
type ErrUsage struct {
	cmd *cobra.Command
}

func (e ErrUsage) Error() string {
	return e.cmd.UsageString()
}

// Indicates we want to exit with a specific error code without printing an error.
type ErrExitCode int

func (e ErrExitCode) Error() string {
	return fmt.Sprintf("exit code: %d", e)
}

func init() {
	err := flags.Init(RootCmd)
	if err != nil {
		stdLog.Fatal(err.Error())
	}
}

func Execute() error {
	return RootCmd.Execute()
}

package cmd

import (
	"go/types"

	"go/clients/auroraclient"
	"go/network"
	"go/services/regulated-assets-approval-server/internal/serve"
	"go/support/config"

	"github.com/spf13/cobra"
)

type ServeCommand struct{}

func (c *ServeCommand) Command() *cobra.Command {
	opts := serve.Options{}
	configOpts := config.ConfigOptions{
		{
			Name:      "issuer-account-secret",
			Usage:     "Secret key of the issuer account.",
			OptType:   types.String,
			ConfigKey: &opts.IssuerAccountSecret,
			Required:  true,
		},
		{
			Name:      "asset-code",
			Usage:     "The code of the regulated asset",
			OptType:   types.String,
			ConfigKey: &opts.AssetCode,
			Required:  true,
		},
		{
			Name:        "database-url",
			Usage:       "Database URL",
			OptType:     types.String,
			ConfigKey:   &opts.DatabaseURL,
			FlagDefault: "postgres://localhost:5432/?sslmode=disable",
			Required:    true,
		},
		{
			Name:        "friendbot-payment-amount",
			Usage:       "The amount of regulated assets the friendbot will be distributing",
			OptType:     types.Int,
			ConfigKey:   &opts.FriendbotPaymentAmount,
			FlagDefault: 10000,
			Required:    true,
		},
		{
			Name:        "aurora-url",
			Usage:       "aurora URL used for looking up account details",
			OptType:     types.String,
			ConfigKey:   &opts.auroraURL,
			FlagDefault: auroraclient.DefaultTestNetClient.auroraURL,
			Required:    true,
		},
		{
			Name:        "network-passphrase",
			Usage:       "Network passphrase of the diamcircle network transactions should be signed for",
			OptType:     types.String,
			ConfigKey:   &opts.NetworkPassphrase,
			FlagDefault: network.TestNetworkPassphrase,
			Required:    true,
		},
		{
			Name:        "port",
			Usage:       "Port to listen and serve on",
			OptType:     types.Int,
			ConfigKey:   &opts.Port,
			FlagDefault: 8000,
			Required:    true,
		},
		{
			Name:      "base-url",
			Usage:     "The base url address to this server",
			OptType:   types.String,
			ConfigKey: &opts.BaseURL,
			Required:  true,
		},
		{
			Name:        "kyc-required-payment-amount-threshold",
			Usage:       "The amount threshold when KYC is required, may contain decimals and is greater than 0",
			OptType:     types.String,
			ConfigKey:   &opts.KYCRequiredPaymentAmountThreshold,
			FlagDefault: "500",
			Required:    true,
		},
	}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the SEP-8 Approval Server",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
			c.Run(opts)
		},
	}
	configOpts.Init(cmd)
	return cmd
}

func (c *ServeCommand) Run(opts serve.Options) {
	serve.Serve(opts)
}

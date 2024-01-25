package cmd

import (
	"go/types"

	"github.com/spf13/cobra"
	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/configureissuer"
	"github.com/diamcircle/go/support/config"
)

type ConfigureIssuer struct{}

func (c *ConfigureIssuer) Command() *cobra.Command {
	opts := configureissuer.Options{}
	configOpts := config.ConfigOptions{
		{
			Name:      "asset-code",
			Usage:     "The code of the regulated asset",
			OptType:   types.String,
			ConfigKey: &opts.AssetCode,
			Required:  true,
		},
		{
			Name:      "base-url",
			Usage:     "The base url to the server where the asset home domain should be. For instance, \"https://test.example.com/\" if your desired asset home domain is \"test.example.com\".",
			OptType:   types.String,
			ConfigKey: &opts.BaseURL,
			Required:  true,
		},
		{
			Name:        "aurora-url",
			Usage:       "Aurora URL used for looking up account details",
			OptType:     types.String,
			ConfigKey:   &opts.AuroraURL,
			FlagDefault: auroraclient.DefaultTestNetClient.AuroraURL,
			Required:    true,
		},
		{
			Name:      "issuer-account-secret",
			Usage:     "Secret key of the issuer account.",
			OptType:   types.String,
			ConfigKey: &opts.IssuerAccountSecret,
			Required:  true,
		},
		{
			Name:        "network-passphrase",
			Usage:       "Network passphrase of the Diamcircle network transactions should be signed for",
			OptType:     types.String,
			ConfigKey:   &opts.NetworkPassphrase,
			FlagDefault: network.TestNetworkPassphrase,
			Required:    true,
		},
	}
	cmd := &cobra.Command{
		Use:   "configure-issuer",
		Short: "Configure the Asset Issuer Account for SEP-8 Regulated Assets",
		Run: func(_ *cobra.Command, _ []string) {
			configOpts.Require()
			configOpts.SetValues()
			c.Run(opts)
		},
	}
	configOpts.Init(cmd)
	return cmd
}

func (c *ConfigureIssuer) Run(opts configureissuer.Options) {
	configureissuer.Setup(opts)
}

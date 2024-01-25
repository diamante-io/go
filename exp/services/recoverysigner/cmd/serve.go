package cmd

import (
	"go/types"

	"go/exp/services/recoverysigner/internal/serve"
	"go/network"
	"go/support/config"
	supportlog "go/support/log"

	"github.com/spf13/cobra"
)

type ServeCommand struct {
	Logger *supportlog.Entry
}

func (c *ServeCommand) Command() *cobra.Command {
	opts := serve.Options{
		Logger: c.Logger,
	}
	configOpts := config.ConfigOptions{
		{
			Name:        "port",
			Usage:       "Port to listen and serve on",
			OptType:     types.Int,
			ConfigKey:   &opts.Port,
			FlagDefault: 8000,
			Required:    true,
		},
		{
			Name:        "db-url",
			Usage:       "Database URL",
			OptType:     types.String,
			ConfigKey:   &opts.DatabaseURL,
			FlagDefault: "postgres://localhost:5432/?sslmode=disable",
			Required:    false,
		},
		{
			Name:        "db-max-open-conns",
			Usage:       "Database max open connections",
			OptType:     types.Int,
			ConfigKey:   &opts.DatabaseMaxOpenConns,
			FlagDefault: 20,
			Required:    false,
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
			Name:      "signing-key",
			Usage:     "diamcircle signing key(s) used for signing transactions comma separated (first key is preferred signer) (will be deprecated with per-account keys in the future)",
			OptType:   types.String,
			ConfigKey: &opts.SigningKeys,
			Required:  true,
		},
		{
			Name:      "sep10-jwks",
			Usage:     "JSON Web Key Set (JWKS) containing one or more keys used to validate SEP-10 JWTs (if the key is an asymmetric key that has separate public and private key, the JWK need only contain the public key) (if multiple keys are provided they will all attempt verification the key ID will be ignored although logged)",
			OptType:   types.String,
			ConfigKey: &opts.SEP10JWKS,
			Required:  true,
		},
		{
			Name:      "sep10-jwt-issuer",
			Usage:     "JWT issuer to verify is in the SEP-10 JWT iss field (not checked if empty)",
			OptType:   types.String,
			ConfigKey: &opts.SEP10JWTIssuer,
			Required:  false,
		},
		{
			Name:      "firebase-project-id",
			Usage:     "Firebase project ID to use for validating Firebase JWTs",
			OptType:   types.String,
			ConfigKey: &opts.FirebaseProjectID,
			Required:  true,
		},
		{
			Name:        "admin-port",
			Usage:       "Port to listen and serve admin functionality including metrics",
			OptType:     types.Int,
			ConfigKey:   &opts.AdminPort,
			FlagDefault: 0,
			Required:    false,
		},
		{
			Name:        "metrics-namespace",
			Usage:       "Namespace to use for metric names prefixed to metrics reported",
			OptType:     types.String,
			ConfigKey:   &opts.MetricsNamespace,
			FlagDefault: "recoverysigner",
			Required:    false,
		},
		{
			Name:      "allowed-source-accounts",
			Usage:     "diamcircle account(s) allowed as source accounts in transactions signed for all users in addition to the registered account comma separated (important: these accounts must never be registered accounts and must never have the signer configured that is a signing key used by this server)",
			OptType:   types.String,
			ConfigKey: &opts.AllowedSourceAccounts,
			Required:  false,
		},
	}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the SEP-30 Recovery Signer server",
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

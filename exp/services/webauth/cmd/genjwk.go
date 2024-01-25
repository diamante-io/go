package cmd

import (
	"encoding/json"

	"go/exp/support/jwtkey"
	supportlog "go/support/log"

	"github.com/spf13/cobra"
	"gopkg.in/square/go-jose.v2"
)

type GenJWKCommand struct {
	Logger *supportlog.Entry
}

func (c *GenJWKCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genjwk",
		Short: "Generate a JSON Web Key (ECDSA/ES256) for JWT issuing",
		Run: func(_ *cobra.Command, _ []string) {
			c.Run()
		},
	}
	return cmd
}

func (c *GenJWKCommand) Run() {
	k, err := jwtkey.GenerateKey()
	if err != nil {
		c.Logger.Fatal(err)
	}

	alg := jose.ES256

	{
		jwk := jose.JSONWebKey{Key: &k.PublicKey, Algorithm: string(alg)}
		bytes, err := json.Marshal(jwk)
		if err == nil {
			c.Logger.Print("Public:", string(bytes))
		} else {
			c.Logger.Print("Public:", err)
		}
	}

	{
		jwk := jose.JSONWebKey{Key: k, Algorithm: string(alg)}
		bytes, err := json.Marshal(jwk)
		if err == nil {
			c.Logger.Print("Private:", string(bytes))
		} else {
			c.Logger.Print("Private:", err)
		}
	}
}

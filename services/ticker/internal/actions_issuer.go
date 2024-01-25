package ticker

import (
	"go/services/ticker/internal/scraper"
	"go/services/ticker/internal/tickerdb"
)

func tomlIssuerToDBIssuer(issuer scraper.TOMLIssuer) tickerdb.Issuer {
	return tickerdb.Issuer{
		PublicKey:        issuer.SigningKey,
		Name:             issuer.Documentation.OrgName,
		URL:              issuer.Documentation.OrgURL,
		TOMLURL:          issuer.TOMLURL,
		FederationServer: issuer.FederationServer,
		AuthServer:       issuer.AuthServer,
		TransferServer:   issuer.TransferServer,
		WebAuthEndpoint:  issuer.WebAuthEndpoint,
		DepositServer:    issuer.DepositServer,
		OrgTwitter:       issuer.Documentation.OrgTwitter,
	}
}

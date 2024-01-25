package gql

import (
	"context"
	"errors"

	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
)

// Issuers resolves the issuers() GraphQL query.
func (r *resolver) Issuers(ctx context.Context) (issuers []*tickerdb.Issuer, err error) {
	dbIssuers, err := r.db.GetAllIssuers(ctx)
	if err != nil {
		// obfuscating sql errors to avoid exposing underlying
		// implementation
		err = errors.New("could not retrieve the requested data")
	}

	for i := range dbIssuers {
		issuers = append(issuers, &dbIssuers[i])
	}

	return issuers, err
}

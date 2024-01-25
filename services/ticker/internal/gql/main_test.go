package gql

import (
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/diamcircle/go/services/ticker/internal/gql/static"
)

func TestValidateSchema(t *testing.T) {
	r := resolver{}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	graphql.MustParseSchema(static.Schema(), &r, opts...)
}

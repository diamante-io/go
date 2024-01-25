package gql

import (
	"testing"

	"go/services/ticker/internal/gql/static"

	"github.com/graph-gophers/graphql-go"
)

func TestValidateSchema(t *testing.T) {
	r := resolver{}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	graphql.MustParseSchema(static.Schema(), &r, opts...)
}

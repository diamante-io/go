/*
This file implements an interface to this application's
custom BigInt GraphQL scalar type. It can receive values
of type int, int32 and int64.
*/

package gql

import (
	"errors"
	"strconv"
)

type BigInt int

func (BigInt) ImplementsGraphQLType(name string) bool {
	return name == "BigInt"
}

func (bigInt *BigInt) UnmarshalGraphQL(input interface{}) error {
	var err error

	switch input := input.(type) {
	case int:
		*bigInt = BigInt(input)
	case int32:
		*bigInt = BigInt(int(input))
	case int64:
		*bigInt = BigInt(int(input))
	default:
		err = errors.New("wrong type")
	}

	return err
}

func (bigInt BigInt) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(bigInt), 10), nil
}

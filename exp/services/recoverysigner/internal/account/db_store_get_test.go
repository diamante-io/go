package account

import (
	"testing"

	"go/exp/services/recoverysigner/internal/db/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	address := "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT"

	// Store the account
	a := Account{
		Address: address,
		Identities: []Identity{
			{
				Role: "sender",
				AuthMethods: []AuthMethod{
					{Type: AuthMethodTypeAddress, Value: "GD4NGMOTV4QOXWA6PGPIGVWZYMRCJAKLQJKZIP55C5DGB3GBHHET3YC6"},
					{Type: AuthMethodTypePhoneNumber, Value: "+10000000000"},
					{Type: AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []AuthMethod{
					{Type: AuthMethodTypeAddress, Value: "GBJCOYGKIJYX3VUEOZ6GVMFP522UO4OEBI5KB5HHWZAZ2DEJTHS6VOHP"},
					{Type: AuthMethodTypePhoneNumber, Value: "+20000000000"},
					{Type: AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	}
	err := store.Add(a)
	require.NoError(t, err)

	// Reading the account out results in the same data.
	aRoundTrip, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, a, aRoundTrip)
}

func TestGet_noIdentities(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	address := "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT"

	// Store the account
	a := Account{
		Address: address,
	}
	err := store.Add(a)
	require.NoError(t, err)

	// Reading the account out results in the same data.
	aRoundTrip, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, a, aRoundTrip)
}

func TestGet_noAuthMethods(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	address := "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT"

	// Store the account
	a := Account{
		Address:    address,
		Identities: []Identity{{Role: "owner"}},
	}
	err := store.Add(a)
	require.NoError(t, err)

	// Reading the account out results in the same data.
	aRoundTrip, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, a, aRoundTrip)
}

func TestGet_notFound(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	address := "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT"

	_, err := store.Get(address)
	assert.Equal(t, ErrNotFound, err)
}

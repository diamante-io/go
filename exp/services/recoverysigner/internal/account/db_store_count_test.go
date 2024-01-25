package account

import (
	"testing"

	"github.com/diamcircle/go/exp/services/recoverysigner/internal/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCount(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	count, err := store.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	a := Account{
		Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
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
	err = store.Add(a)
	require.NoError(t, err)

	count, err = store.Count()
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

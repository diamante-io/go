package account

import (
	"testing"

	"github.com/diamcircle/go/exp/services/recoverysigner/internal/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindWithIdentityAuthMethod(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	// Register an account with two identities
	a1 := Account{
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
	err := store.Add(a1)
	require.NoError(t, err)

	// Register an account with one identity that overlaps with identities in a1
	a2 := Account{
		Address: "GA3ADWA6QWC6D7VSUS4QZCPYC5SYJQGCBIVLIHO4P2WDGPJRJEQO3QNS",
		Identities: []Identity{
			{
				Role: "owner",
				AuthMethods: []AuthMethod{
					{Type: AuthMethodTypeAddress, Value: "GBJCOYGKIJYX3VUEOZ6GVMFP522UO4OEBI5KB5HHWZAZ2DEJTHS6VOHP"},
					{Type: AuthMethodTypePhoneNumber, Value: "+20000000000"},
					{Type: AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	}
	err = store.Add(a2)
	require.NoError(t, err)

	// Check that the first account can be found by its sender auth methods
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeAddress, "GD4NGMOTV4QOXWA6PGPIGVWZYMRCJAKLQJKZIP55C5DGB3GBHHET3YC6")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1}, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypePhoneNumber, "+10000000000")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1}, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeEmail, "user1@example.com")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1}, found)
	}

	// Check that both accounts can be found by the receiver/owner auth methods
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeAddress, "GBJCOYGKIJYX3VUEOZ6GVMFP522UO4OEBI5KB5HHWZAZ2DEJTHS6VOHP")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1, a2}, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypePhoneNumber, "+20000000000")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1, a2}, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeEmail, "user2@example.com")
		require.NoError(t, err)
		assert.Equal(t, []Account{a1, a2}, found)
	}

	// Check that accounts are not found by their own address
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeAddress, "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT")
		require.NoError(t, err)
		assert.Empty(t, found)
	}

	// Check that accounts are not found by an unrelated auth methods
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeAddress, "GBNZT3ZY6QYLIZLHQRQCHJGBEVV4QLR2CAL3WCMAO52PJMPISIKMS7OQ")
		require.NoError(t, err)
		assert.Empty(t, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypePhoneNumber, "+99999999999")
		require.NoError(t, err)
		assert.Empty(t, found)
	}
	{
		found, err := store.FindWithIdentityAuthMethod(AuthMethodTypeEmail, "user9@example.com")
		require.NoError(t, err)
		assert.Empty(t, found)
	}
}

func TestFindWithIdentityAuthMethod_notFound(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	address := "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT"

	accounts, err := store.FindWithIdentityAuthMethod(AuthMethodTypeAddress, address)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

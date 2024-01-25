package account

import (
	"testing"

	"go/exp/services/recoverysigner/internal/db/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
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

	// Update the identities on the account
	b := Account{
		Address: address,
		Identities: []Identity{
			{
				Role: "owner",
				AuthMethods: []AuthMethod{
					{Type: AuthMethodTypeAddress, Value: "GAUZVLZUTB3SE4MTQ6CFAQXOCMMVCG6LCUZNPM4ZS5X3HZ4BZB4RJM2S"},
					{Type: AuthMethodTypePhoneNumber, Value: "+30000000000"},
					{Type: AuthMethodTypeEmail, Value: "user3@example.com"},
				},
			},
		},
	}

	err = store.Update(b)
	require.NoError(t, err)

	updatedAcc, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, b, updatedAcc)

	// Check the account row has not been changed.
	{
		type row struct {
			ID      int64  `db:"id"`
			Address string `db:"address"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT id, address FROM accounts`)
		require.NoError(t, err)
		wantRows := []row{
			{
				ID:      1,
				Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
			},
		}
		assert.Equal(t, wantRows, rows)
	}

	// Check the identity rows have been updated.
	{
		type row struct {
			AccountID int64  `db:"account_id"`
			ID        int64  `db:"id"`
			Role      string `db:"role"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT account_id, id, role FROM identities`)
		require.NoError(t, err)
		wantRows := []row{
			{
				AccountID: 1,
				ID:        3,
				Role:      "owner",
			},
		}
		assert.Equal(t, wantRows, rows)
	}

	// Check the auth method rows have been updated.
	{
		type row struct {
			AccountID  int64  `db:"account_id"`
			IdentityID int64  `db:"identity_id"`
			ID         int64  `db:"id"`
			Type       string `db:"type_"`
			Value      string `db:"value"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT account_id, identity_id, id, type_, value FROM auth_methods ORDER BY id`)
		require.NoError(t, err)
		wantRows := []row{
			{
				AccountID:  1,
				IdentityID: 3,
				ID:         7,
				Type:       "diamcircle_address",
				Value:      "GAUZVLZUTB3SE4MTQ6CFAQXOCMMVCG6LCUZNPM4ZS5X3HZ4BZB4RJM2S",
			},
			{
				AccountID:  1,
				IdentityID: 3,
				ID:         8,
				Type:       "phone_number",
				Value:      "+30000000000",
			},
			{
				AccountID:  1,
				IdentityID: 3,
				ID:         9,
				Type:       "email",
				Value:      "user3@example.com",
			},
		}
		assert.Equal(t, wantRows, rows)
	}
}

func TestUpdate_removeIdentities(t *testing.T) {
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

	// Remove the identities on the account
	b := Account{
		Address: address,
	}

	err = store.Update(b)
	require.NoError(t, err)

	updatedAcc, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, b, updatedAcc)

	{
		type row struct {
			AccountID int64  `db:"account_id"`
			ID        int64  `db:"id"`
			Role      string `db:"role"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT account_id, id, role FROM identities ORDER BY id`)
		require.NoError(t, err)
		assert.Equal(t, []row{}, rows)
	}
	{
		type row struct {
			AccountID  int64  `db:"account_id"`
			IdentityID int64  `db:"identity_id"`
			ID         int64  `db:"id"`
			Type       string `db:"type_"`
			Value      string `db:"value"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT account_id, identity_id, id, type_, value FROM auth_methods ORDER BY id`)
		require.NoError(t, err)
		assert.Equal(t, []row{}, rows)
	}
}

func TestUpdate_noAuthMethods(t *testing.T) {
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

	b := Account{
		Address: address,
		Identities: []Identity{
			{
				Role: "sender",
			},
			{
				Role: "receiver",
			},
		},
	}

	err = store.Update(b)
	require.NoError(t, err)

	updatedAcc, err := store.Get(address)
	require.NoError(t, err)
	assert.Equal(t, b, updatedAcc)

	// check there is no row in auth_methods
	type row struct {
		AccountID  int64  `db:"account_id"`
		IdentityID int64  `db:"identity_id"`
		ID         int64  `db:"id"`
		Type       string `db:"type_"`
		Value      string `db:"value"`
	}
	rows := []row{}
	err = session.Select(&rows, `SELECT account_id, identity_id, id, type_, value FROM auth_methods`)
	require.NoError(t, err)
	assert.Equal(t, []row{}, rows)
}

func TestUpdate_notFound(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	a := Account{
		Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
	}

	err := store.Update(a)
	assert.Equal(t, ErrNotFound, err)
}

func TestUpdate_notFound_properlyClosesDBConnections(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()
	session.SetMaxIdleConns(1)
	session.SetMaxOpenConns(1)

	store := DBStore{
		DB: session,
	}

	a := Account{
		Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
	}

	for range [2]int{} {
		// If the database transaction is not being properly closed when
		// returning an error, the execution will get stuck in the following
		// line of code when the `Update` method tries to start a new DB
		// transaction through `s.DB.Beginx()`:
		err := store.Update(a)
		assert.Equal(t, ErrNotFound, err)
	}
}

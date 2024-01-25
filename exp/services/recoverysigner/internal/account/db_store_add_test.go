package account

import (
	"testing"

	"go/exp/services/recoverysigner/internal/db/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

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
	err := store.Add(a)
	require.NoError(t, err)

	// Check the account row has been added.
	{
		type row struct {
			ID      int64  `db:"id"`
			Address string `db:"address"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT id, address FROM accounts ORDER BY id`)
		require.NoError(t, err)
		wantRows := []row{
			{
				ID:      1,
				Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
			},
		}
		assert.Equal(t, wantRows, rows)
	}

	// Check the identity rows have been added.
	{
		type row struct {
			AccountID int64  `db:"account_id"`
			ID        int64  `db:"id"`
			Role      string `db:"role"`
		}
		rows := []row{}
		err = session.Select(&rows, `SELECT account_id, id, role FROM identities ORDER BY id`)
		require.NoError(t, err)
		wantRows := []row{
			{
				AccountID: 1,
				ID:        1,
				Role:      "sender",
			},
			{
				AccountID: 1,
				ID:        2,
				Role:      "receiver",
			},
		}
		assert.Equal(t, wantRows, rows)
	}

	// Check the auth method rows have been added.
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
				IdentityID: 1,
				ID:         1,
				Type:       "diamcircle_address",
				Value:      "GD4NGMOTV4QOXWA6PGPIGVWZYMRCJAKLQJKZIP55C5DGB3GBHHET3YC6",
			},
			{
				AccountID:  1,
				IdentityID: 1,
				ID:         2,
				Type:       "phone_number",
				Value:      "+10000000000",
			},
			{
				AccountID:  1,
				IdentityID: 1,
				ID:         3,
				Type:       "email",
				Value:      "user1@example.com",
			},
			{
				AccountID:  1,
				IdentityID: 2,
				ID:         4,
				Type:       "diamcircle_address",
				Value:      "GBJCOYGKIJYX3VUEOZ6GVMFP522UO4OEBI5KB5HHWZAZ2DEJTHS6VOHP",
			},
			{
				AccountID:  1,
				IdentityID: 2,
				ID:         5,
				Type:       "phone_number",
				Value:      "+20000000000",
			},
			{
				AccountID:  1,
				IdentityID: 2,
				ID:         6,
				Type:       "email",
				Value:      "user2@example.com",
			},
		}
		assert.Equal(t, wantRows, rows)
	}
}

func TestAdd_conflict(t *testing.T) {
	db := dbtest.Open(t)
	session := db.Open()

	store := DBStore{
		DB: session,
	}

	a := Account{
		Address: "GCLLT3VG4F6EZAHZEBKWBWV5JGVPCVIKUCGTY3QEOAIZU5IJGMWCT2TT",
	}

	err := store.Add(a)
	require.NoError(t, err)

	err = store.Add(a)
	assert.Equal(t, ErrAlreadyExists, err)
}

func TestAdd_conflict_properlyClosesDBConnections(t *testing.T) {
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

	err := store.Add(a)
	require.NoError(t, err)

	for range [2]int{} {
		// If the database transaction is not being properly closed when
		// returning an error, the execution will get stuck in the following
		// line of code when the `Add` method tries to start a new DB
		// transaction through `s.DB.Beginx()`:
		err = store.Add(a)
		require.Equal(t, ErrAlreadyExists, err)
	}
}

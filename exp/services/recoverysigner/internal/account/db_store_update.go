package account

import (
	"database/sql"

	"github.com/lib/pq"
)

func (s *DBStore) Update(a Account) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var accountID int64
	// Delete an identity will delete the associated auth methods because of the ON DELETE CASCADE reference.
	// https://go/blob/b3e0a353a901ce0babad5b4953330e55f2c674a1/exp/services/recoverysigner/internal/db/dbmigrate/migrations/20200311000002-create-auth-methods.sql#L11
	err = tx.Get(&accountID, `
		WITH deleted_identities AS (
			DELETE FROM identities
			USING accounts
			WHERE identities.account_id = accounts.id AND accounts.address = $1
			RETURNING identities.account_id AS account_id
		)
		SELECT DISTINCT account_id FROM deleted_identities
	`, a.Address)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	for _, i := range a.Identities {
		var authTypes, authValues pq.StringArray
		for _, m := range i.AuthMethods {
			authTypes = append(authTypes, string(m.Type))
			authValues = append(authValues, m.Value)
		}
		_, err = tx.Exec(`
			WITH new_identity AS (
				INSERT INTO identities (account_id, role)
				VALUES ($1, $2)
				RETURNING account_id, id
			)
			INSERT INTO auth_methods (account_id, identity_id, type_, value)
			SELECT account_id, id, unnest($3::auth_method_type[]), unnest($4::text[])
			FROM new_identity
		`, accountID, i.Role, authTypes, authValues)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

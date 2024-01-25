package account

import "github.com/lib/pq"

func (s *DBStore) Add(a Account) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	accountID := int64(0)
	err = tx.Get(&accountID, `
		INSERT INTO accounts (address)
		VALUES ($1)
		RETURNING id
	`, a.Address)
	if err != nil {
		// 23505 is the PostgreSQL error for Unique Violation.
		// See https://www.postgresql.org/docs/9.2/errcodes-appendix.html.
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrAlreadyExists
		}
		return err
	}

	for _, i := range a.Identities {
		identityID := int64(0)
		err = tx.Get(&identityID, `
			INSERT INTO identities (account_id, role)
			VALUES ($1, $2)
			RETURNING id
		`, accountID, i.Role)
		if err != nil {
			return err
		}

		for _, m := range i.AuthMethods {
			_, err = tx.Exec(`
				INSERT INTO auth_methods (account_id, identity_id, type_, value)
				VALUES ($1, $2, $3, $4)
			`, accountID, identityID, m.Type, m.Value)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

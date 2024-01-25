package account

func (s *DBStore) Get(address string) (Account, error) {
	accounts, err := s.getAccounts("accounts.address = $1", address)
	if err != nil {
		return Account{}, err
	}

	// There should only ever be at most one account due to the database
	// constraint where an address must be unique.

	if len(accounts) == 0 {
		return Account{}, ErrNotFound
	}

	return accounts[0], nil
}

func (s *DBStore) getAccounts(where string, args ...interface{}) ([]Account, error) {
	query := `SELECT
			accounts.id AS account_id,
			accounts.address AS account_address,
			identities.id AS identity_id,
			identities.role AS identity_role,
			auth_methods.type_ AS auth_method_type,
			auth_methods.value AS auth_method_value
		FROM accounts
		LEFT JOIN identities ON identities.account_id = accounts.id
		LEFT JOIN auth_methods ON auth_methods.identity_id = identities.id
		WHERE ` + where + `
		ORDER BY accounts.id, identities.id, auth_methods.id`

	rows, err := s.DB.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	accountIndexByAccountID := map[int64]int{}
	identityIndexByIdentityID := map[int64]int{}

	for rows.Next() {
		var r struct {
			AccountID       int64   `db:"account_id"`
			AccountAddress  string  `db:"account_address"`
			IdentityID      *int64  `db:"identity_id"`
			IdentityRole    *string `db:"identity_role"`
			AuthMethodType  *string `db:"auth_method_type"`
			AuthMethodValue *string `db:"auth_method_value"`
		}
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}

		accountIndex, ok := accountIndexByAccountID[r.AccountID]
		if !ok {
			a := Account{Address: r.AccountAddress}
			accounts = append(accounts, a)
			accountIndex = len(accounts) - 1
			accountIndexByAccountID[r.AccountID] = accountIndex
		}
		a := accounts[accountIndex]

		// IdentityID and IdentityRole will be nil if the LEFT JOIN results in
		// an account row that joins to no identities.
		if r.IdentityID != nil && r.IdentityRole != nil {
			identityID := *r.IdentityID
			identityRole := *r.IdentityRole

			identityIndex, ok := identityIndexByIdentityID[identityID]
			if !ok {
				i := Identity{Role: identityRole}
				a.Identities = append(a.Identities, i)
				identityIndex = len(a.Identities) - 1
				identityIndexByIdentityID[identityID] = identityIndex
			}
			i := a.Identities[identityIndex]

			// AuthMethodType and AuthMethodValue will be nil if the LEFT JOIN
			// results in an account/identity row that joins to no auth
			// methods.
			if r.AuthMethodType != nil && r.AuthMethodValue != nil {
				authMethodType := *r.AuthMethodType
				authMethodValue := *r.AuthMethodValue

				m := AuthMethod{
					Type:  AuthMethodType(authMethodType),
					Value: authMethodValue,
				}
				i.AuthMethods = append(i.AuthMethods, m)
			}

			a.Identities[identityIndex] = i
		}

		accounts[accountIndex] = a
	}

	return accounts, nil
}

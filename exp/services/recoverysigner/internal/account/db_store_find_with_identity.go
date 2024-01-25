package account

import "github.com/lib/pq"

func (s *DBStore) FindWithIdentityAuthMethod(t AuthMethodType, value string) ([]Account, error) {
	query := `SELECT account_id
		FROM auth_methods
		WHERE type_ = $1
		AND value = $2`
	accountIDs := []int64{}
	err := s.DB.Select(&accountIDs, query, t, value)
	if err != nil {
		return nil, err
	}

	accounts, err := s.getAccounts(
		`accounts.id = ANY($1::bigint[])`,
		pq.Int64Array(accountIDs),
	)
	if err != nil {
		return []Account{}, err
	}

	return accounts, nil
}

func (s *DBStore) FindWithIdentityAddress(address string) ([]Account, error) {
	return s.FindWithIdentityAuthMethod(AuthMethodTypeAddress, address)
}

func (s *DBStore) FindWithIdentityEmail(email string) ([]Account, error) {
	return s.FindWithIdentityAuthMethod(AuthMethodTypeEmail, email)
}

func (s *DBStore) FindWithIdentityPhoneNumber(phoneNumber string) ([]Account, error) {
	return s.FindWithIdentityAuthMethod(AuthMethodTypePhoneNumber, phoneNumber)
}

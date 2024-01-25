package account

func (s *DBStore) Count() (int, error) {
	count := int(0)
	err := s.DB.Get(&count, `SELECT COUNT(*) FROM accounts`)
	if err != nil {
		return 0, err
	}
	return count, nil
}

package db

// NoRowsError is returned when an insert is attempted without providing any
// values to insert.
type NoRowsError struct {
}

func (err *NoRowsError) Error() string {
	return "no rows provided to insert"
}

var _ error = &NoRowsError{}

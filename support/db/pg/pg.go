package pg

import (
	"go/support/errors"

	"github.com/lib/pq"
)

func IsUniqueViolation(err error) bool {
	switch pgerr := errors.Cause(err).(type) {
	case *pq.Error:
		return string(pgerr.Code) == "23505"
	default:
		return false
	}
}

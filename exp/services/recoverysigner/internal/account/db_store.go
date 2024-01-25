package account

import (
	"github.com/jmoiron/sqlx"
)

type DBStore struct {
	DB *sqlx.DB
}

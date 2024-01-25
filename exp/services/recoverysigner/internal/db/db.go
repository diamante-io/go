package db

import (
	"github.com/jmoiron/sqlx"
)

func Open(dataSourceName string) (*sqlx.DB, error) {
	return sqlx.Open("postgres", dataSourceName)
}

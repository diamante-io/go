package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Open(dataSourceName string) (*sqlx.DB, error) {
	return sqlx.Open("postgres", dataSourceName)
}

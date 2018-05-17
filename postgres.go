package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func (db *Database) connectPostgres() (*sql.DB, error) {
	return sql.Open("postgres", db.Dsn)
}

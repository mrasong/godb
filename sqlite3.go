package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) connectSqlite3() (*sql.DB, error) {
	return sql.Open("sqlite3", db.Dsn)
}

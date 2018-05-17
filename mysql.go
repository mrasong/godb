package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) connectMysql() (*sql.DB, error) {
	return sql.Open("mysql", db.Dsn)
}

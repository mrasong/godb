package db

import (
	"fmt"
	"strconv"
	"strings"

	"database/sql"
)

// connect dirver
func (db *Database) connect() (*sql.DB, error) {
	switch db.Driver {

	case "sqlite", "sqlite3":
		return db.connectSqlite3()

	case "postgres":
		return db.connectPostgres()

	default:
		return db.connectMysql()
	}
}

// reset instance vars
func (db *Database) reset() *Database {
	var vars _vars
	db.Vars = vars
	return db
}

// from table
func (db *Database) From(table string) *Database {
	db.reset()
	db.Vars.Table = table

	return db
}

// fields
func (db *Database) Fields(fields string) *Database {
	db.Vars.Fields = fields
	if db.Vars.Fields == "" {
		db.Vars.Fields = "*"
	}

	return db
}

// where
func (db *Database) Where(where string, bind []interface{}) *Database {
	db.Vars.Where = append(db.Vars.Where, where)
	db.Vars.Bind = append(db.Vars.Bind, bind...)

	return db
}

// order
func (db *Database) Order(order string) *Database {
	if order != "" {
		db.Vars.Order = fmt.Sprintf("ORDER BY %s", order)
	}

	return db
}

// limit
func (db *Database) Limit(limit int64) *Database {
	if limit > 0 {
		db.Vars.Limit = fmt.Sprintf("LIMIT %d", limit)
	}

	return db
}

// offset
func (db *Database) Offset(offset int64) *Database {
	if offset > 0 {
		db.Vars.Offset = fmt.Sprintf("LIMIT %d", offset)
	}
	return db
}

func filter(value interface{}) string {
	v := ""
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		v = fmt.Sprintf("%d", value)

	case float32, float64:
		v = fmt.Sprintf("%f", value)

	case string:
		v = fmt.Sprintf("%s", strconv.Quote(fmt.Sprintf("%s", value)))
	}

	return v
}

func (db *Database) buildSelect() *Database {
	if db.Vars.Fields == "" {
		db.Vars.Fields = "*"
	}

	where := ""
	if len(db.Vars.Where) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(db.Vars.Where, " AND "))
	}

	db.Vars.Query = strings.Join([]string{
		"SELECT",
		db.Vars.Fields,
		"FROM",
		db.Vars.Table,
		where,
		db.Vars.Order,
		db.Vars.Limit,
		db.Vars.Offset,
	}, " ")

	return db
}

func (db *Database) buildInsert(data map[string]interface{}) *Database {
	var fields []string
	var values []string

	for k, v := range data {
		fields = append(fields, k)
		values = append(values, filter(v))
	}

	db.Vars.Query = strings.Join([]string{
		"INSERT INTO",
		db.Vars.Table,
		fmt.Sprintf(
			"(%s) VALUES (%s)",
			strings.Join(fields, ","),
			strings.Join(values, ","),
		),
	}, " ")

	return db
}

func (db *Database) buildUpdate(data map[string]interface{}) *Database {
	var values []string

	for k, v := range data {
		values = append(values, fmt.Sprintf("%s = %s", k, filter(v)))
	}

	where := ""
	if len(db.Vars.Where) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(db.Vars.Where, " AND "))
	}

	db.Vars.Query = strings.Join([]string{
		"UPDATE",
		db.Vars.Table,
		"SET",
		strings.Join(values, ", "),
		where,
		db.Vars.Order,
		db.Vars.Limit,
		db.Vars.Offset,
	}, " ")

	return db
}

func (db *Database) buildDelete() *Database {
	where := ""
	if len(db.Vars.Where) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(db.Vars.Where, " AND "))
	}

	db.Vars.Query = strings.Join([]string{
		"DELETE FROM",
		db.Vars.Table,
		where,
		db.Vars.Order,
		db.Vars.Limit,
		db.Vars.Offset,
	}, " ")

	return db
}

func (db *Database) Query(sql string, bind []interface{}) (*sql.Rows, error) {
	db.reset()
	db.Vars.Query = sql
	db.Vars.Bind = append(db.Vars.Bind, bind...)

	// db
	conn, err := db.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// get list
	res, err := conn.Query(db.Vars.Query, db.Vars.Bind...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	return res, nil
}

func (db *Database) Insert(data map[string]interface{}) (int64, error) {
	// connect db
	conn, err := db.buildInsert(data).connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// prepare query
	insert, err := conn.Prepare(db.Vars.Query)
	if err != nil {
		return 0, err
	}
	defer insert.Close()

	// exec query
	rs, err := insert.Exec()
	if err != nil {
		return 0, err
	}

	return rs.LastInsertId()
}

func (db *Database) Update(data map[string]interface{}) (int64, error) {
	// db
	conn, err := db.buildUpdate(data).connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// get list
	update, err := conn.Prepare(db.Vars.Query)
	if err != nil {
		return 0, err
	}
	defer update.Close()

	rs, err := update.Exec(db.Vars.Bind...)
	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}

func (db *Database) Delete() (int64, error) {
	// db
	conn, err := db.buildDelete().connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// get list
	delete, err := conn.Prepare(db.Vars.Query)
	if err != nil {
		return 0, err
	}
	defer delete.Close()

	rs, err := delete.Exec(db.Vars.Bind...)
	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}

func (db *Database) Count() int64 {
	// db
	conn, err := db.Fields("COUNT(*)").buildSelect().connect()
	if err != nil {
		return 0
	}
	defer conn.Close()

	var count int64 = 0
	// query
	err = conn.QueryRow(db.Vars.Query, db.Vars.Bind...).Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

func (db *Database) Find() ([]interface{}, error) {
	return db.FindAll()
}
func (db *Database) FindAll() ([]interface{}, error) {
	// list
	list := []interface{}{}

	// db
	conn, err := db.buildSelect().connect()
	if err != nil {
		return list, err
	}
	defer conn.Close()

	// get list
	res, err := conn.Query(db.Vars.Query, db.Vars.Bind...)
	if err != nil {
		return list, err
	}
	defer res.Close()

	// columns
	columns, err := res.Columns()
	if err != nil {
		return list, err
	}

	// fields and values interface
	fields := make([]interface{}, len(columns))
	values := make([]sql.RawBytes, len(columns))
	for i := range fields {
		fields[i] = &values[i]
	}

	// list
	for res.Next() {
		err = res.Scan(fields...)
		if err != nil {
			continue
		}

		row := make(map[string]interface{})
		for k, v := range values {
			row[columns[k]] = string(v)
		}

		// append to list
		list = append(list, row)
	}

	return list, nil
}

func (db *Database) FindOne() (map[string]interface{}, error) {
	// db
	conn, err := db.buildSelect().connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// get list
	res, err := conn.Query(db.Vars.Query, db.Vars.Bind...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// columns
	columns, err := res.Columns()
	if err != nil {
		return nil, err
	}

	// fields and values interface
	fields := make([]interface{}, len(columns))
	values := make([]string, len(columns))
	for i := range fields {
		fields[i] = &values[i]
	}

	// get list
	err = conn.QueryRow(db.Vars.Query, db.Vars.Bind...).Scan(fields...)
	if err != nil {
		return nil, err
	}

	// rs
	rs := make(map[string]interface{})
	for k, v := range values {
		if intval, err := strconv.Atoi(v); err == nil {
			rs[columns[k]] = intval
		} else {
			rs[columns[k]] = v
		}
	}

	return rs, nil
}

func (db *Database) FindFirst() (map[string]interface{}, error) {
	return db.FindOne()
}

func (db *Database) findFirst() (map[string]interface{}, error) {
	return db.FindOne()
}

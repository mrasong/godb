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

// From("table")
// From("table as a")
func (db *Database) From(table string) *Database {
	db.reset()
	db.Vars.Table = table

	return db
}

// Join("table b", "LEFT")
// Join("table b", "left")
// Join("table c", "right")
// Join("table d", "Inner")
func (db *Database) Join(table string, join string) *Database {
	if table == "" {
		return db
	}

	join = strings.ToUpper(join)
	switch join {
	case "LEFT", "RIGHT", "INNER":
		join = join

	default:
		join = ""

	}
	db.Vars.Join = fmt.Sprintf("%s JOIN %s", join, table)

	return db
}

// On("table.kid = join.id")
// On("a.uid = user.id")
func (db *Database) On(condition string) *Database {
	if condition != "" {
		db.Vars.On = fmt.Sprintf("ON %s", condition)
	}
	return db
}

// Fields("id, name, age")
// Fields("count(*) as count")
// Fields("user.id, company.name")
func (db *Database) Fields(fields string) *Database {
	db.Vars.Fields = fields
	if db.Vars.Fields == "" {
		db.Vars.Fields = "*"
	}

	return db
}

// type Bind []interface{}
// Where("id = ?", Bind{1})
// Where("id = ? and name = ?", Bind{1, "go"})
// Where("name like ?", Bind{"%go%"})
// Where("a.id = ? OR b.name = ?", Bind{1, "go"})
func (db *Database) Where(where string, bind []interface{}) *Database {
	db.Vars.Where = append(db.Vars.Where, where)
	db.Vars.Bind = append(db.Vars.Bind, bind...)

	return db
}

// Group("id")
func (db *Database) Group(group string) *Database {
	if group != "" {
		db.Vars.Group = fmt.Sprintf("GROUP BY %s", group)
	}

	return db
}

// Having("id > 1")
func (db *Database) Having(having string) *Database {
	if having != "" {
		db.Vars.Having = fmt.Sprintf("HAVING %s", having)
	}

	return db
}

// Order("id DESC")
// Order("created_at DESC, id DESC")
func (db *Database) Order(order string) *Database {
	if order != "" {
		db.Vars.Order = fmt.Sprintf("ORDER BY %s", order)
	}

	return db
}

// Limit(3)
func (db *Database) Limit(limit int64) *Database {
	if limit > 0 {
		db.Vars.Limit = fmt.Sprintf("LIMIT %d", limit)
	}

	return db
}

// Offset(100)
func (db *Database) Offset(offset int64) *Database {
	if offset > 0 {
		db.Vars.Offset = fmt.Sprintf("OFFSET %d", offset)
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
		v = strconv.Quote(fmt.Sprintf("%s", value))

	case []byte:
		fmt.Println("byte")
		v = string(value.([]byte)[:])
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
		db.Vars.Join,
		db.Vars.On,
		db.Vars.Group,
		db.Vars.Having,
		where,
		db.Vars.Order,
		db.Vars.Limit,
		db.Vars.Offset,
	}, " ")

	return db
}

func (db *Database) buildInsert(data map[string]interface{}) *Database {
	// k => v
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
	// update data
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
	// where
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
	// reset vars
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

func (db *Database) SetField(field string, value interface{}) (int64, error) {
	// db
	conn, err := db.buildUpdate(map[string]interface{}{
		field: value,
	}).connect()
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

func (db *Database) Count() (int64, error) {
	// db
	conn, err := db.Fields("COUNT(*) as count").buildSelect().connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var count int64 = 0
	// query
	err = conn.QueryRow(db.Vars.Query, db.Vars.Bind...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return int64(count), nil
}

func (db *Database) Find() ([]map[string]interface{}, error) {
	// find all
	return db.FindAll()
}

func (db *Database) FindAll() ([]map[string]interface{}, error) {
	// list
	list := []map[string]interface{}{}

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
	for k, _ := range fields {
		fields[k] = &values[k]
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

func (db *Database) FindFirst() (map[string]interface{}, error) {
	// find one
	return db.FindOne()
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
	values := make([]interface{}, len(columns))
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
		if _v, ok := v.([]byte); ok {
			rs[columns[k]] = string(_v[:])
		} else {
			rs[columns[k]] = v
		}
	}

	return rs, nil
}

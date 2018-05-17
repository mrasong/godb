package db

type _vars struct {
	Table string
	Query string

	Fields string
	Where  []string
	Bind   []interface{}

	Order  string
	Limit  string
	Offset string

	Condition []interface{}
}

type Database struct {
	Driver string
	Dsn    string

	Vars _vars
}

func DB(driver string, dsn string) Database {
	return Database{Driver: driver, Dsn: dsn}
}

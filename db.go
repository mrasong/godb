package db

type _vars struct {
	Table  string
	Join   string
	On     string
	Group  string
	Having string
	Query  string

	Fields     string
	Where      []string
	Bind       []interface{}
	BindInsert []interface{}
	BindUpdate []interface{}

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

func New(driver string, dsn string) Database {
	return Database{Driver: driver, Dsn: dsn}
}

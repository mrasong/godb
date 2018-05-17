## simple database tool for go 

Support 

 - mysql  [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
 - postgresql [github.com/lib/pq](https://github.com/lib/pq)
 - sqlite3 [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)



### Instantiation 


```
import(
	"github.com/mrasong/godb"
)

const (

	DB_DRIVER = "mysql"
	DB_DSN    = "user:password@tcp(localhost:3306)/go?charset=utf8mb4"

	// DB_DRIVER = "postgres"
	// DB_DSN    = "postgresql://user:password@localhost/go"

	// DB_DRIVER = "sqlite3"
	// DB_DSN    = "./sqlite3.db"

)


var Db db.Database = db.Database{Driver: DB_DRIVER, Dsn: DB_DSN}
// var Db := db.DB(DB_DRIVER, DB_DSN)

```


### Functions

```
func (db *Database) Count() (int64, error)

func (db *Database) Delete() (int64, error)

func (db *Database) Fields(fields string) *Database

func (db *Database) Find() ([]interface{}, error)

func (db *Database) FindAll() ([]interface{}, error)

func (db *Database) FindFirst() (map[string]interface{}, error)

func (db *Database) FindOne() (map[string]interface{}, error)

func (db *Database) From(table string) *Database

func (db *Database) Insert(data map[string]interface{}) (int64, error)

func (db *Database) Limit(limit int64) *Database

func (db *Database) Offset(offset int64) *Database

func (db *Database) Order(order string) *Database

func (db *Database) Query(sql string, bind []interface{}) (*sql.Rows, error)

func (db *Database) Update(data map[string]interface{}) (int64, error)

func (db *Database) Where(where string, bind []interface{}) *Database
```


### Example

```
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mrasong/godb"
)

type bind []interface{}
type data map[string]interface{}

const (
	DB_DRIVER = "mysql"
	DB_DSN    = "user:password@tcp(host:port)/database?charset=utf8mb4"

	// DB_DRIVER = "postgres"
	// DB_DSN    = "postgresql://localhost/go"

	// DB_DRIVER = "sqlite3"
	// DB_DSN    = "./sqlite3.db"
)

var Db db.Database = db.Database{Driver: DB_DRIVER, Dsn: DB_DSN}

// var Db = db.DB(DB_DRIVER, DB_DSN)

func JSON(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json error:", data)
		return
	}

	fmt.Println(string(b))
	return
}

func FindAll() {
	findAll, err := Db.From("test").
		// Fields("id,name").
		Where("name like ?", bind{"%go%"}).
		Order("id desc").
		Limit(10).
		FindAll()

	if err != nil {
		JSON(err)
		return
	}
	JSON(findAll)
}

func FindOne() {
	findOne, err := Db.From("test").
		// Fields("id,name").
		Where("id = ?", bind{1}).
		FindOne()

	if err != nil {
		JSON(err)
		return
	}
	JSON(findOne)
}

func Delete() {
	delete, err := Db.
		From("test").
		Where("id = ?", bind{1}).
		Delete()

	if err != nil {
		JSON(err)
		return
	}
	JSON(delete)
}

func Insert() {
	insert, err := Db.
		From("test").
		Insert(data{
			"created_at": time.Now().Unix(),
			"name":       "hi go",
			"text":       `this should be a long text`,
		})

	if err != nil {
		JSON(err)
		return
	}
	JSON(insert)
}

func Update() {
	update, err := Db.
		From("test").
		Where("id = ?", bind{1}).
		Update(data{
			"created_at": time.Now().Unix(),
			"name":       "hello go",
			"text":       `uint is an "unsigned" integer type that is at least '32' bits in size. It is a distinct type, however, and not an alias for, say, uint32.`,
		})

	if err != nil {
		JSON(err)
		return
	}
	JSON(update)
}

func Count() {
	count := Db.
		From("test").
		Where("id > ?", bind{1}).
		Count()

	JSON(count)
}

func Query() {
	query, err := Db.Query(`CREATE TABLE test ...`, nil)
	if err != nil {
		JSON(err)
		return
	}
	JSON(query)

	query2, err := Db.Query(`select * from test where id = ?`, bind{1})
	if err != nil {
		JSON(err)
		return
	}
	JSON(query2)
}

func main() {
	FindAll()
}

```

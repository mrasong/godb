## simple database tool for go 

Support 

 - mysql  [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
 - postgresql [github.com/lib/pq](https://github.com/lib/pq)
 - sqlite3 [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)



### Instantiation 


```go
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
// var Db := db.New(DB_DRIVER, DB_DSN)

```


### Functions

```go
func (db *Database) Count() (int64, error)

func (db *Database) Delete() (int64, error)

func (db *Database) Fields(fields string) *Database
    Fields("id, name, age") 
    Fields("count(*) as count") 
    Fields("user.id, company.name")

func (db *Database) Find() ([]map[string]interface{}, error)

func (db *Database) FindAll() ([]map[string]interface{}, error)

func (db *Database) FindFirst() (map[string]interface{}, error)

func (db *Database) FindOne() (map[string]interface{}, error)

func (db *Database) From(table string) *Database
    From("table") 
    From("table as a")

func (db *Database) Group(group string) *Database
    Group("id")

func (db *Database) Having(having string) *Database
    Having("id > 1")

func (db *Database) Insert(data map[string]interface{}) (int64, error)

func (db *Database) Join(table string, join string) *Database
    Join("table b", "LEFT") 
    Join("table b", "left") 
    Join("table c", "right")
    Join("table d", "Inner")

func (db *Database) Limit(limit int64) *Database
    Limit(3)

func (db *Database) Offset(offset int64) *Database
    Offset(100)

func (db *Database) On(condition string) *Database
    On("table.kid = join.id") 
    On("a.uid = user.id")

func (db *Database) Order(order string) *Database
    Order("id DESC") 
    Order("created_at DESC, id DESC")

func (db *Database) SetField(field string, value interface{}) (int64, error)
    SetField("views", 10086)
    SetField("name", "go")

func (db *Database) Query(sql string, bind []interface{}) (*sql.Rows, error)

func (db *Database) Update(data map[string]interface{}) (int64, error)

func (db *Database) Where(where string, bind []interface{}) *Database
    type Bind []interface{} 
    Where("id = ?", Bind{1}) 
    Where("id = ? and name = ?", Bind{1, "go"}) 
    Where("name like ?", Bind{"%go%"}) 
    Where("a.id = ? OR b.name = ?", Bind{1, "go"})
```


### Example

```go
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

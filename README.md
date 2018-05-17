## simple database tool for go 


support 

 - mysql
 - postgresql
 - sqlite3



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


### Query (*sql.Rows, error)

```
	query, err := Db.Query(`CREATE TABLE test ...`, nil)
	query2, err := Db.Query(`SELECT * FROM test WHERE id = ? or name = ?`, bind{1, "go"})
```


### FindOne (map[string]interface{}, error)

```
	findOne, err := Db.From("test").
		// Fields("id,name").
		Where("id = ? AND name = ?", bind{1, "go"}).
		FindOne()
```


### FindAll ([]interface{}, error)

```
	findAll, err := Db.From("test").
		// Fields("id,name").
		Where("name like ?", bind{"%go%"}).
		Order("id desc").
		Limit(10).
		FindAll()
```


### Insert (int64, error)

```
	insert, err := Db.
		From("test").
		Insert(data{
			"created_at": time.Now().Unix(),
			"name":       "hi go",
			"text":       `this should be a long text`,
		})
```


### Update (int64, error)

```
	update, err := Db.
		From("test").
		Where("id = ?", bind{1}).
		Update(data{
			"created_at": time.Now().Unix(),
			"name":       fmt.Sprintf("the %d", i),
			"text":       `uint is an "unsigned" integer type that is at least '32' bits in size. It is a distinct type, however, and not an alias for, say, uint32.`,
		})
```


### Delete (int64, error)

```
	delete, err := Db.
		From("test").
		Where("id = ?", bind{i}).
		Delete()
```



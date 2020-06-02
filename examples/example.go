package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlxb "github.com/laoqiu/sqlx-builder"
	ex "github.com/laoqiu/sqlx-builder/examples/proto"
)

// Person 用户对象
type Person struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Info    sql.NullString `json:"info"`
	Created *time.Time     `json:"created"`
}

func main() {
	db, err := sqlxb.Connect(
		sqlxb.Driver("mysql"),
		sqlxb.URI("root:123456@tcp(127.0.0.1:3306)/my_app"),
	)
	if err != nil {
		log.Fatal(err)
	}

	schemas := []string{
		`
		CREATE TABLE IF NOT EXISTS my_app.person (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(45) NULL,
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id));
		`,
	}

	for _, s := range schemas {
		if _, err := db.Exec(s); err != nil {
			log.Fatal(err)
		}
	}

	person := &ex.Person{Name: "test name 1"}

	// result, err := sqlxb.NewBuilder(db).Debug(true).Table("person").Insert(person)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(result)

	err = sqlxb.NewBuilder(db).Debug(true).Table("person").One(person)
	if err != nil {
		log.Fatal(err)
	}
}

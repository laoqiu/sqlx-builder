package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	sqlxb "github.com/laoqiu/sqlx-builder"
	ex "github.com/laoqiu/sqlx-builder/examples/proto"
)

func main() {
	db, err := sqlxb.Connect(
		sqlxb.Driver("mysql"),
		sqlxb.URI("root:123456@tcp(127.0.0.1:3306)/my_app"),
	)
	if err != nil {
		log.Fatal(err)
	}
	sqlxb.LoadMapper(db, sqlxb.DefaultMapper)

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

	result, err := sqlxb.NewBuilder(db).Debug(true).Table("person").Insert(person)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	err = sqlxb.NewBuilder(db).Unsafe().Debug(true).Table("person").Fields("id", "name", "create_at").One(person)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(person)
}

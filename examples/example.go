package main

import (
	"database/sql"
	"log"
	"time"

	sqlxb "github.com/laoqiu/sqlx-builder"
	_ "github.com/mattn/go-sqlite3"
)

// Person 用户对象
type Person struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Info    sql.NullString `json:"info"`
	Created *time.Time     `json:"created"`
}

type Place struct {
	PersonID int64  `json:"person_id"`
	Address  string `json:"address"`
}

type PersonPlace struct {
	Person
	Place
}

func main() {
	db, err := sqlxb.Connect(
	// sqlxb.Charset("utf8"),
	// sqlxb.Driver("mysql"),
	// sqlxb.URI("user:password@tcp(127.0.0.1:3306)/hello"),
	)
	if err != nil {
		log.Fatal(err)
	}

	schemas := []string{
		`CREATE TABLE IF NOT EXISTS person (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(20) NOT NULL,
			info VARCHAR(100) DEFAULT '',
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
		`,
		`CREATE TABLE IF NOT EXISTS place (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			person_id INTEGER NOT NULL,
			address VARCHAR(100) DEFAULT '',
			country TEXT DEFAULT '',
			city TEXT DEFAULT '');
		`,
	}

	for _, s := range schemas {
		if _, err := db.Exec(s); err != nil {
			log.Fatal(err)
		}
	}

	person := &Person{
		ID:   1,
		Name: "test",
	}
	result, err := sqlxb.NewBuilder(db).Debug(true).Table("person").Insert(person)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	if err := sqlxb.NewBuilder(db).Debug(true).Table("person").Distinct().One(person); err != nil {
		log.Fatal(err)
	}

	// 链式操作
	data := []PersonPlace{}
	if err := sqlxb.NewBuilder(db).Debug(true).Table("person").Join("place", "person.id = place.person_id").
		Fields("person.id", "person.name", "person.info", "place.address").
		Where("person.name = ?", person.Name).
		All(&data); err != nil {
		log.Fatal(err)
	}
	log.Println(data)
}

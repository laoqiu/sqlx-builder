package main

import (
	"database/sql"
	"log"
	"time"

	sqlxb "github.com/laoqiu/sqlx-builder"
)

// Person 用户对象
type Person struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Address sql.NullString `json:"address"`
	Created *time.Time     `json:"created"`
}

func main() {
	db, _ := sqlxb.Connect() // 返回*sqlx.DB对象
	person := &Person{}
	if err := sqlxb.New(db).Debug(true).Table("person").Distinct().Get(person); err != nil {
		log.Fatal(err)
	}
}

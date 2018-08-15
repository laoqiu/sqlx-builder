package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/laoqiu/sqlx-colt"
)

type Person struct {
	ID      int64  `db:"id" tbl:"PRIMARY_KEY AUTO_INCREMENT"`
	Name    string `tbl:"INDEX"`
	Address sql.NullString
	Created *time.Time
}

func main() {
	var db *sqlx.DB
	db, err := sqlxcolt.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/tms", "utf8mb4", true, 10, 10)
	a := &Person{}
	dbc := sqlxcolt.Query{}
	err = dbc.Bind(db).Table("logs_dispatch").Distinct().
		First(a)
	fmt.Println(err)
}

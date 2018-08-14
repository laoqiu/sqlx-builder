package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/laoqiu/sqlx-query-colt"
)

type Person struct {
	ID      int64  `db:"id" tbl:"PRIMARY_KEY AUTO_INCREMENT"`
	Name    string `tbl:"INDEX"`
	Address sql.NullString
	Created *time.Time
}

func main() {
	var db *sqlx.DB
	db, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/ml_test?charset=utf8mb4&parseTime=true")
	if err != nil {
		fmt.Println(err)
		return
	}
	a := &Person{}
	dbc := xcolt.Query{}
	err = dbc.Bind(db).Table("logs_dispatch").Distinct().
		First(a)
	fmt.Println(err)
}

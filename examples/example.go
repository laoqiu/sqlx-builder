package main

import (
	"database/sql"
	"time"

	"github.com/laoqiu/sqlx-builder"
)

type Person struct {
	ID      int64  `db:"id" tbl:"PRIMARY_KEY AUTO_INCREMENT"`
	Name    string `tbl:"INDEX"`
	Address sql.NullString
	Created *time.Time
}

func main() {
	db, _ := sqlxt.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/tms", "utf8mb4", true, 10, 10)
	person := &Person{}
	query := sqlxt.Table("logs_dispatch").Distinct()
	sqlxt.New(db, query).Get(person)
}

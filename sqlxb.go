package sqlxb

import (
	"database/sql"
	"log"

	// _ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

// Builder 返回sqlx的builder
type Builder struct {
	debug bool
	db    *sqlx.DB
	tx    *sqlx.Tx
	query *Query
}

// New 返回普通模式
func New(db *sqlx.DB) *Builder {
	return &Builder{
		db: db,
	}
}

// NewTx 事务模式
func NewTx(tx *sqlx.Tx) *Builder {
	return &Builder{
		tx: tx,
	}
}

// Debug 设置debug为true
func (b *Builder) Debug(debug bool) *Builder {
	b.debug = debug
	return b
}

// Get 返回单条数据结果
func (b *Builder) Get(dest interface{}) error {
	b.Limit(1)
	query, args, err := b.BuildQuery()
	if err != nil {
		return err
	}
	if b.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	var row *sqlx.Row
	if b.tx != nil {
		row = b.tx.Unsafe().QueryRowx(query, args...)
	} else {
		row = b.db.Unsafe().QueryRowx(query, args...)
	}
	return row.StructScan(dest)
}

// All 返回列表查询结果
func (b *Builder) All(dest interface{}) error {
	var err error
	query, args, err := b.BuildQuery()
	if err != nil {
		return err
	}
	if b.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	if b.tx != nil {
		err = b.tx.Unsafe().Select(dest, query, args...)
	} else {
		err = b.db.Unsafe().Select(dest, query, args...)
	}
	return err
}

// Update 执行 UPDATE 语句
func (b *Builder) Update(data interface{}) (sql.Result, error) {
	return b.Exec("UPDATE", data)
}

// Insert 执行 INSERT 语句
func (b *Builder) Insert(data interface{}) (sql.Result, error) {
	return b.Exec("INSERT", data)
}

// InsertIgnore 执行 INSERT_IGNORE 语句
func (b *Builder) InsertIgnore(data interface{}) (sql.Result, error) {
	return b.Exec("INSERT_IGNORE", data)
}

// InsertOnDuplicateUpdate 执行 INSERT_ON_DUPLICATE_UPDATE 语言
func (b *Builder) InsertOnDuplicateUpdate(data interface{}) (sql.Result, error) {
	return b.Exec("INSERT_ON_DUPLICATE_UPDATE", data)
}

// Delete 执行 DELETE 语言
func (b *Builder) Delete() (sql.Result, error) {
	return b.Exec("DELETE", nil)
}

// Exec 执行sql语句
func (b *Builder) Exec(method string, s interface{}) (sql.Result, error) {
	var err error
	query, args, err := b.BuildExec(method, StructToMap(s))
	if err != nil {
		return nil, err
	}
	if b.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	var result sql.Result
	if b.tx != nil {
		result, err = b.tx.Exec(query, args...)
	} else {
		result, err = b.db.Exec(query, args...)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

package builder

import (

	// _ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// LogTemp 打印日志格式
var LogTemp = "sql output:\nquery: %v\n args: %v\n"

// DB sqlx.DB or sqlx.Tx
type DB interface {
	DriverName() string
	Rebind(query string) string
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
	Exec(query string, args ...interface{}) (sql.Result, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

// Builder 返回sqlx的builder
type Builder struct {
	db    *sqlx.DB
	tx    *sqlx.Tx
	debug bool
	query *Query
}

// NewBuilder return new builder
func NewBuilder(db *sqlx.DB) *Builder {
	return &Builder{db: db, tx: nil, query: nil}
}

// Unsafe Scan Destination Safety
func (b *Builder) Unsafe() *Builder {
	b.db = b.db.Unsafe()
	return b
}

// SetTx 事务支持，如果传nil则退出事务
func (b *Builder) SetTx(tx *sqlx.Tx) *Builder {
	b.tx = tx
	return b
}

// DB 自动选择存储对象(db or tx)
func (b *Builder) DB() DB {
	if b.tx != nil {
		return b.tx
	}
	return b.db
}

// Debug 设置debug值
func (b *Builder) Debug(v bool) *Builder {
	b.debug = v
	return b
}

// One 返回单条数据
func (b *Builder) One(dest interface{}) error {
	query, args, err := b.BuildQuery()
	if err != nil {
		return err
	}
	if b.debug {
		log.Printf(LogTemp, query, args)
	}
	return b.DB().Get(dest, query, args...)
}

// All 返回多条数据
func (b *Builder) All(dest interface{}) error {
	query, args, err := b.BuildQuery()
	if err != nil {
		return err
	}
	if b.debug {
		log.Printf(LogTemp, query, args)
	}
	return b.DB().Select(dest, query, args...)
}

// Update 执行 UPDATE 语句
func (b *Builder) Update(data interface{}) (sql.Result, error) {
	return b._exec("UPDATE", data)
}

// Insert 执行 INSERT 语句
func (b *Builder) Insert(data interface{}) (sql.Result, error) {
	return b._exec("INSERT", data)
}

// InsertIgnore 执行 INSERT_IGNORE 语句
func (b *Builder) InsertIgnore(data interface{}) (sql.Result, error) {
	return b._exec("INSERT_IGNORE", data)
}

// InsertOnDuplicateUpdate 执行 INSERT_ON_DUPLICATE_UPDATE 语言
func (b *Builder) InsertOnDuplicateUpdate(data interface{}) (sql.Result, error) {
	return b._exec("INSERT_ON_DUPLICATE_UPDATE", data)
}

// Delete 执行 DELETE 语言
func (b *Builder) Delete() (sql.Result, error) {
	return b._exec("DELETE", nil)
}

// _exec 执行sql语句
func (b *Builder) _exec(method string, s interface{}) (sql.Result, error) {
	var err error
	query, args, err := b.BuildExec(method, StructToMap(s))
	if err != nil {
		return nil, err
	}
	if b.debug {
		log.Printf(LogTemp, query, args)
	}
	var result sql.Result
	result, err = b.DB().Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

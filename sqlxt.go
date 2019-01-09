package sqlxt

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type Sqlxt struct {
	debug bool
	db    *sqlx.DB
	tx    *sqlx.Tx
	query *Query
}

func New(db *sqlx.DB, q *Query, debug bool) *Sqlxt {
	return &Sqlxt{
		db:    db,
		query: q,
		debug: debug,
	}
}

func NewTx(tx *sqlx.Tx, q *Query, debug bool) *Sqlxt {
	return &Sqlxt{
		tx:    tx,
		query: q,
		debug: debug,
	}
}

func (st *Sqlxt) Get(dest interface{}) error {
	st.query.Limit(1)
	query, args, err := st.query.BuildQuery()
	if err != nil {
		return err
	}
	if st.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	var row *sqlx.Row
	if st.tx != nil {
		row = st.tx.Unsafe().QueryRowx(query, args...)
	} else {
		row = st.db.Unsafe().QueryRowx(query, args...)
	}
	return row.StructScan(dest)
}

func (st *Sqlxt) All(dest interface{}) error {
	var err error
	query, args, err := st.query.BuildQuery()
	if err != nil {
		return err
	}
	if st.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	if st.tx != nil {
		err = st.tx.Unsafe().Select(dest, query, args...)
	} else {
		err = st.db.Unsafe().Select(dest, query, args...)
	}
	return err
}

func (st *Sqlxt) Update(data interface{}) (sql.Result, error) {
	return st.Exec("UPDATE", data)
}

func (st *Sqlxt) Insert(data interface{}) (sql.Result, error) {
	return st.Exec("INSERT", data)
}

func (st *Sqlxt) InsertIgnore(data interface{}) (sql.Result, error) {
	return st.Exec("INSERT_IGNORE", data)
}

func (st *Sqlxt) InsertOnDuplicateUpdate(data interface{}) (sql.Result, error) {
	return st.Exec("INSERT_ON_DUPLICATE_UPDATE", data)
}

func (st *Sqlxt) Delete() (sql.Result, error) {
	return st.Exec("DELETE", nil)
}

func (st *Sqlxt) Exec(method string, s interface{}) (sql.Result, error) {
	var err error
	query, args, err := st.query.BuildExec(method, StructToMap(s))
	if err != nil {
		return nil, err
	}
	if st.debug {
		log.Printf("sql output:\nquery: %v\n args: %v\n", query, args)
	}
	var result sql.Result
	if st.tx != nil {
		result, err = st.tx.Exec(query, args...)
	} else {
		result, err = st.db.Exec(query, args...)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

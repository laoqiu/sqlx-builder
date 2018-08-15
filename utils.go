package sqlxcolt

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/jmoiron/sqlx"
)

func Max(n ...int) int {
	return sort.IntSlice(n)[0]
}

func Min(n ...int) int {
	return sort.IntSlice(n)[len(n)-1]
}

func If(condition bool, v1 interface{}, v2 interface{}) interface{} {
	if condition {
		return v1
	}
	return v2
}

func Filter(iter []string, f func(x string) bool) []string {
	result := []string{}
	for _, i := range iter {
		if f(i) {
			result = append(result, i)
		}
	}
	return result
}

// DefaultMapper 默认的Mapper函数 PersonAddress -> person_address
func DefaultMapper(name string) string {
	var s []byte
	for i, r := range []byte(name) {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
			if i != 0 {
				s = append(s, '_')
			}
		}
		s = append(s, r)
	}
	return string(s)
}

func GetType(v interface{}) string {
	return reflect.TypeOf(v).Name()
}

// Connect 获得数据库连接
func Connect(driver, uri, charset string, parseTime bool, maxOpen, maxIdel int) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driver, fmt.Sprintf("%s?charset=%s&parseTime=%v", uri, charset, parseTime))
	if err != nil {
		return nil, err
	}
	// 配置连接池
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdel)
	return db, nil
}

// LoadMapper 全局替换mapper
func LoadMapper(db *sqlx.DB, mapper func(name string) string) {
	db.MapperFunc(mapper)
	// 使用`sqlx.Named`时生效
	sqlx.NameMapper = mapper
}

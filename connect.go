package sqlxb

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Connect 获得数据库连接
func Connect(opts ...Option) (*sqlx.DB, error) {
	o := NewOptions(opts...)
	if o.Driver == "mysql" {
		// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		o.URI = fmt.Sprintf("%s?charset=%s&parseTime=%v", o.URI, o.Charset, o.ParseTime)
	}
	db, err := sqlx.Connect(o.Driver, o.URI)
	if err != nil {
		return nil, err
	}
	// 配置连接池
	db.SetMaxOpenConns(o.MaxClient)
	db.SetMaxIdleConns(o.MaxClient)
	return db, nil
}

// Max 最大值
func Max(n ...int) int {
	return sort.IntSlice(n)[0]
}

// Min 最小值
func Min(n ...int) int {
	return sort.IntSlice(n)[len(n)-1]
}

// If 判断对象
func If(condition bool, v1 interface{}, v2 interface{}) interface{} {
	if condition {
		return v1
	}
	return v2
}

// Filter 过滤器
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

// GetType 反射取属性名
func GetType(v interface{}) string {
	return reflect.TypeOf(v).Name()
}

// LoadMapper 全局替换mapper
func LoadMapper(db *sqlx.DB, mapper func(name string) string) {
	db.MapperFunc(mapper)
	// 使用`sqlx.Named`时生效
	sqlx.NameMapper = mapper
}

// StructToMap struct转map
func StructToMap(i interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	if i != nil {
		iVal := reflect.ValueOf(i).Elem()
		tp := iVal.Type()
		for i := 0; i < iVal.NumField(); i++ {
			tag := tp.Field(i).Tag.Get("json")
			if len(tag) > 0 {
				name := strings.Split(tag, ",")[0]
				if name != "-" {
					values[name] = iVal.Field(i).Interface()
				}
			}
		}
	}
	return values
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

package sqlxb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Query 返回sql语句主体
type Query struct {
	table    string
	fields   []string
	join     [][]interface{}
	where    [][]interface{}
	distinct bool
	order    string
	group    string
	having   string
	limit    int
	offset   int
	lock     string
}

// Table 生成基本query并赋值builder
func (b *Builder) Table(tablename string) *Builder {
	q := &Query{
		table: tablename,
	}
	b.query = q
	return b
}

// LockInShareMode 共享锁
func (b *Builder) LockInShareMode() *Builder {
	b.query.lock = "LOCK IN SHARE MODE"
	return b
}

// LockForUpdate 写锁
func (b *Builder) LockForUpdate() *Builder {
	b.query.lock = "FOR UPDATE"
	return b
}

// Fields 定义返回字段
func (b *Builder) Fields(fields ...string) *Builder {
	b.query.fields = fields
	return b
}

// AddFields 添加新的返回字段
func (b *Builder) AddFields(fields ...string) *Builder {
	b.query.fields = append(b.query.fields, fields...)
	return b
}

// Select 指向别名Fields
func (b *Builder) Select(fields ...string) *Builder {
	b.Fields(fields...)
	return b
}

// AddSelect 指向别名AddFields
func (b *Builder) AddSelect(fields ...string) *Builder {
	b.AddFields(fields...)
	return b
}

// Join 赋值多表关联join表达式
func (b *Builder) Join(table ...interface{}) *Builder {
	b.query.join = append(b.query.join, []interface{}{"INNER JOIN", table})
	return b
}

// LeftJoin 多表关联leftjoin表达式
func (b *Builder) LeftJoin(table ...interface{}) *Builder {
	b.query.join = append(b.query.join, []interface{}{"LEFT JOIN", table})
	return b
}

// RightJoin 多表关联rightjoin表达式
func (b *Builder) RightJoin(table ...interface{}) *Builder {
	b.query.join = append(b.query.join, []interface{}{"RIGHT JOIN", table})
	return b
}

// UnionJoin 联合查询union join表达式
func (b *Builder) UnionJoin(table ...interface{}) *Builder {
	b.query.join = append(b.query.join, []interface{}{"UNION JOIN", table})
	return b
}

// Where 条件查询
func (b *Builder) Where(query string, args ...interface{}) *Builder {
	b.query.where = append(b.query.where, []interface{}{"AND", query, args})
	return b
}

// GroupBy 分组查询
func (b *Builder) GroupBy(group string) *Builder {
	b.query.group = group
	return b
}

// OrderBy 排序
func (b *Builder) OrderBy(order string) *Builder {
	b.query.order = order
	return b
}

// Distinct 去重
func (b *Builder) Distinct() *Builder {
	b.query.distinct = true
	return b
}

// Limit set limit number
func (b *Builder) Limit(n int) *Builder {
	b.query.limit = n
	return b
}

// Offset set offset number
func (b *Builder) Offset(n int) *Builder {
	b.query.offset = n
	return b
}

func (b *Builder) _parseInsert(data map[string]interface{}) (string, string) {
	var keystr, valstr string
	var key, val []string
	for k, v := range data {
		if len(b.query.fields) == 0 || indexOf(k, b.query.fields) != -1 {
			// 反射找出类型
			switch v.(type) {
			case string:
				key = append(key, k)
				val = append(val, fmt.Sprintf("'%s'", v))
			case int, int8, int32, int64, float32, float64, bool:
				key = append(key, k)
				val = append(val, fmt.Sprintf("%v", v))
			default:
			}
		}
	}
	keystr = strings.Join(key, ", ")
	valstr = strings.Join(val, ", ")
	return keystr, valstr
}

func (b *Builder) _parseUpate(data map[string]interface{}) string {
	var setstr string
	var result []string

	for k, v := range data {
		if len(b.query.fields) == 0 || indexOf(k, b.query.fields) != -1 {
			var _value string
			// 反射找出类型
			switch v.(type) {
			case int, int32, int64:
				_value = fmt.Sprintf("%s = %d", k, v)
			case float32, float64:
				_value = fmt.Sprintf("%s = %f", k, v)
			default:
				_value = fmt.Sprintf("%s = '%s'", k, v)
			}
			result = append(result, _value)
		}
	}

	setstr = strings.Join(result, ", ")
	return setstr
}

// BuildExec 返回需要执行的sql表达式
func (b *Builder) BuildExec(method string, data map[string]interface{}) (string, []interface{}, error) {
	var sqlstr string
	var tablename string

	tablename = b.query.table

	where, args, err := b.parseWhere()
	if err != nil {
		return "", nil, err
	}
	where = If(where == "", "", "WHERE "+where).(string)

	switch method {
	case "INSERT", "INSERT_IGNORE":
		keystr, valstr := b._parseInsert(data)
		ignore := If(method == "INSERT_IGNORE", "IGNORE", "").(string)
		sqlstr = fmt.Sprintf("INSERT %s INTO %s (%s) VALUES (%s)", ignore, tablename, keystr, valstr)
	case "INSERT_ON_DUPLICATE_UPDATE":
		keystr, valstr := b._parseInsert(data)
		setstr := b._parseUpate(data)
		sqlstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s", tablename, keystr, valstr, setstr)
	case "UPDATE":
		setstr := b._parseUpate(data)
		sqlstr = fmt.Sprintf("UPDATE %s SET %s %s", tablename, setstr, where)
	case "DELETE":
		sqlstr = fmt.Sprintf("DELETE FROM %s %s", tablename, where)
	}
	//log.Println("sql output ->", sqlstr)
	return sqlstr, args, nil
}

// BuildQuery 合并query表达式
func (b *Builder) BuildQuery() (string, []interface{}, error) {
	// table
	table := DefaultMapper(b.query.table)
	// join
	join, joinTables, err := b.parseJoin()
	if err != nil {
		return "", nil, err
	}
	// distinct
	distinct := If(b.query.distinct, "DISTINCT", "").(string)
	// fields
	var allFields string
	if len(joinTables) == 0 {
		allFields = "*"
	} else {
		_t := []string{table + ".*"}
		for _, i := range joinTables {
			_t = append(_t, i+".*")
		}
		allFields = strings.Join(_t, ", ")
	}
	fields := If(len(b.query.fields) == 0, allFields, strings.Join(b.query.fields, ",")).(string)
	// where
	where, args, err := b.parseWhere()
	if err != nil {
		return "", nil, err
	}
	// where
	where = If(where == "", "", "WHERE "+where).(string)
	// group
	group := If(b.query.group == "", "", "GROUP BY "+b.query.group).(string)
	// order
	order := If(b.query.order == "", "", "ORDER BY "+b.query.order).(string)
	// having
	having := If(b.query.having == "", "", " HAVING "+b.query.having).(string)
	// limit
	limit := If(b.query.limit == 0, "", fmt.Sprintf("LIMIT %d", b.query.limit)).(string)
	// offset
	offset := If(b.query.offset == 0, "", fmt.Sprintf("OFFSET %d", b.query.offset)).(string)
	// 组合
	sqlstr := strings.Join(Filter([]string{
		"SELECT", distinct, fields, "FROM", table, join, where, group, having, order, limit, offset, b.query.lock},
		func(x string) bool { return x != "" }), " ")
	// log.Println("sql output ->", sqlstr)
	return sqlstr, args, nil
}

func (b *Builder) parseJoin() (string, []string, error) {
	var result []string
	var joinTables []string
	for _, join := range b.query.join {
		var w string
		var ok bool
		var args []interface{}
		var sp string
		sp = join[0].(string)
		if args, ok = join[1].([]interface{}); !ok {
			return "", nil, errors.New("join conditions are wrong")
		}
		switch len(args) {
		case 1:
			w = args[0].(string)
			joinTables = append(joinTables, args[0].(string))
		case 2:
			w = fmt.Sprintf("%s ON %s", args[0].(string), args[1].(string))
			joinTables = append(joinTables, args[0].(string))
		case 3:
			w = fmt.Sprintf("%s AS %s ON %s", args[0].(string), args[1].(string), args[2].(string))
			joinTables = append(joinTables, args[1].(string))
		default:
			return "", nil, errors.New("join format error")
		}
		result = append(result, sp+" "+w)
	}
	return strings.Join(result, " "), joinTables, nil
}

func (b *Builder) parseWhere() (string, []interface{}, error) {
	var result []string
	var args []interface{}
	for _, where := range b.query.where {
		var ok bool
		var wargs []interface{}
		sp := where[0].(string)
		condition := where[1].(string)
		if wargs, ok = where[2].([]interface{}); !ok {
			return "", nil, errors.New("where conditions are wrong")
		}
		if strings.Index(condition, ":") > 0 {
			_query, _args, err := sqlx.Named(condition, wargs[0])
			if err != nil {
				return "", nil, err
			}
			condition = _query
			wargs = _args
		}
		if strings.Index(strings.ToUpper(condition), " IN ") > 0 {
			_query, _args, err := sqlx.In(condition, wargs...)
			if err != nil {
				return "", nil, err
			}
			condition = _query
			wargs = _args
		}
		result = append(result, sp+" "+condition)
		args = append(args, wargs...)
	}
	wherestring := strings.Trim(strings.TrimLeft(strings.TrimLeft(strings.Join(result, " "), "AND"), "OR"), " ")
	return wherestring, args, nil
}

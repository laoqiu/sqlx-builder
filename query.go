package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antlabs/deepcopy"
	"github.com/jmoiron/sqlx"
)

// Query 返回sql语句主体
type Query struct {
	Table        string
	Fields       []string
	IgnoreFields []string
	Join         [][]interface{}
	Where        [][]interface{}
	Distinct     bool
	Order        string
	Group        string
	Having       string
	Limit        int
	Offset       int
	Lock         string
	Comment      string
}

// Table 生成基本query并赋值Query
func (b *Builder) Table(tablename string) *Builder {
	b.query = &Query{
		Table: tablename,
	}
	return b
}

// Copy 复制
func (b *Builder) Copy() *Builder {
	cp := &Builder{
		db:    b.db,
		tx:    b.tx,
		debug: b.debug,
		query: &Query{},
	}
	deepcopy.Copy(cp.query, b.query).Do()
	return cp
}

// Comment 加入sql注释
func (b *Builder) Comment(v string) *Builder {
	b.query.Comment = fmt.Sprintf("/* %s */", v)
	return b
}

// LockInShareMode 共享锁
func (b *Builder) LockInShareMode() *Builder {
	b.query.Lock = "LOCK IN SHARE MODE"
	return b
}

// LockForUpdate 写锁
func (b *Builder) LockForUpdate() *Builder {
	b.query.Lock = "FOR UPDATE"
	return b
}

// Fields 定义返回字段
func (b *Builder) Fields(fields ...string) *Builder {
	b.query.Fields = fields
	return b
}

// IgnoreFields 忽略关键字段(如id)
func (b *Builder) IgnoreFields(fields ...string) *Builder {
	b.query.IgnoreFields = fields
	return b
}

// AddFields 添加新的返回字段
func (b *Builder) AddFields(fields ...string) *Builder {
	b.query.Fields = append(b.query.Fields, fields...)
	return b
}

// Join 赋值多表关联join表达式
func (b *Builder) Join(table ...interface{}) *Builder {
	b.query.Join = append(b.query.Join, []interface{}{"INNER JOIN", table})
	return b
}

// LeftJoin 多表关联leftjoin表达式
func (b *Builder) LeftJoin(table ...interface{}) *Builder {
	b.query.Join = append(b.query.Join, []interface{}{"LEFT JOIN", table})
	return b
}

// RightJoin 多表关联rightjoin表达式
func (b *Builder) RightJoin(table ...interface{}) *Builder {
	b.query.Join = append(b.query.Join, []interface{}{"RIGHT JOIN", table})
	return b
}

// UnionJoin 联合查询union join表达式
func (b *Builder) UnionJoin(table ...interface{}) *Builder {
	b.query.Join = append(b.query.Join, []interface{}{"UNION JOIN", table})
	return b
}

// Where 条件查询
func (b *Builder) Where(query string, args ...interface{}) *Builder {
	b.query.Where = append(b.query.Where, []interface{}{"AND", query, args})
	return b
}

// GroupBy 分组查询
func (b *Builder) GroupBy(group string) *Builder {
	b.query.Group = group
	return b
}

// OrderBy 排序
func (b *Builder) OrderBy(order string) *Builder {
	b.query.Order = order
	return b
}

// Distinct 去重
func (b *Builder) Distinct() *Builder {
	b.query.Distinct = true
	return b
}

// Limit set limit number
func (b *Builder) Limit(n int) *Builder {
	b.query.Limit = n
	return b
}

// Offset set offset number
func (b *Builder) Offset(n int) *Builder {
	b.query.Offset = n
	return b
}

func (b *Builder) _parseInsert(data map[string]interface{}) (string, string) {
	var keystr, valstr string
	var key, val []string
	for k, v := range data {
		if len(b.query.Fields) == 0 || indexOf(k, b.query.Fields) != -1 {
			key = append(key, "`"+k+"`")
			// 反射找出类型
			switch v.(type) {
			case int, int8, int32, int64, float32, float64, bool:
				val = append(val, fmt.Sprintf("%v", v))
			default:
				val = append(val, fmt.Sprintf("'%s'", v))
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
		if len(b.query.Fields) == 0 || indexOf(k, b.query.Fields) != -1 {
			var _value string
			// 反射找出类型
			switch v.(type) {
			case int, int8, int32, int64, float32, float64, bool:
				_value = fmt.Sprintf("`%s` = %v", k, v)
			default:
				_value = fmt.Sprintf("`%s` = '%s'", k, v)
			}
			result = append(result, _value)
		}
	}

	setstr = strings.Join(result, ", ")
	return setstr
}

func (b *Builder) _parseJoin() (string, []string, error) {
	var result []string
	var joinTables []string
	for _, join := range b.query.Join {
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

func (b *Builder) _parseWhere() (string, []interface{}, error) {
	var result []string
	var args []interface{}
	for _, where := range b.query.Where {
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

// BuildExec 返回需要执行的sql表达式
func (b *Builder) BuildExec(method string, data map[string]interface{}) (string, []interface{}, error) {
	var sqlstr string
	var tablename string

	tablename = b.query.Table

	where, args, err := b._parseWhere()
	if err != nil {
		return "", nil, err
	}
	where = If(where == "", "", "WHERE "+where).(string)

	switch method {
	case "INSERT", "INSERT_IGNORE":
		keystr, valstr := b._parseInsert(data)
		ignore := If(method == "INSERT_IGNORE", "IGNORE", "").(string)
		sqlstr = fmt.Sprintf("INSERT %s %s INTO %s (%s) VALUES (%s)", b.query.Comment, ignore, tablename, keystr, valstr)
	case "INSERT_ON_DUPLICATE_UPDATE":
		keystr, valstr := b._parseInsert(data)
		setstr := b._parseUpate(data)
		sqlstr = fmt.Sprintf("INSERT %s INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s", b.query.Comment, tablename, keystr, valstr, setstr)
	case "UPDATE":
		setstr := b._parseUpate(data)
		sqlstr = fmt.Sprintf("UPDATE %s %s SET %s %s", b.query.Comment, tablename, setstr, where)
	case "DELETE":
		sqlstr = fmt.Sprintf("DELETE %s FROM %s %s", b.query.Comment, tablename, where)
	}
	//log.Println("sql output ->", sqlstr)
	return sqlstr, args, nil
}

// BuildQuery 合并query表达式
func (b *Builder) BuildQuery() (string, []interface{}, error) {
	// table
	table := DefaultMapper(b.query.Table)
	// join
	join, joinTables, err := b._parseJoin()
	if err != nil {
		return "", nil, err
	}
	// distinct
	distinct := If(b.query.Distinct, "DISTINCT", "").(string)
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
	fields := If(len(b.query.Fields) == 0, allFields, strings.Join(b.query.Fields, ",")).(string)
	// where
	where, args, err := b._parseWhere()
	if err != nil {
		return "", nil, err
	}
	// where
	where = If(where == "", "", "WHERE "+where).(string)
	// group
	group := If(b.query.Group == "", "", "GROUP BY "+b.query.Group).(string)
	// order
	order := If(b.query.Order == "", "", "ORDER BY "+b.query.Order).(string)
	// having
	having := If(b.query.Having == "", "", " HAVING "+b.query.Having).(string)
	// limit
	limit := If(b.query.Limit == 0, "", fmt.Sprintf("LIMIT %d", b.query.Limit)).(string)
	// offset
	offset := If(b.query.Offset == 0, "", fmt.Sprintf("OFFSET %d", b.query.Offset)).(string)
	// 组合
	sqlstr := strings.Join(Filter([]string{
		"SELECT", b.query.Comment, distinct, fields, "FROM", table, join, where, group, having, order, limit, offset, b.query.Lock},
		func(x string) bool { return x != "" }), " ")
	return sqlstr, args, nil
}

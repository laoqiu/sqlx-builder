package sqlxt

import (
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

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
}

func Table(tablename string) *Query {
	return &Query{
		table: tablename,
	}
}

func (q *Query) Fields(fields ...string) *Query {
	q.fields = fields
	return q
}

func (q *Query) AddFields(fields ...string) *Query {
	q.fields = append(q.fields, fields...)
	return q
}

func (q *Query) Select(fields ...string) *Query {
	q.Fields(fields...)
	return q
}

func (q *Query) AddSelect(fields ...string) *Query {
	q.AddFields(fields...)
	return q
}

func (q *Query) Join(table ...interface{}) *Query {
	q.join = append(q.join, []interface{}{"INNER JOIN", table})
	return q
}

func (q *Query) LeftJoin(table ...interface{}) *Query {
	q.join = append(q.join, []interface{}{"LEFT JOIN", table})
	return q
}

func (q *Query) RightJoin(table ...interface{}) *Query {
	q.join = append(q.join, []interface{}{"RIGHT JOIN", table})
	return q
}

func (q *Query) UnionJoin(table ...interface{}) *Query {
	q.join = append(q.join, []interface{}{"UNION JOIN", table})
	return q
}

func (q *Query) Where(query string, args ...interface{}) *Query {
	q.where = append(q.where, []interface{}{"AND", query, args})
	return q
}

func (q *Query) GroupBy(group string) *Query {
	q.group = group
	return q
}

func (q *Query) OrderBy(order string) *Query {
	q.order = order
	return q
}

func (q *Query) Distinct() *Query {
	q.distinct = true
	return q
}

func (q *Query) Limit(n int) *Query {
	q.limit = n
	return q
}

func (q *Query) Offset(n int) *Query {
	q.offset = n
	return q
}

func (q *Query) _parseInsert(data map[string]interface{}) (string, string) {
	var keystr, valstr string
	var key, val []string
	for k, v := range data {
		if len(q.fields) == 0 || indexOf(k, q.fields) != -1 {
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

func (q *Query) _parseUpate(data map[string]interface{}) string {
	var setstr string
	var result []string

	for k, v := range data {
		if len(q.fields) == 0 || indexOf(k, q.fields) != -1 {
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

func (q *Query) BuildExec(method string, data map[string]interface{}) (string, []interface{}, error) {
	var sqlstr string
	var tablename string

	tablename = q.table

	where, args, err := q.parseWhere()
	if err != nil {
		return "", nil, err
	}
	where = If(where == "", "", "WHERE "+where).(string)

	switch method {
	case "INSERT", "INSERT_IGNORE":
		keystr, valstr := q._parseInsert(data)
		ignore := If(method == "INSERT_IGNORE", "IGNORE", "").(string)
		sqlstr = fmt.Sprintf("INSERT %s INTO %s (%s) VALUES (%s)", ignore, tablename, keystr, valstr)
	case "INSERT_ON_DUPLICATE_UPDATE":
		keystr, valstr := q._parseInsert(data)
		setstr := q._parseUpate(data)
		sqlstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s", tablename, keystr, valstr, setstr)
	case "UPDATE":
		setstr := q._parseUpate(data)
		sqlstr = fmt.Sprintf("UPDATE %s SET %s %s", tablename, setstr, where)
	case "DELETE":
		sqlstr = fmt.Sprintf("DELETE FROM %s %s", tablename, where)
	}
	//log.Println("sql output ->", sqlstr)
	return sqlstr, args, nil
}

func (q *Query) BuildQuery() (string, []interface{}, error) {
	// table
	table := DefaultMapper(q.table)
	// join
	join, joinTables, err := q.parseJoin()
	if err != nil {
		return "", nil, err
	}
	// distinct
	distinct := If(q.distinct, "DISTINCT", "").(string)
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
	fields := If(len(q.fields) == 0, allFields, strings.Join(q.fields, ",")).(string)
	// where
	where, args, err := q.parseWhere()
	if err != nil {
		return "", nil, err
	}
	// where
	where = If(where == "", "", "WHERE "+where).(string)
	// group
	group := If(q.group == "", "", "GROUP BY "+q.group).(string)
	// order
	order := If(q.order == "", "", "ORDER BY "+q.order).(string)
	// having
	having := If(q.having == "", "", " HAVING "+q.having).(string)
	// limit
	limit := If(q.limit == 0, "", fmt.Sprintf("LIMIT %d", q.limit)).(string)
	// offset
	offset := If(q.offset == 0, "", fmt.Sprintf("OFFSET %d", q.offset)).(string)
	// 组合
	sqlstr := strings.Join(Filter([]string{
		"SELECT", distinct, fields, "FROM", table, join, where, group, having, order, limit, offset},
		func(x string) bool { return x != "" }), " ")
	log.Println("sql output ->", sqlstr)
	return sqlstr, args, nil
}

func (q *Query) parseJoin() (string, []string, error) {
	var result []string
	var joinTables []string
	for _, join := range q.join {
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

func (q *Query) parseWhere() (string, []interface{}, error) {
	var result []string
	var args []interface{}
	for _, where := range q.where {
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

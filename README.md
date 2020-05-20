# sqlx-builder
用于实现sqlx的链式操作

## 安装
```
go get github.com/laoqiu/sqlx-builder
```

### 初始化

```
import (
	sqlxb "github.com/laoqiu/sqlx-builder"
)
db, err := sqlxb.Connect()
if err != nil {
	return err
}
// 加载mapper
builder.LoadMapper(db, sqlxb.DefaultMapper)
```

### 支持的链式操作
```
dest := &Person{}
p := map[string]interface{}{
    "phone": "13012345678",
}
debug := true
query = sqlxb.Table("tablename").Join("table2", "table2.id = table1.t_id").
    Where("name like ?", "%name%").
    Where("address = ?", "test").
    Where("phone = :phone", p)
err := sqlxb.New(db, query, debug).First(dest)
if err != nil {
    log.Println(err) 
}
```

### 支持的函数及结构体
* Builder
```
func New(*sqlx.DB, *Query, debug bool) *Builder
func NewTx(*sqlx.Tx, *Query, debug bool) *Builder
func (st *Builder) Insert(data interface{}) error
func (st *Builder) Update(data interface{}) error
func (st *Builder) Delete() error
func (st *Builder) Get(dest interface{}) (sql.Result, error)
func (st *Builder) All(dest interface{}) (sql.Result, error)
```
* Query
```
func Table(tablename) *Query
func (q *Query) Distinct() *Query
func (q *Query) Select(...fields) *Query
func (q *Query) Join(table, condition, label) *Query
func (q *Query) LeftJoin(table, condition, label) *Query
func (q *Query) RightJoin(table, condition, label) *Query
func (q *Query) UnionJoin(table, condition, label) *Query
func (q *Query) Where(condition, ...args) *Query
func (q *Query) Limit(n) *Query
func (q *Query) Offset(n) *Query
```
* func Connect
* func LoadMapper

# sqlxt
用于sqlx简化操作，仅完成基础的一些操作，之后想支持通过struct转table结构
目前使用protobuf生成struct，感觉转table用处不大

不使用sqlx的可以去看 gorose
更多struct功能的可以去看 structable

## 安装
```
go get github.com/laoqiu/sqlxt
```

### 初始化Init函数，允许传入option (在获取flag时更新db的参数)

```
file: db.go
var (
	db *sqlx.DB
)

func Init(opts ...sqlxt.Option) error {
	o := sqlxt.NewOptions(opts...)
	db, err := sqlxt.Connect(o.Driver, o.URI, o.Charset, o.ParseTime, o.MaxClient, o.MaxClient)
	if err != nil {
		return err
	}
	// 加载mapper
	sqlxt.LoadMapper(db, sqlxt.DefaultMapper)
	return nil
}

file: main.go
func main() {
	dbOpts := []sqlxt.Option{}
	service := micro.NewService(
        ...
	dbOpts = append(dbOpts, sqlxt.URI(c.String("database_url")))
	...
	)
	// db init
	if err := db.Init(dbOpts...); err != nil {
		log.Fatal(err)
	}
}
```

### 支持的链式操作
```
dest := &Person{}
p := map[string]interface{}{
    "phone": "13012345678",
}
query = sqlxt.NewQuery().Table("tablename").Join("table2", "table2.id = table1.t_id").
    Where("name like ?", "%name%").
    Where("address = ?", "test").
    Where("phone = :phone", p)
err := sqlxt.New(db, query).First(dest)
if err != nil {
    log.Println(err) 
}
```

### 支持的函数及结构体
* Sqlxt
```
func New(*sqlx.DB, *Query) *Sqlxt
func NewTx(*sqlx.Tx, *Query) *Sqlxt
func (st *Sqlxt) Insert(data) error
func (st *Sqlxt) Update(data) error
func (st *Sqlxt) Delete() error
func (st *Sqlxt) Get(dest) (sql.Result, error)
func (st *Sqlxt) All(dest) (sql.Result, error)
```
* Query
```
func NewQuery() *Query
func (q *Query) Table(tablename) *Query
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

# sqlx-query-colt
用于sqlx简化操作，仅完成基础的一些操作，之后想支持通过struct转table结构
目前使用protobuf生成struct，感觉转table用处不大

不使用sqlx的可以去看 gorose
更多struct功能的可以去看 structable

### 支持的链式操作
```
dest := &Person{}
p := map[string]interface{}{
    "phone": "13012345678",
}
query := xcolt.Query{}
query.Bind(db).Table("tablename").Join("table2", "table2.id = table1.t_id").
    Where("name like ?", "%name%").
    Where("address = ?", "test").
    Where("phone = :phone", p).
    First(dest)
```

### 支持的函数结构体
```
Bind(sqlx.DB)
BindTx(sqlx.Tx)
Table(tablename)
Distinct()
Select(...fields)
Join(table, condition, label)
LeftJoin(table, condition, label)
RightJoin(table, condition, label)
UnionJoin(table, condition, label)
Where(condition, ...args)
Limit(n)
Offset(n)
Insert(data)
Update(data)
Delete()
First(dest)
All(dest)
```
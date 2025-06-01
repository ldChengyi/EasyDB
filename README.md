## EasyDB 内存数据库引擎

EasyDB 是一个轻量级的内存数据存储引擎，支持多种索引结构和查询方式，适用于数据抓取与实时查询场景。
好吧就是一个小轮子，我在写抓包的时候不知道用啥匹配当前使用情况就自己写了个。
不知道有没有人用TAT
然后下面都是GPT写的www

### 特性

- 支持泛型，可存储任意类型数据
- 内置三种索引类型：
  - 精确匹配（Exact Match）
  - 前缀匹配（Prefix Match）
  - 子串匹配（Substring Match）
- 线程安全
- 支持基本的 CRUD 操作
- 高性能的内存存储(划掉)

### 快速开始

1. 安装

```bash
go get github.com/ldChengYi/EasyDB
```

2. 基本使用

```go
package main

import (
    "context"
    "github.com/ldChengYi/EasyDB/core/storage"
)

// 定义数据结构
type User struct {
    Name string
    Age  int
}

func main() {
    // 创建存储实例
    opts := storage.Options{
        InitialCapacity: 1000,
        EnableVersioning: true,
        FieldIndexes: []storage.FieldIndexConfig[User]{
            {
                Field: "name",
                Extractor: func(r *types.Record[User]) string {
                    return r.Data.Name
                },
                Types: []storage.IndexType{
                    storage.IndexExact,    // 精确匹配
                    storage.IndexPrefix,    // 前缀匹配
                    storage.IndexSubstring, // 子串匹配
                },
            },
        },
    }
    
    store := storage.New[User](opts)
    ctx := context.Background()
    
    // 插入数据
    user := User{Name: "张三", Age: 25}
    record, _ := store.Insert(ctx, user)
    
    // 查询数据
    found, _ := store.Get(ctx, record.ID)
    
    // 更新数据
    user.Age = 26
    updated, _ := store.Update(ctx, record.ID, user)
    
    // 删除数据
    _ = store.Delete(ctx, record.ID)
    
    // 列表查询
    records, total, _ := store.List(ctx, 0, 10)
}
```

3. 使用索引查询

```go
// 使用索引管理器进行查询
// 精确查询
results := store.IndexManager.Query("name", "张三")

// 前缀查询
results := store.IndexManager.QueryPrefix("name", "张")

// 子串查询
results := store.IndexManager.QuerySubstring("name", "三")
```

### 配置选项

- `InitialCapacity`: 初始存储容量
- `EnableVersioning`: 是否启用版本控制
- `FieldIndexes`: 字段索引配置
  - `Field`: 索引字段名
  - `Extractor`: 字段值提取函数
  - `Types`: 支持的索引类型

### 性能建议

1. 合理设置初始容量，避免频繁扩容
2. 只为需要查询的字段创建索引
3. 根据查询模式选择合适的索引类型：
   - 精确匹配：适用于等值查询
   - 前缀匹配：适用于自动完成、搜索提示
   - 子串匹配：适用于模糊搜索，但消耗较多内存

### 注意事项

1. 所有数据存储在内存中，重启后数据会丢失
2. 子串索引会占用较多内存，请谨慎使用
3. 建议在单机场景下使用
4. 适合数据量中等的实时查询场景

### 许可证

MIT License


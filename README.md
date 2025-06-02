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
- 支持链式查询，查询封装在API包里
- 使用见下
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
    "github.com/ldChengYi/EasyDB/api"
    "github.com/ldChengYi/EasyDB/storage"
)

// 定义数据结构
type TestData{
    Name string,
    Age int,
    ID int,
    Tags []string,
    CreatedAt, time.Time
}
func setupTestStore(t *testing.T) *storage.Store[TestData] {
	store, err := NewStoreBuilder[TestData]().
		SetCapacity(1000).
		SetVersioning(true).
		AddIndex("Name", func(r *types.Record[TestData]) interface{} {
			return r.Data.Name
		}, storage.IndexExact, storage.IndexPrefix, storage.IndexSubstring).
		AddIndex("Age", func(r *types.Record[TestData]) interface{} {
			return r.Data.Age // 使用固定长度格式化，便于范围查询
		}, storage.IndexExact, storage.IndexPrefix). // 添加前缀索引支持范围查询
		AddIndex("ID", func(r *types.Record[TestData]) interface{} {
			return r.Data.ID
		}, storage.IndexExact).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, store)

	// 插入测试数据
	testData := []TestData{
		{
			ID:        "1",
			Name:      "张三",
			Age:       25,
			Score:     85.5,
			Tags:      []string{"学生", "男"},
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "2",
			Name:      "李四",
			Age:       30,
			Score:     92.0,
			Tags:      []string{"教师", "男"},
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:        "3",
			Name:      "王五",
			Age:       28,
			Score:     88.5,
			Tags:      []string{"学生", "女"},
			CreatedAt: time.Now().Add(-72 * time.Hour),
		},
		{
			ID:        "4",
			Name:      "赵六",
			Age:       35,
			Score:     95.0,
			Tags:      []string{"教师", "女"},
			CreatedAt: time.Now(),
		},
	}

	ctx := context.Background()
	for _, data := range testData {
		_, err := store.Insert(ctx, data)
		assert.NoError(t, err)
	}

	return store
}

func TestQuery_Equals(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// 测试精确匹配
	results, err := NewQuery(store).
		Where("Age").Equals("35").
		Do(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, 35, results[0].Data.Age)
}

func TestQuery_Contains(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// 测试包含匹配
	results, err := NewQuery(store).
		Where("Name").Contains("张").
		Do(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Contains(t, results[0].Data.Name, "张")
}

func TestQuery_In(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// 测试 IN 查询
	results, err := NewQuery(store).
		Where("Age").In(25, 30).
		Do(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Contains(t, []int{25, 30}, r.Data.Age)
	}
}

func TestQuery_Between(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	// 测试范围查询
	results, err := NewQuery(store).
		Where("Age").Between(25, 30).
		Do(ctx)

	t.Log(results)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	for _, r := range results {
		assert.GreaterOrEqual(t, r.Data.Age, 25)
		assert.LessOrEqual(t, r.Data.Age, 30)
	}
}
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


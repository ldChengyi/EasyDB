package api

import (
	"fmt"

	"github.com/ldChengYi/EasyDB/core/storage"
	"github.com/ldChengYi/EasyDB/core/types"
)

// StoreBuilder 是一个用于构建存储实例的构建器。
// 它提供了流式API来配置存储选项和索引。
// 泛型参数 T 可以是任意结构体类型。
type StoreBuilder[T any] struct {
	initialCapacity  int
	enableVersioning bool
	indexBuilder     *IndexBuilder[T]
	built            bool
}

// NewStoreBuilder 创建一个新的存储构建器实例。
// 返回:
//   - *StoreBuilder[T]: 新的存储构建器实例，默认初始容量为1000，启用版本控制
func NewStoreBuilder[T any]() *StoreBuilder[T] {
	return &StoreBuilder[T]{
		initialCapacity:  1000,
		enableVersioning: true,
		indexBuilder:     NewIndexBuilder[T](),
	}
}

// SetCapacity 设置存储的初始容量。
// 参数:
//   - capacity: 初始容量大小
//
// 返回:
//   - *StoreBuilder[T]: 构建器实例，用于链式调用
func (b *StoreBuilder[T]) SetCapacity(capacity int) *StoreBuilder[T] {
	b.initialCapacity = capacity
	return b
}

// SetVersioning 设置是否启用版本控制。
// 参数:
//   - enable: 是否启用版本控制
//
// 返回:
//   - *StoreBuilder[T]: 构建器实例，用于链式调用
func (b *StoreBuilder[T]) SetVersioning(enable bool) *StoreBuilder[T] {
	b.enableVersioning = enable
	return b
}

// AddIndex 添加字段索引配置。
// 参数:
//   - field: 要索引的字段名
//   - extractor: 字段值提取函数
//   - types: 索引类型列表
//
// 返回:
//   - *StoreBuilder[T]: 构建器实例，用于链式调用
func (b *StoreBuilder[T]) AddIndex(field string, extractor func(*types.Record[T]) interface{}, types ...storage.IndexType) *StoreBuilder[T] {
	b.indexBuilder.AddField(field, extractor, types...)
	return b
}

// Build 构建并返回存储实例。
// 返回:
//   - *storage.Store[T]: 构建的存储实例
//   - error: 构建过程中的错误
func (b *StoreBuilder[T]) Build() (*storage.Store[T], error) {
	if b.built {
		return nil, fmt.Errorf("store builder already used")
	}

	if b.initialCapacity <= 0 {
		return nil, fmt.Errorf("initial capacity must be positive")
	}

	opts := storage.Options{
		InitialCapacity:  b.initialCapacity,
		EnableVersioning: b.enableVersioning,
		FieldIndexes:     b.indexBuilder.Build(),
	}

	b.built = true
	return storage.New[T](opts), nil
}

// IndexBuilder 是一个用于构建字段索引配置的构建器。
// 泛型参数 T 可以是任意结构体类型。
type IndexBuilder[T any] struct {
	configs []storage.FieldIndexConfig[T]
}

// NewIndexBuilder 创建一个新的索引构建器实例。
// 返回:
//   - *IndexBuilder[T]: 新的索引构建器实例
func NewIndexBuilder[T any]() *IndexBuilder[T] {
	return &IndexBuilder[T]{
		configs: make([]storage.FieldIndexConfig[T], 0),
	}
}

// AddField 添加字段索引配置。
// 参数:
//   - field: 要索引的字段名
//   - extractor: 字段值提取函数
//   - types: 索引类型列表
//
// 返回:
//   - *IndexBuilder[T]: 构建器实例，用于链式调用
func (b *IndexBuilder[T]) AddField(field string, extractor func(*types.Record[T]) interface{}, types ...storage.IndexType) *IndexBuilder[T] {
	b.configs = append(b.configs, storage.FieldIndexConfig[T]{
		Field:     field,
		Extractor: extractor,
		Types:     types,
	})
	return b
}

// Build 构建并返回索引配置列表。
// 返回:
//   - []storage.FieldIndexConfig[T]: 索引配置列表
func (b *IndexBuilder[T]) Build() []storage.FieldIndexConfig[T] {
	return b.configs
}

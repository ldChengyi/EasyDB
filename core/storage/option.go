package storage

import "github.com/ldChengYi/EasyDB/core/types"

type FieldIndexConfig[T any] struct {
	Field     string                             // 字段名称
	Extractor func(*types.Record[T]) interface{} // 如何从记录中提取字段
	Types     []IndexType                        // 支持的索引类型（精确、前缀、子串）
}

// Options 存储引擎配置选项
type Options struct {
	// InitialCapacity 初始容量
	InitialCapacity int

	// EnableVersioning 是否启用版本控制
	EnableVersioning bool

	// 泛型不支持，需要 Store 初始化时断言
	FieldIndexes any
}

package storage

import (
	"github.com/ldChengYi/EasyDB/core/ds"
	"github.com/ldChengYi/EasyDB/core/types"
)

// IndexType 表示索引类型
type IndexType string

const (
	IndexExact     IndexType = "exact"     // 精确匹配
	IndexPrefix    IndexType = "prefix"    // 前缀匹配
	IndexSubstring IndexType = "substring" // 包含匹配
)

// FieldIndex 表示某字段的索引结构（支持多个类型）
type FieldIndex[T any] struct {
	extractor func(*types.Record[T]) string

	exact    map[string]map[uint64]struct{} // 精确匹配索引
	inverted map[string]map[uint64]struct{} // 子串倒排索引
	trie     *ds.Trie                       // 前缀匹配索引
}

// IndexManager 管理所有字段的索引
type IndexManager[T any] struct {
	indexes map[string]*FieldIndex[T] // fieldName -> 索引结构
}

// NewIndexManager 构造一个空的索引管理器
func NewIndexManager[T any]() *IndexManager[T] {
	return &IndexManager[T]{
		indexes: make(map[string]*FieldIndex[T]),
	}
}

// Register 注册字段的提取器和索引类型
func (im *IndexManager[T]) Register(field string, extractor func(*types.Record[T]) string, types ...IndexType) {
	fi := &FieldIndex[T]{extractor: extractor}

	for _, t := range types {
		switch t {
		case IndexExact:
			fi.exact = make(map[string]map[uint64]struct{})
		case IndexPrefix:
			fi.trie = ds.NewTrie()
		case IndexSubstring:
			fi.inverted = make(map[string]map[uint64]struct{})
		}
	}

	im.indexes[field] = fi
}

// AddIndexByRecord 将记录添加到所有索引中
func (im *IndexManager[T]) AddIndexByRecord(record *types.Record[T]) {
	id := record.ID
	for _, fi := range im.indexes {
		val := fi.extractor(record)

		// 精确索引
		if fi.exact != nil {
			if _, ok := fi.exact[val]; !ok {
				fi.exact[val] = make(map[uint64]struct{})
			}
			fi.exact[val][id] = struct{}{}
		}

		// 前缀索引
		if fi.trie != nil {
			fi.trie.Insert(val, id)
		}

		// 子串索引
		if fi.inverted != nil {
			for i := 0; i <= len(val)-1; i++ {
				for j := i + 1; j <= len(val); j++ {
					sub := val[i:j]
					if _, ok := fi.inverted[sub]; !ok {
						fi.inverted[sub] = make(map[uint64]struct{})
					}
					fi.inverted[sub][id] = struct{}{}
				}
			}
		}
	}
}

// RemoveIndexByRecord 将记录从所有索引中移除
func (im *IndexManager[T]) RemoveIndexByRecord(record *types.Record[T]) {
	id := record.ID
	for _, fi := range im.indexes {
		val := fi.extractor(record)

		// 精确索引
		if fi.exact != nil {
			if idSet, ok := fi.exact[val]; ok {
				delete(idSet, id)
				if len(idSet) == 0 {
					delete(fi.exact, val)
				}
			}
		}

		// 前缀索引
		if fi.trie != nil {
			fi.trie.Delete(val, id)
		}

		// 子串索引
		if fi.inverted != nil {
			for i := 0; i <= len(val)-1; i++ {
				for j := i + 1; j <= len(val); j++ {
					sub := val[i:j]
					if idSet, ok := fi.inverted[sub]; ok {
						delete(idSet, id)
						if len(idSet) == 0 {
							delete(fi.inverted, sub)
						}
					}
				}
			}
		}
	}
}

// UpdateIndexByRecord 用新数据更新旧数据索引
func (im *IndexManager[T]) UpdateIndexByRecord(oldRecord, newRecord *types.Record[T]) {
	im.RemoveIndexByRecord(oldRecord)
	im.AddIndexByRecord(newRecord)
}

// Query 查询索引数据（优先精确 > 前缀 > 子串）
func (im *IndexManager[T]) Query(field string, keyword string) map[uint64]struct{} {
	if fi, ok := im.indexes[field]; ok {
		// 精确匹配
		if fi.exact != nil {
			if set, ok := fi.exact[keyword]; ok {
				return set
			}
		}
		// 前缀匹配
		if fi.trie != nil {
			if set := fi.trie.QueryPrefix(keyword); set != nil {
				return set
			}
		}
		// 子串匹配
		if fi.inverted != nil {
			if set, ok := fi.inverted[keyword]; ok {
				return set
			}
		}
	}
	return nil
}

// QueryPrefix 仅使用前缀索引进行查询
func (im *IndexManager[T]) QueryPrefix(field string, prefix string) map[uint64]struct{} {
	if fi, ok := im.indexes[field]; ok {
		if fi.trie != nil {
			return fi.trie.QueryPrefix(prefix)
		}
	}
	return nil
}

// QuerySubstring 仅使用子串倒排索引进行查询
func (im *IndexManager[T]) QuerySubstring(field string, substr string) map[uint64]struct{} {
	if fi, ok := im.indexes[field]; ok {
		if fi.inverted != nil {
			if set, ok := fi.inverted[substr]; ok {
				return set
			}
		}
	}
	return nil
}

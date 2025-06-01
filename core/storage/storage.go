package storage

import (
	"context"

	"github.com/ldChengYi/EasyDB/core/types"
)

// Storage 定义存储引擎接口
type Storage[T any] interface {
	// Insert 插入一条记录
	Insert(ctx context.Context, data T) (*types.Record[T], error)

	// Get 获取一条记录
	Get(ctx context.Context, id uint64) (*types.Record[T], error)

	// Update 更新一条记录
	Update(ctx context.Context, id uint64, data T) (*types.Record[T], error)

	// Delete 删除一条记录
	Delete(ctx context.Context, id uint64) error

	// List 列出记录
	List(ctx context.Context, offset, limit int) ([]*types.Record[T], int, error)
}

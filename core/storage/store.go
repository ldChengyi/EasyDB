package storage

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ldChengYi/EasyDB/core/errors"
	"github.com/ldChengYi/EasyDB/core/types"
)

// Store 内存存储引擎实现
type Store[T any] struct {
	sync.RWMutex
	data          []*types.Record[T]
	idGen         atomic.Uint64
	idMapIndex    map[uint64]int
	aliveIndexes  []int
	aliveIndexSet map[int]struct{}

	IndexManager *IndexManager[T]
	options      Options
}

// New 创建新的内存存储实例
func New[T any](opts Options) *Store[T] {
	if opts.InitialCapacity <= 0 {
		opts.InitialCapacity = 1000
	}

	store := &Store[T]{
		data:          make([]*types.Record[T], 0, opts.InitialCapacity),
		idMapIndex:    make(map[uint64]int),
		aliveIndexes:  make([]int, 0),
		aliveIndexSet: make(map[int]struct{}),
		IndexManager:  NewIndexManager[T](), // 初始化新的索引管理器
		options:       opts,
	}

	if list, ok := opts.FieldIndexes.([]FieldIndexConfig[T]); ok {
		for _, cfg := range list {
			store.IndexManager.Register(cfg.Field, cfg.Extractor, cfg.Types...)
		}
	}

	return store
}

func (s *Store[T]) Insert(ctx context.Context, data T) (*types.Record[T], error) {
	s.Lock()
	defer s.Unlock()

	id := s.idGen.Add(1)
	now := time.Now().UnixNano()
	record := &types.Record[T]{
		ID:      id,
		Data:    data,
		Version: 1,
		Meta: types.RecordMeta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	index := len(s.data)
	s.data = append(s.data, record)
	s.idMapIndex[id] = index
	s.addAliveIndex(index)

	s.IndexManager.AddIndexByRecord(record)

	return record, nil
}

func (s *Store[T]) Get(ctx context.Context, id uint64) (*types.Record[T], error) {
	s.RLock()
	defer s.RUnlock()

	idx, ok := s.idMapIndex[id]
	if !ok || s.data[idx].Meta.Deleted {
		return nil, errors.ErrNotFound
	}

	return s.data[idx], nil
}

func (s *Store[T]) Update(ctx context.Context, id uint64, data T) (*types.Record[T], error) {
	s.Lock()
	defer s.Unlock()

	idx, ok := s.idMapIndex[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	record := s.data[idx]
	if record.Meta.Deleted {
		return nil, errors.ErrRecordDeleted
	}

	old := *record
	record.Data = data
	record.Meta.UpdatedAt = time.Now().UnixNano()
	if s.options.EnableVersioning {
		record.Version++
	}

	s.IndexManager.UpdateIndexByRecord(&old, record)

	return record, nil
}

func (s *Store[T]) Delete(ctx context.Context, id uint64) error {
	s.Lock()
	defer s.Unlock()

	idx, ok := s.idMapIndex[id]
	if !ok {
		return errors.ErrNotFound
	}

	record := s.data[idx]
	if record.Meta.Deleted {
		return errors.ErrRecordDeleted
	}

	record.Meta.Deleted = true
	record.Meta.UpdatedAt = time.Now().UnixNano()
	s.removeAliveIndex(idx)

	s.IndexManager.RemoveIndexByRecord(record)
	return nil
}

func (s *Store[T]) List(ctx context.Context, offset, limit int) ([]*types.Record[T], int, error) {
	s.RLock()
	defer s.RUnlock()

	total := len(s.aliveIndexes)
	if offset >= total {
		return []*types.Record[T]{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	records := make([]*types.Record[T], 0, end-offset)

	for _, idx := range s.aliveIndexes[offset:end] {
		rec := s.data[idx]
		// 可以保留防御式检查，防止脏数据
		if rec.Meta.Deleted {
			continue
		}
		records = append(records, rec)
	}

	return records, total, nil
}

// 添加活跃项（Insert 时调用）
func (s *Store[T]) addAliveIndex(index int) {
	s.aliveIndexes = append(s.aliveIndexes, index)
	s.aliveIndexSet[index] = struct{}{}
}

// 移除活跃项（Delete 时调用）
func (s *Store[T]) removeAliveIndex(index int) {
	delete(s.aliveIndexSet, index)
	for i, idx := range s.aliveIndexes {
		if idx == index {
			s.aliveIndexes = append(s.aliveIndexes[:i], s.aliveIndexes[i+1:]...)
			break
		}
	}
}

func (s *Store[T]) Data() []*types.Record[T] {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

func (s *Store[T]) AliveIndexes() []int {
	s.RLock()
	defer s.RUnlock()
	return s.aliveIndexes
}

func (s *Store[T]) Size() int {
	return len(s.data)
}

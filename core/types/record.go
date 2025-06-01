package types

// RecordMeta 记录的元数据
type RecordMeta struct {
	CreatedAt int64
	UpdatedAt int64
	Deleted   bool
}

// Record 表示数据库中的一条记录
type Record[T any] struct {
	ID      uint64
	Data    T
	Version uint64
	Meta    RecordMeta
}

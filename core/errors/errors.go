package errors

import "errors"

var (
	// ErrNotFound 记录未找到错误
	ErrNotFound = errors.New("record not found")

	// ErrDataDeleted 数据已经被删除
	ErrRecordDeleted = errors.New("record has been deleted")

	// ErrInvalidInput 无效输入错误
	ErrInvalidInput = errors.New("invalid input")

	// ErrVersionConflict 版本冲突错误
	ErrVersionConflict = errors.New("version conflict: record has been modified by another operation")

	// ErrRecordNotFound 记录未找到
	ErrRecordNotFound = errors.New("记录未找到")

	// ErrFieldNotFound 字段未找到
	ErrFieldNotFound = errors.New("字段未找到")

	// ErrDuplicateKey 重复键值
	ErrDuplicateKey = errors.New("重复的键值")

	// ErrIndexAlreadyExists 索引已存在
	ErrIndexAlreadyExists = errors.New("索引已存在")

	// ErrIndexNotFound 索引未找到
	ErrIndexNotFound = errors.New("索引未找到")

	// ErrNoSnapshot 快照不存在
	ErrNoSnapshot = errors.New("快照不存在")
)

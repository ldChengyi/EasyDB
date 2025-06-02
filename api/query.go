package api

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ldChengYi/EasyDB/core/storage"
	"github.com/ldChengYi/EasyDB/core/types"
	"github.com/ldChengYi/EasyDB/util"
)

// Query 是一个泛型查询构建器，支持链式调用的查询语法。
// 它提供了丰富的查询条件、排序和分页功能。
// 泛型参数 T 可以是任意结构体类型。
type Query[T any] struct {
	store      *storage.Store[T]
	conditions []queryCondition
	limit      int
	offset     int
	orderBy    string
	orderDesc  bool
	timeRange  struct {
		start, end int64
	}
}

// NewQuery 创建一个新的查询构建器实例。
// 参数:
//   - store: 要查询的数据存储实例
//
// 返回:
//   - 新的查询构建器实例，默认限制为100条记录
func NewQuery[T any](store *storage.Store[T]) *Query[T] {
	return &Query[T]{
		store:      store,
		conditions: make([]queryCondition, 0),
		limit:      100,
	}
}

// Where 开始构建字段查询条件。
// 参数:
//   - field: 要查询的字段名
//
// 返回:
//   - 字段查询构建器，用于指定具体的查询条件
func (q *Query[T]) Where(field string) *FieldQuery[T] {
	return &FieldQuery[T]{
		query: q,
		field: field,
	}
}

// FieldQuery 是字段查询构建器，用于构建特定字段的查询条件。
type FieldQuery[T any] struct {
	query *Query[T]
	field string
}

// Equals 添加精确匹配条件。
// 参数:
//   - value: 要匹配的值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) Equals(value interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opEquals,
		value:    value,
	})
	return fq.query
}

// Contains 添加包含匹配条件。
// 参数:
//   - value: 要包含的字符串
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) Contains(value string) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opContains,
		value:    value,
	})
	return fq.query
}

// In 添加集合匹配条件。
// 参数:
//   - values: 要匹配的值列表
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) In(values ...interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opIn,
		value:    values,
	})
	return fq.query
}

// Between 添加范围匹配条件。
// 参数:
//   - start: 范围起始值
//   - end: 范围结束值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) Between(start, end interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opBetween,
		value:    []interface{}{start, end},
	})
	return fq.query
}

// GreaterThan 添加大于条件。
// 参数:
//   - value: 比较值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) GreaterThan(value interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opGt,
		value:    value,
	})
	return fq.query
}

// GreaterThanOrEqual 添加大于等于条件。
// 参数:
//   - value: 比较值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) GreaterThanOrEqual(value interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opGte,
		value:    value,
	})
	return fq.query
}

// LessThan 添加小于条件。
// 参数:
//   - value: 比较值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) LessThan(value interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opLt,
		value:    value,
	})
	return fq.query
}

// LessThanOrEqual 添加小于等于条件。
// 参数:
//   - value: 比较值
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (fq *FieldQuery[T]) LessThanOrEqual(value interface{}) *Query[T] {
	fq.query.conditions = append(fq.query.conditions, queryCondition{
		field:    fq.field,
		operator: opLte,
		value:    value,
	})
	return fq.query
}

// Limit 设置结果数量限制。
// 参数:
//   - limit: 最大返回记录数
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (q *Query[T]) Limit(limit int) *Query[T] {
	q.limit = limit
	return q
}

// Offset 设置结果偏移量。
// 参数:
//   - offset: 跳过的记录数
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (q *Query[T]) Offset(offset int) *Query[T] {
	q.offset = offset
	return q
}

// OrderBy 设置排序规则。
// 参数:
//   - field: 排序字段
//   - desc: 是否降序排序
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (q *Query[T]) OrderBy(field string, desc bool) *Query[T] {
	q.orderBy = field
	q.orderDesc = desc
	return q
}

// InTimeRange 设置时间范围过滤。
// 参数:
//   - start: 起始时间
//   - end: 结束时间
//
// 返回:
//   - 查询构建器实例，用于链式调用
func (q *Query[T]) InTimeRange(start, end time.Time) *Query[T] {
	q.timeRange.start = start.UnixNano()
	q.timeRange.end = end.UnixNano()
	return q
}

// Do 执行查询并返回结果。
// 参数:
//   - ctx: 上下文，用于控制查询超时和取消
//
// 返回:
//   - []*types.Record[T]: 查询结果记录列表
//   - error: 查询过程中的错误
func (q *Query[T]) Do(ctx context.Context) ([]*types.Record[T], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	done := make(chan struct {
		results []*types.Record[T]
		err     error
	})

	go func() {
		results, err := q.executeQuery(queryCtx)
		done <- struct {
			results []*types.Record[T]
			err     error
		}{results, err}
	}()

	select {
	case <-queryCtx.Done():
		return nil, queryCtx.Err()
	case result := <-done:
		return result.results, result.err
	}
}

// executeQuery 执行实际的查询操作。
// 参数:
//   - ctx: 上下文，用于控制查询超时和取消
//
// 返回:
//   - []*types.Record[T]: 查询结果记录列表
//   - error: 查询过程中的错误
func (q *Query[T]) executeQuery(ctx context.Context) ([]*types.Record[T], error) {
	var results []*types.Record[T]
	var matchedIDs map[uint64]struct{}

	for i, cond := range q.conditions {
		currentMatches, err := q.processCondition(ctx, cond)
		if err != nil {
			return nil, fmt.Errorf("failed to process condition: %w", err)
		}

		if i == 0 {
			matchedIDs = currentMatches
		} else {
			for id := range matchedIDs {
				if _, ok := currentMatches[id]; !ok {
					delete(matchedIDs, id)
				}
			}
		}
	}

	for id := range matchedIDs {
		if record, err := q.store.Get(ctx, id); err == nil {
			results = append(results, record)
		}
	}

	if q.orderBy != "" {
		// if err := q.sortResults(results); err != nil {
		// 	return nil, fmt.Errorf("failed to sort results: %w", err)
		// }
	}

	return q.applyPagination(results)
}

// processCondition 处理单个查询条件。
// 参数:
//   - ctx: 上下文
//   - cond: 查询条件
//
// 返回:
//   - map[uint64]struct{}: 匹配的记录ID集合
//   - error: 处理过程中的错误
func (q *Query[T]) processCondition(ctx context.Context, cond queryCondition) (map[uint64]struct{}, error) {
	switch cond.operator {
	case opEquals:
		return q.processEqualCondition(cond)
	case opContains:
		return q.processContainCondition(cond)
	case opIn:
		return q.processInCondition(cond)
	case opBetween, opGt, opGte, opLt, opLte:
		return q.processRangeCondition(ctx, cond)
	default:
		return nil, fmt.Errorf("unsupported operator: %s", cond.operator)
	}
}

func (q *Query[T]) processEqualCondition(cond queryCondition) (map[uint64]struct{}, error) {
	fts := q.store.IndexManager.GetFieldTypes()
	ft, ok := fts[cond.field]
	if !ok {
		return nil, fmt.Errorf("field %s not indexed", cond.field)
	}

	convertedVal, err := convertValueToType(cond.value, ft)
	if err != nil {
		return nil, fmt.Errorf("type conversion failed: %v", err)
	}

	matches := q.store.IndexManager.Query(cond.field, convertedVal)
	if matches == nil {
		return make(map[uint64]struct{}), nil
	}
	return matches, nil
}

// processContainCondition 处理 Contain 条件。
// 参数:
//   - cond: Contain 查询条件
//
// 返回:
//   - map[uint64]struct{}: 匹配的记录ID集合
//   - error: 处理过程中的错误
func (q *Query[T]) processContainCondition(cond queryCondition) (map[uint64]struct{}, error) {
	im := q.store.IndexManager
	field := cond.field

	// 转为 string，用于 prefix/substring 匹配
	valStr, err := util.SafeToString(cond.value)
	if err != nil {
		return nil, fmt.Errorf("field %s: value not string-convertible: %w", field, err)
	}

	// 优先使用前缀索引
	if result := im.QueryPrefix(field, valStr); result != nil {
		return result, nil
	}

	// 再使用子串倒排索引
	if result := im.QuerySubstring(field, valStr); result != nil {
		return result, nil
	}

	// 如果该字段没注册相关索引，返回错误
	return nil, fmt.Errorf("field %s does not support prefix or substring index", field)
}

// processInCondition 处理 IN 条件。
// 参数:
//   - cond: IN 查询条件
//
// 返回:
//   - map[uint64]struct{}: 匹配的记录ID集合
//   - error: 处理过程中的错误
func (q *Query[T]) processInCondition(cond queryCondition) (map[uint64]struct{}, error) {
	// 反射解出 cond.value 的 slice 元素
	val := reflect.ValueOf(cond.value)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("in operator requires a slice value, got %T", cond.value)
	}

	im := q.store.IndexManager
	field := cond.field

	// 检查是否存在对应字段的索引
	if _, ok := im.GetIndexes()[field]; !ok {
		return nil, fmt.Errorf("no index found for field %s", field)
	}

	result := make(map[uint64]struct{})
	foundAny := false

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		set := im.Query(field, item)
		if set != nil {
			foundAny = true
			for id := range set {
				result[id] = struct{}{}
			}
		}
	}

	if !foundAny {
		return nil, fmt.Errorf("no matching entries found in 'IN' condition for field %s", field)
	}

	return result, nil
}

// processRangeCondition 处理范围条件。
// 参数:
//   - ctx: 上下文
//   - cond: 范围查询条件
//
// 返回:
//   - map[uint64]struct{}: 匹配的记录ID集合
//   - error: 处理过程中的错误
func (q *Query[T]) processRangeCondition(ctx context.Context, cond queryCondition) (map[uint64]struct{}, error) {
	result := make(map[uint64]struct{})
	fieldExtractor, ok := q.store.IndexManager.GetExtractor(cond.field)
	if !ok {
		return nil, fmt.Errorf("field extractor not found for field: %s", cond.field)
	}

	all := q.store.Data()
	for _, r := range all {
		val := fieldExtractor(r)

		switch cond.operator {
		case opBetween:
			// between 要求是 [min, max] 两个元素
			bounds, ok := cond.value.([]interface{})
			if !ok || len(bounds) != 2 {
				return nil, fmt.Errorf("between requires [min, max] slice")
			}
			if util.Compare(val, bounds[0]) >= 0 && util.Compare(val, bounds[1]) <= 0 {
				result[r.ID] = struct{}{}
			}
		case opGt:
			if util.Compare(val, cond.value) > 0 {
				result[r.ID] = struct{}{}
			}
		case opGte:
			if util.Compare(val, cond.value) >= 0 {
				result[r.ID] = struct{}{}
			}
		case opLt:
			if util.Compare(val, cond.value) < 0 {
				result[r.ID] = struct{}{}
			}
		case opLte:
			if util.Compare(val, cond.value) <= 0 {
				result[r.ID] = struct{}{}
			}
		default:
			return nil, fmt.Errorf("unsupported operator: %s", cond.operator)
		}
	}

	return result, nil
}

// applyPagination 应用分页。
// 参数:
//   - results: 要分页的记录列表
//
// 返回:
//   - []*types.Record[T]: 分页后的记录列表
//   - error: 分页过程中的错误
func (q *Query[T]) applyPagination(results []*types.Record[T]) ([]*types.Record[T], error) {
	if q.offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}
	if q.limit <= 0 {
		return nil, fmt.Errorf("limit must be positive")
	}
	if q.offset >= len(results) {
		return make([]*types.Record[T], 0), nil
	}

	start := q.offset
	end := start + q.limit
	if end > len(results) {
		end = len(results)
	}

	// 创建新的切片以避免共享底层数组
	result := make([]*types.Record[T], end-start)
	copy(result, results[start:end])
	return result, nil
}

// toString 尝试将任意类型转换为 string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// isSameType 判断两个值是否类型一致
func isSameType(a, b interface{}) bool {
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

func convertValueToType(val interface{}, targetType reflect.Type) (interface{}, error) {
	v := reflect.ValueOf(val)

	if !v.Type().ConvertibleTo(targetType) {
		return nil, fmt.Errorf("cannot convert %v to %v", v.Type(), targetType)
	}

	converted := v.Convert(targetType)
	return converted.Interface(), nil
}

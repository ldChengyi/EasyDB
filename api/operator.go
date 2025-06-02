package api

// operator 定义查询操作符（包内私有）
type operator string

const (
	opEquals   operator = "eq"       // 精确匹配
	opContains operator = "contains" // 包含匹配
	opIn       operator = "in"       // 集合匹配
	opBetween  operator = "between"  // 范围匹配
	opGt       operator = "gt"       // 大于
	opGte      operator = "gte"      // 大于等于
	opLt       operator = "lt"       // 小于
	opLte      operator = "lte"      // 小于等于
)

// queryCondition 表示查询条件（包内私有）
type queryCondition struct {
	field    string      // 字段名
	operator operator    // 操作符
	value    interface{} // 比较值
}

// 增加类型检查辅助函数
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

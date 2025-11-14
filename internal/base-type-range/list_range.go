package basetyperange

import (
	"reflect"
)

func ListValueRangerChecker(fieldChecker baseChecker) *ListRange {
	return &ListRange{
		fieldChecker: fieldChecker,
	}
}

type ListRange struct {
	fieldChecker baseChecker
}

func (lr *ListRange) Check(value any) bool {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Array && valueType.Kind() != reflect.Slice { // 必须是数组的类型
		return false
	}

	// 遍历数组的每个元素进行检测
	valueValue := reflect.ValueOf(value)
	length := valueValue.Len()
	for i := 0; i < length; i++ {
		elemValue := valueValue.Index(i).Interface()
		if !lr.fieldChecker.Check(elemValue) {
			return false
		}
	}

	return true
}

func (lr *ListRange) ToString() string {
	return ""
}

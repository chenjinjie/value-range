package basetyperange

import (
	"fmt"
	"reflect"
)

func MapValueRangerChecker(keyChecker baseChecker, fieldChecker baseChecker) *MapRange {
	return &MapRange{
		keyChecker:   keyChecker,
		fieldChecker: fieldChecker,
	}
}

type MapRange struct {
	keyChecker   baseChecker
	fieldChecker baseChecker
}

func (mr *MapRange) Check(value any) bool {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Map { // 必须是 map 类型
		fmt.Printf("value no struct, is: %s\n", valueType.Kind().String())
		return false
	}

	// 获得这个 map 的key 的类型
	keyType := valueType.Key()
	switch keyType.Kind() { // 只支持部分 key 类型，因为考虑到这个类型范围检测，主要是为了从 配置表、协议等场景中的数据范围检测使用，这些场景下的 map 的 key 的类型一般不会太复杂
	case reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 支持的 key 类型
	default:
		fmt.Printf("map key type no support, is: %s\n", keyType.Kind().String())
		return false
	}

	valueValue := reflect.ValueOf(value)
	for _, key := range valueValue.MapKeys() {
		mapValue := valueValue.MapIndex(key).Interface()
		if !mr.keyChecker.Check(key.Interface()) {
			return false
		}
		if !mr.fieldChecker.Check(mapValue) {
			return false
		}
	}

	return true
}

func (mt *MapRange) ToString() string {
	return ""
}

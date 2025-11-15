package basetyperange

import (
	"fmt"
	"reflect"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func MapValueRangerChecker(options options.Options, keyChecker baseChecker, fieldChecker baseChecker) *MapRange {
	return &MapRange{
		options:      options,
		keyChecker:   keyChecker,
		fieldChecker: fieldChecker,
	}
}

type MapRange struct {
	options      options.Options
	keyChecker   baseChecker
	fieldChecker baseChecker
}

func (mr *MapRange) Check(value any) (bool, *checklog.CheckLog) {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Map { // 必须是 map 类型
		return false, checklog.CheckLogFail(fmt.Sprintf("value no struct, is: %s\n", valueType.Kind().String()))
	}

	// 获得这个 map 的key 的类型
	keyType := valueType.Key()
	switch keyType.Kind() { // 只支持部分 key 类型，因为考虑到这个类型范围检测，主要是为了从 配置表、协议等场景中的数据范围检测使用，这些场景下的 map 的 key 的类型一般不会太复杂
	case reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 支持的 key 类型
	default:
		return false, checklog.CheckLogFail(fmt.Sprintf("map key type no support, is: %s\n", keyType.Kind().String()))
	}

	haveError := false
	checkNodeMap := checklog.CheckLogNodeMap()

	valueValue := reflect.ValueOf(value)
	for _, key := range valueValue.MapKeys() { // map 一次放入两个 log

		mapValue := valueValue.MapIndex(key).Interface()
		okKey, childCheckLogKey := mr.keyChecker.Check(key.Interface())
		okValue, childCheckLogValue := mr.fieldChecker.Check(mapValue)

		if !okKey || !okValue {
			haveError = true

			childCheckLogKey.SetFieldKey("")
			checkNodeMap.AddLog(childCheckLogKey)

			childCheckLogValue.SetFieldKey("")
			checkNodeMap.AddLog(childCheckLogValue)

		} else if !mr.options.OnlySaveCheckFailLog {
			childCheckLogKey.SetFieldKey("")
			checkNodeMap.AddLog(childCheckLogKey)

			childCheckLogValue.SetFieldKey("")
			checkNodeMap.AddLog(childCheckLogValue)
		}
	}

	return !haveError, checklog.CheckLogWithNode(checkNodeMap)
}

func (mt *MapRange) ToString() string {
	return ""
}

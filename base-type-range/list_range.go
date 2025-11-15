package basetyperange

import (
	"fmt"
	"reflect"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func ListValueRangerChecker(options options.Options, fieldChecker baseChecker) *ListRange {
	return &ListRange{
		options:      options,
		fieldChecker: fieldChecker,
	}
}

type ListRange struct {
	options      options.Options
	fieldChecker baseChecker
}

func (lr *ListRange) Check(value any) (bool, *checklog.CheckLog) {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Array && valueType.Kind() != reflect.Slice { // 必须是数组的类型
		return false, checklog.CheckLogFail("value no array or slice")
	}

	haveError := false
	checkNode := checklog.CheckLogNodeCommon()

	// 遍历数组的每个元素进行检测
	valueValue := reflect.ValueOf(value)
	length := valueValue.Len()
	for i := 0; i < length; i++ {
		elemValue := valueValue.Index(i).Interface()
		if ok, childCheckLog := lr.fieldChecker.Check(elemValue); !ok {
			haveError = true
			childCheckLog.SetFieldKey(fmt.Sprintf("[%d]", i))
			checkNode.AddLog(childCheckLog)
		} else if !lr.options.OnlySaveCheckFailLog {
			childCheckLog.SetFieldKey(fmt.Sprintf("[%d]", i))
			checkNode.AddLog(childCheckLog)
		}
	}

	return !haveError, checklog.CheckLogWithNode(checkNode)
}

func (lr *ListRange) ToString() string {
	return ""
}

package basetyperange

import (
	"fmt"
	"reflect"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func StructValueRangerChecker(options options.Options, checker any) *StructRange {
	checkerType := reflect.TypeOf(checker)
	if checkerType.Kind() != reflect.Struct { // checker 必须是 struct 类型
		panic("struct range checker must be struct type")
	}
	count := checkerType.NumField()
	if count == 0 {
		panic("struct range checker has no field")
	}
	checkerValue := reflect.ValueOf(checker)

	mapChecker := make(map[string]baseChecker)
	for i := 0; i < count; i++ {
		fieldName := checkerType.Field(i).Name
		fieldCheckerValue := checkerValue.Field(i).Interface()
		fieldChecker, ok := fieldCheckerValue.(baseChecker)
		if !ok {
			panic("struct range checker field: " + fieldName + " is not a baseChecker")
		}
		mapChecker[fieldName] = fieldChecker
	}

	return &StructRange{
		options:    options,
		mapChecker: mapChecker,
	}
}

// 结构体范围检测
type baseChecker interface {
	Check(value any) (bool, *checklog.CheckLog)
}

type StructRange struct {
	options    options.Options
	mapChecker map[string]baseChecker
}

func (sr *StructRange) prt2OriThenCheck(value any) (bool, *checklog.CheckLog) {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Ptr { // 必须是指针类型
		return false, checklog.CheckLogFail(fmt.Sprintf("value no ptr, is: %s\n", valueType.Kind().String()))
	}
	valueValue := reflect.ValueOf(value)
	oriValue := valueValue.Elem().Interface()

	oriValueType := reflect.TypeOf(oriValue)
	if oriValueType.Kind() != reflect.Struct { // 不要搞指针套指针，只能一层指针指向 struct
		return false, checklog.CheckLogFail(fmt.Sprintf("value no struct, is: %s\n", oriValueType.Kind().String()))
	}

	return sr.Check(oriValue)
}

func (sr *StructRange) Check(value any) (bool, *checklog.CheckLog) {
	valueType := reflect.TypeOf(value)

	if valueType.Kind() == reflect.Ptr { // 如果是指针类型，获得其真正的 struct 再去检测
		return sr.prt2OriThenCheck(value)
	}

	if valueType.Kind() != reflect.Struct { // 必须是结构体类型
		return false, checklog.CheckLogFail(fmt.Sprintf("value no struct, is: %s", valueType.Kind().String()))
	}

	valueNumFieldsLen := valueType.NumField()
	checkerNumFieldsLen := len(sr.mapChecker)
	if checkerNumFieldsLen != valueNumFieldsLen { // cheker 对不上，不用检测就知道不通过了
		return false, checklog.CheckLogFail(fmt.Sprintf("value struct field num not match checker, value fields: %d, checker fields: %d", valueNumFieldsLen, checkerNumFieldsLen))
	}

	haveError := false
	checkNode := checklog.CheckLogNodeCommon()
	// 逐字段检测
	valueValue := reflect.ValueOf(value)
	for fieldName, fieldChecker := range sr.mapChecker {
		_, ok := valueType.FieldByName(fieldName)
		if !ok { // value 中没有对应的字段，检测不通过
			leaf := checklog.CheckLogFail("not exit")
			leaf.SetFieldKey(fieldName)
			checkNode.AddLog(leaf)
			haveError = true
			continue
		}
		valueFieldValue := valueValue.FieldByName(fieldName).Interface()

		if ok, childCheckLog := fieldChecker.Check(valueFieldValue); !ok { // 去检测该字段的值是否符合范围
			haveError = true
			childCheckLog.SetFieldKey(fieldName)
			checkNode.AddLog(childCheckLog)
		} else if !sr.options.OnlySaveCheckFailLog {
			childCheckLog.SetFieldKey(fieldName)
			checkNode.AddLog(childCheckLog)
		}
	}

	return !haveError, checklog.CheckLogWithNode(checkNode)
}

func (sr *StructRange) ToString() string {
	return ""
}

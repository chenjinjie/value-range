package expandtyperange

import (
	"fmt"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func EnumValueStore() *EnumStore {
	return &EnumStore{
		oriEnumData: make(map[string]map[uint64]struct{}),
	}
}

// 枚举的情况比较简单，全部转为 uint64 来做存储和检测就行了
type EnumStore struct {
	oriEnumData map[string]map[uint64]struct{}
}

func (es *EnumStore) LoadOneEnum(enumKey string, enumData map[uint64]struct{}) bool {
	if _, ok := es.oriEnumData[enumKey]; ok {
		fmt.Printf("enum key duplicate load: %s\n", enumKey)
		return false
	}

	es.oriEnumData[enumKey] = enumData
	return true
}

func (es *EnumStore) CheckEnumValue(enumKey string, value uint64) (bool, *checklog.CheckLog) {
	enumData, ok := es.oriEnumData[enumKey]
	if !ok {
		return false, checklog.CheckLogFail(fmt.Sprintf("enum key not exit: %s\n", enumKey))
	}

	_, ok = enumData[value]
	if !ok {
		return false, checklog.CheckLogFail(fmt.Sprintf("%d no in enum %s", value, enumKey))
	}

	return true, checklog.CheckLogSuccess(fmt.Sprintf("%d", value))
}

func (es *EnumStore) EnumRuleExit(enumKey string) bool {
	_, ok := es.oriEnumData[enumKey]
	return ok
}

func EnumValueRangerChecker(options options.Options, enumStore *EnumStore, enumKey string) *EnumRange {
	if enumStore == nil {
		panic("EnumValueRangerChecker enumStore is nil")
	}
	if !enumStore.EnumRuleExit(enumKey) {
		panic("EnumValueRangerChecker enumKey not exit: " + enumKey)
	}

	return &EnumRange{
		enumKey:   enumKey,
		enumStore: enumStore,
	}
}

type EnumRange struct {
	options   options.Options
	enumKey   string
	enumStore *EnumStore
}

func (er *EnumRange) Check(value any) (bool, *checklog.CheckLog) {
	switch v := value.(type) {
	case uint64:
		return er.enumStore.CheckEnumValue(er.enumKey, v)
	case uint32:
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case uint16:
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case uint8:
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case uint:
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int64:
		if v < 0 {
			return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value negative int64: %d", v))
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int32:
		if v < 0 {
			return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value negative int32: %d", v))
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int16:
		if v < 0 {
			return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value negative int16: %d", v))
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int8:
		if v < 0 {
			return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value negative int8: %d", v))
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int:
		if v < 0 {
			return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value negative int: %d", v))
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	default:
		return false, checklog.CheckLogFail(fmt.Sprintf("EnumRange check value type no support, type: %T", v))
	}
}

func (er *EnumRange) ToString() string {
	return ""
}

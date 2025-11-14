package expandtyperange

import (
	"fmt"
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

func (es *EnumStore) CheckEnumValue(enumKey string, value uint64) bool {
	enumData, ok := es.oriEnumData[enumKey]
	if !ok {
		fmt.Printf("enum key not exit: %s\n", enumKey)
		return false
	}

	_, ok = enumData[value]
	if !ok {
		fmt.Printf("enum key: %s value: %d not exit\n", enumKey, value)
		return false
	}

	return true
}

func (es *EnumStore) EnumRuleExit(enumKey string) bool {
	_, ok := es.oriEnumData[enumKey]
	return ok
}

func EnumValueRangerChecker(enumStore *EnumStore, enumKey string) *EnumRange {
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
	enumKey   string
	enumStore *EnumStore
}

func (er *EnumRange) Check(value any) bool {
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
			fmt.Printf("EnumRange check value negative int64: %d\n", v)
			return false
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int32:
		if v < 0 {
			fmt.Printf("EnumRange check value negative int32: %d\n", v)
			return false
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int16:
		if v < 0 {
			fmt.Printf("EnumRange check value negative int16: %d\n", v)
			return false
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int8:
		if v < 0 {
			fmt.Printf("EnumRange check value negative int8: %d\n", v)
			return false
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	case int:
		if v < 0 {
			fmt.Printf("EnumRange check value negative int: %d\n", v)
			return false
		}
		return er.enumStore.CheckEnumValue(er.enumKey, uint64(v))
	default:
		fmt.Printf("EnumRange check value type no support, type: %T\n", v)
		return false
	}
}

func (er *EnumRange) ToString() string {
	return ""
}

package valuerange

import (
	"fmt"

	basetyperange "github.com/chenjinjie/value-range/internal/base-type-range"
	expandtyperange "github.com/chenjinjie/value-range/internal/expand-type-range"
)

type ValueRangerChecker interface {
	Check(value any) bool
	ToString() string
}

func ValueRangeChecker() *ValueRange {
	return &ValueRange{
		refStore:  expandtyperange.RefValueStore(),
		enumStore: expandtyperange.EnumValueStore(),

		checkerStore: make(map[string]ValueRangerChecker),
	}
}

// 对一个系统进行值范围检测的对象
// 因为拓展类型中，ref、enum 等类型的值范围检测，需要把其他数据引入进来并且缓存起来，做一些映射关系等操作
// 想要共用这些缓存，所以用一个对象再包起来
type ValueRange struct {
	refStore  *expandtyperange.RefStore
	enumStore *expandtyperange.EnumStore

	checkerStore map[string]ValueRangerChecker
}

// 提前加载配置表
func (vr *ValueRange) LoadOneCfg(cfgKey string, cfgData any) bool {
	if vr.refStore == nil {
		return false
	}
	return vr.refStore.LoadOneOriData(cfgKey, cfgData)
}

// 提前加载枚举配置
func (vr *ValueRange) LoadOneEnumCfg(enumKey string, enumData map[uint64]struct{}) bool {
	if vr.enumStore == nil {
		return false
	}
	return vr.enumStore.LoadOneEnum(enumKey, enumData)
}

func (vr *ValueRange) RegChecker(key string, checker ValueRangerChecker) {
	if _, ok := vr.checkerStore[key]; ok {
		panic("reg checker duplicate key: " + key)
	}
	vr.checkerStore[key] = checker
}

func (vr *ValueRange) Check(key string, value any) bool {
	checker, ok := vr.checkerStore[key]
	if !ok {
		fmt.Printf("check rule not exit, key: %s", key)
		return false
	}
	return checker.Check(value)
}

func (vr *ValueRange) IntValueRangerChecker(pattern string) ValueRangerChecker {
	return basetyperange.IntValueRangerChecker(pattern)
}

func (vr *ValueRange) StringValueRangerChecker(rangeStr string) ValueRangerChecker {
	return basetyperange.StringValueRangerChecker(rangeStr)
}

func (vr *ValueRange) BoolValueRangerChecker(rangeStr string) ValueRangerChecker {
	return basetyperange.BoolValueRangerChecker(rangeStr)
}

func (vr *ValueRange) StructValueRangerChecker(checker any) ValueRangerChecker {
	return basetyperange.StructValueRangerChecker(checker)
}

func (vr *ValueRange) ListValueRangerChecker(fieldChecker ValueRangerChecker) ValueRangerChecker {
	return basetyperange.ListValueRangerChecker(fieldChecker)
}

func (vr *ValueRange) MapValueRangerChecker(keyChecker ValueRangerChecker, fieldChecker ValueRangerChecker) ValueRangerChecker {
	return basetyperange.MapValueRangerChecker(keyChecker, fieldChecker)
}

func (vr *ValueRange) RefValueRangerChecker(rangeStr string) ValueRangerChecker {
	return expandtyperange.RefValueRangerChecker(vr.refStore, rangeStr)
}

func (vr *ValueRange) EnumValueRangerChecker(enumKey string) ValueRangerChecker {
	return expandtyperange.EnumValueRangerChecker(vr.enumStore, enumKey)
}

package valuerange

import (
	"fmt"

	basetyperange "github.com/chenjinjie/value-range/base-type-range"
	checklog "github.com/chenjinjie/value-range/check-log"
	expandtyperange "github.com/chenjinjie/value-range/expand-type-range"
	"github.com/chenjinjie/value-range/options"
)

type ValueRangerChecker interface {
	Check(value any) (bool, *checklog.CheckLog)
	ToString() string
}

func ValueRangeChecker(options options.Options) *ValueRange {
	return &ValueRange{
		options:   options,
		refStore:  expandtyperange.RefValueStore(),
		enumStore: expandtyperange.EnumValueStore(),

		checkerStore: make(map[string]ValueRangerChecker),
	}
}

// 对一个系统进行值范围检测的对象
// 因为拓展类型中，ref、enum 等类型的值范围检测，需要把其他数据引入进来并且缓存起来，做一些映射关系等操作
// 想要共用这些缓存，所以用一个对象再包起来
type ValueRange struct {
	options options.Options

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

func (vr *ValueRange) Check(key string, value any) (bool, *checklog.CheckLog) {
	checker, ok := vr.checkerStore[key]
	if !ok {
		return false, checklog.CheckLogFail(fmt.Sprintf("check rule not exit, key: %s", key))
	}
	return checker.Check(value)
}

// panic if error
func (vr *ValueRange) IntValueRangerChecker(patternStr string) ValueRangerChecker {
	return basetyperange.IntValueRangerChecker(vr.options, patternStr)
}

// panic if error
func (vr *ValueRange) StringValueRangerChecker(patternStr string) ValueRangerChecker {
	return basetyperange.StringValueRangerChecker(vr.options, patternStr)
}

// panic if error
func (vr *ValueRange) BoolValueRangerChecker(patternStr string) ValueRangerChecker {
	return basetyperange.BoolValueRangerChecker(vr.options, patternStr)
}

// panic if error
func (vr *ValueRange) StructValueRangerChecker(checker any) ValueRangerChecker {
	return basetyperange.StructValueRangerChecker(vr.options, checker)
}

// panic if error
func (vr *ValueRange) ListValueRangerChecker(fieldChecker ValueRangerChecker) ValueRangerChecker {
	return basetyperange.ListValueRangerChecker(vr.options, fieldChecker)
}

// panic if error
func (vr *ValueRange) MapValueRangerChecker(keyChecker ValueRangerChecker, fieldChecker ValueRangerChecker) ValueRangerChecker {
	return basetyperange.MapValueRangerChecker(vr.options, keyChecker, fieldChecker)
}

// panic if error
func (vr *ValueRange) RefValueRangerChecker(refPatternStr string) ValueRangerChecker {
	return expandtyperange.RefValueRangerChecker(vr.options, vr.refStore, refPatternStr)
}

// panic if error
func (vr *ValueRange) EnumValueRangerChecker(enumKey string) ValueRangerChecker {
	return expandtyperange.EnumValueRangerChecker(vr.options, vr.enumStore, enumKey)
}

// panic if error
func (vr *ValueRange) CheckerFromReg(regKey string) ValueRangerChecker {
	checker, ok := vr.checkerStore[regKey]
	if !ok {
		panic(fmt.Sprintf("check rule not exit, key: %s", regKey))
	}
	return checker
}

func (vr *ValueRange) OnlySaveCheckFailLog() bool {
	return vr.options.OnlySaveCheckFailLog
}

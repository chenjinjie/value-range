package valuerange

import (
	"testing"
)

type heroTagCfgChecker struct {
	Free ValueRangerChecker
}

// hero 配置表的 品质 枚举值
const enumHeroCfgQualityKey = "heroCfgQuality"
const (
	heroCfgQuality_1 uint8 = 1 // 品质1
	heroCfgQuality_2 uint8 = 2
	heroCfgQuality_3 uint8 = 3
	heroCfgQuality_4 uint8 = 4
	heroCfgQuality_5 uint8 = 5
)

const enumHeroCfgAttr = "heroCfgAttr"
const (
	heroCfgAttr_hp uint32 = 1
	heroCfgAttr_mp uint32 = 2
)

// 做一个对应的配置检测
type heroCfgChecker struct {
	Id      ValueRangerChecker
	Desc    ValueRangerChecker
	Quality ValueRangerChecker
	Open    ValueRangerChecker
	Tag     ValueRangerChecker
	Skins   ValueRangerChecker
	Attrs   ValueRangerChecker
}

type heroSkinCfgChecker struct {
	Id   ValueRangerChecker
	Desc ValueRangerChecker
}

func TestCfgCheck(t *testing.T) {
	/// 创建一个检测对象，并且预加载数据
	valueRangeChecker := ValueRangeChecker()

	{ /// 预加载 => 不保证多线程安全
		// 预加载配置数据
		if !valueRangeChecker.LoadOneCfg(heroCfgKey, heroCfgList) {
			t.Errorf("load heroCfg data failed")
			return
		}
		if !valueRangeChecker.LoadOneCfg(heroSkinCfgKey, heroSkinCfgList) {
			t.Errorf("load heroSkinCfg data failed")
			return
		}

		// 预先加载枚举
		var enumQuality = map[uint64]struct{}{
			uint64(heroCfgQuality_1): {},
			uint64(heroCfgQuality_2): {},
			uint64(heroCfgQuality_3): {},
			uint64(heroCfgQuality_4): {},
			uint64(heroCfgQuality_5): {},
		}
		if !valueRangeChecker.LoadOneEnumCfg(enumHeroCfgQualityKey, enumQuality) {
			t.Errorf("load enum heroCfgQuality data failed")
			return
		}

		var enumAttr = map[uint64]struct{}{
			uint64(heroCfgAttr_hp): {},
			uint64(heroCfgAttr_mp): {},
		}
		if !valueRangeChecker.LoadOneEnumCfg(enumHeroCfgAttr, enumAttr) {
			t.Errorf("load enum heroCfgAttr data failed")
			return
		}
	}

	{ /// 注册没个配置表的检测规则 => 不保证多线程安全
		regAllCheckers := func() (result bool) {
			// 为了让注册的代码更好看点，不需要每次都判断 chekcer 创建成功，
			// 类似 cheker, err := xxx(); if err != nil { return false } 这种
			// 将 checker 创建失败，都用 panic 的方式抛出，然后在这里统一捕获
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic occurred during checker registration: %v", r)
					result = false
				}
			}()

			valueRangeChecker.RegChecker(heroCfgKey, valueRangeChecker.ListValueRangerChecker(valueRangeChecker.StructValueRangerChecker(heroCfgChecker{
				Id:      valueRangeChecker.IntValueRangerChecker(""),                     // 配置id 不限制范围
				Desc:    valueRangeChecker.StringValueRangerChecker(""),                  // 配置描述 不限制范围
				Quality: valueRangeChecker.EnumValueRangerChecker(enumHeroCfgQualityKey), // 品质范围 => 枚举值
				Open:    valueRangeChecker.BoolValueRangerChecker(""),                    // 是否开放 不限制范围
				Tag: valueRangeChecker.StructValueRangerChecker(heroTagCfgChecker{
					Free: valueRangeChecker.BoolValueRangerChecker(""),
				}),
				Skins: valueRangeChecker.ListValueRangerChecker(valueRangeChecker.RefValueRangerChecker(heroSkinCfgKey + ".Id")), // 可用皮肤列表，引用检测 heroSkinCfg 表的 Id 字段
				Attrs: valueRangeChecker.MapValueRangerChecker(valueRangeChecker.EnumValueRangerChecker(enumHeroCfgAttr), valueRangeChecker.IntValueRangerChecker("(0,-)")),
			})))

			valueRangeChecker.RegChecker(heroSkinCfgKey, valueRangeChecker.ListValueRangerChecker(valueRangeChecker.StructValueRangerChecker(heroSkinCfgChecker{
				Id:   valueRangeChecker.IntValueRangerChecker(""), // 皮肤配置id 不限制范围
				Desc: valueRangeChecker.StringValueRangerChecker(""),
			})))

			return true
		}

		if !regAllCheckers() {
			t.Errorf("register checker failed")
			return
		}
	}

	{ /// 开始检测每张配置表的配置值 是否都符合要求
		if !valueRangeChecker.Check(heroCfgKey, heroCfgList) {
			t.Errorf("cfg check failed. %s ", heroCfgKey)
			return
		}
		if !valueRangeChecker.Check(heroSkinCfgKey, heroSkinCfgList) {
			t.Errorf("cfg check failed. %s ", heroSkinCfgKey)
			return
		}
	}

	t.Log(" === all hero config check passed ===\n")
}

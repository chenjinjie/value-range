package basetyperange

import (
	"fmt"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func BoolValueRangerChecker(options options.Options, patternStr string) *BoolRange {
	if patternStr == "" {
		return &BoolRange{
			originalStr: "",
			noRange:     true,

			needTrue: false,
		}
	}

	if patternStr == "true" {
		return &BoolRange{
			originalStr: patternStr,
			noRange:     false,

			needTrue: true,
		}
	}

	if patternStr == "false" {
		return &BoolRange{
			options:     options,
			originalStr: patternStr,
			noRange:     false,

			needTrue: false,
		}
	}

	panic("BoolRange pattern not illegal: " + patternStr) // 值范围描述字符串不合法
}

type BoolRange struct {
	options options.Options

	originalStr string // 原始字符串表示
	noRange     bool   // 没有数值范围限制，是 bool 即可

	needTrue bool // 需要为 true
}

func (lr *BoolRange) Check(value any) (bool, *checklog.CheckLog) {
	v, ok := value.(bool)
	if !ok { // 不是 bool 类型
		return false, checklog.CheckLogFail("value is not a bool")
	}
	if lr.noRange {
		return true, checklog.CheckLogSuccess(fmt.Sprintf("%t", v))
	}

	if lr.needTrue && v {
		return true, checklog.CheckLogSuccess(fmt.Sprintf("%t", v))
	} else if !lr.needTrue && !v {
		return true, checklog.CheckLogSuccess(fmt.Sprintf("%t", v))
	} else {
		return false, checklog.CheckLogFail(fmt.Sprintf("value need %t", lr.needTrue))
	}
}

func (lr *BoolRange) ToString() string {
	return ""
}

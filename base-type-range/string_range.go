package basetyperange

import (
	"fmt"

	checklog "github.com/chenjinjie/value-range/check-log"
	"github.com/chenjinjie/value-range/options"
)

func StringValueRangerChecker(options options.Options, patternStr string) *StringRange {
	return &StringRange{
		options: options,
	}
}

type StringRange struct {
	options options.Options
}

func (sr *StringRange) Check(value any) (bool, *checklog.CheckLog) {
	switch value.(type) {
	case string:
		return true, checklog.CheckLogSuccess(fmt.Sprintf("%s", value))
	default:
		return false, checklog.CheckLogFail("value is not a string")
	}
}

func (sr *StringRange) ToString() string {
	return ""
}

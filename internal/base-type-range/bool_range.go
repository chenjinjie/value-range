package basetyperange

func BoolValueRangerChecker(rangeStr string) *BoolRange {
	if rangeStr == "" {
		return &BoolRange{
			originalStr: "",
			noRange:     true,

			needTrue: false,
		}
	}

	if rangeStr == "true" {
		return &BoolRange{
			originalStr: rangeStr,
			noRange:     false,

			needTrue: true,
		}
	}

	if rangeStr == "false" {
		return &BoolRange{
			originalStr: rangeStr,
			noRange:     false,

			needTrue: false,
		}
	}

	panic("BoolRange pattern not illegal: " + rangeStr) // 值范围描述字符串不合法
}

type BoolRange struct {
	originalStr string // 原始字符串表示
	noRange     bool   // 没有数值范围限制，是 bool 即可

	needTrue bool // 需要为 true
}

func (lr *BoolRange) Check(value any) bool {
	v, ok := value.(bool)
	if !ok { // 不是 bool 类型
		return false
	}
	if lr.noRange {
		return true
	}
	if lr.needTrue {
		return v
	} else {
		return !v
	}
}

func (lr *BoolRange) ToString() string {
	return ""
}

package basetyperange

func StringValueRangerChecker(rangeStr string) *StringRange {
	return &StringRange{}
}

type StringRange struct {
}

func (sr *StringRange) Check(value any) bool {
	switch value.(type) {
	case string:
		return true
	default:
		return false
	}
}

func (sr *StringRange) ToString() string {
	return ""
}

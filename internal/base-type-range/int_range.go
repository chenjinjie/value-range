package basetyperange

import (
	"fmt"
	"math"
	"regexp"
)

/*
匹配范围模板
 1. ^         - 字符串开始
 2. ([\[\(])  - 左括号，匹配 [ 或 (
 3. (\d+)     - 数字，匹配最小值
 4. ,         - 逗号，这个是必须要有的分隔符
 5. ([\d\-]+) - 数字或负号，匹配最大值
 6. ([\]\)])  - 右括号，匹配 ] 或 )
 7. $         - 字符串结束

匹配：
 1. [0,10]    - 闭区间
 2. (0,10)    - 开区间
 3. [0,-]     - 左闭右无限
 4. (0,-)     - 左开右无限

不匹配：
 1. [0.1,10]  - 小数
 2. [a,b]     - 字母
 3. {0,10}    - 错误的括号类型

.
*/
var intRangePattern = regexp.MustCompile(`^([\[\(])(\d+),([\d\-]+)([\]\)])$`)

func IntValueRangerChecker(rangeStr string) *IntRange {
	if rangeStr == "" {
		return &IntRange{
			originalStr:  "",
			noRange:      true,
			min:          0,
			inclusiveMin: false,
			max:          0,
			inclusiveMax: false,
		}
	}

	matches := intRangePattern.FindStringSubmatch(rangeStr)
	if matches == nil {
		panic(fmt.Sprintf("IntRange pattern not illegal: %s", rangeStr)) // 值范围描述字符串不合法
	}

	// 开始解析范围
	// matches[0]: 是完整匹配的字符串，从 [1] ~ [n] 对应正则表达式中 () 捕获的子串，所以有：
	// matches[0]: 完整匹配
	// matches[1]: 左括号 [ 或 (
	// matches[2]: 最小值数字
	// matches[3]: 最大值数字或 -
	// matches[4]: 右括号 ] 或 )
	originalStr := matches[0]
	leftBracket := matches[1]
	minStr := matches[2]
	maxStr := matches[3]
	rightBracket := matches[4]

	// fmt.Printf("originalStr: %s, leftBracket: %s, minStr: %s, maxStr: %s, rightBracket: %s\n", originalStr, leftBracket, minStr, maxStr, rightBracket)

	// 判断是否包含最小值（[ 表示包含，( 表示不包含）
	inclusiveMin := leftBracket == "["

	// 解析最小值
	var minVal int64 = 0
	if _, err := fmt.Sscanf(minStr, "%d", &minVal); err != nil {
		panic(fmt.Sprintf("IntRange parse min value failed: %s, err: %v", minStr, err))
	}

	// 解析最大值
	var maxVal int64 = 0
	var noLimitMax = false

	if maxStr == "-" {
		noLimitMax = true // 右边界无限大
	} else {
		if _, err := fmt.Sscanf(maxStr, "%d", &maxVal); err != nil {
			panic(fmt.Sprintf("IntRange parse max value failed: %s, err: %v", maxStr, err))
		}
		if maxVal < minVal {
			panic(fmt.Sprintf("IntRange max value less than min value: min=%d, max=%d, originalStr=%s", minVal, maxVal, originalStr))
		}
	}
	inclusiveMax := rightBracket == "]"

	return &IntRange{
		originalStr:  originalStr,
		noRange:      false,
		min:          minVal,
		inclusiveMin: inclusiveMin,
		max:          maxVal,
		inclusiveMax: inclusiveMax,
		noLimitMax:   noLimitMax,
	}
}

type IntRange struct {
	originalStr string // 原始字符串表示
	noRange     bool   // 没有数值范围限制，是 int 即可

	min          int64
	inclusiveMin bool

	max          int64
	inclusiveMax bool
	noLimitMax   bool
}

func (ir *IntRange) Check(value any) bool {
	if ir.noRange { // 没有值范围限制，是 int/uint 即可
		switch value.(type) {
		case int, int8, int16, int32, int64:
			return true
		case uint, uint8, uint16, uint32, uint64:
			return true
		default:
			fmt.Printf("IntRange check value: [%+v] not int type", value)
			return false
		}
	}

	var i64Value int64
	switch v := value.(type) {
	case int:
		i64Value = int64(v)
	case int8:
		i64Value = int64(v)
	case int16:
		i64Value = int64(v)
	case int32:
		i64Value = int64(v)
	case int64:
		i64Value = v
	case uint:
		if v > math.MaxInt64 {
			fmt.Printf("IntRange check value: uint[%d] over int64 max", value)
			return false
		}
		i64Value = int64(v)
	case uint8:
		i64Value = int64(v)
	case uint16:
		i64Value = int64(v)
	case uint32:
		i64Value = int64(v)
	case uint64:
		if v > math.MaxInt64 {
			fmt.Printf("IntRange check value: uint64[%d] over int64 max", value)
			return false
		}
		i64Value = int64(v)
	default:
		fmt.Printf("IntRange check value: [%+v] not int type", value)
		return false
	}

	if ir.inclusiveMin {
		if i64Value < ir.min {
			return false
		}
	} else {
		if i64Value <= ir.min {
			return false
		}
	}

	if !ir.noLimitMax {
		if ir.inclusiveMax {
			if i64Value > ir.max {
				return false
			}
		} else {
			if i64Value >= ir.max {
				return false
			}
		}
	}
	return true
}

func (ir *IntRange) ToString() string {
	return ""
}

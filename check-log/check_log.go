package checklog

import (
	"fmt"
	"strings"
)

//// 一条检查日志，会有自己的字段名信息，字段名信息由 调用 Check 函数的调用方补全
/// 叶子只会有错误消息
/// node 节点，只会有子节点

// 叶子节点
type CheckLogLeaf struct {
	message string
}

// check 日志节点
type CheckLogNode struct {
	logs  []*CheckLog // 检查日志列表
	isMap bool        // 是否为 map
}

func (cln *CheckLogNode) AddLog(log *CheckLog) {
	cln.logs = append(cln.logs, log)
}

// CheckLog 表示一个检查日志
type CheckLog struct {
	fieldKey string        // 字段名信息
	node     *CheckLogNode // node 和 leaf 是互斥的
	leaf     *CheckLogLeaf
}

func (cl *CheckLog) SetFieldKey(fieldKey string) {
	cl.fieldKey = fieldKey
}

func CheckLogSuccess(msg string) *CheckLog {
	return &CheckLog{
		fieldKey: "",
		node:     nil,
		leaf: &CheckLogLeaf{
			message: msg,
		},
	}
}

func CheckLogFail(msg string) *CheckLog {
	return &CheckLog{
		fieldKey: "",
		node:     nil,
		leaf: &CheckLogLeaf{
			message: msg,
		},
	}
}

func CheckLogWithNode(node *CheckLogNode) *CheckLog {
	return &CheckLog{
		fieldKey: "",
		node:     node,
		leaf:     nil,
	}
}

func CheckLogNodeCommon() *CheckLogNode {
	return &CheckLogNode{
		logs:  make([]*CheckLog, 0),
		isMap: false,
	}
}

func CheckLogNodeMap() *CheckLogNode {
	return &CheckLogNode{
		logs:  make([]*CheckLog, 0),
		isMap: true,
	}
}

func (cl *CheckLog) ToString() string {
	if cl == nil {
		return ""
	}
	return cl.formatString(0)
}

func (cl *CheckLog) formatString(depth int) string {
	// 处理叶子节点
	if cl.leaf != nil {
		if cl.fieldKey == "" {
			return cl.leaf.message
		}
		return fmt.Sprintf("%s = %s", cl.fieldKey, cl.leaf.message)
	}

	// 处理 node 节点
	if cl.node == nil {
		return ""
	}

	var builder strings.Builder

	// 如果有 key，先写入 key
	if cl.fieldKey != "" {
		builder.WriteString(cl.fieldKey)
		builder.WriteString(" = ")
	}

	// 处理 map 类型
	if cl.node.isMap {
		builder.WriteString("{\n")
		count := len(cl.node.logs)

		// 验证是否为偶数（key-value 对）
		if count%2 != 0 {
			builder.WriteString(makeIndent(depth + 1))
			builder.WriteString("ERROR: map logs count is not even\n")
			builder.WriteString(makeIndent(depth))
			builder.WriteString("}")
			return builder.String()
		}

		for i := 0; i < count; i += 2 {
			keyLog := cl.node.logs[i]
			valueLog := cl.node.logs[i+1]

			builder.WriteString(makeIndent(depth + 1))
			builder.WriteString("[")
			builder.WriteString(keyLog.formatString(depth + 1))
			builder.WriteString("] = ")
			builder.WriteString(valueLog.formatString(depth + 1))
			// 添加逗号（最后一项除外）
			if i < count-2 {
				builder.WriteString(",")
			}
			builder.WriteString("\n")
		}

		builder.WriteString(makeIndent(depth))
		builder.WriteString("}")
	} else {
		// 处理普通结构体或数组
		builder.WriteString("{\n")

		for i, log := range cl.node.logs {
			builder.WriteString(makeIndent(depth + 1))
			builder.WriteString(log.formatString(depth + 1))

			// 添加逗号（最后一项除外）
			if i < len(cl.node.logs)-1 {
				builder.WriteString(",")
			}
			builder.WriteString("\n")
		}

		builder.WriteString(makeIndent(depth))
		builder.WriteString("}")
	}

	return builder.String()
}

const formatTab = "  "

func makeIndent(depth int) string {
	if depth <= 0 {
		return ""
	}
	return strings.Repeat(formatTab, depth)
}

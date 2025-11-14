package expandtyperange

import (
	"fmt"
	"reflect"
	"regexp"
)

/*
匹配范围模板
 1. ^                     		- 字符串开始
 2. ([a-zA-Z_][a-zA-Z0-9_]*) 	- 第一个字段名
 3. \.                    		- 点号分隔符
 4. ([a-zA-Z_][a-zA-Z0-9_]*) 	- 第二个字段名
 5. $                     		- 字符串结束
匹配： 严格匹配两级，格式为 xxx.xxx ，如果配置表是 struct 的，第二级直接是 struct 的字段名，如果是 map 或是 list 的话，要求其存放的数据的是 struct，然后，第二级直接是 struct 的字段名
 1. hero.id
 2. _hero._id
.
*/

var refRangePattern = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\.([a-zA-Z_][a-zA-Z0-9_]*)$`)

// 可以用来 ref 引用的类型
// ref 更多的是想实现，某一个配置值，是另一个表的 key，或是表中某一个 field 的值
// 所以支持的类型比较有限
// 目前支持 int 系列，uint 系列，string
var allowTypeSet = map[reflect.Kind]struct{}{
	reflect.Int:   {},
	reflect.Int8:  {},
	reflect.Int16: {},
	reflect.Int32: {},
	reflect.Int64: {},

	reflect.Uint:   {},
	reflect.Uint8:  {},
	reflect.Uint16: {},
	reflect.Uint32: {},
	reflect.Uint64: {},

	reflect.String: {},
}

func RefValueStore() *RefStore {
	return &RefStore{
		oriData: make(map[string]any),

		mapStrRefCheckRule:  make(map[string]map[string]struct{}),
		mapUintRefCheckRule: make(map[string]map[uint64]struct{}),
		mapIntRefCheckRule:  make(map[string]map[int64]struct{}),
	}
}

type RefStore struct {
	oriData map[string]any // 引入的原始数据

	// 缓存的 ref 的 id 规则，这样子就是加载的时候慢点，但是 check 的时候就快多了
	// key: originalStr
	mapStrRefCheckRule  map[string]map[string]struct{}
	mapUintRefCheckRule map[string]map[uint64]struct{}
	mapIntRefCheckRule  map[string]map[int64]struct{}
}

func (rs *RefStore) checkRuleExits(originalStr string) bool {
	_, ok1 := rs.mapStrRefCheckRule[originalStr]
	_, ok2 := rs.mapUintRefCheckRule[originalStr]
	_, ok3 := rs.mapIntRefCheckRule[originalStr]
	return ok1 || ok2 || ok3
}

func (rs *RefStore) CheckUintValue(originalStr string, value uint64) bool {
	idSet, ok := rs.mapUintRefCheckRule[originalStr]
	if !ok {
		return false
	}
	_, ok = idSet[value]
	return ok
}

func (rs *RefStore) CheckIntValue(originalStr string, value int64) bool {
	idSet, ok := rs.mapIntRefCheckRule[originalStr]
	if !ok {
		return false
	}
	_, ok = idSet[value]
	return ok
}

func (rs *RefStore) CheckStrValue(originalStr string, value string) bool {
	idSet, ok := rs.mapStrRefCheckRule[originalStr]
	if !ok {
		return false
	}
	_, ok = idSet[value]
	return ok
}

// 加载一条原始数据，用于后续的值范围检测
// data 经常就是整个配置数据对象
// 这边可以把全部客户端配置都 load 进来
func (rs *RefStore) LoadOneOriData(key string, data any) bool {
	if _, ok := rs.oriData[key]; ok {
		fmt.Printf("refrange dup load ori data, key: %s", key)
		return false
	}

	// 支持的数据类型，在 AddRefCheckRule 的时候在检测
	// struct => k-v 的配置表
	// map => k-v 的配置表
	// array/slice => 列表
	// tData := reflect.TypeOf(data)
	// switch tData.Kind() {
	// case reflect.Struct:
	// case reflect.Map:
	// case reflect.Array, reflect.Slice:
	// default:
	// 	fmt.Printf("refrange load ori data type no support, key: %s, type: %s\n", key, tData.Kind().String())
	// 	return false
	// }

	rs.oriData[key] = data
	return true
}

func (rs *RefStore) AddRefCheckRule(rangeStr string) string {
	matches := refRangePattern.FindStringSubmatch(rangeStr)
	if matches == nil {
		panic(fmt.Sprintf("RefRange pattern not illegal: %s", rangeStr)) // 值范围描述字符串不合法
	}

	// 开始解析范围
	// matches[0]: 是完整匹配的字符串，从 [1] ~ [n] 对应正则表达式中 () 捕获的子串，所以有：
	// matches[0]: key
	// matches[1]: key 对应的数据中的某一个 key
	originalStr := matches[0]
	oriDataKey := matches[1]
	fieldKey := matches[2]

	if rs.checkRuleExits(originalStr) {
		return originalStr // 已经存在这个规则了，可以复用，直接返回了
	}
	oriData, ok := rs.oriData[oriDataKey]
	if !ok {
		panic(fmt.Sprintf("RefRange no load ori data key: %s", oriDataKey)) // 值范围描述字符串不合法
	}

	// 解析 oriData，找到 fieldKey 这个值，获得值的类型
	oriDataType := reflect.TypeOf(oriData)
	switch oriDataType.Kind() {
	case reflect.Struct:
		rs.checkStructFeildType(oriData, originalStr, fieldKey)
	case reflect.Map:
		rs.checkMapFeildType(oriData, originalStr, fieldKey)
	case reflect.Array, reflect.Slice:
		rs.checkListFeildType(oriData, originalStr, fieldKey)
	default:
		panic(fmt.Sprintf("RefRange ori data type no support, key: %s, type: %s", oriDataKey, oriDataType.Kind().String()))
	}

	return originalStr
}

// 配置表是个 struct 的
// 检测 struct 类型的 fieldKey 字段类型是否符合 allowTypeSet 的要求
// 并且转化为对应转化类型，加到缓存中
func (rs *RefStore) checkStructFeildType(oriData any, originalStr, fieldKey string) {
	oriDataType := reflect.TypeOf(oriData)
	if oriDataType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("RefRange ori data type no struct, type: %s", oriDataType.Kind().String()))
	}
	oriDataFieldType, ok := oriDataType.FieldByName(fieldKey)
	if !ok {
		panic(fmt.Sprintf("RefRange ori data no field key: %s", fieldKey))
	}
	if _, ok := allowTypeSet[oriDataFieldType.Type.Kind()]; !ok {
		panic(fmt.Sprintf("RefRange ori data field type no support, key: %s, type: %s", fieldKey, oriDataFieldType.Type.Kind().String()))
	}

	oriDataValue := reflect.ValueOf(oriData)
	oriDataFeileValue := oriDataValue.FieldByName(fieldKey).Interface() // 获取字段的值，在 type 中获得过了，肯定会存在的

	switch v := oriDataFeileValue.(type) {
	case uint64:
		if _, ok := rs.mapUintRefCheckRule[originalStr]; !ok {
			rs.mapUintRefCheckRule[originalStr] = make(map[uint64]struct{})
		}
		rs.mapUintRefCheckRule[originalStr][v] = struct{}{}
	case uint32:
		if _, ok := rs.mapUintRefCheckRule[originalStr]; !ok {
			rs.mapUintRefCheckRule[originalStr] = make(map[uint64]struct{})
		}
		rs.mapUintRefCheckRule[originalStr][uint64(v)] = struct{}{}
	case uint16:
		if _, ok := rs.mapUintRefCheckRule[originalStr]; !ok {
			rs.mapUintRefCheckRule[originalStr] = make(map[uint64]struct{})
		}
		rs.mapUintRefCheckRule[originalStr][uint64(v)] = struct{}{}
	case uint8:
		if _, ok := rs.mapUintRefCheckRule[originalStr]; !ok {
			rs.mapUintRefCheckRule[originalStr] = make(map[uint64]struct{})
		}
		rs.mapUintRefCheckRule[originalStr][uint64(v)] = struct{}{}
	case uint:
		if _, ok := rs.mapUintRefCheckRule[originalStr]; !ok {
			rs.mapUintRefCheckRule[originalStr] = make(map[uint64]struct{})
		}
		rs.mapUintRefCheckRule[originalStr][uint64(v)] = struct{}{}
	case int64:
		if _, ok := rs.mapIntRefCheckRule[originalStr]; !ok {
			rs.mapIntRefCheckRule[originalStr] = make(map[int64]struct{})
		}
		rs.mapIntRefCheckRule[originalStr][v] = struct{}{}
	case int32:
		if _, ok := rs.mapIntRefCheckRule[originalStr]; !ok {
			rs.mapIntRefCheckRule[originalStr] = make(map[int64]struct{})
		}
		rs.mapIntRefCheckRule[originalStr][int64(v)] = struct{}{}
	case int16:
		if _, ok := rs.mapIntRefCheckRule[originalStr]; !ok {
			rs.mapIntRefCheckRule[originalStr] = make(map[int64]struct{})
		}
		rs.mapIntRefCheckRule[originalStr][int64(v)] = struct{}{}
	case int8:
		if _, ok := rs.mapIntRefCheckRule[originalStr]; !ok {
			rs.mapIntRefCheckRule[originalStr] = make(map[int64]struct{})
		}
		rs.mapIntRefCheckRule[originalStr][int64(v)] = struct{}{}
	case int:
		if _, ok := rs.mapIntRefCheckRule[originalStr]; !ok {
			rs.mapIntRefCheckRule[originalStr] = make(map[int64]struct{})
		}
		rs.mapIntRefCheckRule[originalStr][int64(v)] = struct{}{}
	case string:
		if _, ok := rs.mapStrRefCheckRule[originalStr]; !ok {
			rs.mapStrRefCheckRule[originalStr] = make(map[string]struct{})
		}
		rs.mapStrRefCheckRule[originalStr][v] = struct{}{}
	default:
		panic(fmt.Sprintf("RefRange ori data field value type no support, key: %s, type: %T", fieldKey, v))
	}
}

// 配置表是个 map 的
func (rs *RefStore) checkMapFeildType(oriData any, originalStr, fieldKey string) {
	oriDataType := reflect.TypeOf(oriData)
	if oriDataType.Kind() != reflect.Map {
		panic(fmt.Sprintf("RefRange ori data type no map, type: %s", oriDataType.Kind().String()))
	}
	oriDataValueType := oriDataType.Elem() // 检测这个 map 的 value 类型是否是 struct
	if oriDataValueType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("RefRange ori data map elem value type no struct, type: %s", oriDataValueType.Kind().String()))
	}

	// 遍历这个 map，把每一个成员去调用 checkStructFeildType 函数
	oriDataValue := reflect.ValueOf(oriData)
	for _, key := range oriDataValue.MapKeys() {
		mapValue := oriDataValue.MapIndex(key).Interface()
		rs.checkStructFeildType(mapValue, originalStr, fieldKey)
	}
}

// 配置表是个 list 的
func (rs *RefStore) checkListFeildType(oriData any, originalStr, fieldKey string) {
	oriDataType := reflect.TypeOf(oriData)
	if oriDataType.Kind() != reflect.Array && oriDataType.Kind() != reflect.Slice {
		panic(fmt.Sprintf("RefRange ori data type no list, type: %s", oriDataType.Kind().String()))
	}
	oriDataValueType := oriDataType.Elem() // 检测这个 list 的元素类型是否是 struct
	if oriDataValueType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("RefRange ori data list elem value type no struct, type: %s", oriDataValueType.Kind().String()))
	}

	// 遍历这个 list，把每一个成员去调用 checkStructFeildType 函数
	oriDataValue := reflect.ValueOf(oriData)
	for i := 0; i < oriDataValue.Len(); i++ {
		listValue := oriDataValue.Index(i).Interface()
		rs.checkStructFeildType(listValue, originalStr, fieldKey)
	}
}

func RefValueRangerChecker(refStore *RefStore, rangeStr string) *RefRange {
	if refStore == nil {
		panic("RefValueRangerChecker refStore is nil")
	}
	originalStr := refStore.AddRefCheckRule(rangeStr)

	return &RefRange{
		originalStr: originalStr,
		refStore:    refStore,
	}
}

type RefRange struct {
	originalStr string
	refStore    *RefStore
}

func (rf *RefRange) Check(value any) bool {
	switch v := value.(type) {
	case uint64:
		return rf.refStore.CheckUintValue(rf.originalStr, v)
	case uint32:
		return rf.refStore.CheckUintValue(rf.originalStr, uint64(v))
	case uint16:
		return rf.refStore.CheckUintValue(rf.originalStr, uint64(v))
	case uint8:
		return rf.refStore.CheckUintValue(rf.originalStr, uint64(v))
	case uint:
		return rf.refStore.CheckUintValue(rf.originalStr, uint64(v))
	case int64:
		return rf.refStore.CheckIntValue(rf.originalStr, v)
	case int32:
		return rf.refStore.CheckIntValue(rf.originalStr, int64(v))
	case int16:
		return rf.refStore.CheckIntValue(rf.originalStr, int64(v))
	case int8:
		return rf.refStore.CheckIntValue(rf.originalStr, int64(v))
	case int:
		return rf.refStore.CheckIntValue(rf.originalStr, int64(v))
	case string:
		return rf.refStore.CheckStrValue(rf.originalStr, v)
	default:
		fmt.Printf("RefRange check value type no support, type: %T\n", v)
		return false
	}
}

func (rf *RefRange) ToString() string {
	return ""
}

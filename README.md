# value-range
value range check in go
https://github.com/chenjinjie/value-range

1. 经常需要读取配置表，或是通过协议来获取/上报数据，数据来源是其他地方，这时候，我们需要对取来的数据进行范围检测，不仅是某个字段是否存在的检测，更重要的是，对字段的值的范围的检查
2. 使用方法
    * 详细可以到仓库下看 test 中的 TestCfgCheck 方法
3. 大量使用了反射，特别是对于 struct 的检测

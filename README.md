# 结构体校验
## 简介
用于校验由json转换而来的数据。支持：

* 数字
* 字符串（含json.Number）
* 布尔
* 切片
* 字典（map）
* 结构体
* 指针
* 接口（通常为any）

## 使用方式
先`NewChecker`，获得checker，然后`add`各种检测函数，就能使用。

## 结构体校验tag
```go
package checker

const TagIgnore = "c-ignore"  // 忽略

const TagBoolMust = "cb-must" // bool必须为指定值

const TagIntIgnore = "ci-ignore" // 忽略检查
const TagIntMax = "ci-max"       // int最大值
const TagIntMin = "ci-min"       // int最小值
const TagIntZero = "ci-zero"     // int允许为0（ignore忽略ci-min等检查，notcheck则忽略全部检查）
const TagIntCheck = "ci-checker" // 检查函数
const TagIntMust = "ci-must"     // 必须值

const TagStringJsonNumber = "cs-json-number"  // 允许是json number
const TagStringLengthMin = "cs-min"  // 最短长度
const TagStringLengthMax = "cs-max"  // 最大长度
const TagStringLength = "cs-length"  // 固定长度
const TagStringZero = "cs-zero"  // 零值
const TagStringIgnore = "cs-ignore"  // 忽略检查
const TagStringChecker = "cs-checker"  // 检查函数
const TagStringMust = "cs-must"  // 必须值
const TagStringRegex = "cs-regex"  // 检查正则

// 备注：应用在slice上的其他tag（除csl的tag），会被用于其子元素的检查
const TagSliceLengthMin = "csl-min"  // 最大长度
const TagSliceLengthMax = "csl-max"  // 最短长度
const TagSliceLength = "csl-length"  // 必须长度
const TagSliceZero = "csl-zero"  // 零值
const TagSliceIgnore = "csl-ignore"  // 忽略
const TagSliceChecker = "csl-checker"  // 检查函数

// 备注：应用在map上的其他tag，**不会**被用于其子元素的检查
const TagMapLengthMin = "cm-min"  // 最大长度
const TagMapLengthMax = "cm-max"  // 最小长度
const TagMapLength = "cm-length"  // 必须长度
const TagMapZero = "cm-zero"  // 零值
const TagMapIgnore = "cm-ignore"  // 忽略
const TagMapChecker = "cm-checker"  // 检查函数
```
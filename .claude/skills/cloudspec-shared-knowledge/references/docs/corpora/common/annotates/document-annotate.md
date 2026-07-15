文档的注解用于在模型中定义文档信息、废弃信息等。

### A3.1 **document annotate**
#### 说明
给模型中的component增加文档。

#### 选择器
```json
*
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| name | string | 属性的中文名称 |
| nameEn | string | 属性的英文名称 |
| zh | string | 属性中文文档描述 |
| en | string | 属性英文文档描述 |
| exampleValue | string | 示例值 |
| required | Boolean | 是否文档必填，默认值为false |
| url | string | 文档的地址。 |


#### 示例
```json
@document ({
  zh: "状态值"
  en: "Status"
  exampleValue: "ONLINE"
})
Status: string
```

### A3.2 sensitive** annotate**
#### 说明
声明component中存储的值是敏感的。

#### 选择器
```json
:is(simpleType, map, array, struct)
```

敏感的字段只能标记在复合结构的顶层。

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无需参数 | | |


#### 示例
```json
@sensitive
Status: string
```

### A3.3 deprecated** annotate**
#### 说明
标记component为废弃。

#### 选择器
```json
:is(simpleType, array, map, struct, member)
```



#### 支持的属性
支持不传入参数或者支持一个struct类型的参数，当支持参数类型为struct时，字段如下：

| 属性 | 类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| message | string | 必填 | 废弃的的提示信息 |
| substitute | [string] | 非必填 | 替代的component ID，在相同的struct中可以使用相对的ID。 |


#### 示例
```json
struct DiskReplicaPairPropertities {
  @deprecated({
    message: "StatusBefore is deprecated since 2022-03-15"
    substitute: [DiskReplicaPairPropertities$Status]
  })
  StatusBefore: string
	@deprecated
  Status: string
  DiskReplicaPairName: string
  All: boolean
  Description: string
  DestinationDiskId: string
  CreateTime: string
  Period: integer
  Bandwidth: integer
  Tags: TagsSchema
  DiskId: string
  VpcIdOther: string
}
```

以上的示例表示StatusBefore已被废弃，其替代的字段是Status。

### A3.4 openStruct** annotate**
#### 说明
声明数据结构对客户公开，这部分会在开发者门户上透出。

#### 选择器
```json
:is(struct, array, map)
```

#### 支持的属性
当没有特殊配置时，必须设置参数。

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| originName | string | 原始名称。 |
| ref | string | 当前结构直接引用其他公共结构时配置。 |


#### 示例
```json
@openStruct
struct InstanceDetail {
  
}
```

> 以上示例说明InstanceDetail这个数据结构是公开的，SDK、开发者门户中会直接对客户展示这个数据结构。
>

### A3.5 untypical** annotate**
#### 说明
非典型设定，通常只出现在转换的数据结构中，请尽量不要使用。

#### 选择器
```json
*
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| keepObjectEmpty | boolean | 将{}处理为{}而不是空的 object。 |
| keepAdditionalProperties | boolean | 是否保留additionalProperties节点。 |
| keepItGovernance | boolean | 保留itGovernance为空结构。 |
| keepPolicies | boolean | 保持转出的Policies={}，不等于 null。 |
| keepResponseNull | boolean | 在部分场景下仅存在 header 出参没有任何响应值出参，此时使用。 |
| keepStaticInfo | boolean | 是否保留staticInfo字段 |
| keepConditionInfoEmpty | boolean | 保持 ram 下的conditionInfo为空。 |
| originTypeMap | boolean | 资源上原始的类型是 map |
| originName | string | 原始名称，用于 API 的字段或者 API 名称和 YAML 中不一致的场景。 |
| noneFormat | boolean | 在转换 yaml 时是否不需要 format。 |
| noneType | boolean | 字段下不声明 type。 |
| noneSchema | boolean | 声明没有 schema 结构。 |


### A3.6 examples** annotate**
#### 说明
声明参数的示例值。

#### 选择器
```json
*
```

#### 支持的属性
支持的属性为数组，数组的条目为：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| summary | string | 概要信息。 |
| description | string | 描述。 |
| value | any | 属性值，任意类型。 |
| externalValue | string | 外部值。 |
| ref |  string | 引用数据结构。 |



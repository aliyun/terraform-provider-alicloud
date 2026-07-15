# CloudSpec IAM 注解列表

本文档列出了 CloudSpec 身份权限（IAM）相关的所有注解及其详细配置说明。

## 目录

- [A7.1 ram annotate](#a71-ram-annotate)
- [A7.4 conditions annotate](#a74-conditions-annotate)
- [A7.2 defineCondition annotate](#a72-definecondition-annotate)
- [A7.3 defineAction annotate](#a73-defineaction-annotate)
- [A7.5 requiredPermission annotate](#a75-requiredpermission-annotate)
- [A7.6 commonPermissions annotate](#a76-commonpermissions-annotate)
- [A7.7 noNeedAuthorization annotate](#a77-noneedauthorization-annotate)
- [A7.8 otherPermissions annotate](#a78-otherpermissions-annotate)
- [A7.10 effectiveCondition annotate](#a710-effectivecondition-annotate)
- [A7.11 conditionValue annotate](#a711-conditionvalue-annotate)
- [A7.12 conditionInfo annotate](#a712-conditioninfo-annotate)

### A7.1 ram annotate
#### 说明
配置服务接入到RAM服务的配置。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到RAM。 |
| reason | string | 当不开启时的原因。 |
| code | string | 服务在RAM中的code。 |
| accessMethod | string | 接入方式，支持：code（自定义接入）、gateway（网关接入） |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@ram({
  enable: true
  code: "ecs"
  releaseStrategy: {
    mode: "synchronous"
    scope: "increment"
  }
  accessMethod: "code"
})
service Ecs {}
```

### A7.4 conditions annotate
#### 说明
配置服务的鉴权条件列表。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [reference] | 引用defineCondition定义的struct列表。 |

#### 示例
```cspec
@conditions([Condition_acs_TestThis, Condition_acs_TestThat])
service Ecs {}
```

### A7.2 defineCondition annotate
#### 说明
定义鉴权条件。

#### 选择器
```cspec
:is(struct)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| conditionKey | string | 条件的key。required |
| type | string | 条件的类型，支持：String、Numeric、Boolean、Date、IpAddress。required |
| description | string | 条件的描述。 |
| descriptionEn | string | 条件的英文描述。 |
| documentUrl | string | 条件的文档链接。 |
| isConditionKeyFixed | boolean | 条件key是否固定。 |
| isMachineTranslation | boolean | 是否机器翻译。 |
| multipleValues | boolean | 是否支持多值。 |
| enums | [string] | 条件的枚举值。 |

#### 示例
```cspec
@defineCondition
struct Condition_acs_TestThis {
  conditionKey: "acs:TestThis"
  type: "Numeric"
  description: "测试"
  descriptionEn: "test"
  isConditionKeyFixed: false
  isMachineTranslation: true
  multipleValues: true
  enums: ["a", "b"]
}
```

### A7.3 defineAction annotate
#### 说明
定义鉴权Action。

#### 选择器
```cspec
:is(struct)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| action | string | RAM Action。required |
| description | string | Action的描述。 |
| operationType | string | 操作类型，支持：Read、Write、List、Tagging |
| conditions | [reference] | 关联的条件列表，引用defineCondition定义的struct。 |
| resources | [resourceStruct] | 鉴权的资源列表。 |
| otherActions | [reference] | 其他关联的Action。 |

#### resourceStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| resource | string | 资源的ARN模式。required |
| required | boolean | 是否必须。 |
| arnVariables | [arnVariableStruct] | ARN中的变量定义。 |

#### arnVariableStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| variable | string | 变量名。 |
| dataSourceType | string | 数据源类型。 |
| expressType | string | 表达式类型。 |
| express | string | 取值表达式。 |

#### 示例
```cspec
@defineAction
struct CommonPermission_FullAccess {
  action: "amp:*"
  conditions: [Condition_acs_TestThis, Condition_acs_TestThat]
  resources: [{
    `resource`: "acs:amp:*:{#accountId}:*"
    required: true
  }]
}
```

### A7.5 requiredPermission annotate
#### 说明
配置操作的必要权限。

#### 选择器
```cspec
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [reference] | 引用defineAction定义的struct列表。 |

#### 示例
```cspec
@requiredPermission([RAMActionPassRole])
operation CreateEcs {
  input: {}
  output: {}
}
```

### A7.6 commonPermissions annotate
#### 说明
配置服务的公共权限。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [reference] | 引用defineAction定义的struct列表。 |

#### 示例
```cspec
@commonPermissions([CommonPermission_FullAccess])
service Ecs {}
```

### A7.7 noNeedAuthorization annotate
#### 说明
配置操作不需要鉴权的Action。

#### 选择器
```cspec
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [reference] | 引用defineAction定义的struct列表。 |

#### 示例
```cspec
@noNeedAuthorization([RAMActionPassRole])
operation CreateEcs {
  input: {}
  output: {}
}
```

### A7.8 otherPermissions annotate
#### 说明
配置操作的其他权限。

#### 选择器
```cspec
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [reference] | 引用defineAction定义的struct列表。 |

#### 示例
```cspec
@otherPermissions([RAMActionPassRole])
operation CreateEcs {
  input: {}
  output: {}
}
```

### A7.10 effectiveCondition annotate
#### 说明
定义鉴权条件的生效条件，当条件为true时生效，否则不生效。可以作用于defineAction中的resources、conditions、otherActions的引用上。

#### 选择器
```cspec
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| groups | [Statement] | 条件组，多个Statement之间是 or 的关系。 |

#### Statement
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| expresses | [Express] | 多个Express是 and 的关系。 |

#### Express
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| dataSourceType | string | 数据源：requestParam（从当前operation的入参中获取）、systemParam（从系统参数中获取）、hostParam（从host参数中获取）、resourceAttribute（从资源属性中获取）、customExpression（自定义）。required |
| expressType | string | 变量/表达式类型：jsonPath（从JSONPath获取，当dataSourceType为requestParam或resourceAttribute支持）、name（指定名称，当dataSourceType为systemParam时支持） |
| express | string | 取值表达式。例如以jsonPath为例：requestParamsContext.requestParams.ImageId；以name为例：callerUid |
| comparisonOperator | string | 表达式运算，枚举：equals、notEquals、equalsIgnoreCase、notEqualsIgnoreCase、isNull、nonNull、in、notIn、inIgnoreCase、notInIgnoreCase |
| comparisonValue | string | 表达式值 |
| product | string | 云产品 |
| resourceType | string | 资源类型 |

> 注意：当dataSourceType为hostParam时，特指从host中获取region信息，此时expressType和express为空，无需指定；其他情况下express必须赋值；当express从多种数据源取值，dataSourceType需要设置为customExpression，expressType留空；当express为多条时，最后必须显式return Boolean条件的语句。

#### 示例
```cspec
@defineAction
struct RAMActionPassRole {
  action: "ram:PassRole"
  operationType: 'Write'
  resources: [
    @effectiveCondition({
      groups: [{
        expresses: [{
          dataSourceType: "requestParam"
          expressType: "jsonPath"
          express: "return requestParamsContext.requestParams.TagKeys = null"
        }]
      }]
    })
    {
      `resource`: "acs:ram::{#accountId}:role/{#RoleName}"
      arnVariables: [{
        variable: "RoleName"
        dataSourceType: "RequestParam"
        expressType: "name"
        express: "requestParamsContext.requestParams.ImageId"
      }]
    }
  ]
  conditions: [
    @effectiveCondition({
      groups: [{
        expresses: [{
          dataSourceType: "requestParam"
          expressType: "jsonPath"
          express: "return requestParamsContext.requestParams.TagKeys = null"
        }]
      }]
    })
    XCondition, YCondition
  ]
}
```

### A7.11 conditionValue annotate
#### 说明
定义引用的 condition 的取值来源。

#### 选择器
```cspec
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| keyVariableName | string | condition key上的变量名，例：acs:RequestTag/${TagKey}中的${TagKey} |
| keyVariableExpression | Express | condition key上的变量取值表达式 |
| valueExpression | Express | condition key取值表达式 |

Express的取值：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| dataSourceType | string | 数据源：requestParam（从当前operation的入参中获取）、systemParam（从系统参数中获取）、hostParam（从host参数中获取）、resourceAttribute（从资源属性中获取）、customExpression（自定义）。required |
| expressionType | string | 变量/表达式类型：jsonPath（从JSONPath获取，当dataSourceType为requestParam或resourceAttribute支持）、name（指定名称，当dataSourceType为systemParam时支持） |
| expression | string | 取值表达式。例如以jsonPath为例：requestParamsContext.requestParams.ImageId；以name为例：callerUid |
| product | string | 产品 code |
| resourceType | string | 资源名称 |

#### 示例
```cspec
@defineAction
struct RAMActionPassRole {
  action: "ram:PassRole"
  operationType: 'Write'
  resources: [{
    `resource`: "acs:ram::{#accountId}:role/{#RoleName}"
    arnVariables: [{
      variable: "RoleName"
      dataSourceType: "RequestParam"
      expressType: "name"
      express: "requestParamsContext.requestParams.ImageId"
    }]
  }]
  conditions: [
    @conditionValue({
      valueExpression: {
        dataSourceType: "resourceAttribute"
        expressionType: "jsonPath"
        expression: "$.Tags[*].TagKey"
      }
    })
    XCondition, YCondition
  ]
}
```

### A7.12 conditionInfo annotate
#### 说明
指定鉴权 condition 的行为。

#### 选择器
```cspec
:is(component_id)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| validationType | string | 指定鉴权行为：always、conditional |
| resources | [conditionInfoResource] | 支持的资源。 |

#### conditionInfoResource
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| arn | string | 资源的 ARN |
| product | string | 产品名称 |
| resourceType | string | 资源类型 |
| validationType | string | 校验模式 |
| variables | [variablesSchema] | 支持的变量 |

#### variablesSchema
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| name | string | 变量名称 |
| dataSourceType | string | 数据源类型，比如systemParam |
| expression | string | 变量表达式 |
| expressionType | string | 变量类型 |

# CloudSpec Operation 注解参考

## Operation级注解

### @for(ResourceName)

关联操作到资源。主要用于以下场景：
- 资源使用`@flagMode`时，操作必须通过`@for`声明关联的资源，系统据此自动生成input/output
- 操作的input/output struct中使用`$PropertyName`语法引用资源属性时，struct上需要`@for`声明引用的资源

```cspec
@for(Instance)
operation CreateInstance { ... }
```

### @visibility("level")

API可见性。

```cspec
@visibility("Public")     // 公开操作
@visibility("Private")    // 私有操作（默认）
```

### @operationInfo({...})

操作元数据。

```cspec
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "create"
  riskType: "none"
  chargeType: "free"
  abilityTreeCode: "239604"
  abilityTreeNodes: ["FEATUREdcdnL8HX1L"]
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `operationTypeOld` | string | 旧操作类型：`"write"` / `"read"` |
| `typeFromOperation` | string | 操作类型：`"create"` / `"get"` / `"update"` / `"delete"` / `"list"` / `"none"` |
| `riskType` | string | 风险级别：`"none"` / `"high"` |
| `chargeType` | string | 计费类型：`"free"` / `"paid"` |
| `autoTest` | boolean | 是否支持自动化测试 |
| `abilityTreeCode` | string | 能力树代码 |
| `abilityTreeNodes` | [string] | 能力树节点列表 |

对应关系：

| typeFromOperation | operationTypeOld |
|-------------------|------------------|
| `"create"` | `"write"` |
| `"update"` | `"write"` |
| `"delete"` | `"write"` |
| `"get"` | `"read"` |
| `"list"` | `"read"` |

### @http({...})

HTTP协议配置。

```cspec
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
  signatureVersions: ["V4"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `schemes` | object | 各环境支持的协议，通常`{ online: ["https"] }` |
| `methods` | [string] | HTTP方法：`["post"]` / `["get", "post"]` / `["get"]` |
| `authenticators` | [string] | 认证方式：`["AK"]` / `["Anonymous"]` |
| `deprecated` | boolean | 是否废弃 |
| `signatureVersions` | [string] | 签名版本 |
| `requestContentType` | [string] | 请求媒体类型 |
| `responseContentType` | [string] | 响应媒体类型 |

### @backendConfigurationHttp({...})

后端服务路由配置。

```cspec
@backendConfigurationHttp({
  applicationName: "ProductName"
  retries: { online: -1 }
  timeout: { online: 5000 }
  backendUrl: {
    online: "http://vpc_online/api/resource/action#vpc"
    pre: "http://vpc_pre/api/resource/action#vpc"
  }
  sign: true
  signPolicy: "Local"
  requestType: "Object"
  responseType: "Object"
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `applicationName` | string | 后端应用名称，必填 |
| `retries` | object | 重试次数，`-1`表示不重试 |
| `timeout` | object | 超时时间（毫秒），必填 |
| `backendUrl` | object | 各环境的后端URL，必填 |
| `sign` | boolean | 是否签名 |
| `signPolicy` | string | 签名策略 |
| `requestType` | string | 请求类型：`"Object"` / `"String"` |
| `responseType` | string | 响应类型：`"Object"` / `"String"` |

### @errorMapping({...})

错误响应字段映射。

```cspec
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
  dynamicCodeField: "dymCode"
  dynamicMessageField: "dymMessage"
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `errorExpression` | string | 错误判断表达式 |
| `codeField` | string | 错误码字段名 |
| `errorMessageField` | string | 错误信息字段名 |
| `httpStatusCodeField` | string | HTTP状态码字段名 |
| `dynamicCodeField` | string | 动态错误码字段名（可选） |
| `dynamicMessageField` | string | 动态错误信息字段名（可选） |

### @document({...})

操作的文档信息。

```cspec
@document({
  name: "创建实例"
})
```

操作上的`@document`通常只需要`name`字段。完整属性：

| 属性 | 类型 | 说明 |
|-----|------|------|
| `name` | string | 操作的中文名称 |
| `nameEn` | string | 操作的英文名称 |
| `zh` | string | 中文描述 |
| `en` | string | 英文描述 |
| `url` | string | 文档地址 |

### @ram({...})

RAM鉴权配置。

```cspec
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `enable` | boolean | 是否启用RAM鉴权 |
| `level` | string | 鉴权级别：`"operate"` |
| `atGateway` | boolean | 是否在网关鉴权 |

### @requiredPermission(StructName)

关联RAM Action定义结构体。

```cspec
@requiredPermission(CreateInstance_RAM_Action_0)
```

### @defineAction

标记结构体为RAM Action定义。

```cspec
@defineAction
struct CreateInstance_RAM_Action_0 {
  action: "ecs:CreateInstance"
  resources: [{
    `resource`: "acs:ecs:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

注意：`resource`是保留字，必须用反引号包裹。

### @numberPaginated({...})

页码分页配置（List操作使用）。

```cspec
@numberPaginated({
  initialPageNumber: 1
  initialPageSize: 20
  inputPageNumber: "PageNumber"
  inputPageSize: "PageSize"
  recordTotal: "$.TotalCount"
  pageTotal: "PageCount"
  items: "Instances"
})
```

### @paginated({...})

Token分页配置（List操作使用）。

```cspec
@paginated({
  inputToken: "NextToken"
  outputToken: "NextToken"
  maxItems: "MaxResults"
  maxItemsDefault: 10
  pageTruncated: "IsTruncated"
  totalCount: "TotalCount"
  items: "Instances"
})
```

### @rootMapping({...})

响应数据路径映射（List操作使用）。

```cspec
@rootMapping({
  responsePath: "$.Items[*]"
})
```

### @actionTrail({...})

ActionTrail审计日志。

```cspec
@actionTrail({ enable: true })
```

### @returnMode({...})

异步返回模式。

```cspec
@returnMode({
  async: true
  callback: QueryInstanceStatus
  interval: 5
  times: 10
})
```

### @async({...})

异步轮询配置。

```cspec
@async({
  resourceProperty: "$.Status"
  interval: 5
  compareType: "assertEqual"
  successValues: ["running"]
  failedValues: ["failed"]
  times: 10
})
```

### @gatewayOptions({...})

API网关选项。

```cspec
@gatewayOptions({
  responseLog: true
  outputParamVersion: 1
  riskLevel: "high"
  fileTransfer: false
  akProvenStatus: false
  showJsonItemName: false
  keepClientResourceOwnerId: true
})
```

### @controlPolicy({...})

访问策略配置。

```cspec
@controlPolicy({
  online: {
    controlPolicyName: 'policy-name'
    rateLimitPolicy: {
      unit: 'second'
      userRateLimit: 100
      apiRateLimit: 1000
    }
  }
})
```

### @eventInfo({...})

事件配置。

```cspec
@eventInfo({
  enable: true
  eventNames: ["StartInstance"]
})
```

### @retryPolicies([...])

重试策略配置。

```cspec
@retryPolicies([{
  Code: "ConcurrentUpdateError"
  Interval: 1
  Times: 10
}])
```

### @systemParameter([...])

系统参数配置。

```cspec
@systemParameter(["callerUid", "AccessKeyId", "RegionId"])
```

### @rateLimit({...})

流控配置。

```cspec
@rateLimit({
  userDefault: 100
  userCountWindow: 60
  apiDefault: 1000
  apiCountWindow: 60
})
```

## 字段级注解

### @parameter({...})

参数行为配置，用于input/output结构体的字段。

```cspec
@parameter({
  in: "body"
  style: "flat"
  name: "instanceName"
  required: false
  deprecated: false
  sensitive: false
  readOnly: false
  hasServerDefaultValue: false
  minItems: 1
  maxLength: "256"
  minLength: 1
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `in` | string | 参数位置：`"query"` / `"body"` / `"path"` / `"header"` / `"formData"` |
| `style` | string | 序列化风格：`"flat"` / `"json"` |
| `name` | string | 参数名称 |
| `responseName` | string | 响应参数名称（用于output） |
| `required` | boolean | 是否必填 |
| `deprecated` | boolean | 是否废弃 |
| `sensitive` | boolean | 是否敏感 |
| `readOnly` | boolean | 是否只读 |
| `hasServerDefaultValue` | boolean | 服务端是否有默认值 |
| `extMonitorInfo` | boolean | 扩展监控信息 |
| `minItems` | number | 最小数组项数 |
| `maxLength` | string/number | 最大长度 |
| `minLength` | number | 最小长度 |

### @required

标记字段为必填。

```cspec
@required
InstanceName: string
```

### @readonly

标记字段为只读。

```cspec
@readonly
Status: string
```

### @default(value)

设置默认值。

```cspec
@default("")
Description: string

@default(10)
PageSize: int32
```

### @document({...})

字段的文档信息。

```cspec
@document({
  name: "实例名称"
  zh: "ECS实例的名称"
  en: "Name of the ECS instance"
  exampleValue: "my-instance"
})
```

### @length({ min, max })

字符串长度约束。

```cspec
@length({ min: 2, max: 128 })
Name: string
```

### @range({ min, max })

数值范围约束。

```cspec
@range({ min: 1, max: 100 })
MaxResults: int32
```

### @regexPattern("pattern")

正则校验。

```cspec
@regexPattern("^[a-zA-Z][a-zA-Z0-9_-]*$")
Name: string
```

### @sensitive

标记敏感字段。

```cspec
@sensitive
Password: string
```

### @deprecated

标记废弃字段。

```cspec
@deprecated
OldField: string

@deprecated({
  message: "Use NewField instead"
  substitute: [NewField]
})
OldField: string
```

### @clientProhibited

禁止客户端传值。

```cspec
@clientProhibited
@readonly
InternalField: string
```

### @clientOptional

客户端可选（与@required配合使用时，表示服务端会生成值）。

```cspec
@required
@clientOptional
ServerGeneratedField: string
```

### @conflictsWith([...])

声明互斥字段。

```cspec
@conflictsWith(['Id'])
Name: string
@conflictsWith(['Name'])
Id: string
```

### @idempotencyToken

标记幂等参数。

```cspec
@idempotencyToken
ClientToken: string
```

### @hasDefaultValue

服务端存在默认值。

```cspec
@readonly
@hasDefaultValue
CreateTime: string
```

### @format("format")

值格式声明。

```cspec
@format("iso8601")
CreateTime: string

@format("json")
ConfigJson: string
```

### @enums(EnumRef)

关联枚举定义。

```cspec
@enums(StatusEnum)
Status: string
```

### @arrayConfig({...})

数组行为配置。

```cspec
@arrayConfig({
  uniqueItems: true
  unordered: true
})
```

### @hidden

隐藏结构体/字段（不在公开文档中展示）。

```cspec
@hidden
struct InternalStruct { ... }
```

### @openStruct

声明数据结构对客户公开。

```cspec
@openStruct
struct InstanceDetail { ... }
```

### @date({...})

时间格式注解。

```cspec
@date({
  format: "YYYY-MM-DDTHH:mm:ss.sssZ"
})
CreateTime: string
```

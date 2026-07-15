# CloudSpec Operation 设计指南

## Operation 基本结构

Operation（操作）是CloudSpec IDL中的核心组件，定义了API的具体行为。每个Operation代表一个可执行的API接口。

```cspec
operation OperationName {
  input: InputStruct
  output: OutputStruct
  errors: [ErrorStruct]
}
```

> **flagMode特殊说明**：如果资源使用`@flagMode`，操作需要添加`@for(ResourceName)`注解，input/output由系统根据资源属性的flag标记自动生成，操作体内只需声明`errors`。

## 完整配置示例

```cspec
@backendConfigurationHttp({
  applicationName: "ECS"
  retries: { online: -1 }
  timeout: { online: 5000 }
  backendUrl: { online: "http://vpc_online/api/instance/create#vpc" }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "创建实例"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "create"
  riskType: "none"
  chargeType: "free"
})
@ram({ enable: true, level: "operate", atGateway: false })
@requiredPermission(CreateInstance_RAM_Action_0)
@visibility("Public")
operation CreateInstance {
  errors: [Error_CreateInstance]
}

error Error_CreateInstance {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: false
}

@defineAction
struct CreateInstance_RAM_Action_0 {
  action: "ecs:CreateInstance"
  resources: [{
    `resource`: "acs:ecs:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

## 注解顺序规范

操作的注解按以下顺序排列（按需选用）：

```
@for(ResourceName)                  // 仅flagMode资源的操作需要
@actionTrail({...})                 // 审计日志（可选）
@backendConfigurationHttp({...})    // 后端服务配置
@document({...})                    // 文档描述
@errorMapping({...})                // 错误映射
@http({...})                        // HTTP配置
@numberPaginated({...})             // 分页配置（仅List操作）
@operationInfo({...})               // 操作信息
@ram({...})                         // RAM权限配置
@requiredPermission(StructName)     // 关联RAM Action
@rootMapping({...})                 // 响应数据映射（仅List操作）
@visibility("Public")              // 可见性
```

## 标准CRUD操作对照表

| 操作类型 | operationType | typeFromOperation | operationTypeOld | HTTP方法 | 文档name | 特有注解 |
|---------|--------------|-------------------|------------------|----------|---------|---------|
| Create | `"Write"` 或 `"create"` | `"create"` | `"write"` | `["post"]` | "创建{资源名}" | — |
| Get | `"Read"` 或 `"get"` | `"get"` | `"read"` | `["get", "post"]` | "查询{资源名}详情" | — |
| Update | `"Write"` 或 `"update"` | `"update"` | `"write"` | `["post"]` | "修改{资源名}" | — |
| Delete | `"Write"` 或 `"delete"` | `"delete"` | `"write"` | `["post"]` | "删除{资源名}" | — |
| List | `"List"` 或 `"list"` | `"list"` | `"read"` | `["get", "post"]` | "查询{资源名}列表" | `@numberPaginated`, `@rootMapping` |

## 非标准操作设计

非CRUD的自定义操作（如AttachDisk、RebootInstance、MoveResourceGroup等）：

- `@for(ResourceName)` 仍然需要关联资源
- `@operationInfo`中的`typeFromOperation`根据实际情况选择最接近的类型，通常用`"none"`或参考同项目写法
- 操作命名遵循`动词+资源名`格式

## 分页配置

### Token分页（@paginated）

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

### 页码分页（@numberPaginated）

```cspec
@numberPaginated({
  initialPageNumber: 1
  initialPageSize: 20
  inputPageNumber: "PageNumber"
  inputPageSize: "PageSize"
  recordTotal: "$.TotalCount"
})
```

## 异步操作配置

### @returnMode

```cspec
@returnMode({
  async: true
  callback: QueryInstanceStatus
  interval: 5
  times: 10
})
```

### @async

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

## 错误处理设计

### 错误定义

```cspec
error Error_OperationName {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

错误属性说明：

| 属性 | 类型 | 说明 |
|-----|------|------|
| `httpCode` | int32 | HTTP状态码 |
| `errorCode` | string | 公开错误码 |
| `backendErrorCode` | string | 后端内部错误码（可选） |
| `errorMessage` | string | 错误信息 |
| `type` | string | `"user"` 或 `"system"` |
| `default` | boolean | 是否为兜底错误码 |

### 多错误码定义

```cspec
operation GetPage {
  errors: [Error_GetPage_InvalidParam, Error_GetPage]
}

error Error_GetPage_InvalidParam {
  httpCode: 400
  errorCode: "InvalidParameters"
  backendErrorCode: "11001"
  errorMessage: "The specified parameters are invalid."
  type: "user"
  default: false
}

error Error_GetPage {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

## RAM权限设计

```cspec
@defineAction
struct OperationName_RAM_Action_0 {
  action: "{service}:OperationName"
  resources: [{
    `resource`: "acs:{service}:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

- `action`格式为`{小写服务名}:{操作名}`
- `resource`是保留字，必须用反引号包裹
- Create操作的资源ARN使用通配符`*`（资源尚未创建）
- Get/Update/Delete操作可以使用具体的资源ARN模式

## 参数设计

### 输入参数

```cspec
struct CreateInstanceInput {
  @required
  @parameter({
    in: 'query'
    name: 'instanceName'
    required: true
  })
  InstanceName: string

  @parameter({
    in: 'body'
    name: 'config'
    style: 'flat'
  })
  Config: InstanceConfig
}
```

`@parameter`中的`in`值：

| 值 | 说明 |
|---|------|
| `query` | URL查询参数 |
| `body` | 请求体参数 |
| `path` | URL路径参数 |
| `header` | 请求头参数 |
| `formData` | 表单数据 |

### 输出参数

```cspec
struct CreateInstanceOutput {
  @parameter({
    responseName: 'instanceId'
  })
  InstanceId: string

  @parameter({
    responseName: 'requestId'
  })
  RequestId: string
}
```

## 重试策略

```cspec
@retryPolicies([{
  Code: "ConcurrentUpdateError"
  Interval: 1
  Times: 10
}, {
  Code: "InternalServerError"
  Interval: 1
  Times: 10
}])
```

## 设计检查清单

### 基础配置

- [ ] 已添加`@visibility`
- [ ] 已配置`@operationInfo`（operationType、riskType、chargeType）
- [ ] 已配置`@http`（schemes、methods、authenticators）
- [ ] 已配置`@backendConfigurationHttp`（applicationName、timeout、backendUrl）
- [ ] 已配置`@errorMapping`

### 参数设计

- [ ] 输入参数结构清晰，命名规范
- [ ] 必填参数已添加`@required`
- [ ] 参数位置（`in`）配置正确

### 功能配置

- [ ] 如需分页，已配置`@numberPaginated`或`@paginated`
- [ ] 如需异步，已配置`@returnMode`和`@async`
- [ ] 如需重试，已配置`@retryPolicies`

### 一致性

- [ ] 与同项目其他Operation的配置风格一致
- [ ] PascalCase命名
- [ ] 2空格缩进

## 常见问题

### 问题1：注解使用等号而非冒号

```cspec
// 错误
@document(zh = "描述", en = "description")

// 正确
@document({ zh: "描述", en: "description" })
```

### 问题2：注解与目标之间有空行

```cspec
// 错误
@required

Name: string #[C,U,R,L]

// 正确
@required
Name: string #[C,U,R,L]
```

### 问题3：RAM结构体中resource未用反引号

```cspec
// 错误
struct Op_RAM_Action_0 {
  action: "svc:Op"
  resources: [{
    resource: "acs:svc:..."
    required: true
  }]
}

// 正确
@defineAction
struct Op_RAM_Action_0 {
  action: "svc:Op"
  resources: [{
    `resource`: "acs:svc:..."
    required: true
  }]
}
```

### 问题4：flagMode资源的操作缺少@for注解

如果资源使用了`@flagMode`，其关联的操作必须通过`@for`声明关联关系：

```cspec
// 错误（flagMode资源的操作）
@document({ name: "创建实例" })
operation CreateInstance {
  errors: [Error_CreateInstance]
}

// 正确
@for(Instance)
@document({ name: "创建实例" })
operation CreateInstance {
  errors: [Error_CreateInstance]
}
```

非flagMode的操作通过显式声明input/output来关联资源，不需要`@for`。

### 问题5：缺少错误定义

每个操作至少需要关联一个error。`default`字段的使用因项目而异，参考同项目其他操作的写法。

### 问题6：配置与同项目不一致

修改操作时，先检查同项目其他操作的`@backendConfigurationHttp`、`@http`配置，保持一致。

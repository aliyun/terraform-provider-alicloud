# CloudSpec Operation 设计指南

## 📋 目录

1. [Operation 设计基础](#operation-设计基础)
2. [Operation 设计原则](#operation-设计原则)
3. [Operation 配置规范](#operation-配置规范)
4. [Operation 注解详解](#operation-注解详解)
5. [Operation 设计最佳实践](#operation-设计最佳实践)
6. [Operation 设计检查清单](#operation-设计检查清单)
7. [常见问题与解决方案](#常见问题与解决方案)

## Operation 设计基础

### 什么是 Operation？

Operation（操作）是 CloudSpec IDL 中的核心组件，定义了 API 的具体行为。每个 Operation 代表一个可执行的 API 接口，包含输入参数、输出结果和错误处理。

### Operation 的基本结构

```cspec
operation OperationName {
  input: InputStruct
  output: OutputStruct
  errors: [ErrorStruct]
}
```

### Operation 的核心要素

1. **输入参数（input）**：定义 API 接收的参数
2. **输出结果（output）**：定义 API 返回的结果
3. **错误处理（errors）**：定义可能的错误类型
4. **注解配置**：定义各种行为和配置

## Operation 设计要点

### 完整配置示例

每个 Operation 都应该包含完整的配置信息，确保功能清晰明确：

```cspec
@visibility(private)
@operationInfo({
  operationType: 'create'
  riskType: 'high'
  chargeType: 'paid'
  autoTest: true
  abilityTreeCode: "ECS_INSTANCE"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  signatureVersions: ["V4"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
})
@backendConfigurationHttp({
  applicationName: "ECS"
  timeout: { online: 30000 }
  backendUrl: { online: "https://ecs.aliyuncs.com" }
  requestType: "Object"
  responseType: "Object"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpCode"
})
@gatewayOptions({
  responseLog: true
  outputParamVersion: 1
  riskLevel: "high"
})
@controlPolicy({
  online: {
    controlPolicyName: 'ecs-create-policy'
    rateLimitPolicy: {
      unit: 'second'
      userRateLimit: 100
      apiRateLimit: 1000
    }
  }
})
operation CreateInstance {
  input: CreateInstanceInput
  output: CreateInstanceOutput
  errors: [CreateInstanceError]
}
```

### 配置结构规范

Operation 配置应按照以下顺序组织：

```cspec
@visibility(private)                    // 1. 可见性配置
@operationInfo({...})                   // 2. 操作信息配置
@http({...})                           // 3. HTTP 配置
@backendConfigurationHttp({...})       // 4. 后端配置
@errorMapping({...})                   // 5. 错误映射配置
@gatewayOptions({...})                 // 6. 网关配置
@controlPolicy({...})                  // 7. 访问策略配置
operation OperationName {              // 8. Operation 定义
  input: InputStruct
  output: OutputStruct
  errors: [ErrorStruct]
}
```

### 必填配置项

每个 Operation 必须包含以下配置：

- **@visibility**：可见性设置，默认为 private
- **@operationInfo**：操作类型、风险级别、计费类型、测试配置
- **@http**：HTTP 协议配置
- **@backendConfigurationHttp**：后端服务配置
- **@errorMapping**：错误处理配置

### 验证检查

完成 Operation 设计后，必须通过以下验证：

1. **语法检查**：`aliyun cspec build`
2. **规范检查**：`aliyun cspec check --name OperationName`
3. **功能测试**：`aliyun cspec auto`
4. **格式检查**：`aliyun cspec format`

## Operation 配置规范

### 1. 可见性配置（@visibility）

**默认值**：`private`

```cspec
@visibility(private)    // 私有操作（默认）
@visibility(public)     // 公开操作
@visibility(SpecialPurpose) // 特殊用途操作
```

**配置原则**：
- 所有新 Operation 默认为私有
- 只有经过审核的 Operation 才能设为公开
- 特殊用途 Operation 需要明确标识

### 2. 操作信息配置（@operationInfo）

**必填配置**：
```cspec
@operationInfo({
  operationType: 'create' | 'update' | 'get' | 'list' | 'delete' | 'none'
  riskType: 'high' | 'none'
  chargeType: 'paid' | 'free'
  autoTest: true | false
  abilityTreeCode: "string"
  abilityTreeNodes: ["string"]
})
```

**配置说明**：
- `operationType`：操作类型，必填
- `riskType`：风险级别，必填
- `chargeType`：计费类型，必填
- `autoTest`：是否支持自动化测试，必填
- `abilityTreeCode`：能力树代码
- `abilityTreeNodes`：能力树节点列表

### 3. HTTP 配置（@http）

**基础配置**：
```cspec
@http({
  schemes: {
    online: ["https"]
    pre: ["https"]
    gray: ["https"]
    daily: ["https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
  signatureVersions: ["V4"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
})
```

**配置说明**：
- `schemes`：支持的传输协议
- `methods`：HTTP 请求方法
- `authenticators`：认证方式
- `signatureVersions`：签名版本
- `requestContentType`：请求媒体类型
- `responseContentType`：响应媒体类型

### 4. 后端配置（@backendConfigurationHttp）

**必填配置**：
```cspec
@backendConfigurationHttp({
  applicationName: "ServiceName"
  timeout: {
    online: 30000
    pre: 30000
    gray: 30000
    daily: 30000
  }
  backendUrl: {
    online: "https://service.aliyuncs.com"
    pre: "https://pre-service.aliyuncs.com"
    gray: "https://gray-service.aliyuncs.com"
    daily: "https://daily-service.aliyuncs.com"
  }
  requestType: "Object"
  responseType: "Object"
})
```

**配置说明**：
- `applicationName`：后端应用名称，必填
- `timeout`：超时时间（毫秒），必填
- `backendUrl`：后端 URL，必填
- `requestType`：请求类型（Object/String）
- `responseType`：响应类型（Object/String）

### 5. 错误映射配置（@errorMapping）

**基础配置**：
```cspec
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpCode"
})
```

**配置说明**：
- `errorExpression`：错误码判断条件
- `codeField`：后端错误码字段
- `errorMessageField`：错误信息字段
- `httpStatusCodeField`：HTTP 状态码字段

### 6. 网关配置（@gatewayOptions）

**基础配置**：
```cspec
@gatewayOptions({
  responseLog: true
  outputParamVersion: 1
  riskLevel: "high"
  userType: "normal"
})
```

**配置说明**：
- `responseLog`：是否记录返回日志
- `outputParamVersion`：出参版本
- `riskLevel`：风险级别
- `userType`：用户类型

### 7. 访问策略配置（@controlPolicy）

**基础配置**：
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

**配置说明**：
- `controlPolicyName`：访问策略名称
- `rateLimitPolicy`：流控策略配置

## Operation 注解详解

### 1. 分页配置注解

#### @paginated（Token 分页）
```cspec
@paginated({
  inputToken: "NextToken"        // 入参的 next token 参数
  outputToken: "NextToken"       // 出参的 next token 参数
  maxItems: "MaxResults"         // 最大条目数量参数
  maxItemsDefault: 10            // 最大条目默认值
  pageTruncated: "IsTruncated"   // 是否截断参数
  totalCount: "TotalCount"       // 总条目数参数
  items: "Instances"             // 数组节点参数
})
```

#### @numberPaginated（页码分页）
```cspec
@numberPaginated({
  initialPageNumber: 1           // 初始页码
  initialPageSize: 10            // 默认每页大小
  inputPageNumber: "PageNumber"  // 页码参数
  inputPageSize: "PageSize"      // 每页大小参数
  recordTotal: "TotalCount"      // 总条目数参数
  pageTotal: "PageCount"         // 总页数参数
  items: "Instances"             // 数组节点参数
})
```

### 2. 系统参数配置注解

#### @systemParameter
```cspec
@systemParameter(["callerUid", "AccessKeyId", "RegionId"])
```

**支持的系统参数**：
- `callerUid`：调用者 UID
- `AccessKeyId`：访问密钥 ID
- `RegionId`：地域 ID
- `ResourceOwnerId`：资源所有者 ID
- `ResourceOwnerAccount`：资源所有者账号

### 3. 参数映射配置注解

#### @parameter
```cspec
struct InputStruct {
  @parameter({
    in: 'query'                  // 参数位置
    name: 'instanceId'           // 参数名称
    required: true               // 是否必填
    style: 'simple'              // 参数风格
    allowEmptyValue: false       // 是否允许空值
    checkBlank: true             // 是否检查空值
    trim: true                   // 是否删除前后空格
  })
  InstanceId: string
}
```

**参数位置选项**：
- `header`：请求头
- `query`：查询参数
- `path`：路径参数
- `body`：请求体
- `system`：系统参数
- `host`：主机参数
- `formData`：表单数据

### 4. 异步配置注解

#### @returnMode
```cspec
@returnMode({
  async: true                    // 是否异步模式
  callback: QueryInstanceStatus  // 异步查询操作
  interval: 5                    // 重试间隔（秒）
  times: 10                      // 最大重试次数
})
```

#### @async
```cspec
@async({
  resourceProperty: "$.Status"   // 轮询检查的资源属性
  interval: 5                    // 轮询间隔（秒）
  compareType: "assertEqual"     // 比较类型
  successValues: ["running"]     // 成功状态值
  failedValues: ["failed"]       // 失败状态值
  times: 10                      // 最大重试次数
})
```

### 5. 重试策略配置注解

#### @retryPolicies
```cspec
@retryPolicies([{
  Code: "ConcurrentUpdateError"  // 待重试的错误码
  Interval: 1                    // 重试间隔（秒）
  Times: 10                      // 最大重试次数
}, {
  Code: "InternalServerError"
  Interval: 1
  Times: 10
}])
```

### 6. 事件配置注解

#### @eventInfo
```cspec
@eventInfo({
  enable: true                   // 是否开启事件驱动
  eventNames: ["StartInstance"]  // 关联的事件名称
})
```

### 7. 流控配置注解

#### @rateLimit
```cspec
@rateLimit({
  userDefault: 100               // 用户默认限制
  userCountWindow: 60            // 用户统计窗口（秒）
  apiDefault: 1000               // API 默认限制
  apiCountWindow: 60             // API 统计窗口（秒）
})
```

## Operation 设计最佳实践

### 1. 命名规范

#### Operation 命名
```cspec
// 好的命名
operation CreateInstance        // 创建实例
operation UpdateInstance       // 更新实例
operation GetInstance          // 获取实例
operation ListInstances        // 列出实例
operation DeleteInstance       // 删除实例

// 不好的命名
operation create               // 小写开头
operation Create               // 动词不明确
operation InstanceCreate       // 顺序错误
```

#### 结构体命名
```cspec
// 好的命名
struct CreateInstanceInput     // 创建实例输入
struct CreateInstanceOutput    // 创建实例输出
struct CreateInstanceError     // 创建实例错误

// 不好的命名
struct Input                   // 过于简单
struct CreateInput             // 不完整
struct InstanceCreateInput     // 顺序错误
```

### 2. 参数设计

#### 输入参数设计
```cspec
struct CreateInstanceInput {
  @required
  @parameter({
    in: 'query'
    name: 'instanceName'
    description: '实例名称'
  })
  InstanceName: string
  
  @parameter({
    in: 'query'
    name: 'instanceType'
    description: '实例类型'
  })
  InstanceType: string
  
  @parameter({
    in: 'body'
    name: 'config'
    description: '实例配置'
  })
  Config: InstanceConfig
}
```

#### 输出参数设计
```cspec
struct CreateInstanceOutput {
  @required
  @parameter({
    responseName: 'instanceId'
    description: '实例ID'
  })
  InstanceId: string
  
  @parameter({
    responseName: 'requestId'
    description: '请求ID'
  })
  RequestId: string
  
  @parameter({
    responseName: 'instanceInfo'
    description: '实例信息'
  })
  InstanceInfo: InstanceInfo
}
```

### 3. 错误处理设计

#### 错误结构设计
```cspec
struct CreateInstanceError {
  @required
  @parameter({
    responseName: 'code'
    description: '错误码'
  })
  Code: string
  
  @required
  @parameter({
    responseName: 'message'
    description: '错误信息'
  })
  Message: string
  
  @parameter({
    responseName: 'requestId'
    description: '请求ID'
  })
  RequestId: string
}
```

#### 错误映射配置
```cspec
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpCode"
  dynamicCodeField: "dymCode"
  dynamicMessageField: "dymMessage"
})
```

### 4. 分页设计

#### Token 分页
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
operation ListInstances {
  input: ListInstancesInput
  output: ListInstancesOutput
  errors: [ListInstancesError]
}
```

#### 页码分页
```cspec
@numberPaginated({
  initialPageNumber: 1
  initialPageSize: 10
  inputPageNumber: "PageNumber"
  inputPageSize: "PageSize"
  recordTotal: "TotalCount"
  pageTotal: "PageCount"
  items: "Instances"
})
operation ListInstances {
  input: ListInstancesInput
  output: ListInstancesOutput
  errors: [ListInstancesError]
}
```

### 5. 异步操作设计

#### 异步创建
```cspec
@returnMode({
  async: true
  callback: QueryInstanceStatus
  interval: 5
  times: 10
})
@async({
  resourceProperty: "$.Status"
  interval: 5
  compareType: "assertEqual"
  successValues: ["running"]
  failedValues: ["failed"]
  times: 10
})
operation CreateInstance {
  input: CreateInstanceInput
  output: CreateInstanceOutput
  errors: [CreateInstanceError]
}
```

#### 异步查询
```cspec
@operationInfo({
  operationType: 'get'
  autoTest: true
})
operation QueryInstanceStatus {
  input: QueryInstanceStatusInput
  output: QueryInstanceStatusOutput
  errors: [QueryInstanceStatusError]
}
```

### 6. 配置一致性

#### 参考现有实现
```cspec
// 参考同空间其他 Operation 的配置方式
@visibility(private)                    // 默认私有
@operationInfo({
  operationType: 'create'              // 参考其他 create 操作
  riskType: 'high'                     // 参考其他高风险操作
  chargeType: 'paid'                   // 参考其他付费操作
  autoTest: true                       // 参考其他测试配置
  abilityTreeCode: "ECS_INSTANCE"      // 参考其他能力树配置
})
@http({
  schemes: { online: ["https"] }       // 参考其他 HTTP 配置
  methods: ["post"]                    // 参考其他 POST 方法
  authenticators: ["AK"]               // 参考其他认证配置
})
@backendConfigurationHttp({
  applicationName: "ECS"               // 参考其他后端配置
  timeout: { online: 30000 }           // 参考其他超时配置
  backendUrl: { online: "https://ecs.aliyuncs.com" } // 参考其他 URL 配置
})
```

## Operation 设计检查清单

### 1. 基础配置检查

- [ ] **可见性配置**：已添加 `@visibility(private)`
- [ ] **操作信息配置**：已配置 `@operationInfo`
  - [ ] `operationType` 已设置
  - [ ] `riskType` 已设置
  - [ ] `chargeType` 已设置
  - [ ] `autoTest` 已设置
- [ ] **HTTP 配置**：已配置 `@http`
  - [ ] `schemes` 已设置
  - [ ] `methods` 已设置
  - [ ] `authenticators` 已设置
- [ ] **后端配置**：已配置 `@backendConfigurationHttp`
  - [ ] `applicationName` 已设置
  - [ ] `timeout` 已设置
  - [ ] `backendUrl` 已设置

### 2. 参数设计检查

- [ ] **输入参数**：结构清晰，命名规范
- [ ] **输出参数**：结构清晰，命名规范
- [ ] **错误处理**：已定义错误结构
- [ ] **参数注解**：已添加 `@parameter` 注解
- [ ] **必填参数**：已添加 `@required` 注解

### 3. 功能配置检查

- [ ] **分页配置**：如需要，已配置分页注解
- [ ] **异步配置**：如需要，已配置异步注解
- [ ] **重试策略**：如需要，已配置重试策略
- [ ] **错误映射**：已配置 `@errorMapping`
- [ ] **系统参数**：如需要，已配置 `@systemParameter`

### 4. 安全配置检查

- [ ] **访问策略**：已配置 `@controlPolicy`
- [ ] **网关配置**：已配置 `@gatewayOptions`
- [ ] **流控配置**：如需要，已配置 `@rateLimit`
- [ ] **事件配置**：如需要，已配置 `@eventInfo`

### 5. 一致性检查

- [ ] **命名规范**：遵循现有命名约定
- [ ] **配置风格**：参考同空间其他 Operation
- [ ] **结构设计**：参考现有实现模式
- [ ] **注解使用**：参考现有注解配置

### 6. 验证检查

- [ ] **语法检查**：通过 `aliyun cspec build`
- [ ] **规范检查**：通过 `aliyun cspec check`
- [ ] **功能测试**：通过 `aliyun cspec auto`
- [ ] **格式检查**：通过 `aliyun cspec format`

## 常见问题与解决方案

### 问题1：Operation 配置不完整

#### 问题描述
Operation 缺少必要的配置注解，导致功能不完整。

#### 解决方案
```cspec
// 问题：缺少基础配置
operation CreateInstance {
  input: CreateInstanceInput
  output: CreateInstanceOutput
  errors: [CreateInstanceError]
}

// 解决：添加完整配置
@visibility(private)
@operationInfo({
  operationType: 'create'
  riskType: 'high'
  chargeType: 'paid'
  autoTest: true
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
})
@backendConfigurationHttp({
  applicationName: "ECS"
  timeout: { online: 30000 }
  backendUrl: { online: "https://ecs.aliyuncs.com" }
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
})
operation CreateInstance {
  input: CreateInstanceInput
  output: CreateInstanceOutput
  errors: [CreateInstanceError]
}
```

### 问题2：参数映射配置错误

#### 问题描述
参数映射配置不正确，导致参数传递失败。

#### 解决方案
```cspec
// 问题：参数映射配置错误
struct CreateInstanceInput {
  InstanceName: string
  InstanceType: string
}

// 解决：添加正确的参数映射
struct CreateInstanceInput {
  @required
  @parameter({
    in: 'query'
    name: 'instanceName'
    required: true
  })
  InstanceName: string
  
  @parameter({
    in: 'query'
    name: 'instanceType'
    required: true
  })
  InstanceType: string
}
```

### 问题3：分页配置不正确

#### 问题描述
分页配置不正确，导致分页功能异常。

#### 解决方案
```cspec
// 问题：分页配置不完整
@paginated({
  inputToken: "NextToken"
  outputToken: "NextToken"
})

// 解决：添加完整的分页配置
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

### 问题4：异步配置不完整

#### 问题描述
异步配置不完整，导致异步操作失败。

#### 解决方案
```cspec
// 问题：异步配置不完整
@returnMode({
  async: true
})

// 解决：添加完整的异步配置
@returnMode({
  async: true
  callback: QueryInstanceStatus
  interval: 5
  times: 10
})
@async({
  resourceProperty: "$.Status"
  interval: 5
  compareType: "assertEqual"
  successValues: ["running"]
  failedValues: ["failed"]
  times: 10
})
```

### 问题5：错误处理配置不完整

#### 问题描述
错误处理配置不完整，导致错误信息不准确。

#### 解决方案
```cspec
// 问题：错误映射配置不完整
@errorMapping({
  codeField: "code"
})

// 解决：添加完整的错误映射配置
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpCode"
  dynamicCodeField: "dymCode"
  dynamicMessageField: "dymMessage"
})
```

### 问题6：配置不一致

#### 问题描述
Operation 配置与同空间其他 Operation 不一致。

#### 解决方案
```cspec
// 问题：配置风格不一致
@http({
  schemes: ["https"]
  methods: ["POST"]
})

// 解决：参考同空间其他 Operation 的配置风格
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  signatureVersions: ["V4"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
})
```

## 总结

Operation 设计是 CloudSpec IDL 开发的核心环节。通过遵循本指南的配置规范和最佳实践，可以设计出高质量、可维护的 Operation，为 CloudSpec 生态的发展奠定坚实基础。

### 关键要点

- **完整配置**：每个 Operation 都应包含完整的注解配置
- **规范命名**：遵循统一的命名约定和结构设计
- **一致性**：参考同空间其他 Operation 的配置方式
- **验证检查**：通过完整的验证流程确保质量

---

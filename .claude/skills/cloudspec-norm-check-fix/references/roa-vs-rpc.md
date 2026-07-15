# RPC vs ROA（RESTful）风格完整对照

## 一、总览

| 维度 | RPC | ROA（RESTful） |
|------|-----|----------------|
| apiStyle声明位置 | service级别`@apiStyle("rpc")` | 每个操作的`@http`中`apiStyle: "restful"` |
| 操作语义 | 通过Action名称区分 | 通过HTTP方法 + URI路径区分 |
| URI路径 | 无 | 必须有，含路径参数占位符 |
| 参数传递 | query / formData / system | path / query / body / header |
| HTTP方法 | post, get | post, get, put, delete |
| 请求格式 | form-urlencoded | application/json |
| 响应格式 | 不声明 | application/json |
| schemes | `["http", "https"]` | `["https"]` |

## 二、@http注解对比

### RPC

```cspec
@http({
  apiStyle: "RPC"
  schemes: {
    online: ["http", "https"]
    daily: ["http", "https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
  variables: {
    daily: [{ key: "$.policies.rateLimitPolicy.apiRateLimit", value: 500 }]
    pre: []
  }
  deprecated: false
})
```

### ROA

```cspec
@http({
  apiStyle: "restful"
  schemes: {
    online: ["https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
  uri: "/2023-03-30/functions/{functionName}/sessions"
  variables: {
    daily: []
    pre: []
  }
  deprecated: false
})
```

**关键差异**：

| 字段 | RPC | ROA |
|------|-----|-----|
| `apiStyle` | `"RPC"` | `"restful"` |
| `uri` | 不存在 | **必须**，RESTful路径 |
| `requestContentType` | 不存在 | `["application/json"]` |
| `responseContentType` | 不存在 | `["application/json"]` |
| `schemes` | 通常`["http", "https"]` | 通常`["https"]` |
| `methods` | `["post"]`或`["get", "post"]` | 按REST语义选择 |

### ROA HTTP方法与操作类型对应

| 操作类型 | HTTP方法 |
|---------|----------|
| Create | `["post"]` |
| Get | `["get"]` |
| Update | `["put"]` |
| Delete | `["delete"]` |
| List | `["get"]` |

## 三、@parameter参数位置对比

### RPC参数位置

| in值 | 说明 | 示例 |
|------|------|------|
| `system` | 系统参数，不对外暴露 | `Action`、`callerUid` |
| `query` | URL查询参数 | `ConfigRuleId` |
| `formData` | POST表单数据 | `ConfigRuleName`、`Description` |

```cspec
struct CreateConfigRule_Input {
  @backendName("action")
  @parameter({
    checkBlank: false
    in: "system"
  })
  Action: string

  @backendName("configRuleName")
  @parameter({
    checkBlank: false
    in: "formData"
  })
  ConfigRuleName: string
}
```

### ROA参数位置

| in值 | 说明 | 示例 |
|------|------|------|
| `path` | URI路径参数，必须在`uri`中有`{name}`占位符 | `functionName` |
| `query` | URL查询参数 | `qualifier`、`limit` |
| `body` | JSON请求体字段 | `sessionTTLInSeconds` |
| `header` | HTTP请求头 | `X-Fc-Log-Type` |

```cspec
struct CreateSession_Input {
  @backendName("functionName")
  @parameter({
    in: "path"
  })
  @required
  functionName: string

  @backendName("qualifier")
  @parameter({
    in: "query"
  })
  qualifier: string

  @backendName("sessionTTLInSeconds")
  @parameter({
    in: "body"
  })
  sessionTTLInSeconds: int64
}
```

**注意**：`path`参数的`@backendName`值必须与`@http`中`uri`里的`{占位符名}`完全一致。

## 四、@backendConfigurationHttp对比

### RPC

```cspec
@backendConfigurationHttp({
  applicationName: "aliyun-config-services"
  retries: { online: -1, pre: -1, daily: -1 }
  timeout: { online: 10000 }
  backendUrl: {
    online: "http://config_vpc_reverse_authorization/config-rule/create#vpc"
  }
  sign: true
  signPolicy: "Local"
})
```

### ROA

```cspec
@backendConfigurationHttp({
  applicationName: "FC"
  requestType: "String"
  responseType: "String"
  retries: { online: -1, pre: -1 }
  timeout: { online: 60000 }
  backendUrl: {
    online: "http://pop.fc3.prod/#vpc"
  }
  signKeyName: {
    online: "FC-API-V20210406"
  }
  sign: true
  signPolicy: "Local"
})
```

**关键差异**：ROA需要`requestType: "String"`和`responseType: "String"`。

## 五、@gatewayOptions对比

### RPC

```cspec
@gatewayOptions({
  responseLog: true
  fileTransfer: false
  akProvenStatus: false
  showJsonItemName: false
  keepClientResourceOwnerId: true
})
```

### ROA

```cspec
@gatewayOptions({
  responseLog: true
  fileTransfer: false
  akProvenStatus: false
  showJsonItemName: false
  keepClientResourceOwnerId: true
  roaRequestBodyLog: true
  roaResponseBodyLog: true
  outputParamVersion: 2
})
```

**关键差异**：ROA需要`roaRequestBodyLog`、`roaResponseBodyLog`和`outputParamVersion: 2`。

## 六、错误处理对比

### RPC - 独立error定义

```cspec
operation CreateConfigRule {
  input: CreateConfigRule_Input
  output: CreateConfigRule_Output_Property
  errors: [Error_CreateConfigRule_ExceedMaxRuleCount, Error_CreateConfigRule_ServiceUnavailable]
}

error Error_CreateConfigRule_ExceedMaxRuleCount {
  httpCode: 400
  errorCode: "ExceedMaxRuleCount"
  backendErrorCode: "ExceedMaxRuleCount"
  errorMessage: "The maximum number of rules is exceeded."
  type: "user"
  default: false
}
```

### ROA - @multipleOutput按状态码区分

ROA风格可能使用`@multipleOutput`按HTTP状态码定义不同的响应结构：

```cspec
@multipleOutput([{
  httpCode: "204"
  outputStructure: DeleteSession_Struct_204_0
  empty: true
}, {
  httpCode: "400"
  outputStructure: DeleteSession_Struct_400_0
  empty: true
}, {
  httpCode: "500"
  outputStructure: DeleteSession_Struct_500_0
  empty: true
}])
```

ROA也可以使用独立error定义，取决于项目惯例。关键是参考同项目已有操作的写法。

## 七、service定义对比

### RPC

```cspec
@apiStyle("rpc")
@commonPermissions([Config_Common_Action_FullAccess])
@runtimeType("pop")
@visibility("Public")
service Config {
  version: "2020-09-07"
  resources: [Rule, Aggregator, ...]
  operations: [ActiveAggregateConfigRules, ...]
}
```

### ROA

```cspec
@runtimeType("pop")
@visibility("Public")
service FC {
  version: "2023-03-30"
}
```

**关键差异**：
- RPC在service上声明`@apiStyle("rpc")`，ROA不在service上声明
- RPC通常在service中列出所有resources和operations，ROA可能不列出

## 八、RAM Action ARN对比

### RPC

```cspec
@defineAction
struct CreateConfigRule_RAM_Action_0 {
  action: "config:CreateConfigRule"
  resources: [{
    `resource`: "acs:config:*:{#accountId}:rule/*"
    required: true
    serviceCode: "Config"
    resourceCode: "Rule"
  }]
}
```

### ROA

```cspec
@defineAction
struct CreateSession_RAM_Action_0 {
  action: "fc:CreateSession"
  resources: [{
    `resource`: "acs:fc:{#regionId}:{#accountId}:functions/{#functionName}"
    required: true
    serviceCode: "FCV3"
    resourceCode: "Function"
  }]
}
```

**关键差异**：
- RPC使用`*`通配符（如`acs:config:*:{#accountId}:rule/*`）
- ROA使用具体变量引用（如`acs:fc:{#regionId}:{#accountId}:functions/{#functionName}`）
- ROA的RAM Action可能包含`serviceCode`和`resourceCode`字段

## 九、ROA特有构造

### @openStruct

ROA项目常使用`@openStruct`标记的共享结构体，通常放在`./openStructs/`目录：

```cspec
@document({ name: "Session" })
@openStruct({ originName: "Session" })
struct open_struct_Session {
  @backendName("sessionId")
  @parameter({ responseName: "sessionId" })
  sessionId: string
  // ...
}
```

操作可以直接引用这些struct作为output：

```cspec
operation CreateSession {
  input: CreateSession_Input
  output: open_struct_Session
}
```

### @multipleOutput

ROA操作可能使用`@multipleOutput`注解在input struct上，按HTTP状态码定义不同响应：

```cspec
@multipleOutput([{
  httpCode: "200"
  outputStructure: GetFunction_Output
}, {
  httpCode: "404"
  outputStructure: Error_Output
  empty: true
}])
struct GetFunction_Input { ... }
```

## 十、书写检查清单

### 新建操作时的风格检查

**RPC项目：**

- [ ] `@http`中`apiStyle`为`"RPC"`
- [ ] `@http`中无`uri`字段
- [ ] `@http`中无`requestContentType`/`responseContentType`
- [ ] 参数使用`in: "query"`或`in: "formData"`
- [ ] 有`in: "system"`的`Action`参数
- [ ] `schemes`包含`["http", "https"]`

**ROA项目：**

- [ ] `@http`中`apiStyle`为`"restful"`
- [ ] `@http`中有`uri`字段，路径正确
- [ ] `@http`中有`requestContentType`和`responseContentType`
- [ ] URI中的`{placeholder}`与`in: "path"`参数的`@backendName`一致
- [ ] 请求体参数使用`in: "body"`
- [ ] 无`in: "system"`参数
- [ ] `@backendConfigurationHttp`包含`requestType`和`responseType`
- [ ] `@gatewayOptions`包含`roaRequestBodyLog`、`roaResponseBodyLog`、`outputParamVersion`
- [ ] `schemes`为`["https"]`
- [ ] HTTP方法与操作语义匹配（POST/GET/PUT/DELETE）

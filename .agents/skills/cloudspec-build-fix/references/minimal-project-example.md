# CloudSpec 最小完整项目示例

> **重要**：以下示例仅展示一个可编译通过的最小项目骨架，用于理解各组件之间的引用关系和必填注解。**实际修复编译错误时，必须以用户项目中已有文件的风格、注解配置为准**，不要照搬本示例的具体值。

---

## 一、项目目录结构

```
project/
├── main.cspec              # 服务定义 + import
├── resources/
│   └── Item.cspec           # 资源定义
└── operations/
    └── GetItem.cspec         # 操作定义
```

---

## 二、RPC 风格最小项目

### main.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@apiStyle("rpc")
@runtimeType("pop")
@visibility("Public")
service ExampleService {
  version: "2024-01-01"
  resources: [Item]
  operations: []
}
```

### resources/Item.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@arn("acs:example:{#regionId}:{#accountId}:item/{#ItemId}")
@document({
  name: "示例资源"
  zh: "示例资源"
  en: "Example Item"
})
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Item {
  identifyDefinition: {
    ItemId: string
  }
  properties: {
    @readonly
    @document({ name: "资源名称", zh: "资源名称", en: "Item name" })
    ItemName: string

    @readonly
    @document({ name: "状态", zh: "资源状态", en: "Status" })
    Status: string
  }
  get: GetItem
}
```

### operations/GetItem.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@backendConfigurationHttp({
  applicationName: "example-service"
  retries: { online: -1 }
  timeout: { online: 10000 }
  backendUrl: {
    online: "http://example-service/item/get#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询示例资源详情"
})
@errorMapping({
  errorCodeField: "Code"
  errorMessageField: "Message"
  successCode: "200"
})
@http({
  apiStyle: "RPC"
  schemes: {
    online: ["http", "https"]
  }
  methods: ["get", "post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationType: "Read"
})
@ram({
  ramCode: "example"
})
@requiredPermission(GetItem_RAM_Action_0)
@visibility("Public")
operation GetItem {
  input: GetItemInput
  output: GetItemOutput
  errors: [Error_GetItem]
}

struct GetItemInput {
  @required
  @parameter({
    in: "query"
    name: "ItemId"
    required: true
  })
  @document({ name: "资源ID", zh: "资源ID", en: "Item ID" })
  ItemId: string
}

struct GetItemOutput {
  @parameter({ responseName: "ItemId" })
  @document({ name: "资源ID", zh: "资源ID", en: "Item ID" })
  ItemId: string

  @parameter({ responseName: "ItemName" })
  @document({ name: "资源名称", zh: "资源名称", en: "Item name" })
  ItemName: string

  @parameter({ responseName: "Status" })
  @document({ name: "状态", zh: "资源状态", en: "Status" })
  Status: string
}

error Error_GetItem {
  httpCode: 404
  errorCode: "ItemNotFound"
  backendErrorCode: "ItemNotFound"
  errorMessage: "The specified item does not exist."
  type: "user"
  default: false
}

@defineAction
struct GetItem_RAM_Action_0 {
  action: "example:GetItem"
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:item/{#ItemId}"
    required: true
  }]
}
```

---

## 三、ROA 风格最小项目

### main.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@runtimeType("pop")
@visibility("Public")
service ExampleService {
  version: "2024-01-01"
  resources: [Item]
  operations: []
}
```

> **注意**：ROA 风格不在 service 上声明 `@apiStyle`，而是在每个操作的 `@http` 中声明 `apiStyle: "restful"`。

### resources/Item.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@arn("acs:example:{#regionId}:{#accountId}:item/{#ItemId}")
@document({
  name: "示例资源"
  zh: "示例资源"
  en: "Example Item"
})
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Item {
  identifyDefinition: {
    ItemId: string
  }
  properties: {
    @readonly
    @document({ name: "资源名称", zh: "资源名称", en: "Item name" })
    ItemName: string

    @readonly
    @document({ name: "状态", zh: "资源状态", en: "Status" })
    Status: string
  }
  get: GetItem
}
```

### operations/GetItem.cspec

```cspec
$version: 1
namespace: alicloud.Example.ExampleService.v20240101

@backendConfigurationHttp({
  applicationName: "example-service"
  requestType: "String"
  responseType: "String"
  retries: { online: -1 }
  timeout: { online: 10000 }
  backendUrl: {
    online: "http://example-service-pop/#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询示例资源详情"
})
@errorMapping({
  errorCodeField: "Code"
  errorMessageField: "Message"
  successCode: "200"
})
@gatewayOptions({
  responseLog: true
  roaRequestBodyLog: true
  roaResponseBodyLog: true
  outputParamVersion: 2
})
@http({
  apiStyle: "restful"
  schemes: {
    online: ["https"]
  }
  methods: ["get"]
  authenticators: ["AK"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
  uri: "/2024-01-01/items/{itemId}"
  deprecated: false
})
@operationInfo({
  operationType: "Read"
})
@ram({
  ramCode: "example"
})
@requiredPermission(GetItem_RAM_Action_0)
@visibility("Public")
operation GetItem {
  input: GetItemInput
  output: GetItemOutput
  errors: [Error_GetItem]
}

struct GetItemInput {
  @required
  @backendName("itemId")
  @parameter({
    in: "path"
  })
  @document({ name: "资源ID", zh: "资源ID", en: "Item ID" })
  ItemId: string
}

struct GetItemOutput {
  @parameter({ responseName: "itemId" })
  @document({ name: "资源ID", zh: "资源ID", en: "Item ID" })
  ItemId: string

  @parameter({ responseName: "itemName" })
  @document({ name: "资源名称", zh: "资源名称", en: "Item name" })
  ItemName: string

  @parameter({ responseName: "status" })
  @document({ name: "状态", zh: "资源状态", en: "Status" })
  Status: string
}

error Error_GetItem {
  httpCode: 404
  errorCode: "ItemNotFound"
  backendErrorCode: "ItemNotFound"
  errorMessage: "The specified item does not exist."
  type: "user"
  default: false
}

@defineAction
struct GetItem_RAM_Action_0 {
  action: "example:GetItem"
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:item/{#ItemId}"
    required: true
  }]
}
```

---

## 四、RPC 与 ROA 关键差异速查

修复编译错误时，注意以下差异点（完整对照见 [roa-vs-rpc.md](roa-vs-rpc.md)）：

| 组件 | RPC | ROA |
|------|-----|-----|
| **service** | 有 `@apiStyle("rpc")` | 无 `@apiStyle` |
| **@http** | 无 `uri`、无 `requestContentType` | 必须有 `uri`、`requestContentType`、`responseContentType` |
| **@http.apiStyle** | `"RPC"` | `"restful"` |
| **@http.schemes** | `["http", "https"]` | `["https"]` |
| **@http.methods** | `["post"]` 或 `["get", "post"]` | 按 REST 语义选择 |
| **@backendConfigurationHttp** | 无 `requestType`/`responseType` | 必须有 `requestType: "String"` 和 `responseType: "String"` |
| **@gatewayOptions** | 无 ROA 专属字段 | 必须有 `roaRequestBodyLog`、`roaResponseBodyLog`、`outputParamVersion: 2` |
| **@parameter.in** | `"query"` / `"formData"` / `"system"` | `"path"` / `"query"` / `"body"` / `"header"` |
| **path 参数** | 不适用 | `@backendName` 必须与 `uri` 中 `{占位符}` 一致 |

---

## 五、使用原则

1. **以用户项目为准**：本示例的注解值（`applicationName`、`backendUrl`、`ramCode` 等）都是占位值。修复时必须参考用户项目中同类操作的实际配置。
2. **参考同项目已有文件**：新建或修复操作时，优先参考同项目中其他已有操作的注解风格和配置模式，保持一致性。
3. **骨架结构是固定的**：虽然具体值因项目而异，但注解的组合和引用关系（operation → struct → error → RAM Action）是固定的。
4. **保留字转义**：RAM Action 中的 `resource` 是保留字，必须用反引号包裹为 `` `resource` ``。

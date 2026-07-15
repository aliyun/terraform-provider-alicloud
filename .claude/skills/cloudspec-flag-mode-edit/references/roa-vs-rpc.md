# RPC vs ROA（RESTful）风格完整对照

## 一、总览

| 维度 | RPC | ROA（RESTful） |
|------|-----|----------------|
| apiStyle 声明位置 | service 级别 `@apiStyle("rpc")` | 每个操作的 `@http` 中 `apiStyle: "restful"` |
| 操作语义 | 通过 Action 名称区分 | 通过 HTTP 方法 + URI 路径区分 |
| URI 路径 | 无 | 必须有，含路径参数占位符 |
| 参数传递 | query / formData / system | path / query / body / header |
| HTTP 方法 | post, get | post, get, put, delete |
| 请求格式 | form-urlencoded | application/json |
| 响应格式 | 不声明 | application/json |
| schemes | `["http", "https"]` | `["https"]` |

## 二、如何识别 API 风格

1. 检查 `main.cspec` 是否有 `@apiStyle("rpc")` → RPC 风格
2. 检查操作文件 `@http` 中的 `apiStyle` 字段：`"RPC"` 或 `"restful"`
3. 检查 `@http` 中是否有 `uri` 字段 → 有则为 ROA

## 三、@http 注解对比

### RPC

```cspec
@http({
  apiStyle: "RPC"
  schemes: {
    online: ["http", "https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
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
  deprecated: false
})
```

**关键差异**：

| 字段 | RPC | ROA |
|------|-----|-----|
| `apiStyle` | `"RPC"` | `"restful"` |
| `uri` | 不存在 | **必须**，RESTful 路径 |
| `requestContentType` | 不存在 | `["application/json"]` |
| `responseContentType` | 不存在 | `["application/json"]` |
| `schemes` | 通常 `["http", "https"]` | 通常 `["https"]` |

### ROA HTTP 方法与操作类型对应

| 操作类型 | HTTP 方法 |
|---------|----------|
| Create | `["post"]` |
| Get | `["get"]` |
| Update | `["put"]` |
| Delete | `["delete"]` |
| List | `["get"]` |

## 四、@backendConfigurationHttp 对比

### RPC

```cspec
@backendConfigurationHttp({
  applicationName: "ProductName"
  retries: { online: -1 }
  timeout: { online: 10000 }
  backendUrl: {
    online: "http://vpc_online/api/resource/create#vpc"
  }
  sign: true
  signPolicy: "Local"
})
```

### ROA

```cspec
@backendConfigurationHttp({
  applicationName: "ProductName"
  requestType: "String"
  responseType: "String"
  retries: { online: -1 }
  timeout: { online: 60000 }
  backendUrl: {
    online: "http://pop.product.prod/#vpc"
  }
  sign: true
  signPolicy: "Local"
})
```

**关键差异**：ROA 需要 `requestType: "String"` 和 `responseType: "String"`。

## 五、Flag 模式下的风格差异

Flag 模式下，RPC 和 ROA 的主要差异体现在操作注解上，资源定义（flag 标记）是相同的。

### 操作注解差异

| 注解 | RPC | ROA |
|------|-----|-----|
| `@http` | 无 `uri`，`apiStyle: "RPC"` | 有 `uri`，`apiStyle: "restful"` |
| `@backendConfigurationHttp` | 无 `requestType`/`responseType` | 需要 `requestType`/`responseType` |
| `@gatewayOptions` | 基础字段 | 额外需要 `roaRequestBodyLog`、`roaResponseBodyLog`、`outputParamVersion` |

### 书写检查清单

**RPC 项目：**

- [ ] `@http` 中 `apiStyle` 为 `"RPC"`
- [ ] `@http` 中无 `uri` 字段
- [ ] `schemes` 包含 `["http", "https"]`

**ROA 项目：**

- [ ] `@http` 中 `apiStyle` 为 `"restful"`
- [ ] `@http` 中有 `uri` 字段
- [ ] `@http` 中有 `requestContentType` 和 `responseContentType`
- [ ] `@backendConfigurationHttp` 包含 `requestType` 和 `responseType`
- [ ] `schemes` 为 `["https"]`

# 操作映射（Operation Mapping）参考

## 概述

操作映射定义了**资源属性**与**操作（API）的入参/出参**之间的关系。非 Flag 模式的资源需要通过 `@operationMapping` + `@defineOperationMapping` 显式配置映射关系。

> **Flag 模式不需要手动配置映射**：使用 `@flagMode` 的资源通过 `#[C,U,R,D,L]` 标记自动推导映射，无需本文档中的配置。

## 基本语法

### 在资源中关联映射

在资源的生命周期操作关联处，使用 `@operationMapping(MappingStructName)` 指定映射结构体：

```cspec
resource Instance {
  identifyDefinition: {
    @readonly
    InstanceId: string
  }
  properties: {
    Name: string
    Status: string
  }
  create: CreateInstance
  get: @operationMapping(Instance_get_mapping) GetInstance
  update: @operationMapping(Instance_update_mapping) UpdateInstance
  delete: @operationMapping(Instance_delete_mapping) DeleteInstance
  list: @operationMapping(Instance_list_mapping) ListInstances
}
```

**说明**：
- `create` 操作通常不需要 `@operationMapping`，因为创建时入参和资源属性通常一一对应
- `get` 操作**最常需要**映射，用于配置出参路径、资源不存在判定等
- 如果操作的入参/出参名称与资源属性名称完全一致，可以省略映射

### 定义映射结构体

使用 `@defineOperationMapping` 注解定义映射结构体：

```cspec
@defineOperationMapping
struct Instance_get_mapping {
  resourceAttributeMappings: [{
    resourceProperty: "$.InstanceId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "InstanceId"
    responsePathType: "normal"
    responsePath: "$.InstanceId"
  }, {
    resourceProperty: "$.Name"
    responsePathType: "normal"
    responsePath: "$.InstanceName"
  }]
}
```

## @defineOperationMapping 完整字段

| 字段 | 类型 | 说明 |
|------|------|------|
| **resourceAttributeMappings** | [resourceAttributeMappingsStruct] | **核心字段**：资源属性与 API 入参/出参的映射关系 |
| **retryPolicies** | [retryPoliciesStruct] | 按错误码的重试策略 |
| **resourceNotExistCondition** | resourceNotExistConditionStruct | 资源不存在的判定条件（用于 Get/Delete） |
| **asyncPollingByProperty** | [asyncPollingByPropertyStruct] | 异步操作轮询配置（通过资源属性判断） |
| **asyncPollingByAPI** | asyncPollingByAPIStruct | 异步操作轮询配置（通过其他 API 判断） |
| rootMapping | rootMappingStruct | 出参根节点配置 |
| constInputParameters | [constInputParametersStruct] | 常量入参配置 |
| pagination | paginationStruct | 分页配置 |
| idempotencyToken | string | 幂等参数的入参路径 |
| errorCodePosition | string | 出参中错误码的路径 |
| resultJudgment | resultJudgmentStruct | API 调用结果判断条件 |
| dependencies | [dependenciesStruct] | 操作的属性依赖配置 |
| incremental | incrementalStruct | 增量操作配置（ADD/REMOVE/TRUNCATE） |
| repeatList | [repeatListStruct] | repeatList 类型入参的展示形式 |
| paramSerializers | [paramSerializersStruct] | JSON 序列化参数配置 |

## resourceAttributeMappings（属性映射）

这是最核心的配置，定义资源属性与 API 参数的对应关系。

### 完整字段

| 字段 | 类型 | 说明 |
|------|------|------|
| **resourceProperty** | string | 资源属性路径，如 `"$.Name"` |
| requestMappingType | string | 入参映射类型：`param`（映射到具体参数）/ `api`（仅关联到 API）/ `computed`（计算属性） |
| requestPathType | string | 入参路径类型：`normal` / `jsonPath` / `repeatList` / `jsonArray` / `commaSeparated` / `kvPairs` |
| requestPath | string | 入参参数路径 |
| requestIn | string | 入参位置：`formData`（RPC）/ `body`（ROA）/ `query` |
| requestRequired | boolean | 入参是否必填 |
| requestConstValue | string | 固定常量值 |
| requestDefaultValue | string | 入参默认值 |
| requestEnumValues | [string] | 入参枚举值 |
| requestTransform | string | 入参到资源属性的 UDF/JSONata 转换 |
| responsePathType | string | 出参路径类型：`normal` / `jsonPath` |
| responsePath | string | 出参参数路径 |
| responseMappingType | string | 出参映射类型：`param`（默认）/ `computed` |
| responseTransform | string | 出参到资源属性的 UDF/JSONata 转换 |
| targetValue | string | 当 requestMappingType 为 `api` 时，API 调用后属性的目标值 |
| targetValueType | string | 目标值类型：`null` / `empty` / `specificValue`（默认） |
| hasSystemDefault | boolean | 是否有系统默认值 |
| constType | boolean | 常量参数类型：`true`（动态资源标识）/ `false`（普通常量） |
| penetrate | boolean | 是否透传 |
| mappingType | string | 映射类型：`property`（属性映射，默认）/ `response`（API 出参映射）/ `request`（API 入参映射） |
| requestValueMappings | [valueMappingsItem] | 入参和资源属性的值映射 |
| responseValueMappings | [valueMappingsItem] | 出参和资源属性的值映射 |
| dependencies | [dependenciesStruct] | 属性依赖配置 |

### 映射路径格式

- **`resourceProperty`**：资源属性路径，以 `$.` 开头，如 `"$.Status"`、`"$.ComputeNodes[*].ImageId"`
- **`requestPath`**：入参路径，直接使用 API 参数名
- **`responsePath`**：出参路径，以 `$.` 开头

### RPC 风格映射示例

RPC 风格下，参数名为 **PascalCase**，入参通过 `query` 或 `formData` 传递。

以下是 NAS FileSystem 的 Create 映射（RPC 项目）：

```cspec
@defineOperationMapping
struct FileSystem_create_CreateFileSystem_mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.StorageType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestRequired: true
    requestPath: "StorageType"
    requestIn: "query"
  }
  , {
    resourceProperty: "$.ProtocolType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestRequired: true
    requestPath: "ProtocolType"
    requestIn: "query"
  }
  , {
    resourceProperty: "$.FileSystemId"
    responsePathType: "normal"
    responsePath: "$.FileSystemId"
  }]
  asyncPollingByProperty: [{
    ResourceProperty: "$.Status"
    DelayedTime: 5
    Interval: 30
    TargetValue: "Running"
    TargetValueType: "assertEqual"
    Times: 20
  }]
}
```

以下是 EHPC Queue 的 Get 映射（RPC 项目）：

```cspec
@defineOperationMapping
struct Queue_get_GetQueue_mapping_info {
  resourceAttributeMappings: [{
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
  }
  , {
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "QueueName"
    responsePathType: "normal"
    responsePath: "$.Queue.QueueName"
  }
  , {
    resourceProperty: "$.EnableScaleIn"
    responsePathType: "normal"
    responsePath: "$.Queue.EnableScaleIn"
  }
  , {
    resourceProperty: "$.Status"
    responsePathType: "normal"
    responsePath: "$.Queue.Status"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$.Queue"
  }
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["ClusterNotFound", "InvalidParams"]
    }]
  }
}
```

**RPC 要点**：
- `requestPath`：PascalCase 参数名（如 `"StorageType"`、`"ClusterId"`）
- `requestIn`：通常为 `"query"`（读操作）或 `"formData"`（写操作），可省略
- `responsePath`：以 `$.` 开头，PascalCase 路径（如 `"$.Queue.QueueName"`）

### ROA 风格映射示例

ROA 风格下，参数名通常为 **camelCase**，body 参数的 `requestPath` 需要加 **`body.` 前缀**。

以下是 OpenAPIExplorer ApiMcpServerCore 的 Get 映射（ROA 项目，`@apiStyle("ROA")`）：

```cspec
@defineOperationMapping
struct ApiMcpServerCore_GetApiMcpServerCore_Mapping {
  resourceAttributeMappings: [{
    resourceProperty: "$.Id"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "id"
  }
  , {
    resourceProperty: "$.OauthClientId"
    responsePath: "$.oauthClientId"
  }
  , {
    resourceProperty: "$.PublicAccessType"
    responsePath: "$.publicAccessType"
  }
  , {
    resourceProperty: "$.EnableAssumeRole"
    responsePath: "$.enableAssumeRole"
  }
  , {
    resourceProperty: "$.Urls.Sse"
    responsePath: "$.urls.sse"
  }
  , {
    resourceProperty: "$.CreateTime"
    responsePath: "$.createTime"
  }]
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["InvalidApiMcpServerCore.NotFound"]
    }]
  }
}
```

以下是 Create 映射，注意 body 参数的 `requestPath` 使用 **`body.` 前缀**：

```cspec
@defineOperationMapping
struct ApiMcpServerCore_CreateApiMcpServerCore_Mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.OauthClientId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.oauthClientId"
  }
  , {
    resourceProperty: "$.PublicAccessType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.publicAccessType"
  }
  , {
    resourceProperty: "$.EnableAssumeRole"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.enableAssumeRole"
  }
  , {
    resourceProperty: "$.Id"
    responsePath: "$.id"
  }]
}
```

**ROA 要点**：
- `requestPath`（body 参数）：**`body.` 前缀** + camelCase（如 `"body.oauthClientId"`）
- `requestPath`（path/query 参数）：直接 camelCase 参数名（如 `"id"`）
- `responsePath`：以 `$.` 开头，camelCase 路径（如 `"$.oauthClientId"`、`"$.urls.sse"`）
- path 参数需与操作 URI 中 `{placeholder}` 一致

### RPC 与 ROA 映射差异对照

| 维度 | RPC 风格 | ROA 风格 |
|------|---------|---------|
| **参数命名** | PascalCase（如 `FileSystemId`） | camelCase（如 `oauthClientId`） |
| **requestPath（body 参数）** | 直接参数名（如 `"Queue.QueueName"`） | **`body.` 前缀**（如 `"body.oauthClientId"`） |
| **requestPath（非 body 参数）** | 直接参数名（如 `"ClusterId"`） | 直接参数名（如 `"id"`） |
| **入参位置 requestIn** | `query`（读）/ `formData`（写） | `body`（写）/ `query`（读）/ `path`（URI 参数） |
| **responsePath 命名** | PascalCase（如 `$.Queue.QueueName`） | camelCase（如 `$.oauthClientId`） |
| **分页方式** | 通常 `pageNumber`（页码分页） | 通常 `nextToken`（游标分页） |

### 常见映射场景

#### 场景 1：属性名与 API 参数名不同

资源属性 `Name` 对应 API 出参 `InstanceName`：

```cspec
{
  resourceProperty: "$.Name"
  responsePathType: "normal"
  responsePath: "$.InstanceName"
}
```

#### 场景 2：仅出参映射（Get 操作的只读字段）

只读字段只需要配置出参映射，不需要入参：

```cspec
{
  resourceProperty: "$.Status"
  responsePathType: "normal"
  responsePath: "$.Status"
}
```

#### 场景 3：API 类型映射（无直接参数对应）

某些操作（如 StartInstance、StopInstance）没有直接对应的参数，但会影响资源属性值：

```cspec
{
  resourceProperty: "$.Status"
  requestMappingType: "api"
  targetValue: "Running"
}
```

#### 场景 4：嵌套对象/数组映射

API 出参嵌套在对象或数组中：

```cspec
{
  resourceProperty: "$.ComputeNodes[*].ImageId"
  responsePathType: "normal"
  responsePath: "$.Queue.ComputeNodes[*].ImageId"
}
```

#### 场景 5：jsonArray 类型入参

Delete/List 操作中，资源属性映射到数组类型入参：

```cspec
{
  resourceProperty: "$.QueueName"
  requestMappingType: "param"
  requestPathType: "jsonArray"
  requestPath: "QueueNames[*]"
}
```

## retryPolicies（重试策略）

配置操作遇到特定错误码时的重试行为。常用于 Delete/Update 操作，当资源处于中间状态时自动重试。

### 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| **Code** | string | 触发重试的错误码 |
| **Interval** | int64 | 重试间隔（秒） |
| **Times** | int64 | 最大重试次数 |

### 示例

```cspec
@defineOperationMapping
struct Instance_delete_mapping {
  retryPolicies: [{
    Code: "IncorrectInstanceStatus"
    Interval: 5
    Times: 20
  }, {
    Code: "DependencyViolation"
    Interval: 10
    Times: 10
  }]
}
```

**说明**：当 Delete 操作返回 `IncorrectInstanceStatus` 错误码时，每隔 5 秒重试，最多重试 20 次。

## resourceNotExistCondition（资源不存在判定）

配置如何判断资源已不存在。通常用于 Get 或 Delete 操作的映射中。

### 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| anyOf | [ResourceNotExistInfo] | 满足任意一个条件即认为不存在 |
| allOf | [ResourceNotExistInfo] | 需满足全部条件 |
| oneOf | [ResourceNotExistInfo] | 只能满足其中一个条件 |

### ResourceNotExistInfo 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| notExistCheckType | string | 检查方式：`checkErrorCode`（检查错误码）/ `checkProperty`（检查属性值） |
| resourceNotExistErrorCodes | [string] | 资源不存在的错误码列表（checkErrorCode 时使用） |
| notExistCheckProperty | string | 待检查的资源属性路径（checkProperty 时使用） |
| notExistCheckTargetValueType | string | 值比较方式：`assertEqual` / `assertNull` / `assertEmpty` |
| notExistCheckTargetValue | string | 目标值（assertEqual 时使用） |

### 示例：通过错误码判断

```cspec
@defineOperationMapping
struct Instance_get_mapping {
  resourceNotExistCondition: {
    anyOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["InvalidInstanceId.NotFound", "InstanceNotFound"]
    }]
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.InstanceId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "InstanceId"
    responsePathType: "normal"
    responsePath: "$.InstanceId"
  }]
}
```

### 示例：通过属性值判断（软删除场景）

某些资源删除后不会立即消失，而是将 Status 置为 "Deleted"：

```cspec
@defineOperationMapping
struct Record_get_mapping {
  resourceNotExistCondition: {
    anyOf: [
      {
        notExistCheckType: "checkErrorCode"
        resourceNotExistErrorCodes: ["InvalidResourceId.NotFound"]
      },
      {
        notExistCheckType: "checkProperty"
        notExistCheckProperty: "$.Status"
        notExistCheckTargetValueType: "assertEqual"
        notExistCheckTargetValue: "Deleted"
      }
    ]
  }
}
```

## asyncPollingByProperty（异步操作轮询）

配置异步操作的轮询策略，通过检查资源属性值判断操作是否完成。

### 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| **ResourceProperty** | string | 轮询检查的资源属性路径 |
| **TargetValue** | string | 成功时的目标值 |
| **FailedValues** | [string] | 失败状态值列表 |
| **Interval** | int64 | 轮询间隔（秒） |
| **Times** | int64 | 最大轮询次数 |
| DelayedTime | int64 | 首次轮询延迟时间（秒） |
| TargetValueType | string | 值比较方式：`assertEqual`（默认）/ `assertNull` / `assertEmpty` / `assertNotEmpty` |

### 示例

```cspec
@defineOperationMapping
struct Instance_create_mapping {
  asyncPollingByProperty: [{
    ResourceProperty: "$.Status"
    DelayedTime: 10
    FailedValues: ["CreateFailed", "Error"]
    Interval: 10
    TargetValue: "Running"
    Times: 60
  }]
}
```

**说明**：Create 操作后，等待 10 秒，然后每 10 秒检查一次 Status 属性。当 Status 变为 `Running` 时认为创建成功；当 Status 变为 `CreateFailed` 或 `Error` 时认为创建失败。最多轮询 60 次。

## constInputParameters（常量入参）

配置操作中非资源属性的常量入参。

### 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| requestMappingType | string | 映射类型，通常为 `"param"` |
| requestPathType | string | 路径类型，通常为 `"normal"` |
| **requestPath** | string | 入参路径 |
| **requestConstValue** | string | 常量值 |
| constType | boolean | `true`（动态资源标识）/ `false`（普通常量） |

### 示例

```cspec
@defineOperationMapping
struct WafInstance_create_mapping {
  constInputParameters: [{
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ProductCode"
    requestConstValue: "waf"
  }]
}
```

## pagination（分页配置）

配置 List 操作的分页方式。

### NextToken 分页（ROA 常用）

```cspec
@defineOperationMapping
struct Instance_list_mapping {
  pagination: {
    type: "nextToken"
    inputToken: "NextToken"
    outputToken: "NextToken"
    maxItems: "MaxResults"
    maxItemsDefault: 50
    totalCount: "TotalCount"
  }
}
```

### PageNumber 分页（RPC 常用）

```cspec
@defineOperationMapping
struct Instance_list_mapping {
  pagination: {
    type: "pageNumber"
    initialPageNumber: 1
    initialPageSize: 20
    inputPageNumber: "PageNumber"
    inputPageSize: "PageSize"
    recordTotal: "$.TotalCount"
    items: "$.Instances.Instance"
  }
}
```

### RPC 与 ROA 分页差异

| 维度 | RPC 风格 | ROA 风格 |
|------|---------|---------|
| **分页类型** | `pageNumber`（页码分页） | `nextToken`（游标分页） |
| **入参** | `PageNumber` + `PageSize` | `NextToken` + `MaxResults` |
| **出参** | `TotalCount` + `PageNumber` + `PageSize` | `NextToken` + `TotalCount` |

## rootMapping（根节点映射）

配置出参中获取数据的根节点路径：

```cspec
@defineOperationMapping
struct Instance_get_mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
}
```

## 文件组织结构

在实际项目中，mapping 结构体通常**单独放在 `resources/{ResourceName}Mappings/` 子目录下**，每个操作一个文件。

**RPC 项目**（扁平目录）：

```
resources/
├── Queue.cspec                          # 资源定义
├── QueueMappings/                       # 映射文件目录
│   ├── CreateQueueMappingInfo.cspec     # Create 映射
│   ├── GetQueueMappingInfo.cspec        # Get 映射
│   ├── UpdateQueueMappingInfo.cspec     # Update 映射
│   ├── DeleteQueueMappingInfo.cspec     # Delete 映射
│   └── ListQueueMappingInfo.cspec       # List 映射
```

**ROA 项目**（按操作类型分子目录）：

```
resources/
├── ApiMcpServerCore.cspec                                    # 资源定义
├── ApiMcpServerCoreMappings/                                 # 映射文件目录
│   ├── Create/
│   │   └── CreateApiMcpServerCoreMappingInfo.cspec           # Create 映射
│   ├── Get/
│   │   └── GetApiMcpServerCoreMappingInfo.cspec              # Get 映射
│   ├── Update/
│   │   └── UpdateApiMcpServerCoreMappingInfo.cspec           # Update 映射
│   ├── Delete/
│   │   └── DeleteApiMcpServerCoreMappingInfo.cspec           # Delete 映射
│   └── List/
│       └── ListApiMcpServerCoresMappingInfo.cspec            # List 映射
```

每个映射文件必须以 `$version: 1` 和 `namespace:` 开头，与资源文件保持相同的 namespace。

## 完整示例（基于实际项目）

### RPC 完整示例（EHPC Queue）

以下示例基于 EHPC 项目的 Queue 资源，展示 RPC 风格的完整 5 个操作映射。

### 资源定义（resources/Queue.cspec）

```cspec
resource Queue {
  identifyDefinition: {
    @readonly
    ClusterId: string
    @readonly
    QueueName: string
  }
  properties: {
    @readonly
    CreateTime: string
    @readonly
    RegionId: string
    EnableScaleOut: boolean
    EnableScaleIn: boolean
    MinCount: int32
    MaxCount: int32
    @readonly
    InitialCount: int32
    InterConnect: string
    VSwitchIds: array<string>
    ComputeNodes: array<{ InstanceType: string, ImageId: string, ... }>
    HostnamePrefix: string
    HostnameSuffix: string
    @readonly
    UpdateTime: string
  }
  create: @operationMapping(Queue_create_CreateQueue_mapping_info) CreateQueue
  update: @operationMapping(Queue_update_UpdateQueue_mapping_info) UpdateQueue
  delete: @operationMapping(Queue_delete_DeleteQueues_mapping_info) DeleteQueues
  get: @operationMapping(Queue_get_GetQueue_mapping_info) GetQueue
  list: @operationMapping(Queue_list_ListQueues_mapping_info) ListQueues
}
```

### Get 映射（resources/QueueMappings/GetQueueMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.Ehpc.EHPC.v20240730

@defineOperationMapping
struct Queue_get_GetQueue_mapping_info {
  resourceAttributeMappings: [{
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
  }
  , {
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "QueueName"
    responsePathType: "normal"
    responsePath: "$.Queue.QueueName"
  }
  , {
    resourceProperty: "$.EnableScaleIn"
    responsePathType: "normal"
    responsePath: "$.Queue.EnableScaleIn"
  }
  , {
    resourceProperty: "$.EnableScaleOut"
    responsePathType: "normal"
    responsePath: "$.Queue.EnableScaleOut"
  }
  , {
    resourceProperty: "$.ComputeNodes[*].InstanceType"
    responsePathType: "normal"
    responsePath: "$.Queue.ComputeNodes[*].InstanceType"
  }
  , {
    resourceProperty: "$.ComputeNodes[*].ImageId"
    responsePathType: "normal"
    responsePath: "$.Queue.ComputeNodes[*].ImageId"
  }
  , {
    resourceProperty: "$.CreateTime"
    responsePathType: "normal"
    responsePath: "$.Queue.CreateTime"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$.Queue"
  }
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["ClusterNotFound", "InvalidParams"]
    }]
  }
}
```

### Create 映射（resources/QueueMappings/CreateQueueMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.Ehpc.EHPC.v20240730

@defineOperationMapping
struct Queue_create_CreateQueue_mapping_info {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "Queue.QueueName"
    responsePathType: "normal"
    responsePath: "$.Name"
  }
  , {
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
  }
  , {
    resourceProperty: "$.EnableScaleIn"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "Queue.EnableScaleIn"
  }
  , {
    resourceProperty: "$.ComputeNodes[*].InstanceType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "Queue.ComputeNodes[*].InstanceType"
  }]
}
```

### Update 映射（resources/QueueMappings/UpdateQueueMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.Ehpc.EHPC.v20240730

@defineOperationMapping
struct Queue_update_UpdateQueue_mapping_info {
  resourceAttributeMappings: [{
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
  }
  , {
    resourceProperty: "$.EnableScaleIn"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "Queue.EnableScaleIn"
  }
  , {
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "Queue.QueueName"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
}
```

### Delete 映射（resources/QueueMappings/DeleteQueueMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.Ehpc.EHPC.v20240730

@defineOperationMapping
struct Queue_delete_DeleteQueues_mapping_info {
  resourceAttributeMappings: [{
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
  }
  , {
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "jsonArray"
    requestPath: "QueueNames[*]"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
}
```

### List 映射（resources/QueueMappings/ListQueueMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.Ehpc.EHPC.v20240730

@defineOperationMapping
struct Queue_list_ListQueues_mapping_info {
  resourceAttributeMappings: [{
    resourceProperty: "$.ClusterId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "ClusterId"
    responseMappingType: "normal"
    responsePathType: "normal"
    responsePath: "$.ClusterId"
  }
  , {
    resourceProperty: "$.QueueName"
    requestMappingType: "param"
    requestPathType: "jsonArray"
    requestPath: "QueueNames[*]"
    responsePathType: "normal"
    responsePath: "$.Queues[*].QueueName"
  }
  , {
    resourceProperty: "$.EnableScaleIn"
    responsePathType: "normal"
    responsePath: "$.Queues[*].EnableScaleIn"
  }
  , {
    resourceProperty: "$.CreateTime"
    responsePathType: "normal"
    responsePath: "$.Queues[*].CreateTime"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$.Queues[*]"
  }
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["ClusterNotFound"]
    }]
  }
}
```

### ROA 完整示例（OpenAPIExplorer ApiMcpServerCore）

以下示例基于 OpenAPIExplorer 项目（`@apiStyle("ROA")`）的 ApiMcpServerCore 资源，展示 ROA 风格的完整 5 个操作映射。

#### 资源定义（resources/ApiMcpServerCore.cspec）

```cspec
resource ApiMcpServerCore {
  identifyDefinition: {
    @readonly
    Id: string
  }
  properties: {
    OauthClientId: string
    @readonly
    CreateTime: string
    VpcWhitelists: array<string>
    EnableAssumeRole: boolean
    AssumeRoleOverridePolicy: string
    AssumeRoleName: string
    PublicAccessType: string
    EnableCustomVpcWhitelist: boolean
    @readonly
    Urls: { Sse: string, Mcp: string, VpcSse: string, VpcMcp: string }
    @readonly
    RequiredRamPolicy: string
  }
  create: @operationMapping(ApiMcpServerCore_CreateApiMcpServerCore_Mapping) CreateApiMcpServerCore
  update: @operationMapping(ApiMcpServerCore_UpdateApiMcpServerCore_Mapping) UpdateApiMcpServerCore
  delete: @operationMapping(ApiMcpServerCore_DeleteApiMcpServerCore_Mapping) DeleteApiMcpServerCore
  get: @operationMapping(ApiMcpServerCore_GetApiMcpServerCore_Mapping) GetApiMcpServerCore
  list: @operationMapping(ApiMcpServerCore_ListApiMcpServerCore_Mapping) ListApiMcpServerCores
}
```

#### Get 映射（resources/ApiMcpServerCoreMappings/Get/GetApiMcpServerCoreMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.OpenAPIExplorer.OpenAPIExplorer.v20241130

@defineOperationMapping
struct ApiMcpServerCore_GetApiMcpServerCore_Mapping {
  resourceAttributeMappings: [{
    resourceProperty: "$.Id"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "id"
  }
  , {
    resourceProperty: "$.OauthClientId"
    responsePath: "$.oauthClientId"
  }
  , {
    resourceProperty: "$.AssumeRoleOverridePolicy"
    responsePath: "$.assumeRoleOverridePolicy"
  }
  , {
    resourceProperty: "$.PublicAccessType"
    responsePath: "$.publicAccessType"
  }
  , {
    resourceProperty: "$.EnableAssumeRole"
    responsePath: "$.enableAssumeRole"
  }
  , {
    resourceProperty: "$.EnableCustomVpcWhitelist"
    responsePath: "$.enableCustomVpcWhitelist"
  }
  , {
    resourceProperty: "$.AssumeRoleName"
    responsePath: "$.assumeRoleName"
  }
  , {
    resourceProperty: "$.VpcWhitelists[*]"
    responsePath: "$.vpcWhitelists[*]"
  }
  , {
    resourceProperty: "$.Urls.Sse"
    responsePath: "$.urls.sse"
  }
  , {
    resourceProperty: "$.Urls.Mcp"
    responsePath: "$.urls.mcp"
  }
  , {
    resourceProperty: "$.Urls.VpcSse"
    responsePath: "$.urls.vpcSse"
  }
  , {
    resourceProperty: "$.Urls.VpcMcp"
    responsePath: "$.urls.vpcMcp"
  }
  , {
    resourceProperty: "$.CreateTime"
    responsePath: "$.createTime"
  }
  , {
    resourceProperty: "$.RequiredRamPolicy"
    responsePath: "$.requiredRamPolicy"
  }
  , {
    resourceProperty: "$.ResourceGroupId"
    responsePath: "$.resourceGroupId"
  }]
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["InvalidApiMcpServerCore.NotFound"]
    }]
  }
}
```

#### Create 映射（resources/ApiMcpServerCoreMappings/Create/CreateApiMcpServerCoreMappingInfo.cspec）

> 注意 body 参数的 `requestPath` 使用 **`body.` 前缀**。

```cspec
$version: 1
namespace: alicloud.OpenAPIExplorer.OpenAPIExplorer.v20241130

@defineOperationMapping
struct ApiMcpServerCore_CreateApiMcpServerCore_Mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.OauthClientId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.oauthClientId"
  }
  , {
    resourceProperty: "$.PublicAccessType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.publicAccessType"
  }
  , {
    resourceProperty: "$.VpcWhitelists[*]"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.vpcWhitelists[*]"
  }
  , {
    resourceProperty: "$.EnableAssumeRole"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.enableAssumeRole"
  }
  , {
    resourceProperty: "$.EnableCustomVpcWhitelist"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.enableCustomVpcWhitelist"
  }
  , {
    resourceProperty: "$.AssumeRoleOverridePolicy"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.assumeRoleOverridePolicy"
  }
  , {
    resourceProperty: "$.AssumeRoleName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.assumeRoleName"
  }
  , {
    resourceProperty: "$.Id"
    responsePath: "$.id"
  }]
}
```

#### Update 映射（resources/ApiMcpServerCoreMappings/Update/UpdateApiMcpServerCoreMappingInfo.cspec）

> path 参数（如 `id`）直接用参数名，body 参数用 `body.` 前缀。

```cspec
$version: 1
namespace: alicloud.OpenAPIExplorer.OpenAPIExplorer.v20241130

@defineOperationMapping
struct ApiMcpServerCore_UpdateApiMcpServerCore_Mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.Id"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "id"
  }
  , {
    resourceProperty: "$.OauthClientId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.oauthClientId"
  }
  , {
    resourceProperty: "$.PublicAccessType"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.publicAccessType"
  }
  , {
    resourceProperty: "$.VpcWhitelists[*]"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.vpcWhitelists[*]"
  }
  , {
    resourceProperty: "$.EnableAssumeRole"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.enableAssumeRole"
  }
  , {
    resourceProperty: "$.EnableCustomVpcWhitelist"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.enableCustomVpcWhitelist"
  }
  , {
    resourceProperty: "$.AssumeRoleOverridePolicy"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.assumeRoleOverridePolicy"
  }
  , {
    resourceProperty: "$.AssumeRoleName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.assumeRoleName"
  }]
}
```

#### Delete 映射（resources/ApiMcpServerCoreMappings/Delete/DeleteApiMcpServerCoreMappingInfo.cspec）

```cspec
$version: 1
namespace: alicloud.OpenAPIExplorer.OpenAPIExplorer.v20241130

@defineOperationMapping
struct ApiMcpServerCore_DeleteApiMcpServerCore_Mapping {
  resourceAttributeMappings: [{
    resourceProperty: "$.Id"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "id"
  }]
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["InvalidApiMcpServerCore.NotFound"]
    }]
  }
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
}
```

#### List 映射（resources/ApiMcpServerCoreMappings/List/ListApiMcpServerCoresMappingInfo.cspec）

> ROA 项目通常使用 `nextToken` 分页。

```cspec
$version: 1
namespace: alicloud.OpenAPIExplorer.OpenAPIExplorer.v20241130

@defineOperationMapping
struct ApiMcpServerCore_ListApiMcpServerCore_Mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$.apiMcpServerCores[*]"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.CreateTime"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "createTime"
    responsePath: "$.apiMcpServerCores[*].createTime"
  }
  , {
    resourceProperty: "$.Id"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "id"
    responsePath: "$.apiMcpServerCores[*].id"
  }]
  pagination: {
    type: "nextToken"
    inputToken: "nextToken"
    outputToken: "$.nextToken"
    maxItems: "$.maxResults"
    maxItemsDefault: 20
    totalCount: "$.totalCount"
  }
}
```

## 注意事项

- **映射结构体命名**：推荐格式为 `{ResourceName}_{type}_{OperationName}_mapping_info`（如 `Queue_get_GetQueue_mapping_info`），也可简写为 `{ResourceName}_{type}_mapping`
- **映射文件组织**：映射结构体应放在 `resources/{ResourceName}Mappings/` 子目录下，每个操作一个文件，文件名如 `GetQueueMappingInfo.cspec`
- **Get 操作映射最重要**：Get 操作的映射决定了资源属性如何从 API 出参中提取，必须确保所有资源属性都有正确的出参映射
- **避免 Get 操作多余入参映射**：Get 操作通常只需要主键作为入参，多余的入参映射可能导致查询失败（参见 FAQ：12-invalid-get-mapping-error）
- **路径格式**：`responsePath` 以 `$.` 开头；`requestPath` 直接用参数名，但 **ROA 的 body 参数需加 `body.` 前缀**（如 `"body.oauthClientId"`），RPC 无此前缀
- **属性名映射**：当资源属性名与 API 参数名不同时，必须通过 `requestPath`/`responsePath` 显式映射
- **异步操作**：创建/删除等耗时操作需要配置 `asyncPollingByProperty`，否则测试会因超时失败
- **重试策略**：Delete/Update 操作建议配置 `retryPolicies`，处理资源处于中间状态的场景
- **rootMapping**：当 API 出参的数据嵌套在某个节点下时（如 `$.Queue`），使用 `rootMapping` 指定根节点，简化各属性的 `responsePath`

# Flag 模式使用指南（含完整项目示例）

## 1. 什么是 Flag 模式

Flag 模式用于从资源定义自动推导 CURDL（Create/Update/Get/List/Delete）能力。
核心思想是：在资源属性上使用 `#[...]` 标记该属性参与哪些操作，系统据此生成或补全对应操作的参数与映射。

最常见入口是给资源加上：

```cspec
@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Item {
  // ...
}
```

## 2. Flag 语义速查

### 2.1 基础 Flag 含义

| Flag | 全称 | 含义 |
|------|------|------|
| C | CREATE | 字段参与 Create 操作入参 |
| U | UPDATE | 字段参与 Update 操作入参 |
| R | READ | 字段在 Get 操作输出中出现 |
| D | DELETE | 字段用于 Delete 操作定位 |
| L | LIST_OUTPUT | 字段在 List 操作输出中出现 |
| S | SEARCH | 字段可作为检索条件 |
| F | FILTER_SUPPORT | 字段可作为 List 过滤入参 |
| A | CREATE_AFTER_OUTPUT | 创建后额外返回的字段 |

### 2.2 关键 Flag 区别

**F（FILTER_SUPPORT）vs L（LIST_OUTPUT）**
- **F**：控制字段是否可作为 List 操作的过滤入参
- **L**：控制字段是否在 List 操作的输出中出现
- **区别**：一个字段可以同时拥有 F 和 L，也可以只拥有其中一个
  - 可过滤但不在列表输出：`#[C,F]`
  - 在列表输出但不可过滤：`#[R,L]`

**A（CREATE_AFTER_OUTPUT）**
- 用于标记在 Create 操作后额外返回的字段
- 通常用于服务端生成的字段，如创建时间、实例 ID 等
- 示例：`#[C,R,A]` 表示创建时作为入参，同时也在创建后返回

## 3. 默认 Flag 行为

根据资源分类（classification），系统会为字段应用不同的默认 Flag：

### 3.1 资源标识字段（identifyDefinition）

- **normal_id**：`#[!C,!U,D,R,F,S]`
  - 不参与创建和更新
  - 用于删除、查询、列表输出
  - 可作为过滤和检索条件
- **normal_props**：`#[C,U,!D,R,F,!S]`
  - 参与创建和更新
  - 不用于删除定位
  - 在查询和列表中输出
  - 可作为过滤条件，不可作为检索条件

### 3.2 属性字段（properties）

默认行为会根据字段的注解（如 `@readonly`、`@required`）自动调整：
- `@readonly` 字段：默认不含 C、U
- `@required` 字段：创建时必填
- 其他字段：根据资源分类应用默认 Flag

## 4. 最小可用示例

下面示例展示一个最小的 Flag 模式资源（含主键与属性标记）：

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Item {
  identifyDefinition: {
    @readonly
    ItemId: string #[R,D,L]
  }
  properties: {
    @required
    Name: string #[C,U,R,L]
    Description: string #[C,U,R,L]
    @readonly
    Status: string #[R,L]
    @required
    Type: string #[C,R,L]
    Priority: int32 #[C,U,R]
  }
}
```

### 这个资源会带来什么

- 资源具备 Flag 模式自动推导能力。
- `ItemId` 可用于删除定位（`D`），并在查询/列表可见（`R/L`）。
- `Name`、`Description` 支持创建和更新，同时可在读取中返回。
- `Status` 只读，不参与创建/更新。

## 3. 完整项目结构示例（推荐）

实际工程建议拆成 `main + resources + operations`。

### 3.1 目录结构

```text
flag-mode-demo/
├── main.cspec
├── resources/
│   └── Item.cspec
└── operations/
    ├── CreateItem.cspec
    ├── GetItem.cspec
    ├── ListItems.cspec
    ├── UpdateItem.cspec
    └── DeleteItem.cspec
```

### 3.2 main.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

import "./resources/*"
import "./operations/*"

@runtimeType("pop")
service TestProduct {
  version: "2024-01-01"
}
```

### 3.3 resources/Item.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@arn({
  standard: "acs:testproduct:${Region}:${AccountId}:item/${ItemId}"
})
@document({
  zh: "测试项目"
  en: "Item"
  name: "测试项目"
  nameEn: "Item"
})
@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Item {
  identifyDefinition: {
    @readonly
    ItemId: string #[R,D,L]
  }
  properties: {
    @readonly
    RegionId: string
    @required
    Name: string #[C,U,R,L]
    Description: string #[C,U,R,L]
    @readonly
    Status: string #[R,L]
    @readonly
    CreateTime: string #[R]
    @required
    Type: string #[C,R,L]
    Priority: int32 #[C,U,R]
  }
}
```

### 4.4 operations/CreateItem.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/create#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "创建测试项目"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: {
    online: ["https"]
  }
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
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@visibility("Public")
operation CreateItem {
  errors: [Error_CreateItem]
}

error Error_CreateItem {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

### 4.5 operations/GetItem.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/get#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询测试项目详情"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: {
    online: ["https"]
  }
  methods: ["get", "post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "read"
  typeFromOperation: "get"
  riskType: "none"
  chargeType: "free"
})
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@visibility("Public")
operation GetItem {
  errors: [Error_GetItem, Error_GetItem_NotFound]
}

error Error_GetItem {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_GetItem_NotFound {
  httpCode: 404
  errorCode: "Item.NotFound"
  errorMessage: "The specified item does not exist."
  type: "user"
  default: false
}
```

### 4.6 operations/ListItems.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/list#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询测试项目列表"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: {
    online: ["https"]
  }
  methods: ["get", "post"]
  authenticators: ["AK"]
  deprecated: false
})
@numberPaginated({
  initialPageNumber: 1
  initialPageSize: 20
  inputPageNumber: "PageNumber"
  inputPageSize: "PageSize"
  recordTotal: "$.TotalCount"
})
@operationInfo({
  operationTypeOld: "read"
  typeFromOperation: "list"
  riskType: "none"
  chargeType: "free"
})
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@rootMapping({
  responsePath: "$.Items[*]"
})
@visibility("Public")
operation ListItems {
  errors: [Error_ListItems]
}

error Error_ListItems {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

### 4.7 operations/UpdateItem.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/update#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "修改测试项目"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: {
    online: ["https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "update"
  riskType: "none"
  chargeType: "free"
})
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@visibility("Public")
operation UpdateItem {
  errors: [Error_UpdateItem, Error_UpdateItem_NotFound]
}

error Error_UpdateItem {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_UpdateItem_NotFound {
  httpCode: 404
  errorCode: "Item.NotFound"
  errorMessage: "The specified item does not exist."
  type: "user"
  default: false
}
```

### 4.8 operations/DeleteItem.cspec

```cspec
$version: 1
namespace: alicloud.TestProduct.TestProduct.v20240101

@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/delete#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "删除测试项目"
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: {
    online: ["https"]
  }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "delete"
  riskType: "none"
  chargeType: "free"
})
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@visibility("Public")
operation DeleteItem {
  errors: [Error_DeleteItem, Error_DeleteItem_NotFound]
}

error Error_DeleteItem {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_DeleteItem_NotFound {
  httpCode: 404
  errorCode: "Item.NotFound"
  errorMessage: "The specified item does not exist."
  type: "user"
  default: false
}
```

## 5. 操作层覆盖（@for）

Flag 模式支持保留自动能力，同时对具体操作做显式覆盖。

### 5.1 当前限制：不支持自定义 Input/Output

以下写法当前会触发校验错误：

```cspec
@for(Item)
operation CreateItem {
  input: CreateItemInput
}
```

典型报错语义：

- Flag 模式不允许用户自定义 API 的 input/output
- Operation 定义了 input/output，请移除该定义

### 5.2 可覆盖内容：操作注解

可以覆盖后端配置、鉴权、错误映射等注解：

```cspec
@for(Item)
@backendConfigurationHttp({
  applicationName: "TestProduct"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/item/create#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@visibility("Public")
operation CreateItem {
  errors: [Error_CreateItem]
}
```

使用 `@for(Resource)` 对单个操作覆盖时，其它未覆盖操作仍保持 Flag 模式默认推导。

## 6. 字段设计建议

### 6.1 资源标识字段

```cspec
@readonly
ItemId: string #[R,D,L]
```

- 删除通常依赖 ID，因此要有 `D`
- 查询与列表常常返回该字段，因此保留 `R/L`
- 创建时通常由服务端分配，不放 `C/U`

### 6.2 可修改业务字段

```cspec
Name: string #[C,U,R,L]
```

适合名称、描述、备注等标准业务字段。

### 6.3 只读状态字段

```cspec
@readonly
Status: string #[R,L]
```

由后端计算或状态机驱动，不建议暴露到创建/更新入参。

### 6.4 不可变字段（创建后不可修改）

```cspec
@required
Type: string #[C,R,L]
```

创建时必填，之后不可修改，查询/列表可见。

### 6.5 检索条件字段

```cspec
Keyword: string #[C,S]
```

如果字段仅用于筛选，不希望进入资源回包，可通过 Flag 组合精细控制。

### 6.6 服务端生成字段

```cspec
@readonly
@clientProhibited
UpdateTime: string #[!C]
```

使用否定标记 `#[!C]` 显式排除 Create 操作。

### 6.7 可过滤但不在列表输出的字段

```cspec
Keyword: string #[C,F]
```

- 字段可作为 List 操作的过滤条件（F）
- 但不在 List 操作的输出中（没有 L）
- 适用于内部筛选参数，不需要在列表中展示

### 6.8 在列表输出但不可过滤的字段

```cspec
Status: string #[R,L]
```

- 字段在 List 操作的输出中出现（L）
- 但不可作为过滤条件（没有 F）
- 适用于只读的统计字段、状态字段等

### 6.9 创建后额外返回的字段

```cspec
CreateTime: string #[C,R,A]
```

- 创建时作为入参（C）
- 在 Get 操作中返回（R）
- 在 Create 操作后额外返回（A）
- 适用于服务端生成但需要立即返回给客户端的字段

### 6.10 操作私有属性

当需要在操作中暴露资源中未定义的私有属性时，使用 `@rac` 注解：

```cspec
@rac({
  operatePrivateType: ["PrivateField1", "PrivateField2"]
})
operation CreateItem {
  errors: [Error_CreateItem]
}
```

- `operatePrivateType`：声明该操作需要使用的私有属性列表
- 这些属性不会出现在资源定义中，但在操作中可用

## 7. 可直接改造的模板

```cspec
$version: 1
namespace: alicloud.Demo.Product.v20260407

@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource DemoResource {
  identifyDefinition: {
    @readonly
    DemoResourceId: string #[R,D,L]
  }
  properties: {
    @readonly
    RegionId: string
    @required
    Name: string #[C,U,R,L]
    Description: string #[C,U,R,L]
    @readonly
    Status: string #[R,L]
    @readonly
    CreateTime: string #[R,L]
  }
}
```

把 `DemoResource`、字段名以及 `@resourceBaseInfo` 中业务元数据替换为真实产品值即可开始使用。

## 8. 推荐落地流程

1. 先完成资源与属性 Flag 标记（资源层）。
2. 运行 `aliyun cspec build` 验证 CURDL 推导结果。
3. 只对确实需要差异化的操作做 `@for` 覆盖。
4. 对覆盖操作补齐 `@backendConfigurationHttp`、`@http`、`@errorMapping`、`@ram` 等生产注解。
5. 运行 `aliyun cspec check --name <ResourceName>` 进行规范检查。

## 9. 常见问题

### Q1：`@readonly` 和 `#[C,U]` 冲突怎么办？

优先统一语义，避免"既只读又可写"的定义。通常推荐：

- 只读字段：`@readonly + #[R,L]`
- 可写字段：去掉 `@readonly`，按需给 `C/U`

### Q2：我只想自定义某一个操作，是否会破坏其他自动生成能力？

不会。使用 `@for(Resource)` 对单个操作覆盖时，其它未覆盖操作仍可保持 Flag 模式默认推导。

### Q3：额外入参（非资源属性）怎么加？

当前 Flag 模式流程不支持通过 operation 自定义 input struct 注入额外入参。
建议把字段建模在资源属性层，通过 `#[...]` + 注解组合控制参与的操作。

### Q4：F 和 L 有什么区别？

- **F（FILTER_SUPPORT）**：控制字段是否可作为 List 操作的过滤入参
- **L（LIST_OUTPUT）**：控制字段是否在 List 操作的输出中出现
- **使用场景**：
  - 可过滤但不在列表输出：`#[C,F]` - 适用于内部筛选参数
  - 在列表输出但不可过滤：`#[R,L]` - 适用于只读的统计字段
  - 既可过滤又在列表输出：`#[R,F,L]` - 适用于普通业务字段

### Q5：A flag 什么时候使用？

A（CREATE_AFTER_OUTPUT）用于标记在 Create 操作后额外返回的字段，通常用于：
- 服务端生成的字段，如创建时间、实例 ID 等
- 需要在创建成功后立即返回给客户端的字段
- 示例：`#[C,R,A]` 表示创建时作为入参，同时也在创建后返回

### Q6：Update 操作中 @required 字段是否必填？

**不是**。Update 操作采用部分更新语义：
- 只有在请求中明确提供的字段才会被更新
- 即使字段标记了 `@required`，在 Update 中也非必填
- `@required` 仅对 Create 操作生效
- ID 属性用于定位资源，通常不需要在 Update 中提供

### Q7：如何标记操作私有属性？

使用 `@rac` 注解在操作上声明：

```cspec
@rac({
  operatePrivateType: ["PrivateField1", "PrivateField2"]
})
operation CreateItem {
  errors: [Error_CreateItem]
}
```

这些属性不会出现在资源定义中，但在操作中可用。

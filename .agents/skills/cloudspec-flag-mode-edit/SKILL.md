---
name: cloudspec-flag-mode-edit
description: |
  CloudSpec Flag 模式编辑 skill。创建/编辑 @flagMode 资源，使用 #[C,U,R,D,F,L,S,A] 标记属性，创建 @for(Resource) 操作，编辑 flag 组合，转换为 flag 模式。
  Triggers: "flag模式", "flagMode", "@flagMode", "flag标记", "#[C,U,R,D,L]", "@for", "flag模式资源", "flag模式操作", "创建flag模式资源", "编辑flag标记", "flag模式转换", "convert to flag mode".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec Flag 模式 IDL 元数据编辑

给定一个 CloudSpec 项目和 Flag 模式相关的编辑需求，**必须严格按照以下步骤执行**。

> **前置**：若尚未了解整体编写流程与 skill 调用顺序，请先调用 **cloudspec-idl-guide**。

## Step1 了解 Flag 模式基础

Flag 模式（`@flagMode`）是 CloudSpec 中一种声明式的属性映射模式。启用后，资源属性通过 `#[C,U,R,D,L,S]` 标记声明自己参与哪些 CRUD 操作，系统会自动为操作生成对应的 Input/Output 结构体，无需手动定义。

### 核心概念

- **Flag 标记**：在资源属性类型后使用 `#[...]` 声明该属性参与哪些操作
- **自动推导**：系统根据 flag 标记自动生成 Create/Update/Get/Delete/List 操作的 input/output
- **操作关联**：Flag 模式下的操作使用 `@for(ResourceName)` 关联资源，操作体内只需声明 `errors`
- **当前限制**：Flag 模式下**不支持**在 operation 中自定义 `input`/`output`

#### Flag 完整列表

| Flag | 全称 | 说明 |
|------|------|------|
| C | CREATE_SUPPORT | 创建操作支持 |
| U | UPDATE_SUPPORT | 更新操作支持 |
| R | READ_SUPPORT | 查询操作支持（Get） |
| D | DELETE_SUPPORT | 删除操作支持 |
| F | FILTER_SUPPORT | List 过滤入参支持（可筛选字段） |
| L | LIST_OUTPUT_SUPPORT | List 输出支持（列表返回字段） |
| S | SEARCH_SUPPORT | 搜索操作支持 |
| A | CREATE_AFTER_OUTPUT | 创建后额外返回字段 |

**重要区分**：
- **F (FILTER_SUPPORT)**：标记字段可作为 List 操作的过滤条件（入参），如 `Status`、`Name` 等筛选条件
- **L (LIST_OUTPUT_SUPPORT)**：标记字段会在 List 操作的输出中返回，如 `Name`、`Status` 等列表展示字段
- **A (CREATE_AFTER_OUTPUT)**：标记字段在 Create 操作成功后会额外返回（通常用于服务端生成的字段）

### 识别 API 风格（前置必做）

CloudSpec 项目分为 **RPC** 和 **ROA（RESTful）** 两种风格，两者在注解配置、参数传递等方面有根本性差异。**在编辑前必须先识别风格**。

**优先使用 CLI 命令获取项目信息**：

```bash
aliyun cspec baseinfo
```

该命令输出 JSON 格式的项目基本信息，包含 `apiStyle`（rpc/roa）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`（资源列表）、`operations`（操作列表）等字段。从 `apiStyle` 字段即可判断 API 风格。

仅当 CLI 不可用时，再手动检查 `main.cspec` 或操作文件，详见 [roa-vs-rpc.md](references/roa-vs-rpc.md)。

### 相关文档

- Flag 模式使用指南（含完整项目示例）：[flag-mode-usage-guide.md](references/flag-mode-usage-guide.md)
- Flag 语法参考（flag 定义、映射规则、选择指南）：[flag-mode-reference.md](references/flag-mode-reference.md)
- RPC/ROA 风格对照：[roa-vs-rpc.md](references/roa-vs-rpc.md)

> **需要更详细的注解说明？** 以下完整文档位于共享知识库 `cloudspec-shared-knowledge`：
> - Resource 注解完整规范：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/resource-annotate.md`
> - 传值约束注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/assignment-constraint-annotate.md`
> - 值约束注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/value-constraint-annotate.md`
> - 文档注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/document-annotate.md`

## Step2 确认任务类型

**必须**先明确用户的需求属于以下哪种任务：

| 任务类型 | 说明 |
|---------|------|
| **新建 Flag 模式项目** | 从零搭建一个使用 @flagMode 的完整 CloudSpec 项目（main + resources + operations） |
| **新建 Flag 模式资源** | 在已有项目中创建一个使用 @flagMode 的资源定义 |
| **新建 Flag 模式操作** | 为 Flag 模式资源创建关联操作（使用 @for 注解） |
| **编辑资源属性 Flag 标记** | 修改已有 Flag 模式资源中属性的 `#[C,U,R,D,L,S]` 标记 |
| **添加/修改/删除资源属性** | 在 Flag 模式资源中添加、修改或删除属性（含 flag 标记） |
| **编辑资源级注解** | 修改 Flag 模式资源的 @arn、@document、@resourceBaseInfo 等注解 |
| **编辑操作注解** | 修改 Flag 模式操作的 @backendConfigurationHttp、@http、@errorMapping 等注解 |
| **转换为 Flag 模式** | 将已有的非 Flag 模式资源转换为 Flag 模式 |

如果用户的需求不够明确，**必须**询问以下信息：

- 目标资源名称
- 具体要修改的字段或注解
- 新字段的类型、是否必填、是否只读
- 新字段应参与哪些操作（C/U/R/D/L/S）

## Step3 执行编辑

### Step3.1 阅读参考文档

根据任务类型，阅读相应的参考文档：

- **所有任务**都需要先阅读 [flag-mode-reference.md](references/flag-mode-reference.md) 了解 flag 语法和映射规则
- **新建项目或资源**时需要阅读 [flag-mode-usage-guide.md](references/flag-mode-usage-guide.md) 了解完整项目结构和示例
- **涉及操作编辑**时需要阅读 [flag-mode-usage-guide.md](references/flag-mode-usage-guide.md) 中的操作覆盖部分

### Step3.2 理解当前项目上下文

在执行编辑前，**必须**：

1. 阅读当前 CloudSpec 项目的 `main.cspec`，获取 namespace 信息。
2. 确认项目是否已使用 `@flagMode`——查看 `./resources` 目录下的资源定义。
3. 如果是编辑已有资源，阅读目标资源的 `.cspec` 文件，完整理解当前资源的结构、注解和 flag 配置。
4. 查看 `./operations` 目录下该资源的操作定义，确认操作是否使用 `@for(ResourceName)` 关联。
5. 参考同项目中其他资源和操作的配置风格，确保一致性。

### Step3.3 新建 Flag 模式项目

创建完整的 Flag 模式项目需要以下文件结构：

```text
project/
├── main.cspec
├── resources/
│   └── {ResourceName}.cspec
└── operations/
    ├── Create{ResourceName}.cspec
    ├── Get{ResourceName}.cspec
    ├── List{ResourceName}s.cspec
    ├── Update{ResourceName}.cspec
    └── Delete{ResourceName}.cspec
```

#### main.cspec 模板

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

import "./resources/*"
import "./operations/*"

@runtimeType("pop")
service {PopCode} {
  version: "{YYYY-MM-DD}"
}
```

#### 资源文件模板

详见 Step3.4。

#### 操作文件模板

详见 Step3.6。

### Step3.4 新建 Flag 模式资源

在 `./resources/` 目录下创建 `{ResourceName}.cspec` 文件。

#### 资源注解顺序

```
@arn(...)                   // 1. ARN
@references([...])          // 2. 引用关系（可选）
@document({...})            // 3. 文档描述
@flagMode                   // 4. Flag 模式标记（必须）
@resourceBaseInfo({...})    // 5. 资源基本信息
@terraform({...})           // 6. 集成配置（可选）
@config({...})              // 7. 集成配置（可选）
@ros({...})                 // 8. 集成配置（可选）
@notDestroy({...})          // 9. 不销毁标记（可选）
```

#### 资源定义模板

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@arn("acs:{service}:${Region}:${AccountId}:{resourceType}/${ResourceId}")
@document({
  zh: "资源中文描述"
  en: "Resource English description"
  name: "资源中文名"
  nameEn: "Resource English Name"
})
@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource {ResourceName} {
  identifyDefinition: {
    @readonly
    {ResourceName}Id: string #[R,D,L]
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
  }
}
```

#### Flag 标记选择规则

根据字段语义选择正确的 flag 组合（详见 [flag-mode-reference.md](references/flag-mode-reference.md)）：

| 字段类型 | 推荐 Flag | 注解 | 说明 |
|---------|----------|------|------|
| **资源主键**（ID） | `#[R,D,L]` | `@readonly` | 服务端生成，用于删除定位和查询返回 |
| **用户可编辑字段** | `#[C,U,R,L]` | — | 创建/更新可传，查询/列表可见 |
| **不可变字段**（创建后不可改） | `#[C,R,L]` | — | 仅创建时传入 |
| **只读系统字段** | `#[R,L]` | `@readonly` | Status、CreateTime 等，L 表示 LIST_OUTPUT_SUPPORT（列表输出支持） |
| **仅详情返回字段** | `#[R]` | — | 列表不返回的详细信息 |
| **可编辑但列表不返回** | `#[C,U,R]` | — | Priority 等 |
| **仅创建时传入** | `#[C]` | — | 一次性初始化参数 |
| **检索条件字段** | `#[C,S]` | — | 仅用于筛选 |
| **可过滤字段** | `#[C,F]` | — | List 操作的过滤入参（FILTER_SUPPORT），如 Status、Name 等筛选条件 |
| **创建后额外返回字段** | `#[C,R,A]` | — | Create 操作成功后额外返回的字段（CREATE_AFTER_OUTPUT） |
| **RegionId** | 无 flag | `@readonly` | 平台特殊处理，唯一不需要 flag 的字段 |

#### 字段注解

Flag 模式下的字段注解与普通模式一致，注解紧贴字段，注解之间不能有空行：

```cspec
@required
@document({
  name: "名称"
  zh: "资源名称"
  en: "Resource name"
})
@length({ min: 2, max: 128 })
Name: string #[C,U,R,L]
```

**注意**：`@readonly` 和 `#[C,U]` 语义冲突，避免"既只读又可写"的定义。

### Step3.5 编辑资源属性 Flag 标记

修改已有 Flag 模式资源的 flag 标记时：

1. **理解每个 flag 的含义**——参考 [flag-mode-reference.md](references/flag-mode-reference.md) 中的 Flag 定义表。
2. **确保 flag 与字段语义匹配**——必填字段应有 `C`，只读字段不应有 `C` 和 `U`。
3. **只修改用户指定的部分**，保留所有未涉及的注解和结构。
4. **否定标记**：使用 `#[!X]` 显式排除某个操作（如 `#[!C]` 排除 Create），仅用于服务端生成字段（如 UpdateTime）。
5. **F 和 L 的区别**：
   - **F (FILTER_SUPPORT)**：标记字段可作为 List 操作的**入参过滤条件**，如 `Status`、`Name` 等，用户可以通过这些字段筛选列表结果
   - **L (LIST_OUTPUT_SUPPORT)**：标记字段会在 List 操作的**输出中返回**，如 `Name`、`Status` 等，用于列表展示
   - **常见组合**：可筛选且列表返回的字段通常使用 `#[C,F,R,L]`（如 `Status`、`Name`）

### Step3.6 新建 Flag 模式操作

Flag 模式下的操作使用 `@for(ResourceName)` 关联资源，操作体内**只需声明 `errors`**，不能自定义 `input`/`output`。

#### 操作注解顺序

```
@for(ResourceName)                  // 1. 关联资源（必须）
@backendConfigurationHttp({...})    // 2. 后端服务配置
@document({...})                    // 3. 文档描述
@errorMapping({...})                // 4. 错误映射
@http({...})                        // 5. HTTP 配置
@numberPaginated({...})             // 6. 分页配置（仅 List 操作）
@operationInfo({...})               // 7. 操作信息
@ram({...})                         // 8. RAM 权限配置
@rootMapping({...})                 // 9. 响应数据映射（仅 List 操作）
@visibility("Public")              // 10. 可见性
```

#### RPC/ROA 风格差异

**`@backendConfigurationHttp` 差异**：

| 字段 | RPC 风格 | ROA 风格 |
|------|---------|---------|
| `requestType` | 不需要 | 需要（如 `"Object"`） |
| `responseType` | 不需要 | 需要（如 `"Object"`） |
| `httpMethod` | 不需要 | 需要（如 `"post"`、`"get"`） |
| `consume` | 不需要 | 需要（如 `"application/json"`） |
| `sign` / `signPolicy` | 需要 | 不需要 |

ROA 风格的 `@backendConfigurationHttp` 示例：

```cspec
@backendConfigurationHttp({
  applicationName: "{Product}"
  requestType: "Object"
  responseType: "Object"
  consume: "application/json"
  httpMethod: "post"
  retries: {
    online: -1
  }
  timeout: {
    online: 10000
  }
  backendUrl: {
    online: "http://example.aliyun-inc.com/api/{resource}/create"
  }
})
```

> **重要**：以上字段值（`requestType`、`responseType`、`httpMethod`、`consume`）必须参考同项目中已有操作的实际配置，不要随意填写。

**入参位置差异**：

**RPC 风格**：
- Create 操作：入参放在 `formData`（form 表单）
- List 操作：入参放在 `query`（URL 查询参数）
- Update 操作：入参放在 `formData`（form 表单）
- Get/Delete 操作：入参放在 `query`（URL 查询参数）

**ROA 风格**：
- 所有操作：入参根据 HTTP 方法决定，GET/DELETE 用 `query`，POST/PUT 用 `formData` 或 `body`

**操作私有属性**：
- 标记 `@rac(operatePrivateType)` 的属性始终放在 `query` 参数中，不受上述规则影响

#### Update 操作部分更新语义

Update 操作支持**部分更新**：
- 非主键（ID）属性在 Update 入参中**一律非必填**
- 只传入需要更新的字段，未传入的字段保持原值
- 主键字段（如 `{ResourceName}Id`）用于定位资源，不参与更新

#### 标准操作模板

**Create 操作**：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@for({ResourceName})
@backendConfigurationHttp({
  applicationName: "{Product}"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/{resource}/create#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "创建{资源中文名}"
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
operation Create{ResourceName} {
  errors: [Error_Create{ResourceName}]
}

error Error_Create{ResourceName} {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

**Get 操作**：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@for({ResourceName})
@backendConfigurationHttp({
  applicationName: "{Product}"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/{resource}/get#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询{资源中文名}详情"
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
operation Get{ResourceName} {
  errors: [Error_Get{ResourceName}, Error_Get{ResourceName}_NotFound]
}

error Error_Get{ResourceName} {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_Get{ResourceName}_NotFound {
  httpCode: 404
  errorCode: "{ResourceName}.NotFound"
  errorMessage: "The specified {resource} does not exist."
  type: "user"
  default: false
}
```

**List 操作**（注意额外的 `@numberPaginated` 和 `@rootMapping`）：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@for({ResourceName})
@backendConfigurationHttp({
  applicationName: "{Product}"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/{resource}/list#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "查询{资源中文名}列表"
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
  responsePath: "$.{ResourceName}s[*]"
})
@visibility("Public")
operation List{ResourceName}s {
  errors: [Error_List{ResourceName}s]
}

error Error_List{ResourceName}s {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

**Update 操作**：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@for({ResourceName})
@backendConfigurationHttp({
  applicationName: "{Product}"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/{resource}/update#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "修改{资源中文名}"
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
operation Update{ResourceName} {
  errors: [Error_Update{ResourceName}, Error_Update{ResourceName}_NotFound]
}

error Error_Update{ResourceName} {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_Update{ResourceName}_NotFound {
  httpCode: 404
  errorCode: "{ResourceName}.NotFound"
  errorMessage: "The specified {resource} does not exist."
  type: "user"
  default: false
}
```

**Delete 操作**：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{YYYYMMDD}

@for({ResourceName})
@backendConfigurationHttp({
  applicationName: "{Product}"
  retries: {
    online: -1
  }
  timeout: {
    online: 5000
  }
  backendUrl: {
    online: "http://vpc_online/api/{resource}/delete#vpc"
  }
  sign: true
  signPolicy: "Local"
})
@document({
  name: "删除{资源中文名}"
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
  riskType: "high"
  chargeType: "free"
})
@ram({
  enable: true
  level: "operate"
  atGateway: false
})
@visibility("Public")
operation Delete{ResourceName} {
  errors: [Error_Delete{ResourceName}, Error_Delete{ResourceName}_NotFound]
}

error Error_Delete{ResourceName} {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_Delete{ResourceName}_NotFound {
  httpCode: 404
  errorCode: "{ResourceName}.NotFound"
  errorMessage: "The specified {resource} does not exist."
  type: "user"
  default: false
}
```

### Step3.7 转换为 Flag 模式

将已有非 Flag 模式资源转换为 Flag 模式时：

1. **在资源上添加 `@flagMode` 注解**（放在 `@document` 之后、`@resourceBaseInfo` 之前）。
2. **为所有属性添加 flag 标记**（`RegionId` 除外），根据字段语义选择正确的 flag 组合。
3. **修改操作定义**：
   - 为每个操作添加 `@for(ResourceName)` 注解。
   - 移除操作体中的 `input` 和 `output` 声明，只保留 `errors`。
   - 删除不再需要的 input/output struct 定义。
4. **验证转换结果**：运行 `aliyun cspec build` 确保编译通过。

#### 转换前后对比

**转换前**（非 Flag 模式）：

```cspec
resource Item {
  identifyDefinition: {
    @readonly
    ItemId: string
  }
  properties: {
    @required
    Name: string
    @readonly
    Status: string
  }
  create: CreateItem
  get: GetItem
}

operation CreateItem {
  input: CreateItemInput
  output: CreateItemOutput
  errors: [Error_CreateItem]
}

struct CreateItemInput {
  @required
  Name: string
}

struct CreateItemOutput {
  ItemId: string
}
```

**转换后**（Flag 模式）：

```cspec
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
    @readonly
    Status: string #[R,L]
  }
}

@for(Item)
operation CreateItem {
  errors: [Error_CreateItem]
}
```

## Step4 验证

编辑完成后，执行以下验证：

### 4.1 语法检查

在 CloudSpec 项目根目录下运行：

```bash
aliyun cspec build
```

确保编译通过。如果编译失败，根据错误信息修复后重新编译。

编译通过后，对涉及的资源运行规范检查：

> **⚠️ Inner API 豁免**：运行 `aliyun cspec baseinfo`，若 `isInnerApi` 为 `true`，则该项目为 inner API，**不需要运行 `aliyun cspec check`**，只需确保 `aliyun cspec build` 通过即可。

```bash
aliyun cspec check --name <ResourceName>
```

### 4.2 规范检查清单

**通用检查：**

- [ ] 文件以 `$version: 1` 和 `namespace:` 开头
- [ ] namespace 与 `main.cspec` 保持一致
- [ ] 资源有 `@arn`、`@document`、`@flagMode`、`@resourceBaseInfo`
- [ ] 主键字段在 `identifyDefinition` 中，且标记 `@readonly`
- [ ] 系统字段（`CreateTime`、`Status` 等）标记 `@readonly`
- [ ] 注解与目标之间没有空行
- [ ] 注解使用冒号语法（`key: value`），不使用等号
- [ ] PascalCase 命名（资源名、字段名）
- [ ] 2 空格缩进

**Flag 模式专项检查：**

- [ ] 资源有 `@flagMode` 注解
- [ ] 所有属性（`RegionId` 除外）都有 flag 标记 `#[...]`
- [ ] flag 标记与字段语义匹配（必填字段有 `C`、只读字段无 `C` 和 `U`）
- [ ] `@readonly` 和 `#[C,U]` 不冲突
- [ ] 主键 flag 为 `#[R,D,L]`
- [ ] F 和 L flag 不混淆：F 用于 List 入参过滤，L 用于 List 输出返回
- [ ] 可筛选字段（如 Status、Name）应同时有 F 和 L flag：`#[C,F,R,L]`

**Flag 模式操作检查：**

- [ ] 每个操作有 `@for(ResourceName)` 注解
- [ ] 操作体内只有 `errors`，没有 `input`/`output`
- [ ] 每个操作至少关联一个 error 定义
- [ ] List 操作有 `@numberPaginated` 和 `@rootMapping`
- [ ] `@document` 中 `name` 字段使用中文："创建/查询/修改/删除 + 资源中文名"
- [ ] 操作命名遵循 `动词+资源名` 格式，List 操作使用复数
- [ ] 与同项目其他操作的配置风格保持一致

### 4.3 常见编译错误

| 错误信息 | 原因 | 修复方式 |
|---------|------|---------|
| Flag 模式不允许自定义 input/output | 操作定义了 input/output | 移除 operation 中的 input/output 声明 |
| 属性缺少 flag 标记 | @flagMode 资源的属性没有 #[...] | 为属性添加合适的 flag 标记 |
| @readonly 与 C/U flag 冲突 | 只读字段不应参与创建/更新 | 移除 C/U flag 或移除 @readonly |

## 注意事项

- **Flag 模式下操作不能自定义 input/output**，input/output 由系统根据资源属性的 flag 标记自动生成。
- **`RegionId` 是唯一不需要 flag 标记的特殊字段**，标记 `@readonly`，无 flag。
- **flag 标记放在类型（或 `>` 闭合符号）之后**，如 `Name: string #[C,U,R,L]`。
- **额外入参（非资源属性）不能通过 operation 自定义 input struct 注入**，建议把字段建模在资源属性层，通过 `#[...]` + 注解组合控制参与的操作。
- **使用 `@for(Resource)` 对单个操作覆盖时，其它未覆盖操作仍保持 Flag 模式默认推导**。
- 编辑资源时，**只修改用户指定的部分**，保留所有未涉及的注解、字段和结构。
- 严禁修改 API 元数据（`./operations` 目录下的文件），除非用户明确要求。
- 注解值使用冒号语法（`@document({ zh: "描述" })`），不使用等号语法。
- 复杂类型字段（array、struct 嵌套）的 flag 标记放在类型闭合符号 `>` 之后。
- **Update 操作部分更新语义**：`@required` 注解仅表示创建时必填，Update 入参中非主键（ID）属性一律非必填，支持部分更新。
- **操作私有属性**：标记 `@rac(operatePrivateType)` 的属性始终放在 `query` 参数中，不受 RPC/ROA 风格参数位置规则影响。
- **Delete 操作 riskType**：Delete 操作的 `riskType` 应设置为 `"high"`，表示高风险操作。
- **默认 visibility**：代码自动生成时，操作的默认 `visibility` 为 `"Private"`，如需公开需显式设置为 `"Public"`。
- **用户覆盖合并机制**：用户可以通过 `@for(Resource)` 和自定义注解覆盖系统默认推导，未覆盖部分保持默认行为。

---
name: cloudspec-operation-edit
description: |
  CloudSpec 操作（Operation）创建与编辑 skill。创建/编辑 operations/ 下的 .cspec 文件，定义 input/output/error，编辑 @http/@operationInfo/@backendConfigurationHttp 注解，管理 RAM 权限（@defineAction/@requiredPermission），支持 curl 转 cspec、重命名操作、RPC/ROA 风格。
  Triggers: "新建操作", "创建API", "添加参数", "修改操作", "operation", "curl转cspec", "create operation", "add parameter", "edit annotation", "编辑操作", "cspec操作", "操作定义", "@http", "@operationInfo", "CRUD操作", "flagMode操作", "RAM权限", "重命名操作", "rename operation", "curl示例", "根据curl", "新建API接口", "创建操作文件", "顺手做一次规范检查".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec Operation 编辑

给定一个CloudSpec项目和操作（Operation）编辑需求，**必须严格按照以下步骤执行**。

> **前置**：若尚未了解整体编写流程与 skill 调用顺序，请先调用 **cloudspec-idl-guide**。

## Step1 了解CloudSpec Operation基础

一个完整的CloudSpec项目分为3部分：资源（resources）、操作（operations）、资源测试（tests）。
CloudSpec IDL语法简介、快速开始查看[quick-start.md](references/quick-start.md)

### 识别API风格（前置必做）

CloudSpec项目分为**RPC**和**ROA（RESTful）**两种风格，两者在注解配置、参数传递、HTTP方法等方面有根本性差异。**在编辑操作前必须先识别风格**。

**优先使用 CLI 命令获取项目信息**：

```bash
aliyun cspec baseinfo
```

该命令输出 JSON 格式的项目基本信息，包含 `apiStyle`（rpc/roa）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`（资源列表）、`operations`（操作列表）等字段。从 `apiStyle` 字段即可判断 API 风格。

仅当 CLI 不可用时，再手动检查文件：
1. 检查`main.cspec`是否有`@apiStyle("rpc")` → RPC风格
2. 检查操作文件`@http`中的`apiStyle`字段：`"RPC"`或`"restful"`

两种风格的详细差异参考[roa-vs-rpc.md](references/roa-vs-rpc.md)。

### ROA 高频易错配置

ROA 操作编辑时，body 参数和后端配置最容易出错。新增或修改前先找同项目 ROA 操作作为模板，再按以下形态校对。

请求体字段必须在 input 字段上声明 `@parameter({ in: "body" })`；path 参数要同时出现在 `@http.uri` 和字段的 `@backendName` 中：

```cspec
@http({
  apiStyle: "restful"
  uri: "/2021-04-06/functions/{functionName}/sessions"
  methods: ["post"]
  requestContentType: "application/json"
  responseContentType: "application/json"
})
operation CreateSession {
  input: CreateSessionInput
  output: CreateSessionOutput
  errors: [Error_CreateSession]
}

struct CreateSessionInput {
  @backendName("functionName")
  @parameter({ in: "path" })
  @required
  FunctionName: string

  @backendName("sessionTTLInSeconds")
  @parameter({ in: "body" })
  SessionTTLInSeconds: int64
}
```

ROA 的 `@backendConfigurationHttp` 通常需要 `requestType: "String"` 和 `responseType: "String"`，并保留同项目已有的 `signKeyName`、`signPolicy`、`backendUrl` 风格：

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

### Operation 概述

Operation（操作）定义了API的具体行为，每个Operation对应一个可执行的API接口。

- `./operations`目录存放所有的操作元数据，每个`.cspec`文件对应一个完整的Operation定义。
- Operation由三部分组成：输入参数（input）、输出结果（output）、错误处理（errors）。

### 相关文档

- Operation设计指南：[operation-design-guide.md](references/operation-design-guide.md)
- Operation相关注解详解：[operation-annotate.md](references/operation-annotate.md)
- RPC/ROA风格对照：[roa-vs-rpc.md](references/roa-vs-rpc.md)

> **需要更详细的注解说明？** 以下完整文档位于共享知识库 `cloudspec-shared-knowledge`：
> - Operation 注解完整规范：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/operation-annotate.md`
> - IAM/RAM 权限注解：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/iam-annotate.md`
> - 企业级能力注解（@terraform/@ros/@config 等）：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/enterprise-annotate.md`
> - 事件类注解（@eventConfig/@events 等）：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/event-annotate.md`
> - Operation 设计指南完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/how-to-design-operation.md`

## Step2 确认任务类型

**必须**先明确用户的需求属于以下哪种任务：

| 任务类型 | 说明 |
|---------|------|
| **新建标准操作** | 创建CRUD中的某一个标准操作（Create/Get/Update/Delete/List） |
| **新建非标准操作** | 创建与资源生命周期无关的自定义操作（如AttachDisk、RebootInstance等） |
| **编辑操作注解** | 修改已有操作的注解配置（如@http、@backendConfigurationHttp、@errorMapping等） |
| **编辑输入输出参数** | 修改操作的input/output结构体中的字段 |
| **编辑错误定义** | 添加、修改或删除操作关联的error定义 |
| **编辑RAM权限** | 修改操作关联的@defineAction结构体 |
| **删除字段或注解** | 从操作定义中移除字段或注解 |

如果用户的需求不够明确，**必须**询问以下信息：

- 目标操作名称（新建时需确认操作关联的资源名称）
- 具体要修改的内容

## Step3 执行编辑

### Step3.1 阅读参考文档

根据任务类型，阅读相应的参考文档：

- **所有任务**都需要先阅读[operation-design-guide.md](references/operation-design-guide.md)了解Operation的结构和配置规范
- **涉及注解**时需要阅读[operation-annotate.md](references/operation-annotate.md)了解注解的详细用法

### Step3.2 理解当前项目上下文

在执行编辑前，**必须**：

1. **运行 `aliyun cspec baseinfo` 获取项目基本信息**（namespace、apiStyle、资源列表、操作列表等）。如 CLI 不可用，则阅读 `main.cspec` 获取 namespace 信息（格式为`alicloud.{Product}.{PopCode}.v{PopVersion}`）。
2. 如果是编辑已有操作，先阅读目标操作的`.cspec`文件，理解当前配置。
3. **（新建操作时必做）选择一个同项目中与目标操作类型最接近的已有操作文件，完整阅读其内容作为模板**。例如：新建 List 操作就找一个已有的 List 操作，新建 Create 操作就找一个已有的 Create 操作。**必须完整复制该模板的所有注解结构**（包括 `@backendConfigurationHttp`、`@returnMode`、`@rootMapping`、`@errorMapping`、`@ram`、`@requiredPermission` 等），然后再修改具体字段值。这是确保注解完整性的关键步骤，**严禁跳过**。
4. 如果操作关联资源，阅读`./resources`目录下对应资源的`.cspec`文件，理解资源的属性定义。

### Step3.3 新建操作

如果是新建操作，**优先使用CLI脚手架**：

```bash
aliyun cspec create operation -m <OperationName> -v http
```

此命令会根据当前项目的API风格（RPC/ROA）自动生成操作模板文件，包含正确的注解结构。生成后再根据需求补充具体配置。

`-v`支持的后端类型：`http`、`hsf`、`dubbo`、`http_hsf`。

如果CLI不可用或需要手动创建，按以下规则生成：

#### 文件头

每个`.cspec`文件必须以版本和命名空间开头：

```cspec
$version: 1
namespace: alicloud.{Product}.{Service}.v{YYYYMMDD}
```

namespace需与当前项目的`main.cspec`中的namespace保持一致。

#### 注解顺序

操作的注解按以下顺序组织（按需选用，非全部必填）：

```
@for(ResourceName)                  // 仅flagMode资源的操作需要
@actionTrail({...})                 // 审计日志（可选）
@backendConfigurationHttp({...})    // 后端服务配置
@controlPolicy({...})              // 流控策略（可选）
@document({...})                    // 文档描述
@errorMapping({...})                // 错误映射
@gatewayOptions({...})             // 网关配置（可选）
@http({...})                        // HTTP配置
@numberPaginated({...})             // 分页配置（仅List操作）
@operationInfo({...})               // 操作信息
@ram({...})                         // RAM权限配置
@requiredPermission(StructName)     // 关联RAM Action
@returnMode({...})                 // 返回模式（可选）
@rootMapping({...})                 // 响应数据映射（仅List操作）
@visibility("Public")              // 可见性
```

#### 必填注解

每个Operation必须包含以下注解，具体配置说明参考[operation-annotate.md](references/operation-annotate.md)：

| 注解 | 说明 |
|-----|------|
| `@backendConfigurationHttp({...})` | 后端路由配置 |
| `@document({...})` | 操作文档，至少包含`name`字段 |
| `@errorMapping({...})` | 错误字段映射 |
| `@http({...})` | HTTP协议配置 |
| `@operationInfo({...})` | 操作类型元数据 |
| `@ram({...})` | RAM鉴权配置 |
| `@requiredPermission(StructName)` | 关联RAM Action结构体 |
| `@visibility("Public")` | 可见性 |

#### 标准操作的operationType对应关系

`@operationInfo`中有多个类型字段，以项目实际使用的为准：

| 操作类型 | operationType | typeFromOperation | operationTypeOld | HTTP方法 | 文档name前缀 |
|---------|--------------|-------------------|------------------|----------|------------|
| Create | `"Write"` 或 `"create"` | `"create"` | `"write"` | `["post"]` | "创建" |
| Get | `"Read"` 或 `"get"` | `"get"` | `"read"` | `["get", "post"]` | "查询...详情" |
| Update | `"Write"` 或 `"update"` | `"update"` | `"write"` | `["post"]` | "修改" |
| Delete | `"Write"` 或 `"delete"` | `"delete"` | `"write"` | `["post"]` | "删除" |
| List | `"List"` 或 `"list"` | `"list"` | `"read"` | `["get", "post"]` | "查询...列表" |

> `@operationInfo`还可能包含`riskType`、`chargeType`、`abilityTreeNodes`、`tenantRelevance`等字段，请参考同项目中已有操作的配置。

#### 操作体

操作体需要声明input、output和errors：

```cspec
operation CreateItem {
  input: CreateItemInput
  output: CreateItemOutput
  errors: [Error_CreateItem]
}
```

> **flagMode特殊说明**：如果资源使用了`@flagMode`，操作需要添加`@for(ResourceName)`注解，且操作体内只需声明`errors`（input/output由系统根据资源属性的flag标记自动生成）：
>
> ```cspec
> @for(Item)
> ...
> operation CreateItem {
>   errors: [Error_CreateItem]
> }
> ```

#### 错误定义

每个操作需要关联error定义：

```cspec
error Error_CreateItem_ServiceUnavailable {
  httpCode: 503
  errorCode: "ServiceUnavailable"
  backendErrorCode: "ServiceUnavailable"
  errorMessage: "The request has failed due to a temporary failure of the server."
  type: "user"
  default: false
}
```

> **关于`default`字段**：部分项目使用`default: true`标记兜底错误，部分项目所有error均为`default: false`。请参考同项目中已有操作的写法保持一致。

#### RAM Action结构体

每个操作需要关联一个`@defineAction`结构体，注意`resource`是保留字，需要用反引号包裹：

```cspec
@defineAction
struct CreateItem_RAM_Action_0 {
  action: "{service}:CreateItem"
  resources: [{
    `resource`: "acs:{service}:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

#### 文件位置

新建的操作文件放在`./operations/`目录下，命名为`{OperationName}.cspec`。

### Step3.4 编辑已有操作

如果是编辑已有操作，遵循以下原则：

1. **只修改用户指定的部分**，保留所有未涉及的注解和结构。
2. 修改注解时，保持注解顺序不变（参考Step3.3中的注解顺序）。
3. 添加新字段到input/output结构体时，需要添加正确的`@parameter`注解。
4. 删除字段时，如果该字段被其他地方引用，需要同步清理引用。
5. 修改`@backendConfigurationHttp`时，参考同项目其他操作的配置保持一致。

### Step3.5 编辑输入输出参数

操作的输入输出参数定义在struct中，编辑参数时需要注意：

#### 添加输入参数

```cspec
struct CreateItemInput {
  @required
  @parameter({
    in: 'query'
    name: 'itemName'
    required: true
  })
  ItemName: string
}
```

参数的`in`值选项（按API风格使用）：

| 值 | 适用风格 | 说明 |
|---|---------|------|
| `query` | RPC + ROA | URL查询参数 |
| `formData` | RPC | POST表单数据 |
| `system` | RPC | 系统参数（如Action），不对外暴露 |
| `path` | ROA | URI路径参数，必须在`@http.uri`中有`{name}`占位符 |
| `body` | ROA | JSON请求体字段 |
| `header` | ROA | HTTP请求头 |

#### 添加输出参数

```cspec
struct CreateItemOutput {
  @parameter({
    responseName: 'itemId'
  })
  ItemId: string
}
```

#### 字段约束注解

| 注解 | 用途 | 示例 |
|-----|------|------|
| `@required` | 必填字段 | `@required` |
| `@readonly` | 只读字段 | `@readonly` |
| `@default(value)` | 默认值 | `@default("")` |
| `@length({min, max})` | 字符串长度约束 | `@length({ min: 2, max: 128 })` |
| `@range({min, max})` | 数值范围约束 | `@range({ min: 1, max: 100 })` |
| `@regexPattern("pattern")` | 正则校验 | `@regexPattern("^[a-zA-Z].*$")` |
| `@document({...})` | 文档描述 | `@document({ zh: "名称", en: "Name" })` |
| `@sensitive` | 敏感字段 | `@sensitive` |
| `@deprecated` | 废弃字段 | `@deprecated` |
| `@clientProhibited` | 禁止客户端传值 | `@clientProhibited` |
| `@conflictsWith([...])` | 互斥字段 | `@conflictsWith(['Id'])` |
| `@idempotencyToken` | 幂等参数 | `@idempotencyToken` |

#### ⚠️ 常见类型与语法陷阱

**类型名陷阱**：

| 错误写法 | 正确写法 | 说明 |
|----------|----------|------|
| `MaxResults: integer` | `MaxResults: int32` | CloudSpec 没有 `integer` 类型，整数用 `int32` 或 `int64` |
| `Count: long` | `Count: int64` | 没有 `long` 类型，用 `int64` |
| `Score: number` | `Score: float` | 没有 `number` 类型，用 `float` 或 `double` |
| `list Items { member: Item }` | `Items: array<{ Name: string }>` 或独立定义 `array Items { item: Item }` | 顶层类型用 `array` 关键字，不用 `list`（`list` 不是合法关键字） |

**注解陷阱 — 以下注解不是 CloudSpec 内置注解**，使用会报 `自定义 annotate xxx 未定义`：

| 错误写法 | 错误信息 | 正确做法 |
|----------|----------|----------|
| `@backend({...})` | `自定义 annotate backend 未定义` | 用 `@backendConfigurationHttp({...})`（操作级后端配置的正确注解） |
| `@idempotent` | `自定义 annotate idempotent 未定义` | 移除（CloudSpec 没有此注解） |
| `@pattern("regex")` | `自定义 annotate pattern 未定义` | 移除，如需正则校验用 `@regexPattern("regex")` |
| `@maxLength(n)` | `自定义 annotate maxLength 未定义` | 移除，如需长度约束用 `@length({ min: 0, max: n })` |
| `@min(n)` / `@max(n)` | `自定义 annotate min/max 未定义` | 移除，如需范围约束用 `@range({ min: n, max: m })` |

**语法陷阱**：

| 错误写法 | 正确写法 | 说明 |
|----------|----------|------|
| `@enums([{value:"A"}{value:"B"}])` | `@enums([{value:"A"},{value:"B"}])` | 注解数组中对象项之间**必须有逗号** |

> 操作的列表类型输出推荐用 `array<{...}>` 内联语法（同文件内定义，结构清晰），参考同项目已有操作的写法。

## Step4 验证

编辑完成后，执行以下验证：

### 4.1 更新索引与语法检查

**新建或重命名操作后，必须先更新导入索引**：

```bash
aliyun cspec fix index
```

然后在CloudSpec项目根目录下运行编译：

```bash
aliyun cspec build
```

确保编译通过。如果编译失败，根据错误信息修复后重新编译。

编译通过后，对涉及的组件运行规范检查：

> **⚠️ Inner API 豁免**：运行 `aliyun cspec baseinfo`，若 `isInnerApi` 为 `true`，则该项目为 inner API，**不需要运行 `aliyun cspec check`**，只需确保 `aliyun cspec build` 通过即可。

```bash
aliyun cspec check --name <OperationName>
```

### 4.2 规范检查清单

- [ ] 文件以`$version: 1`和`namespace:`开头
- [ ] 注解与目标之间没有空行
- [ ] 注解使用冒号语法（`key: value`），不使用等号
- [ ] 所有操作都有`@http`、`@document`、`@operationInfo`、`@ram`、`@requiredPermission`、`@visibility`
- [ ] flagMode资源的操作有`@for(ResourceName)`注解
- [ ] 每个操作至少关联一个error定义
- [ ] RAM Action结构体使用`@defineAction`，且`resource`键使用反引号包裹
- [ ] PascalCase命名（操作名、字段名、结构体名）
- [ ] 2空格缩进
- [ ] operation、error、struct块之间有一个空行
- [ ] 操作命名遵循`动词+资源名`格式（如`CreateInstance`、`ListInstances`）
- [ ] List操作名称使用复数形式
- [ ] 与同项目其他操作的配置风格保持一致

## 注意事项

- 编辑操作时，**必须**参考同项目中其他操作的配置方式，确保`@backendConfigurationHttp`、`@http`、`@errorMapping`等注解的配置风格一致。
- `@document`中的`name`字段使用中文，标准操作遵循："创建/查询/修改/删除 + 资源中文名"的命名方式，不使用"新建/获取/更新"。
- 严禁修改资源元数据（`./resources`目录下的文件），除非用户明确要求。
- 注解值使用冒号语法（`@document({ zh: "描述" })`），不使用等号语法（~~`@document(zh = "描述")`~~）。
- `resource`是CloudSpec IDL保留字，在RAM Action结构体中必须使用反引号包裹（`` `resource` ``）。
- 非标准操作（非CRUD）的`@operationInfo`中，`typeFromOperation`通常设为与操作行为最接近的类型，如果都不匹配则根据实际情况参考同项目中的写法。
- **RPC风格**：`@http`中无`uri`，参数使用`in: "query"`或`in: "formData"`，schemes通常为`["http", "https"]`。
- **ROA风格**：`@http`中必须有`uri`和`requestContentType`/`responseContentType`，参数使用`in: "path"`/`in: "body"`/`in: "query"`，`@backendConfigurationHttp`需要`requestType`和`responseType`，`@gatewayOptions`需要ROA专属字段。详见[roa-vs-rpc.md](references/roa-vs-rpc.md)。

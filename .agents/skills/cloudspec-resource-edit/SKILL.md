---
name: cloudspec-resource-edit
description: |
  CloudSpec 资源（Resource）定义编辑 skill。添加/修改/删除 resources/ 下的资源属性，编辑资源级注解（@arn, @document, @resourceBaseInfo, @flagMode），调整字段注解（@required, @readonly, @length, @enums），修改 flag 标记 #[C,U,R,D,L]，配置 @operationMapping/@defineOperationMapping。
  Triggers: "添加属性", "添加字段", "修改资源", "删除属性", "resource", "资源编辑", "flag标记", "operationMapping", "add property", "modify resource", "delete field", "edit resource", "cspec资源", "资源定义", "@arn", "@flagMode", "#[C,U,R,D,L]", "identifyDefinition", "主键", "子资源", "加一个字段", "新增属性", "修改资源注解", "文档描述改成".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec 资源编辑

给定一个CloudSpec项目和资源（Resource）编辑需求，**必须严格按照以下步骤执行**。

> **前置**：若尚未了解整体编写流程与 skill 调用顺序，请先调用 **cloudspec-idl-guide**。

## Step1 了解CloudSpec资源基础

一个完整的CloudSpec项目分为3部分：资源（resources）、操作（operations）、资源测试（tests）。
CloudSpec IDL语法简介、快速开始查看[quick-start.md](references/quick-start.md)

### 识别API风格（前置必做）

CloudSpec项目分为**RPC**和**ROA（RESTful）**两种风格，两者对资源关联的操作在参数传递、错误处理等方面有根本性差异。**在编辑资源前必须先识别风格**。

**优先使用 CLI 命令获取项目信息**：

```bash
aliyun cspec baseinfo
```

该命令输出 JSON 格式的项目基本信息，包含 `apiStyle`（rpc/roa）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`（资源列表）、`operations`（操作列表）等字段。从 `apiStyle` 字段即可判断 API 风格。

仅当 CLI 不可用时，再手动检查 `main.cspec` 或操作文件，详见[roa-vs-rpc.md](references/roa-vs-rpc.md)。

### 资源概述

Resource（资源）是CloudSpec IDL中的核心组件，描述云上长生命周期对象的元数据。

- `./resources`目录存放所有的资源元数据，每个`.cspec`文件对应一个资源定义。
- 资源由三部分组成：主键（identifyDefinition）、属性（properties）、生命周期操作关联（create/get/update/delete/list）。
- 属性可以通过引用外部struct定义（`properties: StructName`），也可以内联定义（`properties: { ... }`）。
- 资源属性和操作的映射关系通过`@operationMapping`配置，或通过`@autoMapping`自动推断。

> **flagMode特殊说明**：部分项目使用`@flagMode`，资源属性通过`#[C,U,R,D,L]`标记声明与各操作的映射关系，此时操作的input/output由系统自动生成。使用前需确认当前项目是否采用此模式。详见[flag-mode.md](references/flag-mode.md)。

### 相关文档

- 资源注解详解：[resource-annotate.md](references/resource-annotate.md)
- 字段约束注解：[field-annotate.md](references/field-annotate.md)
- 操作映射参考（@operationMapping/@defineOperationMapping）：[operation-mapping.md](references/operation-mapping.md)
- Flag模式参考（仅flagMode项目适用）：[flag-mode.md](references/flag-mode.md)
- RPC/ROA风格对照：[roa-vs-rpc.md](references/roa-vs-rpc.md)

> **需要更详细的注解说明？** 以下完整文档位于共享知识库 `cloudspec-shared-knowledge`：
> - Resource 注解完整规范：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/resource-annotate.md`
> - 企业级能力注解（@terraform/@ros/@config/@rmc 等）：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/enterprise-annotate.md`
> - 传值约束注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/assignment-constraint-annotate.md`
> - 值约束注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/value-constraint-annotate.md`
> - 文档注解完整版：`../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/document-annotate.md`

## Step2 确认任务类型

**必须**先明确用户的需求属于以下哪种任务：

| 任务类型 | 说明 |
|---------|------|
| **添加字段** | 向已有资源的properties中添加新的属性 |
| **修改字段** | 修改已有字段的类型或注解 |
| **删除字段** | 从资源properties中移除字段 |
| **修改资源注解** | 修改资源级注解（@arn、@document、@resourceBaseInfo、@references等） |
| **修改字段注解** | 修改字段级注解（@required、@readonly、@document、@length等） |
| **修改flag标记** | 调整字段的`#[C,U,R,D,L]`标记（仅@flagMode资源） |
| **修改主键** | 修改identifyDefinition中的主键定义 |
| **修改操作映射** | 修改@operationMapping或@defineOperationMapping的映射配置 |

如果用户的需求不够明确，**必须**询问以下信息：

- 目标资源名称
- 具体要修改的字段或注解
- 新字段的类型、是否必填、是否只读

## Step3 执行编辑

### Step3.1 阅读参考文档

根据任务类型，阅读相应的参考文档：

- **所有任务**都需要先阅读[resource-annotate.md](references/resource-annotate.md)了解资源注解体系
- **涉及字段注解**时需要阅读[field-annotate.md](references/field-annotate.md)了解字段约束注解的用法
- **涉及操作映射**（@operationMapping/@defineOperationMapping）时需要阅读[operation-mapping.md](references/operation-mapping.md)了解映射语法、RPC/ROA 差异和常见配置
- **涉及flag标记**（仅@flagMode项目）时需要阅读[flag-mode.md](references/flag-mode.md)了解flag的含义和选择规则

### Step3.2 理解当前项目上下文

在执行编辑前，**必须**：

1. **运行 `aliyun cspec baseinfo` 获取项目基本信息**（namespace、apiStyle、资源列表、操作列表等）。如 CLI 不可用，则阅读 `main.cspec` 获取 namespace 信息。
2. 阅读目标资源的`.cspec`文件，完整理解当前资源的结构、注解和flag配置。
3. 确认资源是否使用`@flagMode`——这决定了字段是否需要flag标记。
4. 查看`./operations`目录下该资源的操作定义，理解资源属性与操作的映射关系。
5. 如果涉及资源间引用关系，查看`@references`注解中引用的其他资源。

### Step3.3 添加字段

向资源的properties中添加新字段时，需要按以下步骤操作：

#### 确定字段类型

CloudSpec支持的基本类型：

| 类型 | 说明 |
|-----|------|
| `string` | 字符串 |
| `int32` | 32位整数 |
| `int64` | 64位整数 |
| `float` | 单精度浮点 |
| `double` | 双精度浮点 |
| `boolean` | 布尔值 |
| `any` | 任意类型（尽量避免） |

复合类型：

| 类型 | 语法 |
|-----|------|
| 内联struct | `FieldName: { Key: string, Value: string }` |
| 引用struct | `FieldName: StructName` |
| 内联数组 | `FieldName: array<string>` 或 `FieldName: array<{ Key: string }>` |
| 引用数组 | `FieldName: ArrayName` |
| Map | `FieldName: map<string>` |
| 枚举引用 | `FieldName: EnumName` 或 `@enums(EnumName) FieldName: string` |

> ⚠️ **常见类型与注解陷阱**：
> - **类型名**：没有 `integer`（用 `int32`/`int64`）、没有 `long`（用 `int64`）、没有 `number`（用 `float`/`double`）、没有 `list`（用 `array`）。
> - **注解数组语法**：对象项之间**必须有逗号**（如 `@enums([{...}, {...}])`）。
> - **不存在的注解**：`@backend`（用 `@backendConfigurationHttp`）、`@idempotent`（移除）、`@pattern`（用 `@regexPattern`）、`@maxLength`（用 `@length`）、`@min`/`@max`（用 `@range`）— 这些会报 `自定义 annotate xxx 未定义`。

#### 确定字段注解

根据字段特性添加注解（详细说明参考[field-annotate.md](references/field-annotate.md)）：

```cspec
@required                                    // 必填
@readonly                                    // 只读
@default("")                                 // 默认值
@document({ name: "名称", zh: "中文描述", en: "English" })  // 文档
@length({ min: 2, max: 128 })               // 字符串长度约束
@range({ min: 1, max: 100 })                // 数值范围约束
@regexPattern("^[a-zA-Z].*$")               // 正则校验
@format("iso8601")                           // 值格式
@sensitive                                   // 敏感信息
@clientProhibited                            // 禁止客户端传值
@enums(EnumName)                             // 关联枚举
@initialStatus("Pending")                    // 初始状态（Status字段用）
```

#### 字段添加示例

**用户可编辑字段**：

```cspec
@document({
  name: "资源名称"
  zh: "资源名称"
  en: "Resource name"
})
@length({ min: 2, max: 128 })
Name: string
```

**不可变字段**（创建后不可修改）：

```cspec
@required
@readonly
@document({
  name: "实例规格"
  zh: "实例规格"
  en: "Instance type"
})
InstanceType: string
```

**只读系统字段**：

```cspec
@readonly
@hasDefaultValue
@document({
  name: "创建时间"
  zh: "资源创建时间"
  en: "Creation time"
})
@format("iso8601")
CreateTime: string
```

**复杂类型字段（Tags）**：

```cspec
@document({
  name: "标签"
  zh: "资源标签"
  en: "Resource tags"
})
Tags: array<
  @document({ name: "标签", zh: "标签", en: "Tag" })
  @arrayConfig({ unordered: true })
  {
    @required
    @document({ name: "标签键", zh: "标签键", en: "Tag key" })
    Key: string
    @document({ name: "标签值", zh: "标签值", en: "Tag value" })
    Value: string
  }
>
```

**服务端生成字段**：

```cspec
@readonly
@clientProhibited
@document({
  name: "更新时间"
  zh: "更新时间"
  en: "Update time"
})
@format("iso8601")
UpdateTime: string
```

> **flagMode下的字段写法**：如果资源使用了`@flagMode`，每个字段（`RegionId`除外）需要在类型后附加flag标记，如 `Name: string #[C,U,R,L]`、`Status: string #[R,L]`。flag的含义和选择规则详见[flag-mode.md](references/flag-mode.md)。

### Step3.4 修改字段

修改已有字段时：

1. **只修改用户指定的部分**，保留所有未涉及的注解和结构。
2. 修改字段类型时，需要检查该字段在操作的input/output中是否有引用，确保类型一致。
3. 如果资源使用`@flagMode`且需要修改flag标记，需理解每个flag的含义（参考[flag-mode.md](references/flag-mode.md)），确保与字段语义匹配。
4. 字段注解之间不能有空行，注解紧贴字段。

### Step3.5 删除字段

删除字段前**必须**全面检查并清理所有引用，遗漏引用会导致编译失败：

1. **操作 input/output 结构体**：使用 `grep -r "FieldName" ./operations/` 搜索所有引用，删除或注释掉引用该字段的行。
2. **`@operationMapping` 映射**：检查 `resourceAttributeMappings` 中是否有 `resourceProperty: "$.FieldName"` 的映射，如有则删除。
3. **测试用例**：检查 `./tests` 目录下是否使用了该字段，如有则同步删除。
4. **资源内其他引用**：检查同一资源文件中是否有 `@conflictsWith`、`@references` 等注解引用了该字段。

**删除后必须运行 `aliyun cspec build` 验证编译通过**，如果仍有引用错误，继续清理直至编译成功。

### Step3.6 修改资源级注解

资源级注解的修改规则：

#### @arn

新格式使用string参数：

```cspec
@arn("acs:{service}:${Region}:${AccountId}:{resourceType}/${ResourceId}")
```

center/global scope的资源`${Region}`可为空：

```cspec
@arn("acs:{service}::${AccountId}:{resourceType}/${ResourceId}")
```

#### @document

```cspec
@document({
  zh: "中文描述"
  en: "English description"
  name: "中文显示名"
  nameEn: "English Name"
})
```

#### @resourceBaseInfo

```cspec
@resourceBaseInfo({
  classification: "normal"          // normal | sub | singleton | virtual | readonly | relation
  deliveryScope: "region"           // region | global | Center | Zone
  paidType: "Free"                  // Free | PayAsYouGo | Subscription | SpecifiedByParameter
  hozComponentList: ["RAM", "TERRAFORM", "ROS"]
  consoleListUrl: ""
  consoleDetailUrl: ""
  availableSites: ["china"]
})
```

#### @references

声明与其他资源的引用关系：

```cspec
@references([{
  relatedResource: alicloud.VPC.Vpc.v20140526#VPC
  localProperty: '$.VpcId'
  remoteProperty: '$.VpcId'
}])
```

#### @flagMode（可选）

启用flag模式（标记注解，无参数）。一旦启用，所有属性（`RegionId`除外）都必须有flag标记`#[C,U,R,D,L]`。此模式下操作的input/output由系统根据flag自动生成。不是所有项目都使用此模式，编辑前须确认当前资源是否已启用。

#### 资源注解顺序

```
@arn(...)                   // 1. ARN
@references([...])          // 2. 引用关系（可选）
@document({...})            // 3. 文档描述
@flagMode                   // 4. Flag模式（可选，非所有项目使用）
@resourceBaseInfo({...})    // 5. 资源基本信息
@terraform({...})           // 6. 集成配置（可选）
@config({...})              // 7. 集成配置（可选）
@ros({...})                 // 8. 集成配置（可选）
@notDestroy({...})          // 9. 不销毁标记（可选）
```

### Step3.7 修改操作映射

> **前置必读**：操作映射的完整语法、字段定义和 RPC/ROA 差异详见 [operation-mapping.md](references/operation-mapping.md)。
>
> **Flag 模式不需要手动配置映射**：使用 `@flagMode` 的资源通过 `#[C,U,R,D,L]` 标记自动推导映射，无需以下配置。

#### 操作关联

```cspec
resource MyResource {
  identifyDefinition: { ... }
  properties: { ... }
  create: CreateMyResource
  get: @operationMapping(MyResource_get_mapping) GetMyResource
  update: @operationMapping(MyResource_update_mapping) UpdateMyResource
  delete: @operationMapping(MyResource_delete_mapping) DeleteMyResource
  list: @operationMapping(MyResource_list_mapping) ListMyResources
  operations: [AttachDisk, DetachDisk]      // 非生命周期操作
  resources: [SubResource]                   // 子资源
}
```

#### @defineOperationMapping 核心配置

映射结构体通过 `@defineOperationMapping` 注解定义，常用配置项：

| 配置项 | 用途 | 常用操作 |
|--------|------|---------|
| **resourceAttributeMappings** | 资源属性与 API 入参/出参的映射 | Get、Update、List |
| **retryPolicies** | 按错误码重试 | Delete、Update |
| **resourceNotExistCondition** | 资源不存在判定 | Get、Delete |
| **asyncPollingByProperty** | 异步操作轮询 | Create、Delete |
| **pagination** | 分页配置 | List |
| **constInputParameters** | 常量入参 | 各操作 |

#### RPC 风格映射示例（NAS FileSystem）

RPC 风格下，入参通过 `query` 或 `formData` 传递，参数名为 PascalCase：

```cspec
@defineOperationMapping
struct FileSystem_get_DescribeFileSystems_mapping {
  resourceAttributeMappings: [{
    resourceProperty: "$.FileSystemId"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "FileSystemId"
    requestIn: "query"
    responsePathType: "normal"
    responsePath: "$.FileSystems.FileSystem[*].FileSystemId"
  }
  , {
    resourceProperty: "$.Status"
    responsePathType: "normal"
    responsePath: "$.FileSystems.FileSystem[*].Status"
  }]
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$.FileSystems.FileSystem[*]"
  }
  resourceNotExistCondition: {
    allOf: [{
      notExistCheckType: "checkErrorCode"
      resourceNotExistErrorCodes: ["InvalidFileSystem.NotFound"]
    }]
  }
}
```

#### ROA 风格映射示例（OpenAPIExplorer ApiMcpServerCore）

ROA 风格下，参数名为 camelCase，body 参数的 `requestPath` 需加 **`body.` 前缀**：

```cspec
// Get 映射：path/query 参数直接用参数名，responsePath 为 camelCase
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

// Create 映射：body 参数用 body. 前缀
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

#### RPC 与 ROA 映射关键差异

| 维度 | RPC 风格 | ROA 风格 |
|------|---------|---------|
| **参数命名** | PascalCase（如 `FileSystemId`） | camelCase（如 `oauthClientId`） |
| **requestPath（body 参数）** | 直接参数名（如 `"Queue.QueueName"`） | **`body.` 前缀**（如 `"body.oauthClientId"`） |
| **入参位置 requestIn** | `query`（读）/ `formData`（写） | `body`（写）/ `query`（读）/ `path` |
| **responsePath 命名** | PascalCase（如 `$.Queue.QueueName`） | camelCase（如 `$.oauthClientId`） |
| **分页类型** | `pageNumber`（页码分页） | `nextToken`（游标分页） |

更多映射场景（属性名不同、仅出参映射、异步轮询、常量入参等）详见 [operation-mapping.md](references/operation-mapping.md)。

## Step4 验证

编辑完成后，执行以下验证：

### 4.1 语法检查

在CloudSpec项目根目录下运行：

```bash
aliyun cspec build
```

确保编译通过。如果编译失败，根据错误信息修复后重新编译。

编译通过后，对涉及的资源运行规范检查。资源编辑必须运行 aliyun cspec check --name <ResourceName>（实际命令如下），除非 `aliyun cspec baseinfo` 显示 `isInnerApi: true`：

> **⚠️ Inner API 豁免**：运行 `aliyun cspec baseinfo`，若 `isInnerApi` 为 `true`，则该项目为 inner API，**不需要运行 `aliyun cspec check`**，只需确保 `aliyun cspec build` 通过即可。

```bash
aliyun cspec check --name <ResourceName>
```

### 4.2 规范检查清单

**通用检查：**

- [ ] 文件以`$version: 1`和`namespace:`开头
- [ ] 资源有`@arn`、`@document`、`@resourceBaseInfo`
- [ ] 主键字段在`identifyDefinition`中，且标记`@readonly`
- [ ] 系统字段（`CreateTime`、`Status`等）标记`@readonly`
- [ ] 注解与目标之间没有空行
- [ ] 注解使用冒号语法（`key: value`），不使用等号
- [ ] PascalCase命名（资源名、字段名）
- [ ] 2空格缩进
- [ ] 新增字段的`@document`包含`name`、`zh`、`en`
- [ ] 删除字段后没有遗留引用

**仅@flagMode资源：**

- [ ] 所有属性（`RegionId`除外）都有flag标记`#[...]`
- [ ] flag标记与字段语义匹配（必填字段有`C`、只读字段无`C`和`U`）

## 注意事项

- 编辑资源时，**只修改用户指定的部分**，保留所有未涉及的注解、字段和结构。
- 标准字段命名：`RegionId`、`Status`、`CreateTime`、`UpdateTime`，不要改变大小写。
- 注解值使用冒号语法（`@document({ zh: "描述" })`），不使用等号语法。
- 严禁修改API元数据（`./operations`目录下的文件），除非用户明确要求。
- 添加`@required`字段时要考虑对已有测试用例和操作input的影响，可能需要同步更新。
- `@arn`新格式使用string参数而非struct，旧的struct格式（`standard`、`ram`、`rmc`）已不推荐。
- 复杂类型字段（array、struct嵌套）的注解放在类型声明的角括号`<>`内部。
- `@flagMode`是可选模式，编辑前务必确认当前资源是否使用。如果使用，`RegionId`是唯一不需要flag标记的特殊字段，flag标记放在类型（或`>`闭合符号）之后。

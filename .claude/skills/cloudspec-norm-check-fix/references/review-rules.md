# CloudSpec 审查规则详解

## 一、语法和格式规则

### 1.1 文件头

每个`.cspec`文件必须以版本和namespace开头：

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{PopVersion}
```

**常见错误：**
- 缺少`$version: 1`
- namespace格式不正确（缺少`alicloud.`前缀、版本格式不对）
- namespace与项目`main.cspec`中不一致

### 1.2 注解语法

注解使用冒号语法，不使用等号：

```cspec
// 正确
@document({ zh: "描述", en: "description" })

// 错误
@document(zh = "描述", en = "description")
```

注解与目标之间不能有空行：

```cspec
// 正确
@required
Name: string

// 错误
@required

Name: string
```

### 1.3 缩进和空行

- 统一使用2空格缩进
- 不同的顶层块（resource、operation、error、struct、enum）之间保留一个空行
- 同一块内部的字段之间不应有空行（注解除外）

### 1.4 命名规范

- 资源名、操作名、结构体名、字段名：PascalCase（`CreateInstance`、`InstanceName`）
- 操作命名：`动词+资源名`（`CreateInstance`、`ListInstances`）
- List操作使用复数形式（`ListInstances`，不是`ListInstance`）
- 命名空间中的Product和PopCode通常使用PascalCase

### 1.5 保留字转义

`resource`、`service`、`operation`等是保留字，作为struct的key时必须用反引号转义：

```cspec
// 正确
@defineAction
struct Op_RAM_Action_0 {
  action: "svc:Op"
  resources: [{
    `resource`: "acs:svc:..."
    required: true
  }]
}

// 错误 - resource未转义
struct Op_RAM_Action_0 {
  action: "svc:Op"
  resources: [{
    resource: "acs:svc:..."
    required: true
  }]
}
```

## 二、资源审查规则

### 2.1 结构完整性

资源必须包含`identifyDefinition`和`properties`：

```cspec
resource MyResource {
  identifyDefinition: {
    @readonly
    MyResourceId: string
  }
  properties: MyResourceProperties
}
```

或内联定义properties：

```cspec
resource MyResource {
  identifyDefinition: {
    @readonly
    MyResourceId: string
  }
  properties: {
    Name: string
    Status: string
  }
}
```

**检查要点：**
- 主键字段在`identifyDefinition`中，不在`properties`中
- 主键字段标记`@readonly`
- properties不为空

### 2.2 必要注解

| 注解 | 必填 | 检查要点 |
|-----|------|---------|
| `@arn` | 是 | ARN格式正确，包含`${AccountId}`和资源标识符 |
| `@document` | 是 | 包含`name`字段 |
| `@resourceBaseInfo` | 是 | 包含`classification`、`deliveryScope`、`paidType` |
| `@flagMode` | 否 | 可选，启用后所有属性需有flag标记 |
| `@references` | 否 | 引用关系中的资源component ID格式正确 |

### 2.3 @arn格式

新格式（推荐）：

```cspec
@arn("acs:{service}:${Region}:${AccountId}:{resourceType}/${ResourceId}")
```

center/global scope的资源`${Region}`可为空：

```cspec
@arn("acs:{service}::${AccountId}:{resourceType}/${ResourceId}")
```

旧格式（仍可用）：

```cspec
@arn({
  standard: "acs:{service}:${Region}:${AccountId}:{resourceType}/${ResourceId}"
})
```

**检查要点：**
- 包含`${AccountId}`
- 资源标识符变量与`identifyDefinition`中的主键对应
- `${Region}`是否存在取决于`@resourceBaseInfo`中的`deliveryScope`（`region`/`Zone`需要，`center`/`global`可省略）

### 2.4 @resourceBaseInfo

```cspec
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
```

| 字段 | 合法值 |
|-----|-------|
| `classification` | `"normal"`, `"sub"`, `"singleton"`, `"virtual"`, `"readonly"`, `"relation"` |
| `deliveryScope` | `"region"`, `"global"`, `"Center"`, `"Zone"` |
| `paidType` | `"Free"`, `"PayAsYouGo"`, `"Subscription"`, `"SpecifiedByParameter"`, `"free"`, `"Paid"` |

### 2.5 字段语义检查

| 字段类型 | 应有注解 | 不应有注解 |
|---------|---------|----------|
| 系统生成只读字段（CreateTime等） | `@readonly`、`@hasDefaultValue` | `@required`（除非同时有`@clientOptional`） |
| 服务端内部字段 | `@readonly`、`@clientProhibited` | — |
| 用户必填字段 | `@required` | `@readonly`（除非是创建后不可改） |
| 创建后不可改字段 | `@readonly` | — |

## 三、操作审查规则

### 3.1 必要注解

| 注解 | 必填 | 检查要点 |
|-----|------|---------|
| `@backendConfigurationHttp` | 是 | 包含`applicationName`、`timeout`、`backendUrl` |
| `@document` | 是 | 包含`name`字段 |
| `@errorMapping` | 是 | 包含`errorExpression`、`codeField`、`errorMessageField` |
| `@http` | 是 | 包含`schemes`、`methods`、`authenticators` |
| `@operationInfo` | 是 | 包含操作类型信息 |
| `@ram` | 是 | 包含`enable`、`level` |
| `@requiredPermission` | 是 | 引用有效的`@defineAction`结构体 |
| `@visibility` | 是 | `"Public"` 或 `"Private"` |
| `@for` | 仅flagMode | 仅当关联资源使用`@flagMode`时必填 |

### 3.2 操作类型一致性

`@operationInfo`中的操作类型字段有多种写法，以项目实际使用的为准：

| 操作类型 | operationType | typeFromOperation | operationTypeOld | HTTP方法 | @document.name前缀 |
|---------|--------------|-------------------|------------------|----------|------------------|
| Create | `"Write"` 或 `"create"` | `"create"` | `"write"` | `["post"]` | "创建..." |
| Get | `"Read"` 或 `"get"` | `"get"` | `"read"` | `["get", "post"]`或`["get"]` | "查询...详情" |
| Update | `"Write"` 或 `"update"` | `"update"` | `"write"` | `["post"]` | "修改..." |
| Delete | `"Write"` 或 `"delete"` | `"delete"` | `"write"` | `["post"]` | "删除..." |
| List | `"List"` 或 `"list"` | `"list"` | `"read"` | `["get", "post"]`或`["get"]` | "查询...列表" |

**检查要点：**
- 操作类型字段（`operationType`/`typeFromOperation`/`operationTypeOld`）与操作语义一致
- HTTP方法与操作类型匹配
- `@document.name`使用正确的中文动词（创建/查询/修改/删除，不使用新建/获取/更新）
- `@operationInfo`中的其他字段（`riskType`、`chargeType`等）与同项目其他操作保持一致

### 3.3 错误定义

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

**检查要点：**
- 每个操作至少关联一个error
- error包含`httpCode`、`errorCode`、`errorMessage`
- `type`为`"user"`或`"system"`
- `default`字段的使用与同项目其他操作保持一致（有的项目使用`default: true`标记兜底错误，有的项目所有error均为`default: false`）
- 如果项目使用`default: true`，同一操作只能有一个`default: true`的error

### 3.4 RAM Action结构体

```cspec
@defineAction
struct CreateItem_RAM_Action_0 {
  action: "service:CreateItem"
  resources: [{
    `resource`: "acs:service:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

**检查要点：**
- 有`@defineAction`注解
- `action`格式为`{小写服务名}:{操作名}`
- `resource`键使用反引号包裹
- `@requiredPermission`引用的结构体名称存在

### 3.5 List操作专项

- [ ] 有`@numberPaginated`或`@paginated`注解
- [ ] 有`@rootMapping`注解，`responsePath`格式正确（通常为`$.{ResourceName}s[*]`）
- [ ] 操作名称使用复数形式

### 3.6 注解顺序

推荐的注解顺序：

```
@for(ResourceName)                  // 仅flagMode
@actionTrail({...})                 // 可选
@backendConfigurationHttp({...})
@controlPolicy({...})              // 可选
@document({...})
@errorMapping({...})
@gatewayOptions({...})             // 可选
@http({...})
@numberPaginated({...})             // 仅List
@operationInfo({...})
@ram({...})
@requiredPermission(StructName)
@returnMode({...})                 // 可选
@rootMapping({...})                 // 仅List
@visibility("...")
```

注解顺序不一致归类为**建议**级别，不是错误。

## 四、配置一致性规则

同一项目内的操作应保持配置风格一致：

| 配置项 | 检查内容 |
|-------|---------|
| `@backendConfigurationHttp` | `applicationName`相同、`timeout`格式一致、`backendUrl`格式一致 |
| `@http` | `schemes`格式一致、`authenticators`相同 |
| `@errorMapping` | 字段名称一致（`codeField`、`errorMessageField`等） |
| `@ram` | `level`和`atGateway`配置一致 |

一致性问题归类为**建议**级别。

## 五、flagMode专项规则

仅当资源使用`@flagMode`时适用：

### 5.1 flag标记完整性

- 所有properties中的字段（`RegionId`除外）必须有flag标记
- identifyDefinition中的字段也需要flag标记

### 5.2 flag语义正确性

| 字段特性 | flag应包含 | flag不应包含 |
|---------|----------|------------|
| `@readonly`系统字段 | `R` | `C`, `U` |
| `@required`创建字段 | `C` | — |
| 主键字段 | `R`, `D` | `C`, `U` |
| 用户可编辑字段 | `C`, `U`, `R` | — |

### 5.3 关联操作检查

- 操作必须有`@for(ResourceName)`注解
- 操作体只声明`errors`（input/output由系统生成）

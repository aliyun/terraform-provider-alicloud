# CloudSpec IDL 语法参考

本文档基于 CloudSpec IDL 的正式语法定义（ANTLR4 Grammar），提供完整的语法规则说明。用于在编译修复时准确判断语法是否合法、定位语法错误根因。

---

## 一、文件结构

每个 `.cspec` 文件必须按以下顺序组织：

```
$version: 1                          ← 版本声明（必须在最前面）
namespace: alicloud.Product.PopCode.vYYYYMMDD   ← 命名空间声明

<类型定义>*                           ← 零个或多个类型定义，顺序无关
```

### 1.1 版本声明

```cspec
$version: 1
```

- 必须是文件的第一个非注释/非注解内容
- `$version:` 后跟数字

### 1.2 命名空间声明

```cspec
namespace: alicloud.ECS.Ecs.v20140526
```

- 格式为 `namespace:` 后跟点分隔的标识符
- 每段标识符必须以字母开头，可包含字母、数字、下划线、连字符
- 典型格式：`alicloud.{Product}.{PopCode}.v{PopVersion}`

---

## 二、顶层类型定义

以下类型定义可以在文件中以**任意顺序**出现，互不依赖顺序：

| 关键字 | 类型 | 说明 |
|--------|------|------|
| `service` | 服务定义 | 包含资源和操作的集合 |
| `resource` | 资源定义 | 云上长生命周期对象 |
| `operation` | 操作定义 | OpenAPI 接口 |
| `struct` | 结构体定义 | 键值对集合 |
| `array` | 数组定义 | 单一类型的有序集合 |
| `map` | Map 定义 | 字符串键的映射 |
| `error` | 错误定义 | 操作的错误码 |
| `enum` | 枚举定义 | 字符串枚举值 |
| `intEnum` | 整数枚举定义 | 整数枚举值 |
| `constraint` | 约束定义 | 模型校验规则 |
| `import` | 导入声明 | 引入外部命名空间 |
| `$apply` | Apply 声明 | 应用配置 |
| `$test` | 测试定义 | 资源测试用例 |
| 基础类型别名 | 类型别名 | 如 `string MyType` |

---

## 三、基础类型

### 3.1 基础类型列表

| 类型 | 说明 |
|------|------|
| `string` | 字符串 |
| `binary` | 字节流 |
| `byte` | Base64 编码的字节流 |
| `int32` | 32 位有符号整数 |
| `int64` | 64 位有符号整数 |
| `float` | 单精度浮点数 |
| `double` | 双精度浮点数 |
| `boolean` | 布尔值（`true` / `false`） |
| `any` | 任意类型 |
| `$null` | 空值 |

### 3.2 复合类型关键字

| 类型 | 说明 |
|------|------|
| `struct` | 结构体 |
| `array` | 数组 |
| `map` | 映射 |

### 3.3 枚举类型关键字

| 类型 | 说明 |
|------|------|
| `enum` | 字符串枚举 |
| `intEnum` | 整数枚举 |

---

## 四、注解语法

注解以 `@` 开头，有三种形式：

### 4.1 标记注解（无参数）

```cspec
@required
@readonly
@flagMode
@sensitive
@deprecated
@clientProhibited
@notResourceProperty
@defineAction
@defineCondition
```

### 4.2 单参数注解

```cspec
@visibility("Public")
@default("")
@regexPattern("^[a-zA-Z].*$")
@format("iso8601")
@apiStyle("rpc")
@arn("acs:{service}:${Region}:${AccountId}:{resourceType}/${ResourceId}")
@for(ResourceName)
@requiredPermission(StructName)
@initialStatus("Pending")
@retryable(true)
```

- 参数可以是：字符串、数字、布尔值、数组、结构体、类型名
- 字符串支持单引号 `'...'` 和双引号 `"..."`

### 4.3 结构体参数注解

```cspec
@document({
  name: "名称"
  zh: "中文描述"
  en: "English description"
})

@http({
  method: ["post"]
  uri: "/v1/resources"
  apiStyle: "restful"
})

@length({ min: 2, max: 128 })
@range({ min: 1, max: 100 })
```

- 使用花括号 `{...}` 包裹
- 键值对使用**冒号** `key: value`，**不使用等号**
- 多个键值对之间可以用逗号分隔，也可以换行分隔

### 4.4 多参数注解

```cspec
@enums(EnumName)
@conflictsWith(['Id'])
@conditions([RAMConditionEncrypted])
```

- 多个参数用逗号分隔

### 4.5 注解规则

- **注解必须紧贴目标**：注解与它修饰的目标（字段、类型定义）之间**不能有空行**
- 多个注解可以连续堆叠在目标上方
- 注解可以出现在：类型定义之前、字段之前

```cspec
// ✅ 正确：注解紧贴目标
@required
@document({ name: "名称" })
Name: string

// ❌ 错误：注解与目标之间有空行
@required

Name: string
```

---

## 五、Service 定义

```cspec
@apiStyle("rpc")
service DessertFactoryService {
  version: "2023-04-01"
  resources: [Candy, Cake]
  operations: [GetWorkshopsNumber]
}
```

- 使用 `service` 关键字
- 花括号内是键值对形式的属性
- `resources` 和 `operations` 的值是数组引用

---

## 六、Resource 定义

```cspec
@arn("acs:example:{#regionId}:{#accountId}:candy/{#CandyId}")
@document({ name: "糖果", zh: "糖果资源" })
@resourceBaseInfo({ classification: "normal", deliveryScope: "region" })
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  get: GetCandy
  update: [UpdateCandy]
  delete: DeleteCandy
  list: ListCandies
  operations: [AttachLabel]
  resources: [SubResource]
}
```

- 使用 `resource` 关键字
- 花括号内的每一项都是 `key: value` 形式
- `identifyDefinition` 的值是内联结构体
- `properties` 的值可以是引用名或内联结构体
- `create`/`get`/`delete` 的值是操作名（单个）
- `update` 的值可以是操作名或操作名数组
- `list` 的值是操作名（单个）
- `operations` 是非生命周期操作的数组
- `resources` 是子资源的数组
- 生命周期操作前可以加 `@operationMapping(MappingStructName)`

---

## 七、Operation 定义

```cspec
@http({ method: ["post"] })
@document({ name: "创建糖果" })
@visibility("Public")
operation CreateCandy {
  input: CreateCandyInput
  output: CreateCandyOutput
  errors: [Error_CreateCandy]
}
```

- 使用 `operation` 关键字
- 花括号内包含 `input`、`output`、`errors` 三个属性
- `input` 和 `output` 引用 struct 名称
- `errors` 是 error 名称的数组
- 如果资源使用 `@flagMode`，操作体内可以只声明 `errors`（input/output 自动生成）

---

## 八、Struct 定义

### 8.1 独立 struct

```cspec
struct CreateCandyInput {
  @required
  @parameter({ in: 'query', name: 'CandyName', required: true })
  CandyName: string

  @document({ name: "颜色" })
  Color: string

  Tags: array<{ Key: string, Value: string }>
}
```

### 8.2 struct 继承

```cspec
struct ExtendedInput extend BaseInput, AnotherBase {
  ExtraField: string
}
```

- 使用 `extend` 关键字，可以继承多个父 struct（逗号分隔）

### 8.3 struct 字段语法

每个字段的完整语法：

```
[注解]*
字段名[?]: 类型 [= 默认值] [#[flags]]
```

- **字段名**：PascalCase 标识符，保留字需用反引号包裹（如 `` `resource` ``）
- **可选标记**：字段名后加 `?` 表示可选
- **类型**：基础类型、引用类型、内联类型
- **默认值**：`= value` 形式（可选）
- **Flags**：`#[C,U,R,D,L]` 形式（仅 `@flagMode` 资源使用）

### 8.4 资源属性引用

在操作的 input/output struct 中，可以用 `$` 引用资源属性：

```cspec
struct GetCandyInput {
  @required
  $CandyId
}

struct GetCandyOutput {
  $CandyId
  $CandyName
  @notResourceProperty
  RequestId: string
}
```

- `$FieldName` 引用资源的同名属性，自动继承类型
- `$FieldName = defaultValue` 可以指定默认值
- 非资源属性需要用 `@notResourceProperty` 标注

---

## 九、Array 定义

### 9.1 独立 array

```cspec
array Candies {
  item: CandyItem
}
```

- 使用 `array` 关键字
- 花括号内只有一个字段，定义数组元素的类型

### 9.2 内联 array

```cspec
Tags: array<
  @document({ name: "标签" })
  @arrayConfig({ unordered: true })
  {
    Key: string
    Value: string
  }
>
```

- 使用 `array<类型>` 语法
- 尖括号内可以包含注解和类型定义
- 类型可以是基础类型、struct 引用或内联 struct

---

## 十、Map 定义

### 10.1 独立 map

```cspec
map ConfigMap {
  key1: string
  key2: int32
}
```

- 使用 `map` 关键字
- 花括号内是键值对，key 可以是标识符或字符串

### 10.2 内联 map

```cspec
Config: map<string>
```

- 使用 `map<值类型>` 语法

---

## 十一、Error 定义

```cspec
@retryable(true)
error Error_CreateCandy_ServiceUnavailable {
  httpCode: 503
  errorCode: "ServiceUnavailable"
  backendErrorCode: "ServiceUnavailable"
  errorMessage: "The request has failed due to a temporary failure of the server."
  type: "user"
  default: false
}
```

- 使用 `error` 关键字
- 花括号内是键值对形式的属性
- `httpCode` 是数字，其他字段通常是字符串或布尔值

---

## 十二、Enum 定义

### 12.1 字符串枚举

```cspec
enum InstanceStatus {
  "Running"
  "Stopped"
  "Pending"
}
```

- 使用 `enum` 关键字
- 花括号内是字符串值列表
- 枚举项可以有注解（如 `@document`）

### 12.2 整数枚举

```cspec
intEnum Priority {
  1
  2
  3
}
```

- 使用 `intEnum` 关键字
- 花括号内是数字值列表

---

## 十三、Constraint 定义

```cspec
constraint MyConstraint {
  selector: "..."
  condition: "..."
}
```

- 使用 `constraint` 关键字
- 花括号内是键值对形式的属性

---

## 十四、Import 声明

```cspec
import "alicloud.VPC.Vpc.v20140526"
```

- 使用 `import` 关键字后跟字符串
- 引入外部命名空间，使得可以通过组件 ID 引用其类型

---

## 十五、组件 ID 引用

跨命名空间引用组件时使用组件 ID 语法：

```cspec
alicloud.VPC.Vpc.v20140526#VPC
```

- 格式：`命名空间#组件名`
- 可以进一步引用子路径：`命名空间#组件名$子路径`

```cspec
@references([{
  relatedResource: alicloud.VPC.Vpc.v20140526#VPC
  localProperty: '$.VpcId'
  remoteProperty: '$.VpcId'
}])
```

---

## 十六、Flag 标记语法

仅在 `@flagMode` 资源中使用，标记字段与 CRUD+List 操作的映射关系：

```cspec
Name: string #[C,U,R,L]
Status: string #[R,L]
InstanceId: string #[!C,R,D,L]
```

- 语法：`#[flag1,flag2,...]`
- 放在字段类型（或 `>` 闭合符号）之后
- 可用的 flag：`C`（Create）、`U`（Update）、`R`（Read/Get）、`D`（Delete）、`L`（List）
- `!` 前缀表示取反（如 `!C` 表示 Create 时不包含）
- `RegionId` 是唯一不需要 flag 标记的特殊字段

---

## 十七、值类型语法

在注解参数、字段默认值等位置可以使用以下值类型：

| 值类型 | 语法 | 示例 |
|--------|------|------|
| 字符串 | `"..."` 或 `'...'` | `"hello"`, `'world'` |
| 数字 | 整数或浮点数 | `42`, `-1`, `3.14`, `1e10` |
| 布尔值 | `true` / `false` | `true` |
| 数组 | `[item1, item2]` | `["a", "b"]`, `[1, 2]` |
| 结构体 | `{ key: value }` | `{ name: "test", count: 1 }` |
| 引用 | 标识符 | `StructName`, `EnumName` |
| 内联 array | `array<类型>` | `array<string>`, `array<{ K: string }>` |
| 内联 map | `map<类型>` | `map<string>` |

---

## 十八、Apply 声明

```cspec
$apply "alicloud.Common.v1" @someAnnotation ApplyName
```

- 以 `$apply` 开头
- 后跟组件 ID 或字符串标识
- 可以附带注解
- 可选的名称标识

---

## 十九、Test 定义

```cspec
$test ComponentId TestName {
  ...
}
```

- 以 `$test` 开头
- 后跟组件 ID 和测试名称
- 花括号内是 struct 形式的测试数据

---

## 二十、注释

```cspec
// 单行注释

/* 多行注释
   可以跨行 */

/// 文档注释（单行）
```

- `//` 单行注释
- `/* ... */` 多行注释
- `///` 文档注释

---

## 二十一、保留字列表

以下关键字在 CloudSpec IDL 中有特殊含义，作为 struct 的 key 时**必须用反引号包裹**：

| 保留字 | 说明 |
|--------|------|
| `resource` | 资源定义关键字 |
| `service` | 服务定义关键字 |
| `operation` | 操作定义关键字 |
| `error` | 错误定义关键字 |
| `struct` | 结构体定义关键字 |
| `array` | 数组定义关键字 |
| `map` | Map 定义关键字 |
| `enum` | 枚举定义关键字 |
| `intEnum` | 整数枚举定义关键字 |
| `constraint` | 约束定义关键字 |
| `import` | 导入关键字 |
| `string` | 基础类型 |
| `binary` | 基础类型 |
| `byte` | 基础类型 |
| `int32` | 基础类型 |
| `int64` | 基础类型 |
| `float` | 基础类型 |
| `double` | 基础类型 |
| `boolean` | 基础类型 |
| `any` | 基础类型 |
| `true` | 布尔值 |
| `false` | 布尔值 |

反引号包裹示例：

```cspec
// RAM Action 中 resource 是保留字
struct CreateItem_RAM_Action {
  action: "svc:CreateItem"
  resources: [{
    `resource`: "acs:svc:{#regionId}:{#accountId}:*"
    required: true
  }]
}
```

---

## 二十二、命名规则

- **标识符**：由字母、数字、下划线、连字符组成，支持中文字符
- **反引号标识符**：`` `name` `` 形式，可包含特殊字符（如 `/`、`:`、`.`、`,`）
- **命名空间**：点分隔的标识符，每段以字母开头
- **组件名**：标识符或数字

---

## 二十三、格式规范

| 规则 | 说明 |
|------|------|
| 缩进 | 2 空格，不使用 Tab |
| 键值对 | 使用冒号 `key: value`，不使用等号 |
| 分隔 | 键值对之间可用逗号或换行分隔 |
| 空行 | 类型定义之间保留一个空行；注解与目标之间不能有空行 |
| 括号 | `{}`、`()`、`[]`、`<>` 必须成对匹配 |

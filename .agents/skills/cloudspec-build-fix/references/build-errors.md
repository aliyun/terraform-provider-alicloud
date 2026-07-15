# CloudSpec 编译错误参考

## 一、语法错误

### 1.1 注解格式错误

**错误表现**：注解参数使用等号而非冒号。

```cspec
// 错误
@document(zh = "描述")

// 正确
@document({ zh: "描述" })
```

**修复方式**：将`=`改为`:`，确保使用花括号包裹多参数注解。

### 1.2 缩进不一致

**错误表现**：混用tab和空格，或缩进层级错误。

**修复方式**：统一使用2空格缩进。可运行`aliyun cspec format`自动格式化。

### 1.3 括号/花括号不匹配

**错误表现**：struct、array、注解参数的括号未闭合。

**修复方式**：检查对应的开闭括号是否匹配，特别注意嵌套的struct和array定义。

### 1.4 注解与目标之间有空行

**错误表现**：注解和它修饰的目标之间插入了空行。

```cspec
// 错误
@required

Name: string

// 正确
@required
Name: string
```

**修复方式**：删除注解与目标之间的空行。

### 1.5 使用了不存在的类型关键字

**错误表现**：`The component MaxResults/integer is not defined`

**常见场景**：误用 `integer`、`long`、`number` 等其他语言的类型名。

```cspec
// 错误
MaxResults: integer

// 正确
MaxResults: int32
```

**修复方式**：使用 CloudSpec 合法基础类型：`string`、`int32`、`int64`、`float`、`double`、`boolean`、`any`、`byte`、`binary`。

### 1.6 使用了 `list` 关键字定义列表

**错误表现**：`extraneous input 'list' expecting {<EOF>, ...}`

**常见场景**：受其他 IDL（如 Smithy/Thrift）影响，在顶层使用 `list` 关键字。

```cspec
// 错误 — list 不是 CloudSpec 的合法顶层关键字
list TestItems {
  member: TestItem
}

// 正确方式 1 — 独立 array 定义
array TestItems {
  item: TestItem
}

// 正确方式 2 — 内联 array（推荐，减少文件间引用）
TestItems: array<
  {
    TestName: string
    Result: string
  }
>
```

**修复方式**：将 `list` 替换为 `array`；如果是定义在操作 output 中，优先用 `array<{...}>` 内联语法。

### 1.7 注解数组中对象项缺少逗号

**错误表现**：`no viable alternative at input '...}{...'`

**常见场景**：`@enums`、`@conditions`、`resources` 等注解的数组参数中，多个对象项之间遗漏逗号。

```cspec
// 错误 — } 和 { 之间缺少逗号
@enums([
  { value: "Basic", description: "基础" }
  { value: "Detail", description: "详细" }
])

// 正确 — 对象项之间必须有逗号
@enums([
  { value: "Basic", description: "基础" },
  { value: "Detail", description: "详细" }
])
```

**修复方式**：在注解数组的相邻对象 `}` 和 `{` 之间添加逗号。同一规则适用于 `resources: [{...},{...}]` 等场景。

### 1.8 使用了不存在的注解（`自定义 annotate xxx 未定义`）

**错误表现**：`自定义 annotate xxx 未定义，请使用 struct 定义，并增加@defineAnnotateGroup注解`

**常见场景**：以下注解**不是 CloudSpec 内置注解**，但在其他 IDL/框架中常见，容易误用：

| 错误写法 | 正确做法 |
|----------|----------|
| `@backend({...})` | 用 `@backendConfigurationHttp({...})`（操作级后端配置） |
| `@idempotent` | 移除（CloudSpec 无此注解） |
| `@pattern("regex")` | 用 `@regexPattern("regex")` |
| `@maxLength(n)` | 用 `@length({ min: 0, max: n })` |
| `@min(n)` / `@max(n)` | 用 `@range({ min: n, max: m })` |

```cspec
// 错误 — @backend 不是内置注解
@backend({
  protocol: "http"
  timeout: 15000
  serviceName: "my-service"
})
operation MyOp { ... }

// 正确 — 使用 @backendConfigurationHttp
@backendConfigurationHttp({
  applicationName: "my-service"
  responseType: "Object"
  timeout: { online: 15000 }
  backendUrl: { online: "http://my-service.aliyuncs.com" }
  sign: true
  signPolicy: "Local"
})
operation MyOp { ... }
```

**修复方式**：将错误注解替换为正确的 CloudSpec 内置注解，或直接移除（如 `@idempotent`）。

### 1.9 使用了不存在的约束注解

**错误表现**：`自定义 annotate min 未定义` / `自定义 annotate maxLength 未定义`

**常见场景**：受 JSON Schema / OpenAPI 影响，使用 `@min`、`@max`、`@maxLength`、`@minLength` 等注解。

```cspec
// 错误 — @min/@max/@maxLength 不是内置注解
@min(1)
@max(100)
NodeCount: int32

@maxLength(256)
Description: string

// 正确 — 使用 @range 和 @length
@range({ min: 1, max: 100 })
NodeCount: int32

@length({ min: 0, max: 256 })
Description: string
```

**修复方式**：`@min`/`@max` → `@range({ min: n, max: m })`；`@maxLength` → `@length({ min: 0, max: n })`。

## 二、引用错误

### 2.1 operation引用了未定义的struct

**错误表现**：`[o-001]operation: XXX referenced struct: YYY is undefined`

**可能原因**：
1. struct名称拼写错误
2. struct定义在其他文件中，但未通过import引入
3. struct确实不存在，需要创建

**修复方式**：
1. 检查struct名称是否拼写正确（PascalCase）
2. 检查是否在同文件或已import的文件中有该struct定义
3. 如果struct确实不存在，根据operation的语义创建对应的struct定义

### 2.2 operation引用了未定义的error

**错误表现**：`[o-001]operation: XXX referenced error: YYY is undefined`

**可能原因**：
1. error名称拼写错误
2. error定义缺失

**修复方式**：
1. 检查error名称是否拼写正确
2. 在同文件中添加缺失的error定义：

```cspec
error Error_OperationName {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: false
}
```

### 2.3 引用的namespace未定义

**错误表现**：`[o-001]operation: XXX referenced namespace: YYY is undefined`

**可能原因**：
1. 使用了组件ID引用（`namespace#ComponentName`），但namespace未import
2. namespace拼写错误

**修复方式**：
1. 在`main.cspec`中检查是否有对应的import声明
2. 修正namespace拼写

## 三、资源生命周期错误

### 3.1 资源引用了不存在的操作

**错误表现**：资源的create/get/update/delete/list引用的操作名在项目中不存在。

**修复方式**：
1. 确认操作名称拼写正确
2. 确认操作文件存在于`./operations`目录
3. 确认操作文件的namespace与资源文件一致

### 3.2 $.XX引用的资源属性不存在

**错误表现**：映射关系中`$.PropertyName`引用的属性在资源定义中找不到。

**可能原因**：
1. 属性名拼写错误
2. 属性已被删除或重命名
3. 资源属性类型与引用方式不匹配（如string类型被当作struct引用内部字段）

**修复方式**：
1. 检查属性名拼写（注意PascalCase）
2. 在资源的`identifyDefinition`和`properties`中查找正确的属性名
3. 如果属性确实不存在，考虑是否需要添加，或修正引用路径

### 3.3 资源测试不允许传入create/update操作不可传入的属性

**错误表现**：测试用例中设置了`@readonly`且`@clientProhibited`的属性。

**可能原因**：
1. 属性是`@readonly`且`@clientProhibited`的系统字段，客户端不应传值
2. 属性虽然`@readonly`但允许客户端传值，只是缺少操作映射

**修复方式**：
1. 如果属性确实是系统字段（`@readonly` + `@clientProhibited`），从测试用例中删除该属性
2. 如果属性应该可写但缺少映射，需要在资源定义中添加对应的`@operationMapping`

### 3.4 类型不匹配

**错误表现**：资源属性类型为string，但测试中使用了struct类型的值（或反之）。

**修复方式**：
1. 以资源定义中的类型为准
2. 修改测试用例中的属性值类型，使其与资源定义一致
3. 如果是array类型，确保测试中的值也用数组语法`[...]`

## 四、保留字错误

### 4.1 保留字未转义

**错误表现**：`resource`、`service`、`operation`等保留字作为struct的key时未用反引号。

```cspec
// 错误
struct RAM_Action {
  action: "svc:Op"
  resources: [{
    resource: "acs:svc:..."    // resource是保留字
    required: true
  }]
}

// 正确
struct RAM_Action {
  action: "svc:Op"
  resources: [{
    `resource`: "acs:svc:..."  // 反引号转义
    required: true
  }]
}
```

**修复方式**：用反引号包裹保留字。

常见保留字列表：`resource`、`service`、`operation`、`error`、`struct`、`enum`、`intEnum`、`map`、`import`、`constraint`、`apply`

## 五、文件头错误

### 5.1 缺少$version声明

**修复方式**：在文件第一行添加`$version: 1`。

### 5.2 namespace不一致

**错误表现**：文件中的namespace与`main.cspec`中声明的不一致。

**修复方式**：以`main.cspec`中的namespace为准，修正当前文件的namespace声明。

## 六、约束错误

约束错误来自CloudSpec内置的规范校验规则，常见的规则编号和含义：

| 规则 | 含义 |
|------|------|
| M-RL-0004~0008 | 资源的CRUD操作应只有一个主键标识符 |
| M-RL-0009 | 资源映射关系校验 |
| M-RL-0012 | 资源创建出参应包含主键 |
| M-RL-0016 | 资源生命周期通用校验 |
| M-RL-0017/0020/0028 | Update/Get/Delete的主键标识符校验 |
| M-RL-0021 | List出参应包含所有主键 |
| M-RL-0022 | Get出参应包含List出参的所有字段 |
| M-RL-0023 | List入参不应包含required参数 |
| M-RL-0025 | List出参应包含资源所有属性 |
| M-RL-0034 | Get出参应包含资源所有属性 |
| M-RL-0036 | List操作应配置分页信息 |
| M-RL-0037 | Update操作不应修改readonly字段 |
| M-RL-0039 | Update操作应包含可更新属性 |
| M-RL-0040 | Update操作应包含所有可更新属性 |
| M-RL-0041 | API参数应在资源属性中有定义 |
| M-RL-0042 | 资源属性的客户端可写性校验 |
| M-RL-0043/0044 | 资源生命周期映射校验 |
| M-AS-0010 | ARN定义校验 |
| M-RT-0015~0047 | 资源类型相关校验 |
| E-PS-0004 | 属性约束校验 |
| R-AC-0010/0012 | 访问控制校验 |
| R-BR-0050~0104 | 破坏性变更校验 |

**修复方式**：
1. 根据规则编号理解校验意图
2. 检查对应的资源定义、操作定义和映射关系
3. 按照规则要求调整，通常涉及添加缺失的映射、修正属性标记或补充操作参数

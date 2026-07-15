# CloudSpec Flag Mode 语法参考

## 概述

Flag 模式（`@flagMode`）是一种声明式的属性映射模式。启用后，资源属性通过 flag 标记声明自己参与哪些 CRUD 操作，系统会自动为操作生成对应的 Input/Output 结构体。

## Flag 语法

Flag 标记紧跟在字段类型之后：

```cspec
FieldName: FieldType #[C,U,R,F,L,A]
```

对于复合类型（array、struct），flag 标记放在类型闭合符号之后：

```cspec
Tags: array<{
  Key: string
  Value: string
}> #[C,U,R,F,L,A]
```

## Flag 定义

| Flag | 全称 | 含义 |
|------|------|------|
| `C` | CREATE_SUPPORT | 创建操作的输入参数 |
| `U` | MODIFY_SUPPORT | 修改操作的输入参数 |
| `D` | DELETE_SUPPORT | 删除操作的必要参数 |
| `R` | READ_SUPPORT | 查询操作的输出字段 |
| `F` | FILTER_SUPPORT | 字段可作为 List 操作的过滤入参 |
| `L` | LIST_OUTPUT_SUPPORT | 字段在 List 操作的输出中出现 |
| `S` | SERVER_GENERATED_DEFAULT | 服务端提供默认值（设置 `hasServerDefaultValue: true`） |
| `A` | CREATE_AFTER_OUTPUT | 创建后额外返回的字段 |
| `!X` | 否定标记 | 显式排除某个操作（如 `#[!C]` 表示排除 Create） |

### F 和 L 的区别

- **F (FILTER_SUPPORT)**：字段可作为 List 操作的过滤入参
- **L (LIST_OUTPUT_SUPPORT)**：字段在 List 操作的输出中出现
- 如果字段没有 `L` 标记但有 `R` 标记，默认也会在 List 输出中出现

## 常用 Flag 组合

| 组合 | 含义 | 典型字段 |
|------|------|---------|
| `#[C,U,R,F,L]` | 用户可编辑字段 | Name、Description、Tags |
| `#[C,R,F,L]` | 创建时设定，不可修改 | Type、InstanceType、ImageId |
| `#[R,F,L]` | 只读系统字段 | Status、CreateTime、UpdateTime |
| `#[R,D,F,L]` | 主键标识字段 | ResourceId（在 identifyDefinition 中） |
| `#[R]` | 仅查询详情返回 | 详细配置信息 |
| `#[C,U,R,F]` | 可编辑但列表不返回 | Priority、详细设置 |
| `#[C]` | 仅创建时传入 | 一次性初始化参数 |
| `#[C,S]` | 创建时传入 + 服务端默认值 | 检索条件字段 |
| `#[!C]` | 排除创建操作 | UpdateTime 等服务端生成字段 |
| 无 flag | 无 CRUD 映射 | RegionId（平台特殊处理） |

### F 和 L 的区分使用

- **带 F 不带 L**：字段可作为 List 过滤条件，但不作为 List 输出字段（如搜索关键字）
- **带 L 不带 F**：字段在 List 输出中出现，但不能作为过滤条件（如计算字段）
- **同时带 F 和 L**：字段既可作为过滤条件，也在 List 输出中出现（如 Name、Status）
- **都不带但有 R**：字段在 Get 输出中出现，List 输出中也会默认出现

## 默认 Flag

系统根据资源字段的分类自动分配默认 flag，无需显式标记：

| 资源分类 | 默认 Flag | 说明 |
|---------|----------|------|
| normal_id | `#[!C,!U,D,R,F,S]` | 主键 ID，不可创建修改，可删除，可读，可过滤，服务端生成 |
| normal_props | `#[C,U,!D,R,F,!S]` | 普通属性，可创建修改，不可删除，可读，可过滤，无服务端默认值 |
| readonly_id | `#[!C,!U,!D,R,F,!S]` | 只读 ID，不可创建修改删除，可读，可过滤 |
| readonly_props | `#[!C,!U,!D,R,F,!S]` | 只读属性，不可创建修改删除，可读，可过滤 |
| relation_id | `#[C,!U,D,R,F,!S]` | 关联 ID，可创建，不可修改，可删除，可读，可过滤 |
| relation_props | `#[C,U,!D,R,F,!S]` | 关联属性，可创建修改，不可删除，可读，可过滤 |

## Flag 冲突检查

系统会检查以下 flag 冲突：

1. **正反义 flag 不能同时出现**：
   - `C` 和 `!C` 不能同时标记
   - `U` 和 `!U` 不能同时标记
   - `D` 和 `!D` 不能同时标记
   - `R` 和 `!R` 不能同时标记
   - `F` 和 `!F` 不能同时标记
   - `L` 和 `!L` 不能同时标记
   - `S` 和 `!S` 不能同时标记

2. **同一 flag 不能重复**：
   - 同一个字段不能重复标记相同的 flag（如 `#[C,C]`）

## Flag 合并规则

在嵌套结构中，flag 会从父级到子级进行合并：

1. **子级 flag 优先于父级 flag**：如果子级显式标记了某个 flag，则以子级为准
2. **父级未出现的 flag 会补充到子级**：父级有但子级未显式标记的 flag，会自动继承到子级
3. **否定标记会覆盖肯定标记**：如果父级有 `C`，子级标记 `!C`，则以子级的否定为准

示例：

```cspec
@flagMode
resource Widget {
  properties: {
    // 父级标记了 C,U,R,F,L
    Config: struct<{
      // 子级显式标记了 R，其他 flag 从父级继承
      Key: string #[R]
      // 子级未标记任何 flag，完全继承父级的 C,U,R,F,L
      Value: string
    }> #[C,U,R,F,L]
  }
}
```

## Flag 到操作的映射规则

### Create 操作

- **输入**：所有带 `C` flag 的 properties 字段
  - identifyDefinition 中的 ID 字段跳过（服务端生成）
  - `@required` 的字段在 Create 输入中也是必填
  - **ROA 风格参数位置**：`body`
  - **RPC 风格参数位置**：`formData`
- **输出**：
  - identifyDefinition 中的 ID 字段
  - 带 `A` (CREATE_AFTER_OUTPUT) 标记的字段
  - `RequestId`

### Get 操作

- **输入**：identifyDefinition 中的 ID 字段（`required: true`）
  - **ROA 风格参数位置**：`query`
  - **RPC 风格参数位置**：`formData`
- **输出**：所有带 `R` flag 的字段 + `RequestId`

### Update 操作

- **部分更新语义**：非 ID 属性一律非必填（`@required` 仅表示创建时必填）
- **输入**：
  - identifyDefinition 中的 ID 字段（`required: true`）
  - 所有带 `U` flag 的 properties 字段（`required: false`）
  - **ROA 风格参数位置**：
    - ID 字段：`query`
    - 属性字段：`body`
  - **RPC 风格参数位置**：
    - ID 字段：`query`
    - 属性字段：`formData`
- **输出**：`RequestId`

### Delete 操作

- **输入**：identifyDefinition 中的 ID 字段（`required: true`）
  - **ROA 风格参数位置**：`query`
  - **RPC 风格参数位置**：`formData`
- **输出**：`RequestId`

### List 操作

- **输入**：
  - 所有带 `F` (FILTER_SUPPORT) flag 的字段（可选筛选条件）
  - 自动生成分页参数
  - **ROA 风格参数位置**：`query`
  - **RPC 风格参数位置**：`query`
- **输出**：
  - `{ResourceName}s` 数组（包含所有带 `R` flag 的字段，或带 `L` flag 的字段）
  - **ROA 风格分页字段**：`NextToken`、`TotalCount`
  - **RPC 风格分页字段**：`PageNumber`、`PageSize`、`TotalCount`
  - `RequestId`

### 操作私有属性

- 通过 `@rac({ operatePrivateType: [...] })` 标记的私有属性
- **始终放在 `query` 参数位置**（无论 ROA 还是 RPC 风格）

## RPC 与 ROA 风格差异

### 参数位置差异

| 操作 | 参数类型 | ROA 风格位置 | RPC 风格位置 |
|-----|---------|-------------|-------------|
| Create | 属性字段 | `body` | `formData` |
| Update | ID 字段 | `query` | `query` |
| Update | 属性字段 | `body` | `formData` |
| Get | ID 字段 | `query` | `formData` |
| Delete | ID 字段 | `query` | `formData` |
| List | 过滤字段 | `query` | `query` |
| 操作私有属性 | 私有字段 | `query` | `query` |

### List 输出差异

- **ROA 风格**：`NextToken` 分页模式，输出字段为 `NextToken`、`TotalCount`
- **RPC 风格**：页码分页模式，输出字段为 `PageNumber`、`PageSize`、`TotalCount`

### 注解透传差异

- **RPC 风格不透传**：`readonly` 和 `clientProhibited` 注解不会透传到生成的参数中
- **ROA 风格**：正常透传注解

## 操作注解自动生成

系统会自动为 Flag 模式生成的操作添加以下注解：

1. **Delete 操作**：
   - `riskType: "high"`（默认为 "high"）

2. **其他操作**（Create、Update、Get、List）：
   - `riskType: "none"`（默认为 "none"）

3. **chargeType**：
   - 从 `@resourceBaseInfo.paidType` 推断

4. **visibility**：
   - 默认为 `"Private"`

5. **List 操作**：
   - 自动生成 `@paginated` 注解（支持 NextToken 分页模式）

## 用户覆盖合并

支持通过 `@for` 标记的 struct 覆盖系统生成的 struct：

### 覆盖机制

1. **用户自定义参数优先**：用户在 `@for` struct 中定义的参数会覆盖系统生成的同名参数
2. **系统生成参数作为补充**：系统生成的参数会作为补充添加到最终 struct 中
3. **额外参数支持**：通过 `@notResourceProperty` 标记可以添加额外参数（不在资源属性中的参数）

示例：

```cspec
// 用户自定义的 Create 输入 struct
@for(Widget)
@notResourceProperty({
  ExtraParam: string
})
struct CreateWidgetInput {
  // 覆盖系统生成的 Name 参数
  Name: string #[C,U,R,F,L]
  // 系统会自动补充其他带 C flag 的字段
}
```

### 覆盖规则

- 用户定义的参数会完全覆盖系统生成的同名参数
- 未被覆盖的系统生成参数会自动添加
- `@notResourceProperty` 标记的参数会被视为额外参数，不会从资源属性中生成

## Flag 选择指南

1. **ID 字段**：`#[R,D,F,L]`，放在 `identifyDefinition` 中，标记 `@readonly`
2. **用户可编辑字段**：`#[C,U,R,F,L]`
3. **不可变字段**（创建后不可改）：`#[C,R,F,L]`
4. **系统字段**（Status、CreateTime）：`#[R,F,L]`，标记 `@readonly`
5. **内部字段**（仅详情返回）：`#[R]`
6. **RegionId**：无 flag，标记 `@readonly`
7. **否定标记**：`#[!C]` 仅用于服务端生成字段（如 UpdateTime）
8. **F 和 L 的区分**：
   - 需要作为 List 过滤条件的字段添加 `F`
   - 需要在 List 输出中出现的字段添加 `L`
   - 既可过滤又需要输出的字段同时添加 `F` 和 `L`

## 完整示例

```cspec
@flagMode
resource Widget {
  identifyDefinition: {
    @readonly
    WidgetId: string #[R,D,L]
  }
  properties: {
    @readonly
    RegionId: string
    @required
    @length({ min: 2, max: 128 })
    Name: string #[C,U,R,L]
    @default("")
    Description: string #[C,U,R,L]
    @required
    Type: string #[C,R,L]
    Priority: int32 #[C,U,R]
    @readonly
    Status: string #[R,L]
    @readonly
    CreateTime: string #[R,L]
    @readonly
    @clientProhibited
    UpdateTime: string #[!C]
  }
}
```

## @for 注解

Flag 模式下的操作使用 `@for(ResourceName)` 关联资源。操作体内只需声明 `errors`，不能自定义 `input`/`output`：

```cspec
@for(Widget)
@document({
  name: "创建小组件"
})
operation CreateWidget {
  errors: [Error_CreateWidget]
}
```

### @for 的限制

- Flag 模式下**不支持**在 operation 中自定义 `input`/`output`
- 以下写法会触发校验错误：

```cspec
// 错误写法
@for(Widget)
operation CreateWidget {
  input: CreateWidgetInput    // 不允许
  output: CreateWidgetOutput  // 不允许
}
```

### @for 的覆盖能力

可以覆盖操作的注解（后端配置、鉴权、错误映射等），但不能覆盖 input/output：

```cspec
@for(Widget)
@backendConfigurationHttp({...})
@http({...})
@errorMapping({...})
operation CreateWidget {
  errors: [Error_CreateWidget]
}
```

未使用 `@for` 覆盖的操作仍保持 Flag 模式默认推导。

## 注意事项

- `@flagMode` 是资源级注解，放在 `@document` 之后、`@resourceBaseInfo` 之前
- 启用 `@flagMode` 后，所有属性（`RegionId` 除外）都**必须**有 flag 标记
- `@readonly` 和 `#[C,U]` 语义冲突，避免同时使用
- 额外入参不能通过 operation 自定义 input struct 注入，应建模在资源属性层或使用 `@notResourceProperty`
- flag 标记放在类型（或 `>` 闭合符号）之后，如 `Name: string #[C,U,R,F,L]`
- 正反义 flag 不能同时出现（如 `C` 和 `!C`），同一 flag 不能重复标记
- Update 操作中，非 ID 属性一律非必填（`@required` 仅表示创建时必填）
- 操作私有属性通过 `@rac({ operatePrivateType: [...] })` 标记，始终放在 `query` 参数位置
- RPC 风格不透传 `readonly` 和 `clientProhibited` 注解
- List 操作中，带 `R` flag 但不带 `L` flag 的字段默认也会在输出中出现
- A (CREATE_AFTER_OUTPUT) flag 仅用于标记创建后需要额外返回的字段

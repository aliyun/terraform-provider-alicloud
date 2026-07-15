# CloudSpec Flag Mode 参考

## 概述

Flag模式（`@flagMode`）是一种声明式的属性映射模式。启用后，资源属性通过flag标记声明自己参与哪些CRUD操作，系统会自动为操作生成对应的Input/Output结构体。

## Flag语法

Flag标记紧跟在字段类型之后：

```cspec
FieldName: FieldType #[C,U,R,L]
```

## Flag定义

| Flag | 全称 | 含义 |
|------|------|------|
| `C` | CREATE_SUPPORT | 创建操作的输入参数 |
| `U` | MODIFY_SUPPORT | 修改操作的输入参数 |
| `R` | READ_SUPPORT | 查询操作的输出字段 |
| `D` | DELETE_SUPPORT | 删除操作的必要参数 |
| `L` | FILTER_SUPPORT | 列表操作中返回/可筛选 |
| `S` | SERVER_GENERATED_DEFAULT | 服务端提供默认值（设置`hasServerDefaultValue: true`） |
| `!X` | 否定标记 | 显式排除某个操作（如`#[!C]`表示排除Create） |

## 常用Flag组合

| 组合 | 含义 | 典型字段 |
|------|------|---------|
| `#[C,U,R,L]` | 用户可编辑字段 | Name、Description、Tags |
| `#[C,R,L]` | 创建时设定，不可修改 | Type、InstanceType、ImageId |
| `#[R,L]` | 只读系统字段 | Status、CreateTime、UpdateTime |
| `#[R,D,L]` | 主键标识字段 | ResourceId（在identifyDefinition中） |
| `#[R]` | 仅查询详情返回 | 详细配置信息 |
| `#[C,U,R]` | 可编辑但列表不返回 | Priority、详细设置 |
| `#[C]` | 仅创建时传入 | 一次性初始化参数 |
| `#[!C]` | 排除创建操作 | UpdateTime等服务端生成字段 |
| 无flag | 无CRUD映射 | RegionId（平台特殊处理） |

## Flag到操作的映射规则

### Create操作

- **输入**：所有带`C` flag的properties字段
  - identifyDefinition中的ID字段跳过（服务端生成）
  - `@required`的字段在Create输入中也是必填
  - 参数位置：`body`
- **输出**：identifyDefinition中的ID字段 + `RequestId`

### Get操作

- **输入**：identifyDefinition中的ID字段（`required: true`，`in: query`）
- **输出**：所有带`R` flag的字段 + `RequestId`

### Update操作

- **输入**：
  - identifyDefinition中的ID字段（`required: true`，`in: query`）
  - 所有带`U` flag的properties字段（`required: false`，`in: body`）
- **输出**：`RequestId`

### Delete操作

- **输入**：identifyDefinition中的ID字段（`required: true`，`in: query`）
- **输出**：`RequestId`

### List操作

- **输入**：
  - 所有带`L` flag的字段（可选筛选条件）
  - 自动生成分页参数
- **输出**：
  - `{ResourceName}s`数组（包含所有带`R` flag的字段）
  - 分页字段（`NextToken`、`TotalCount`）
  - `RequestId`

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

## Flag选择指南

1. **ID字段**：`#[R,D,L]`，放在`identifyDefinition`中，标记`@readonly`
2. **用户可编辑字段**：`#[C,U,R,L]`
3. **不可变字段**（创建后不可改）：`#[C,R,L]`
4. **系统字段**（Status、CreateTime）：`#[R,L]`，标记`@readonly`
5. **内部字段**（仅详情返回）：`#[R]`
6. **RegionId**：无flag，标记`@readonly`
7. **否定标记**：`#[!C]`仅用于服务端生成字段（如UpdateTime）

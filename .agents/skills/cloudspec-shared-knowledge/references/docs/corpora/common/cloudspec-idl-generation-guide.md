# CloudSpec IDL 生成指南

## 1. 概述

本指南旨在为开发人员提供一套标准化的流程和规范，用于根据产品需求文档（PRD）生成符合 CloudSpec 规范的 IDL 文件。CloudSpec（Cloud Specification）是一种用于统一描述云产品服务、资源、OpenAPI 及配套能力的模型规范。

## 2. 基本概念

### 2.1 核心组件

1. **Namespace（命名空间）**：服务、资源、操作的唯一标识空间，格式为 `alicloud.{PRODUCT}.{POP-CODE}.v{POP-VERSION}`
2. **Service（服务）**：云产品提供的资源和操作集合
3. **Resource（资源）**：服务下能被自动化管理的长期存在对象
4. **Operation（操作）**：资源的生命周期接口或相关操作接口（OpenAPI）
5. **Error（错误码）**：操作出错时的信息定义
6. **Annotate（元注解）**：为模型元素添加额外信息表达的注解

### 2.2 模型类型

| 类型 | 说明             |
|------|----------------|
| string | 字符串，不限制长度      |
| byte | 经过base64编码的字节流 |
| binary | 字节流            |
| boolean | 布尔值            |
| int32 | 32位有符号整数       |
| int64 | 64位有符号整数       |
| float | 单精度浮点数         |
| double | 双精度浮点数         |
| enum | 枚举值（string类型）  |
| intEnum | 数值型枚举（int32类型） |
| struct | 结构体            |
| array | 数组结构           |
| map | map类型          |
| service | 服务类型           |
| resource | 资源类型           |
| operation | 操作类型           |

## 3. PRD文档分析流程

### 3.1 资源识别

从PRD文档中识别出以下信息：
1. **资源名称**：明确的资源实体名称
2. **资源主键**：唯一标识资源的属性
3. **资源属性**：资源的所有可配置属性
4. **生命周期操作**：创建、获取、更新、删除、列举操作
5. **关联资源**：与其他资源的依赖关系

### 3.2 操作识别

从PRD文档中识别出以下信息：
1. **操作名称**：API接口名称
2. **HTTP方法**：GET、POST、PUT、DELETE等
3. **URI路径**：API的访问路径
4. **请求参数**：入参定义
5. **响应参数**：出参定义
6. **错误码**：可能的错误情况

## 4. CloudSpec文件结构

### 4.1 项目结构

```
project/
├── main.cspec                 # 主文件
├── resources/                 # 资源定义目录
│   ├── ResourceName.cspec     # 资源定义文件
│   └── ResourceMappings/      # 资源映射目录
├── operations/                # 操作定义目录
│   └── OperationName.cspec    # 操作定义文件
└── openStructs/               # 开放结构体目录
    └── StructName.cspec       # 结构体定义文件
```

### 4.2 文件命名规范

1. **资源文件**：使用大驼峰命名法，如 `Function.cspec`
2. **操作文件**：使用大驼峰命名法，如 `CreateFunction.cspec`
3. **结构体文件**：使用大驼峰命名法，如 `FunctionProperties.cspec`

## 5. 资源定义规范

### 5.1 基本结构

```json
$version: 1
namespace: alicloud.PRODUCT.CODE.v20230330

@document({
  zh: "资源中文描述"
  en: "Resource English Description"
})
resource ResourceName {
  identifyDefinition: {
    @document({
      name: "主键名称"
      zh: "主键中文描述"
      en: "Primary Key English Description"
    })
    @required
    PrimaryKeyId: string
  }
  properties: ResourceProperties
  create: CreateResourceOperation
  get: GetResourceOperation
  update: [UpdateResourceOperation]
  delete: DeleteResourceOperation
  list: ListResourceOperation
}
```

### 5.2 属性定义

1. **主键定义**：必须在 `identifyDefinition` 中定义，使用 `@required` 注解
2. **属性定义**：在 `properties` 结构体中定义所有非主键属性
3. **文档注解**：使用 `@document` 为每个属性添加中英文描述

### 5.3 常用注解

1. **@required**：标记必填属性
2. **@document**：添加文档描述
3. **@references**：声明资源间关系
4. **@resourceBaseInfo**：资源基础信息

## 6. 操作定义规范

### 6.1 基本结构

```json
@document({
  name: "操作名称"
  zh: "操作中文描述"
})
@http({
  methods: ["post"]
  uri: "/2023-03-30/resources"
})
operation CreateResource {
  input: CreateResourceInput
  output: CreateResourceOutput
}
```

### 6.2 HTTP注解

1. **methods**：HTTP请求方法
2. **uri**：API访问路径
3. **authenticators**：认证方式
4. **requestContentType**：请求内容类型
5. **responseContentType**：响应内容类型

### 6.3 入参出参定义

1. **输入结构体**：定义所有请求参数
2. **输出结构体**：定义所有响应参数
3. **引用资源属性**：使用 `$PropertyName` 引用资源属性

## 7. 映射关系定义

### 7.1 操作映射结构

```json
@defineOperationMapping
struct Resource_create_CreateResource_mapping {
  rootMapping: {
    responsePathType: "jsonPath"
    responsePath: "$"
  }
  resourceAttributeMappings: [{
    resourceProperty: "$.PropertyName"
    requestMappingType: "param"
    requestPathType: "normal"
    requestPath: "body.propertyName"
    responsePathType: "normal"
    responsePath: "$.propertyName"
  }]
}
```

### 7.2 映射字段说明

1. **resourceProperty**：资源属性路径
2. **requestMappingType**：请求映射类型（param/api）
3. **requestPath**：请求参数路径
4. **responsePath**：响应参数路径

## 8. 生成步骤

### 8.1 资源文件生成

1. **创建资源定义**：
   - 确定资源名称和主键
   - 定义资源属性结构体
   - 关联CRUDL操作

2. **添加必要注解**：
   - @document 添加文档描述
   - @references 声明资源关系
   - @resourceBaseInfo 添加资源基础信息

### 8.2 操作文件生成

1. **创建操作定义**：
   - 定义操作名称和HTTP信息
   - 创建入参和出参结构体
   - 关联错误码定义

2. **添加操作注解**：
   - @http 配置HTTP信息
   - @document 添加操作描述
   - @ram 配置RAM权限

### 8.3 映射文件生成

1. **创建映射结构**：
   - 使用 @defineOperationMapping 定义映射
   - 配置 rootMapping 根映射
   - 定义 resourceAttributeMappings 属性映射

2. **配置映射关系**：
   - 映射请求参数到资源属性
   - 映射响应参数到资源属性
   - 配置异步操作和重试策略

## 9. 验证流程

### 9.1 语法检查

```bash
# 语法检查
aliyun cspec check

```

### 9.2 构建验证

```bash
# 构建验证
aliyun cspec build

# 检查破坏性变更
aliyun cspec bc
```

### 9.3 功能验证

```bash
# 生成Terraform代码
aliyun cspec terraform

# 生成示例代码
aliyun cspec sample
```

## 10. 最佳实践

### 10.1 命名规范

1. **一致性**：保持命名风格一致
2. **清晰性**：名称应清晰表达用途
3. **简洁性**：避免冗余和过长的名称

### 10.2 文档规范

1. **完整性**：为所有组件添加文档注解
2. **准确性**：确保文档描述与功能一致
3. **双语支持**：提供中英文文档

### 10.3 结构规范

1. **模块化**：合理拆分文件结构
2. **可维护性**：保持结构清晰易维护
3. **可扩展性**：预留扩展空间

## 11. 常见问题

### 11.1 映射关系问题

1. **属性不匹配**：确保资源属性与API参数正确映射
2. **路径错误**：检查JSON路径表达式是否正确
3. **类型不一致**：确保数据类型匹配

### 11.2 注解使用问题

1. **注解遗漏**：确保必要的注解都已添加
2. **注解错误**：检查注解参数是否正确
3. **注解冲突**：避免注解间的冲突

### 11.3 验证问题

1. **语法错误**：检查IDL语法是否正确
2. **规范不符**：确保符合CloudSpec规范
3. **映射缺失**：检查操作映射是否完整

## 12. 附录

### 12.1 常用注解列表

| 注解 | 用途 | 适用对象 |
|------|------|----------|
| @document | 添加文档描述 | 所有组件 |
| @required | 标记必填属性 | struct成员 |
| @references | 声明资源关系 | resource |
| @http | 配置HTTP信息 | operation |
| @ram | 配置RAM权限 | operation |
| @defineOperationMapping | 定义操作映射 | struct |

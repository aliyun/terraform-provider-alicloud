# CloudSpec 字段级注解参考

## 传值约束注解

### @required

标记字段为必填。

```cspec
@required
Name: string #[C,U,R,L]
```

### @readonly

标记字段只读。作用在不同component时含义不同：
- 修饰字段：该字段不可被修改（只能在创建阶段设置）
- 修饰operation：该操作不允许修改任何资源属性

```cspec
@readonly
Status: string #[R,L]
```

### @clientOptional

客户端可选。与`@required`配合使用时，表示该字段服务端会生成值，客户端可以不传。

```cspec
@required
@clientOptional
ServerGeneratedField: string #[C,R,L]
```

### @clientProhibited

禁止客户端传值。与`@readonly`配合使用时，表示该值只能由服务端生成且不可修改。

```cspec
@readonly
@clientProhibited
InternalField: string #[R]
```

### @conflictsWith([...])

声明互斥字段，传值时不能同时传递。

```cspec
@conflictsWith(['Id'])
Name: string #[C,U,R,L]
@conflictsWith(['Name'])
Id: string #[C,R,L]
```

### @idempotencyToken

标记幂等参数。

```cspec
@idempotencyToken
ClientToken: string #[C]
```

### @hasDefaultValue

标记服务端存在默认值。客户端不传值时，服务端会生成默认值。

```cspec
@readonly
@hasDefaultValue
CreateTime: string #[R,L]
```

## 值约束注解

### @default(value)

设置默认值。

```cspec
@default("")
Description: string #[C,U,R,L]

@default(10)
PageSize: int32 #[C,U,R]

@default(true)
Enabled: boolean #[C,U,R,L]
```

### @length({ min, max })

字符串长度约束。

```cspec
@length({ min: 2, max: 128 })
Name: string #[C,U,R,L]
```

### @range({ min, max })

数值范围约束。

```cspec
@range({ min: 1, max: 100 })
Priority: int32 #[C,U,R]
```

### @regexPattern("pattern")

正则校验。

```cspec
@regexPattern("^[a-zA-Z][a-zA-Z0-9_-]*$")
Name: string #[C,U,R,L]
```

### @enums(EnumRef)

关联枚举定义。

```cspec
@enums(StatusEnum)
Status: string #[R,L]

enum StatusEnum {
  "Active"
  "Inactive"
  "Creating"
}
```

### @initialStatus("value")

声明Status字段的初始状态值。

```cspec
@enums(StatusEnum)
@initialStatus("Pending")
Status: string #[R,L]
```

## 文档注解

### @document({...})

为字段添加文档信息。

```cspec
@document({
  name: "实例名称"
  nameEn: "Instance Name"
  zh: "ECS实例的名称"
  en: "Name of the ECS instance"
  exampleValue: "my-instance"
})
Name: string #[C,U,R,L]
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `name` | string | 中文名称 |
| `nameEn` | string | 英文名称 |
| `zh` | string | 中文描述 |
| `en` | string | 英文描述 |
| `exampleValue` | string | 示例值 |

### @sensitive

标记敏感字段（如密码）。

```cspec
@sensitive
Password: string #[C]
```

### @deprecated

标记废弃字段。

```cspec
@deprecated
OldField: string #[R]

@deprecated({
  message: "Use NewField instead"
  substitute: [NewField]
})
OldField: string #[R]
```

## 格式注解

### @format("format")

声明值的格式。

```cspec
@format("iso8601")
CreateTime: string #[R,L]

@format("json")
ConfigJson: string #[C,U,R]
```

### @date({ format })

声明时间格式。

```cspec
@date({
  format: "YYYY-MM-DDTHH:mm:ss.sssZ"
})
CreateTime: string #[R,L]
```

## 数组相关注解

### @arrayConfig({...})

数组行为配置。

```cspec
@arrayConfig({
  uniqueItems: true
  unordered: true
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `uniqueItems` | boolean | 数组元素是否唯一 |
| `unordered` | boolean | 是否无序数组 |

## 参数映射注解

### @parameter({...})

控制API参数行为（主要用于operation的input/output结构体中的字段）。

```cspec
@parameter({
  in: "body"
  style: "flat"
  required: false
  deprecated: false
  sensitive: false
  readOnly: false
  hasServerDefaultValue: false
  maxLength: "256"
  minLength: 1
  minItems: 1
})
FieldName: string
```

### @resourceProperty('PropertyName')

当操作的出入参字段名与资源属性名不一致时，声明映射关系。

```cspec
@resourceProperty('HouseId')
HouseId2: string
```

### @notResourceProperty

声明操作出入参中的字段不是资源属性。

```cspec
@notResourceProperty
RequestId: string
```

## 资源属性配置注解

### @rac({...})

资源属性的额外配置。

```cspec
@rac({
  filterProperty: true
  operatePrivateType: ["create", "update"]
  writeOnly: true
  hasServerDefaultValue: true
})
```

| 属性 | 类型 | 说明 |
|-----|------|------|
| `filterProperty` | boolean | 是否为过滤条件 |
| `processVariable` | boolean | 是否为过程变量 |
| `operatePrivateType` | [string] | 操作私有属性类型 |
| `writeOnly` | boolean | 是否只写（不出现在查询返回中） |
| `hasServerDefaultValue` | boolean | 是否有服务端默认值 |
| `constType` | boolean | 动态资源标识 |
| `createMaxItems` | integer | Create操作的数组元素最大个数 |

### @conditionalDependency({...})

资源属性的条件依赖配置。

```cspec
@conditionalDependency({
  mustExist: ["Name"]
})
Id: string #[C,R,L]
```

详细结构参见resource-annotate.md中的说明。

### @secondUniqueKey(order)

标记备选主键。

```cspec
@secondUniqueKey(1)
Name: string #[C,R,L]
```

## 其他注解

### @hidden

隐藏结构体或字段（不在公开文档中展示）。

```cspec
@hidden
struct InternalStruct { ... }
```

### @openStruct

声明数据结构对客户公开。

```cspec
@openStruct
struct InstanceDetail { ... }
```

### @for(ResourceName)

当struct中引用资源属性（`$PropertyName`语法）时，需要在struct上声明引用的资源。

```cspec
@for(Instance)
struct GetInstanceOutput {
  $Name
  $Status
  @notResourceProperty
  RequestId: string
}
```

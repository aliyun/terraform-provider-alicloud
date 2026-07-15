# CloudSpec 自动化生成注解列表

本文档列出了 CloudSpec 自动化生成相关的所有注解及其详细配置说明，可以指定由模型自动生成和资源相关的出入参。

## 目录

- [A12.1 autoGenerateOperations annotate](#a121-autogenerateoperations-annotate)
- [A12.2 autoGenerateResource annotate](#a122-autogenerateresource-annotate)
- [A12.3 autoGenerateResourceTest annotate](#a123-autogenerateresourcetest-annotate)

### A12.1 autoGenerateOperations annotate
#### 说明
由模型按规范自动生成OpenAPI的定义。

#### 选择器
```cspec
:is(resource)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无 | | |

#### 示例
```cspec
$version: 1
namespace: alicloud.ECS.Ecs.v20140526

@apiStyle("rpc")
@ram({
  enable: true
  ramCodes: [{
    popCode: 'Ecs'
    codes: ['ecs']
  }]
})
service Ecs {
  version: "2014-05-26"
  resources: [Instance]
}

@autoGenerateOperations
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "Region"
  getRegionIdByEndpoint: true
})
resource Instance {
  identifyDefinition: resourceInstanceIdentifyDefinition
  properties: resourceInstanceProperties
}

struct resourceInstanceIdentifyDefinition {
  @required
  Id: string
}

struct resourceInstanceProperties {
  Name: string
  Name1: string
  RegionId: string
  @required
  ComplexStruct: {
    ComplexA: string
    ComplexB: int32
    ComplexC: ComplexCStruct
  }
  ComplexArray: ComplexArrayDefine
  ComplexMap: ComplexMapDefine
  Tags: TagsArrayDefine
}

array TagsArrayDefine {
  item: {
    TagKey: string
    TagValue: string
  }
}

struct ComplexCStruct {
  ComplexCA: string
  ComplexCB: int32
}

array ComplexArrayDefine {
  item: {
    ComplexA: string
    ComplexB: int32
    ComplexC: float
  }
}

map ComplexMapDefine {
  key: string
  value: {
    ComplexA: string
    ComplexB: int32
    ComplexC: float
  }
}
```

指定后，模型将自动生成CreateFloor等OpenAPI的定义，适用于首次生产。CLI已经提供了auto指令帮助你快速完成operation的生成。

### A12.2 autoGenerateResource annotate
#### 说明
由模型按规范自动生成资源的定义。

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| resourceName | string | 待生成的资源名称 |
| type | string | 资源生命周期的分类，支持枚举：create、update、delete、get、list |
| isAssociate | boolean | true 表示不是资源的编排API；false 表示是资源的编排API。默认值为false |

### A12.3 autoGenerateResourceTest annotate
#### 说明
自动生成资源的测试用例。

#### 选择器
```cspec
:is(resource)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| associateApi | boolean | true 生成放置在资源下的非编排类的 API 的用例；false 不生成放置在资源下的非编排类的 API 的用例。默认值为true，表示非编排类的 API 也生成测试用例。 |

关于associateApi的具体说明如下面的例子：

```cspec
$version: 1
namespace: a.b

resource A {
  create: CreateA
  get: GetA
  operations: [AttachA]
}
```

当associateApi = true 时，测试用例自动生成会给 operations 下的operation 生成测试用例。当为 false 时，只生成资源本身的测试用例。

# CloudSpec 企业级能力注解列表

本文档列出了 CloudSpec 企业级能力相关的所有注解及其详细配置说明，可以给 service、resource、operation 配置接入企业级能力的额外信息。

## 目录

- [A8.1 terraform annotate](#a81-terraform-annotate)
- [A8.2 rmc annotate](#a82-rmc-annotate)
- [A8.3 tagService annotate](#a83-tagservice-annotate)
- [A8.4 resourceGroup annotate](#a84-resourcegroup-annotate)
- [A8.5 actionTrail annotate](#a85-actiontrail-annotate)
- [A8.6 ros annotate](#a86-ros-annotate)
- [A8.7 config annotate](#a87-config-annotate)
- [A8.8 ccApi annotate](#a88-ccapi-annotate)
- [A8.9 sdk annotate](#a89-sdk-annotate)
- [A8.10 apiDocument annotate](#a810-apidocument-annotate)
- [A8.11 rmGw annotate](#a811-rmgw-annotate)

### A8.1 terraform annotate
#### 说明
配置资源/服务接入到Terraform服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
作用于资源生效字段：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到Terraform。 |
| resourceName | string | 接入到Terraform中的资源名称。 |
| dataSourceName | string | 接入到Terraform中的数据源名称。 |
| reason | string | 当不开启时的原因。 |
| unnecessary | boolean | 是否无必要接入。 |

作用于service生效字段：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到Terraform。 |
| code | string | 服务在Terraform中的code。 |
| reason | string | 当不开启时的原因。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
作用于资源示例：

```cspec
@terraform({
  enable: true,
  resourceName: 'dessert_candy'
  dataSourceName: 'dessert_candies'
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到Terraform中，且直接定义了接入到Terraform的资源和数据源的名称。

作用于service示例：

```cspec
@terraform({
  enable: true
  code: "B"
  releaseStrategy: {
    mode: "delay"
    delay: 30
    scope: "all"
  }
})
service B {}
```

### A8.2 rmc annotate
#### 说明
配置资源/服务接入到RMC服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到RMC。 |
| name | string | 接入到RMC中的资源名称，默认和资源名称一致，不需要改写。 |
| reason | string | 当不开启时的原因。 |
| unnecessary | boolean | 是否无必要接入。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在RMC中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
作用于资源示例：

```cspec
@rmc({
  enable: true,
  name: 'CandyOther'
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到RMC中，且直接定义了接入到RMC的资源名称。

作用于service示例：

```cspec
@rmc({
  enable: true
  code: "B"
  releaseStrategy: {
    mode: "delay"
    delay: 7
    scope: "all"
  }
})
service B {}
```

### A8.3 tagService annotate
#### 说明
配置资源/服务接入到标签服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到标签服务。 |
| name | string | 接入到标签服务中的资源名称，默认和资源名称一致，不需要改写。 |
| reason | string | 当enable为false时，表示不接入的原因。 |
| unnecessary | boolean | 是否无必要接入。 |
| accessMethod | string | 接入方式，支持：gateway（网关接入）、code（自定义接入）、hybrid（混合） |
| innerPopCode | string | 内部POP Code。 |
| innerPopVersion | string | 内部POP版本。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在标签服务中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@tagService({
  enable: true
  code: "ecs"
  accessMethod: "hybrid"
  innerPopCode: "Ecs-inc"
  innerPopVersion: "2020-02-02"
  releaseStrategy: {
    mode: "synchronous"
    scope: "all"
  }
})
service Ecs {}
```

### A8.4 resourceGroup annotate
#### 说明
配置资源/服务接入到资源组服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到资源组。 |
| name | string | 接入到资源组中的资源名称，默认和资源名称一致，不需要改写。 |
| enableAutoTest | boolean | 是否开启自动化测试。 |
| bizMode | string | 业务模式。 |
| accessMode | string | 接入模式。 |
| reason | string | 当enable为false时，表示不接入的原因。 |
| unnecessary | boolean | 是否无必要接入。 |
| arnResourceType | string | 资源 ARN。 |
| odpsTableName | string | ODPS中对账的表名称。 |
| odpsTaskName | string | ODPS中的任务名称。 |
| accessMethod | string | 接入方式，支持：gateway（网关接入）、decentralization（自定义接入）、hybrid（混合） |
| needSplitBill | boolean | 是否分账。 |
| needSplitBillInFuture | boolean | 未来是否分账。 |
| odpsNodeOutputName | string | 资源-rg关系odps表节点输出名称 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在资源组中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@resourceGroup({
  enable: true
  name: 'CandyOther'
  enableAutoTest: false
  bizMode: 'NoLimit'
  accessMode: 'RmGw'
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到资源组服务中，且定义了接入资源网关的一些配置。

### A8.5 actionTrail annotate
#### 说明
配置资源/服务接入到ActionTrail服务的配置。

注意：当不配置该annotate时，operation默认行为为打开。

#### 选择器
```cspec
:is(resource, operation, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到操作审计。required |
| name | string | 接入到操作审计中的资源名称，默认和资源名称一致，不需要改写。 |
| reason | string | 当不开启时的原因。 |
| unnecessary | boolean | 是否无必要接入。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在ActionTrail中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@actionTrail({
  enable: true
  name: 'CandyOther'
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到操作审计服务中，且直接定义了接入到操作审计服务的资源名称。

### A8.6 ros annotate
#### 说明
配置资源/服务接入到ROS服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到ROS。 |
| reason | string | 当enable为false时，表示不接入的原因。 |
| name | string | 接入到ROS中的资源名称，默认和资源名称一致，不需要改写。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在ROS中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@ros({
  enable: true
  name: 'CandyOther'
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到ROS服务中，且直接定义了接入到ROS服务的资源名称。

### A8.7 config annotate
#### 说明
配置资源/服务接入到配置审计服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到配置审计。 |
| reason | string | 当enable为false时，表示不接入的原因。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在Config中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@config({
  enable: true
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到配置审计服务。

### A8.8 ccApi annotate
#### 说明
配置资源/服务接入到Cloud Control API服务的配置。

#### 选择器
```cspec
:is(resource, service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否接入到CC API。 |
| reason | string | 当enable为false时，表示不接入的原因。 |
| unnecessary | boolean | 是否无必要接入。 |

以下属性当注解作用于service时生效：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| code | string | 服务在CC API中的code。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@ccApi({
  enable: true
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到CC API服务。

### A8.9 sdk annotate
#### 说明
配置OpenAPI参数在SDK中的特性，服务的SDK特性。

#### 选择器
```cspec
:is(struct > member, service)
```

#### 支持的属性
作用于API参数生效字段：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| name | string | 替换属性在SDK中的名称。 |
| ignore | boolean | SDK是否隐藏该object类型的根字段，直接生成子属性，默认为false。 |
| extendType | string | 旧版SDK JSON字段类型。 |

作用于service生效字段：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否开启。 |
| code | string | 服务在SDK中的code。 |
| reason | string | 当不开启时的原因。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
作用于参数示例：

```cspec
struct A {
  @sdk({
    name: "cpuUsageOther"
  })
  CpuUsage: string
}
```

作用于服务示例：

```cspec
@sdk({
  enable: true
  code: "otherCode"
  releaseStrategy: {
    mode: "synchronous"
    scope: "all"
  }
})
service ECS {}
```

### A8.10 apiDocument annotate
#### 说明
配置服务API文档。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否开启。 |
| reason | string | 当不开启时的原因。 |
| releaseStrategy | releaseStrategyStruct | 发布的策略配置。 |

#### 示例
```cspec
@apiDocument({
  enable: true
  releaseStrategy: {
    mode: "synchronous"
    scope: "all"
  }
})
service ECS {}
```

### A8.11 rmGw annotate
#### 说明
配置资源接入到资源网关，主要用于网关资源级鉴权。

#### 选择器
```cspec
:is(resource)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | true表示接入。默认情况下，当不存在这个annotate时，与false等价 |
| serviceAccount | string | 云产品服务账号 |
| odpsTableName | string | 资源ODPS表名 |
| unnecessary | boolean | 是否无必要接入。 |

#### 示例
```cspec
@rmGw({
  enable: true
})
resource Candy {

}
```

以上的示例表示Candy这个资源接入到网关资源级鉴权服务。

### 公共数据结构
#### releaseStrategyStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| mode | string | 发布的模式，支持：synchronous（实时发布）、delay（在API/资源发布后延迟发布）、manual（手工接入） |
| delay | integer | 当mode为delay时，延迟发布的天数。 |
| scope | string | 作用的范围：increment（增量的API/资源）、all（全量的API/资源） |

# CloudSpec 约束器注解列表

本文档列出了 CloudSpec 约束器相关的所有注解及其详细配置说明。

## 目录

- [A13.1 skipConstraint annotate](#a131-skipconstraint-annotate)

### A13.1 skipConstraint annotate
#### 说明
配置部分忽略约束器。

> 特别注意：当规则被设定为不可跳过时，配置skipConstraint会被忽略。

#### 选择器
```cspec
*
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [string] | 在component上配置忽略约束。 |

#### 示例
```cspec
import api.apiyun-inc.com/buildin/constraints/1.0.cspec

@skipConstraint([aliyun.cloudspec#RT0001])
resource Floor {
  identifyDefinition: {
    HouseId: string
    FloorId: string
  }
  properties: FloorProperties
  create: CreateFloor
  update: UpdateFloor
  delete: DeleteFloor
  get: GetFloor
  list: ListFloor
}

struct FloorProperties {
  WindowsNumber: int32
  StoreyHeight: float
  ToiletsNumber: int32
}
```

从api.apiyun-inc.com/buildin/constraints/1.0.cspec中会导入一组约束器，在资源上配置@skipConstraint 可忽略导入的约束器中指定名称的检查。

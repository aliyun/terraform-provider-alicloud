### 概览
![](https://intranetproxy.alipay.com/skylark/lark/0/2023/png/309278/1681303874260-b57763bf-aab3-46f6-b86a-dc645cfcedb6.png)

CloudSpec模型由基本类型、复合类型、服务类型、选择器、约束组成，用Interface Define Language的方式描述云产品的服务、资源、OpenAPI及各配套能力的关系。模型支持自动转换为Resource/RAM/OpenAPI Meta及反转。

### 组成
组成CloudSpec的每个类型称之为Component。

```puml
@startuml

class CloudSpec {
  $version: int32
  services: [service]
  constraints: [constraint]
}
note left of CloudSpec::constraints
  约束定义
end note
object BasicType
object StringType
entity string
enum enum
entity byte
entity binary
object BooleanType
object NumberType
entity int32
enum intEnum
entity int64
entity float
entity double
object AggregateType
class struct {
}
class array {
  item: any
}
class map {
  key: string
  value: any
}
class service {
  version: string
  resources: [resource]
  operations: [operation]
  errors: [error]
}
class resource {
  identifyDefinition: struct
  properties: struct
  get: operation
  list: operation
  update: [operation]
  create: operation
  delete: operation
  batchCreate: operation
  batchDelete: operation
  batchUpdate: operation
  resources: [resource]
}
class operation {
  input: struct
  output: struct
  errors: [error]
}
class error {
  httpCode: int32
  message: string
}

class constraint {
  name: string
  id: string
  notice: string 
  severity: string
  namespaces: [string]
  selector: string
  annotate: [string]
  skipPrelude: boolean
}

class $test {
}

class resourceTest {
	init: struct
	modifies: array<struct>
	destroy: struct
}

class operationTest {
	input: struct
	asserts: array<struct>
}

$test <-down- resourceTest
$test <-down- operationTest


CloudSpec <-down- BasicType
CloudSpec <-down- AggregateType
CloudSpec <-down- service
CloudSpec <-down- resource
CloudSpec <-down- operation
CloudSpec <-down- error
service .. resource
service .. operation
resource .. operation
operation .. error
BasicType <-down- StringType
StringType <-down- string
string <-down- enum
StringType <-down- byte
StringType <-down- binary
BasicType <-down- BooleanType
BasicType <-down- NumberType
NumberType <-down- int32
int32 <-down- intEnum
NumberType <-down- int64
NumberType <-down- float
NumberType <-down- double
AggregateType <-down- struct
AggregateType <-down- array
AggregateType <-down- map
CloudSpec <-down- constraint
CloudSpec <-down- $test

service <|-- constraint: 约束
resource <|-- constraint: 约束
operation <|-- constraint: 约束
error <|-- constraint: 约束
@enduml
```

### Component ID ABNF
```json
ComponentId =
    RootComponentId [ComponentIdMember];

RootComponentId =
    AbsoluteRootComponentId / Identifier;

AbsoluteRootComponentId =
    Namespace "#" Identifier;

Namespace =
    Identifier *("." Identifier);

Identifier =
    IdentifierStart *IdentifierChars;

IdentifierStart =
    (1*"_" (ALPHA / DIGIT)) / ALPHA;

IdentifierChars =
    ALPHA / DIGIT / "_";

ComponentIdMember =
    "$" Identifier;

ALPHA = %x41-5A / %x61-7A;

DIGIT = %x30-39;
```

### Component ID 组成说明
#### 跨namespace的引用表达
```json
$version: 1
namespace: a.b
struct A {
  B: a.c#D
}
```

示例中的B引用的是a.cnamespace下的名称为D的component。component的类型可以是基本类型、复合类型、服务类型。

:::warning
特别注意，如果引用本地的对象，namespace可以省略，但#不可省略。

:::

#### 带属性的引用
```json
$version: 1
namespace: a.b

struct A1 {
  A11: int32
  A12: string
}

struct A {
  B: string
  B1: A1
}

map M {
  key: string
  value: A
}

array ArrayA {
  item: A1
}
```

在a.c中引用的格式如下：

```json
$version: 1
namespace: a.c
// 需要显示的导入namespace为a.b的文件
import "./ab.cspec"

struct AC {
  // 这里a的类型为int32
  a: a.b#A1$A11
  // 这里b的类型和a.b下的map M是一致的，也是一个map类型
  b: a.b#M
  // 这里c的类型和a.b下的map M的value的类型一致，是一个struct A的类型
  c: a.b#M$value 
  // 这里d的类型其实是a.b夏的A1的A11结构
  d: a.b#ArrayA$item/A11
}
```


### CloudSpec模型是什么
CloudSpec（Cloud Specification）是以CloudSpec规范为核心，以镇元平台为实现，标准化生产CloudSpec模型，以模型驱动对客的高质量OpenAPI建设的一套综合体系。CloudSpec模型是用于统一描述云产品服务、资源、OpenAPI及配套能力（如RAM鉴权、操作审计、流控等）的具象表达。

### 基本概念介绍
+ Namespace，命名空间，服务、资源、操作都在唯一的命名空间下，命名空间用于隔离不同的产品的资源和操作，在同一个命名空间下，服务、资源、操作的名称是唯一的，在不同的命名空间下，这些是允许重复的。
+ Service，服务，云产品提供的资源、操作的集合，例如弹性计算提供的ECS是一个服务、提供的计算巢又是另外一个服务，这两个服务的命名空间（namespace）可以一致。
+ Resource，资源，是服务下能被自动化管理的对象。这个对象是一个长期存在的实体。是资源的例子如ECS的主机实例（Instance）、短信的模板（Template），不是资源的例子如短信服务的消息，消息是一次性的对象，这类对象在设计时不需要按资源化的方式去设计OpenAPI。
+ Operation，操作，包含资源完成生命周期的接口及不和资源相关的其他操作接口，表现形式就是OpenAPI。
+ Error，错误码，是操作在出错时的信息，错误包含HTTP Code定义，错误码和错误信息定义。
+ Annotate，元注解，是作用于以上类型的注解信息，给每一个类型增加更多的额外的信息表达。
+ Constraint，约束，通过选择器选择模型下的元素，对这些元素执行判定，检查是否满足定义的要求。

### 模型类型
| 类型 | 说明 |
| --- | --- |
| string | 字符串，不限制长度。 |
| byte | 表示经过base64编码的字节流。 |
| binary | 字节流。 |
| boolean | 布尔值 |
| int32 | 整形，32位有符号整数，-2^31 到 (2^31)-1 。 |
| int64 | 长整形，32位有符号整数，-2^63 到 (2^63)-1 。 |
| float | IEEE-754规定的单精度浮点数。 |
| double | IEEE-754规定的双精度浮点数。 |
| any | 任意类型。 |
| enum | 枚举，枚举的值是string。 |
| intEnum | 数值型枚举，枚举值的类型是int32。 |
| struct | 结构体，key为string类型，value可以不限定。 |
| array | 数组结构，值必须为一种类型。 |
| map | 地图类型，key为string类型，value的类型需要定义出来，且为一种。 |
| service | 服务类型，包含一组资源、操作。 |
| resource | 资源类型，包含资源属性、操作及子资源。 |
| operation | 操作类型（OpenAPI），包含入参、出参、错误码。 |


### 元注解
通过以上类型能够表达全部的模型元素，但如果全部的信息都通过属性表达，整体会比较冗余。举个例子，对于一个string类型需要定义一个枚举值，那么表达方式可以选择：

```json
struct valueString {
  type: string
  regex: '^[a-zA-Z0-9]{2,16}$'
}
```

以上的方式通常出现在以JSON方式表达模型元素中，对于CloudSpec IDL，设计了更精简直观的表达：

```json
@regexPattern('^[a-zA-Z0-9]{2,16}$')
string valueString
```

regexPattern就是一个元注解。

### 甜品工厂服务（DessertFactoryService）
为了更好的说明如何使用CloudSpec模型，我们以甜品工厂为例，这个工厂本身就是一个服务，工厂生产糖果（Candy），蛋糕（Cake），那么对于工厂（Service）而言，糖果和蛋糕都是它的资源（resource）。可以使用如下的方式定义出service：

```json
$version: "1"
namespace example.DessertFactory

service DessertFactoryService {
  resources: [Candy, Cake]
}
```

当然工厂不一定全部的操作都包含在资源中，比如工厂提供了一个接口查询当前实时在运转的车间数量：

```json
operation GetWorkshopsNumberInOperation {}
```

这个GetWorkshopsNumberInOperation操作本身是工厂自身，而不是其资源，模型也是支持这种定义，将其定义在service的operations属性中即可：

```json
service DessertFactoryService {
  resources: [Candy, Cake]
  operations: [GetWorkshopsNumberInOperation]
}
```

从这里不难看出，operations中定义的操作（OpenAPI）前提是无法关联到任何一个资源。

### 定义资源
资源在服务（Service）或者其他的资源（Resource）中被定义，在服务中定义的资源是主资源，在其他资源下定义的是子资源。对于子资源有约束限制，其主键必须包含全部的父资源的主键，假定资源的层级如下：

![画板](https://intranetproxy.alipay.com/skylark/lark/0/2023/jpeg/309278/1680835772292-fc58eb19-9278-40fa-8240-254f5f544b48.jpeg)

+ 资源B和C主键列表必须包含资源A的全部主键；
+ 资源D的主键列表中必须包含资源B的全部主键（B的主键列表中也必须包含资源A的全部主键）。

通过模型定义资源如下：

```json
$version: 1
namespace example.DessertFactory

service DessertFactoryService {
  version: "2023-04-01"
  resources: [Candy, Cake]
}

struct CandyProperties {
  CandyName: string
  Color: string
  CreationTime: string
}

resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: [UpdateCandy]
  get: GetCandy
  list: ListCandies
}
```

通过resource关键字定义资源Candy，identifyDefinition用于定义资源的主键，糖果的主键是CandyId。properties用于定义资源的属性，属性包含CandyName、Color及CreationTime。create、delete、update、get、list用于指定糖果的生命周期操作（OpenAPI）。分别对应的含义如下：

| create | 创建操作，只允许存在一个。 |
| --- | --- |
| update | 修改资源属性操作，允许存在多个，但是建议越少越好。 |
| delete | 删除资源操作，只允许存在一个。 |
| get | 查询资源操作，只允许存在一个，出参需要返回全部的资源属性。 |
| list | 获取资源列表操作，只允许存在一个，出参需要返回全部的资源主键及全量的资源属性。 |


### 定义操作
操作（operation）存在于service或者resource中，通过operation关键字可以定义一个操作。操作包含三个元素，入参的定义、出参的定义及错误。不论是ROA还是RPC风格，也不论参数的位置是在query、path或者body中，都需要定义出来。

```json
operation GetCandy {
  input: GetCandyInput
  output: GetCandyOutput
  errors: [CandyNotExist]
}

struct GetCandyInput {
  CandyId: string
}

struct GetCandyOutput {
  CandyId: string
  CandyName: string
  Color: string
  CreationTime: string
}

error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  errorMessage: "The queried candy does not exist"
}
```

操作中包含三个属性，input定义操作的入参，output定义操作的出参， errors定义操作中出现的所有错误码。其中input和output是struct结构，直接通过struct关键字定义，错误通过error关键字定义。

资源的生命周期对于各类型操作的出入参存在要求，以下是基本的要求：

| 操作分类 | 出、入参分类 | 主键 | 其他全部属性 |
| --- | --- | --- | --- |
| create | 入参 |  | |
| | 出参 | ✅ | |
| delete | 入参 | ✅ | |
| | 出参 | | |
| get | 入参 | ✅ | |
| | 出参 | | ✅ |
| list | 入参 | | |
| | 出参 | | ✅ |
| update | 入参 | ✅ | |
| | 出参 | | |


这种约束也可以从模型上通过annotate定义出来，写法如下：

```json
struct GetCandyInput {
  @required
  CandyId: string
}
```

当operation是资源的生命周期操作时，在operation的入参和出参中可以直接通过`$`引用，这样就可以不用定义属性和类型。写法如下：

```json
struct GetCandyInput {
  @required
  $CandyId
}
```

模型知晓GetCandyInput这个结构是operation GetCandy的入参，而GetCandy这个operation是资源Candy的查询接口，因此，这里引用的`$CandyId`就是资源的属性。

通过`$`引用的方式资源属性名称和操作的出入参结构中的名称必须一致，如果存在例外，也可以使用`resourceProperty`annotate来标注对应的属性，这个标注也适用于复合结构，对于复合结构应用时层级保持不变。

```json
struct GetCandyInput {
  @required
  @resourceProperty("CandyId")
  CandyIdOther:string
}
```

这个定义方式表示入参中的`CandyIdOther`等同于资源属性中的`CandyId`。

模型中约定，资源关联的operation中所有的出入参都必须是资源的属性，否则需要通过`notResourceProperty`annotate标出，定义如下：

```json
operation GetCandy {
  input: GetCandyInput
  output: GetCandyOutput
  errors: [CandyNotExist]
}

struct GetCandyInput {
  CandyId: string
}

struct GetCandyOutput {
  CandyId: string
  CandyName: string
  Color: string
  CreationTime: string
  @notResourceProperty
  RequestId: string
}

error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  errorMessage: "The queried candy does not exist"
}
```

对于GetCandyOutput结构，RequestId就被标记为不是资源的属性。强烈建议在资源相关的操作中尽可能少的出现notResourceProperty。

### 定义错误码
错误码包含在operation中，描述operation出错时的错误信息格式，使用error关键字可以定义一个错误，如下：

```json
error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  errorMessage: "The queried candy does not exist"
}
```

其中httpCode表示出错时返回的HTTP状态码，errorCode和errorMessage分别在出错时的字段中呈现。对于错误码，需要显示通过annotate标注是否可以重试，默认为false，表示不可重试，通过`retryable` annotate标注，如下：

```json
@retryable(true)
error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  errorMessage: "The queried candy does not exist"
}
```

### 定义RAM鉴权
RAM的鉴权信息是operation上的额外信息，在模型中可以通过`requiredPermission`annotate完成定义。在使用`requiredPermission`前，需要从service上定义当前支持的condition列表，定义支持的condition列表用`conditions`annotate来完成，如：

```json
@defineCondition
struct RAMConditionEncrypted {
  conditionKey: "ebs:Encrypted"
  type: "Bool"
  description: "是否加密云盘"
  documentUrl: "https://help.aliyun.com/x.html"
}

@conditions([RAMConditionEncrypted])
service DessertFactoryService {
  version: "2023-04-01"
  resources: [Candy, Cake]
}
```

`defineCondition`用于定义一个额外的condition，这个condition区别于RAM_BASIC_CONDITIONS基础的condition条件。

此外还可以使用`defineAction`annotate单独定义一个action，例如：

```json
@defineAction
struct RAMActionPassRole {
  description: "使用某个特定的RAM Role"
  action: "ram:PassRole"
  resources: [ {
      `resource`: "acs:ram::{#accountId}:role/{#RoleName}"
      required: true
    }
  ]
  conditions: [
    RAMConditionEncrypted
  ]
}
```

在定义完service支持的condition及需要的额外的action后，就可以从operation上使用`requiredPermission`定义操作需要的RAM鉴权信息，如下：

```json
@requiredPermission({
  action: "dessert:GetCandy"
  otherActions: [RAMActionPassRole]
  resources: [ { 
    `resource`: "acs:example:{#regionId}:{#accountId}:candy/{#CandyId}"
    required: true
  }]
  conditions: [RAMConditionEncrypted]
}
)
operation GetCandy {
  input: GetCandyInput
  output: GetCandyOutput
  errors: [CandyNotExist]
}

```

其中` dessert `是服务的RAM code，ram code需要保持全小写。

### 接入RMC
通过`rmc`annotate注解，可以声明资源接入到RMC，声明后模型定义需要进一步满足RMC的要求。模型定义如下：

```json
@rmc({
  enable: true
})
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}
```

### 接入Terraform
通过`terraform`annotate注解，可以声明资源接入到Terraform，声明后模型定义需要进一步满足Terraform的要求，模型定义如下：

```json
@terraform({
  enable: true
})
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}
```

### 接入审计
通过`actionTrail`annotate注解，可以声明资源接入到ActionTrail，声明后模型定义需要进一步满足ActionTrail的要求，模型定义如下：

```json
@actionTrail({
  enable: true
})
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}
```

### 接入TAG
通过`tagService`annotate注解，可以声明资源接入到TAG，声明后模型定义需要进一步满足TAG的要求，模型定义如下：

```json
array TagSchema {
  item: TagSchemaItem
}

struct TagSchemaItem {
  TagKey: string
  TagValue: string
}

struct CandyProperties {
  CandyName: string
  Color: string
  CreationTime: string
  Tags: TagSchema
}

@tagService({
  enable: true
})
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}

```

当接入TAG时，资源的属性上必须存在Tags属性。

### 接入资源组
通过`resourceGroup`annotate注解，可以声明资源接入到资源组，声明后模型定义需要进一步满足资源组的要求，模型定义如下：

```json
struct CandyProperties {
  CandyName: string
  Color: string
  CreationTime: string
  ResourceGroupId: string
}

@resourceGroup({
  enable: true
})
resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}
```

当接入资源组时，资源属性必须存在ResourceGroupId，类型为string。

### 完整示例
```json
$version: 1
namespace example.DessertFactory

service DessertFactoryService {
  version: "2023-04-01"
  resources: [Candy]
}

array TagSchema {
  item: TagSchemaItem
}

struct TagSchemaItem {
  TagKey: string
  TagValue: string
}

struct CandyProperties {
  CandyName: string
  Color: string
  CreationTime: string
  Tags: TagSchema
  ResourceGroupId: string
}

resource Candy {
  identifyDefinition: {
    CandyId: string
  }
  properties: CandyProperties
  create: CreateCandy
  delete: DeleteCandy
  update: UpdateCandy
  get: GetCandy
  list: ListCandies
}

@requiredPermission({
  action: "example:GetCandy"
  otherActions: [RAMActionPassRole]
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:candy/{#CandyId}"
    required: true
  }]
  conditions: [RAMConditionEncrypted]
}
)
operation GetCandy {
  input: GetCandyInput
  output: GetCandyOutput
  errors: [CandyNotExist]
}

struct GetCandyInput {
  CandyId: string
}

struct GetCandyOutput {
  CandyId: string
  CandyName: string
  Color: string
  CreationTime: string
  @notResourceProperty
  RequestId: string
  @optional
  $Tags
  @optional
  $ResourceGroupId
}

error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  errorMessage: "The queried candy does not exist"
}

@requiredPermission({
  action: "example:CreateCandy"
  otherActions: [RAMActionPassRole]
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:candy/*"
    required: true
  }]
  conditions: [RAMConditionEncrypted]
}
)
operation CreateCandy {
  input: CreateCandyInput
  output: CreateCandyOutput
  errors: [MachineFailure]
}

error MachineFailure {
  httpCode: 500
  errorCode: "MachineFailure"
  errorMessage: "Failed because of machine failure"
}

struct CreateCandyInput {
  @required
  CandyName: string
  @required
  Color: string
  Tags: TagSchema
  ResourceGroupId: string
}

struct CreateCandyOutput {
  @required
  CandyId: string
}

@requiredPermission({
  action: "example:DeleteCandy"
  otherActions: [RAMActionPassRole]
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:candy/${CandyId}"
    required: true
  }]
  conditions: [RAMConditionEncrypted]
}
)
operation DeleteCandy {
  input: DeleteCandyInput
  output: DeleteCandyOutput
  errors: [InternalError]
}

error InternalError {
  httpCode: 500
  errorCode: "InternalError"
  errorMessage: "Internal server error"
}

struct DeleteCandyInput {
  @required
  CandyId: string
}

struct DeleteCandyOutput {
  @notResourceProperty
  Success: boolean
}

operation UpdateCandy {
  input: UpdateCandyInput
  output: UpdateCandyOutput
  errors: [InternalError]
}

struct UpdateCandyInput {
  @required
  CandyId: string
  @required
  $CandyName
  @optional
  $Color
  @optional
  $Tags
  @optional
  $ResourceGroupId
}

struct UpdateCandyOutput {
  @notResourceProperty
  Success: boolean
}

@paginated({
  inputToken: "NextToken",
  outputToken: "NextToken",
  pageSize: "MaxResults",
  items: "Candies"
})
@requiredPermission({
  action: "example:ListCandies"
  otherActions: [RAMActionPassRole]
  resources: [{
    `resource`: "acs:example:{#regionId}:{#accountId}:candy/*"
    required: true
  }]
  conditions: [RAMConditionEncrypted]
}
)
operation ListCandies {
  input: ListCandiesInput
  output: ListCandiesOutput
  errors: [InternalError]
}

struct ListCandiesInput {
  @optional
  ResourceGroupId: string
  @optional
  CandyId: string
  @optional
  $CandyName
}

struct ListCandiesOutput {
  @required
  Candies: Candies
}

array Candies {
  item: CandyItem
}

struct CandyItem {
  @required
  $CandyId
  @required
  $CandyName
  @required
  $Color
  @required
  $CreationTime
  @optional
  $Tags
  @optional
  $ResourceGroupId
}
```


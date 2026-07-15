## service
服务定义。服务是资源、操作（API）、错误码组合而成的对外服务的类型。通过service关键字定义，具体的定义语法如下：

```json
service DessertFactoryService {
  
}
```

服务支持如下的属性：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| version | string | 服务的版本号，对外等同于OpenAPI的版本，例如2020-01-01。  服务的版本号不一定非得是日期，可以是任意类型的字符串。 |
| resources | [string] | 服务下的资源列表，包含的每一个资源需要使用resources XXX在相同的namespace下定义。 |
| operations | [string] | 服务下的操作列表，这些操作必须使用operation XXX在相同的namespace下定义。  资源下的操作不需要在服务下单独定义，这些定义的操作是和资源无关的部分。 |


一个具体的定义示例：

```json
$version: 1
namespace example.DessertFactory

service DessertFactoryService {
  version: "2023-04-01"
  resources: [Candy]
  operations: [ProbeFactoryStatus]
}
```

这个示例中，使用service关键字定义了名字为DessertFactoryService的服务，其版本是2023-04-01，包含一个Candy的资源及一个操作ProbeFactoryStatus。

## resource
资源定义。资源是若干个operation操作的对象，包含围绕这个对象的生命周期的若干个operation及相关的配置定义。资源必须在服务下，使用resources关键字定义，方式如下：

```json
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

资源支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| identifyDefinition | struct | 资源的主键列表，资源的主键必须是顶级的元素，不能是符合结构下的子元素。  分类为子资源时，主键列表必须包含全部父资源的主键。 |
| properties | struct | 资源的属性列表（除主键外），资源支持的全部属性必须显式定义。 |
| create | string | 指定创建operation，该operation调用后完成资源的创建，主键在create之后从出参中给出。 |
| get | string | 指定查询operation，查询的返回的出参中需要包含全部的资源属性。 |
| update | [string] | 指定修改的operation，用于修改所有可改的资源属性。 |
| list | string | 用户指定批量查询资源列表的operation，查询的返回的出参中需要包含全部的资源属性。 |
| delete | string | 指定删除的operation。 |
| operations | [string] | 不是资源的生命周期OpenAPI，但是和资源相关的操作集合。 |
| resources | [string] | 指定当前资源的子资源集合。 |


一个具体的示例如下：

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
```

在这个示例中，给服务DessertFactoryService定义了一个资源Candy，包含了资源主键、属性、指定了一系列资源生命周期操作的operation。

## operation
定义操作。操作包含入参、出参及一个OpenAPI操作可能的错误码。操作包含在资源或者服务下，当包含在资源下时，表示资源的生命周期操作。通过关键字operation可以定义操作，定义的格式如下：

```json
operation CreateCandy {
  input: CreateCandyInput
  output: CreateCandyOutput
  errors: [CreateCandyError]
}
```

操作支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| input | string | 操作的入参，每一个operation定义的出参结构需要是唯一的component。不允许多个不同的操作使用相同的入参格式定义。 |
| output | string | 操作的出参，每一个操作的出参定义需要一个唯一的component，不允许多个不同的操作使用相同的出参格式定义。 |
| errors | [string] | 操作可能出现的错误码集合定义，error使用单独的error关键字定义。 |


一个具体的示例如下：

```json
operation GetCandy {
  input: GetCandyInput
  output: GetCandyOutput
  errors: [CandyNotExist]
}

struct GetCandyInput {
  CandyId: string
}

// 使用$引用资源属性模式，必须使用@for显式声明引用的资源名称
@for(Candy)
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
```

示例中定义了一个操作GetCandy，声明了入参的结构为GetCandyInput，出参结构为GetCandyOutput，并定义了一个错误码CandyNotExist。

特别的，从GetCandyOutput中有了一个特殊的写法`$ResourceGroupId`，这个写法表示引用资源的属性ResourceGroupId，表示和资源属性的名称和类型表示一致。此外，还有`@optional`这种写法，这个是模型支持的annotate，在后续的章节中会写介绍，模型通过注解的方式给组成元素更多的信息表达。

## error
定义错误。操作中可能碰到的错误都应该通过错误标准化定义出来。通过error关键字可以定义错误，定义的方式如下：

```json
error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  backendErrorCode: "candy-not-exist"
  errorMessage: "The queried candy does not exist"
}
```

错误支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| httpCode | int32 | HTTP Code，当操作以OpenAPI方式交付时的HTTP状态码。 |
| errorCode | string | 错误码。 |
| errorMessage | string | 错误详情。 |
| backendErrorCode | string | 后端的错误码。 |
| type | string | user - 用户自定义 |
| default | boolean | 是否为默认错误码，默认为false |
| extendedErrorCode | string | 继承的错误码 |


一个具体的示例如下：

```json
error InternalError {
  httpCode: 500
  errorCode: "InternalError"
  backendErrorCode: "internal-error"
  errorMessage: "Internal server error"
}
```

## $test
定义 resource 和 operation 的测试用例，示例：

```bash
$test COMPONENT_ID TEST_NAME {
}
```

第一个参数COMPONENT_ID声明需要测试的对象，支持 operation 或者 instance 定义；第二个参数的是测试用例的名称。

当测试用例依赖其他的 case 先执行时，在参数中可使用标准的 JSONPATH 的方式引用其他资源测试用例的属性值，一个实际的 case 如下：

```bash
resource Instance {
  identifyDefinition: {
    InstanceId: string
  }
  properties: {
    SourceInstanceId: string
    InstanceName: string
    Complex: {
      A: string
    }
  }
}

// 前置步骤创建一个实例
$test Instance InstanceTest {
  init: {
    InstanceName: "name"
    Complex: {
      A: "a"
    }
  }
}

$test Instance OtherInstanceCreate {
  init: {
    InstanceName: "otherName"
    // InstanceTest 是其他测试的名称，InstanceId是其他测试的资源属性
    SourceInstanceId: "{{$.InstanceTest.InstanceId}}"
    Complex: {
      A: "{{$.InstanceTest.Complex.A}}"
    }
  }
}
```

示例中，InstanceTest 为前置步骤，创建了一台实例，OtherInstanceCreate为另外一个测试，引用了测试InstanceTest，因此在调度时，InstanceTest先创建，创建后资源Instance的属性会被保存下来，OtherInstanceCreate中可以使用`[{{$.InstanceTest.Complex.A}} ](about:blank)` 这种JSONPATH 的格式来引用资源属性。

### 资源测试支持的属性
支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| init | struct | 资源创建阶段声明的属性目标值。  示例：  ```bash {   // 声明创建需要的目标结果，其中名称为a，类型为b   init: {     Name: "a"     Type: "b"   } } ```  :::warning **此步骤不可为空**，对于只读资源，可以填入 filter 的条件。  :::  |
| modifies | array<struct> | 修改节点资源属性的变化，为数组，数组的每个条目声明资源属性的变化值，例如:  ```bash {   // 声明创建需要的目标结果，其中名称为a，类型为b   init: {     Name: "a"     Type: "b"   }   // 两个属性变化步骤   modifies: [{     // 第一步，将 name 由a修改为 a1     Name: "a1"   }, {     // 第二步，将 type 由b修改为c     Type: "c"   }] } ```  :::warning + 如果您的资源不支持修改，这个步骤可以缺省；  + "$_null_" 会被识别为清空该字段，请不要使用为修改的 string 值。  + "$_empty_" 会被识别为清空该字段，请不要使用为修改的 string 值。  :::  |
|  destroy | struct |  通常，销毁阶段不需要辅助参数，直接声明为{}即可，例如：  ```bash {   // 声明创建需要的目标结果，其中名称为a，类型为b   init: {     Name: "a"     Type: "b"   }   // 两个属性变化步骤   modify: [{     // 第一步，将 name 由a修改为 a1     Name: "a1"   }, {     // 第二步，将 type 由b修改为c     Type: "c"   }]   // 销毁   destroy: {} } ```  如果资源的删除存在操作私有属性，那么在销毁阶段允许传递，例如：  ```bash {   // 声明创建需要的目标结果，其中名称为a，类型为b   init: {     Name: "a"     Type: "b"   }   // 两个属性变化步骤   modify: [{     // 第一步，将 name 由a修改为 a1     Name: "a1"   }, {     // 第二步，将 type 由b修改为c     Type: "c"   }]   // 销毁   destroy: {     // 传递一个强制删除     Force: true   } } ```  :::warning 如果您的资源不支持主动销毁，这个步骤可以置空。  对于支持销毁的资源，如果此步骤未填写，则用例在执行后不会主动销毁，因此如果您需要在测试后删除资源，请保留destroy步骤，参数可以传空。  :::  |


### operation 测试支持的属性
支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| input | struct | 调用 API 的入参值。  例如：  ```bash {   input: {     // 调用 API 的入参传值为a1     Name: "a1"   } } ```  |


### 测试入口定义
使用$test 定义的测试用例默认都作为被依赖的用例存在，在调度执行时不会直接运行。因此您需要从 [@testConfig](https://aliyuque.antfin.com/cloudspec/model/kpctnvme370fzyag) annotate中显式配置用例为入口，示例如下：

```bash
@testConfig({
    main: true
})
$test Instance test_1 {
    init: {
        Name: "test"
        Complex: [
            {
                C1: {
                    C2: "a"
                }
            }, {
                C1: {
                    C2: "b"
                }
            }
        ]
    }
    modifies: []
}
```

作为依赖的 case：

```bash
// 这里没有声明 main=true，因此此用例只会在有用例标记为 main=true的依赖时才会被调用执行。
$test Instance test_before {
  init: {
    Name: "before"
  }
  destroy: {}
}

@testConfig({
    main: true
})
$test Instance test_1 {
    init: {
        // 资源属性不存在$.FlowStoragePath2，解析需要报错
        Name: "{{$.test_before.Name}}"
        Complex: [
            {
                C1: {
                    C2: "a"
                }
            }, {
                C1: {
                    C2: "b"
                }
            }
        ]
    }
    modifies: []
    destroy: {}
}
```

由于 test_before 没有声明 main=true，因此此用例只会在有用例标记为 main=true的依赖时才会被调用执行。在示例中test_1依赖test_before的资源属性，因此在测试运行时，会先调用test_before，然后再执行test_1，由于test_before和test_1均声明了要销毁，在调度时执行的顺序如下：

```bash
test_before(create) -> test_1(create) -> test_1(destroy) -> test_before(destroy)
```

### 支持的函数
$test中属性的值支持通过函数传入，声明方式为：

```plain
{{function.functionName(parameter)}}
```

目前支持的函数：

```plain
time()：返回当前的毫秒时间戳
randomNum(min,max)：返回一个介于 min 和 max 之间的随机整数
randomString(len)：返回一个随机字符串，长度为 len
randomIntString(min,max)：返回返回一个介于 min 和 max 之间的随机数，类型为String
```

其中，返回值为string的函数支持在字符串中拼接。

示例

```plain
$test Instance test_1 {
    init: {
        CreateTime: "{{function.time()}}"
        // 注意，仅randomString和randomIntString这种返回值为String的函数支持拼接
        Name："my-test-{{function.randomIntString(1,100)}}"
        Complex: [
            {
                C1: {
                    // 注意这里C2的值是int
                    C2: "{{function.randomNum(1,100)}}"
                }
            }, {
                C1: {
                    C2: "{{function.randomString(10)}}"
                }
            }
        ]
    }
}
```

### A4.4 runtimeType annotate
#### 说明
服务的运行时网关类型注册。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 支持pop，none，默认值为pop。 |

### A4.7 isolationType annotate
#### 说明
服务的隔离级别的配置，当实例化的service上不配置该annotate时，表示当前版本的服务是outer类别。

> 注意：只在实例化的service上定义时才生效，服务级的IDL定义中的isolationType请使用@serviceConfig中的isolationType字段配置。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 服务隔离的类别，支持：inner（内部服务）、outer（外部服务）。required |

#### 示例
```cspec
$version: 1
namespace: a.b

@isolationType("inner")
service A {}
```

### A4.2 horizontalCodeMapping annotate
#### 说明
配置服务和横向组件的code映射关系。

#### 选择器
```cspec
:is(service) [service|version ?= false]
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| pipCode | string | 服务映射的云知code。required |
| locationCode | string | 服务映射的瑶池code。 |
| ram | string | 云产品在RAM服务中的标识。 |
| tag | string | 云产品在TAG服务中的标识。 |
| resourceGroup | string | 云产品在资源组服务中的标识。 |
| terraform | string | 云产品在Terraform服务中的标识。 |
| actionTrail | string | 云产品在ActionTrail服务中的标识。 |
| rmc | string | 云产品在资源管理服务中的标识。 |
| config | string | 云产品在配置审计服务中的标识。 |
| ccApi | string | 云产品在CC API服务中的标识。 |
| relatedAwsServiceCode | [string] | 映射的AWS的code。 |
| lxCode | string | 映射的凌霄code。 |

#### 示例
```cspec
@horizontalCodeMapping({
  pipCode: 'pipCode1',
  locationCode: 'locationCode1'
})
service Utopia {
  resources: [House]
}
```

### A4.5 serviceConfig annotate
#### 说明
服务的配置信息。

#### 选择器
```cspec
:is(service)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 服务定义的struct。 |

通过defineServiceConfig定义struct后配置，例如：

```cspec
@defineServiceConfig
struct ServiceConfig {}

@serviceConfig(ServiceConfig)
service A {}
```

### A4.6 defineServiceConfig annotate
#### 说明
定义服务的配置信息。

#### 选择器
```cspec
:is(struct)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| deploymentType | string | 服务的部署级别（已废弃）：region 地域级部署、center 中心化部署 |
| isolationType | string | 服务隔离的类别，支持：inner 内部服务、outer 外部服务 |
| gatewayInfos | [GatewayInfo] | 网关的配置。 |

#### GatewayInfo
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| gatewayType | string | 网关类别，支持pop、none。required |
| duplicated | boolean | 是否重复定义的POP Code。特殊场景使用，请勿自行使用。 |
| code | string | popCode。required |
| qps | integer | 服务支持的QPS。required |
| domains | [string] | 服务支持的域名。required |
| isolationType | string | 服务隔离的类别，支持：inner 内部服务、outer 外部服务。required |
| versions | [versionScheme] | 服务的版本配置。 |
| centers | [string] | 部署站点，支持：china 中国、intl 国际 |
| deploymentType | string | 服务的部署级别：region 地域级部署、center 中心化部署。required |
| defaultRoaVersion | string | 默认ROA版本 |

#### versionScheme
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| version | string | 服务的版本。required |
| apiStyle | string | 服务的类别，支持rpc/restful。required |
| isolationType | string | 服务隔离的类别，支持：inner 内部服务、outer 外部服务。required |
| errorMapping | errorMappingStruct | 错误映射配置。 |

#### errorMappingStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| exceptionExpress | string | 错误的判断条件。required。例如：`success=true` |
| errors | [reference] | 错误列表，使用error定义。 |
| unknownError | unknownErrorStruct | 未知的错误映射 |

#### unknownErrorStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| errorCode | string | 错误码。required |
| errorMessage | string | 错误信息。required |
| httpCode | integer | HTTP code。required |

#### 示例
```cspec
$version: 1
namespace: alicloud.Ecs

@serviceConfig(EcsConfig)
service Ecs {}

@defineServiceConfig
struct EcsConfig {
  deploymentType: "region"
  isolationType: "outer"
  gatewayInfos: [{
    gatewayType: "pop"
    duplicated: true
    code: "Ecs"
    qps: 100
    domains: ["ecs.aliyun.com", "ec2.aliyun.com"]
    centers: ["intl", "china"]
    isolationType: "outer"
    versions: [{
      version: "2014-05-26"
      apiStyle: "rpc"
      isolationType: "outer"
      errorMapping: {
        exceptionExpress: "success=true"
        unknownError: {
          errorCode: "wsda2"
          errorMessage: "aw212asd"
          httpCode: 404
        }
        errors: [Abv, Abv12]
      }
    }]
  }]
}
```


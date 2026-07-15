测试用例相关的 annotate。

### A14.1  testConfig** annotate**
#### 说明
测试的配置信息。

#### 选择器
```json
:is(test)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| main | boolean | 是否是待执行的测试，默认为 false。  为 false 时，表示当前测试不会被主动执行，而是作为其他用例的依赖存在。 |
| execConfigUuid | string | 用例的执行账号UID，从镇元配置后获取。 |
| forceCreate | Boolean | 对于作为依赖资源时生效，默认为 false。不同的值表示：  + false，执行引擎根据 init 步骤执行的状态值进行判断，如果账号下存在完全匹配的资源时不再创建，复用之前的资源；  + true，执行引擎会根据 init 步骤的状态值尝试销毁已经存在的资源，销毁后再执行创建。 |
| uses | [reference] | :::warning 特别注意，资源的测试不允许依在 uses 配置operation分类的测试作为输入。  请尽量避免前置调用 API，如果需要前置调用某个 API，需要使用 @before 注解中的 invokeOperations 字段声明，且invokeOperations中的API 调用有先后顺序。  **您通常不需要显式配置 uses。**  :::  当前用例依赖的用例列表，use 后可以从测试用例引用其状态/出参值。不同类别的互相引用有约束限制，具体限制请参考：[不同测试类别互相依赖情况](#Of8Rn)  以下是几个不同模式引用的示例：  #### 引用 struct 作统一变量定义  可以将变量定义到一个 struct 中，uses 可以引用 struct，示例：[use引用 struct 定义变量](#gsjW6)。  #### 资源测试依赖其他资源的 state  资源测试可以依赖其他其他资源的 state。被依赖资源会在当前资源创建前创建，在当前资源销毁后销毁。  请参考示例：[资源测试依赖其他资源的 state](#a2UA9)  #### operation 依赖其他资源的 state  operation的测试允许依赖其他的资源属性，示例：[operation 依赖其他资源的 state](#W910s)  #### operation 依赖其他 operation 的出参  operation的测试也允许引用其他 operation 的出参，示例：[operation 测试依赖其他 operation 的出参](#Ol8eD) |
| runtime | [runtime](#nJHLy) | 配置运行时配置。  :::warning 请注意，execConfigUuid在镇元 2.0均可以配置这些参数，本地配置了这些参数时，以本地配置优先。  :::  |


#### runtime 配置
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| regionId | string |  测试时的 regionId。中心化产品无需配置。 |
| bindHost | string | 测试环境为 pre 或者 daily 时绑定的host 地址，例如 pop 的预发地址：`popunify-pre.aliyuncs.com`。 |
| endpoint | string | endpoint 的拼接规则，regionId 部分请使用固定变量 [RegionId]，示例：`amp-2-vpc.[RegionId].aliyuncs.com`，中心化域名无需保留变量。 |
| env | string | API 调用的环境，支持 online、pre、daily，默认为 online。 |
| envOverrideSupport | boolean | 支持测试环境被改写，测试运行时会统一指定运行环境：  + 允许被改写时，不论是 main 还是依赖资源都会被指定为特定的环境；  + 如果不允许改写，则以此处声明的环境为准。  默认为 true。 |
| accountOverrideSupport | boolean | 是否支持改写测试账号配置，极端情况下，资源的创建需要使用指定的账号，如果标记为不可被覆盖，在执行换账号测试时，该测试执行不会被更换，该配置项仅对非 main 的测试生效。  默认值为 true，允许被覆盖。 |


#### runtimeOthers
当您的资源操作API需要通过其他namespace下的API进行异步检查的时候，需要您在main下的testConfig配置相应的信息

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| popCode | string | 异步API所属popCode |
| regionId | string | 异步API所对应region |
| host | string | 异步API绑定的host |
| endpoint | string | 异步API所属产品的endpoint |


示例如下：  


```json
@testConfig({
  runtimeOthers:[{
    popCode:"Ecs"
    regionId:"cn-hangzhou"
    endpoint:"esa-sp-pre.ecs.aliyuncs.com"
  },{
    popCode:"Vpc"
    regionId:"cn-hangzhou"
    endpoint:"esa-sp-pre.ecs.aliyuncs.com"
  }]
})
$test Instance testA{}
```

#### 从本地指定 AK/SK
支持您从本地环境变量设置测试执行的 AK/SK/STS TOKEN信息，**本地环境变量的优先级高于代码中配置的 execConfigUuid中配置的信息**，支持的环境变量如下：

| 环境变量 | 说明 |
| --- | --- |
| ALIBABA_CLOUD_ACCESS_KEY_ID | AccessKey ID |
| ALIBABA_CLOUD_ACCESS_KEY_SECRET | AccessKey Secret |
| ALIBABA_CLOUD_ROLE_ARN |  用于扮演角色时的角色 ARN。 |


当测试配置中的 runtime.accountOverrideSupport为 false 时，针对该测试，环境变量的设置会失效，通常作用在交叉使用多个账号测试时，请谨慎使用。

##### 不同测试类别uses依赖情况
| 测试类别 | 依赖资源 | 依赖operation |  依赖 operation 出参 |
| --- | --- | --- | --- |
| 资源测试 | ✅ | ❌ | ❌ |
| operation 测试 | ✅ | ✅ | ✅ |


#### 示例
```json
operation GetInstance {
  input: {}
  output: {}
}

@testConfig({
  main: true
  execConfigUuid: "xxx"
})
$test GetInstance GetInstanceTest1 {
  input: {}
  asserts: []
}
```

##### use引用 struct 定义变量
```bash
// 统一的变量定义，这里定义一个变量a，并赋值为 1。
struct varDefine {
  a: "1"
}

operation GetInstance {
  input: {
    a: string
  }
  output: {}
}

@testConfig({
  main: true
  execConfigUuid: "xxx"
})
$test GetInstance GetInstanceTest1 {
  input: {
    // 使用变量的值，最终传入的值为 1
    a: "{{$.varDefine.a}}"
  }
  asserts: []
}
```

##### 资源测试依赖其他资源的 state
```bash
$version: 1
namespace: a.b

resource Vpc {
  identifyDefinition: {
    VpcId: string
  }
  properties: {
    VpcName: string
    CidrBlock: string
  }
}

resource VSwitch {
  identifyDefinition: {
    VpcId: string
    VSwitchId: string
  }
  properties: {
    VSwitchName: string
  }
}

// 依赖 VPC 资源的创建
$test Vpc VpcCreate {
  init: {
    VpcName: "test"
    CidrBlock: "10.0.0.0/8"
  }
  destroy: {
  }
}

@testConfig({
  main: true
})
$test VSwitch VSwitchCreate {
  init: {
    VpcId: "{{$.VpcCreate.VpcId}}"
    VSwitchName: "test"
  }
  destroy: {
  }
}
```

##### operation 测试依赖其他资源的 state
```bash
$version: 1
namespace: a.b

resource Vpc {
  identifyDefinition: {
    VpcId: string
  }
  properties: {
    VpcName: string
    CidrBlock: string
  }
}

operation CreateVswitch {
  input: {
    VpcId: string
    VSwitchName: string
  }
  output: {
    VSwitchId: string
  }
}

// 依赖 VPC 资源的创建
$test Vpc VpcCreate {
  init: {
    VpcName: "test"
    CidrBlock: "10.0.0.0/8"
  }
  destroy: {
  }
}

@testConfig({
  main: true
})
$test CreateVswitch VSwitchCreate {
  input: {
    VpcId: "{{$.VpcCreate.VpcId}}"
    VSwitchName: "test"
  }
}
```

##### operation 测试依赖其他 operation 的出参
```bash
$version: 1
namespace: a.b

operation CreateVswitch {
  input: {
    VpcId: string
    VSwitchName: string
  }
  output: {
    VSwitchId: string
  }
}

operation GetVswitch {
  input: {
    VSwitchId: string
  }
  output: {
    VpcId: string
    VSwitchName: string
  }
}

// 先调用 CreateVswitch 创建一个 VSwitch
$test VSwitch VSwitchCreate {
  input: {
    VpcId: "test"
    VSwitchName: "myVSwitch"
  }
}

@testConfig({
  main: true
})
$test GetVswitch GetVswitchTest {
  input: {
    VSwitchId: "{{$.VSwitchCreate.VSwitchId}}"
  }
}
```

### A14.2  before** annotate**
#### 说明
配置测试用例前置的信息。

#### 选择器
```json
:is(test)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| invokeOperations | [reference] | 用例运行前调用的 operation 测试，通常用于资源测试前置需要开通的场景。  reference 的对象必须是$test的名称，有严格的调用顺序。 |


#### 示例
```bash
$version: 1
namespace: a.b

operation OpenVpc {
  input: {}
  output: {}
}

resource Vpc {
  identifyDefinition: {
    VpcId: string
  }
  properties: {
    VpcName: string
    CidrBlock: string
  }
}

// 开通接口
$test OpenVpc open_vpc {
}

@before({
  invokeOperations: [open_vpc]
})
$test Vpc vpc_test {
  input: {
    VpcName: "test"
    CidrBlock: "10.0.0.0/8"
  }
}
```


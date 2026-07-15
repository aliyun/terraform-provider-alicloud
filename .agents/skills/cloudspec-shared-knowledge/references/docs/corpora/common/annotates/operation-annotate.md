# CloudSpec Operation 注解列表

本文档列出了 CloudSpec Operation 支持的所有注解及其详细配置说明。

## 目录

- [A6.1 http 注解](#a61-http-注解)
- [A6.2 paginated 注解](#a62-paginated-注解)
- [A6.2.1 numberPaginated 注解](#a621-numberpaginated-注解)
- [A6.3 operationInfo 注解](#a63-operationinfo-注解)
- [A6.4 backendConfigurationHttp 注解](#a64-backendconfigurationhttp-注解)
- [A6.5 backendConfigurationHsf 注解](#a65-backendconfigurationhsf-注解)
- [A6.6 backendConfigurationDubbo 注解](#a66-backendconfigurationdubbo-注解)
- [A6.7 backendConfigurationHttpHsf 注解](#a67-backendconfigurationhttphsf-注解)
- [A6.8 backendConfiguration 注解（已废弃）](#a68-backendconfiguration-注解已废弃)
- [A6.9 errorMapping 注解](#a69-errormapping-注解)
- [A6.10 systemParameter 注解](#a610-systemparameter-注解)
- [A6.11 systemParameterConfiguration 注解](#a611-systemparameterconfiguration-注解)
- [A6.12 parameter 注解](#a612-parameter-注解)
- [A6.13 visibility 注解](#a613-visibility-注解)
- [A6.14 gatewayOptions 注解](#a614-gatewayoptions-注解)
- [A6.15 controlPolicy 注解](#a615-controlpolicy-注解)
- [A6.16 codec 注解](#a616-codec-注解)
- [A6.17 valueMapping 注解](#a617-valuemapping-注解)
- [A6.18 additionalJsonConfiguration 注解](#a618-additionaljsonconfiguration-注解)
- [A6.19 backendName 注解](#a619-backendname-注解)
- [A6.20 successHttpCode 注解](#a620-successhttpcode-注解)
- [A6.21 resourcePropertyMapping 注解](#a621-resourcepropertymapping-注解)
- [A6.22 repeatListParameter 注解](#a622-repeatlistparameter-注解)
- [A6.23 dependency 注解](#a623-dependency-注解)
- [A6.24 incremental 注解](#a624-incremental-注解)
- [A6.25 rootMapping 注解](#a625-rootmapping-注解)
- [A6.26 retryPolicies 注解](#a626-retrypolicies-注解)
- [A6.27 nested 注解](#a627-nested-注解)
- [A6.28 returnMode 注解](#a628-returnmode-注解)
- [A6.28 extErrorMapping 注解](#a628-exterrormapping-注解)
- [A6.29 defineMultipleOutput 注解](#a629-definemultipleoutput-注解)
- [A6.30 multipleOutput 注解](#a630-multipleoutput-注解)
- [A6.31 headerOutput 注解](#a631-headeroutput-注解)
- [A6.32 rateLimit 注解](#a632-ratelimit-注解)
- [A6.33 autoBackend 注解](#a633-autobackend-注解)
- [A6.34 async 注解](#a634-async-注解)
- [A6.35 eventInfo 注解](#a635-eventinfo-注解)
- [A6.36 asyncApi 注解](#a636-asyncapi-注解)
- [A6.37 serializeFormat 注解](#a637-serializeformat-注解)
- [A6.38 backendConfigurationWebsocket 注解](#a638-backendconfigurationwebsocket-注解)

## A6.1 http 注解
#### 说明
配置OpenAPI的基本信息。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| ~~apiStyle(已废弃)~~ | ~~string~~ | ~~OpenAPI风格：~~  + `~~RPC~~`  + `~~RESTFUL~~` |
| schemes | [schemeStruct](#schemeStruct) | 支持的传输协议，支持多环境配置。 |
| methods | [string] | HTTP请求方法，线上环境的配置（默认配置）  `RPC`风格可选值：  + `get`  + `post`  `RESTFUL`风格可选值：  + `get`  + `head`  + `post`  + `put`  + `delete`  + `patch`  + `options` |
| authenticators | [string] | 认证方式：  + `AK`  + `APP`  + `PrivateKey`  + `BearerToken`  + `Anonymous` |
| signatureVersion | [string] | 认证方式：  + `V1`  + `V4` |
| requestContentType | [string] | 请求时支持的媒体类型：  + `application/json`  + `application/xml`  + `application/x-www-form-urlencoded`  + `multipart/form-data`  + `application/octet-stream` |
| responseContentType | [string] | 返回时支持的媒体类型：  + `application/json`  + `application/xml`  + `application/octet-stream` |
| uri | string | 操作的URI模式，当风格是ROA时必须设置。 |
| servers | [[ServiceStruct](#serviceStruct)] | 定义服务列表，在ROA模式下有效。 |
| variables | struct | 结构如variablesStruct。  配置诸如$.policies.controlPolicyName等在灰度、预发、日常覆盖默认值的条目。 |
| tags | [{  key: string  description: string  }] | 支持的标签 |
| upStreamEventTypeExpression | string | 代表上行方向上的消息类型，格式为jsonPat，需配合events使用 |
| downStreamEventTypeExpression | string | 代表下行方向上的消息类型，格式为jsonPath，需配合events使用 |


#### variablesStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| gray | array<struct> | {  key: "$.policies.controlPolicyName（示例）"  value: "__null__"(示例，也可以为整形、Boolean)  } |
| pre | array<struct> | {  key: "$.policies.controlPolicyName（示例）"  value: "__null__"(示例，也可以为整形、Boolean)  } |
| daily | array<struct> | {  key: "$.policies.controlPolicyName（示例）"  value: "__null__"(示例，也可以为整形、Boolean)  } |
|  |  |  |


##### ServiceStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| host | host | 后端的host值。 |
| basePath | string | 路径。 |
| description | string | 描述信息。 |


#### schemeStruct 结构
假定scheme的取值定义为

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| online | [[SCHEME_VALUE](#schemeValue)] | 线上环境的配置，当其他环境不配置时，默认同线上环境的配置。 |
| daily | [[SCHEME_VALUE](#schemeValue)] | 日常环境的配置。  如果日常环境不配置，那么日常环境默认使用online的配置。 |
| pre | [[SCHEME_VALUE](#schemeValue)] | 预发环境的配置。  如果预发环境不配置，那么预发环境默认使用online的配置。 |
| gray | [[SCHEME_VALUE](#schemeValue)] | 灰度环境的配置。  如果灰度环境不配置，那么灰度环境默认使用online的配置。 |

#### Schema支持的值
其取值支持：

+ `http`
+ `https`
+ `sse`

#### 示例
```json
@http({
  schemes: {
    online: ["http", "https"]
	}
  methods: ["get", "post"]
  authenticators: ["AK"]
  signatureVersions: ["V4"]
})
operation CreateEcs {
  input: {}
  output: {}
}
```

## A6.2 paginated 注解
#### 说明
配置OpenAPI的分页信息。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| inputToken | string | NextToken，指定入参的next token参数。**required** |
| outputToken | string | NextToken，指定出参的next token参数。**required** |
| maxItems | string | 例如MaxResults，指定入参中的最大条目数量的参数。**required** |
| maxItemsDefault | int32 | 最大条目的默认值**required** |
| pageTruncated | string | 指定出参中是否截断的参数 |
| totalCount | string | 指定出参中总条目的参数 |
| items | string | 指定出参中的数组节点的值。 |
| pageTruncatedExtra | struct |  |
| totalCountExtra | struct |  |


#### 示例
```json
@paginated({
  inputToken: "NextToken"
  outputToken: "NextToken"
  maxItems: "MaxResults"
  items: "ReplicaPairs"
  maxItemsDefault: 10
})
operation GetCandy {
  
}
```

## A6.2.1 numberPaginated 注解
#### 说明
配置OpenAPI的分页信息，按页码分页。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| initialPageNumber | int32 | 初始化的页码。 |
| initialPageSize | int32 | 默认每页的大小配置。 |
| inputPageNumber | string | 入参中指定页码的参数。 |
| inputPageSize | string | 入参中指定每页大小的参数。 |
| recordTotal | string | 出参中指定总条目的参数。 |
| pageTotal | string | 指定出参中指定总页面数的参数。 |
| items | string | 指定出参中的数组节点的值。 |


#### 示例
```json
@numberPaginated({
  initialPageNumber: 1
  initialPageSize: 10
  inputPageNumber: "PageNumber"
  inputPageSize: "PageSize"
  recordTotal: "$.RecordTotal"
})
operation GetCandy {
  
}
```

## A6.3 operationInfo 注解
#### 说明
配置 operation 的读写、能力树等信息。配置操作的读写分类，resource中的生命周期OpenAPI不需要指定这个分类，resources中的operations和service中的operations支持配置此分类。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
annotate支持一个struct的入参，struct中支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| operationType | string | create、update、get、list、delete、none中的一个值。none代表当前API操作是非资源化的，标识为"不涉及"。 |
| riskType | string | high - 高风险  none - 低风险 |
| chargeType | string | paid - 涉及到付费  free - 免费 |
| operationTypeOld | string | **已废弃，请不要使用。**  readAndWrite - 读写  read - 只读  write - 只写 |
| typeFromOperation | string | 来自operation上的标记分类。**已废弃，请不要使用。** |
| abilityTreeCode | string | 能力树对应的 code |
| abilityTreeNodes | [string] | 能力树对应的 codes |
| apiTrailType | string | 例如controlEvent |
| autoTest | boolean | 是否支持自动化测试 |
| notSupportAutoTestReason | string | 不支持自动化测试的原因 |
| tenantRelevance | string |  租户相关性支持：  + tenant 租户型 API，API 访问或者操作租户内的信息、资源，需要具备"租户隔离"能力  + publicInformation 公共信息 API，API访问公开的信息，无需具备"租户隔离"能力 |


#### 示例
```json
@operationInfo({
  operationType: 'get'
  
})
operation GetCandy {
  
}
```

## A6.4 backendConfigurationHttp 注解
#### 说明
定义HTTP模式的后端地址配置。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
该注解本身没有参数，通过结构体struct信息表达，struct支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| applicationName | string | 后端应用的名称。 |
| retries | struct | struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  取值相同，默认为-1表示不限制。  例如：  ```json {   online: -1   pre: 2 } ```  |
| timeout | struct | 线上的超时时间，单位毫秒，value类型为int32，最大值60000ms，struct支持的key如下：  + online **required**  + pre  + gray  + daily  例如：  ```json {   online: 3000   pre: 60000 } ```  |
| requestType | string | 请求类型，支持Object和String  + Object 增加映射  + String 默认映射 |
| backendUrl | struct | 后端的URL地址，struct支持的key如下：  + online **required**  + pre  + gray  + daily  例如：  ```json {   online: '/a'   pre: '/b' } ```  |
| responseType | string | 返回结果类型，支持Object和String，当API风格为RESTFUL时必填。  + Object 映射  + String 透传 |
| signKeyName | struct | 后端签名，值为string类型，支持的key如下：  + online   + pre  + gray  + daily |
| dynamicBackendAddress | struct | 动态后端地址，值的结果为string，支持的key如下：  + online   + pre  + gray  + daily |
| statusCodeTransparent | boolean | 状态码透传，默认为false。RESTFUL风格时需要。 |
| sign | Boolean | 是否签名 |
| signPolicy | string | 签名的策略，当前支持固定值Local |
| consume | string | 增强映射模式下使用，后端请求类型：  支持：  + application/x-www-form-urlencoded  + application/json |
| httpMethod | string | 增强映射模式下使用，后端http方法  支持：  + get  + post  + delete  + put |
| httpsValidation | boolean | 是否校验HTTPS证书。 |
| invokeType | string | 调用类型 |
| httpsValidationGray | string | 灰度环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| httpsValidationPre | string | 预发环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| httpsValidationDaily | string | 日常环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| sessionTimeout | long | 代表用户自定义的会话持续超时时间，单位毫秒，需要配合events使用 |
| sessionIdleTimeout | long | 代表会话空闲超时时间，单位毫秒，需要配合events使用 |


#### 示例
```json
@backendConfigurationHttp({
  applicationName: "Demo"
  retries: {
    online: -1
  }
  timeout: {
    online: 3000
  }
  requestType: "Object"
  backendUrl: {
    online: "http://example.com"
  }
  responseType: "Object"
  signKeyName: {
    online: "key1234"
  }
	dynamicBackendAddress: {
    online: 'xxxx'
  }
})
operation A {}
```

## A6.5 backendConfigurationHsf 注解
#### 说明
定义HSF模式的后端地址配置。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
该注解本身没有参数，通过结构体struct信息表达，struct支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| applicationName | string | 后端应用的名称。 |
| invokeType | string | 入参的形式，支持单个和多个入参：  + Multi 多个入参  + Single 单个入参 |
| retries | struct | 重试次数，不开启为-1，默认值-1。  struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  例如：  ```json {   online: -1   pre: 2 } ```  |
| timeout | struct | struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  取值相同，默认为-1表示不限制。  例如：  ```json {   online: 6000   pre: 7000 } ```  |
| group | struct | HSF分组信息  struct的key如下，value的取值为string：  + online **required**  + pre  + gray  + daily |
| version | struct | HSF的版本号  struct的key如下，value的取值为string：  + online **required**  + pre  + gray  + daily |
| service | string | HSF接口名称 |
| method | string | HSF方法名称 |
| paramTypes | [string] | HSF的参数格式列表 |
| dynamicBackendAddress | struct | 动态后端地址，值的结果为string，支持的key如下：  + online   + pre  + gray  + daily |
| serviceProtocol | string | 服务协议类型，支持：  + dubbo  + tri |


#### 示例
```json
@backendConfigurationHsf({
  applicationName: "Demo"
  retries: {
    online: -1
    pre: 1
    daily: 2
    gray: 3
  }
  timeout: {
    online: 3000
    pre: 1
    daily: 2
    gray: 3
  }
  invokeType: "Single"
  `service`: "com.aliyun.amp"
   method: "GetCandy"
  paramTypes: ["java.lang.String", "java.lang.Integer"]
  group: {
      online: "hsf"
        pre: "pre"
        daily: "daily"
        gray: "gray"
    }
    version: {
      online: "1.0.0"
      pre: "pre-1.0.0"
        daily: "daily-1.0.0"
        gray: "gray-1.0.0"
    }
  dynamicBackendAddress: {
    online: 'online-xxxx'
    pre: 'pre-xxxx'
    daily: 'daily-xxxx'
    gray: 'gray-xxxx'
  }
})
operation BB {}
```

### A6.6 backendConfigurationDubbo annotate
#### 说明
定义Dubbo模式的后端地址配置。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
该注解本身没有参数，通过结构体struct信息表达，struct支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| applicationName | string | 后端应用的名称。 |
| invokeType | string | 入参的形式，支持单个和多个入参：  + Multi 多个入参  + Single 单个入参 |
| retries | struct | 重试次数，不开启为-1，默认值-1。  struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily |
| timeout | struct | struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  取值相同，默认为-1表示不限制。 |
| version | struct | Dubbo的版本号  struct的key如下，value的取值为string：  + online **required**  + pre  + gray  + daily |
| service | string | Dubbo接口名称 |
| method | string | DubboF方法名称 |
| paramTypes | [string] | Dubbo的参数格式列表 |
| dynamicBackendAddress | string | 动态后端地址，值的结果为string，支持的key如下：  + online   + pre  + gray  + daily |
| serviceProtocol | string | 服务协议类型，支持：  + dubbo  + tri |
| registryProtocol | string | 注册中心协议类型，支持：  + dubbo（适用于Dubbo2.x）  + nacos（适用于Dubbo3.x） |


#### 示例
```json
@backendConfigurationDubbo({
  applicationName: "Demo"
  retries: {
    online: -1
    pre: 1
    daily: 2
    gray: 3
  }
  timeout: {
    online: 3000
    pre: 1
    daily: 2
    gray: 3
  }
  invokeType: "Single"
  `service`: "com.aliyun.amp"
   method: "GetCandy"
  paramTypes: ["java.lang.String", "java.lang.Integer"]
  version: {
    online: "1.0.0"
    pre: "pre-1.0.0"
      daily: "daily-1.0.0"
      gray: "gray-1.0.0"
  }
  dynamicBackendAddress: {
    online: 'online-xxxx'
    pre: 'pre-xxxx'
    daily: 'daily-xxxx'
    gray: 'gray-xxxx'
  }
})
operation CC {}
```

### A6.7 backendConfigurationHttpHsf annotate
#### 说明
定义HTTP-HSF模式的后端地址配置。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
该注解本身没有参数，通过结构体struct信息表达，struct支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| applicationName | string | 后端应用的名称。 |
| invokeType | string | 入参的形式，支持单个和多个入参：  + Multi 多个入参  + Single 单个入参 |
| retries | struct | 重试次数，不开启为-1，默认值-1。  struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily |
| timeout | struct | struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  取值相同，默认为-1表示不限制。 |
| group | struct | HSF分组信息  struct的key如下，value的取值为string：  + online **required**  + pre  + gray  + daily |
| version | string | HSF的版本号  struct的key如下，value的取值为string：  + online **required**  + pre  + gray  + daily |
| service | string | HSF接口名称 |
| method | string | HSF方法名称 |
| paramTypes | [string] | HSF的参数格式列表 |
| dynamicBackendAddress | struct | 动态后端地址，值的结果为string，支持的key如下：  + online   + pre  + gray  + daily |
| backendUrl | struct | 后端的URL地址，struct支持的key如下：  + online **required**  + pre  + gray  + daily |
| responseType | string | 支持Object和String  + Object 映射  + String 透传 |
| signKeyName | struct | 后端签名，值为string类型，支持的key如下：  + online   + pre  + gray  + daily |
| requestType | string | 请求类型  + 增强类型 Object  + 默认映射 String |
| sign | Boolean | 是否签名 |
| signPolicy | string | 签名的策略，当前支持固定值Local |


#### 示例
```json
@backendConfigurationHttpHsf({
  applicationName: "Demo"
  retries: {
    online: -1
  }
  timeout: {
    online: 3000
  }
  group: {
    online: "hsf"
  }
  version: {
    online: "1.0.0"
  }
  service: "com.aliyun.amp"
  method: "GetCandy"
  paramTypes: ["java.lang.String"]
  requestType: "Object"
  backendUrl: {
    online: "http://example.com:10220"
  }
  responseType: "Object"
  signKeyName: {
    online: "signKey1234"
  }
})
```

### A6.8 backendConfiguration annotate（已废弃）
#### 说明
定义操作的后端配置。

:::danger
已废弃，不再支持。

废弃时间：2023.09.08。

请直接使用backendConfigurationHttp等注解。

:::

#### 选择器
```json
:is(operation)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| ~~online~~ | ~~string~~ | ~~默认的后端配置。~~ |
| ~~daily~~ | ~~string~~ | ~~日常的额外配置，非必填，存在值时会覆盖online的配置~~ |
| ~~pre~~ | ~~string~~ | ~~预发的额外配置，非必填，存在值时会覆盖online的配置~~ |
| ~~grey~~ | ~~string~~ | ~~灰度环境的额外配置，非必填，存在值时会覆盖online的配置~~ |


### A6.9 errorMapping annotate
#### 说明
定义后端返回的错误映射配置。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| errorExpression | string | 错误码判断条件。 |
| codeField | string | 后端错误码的获取字段 |
| errorMessageField | string | 获取后端错误信息的字段 |
| httpStatusCodeField | string | 获取后端HTTP状态码的字段 |
| dynamicCodeField | string | 获取动态错误码的字段 |
| dynamicMessageField | string | 获取动态错误信息的字段 |
| extendedCodeField | string | 用于表示扩展后端错误码的字段 |


#### 示例
```json
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpCode"
  dynamicCodeField: "dymCode"
  dynamicMessageField: "dymMessage"
  extendedCodeField: "extendCode"
})
operation GetCandy {
}
```

### A6.10 systemParameter annotate
#### 说明
定义operation需要接受的系统参数，所有支持的系统参数见：[https://aliyuque.antfin.com/cloudspec/model/ocave330kb14qrbk](https://aliyuque.antfin.com/cloudspec/model/ocave330kb14qrbk) 。请注意，系统参数的接收参数固定，请从以上表格查阅，不支持修改。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
[string]，数组的值为所有支持的系统参数，[https://aliyuque.antfin.com/cloudspec/model/ocave330kb14qrbk](https://aliyuque.antfin.com/cloudspec/model/ocave330kb14qrbk) 

#### 示例
```json
@systemParameter(["callerUid", "AccessKeyId"])
operation A {
 
}
```

### A6.11 systemParameterConfiguration annotate
#### 说明
在操作上定义后端需要的系统参数。

:::color4
已废弃，不再支持。

:::

### A6.12 parameter annotate
#### 说明
定义后端接受参数的字段名称，如果不标注，则默认使用相同的名称传递。

同时适用于入参和出参的映射。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| in | string | 入参接收参数的位置，支持：  + header  + query  + path  + body  + system  + host  + formData |
| allowEmptyValue | Boolean | 是否允许空值 |
| style | string | 风格，不同参数类型支持的 style 如下：  **array 类型：**  + repeatList   > repeatList 的序列化方式为XXX.N的形式，例如`Instance.1=i-instance1&&Instance.2=i-instance2`   + simple   > simple 数组的序列化方式为逗号分隔，例如`i-instance 1, i-instance2`   + spaceDelimited   > spaceDelimited 数组的序列化方式为空格分隔，例如`i-instance1 i-instance2`   + pipeDelimited   > pipeDelimited 数组的序列化方式为竖杠| 分割，例如i-instance1|i-instance2   + json  > json：数组的序列化方式为json，例如［"i-instance1", "i-instance2"]   + flat  > flat 数组的序列化方式为XXX.N的形式，例如Instance.1=i-  > instance1&Instance.2=i-instance2，相对于repeatList而言功能更强大，支持数组  > 嵌套数组等形式   **map 类型：**  + json  > json 序列化方式为json，例如{"a": 1, "b": 2}   + flat  > flat 序列化方式为A.B的形式，例如Param.A=1&Param.B=2   **object类型：**  + json  > json 序列化方式为json，例如{"a": 1, "b": 2}   + flat  > flat 序列化方式为A.B的形式，例如Param.A=1&Param.B=2   |
| responseName | string | 出参从后端那个字段获取。 |
| checkBlank | Boolean | 是否检验空值。 |
| trim | Boolean | 是否删除前后空格。 |
| nullable | Boolean | 是否可为null |
| isFileTransferUrl | Boolean | 是否为文件上传url |
| parseType | string | 支持struct、json、string。 |
| repeatListMinItems废弃 | int32 | 当style = repeatList时最少元素个数（含） |
| repeatListMaxItems废弃 | int32 | 当style = repeatList时最多元素个数（含） |
| repeatListIndexName | string | 当style = repeatList时，后端序号字段 |
| repeatListSubName | string | 当style = repeatList时，下级参数名称 |
| checkRepeatList | Boolean | 是否开启repeatList参数校验，默认为false |
| repeatListSequence | Boolean | 当style = repeatList时，其序列是否连续 |
| repeatListDataType | string | 支持Json/Map，当style = repeatList时，后端接收类型 |
| itemName | string | 当开启了showItemName时，且类型为array时，设置节点的名称。 |
| mappingType | string | roa风格，后端为hsf时，body入参到后端hsf入参的映射类型 |
| bizType | string | 后端参数类型：  + header  + query  + path  + cookie  + system  + bizSystem  + const  + body  + bodyMember  + formData  + host |
| timestampToDate | boolean | HSF转HSF_http之后需要网关来转换timestamp 到日期类型,日期类型只有一个默认值，不需要与其它字段结合使用。 |
| index | int64 | 后端简单入参的参数位置。 |
| groupIndex | int64 | 后端复杂入参参数组的参数位置。 |
| backendParamType | string | 后端参数类型，前期设计用于标识HSF、Dubbo的后端服务下前端传入日期时的标识 |
| minItems | int64 | 最少的元素数量 |
| maxItems | int64 | 最多的元素数量 |
| backendPosition | string | 后端参数位置。 |
| docRequired | boolean | 是否文档必填。 |
| nullToEmpty | Boolean | null是否转换为空 |
| sdkPropertyName | string | 在SDK中的名称 |
| sdkIgnore | Boolean | SDK中是否忽略 |
| emptyToNull | boolean | 空数据是否转换为null |
| extendType | string | 继承的类型，例如 Map(String,Object) |


#### 示例
```json
struct DeleteCandyInput {
  @required
  @parameter({
    name: 'candyId'
  	in: 'query'
	})
  CandyId: string
}
```

### A6.13 visibility annotate
#### 说明
定义operation的可见性，默认为Private私有。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | operation的可见性。  + SpecialPurpose 客户专用  + Private 私有  + Public 公开 |


#### 示例
```json
@visibility("Public")
operation test {
  
}
```

### A6.14 gatewayOptions annotate
#### 说明
定义operation的网关配置。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| responseLog | Boolean | 记录返回日志 |
| roaRequestBodyLog | Boolean | RESTFUL API记录请求体日志 |
| roaResponseBodyLog | Boolean | RESTFUL API记录响应体日志 |
| fileTransfer | Boolean | 支持文件中转传输 |
| akProvenStatus | Boolean | 开启还是关闭 |
| showJsonItemName | Boolean | 是否开启JSON Item |
| keepClientResourceOwnerId | Boolean | 获取资源所有者 |
| dataUpload | Boolean | 是否是数据类上传类型API |
| dataDownload | Boolean | 是否是数据类下载类型API |
| repeatListConflictValidation | Boolean | 验证repeatList中是否有冲突参数 |
| badInputJsonValidation | Boolean | 验证请求中的非法json参数 |
| badOutputJsonValidation | Boolean | 验证响应中是否含有非法json参数 |
| maxRequestBodyConstraints | integer | 约束请求Body的最大值（字节） |
| remainEmptyArrayInJsonParameter | boolean |  |
| tagOptions | [TagOptions](#lqAWy) | Tag在网关的一些配置。 |
| isolationType | string | 隔离类别。 |
| rgOptions | [RgOptions](#jDWAm) | 资源网关的配置。 |
| outputParamVersion | integer | 出参版本 |
| riskLevel | string | 风险级别 |
| userType | string |  |
| afterAuthService | string |  |
| inputArrayAllowNull | boolean |  |
| outputArrayAllowNull | boolean |  |
| defaultResponseFormat | string |  |
| flowControlResponseHeader | boolean |  |
| signatureComposer | string |  |
| optionalResponseFormats | List<string> |  |
| parameterSensitives | List<SensitiveRegular> |  |
| responseSensitives | List<SensitiveRegular> |  |
| ramAuth | RamAuth | RAM鉴权信息，已废弃请勿使用。 |
| gatewayProtocols | List<string> | 非必填。前端界面上不需要体现，用户直接在代码中编写 |
| originalOutputParamVersion | string | 取值范围："1.0" / "2.0" / "2.5" / null |
| outputParamAlignEnable | boolean | 取值范围：true / false / null |
| requestTransparentTransmission | boolean | 取值范围：true / false / null |


#### RamAuth
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| enable | boolean | 是否启用 |
| authLevel | string | 鉴权级别 |
| authTypes | array<string> | 鉴权类型 |
| ramCode | string | ram code |
| regionIdParamName | string | region ID 在参数中的名称 |
| ramAction | string | ram action |
| ramArn | string | ram 鉴权 ARN |


#### SensitiveRegular
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| prefix | string | 前缀 |
| suffix | string | 后缀 |


#### TagOptions
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| tagPolicyEnable | Boolean | 是否开启 |
| tagCode | string | 例如ecs，接入tag的名称 |
| tagParameterName | string | 参数名称，例如Tag |
| tagParameterType | string | 参数类型，例如list_map |
| tagResourceType | string | 类型，例如static |
| tagResourceTypes | [string] | 例如keypair，支持的资源类型 |
| tagResourceTypeParameterName | String |  |
| tagParameterPath | String |  |
| verifyNoTag | Boolean |  |
| tagResourceTypeParameterMap | array<array<[tagResourceTypeParameter](#IP7l3)>> | 注意：是双层数组。 |


##### tagResourceTypeParameter
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| key | string | 键。 |
| value | string | 值。 |


#### RgOptions
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| rgInterceptEnable | Boolean | 是否开启拦截。 |


#### 示例
```json
@gatewayOptions({
  responseLog: true
})
operation test {
  
}
```

### A6.15 controlPolicy annotate
#### 说明
定义operation的访问策略。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| online | [policyStruct](#EhScH) | 线上环境的策略名称 |
| gray | [policyStruct](#EhScH) | 灰度环境的策略名称（可复用线上） |
| pre | [policyStruct](policyStruct) | 预发环境的策略名称 |
| daily | [policyStruct](#EhScH) | 日常环境的策略名称 |


##### policyStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| controlPolicyName | string | 访问策略名称 |
| grayScalePolicyName | string | 灰度扩展策略 |
| rateLimitPolicy | rateLimitPolicyStruct | 流控策略 |


##### rateLimitPolicyStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| unit | string | 访问策略名称 |
| userRateLimit | Integer | 单用户限制 |
| apiRateLimit | Integer | API限制 |
| ipRateLimit | Integer | IP限制 |
| specialRateLimitPolicyName | string | 特殊流控策略 |


#### 示例
```json
@controlPolicy({
	online: {
    controlPolicyName: 'policytest'
  }
})
operation test {
  
}
```

### A6.16 codec annotate
#### 说明
定义属性的POP编解码方式。默认为空，无需设置。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | + dingEncodeV1  + dingDecodeV1 |


#### 示例
```json
struct test {
  @codec('dingEncodeV1')
  string: string
}
```

### A6.17 valueMapping annotate
#### 说明
值映射配置。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
valueMapping支持的参数为struct，格式如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| defaultValue | string | 值映射的默认值。 |
| cases | [case] | 值映射配置。 |


case（值映射配置）参数支持若干个数组，每个数组是一个struct，struct支持的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| targetValue | string | 展示值。 |
| value | string | 后端值。 |


#### 示例
```json
struct DeleteCandyInput {
  @required
  @valueMapping({
    defaultValue: "test"
  	cases: [{
    	targetValue: "candyId"
  		value: "query"
		}]
	})
  CandyId: string
}
```

### A6.18 additionalJsonConfiguration annotate
#### 说明
配置属性的额外JSON配置。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 额外的JSON配置。 |


#### 示例
```json
struct test {
  @additionalJsonConfiguration('{"a": 1}')
  string: string
}
```

### A6.19 backendName annotate
#### 说明
配置后端参数名称。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 后端参数名称。 |
| - | string | 后端接受参数的位置。 |


#### 示例
```json
struct test {
  @backendName("test", "query")
  a: string
}
```

### A6.20 successHttpCode annotate
#### 说明
配置接口请求成功时的HTTP状态码，默认为200。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 成功时的HTTP状态码。 |


#### 示例
```json
@successHttpCode("200")
operation test {
}
```

### A6.21 resourcePropertyMapping annotate
#### 说明
配置operation上的出入参字段和资源属性的映射关系。

#### 选择器
```json
:is(operation, struct > member)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| requestMappingType | string | 入参的映射类型  枚举:param/api/computed  指映射到具体的API参数上，还是只关联到API上。对于启动ECS 实例和停止ECS实例的API,API参数中没有体现Status属性，但是影响了Status的值，这种场景需要用【api】这种映射类型 |
| resourceProperty | string | 仅当 requestMappingType = api 时生效，表示当前整个 operation 调用后更新的资源属性。  当requestMappingType = api 时，整个配置生效的条目为：  + resourceProperty  + targetValue  + propertyDependencies |
| requestTransform | string | 入参到资源属性映射的UDF配置。 |
| responseTransform | string | 出参到资源属性映射的UDF配置。 |
| requestPathType | string | 入参路径表达式类型  枚举:  normal/jsonPath/repeatList/jsonArray ([1,2,3])/commaSeparated(1,2,3) /kvPairs(用于适配凌霄API，需要通过两个repeatList入参来指定的参数 |
| requestIn | string | 参数的位置，支持formData等 |
| requestPath | string | 入参配置（在operation上配置表示当前参数，不用单独配置，配置会被忽略） |
| requestRequired | Boolean | 入参是否必填 |
| requestConstValue | string | 当入参需要固定传常量时，常量值的配置 |
| targetValue | string | 入参数映射类型为api时，api调用后，该属性变为的目标值是什么 |
| requestDefaultValue | string | 入参的默认值配置 |
| requestMinValue | string | 入参的最小配置 |
| requestMaxValue | string | 入参的最大配置 |
| requestEnumValues | [string] | 入参的枚举值配置 |
| hasSystemDefault | boolean | 是否存在系统默认值 |
| responsePathType | string | 响应参数提取方式 枚举:jsonPath |
| responseMappingType | string | 出参的映射类型  枚举:param/computed  指映射到具体的API参数上，计算属性。默认为param，为空表示为param。 |
| responsePath | string | 响应参数的路径表达式。  在operation上的参数上配置时此条目不用配置，如果配置会被忽略。 |
| constType | Boolean | 常量参数类型   + true(动态资源标识)   + false(普通常量) |
| penetrate | Boolean | 是否是透传   + true(是)   + false(不是) |
| mappingType | string | 映射配置的类型，枚举：  + property(属性，表示属性与xx的映射关系)  + response(API出参，表示API参数和xx的映射关系)  + request（API入参），默认为属性映射配置 |
| requestValueMappings | [valueMappingsItem] | 入参和资源属性资映射配置 |
| responseValueMappings | [valueMappingsItem] | 出参和资源属性的映射配置 |
| propertyDependencies | [propertyDependency] | 属性依赖配置 |


valueMappingsItem的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| resourceValue | string | 作用于资源属性的值。 |
| operationValue | string | operation上的入参或者出参值。 |


propertyDependency的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| conditions | [condition] | 条件。 |
| action | string | 属性依赖的动作，枚举值：  + invalid-无效，  + required-必选，  + valid-生效，  + conflicted-冲突 |


condition的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| property | string | 属性名称 |
| value | string or array | 属性的值，默认为所有有效值 |
| requestProperty | boolean | 是否为操作入参属性，如果是的话需要从requestModel中获取判断值，否则从currentModule中获取，true-是，false-否，默认为false |
| valueType | string | 属性值类型，specifiedValue-指定值，null-空值，empty-空，默认 = specifiedValue，specifiedValue时 value 生效 |


#### 示例
```json
@resourcePropertyMapping({
  requestTransform: "AA"
  responseTransform: "BB"
})
operation test {
}
```

### A6.22 repeatListParameter annotate
#### 说明
配置operation上的入参字段style为repeatList时，当需要按实际传值时的注解。

#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
支持的参数为数组，每个数组下的条目为struct，struct的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| requestPath | string | 入参值的路径，例如`Parameter.10.Value`。 |
| keyPath | string | 入参的key的路径，例如`Parameter.10.Code`。 |
| key | string | 入参的key的值，例如 `BigScreen`。 |
| requestRequired | Boolean | 入参是否必填 |


#### 示例
```json
struct test {
  @repeatListParameter([{
    requestPath: "Parameter.10.Value"
    keyPath: "Parameter.10.Code"
    key: "BigScreen"
    requestRequired: true
  }])
	@parameter({
    checkBlank: false
    repeatListMinItems: 1
    repeatListMaxItems: 100
    repeatListIndexName: "componentIndex"
    repeatListSubName: "test2"
    checkRepeatList: true
    repeatListSequence: true
    repeatListDataType: "Json"
    in: "query"
    style: "repeatList"
  })
  Parameter: ParamterArray
}

array ParamterArray {
  item: ParamterArrayItem
}

struct ParamterArrayItem {
  Code: string
  Value: string
}
```

### A6.23 dependency annotate
#### 说明
配置operation上的条件依赖。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
dependency支持的参数格式为数组，数组中值结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| conditions | [condition] | 条件。 |
| action | string | 属性依赖的动作，枚举值：  + invalid-无效，  + required-必选，  + valid-生效，  + conflicted-冲突 |


condition的结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| property | string | 属性名称 |
| value | array | 属性的值，默认为所有有效值 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@dependency([{
  action: "invalid"
  conditions: [{
    property: "$.InstanceName"
    value: ["1", "2"]
  }
  ]
}
])
operation UpdateFlow {
}


```

### A6.24 incremental annotate
#### 说明
标识operation是增量操作的类型，增量操作的对象通过映射完成。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持struct类型参数，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| type | string | 类型，必选，支持的值如下：  + ADD 增加  + REMOVE 删除  + TRUNCATE 清空 |
| maxBatchSize | int32 | 单次操作数量最大的限制。 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@incremental({
  type: "ADD"
})
operation UpdateFlow {
}


```

### A6.25 rootMapping annotate
#### 说明
标识operation的出参中获取数据的顶级节点。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持struct类型参数，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| responsePathType | string | 出参的序列化类型，例如jsonPath |
| responsePath | string | 例如$，表示出参的struct的根节点就是获取数据的节点。 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@rootMapping({
  responsePathType: "jsonPath"
  responsePath: "$"
})
operation UpdateFlow {
}


```

### A6.26 retryPolicies annotate
#### 说明
标识operation的按错误码的重试策略。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持数组参数，数组的值为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| Code | string | 待重试的错误码。 |
| Interval | int64 | 重试间隔（单位：秒）。 |
| Times | int64 | 最大的重试次数。 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@retryPolicies([{
  Code: "ConcurrentUpdateError"
  Interval: 1
  Times: 10
}
, {
  Code: "InternalServerError"
  Interval: 1
  Times: 10
}
])
operation UpdateFlow {
}


```

### A6.27 nested annotate
#### 说明
默认情况下，出入参到资源属性映射都是从顶级元素开始的，当资源的顶级属性映射到出入参的嵌套子结构时，可以使用该注解。
1. nested修饰的参数类型为object时，表示这个object的元素与资源属性做映射，此时rootMapping取nested修饰的节点的路径2. nested修饰的参数类型为array时，表示这个array的item（item类型必须为object）内部元素与资源属性做映射，此时rootMaping取array的item的路径（路径和showJsonItem选项有关）
#### 选择器
```json
:is(struct > member)
```

#### 支持的属性
没有属性选项，无参注解

#### 示例
```json
resource Candy{
  properties: CandyPro
}

struct CandyPro {
    name: string
    price: int32
}

// 情况1，nested修饰object，rootMapping为  $.date
struct GetCandyOutput{
  success: boolean
  codes: string
  @nested
  date: CandyPro
}
```

```json
resource Candy{
  properties: CandyPro
}

struct CandyPro {
    name: string
    price: int32
}


// 情况2，nested修饰array
// showJsonItem打开，rootMapping为  $.date.item[*]
// showJsonItem关闭，rootMapping为  $.date[*]
struct GetCandyOutput{
  success: boolean
  codes: string
  @nested
  date: ManyCandyPro
}

array ManyCandyPro {
  item:CandyPro
}

```

### A6.28 returnMode annotate
#### 说明
配置operation的异步配置信息。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持数组参数，数组的值为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| async | boolean | 是否是异步模式。**required**默认为false，表示同步。 |
| callback | componentId | 异步查询的operation。当模式为异步时必须配置。 |
| interval | int64 | 重试间隔（单位：秒）。当模式为异步时必须配置。 |
| times | int64 | 最大的重试次数。当模式为异步时必须配置。 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@returnMode({
  async: true
  callback: API2
  interval: 1
  times: 10
})
operation UpdateFlow {
}


```

### A6.28 extErrorMapping annotate
#### 说明
配置operation旁路错误码。该配置不影响API运行时返回结果，仅用于在网关日志中提取正确的错误信息。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持数组参数，数组的值为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| codeField | string | 错误码字段，jsonpath格式，例如 $.Code。**required** |
| errorMessageField | string | 错误信息字段名称，jsonpath格式，例如 $.Message。**required** |
| httpCodeExpression | struct | 根据 http 状态码进行判断，支持的字段为：  + left：string类型，此处固定为 httpStatusCode  + condition：string类型，判断条件，支持      - =      - !=      - >=      - >      - <=      - <  + right ：string类型，表示条件的具体值。  :::warning httpCodeExpression和bodyExpression至少设置一个。  :::  |
| bodyExpression | struct | 根据返回参数进行判断。  + left：string类型，通过jsonpath获取某个返回值字段。  + condition：string类型，判断条件，支持      - =      - !=      - >=      - >      - <=      - <  + right ：string类型，表示条件的具体值。  :::warning httpCodeExpression和bodyExpression至少设置一个。  :::  |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@extErrorMapping({
  codeField: "$.Code"
  errorMessageField: "$.Message"
  httpCodeExpression: {
    condition: "="
    right: "500"
    left: "httpStatusCode"
  }
  bodyExpression: {
    left: "$.requestId"
    condition: ">="
    right: "600"
  }
})
operation UpdateFlow {
}


```

### A6.29 defineMultipleOutput annotate
#### 说明
定义一个复合结构，标识是不同的状态码的operation返回值。

#### 选择器
```json
:is(struct, array)
```

#### 支持的属性
无需参数。

#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@defineMultipleOutput
struct Response201 {
  name: string
}

@defineMultipleOutput
array Response202 {
  item: string
}

@multipleOutput([{
  httpCode: "201"
  outputStructure: Response201
}, {
  httpCode: "202"
  outputStructure: Response202
}])
operation UpdateFlow {
}


```

### A6.30 multipleOutput annotate
#### 说明
配置operation不同状态码的返回值结构，注意operation中定义的output是默认200的返回值，使用successHttpCode注解可以声明成功的状态码。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持数组参数，数组的值为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| httpCode | string | HTTP状态码的值。**required** |
| outputStructure | reference | 指定状态码定义返回时的结构。**required** |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@defineMultipleOutput
struct Response201 {
  name: string
}

@multipleOutput([{
  httpCode: "201"
  outputStructure: Response201
}])
operation UpdateFlow {
}


```

### A6.31 headerOutput annotate
#### 说明
配置从header上输出的信息。

#### 选择器
```json
:is(struct)
```

#### 支持的属性
支持数组参数，数组的值为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| key | string | header的key。**required** |
| type | string | 数据类型。**required**  支持：  + any 任意类型  + string  + number  + integer  + boolean  + array  :::warning 注意，此处传入的类型需要使用string类型，不是IDL中的build-in type。  :::  |
| title | string | 标题。 |
| backendName | string | 后端参数名称。 |
| sdkPropertyName | string | SDK中使用的名称。 |
| format | string | 格式。 |
| example | string | 示例值。 |
| required | Boolean | 是否必填 |
| docRequired | boolean | 是否文档必填 |
| name | string | 名称 |
| description | string | 描述 |
| style | string | 风格 |
| parseType | string | 解析类型 |
| mapValueType | string | 当值的类型是 map 时，指定 map 条目类型，通常在设置为 map 时，style应等于 JSON。  因为从 header 传参，其值是 key:value的键值对。 |


#### 示例
```json
// This file is automatically generated.
$version: 1
namespace: alicloud.FnF.fnf.v20190315

@defineMultipleOutput
@headerOutput([{
  key: 'x-acs-ext-1'
  type: "string"
  title: "示例的参数"
  backendName: "ext1"
}])
struct Response201 {
  name: string
}

@multipleOutput([{
  httpCode: "201"
  outputStructure: Response201
}])
operation UpdateFlow {
}


```

### A6.32 rateLimit annotate
#### 说明
配置operation的流控信息。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持参数格式为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| userDefault | int32 | 用户默认值 |
| userCountWindow | int64 | 用户统计窗口值 |
| apiDefault | int32 | API默认值。 |
| apiCountWindow | int64 | API统计窗口值 |


### A6.33 autoBackend annotate
#### 说明
配置operation出入参默认后端接收参数行为。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持参数格式为struct，结构如下：

| 属性 | 类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| in | string | 否 | 默认对外的参数位置，支持：  + header  + query  + path  + body  + system  + host  + formData  **默认值为 query。** |
| defaultPosition | string | 否 | 参数默认的后端参数位置，支持：  + header  + query  + path  + body  **默认值为 query。** |
| naming | string | 否 | 命名规则，支持：  + default 默认值，参数名称改为小驼峰，对于出参的数组节点，默认为同名的itemName，优先级低于backendName和 parameter，如果从backendName何 parameter 中做了相关的配置，那么默认的行为被覆盖。 |
| arrayStyle | string | 否 | 指定入参类型为 array 默认的风格，支持：  + repeatList   > repeatList 的序列化方式为XXX.N的形式，例如`Instance.1=i-instance1&&Instance.2=i-instance2`   + simple   > simple 数组的序列化方式为逗号分隔，例如`i-instance 1, i-instance2`   + spaceDelimited   > spaceDelimited 数组的序列化方式为空格分隔，例如`i-instance1 i-instance2`   + pipeDelimited   > pipeDelimited 数组的序列化方式为竖杠| 分割，例如i-instance1|i-instance2   + json  > json：数组的序列化方式为json，例如［"i-instance1", "i-instance2"]   + flat  > flat 数组的序列化方式为XXX.N的形式，例如Instance.1=i-  > instance1&Instance.2=i-instance2，相对于repeatList而言功能更强大，支持数组  > 嵌套数组等形式   不指定时，默认值为 json。 |
| mapStyle | string | 否 | 指定入参类型为  map 默认的风格，支持：  + json  > json 序列化方式为json，例如{"a": 1, "b": 2}   + flat  > flat 序列化方式为A.B的形式，例如Param.A=1&Param.B=2   不指定时，默认为 json。 |
| objectStyle | string | 否 | 指定入参类型为 object 时默认的风格，支持：  + json  > json 序列化方式为json，例如{"a": 1, "b": 2}   + flat  > flat 序列化方式为A.B的形式，例如Param.A=1&Param.B=2   不指定时，默认为json。 |


示例：

```bash
$version: 1
namespace: alicloud.oss

@autoBackend
operation A {
  input: AInput
  output: AOutput
  errors: []
}

struct AInput {
  // 后端接收参数为 name，参数位置为 query
  Name: string
  // 后端接收参数为 age，参数位置为 query
  Age: int32
  // 后端接收参数为 sex2，因为backendName 的优先级更高，参数位置为 query
  @backendName("sex2")
  Sex: string
}
```

### A6.34 async annotate
#### 说明
当 operation 关联到资源时，操作异步完成需要检查资源属性的状态时配置。

#### 选择器
```json
:is(operation)
```

#### 支持的属性
支持参数格式为struct，结构如下：

| 属性 | 是否必填 | 类型 | 说明 |
| --- | --- | --- | --- |
| resourceProperty | 是 | String | 轮询检查的资源属性。 |
| delayedTime | 否 | int64 | 第一次请求延迟时间，单位秒，不填，delay=0 |
| failedValues | 否 | [string] | 当扫描到这些状态值时认为失败。 |
| interval | 是 | int64 | 轮询的间隔时间，单位秒。 |
| compareType | 是 | string | 待查询资源属性的类型，支持assertNull/assertEmpty/assertNotEmpty/assertEqual，默认为assertEqual  + assertNull 值为null  + assertEmpty 值存在，但是内容为空  + assertNotEmpty 值存在且非空  + assertEqual 值存在，且等于TargetValue配置的值。 |
| successValues | 否 | [string] | 当检查属性的值等于配置值时，认为成功。 |
| times | 是 | int64 | 最大重试的次数。 |


示例：

```bash
$version: 1
namespace: alicloud.oss

@async({
  resourceProperty: "$.Status"
  interval: 5
  compareType: "assertEqual"
  successValues: "running"
  times: 10
})
operation A {
  input: AInput
  output: AOutput
  errors: []
}
```

### A6.35 eventInfo annotate
#### 说明
配置 operation 的事件驱动信息。

#### 支持的属性
支持参数格式为struct，结构如下：

| 属性 | 是否必填 | 类型 | 说明 |
| --- | --- | --- | --- |
| enable | 否 | boolean | 是否开启事件驱动。 |
| eventNames | 否 | array<string> | API 关联的事件。 |


示例：

```bash
@eventInfo({
  enable: true
  eventNames: ["StartInstance"]
})
operation A {
  input: AInput
  output: AOutput
  errors: []
}
```

### A6.36 asyncApi annotate
#### 说明
当 operation 关联到资源时，操作异步完成需要检查API时配置。

#### 支持的属性
支持参数格式为struct，结构如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| AsyncApi | String | 轮询的异步OpenAPI，例如 `alicloud.SpecTest.SpecTest.v20220101#CreateBook`。 |
| RequestAttributeMappings | [[RequestAttributeMappingsStruct](#LCrg2)] | 异步轮询OpenAPI的入参配置。 |
| ResponseAttributeMappings | [[ResponseAttributeMappingsStruct](#WUGqf)] | 异步轮询OpenAPI的出参到资源属性的映射配置。 |
| AsyncDetections | [[asyncPollingByPropertyStruct](#Xcq0P)] | 异步轮询成功、失败的策略配置。 |


#### RequestAttributeMappingsStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| ApiRequestPath | String | API的入参路径。 |
| RequestRequired | Boolean | 参数是否必填。 |
| MappingType | string | 映射类型  + response 从关联的OpenAPI（资源上关联的OpenAPI）的出参中取值  + const 取常量值  + property 从资源属性上取值 |
| RequestPathType | string | 入参参数类型，例如normal |
| RequestPath | string | 当MappingType = response时，表示关联到的资源操作的OpenAPI的出参值，当MappingType = property时，表示当前资源的属性值。 |
| RequestConstValue | string | 常量值。 |
| RequestTransform | string | 入参的JSONata转换。 |


#### ResponseAttributeMappingsStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| ResourceProperty | String | 资源属性的路径。 |
| MappingType | string | 映射类型，例如property。 |
| ResponsePath | string | 异步轮询OpenAPI的出参表达式，例如`$.BookId`。 |
| ResponsePathType | string | 出参的路径类型，例如normal。 |
| ResponseTransform | string | 出参的JSONata表达式。 |


#### asyncPollingByPropertyStruct
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| ResourceProperty | String | 轮询检查的资源属性。 |
| DelayedTime | int64 | 第一次请求延迟时间，单位秒。 |
| FailedValues | [string] | 当扫描到这些状态值时认为失败。 |
| Interval | int64 | 轮询的间隔时间，单位秒。 |
| TargetValueType | string | 待查询资源属性的类型，支持assertNull/assertEmpty/assertNotEmpty/assertEqual，默认为assertEqual  + assertNull 值为null  + assertEmpty 值存在，但是内容为空  + assertNotEmpty 值存在且非空  + assertEqual 值存在，且等于TargetValue配置的值。 |
| TargetValue | string | 当检查属性的值等于配置值时，认为成功。 |
| Times | int64 | 最大重试的次数。 |


示例：

```java
@asyncApi({
    AsyncApi: a.b#B
    RequestAttributeMappings: [{
      ApiRequestPath: "Category"
      RequestRequired: false
      MappingType: "response"
      RequestPathType: "normal"
      RequestTransform: "xxxxxxxx"
      RequestPath: "$.RequestId"
    }
    ]
    ResponseAttributeMappings: [{
      ResourceProperty: "$.InstanceName"
      MappingType: "property"
      ResponsePath: "$.BookId"
      ResponsePathType: "normal"
      ResponseTransform: "xxxxxxxx"
    }
    ]
    AsyncDetections: [{
      ResourceProperty: "$.BookName"
      DelayedTime: 5
      FailedValues: ["22222", "failed"]
      Interval: 5
      TargetValue: "1111"
      TargetValueType: "assertEqual"
      Times: 10
    }
    ]
  })
operation A {
}
```

### A6.37 serializeFormat annotate
#### 说明
出入参为JSON，描述其序列化后的结构。

#### 支持的属性
参数为string，描述序列化后结构的复合类型结构体名

需要注意的是，该结构体作为描述json结构体，内部不得引用其他结构体，也不应该引用其他结构体

#### 示例
```java
// struct
struct a_input{
    @serializeFormat(test)
    a:any
}

struct test{
    a: string
    b: array<int32>
    c: map<{
        c1:string
        c2:int32
    }>
}

// array
struct a_input{
    @serializeFormat(test)
    a:any
}

array test{
    item: {
        a: string
        b: array<int32>
        c: map<{
            c1:string
            c2:int32
        }>
    }
}



// map
struct a_input{
    @serializeFormat(test)
    a:any
}

map test{
    key:string
    value: {
        a: string
        b: array<int32>
        c: map<{
            c1:string
            c2:int32
        }>
    }
}

```

以struct为例子，上述描述等价于：

```json
{
  "type":"object",
  "properties":{
    "a":{
      "type":"string"
    },
    "b":{
      "type":"array"
      "items":{
        "type":"integer"
        "format":"int32"
      }
    }
    "c":{
      "type":"object"
      "additionalProperties":{
        "type":"object"
        "properties":{
          "c1":{
            "type":"string"
          }
          "c2":{
            "type":"integer"
            "format":"int32"
          }
        }
      }
    }
  }
}
```

### A6.38 backendConfigurationWebsocket annotate
#### 说明
定义Websocket模式的后端地址配置。

#### 选择器
```json
:is(operation)
```



#### 支持的属性
该注解本身没有参数，通过结构体struct信息表达，struct支持的属性如下：

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| applicationName | string | 后端应用的名称。 |
| retries | struct | struct的key如下，value的取值为int32：  + online **required**  + pre  + gray  + daily  取值相同，默认为-1表示不限制。  例如：  ```json {   online: -1   pre: 2 } ```  |
| timeout | struct | 线上的超时时间，单位毫秒，value类型为int32，最大值60000ms，struct支持的key如下：  + online **required**  + pre  + gray  + daily  例如：  ```json {   online: 3000   pre: 60000 } ```  |
| requestType | string | 请求类型，支持Object和String  + Object 增加映射  + String 默认映射 |
| backendUrl | struct | 后端的URL地址，struct支持的key如下：  + online **required**  + pre  + gray  + daily  例如：  ```json {   online: '/a'   pre: '/b' } ```  |
| responseType | string | 返回结果类型，支持Object和String，当API风格为RESTFUL时必填。  + Object 映射  + String 透传 |
| signKeyName | struct | 后端签名，值为string类型，支持的key如下：  + online   + pre  + gray  + daily |
| dynamicBackendAddress | struct | 动态后端地址，值的结果为string，支持的key如下：  + online   + pre  + gray  + daily |
| statusCodeTransparent | boolean | 状态码透传，默认为false。RESTFUL风格时需要。 |
| sign | Boolean | 是否签名 |
| signPolicy | string | 签名的策略，当前支持固定值Local |
| consume | string | 增强映射模式下使用，后端请求类型：  支持：  + application/x-www-form-urlencoded  + application/json |
| httpMethod | string | 增强映射模式下使用，后端http方法  支持：  + get  + post  + delete  + put |
| httpsValidation | boolean | 是否校验HTTPS证书。 |
| invokeType | string | 调用类型 |
| httpsValidationGray | string | 灰度环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| httpsValidationPre | string | 预发环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| httpsValidationDaily | string | 日常环境支持校验 HTTPS  必须为 string 类型：  + true 开启  + false 不开启  + `**__null__**`** 不设置** |
| sessionTimeout | long | 代表用户自定义的会话持续超时时间，单位毫秒，需要配合events使用 |
| sessionIdleTimeout | long | 代表会话空闲超时时间，单位毫秒，需要配合events使用 |


#### 示例
```json
@backendConfigurationWebsocket({
  applicationName: "Demo"
  retries: {
    online: -1
  }
  timeout: {
    online: 3000
  }
  requestType: "Object"
  backendUrl: {
    online: "http://example.com"
  }
  responseType: "Object"
  signKeyName: {
    online: "key1234"
  }
	dynamicBackendAddress: {
    online: 'xxxx'
  }
})
operation A {}
```

### 

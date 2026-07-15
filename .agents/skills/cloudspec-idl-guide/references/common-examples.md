# CloudSpec 高频示例

本页收敛容易写错的典型配置。实际编辑时仍以当前项目已有同类文件为第一模板。

## ROA body 参数

ROA 请求体字段必须声明 `in: "body"`；path 参数必须与 `@http.uri` 中的占位符一致。

```cspec
@http({
  apiStyle: "restful"
  uri: "/2021-04-06/functions/{functionName}/sessions"
  methods: ["post"]
  requestContentType: "application/json"
  responseContentType: "application/json"
})
operation CreateSession {
  input: CreateSessionInput
  output: CreateSessionOutput
  errors: [Error_CreateSession]
}

struct CreateSessionInput {
  @backendName("functionName")
  @parameter({ in: "path" })
  @required
  FunctionName: string

  @backendName("sessionTTLInSeconds")
  @parameter({ in: "body" })
  SessionTTLInSeconds: int64
}
```

## ROA 后端配置

ROA 的 `@backendConfigurationHttp` 通常需要 `requestType: "String"` 和 `responseType: "String"`。`applicationName`、`backendUrl`、`signKeyName`、`signPolicy` 必须参考同项目已有操作。

```cspec
@backendConfigurationHttp({
  applicationName: "FC"
  requestType: "String"
  responseType: "String"
  retries: { online: -1, pre: -1 }
  timeout: { online: 60000 }
  backendUrl: {
    online: "http://pop.fc3.prod/#vpc"
  }
  signKeyName: {
    online: "FC-API-V20210406"
  }
  sign: true
  signPolicy: "Local"
})
```

## ROA operationMapping body requestPath

ROA 资源映射中，body 入参的 `requestPath` 需要带 `body.` 前缀。

```cspec
@defineOperationMapping
struct HttpApiOperationMapping {
  Create: {
    operation: CreateHttpApi
    request: [{
      property: HttpApiName
      requestPath: "body.httpApiName"
      requestIn: "body"
    }]
  }
}
```

RPC 项目通常直接使用参数名，不加 `body.` 前缀。

## resource 测试

```cspec
@testConfig({
  main: true
})
$test Instance test_1 {
  init: {
    Name: "my-test-{{function.randomIntString(1,100)}}"
  }
  modifies: [
    {
      Description: "$_empty_"
    }
  ]
  destroy: {}
}
```

要点：

- 入口测试必须显式 `main: true`。
- 支持销毁的资源保留 `destroy: {}`。
- `"$_null_"` 和 `"$_empty_"` 会被识别为清空字段，不要当普通字符串值使用。

## operation 测试

operation 测试使用 `input`，不要写成 resource 测试的 `init/modifies/destroy`。

```cspec
$test CreateInstance create_instance_case {
  input: {
    RegionId: "cn-hangzhou"
    Name: "my-test-{{function.randomIntString(1,100)}}"
  }
}
```

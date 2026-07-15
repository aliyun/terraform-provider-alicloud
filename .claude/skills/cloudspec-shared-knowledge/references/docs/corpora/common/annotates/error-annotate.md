# CloudSpec Error 注解列表

本文档列出了 CloudSpec 错误码相关的所有注解及其详细配置说明。

## 目录

- [A11.1 retryable annotate](#a111-retryable-annotate)
- [A11.2 serviceAlreadyEnabled annotate](#a112-servicealreadyenabled-annotate)

### A11.1 retryable annotate
#### 说明
确认错误码是否可以重试。

#### 选择器
```cspec
:is(error)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | boolean | 是否可重试，默认为true |

#### 示例
```cspec
@retryable(true)
error CandyNotExist {
  httpCode: 404
  errorCode: "CandyNotExist"
  backendErrorCode: "candy-not-exist"
  errorMessage: "The queried candy does not exist"
}
```

### A11.2 serviceAlreadyEnabled annotate
#### 说明
错误码是否标识服务已经开通。

#### 选择器
```cspec
:is(error)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| desc | string | 描述 |
| descEn | string | 英文描述 |

#### 示例
```cspec
@serviceAlreadyEnabled({
  desc: "本错误码代码服务已经开通"
  descEn: "this errorCode identity to service has already enabled"
})
error MyErrorCode {
  httpCode: 404
  errorCode: "MyErrorCode"
  backendErrorCode: "candy-not-exist"
  errorMessage: "The queried candy does not exist"
}
```

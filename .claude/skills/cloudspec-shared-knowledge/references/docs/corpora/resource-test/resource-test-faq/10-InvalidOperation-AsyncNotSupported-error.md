# 问题描述
异步检查失败，提示`InvalidOperation.AsyncNotSupported`

# 问题表现
异步检查失败，提示`InvalidOperation.AsyncNotSupported`

# 错误原因
由于指定API不支持异步检查，但是资源配置中又配置了异步检查，所以会出现错误。

# 解决方案
删除资源和API映射配置中的异步检查配置信息。

# 示例
原本的映射信息：
```
// 资源定义
resource Record{
   // ...
   delete: @operationMapping(Record_delete_mapping) DeleteRecord
}

// 映射信息
@defineOperationMapping
struct Record_delete_mapping {
  // ...
  errorCodePosition: "$.Activation.ActivationId"
  retryPolicies: []
  resourceNotExistErrorCodes: ["ActivationId.NotFound"]
  notExistCheckType: "checkErrorCode"
}

```

修复后的：
```
// 资源定义
resource Record{
   // ...
   delete: @operationMapping(Record_delete_mapping) DeleteRecord
}

// 映射信息
@defineOperationMapping
struct Record_delete_mapping {
  // ...
}
```
修复后的映射信息Record_delete_mapping中删除了异步检查配置信息。

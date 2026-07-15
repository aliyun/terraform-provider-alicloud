# 问题描述
API调用成功，但是资源测试提示：执行资源操作失败，错误码：OperationFailure.OperationFailed。

# 问题表现
对资源进行增删改查时，对应的API调用成功了，但是资源测试提示执行资源操作xx失败。比如删除资源时，Delete API调用成功，但是资源测试提示：deleteXXX执行失败,错误码:OperationFailure.OperationFailed。

# 错误原因
资源和API中配置了不正确的错误码位置信息（errorCodePosition），导致API调用成功，并且返回了正确的出参，但是测试引擎将正确出参识别为错误码，进而认为执行资源操作失败。

# 解决方案
删除或订正映射信息中的错误码位置信息。

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
  errorCodePosition: "$.RecordId"
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

比如对于上面例子中的DeleteRecord这个API，调用成功后，出参会返回删除的Record的主键ID，
如果Record_delete_mapping中errorCodePosition配置了"$.RecordId"，测试引擎会将RecordId字段识别为错误码，进而认为执行资源操作失败。所以这里需要删除。

# 问题描述
缺少重试策略导致测试用例失败

# 问题表现
资源测试用例在modifies阶段或者destroy阶段失败，失败原因是因为资源状态还没达到预期，导致失败； 
比如destroy阶段失败，错误日志显示资源正处于更改状态中，或者modifies阶段失败，错误日志显示资源正处于创建状态中。

# 解决方案
在CloudSpec IDL的资源与API映射的结构体中添加或订正`retryPolicies`，详细语法参考`.docs/common/annotates/resource-annotate.md`中的
`A5.5 defineOperationMapping annotate`部分中的`retryPolicies`参数和`retryPoliciesStruct`

# 示例
原本的的CloudSpec IDL结构体：
```
// 资源定义
resource Record{
   // ...
   delete: @operationMapping(DeleteRecord_Mapping_Struct) DeleteRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct DeleteRecord_Mapping_Struct{
    // ...
    retryPolicies: []
}
```

修复后的：
```
// 资源定义
resource Record{
   // ...
   delete: @operationMapping(DeleteRecord_Mapping_Struct) DeleteRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct DeleteRecord_Mapping_Struct{
    // ...
   retryPolicies: [{
       Code: "Record.Modifying"
       Interval: 5
       Times: 20
   }]
}
```
修复后的结构体中，`retryPolicies`中添加了`Code`、`Interval`、`Times`三个参数，分别表示错误码、重试间隔和重试次数。
这样配置以后，当执行完modifies步骤后，进入destroy阶段，调用DeleteRecord这个API删除资源时，若返回错误码为`Record.Modifying`，测试并不会认为失败，而是会按照配置进行重试。
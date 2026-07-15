# 问题描述
缺少资源不存在的判定条件导致资源destroy失败

# 问题表现
资源测试用例在destroy阶段，delete API调用成功了，但是出现下面情况：
1. get API调用返回错误，但其实这个get的错误是预期内的，因为资源已经不存在了，但是因为缺少资源不存在的判定条件，导致cli认为Get API调用失败了，进而执行失败
2. delete API调用成功后，将资源的某个标记位属性置为删除状态（eg，Status=Deleted），但是因为缺少资源不存在的判定条件，get API调用一只能查到资源，导致cli认为delete API调用失败了，进而执行失败


# 解决方案
在CloudSpec IDL的资源与Get、或者Delete API映射的结构体中添加或订正`resourceNotExistCondition`，详细语法参考`.docs/common/annotates/resource-annotate.md`中的
`A5.5 defineOperationMapping annotate`部分中的`resourceNotExistCondition`参数和`resourceNotExistConditionStruct`，具体在Get还是Delete中添加需要根据实际情况选择

# 示例
原本的的CloudSpec IDL结构体：
```
// 资源定义
resource Record{
   // ...
   properties:{
       // ...
       Status:string
   }
   get: @operationMapping(GetRecord_Mapping_Struct) GetRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct GetRecord_Mapping_Struct{
    // ...
    resourceNotExistCondition: {}
}
```

修复后的：
```
// 资源定义
resource Record{
   // ...
   properties:{
       // ...
       Status:string
   }
   get: @operationMapping(GetRecord_Mapping_Struct) GetRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct GetRecord_Mapping_Struct{
    // ...
   resourceNotExistCondition:{
       anyOf: [
           {
               notExistCheckType: "checkErrorCode"
               resourceNotExistErrorCodes: ["InvalidResourceId.NotFound"]
           },
           {
               notExistCheckType: "checkProperty"
               notExistCheckProperty: "$.Status"
               notExistCheckTargetValueType: "assertEqual"
               notExistCheckTargetValue: "Deleted"
           }
       ]
   }
}
```
修复后的映射结构体中，`resourceNotExistCondition`中添加了`anyOf`条件，表示：当调用GetRecord这个API时，如果返回错误码为`InvalidResourceId.NotFound` 
或者GetRecord这个API调用成功，并且当前Record资源的属性`Status`的值为`Deleted`；这两种情况都会都可以被认为当前资源不存在了。
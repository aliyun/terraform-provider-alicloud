# 问题描述
缺少异步配置信息

# 问题表现
异步操作执行后，在进行异步结果是否完成的检查中，在设置的最大超时时间之内，异步操作仍然没有执行结束（结束标准：异步操作成功/失败）


# 解决方案
在CloudSpec IDL的资源与API映射的结构体中添加或订正`asyncPollingByProperty`或`asyncPollingByAPI`，详细语法参考`.docs/common/annotates/resource-annotate.md`中的
`A5.5 defineOperationMapping annotate`部分中的`asyncPollingByProperty`参数和`asyncPollingByAPI`参数，以及`asyncPollingByPropertyStruct`结构和`asyncPollingByAPIStruct`结构


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
   create: @operationMapping(CreateRecord_Mapping_Struct) CreateRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct CreateRecord_Mapping_Struct{
    // ...
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
   create: @operationMapping(CreateRecord_Mapping_Struct) CreateRecord
}

// API与资源映射关系及配置
@defineOperationMapping
struct CreateRecord_Mapping_Struct{
    // ...
    asyncPollingByProperty: [{
        ResourceProperty: "$.Status"
        DelayedTime: 10
        FailedValues: ["PermanentFailure"]
        Interval: 10
        TargetValue: "Available"
        Times: 60
    }]
}
```
修复后的映射结构体中，添加了`asyncPollingByProperty`条目，表示：当调用CreateRecord这个API后，会等待10s，然后开始进行间隔10s、最多重复次数60次的异步检查。
异步检查判断资源属性`Status`，如果`Status`的值变为`Available`，就表示创建成功了，异步结束。
# 问题描述

API调用不符合预期，表现：
1. 调用失败，提示某个参数有问题，或者不存在，比如`The specified disk does not exist in the region.`
2. 调用成功，但是不符合预期，比如差不到任何资源。

一种可能的原因是因为API调用入参格式不对。

# 错误原因

先看API调用时的请求，发现请求参数`diskIds`的值为`d-bp1cxsrc1gzk351ticcf`，但是看API文档示例，正确的API入参值应该为: `["d-bp14k9cxvr5uzy54****"]`，可以看到，正确的入参值是数组，并且要用引号包围。
但是资源属性中`DiskId`就是普通的String，说明资源属性到API的映射配置，需要进行一个动态计算，不能单纯直接映射，需要配置一个Jsonate表达式。

# 解决方案

配置Jsonate表达式，将资源属性`DiskId`动态计算转换后，再映射到API入参`diskIds`

# 示例


原本的映射信息：
```
// 资源定义
resource Record{
   // ...
   create: @operationMapping(Record_create_mapping) CreateRecord
}

// 映射信息
@defineOperationMapping
struct Record_create_mapping {
  resourceAttributeMappings: [
      // 错误的请求映射
      {
          resourceProperty: "$.DiskId"
          requestMappingType: "param"
          requestPathType: "normal"
          requestPath: "DiskIds"
      }
  ]
}

```

修复后的：

```
// 资源定义
resource Record{
   // ...
   create: @operationMapping(Record_create_mapping) CreateRecord
}

// 映射信息
@defineOperationMapping
struct Record_create_mapping {
  resourceAttributeMappings: [
      // 动态计算后的请求映射
      {
          resourceProperty: "$.DiskId"
          // 改为compute
          requestMappingType: "computed"
          // 配置Jsonate表达式
          requestTransform: "[$string(ResourceProperties.DiskId)]"
          requestPathType: "normal"
          requestPath: "diskIds"
      }
  ]
}
```

比如对于上面例子中，`[$string(ResourceProperties.DiskId)]` 表示把资源属性的`DiskId`转换为字符串数组，再赋值给API的入参 `diskIds`。
引用根节点：
```
{
OldResourceProperties:{
-- 资源操作之前的属性（通常情况下无需初始化）
}
ResourceProperties:{
-- 需要操作的资源的属性
}
ApiInput:{
-- 当前API的入参JSON
}
ApiOutput:{
// 这里未获取，不能使用
},
ParentApiInput: {
// 为空，异步场景下才有
},
ParentApiOutput: {
// 为空，异步场景下才有
}
}
```

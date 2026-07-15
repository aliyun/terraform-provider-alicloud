# 问题描述

Get API调用失败，提示某个参数有问题，或者查不到资源，或者查到的资源为空。比如Get调用失败，提示`The specified parameter \"DiskIds\" is not supported.`



# 错误原因

有可能是Get API多映射了一些的请求入参到资源属性，导致的最终会有除了主键外的其他的资源属性被映射到Get的入参，这些参数变成了filter，导致Get API调用失败、或者查不到资源。

# 解决方案

可以尝试删除除了主键和必填参数外的请求映射关系。

# 示例


原本的映射信息：
```
// 资源定义
resource Record{
   // ...
   get: @operationMapping(Record_get_mapping) GetRecord
}

// 映射信息
@defineOperationMapping
struct Record_get_mapping {
  
  resourceAttributeMappings: [
      // 主键映射
      {
          resourceProperty: "$.RecordId"
          requestMappingType: "param"
          requestPathType: "normal"
          requestPath: "RecordId"
      } ,
      // 多余的请求映射
      {
          resourceProperty: "$.DiskId"
          requestMappingType: "param"
          requestPathType: "normal"
          requestPath: "DiskIds"
          responsePathType: "normal"
          responsePath: "$.Disks.Disk[*].DiskId"
      }
  ]

```

修复后的：

```
// 资源定义
resource Record{
   // ...
   get: @operationMapping(Record_get_mapping) GetRecord
}

// 映射信息
@defineOperationMapping
struct Record_get_mapping {
  
  resourceAttributeMappings: [
      // 主键映射
      {
          resourceProperty: "$.RecordId"
          requestMappingType: "param"
          requestPathType: "normal"
          requestPath: "RecordId"
      } ,
      // 删去多余的请求映射
      {
          resourceProperty: "$.DiskId"
          responsePathType: "normal"
          responsePath: "$.Disks.Disk[*].DiskId"
      }
  ]
```

比如对于上面例子中，删除了多余的DiskId的请求映射，注意这样不一定有效，需要尝试，但是绝大多数资源的Get只需要主键就可以。

# 问题描述
调用ListTagResources这个API，提示缺少必填参数ResourceType

# 问题表现
资源的get方法有一个ListTagResources，调用失败，返回错误提示缺少必填参数ResourceType


# 解决方案
在ListTagResources和资源对应的mapping结构体中添加常量入参配置（详细语法参考`.docs/common/annotates/resource-annotate.md`中的`constInputParameters`），ResourceType的常量入参一般是全大写资源名。

# 示例
错误写法如下：
```
// 资源定义
resource Record{
   properties:{
       ResourceType:string
   }
   get: @operationMapping(ListTagResources_Mapping_Struct) ListTagResources
}

// API与资源映射关系及配置
@defineOperationMapping
struct ListTagResources_Mapping_Struct{
    // ...  
}
```

正确写法如下：
```
// 资源定义
resource Record{
   properties:{
       ResourceType:string
   }
   get: @operationMapping(ListTagResources_Mapping_Struct) ListTagResources
}

// API与资源映射关系及配置
@defineOperationMapping
struct ListTagResources_Mapping_Struct{
    // ...  
}
```

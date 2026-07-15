# 问题描述
Create/Update操作后，会用Create/Update操作的测试入参值和操作后的Get到的资源属性值比较，如果存在Diff信息，则表示本次操作没有达到预期值。

# 问题表现
init或者modifies阶段后提示不符合预期错误


# 解决方案
1. 在CloudSpec IDL的资源与API映射的结构体中添加或订正API参数与资源属性的映射信息`resourceAttributeMappings`，详细语法参考`.docs/common/annotates/resource-annotate.md`中的
`A5.5 defineOperationMapping annotate`部分中的`resourceAttributeMappings`参数和`resourceAttributeMappingsStruct`
2. 如果确认Get API不返回这个参数，将这个资源属性定义为操作私有，详细参考`/docs/cloudspec-idl-guide/Resource Annotate.md`中的
`A5.7 rac annotate`部分中的`operatePrivateType`参数


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
    resourceAttributeMappings: [
        // ...
    ]
}
```

按照解决方案1修复后的：
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
    resourceAttributeMappings: [
        // ...
        {
            resourceProperty: "$.Status"
            responsePathType: "normal"
            responsePath: "$.body.Status"
        }
    ]
}
```
修复后的映射结构体中，在`resourceAttributeMappings`条目里新增关于`Status`的映射，表示：将GetRecord这个API的出参的body中的Status参数映射到资源属性`Status`。
这样调用完GetRecord以后就能获得资源属性的Status了。


按照解决方案2修复后：
```
// 资源定义
resource Record{
   // ...
   properties:{
       // ...
       @rac({
           operatePrivateType: ["create"]
       })
       Status:string
   } 
}

```
在资源属性的`Status`上加上了`@rac`注解，标识这个属性是creat操作私有的，仅在create阶段会被使用，get可以不返回，这样CLI执行资源测试时，会忽略这个属性的diff。
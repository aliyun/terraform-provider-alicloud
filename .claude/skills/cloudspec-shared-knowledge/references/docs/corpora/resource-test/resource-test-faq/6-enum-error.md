# 问题描述
资源属性枚举值缺失，和API的可传入枚举对不上

# 问题表现
cli执行失败，提示`InvalidParameter.EnumCheckFailure`


# 解决方案
在对应的资源属性上补充枚举，详细语法参考`.docs/common/annotates/assignment-constraint-annotate.md`中的`A1.6 enums annotate`，和
`.docs/common/basic-grammer/basic-types.md`中的`enum`



# 示例
原本的的CloudSpec IDL结构体以及资源测试用例，由于测试用例中`Status`属性的值为`Modifing`，所以需要补充枚举。
```
// 资源定义
resource Record{
   // ...
   properties:{
       // ...
       @enums(Record_Status_Enum)
       Status:string
   }
   
}

// 枚举定义
enum Record_Status_Enumm {
  "Created"
  "Deleted"
}


// 测试用例
$test Record test_Record{
    init:{
        Status:"Modifing"
    }
}

```

修复后的：
```
// 资源定义
resource Record{
   // ...
   properties:{
       // ...
       @enums(Record_Status_Enum)
       Status:string
   }
   
}

// 枚举定义
enum Record_Status_Enumm {
  "Created"
  "Deleted"
  "Modifing"
}


// 测试用例
$test Record test_Record{
    init:{
        Status:"Modifing"
    }
}
```
修复后的`Record_Status_Enumm`枚举定义，添加了`Modifing`，这样CLI就不会报错了。

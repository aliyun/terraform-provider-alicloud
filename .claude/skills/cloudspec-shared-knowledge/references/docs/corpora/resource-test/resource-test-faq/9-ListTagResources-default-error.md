# 问题描述

调用ListTagResources这个API，提示The specified parameter Tag is not support for default resource group.

# 问题表现

资源的get方法有一个ListTagResources，调用失败，返回错误提示The specified parameter Tag is not support for default resource
group.

# 问题原因

创建资源的时候没有指定资源组，则资源归属于默认资源组，默认资源组不支持标签操作。

# 解决方案

创建一个资源组的前置依赖，并且在资源测试用例的init阶段引用资源组ID。一个示例如下：

原本测试用例如下：

```
@testConfig({
  main: true
  // ...
})
$test Record test_Record{
    init:{
        Type: "AAAA"
    }
}
```

订正后的测试用例如下：

```

$test alicloud.ResourceManager.ResourceManager.v20200331#ResourceGroup resource_ResourceGroup_test {
  init: {
    DisplayName: "测试资源组"
    ResourceGroupName: "test-rg"
  }
  destroy: {
  }
}


@testConfig({
  main: true
  // ...
})
$test Record test_Record{
    init:{
        Type: "AAAA"
        ResourceGroupId: "{{$.resource_ResourceGroup_test.ResourceGroupId}}"
    }
}
```

# 注意事项

需要注意的是`resource_ResourceGroup_test`这个名字需要全局唯一，不能重复，所以具体叫什么名字就需要根据实际场景来定，
建议命名为`{依赖的测试名称}_test_dependency`，比如本例中可以命名为`test_Record_test_dependency`。

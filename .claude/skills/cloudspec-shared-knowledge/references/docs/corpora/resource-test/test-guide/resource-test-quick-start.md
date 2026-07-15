# Aliyun CLI cspec plugin中执行资源测试说明书

## 测试用例语法说明

测试定义是CloudSpec IDL定义的一部分，每个测试用例都是一个`.cspec`文件，每个测试用例包含一个或多个`$test`。`$test`间可以互相依赖，CLI执行引擎会自动解析测试对象间依赖关系，构建运行图，按顺序调度执行。一个具体的示例如下：

```
@testConfig({
  execConfigUuid: "bac6d542870*****2374e7a6d2"
  main: true
  runtime: {
    regionId: "cn-hangzhou"
    endpoint: "esa.[RegionId].aliyuncs.com"
  }
})
$test Record record_test {
  init: {
    SiteId:{{$.record_pre_Site.SiteId}}
    RecordName: "music.changes.com.cn"
    RecordType: "CNAME"
    Data: {
      Value: "music.alyun.com"
    }
    Ttl: 10
    Comment: "test"
    Proxied: true
  }
  modifies: [{
    Comment: "test11"
  }, {
    Comment: "test22"
  }]
  destroy: {}
}


// 前置依赖测试
$test Site record_pre_Site {
  init: {
    SiteName: "mytest"
  }
  destroy: {}
}
```

### testConfig注解
`@testConfig`是测试运行的配置注解。
- execConfigUuid是测试执行账号，如果未配置，会从本地环境量读取。
- main表示是否是主测试对象，如果为true，则表示当前`$test`是主测试对象，其余的`$test`都是前置依赖测试对象。一个测试用例(`.cspec`文件)只能有一个main test。
- runtime表示测试运行时环境，包含endpoint信息等。
- 通常只有main test需要配置`@testConfig`。


### $test 头部
- 如上面示例，`$test`为固定语法标记，`Record`是资源名称，表示测试的对象，`record_test`是这个`$test`对象的名称，必须全局唯一，不可重名。

### $test 主体
- `init`是资源的初始化步骤，表示初始化一个具有声明属性的资源。
- `modifies`是资源的修改步骤，支持数组，每个步骤支持传入资源支持的可改属性，从上到下依次执行。
- `destroy`是资源的销毁步骤，支持传入销毁的操作私有参数，但是通常不应该这么设计OpenAPI。

### $test间的相互依赖
`$test`支持依赖其他`$test`，例如上面的例子，`record_test`的init步骤中需要依赖`record_pre_Site`的`SiteId`属性。引用的语法为`{{$.前置测试名称.前置资源属性JSONPATH}}`。

也可以依赖跨空间下的资源，例如前置依赖VPC时：

```
$test alicloud.VPC.Vpc.v20160428#VPC vpc_pre {
  CidrBlock: "10.0.0.0/8"
  Description: "test"
  Ipv6Isp: "BGP"
}

@testConfig({
  // ...
  main :true
})
$test Record record_test {
  init: {
    VPCId: "{{$.vpc_pre.VpcId}}"
  }
}

```

注意，如果依赖的资源不在当前空间，则需要指定namespace，如上面vpc_pre例子，VPC这个资源需要指定namespace，下面的写法是错误的：
```
$test Vpc vpc_pre {
  CidrBlock: "10.0.0.0/8"
  Description: "test"
  Ipv6Isp: "BGP"
}
```

## 执行测试用例
在当前CLoudSpec IDL的项目运行

```
aliyun cspec test run -n {测试名}
```

注意`-n `参数后面跟的是main test的名称，比如对于上面例子，运行命令应该为：
```
aliyun cspec test run -n record_test
```
# 问题描述
测试用例缺失前置依赖资源的创建

# 问题表现
测试用例某个属性（往往是xxxId）依赖另一个资源的主键，但是这个属性在当前测试用例里面直接硬编码（往往是假的），导致当前资源创建失败


# 解决方案
1. 完善资源测试用例，创建前置依赖，并在测试用例中引用。比如对于下面例子，就是很经典的前置依赖缺失，DnsRecord资源创建需要依赖Site资源主键，不能直接写死。
2. 前置依赖可以从资源的`@references`注解获取，如果从`@references`注解也获取不到，可以从当前namespace下资源找一下，做合理推断，还是不行就放弃该资源，但是要明确记录下来。
3. 前置依赖的$test的写法可以去`/tests`目录下找，看看依赖资源有没有测试用例可以照抄，抄的时候注意命名以及依赖的依赖的命名，不要怕麻烦！


# 示例
原本的测试用例
```
$version: 1
namespace: alicloud.ESA.ESA.v20240910

$test DnsRecord test_record{
  init:{
    SiteId: "site-ahjshasasiasjakjsoas"
  }
}
```

正确写法如下：

## 假设Site是当前namespace下的资源
```
$version: 1
namespace: alicloud.ESA.ESA.v20240910


$test Site my_site{
  init:{
    SiteName: "myTestSite" 
  }
}


$test DnsRecord test_record{
  init:{
    SiteId: "{{$.my_site.SiteId}}"
  }
}
```

## 假设Site是其他namespace下的资源
```
$version: 1
namespace: alicloud.ESA.ESA.v20240910


$test alicloud.VPC.Vpc.v20160428#Site my_site{
  init:{
    SiteName: "myTestSite" 
  }
}


$test DnsRecord test_record{
  init:{
    SiteId: "{{$.my_site.SiteId}}"
  }
}
```
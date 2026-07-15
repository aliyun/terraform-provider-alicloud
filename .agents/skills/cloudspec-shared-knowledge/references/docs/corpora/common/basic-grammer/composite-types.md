### struct
#### 标准定义
定义结构体，结构体需要显式定义全部的keys集合及每个value的类型。

```json
struct MyCustomExampleStruct {
  A1: int32
  A2: string
  StructExample: StructExampleStruct
}

struct StructExampleStruct {
  B1: string
  B2: int64
}
```

#### extend
struct支持extend语法简化表达，如果当前结构体和extend中存在同名的属性，则以当前定义的属性为准，例如：

```json
struct StructExampleStruct {
  B1: string
  B2: int64
}

struct StructExampleStructExample extend StructExampleStruct {
  B3: string
  B4: int64
}
```

最终StructExampleStructExample的字段等同于：

```json
struct StructExampleStructExample{
  B1: string
  B2: int64
  B3: string
  B4: int64
}
```

#### 内联定义
不推荐在struct中重复定义struct，这样无法在选择器中灵活的选择到每一个层级的定义。例如：

```json
struct RootStruct {
  A: {
    B1: string
    B2: int32
  }
}
```

> 使用{}默认包裹的对象默认为struct。
>

struct 也支持内联模式定义，固定语法如下：

```bash
struct a {
  // c以下的结构都是内联的 struct 定义
  c: {
    d: {
      e: string
    }
  }
}
```

#### 特殊 key 转义
如果 struct 的 key 包含了 service、resource 等[语言保留字](https://aliyuque.antfin.com/cloudspec/model/qm0s6ccgbg0p3zal)，请使用````转义，例如：

```bash
struct A {
  `resource`: string
   // 极少出现，慎用
  `${{a.b}}`: int32
}
```

### array
#### 标准定义
数组，数组是同类型元素的集合，array类型下的值必须是相同的类型。

```json
array MyArray {
  item: string
}
```

其中item是固定写法，表示array下的子节点。

#### 内联定义
也支持使用内联的模式定义数组，内联模式定义数组格式如下：

```bash
array<TYPE>
```

其中<>表示数组 item 节点的类型，具体示例：

```bash
struct a {
  // b是一个数组，其数组的节点为 string
  b: array<string>
}

struct c {
  // c1是一个数组，其item 节点是一个 struct，struct 存在一个条目d，类型是 Boolean
  c1: array<{
    d: boolean
  }>
}
```

### map
#### 标准定义
地图类型，map类型的keys集合是无限制的。以下示例定义一个键值为string，其值也是string的map结构。

```json
map MapExample {
  key: string
  value: string
}
```

> Warning：应尽量少的使用Map结构，Map结构会带来众多的不确定性，会破坏集成的体验。
>

#### 内联定义
map 也支持使用内联模式定义，定义的语法为：

```bash
map<TYPE>
```

 其中 TYPE 表示 map 的 value 的类型，当前 key 的类型固定等于 string，因此不用单独定义。

具体示例如下：

```bash
struct a {
  // c是一个 map，其 map 的 value 值类型为任意类型。
  c: map<any>
}
```


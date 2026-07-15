传值约束，用于约束模型和后端双向交互传值的约束。

### A2.1 **required annotate**
#### 说明
约束结构中的字段在入参、出参中为必填。除了显式的标注外，_默认的结构都是非必填的_。

#### 选择器
```json
:is(struct > member)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无参数 | | |


#### 示例
```json
struct MyStructure {
    @required
    foo: string
}
```

### A2.2 client**Optional annotate**
#### 说明
约束结构中的字段在入参为客户端非必填，默认值为false，表示客户端必填，存在则表示该资源客户端非必选。

:::warning
+ 字段标记为required，且为客户端可选clientOptional，这说明这个字段的赋值是在服务端生成的，且一定返回值。
+ 字段没有标记为required，且标记了clientOptional，则说明这个字段可以传值也可以不传，不一定存在返回值。

:::

#### 选择器
```json
:is(struct > member)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无参数 | | |


#### 示例
```json
struct MyStructure {
    @required
    @clientOptional
    foo: string
}
```

### A2.3 clientProhibited** annotate**
#### 说明
默认值false，表示该字段客户端允许传值。当设置了该注解时，则不允许客户端传递此值。

:::warning
注意，readonly是指字段不可改，可以在create阶段传值。

如果同时被标记为clientProhibited，则说明这个值只能从服务端生成，且不可修改。

:::

#### 选择器
```json
:is(struct > member)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无参数 | | |


#### 示例
```json
struct MyStructure {
    @required
    @clientProhibited
    foo: string
}
```

### A2.4 **conflictsWith annotate**
#### 说明
约束传值时和其他的字段不能同时传递。

#### 选择器
```json
:is(struct > member)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [string] | 和其他字段冲突的component ID值，当前struct写相对路径即可。 |


#### 示例
```json
struct MyStructure {
  @conflictsWith(['Id'])
  Name: string
  @conflictsWith(['Name'])
  Id: string
}
```

该示例说明对于MyStructure对象，Name和Id是冲突的，在传值的时候只能选择其中一个。



### A2.5 idempotencyToken** annotate**
#### 说明
说明参数是幂等参数。

#### 选择器
```json
:is(struct > member :test(> string))
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无需参数 | | |


#### 示例
```json
struct MyStructure {
  @conflictsWith(['Id'])
  Name: string
  @conflictsWith(['Name'])
  Id: string
  @idempotencyToken
  ClientToken: string
}
```

### A2.6 readonly** annotate**
#### 说明
标记修饰的对象只读，当作用不同的component时，含义不同：

+ 修饰基本数据类型，表示该字段不可被修改；
+ 修饰operation时，表示该操作不允许修改任何资源属性。

#### 选择器
```json
:is(struct > member, array > member, map > member, simpleType, map, array, struct, operation)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | Boolean |  非必填，默认值为 true。 |


#### 示例
```json
struct MyStructure {
  @readonly
  Name: string
  @conflictsWith(['Name'])
  Id: string
  @idempotencyToken
  ClientToken: string
}
```

:::warning
用readonly描述Name，表示Name是一个只读的属性，这个属性只允许在创建阶段指定，在资源创建后字段不允许被修改。

:::

### A2.7 hasDefaultValue** annotate**
#### 说明
是否存在默认值注解，属性存在这个注解时，当客户端默认不传值时，服务端会生成一个默认值。

例如CreateTime，这个字段客户端不能传递，但服务端会默认生成一个时间错。

#### 选择器
```json
:is(struct > member, array > member, map > member, simpleType, map, array, struct)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| 无需参数 | | |


#### 示例
```json
struct MyStructure {
  @readonly
  @hasDefaultValue
  CreateTime: string
}
```

> readonly指这个参数只读，除了创建阶段不可修改。
>

### A2.8 date** annotate**
#### 说明
表示当前属性是时间，支持配置时间的格式。

#### 选择器
```json
:is(struct > member, array > member, map > member, simpleType, map, array, struct)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| format | string | 完整日期和时间格式：YYYY-MM-DDTHH:mm:ss.sssZ  YYYY：四位数的年份  MM：两位数的月份  DD：两位数的日期  T：日期和时间之间的分隔符  HH：两位数的小时  mm：两位数的分钟  ss：两位数的秒钟  sss：三位数的毫秒（可选）  Z：时区信息（例如，+08:00或Z表示UTC） |


#### 示例
```json
struct MyStructure {
  @date({
    format: "YYYY-MM-DDTHH:mm:ss.sssZ"
  })
  CreateTime: string
}
```

### A2.9 arrayConfig **annotate**
#### 说明
数组的额外配置。

#### 选择器
```json
:is(array)
```



#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| uniqueItems | Boolean | 是否校验数组下的值是唯一的，默认值为false。 |
| unordered | Boolean | 是否为无序的数组，表示返回值中数组元素的值可能会变换顺序，默认值为false，表示有序数组。 |


#### 示例
```json
@arrayConfig({
  uniqueItems: true
  unordered: true
})
array arrayTest {
  item: string
}
```


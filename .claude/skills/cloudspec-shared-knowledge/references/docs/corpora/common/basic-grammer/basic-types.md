## 契约类型
数据类型定义以小写的字符开头，例如 string，boolean，integer。

### string
string，表示变量A是一个字符类型，字符类型不限制长度。

```json
struct Candy {
	Name: string
}
```



也可以使用string 声明类型：

```json
string CustomString

struct CustomStruct {
  // 这里CustomString变成一个单独的类型
  Name: CustomString
}
```

### byte
表示经过base64编码的字节流。

### binary
字节流

### boolean
布尔值。

### int32
整形，32位有符号整数，-2^31 到 (2^31)-1 。

### int64
长整形，32位有符号整数，-2^63 到 (2^63)-1 。

### float
IEEE-754规定的单精度浮点数。

### double
IEEE-754规定的双精度浮点数。

### enum
定义枚举值，枚举值类型是string，如果枚举值是数字需要使用intEnum定义。

```json
enum EnumType {
  "PostPaid"
  "PrePaid"
}
```

### intEnum
定义数字型枚举值，

```json
intEnum IntEnumExample {
  1
  2
  3
}
```

注意，使用以上这些关键字再次定义的类型都是自定义类型，这些类型可以直接被二次引用。

### any
表示任意类型，例如：

```yaml
struct S {
  a: any
}
```

a字段为任意类型。

## 值类型
IDL 中绝大场景均是契约定义，声明一个字段可接收的值类型。在 annotate 或者测试场景，需要使用值定义，IDL 不仅可以进行契约定义，也可以进行值定义。

### 字符
字符类型使用用单引号或者双引号包裹，例如：

```yaml
struct a {
  a1: "1"
  a2: '1'
}
```

### 数值
数值不允许使用引号包裹，允许定义整形、浮点，例如：

```yaml
struct a {
  a1: 1
  a2: 1.0
}
```

### 布尔
boolean类型支持 true、false 两个值，不需要引号包裹，例如：

```bash
struct a {
  a1: true
  a2: false
}
```

### 数组
值定义允许使用数组，数组使用[]表达，例如：

```yaml
struct a {
  a1: [1, 2]
  a2: ["1", "3"]
  a3: [true, false]
}
```

### null
允许使用$null定义值类型，例如：

```yaml
struct a {
  a1: $null
  a2: {
    a21: $null
  }
}
```


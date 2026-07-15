CloudSpec IDL定义了一些保留字，列表如下：

```bash
struct
map
array
resource
operation
service
error
string
byte
binary
boolean
double
float
true
false
constraint
int32
int64
enum
intEnum
extend
import
$version
namespace
any
$apply
```

当在结构中使用上述关键字时，需要使用 ```` 转义，语法如下：

```json
struct CreateDiskReplicaGroupOutput {
  @required
  @backendParameterName({responseName: 'groupId'})
  $ReplicaGroupId

  @notResourceProperty
  @backendParameterName({responseName: 'requestId'})
  RequestId: string
  `error`: string
  `struct`: string
  `map`: string
  `string`: hjahah
  RequestId2: 2
  `true`: 333
  `false`: 222
  `boolean`: xx
  `enum`: xxx
}
```

另外`$xxx`会用于后续模型扩展的关键字定义，请勿在IDL中直接使用$开头的命名。


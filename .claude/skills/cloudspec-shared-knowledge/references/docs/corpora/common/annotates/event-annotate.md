# CloudSpec 事件类注解列表

本文档列出了 CloudSpec 事件类相关的所有注解及其详细配置说明，用于描述流式传输事件的配置信息。

## 目录

- [A15.1 eventConfig annotate](#a151-eventconfig-annotate)
- [A15.2 gatewayTransportType annotate](#a152-gatewaytransporttype-annotate)
- [A15.3 events annotate](#a153-events-annotate)

### A15.1 eventConfig annotate
#### 说明
描述流式传输事件配置信息，使用默认值不会生成cspec文件对应字段或注解。

#### 选择器
```cspec
:is(struct)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| errorEvent | boolean | 是否标记为错误码事件（default false） |
| bodyContentType | string | 事件的是结构化的还是二进制的（default json）json/binary |
| eventId | string | 消息标识符。用于标识事件 |
| ackEventId | string | ack的响应事件 |

#### 示例：定义一个事件结构
```cspec
@eventConfig({
    errorEvent: false
    bodyContentType: "json"
    ackEventId: "FinsihedEvent"
    eventId: "CreateDBEvent"
})
@gatewayTransportType("proxy")
struct event1 {
    @required
    roomId: string
    @required
    content: string
}
```

#### 示例：全文件
```cspec
$version: 1
namespace: alicloud.EBS

service chatService {
  version: "1.0.0"
  operations:[Connect]
}

@http({
  schemes: {
    online: ['https', 'websocket']
  }
  methods: ["get"]
  authenticators: ["AK"]
  requestContentType: ["application/json"]
  responseContentType: ["application/json"]
  uri: "/connect"
  deprecated: false
  protocolOverWebsocket: "awap"
})
operation Connect {
  @events([event1])
  input: {
    @required
    roomId: string
    @required
    content: string
  }

  @events([event2])
  output: {
  }
}

@eventConfig({
    errorEvent: false
    bodyContentType: "json"
    ackEventId: "FinsihedEvent"
    eventId: "CreateDBEvent"
})
@gatewayTransportType("proxy")
struct event1 {
    @required
    roomId: string
    @required
    content: string
}

@eventConfig({
    errorEvent: false
    bodyContentType: "json"
    ackEventId: "FinsihedEvent"
    eventId: "DeleteDBEvent"
})
@gatewayTransportType("proxy")
struct event2{
    @required
    ok: string
}
```

### A15.2 gatewayTransportType annotate
#### 说明
描述流式传输事件配置信息，需要配合@eventConfig使用。

#### 选择器
```cspec
:is(struct)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | string | 可填内容：proxy/mapping |

#### 示例
```cspec
@eventConfig({
    errorEvent: false
    bodyContentType: "json"
    ackEventId: "FinsihedEvent"
    eventId: "CreateDBEvent"
})
@gatewayTransportType("proxy")
struct event1 {
    @required
    roomId: string
    @required
    content: string
}
```

### A15.3 events annotate
#### 说明
将流式传输事件与传输方向：上行（input）或下行（output）绑定。

#### 选择器
```cspec
:test(operation -[input,output]->)
```

#### 支持的属性
| 属性 | 类型 | 说明 |
| --- | --- | --- |
| - | [struct] | struct类型的引用的数组，聚合需要绑定在input或output的struct |

#### 示例
```cspec
operation Connect {
  @events([event1])
  input: {
    @required
    roomId: string
    @required
    content: string
  }

  @events([event2])
  output: {
  }
}
```

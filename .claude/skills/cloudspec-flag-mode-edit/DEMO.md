# cloudspec-flag-mode-edit Skill 效果展示

本文档通过 5 个真实场景展示 `cloudspec-flag-mode-edit` skill 的能力。每个场景包含**用户输入**、**skill 执行过程**和**最终输出**。

---

## 场景一：从零搭建 Flag 模式项目

### 用户输入

> 帮我创建一个 Flag 模式的 CloudSpec 项目，产品名 CloudMonitor，资源是 AlarmRule（告警规则），包含名称、描述、状态、告警等级、创建时间等字段。

### Skill 执行过程

1. **识别任务类型**：新建 Flag 模式项目
2. **阅读参考文档**：`flag-mode-usage-guide.md`（完整项目结构示例）、`flag-mode-reference.md`（flag 选择指南）
3. **确定字段 flag 组合**：
   - `AlarmRuleId`：主键，`#[R,D,F,L]` + `@readonly`
   - `Name`：用户可编辑，`#[C,U,R,F,L]` + `@required`
   - `Description`：用户可编辑，`#[C,U,R,F,L]`
   - `Status`：只读系统字段，`#[R,F,L]` + `@readonly`
   - `Level`：创建后不可改，`#[C,R,F,L]` + `@required`
   - `CreateTime`：只读，`#[R,L]` + `@readonly`
4. **生成完整项目**：main.cspec + resources/AlarmRule.cspec + 5 个操作文件

### 最终输出

**resources/AlarmRule.cspec**：

```cspec
$version: 1
namespace: alicloud.CloudMonitor.CloudMonitor.v20240101

@arn("acs:cloudmonitor:${Region}:${AccountId}:alarmrule/${AlarmRuleId}")
@document({
  zh: "告警规则"
  en: "Alarm Rule"
  name: "告警规则"
  nameEn: "Alarm Rule"
})
@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource AlarmRule {
  identifyDefinition: {
    @readonly
    AlarmRuleId: string #[R,D,F,L]
  }
  properties: {
    @readonly
    RegionId: string
    @required
    @length({ min: 2, max: 128 })
    Name: string #[C,U,R,F,L]
    Description: string #[C,U,R,F,L]
    @readonly
    @enums(AlarmRuleStatusEnum)
    @initialStatus("Pending")
    Status: string #[R,F,L]
    @required
    Level: string #[C,R,F,L]
    @readonly
    @format("iso8601")
    CreateTime: string #[R,L]
  }
}

enum AlarmRuleStatusEnum {
  "Pending"
  "Running"
  "Disabled"
}
```

**operations/CreateAlarmRule.cspec**：

```cspec
$version: 1
namespace: alicloud.CloudMonitor.CloudMonitor.v20240101

@for(AlarmRule)
@backendConfigurationHttp({
  applicationName: "CloudMonitor"
  retries: { online: -1 }
  timeout: { online: 5000 }
  backendUrl: { online: "http://vpc_online/api/alarmrule/create#vpc" }
  sign: true
  signPolicy: "Local"
})
@document({ name: "创建告警规则" })
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "create"
  riskType: "none"
  chargeType: "free"
})
@ram({ enable: true, level: "operate", atGateway: false })
@visibility("Public")
operation CreateAlarmRule {
  errors: [Error_CreateAlarmRule]
}

error Error_CreateAlarmRule {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

### 体现的 Skill 能力

- ✅ 根据字段语义自动选择正确的 flag 组合（区分 F 和 L）
- ✅ 生成完整的项目结构（main + resources + operations）
- ✅ 正确使用 `@flagMode` + `@for(Resource)` 模式
- ✅ 操作体内只有 `errors`，不自定义 input/output

---

## 场景二：为已有资源添加新字段并设置 Flag

### 用户输入

> 在 AlarmRule 资源中添加一个 Tags 字段（标签数组），支持创建和更新，查询和列表都返回，同时可以作为列表过滤条件。再添加一个 ClientToken 字段，仅创建时传入，作为幂等参数。

### Skill 执行过程

1. **识别任务类型**：添加资源属性（含 flag 标记）
2. **阅读参考文档**：`flag-mode-reference.md`（flag 选择指南、复合类型 flag 语法）
3. **确定字段 flag 组合**：
   - `Tags`：可编辑 + 可过滤 + 列表输出 → `#[C,U,R,F,L]`
   - `ClientToken`：仅创建 + 幂等 → `#[C]` + `@idempotencyToken`
4. **读取当前资源文件**，在 properties 末尾添加字段

### 最终输出

在 `resources/AlarmRule.cspec` 的 properties 中添加：

```cspec
    @document({
      name: "标签"
      zh: "资源标签"
      en: "Resource tags"
    })
    Tags: array<
      @document({ name: "标签", zh: "标签", en: "Tag" })
      @arrayConfig({ unordered: true })
      {
        @required
        @document({ name: "标签键", zh: "标签键", en: "Tag key" })
        Key: string
        @document({ name: "标签值", zh: "标签值", en: "Tag value" })
        Value: string
      }
    > #[C,U,R,F,L]
    @idempotencyToken
    @rac({
      operatePrivateType: ["create"]
    })
    ClientToken: string #[C]
```

### 体现的 Skill 能力

- ✅ 复合类型（array）的 flag 标记放在 `>` 闭合符号之后
- ✅ 正确使用 `@idempotencyToken` 和 `@rac({ operatePrivateType })` 标记操作私有属性
- ✅ 区分 F（过滤入参）和 L（列表输出），Tags 同时需要两者

---

## 场景三：修改 Flag 标记 + 理解 F 和 L 的区别

### 用户输入

> 我有一个资源 Instance，里面有个 InternalIp 字段当前标记是 `#[C,R,L]`。我希望它不在列表中返回（列表太长了），但仍然可以作为列表的过滤条件来搜索。另外 Status 字段当前是 `#[R,L]`，我希望它也能作为列表过滤条件。

### Skill 执行过程

1. **识别任务类型**：修改字段 flag 标记
2. **阅读参考文档**：`flag-mode-reference.md`（F 和 L 的区别）
3. **分析需求**：
   - `InternalIp`：不在列表输出（去掉 L）+ 可过滤（加 F）→ `#[C,R,F,!L]`
   - `Status`：保留列表输出（L）+ 加过滤（F）→ `#[R,F,L]`
4. **读取当前资源文件**，修改对应字段的 flag

### 最终输出

```cspec
// 修改前
InternalIp: string #[C,R,L]
@readonly
Status: string #[R,L]

// 修改后
InternalIp: string #[C,R,F,!L]
@readonly
Status: string #[R,F,L]
```

### 体现的 Skill 能力

- ✅ 精确区分 F（FILTER_SUPPORT，List 过滤入参）和 L（LIST_OUTPUT_SUPPORT，List 输出）
- ✅ 正确使用否定标记 `!L` 显式排除 List 输出
- ✅ 只修改用户指定的部分，保留其他注解和结构不变

---

## 场景四：将非 Flag 模式资源转换为 Flag 模式

### 用户输入

> 帮我把这个资源转换为 Flag 模式：
> ```cspec
> resource Bucket {
>   identifyDefinition: {
>     @readonly
>     BucketName: string
>   }
>   properties: {
>     @required
>     StorageClass: string
>     @readonly
>     CreationDate: string
>     Acl: string
>     @readonly
>     RegionId: string
>   }
>   create: CreateBucket
>   get: GetBucket
>   delete: DeleteBucket
>   list: ListBuckets
> }
> ```

### Skill 执行过程

1. **识别任务类型**：转换为 Flag 模式
2. **阅读参考文档**：`flag-mode-reference.md`（flag 选择指南、默认 flag）、`flag-mode-usage-guide.md`（转换流程）
3. **分析每个字段的语义**（根据原始操作列表 create/get/delete/list，无 update 操作）：
   - `BucketName`：主键 → `#[R,D,F,L]`
   - `StorageClass`：创建后不可改（无 Update 操作）→ `#[C,R,F,L]`
   - `CreationDate`：只读 → `#[R,L]`
   - `Acl`：无 Update 操作，仅创建时设定 → `#[C,R,F,L]`
   - `RegionId`：无 flag（平台特殊处理）
4. **添加 `@flagMode` 注解**
5. **修改操作定义**：添加 `@for(Bucket)`，移除 input/output
6. **移除资源中的生命周期操作关联**（create/get/delete/list）

### 最终输出

```cspec
@flagMode
@resourceBaseInfo({
  classification: "normal"
  deliveryScope: "region"
  paidType: "Free"
})
resource Bucket {
  identifyDefinition: {
    @readonly
    BucketName: string #[R,D,F,L]
  }
  properties: {
    @readonly
    RegionId: string
    @required
    StorageClass: string #[C,R,F,L]
    @readonly
    CreationDate: string #[R,L]
    Acl: string #[C,R,F,L]
  }
}
```

操作文件转换示例（CreateBucket.cspec）：

```cspec
// 转换前
operation CreateBucket {
  input: CreateBucketInput
  output: CreateBucketOutput
  errors: [Error_CreateBucket]
}

// 转换后
@for(Bucket)
@document({ name: "创建存储桶" })
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "create"
  riskType: "none"
  chargeType: "free"
})
operation CreateBucket {
  errors: [Error_CreateBucket]
}
```

### 体现的 Skill 能力

- ✅ 完整的转换流程：添加 @flagMode → 标记 flag → 修改操作 → 清理 struct
- ✅ 正确识别字段语义并选择 flag（无 Update 操作时，所有可写字段均为 `#[C,R,F,L]`）
- ✅ RegionId 不加 flag（平台特殊处理）
- ✅ 操作转换：添加 @for，移除 input/output，只保留 errors

---

## 场景五：为 Flag 模式操作覆盖注解 + 使用 A flag

### 用户输入

> 我有一个 Flag 模式资源 Order，Create 操作需要自定义后端配置（超时 10 秒、特殊的后端 URL），并且创建后需要额外返回 OrderNo（订单号，服务端生成）。Delete 操作需要标记为高危操作。

### Skill 执行过程

1. **识别任务类型**：编辑操作注解 + 添加资源属性
2. **阅读参考文档**：`flag-mode-reference.md`（A flag、操作注解覆盖）、`flag-mode-usage-guide.md`（操作覆盖示例）
3. **分析需求**：
   - `OrderNo` 字段：服务端生成 + 创建后返回 → `#[R,L,A]` + `@readonly`
   - Create 操作：覆盖 `@backendConfigurationHttp`
   - Delete 操作：`riskType: "high"`（代码实现中 Delete 默认就是 high）
4. **修改资源文件**：添加 OrderNo 字段
5. **修改操作文件**：覆盖注解

### 最终输出

**resources/Order.cspec 中添加字段**：

```cspec
    @readonly
    @document({
      name: "订单号"
      zh: "服务端生成的订单号"
      en: "Server-generated order number"
    })
    OrderNo: string #[R,L,A]
```

**operations/CreateOrder.cspec**（覆盖后端配置）：

```cspec
$version: 1
namespace: alicloud.OrderService.OrderService.v20240101

@for(Order)
@backendConfigurationHttp({
  applicationName: "OrderService"
  retries: { online: 3 }
  timeout: { online: 10000 }
  backendUrl: { online: "http://order-vpc.internal/api/order/create#vpc" }
  sign: true
  signPolicy: "Local"
})
@document({ name: "创建订单" })
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "create"
  riskType: "none"
  chargeType: "paid"
})
@ram({ enable: true, level: "operate", atGateway: false })
@visibility("Public")
operation CreateOrder {
  errors: [Error_CreateOrder]
}

error Error_CreateOrder {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}
```

**operations/DeleteOrder.cspec**（高危操作）：

```cspec
$version: 1
namespace: alicloud.OrderService.OrderService.v20240101

@for(Order)
@backendConfigurationHttp({
  applicationName: "OrderService"
  retries: { online: -1 }
  timeout: { online: 5000 }
  backendUrl: { online: "http://order-vpc.internal/api/order/delete#vpc" }
  sign: true
  signPolicy: "Local"
})
@document({ name: "删除订单" })
@errorMapping({
  errorExpression: "success=false"
  codeField: "code"
  errorMessageField: "message"
  httpStatusCodeField: "httpStatusCode"
})
@http({
  schemes: { online: ["https"] }
  methods: ["post"]
  authenticators: ["AK"]
  deprecated: false
})
@operationInfo({
  operationTypeOld: "write"
  typeFromOperation: "delete"
  riskType: "high"
  chargeType: "paid"
})
@ram({ enable: true, level: "operate", atGateway: false })
@visibility("Public")
operation DeleteOrder {
  errors: [Error_DeleteOrder, Error_DeleteOrder_NotFound]
}

error Error_DeleteOrder {
  httpCode: 400
  errorCode: "InternalException"
  errorMessage: "Failed to call the service."
  type: "user"
  default: true
}

error Error_DeleteOrder_NotFound {
  httpCode: 404
  errorCode: "Order.NotFound"
  errorMessage: "The specified order does not exist."
  type: "user"
  default: false
}
```

**系统自动生成的 Create 输出**（基于 flag 推导）：

```
CreateOrder_Output = {
  OrderId: string      // identifyDefinition 中的 ID
  OrderNo: string      // A flag 标记的额外返回字段
  RequestId: string    // 系统字段
}
```

### 体现的 Skill 能力

- ✅ 正确使用 A (CREATE_AFTER_OUTPUT) flag 标记创建后额外返回的字段
- ✅ 使用 `@for(Resource)` 覆盖操作注解，不破坏其他操作的自动推导
- ✅ Delete 操作 riskType 设为 "high"（与代码实现一致）
- ✅ chargeType 根据 paidType 推断（非 Free → "paid"）
- ✅ 理解 Create 输出的组成：ID + A 标记字段 + RequestId

---

## 能力总结

| 能力维度 | 覆盖场景 |
|---------|---------|
| **项目脚手架** | 场景一：从零搭建完整 Flag 模式项目 |
| **字段编辑** | 场景二：添加复合类型字段 + 操作私有属性 |
| **Flag 精确控制** | 场景三：区分 F 和 L、使用否定标记 |
| **模式转换** | 场景四：非 Flag 模式 → Flag 模式完整转换 |
| **操作覆盖** | 场景五：自定义注解 + A flag + 高危操作 |
| **代码实现对齐** | 所有场景：flag 定义、映射规则、参数位置均与代码实现一致 |

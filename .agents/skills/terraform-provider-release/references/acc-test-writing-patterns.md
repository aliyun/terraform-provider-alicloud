# ACC 测试写作 patterns

针对 Step 9 手写测试（generator 没有可用镇元测试用例，或复杂依赖需要人为设计）时容易踩的坑。本文只覆盖**跨资源共性**，硬编码规范另见 `provider-resource-review/SKILL.md` §6。

## 1. 复杂前置资源：AI 网关这类需要 VPC/VSwitch/Gateway 的场景

**问题**：AiModelProvider 需要 AI 类型的 APIG Gateway，Gateway 又需要 VPC + 2 个 AZ 的 vswitch。全部 inline 创建的话，测试结束 destroy 阶段大概率被上游资源的清理 bug 阻塞（e.g. APIG Gateway 自动创建的 SG 没级联清理 → `DependencyViolation.SecurityGroup` 挡住 VPC 删）。

**patterns**：

### 1.1 优先用共享 `default-NODELETING` VPC

框架不会尝试删 data source 引用的资源，SG 遗留就不阻塞测试收尾：

```hcl
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "j" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

data "alicloud_vswitches" "k" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-k"
}
```

`cn-hangzhou-j` / `cn-hangzhou-k` 是 AI 网关能落的两个 zone。若产品支持的 zone 不同，先查 `data.alicloud_zones` 或参考同产品其它测试。

### 1.2 AI 网关配置模板

```hcl
resource "alicloud_apig_gateway" "default" {
  gateway_name = var.name
  spec         = "aigw.small.x1"
  gateway_type = "AI"
  payment_type = "PayAsYouGo"

  vpc { vpc_id = data.alicloud_vpcs.default.ids.0 }
  network_access_config { type = "Intranet" }
  zone_config { select_option = "Manual" }
  vswitch { vswitch_id = data.alicloud_vswitches.j.ids.0 }
  log_config { sls { enable = false } }

  zones {
    vswitch_id = data.alicloud_vswitches.j.ids.0
    zone_id    = "cn-hangzhou-j"
  }
  zones {
    vswitch_id = data.alicloud_vswitches.k.ids.0
    zone_id    = "cn-hangzhou-k"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

### 1.3 若上游资源 destroy 有已知遗留 bug

在 Aone 单独记一条 `upstream_bug`（资源名 + 现象 + tf-debug.log 定位），**不在本 test 里绕**——workaround 掩盖问题，让下一个受影响的资源重踩。

## 2. Datasource test：不要滥用 `depends_on`

**症状**：datasource 引用的资源 `.id` 已经建立隐式依赖，若额外加 `depends_on = [<resource>]`，Terraform 会在 apply 后的 refresh 阶段把 datasource state 从空刷到非空，被 test framework 判定为"plan not empty after apply"→ Step 0 挂。

**错**：

```hcl
data "alicloud_apig_ai_model_providers" "default" {
  gateway_id = alicloud_apig_ai_model_provider.default.gateway_id
  ids        = [alicloud_apig_ai_model_provider.default.id]
  depends_on = [alicloud_apig_ai_model_provider.default]   # ← 多余，且有害
}
```

**对**：

```hcl
data "alicloud_apig_ai_model_providers" "default" {
  gateway_id = alicloud_apig_ai_model_provider.default.gateway_id
  ids        = [alicloud_apig_ai_model_provider.default.id]
}
```

`.id` 引用足以让 Terraform 先建资源再读 datasource。

## 3. TestingCoverageRate：Optional 字段也必须在某个 Step Config 出现

**症状**：CI `TestingCoverageRate` 报 `service_ids missing test cases`——虽然放了 `ImportStateVerifyIgnore: []string{"service_ids"}`，coverage rate check 只免"modified"部分，**不免**"testing"部分。

**修法**：在**任一** Step 的 Config 里显式给 Optional 字段一个值——空数组也算：

```go
{
  Config: testAccConfig(map[string]interface{}{
    "gateway_id":     "${alicloud_apig_gateway.default.id}",
    "model_provider": "openai",
    "display_name":   name,
    "service_ids":    []string{},   // ← 覆盖 Optional
  }),
  Check: resource.ComposeTestCheckFunc(
    testAccCheck(map[string]string{
      "gateway_id":     CHECKSET,
      "model_provider": "openai",
      "display_name":   name,
      "service_ids.#":  "0",
    }),
  ),
},
```

如果字段值真难构造（e.g. 需要另一个真实资源做绑定目标），空值/占位值也比缺席强。

## 4. `ImportStateVerifyIgnore` 何时用

- Read 无法读回的字段（生成产物 Bug 2 场景，Read 未反向映射的 write-only 字段）→ 放这里免 Import 校验挂
- Create-only 的敏感字段（e.g. `password` / `secret` 服务端不回读）→ 放这里
- **不要**放已经通过 `d.Set` 正确回读的字段——正确回读的字段应通过 Import 校验，掩盖只会积累技术债

## 5. Test 函数命名 + 一处地雷

- `func TestAccAliCloud<Namespace><Resource>_basic(t *testing.T)` — 单资源基础测试
- Datasource：`func TestAccAliCloud<Namespace><Resource>sDataSource(t *testing.T)` — 复数
- 命名规则决定 `invoke-terraform-acc-test-remote` 的自动匹配（`TestAccAliCloud{Namespace}{ResourceTypeCode}*`），拼错就 0 用例
- **Legacy 命名雷**：`TestAccAlicloud`（小写 c）或 `TestAccAliCloudXxx_basic` vs `TestAcc_AliCloudXxx` 都跑不到——统一 `TestAccAliCloud`

## 6. 收尾

- 每个新写的 test 都要有配对的 `AlicloudXxxMap` + `AlicloudXxxBasicDependence` 常量/函数
- Datasource test 复用 resource test 的 `AlicloudXxxBasicDependence`——只加一个 `resource "..." "default"` block + `data "..." "default"` block
- 一切依赖跨账号/跨产品的 env var（`ALICLOUD_ACCESS_KEY_1` 等），走 `testAccPreCheckWithXxx` gate，**别写死 ID**

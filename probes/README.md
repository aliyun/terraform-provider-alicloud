# probes/ — 合成客户场景语料库

`bootstrap/probe.sh` 的输入。每个场景是一段**客户视角**的 Terraform 配置：像真实客户一样，参考
aliyun/alicloud 官方 website/docs 示例，把资源拼成一套真实用法，用来主动探测
terraform-provider-alicloud 潜在且未暴露的问题。

> 能力全景、危害分级→优先级映射、tier 分层与路线图见 `escalation/cap-tf-customer-probe.md`。
> 探测循环 runbook 见 `loops/tf-probe.md`。技能与判定规则见 `.claude/skills/tf-customer-probe/`。

## 分层(2026-07-03 重定义)

- **tier-0 = 静态三方一致性扫描**（TF 文档 ↔ OpenAPI 文档 ↔ provider 源码），**以资源为单位**跑
  `probe.sh tier0`，不跑 terraform。机械部分只做本地 文档↔源码 diff;OpenAPI 一侧留 `judgment_queue`
  交 skill 层查证。**范围红线:只测已接入 TF 的面,未接入的资源/参数一律不报 gap**（那是需求不是 bug）。
- **tier-1 = 真实 apply 全生命周期探测**（默认开启），**以场景为单位**跑 `probe.sh run <id>`。
  本目录下的场景全部是 tier-1 场景，故 `scenario.yaml` **不再有 `tier` 键**。

## 目录结构

```
probes/
  README.md                     ← 本文件
  scenarios/
    <id>/
      scenario.yaml             ← 扁平元数据(probe.sh 用 grep/sed 解析)
      main.tf                   ← 场景 HCL(以官方文档示例为骨架)
      checks.md                 ← 期望行为清单 + 文档链接 + 探测点
      step2/                    ← 可选:update 场景的覆盖层(完整 tf,覆盖 main.tf)
        main.tf
```

### scenario.yaml 键

| 键 | 说明 |
|----|------|
| `id` | 场景 id,**必须与目录名一致且全局唯一** |
| `title` | 一句话标题 |
| `persona` | `beginner` / `composer` / `updater` / `importer`(见 skill scenario-authoring.md) |
| `products` | 涉及云产品(逗号分隔) |
| `resources` | 涉及的 **managed** 资源类型(逗号分隔;data 源不列)。也是 `probe.sh tier0` 无参时的扫描并集来源 |
| `cost` | 成本提示(`free` 等);**不再是硬门**——门是 prepaid 守门(见下) |
| `detect` | 本场景意在触发的 finding code(逗号分隔) |
| `update_step` | `true` 则必须有 `step2/` |
| `import_check` | `true` 则必须成对声明 `import_address` + `import_id_output` |
| `import_address` | import 目标地址(如 `alicloud_vpc.main`) |
| `import_id_output` | 提供 import id 的 output 名 |
| `region`（可选） | 覆盖本场景 region;不写则走 config `regions.focus`（当前 5 场景都不写） |
| `allow_prepaid`（可选） | `true` 显式豁免 prepaid 守门（本轮场景都不需要,留空即默认守门） |
| `source_docs` | 场景骨架来源文档链接 |

> **无 `tier` 键**：分层重定义后场景天然都是 tier-1;tier-0 以资源为单位,不以场景为单位。

## 工单回灌规则（regression-<aone-id>）—— P1 落地，先立规则

每一个被 aone-triage 处理过的**真实客户 TF 问题**，修复合入后应回灌为一个 `regression-<aone-id>` 场景：

- 目录名 `scenarios/regression-<aone-id>/`；`scenario.yaml` 的 `source_docs` 换成 Aone 工单链接；
- `main.tf` 为该工单的最小复现配置；`checks.md` 记录“修复前症状 / 修复后期望”；
- 用途:provider 发新版后复跑该场景 = 回归验证,确认老问题没复发（five-layer 防线第③层）。

## 纪律（硬规则）

1. **绝不在 `probes/` 目录里直接跑 terraform**。runner 会把 `*.tf` 拷到隔离工作目录
   `.my-day/probe/<ts>-<id>/`（已 gitignore）再跑；state / plan / `.terraform/` 全部落在那里，永不入库。
2. **场景必须 pin provider 版本**：每个 `main.tf` 头部固定
   `alicloud = { source = "aliyun/alicloud", version = "1.284.0" }`，与 `config/probe.json` 一致。
3. **命名/标签带 run_id**：可命名资源名称含 `${var.run_id}`；支持 tags 的资源统一加
   `managed_by = "jarvis-probe"`（外加 `run_id` 标签），便于孤儿排查与 `probe.sh sweep` / aliyun CLI 按标签清理。
4. **prepaid 守门(销毁性,不是成本)**：本机是测试账号,付费不设限;但**包年包月/订阅(PrePaid/Subscription)资源
   多数无法 API 销毁**,会破坏「零残留」纪律。runner 在 apply 前扫 plan 的 `*charge_type`/`*payment_type`
   字段,命中 PrePaid/Subscription 默认阻断(`prepaid_block`)。确需探测的场景显式声明 `allow_prepaid: true`。
5. **region**：provider block 不写显式 region(沿用 `ALICLOUD_REGION`);运行时 region 由 runner 解析后注入。
   默认 `config.regions.focus`（eu-central-1,重点探测方向),`--region` 可切,`regions.matrix` 是未来矩阵候选。
6. **以文档为骨架**：HCL 以 `website/docs` 官方示例为骨架做客户式组合,宁可保守照抄,不要凭记忆编字段。

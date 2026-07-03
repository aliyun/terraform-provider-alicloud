# probes/ — 合成客户场景语料库

`bootstrap/probe.sh` 的输入。每个场景是一段**客户视角**的 Terraform 配置：像真实客户一样，参考
aliyun/alicloud 官方 website/docs 示例，把资源拼成一套真实用法，用来主动探测
terraform-provider-alicloud 潜在且未暴露的问题。

> 能力全景、危害分级→优先级映射、tier 分层与路线图见 `escalation/cap-tf-customer-probe.md`。
> 探测循环 runbook 见 `loops/tf-probe.md`。技能与判定规则见 `.claude/skills/tf-customer-probe/`。

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
| `tier` | 场景自然 tier(0=仅静态校验;1=需真实 apply 免费资源) |
| `products` | 涉及云产品(逗号分隔) |
| `resources` | 涉及的 **managed** 资源类型(逗号分隔;data 源不列)。tier-1 场景须 ⊆ `config/probe.json` 的 `tier1_allowlist` |
| `cost` | `free`(P0 只允许 free) |
| `detect` | 本场景意在触发的 finding code(逗号分隔) |
| `update_step` | `true` 则必须有 `step2/` |
| `import_check` | `true` 则必须成对声明 `import_address` + `import_id_output` |
| `import_address` | import 目标地址(如 `alicloud_vpc.main`) |
| `import_id_output` | 提供 import id 的 output 名 |
| `source_docs` | 场景骨架来源文档链接 |

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
4. **只用免费资源**：P0 仅允许 `cost: free`；付费资源（tier-2）绝不进语料库。
5. **以文档为骨架**：HCL 以 `website/docs` 官方示例为骨架做客户式组合，宁可保守照抄，不要凭记忆编字段。

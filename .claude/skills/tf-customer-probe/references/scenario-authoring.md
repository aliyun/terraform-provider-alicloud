# scenario-authoring —— 怎么写一个 probe 场景

写新场景时按此规范。核心信念：**像真实客户一样参考文档**——宁可保守照抄官方示例，不要凭记忆编字段。

## 场景来源优先级

1. **provider 仓 `website/docs` 示例**（最高）：`bootstrap/workspace.sh dir terraform_provider` 拿本地仓，
   逐资源读 `website/docs/r/<name>.html.markdown` 的 Example Usage + Argument Reference，确认目标版本下
   资源名存在、必填参数、字段是否废弃/改名、块是内联还是拆分资源。
2. **阿里云官方 TF 教程 / 最佳实践**：registry 上的 how-to、云产品 IaC 文档。
3. **terraform-alicloud registry modules**（alibaba/* 官方 module）的组合方式。
4. **真实工单回灌**：aone-triage 处理过的客户问题 → `regression-<aone-id>` 场景（见 probes/README.md）。

## persona 定义

| persona | 客户画像 | 探测侧重 |
|---------|----------|----------|
| `beginner` | 新手，照抄单资源/最小组合示例 | 文档示例能否直接跑通（validate/plan/apply） |
| `composer` | 进阶，把多资源/多内联块拼成真实用法 | 组合下的 schema 冲突、JSON 归一化、永久 diff |
| `updater` | 已上线后改配置重跑（step2 覆盖层） | 更新是否生效、改/删字段后是否永久 diff、误重建 |
| `importer` | 用 `terraform import` 纳管存量资源 | import 是否还原完整 state（import 断链） |

## 免费资源判定与清单

- **判定**：资源本身不产生按量/包年计费（创建即计费的实例类一律排除）。空壳/元数据/策略类通常免费。
- **P0 已用免费清单**（也是 `config.tier1_allowlist` 起点）：
  `alicloud_vpc`、`alicloud_vswitch`、`alicloud_security_group`、`alicloud_security_group_rule`、
  `alicloud_ram_role`、`alicloud_ram_policy`、`alicloud_ram_role_policy_attachment`、`alicloud_oss_bucket`（空桶）。
- **存疑即排除**：不确定是否免费的资源，先留在 tier-0（只 validate/plan，不 apply），或干脆不入库。
- 新增 tier-1 资源**必须同步进 `config.tier1_allowlist`**（否则 allowlist 硬门会拦截，探测被降级）。

## 命名 / 标签 / pin 纪律

- 每个 `main.tf` 头部固定：
  ```hcl
  terraform {
    required_version = ">= 1.5.0"
    required_providers {
      alicloud = { source = "aliyun/alicloud", version = "1.284.0" }
    }
  }
  variable "run_id" { type = string }
  ```
- provider block **不写显式 region**（沿用 `ALICLOUD_REGION` env，贴近客户用法）。
- 可命名资源名称含 `${var.run_id}`；支持 tags 的资源统一加 `managed_by = "jarvis-probe"`（外加 `run_id` 标签）。
- OSS 桶名等有格式约束的（小写、全局唯一、长度）直接用 `var.run_id`（runner 传入 `jvp-<id>-<短哈希>`，天然合规）。

## update / import / tier 声明法

- **update 场景**：`scenario.yaml` 设 `update_step: true`，并放 `step2/main.tf`——**完整配置**（含 terraform 块 +
  variable + 资源），runner 会用它覆盖工作目录的 `main.tf` 再 apply。改动尽量原地更新（避免改 ForceNew 字段，
  否则测的是重建不是更新）。
- **import 场景**：`import_check: true` + 成对声明 `import_address`（如 `alicloud_vpc.main`）与 `import_id_output`
  （提供真实 id 的 output 名）。runner 会 `state rm` → `import`（id 取自该 output）→ `plan`。
- **tier 声明**：需真实 apply 才有价值的场景写 `tier: 1`；纯静态（只想验 validate/plan，不建资源）写 `tier: 0`。
  实际执行 tier = min(场景 tier, `--tier` 请求, 配置允许最高 tier)——`tier1_enabled=false` 时一律降级 tier-0。

## 写完自检

- `bootstrap/probe.sh list` 能看到新场景。
- `bootstrap/probe.sh run <id> --dry` 退 0 且步骤计划符合预期。
- `bash test/probe_test.sh` 全绿（校验 yaml 键齐全、id 唯一、tier-1 resources ⊆ allowlist、pin 版本等）。

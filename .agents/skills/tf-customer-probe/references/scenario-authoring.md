# scenario-authoring —— 怎么写一个 probe 场景

写新场景时按此规范。核心信念：**像真实客户一样参考文档**——宁可保守照抄官方示例，不要凭记忆编字段。

> 分层重定义(2026-07-03)后,`probes/scenarios/` 里全是 **tier-1 生命周期场景**(真实 apply),`scenario.yaml`
> **无 `tier` 键**。tier-0 是**资源级**静态扫描(`probe.sh tier0`),不写成场景。

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

## 资源选择:prepaid 守门(不是成本白名单)

- 本机是**测试账号**,付费不设限(成本白名单成本门已撤销)。
- **真正的门是「可销毁性」**:包年包月/订阅(PrePaid/Subscription)资源多数无法 API 销毁,会破坏「零残留」纪律。
  runner apply 前扫 plan 的 `*charge_type`/`*payment_type` 字段,命中 PrePaid/Subscription 默认阻断。
- **写场景时**:优先按量付费(PostPaid)或本身不计费的资源;若资源默认包年包月,显式在 HCL 里设按量付费字段
  (如 `instance_charge_type = "PostPaid"` / `payment_type = "PayAsYouGo"`)。
- **确需探测 prepaid 生命周期**:`scenario.yaml` 声明 `allow_prepaid: true` 显式豁免(慎用,须能确认可销毁)。

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
- provider block **不写显式 region**（沿用 `ALICLOUD_REGION` env，贴近客户用法;runner 注入解析后的 region）。
- 可命名资源名称含 `${var.run_id}`；支持 tags 的资源统一加 `managed_by = "jarvis-probe"`（外加 `run_id` 标签）。
- OSS 桶名等有格式约束的（小写、全局唯一、长度）直接用 `var.run_id`（runner 传入 `jvp-<id>-<短哈希>`，天然合规）。

## region 声明

- 默认走 config `regions.focus`（eu-central-1,重点探测方向)——**多数场景不写 `region:`**。
- 只有当场景**必须**在某 region(资源可用性/配额)时才在 `scenario.yaml` 写 `region: <r>`。
- 命令行 `probe.sh run <id> --region <r>` 临时覆盖(优先级最高)。`regions.matrix` 是未来矩阵探测的候选集。

## update / import 声明法

- **update 场景**：`scenario.yaml` 设 `update_step: true`，并放 `step2/main.tf`——**完整配置**（含 terraform 块 +
  variable + 资源），runner 会用它覆盖工作目录的 `main.tf` 再 apply。改动尽量原地更新（避免改 ForceNew 字段，
  否则测的是重建不是更新）。
- **import 场景**：`import_check: true` + 成对声明 `import_address`（如 `alicloud_vpc.main`）与 `import_id_output`
  （提供真实 id 的 output 名）。runner 会 `state rm` → `import`（id 取自该 output）→ `plan`。

## tier-0 扫描范围红线(资源级,不写成场景)

- tier-0 是 `probe.sh tier0 [alicloud_xxx ...]`(无参 = 场景 resources 并集),做文档↔源码机械 diff。
- **只测已接入 TF 的面**:某云产品资源/参数**没接入 provider**,不在 tier-0 检查范围——那是「需求」不是「bug」,
  真要提需求走 tf_customer 需求路径,不当 gap 报。

## 写完自检

- `bootstrap/probe.sh list` 能看到新场景。
- `bootstrap/probe.sh run <id> --dry` 退 0 且步骤计划符合预期(region 解析正确、prepaid 守门说明合理)。
- `bash test/probe_test.sh` 全绿（校验 yaml 键齐全、id 唯一、无 tier 键、pin 版本、prepaid 守门、region 优先级等）。

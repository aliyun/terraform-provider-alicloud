---
name: tf-customer-probe
description: Use when Jarvis should PROACTIVELY hunt for latent, not-yet-reported bugs in terraform-provider-alicloud by acting like a real customer — reference the official aliyun/alicloud website/docs, write HCL for different personas, run the terraform lifecycle on free resources, and turn failures into severity-ranked Aone ticket drafts. Triggers 主动探测 / 合成客户 / synthetic customer / tf-probe / probe provider / 防患于未然 / 跑一轮场景探测 / 发现潜在问题 / provider 体检. NOT for 已有工单处理(用 aone-triage)、GitHub PR 评审(用 terraform-pr-review)、跑远程 AccTest(用 invoke-terraform-acc-test-remote)、从零开发某个资源(用 provider-resource-dev)。
---

# tf-customer-probe —— 合成客户探测

像真实客户一样参考官方文档、用 terraform 跑真实生命周期，主动发现 provider 潜在问题；按危害分级
产出 Aone 工单草稿，接回 aone-triage loop 形成「自己发现→自己建单→自己修→场景回归复跑」闭环。

- 能力全景/路线图:`escalation/cap-tf-customer-probe.md`
- 循环 runbook:`loops/tf-probe.md`
- 场景语料库:`probes/`（`probes/README.md`）
- runner:`bootstrap/probe.sh`；配置:`config/probe.json`

## 红线（先读）

- **tier-1 双门**：真实 apply 只在 `config.tiers.tier1_enabled=true` **且** plan 资源 ⊆ `tier1_allowlist` 时才放行；两门缺一即降级 tier-0。
- **绝不 tier-2**：付费资源永不探测。
- **凭证零泄漏**：AK/SK 绝不落日志 / verdict / draft / 工单；`doctor` 只报 set/unset。
- **probe 会话不 claim 工单**：本能力只**产出** draft/新单；这些单由 aone-triage loop 后续认领修复（避免既当发现方又当认领方）。
- **只 destroy 本 run 自己建的资源**：state 隔离在 `.my-day/probe/<ts>-<id>/`，绝不碰生产存量。

## 流程

### 1. 预检
```bash
bootstrap/probe.sh doctor
```
- 缺 terraform → env 问题，走 self-improve/escalation，不硬闯（P0 无 terraform 时只能 --dry）。
- 缺凭证 → 只能 tier-0（init/validate/plan 中 plan 亦需凭证；无凭证 plan 记 `no_creds` env_issue）。

### 2. 选场景
```bash
bootstrap/probe.sh list
```
- 默认**最久未跑优先**：看 `runs/probe/` 里各场景最新日期，挑最旧的。
- 单轮 ≤ `config.limits.max_scenarios_per_run`（默认 3）。
- provider 发新版后应尽快全量跑一轮（版本升级易引入 state 不兼容 / 新永久 diff）。

### 3. 逐场景执行
```bash
bootstrap/probe.sh run <id>            # 自然 tier,被配置封顶
bootstrap/probe.sh run <id> --dry      # 只看步骤计划,不需 terraform
```
读产出 `verdict.json`（工作目录 + `runs/probe/<日期>-<id>.json`）。退出码：0 无 findings / 1 有 findings / 2 env 阻断 / 3 清理失败（最高优先级人工介入）。

### 4. 逐 finding 判定（Claude 的判断职责，不是脚本）
对 `verdict.json.findings` 每一条：
1. **复核是否真 provider 问题**：对照 `evidence` 日志、`references/severity-rubric.md`。
2. **查 provider 仓 CHANGELOG Unreleased 段**：若已在 master 修掉，标注「已修复未发布」而**不建单**。
3. `env_issues` 一律不建单（凭证/网络/allowlist/降级）——它们是环境噪声，不是 provider bug。

### 5. 去重
- a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词；
- GitHub `aliyun/terraform-provider-alicloud` open issues 只读检索（同症状是否已被上游报告）。
- 命中重复 → 在已有 draft/单上**追加 evidence**，不新建。

### 6. 产出工单（P0=draft）
- `config.ticket.mode=draft`（当前）→ 按 `references/ticket-template.md` 写
  `escalation/probe-drafts/<日期>-<场景>-<code>.md`，文件头 `status: pending-review`。**不写 Aone**。
- `mode=file`（未来毕业后）→ 走 adhoc-intake 建单纪律：category `req`、project/assignee/tag/priority 按 config，
  受 `config.limits.daily_new_tickets` 上限（默认 3）。毕业条件见 cap 路线图。

### 7. 清理核查
```bash
bootstrap/probe.sh sweep
```
扫 `.my-day/probe/*/terraform.tfstate` 残留；**有残留立即停并升级**（可能有计费资源没删干净）。

### 8. 审计与汇报
- `runs/probe/` 已由 runner 落盘（json + 人读 md）。
- 会话汇报模板：`场景数 N / 发现数 F / draft 数 D / env 问题数 E`，逐 draft 附路径与建议优先级。

## references

| 文件 | 内容 |
|------|------|
| `references/severity-rubric.md` | S1–S4 危害分级 + detector 默认 severity_hint + 升降级判则 + Aone 优先级映射 |
| `references/ticket-template.md` | draft/工单标题与正文骨架、建单参数块、禁 AK/SK 与 AI 署名硬规则 |
| `references/scenario-authoring.md` | 场景来源优先级、persona 定义、免费资源判定、命名/pin 纪律、update/import/tier 声明法 |

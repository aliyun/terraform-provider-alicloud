---
name: tf-customer-probe
description: Use when Jarvis should PROACTIVELY hunt for latent, not-yet-reported bugs in terraform-provider-alicloud by acting like a real customer — statically cross-check TF docs vs OpenAPI docs vs provider source (tier-0), and really apply real resources across the terraform lifecycle (tier-1), turning failures into severity-ranked Aone ticket drafts. Triggers 主动探测 / 合成客户 / synthetic customer / tf-probe / probe provider / 防患于未然 / 跑一轮场景探测 / 三方一致性扫描 / 发现潜在问题 / provider 体检. NOT for 已有工单处理(用 aone-triage)、GitHub PR 评审(用 terraform-pr-review)、跑远程 AccTest(用 invoke-terraform-acc-test-remote)、从零开发某个资源(用 provider-resource-dev)。
---

# tf-customer-probe —— 合成客户探测

像真实客户一样,既**静态**核对文档与实现的一致性,又**真实 apply** 跑生命周期,主动发现 provider
潜在问题;按危害分级产出 Aone 工单草稿,接回 aone-triage loop 形成「自己发现→自己建单→自己修→复跑回归」闭环。

- 能力全景/路线图:`escalation/cap-tf-customer-probe.md`
- 循环 runbook:`loops/tf-probe.md`
- 场景语料库(**外置 jarvis 仓外**):`terraform_playground/`,按云产品维度两级归档 `<product>/<id>/`
  (scenario.yaml + main.tf + checks.md,可选 step2/)。默认路径 `<jarvis 根目录的父目录>/terraform_playground`,
  可用 `JARVIS_TF_PLAYGROUND`(env,最高优先)或 `config/probe.json` 的 `paths.playground_dir` 覆盖;
  README 与场景规范见 playground 仓 README + `references/scenario-authoring.md`。
- runner:`bootstrap/probe.sh`;配置:`config/probe.json`

## 分层(2026-07-03 重定义)

- **tier-0 = 静态三方一致性扫描**(TF 文档 ↔ OpenAPI 文档 ↔ provider 源码),**以资源为单位**,不跑 terraform。
  机械部分只做本地 文档↔源码 diff(确定性强);OpenAPI 一侧留 `judgment_queue` 交 skill 层查证。
- **tier-1 = 真实 apply 全生命周期探测**(默认开启),**以场景为单位**,region 默认 focus=eu-central-1(重点方向,可切)。

## 红线（先读）

- **tier-1 默认真实 apply**:`config.tiers.tier1.enabled=true`。关掉则**封顶 plan-only**(init/validate/plan,不 apply)。
- **prepaid 销毁性守门**:apply 前扫 plan,命中 PrePaid/Subscription 计费类型默认阻断(`prepaid_block`)——原因不是钱
  (测试账号),是包年包月资源多数无法 API 销毁,破坏「零残留」。场景 `allow_prepaid: true` 或 config `prepaid_guard=false` 可豁免。
- **零残留 + 只 destroy 本 run 自建资源**:state 隔离在 `.my-day/probe/<ts>-<id>/`,绝不碰生产存量;`sweep` 残留即停并升级。
- **测试账号边界**:只用环境注入的**测试** AK/SK,绝不用生产账号;凭证值绝不落日志/verdict/draft/工单。
- **tier-0 范围红线**:只核对**已接入 TF** 的面。云产品未接入 provider 的资源/参数**不报 gap**(那是需求不是 bug,走 tf_customer 需求路径)。
- **probe 会话不 claim 工单**:本能力只**产出** draft/新单;由 aone-triage loop 后续认领修复(避免既当发现方又当认领方)。

## 流程

### 0. 预检
```bash
bootstrap/probe.sh doctor
```
terraform/jq/凭证/config/**本地 provider 仓**(tier-0 依赖)。缺 terraform → tier-1 不可跑;缺 provider 仓 → tier-0 不可跑。env 问题走 self-improve/escalation,不硬闯。

### A. tier-0 静态三方一致性扫描(建议先跑)
```bash
bootstrap/probe.sh tier0                 # 无参 = 全部场景 resources 并集
bootstrap/probe.sh tier0 alicloud_vpc    # 指定资源
```
1. runner 机械产出 `findings`(五类 doc↔source gap:`doc_gap_phantom/undocumented/flag_mismatch/forcenew/deprecated`)+ 每资源一条 `judgment_queue`。
2. **OpenAPI 侧判定接棒(Claude/verifier 子代理)**:对 `judgment_queue` 每资源,走 aone-triage/verifier 同款双层查证惯例——
   对照 `judgment_queue.api_actions` 指向的 OpenAPI 文档(next.api.aliyun.com),核验 provider 实际调用的 API 的请求/响应参数、
   枚举值、行为是否与 TF 文档一致;不一致产 `doc_api_gap`(severity 按影响面 S2–S4 由判定给)。
   **范围红线**:只核对已接入的 API/参数面,未接入 TF 的一律不报。
3. findings 定级(见 `references/severity-rubric.md`)→ 去重 → draft。

### B. tier-1 真实 apply 生命周期探测
```bash
bootstrap/probe.sh list                       # 选场景(最久未跑优先,看 runs/probe/ 日期);单轮 ≤ limits.max_scenarios_per_run
bootstrap/probe.sh run <id>                    # region 默认 config.regions.focus(eu-central-1)
bootstrap/probe.sh run <id> --region cn-hangzhou   # 切 region(重点方向 eu-central-1,但非唯一)
bootstrap/probe.sh run <id> --dry              # 只看步骤计划(region 解析、prepaid 守门),不需 terraform
```
读产出 `verdict.json`(工作目录 + `runs/probe/<日期>-<id>.json`)。退出码:0 无 findings / 1 有 findings / 2 env 阻断(含 `prepaid_block`/`tier1_disabled_plan_only`) / 3 清理失败(最高优先级人工介入)。

### C. 逐 finding 判定(Claude 的判断职责,不是脚本)
1. **复核是否真 provider 问题**:对照 `evidence`、`references/severity-rubric.md`;tier-0 findings 还要看 doc/source 上下文是否确为不一致。
2. **查 provider 仓 CHANGELOG Unreleased 段**:已在 master 修掉的标「已修复未发布」**不建单**。
3. `env_issues` 一律不建单(凭证/网络/prepaid/plan-only 都是环境噪声)。

### D. 去重 → 产出工单(当前 mode=file,2026-07-05 已毕业)
- 去重:a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词;GitHub `aliyun/terraform-provider-alicloud` open issues 只读检索;重复则**追加 evidence**不新建。
- `config.ticket.mode=file`(**当前默认**,2026-07-05 主人拍板毕业:采纳率 7/8=87.5% 提前毕业)→ 走 adhoc-intake 建单纪律(category `req`、project/assignee/tag 按 config),受 `daily_new_tickets` 上限。**直接建 Aone 单,不再走 draft 人审。**
- `config.ticket.mode=draft`(**保留为可回退开关**)→ 按 `references/ticket-template.md` 写 `escalation/probe-drafts/<日期>-<资源或场景>-<code>.md`,头加 `status: pending-review`,**不写 Aone**;仅在需要临时收回自动建单权时切回。

### E. 清理核查 + 审计汇报
```bash
bootstrap/probe.sh sweep     # 有残留立即停并升级
```
`runs/probe/` 已由 runner 落盘。会话汇报:`tier0 扫描资源数 / findings / judgment_queue 数`,`tier1 场景数 / findings / draft 数 / env 问题数`,逐 draft 附路径与建议优先级。

## references

| 文件 | 内容 |
|------|------|
| `references/severity-rubric.md` | tier-0/tier-1 finding 码默认 severity_hint + S1–S4 危害分级 + 升降级判则 + Aone 优先级映射 |
| `references/ticket-template.md` | draft/工单标题与正文骨架、建单参数块、禁 AK/SK 与 AI 署名硬规则 |
| `references/scenario-authoring.md` | 场景来源优先级、persona 定义、prepaid 守门、region 声明、update/import 声明、tier-0 范围红线 |

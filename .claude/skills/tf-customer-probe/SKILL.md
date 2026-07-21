---
name: tf-customer-probe
description: Use when Jarvis should PROACTIVELY hunt for latent, not-yet-reported bugs in terraform-provider-alicloud by acting like a real customer — statically cross-check TF docs vs OpenAPI docs vs provider source (tier-0), and really apply real resources across the terraform lifecycle (tier-1), turning failures into severity-ranked Aone tickets. Triggers 主动探测 / 合成客户 / synthetic customer / tf-probe / probe provider / 防患于未然 / 跑一轮场景探测 / 三方一致性扫描 / 发现潜在问题 / provider 体检. NOT for 已有工单处理(用 aone-triage)、GitHub PR 评审(用 terraform-pr-review)、跑远程 AccTest(用 invoke-terraform-acc-test-remote)、从零开发某个资源(用 provider-resource-dev)。
---

# tf-customer-probe —— 合成客户探测

像真实客户一样,既**静态**核对文档与实现的一致性,又**真实 apply** 跑生命周期,主动发现 provider
潜在问题;按危害分级产出 Aone 工单,接回 aone-triage loop 形成「自己发现→自己建单→自己修→复跑回归」闭环。

- 循环 runbook:`loops/tf-probe.md`
- 场景语料库=**独立 git 数据仓** `tf_playground`(gitlab `terraflow/tf_playground`,**直推 master + 工单报备**,非代码不走 MR),
  按云产品维度两级归档 `<product>/<id>/`(scenario.yaml + main.tf + checks.md,可选 step2/)。根解析优先级:
  env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` > `bootstrap/workspace.sh dir tf_playground`
  (数据仓已登记 `config/workspaces.json`,多机 clone 后零配置)> 默认 `<jarvis 根目录的父目录>/terraform_playground`;
  场景规范见 `references/scenario-authoring.md`。
- runner:`bootstrap/probe.sh`;配置:`config/probe.json`

## 分层

- **tier-0 = 静态三方一致性扫描**(TF 文档 ↔ OpenAPI 文档 ↔ provider 源码),**以资源为单位**,不跑 terraform。
  机械部分做本地 文档↔源码 diff(五类 `doc_gap_*`)**+ OpenAPI 侧机械三方 diff**(T0-mech:`probe-meta.sh` 拉 amp
  元数据,六类 `api_gap_*` 预筛);机械层**拿不准的项才留 `judgment_queue`**(映射不上/纯 prose 约束/被抑制存疑/OSS 无 action),
  terraform-pd 只判疑点。**精度命门:拿不准一律 queue,绝不硬报。**
- **tier-1 = 真实 apply 全生命周期探测**(默认开启),**以场景为单位**,region 默认 focus=eu-central-1(重点方向,可切)。

## 红线（先读）

- **tier-1 默认真实 apply**:`config.tiers.tier1.enabled=true`。关掉则**封顶 plan-only**(init/validate/plan,不 apply)。
- **prepaid 销毁性守门**:apply 前扫 plan,命中 PrePaid/Subscription 计费类型默认阻断(`prepaid_block`)——原因不是钱
  (测试账号),是包年包月资源多数无法 API 销毁,破坏「零残留」。场景 `allow_prepaid: true` 或 config `prepaid_guard=false` 可豁免。
- **零残留 + 只 destroy 本 run 自建资源**:state 隔离在 `.my-day/probe/<ts>-<id>/`,绝不碰生产存量;`sweep` 残留即停并升级。
- **测试账号边界**:只用环境注入的**测试** AK/SK,绝不用生产账号;凭证值绝不落日志/verdict/工单。
- **tier-0 范围红线**:只核对**已接入 TF** 的面。云产品未接入 provider 的资源/参数**不报 gap**(那是需求不是 bug,走 tf_customer 需求路径)。
- **probe 会话不 claim 工单**:本能力只**产出**新单;由 aone-triage loop 后续认领修复(避免既当发现方又当认领方)。

## 流程

### 0. 预检
```bash
bootstrap/probe.sh doctor
```
terraform/jq/凭证/config/**本地 provider 仓**(tier-0 依赖)/**probe-meta**(T0-mech OpenAPI 元数据获取层)。缺 terraform → tier-1 不可跑;缺 provider 仓 → tier-0 不可跑;缺 probe-meta(无 venv/凭证)→ tier-0 **自动降级**为纯 doc↔source + 全 queue(不阻断)。env 问题走 self-improve；需人工介入时建 Aone 或挂起当前 Task，不硬闯。

> **T0-mech 元数据层**:`bootstrap/probe-meta.sh`(fetch/cached-fetch/clear/available)薄封装 `amp-resource-metadata` skill 的
> `get_api_definition.py`,带 `cache.sh` 7d 缓存(全量巡检才跑得起)。启用需 `amp-resource-metadata/scripts/setup.sh` 建 venv +
> 配 `AMP_`/`ALIBABA_CLOUD_` 凭证(APISpecData 白名单见该 skill SKILL.md)。probe.sh 不直接调 python。

### A. tier-0 静态三方一致性扫描(建议先跑)
```bash
bootstrap/probe.sh tier0                        # 无参 = 全部场景 resources 并集
bootstrap/probe.sh tier0 alicloud_vpc           # 指定资源
bootstrap/probe.sh tier0 --all --rotate 20      # website/docs/r 全量,取 20 个最久未扫轮换巡检
bootstrap/probe.sh tier0 --no-mech              # 关机械层(纯 doc↔source + 全 queue,api_gap 检测关闭)
```
1. **机械三方 diff 先行**(runner 全机械,确定性强):
   - 本地 文档↔源码:五类 `doc_gap_*`(phantom/undocumented/flag_mismatch/forcenew/deprecated)。
   - OpenAPI 侧(T0-mech,`probe-meta.sh available` 时自动开;不可用则降级为纯 doc↔source + 全 queue):从源码抽
     (product,version,action) 三元组 + 解析 `StringInSlice`/`IntBetween`/`Default`,对 amp 元数据机械 diff → 六类
     `api_gap_*`(deprecated_action/enum_superset/required/type/range/default)。被抑制项入 `suppressed[]`、TF 更严项入 `coverage_notes[]`。
2. **terraform-pd 只判疑点**(收窄后的 `judgment_queue`,每条带 `reason`):`prose_review`(长度/字符集/基数等纯 prose 约束、
   行为一致性)/`unmapped_params`(snake→Camel 映射不上,如 convert 改名)/`enum_unparsed`(枚举非字面 slice)/
   `no_triple`(OSS SDK 风格抽不到 action)/`meta_unavailable`。对这些走 aone-triage/terraform-pd 双层查证,不一致产 `doc_api_gap`。
   **范围红线**:只核对已接入的 API/参数面,未接入 TF 的一律不报。
3. findings 定级(见 `references/severity-rubric.md`,含 `api_gap_*` 表)→ 去重 → 建单。

### B. tier-1 真实 apply 生命周期探测
```bash
bootstrap/probe.sh list                       # 选场景(最久未跑优先,LAST_RUN 列由 t1-last-run.json 索引);单轮 ≤ limits.max_scenarios_per_run
bootstrap/probe.sh run <id>                    # region 默认 config.regions.focus(eu-central-1)
bootstrap/probe.sh run <id> --region cn-hangzhou   # 切 region(重点方向 eu-central-1,但非唯一)
bootstrap/probe.sh run <id> --dry              # 只看步骤计划(region 解析、prepaid 守门),不需 terraform
```
读产出 `verdict.json`(工作目录 + `runs/probe/<YYYYMMDD>-<HHMMSS>-<id>.json`)。退出码:0 无 findings / 1 有 findings / 2 env 阻断(含 `prepaid_block`/`tier1_disabled_plan_only`) / 3 清理失败(最高优先级人工介入)。

`scenario.yaml` 可选新键(全部由 runner 支撑,persona 一行带过,详细语义/护栏见 `references/scenario-authoring.md`):
- `steps: step2,step3`(CSV)+ `step<N>_expect: changed|no_changes|fail`——泛化 update_step,逐步声明期望;
- `expect_fail: validate|plan|apply` + `expect_error_contains`——四态判定(expected/expected_but_error_mismatch/late_validation/expected_fail_missed);
- `provider_version_from: <old-pin>`——upgrader 场景 A→B state 兼容;
- `drift_cli: aliyun <product> <Action> ...`——drifter 场景(默认 `drift_enabled=false` 关闭,转正走 config MR)。

### C. 逐 finding 判定(Claude 的判断职责,不是脚本)
1. **复核是否真 provider 问题**:对照 `evidence`、`references/severity-rubric.md`;tier-0 findings 还要看 doc/source 上下文是否确为不一致。
2. **查 provider 仓 CHANGELOG Unreleased 段**:已在 master 修掉的标「已修复未发布」**不建单**。
3. `env_issues` 一律不建单(凭证/网络/prepaid/plan-only 都是环境噪声)。

### D. 去重 → 产出工单
- 去重:a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词;GitHub `aliyun/terraform-provider-alicloud` open issues 只读检索;重复则**追加 evidence**不新建。
- **generated 场景先过校订门**:`origin: generated` 且未经人工校订的场景跑出 `apply_fail`/`validate_fail` **先归「场景质量疑点」入人工校订队列,不直接建 bug 单**(机器抽文档+机械改造的配置本身可能有缺陷);须人工确认确系 provider 行为才升级。手写/回灌场景无此限。
- `config.ticket.mode=file` 是唯一支持模式：走 adhoc-intake 建单纪律(category `req`、project/assignee/tag 按 config),受 `daily_new_tickets` 上限，直接建 Aone 单。低置信 finding 不落本地草稿，进入 generated 场景校订门或创建 needs-attention Aone。

### D.5. 建单后回写 ledger(kind=ticket)
每建成一张 probe 工单,追加一条 `runs/probe/ledger.jsonl` 记录,供 board/`archive` 后续对账与消费:
```bash
jq -nc --arg ts "$(date -u +%FT%TZ)" \
       --arg tid "<新建工单 id>" \
       --arg code "<finding code>" \
       --arg v "<verdict path,如 runs/probe/20260708-193012-<sid>.json>" \
       '{ts:$ts, kind:"ticket", ticket_id:$tid, finding_code:$code, verdict:$v}' \
    >> runs/probe/ledger.jsonl
```
### E. 清理核查 + 归档 + 蒸馏 + 审计汇报
```bash
bootstrap/probe.sh sweep     # 有残留立即停并升级
bootstrap/probe.sh archive   # 幂等归档:verdict retention/workdir gc/plugin-cache 报告/待办清单(--dry 演习)
```
`runs/probe/` 已由 runner 落盘。**收尾必做**：按 `references/knowledge-distillation.md` 契约把本轮学到的
产品级知识蒸馏进 `<playground>/<product>/KNOWLEDGE.md`(触发点①probe 轮 Step E 收尾)。

会话汇报:`tier0 扫描资源数 / findings / judgment_queue 数`,`tier1 场景数 / findings / 工单数 / env 问题数`,
归档件数(verdicts moved / workdirs gc),蒸馏条目数(逐产品);逐工单附链接与建议优先级。

## references

| 文件 | 内容 |
|------|------|
| `references/severity-rubric.md` | tier-0/tier-1 finding 码默认 severity_hint + S1–S4 危害分级 + 升降级判则 + Aone 优先级映射 |
| `references/ticket-template.md` | 工单标题与正文骨架、建单参数块、禁 AK/SK 与 AI 署名硬规则 |
| `references/scenario-authoring.md` | 场景来源优先级、persona 定义(beginner/composer/updater/importer/migrator/upgrader/refactorer/drifter/ds-checker/ci-runner)、新键(steps/expect_fail/upgrader/drifter)、prepaid 守门、region 声明、tier-0 范围红线 |
| `references/knowledge-distillation.md` | 云产品 KNOWLEDGE.md 蒸馏契约(跨 skill 单点):五节结构 / 条目格式 / 触发点(probe/triage/dev)/ 收录判据 / sanitize / 毕业标准 |

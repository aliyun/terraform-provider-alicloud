# tf-probe 合成客户探测循环

> 主动探测 terraform-provider-alicloud 潜在问题的 runbook。与 aone-triage 形成闭环：
> probe 产出的工单被 triage loop 扫到后按正常流程认领修复；修复发布后对应场景/资源复跑即回归验证。
> 能力全景/路线图见 `escalation/cap-tf-customer-probe.md`；全流程技能见 `.claude/skills/tf-customer-probe`。

## 分层

- **tier-0 = 静态三方一致性扫描**(TF 文档 ↔ OpenAPI 文档 ↔ provider 源码),**以资源为单位**,不跑 terraform。
- **tier-1 = 真实 apply 全生命周期探测**(默认开启),**以场景为单位**,region 默认 focus=eu-central-1(重点方向)。
- **执行顺序建议:先 tier-0 扫本轮涉及资源(快、零成本)→ 再 tier-1 跑场景生命周期。**

---

## 一、触发

| 方式 | 说明 |
|------|------|
| 手动触发 | 用户在会话中执行 `/tf-probe` 或直接发指令「跑一轮场景探测 / 三方一致性扫描」 |
| provider 新版 | 新版本发布后应尽快全量跑一轮(版本升级易引入 state 不兼容 / 新永久 diff / 文档漂移) |
| cron / bridge | P2 才接入,本轮**不接** |

---

## 二、预检（doctor）

```bash
bootstrap/probe.sh doctor
```

- 全绿退 0 才继续。
- 缺 **terraform** → tier-1 不可跑;env 问题,**不硬闯**,按 loops/self-improve.md 记缺口 / escalation。
- 缺 **本地 provider 仓**(website/docs + alicloud) → tier-0 静态扫描不可跑(WARN)。
- 缺 **probe-meta**(T0-mech OpenAPI 元数据层:无 venv/凭证)→ tier-0 **自动降级**为纯 doc↔source + 全 queue(WARN,不阻断)。
- 缺凭证 → tier-1 只能 plan-only(plan 亦需凭证;无凭证 plan 记 `no_creds`,不是 finding)。

---

## 三、tier-0 静态扫描(先跑)

```bash
bootstrap/probe.sh tier0                        # 无参 = 全部场景 resources 并集(去重)
bootstrap/probe.sh tier0 alicloud_vpc           # 指定资源
bootstrap/probe.sh tier0 --all --rotate 20      # website/docs/r 全量,取 20 个最久未扫轮换巡检(状态落 .my-day/probe/t0mech-scanned.json)
bootstrap/probe.sh tier0 --no-mech              # 关机械层(纯 doc↔source + 全 queue,= T0-mech 前现行为)
bootstrap/probe.sh tier0 --dry                  # 只列将扫资源 + mech 模式 + 文档/源码存在性
```

- **机械三方 diff 先行**:五类 `doc_gap_*`(doc↔source)**+ 六类 `api_gap_*`**(T0-mech,OpenAPI 元数据机械 diff:
  deprecated_action/enum_superset/required/type/range/default);被抑制记 `suppressed[]`、TF 更严记 `coverage_notes[]`。
  `probe-meta` 不可用 → 自动降级为纯 doc↔source + 全 queue(现行为)。
- **verifier 只判疑点**:收窄后的 `judgment_queue`(每条带 `reason`:prose_review/unmapped_params/enum_unparsed/no_triple/
  meta_unavailable)走双层查证,产 `doc_api_gap`。**范围红线:只核对已接入 TF 的面,未接入的资源/参数不报 gap(需求非 bug)**。
- 落盘 `runs/probe/<YYYYMMDD>-<HHMMSS>-tier0.json`(verdict `mech` 字段标 on/off/degraded) + 人读 md。退出码:0 无 findings / 1 有 findings / 2 runner 错误。

---

## 四、tier-1 场景生命周期探测

### 4.1 场景选择
```bash
bootstrap/probe.sh list
```
- 默认**最久未跑优先**:LAST_RUN 列由 `.my-day/probe/t1-last-run.json` 索引给出(仅推进到 plan 完成及以后才更新),
  挑最旧的;单轮 ≤ `config.limits.max_scenarios_per_run`(默认 3)。索引缺失回退扫 `runs/probe/` 文件名(兼容旧
  `<YYYYMMDD>-<sid>.json` 与新 `<YYYYMMDD>-<HHMMSS>-<sid>.json`)。
- persona 分类(scenario.yaml `persona:` 键,详细定义见 `.claude/skills/tf-customer-probe/references/scenario-authoring.md`):
  `beginner` / `composer` / `updater` / `importer` / `migrator`(配置形态迁移)/`upgrader`(provider 版本 A→B)/
  `refactorer`(moved block/资源改名)/`drifter`(带外改动检测)/`ds-checker`(数据源读回一致性,纯 HCL 零 runner 改动)/
  `ci-runner`(taxonomy 立骨架,本轮不写场景)。负路径断言 = `expect_fail` 键语义(不是 persona)。

### 4.2 执行
```bash
bootstrap/probe.sh run <id>                     # region 默认 config.regions.focus(eu-central-1)
bootstrap/probe.sh run <id> --region cn-hangzhou # 切 region(重点 eu-central-1,但保留多 region)
bootstrap/probe.sh run <id> --dry               # 只看步骤计划(region 解析、prepaid 守门)
```
- region 优先级:`--region` > scenario.yaml `region:` > config `regions.focus` > 环境 `ALICLOUD_REGION`。
- **prepaid 销毁性守门**:apply 前扫 plan,命中 PrePaid/Subscription 计费类型默认阻断(`prepaid_block`);场景 `allow_prepaid: true` 豁免。
- `tiers.tier1.enabled=false` → **封顶 plan-only**(不 apply,env_issue `tier1_disabled_plan_only`)。
- 读产出 `verdict.json`(工作目录 + `runs/probe/<YYYYMMDD>-<HHMMSS>-<id>.json` + md)。退出码:0/1/2/**3 清理失败(最高优先级人工介入)**。
- **findings 判定是 Claude 的职责**:逐 finding 对照 evidence、`references/severity-rubric.md`、CHANGELOG Unreleased 段(已修掉的标「已修复未发布」不建单)。

---

## 四点五、RC 门禁（发版前全量过闸，rc-gate.sh）

> 把 tier-0（全量机械层）+ tier-1（全场景生命周期）串成**一条发版前门禁线**,一次跑完给红/黄/绿判定与退码,
> 供 `terraform-provider-release` SOP 在 PR/合并前当闸门用。**它是编排层,不复制 probe 逻辑——每步都调 probe.sh 子命令。**

### 何时跑

- **provider 发版前(RC 阶段)**:切好发版分支、准备提 PR / 合并前,跑一遍确认全量探测过闸。
- **新版本落地后**:版本升级易引入 state 不兼容 / 新永久 diff / 文档漂移,全量跑一轮兜底。
- 日常增量开发**不必**每次跑全量;那是 tier-0 `--rotate` 巡检 + tier-1 单场景的活。

### 怎么跑

```bash
bootstrap/rc-gate.sh <provider-dir>            # 完整模式:tier-1 真实 apply 顺跑(并发 1)
bootstrap/rc-gate.sh <provider-dir> --quick    # 快扫:tier-1 改 plan 为止(run --dry,零创建),不 apply
```

- `<provider-dir>` = 本地 terraform-provider-alicloud 仓(经 `JARVIS_PROBE_PROVIDER_DIR` 传给 tier0,需 website/docs/r + alicloud)。
- 三步:① `probe.sh tier0 --all --limit 200`(机械层全量,**降级容忍**:mech=degraded 记黄不记红)
  → ② `probe.sh list` 枚举全场景 → 逐场景 `probe.sh run`(完整=真实 apply;`--quick`=`run --dry` plan 为止)
  → ③ 汇总 `runs/rc-gate/<date>-report.md` 并按判定退码。
- 可调 env:`RC_GATE_QUEUE_YELLOW`(queue 激增黄线,默认 40)、`RC_GATE_TOTAL_TIMEOUT`(tier-1 阶段总超时预算秒,默认 0=不限,超预算剩余场景跳过记黄)。

### 怎么读报告 / 判定与退码

报告落 `runs/rc-gate/<date>-report.md`,顶部 `## VERDICT: <色>  (exit N)` 一眼定结论:

| 判定 | 退码 | 触发 | 动作 |
|------|------|------|------|
| 🔴 **RED** | 1（阻断） | tier-0 `api_gap` 严重度 **S3+** / 场景 `run` 退 1（provider finding）/ 场景 `run` 退 3（destroy 失败或 state 残留） | **禁发**,先修红项再复跑 |
| 🟡 **YELLOW** | 0（放行但标注） | 机械层降级 mech=degraded / judgment_queue 激增 / tier-0 有非 S3+ finding（doc_gap、S4 api_gap）/ 场景 `run` 退 2（env 阻断,非 bug）/ `--quick` 未跑 apply / 无场景 / 超时跳过 | 可放行,但报告显著标注,发版前尽量清 |
| 🟢 **GREEN** | 0 | 以上皆无 | 过闸 |
| ⚪ **CANNOT_CERTIFY** | 2（门禁不完整） | tier-0 无法运行（`probe.sh tier0` 退 2,provider 仓不可用） | 修环境后重跑,**勿把不可判当绿放行** |

- 优先级 **RED > CANNOT_CERTIFY > YELLOW > GREEN**;stdout 末行同样打 `VERDICT: <色>` 便于 CI/脚本 grep。
- **降级不误伤**:probe-meta（OpenAPI 元数据层）不可用时 tier-0 自动降级,api_gap 检测关闭→api_gap S3+ 恒 0,只靠 doc↔source 出黄项,门禁**不会因环境降级而误判红**。

### 与 terraform-provider-release SOP 的衔接

```
terraform-provider-release SOP: 需求澄清 → gap 分析 → (生成/改码) → [RC 门禁: rc-gate.sh] → 远程 ACC → PR → 人工合并
                                                                        ↑
                                          绿/黄 → 继续 SOP 的 PR/ACC/合并环节(黄项知情放行)
                                          红   → 回到改码,修红项后复跑门禁
                                          不可判 → 修 provider 仓/probe-meta 环境后重跑
```

release skill 侧仅加了 additive「发版前强化门禁(可选)」指引(见其 `references/rc-gate.md`),不改 SOP 既有步骤;是否插入门禁由跑者按发版风险自定。

---

## 五、产物分流

| 产物 | 去向 |
|------|------|
| `findings`(provider 疑似 bug,tier-0 或 tier-1) | 去重后 → **当前 `mode=file` 直接 adhoc-intake 建单**(tf_provider 528766,见 skill);`mode=draft`(可回退开关)时改落 `escalation/probe-drafts/` |
| `env_issues`(凭证/网络/`prepaid_block`/`tier1_disabled_plan_only`) | **不建单**;缺 terraform / 缺 provider 仓等能力类走 loops/self-improve.md 或 escalation |

- 去重:a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词;GitHub 上游 open issues 只读检索;重复则追加 evidence 不新建。
- **probe 会话本身不 claim 工单**:产出的新单由 aone-triage loop 后续接管。

---

## 五点五、归档与台账（archive + ledger + LRU 索引）

```bash
bootstrap/probe.sh archive          # 幂等归档（真跑）
bootstrap/probe.sh archive --dry    # 演习（不动文件,ledger 追加 kind=archive_dry 记审计）
```

`probe.sh archive` 五件事一次做完(幂等,dry 也追加 ledger 便于对账):

1. **draft 归档**:`escalation/probe-drafts/*.md` frontmatter `status` 为 `filed` / `rejected-*` → 移入
   `escalation/probe-drafts/archived/`；`pending-review`/未知 status 留原地进「待办清单」。
2. **verdict retention**:`runs/probe/` 顶层超 `limits.audit_retention_days`(默认 60)的 verdict json+md
   → `runs/probe/archive/<YYYYMM>/`。**排除** `ledger.jsonl` 与 `*-summary.md`。
3. **工作目录 gc**:`.my-day/probe/<ts>-<sid>/`(runner 实际用 `YYYYMMDDTHHMMSSZ-<sid>` 形态,归档扫码兼容
   `^[0-9]{8,}[A-Za-z0-9]*-.+$`)且 tfstate 不存在或 resources 为空、mtime 超 `limits.workdir_retention_days`
   (默认 7)→ rm -rf。**显式排除** `.plugin-cache/`、`t0mech-scanned.json`、`t1-last-run.json`、`manual-*`。
   tfstate 非空绝不删(sweep 残留即停语义不变)。
4. **plugin-cache 报告**:与 `config .provider.version` 及任何场景声明过的 `provider_version_from` 都不符的
   registry 版本目录 → **只报体积不删**(人工评估)。
5. **待办清单输出**:pending-review drafts、playground `_quarantine/` 存量、`origin: generated` 未校订计数——
   archive 每次跑都打印,作为「校订队列」的可见性载体。

**台账 `runs/probe/ledger.jsonl`**(append-only,本机加速索引,真源在 Aone):
- tier0/tier1 finalize 时 runner 追加 `{ts, kind:"tier0"|"tier1", scenario|resources, findings, verdict}`;
- skill Step D.5 建单后 Claude 追加 `{ts, kind:"ticket", ticket_id, finding_code, verdict}`;
- archive 追加 `{ts, kind:"archive"|"archive_dry", moved:{drafts,verdicts,workdirs}}`。

**LRU 索引 `.my-day/probe/t1-last-run.json`**(scenario→epoch,`probe.sh list` 行尾 LAST_RUN 列):
- 仅在真实执行**推进到 plan 完成及以后**才更新(env 阻断如 no_creds/prepaid_block/tier1_disabled 不更新,
  避免被阻断场景永远失去 LRU 优先权);
- 索引缺失时 `_t1_last_run_get` 会回退扫 `runs/probe/` 文件名(兼容旧 `<YYYYMMDD>-<sid>.json` 与新
  `<YYYYMMDD>-<HHMMSS>-<sid>.json` 正则,sid 含连字符禁 naive cut)。

---

## 六、清理核查（sweep，残留即停）

```bash
bootstrap/probe.sh sweep
```

- 无残留退 0。
- **有残留退 1 → 立即停并升级**:`.my-day/probe/*/terraform.tfstate` 有非空 `resources`,按提示手动 `terraform destroy`
  或按 `managed_by=jarvis-probe` 标签用 aliyun CLI 清理(prepaid 守门已尽量把不可销毁资源挡在 apply 前)。

---

## 七、Done — 本轮结束标准

- tier-0:本轮涉及资源都扫过,`runs/probe/<YYYYMMDD>-<HHMMSS>-tier0.json` 落盘,judgment_queue 已交判定。
- tier-1:每个跑过的场景都有 `verdict.json`;`sweep` 零残留。
- 每条 provider finding 都已去重并落 draft/工单(或标注「已修复未发布」/「上游已报」跳过);建单后 D.5 回写
  `runs/probe/ledger.jsonl`(`kind:"ticket"`)。
- **`probe.sh archive` 已跑过一次**——draft 终态归档 + verdict retention + workdir gc + 待办清单打印;ledger 追加。
- **知识蒸馏完成**——按 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md` 契约,本轮涉及
  的每个产品都过一遍 KNOWLEDGE.md 追加(命中收录判据的条目;无可收条目则显式记「无产品级新增知识」)。
- 会话汇报:`tier0 资源数/findings/judgment 数` + `tier1 场景数/findings/draft 数/env 数` + 归档件数
  (drafts/verdicts/workdirs) + 蒸馏条目数(逐产品),逐 draft/工单附路径与建议优先级。

---

## 八、与 aone-triage 的闭环

```
tf-probe 发现问题(tier-0 doc gap / tier-1 生命周期 bug) → draft/建单(tf_provider 528766, tag jarvis-probe)
        ↓
aone-triage loop 扫到该单 → claim → provider-resource-dev/修复 → PR → 发布
        ↓
发布后:tf-probe 复跑对应资源(tier0)/场景(tier1) → 无 finding = 回归通过(闭环)
```

真实客户工单也应回灌为 `regression-<aone-id>` 场景(直落 git 数据仓 `tf_playground/<product>/regression-<aone-id>/`
+ 直推 master + 工单评论报备路径,无需 worktree/MR;规范见 `.claude/skills/tf-customer-probe/references/scenario-authoring.md`),成为发版前回归项。

### 修复侧衔接

probe 单指派 jarvis,会被 `scan.sh` 自然扫到进 triage;aone-triage 按 **四分类路由**(provider 代码修 / TF 文档修 /
上游协作 / 需实验定性)修复,PR 合并后按溯源**复验关单**,并回灌 regression 场景——完整飞轮见
`escalation/cap-probe-fix-flywheel.md`,路由与复验状态机见 `.claude/skills/aone-triage/references/probe-ticket-routing.md`。

---

## 九、工具链速查

| 工具 | 作用 |
|------|------|
| `bootstrap/probe.sh doctor` | 环境预检(terraform/jq/凭证/config/本地 provider 仓/probe-meta/aliyun CLI(drifter)) |
| `bootstrap/probe.sh list` | 列全部场景(PRODUCT/id/persona/resources/detect + 行尾 LAST_RUN 列,LRU 索引 `.my-day/probe/t1-last-run.json` 优先,索引缺回退扫 `runs/probe/` 文件名兼容新旧两种正则) |
| `bootstrap/probe.sh tier0 [--no-mech] [--all] [--limit N] [--rotate N] [alicloud_xxx ...] [--dry]` | tier-0 静态三方一致性扫描(含 T0-mech OpenAPI 机械 diff) |
| `bootstrap/probe-meta.sh {fetch\|cached-fetch\|clear\|available}` | T0-mech OpenAPI 元数据获取层(薄封装 amp skill + cache.sh 7d) |
| `bootstrap/probe.sh run <id> [--region r] [--dry] [--keep]` | tier-1 真实 apply 生命周期探测(支持新键:steps CSV / step\<N\>_expect / expect_fail(+expect_error_contains) / provider_version_from / drift_cli) |
| `bootstrap/probe.sh sweep` | 扫残留 state,残留退 1 |
| `bootstrap/probe.sh archive [--dry]` | **归档 + 台账**:draft 终态 → archived/;verdict retention (60d,排 ledger/summary) → `runs/probe/archive/<YYYYMM>/`;`.my-day/probe/<ts-sid>` 空 state 且过期 → rm(排 .plugin-cache/manual-*/索引文件);plugin-cache 陌生版本只报体积不删;pending drafts + _quarantine + origin:generated 待办清单;ledger 追加 kind=archive(dry=archive_dry) |
| `bootstrap/rc-gate.sh <provider-dir> [--quick]` | **RC 门禁线**:tier-0 全量 + tier-1 全场景一次过闸,红/黄/绿判定,报告落 `runs/rc-gate/<date>-report.md`(退码 红1/黄0/绿0/不可判2);见「四点五」 |
| `config/probe.json` | regions / tiers(tier1.enabled, prepaid_guard, drift_enabled, drift_action_allow) / limits(audit_retention_days=60, workdir_retention_days=7) / ticket / paths(drafts / drafts_archived / playground_dir) |
| `tf_playground/<product>/<id>/`(独立 git 数据仓) | tier-1 场景语料库(云产品维度两级布局);根解析 env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` > `bootstrap/workspace.sh dir tf_playground` > 默认 `<jarvis 父目录>/terraform_playground` |
| `tf_playground/<product>/KNOWLEDGE.md` | jarvis 蒸馏的产品级可执行知识(五节结构);契约见 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md`(跨 skill 单点,probe/triage/dev 三处消费) |
| `runs/probe/ledger.jsonl` | 本机台账 append-only(tier0/tier1/ticket/archive);真源在 Aone |
| `.my-day/probe/t1-last-run.json` | LRU 索引 scenario→epoch(list 行尾 LAST_RUN 显示;仅推进到 plan 完成及以后才更新) |
| `.claude/skills/tf-customer-probe` | 全流程技能 + references(scenario-authoring/severity-rubric/ticket-template/knowledge-distillation) |
| `escalation/cap-tf-customer-probe.md` | 能力路线图 |

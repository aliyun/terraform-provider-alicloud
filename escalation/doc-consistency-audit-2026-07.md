# doc-consistency-audit-2026-07 —— 执行逻辑冲突审计 + skill/loop 当前态化清理

> Aone 83956319。审计范围:`.claude/skills/**`(SKILL.md + references/)、`loops/*.md`、`CLAUDE.md`、
> `autonomy.md`;cap-* 仅审「当前态陈述错误」不清历史。判别标准:执行指令类文件只保留当前态,
> 「为什么」的历史归 cap,确属护栏理由的一句话 why 可留。
>
> **真源基准**(逐主题以代码/config 为准做跨文件断言对齐):
> - probe 建单模式:`config/probe.json` `.ticket.mode=file`。
> - 分层:tier-0=三方一致性静态扫描 / tier-1=真实 apply 全生命周期(`config .tiers`)。
> - playground 路径链:`bootstrap/probe.sh` `probe_playground_dir()` = env `JARVIS_TF_PLAYGROUND` > config `.paths.playground_dir` > `workspace.sh dir tf_playground` > 父目录默认(**4 步**)。
> - 成本/订阅门:`bootstrap/probe-corpus.sh` L180-401 值敏感 + 有/无 ds 分岔。
> - dispatcher 语义:`bridge/jarvis_dingtalk_bot.py` ScanScheduler `_tick`/`_decide`/`_in_scope`/`_human_touched`。

## 结论摘要

- **真冲突 10 处 / 归 6 类**(文档陈述与当前实现/config 矛盾,AI 照做会错):见下表 ★。
  6 类 = ① probe mode 现/未来倒置(2 文件)② playground 路径链漏 `workspace.sh`(3 文件)③ 订阅门漏「无 ds→apply:true」分岔(2 文件)④ `probe-corpus.sh` 注释命名已删字段 ⑤ dispatcher「只新单+无条件跳 idle」失真 ⑥ cap-flywheel F3「待合并」过期。
- **冗余/历史叙事 5 类**(不算错但违「当前态」原则,归 cap):见下表 ○。
- **核对无冲突 3 项**(verify.sh 新检查 / rc-gate 产物路径 / run.sh 单一入口):文档已一致,无需改。
- 改动 16 文件(5 skill mirror 对 + `probe-corpus.sh` + `config/probe.json` + 2 cap + 2 loop)。测试见文末。

---

## 冲突矩阵(主题 × 文件 × 冲突内容 × 处置)

### 主题 1 — probe 建单模式(现行 `mode=file` 直接建单,config 为真源)

| 类 | 文件:位置 | 冲突内容 | 处置 |
|----|-----------|----------|------|
| ★真冲突 | `loops/tf-probe.md` 产物分流表 | 「去重后 → **draft(当前 mode)**…毕业后 → 建单」——把当前 mode 说成 draft,与 config `mode=file` 直接矛盾 | 改为「当前 `mode=file` 直接 adhoc-intake 建单;`mode=draft`(可回退)时改落 drafts」 |
| ★真冲突 | `references/ticket-template.md` L3 / L58 | 「或**未来毕业后**按此建 Aone 单」「建单参数块(**毕业后** mode=file 用)」——把已生效的 file 当未来事件(现/未来倒置) | 翻正:当前 file 直接建单;draft 为可回退开关 |
| ○冗余 | `tf-customer-probe/SKILL.md` §D 标题+正文 | 「(当前 mode=file,**2026-07-05 已毕业**)」「主人拍板毕业:采纳率 7/8=87.5% 提前毕业」——毕业叙事 | 标题留 `mode=file`;draft 压成一句「可回退开关,切回即恢复人审」;毕业史归 cap |
| ○冗余 | `references/severity-rubric.md` L24 | 「首次真实建单(`ticket.mode=file` **毕业后**)」 | 去「毕业后」,改「首次真实建单(`mode=file`)时」 |
| 补全 | `tf-customer-probe/SKILL.md` §D | skill 未写「generated 场景 apply_fail 先归校订不建单」(只在 authoring.md) | §D 补一行 generated 校订门,使 skill 建单流程自洽 |

### 主题 2 — tier 定义(清历史演变叙事,留 cap)

| 类 | 文件:位置 | 冲突内容 | 处置 |
|----|-----------|----------|------|
| ○冗余 | `SKILL.md` L19 / `loops/tf-probe.md` L7 / `authoring.md` L5 | 章节标题/引言带「(2026-07-03 重定义)」历史标记 | 删日期标记,只留「## 分层」当前态陈述 |
| ○冗余 | `SKILL.md` L22 | 「(T0-mech,**2026-07-05**:probe-meta.sh…)」日期标记 | 删日期,保留机制描述 |
| 核对无冲突 | `SKILL.md` L29 / `tf-probe.md` L74 / `authoring.md` L85 「init/validate/plan」 | 这是 **plan-only 封顶**的当前描述(非旧 tier-0 定义),正确 | 不动 |

> 「原 tier-0=init/validate/plan」「2026-07-03 重定义」的演变史保留在 `cap-tf-customer-probe.md` 决策记录(2026-07-03 段)。

### 主题 3 — playground 路径链(4 步链一致 + git 数据仓 + 清 probes/ 残留)

| 类 | 文件:位置 | 冲突内容 | 处置 |
|----|-----------|----------|------|
| ★真冲突 | `SKILL.md` L13-16 | 路径链只写 3 步(env > config > 默认),**漏 `workspace.sh dir tf_playground`**;「外置 jarvis 仓外」旧表述 | 补 4 步链;改「独立 git 数据仓 `tf_playground`(直推 master + 工单报备)」 |
| ★真冲突 | `loops/tf-probe.md` L198(工具表) | 同上 3 步链漏 workspace.sh;「(外置,仓外)」 | 补 4 步链 + git 数据仓表述 |
| ★真冲突 | `cap-tf-customer-probe.md` L44(架构表) | 架构表 3 步链漏 workspace.sh(当前态陈述错误) | 补 4 步链(cap 架构/现状表对齐现行) |
| ○冗余 | `probe-ticket-routing.md` L37/L41 | 「外置在 jarvis 仓外」;L41「原…周批 MR 入 **probes/**」已废弃残留 | 改 git 数据仓表述;删 L41 probes/ 废弃史(归 cap) |
| ○冗余 | `authoring.md` L158 / `loops/tf-probe.md` L175 | 「取代原 escalation/scenario-drafts + 周批 MR(已废弃)」;「直落外置」 | 删废弃史尾巴;「外置」→「git 数据仓」 |
| 核对无冲突 | `authoring.md` L25-27 | 4 步链**已正确**(含 workspace.sh) | 不动(基准样板) |
| 保留历史 | `cap-tf-customer-probe.md` L129 | 2026-07-03 决策记录里的 3 步链 | 不动(历史,当时确为 3 步;workspace.sh 登记是 07-06 git 化后新增) |

### 主题 4 — 成本/订阅门(值敏感 + 有 ds→apply:false+ds-变体 / 无 ds→apply:true+allow_prepaid)

| 类 | 文件:位置 | 冲突内容 | 处置 |
|----|-----------|----------|------|
| ★真冲突 | `authoring.md` §apply 门 / §订阅规范 | 「命中订阅门 → `apply:false`(一律)」——漏 07-06「**无 ds→apply:true+allow_prepaid**」分岔,与生成器 `probe-corpus.sh` L343-401 矛盾 | 加有/无 ds 分岔;标题「有 data source 引存量,无则放行真跑」 |
| ★真冲突 | `bootstrap/probe-corpus.sh` L20-21(头注释) | 注释仍写「命中 `resource_patterns` / `charge_field_patterns`→apply:false」——命名**已删字段**,与自身 L180+ 实现矛盾 | 注释改值敏感 + 有/无 ds 分岔(仅注释,不动逻辑) |
| ★真冲突 | `config/probe.json` `.tier1_risk_denylist.desc` | 真源 desc「命中即写 apply:false」漏 ds 分岔 | desc 补有/无 ds 分岔,使真源自洽 |
| 核对无冲突 | `authoring.md` §资源选择 prepaid 守门 | 「成本白名单已撤销 / prepaid 销毁性守门」当前正确 | 不动 |

### 主题 5 — dispatcher 语义(新单+更新单可派发 + 人工介入门 + 灰度 + 暂停 + idle skip)

| 类 | 文件:位置 | 冲突内容 | 处置 |
|----|-----------|----------|------|
| ★真冲突 | `loops/aone-triage.md` §一 ScanScheduler 行 | 「diff 出**新工作项**」(只新单)+「跳过 `jarvis-idle`」(无条件)——漏当前:更新单也派发、jarvis-idle 过**人工介入门**(activity 作者判据)才 skip/重派、灰度安全阀、暂停开关 | 重写:新单+外部更新单一并派发;`_decide` 逐单判定(终态/done/claimed skip;idle 过 `_human_touched`);`_in_scope` 灰度;`.my-day/bridge/pause` 暂停;附「无钉钉降级」 |
| 核对无冲突 | `cap-probe-fix-flywheel.md` 当前态段 | grep「更新单只播报」**无残留**;「钉钉播报」指卡片信息化(现行正确),非「更新单只播报」 | 不动(cap 无被取代语义) |
| ★真冲突 | `cap-probe-fix-flywheel.md` §F3 阶段计划 | MR-4~8「⏳/待仓库主人合并」——**全部已合入**(07-06 `worktree-f3-consolidated` 整合) | ⏳→✅;合并状态改「全数合入 master」;删「本收尾单不碰」过期提示(cap 现状表对齐现行,决策记录不动) |

> cap-probe-fix-flywheel 决策记录未单列「新单+更新单+人工介入门」增强(该特性晚于 07-06 记录落地)。因 cap 无处
> **断言旧语义为当前**(仅 F2 决策记录如实描述 07-03 快照),按「决策记录里的历史不动」不改;当前态真源是已更新的 loop。

### 主题 6 — 杂项

| 类 | 文件:位置 | 结论 | 处置 |
|----|-----------|------|------|
| 核对无冲突 | `bootstrap/verify.sh` 新检查(`jarvis-github-token` / `jarvis-a1-session`,均 WARN) | loop/skill 未逐项枚举 verify 子检查;与 `CLAUDE.md #6`(github-identity 硬门)方向一致,无矛盾 | 不动 |
| 核对无冲突 | rc-gate 产物路径 `runs/rc-gate/<date>-report.md` | `loops/tf-probe.md` §四点五 与 `terraform-provider-release/references/rc-gate.md` 一致 | 不动 |
| 核对无冲突 | `bridge/run.sh` 单一入口 | `CLAUDE.md #3` 与 `loops/aone-triage.md` 表述一致(本次 loop 补「无钉钉降级」) | 已随主题 5 一并对齐 |
| 核对无冲突 | 旧 MR 号 / 「待合并」残留 | 执行指令类文件**无**(仅 `CLAUDE.md #17` 泛化 PR/MR 政策,非过期);过期「待合并」仅在 cap-flywheel F3(已随主题 5 修) | 已修 |

---

## 改动文件清单(16)

**skill(mirror 对,各 2 份 .claude/.agents):**
- `tf-customer-probe/SKILL.md`、`references/{scenario-authoring,ticket-template,severity-rubric}.md`
- `aone-triage/references/probe-ticket-routing.md`

**loop / config / 脚本注释:**
- `loops/aone-triage.md`(dispatcher)、`loops/tf-probe.md`(mode/tier/playground)
- `config/probe.json`(订阅门 desc)、`bootstrap/probe-corpus.sh`(头注释)

**cap(仅当前态陈述纠正,历史决策记录不动):**
- `cap-tf-customer-probe.md`(架构表 playground 4 步链)、`cap-probe-fix-flywheel.md`(F3 合并状态)

## 测试全景

`.claude/skills` 改动经 `bootstrap/mirror.sh to-codex` 同步 `.agents`,5 个改动 skill 文件 mirror-check 全 CLEAN。

| 套件 | 结果 |
|------|------|
| aone_triage_templates_sync / html_report_preview_skill_sync / terraform_pr_review_skill_rules / provider_internal_aone_rules | PASS |
| probe / verify / workspaces_config / claim / board_probe | PASS |
| **(baseline 既红,与本次改动无关)** provider_resource_dev_skill_sync / acctest_remote_skill_rules / sync_provider_dedup | FAIL(baseline 同样 FAIL) |
| **(baseline 既红)** `mirror.sh check` 全量 | DRIFT(缺 `AGENTS.md` + provider-resource-dev `.Codex` token 漂移,均本次未触文件) |

> 3 个红套件 + 全量 mirror drift 均为**先存缺陷**,涉及 `provider-resource-dev` / `invoke-terraform-acc-test-remote` /
> `sync-provider.sh` / 顶层 `AGENTS.md`——**均非本次改动文件**,base commit 4416364 上同样红。本次改动**零新增失败**,
> 所改文件 mirror 全绿。先存缺陷建议另开单修(共享工具/mirror bug 走 fix+CR,不并入本 doc-consistency commit 以免混淆)。

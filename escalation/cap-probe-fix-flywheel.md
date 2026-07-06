# cap-probe-fix-flywheel

## 背景

项目终态 = 数字员工自治运转。探测能力（`cap-tf-customer-probe`）P0 已建成并实战验证：2026-07-03 首轮
产出 **7 张 probe 工单**（83881282 / 83881291 / 83881297 / 83881301 / 83881314 / 83882090 / 83882094,
均指派 jarvis、打 `jarvis-probe` 标签）+ **1 份 verify 否决样本**（判定层挡下的伪 finding)。

「探测出来」只是半程。本 cap 设计**「发现的问题被 jarvis 自己修完并回归验证」**的完整飞轮——从发现到发布再到
回灌复跑,形成自闭环,让整个项目在三个人工硬门之外自动运转。

## 探测侧完成度判定（如实）

**已建成（实战验证）**：

- tier-0 机械扫描 + OpenAPI 判定管道（文档↔源码 diff + `judgment_queue`）
- tier-1 真实 apply 全生命周期（init→…→destroy）
- findings / env_issues 严格分流
- draft → 人审 → 建单管道（`escalation/probe-drafts/`）
- 去重台账（`jarvis-probe` 标签 + draft filed 回写）
- 护栏（prepaid 销毁性守门 / 强制 destroy / sweep 残留核查 / a1id 身份 / 测试账号边界）
- 审计与 **154 项 hermetic 测试**

**未完成（自动运转三缺口）**：

1. **调度**——✅ **已落地（F2，2026-07-03）**：bridge ScanScheduler auto-dispatch 并发派发（每单一 headless、上限 `JARVIS_DISPATCH_MAX`、软去重台账、钉钉播报）+ ProbeScheduler 每日探测轮 + RevisitScheduler 每日人工门重访；`--dry-run-once` 可离线验证决策。
2. **规模**——✅ **首发线已落地（F3 T0-mech，2026-07-05）**：`tier0 --all/--limit/--rotate` 支持 website/docs/r 全资源清单 + LRU 轮换巡检；未竟：场景语料生成器 + cloudspec 覆盖矩阵（tier-1 侧）。
3. **判定成本**——✅ **首发线已落地（F3 T0-mech，2026-07-05）**：OpenAPI 侧从「verifier 人海判定」升级为「机械 diff 预筛（六类 `api_gap_*`）+ verifier 只判 `judgment_queue` 疑点」；`bootstrap/probe-meta.sh` 薄封装 amp 元数据 + `cache.sh` 7d 缓存，`probe.sh tier0` 抽 (product,version,action) 三元组 + `StringInSlice`/`IntBetween`/`Default` 解析做机械 diff。精度护栏:拿不准一律 queue + 抑制表 + 容差表 + coverage_note。首轮 8 资源标定复现 ClassicLink deprecated、RAM/SG 零新增误报、OSS SDK 风格正确降级为 queue。

## 飞轮六段架构

| 段 | 触发 | 执行机件 | 现状 | 缺口 | 人工门 |
|----|------|----------|------|------|--------|
| ① 发现 Probe | 定时 + 新版本发布 | `tf-customer-probe` skill + `probe.sh`(tier-0/tier-1) + bridge ProbeScheduler(每日轮) | 已建 | ✅ bridge 每日探测轮已接入(F2);未竟:provider 新版本事件触发 | 无 |
| ② 立项 File | findings → 分级 → 去重 → 建单 528766 指派 jarvis | draft 管道 + 建单模板/标签/上限(100) | ✅ `mode=file` **已毕业**(2026-07-05 主人拍板提前毕业,首轮采纳率 7/8=87.5%);直接建单 | draft 保留可回退开关(`ticket.mode=draft`) | 无(临时人审门已撤) |
| ③ 修复 Fix | bridge ScanScheduler 扫池(probe 单指派 jarvis 被 `scan.sh` 自然扫到) → headless dispatch → aone-triage 认领 → **按 `probe-ticket-routing` 路由** | (a) provider 代码修 → `provider-resource-dev` → fork+UT → `invoke-terraform-acc-test-remote` 验收 → GitHub PR(`github-identity` 硬门,`api-tool-agent`);(b) TF 文档修 → 同 PR 路径 docs-only;(c) 上游协作 → cloudspec_gap 等 `submit_only` 转发;(d) 需实验定性 → 先跑 tier-1 变体场景再归入 a/b/c | 机件**全部已存在** | 只缺**路由规范**(本 commit 落地) | upstream PR merge(maintainer) |
| ④ 验证 Verify-fix | PR 合并 → master 复验 | tier-0 重扫该项应消失 / tier-1 场景复跑应绿 →「已修未发布」→ 发布后复跑绿 → `claim.sh finish`(jarvis-done) | 靠工单溯源字段映射回场景/资源,无需新 runner 子命令 | 复验编排规范(routing reference 内定义)+ 状态机落 tag/评论 | 无 |
| ⑤ 发布 Release | changelog 聚合 → 发版 | `terraform-changelog` skill(已有) | 已有 | 发布前 RC 门禁 = 全场景语料过一遍(P2) | release_prod(autonomy.md 永久停止项) |
| ⑥ 回灌 Regress | 每张修复完成的 probe 单 + 每张真实客户单 | `regression-<aone-id>` 场景落 **git 化 `terraflow/tf_playground`**(2026-07-06:直推 master + 工单报备,取代原仓外裸目录)+ 工单评论报备场景路径 | 规则已立(2026-07-03 外置 → 2026-07-06 git 化) | 收尾清单挂钩 | 工单报备(仓库主人查验;git 直推 master 无 MR 门) |

## 三个永久人工硬门（其余全自动）

1. **jarvis 仓 MR 合并**（仓库主人）。
2. **terraform-provider-alicloud upstream PR merge**（团队 maintainer;jarvis 可自动 remind / 催）。
3. **release_prod**（正式发布,autonomy.md 永久停止项）。

外加 **draft 毕业前的人审 draft**（临时门,`mode` 毕业翻 `file` 后撤）。

## 兜底与健康度

- **escalate 面不变**：`low_conf` / `verify_fail` / `redline` / `missing_capability` + probe 新增 `destroy_fail` / 残留。
- **僵死收敛**：`reconcile.sh` / watchdog / coord 孤儿收养。
- **飞轮度量**（P3 接 `board.sh`）：每轮发现数、draft 采纳率、建单→修复→发布周期、回归通过率、场景覆盖率。

## 阶段计划

- **F0（本 commit）**：飞轮 cap 设计 + aone-triage probe 单路由 reference——已有的 7 张单下一轮 triage 即可正确消化。
- **F1（首证）**：选 **83881291**（dns_hostname_status,最干净的代码修）人工触发走完整链一遍:
  `provider-resource-dev` 修复 → acc test → PR → merge 后 master 复验 → 回灌 regression 场景 → finish。
  产出飞轮首个端到端样本与卡点清单。
- **F2（自动化）**：**已落地（2026-07-03，本 commit）**——bridge ScanScheduler auto-dispatch 并发派发（每单一 headless、
  上限 `JARVIS_DISPATCH_MAX`、软去重台账 `.my-day/bridge/dispatched.json`、钉钉播报、授权前置降为 `JARVIS_AUTO_DISPATCH=0` 回退）
  + DispatchPool（并发/排队/软去重复用 `_dispatch_bg` 核心与 SUSPEND/WaitWatcher）+ ProbeScheduler 每日探测轮（跑 `loops/tf-probe.md`）
  + RevisitScheduler 每日 `jarvis-idle` 人工门重访 + `--dry-run-once` 验证入口 + hermetic 单测 `test/bridge_dispatch_test.sh`。
  **未竟**：provider 新版本事件触发；复验步骤进 triage 收尾清单。（`mode=file` 毕业已于 2026-07-05 主人拍板翻开关。）
- **F3（规模化）已落地清单**——机械化 / 语料 / 门禁 / 度量 / 降级 / 单一入口全数成 MR（2026-07-05～06），**已全数合入 master**（2026-07-06 经 `worktree-f3-consolidated` 整合合并落地）：
  - ✅ **T0-mech**（`worktree-f3-t0mech`）—— tier-0 OpenAPI 侧机械化（`probe-meta.sh` + `tier0` 六类
    `api_gap_*` 机械 diff 预筛 + `--all/--rotate` 全资源轮换巡检 + doctor 探针可用性降级）。**已合入 master**（2026-07-05）。
  - ✅ **Corpus-gen 场景语料生成器**（`worktree-f3-corpus-gen`）—— website docs 全量 → tier-1 语料。
  - ✅ **RC 门禁**（`worktree-f3-rc-gate`）—— 发布前全资源 tier-0 + 全场景 tier-1 过一遍，产物落 `runs/rc-gate/`。
  - ✅ **度量看板**（`worktree-f3-board`）—— 发现数 / 采纳率 / 发现→修复周期 / 回归通过率 / 覆盖率，接 `board.sh`。
  - ✅ **无钉钉降级**（`worktree-f3-nodingtalk`）—— bridge 缺钉钉凭证时干净降级，不阻断 scan/dispatch/probe 调度。
  - ✅ **run.sh 单一入口**（`worktree-bridge-run-entry`）—— `bridge/run.sh` 收敛后台启动入口。
  - **合并状态**：全数合入 master——T0-mech（2026-07-05）在先，规模化四件（Corpus-gen / RC 门禁 / 度量看板 /
    无钉钉降级）+ 单一入口 run.sh 于 2026-07-06 经 `worktree-f3-consolidated` 整合合并落地。
- **F4（目标态）**：无人值守运转,人只守三硬门与 escalation 队列。**多机运营不单立 cap**
  （2026-07-06 主人定调,兜底靠 claim 点读锁 + `tf_playground` git 收敛,见决策记录 2026-07-06 附注）。

## 决策记录

- **2026-07-03**：仓库主人提问「探测部分是否建成?如何让 jarvis 自闭环解决探测出的问题,让整个项目自动运转」并指示继续设计
  → 本 cap（飞轮总设计）+ F0 落地（aone-triage probe 单路由 reference)。
- **2026-07-03（调度指令，方案 jarvis 设计）**：仓库主人指示「定时检查 Aone 池状态，发现新工单直接起独立 headless jarvis
  实例并发处理，授权前置降级为可配回退；同时补齐 probe 探测轮与人工门重访两个定时器」。developer 子代理实现（Aone 83902495）：
  ScanScheduler auto-dispatch（`JARVIS_AUTO_DISPATCH` 默认开）+ DispatchPool（并发上限 `JARVIS_DISPATCH_MAX` 默认 2 / 队列 20 /
  软去重 24h，持久化 `.my-day/bridge/dispatched.json`，复用 `_dispatch_bg` 核心与 SUSPEND/WaitWatcher）+ ProbeScheduler（默认 10 点）
  + RevisitScheduler（默认 9 点，扫各池 `jarvis-idle` 的 `[probe]`/待续条件工单）+ `--dry-run-once` 验证入口。红线遵守：
  dispatcher 永不直接写 Aone、不 adopt、不动 TataPool 对话行为；claim.sh 仍是竞争互斥真源，dispatcher 只做软去重。
  **决策补充**：scan auto 除跳过 `jarvis-claimed`/`jarvis-done` 外也跳过 `jarvis-idle`（人工门/已 park 项归 RevisitScheduler 每日节拍，
  避免 30 分钟 scan 反复 spin 空耗实例）——较原规格新增 idle 过滤（2026-07-05 主人追认转正式设计，见下）。
- **2026-07-05（主人批准两项）**：① **`mode=file` 毕业**——首轮探测采纳率 7/8=87.5%，主人拍板**提前毕业**（不再等
  「累计 ≥10 draft 且 ≥90%」门槛）：`config/probe.json` `ticket.mode` `draft`→`file`，findings 直接建 Aone 单、
  不再走 draft 人审；draft 模式保留为**可回退开关**（临时收回自动建单权时切回）。② **dispatcher idle 过滤转正**——
  scan auto 跳过 `jarvis-idle`（归 RevisitScheduler 每日重访）由上轮「新增待追认」转为**正式设计**。
- **2026-07-05（F3 首发线 T0-mech，developer 子代理实现，Aone 83929676）**：tier-0 OpenAPI 侧机械化落地。新增
  `bootstrap/probe-meta.sh`（薄封装 amp `get_api_definition.py` + `cache.sh` 7d 缓存，fetch/cached-fetch/clear/available，
  probe.sh 不直接调 python）;`probe.sh tier0` 升级——(product,version,action) 三元组抽取 + 源码 `StringInSlice`/`IntBetween`/`Default`
  解析 + 六类 `api_gap_*` 机械 diff + 抑制表/容差表/coverage_note 精度护栏 + `--all/--limit/--rotate` 全资源 LRU 轮换巡检;doctor 增
  probe-meta 可用性检查（不可用 WARN + 自动降级为纯 doc↔source 现行为）。hermetic 单测 154→206 断言（PATH 桩替底层 python + fixture
  五类陷阱）。**首轮 8 资源真实标定**：机械层复现 `alicloud_vpc` ClassicLink `api_gap_deprecated_action` ×2（Enable/DisableVpcClassicLink，
  = 07-03 单 83881282）;`alicloud_ram_role`/`alicloud_security_group` 零新增误报（三方一致）;`alicloud_oss_bucket` SDK 风格抽不到三元组 →
  正确降级为 `no_triple` queue（enum 超集需元数据可得，符合规格）;`dns_hostname_status`（convert 改名）→ `unmapped_params` queue 不硬猜。
  **红线遵守**:精度命门=拿不准一律 queue 不硬报;本机无 AMP 凭证时自动降级。
- **2026-07-06（主人指令四条 + 多机运营定调）**：
  1. **成本门收窄为「值敏感」**：prepaid / 成本守门不再对所有 PrePaid/Subscription 一刀切阻断，收窄为**只对真正值敏感
     （产生不可逆计费 / 无法 API 销毁）的场景**把关；**订阅类资源若无对应 data source（无 ds 可读回校验）则允许 `apply`**
     （等价 `allow_prepaid`），避免探测面被过度锁死漏探。
  2. **场景语料库 git 化**：`terraform_playground` 从「仓外裸目录」升级为 git 仓 **`terraflow/tf_playground`**，采用
     **直推 master + 工单报备**模型（jarvis 直接 push master、工单评论报备场景路径，无 MR 人工门）。取代原「外置无 MR」
     简化方案，同时天然收敛多机语料分叉（见下附注）。
  3. **单一入口 `bridge/run.sh`**：bridge 后台启动收敛到单一入口脚本 `bridge/run.sh`（MR-8, `worktree-bridge-run-entry`），
     统一 env / 日志 / 各调度器拉起。
  4. **AMP 凭证已配置 → 机械面全量点亮**：本机已配置 AMP 白名单凭证，T0-mech 机械 diff 不再受「本机无凭证」约束，
     `tier0 --all` 全量实拉元数据可端到端跑（此前 cap 记的「待有凭证环境验证」约束解除）。

  **附注（多机运营，主人定调）**：**不为「多机运营」单立 F4 cap**。理由 / 兜底要点：
  - **修复类多机安全**靠 `claim.sh` 的**点读锁**（claim 竞争互斥 + owner_instance）兜住——多实例并发认领同一工单不会撞车；
  - **bridge / 日轮**（scan 派发、probe 轮、revisit）**单主机跑**，不做多主机分布式调度；
  - **语料分叉**已由 **`tf_playground` git 化**（指令 2）解决——各机 push 同一 git 仓，天然收敛；
  - **跨机去重台账**（dispatched 软去重当前 per-machine 落 `.my-day/`）**待真有多机需求再做**，现阶段不预建。
- **待主人决策**：
  1. ~~`mode=file` 翻开关时点~~ → **已定（2026-07-05）**：主人拍板提前毕业（采纳率 87.5%），draft 留作回退。
  2. ~~bridge probe cron 频率~~ → **已定（2026-07-03）**：probe 每日 `JARVIS_PROBE_HOUR`（默认 10）一轮、revisit 每日
     `JARVIS_REVISIT_HOUR`（默认 9）一轮；频率/并发均可 env 调，运行一段后按实际负载再校准。
  3. F1 首证启动授权。

## 关联

- 探测能力本体：`cap-tf-customer-probe.md`。
- 能力建设单：https://project.aone.alibaba-inc.com/v2/project/2100304/req/83879813
- MR：https://code.alibaba-inc.com/terraflow/jarvis/codereview/28364380
- 首轮 7 张 probe 单：83881282 / 83881291 / 83881297 / 83881301 / 83881314 / 83882090 / 83882094。
- F0 路由 reference：`.claude/skills/aone-triage/references/probe-ticket-routing.md`。

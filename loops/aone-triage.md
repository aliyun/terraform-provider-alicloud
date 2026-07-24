# aone-triage 无人值守到预发 triage loop

> 单条工单怎么处理在 `.claude/skills/aone-triage`(读单/归类/查证/回复/bookend);本文件只管 **loop 编排**:控制面 Session、触发入口、认领竞争、autonomy 判定、收敛维护。

---

## 零、实例协调（控制面）

每个持久任务只通过控制面的 Worker → Task → fenced Session 生命周期执行。bridge 注册稳定 Worker，`PersistenceExecutor` lease Task，`SessionController` 负责 start/heartbeat/suspend/complete；交互式入口由 `jarvis-interactive-worker.py` 使用同一套 API。业务上下文、分支、transcript 与 result refs 随 Session 外化，不再写本地 coord task/checkpoint 文件。

- Worker/Session 心跳和 fence 是存活与写权限真源；本机 PID 只用于终止当前进程，不参与跨机接管判定。
- 进程或机器中断后由服务端 reaper 收敛 Session，`RECOVERY_REQUIRED`/`RESUMABLE` Task 再由 bridge 恢复调度；不存在本地 orphan scan/adopt。
- `SUSPENDED` Task 保留 wait/affinity 与续跑上下文，不占执行槽；换机后从控制面恢复。
- 控制面不可用时持久任务 fail-closed，禁止降级为本地无追踪执行。详见 [docs/execution-architecture.md](../docs/execution-architecture.md) 与 [docs/multi-worker-deployment.md](../docs/multi-worker-deployment.md)。

---

## 一、触发与输入

| 方式 | 说明 |
|------|------|
| bridge 定时扫池（自动派发） | Scheduler 的 **`scan` runner** 周期做 Python 直查并集探测（每池 `assignedTo∪workitem.tracker∪tag=jarvis-idle` × `DIGITAL_WORKER_IDS`，池间+池内并行），把新单与外部更新单统一 upsert 为控制面 Task；`bootstrap/scan.sh` 只保留人工审计/兜底。Persistent Worker 按 `JARVIS_DISPATCH_MAX` lease 执行，控制面 Task/Session/lease/fence 是唯一真源，失败时不回退本地无状态执行。`_decide` 对终态、`jarvis-done`、`jarvis-claimed`、`jarvis-npe` 跳过；`jarvis-idle` 仅在检测到上轮之后的人工介入时唤醒。范围安全阀为 `JARVIS_DISPATCH_POOLS`、`JARVIS_DISPATCH_CREATED_BEFORE` 与 `.my-day/bridge/pause`。daily_probe、claim_health、daily_nudge、reply、pr_watch、recovery 分别由独立 runner 执行。唯一进程入口是 **`bridge/run.sh start`**。 |
| bridge dispatch | Tata 委派单工单，headless 执行（autonomy.md headless 模式：auto 列表免授权、遇阻 `[[SUSPEND:...]]` 挂起） |
| 用户指令 | 会话里给 Aone URL / 工单 id → 直接进「二、逐项执行」单条流程 |
| 手动兜底 | `/aone-triage` 或手动跑 `bootstrap/scan.sh`（排查/对账用；`plan.sh` 供 bridge/serve 流程出计划） |

- **supervised（默认交互模式）**：Aone 写操作逐项等用户授权后执行；unattended/headless 按 `autonomy.md` auto 列表放行。
- **分诊是 Claude 做，不是脚本扫表**：按优先级/标题/标签/状态排序 + 折叠噪声，脚本只出数据。

### 产物落点纪律

- 计划/审计走 `runs/`（run_done/escalate 记录）。
- 临时数据用 `.my-day/`（gitignored）或 `mktemp`，**禁止**往仓库根甩 `scratch_*`/`*.tmp`。
- 开发不在 master：worktree 切分支 → MR → 人工合并；`.gitignore` 兜底脏文件。

---

## 二、逐项执行

对每条已授权（或 headless auto）的工作项，按以下顺序处理：

> **Terraform 单写者约束**：PD/RD/QA 只作为同一 headless run 内的 internal subagent。
> PD/QA 纯内部、禁止外写，开发阶段 RD 也不发工单进展；最后由 terraform-rd finalizer 汇总
> 全部证据并回复一次。Terraform 主处理 run 不做中途 `wrap.sh sync`，MR/CR 链接放最终聚合
> 回复；后续重要事件由 bridge 以 TerraformRD 身份幂等更新，semantic source 仅以短摘要
> 进入 ledger/marker，正文统一 sanitize，Aone `post_uncertain` 只查 marker 不重发。满 8 天
> 无实质进展的重访催办同时走 Aone 评论区 @ 与钉钉私信，两个通道独立补偿；无变化/重复事件静默。
> claim/final wrap/release/event 命令必须走严格 RD 身份，RD 未登录即阻断。

### 2.1 去重检查

```bash
bootstrap/log.sh seen <id>
```

- `seen` 返回 0（已处理）→ 跳过本条，继续下一条。
- `seen` 返回 1（未处理）→ 进入认领。

### 2.2 认领（并发竞争锁）

```bash
# 非 Terraform
bootstrap/claim.sh claim <id> <pool-project>

# Terraform
JARVIS_A1_IDENTITY=terraform-rd bootstrap/claim.sh claim <id> <pool-project>
```

- 退出码 0 → 认领成功，继续 triage。
- 退出码 1 → 其他实例已抢先认领，**跳过本条**，继续下一条。
- 退出码 3 → 工单缺必填字段（结构化 `【name】(fieldId)不能为空`，非 lost race）：运行
  `bootstrap/aone-fields.sh missing <id>` 获取合法候选，agent 结合工单语义明确选值后执行
  `bootstrap/aone-fields.sh fill <id> <fieldId>=<value>`，再重试 claim；禁止盲填分类。
- 其它非零 → claim 真失败，直接 escalate，不进入 readback 或 SKIP。

### 2.3 单条 triage（技能调用）

调用 `.claude/skills/aone-triage` 技能，传入当前工作项 id。

Terraform 线由编排层在同一 run 内依次 Task 起 PD→RD→QA：PD 返回查证和路由提案，RD 开发，
QA 独立验收；QA fail 内部退回 RD 修复后重跑。最后再 Task 起 RD finalizer，执行允许的路由动作
并一次性回复。非 Terraform 线仍按技能原流程读取、查证、回复、打标、建需求或建 CR。

Terraform PD 的三层查证还必须调用 `screenshot-evidence` 生成本地截图和
`evidence-manifest.md`，但不得上传或回贴；RD finalizer 校验 manifest 后统一上传一次报告，
上传命令不传 `--comment`，报告链接只进入本轮唯一聚合回复。三层缺失且无明确 N/A 原因时不得
静默收尾。

若单条工单进入 Terraform Provider 资源开发且不走自动化生成链路，按 `tf_provider` 池创建或复用
**terraform-alicloud** 内部研发单（项目 `528766`），指派按 aone-triage skill
`references/tf-customer-request-routing.md` 分工表路由到具体人，并与客户主单双向关联。路由动作由
最终 RD 审查执行；每个被 claim 的 Terraform 工单在本轮主处理 run 只于最终聚合时回复一次，
不按 PD/RD/QA 阶段同步。后续 gate/PR/终态失败的重要事件可按 `loops/persona-collab.md`
§五追加 RD-only 幂等更新。需要 cloudspec_gap 或云产品上游协助时，把详细问题放入最终聚合及对应依赖单。

GitHub PR/评论/推分支的身份纪律见 CLAUDE.md 工作纪律 #6（`bootstrap/github-identity.sh`，账号必须 `api-tool-agent`）。

### 2.4 autonomy 判定（`autonomy.md` 策略）

技能执行后，依据 `autonomy.md` 策略块判断下一步：

| 条件 | 动作 |
|------|------|
| 高置信（high_conf）+ 操作可逆（reply / create_req / tag / create_cr / worktree / prestage） | 自动执行至**预发（prestage）**，完成后 release |
| 低置信（low_conf）/ 验证失败（verify_fail）/ 红线（redline） | 外化上下文，将 Task `SUSPENDED` 并发布 needs-attention 事件 |

```bash
# 自动走到预发
# → prestage 完成后记录 run_done，并释放认领
bootstrap/log.sh run_done <id> "自动部署至预发"
JARVIS_A1_IDENTITY=terraform-rd bootstrap/claim.sh release <id> <pool-project>  # Terraform
bootstrap/claim.sh release <id> <pool-project>                                # 非 Terraform

# needs-attention 路径由当前 fenced Session suspend；事件发布器幂等写回 Aone/钉钉。
```

---

## 三、定期维护（控制面收敛）

bridge scheduler 统一承担收敛，不再运行本地 reconcile 脚本：

- 服务端 reaper 按 Worker/Session heartbeat、lease 和 fence 收敛死亡执行，并把可恢复 Task 放回恢复链路。
- `claim_health` runner 最多每 5 分钟对账 `jarvis-claimed` 与控制面 Task/Session/timeline：健康 RUNNING/LEASED 不按总时长告警，未过期 wait 静默，心跳失联留 15 分钟给 lease/reaper 收敛；无 Task/终态残留/结构异常需间隔至少 5 分钟的两次确认。单次查询失败静默重试，告警不自动 release。
- `scan` runner 对账 `jarvis-done` 标签与合法完成态集合（`.claim.done_statuses` ∪ 各池 `.done_status`，含 tf_provider `待发布`），漂移时发布一次状态告警。
- `claim.sh finish` 若 done_status 未能落地，立即降级 `jarvis-done`→`jarvis-idle` 并返回失败，由当前 Task 进入 needs-attention，避免制造跳过黑洞。

---

## 四、Done — 本轮结束标准

| 结果 | 说明 |
|------|------|
| **到预发（prestage）** | 高置信可逆操作已自动完成，`run_done` 已写入 `runs/` |
| **挂起（needs-attention）** | 低置信/红线 Task 已 `SUSPENDED`，人工决策事件已发布 |

每条工作项最终落入上述两个状态之一，本轮 loop 结束。

---

## 五、仅人工步骤（正式发布 release_prod）

`autonomy.md` 策略永久停止项：

```
stop: ["release_prod"]
```

**正式发布（release_prod / release）无论任何模式，必须人工在 Aone 或 a1 CLI 中确认后才能执行。**

Jarvis 不会自动触发 release_prod。预发验收通过后，由工程师手动发布到生产环境。

---

## 六、工具链速查

| 工具 | 作用 |
|------|------|
| `bootstrap/preflight.sh` | 开局自检日级闸门：install+verify 24h 跑一次,`--force` 强制重跑 |
| `bootstrap/scan.sh` | 人工入箱扫描与审计兜底；Scheduler 不调用它 |
| `bootstrap/aone-get.sh <id>` | 取工单详情(3h 缓存,写后失效);`JARVIS_CACHE_TTL=0` 强制重取 |
| `bootstrap/cache.sh` | 通用 TTL 缓存(get/bust/fresh),落 `.my-day/cache/` |
| `bootstrap/plan.sh` | 出执行计划；supervised 退码 2 等待授权（bridge/serve 流程用） |
| `bootstrap/log.sh seen` | 去重检查 |
| `bootstrap/log.sh run_done` | 记录完成 |
| `bootstrap/wrap.sh sync/done` | Aone 回填与收尾；控制面 Terraform 主处理 run 内禁用（RD finalizer 返回 `AONE_RESULT.reply_body`，由 executor 单次落账）；独立非 executor finalizer 才 done 一次，后续重要事件走 bridge RD-only event publisher |
| `bootstrap/html-report-preview.sh upload/from-aone` | 非 Terraform 可端到端上传并按需使用 `--comment`；Terraform PD/QA 只返回本地路径，不上传、不回贴，RD finalizer 可统一上传一次但不得传 `--comment`，预览链接只进入唯一聚合回复 |
| `bootstrap/claim.sh claim <id> <project>` | 认领工作项；退码 1 = 输了跳过，退码 3 = 缺必填字段，需经 `aone-fields.sh` 挑合法值回填后重试；其它 update 失败直接上抛，不误报 lost race。认领成功还会把 Aone status 从起始态推进到该池进行中状态，best-effort 非阻断 |
| `bootstrap/aone-fields.sh missing <id>` | 列出当前为空的必填自定义字段；field-list options 为空时补查 field options API，输出合法候选，不自动选值 |
| `bootstrap/aone-fields.sh fill <id> <fieldId>=<value> …` | 回填 agent 已明确选择的字段值（重复 `--cfs`）；拒绝空值/非法参数 |
| `bootstrap/claim.sh release <id> <project>` | 释放认领（打 jarvis-idle 标签：本轮处理完，等待人或下一个 jarvis 接手；不动 Aone status） |
| `bootstrap/claim.sh finish <id> <project>` | jarvis 判断真完成（打 jarvis-done 标签 + status 改为 `pools.json` 里该池 × workitemType 的 `done_status`；`.claim.done_status` = `已发布待需求排期` 只是**全局兜底**，主流走 per-池 per-category。tf_provider(528766) 产品类需求 → **`待发布`**，不是 `已发布`——workflow 不允许 `已选择` 直跳 `已发布`。被拒时先 `bin/a1id -- project workitem field options status --project <id> --type <workitemType>` 查合法枚举再改 pools.json。**若 done_status 终究落不到合法完成态**，finish 会降级 `jarvis-done`→`jarvis-idle` 并返回失败，由 Task 进入 needs-attention） |
| `bootstrap/triage-one.sh <id> <pool> <project> "<summary>" <status>` | 单条工单收尾 bookend：claim→wrap done→release；claim 输竞争则 SKIP，退码 3 则输出合法字段候选并返回失败（不 wrap/release） |
| `bootstrap/wrap-check.sh` | Stop 闸门：会话结束时校验未完工工单是否已回填，失败则阻断 |
| `.claude/skills/aone-triage` | 单条工单全流程技能 |
| `autonomy.md` | 模式/置信度/停止项策略 |

# aone-triage 无人值守到预发 triage loop

> 单条工单怎么处理在 `.claude/skills/aone-triage`(读单/归类/查证/回复/bookend);本文件只管 **loop 编排**:实例协调、触发入口、认领竞争、autonomy 判定、收敛维护。

---

## 零、实例协调（coord.sh）

每个 triage 实例启动后先通过 `coord.sh` 扫孤儿、续跑；dispatch（Tata 委派）实例只做心跳与 checkpoint，不做 adopt。

```bash
# triage 实例开局：注册自身，扫孤儿并续跑
COORD_ID=$(bootstrap/coord.sh register triage "$$")   # 传自身 pid,coord.sh dead 才能按 kill -0 判活
for oid in $(bootstrap/coord.sh list-orphans); do
  COORD_ID=$COORD_ID bootstrap/coord.sh adopt "$oid"
done

# dispatch 实例（Tata 委派）：只心跳 + checkpoint，跳过 adopt
# nohup bootstrap/heartbeat.sh "$COORD_ID" $$ >/dev/null 2>&1 &
```

- `list-orphans`：列出任务文件中 owner_instance 已死（`coord.sh dead` 返回 0）的工单 id。
- `adopt <aone_id>`：将孤儿任务的 `owner_instance` 改为当前实例，使其进入本轮 triage 队列续跑。
- dispatch 实例（如 DingTalk bridge）不执行 adopt，避免重复接管正常分派的工单；仅保持心跳和阶段 checkpoint，供 watchdog 监控存活。

---

## 一、触发与输入

| 方式 | 说明 |
|------|------|
| bridge 定时扫池（自动派发） | bridge ScanScheduler 周期跑 `bootstrap/scan.sh --force`，diff 出**新单 + 外部更新单**(gmtModified 变化)一并投**并发 DispatchPool 起 headless jarvis**（每单一实例、并发上限 `JARVIS_DISPATCH_MAX`、软去重台账 `.my-day/bridge/dispatched.json`；`claim.sh` 仍是竞争互斥真源）。**派发判定**(`_decide`)逐单：终态 / `jarvis-done` / `jarvis-claimed` → skip；`jarvis-npe`（路由不明，人工标记）→ skip（**优先于 idle 门**：idle+npe 就算有人评论也不重派，直到人工澄清路由、摘掉标签，经 gmtModified 更新路径自然恢复派发）；`jarvis-idle` 过**人工介入门**（activity 作者判据 `_human_touched`——人工在 jarvis 上轮动作**之后**介入过才 force 重派，否则 skip 等每日 Revisit）；其余（含新单/外部更新）→ 派发。**灰度安全阀**（`_in_scope`，默认全放；`JARVIS_DISPATCH_POOLS` 池白名单 + `JARVIS_DISPATCH_CREATED_BEFORE` 创建上限可收窄）+ **运行时暂停**（`touch .my-day/bridge/pause` 停扫+派、`rm` 恢复，在跑实例不受影响）。钉钉卡片语义从「求授权」改为「播报：已自动派发 #id」。**授权前置=`JARVIS_AUTO_DISPATCH=0` 的回退模式**（新单入 pending 待钉钉「处理 #id / 全部处理」，更新单仅作「有更新」感知卡片段落）。另有 ProbeScheduler（每日 `JARVIS_PROBE_HOUR` 探测轮，跑 `loops/tf-probe.md`）与 RevisitScheduler（每日 `JARVIS_REVISIT_HOUR` 重访 `jarvis-idle` 人工门工单，**排除 `jarvis-npe`**——路由不明单不投复查）。见 `bridge/jarvis_dingtalk_bot.py`；启动入口统一 **`bridge/run.sh start`**（自动 source `bootstrap/.env`+`bridge/jarvis.env`、判定钉钉/降级模式、缺钉钉凭证干净降级不阻断扫描/派发/调度、pidfile 守护；`bridge/run.sh dry-run` 透传 `--dry-run-once` 离线看派发/跳过决策）。**Jarvis 不再主动扫**（CLAUDE.md 开局动作 #3） |
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

### 2.1 去重检查

```bash
bootstrap/log.sh seen <id>
```

- `seen` 返回 0（已处理）→ 跳过本条，继续下一条。
- `seen` 返回 1（未处理）→ 进入认领。

### 2.2 认领（并发竞争锁）

```bash
bootstrap/claim.sh claim <id> <pool-project>
```

- 退出码 0 → 认领成功，继续 triage。
- 退出码 1 → 其他实例已抢先认领，**跳过本条**，继续下一条。

### 2.3 单条 triage（技能调用）

调用 `.claude/skills/aone-triage` 技能，传入当前工作项 id。

技能完成：读取工单 → 查证（OpenAPI + Cloudspec 映射 + provider 源码）→ 回复 / 打标 / 建需求 / 建 CR。

若单条工单进入 Terraform Provider 资源开发且不走自动化生成链路，先按 `tf_provider` 池创建或复用 **terraform-alicloud** 内部研发单（项目 `528766`），**指派按 aone-triage skill `references/tf-customer-request-routing.md` 分工表路由到具体人**——即便由 jarvis 代为开发，关联单也挂具体人名下，方便其注意到；与客户主单双向关联。**指派给过载（484483）的关联单，jarvis 直接 claim 跟进解决，bookend 同时处理客户主单与关联单**（研发细节 wrap 关联单、客户主单只 wrap 关键节点，收尾两边各自 done+release）；指派给其他人的关联单不 claim，建单 + @对方等接手。研发细节、验证、PR/CI/验收信息写内部研发单，客户主单仅同步关键节点和卡点。需要 cloudspec_gap 或云产品上游协助时，把详细协作问题同步到对应依赖单。

GitHub PR/评论/推分支的身份纪律见 CLAUDE.md 工作纪律 #6（`bootstrap/github-identity.sh`，账号必须 `api-tool-agent`）。

### 2.4 autonomy 判定（`autonomy.md` 策略）

技能执行后，依据 `autonomy.md` 策略块判断下一步：

| 条件 | 动作 |
|------|------|
| 高置信（high_conf）+ 操作可逆（reply / create_req / tag / create_cr / worktree / prestage） | 自动执行至**预发（prestage）**，完成后 release |
| 低置信（low_conf）/ 验证失败（verify_fail）/ 红线（redline） | escalate 后 release |

```bash
# 自动走到预发
# → prestage 完成后记录 run_done，并释放认领
bootstrap/log.sh run_done <id> "自动部署至预发"
bootstrap/claim.sh release <id> <pool-project>

# escalate 路径
bootstrap/log.sh escalate <id> "<reason>"
bootstrap/claim.sh release <id> <pool-project>
```

---

## 三、定期维护（僵尸清扫）

定期（建议每次 loop 结束后，或按 cron 独立触发）运行：

```bash
bootstrap/reconcile.sh all      # stale + orphan + drift 三路并跑
bootstrap/reconcile.sh stale    # 只跑僵尸 claim(原 sweep.sh)
```

`reconcile.sh stale` 查找所有 `jarvis-claimed` 超过 `claim.ttl_min`（默认 45 分钟）的工作项，将其写入 `escalation/` 目录，防止僵尸认领长期占用工单。`orphan` 处理 owner_instance 已死的 task；`drift` 补漏 release。

---

## 四、Done — 本轮结束标准

| 结果 | 说明 |
|------|------|
| **到预发（prestage）** | 高置信可逆操作已自动完成，`run_done` 已写入 `runs/` |
| **入队（escalation）** | 低置信/红线条目已写入 `escalation/`，等待人工决策 |

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
| `bootstrap/scan.sh` | 入箱扫描 → JSON 工作项列表（bridge ScanScheduler 定时调；手动兜底） |
| `bootstrap/aone-get.sh <id>` | 取工单详情(3h 缓存,写后失效);`JARVIS_CACHE_TTL=0` 强制重取 |
| `bootstrap/cache.sh` | 通用 TTL 缓存(get/bust/fresh),落 `.my-day/cache/` |
| `bootstrap/plan.sh` | 出执行计划；supervised 退码 2 等待授权（bridge/serve 流程用） |
| `bootstrap/log.sh seen` | 去重检查 |
| `bootstrap/log.sh run_done` | 记录完成 |
| `bootstrap/wrap.sh sync/done` | 进展回填 Aone（唯一真源）+收尾审计；多行正文用 `--summary-stdin`/`--summary-file` |
| `bootstrap/html-report-preview.sh upload/from-aone` | 将 HTML/zip/Aone 附件报告上传到 AutomationAgent，返回在线预览链接；`--comment` 可直接回贴 Aone |
| `bootstrap/log.sh escalate` | 记录上报 |
| `bootstrap/claim.sh claim <id> <project>` | 认领工作项（输赢竞争锁）；退码 1 = 输了跳过 |
| `bootstrap/claim.sh release <id> <project>` | 释放认领（打 jarvis-idle 标签：本轮处理完，等待人或下一个 jarvis 接手；不动 Aone status） |
| `bootstrap/claim.sh finish <id> <project>` | jarvis 判断真完成（打 jarvis-done 标签 + status 改为 .claim.done_status，默认"已发布待需求排期"） |
| `bootstrap/triage-one.sh <id>` | 单条工单 bookend 编排：claim→子代理→wrap done（**status 必填**）→release |
| `bootstrap/wrap-check.sh` | Stop 闸门：会话结束时校验未完工工单是否已回填，失败则阻断 |
| `bootstrap/reconcile.sh [stale\|orphan\|drift\|all]` | 收敛族入口(P1.c 合并 sweep+watchdog+原 reconcile):stale=超时 claim→escalate;orphan=owner dead→escalate;drift=台账 vs Aone 对账;all=顺跑三者(默认) |
| `.claude/skills/aone-triage` | 单条工单全流程技能 |
| `autonomy.md` | 模式/置信度/停止项策略 |

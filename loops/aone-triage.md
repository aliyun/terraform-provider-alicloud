# aone-triage 无人值守到预发 triage loop

> Plan #2 Task 4 — 将各组件串成完整的 unattended-to-prestage 闭环。

---

## 零、实例协调（coord.sh）

每个 triage 实例启动后先通过 `coord.sh` 扫孤儿、续跑；dispatch（Tata 委派）实例只做心跳与 checkpoint，不做 adopt。

```bash
# triage 实例开局：注册自身，扫孤儿并续跑
COORD_ID=$(bootstrap/coord.sh register triage)
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

## 一、触发

| 方式 | 说明 |
|------|------|
| 定时触发 | cron / 调度器定期调用本 runbook |
| 手动触发 | 用户在 Claude Code 会话中执行 `/aone-triage` 或直接发指令 |

---

## 二、输入扫描

```bash
bootstrap/scan.sh
```

输出：标准 JSON 数组，格式 `[{id, title, type, status}]`，代表当前分配给本账号的所有工作项。

- 若 scan.sh 返回非零退出码，终止本轮，输出错误信息，等待下次触发。
- 若结果为空数组，本轮结束，不做任何写操作。

---

## 三、出执行计划（supervised 授权关口）

```bash
bootstrap/plan.sh
```

plan.sh 分析扫描结果，输出本轮拟执行的逐项计划。

- **supervised 模式（默认）**：plan.sh 退出码 2 表示计划等待授权；输出计划后**暂停**，逐项向用户请求授权。
  - 用户回复"可行 / 允许 / 授权"后，该条目进入下一阶段。
  - 未授权条目本轮跳过。
- **unattended 模式**：需用户显式以 `--mode unattended` 开启，计划直接进入执行阶段，无需逐条等待。

---

### 阶段3说明：分诊与授权的流程

流程是：
1. **scan（全集，带 priority/tag）** → 扫描所有工作项，带优先级和标签信息
2. **plan（去重，不重扫）** → 生成计划，仅去重不重新扫描
3. **Claude 按优先级/标题/标签/状态排序+折叠噪声** → 由 Claude 进行分诊和排序，推荐 N 条工作项
4. **逐条授权** → 用户对推荐项逐条确认

**关键点：分诊是 Claude 做，不是脚本扫表；plan 只去重不重扫。**

保留 **supervised 门**与 **release_prod 硬门**两道关口。

### 产物落点纪律

- 计划/审计走 `runs/`（scan→plan 的 stdout，run_done/escalate 记录）。
- 临时数据用 `.my-day/`（gitignored）或 `mktemp`，**禁止**往仓库根甩 `scratch_*`/`*.tmp`。
- 开发不在 master：worktree 切分支 → MR → 人工合并；`.gitignore` 兜底脏文件。

---

## 四、逐项执行

对每条**已授权**的工作项，按以下顺序处理：

### 4.1 去重检查

```bash
bootstrap/log.sh seen <id>
```

- `seen` 返回 0（已处理）→ 跳过本条，继续下一条。
- `seen` 返回 1（未处理）→ 进入认领。

### 4.1.5 认领（并发竞争锁）

```bash
bootstrap/claim.sh claim <id> <pool-project>
```

- 退出码 0 → 认领成功，继续 triage。
- 退出码 1 → 其他实例已抢先认领，**跳过本条**，继续下一条。

### 4.2 单条 triage（技能调用）

调用 `.claude/skills/aone-triage` 技能，传入当前工作项 id。

技能完成：读取工单 → 查证（OpenAPI + Cloudspec 映射 + provider 源码）→ 回复 / 打标 / 建需求 / 建 CR。

若单条工单进入 Terraform Provider 资源开发且不走自动化生成链路，先按 `tf_provider` 池创建或复用 **terraform-alicloud** 内部研发单（项目 `528766`，指派 `WORKER_1782379562571`），与客户主单双向关联；研发细节、验证、PR/CI/验收信息写内部研发单，客户主单仅同步关键节点和卡点。需要 cloudspec_gap 或云产品上游协助时，把详细协作问题同步到对应依赖单。

若单条工单会产生 GitHub PR/评论/推分支，写操作前必须执行 `bootstrap/github-identity.sh check`；`gh` 写操作通过 `bootstrap/github-identity.sh gh ...` 执行，推分支通过 `bootstrap/github-identity.sh push <owner/repo> <local-ref> <remote-ref>` 执行；`JARVIS_GITHUB_TOKEN` 登录名必须是 `api-tool-agent`，PR head 使用 `api-tool-agent:<branch>`。

### 4.3 autonomy 判定（`autonomy.md` 策略）

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

## 四点五、定期维护（僵尸清扫）

定期（建议每次 loop 结束后，或按 cron 独立触发）运行：

```bash
bootstrap/sweep.sh
```

sweep.sh 查找所有 `jarvis-claimed` 超过 `claim.ttl_min`（默认 45 分钟）的工作项，将其写入 `escalation/` 目录，防止僵尸认领长期占用工单。

---

## 五、Done — 本轮结束标准

| 结果 | 说明 |
|------|------|
| **到预发（prestage）** | 高置信可逆操作已自动完成，`run_done` 已写入 `runs/` |
| **入队（escalation）** | 低置信/红线条目已写入 `escalation/`，等待人工决策 |

每条工作项最终落入上述两个状态之一，本轮 loop 结束。

---

## 六、仅人工步骤（正式发布 release_prod）

`autonomy.md` 策略永久停止项：

```
stop: ["release_prod"]
```

**正式发布（release_prod / release）无论任何模式，必须人工在 Aone 或 a1 CLI 中确认后才能执行。**

Jarvis 不会自动触发 release_prod。预发验收通过后，由工程师手动发布到生产环境。

---

## 七、工具链速查

| 工具 | 作用 |
|------|------|
| `bootstrap/preflight.sh` | 开局自检日级闸门：install+verify 24h 跑一次,`--force` 强制重跑 |
| `bootstrap/scan.sh` | 入箱扫描 → JSON 工作项列表 |
| `bootstrap/aone-get.sh <id>` | 取工单详情(3h 缓存,写后失效);`JARVIS_CACHE_TTL=0` 强制重取 |
| `bootstrap/cache.sh` | 通用 TTL 缓存(get/bust/fresh),落 `.my-day/cache/` |
| `bootstrap/plan.sh` | 出执行计划；supervised 退码 2 等待授权 |
| `bootstrap/log.sh seen` | 去重检查 |
| `bootstrap/log.sh run_done` | 记录完成 |
| `bootstrap/wrap.sh sync/done` | 进展回填 Aone（唯一真源）+收尾审计 |
| `bootstrap/log.sh escalate` | 记录上报 |
| `bootstrap/claim.sh claim <id> <project>` | 认领工作项（输赢竞争锁）；退码 1 = 输了跳过 |
| `bootstrap/claim.sh release <id> <project>` | 释放认领（打 jarvis-idle 标签：本轮处理完，等待人或下一个 jarvis 接手；不动 Aone status） |
| `bootstrap/claim.sh finish <id> <project>` | jarvis 判断真完成（打 jarvis-done 标签 + status 改为 .claim.done_status，默认"已发布待需求排期"） |
| `bootstrap/sweep.sh` | 清扫超时认领（>45min）→ 写入 `escalation/` |
| `bootstrap/triage-one.sh <id>` | 单条工单 bookend 编排：claim→子代理→wrap done（**status 必填**）→release |
| `bootstrap/wrap-check.sh` | Stop 闸门：会话结束时校验未完工工单是否已回填，失败则阻断 |
| `bootstrap/reconcile.sh` | 漂移对账：比对 runs/ 台账与 Aone 在线状态，输出差异清单 |
| `.claude/skills/aone-triage` | 单条工单全流程技能 |
| `autonomy.md` | 模式/置信度/停止项策略 |

# Poison source-entry guards：让不可读的候选条目退出热路径

**日期**：2026-08-04
**分支**：`worktree-poison-source-guards`
**触发**：macmini（`AgenticTools-Macmini.local`）任务处理变慢排查，发现调度线程被注定失败的 Aone 读取长期占用。

---

## 一、问题

`list_source_status_candidates` 返回的候选条目里，有 5 条**永远读不成功**，但每个扫描轮次都会被重新读一遍，无限期。

### 观测数据（2026-08-04，截至 17:38，非完整一天）

| 失败签名 | 次数 |
|---|---|
| `workitem get failed (403)` | 923 |
| `workitem list failed (403)` | 654 |
| `workitem get failed (404)` | 306 |
| **合计** | **1883** |

实测单次 `a1 project workitem get` 约 **2.9s**（3 次采样：2.99 / 2.85 / 2.76）→ 约 **91 分钟 ≈ 1.5 小时**的调度线程时间，花在不可能成功的调用上。`owner-health`(4470 行) 与 `aone-workitem-ownership`(3080 行) 因此成为当日日志量前二。

### 5 条 poison 条目

```
# 404 组（aoneId 不是合法工单号，真实为 8 位）
taskId=1130  sourceProjectKey=1086837  aoneId=9001  sourceStatus=待处理
taskId=1135  sourceProjectKey=1086837  aoneId=7710  sourceStatus=待处理
taskId=4451  sourceProjectKey=1086837  aoneId=779   sourceStatus=待处理

# 403 组（项目未登记，jarvis 无成员权限）
taskId=2091  sourceProjectKey=709564   aoneId=84574563  sourceStatus=待处理
taskId=3813  sourceProjectKey=1073966  aoneId=68053511  sourceStatus=已发布待需求方验收
```

三条 404 的 aoneId 各被 point-read **正好 167 次**，等于当日扫描轮数 —— 每轮一次，永不收敛。

### 两个根因

**根因 1：404 与临时失败被折叠成同一种结果。**
`bridge/scheduler/runners/scan.py::_point_read_source_status`（`:261` 起，`@staticmethod`）：

```python
if result.returncode != 0:
    log.warning("ScanRunner: source status point-read #%s rc=%d: %s", ...)
    return task, None      # ← 404（永久）与超时/抖动（临时）无法区分
```

返回 `None` 即不上报 source status，条目于是永远停在非终态，下轮再读。方法内的 `aone_id.isdigit()` 门只校验是否为数字、不校验位数，`779` 能通过。

注意该方法的 docstring 明确写着 *"Failures return `None` and leave the persisted status untouched."* —— 对**临时**失败而言这是有意为之，本设计必须保留；要改的只是把**永久** 404 从这个笼子里拿出来。

**根因 2：`aone_workitem_ownership.py` 缺少两个 sibling 都有的护栏。**

`list_source_status_candidates` 有三个消费者：

| 消费者 | terminal skip | 池过滤 |
|---|---|---|
| `scan.py` | 有（3 处） | 有（`_read_pools` + `JARVIS_DISPATCH_POOLS`） |
| `owner_health.py` | 有（`:239`，commit `eb04a2c`） | — |
| **`aone_workitem_ownership.py`** | **无** | **无** |

`_list_candidates`（`:450`）直接翻页、不做任何过滤，所以未登记项目和已终态条目照样进入读取。taskId=3813 的 `sourceStatus` 已是终态（`已发布待需求方验收` 在 `TERMINAL_STATUSES` 内），owner_health 会跳过它，但 ownership 仍每轮 batch list 其项目 327 次。

---

## 二、范围

### 本 spec 包含

- **变更 a**：ownership runner 补池过滤 + terminal skip，与两个 sibling 对齐
- **变更 b**：404 永久失败熔断 —— 发一次钉钉通知 + 自动盖终态
- **变更 c**：项目级 403 后不再回退逐条 detail

### 明确不包含（deferred）

- **403 指数退避 + 一次性告警**：变更 a 的池过滤覆盖了全部现存 403 案例（两条都在未登记项目）。剩余场景是「已登记池被撤权」，当前不存在，按 YAGNI 不做；真出现时变更 c 已挡掉大部分成本。
- **3 个 `RECOVERY_REQUIRED` 卡单**（task 624/2778/3312，工单 #82716072/#84856234/#84873148）：根因是 payload `policyRevision` 为 `terraform-rd-single-writer-v4`、当前为 `v6`，与本 spec 的读取浪费无关，另开跟踪。
- **`JARVIS_DISPATCH_MAX` 10→5 未生效**：运维动作（bridge 于 13:34 启动时带 `=10`，配置 16:32 才改为 `5` 且未重启），无需改码。

---

## 三、设计

### 变更 a：ownership runner 护栏对齐

**位置**：`bridge/scheduler/runners/aone_workitem_ownership.py::_list_candidates`（`:450`）

翻页取得候选后，过滤两类条目：

1. `sourceProjectKey` 不在已登记项目集合内 —— 复用 `AoneQueryMixin._read_pools()`（`bridge/aone_tasks.py:715` 类、`:722` 方法，`@staticmethod`，返回 `[(key, project, exclude_status[], pr_merged_status|None)]`），取其中的 `project` 构成集合，不另写读取逻辑。该文件目前**未从 `bridge.aone_tasks` 导入任何东西**，需新增 import（`TERMINAL_STATUSES` 同样来自该模块）。
2. `sourceStatus` ∈ `TERMINAL_STATUSES`

**fail-open 是硬要求**：`_read_pools()` 返回空时**不得**过滤（否则一次配置读失败会静默清空全部候选，把浪费问题变成停摆问题）。照 `scan.py:246` 的 `if not pools: return` 同形处理。

**日志**：每轮只出一行聚合计数，不逐条打。`eb04a2c` 的注释已经量化过这个取舍 —— 逐条打的日志量比它节省的告警量更大。

**效果**：taskId=2091、3813 离开热路径；3813 同时被 terminal 门命中。

### 变更 b：404 熔断（notify once + 自动盖终态）

**位置**：`bridge/scheduler/runners/scan.py::_point_read_source_status`（`:261`）+ `_resolve_source_statuses` + 其调用循环（`:449-470`）

**职责划分（因为 `_point_read_source_status` 是 `@staticmethod`，拿不到实例状态）**：
- `_point_read_source_status`：只做分类，把「永久 404」作为一个可区分的标记返回，不碰状态、不发通知
- `_resolve_source_statuses`：把标记向上透传
- 调用循环（`:449` 起）：持有 runner 实例，负责读写 episode 状态文件、判护栏、发通知，并复用循环内**已存在**的 `client.update_source_status(task_id, aone_id, source_status, request_id="source-status:%s:%s" % (task_id, digest))` 调用 —— 把 `source_status` 传成 `"Invalid"` 即可，`request_id` 的 digest 机制原样复用，无需新增调用点

现状该循环遇到 falsy 的 `source_status` 会直接 `continue`，这正是 poison 条目每轮逃逸的地方；改动即在此处插入护栏判定。

选择这个落点的理由：`update_source_status` **全仓唯一调用者是 `scan.py:465`**，而 404 恰好在同一组件被观测到。改动因此是「把一个错误的错误分类修对」，不新增控制面调用点、不新增 scheduler job。

（备选方案 `owner_health` 被否：它语义上是「告警一次 + 补救」那个 job，但它读控制面 `get_task_by_aone`、**不在 Aone 404 的观测路径上**，塞进去等于给它加一条本不属于它的 Aone 读取职责。）

**分类**：解析 a1 stderr，命中 `workitem get failed (404)` 判为永久失败，与临时失败分开返回。

**三个护栏，全部成立才盖终态**：

1. `len(aone_id) < 8` —— aoneId 位数不合法。779/7710/9001 全中；真实 8 位号不会误伤。
2. 同一条目跨 **≥2 轮**观测，且首末观测相隔 **≥24h** —— 单日抖动不足以触发。
3. 本轮**同项目**至少有 1 次 Aone 读取成功 —— 排除 Aone 整体故障期间误判。

**动作**：`update_source_status(task_id, aone_id, "Invalid")`。`Invalid` 本就在 `TERMINAL_STATUSES` 内（来自 `config/pools.json` 的 `claim.done_statuses`），因此盖完之后**现有**的 terminal skip 会在下一轮自动接管 —— 不需要新写任何跳过逻辑。

选 `Invalid` 而非 `已取消`：语义是「这不是一个真实工单号」，而不是「工作被取消」。

**通知**：`_dingtalk_event_enqueue(ticket=aone_id, project=sourceProjectKey, event_key=f"source-poison:not-found:{task_id}", ...)` + `_dingtalk_event_flush()`。ledger（`.my-day/bridge/dingtalk-event-ledger.json`）按 key 去重，天然只发一次。正文列出 taskId / aoneId / project / 盖入的状态。

**跨轮状态**：新文件 `.my-day/bridge/source-poison-health.json`，照 `owner-unavailable-health.json` 的 episodes 形状（`{taskId: {aoneId, firstSeenAt, lastSeenAt, count, lastAlertAt}}`）。不复用 `scan-snapshot.json` —— 那个文件的语义是 discovery lag，混入熔断状态会让两者都说不清。

**为什么不是删除**：`cleanup_legacy_kind_tasks` 需要 admin token，bridge 的 token 没有（实测 `HTTP 401: Jarvis admin endpoints require the admin token`），且它只针对 deprecated task_type。盖终态反而更好：可逆、可审计。

### 变更 c：项目级 403 不再回退逐条

**位置**：`bridge/scheduler/runners/aone_workitem_ownership.py:869-879`

现状：`_fetch_project_batch` 抛异常 → 整个 batch 的条目全部追加进 `detail_reads` 逐条重读。当失败原因是**项目级权限**（403 `您不是项目成员`）时，逐条 detail **必然同样 403** —— 这层回退是纯粹的逻辑错误。

改法：对异常分类，命中项目权限类 403 时不追加 `detail_reads`，只出一行日志。

**分类判据要写明确**（避免实现时各自解释）：`_fetch_project_batch` 抛出的异常字符串形如
`a1 list failed rc=1: Error: workitem list failed (403): 您不是项目成员，没有项目权限，因此不能访问该项目。请联系项目管理员：…`
判据 = 异常文本同时含 `workitem list failed (403)`。仅凭 `403` 子串不够（工单号里可能出现 `403`）；也不用中文文案做主判据（文案会变），中文只作日志说明。

**不动** `:897-903` 那条「batch omitted 单条 → 回退 detail」的分支 —— 那是合理的（batch 结果缺行，单条确实可能读到）。

---

## 四、错误处理

| 风险 | 处理 |
|---|---|
| `_read_pools()` 读空 → 清空全部候选 | fail-open：空则不过滤（变更 a 的首要风险点） |
| Aone 整体故障期间大批 404 → 误盖终态 | 护栏 3（同项目须有成功读取）+ 护栏 2（≥24h 跨轮） |
| 真实 8 位工单被删除 → 持续 404 | 护栏 1 拦住，不自动盖；保持现状每轮重读（不退化） |
| `update_source_status` 失败 | 不写 episode 的 alert 标记，下轮重试；不宣称已摘除 |
| 钉钉发送失败 | ledger 保持 pending，下轮 `_flush` 补偿（既有语义） |

护栏 1 有意保守：它把「工单真实存在过但被删除」这种情况排除在自动摘除之外。这类条目会继续每轮重读 —— 相比误摘真实工作，可接受。

---

## 五、测试

落 `bridge/scheduler/tests/`（与 `test_daily_probe.py` 同形）。

**变更 a**
- 未登记项目的候选被过滤
- 终态候选被过滤
- `_read_pools()` 返回空时**不**过滤（fail-open 回归）
- 聚合日志只出一行

**变更 b**
- 三个护栏各自单独不成立时**均不得**盖终态（3 个独立用例）
- 三者同时成立时调用 `update_source_status(..., "Invalid")` 恰好一次
- 通知按 `event_key` 幂等，重复轮次不重发
- 临时失败（超时、非 404 的 rc≠0）不写 episode、不盖终态
- `update_source_status` 抛错时不标记已告警

**变更 c**
- 项目级 403 异常 → `detail_reads` 不增长
- 其它异常 → 保持现有逐条回退
- batch omitted 单条 → 仍回退 detail（不受影响）

---

## 六、上线

三个变更互不依赖，可同一 PR。合并后需重启 bridge 生效 —— 与待决的 `JARVIS_DISPATCH_MAX` 重启可合并为同一次操作。

**验收信号**（重启后一个完整日）：`workitem get failed (403|404)` 与 `workitem list failed (403)` 三个签名的日计数应从 1883 降到接近 0；`aone-workitem-ownership` 日志行数显著下降；`.my-day/bridge/source-poison-health.json` 出现 3 条 episode 且各带一次 `lastAlertAt`。

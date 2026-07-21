# Aone external-write operation receipts

Ticket 84386065 需求一：把 `wrap.sh sync/done` 的 comment/status 写与 `claim.sh
release/finish` 的 tag 写接入控制面 operation receipt（begin/ack/fail/reconcile）
两阶段状态机。此前只有 `claim.sh claim` 的 tag 写有回执。

适用范围：当前会话是 database-fenced worker（交互 Claude/Codex，或预注册的
headless worker）且持有目标工单（`has-current` 为真）。其余上下文（无 worker
state、非当前工单——例如 bookend 同时回填的客户主单）保持原有裸写行为不变。

## 数据面契约（AutomationAgent，已在 master）

- `POST operations/begin`：`operationKey` 幂等（同 task+generation 内唯一）；已
  ACKED 返回 `proceed=false`；状态 UNKNOWN 时 409 `UNKNOWN operation must be
  reconciled`。
- `POST operations/ack`：携 workerKey+fenceToken+externalRef；旧 fence 412。
- `POST operations/fail`：`unknown=true` 冻结为 UNKNOWN（副作用不可证），否则按
  retryAllowed 进 RETRY_WAIT/DEAD。
- `GET operations/{id}` / `GET operations/by-key`：machine-token 保护的 point-read，
  返回 operation、Task、Session/fence、Worker 与最小化 readbackSpec；不返回
  requestPayload 或凭证。
- `POST operations/not-started`：仅当前 source worker+process+fence 可把普通
  SENDING 标为 FAILED_NOT_STARTED。它证明控制面已创建 intent、但外部副作用尚未
  开始；同 key 可再次 begin。recovery lease release 只表示 reconciler 没拿到权威
  readback，不能冒充 abort/not-started。
- `POST operations/reconcile`：`{operationId, workerKey, found, externalRef,
  retryAllowed, retryAfterSeconds}`。`found=true`→ACKED（副作用已在，恰好一次收
  敛）；`found=false`+retryAllowed→RETRY_WAIT（之后可重新 begin→proceed 重发）。
- `JarvisOperationType`：`AONE_CLAIM / AONE_COMMENT / AONE_STATUS / AONE_RELEASE
  / OTHER`。

## Worker CLI 扩展（bootstrap/jarvis-interactive-worker.py）

沿用单槽 `pendingOperation` 本地意图记录（persist intent before begin），在
claim 专用循环之外泛化出一组"当前 assignment 上的外部写回执"命令：

- `operation-begin <aone_id> <kind> <material> [--payload-json <json>]
  [--not-required] [--replay-safe]`
  - `kind ∈ {comment, status, release-tag, finish-tag}` → operationType 映射
    `AONE_COMMENT / AONE_STATUS / AONE_RELEASE / AONE_RELEASE`；
  - 要求 `current` assignment 匹配 aone_id；不创建/恢复 session（区别于
    prepare-claim）；已有**不同 key** 的 pendingOperation → conflict 退出；
  - `operationKey = "<kind>:%s:%s:%s:%s" % (task_id, generation, attempt12,
    digest12)`，attempt12 = sha256("<sessionId>:<fenceToken>") 前 12 位（取自
    current assignment），digest12 = sha256(material) 前 12 位。comment 的
    material=最终正文 sha256，status 的 material=目标状态值，tag 的 material=
    目标标签集。attempt 分量把回执按 lease attempt 隔离：同 generation 任务复活
    后的新 Session/fence 派生新 key，不会命中上一 attempt 的 ACKED 回执而跳过
    本轮该写的标签/状态；同 attempt 内重试（fence 不变）仍复用同 key 去重；
  - 决策输出 JSON `{accepted, proceed, needsReadback, operationId,
    operationStatus}`：
    - 服务端 ACKED → `proceed=false`（恰好一次：跳过 Aone 写）；
    - 服务端 proceed → `proceed=true`；
    - SENDING 且本地有同 key 意图记录：`--replay-safe`（tag/status 等幂等写）
      → `proceed=true` 直接重放；非 replay-safe（comment 不幂等）→
      `needsReadback=true`，调用方必须 readback 后走 reconcile；
    - SENDING 且无本地意图 → fail-closed conflict（丢失 begin 响应，无法自证）；
    - 本地意图记录状态为 UNKNOWN（上轮 abort --unknown 留下）→ 不调 begin，
      直接 `needsReadback=true` 返回存量 operationId。
  - begin HTTP 结果分类：2xx 才接受服务端 receipt；400/401/403 是明确拒绝，操作
    确定未开始，清本地 BEGINNING intent；409 保留 key/幂等线索供 readback；412
    作为 lost fence 处理；timeout、连接中断和无法判断的 5xx 冻结为 UNKNOWN。
- `operation-ack <aone_id> <external_ref>`：既有实现即通用（ACK 当前
  pendingOperation 并清槽）。
- `operation-abort <aone_id> <message> [--unknown]`：仅终结**回执**，不 fail
  session、不写 pendingClaim（区别于 claim 失败路径的 operation-fail）：
  - 定性未开始：not-started 后清本地槽（下轮重新 begin 重试）；
  - `--unknown`：fail_operation(unknown=true) 且**保留**本地槽（记 status=
    UNKNOWN），供下轮 `operation-begin` 短路到 readback/reconcile。
- `operation-reconcile <aone_id> --found <external_ref> | --not-found
  [--no-retry]`：调 reconcile 端点；found→清槽返回 `{proceed:false}`；
  not-found+retry→保留 key（status=RETRY_WAIT）返回 `{proceed:false,
  retryScheduled:true}`，调用方随后重新 `operation-begin` 拿 proceed。
  若 UNKNOWN 发生在 begin 响应前、operationId 为空，CLI 先按
  taskId+generation+operationKey point-read：找到后补齐 operationId 再 readback；
  404 则权威证明 begin 未提交、清本地 intent；网络错误仍保持 UNKNOWN。

会话心跳：mid-task 回执不动 `heartbeatEnabled`（session 本来就在跑）；这与
claim 期间"ACK 前不许续租"的既有语义不同且刻意为之。

## wrap.sh 协议

`_receipted_comment <id> <final_text>`（sync 与 write_done 的 comment 都走）：

1. 非 fenced 上下文 → 原样 `a1 comment create`（现状不变）。
2. fenced：`digest=sha256(final_text 规范化)`（规范化=去 \r、去尾空白）；
   `operation-begin <id> comment <digest> --payload-json {digest,preview}`。
   begin 调用失败（控制面不可用/conflict）→ **fail-closed 不写 Aone**，退出非
   零（sync 同样阻断——fenced 会话的写必须有回执）。
3. `proceed=false` 且 ACKED → 跳过发送（上轮已恰好一次落地），继续后续步骤。
4. `needsReadback` → `a1 comment list -f json` 取评论流，按规范化 sha256 匹配。
   readback 三态：
   - **found**（拉取+解析成功且命中）→ `operation-reconcile --found
     aone:<id>:comment:<cid>` → 跳过发送；
   - **definitely-not-found**（拉取成功、JSON 完整解析为列表、确认无匹配）→
     `operation-reconcile --not-found` → 重新 `operation-begin` →
     `proceed=true` 进 5；
   - **unavailable**（a1 非零退出 / JSON 解析失败 / 命中但无 comment id）→
     **不得**当 not-found：槽仍是 SENDING 则 `operation-abort --unknown` 冻结，
     已是 UNKNOWN 则保持原状不再 abort；输出提示、非零退出等重试。首发成功仅
     ACK 丢失时，一次 503/超时/坏 JSON 不会被错判成「副作用不存在」而重复评论。
5. `proceed=true` → `a1 comment create -f json`（或解析 `ID:` 行）拿 comment
   id：
   - 成功 → `operation-ack <id> aone:<id>:comment:<cid>`；ack 失败（CLI 内部已
     自动冻结 UNKNOWN）→ 退出非零，提示重跑 wrap 收敛；
   - a1 退出非零 → `operation-abort <id> "<msg>"`（定性失败，可重试）→ 退出非
     零；
   - a1 超时/结果不明 → `operation-abort <id> "<msg>" --unknown` → 退出非零。

`write_done` 的 status 写同理：`operation-begin <id> status <目标状态>
--replay-safe` → `a1 update --status` → ack `aone:<id>:status:<值>`；失败
abort；重放安全（幂等 set），readback 不需要（SENDING+本地意图直接重放）。

comment 与 status 是两个串行回执（单槽约束）；comment 先行，status 随后。

## claim.sh release / finish

- `release`（fenced 分支）：has-current 通过后：
  `operation-begin <id> release-tag idle --replay-safe` → `_update_tags_merged`
  （idle ∪ −claimed）→ 成功 `operation-ack <id>
  aone:<proj>:<id>:tag:jarvis-idle` → 既有 `suspend`。tag 写失败 →
  `operation-abort`；readback 不确定（point-read 不可达）→ `--unknown`。
  begin 失败 → fail-closed 退出 2（不碰 Aone、session 保持）。
- `finish`（fenced 分支）：同 pattern，kind=finish-tag、material=done；status
  写沿用 finish 现有降级逻辑（status 落地失败已有 jarvis-idle 降级+escalate 兜
  底，不再叠一层回执；降级 tag 写维持 best-effort，注释说明）。收尾
  complete/suspend 前回执必须已清（`transition` 既有 pendingOperation 闸门天然
  强制顺序）。

## 不变量

- 单槽：任一时刻至多一个在途外部写回执；违反 → conflict 退出。
- fenced 会话内：没有回执（begin proceed / reconcile found 之外）不允许产生任
  何 Aone 可见效果；控制面不可用一律 fail-closed。
- UNKNOWN 永不盲重放：只能经 readback（comment 按内容 digest 匹配、tag/status
  按 point-read）→ reconcile 收敛。
- 旧 fence 的 ack/fail/reconcile 由服务端 412/409 拒绝，本地槽随 conflict 清
  理交恢复路径。

## 恢复命令白名单（PreToolUse guard）

`_local_tool_block_reason` 在 pendingOperation 存在时阻断可能产生副作用的工具调用。对
claim 回执这是自洽的（恢复入口 = 标准 `claim.sh claim`，本就放行）；但外部写
回执 abort --unknown / reconcile not-found 后若进程丢失，收敛入口本身就是
wrap.sh / claim.sh 的 Bash 调用——不开白名单会把会话锁死。故 guard
（`_exact_receipt_recovery` + `_external_receipt_recovery_target`）在且仅在以
下条件全部成立时放行**精确恢复命令**：

- pendingOperation 是外部写回执（`kind ∈ comment/status/release-tag/
  finish-tag`；AONE_CLAIM 槽无 kind 字段，不受影响，仍只放行 claim.sh claim）；
- 槽状态 ∈ {UNKNOWN, RETRY_WAIT}（BEGINNING/SENDING 崩溃残留不放行——先走
  SessionStart 恢复路径）；
- 槽的 aoneId == current assignment 的 aoneId，且 worker 未 stopped、非
  subagent 发起；
- 命令是 standalone 单命令（无换行/反引号/`$(`/未引号 `;&|<>` 组合），且形如：
  - `/bin/bash <本仓绝对路径>/bootstrap/wrap.sh sync|done|done-no-status
    <同一 aone_id> …`（命令头/脚本绝对路径/子命令/aone id 四段精确，尾参自由
    ——正文是数据不是 shell）；
  - `/bin/bash <本仓绝对路径>/bootstrap/claim.sh release|finish <同一 aone_id>
    <同一 project_id>`（四段全精确，finish 不放行 status override 第五参）。

组合命令、env 前缀、相对路径、lookalike 仓路径、非 `/bin/bash`、别的 aone id
一律维持阻断。放行只解锁「能重跑收敛命令」；实际收敛仍由 worker CLI 的
begin/readback/reconcile fail-closed 把关。

wrapper 恢复命令按 kind 匹配（comment/status→wrap，release-tag→release，
finish-tag→finish）；同 Aone 的精确 `operation-abort` / `operation-reconcile`
命令也可通过恢复门，错误 kind 的写命令继续阻断。

无论本地槽或 Task 处于何种恢复/终态，以下精确、参数级校验的只读命令始终放行：
Worker/READY/Task/Session timeline、operation point-read、当前本地 permit 状态、
runtime config 来源诊断。诊断只显示 token configured/missing，绝不显示 token 值。
`FAILED_FINAL/CANCELED/SUCCEEDED` 只给出终态诊断，不再循环提示 claim；
`SUSPENDED` 等待或 wake，`RECOVERY_REQUIRED` 按 operation kind readback/reconcile，
只有 `READY` 才提示 claim。

## 机器级 runtime config

bridge、claim、wrap、interactive worker、task client、status 与 hooks 统一经
`bootstrap/runtime-config.sh` 加载控制面运行时配置。优先级为：非空进程环境 >
`JARVIS_RUNTIME_ENV` > `${XDG_CONFIG_HOME:-~/.config}/jarvis/runtime.env` > git
common dir 主 checkout 的 legacy `bootstrap/.env` / `bridge/jarvis.env`。新机器级或
显式 secret 文件要求 0600（或更严格）；凭证不复制进 worktree、tracked 文件或日志。
`runtime-config.sh diagnose` 只打印生效来源、base URL 与 token 是否已配置。

## 重启恢复：mid-task 回执孤儿化（orphanOperations）

replacement incarnation（进程重启/换代）构造 recovery pendingClaim 时，只有
**无 kind 的 AONE_CLAIM 槽**才把 operationKey/receiptUnknown 继承进
pendingClaim（下次 prepare-claim 复用同 key 幂等收敛）。mid-task 回执槽
（kind ∈ comment/status/release-tag/finish-tag）**绝不**继承：其 operationKey
属于不同 operationType/payload，复用为 AONE_CLAIM key 会被服务端按
「operationKey was reused with a different request」拒绝，任务陷入恢复循环。

这类槽改写入 `state["orphanOperations"]`（`{operationId, operationKey, kind,
aoneId, status, orphanedAt}` 记录列表），跨 incarnation 持久保留、重启不丢，
但它只用于本机诊断，不是恢复候选真源。scheduler 角色的
`ExternalOperationRecoveryScheduler` 从控制面
`/operations/external-recovery-candidates` 分页取 current-generation 候选，按
operationId 获取 token lease，续租后只读核验 Aone comment/status/tag；只有完整读回
能确定 found/not-found 才 reconcile，网络错误/坏 JSON/歧义一律 release fail-closed。
pendingClaim 不带 key →
prepare-claim 派生全新 aone-claim key。orphan 记录不参与 PreToolUse guard
放行判定（pendingOperation 槽已清）。prepare-claim 侧另有兜底：saved
operationKey 非 `aone-claim:` 前缀一律弃用改派新 key（防 pre-fix 存量 state）。

## 已知边界

- **本地 orphan 不驱动收敛**：replacement 仍保留 operationId/kind 供排障，但 bridge
  只信控制面候选 API。这样迁机、重装或本地 ledger 丢失不会漏恢复，也不会因陈旧
  orphan 误收敛已经 ACKED/换代的 operation。
- **comment 跨 attempt 相同正文会各发一次**：operationKey 含 attempt 分量后，
  exactly-once 的粒度是「每个 lease attempt」——同 generation 复活的新
  Session/fence 若生成字节级相同的评论正文，会再发一条。这是有意语义
  （per-attempt exactly-once），换取 release/finish/status 不被旧 attempt 的
  ACKED 回执错误去重。
- **comment 正文可重放性依赖 code_footer 稳定**：wrap.sh 的 claim prefix 已改
  为回执收敛成功后才消费（peek/clear），但页脚含 commit sha + dirty 位——失败
  重跑之间开发库若有新提交，digest/operationKey 漂移 → 与遗留 UNKNOWN 槽
  conflict，仍需人工清理。残留风险，暂接受。

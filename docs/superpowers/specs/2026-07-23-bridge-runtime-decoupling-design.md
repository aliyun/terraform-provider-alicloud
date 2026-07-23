# Bridge Runtime 去耦设计：Scheduler、执行器与 DingTalk Adapter

## 1. 结论与目标

当前 `SchedulerService` 的时钟、slot、admission 与 restart 语义已经独立，但业务 runner
仍通过 `LegacyBridgeContext -> JarvisHandler(no_dingtalk=True)` 进入
`jarvis_dingtalk_bot.py`。这不是可接受的最终边界：Scheduler 进程会加载并构造执行器、旧
周期对象和交互层依赖，worker 也继续把实际 Task 执行回调留在 Bot 中。

本设计的目标是将三个运行时完全分离：

```text
                 DingTalk / Tata inbound
                           |
                           v
                 bridge/dingtalk/handler.py
                           |
                 control-plane Task submit
                           v
  Scheduler ----> bridge/scheduler/* <---- jobs.yaml
      |                    |
      |             periodic Job runners
      |                    |
      +--------------------+----> Aone / GitHub / event clients

  Scheduler host / worker host
             |
             v
    bridge/execution/persistent_task.py
             |
             v
      PersistenceExecutor lease callback
```

最终不得存在以下依赖：

- `bridge/scheduler/** -> jarvis_dingtalk_bot`
- `bridge/execution/** -> jarvis_dingtalk_bot`
- Scheduler 进程构造 `JarvisHandler`、`PersistenceExecutor` 或 `EphemeralExecutor`
- worker 进程导入 Scheduler registry、engine 或 runner

`jarvis_dingtalk_bot.py` 在完成迁移后应移除或降级为兼容入口；不能继续作为领域服务的
composition root。

## 2. 目标目录与依赖规则

```text
bridge/
├── scheduler/
│   ├── engine.py                 # 时间规划、slot admission、恢复、并发、drain
│   ├── registry.py               # jobs.yaml 解析、definition 校验、runner registry
│   ├── service.py                # Scheduler Worker 生命周期和控制面协议
│   ├── composition.py            # Scheduler 进程唯一 composition root
│   ├── jobs.yaml
│   └── runners/
│       ├── aone.py               # aone.scan / aone.claim-health / aone.reply
│       ├── daily.py              # daily.nudge / daily.probe
│       ├── pr.py                 # pr.watch
│       └── recovery.py           # external.recovery
├── execution/
│   ├── composition.py            # worker / scheduler-host executor composition
│   └── persistent_task.py        # PersistenceExecutor 的 lease callback
├── dingtalk/
│   ├── handler.py                # 入站消息、卡片、交互展示
│   └── composition.py            # DingTalk stream composition root
├── clients/
│   ├── aone.py                   # Aone 查询、认领/状态操作、事件投递窄接口
│   ├── github.py                 # PR / CI / 评论窄接口
│   └── notifications.py          # DingTalk 发送与 durable-event adapter
├── runtime/
│   ├── headless.py               # Jarvis CLI / transcript / process guard
│   └── event_ledger.py           # Aone、DingTalk 双通道幂等投递
└── main.py                       # 只按 JARVIS_BRIDGE_ROLE 装配进程
```

依赖方向只允许从 composition 到具体服务：

```text
dingtalk -> execution / clients / runtime
scheduler -> clients / runtime
execution -> clients / runtime
clients, runtime -> 标准库及窄外部 SDK
```

`dingtalk` 不被 `scheduler` 或 `execution` 导入。`scheduler` 与 `execution` 彼此不导入；
它们只共享协议、clients 与 runtime 中的无 UI 基础能力。

## 3. 三个独立进程

### 3.1 Scheduler host

`JARVIS_BRIDGE_ROLE=scheduler bridge/run.sh start` 启动两个并列进程：

1. Scheduler Worker：`scheduler/composition.py` -> `SchedulerService`。
2. Persistent Task Worker：`execution/composition.py` -> `PersistenceExecutor`。

两者由同一个 `run.sh` 生命周期管理，但不能互相构造。计划 restart 时先停止 Scheduler 的
admission 并 drain 已准入的 Job；仅在 Scheduler drain 成功后，才 drain/stop Task Worker。
Scheduler drain 超时保持原进程，拒绝启动 successor。

### 3.2 Worker host

`JARVIS_BRIDGE_ROLE=worker bridge/run.sh start` 仅启动
`execution/composition.py`。它不读取 `jobs.yaml`、不注册 `bridge-scheduler` Worker、
不创建定时线程，也不导入 `scheduler.*`。

### 3.3 DingTalk host（可与 Scheduler host 共置）

DingTalk stream 仅初始化 `dingtalk/composition.py` 与 `dingtalk/handler.py`。
收到卡片或人工消息后只做认证、交互状态维护与 Task submit；需要执行的动作一律写入控制面。
它不能持有 Scanner、PR watcher 或 Daily job，也不能通过本地线程绕过 Task fence。

## 4. Scheduler runner 的实际职责

每个 runner 都是一次性、无常驻线程的 `run(definition, scheduled_for)`，返回
`JobResult`。SchedulerEngine 是唯一 cadence、misfire、retry、overlap 和 drain owner。

| Job | 新 runner | 最小依赖 | 不允许依赖 |
| --- | --- | --- | --- |
| `aone.scan` | `runners/aone.py:AoneScanRunner` | AoneClient、TaskPublisher、field preflight | DingTalk handler、Executor |
| `aone.claim-health` | `runners/aone.py:ClaimHealthRunner` | AoneClient、TaskClient、EventLedger | Scanner 对象、Bot |
| `aone.reply` | `runners/aone.py:AoneReplyRunner` | AoneClient、TaskClient、wake publisher | `JarvisHandler._wake` |
| `daily.nudge` | `runners/daily.py:DailyNudgeRunner` | AoneClient、EventLedger、clock | `DailyScheduler` |
| `daily.probe` | `runners/daily.py:DailyProbeRunner` | HeadlessRuntime、summary store | Bot module / handler |
| `pr.watch` | `runners/pr.py:PrWatchRunner` | GitHubClient、AoneClient、TaskPublisher、EventLedger | `PrWatchScheduler` |
| `external.recovery` | `runners/recovery.py:ExternalRecoveryRunner` | TaskClient、AoneClient、EventLedger | `ExternalOperationRecoveryScheduler` |

Runner 只可保存可重建的业务 cursor（例如 PR 观察记录）；slot 当前态、重试和运行中状态必须由
控制面 scheduled-job 协议拥有。不得把通用 checkpoint 重新写回旧 JSON 状态文件。

## 5. Persistent Task 执行边界

`execution/persistent_task.py` 提供显式的 `execute_lease(lease, controller)`：

- 构造 prompt、选择 Lane、调用 `HeadlessRuntime`；
- 按结构化 Task result 执行 claim / wrap / suspend / release / finish；
- 负责 task fence、子进程 bind/stop、session progress 和错误归类；
- 使用 `clients` 与 `runtime`，但不导入 `scheduler` 或 `dingtalk`。

这保留现有 PersistentExecutor 的 lease/fence 语义，且让 worker 在没有 Scheduler 的机器上可
独立运行。`aone.reply` 只写 wake Task；真正的 resumed headless execution 仍由任意可用 worker
lease 后执行。

## 6. 迁移步骤

每一步只迁移一个可验证边界；旧路径与新路径不能同时拥有同一个 Job 的 cadence。

### M0：建立共享窄接口

先从 Bot 中提取 `AoneClient`、`GitHubClient`、`EventLedger`、`HeadlessRuntime` 和
`TaskPublisher` 的最小接口及 fake。此阶段不移动业务行为，禁止引入新的常驻线程。

验收：Scheduler unit test 可使用 fake clients 构造所有 runner，而不 import Bot。

### M1：迁移 Aone 三个 Job

依次迁移 `aone.scan`、`aone.claim-health`、`aone.reply` 到 `runners/aone.py`。先固定现有
输入/输出、过滤条件、event key、Task payload 和错误重试语义，再删除对应 `AoneScheduler` 与
`AoneReplyScheduler` 的单次 tick API。

验收：`scheduler/runners/aone.py` 不 import Bot；Bot 不再构造 scanner/reply scheduler；真实
Scheduler 启动时 module import graph 不含 `jarvis_dingtalk_bot`。

### M2：迁移 daily Job

迁移 `daily.nudge` 与 `daily.probe` 到 `runners/daily.py`。Probe 的 Jarvis CLI adapter 同时从
`headless/jarvis_adapter.py` 中移出 Bot，改为 `runtime/headless.py` 的窄 process/transcript
adapter。

验收：daily runner 在无 DingTalk 凭证、无 Tata 依赖下可启动；重试不会跨 daily slot。

### M3：迁移 PR 与外部恢复

迁移 `pr.watch` 到 `runners/pr.py`，迁移 `external.recovery` 到 `runners/recovery.py`。保留
PR registry 和 event-ledger 的兼容读取，但在首个新格式写入后由迁移函数原子升级。

验收：Bot 不再持有 `PrWatchScheduler` 或 `ExternalOperationRecoveryScheduler`；Scheduler
重启后 PR dedup、CI 上限和双通道事件账本不回退。

### M4：抽取 Persistent Task 执行器

将 `JarvisHandler._execute_task_lease`、任务结果处理、session/process guard 和 task bookend
移入 `execution/persistent_task.py`。`execution/composition.py` 显式创建 TaskClient、
CapacityManager、PersistenceExecutor 与 callback。

验收：worker role 启动时不 import Bot、Scheduler 或 DingTalk SDK；现有 task fence 和
resume/suspend 端到端测试不变。

### M5：收缩 DingTalk adapter

将入站处理移入 `dingtalk/handler.py`，把卡片/广播抽为 NotificationPort。删除
`JarvisHandler(no_dingtalk=True)`、`LegacyBridgeContext`、`scheduler/runners/legacy.py` 与
`jarvis_dingtalk_bot.py` 中所有周期类定义。

验收：全仓不存在 `LegacyBridgeContext`、`no_dingtalk=True` 的 Scheduler 使用，且
`scheduler/**`、`execution/**` 对 Bot 的 import 为零。

## 7. 迁移安全门

每个 M 阶段都必须同时满足：

1. 新旧 runner 不能并行注册同一 `job_key`；registry 明确只留一个 owner。
2. Scheduler restart：DRAINING 期间拒绝新的 start/recover，已准入 slot 的 complete/fail
   只能由同一 process UUID 提交。
3. worker restart：lease/session fence 仍能终止失主进程，不能因拆分退回本地 inflight 状态。
4. Aone/DingTalk 重要事件继续各自独立幂等；迁移不允许直接 comment 绕过 ledger。
5. 每个迁移 Job 必须有 fake-client unit、Scheduler service integration、一次真实控制面
   register/start/complete 观察；不把 stub READY 当成预发证明。
6. 每阶段完成后执行 import-boundary test：scheduler、execution、worker entrypoint 均不得
   import DingTalk adapter 或 Bot compatibility module。

## 8. 完成定义

只有当以下条件全部成立，才可宣称“Scheduler 已脱离 Bot”：

- `bridge/scheduler/runners/legacy.py` 已删除；
- `build_runners()` 不再引用任何 `Legacy*` 类型；
- `scheduler/main.py` 的构造路径不加载 Bot、DingTalk SDK、PersistenceExecutor 或
  EphemeralExecutor；
- `execution/persistent_task.py` 是唯一 PersistentExecutor callback owner；
- DingTalk handler 不构造任一定时 Job；
- scheduler host、worker host、DingTalk host 的 import-boundary test 和 restart/drain
  验收均为绿；
- 旧 Bot 模块删除或仅保留无业务逻辑的兼容 re-export，且无生产入口依赖它。

# Bridge Runtime 去耦设计：Scheduler、执行器与 DingTalk Adapter

## 1. 结论与目标

当前已完成第一段运行时拆分：`SchedulerService` 的时钟、slot 与 admission 由独立进程拥有；
`bridge/task_worker.py` 是 Persistent Task 执行的独立入口。`bridge/run.sh` 为 scheduler 角色分别
管理 Bot、task worker 与 Scheduler 的 PID/log 生命周期，而 worker 角色只启动 task worker。

周期 Job 也不再经由 `scheduler/runners/legacy.py` 取得 runtime catalogue。reply、recovery、Aone、
nudge 与 PR 各有专用 runner 模块，且 `build_runners()` 直接装配它们。`legacy.py` 仅保留
import-safe 兼容标记。

这仍只是运行时编排拆分，不是完整领域代码迁移：部分 runner 会按需导入 Bot 中暂存的业务状态机
（例如 `AoneScheduler`、`_NudgeJob`、`PrWatchScheduler`）；task worker 也通过
`JarvisHandler(no_dingtalk=True)` 取得当前的 Task callback。因此当前不能宣称 Scheduler 或 worker
已完全脱离 `jarvis_dingtalk_bot.py`，更不能以此文档宣称已部署、已预发验证或已验证 restart 行为。

最终目标是将三个运行时完全分离：

```text
                 DingTalk / Tata inbound
                           |
                           v
                 bridge/dingtalk/handler.py
                           |
                           v
                 control-plane Task submit

  Scheduler ----> bridge/scheduler/* <---- jobs.yaml
      |                    |
      |             periodic Job runners
      |                    |
      +--------------------+----> Aone / GitHub / event clients

  Scheduler host / worker host
             |
             v
     independent Persistent Task callback
```

最终不得存在以下依赖（当前尚未全部满足）：

- `bridge/scheduler/** -> jarvis_dingtalk_bot`
- Task worker callback -> `jarvis_dingtalk_bot`
- Scheduler 进程构造 `JarvisHandler`、`PersistenceExecutor` 或 `EphemeralExecutor`
- worker 进程导入 Scheduler registry、engine 或 runner

## 2. 当前目录与依赖状态

```text
bridge/
├── main.py                       # SchedulerService 入口
├── task_worker.py                # 独立 Persistent Task Worker 入口
├── scheduler/
│   ├── engine.py                 # 时间规划、slot admission、恢复、并发、drain
│   ├── registry.py               # jobs.yaml 解析、definition 校验、runner registry
│   ├── service.py                # Scheduler Worker 生命周期和控制面协议
│   └── runners/
│       ├── aone.py               # aone.scan / aone.claim-health
│       ├── reply.py              # aone.reply
│       ├── nudge.py              # daily.nudge
│       ├── daily_probe.py        # daily.probe
│       ├── pr.py                 # pr.watch
│       ├── recovery.py           # external.recovery
│       └── legacy.py             # 空的兼容标记，非 runtime catalogue
└── jarvis_dingtalk_bot.py        # 当前 Bot/领域状态机与部分 callback 依赖
```

当前允许的临时依赖是专用 runner 在实际 slot 运行时按需加载遗留业务实现；它们不构造
`JarvisHandler` 或遗留后台 loop。最终依赖方向应收敛为：

```text
dingtalk -> task execution / clients / runtime
scheduler -> clients / runtime
task execution -> clients / runtime
clients, runtime -> 标准库及窄外部 SDK
```

## 3. 当前进程与入口边界

### 3.1 Scheduler host

`JARVIS_BRIDGE_ROLE=scheduler bridge/run.sh start` 当前启动三个独立进程：

1. 兼容 Bot / DingTalk stream：`jarvis_dingtalk_bot.py`。
2. Persistent Task Worker：`task_worker.py` -> `PersistenceExecutor`。
3. Scheduler Worker：`main.py` -> `SchedulerService`。

`run.sh` 为 task worker 使用单独 pidfile/log。在新启动路径中，Bot 启动后再启动 task worker，
然后启动 Scheduler；若 Scheduler 启动失败，脚本停止刚启动的 Bot，但保留 task worker 以保护既有
lease。停止路径分别向 Scheduler 与 worker 发送终止请求，超时时保留对应进程。上述仅为代码路径，
尚无真实部署、重启或预发验证结论。

### 3.2 Worker host

`JARVIS_BRIDGE_ROLE=worker bridge/run.sh start` 仅启动 `task_worker.py`，不启动 DingTalk listener
或 SchedulerEngine。当前 worker 仍导入 Bot 以构造无 DingTalk handler 和 Task callback，因此
“worker 不 import Bot”仍是后续目标。

### 3.3 DingTalk host

当前 DingTalk stream 仍由 Bot 兼容入口承载。其进一步收缩到独立 adapter，不属于本次已完成范围。

## 4. Scheduler runner 的实际职责

每个 runner 均是一次性、无常驻线程的 `run(definition, scheduled_for)`，返回 `JobResult`。
SchedulerEngine 拥有 cadence、misfire、retry、overlap 与 drain。`legacy.py` 不再导出 runner，
也不被 `build_runners()` 用作 catalogue。

| Job | 当前专用 runner | 当前最小边界 | 未完成的领域迁移 |
| --- | --- | --- | --- |
| `aone.scan` | `runners/aone.py:AoneScanRunner` | TaskClient、受限 Aone context | 按需复用 Bot 的 `AoneScheduler` |
| `aone.claim-health` | `runners/aone.py:AoneClaimHealthRunner` | TaskClient、受限 Aone context | 按需复用 Bot 状态机 |
| `aone.reply` | `runners/reply.py:ReplyRunner` | TaskClient、`WakePersistence` | 使用 Bot 的 Aone/Task 辅助逻辑 |
| `daily.nudge` | `runners/nudge.py:DailyNudgeRunner` | 受限 nudge context | 按需复用 Bot 的 `_NudgeJob` |
| `daily.probe` | `runners/daily_probe.py:DailyProbeRunner` | HeadlessRuntime、summary store | 不在本次拆分范围 |
| `pr.watch` | `runners/pr.py:PrWatchRunner` | TaskClient、受限 PR context | 按需复用 Bot 的 `PrWatchScheduler` |
| `external.recovery` | `runners/recovery.py:ExternalRecoveryRunner` | TaskClient、worker key、repo root | 使用独立 `ExternalOperationRecoveryScheduler` |

Runner 只可保存可重建的业务 cursor（例如 PR 观察记录）；slot 当前态、重试和运行中状态仍必须由
控制面 scheduled-job 协议拥有。

## 5. Persistent Task 执行边界

当前 `task_worker.py` 独立拥有 `PersistenceExecutor` 的创建、启动、信号停止与 executor stop
生命周期，并由 `run.sh` 以独立 pidfile/log 管理。它通过 handler 的
`persistent_task_execution.execute` 提供 lease callback，且停止时清理该 handler 持有的 ephemeral
executor。

因此 worker 可在没有 SchedulerEngine 或 DingTalk stream 的机器上运行，但 task fence、任务结果
处理、session/process guard 和 task bookend 仍在 Bot 依赖链。`aone.reply` 只持久化 wake；真正的
resumed headless execution 仍由任意可用 worker lease 后执行。

## 6. 迁移步骤与状态

每一步只迁移一个可验证边界；旧路径与新路径不能同时拥有同一个 Job 的 cadence。

### M0：建立共享窄接口（未完成）

从 Bot 中提取 Aone、GitHub、event ledger、HeadlessRuntime 和 TaskPublisher 的最小接口及 fake，
使 runner 的业务实现不再依赖 Bot。

### M1：Aone 与 reply Scheduler 接入（已完成接入，业务迁移未完成）

`aone.scan`、`aone.claim-health` 已接入 `runners/aone.py`，`aone.reply` 已接入
`runners/reply.py`；registry 只从这些专用模块装配。现阶段保留 Bot 中的单次业务 tick，供受限
context 按需调用；尚未删除其领域实现。

### M2：nudge Scheduler 接入（已完成接入，业务迁移未完成）

`daily.nudge` 已接入 `runners/nudge.py`，`daily.probe` 保持 `runners/daily_probe.py`。nudge 仍
复用 Bot 的 `_NudgeJob`；HeadlessRuntime 的目录收敛不属于本次范围。

### M3：PR 与外部恢复 Scheduler 接入（已完成接入，业务迁移未完成）

`pr.watch` 与 `external.recovery` 分别由 `runners/pr.py` 和 `runners/recovery.py` 装配。PR runner
仍复用 Bot 中的 `PrWatchScheduler`；恢复 runner 使用独立的
`ExternalOperationRecoveryScheduler`。PR registry/event-ledger 的格式迁移不在本次范围。

### M4：独立 Task Worker 入口与生命周期（已完成入口拆分，callback 迁移未完成）

`task_worker.py` 已独立创建、启动和停止 `PersistenceExecutor`，并使 worker role 不启动
Scheduler/Bot stream。下一步是将任务结果处理、session/process guard 和 task bookend 迁至独立
模块，使其不再通过 `JarvisHandler(no_dingtalk=True)` 取得 callback。

### M5：收缩 DingTalk adapter（未开始）

将入站处理、卡片与广播从 Bot 兼容入口收敛到独立 adapter；随后删除遗留业务类和 Bot runtime
依赖。

## 7. 后续迁移安全门

后续每个领域迁移阶段都必须同时满足：

1. 新旧 runner 不能并行注册同一 `job_key`；registry 明确只留一个 owner。
2. Scheduler drain 期间拒绝新的 start/recover，已准入 slot 的 complete/fail 只能由同一 process UUID 提交。
3. worker 拆分不能退回本地 inflight 状态，lease/session fence 仍须能终止失主进程。
4. Aone/DingTalk 重要事件继续各自独立幂等；迁移不允许直接 comment 绕过 ledger。
5. 每个迁移 Job 必须有 fake-client unit、Scheduler service integration、一次真实控制面
   register/start/complete 观察；不把 stub READY 当成预发或部署证明。
6. 在 Bot 依赖实际移除前，不以 import-boundary test 绿灯宣称该依赖已消失。

## 8. 完成定义

当前实现尚不能宣称“Scheduler 已脱离 Bot”，也没有本次部署、预发或真实 restart 验证。只有当以下
条件全部成立，才可作该宣称：

- `legacy.py` 可删除，且 `build_runners()` 不引用任何 Legacy 类型或延迟 Bot 领域实现；
- Scheduler 的构造路径不加载 Bot、DingTalk SDK、`PersistenceExecutor` 或 `EphemeralExecutor`；
- 独立 Task callback 成为唯一 `PersistenceExecutor` callback owner；
- DingTalk adapter 不构造任一定时 Job；
- scheduler host、worker host、DingTalk host 的 import-boundary test 与 restart/drain 验收均为绿；
- 旧 Bot 模块删除或仅保留无业务逻辑的兼容 re-export，且无生产入口依赖它。

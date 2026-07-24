# Bridge Runtime 去耦设计：Scheduler、执行器与 DingTalk Adapter

## 1. 结论

Bridge 现有三个独立进程边界：DingTalk inbound Bot、Persistent Worker 与
SchedulerService。唯一生命周期入口是 `bridge/run.sh`：scheduler 角色启动三者，worker
角色只启动 Persistent Worker。重启 Scheduler 或 Bot 不会停止已 lease 的 Persistent Worker 进程。

周期任务只由 `bridge/scheduler/jobs.yaml` 的显式条目注册，`jobs.py` 只负责校验和加载；`SchedulerEngine` 拥有
cadence、misfire、retry、overlap 与 drain。Bot 不构造周期 Scheduler，也不持有周期 Job 的
线程或 cadence。

```text
                         DingTalk inbound
                               |
                               v
                    jarvis_dingtalk_bot.py
                               |
                               v
                      control-plane Task submit

 jobs.yaml -> jobs.py -> SchedulerService -> SchedulerEngine -> runners -> Aone / GitHub

 persistent_worker.py -> PersistenceExecutor -> persistent_tasks.py -> leased Task process
```

Scheduler 和 Persistent Worker 均不导入对方的运行时；Scheduler 不构造 DingTalk stream、
`PersistenceExecutor` 或 `EphemeralExecutor`。Persistent Worker 不导入 Scheduler engine 或 runner。
这份设计只记录本地实现边界，不代表已部署、预发验证或完成真实 restart 验收。

## 2. 当前目录与依赖边界

```text
bridge/
├── run.sh                        # 唯一进程监督入口
├── main.py                       # SchedulerService composition root
├── persistent_worker.py          # Persistent Worker composition root
├── jarvis_dingtalk_bot.py        # DingTalk inbound adapter
├── headless_runtime.py           # Headless 执行窄接口
├── persistent_tasks.py           # Persistent Task / wake / bookend 执行边界
└── scheduler/
    ├── jobs.yaml                 # 七个显式 Job definition
    ├── jobs.py                   # YAML 校验与加载
    ├── model.py                  # schedule、definition、runner、result
    ├── engine.py                 # 计划、准入、执行与终态提交
    ├── control_plane_client.py   # scheduled-job 控制面边界
    ├── service.py                # Worker 生命周期、心跳、READY 与 drain
    └── runners/
        ├── scan.py               # aone.scan
        ├── claim_health.py       # aone.claim-health
        ├── daily_nudge.py        # daily.nudge
        ├── pr_watch.py           # pr.watch
        ├── reply.py              # aone.reply
        ├── recovery.py           # external.recovery
        └── daily_probe.py        # daily.probe
```

`headless_runtime.py` 和 `persistent_tasks.py` 是扁平共享模块，不再以目录层级制造额外适配边界。
Runner 只保留可重建的业务 cursor；slot 当前态、重试与运行中状态由 scheduled-job 控制面拥有。

## 3. 进程生命周期

### Scheduler host

`JARVIS_BRIDGE_ROLE=scheduler bridge/run.sh start` 启动：

1. `jarvis_dingtalk_bot.py` 的 DingTalk stream；
2. `persistent_worker.py` 的 `PersistenceExecutor`；
3. `main.py` 的 `SchedulerService`。

`run.sh` 为三个进程维护独立 pidfile/log。计划内 Scheduler 停止先关闭 admission、登记
`DRAINING` 并等待已准入 slot；超时不以 SIGKILL 替换仍在 drain 的 Scheduler。Persistent Worker
独立持有 Task lease，因此 Scheduler/Bot 重启不会中断其 Session。

### Worker host

`JARVIS_BRIDGE_ROLE=worker bridge/run.sh start` 只启动 `persistent_worker.py`。它从控制面 lease
Task、启动 headless 子进程并回写结果；不启动 DingTalk listener 或 SchedulerService。

## 4. 周期 Job 职责

每个 runner 均为一次性 `run(definition, scheduled_for)`，返回 `JobResult`，不得创建常驻线程或
自行 sleep。

| Job | runner | 职责 |
| --- | --- | --- |
| `aone.scan` | `runners/scan.py` | Aone 并集扫描、Task 写入与状态观察 |
| `aone.claim-health` | `runners/claim_health.py` | claimed 健康核验与事件 ledger 补偿 |
| `daily.nudge` | `runners/daily_nudge.py` | Terraform idle 停滞双通道催办 |
| `aone.reply` | `runners/reply.py` | 回复后持久化 wake，供 worker lease 执行 |
| `daily.probe` | `runners/daily_probe.py` | 独立 Terraform Headless 探测；默认停用 |
| `pr.watch` | `runners/pr_watch.py` | PR、CI、评审与合并生命周期；按活跃度调整 next due |
| `external.recovery` | `runners/recovery.py` | 以 recovery token 只读核验外部回执 |

`daily.probe` 是唯一刻意停用的 Job；启用真实 Terraform 操作前必须获得单独授权。

## 5. 安全门与完成定义

1. 同一 `job_key` 只有一个 runner owner，不能存在第二个 cadence。
2. Scheduler drain 期间拒绝新的 start/recover；已准入 slot 的 complete/fail 只能由同一
   process UUID 提交。
3. Worker 只依赖控制面 lease/session fence，不回退为本地 inflight 状态。
4. Aone 与 DingTalk 重要事件分别通过幂等 ledger 发布；迁移不得绕过 ledger 直接评论。
5. 每个 Job 必须有 fake-client unit 与 Scheduler service integration；stub READY 不能作为
   预发或部署证明。
6. 导入边界和 restart/drain 测试通过后，仍需单独授权真实控制面及预发观察。

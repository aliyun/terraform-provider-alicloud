# Bridge 定时任务模块重构 CR 总结

> 汇总日期：2026-07-24
>
> CR：[28675904](https://code.alibaba-inc.com/terraflow/jarvis-preview/codereview/28675904)
>
> Aone：[重构 Bridge 定时任务模块：统一调度规范与安全重启设计](https://project.aone.alibaba-inc.com/v2/project/2100304/req/84352841)
>
> 对比基线：`origin/master@756c986c53e1b50a6b88475ae0bf5a4cf258cb95`
>
> 代码汇总基线：`d4f3084bc7252bbdd989e7b91407b38bc464db17`

## 1. 结论

本 CR 将 Bridge 的周期调度、持久任务执行和钉钉入站处理拆成三个明确的 Python
入口，并由 `bridge/main.py` 统一监督。`bridge/run.sh` 只管理一个
`bridge.main` 外层进程，不再复制三套组件的 start/stop/status 分支。

最终进程边界如下：

```text
bridge/run.sh
└── python -m bridge.main
    ├── python -m bridge.scheduler.scheduler
    ├── python -m bridge.persistent_worker
    └── python -m bridge.jarvis_dingtalk_bot
```

- `scheduler` 主机运行 Scheduler、Persistent Worker 和 DingTalk Bot。
- `worker` 主机只运行 Persistent Worker。
- Scheduler 或 Bot 启动失败、运行时退出或单独重启，不会停止已经 READY
  且可能持有 Task lease 的 Persistent Worker。
- 计划内 Bridge 重启采用 PID 绑定的 worker 保留标记，新 supervisor
  只接管该标记指向的存活 worker，避免误接管陈旧或无关进程。

## 2. 模块职责

| 文件 | 最终职责 |
| --- | --- |
| `bridge/run.sh` | source 环境、解析 `start/stop/status/restart`、选择 role，并机械管理 `bridge.main` |
| `bridge/main.py` | 统一监督三个入口、汇总 READY、隔离组件失败、执行安全停止和 worker 接管 |
| `bridge/scheduler/scheduler.py` | SchedulerService 的纯 composition root；负责配置校验、启动、READY 和 drain |
| `bridge/persistent_worker.py` | Persistent Worker composition root；只负责 lease 和持久 Task 执行 |
| `bridge/jarvis_dingtalk_bot.py` | DingTalk inbound adapter；不再拥有任何周期 Scheduler |
| `bridge/pending_dispatch.py` | `JARVIS_AUTO_DISPATCH=0` 时的持久化待派发状态和原子消费 |
| `bridge/pr_watch_registry.py` | Bot 写入、`pr.watch` runner 读取的共享持久注册表；不是第四个进程入口 |
| `bridge/install-launchd.sh` | 安装前校验、launchd 启停、失败回滚和安全重启 |
| `bridge/scheduler/jobs.yaml` | 七个周期 Job 的唯一显式定义 |

`bridge/pr_watch_registry.py` 的存在是为了解耦数据所有权：Bot 接收入站命令并写入
PR watch 注册信息，Scheduler 的 `pr.watch` runner 按计划读取并处理。它不创建线程、不
管理 cadence，也不参与 `main.py` 的进程监督。

## 3. Scheduler 统一调度

所有周期 Job 都由 `bridge/scheduler/jobs.yaml` 显式注册，runner 每次只执行一个
有界 tick，不自行创建常驻线程或 sleep：

| Job | Runner | 作用 |
| --- | --- | --- |
| `aone.scan` | `runners/scan.py` | Aone 并集扫描、Task 写入与状态观察 |
| `aone.claim-health` | `runners/claim_health.py` | claim 健康核验与事件 ledger 补偿 |
| `daily.nudge` | `runners/daily_nudge.py` | Terraform idle 停滞双通道催办 |
| `aone.reply` | `runners/reply.py` | 回复后的持久化 wake |
| `daily.probe` | `runners/daily_probe.py` | 独立 Terraform Headless 探测；默认停用 |
| `pr.watch` | `runners/pr_watch.py` | PR、CI、评审和合并生命周期 |
| `external.recovery` | `runners/recovery.py` | 以 recovery token 核验不确定外部回执 |

`SchedulerEngine` 统一拥有 cadence、misfire、retry、overlap、admission 和 drain；
AutomationAgent scheduled-job 控制面保存 definition、slot 和执行终态。业务 cursor
仍保留在对应 runner 的持久化边界中，不能代替控制面状态。

## 4. Review 发现的四个兼容性问题及修复

### 4.1 `JARVIS_AUTO_DISPATCH=0` 在 supervisor 模式下失效

原问题：入口拆分后，Bot 和 Worker 不再共享进程内对象，关闭自动派发时产生的待派发
状态可能只存在于 Bot 内存，Worker 无法消费。

修复：新增 `bridge/pending_dispatch.py`，将待派发信息原子写入持久化状态；Bot
负责登记，Persistent Worker 在控制面 lease 路径中消费。关闭自动派发的既有语义得到保留，
没有回退为隐式自动派发。

### 4.2 Scheduler/Bot 故障会连带停止 READY Worker

原问题：外层 supervisor 如果把任一子进程失败视为整体失败，会中断已经持有 Task lease
的 Worker，扩大非执行组件故障的影响范围。

修复：`bridge/main.py` 独立监督三个入口。Scheduler 和 Bot 的启动失败或运行时退出只进入
各自的重启路径；只要 Worker 已 READY，外层服务仍保持执行能力。外层 READY 以 Scheduler
和 Worker 为必要条件，不等待 Bot 建立长连接。

### 4.3 重启期间无法安全保留正在执行的 Worker

原问题：普通 stop/start 会终止 Persistent Worker，可能打断同机正在执行的普通 Task；
宽泛的 pidfile 接管又可能误接管陈旧进程。

修复：

1. 计划内 restart 先让 Scheduler 关闭 admission 并 drain 已准入 slot。
2. Bot 正常停止。
3. Persistent Worker 保持运行，继续维护已有 lease。
4. `run.sh` 写入带旧 supervisor PID、worker PID 和身份信息的短期接管标记。
5. 新 `bridge.main` 只有在父进程和 worker 身份均匹配、目标进程仍存活时才接管。
6. 标记无效、过期或 PID 不匹配时 fail closed，不把任意进程当作受管 Worker。

local 和 launchd 两条 restart 路径复用同一接管协议。`launchctl kill` 失败时会重新 enable
服务，避免 launchd 被留在禁用状态。

### 4.4 安装/重启在停旧服务后才发现新运行时无效

原问题：依赖、解释器或配置错误若在旧服务停止后才暴露，会造成不必要停机。

修复：`bridge/install-launchd.sh` 和 `bridge/run.sh` 在发停止信号前完成 Python
解释器、依赖、role、Scheduler definition 和关键配置校验。校验失败时保留旧服务，不进入
重启窗口。

## 5. 与最新 master 的兼容同步

CR 在开发 worktree 中持续以 `origin/master` 为约束，并已合入当时最新的
`756c986c53e1b50a6b88475ae0bf5a4cf258cb95`。同步时保留了以下 master 行为：

- 字段修复继续原地完成，不创建单独的 `field_repair` Task。
- Aone 评论 cursor 保持幂等，避免重访时重复处理同一回复。
- macOS boot id 使用稳定来源，避免进程重启被误判为主机重启。
- 提交人钉钉私信通知逻辑仍在新的职责边界中生效。
- Terraform 截图证据与最终回复契约继续由现有工作流校验。
- `bridge.jarvis_dingtalk_bot._ticket_prompt` 保留 re-export，兼容已有测试和调用方。

本 CR 没有在 master checkout 上开发；所有 tracked 修改均位于
`worktree-scheduler-refactor-v14` worktree，master 只作为兼容基线。

## 6. 重启与 READY 时序

```text
预校验新 runtime / 配置
        │
        ▼
Scheduler 停止 admission 并 drain
        │
        ├── Bot 正常停止
        │
        └── Persistent Worker 保持运行和 lease
        │
        ▼
旧 main 写入 PID 绑定的接管标记并退出
        │
        ▼
新 main 校验并接管原 Worker
        │
        ├── 启动新 Scheduler
        └── 启动新 Bot
        │
        ▼
Scheduler + Worker READY → 外层 READY
```

异常崩溃与计划内 restart 的语义不同：Scheduler 的 interrupted slot 由
`recover-interrupted` 依据控制面 fence 恢复；它不承诺恢复 Python 栈或外部调用的半进度。
Persistent Worker 则依靠 Task/Session lease 继续或收敛执行。

## 7. 验证结果

| 验证项 | 结果 |
| --- | --- |
| `test/bridge_run_test.sh` | 87/87 通过 |
| `test/bridge_dispatch_test.sh` | 299/299 通过 |
| 最新 master 同步后的 focused tests | 36/36 通过 |
| `test/interactive_claim_test.sh` | 18/18 通过 |
| Aone screenshot prompt/evidence tests | 通过 |
| Bridge 全量 discovery | 587 项；3 failures + 3 errors |

Bridge 全量 discovery 的 6 个异常全部位于 `bridge.test_interactive_worker`，并且在当前
master checkout 上以相同方式复现，因此记录为 master 基线问题，不是本 CR 新增回归。
本地/静态测试和 stub 控制面验证不等同于预发行为或正式发布验收。

## 8. 变更边界与待审结论

- 不保留旧 `bridge/task_worker.py` 兼容壳；持久执行入口统一为
  `bridge.persistent_worker`。
- `bridge/main.py` 只做入口监督，不接管 Scheduler 的 cadence，也不承载 runner 业务。
- `bridge/run.sh` 不再实现三套重复组件生命周期逻辑。
- `bridge/pr_watch_registry.py` 是共享状态模块，不是额外 daemon。
- `daily.probe` 仍默认停用；启用真实 Terraform 操作需要独立授权。
- 本次只完成 CR 代码、兼容性修复和本地验证，没有合并 master，也没有执行正式发布。

当前结论：CR 已具备人工评审条件；合并与正式发布仍由仓库负责人决定。

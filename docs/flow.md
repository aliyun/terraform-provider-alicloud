# 主流程（bridge 自动扫派）

## 0. 自举 (Bootstrap)
- `bootstrap/install.sh` 装依赖
- `bootstrap/preflight.sh` 24h 闸门（install + verify），全绿才干活

## 1. 扫描与派发 (bridge)
- `bridge/run.sh start` 常驻：ScanScheduler 定时跑 `scan.sh --force`，扫 `config/pools.json` 登记的池
- `ExecutionRouter` 只按可恢复性分类：业务工单/重访/唤醒/PR 跟进 → `Task`；probe/本地检查/一次性命令 → `EphemeralJob`
- `PersistenceExecutor` 从控制面 lease Task；`EphemeralExecutor` 执行本地一次性作业；两者共享 `CapacityManager` 与 `ExecutionRuntime`
- Task 必须先进入控制面；控制面不可用时 fail-closed，不允许回退到本地无状态执行
- `_decide` 逐单判定：终态 / `jarvis-done` / `jarvis-claimed` / `jarvis-npe` → skip；`jarvis-idle` 过人工介入门
- **[人工点·可选]** 回退模式 `JARVIS_AUTO_DISPATCH=0`：新单入 pending，钉钉授权后才派（`plan.sh` 出计划）
- 后台调度：统一 Scheduler 从 `bridge/scheduler/jobs.yaml` 注册 Headless/Handler Job（`daily.probe` 当前默认关闭）；RevisitScheduler 负责每日 idle 重访，AoneScheduler 负责扫描及 stale/done 状态对账，PersonaScheduler 负责评论区数字人接力（默认关）

## 2. 单工单 Triage (headless)
- `claim.sh claim` 认领（竞争锁，输了 SKIP）
- `aone-triage` skill 查证（OpenAPI + Cloudspec 映射 + provider 源码）
- autonomy 判定：auto 列表内自动执行到预发/CR；低置信/红线 → Task `SUSPENDED` + needs-attention 事件
- 遇必须人类决策 → Aone 评论 @人 + `[[SUSPEND:...]]` 挂起，WaitWatcher 收到回复后 `--resume` 唤醒
- 收尾 bookend：`wrap.sh done`（评论+状态）→ `claim.sh release`（打 `jarvis-idle`）

## 3. 硬门 (Hard Gate)
- **[人工点]** `release_prod` 永停，正式发布必须人工审批

## 4. 收敛 (Convergence)
- 服务端 reaper：按 Worker/Session heartbeat、lease 与 fence 收敛中断执行
- AoneScheduler：扫描/派发并检查 `jarvis-done`/Aone 完成态漂移；ClaimHealthScheduler 独立以 ≤5min 节奏对账 claimed 工单与 Task/Session 心跳，异常经 Aone/钉钉双通道幂等发布
- `wrap-check.sh` Stop 闸门：会话结束校验未完工工单已回填
- `SUSPENDED` Task + Aone 是人工决策真源；`runs/` 保留运行审计

---

## 人工审核点
1. **正式发布**：`release_prod` 永停，人工审批
2. **回退模式派发授权**：`JARVIS_AUTO_DISPATCH=0` 时钉钉逐单授权
3. **needs-attention 事件**：低置信 / 红线 / 缺能力，人工拍板后恢复 Task

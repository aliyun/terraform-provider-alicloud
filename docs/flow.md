# 主流程（bridge 自动扫派）

## 0. 自举 (Bootstrap)
- `bootstrap/install.sh` 装依赖
- `bootstrap/preflight.sh` 24h 闸门（install + verify），全绿才干活

## 1. 扫描与派发 (bridge)
- `bridge/run.sh start` 是唯一机械生命周期入口；它只管理 `bridge.main`。Python supervisor 在 scheduler 角色独立监督 Bot、SchedulerService 与 Persistent Worker，在 worker 角色只监督 Persistent Worker
- `ExecutionRouter` 只按可恢复性分类：业务工单/重访/唤醒/PR 跟进 → `Task`；probe/本地检查/一次性命令 → `EphemeralJob`
- `PersistenceExecutor` 从控制面 lease Task；`EphemeralExecutor` 执行本地一次性作业；两者共享 `CapacityManager` 与 `ExecutionRuntime`
- Task 必须先进入控制面；控制面不可用时 fail-closed，不允许回退到本地无状态执行
- `_decide` 逐单判定：终态 / `jarvis-done` / `jarvis-claimed` / `jarvis-npe` → skip；`jarvis-idle` 过人工介入门
- 后台调度：统一 Scheduler 从 `bridge/scheduler/jobs.yaml` 的七项配置注册 Job（`daily.probe` 每日 10:00，Asia/Shanghai）；每个 job 由同名语义的独立 runner 执行：scan、claim_health、daily_nudge、reply、pr_watch、recovery、daily_probe

## 2. 单工单 Triage (headless)
- bridge executor 在模型进程外持有 lease 并托管 `claim`
- `aone-triage` skill 查证（OpenAPI + Cloudspec 映射 + provider 源码）；Terraform PD 同步生成三层本地截图 manifest，最终 RD 无 `--comment` 上传报告并在唯一聚合回复中贴链接
- autonomy 判定：auto 列表内自动执行到预发/CR；低置信/红线 → Task `SUSPENDED` + needs-attention 事件
- 遇必须人类决策 → Aone 评论 @人 + `[[SUSPEND:...]]` 挂起，`reply` runner 收到回复后持久化 wake Task
- 收尾 bookend：`wrap.sh done`（评论+状态）→ `claim.sh release`（打 `jarvis-idle`）

## 3. 硬门 (Hard Gate)
- **[人工点]** `release_prod` 永停，正式发布必须人工审批

## 4. 收敛 (Convergence)
- 服务端 reaper：按 Worker/Session heartbeat、lease 与 fence 收敛中断执行
- `aone.scan`：扫描/派发并检查完成标记/Aone 完成态漂移；`aone.claim-health` 以 5 分钟节奏对账已认领工单与 Task/Session 心跳，异常经 Aone/钉钉双通道幂等发布
- `wrap-check.sh` Stop 闸门：会话结束校验未完工工单已回填
- `SUSPENDED` Task + Aone 是人工决策真源；`runs/` 保留运行审计

---

## 人工审核点
1. **正式发布**：`release_prod` 永停，人工审批
2. **needs-attention 事件**：低置信 / 红线 / 缺能力，人工拍板后恢复 Task

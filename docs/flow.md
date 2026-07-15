# 主流程（bridge 自动扫派）

## 0. 自举 (Bootstrap)
- `bootstrap/install.sh` 装依赖
- `bootstrap/preflight.sh` 24h 闸门（install + verify），全绿才干活

## 1. 扫描与派发 (bridge)
- `bridge/run.sh start` 常驻：ScanScheduler 定时跑 `scan.sh --force`，扫 `config/pools.json` 登记的池
- diff 出新单 + 外部更新单 → DispatchPool 并发起 headless jarvis（每单一实例，上限 `JARVIS_DISPATCH_MAX`）
- `_decide` 逐单判定：终态 / `jarvis-done` / `jarvis-claimed` / `jarvis-npe` → skip；`jarvis-idle` 过人工介入门
- **[人工点·可选]** 回退模式 `JARVIS_AUTO_DISPATCH=0`：新单入 pending，钉钉授权后才派（`plan.sh` 出计划）
- 后台调度：ProbeScheduler（每日探测轮）/ RevisitScheduler（每日 idle 重访）/ ReconcileScheduler（周期对账）/ PersonaScheduler（评论区数字人接力，默认关）

## 2. 单工单 Triage (headless)
- `claim.sh claim` 认领（竞争锁，输了 SKIP）
- `aone-triage` skill 查证（OpenAPI + Cloudspec 映射 + provider 源码）
- autonomy 判定：auto 列表内自动执行到预发/CR；低置信/红线 → escalate
- 遇必须人类决策 → Aone 评论 @人 + `[[SUSPEND:...]]` 挂起，WaitWatcher 收到回复后 `--resume` 唤醒
- 收尾 bookend：`wrap.sh done`（评论+状态）→ `claim.sh release`（打 `jarvis-idle`）

## 3. 硬门 (Hard Gate)
- **[人工点]** `release_prod` 永停，正式发布必须人工审批

## 4. 收敛 (Reconcile)
- `reconcile.sh all`：stale（超时 claim）/ orphan（实例已死）/ drift（台账对账）→ escalate
- `wrap-check.sh` Stop 闸门：会话结束校验未完工工单已回填
- `escalation/` 人工决策队列；`runs/` 审计

---

## 人工审核点
1. **正式发布**：`release_prod` 永停，人工审批
2. **回退模式派发授权**：`JARVIS_AUTO_DISPATCH=0` 时钉钉逐单授权
3. **escalation/ 队列**：低置信 / 红线 / 缺能力，人工拍板

# cap-persona-dispatch-crash-resume

## 缺口类型

`missing_capability` — bridge dispatch 可靠性两个结构性缺口，叠加导致「被打断的派发无人接手，直到人工再次评论才恢复」。

## 触发

排查某客户单（2026-07-14）时发现：该单被 persona 门派发后，数字人回复「将进入 respond 阶段……结论稍后跟进」，随后沉默 **73 分钟**，直到人工再次评论才恢复。根因非偶发延迟，而是下述两缺口叠加。

## 缺口 A：数字人「稍后跟进」空头支票反模式

- headless run 是**单发**执行：run 退出后没有后台任务替它续跑。
- 本次数字人贴出「结论稍后跟进」后，同一 run 内未产出结论即以 `no_result` 结束 = **开空头支票**。
- 约束缺失：数字人 SOP 未禁止「宣告异步跟进后退出」。应二选一：
  1. 同一 run 内把查证做完、当场贴出结论再退出；
  2. 发 `[[SUSPEND:{...}]]` 哨兵进入正式挂起，由 WaitWatcher 唤醒续跑。

**待审补丁**(本分支)：
- `bridge/jarvis_dingtalk_bot.py`: `_persona_prompt` 注入「单发纪律」段（二选一约束）
- `loops/persona-collab.md`: §4.3 写入 SOP 约束（与 §4.4 并列）
- `.claude/agents/terraform-{pd,rd,qa}.md`: 各 agent 文件加「单发纪律(严禁空头支票)」小节

## 缺口 B：graceful-stop 打断的 dispatch 缺 crash-resume

时间线（2026-07-14，见 `.my-day/bridge/bot.log`）：

- `15:09:32` dispatch 以 `no_result` 结束 → 判定 transient，排 `retry 1/2`，进入 30s backoff。
- `15:09:42` bridge 收到 graceful stop 信号 → 清掉在跑 worker（含本单）→ **backoff 中的 retry 被一并杀掉，从未执行**。
- 重启后：`persona-ledger` 已把触发评论标 `processed`、`dispatched.json` 已记录 → persona 门只认新人工评论，不再重派。
- 结果：被打断的活无人接手，沉默至 `16:16` 人工再评论才恢复。

**已闭环**:commit `e6ffd98` (feat/bridge-restart-resume) 落地——
`inflight.json` 登记表 + `dispatch_item` 加 kind 参数 + `terminate_all` 方案B(对有 inflight
记录的 ticket 不释放 claim) + `_resume_inflight()` 重启后自动续跑。

## 置信度

high_conf —— 时间线与代码（`bridge/jarvis_dingtalk_bot.py` retry 循环、graceful stop 清理）、`.my-day/bridge/bot.log`、`persona-ledger.json`、`dispatched.json` 四路证据一致。

## 关联工单

- 本单登记：api_toolkit 池 · 需求 84281589
- 触发来源：客户单 67247576
- 上游 Gap B 已落地：需求 84286464

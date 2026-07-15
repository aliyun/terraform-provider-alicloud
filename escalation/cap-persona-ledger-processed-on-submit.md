# cap: PersonaScheduler ledger 在 submit 入队时即标 processed → 重启丢队列后永不重派

- **类型**：并发/持久化时序缺陷（work 未执行却已记"已处理"）
- **关联**：[[cap-claim-ledger-owner-scoping]]、工单 84297352（原根 @Terraform-研发数字人 两次无人理）；
  同一报单人上一相邻黑洞见 `bridge/test_persona_tracker_scan.py`（84240297，候选集没捞到，与本次不同环节）

## 背景（触发工单 84297352）

客户 bug 单（GitHub #9971，TDE 打挂 alicloud_redis_tair_instance / alicloud_kvstore_account），
原根在 07-15 09:36 和 13:50 两次评论 @Terraform-研发数字人（= 我们的 terraform-rd persona，
WORKER_1783582458263）。PersonaScheduler 两次都识别并 `dispatch #84297352 → terraform-rd`，
但工单至今零回复、状态仍 New。

## 根因

`PersonaScheduler._apply_decisions` 旧逻辑：`ok, _ = self.pool.submit(...)` 返回 `ok=True`
后**同步**把评论写进 `ledger.processed` + 推 `last_seen`。但 `DispatchPool.submit` 的
`ok=True` 只表示"入队接受"——池满（`JARVIS_DISPATCH_MAX` 默认 3）时 future 仅在
ThreadPoolExecutor 内部排队，`_work`（真正 spawn headless 的闭包）**尚未执行**。

时序黑洞：
1. 09:38:01 submit 接受 84297352（3 槽被 84225083/84251052/84255524 占满）→ ledger 标 processed。
   **bot.log 全表只有两行 `persona: dispatch #84297352`，从无对应 `dispatch_item #84297352 start`**
   —— worker 一次都没跑起来（铁证）。
2. 09:54:27 换 token 优雅重启（signal 15）→ `terminate_all` 调 `shutdown(cancel_futures=True)`
   丢弃排队未启动的 future → 84297352 的 persona work 被静默丢弃、评论从未回。
3. ledger 已标 processed → 重启后 PersonaScheduler `skip(reason=processed)` 永久不再重派。
4. 13:50 原根再评论 → 13:52 再 dispatch → 14:03 又一次换 token 重启，重复步骤 1-3。
5. 叠加 ScanScheduler **冷启动**把这单当存量积压播种进基线（gmtModified 13:50 < 重启 14:03，
   不算"新单/外部更新"）→ 主扫描也永不碰它。**双重盲区**：persona 认为已处理、scan 认为是旧积压。

## 本次改动（落盘时点：submit 接受 → worker 真正启动）

`bridge/jarvis_dingtalk_bot.py`：
- `_dispatch_persona` 加 `on_start` 回调参数，包进 `_work`：worker 拿到槽位、真正执行的
  第一行才触发 `on_start()`（排队被 cancel → `_work` 从不执行 → 回调不触发）。
- `_apply_decisions` 把 `last_seen`/`processed`/`dispatch_count`/`escalated` 落盘逻辑封成
  `_commit_ledger` 闭包（绑定当前 cid/d），作为 `on_start` 传入；删除 submit 返回处的同步落盘。

正确性：排队期间 future 已在 `DispatchPool._active`（submit 立即登记），故下一 tick
`_iid_in_flight` 返回 True → `skip(in_flight_active)`、不重复派也不落盘；worker 真启动才
落盘防同评论重派；重启 cancel 队列后 `_active` 随进程消失、ledger 干净 → 下一 tick 自然重派。

补 `bridge/test_persona_ledger_on_start.py`：queued-未启动→ledger 干净、worker 启动→落盘、
skip→仍即时推 last_seen。全绿；既有 persona/scan/handoff 6 套测试全过。

## 残余风险（已知，未在本次消除）

worker **已启动**（on_start 已落盘）但在 headless 回评论**之前**被重启 SIGKILL：processed 已标、
评论未回，仍会烂尾。窗口远小于原 bug（原 bug 在"仅排队"整段都丢）。`terminate_all` 的
`_resume_inflight`（`--resume` 续跑 + 保 claim）目前只覆盖 `kind=="ticket"`，未覆盖 persona。
后续可把 persona-kind 也纳入 in-flight 续跑，或把落盘再推后到"评论确认已回"。

## 置信度

high_conf —— bot.log 时序 + persona-ledger.json（processed=[124851623,124871006]）+ 源码
`submit`/`terminate_all(cancel_futures=True)`/`_apply_decisions` 三处交叉印证。

## 落地

worktree `worktree-persona-ledger-processed-on-start` → 改 `bridge/jarvis_dingtalk_bot.py`
+ 补 `bridge/test_persona_ledger_on_start.py` → MR 待仓库主人合并；绝不自动合 master。
部署 = 合并后重启 bridge。**当轮 84297352 的即时补救**走人工直接接单处理（已 claim + 查证 + 回复），
不依赖本 fix 上线。

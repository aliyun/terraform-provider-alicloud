# Jarvis 主动群消息策略

群聊只承载“没有明确责任人、且需要多人立即介入”的事件。任务状态以控制面、Aone 和运行日志为准；可定向处理的事件进入 Aone 或原会话，不用群消息重复广播。

| 时机 | 默认处理 | 解决方案 |
| --- | --- | --- |
| 自动扫描发现工单并写入 Task | 不发群消息 | 控制面记录 `READY`；日志写 `persisted`。禁止再使用“已自动派发”，因为 Task 入库不代表 Worker 已 lease。 |
| 本地临时队列满 | 不发群消息 | warning 日志；下轮自动重试。持续容量异常由监控告警承接。 |
| `JARVIS_AUTO_DISPATCH=0` 等待人工授权 | 保留群卡片 | 这是显式的人机审批入口，必须让群内授权人看到；回复“处理 #ID/全部处理”后进入队列。 |
| stale `jarvis-claimed` 检查 | 不发群消息 | 控制面 lease/reaper 负责恢复，扫描器只记 warning，避免每轮重复告警。 |
| PR 自动补登记 | 不发群消息 | 登记和 CI/review 游标持久化到 `.my-day/bridge/pr-watch.json`，重启后直接恢复。 |
| PR CI 失败并写入修复 Task | 不发群消息 | 控制面记录 Task；超过自动修复上限才进入“需人工介入”。 |
| PR 新评审评论并写入回复 Task | 不发群消息 | 控制面记录 Task，PR-watch 持久游标防重。 |
| PR 合并并自动收尾工单 | 不发群消息 | Aone 状态/幂等事件已经是交付记录，不再二次广播。 |
| PR 关闭、`jarvis-npe` 或自动修复耗尽 | 保留群告警 | 没有可自动推进路径，需要人工决定；同一 PR 的上限状态由持久台账去重。 |
| DailyScheduler nudge | 不发群消息 | 继续使用 Aone 评论 + 责任人钉钉私信；每日调度状态持久化，且只在配置小时内运行。 |
| 每日 probe 入队、开始、完成 | 不发群消息 | 结果进入 probe 归档/新建工单/日志；群里不播生命周期。 |
| 普通 Task 完成 | 不发群消息 | Aone bookend 和控制面终态是唯一状态源。 |
| 普通 Task 失败（含重试耗尽） | 不发群消息 | 非 Terraform 写 Aone 死因；Terraform 走 RD 幂等事件；原始错误只留日志。 |
| Task 挂起等待某人回复 | 不主动发公共群 | Aone 中 @责任人，交互请求则回复原会话；Daily nudge 走私信。 |
| Aone 回复触发 Task 唤醒 | 不主动发公共群 | 控制面从 `SUSPENDED` 转 `READY`，生命周期静默；原交互会话可继续显示结果。 |
| 用户在群里主动 @机器人/授权/委派 | 按原会话回复 | 这是请求-响应，不属于主动群消息。 |

实现约束：扫描器、PR-watch、DailyScheduler 和持久化执行器统一使用 routine 日志 sink；只有审批卡片和 `PrWatchScheduler._escalate` 可以主动调用群 broadcast。

# 池级状态集合（pools.json 展开参考）

> **真源是 `config/pools.json`**（各池 `delivery_metrics_status.closed/solution/excluded`、
> `done_status`、`progress_status`）。本文件是常见池的固化交叉参考，pools.json 变更后以此为准；
> 本文件与 pools.json 不一致时，以 pools.json 为准。先读 pools.json，别只信这里。

## tf_customer（project 1086837「Terraform - 客户需求」）

workitem types：需求问题(Req)、功能缺陷(Bug)、线上问题(Bug)、任务(Task)、测试用例执行(Task)。
**不存在"性能瓶颈"类型**——别去查它。

`delivery_metrics_status`：closed=`[验收通过]`、solution=`[已合入主线, 已发布待需求方验收]`、
excluded=`[需求撤回, 已拒绝, 客户未响应]`。统计"关单"时 closed+solution 通常都算已交付，excluded
算"未成交关闭"，按汇报口径选要不要并入。

### Terminal（已关单）集合（跨类型并集）

需求问题(Req)：
`已合入主线` / `已发布待需求方验收` / `验收通过` / `方案功能已存在` / `需求撤回` / `已拒绝` / `客户未响应` / `已取消`

功能缺陷 / 线上问题(Bug)：
`Fixed` / `Closed` / `Won'tfix` / `Later` / `Worksforme` / `Duplicate` / `Invalid` / `External` / `ByDesign`

任务 / 测试用例执行(Task)：
`已完成`

> `closedDate` 对 `Fixed`、`已合入主线`、部分 `验收通过` **为空**——这是坑 #2 的根因。关单时间一律用
> `updateStatusAt`（状态更新于），terminal 集合做判据。

### Active（仍在流转 / 积压）集合（跨类型并集）

`New` / `需求待补充` / `待处理` / `评估中` / `待上游排期` / `问题讨论` / `长期跟进` / `待排期` / `已排期` /
`问题解决中` / `验收中` / `验收不通过` / `Open` / `Reopen`

> 积压统计用**服务端正向过滤**（`--filter "status=评估中,问题解决中,..."`）取活跃全集，别从"all"
> 里 client-side 排除 terminal——all 在大池上本身就被坑 #1 的 25 截断过。

## 新开（created）口径

用 `gmtCreate`（API 字段）/ `create-date`（columns 名）落在窗口内做计数。`gmtModified` 不是创建时间，
别误用。新开口径不受坑 #2 影响（创建时间字段可靠），但仍受坑 #1/#3/#4 影响。

## 跨池通用提醒

- 不同池的 done/active 集合**差异很大**（tf_provider 的 `待发布` 是合法完成态，tf_customer 不是）——
  换池前重读该池的 pools.json，别把 tf_customer 的集合套到别的池。
- `Won'tfix` 含撇号（apostrophe），内联 jq 会断；写进文件用 `jq -f`（见 SKILL.md jq 小抄）。
- 字段值缺失做 `((.field // "") | tostring)` 再 `test`，否则 `test` 作用于 null 报错。

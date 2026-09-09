---
name: aone-stats
description: >-
  Use when counting or tallying Aone work items at the pool/project level — closed/resolved
  counts in a time window, newly-created counts, backlog or 需求积压 by status / product /
  workitemType, or any aggregate statistic over an Aone project (e.g. 1086837 tf_customer).
  Not for single-ticket triage (use aone-triage), GitHub PR review, or personal weekly-report
  attribution (use writing-weekly-report). Triggers: 「统计/数一下/看下积压/最近一周关单/新开了多少/
  backlog」 over a pool, not over one ticket.
---

# Aone 池级统计

> 全流程走 `a1`（非 Terraform 默认 `bin/a1id --`，Terraform 写身份见 AGENTS.md #6）。
> 本技能只读，不改工单，免 claim/bookend。

## 核心原则

**先问"时间信号用哪个字段"，再问"状态集合怎么分"，最后才跑 list。** 池级统计的错误几乎全来自
(1) 分页默认值截断、(2) 关单时间字段选错、(3) list 默认输出不返回字段值。这三步顺序错了，数字一定错。

## 正确查询配方（先照这个跑）

```bash
# 1) 取池内"本周状态变更"全集——大 --page-size + 翻页，绝不依赖默认 25
bin/a1id -- project workitem list \
  --project <池id> \
  --filter "updateStatusAt>=<ISO起>,updateStatusAt<=<ISO止>" \
  --columns id,status,subject,workitemType,identifier,closedDate,status-update-date,modified,create-date \
  --page-size 500 --page 1   # >25 全量需 --page 2,3,... 翻完，或一次给够大 page-size

# 2) 取"活跃积压"全集——服务端正向过滤，别从 all 里 client-filter
bin/a1id -- project workitem list \
  --project <池id> \
  --filter "status=评估中,问题解决中,待处理,..."   # 用各池 active 集合，见 references/status-sets.md
  --columns id,status,subject,workitemType,identifier,create-date \
  --page-size 500 --page 1

# 3) 关单计数 = 步骤1 结果里命中 terminal 集合的条数（按 updateStatusAt 落在窗口内）
#    新开计数 = 同窗口按 gmtCreate/create-date 过滤
#    积压 = 步骤2 全集，再按 product/type/age 二级切分
```

时间窗口的 ISO 起止用 `updateStatusAt`（状态更新于），**不是** `closedDate`（见坑 #2）。

## 五个坑（全部实战踩过，RED 基线）

| # | 错法 | 现象 | 正解 |
|---|---|---|---|
| 1 | 用默认 `--page-size` | `list` 默认 page-size=**25**、`--page` 默认 1，**不自动翻页**；"关单 25"其实只是第 1 页 | 显式 `--page-size 300~500`；超量再 `--page 2,3,...` 翻完，或循环聚拢 |
| 2 | 用 `closedDate` 当关单信号 | `closedDate` 对 Fixed / 已合入主线 / 部分 验收通过 **为空**，漏报严重 | 用 `updateStatusAt`（状态更新于）+ **terminal 状态集合**做关单判据；`closedDate` 仅作交叉参考 |
| 3 | 不带 `--columns` | 默认 list 输出把 `closedDate`/`updateStatusAt` **返回成空**，看着像"都没关" | 显式 `--columns id,status,closedDate,status-update-date,modified,...` 取真实值。filter 服务端照样生效，但值要 columns 才回 |
| 4 | `--filter "identifier=a,b,...,z"` 塞 45 个 | 多值 identifier filter 上限约 25，**只回后 25 个**，前段静默丢失 | 拆成 ≤20 一批，循环取并集；大批量宁可改用 status/updateStatusAt 服务端过滤 |
| 5 | 信 field-options 做"涉及云产品"分类 | 字段（如 140097）的 field-options 查询返回**不全**（漏 RDS/PolarDB/Lindorm/OceanBase/Redis/MongoDB），但工单实际值更丰富（"Redis/云数据库 Tair"） | 用工单**实际字段值 + 标题(subject)** 双正则匹配做产品分类；field-options 只读不可当全集 |

## 状态集合（关单 vs 活跃）

terminal / active 集合**按池 × workitemType 不同**，真源是 `config/pools.json` 各池的
`delivery_metrics_status`（closed/solution/excluded）与 `done_status`/`progress_status`。
**先读 pools.json 拿当前池的集合**，别凭记忆。常见池的固化集合作交叉参考见
`references/status-sets.md`（tf_customer 的 需求问题/功能缺陷/线上问题/任务 各自 terminal 与 active 列表）。

## jq 小抄（值/集合交叉对账）

```bash
# apostrophe（Won'tfix / Won't）会打断内联 jq —— 写进文件用 -f
cat > /tmp/term.jq <<'EOF'
def terminal: ["已合入主线","已发布待需求方验收","验收通过","Fixed","Closed","Won'tfix","Later","Duplicate","Invalid","ByDesign","已完成"];
{ closed: [.[] | select(.status as $s | terminal | index($s))] | length,
  active: [.[] | select(.status as $s | terminal | index($s) | not)] | length }
EOF
jq -f /tmp/term.jq pool_thisweek.json

# 取 id 一行一个（给 comm/comm -12 交集用）——千万别 jq '[.[]|.identifier]|sort'，
# 那产出单行 JSON 数组，comm 找到 0 交集是假象
jq -r '.[].identifier' a.json | sort > /tmp/a.ids
jq -r '.[].identifier' b.json | sort > /tmp/b.ids
comm -12 /tmp/a.ids /tmp/b.ids   # 交集
```

null 安全：字段可能缺失，`test` 作用于 null 会报错，先 `((.["140097"] // "") | tostring)` 兜底再 test。

## 交叉验证（防漏报）

单源为 0 ≠ 没数据。关单数至少两源交叉：(a) `updateStatusAt∈窗口 ∧ status∈terminal` 的 list；
(b) 派发/改派回执或上一轮已知终态集合。两源差值 > 阈值时换源重查，不要辩护旧数字（同
writing-weekly-report §3.5.3 的教训：曾把一周 35 误统计成更少）。

## references

| 文件 | 内容 |
|------|------|
| `references/status-sets.md` | tf_customer 等 common 池 × workitemType 的 terminal/active 状态固化集合（pools.json 的展开参考） |

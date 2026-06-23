# 阶段2 重定义：plan 瘦身 + Claude 分诊 + 五阶段图

> 设计 v0 · 2026-06-23 · 待评审

## 目标
阶段1(scan)已给全集；阶段2 不再重扫，plan.sh 瘦成"去重+按桶汇总(秒级)"，优先级分诊上移给 Claude（按 priority/标题/标签/状态），supervised 门保留。仓库加一张五阶段流程图。

## 已验事实
- list `--columns id,title,status,priority,tag` 可一次取 priority(中/高)+tag，无需逐条 get
- 默认 list 不含 priority/tag → 必须加 columns
- plan.sh 现状：重跑 scan(15s)+逐条 fork seen(O(n²))+置信写死

## 改动
1. **scan.sh**：list 加 `--columns ...priority,tag`，输出加 `priority,tag` 字段
2. **plan.sh 瘦身**：不再调 scan，从 stdin/文件吃 scan 输出；批量去重(runs/ 一次读成 set，jq 一遍 reject)；按 pool+priority 汇总写 runs/plan-日.md；supervised 退2 守门
3. **分诊交 Claude**：loop 文档说明——授权前 Claude 读清单按 优先级/标题/标签/状态 排序+折叠噪声→荐 N 条，置信初判
4. **五阶段图**：docs/flow.md，0自举/1扫描/2去重+分诊门/3逐条triage/4硬门/5收工

## 验收
scan 出 priority+tag；plan 不重扫、批量去重秒级;计划按优先级排;门保留;flow.md 在仓库。

## YAGNI
不做:真置信查证(仍在3c)、自动正式发布、跨池并行调度改写。

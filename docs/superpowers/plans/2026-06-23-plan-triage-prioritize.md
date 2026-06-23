# Plan Triage + Prioritize Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: superpowers:subagent-driven-development. Steps use checkbox syntax.

**Goal:** scan 带 priority/tag；plan 不重扫、批量去重秒级、按优先级汇总；分诊交 Claude；加五阶段图。

**Architecture:** scan 加 --columns；plan 从 stdin 吃 scan、一次读 runs/ 去重；loop 文档说明 Claude 分诊；docs/flow.md。

**Tech Stack:** bash, jq, a1。

## Global Constraints
- 阶段1 全集已出，阶段2 不重扫
- plan 批量去重(runs/ 一次读),no per-item fork
- supervised 退2 守门不变；release_prod 永停
- 真置信查证仍在阶段3，plan 只机械汇总

---

### Task 1: scan 输出加 priority+tag
**Files:** Modify bootstrap/scan.sh, test/scan_test.sh
- [ ] S1 测：stub 返回含 priority/tag → 输出对象含这两字段；列含 priority,tag
- [ ] S2 FAIL；S3 list 加 `--columns id,title,status,priority,tag`,jq map 加 priority/tag;S4 旧+新过;**Commit** `feat: scan emits priority+tag`

### Task 2: plan 吃 scan 输入,批量去重
**Files:** Modify bootstrap/plan.sh, test/plan_test.sh
- [ ] S1 测:plan 读 stdin JSON(不调 scan);去重一次读 runs/;1000 条无 per-item fork
- [ ] S2 FAIL;S3 plan 从 stdin/文件取 items(无 scan.sh 调用),`ls runs/`一次成 set,jq reject 已见,排表加 priority 列;退2;S4 PASS;**Commit** `perf: plan consumes scan stdin + batch dedup`

### Task 3: loop 文档分诊说明
**Files:** Modify loops/aone-triage.md
- [ ] S1 阶段3前加:scan→plan去重→Claude 按 priority/标题/标签/状态 排序折叠荐N→授权;**Commit** `docs: claude prioritizes before auth`

### Task 4: 五阶段流程图
**Files:** Create docs/flow.md
- [ ] S1 写 0自举/1扫描/2去重+分诊门/3triage/4硬门/5收工,标注人工点;S2 commit `docs: five-stage flow`

## 验收
scan 出 priority+tag;plan 不重扫秒级;分诊在 Claude;flow.md 在仓库。

# Aone Pools Config Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** 把 Aone 项目事实抽成 `config/pools.json` 单一数据源，技能/scan 引用它，并加逐池视检 pools.sh。

**Architecture:** JSON 数据（jq 直读，不加 yq）；技能 prose 去 ID 指 config；scan 读 claim.tag；verify+pools.sh 用 config。

**Tech Stack:** bash, jq, a1 CLI。

## Global Constraints
- ID 全部 verbatim 抽自现有技能文件，不臆造
- JSON，不引入 yq
- 唯一硬门=正式发布；scan 跳过 `NOT tag=jarvis-claimed`
- 技能保留流程+坑，只移数据

---

### Task 1: config/pools.json
**Files:** Create: `config/pools.json`; Test: `test/pools_config_test.sh`
- [ ] **S1 测**：jq 解析通过；`.pools|length>=3`；含 2100304/2124589/2165097、claim.tag=jarvis-claimed
- [ ] **S2 FAIL**
- [ ] **S3** 写 pools.json：terraform(2100304/cfs107239=906688/320687)、agent_portal(2124589/app283346/66·67/320687/cfs100340/delivery)、cloudspec_gap(2165097/479782)、routing、claim{tag,fallback:title}
- [ ] **S4 PASS**；**Commit** `feat: pools.json data source`

### Task 2: scan.sh 读 claim.tag 跳过已认领
**Files:** Modify: `bootstrap/scan.sh`; `test/scan_test.sh`
- [ ] **S1 测**：stub list, 加 `--filter NOT tag=<config.tag>`；config 缺则不过滤
- [ ] **S2 FAIL**；**S3** scan 读 `jq -r .claim.tag config/pools.json` 拼 filter；**S4 PASS**；**Commit** `feat: scan skips claimed`

### Task 3: pools.sh 逐池视检
**Files:** Create: `bootstrap/pools.sh`; Test: `test/pools_test.sh`
- [ ] **S1 测**：stub a1, 2 池→输出每池 `name project open=N`，退0
- [ ] **S2 FAIL**；**S3** 读 pools.json 每池 `a1 list --project -f json|jq length`；**S4 PASS**；**Commit** `feat: per-pool inspection`

### Task 4: 技能去 ID 指 config
**Files:** Modify: `.claude/skills/aone-triage/references/routing.md`, `SKILL.md`, `delivery-aliyun-automation-agent.md`
- [ ] **S1** 三文件硬 ID 行加“见 config/pools.json”，保留流程/坑/cfs 用法
- [ ] **S2 验** `bash bootstrap/verify.sh` exit0；**Commit** `docs: skill points to pools.json`

## 验收
改池=改1行；scan 跳认领；pools.sh 列各池 open；技能无散落 ID。

# Claim Concurrency Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: superpowers:subagent-driven-development. Steps use checkbox syntax.

**Goal:** 池级 scan + tag 认领(带时间戳)+ 僵尸清扫，让多机并行不抢不重、挂掉可回收。

**Architecture:** config 加 done_tag/ttl；claim.sh 认领/释放；scan.sh 改逐池 NOT-tag 合并；sweep.sh 超时僵尸入 escalation/。

**Tech Stack:** bash, jq, a1。

## Global Constraints
- tag 过滤须带 --project；无 untag → release=retag jarvis-done
- 认领=tag jarvis-claimed + 评论 `jarvis-claim <host> <UTC>`；回读确认
- sweep 不自动回收，僵尸入 escalation/ 待人；TTL 默认 45min
- 读多写少；release 唯一硬门 release_prod 不受影响

---

### Task 1: config 加认领字段
**Files:** Modify `config/pools.json`, `test/pools_config_test.sh`
- [ ] S1 测：`.claim.done_tag=="jarvis-done"` 且 `.claim.ttl_min==45`
- [ ] S2 改 config claim 块加 done_tag/ttl_min；S3 测过；**Commit** `feat: claim done_tag+ttl config`

### Task 2: claim.sh 认领/释放
**Files:** Create `bootstrap/claim.sh`, `test/claim_test.sh`
- [ ] S1 测：stub a1，`claim 1 P`→调 update --tag jarvis-claimed + comment create 含 jarvis-claim；`release 1 P`→update --tag jarvis-done;回读非己→非0
- [ ] S2 FAIL；S3 claim.sh(claim/release, host=$(hostname) UTC, 回读校验)；S4 PASS；**Commit** `feat: claim.sh`

### Task 3: scan.sh 改池级合并
**Files:** Modify `bootstrap/scan.sh`, `test/scan_test.sh`
- [ ] S1 测：stub 2 池各返 1 项→输出合并 2 项；--project+--filter NOT tag=jarvis-claimed 传入
- [ ] S2 FAIL；S3 循环 pools[].project，list --project --filter，jq -s add 合并 {id,title,type,status,pool}；空池跳；S4 旧+新测过；**Commit** `feat: pool-scoped scan with claim filter`

### Task 4: sweep.sh 僵尸清扫
**Files:** Create `bootstrap/sweep.sh`, `test/sweep_test.sh`
- [ ] S1 测：stub 一池含 claimed 项+老评论(>ttl)→escalation/ 有该 id；新评论→不入
- [ ] S2 FAIL；S3 逐池 list --tag jarvis-claimed，comment list 取 jarvis-claim UTC，超 ttl→escalate；S4 PASS；**Commit** `feat: sweep stale claims`

### Task 5: 接线 verify+loop
**Files:** Modify `bootstrap/verify.sh`, `loops/aone-triage.md`
- [ ] S1 verify 加 done_tag 校验；loop 串 claim→done+sweep；S2 verify exit0；**Commit** `feat: wire claim/sweep into loop`

## 验收
2 clone 不抢；领→done 闭环；僵尸 45min 入队;scan 跳已领。

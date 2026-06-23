# jarvis Triage Run-Loop Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让 jarvis 无人值守扫 Aone 入箱、逐项 triage、全链跑到预发/CR、低置信入 escalation/、全程写 runs/，唯一硬门=正式发布。

**Architecture:** 不重写 aone-triage（vendored 技能负责单条工单全流程）；本计划建"驱动+护栏"：scan 拉清单 → 逐项交技能 → autonomy.md 决定自动/停 → runs/审计、escalation/停队。bash 脚本做幂等/可测的扫描与记账，决策权与 loop 用文档约束。

**Tech Stack:** bash, a1 CLI(`a1 project workitem`), jq, vendored aone-triage skill。

## Global Constraints

- 默认 supervised：跑批先出执行计划→人逐条判可行/允许/授权→确认后才执行；未确认不写 Aone
- 唯一硬门=正式发布；授权后预发/CR 以下自动；低置信/验收不过→起草不发出，入 escalation/
- 档位 unattended 仅在用户显式切换后启用（信任建立后）
- 红线：代码改动必新分支+CR，禁直 push master，禁零 diff 空 CR
- 真源=Aone；本地 runs//escalation/ 可重建，不作真相
- 全程 a1 CLI；身份 `a1 auth whoami`；cloudspec 无需凭证
- 每动作幂等：同一 workitemId 不重复处理（runs/ 去重）

---

### Task 1: autonomy.md 决策权策略

**Files:** Modify: `autonomy.md`

**Interfaces:** Produces: 一份机读+人读策略，loop 与技能据此判自动 vs 停。

- [ ] **Step 1:** 写两档：`supervised`(默认,出计划等授权再做) / `unattended`(只 escalate 才找人,需显式切)
- [ ] **Step 2:** 写自动档：授权后高置信+可逆(回复/建需求/打标/建 CR/worktree/预发)→自动；正式发布永停
- [ ] **Step 3:** 写置信判定：查证两层(OpenAPI+源码)一致=高；冲突/缺源码/规则未命中=低
- [ ] **Step 4:** 末尾机读 ```json：`{"mode":"supervised","auto":[...],"stop":["release_prod"],"escalate_if":["low_conf","verify_fail","redline"]}`
- [ ] **Step 5: 验** `grep -q '"mode":"supervised"' autonomy.md && grep -q release_prod autonomy.md`；**Commit** `feat: autonomy policy + supervised gate`

---

### Task 2: scan 脚本（入箱→items.json）

**Files:** Create: `bootstrap/scan.sh`; Test: `test/scan_test.sh`

**Interfaces:** Produces: `scan.sh` → stdout JSON 数组 `[{id,title,type,status}]`；失败退非零。

- [ ] **Step 1: 测** stub `a1` 返固定 json，跑 scan.sh 断言输出是数组且含 id；空入箱→`[]` exit0
- [ ] **Step 2: FAIL**（无脚本）
- [ ] **Step 3:** scan.sh：`a1 project workitem list --assignee me --status open -f json` | jq 规整为 `[{id,title,type,status}]`；`set -u`
- [ ] **Step 4: PASS**；real run 打印当前入箱
- [ ] **Step 5: Commit** `feat: inbox scan.sh`

---

### Task 3: runs/ 审计 + escalation/ 去重记账

**Files:** Create: `bootstrap/log.sh`; Test: `test/log_test.sh`

**Interfaces:** Produces: `run_done <id> <summary>` 写 `runs/<date>-<id>.md`；`escalate <id> <reason>` 写 `escalation/<id>.md`；`seen <id>` exit0 if runs/ 已有该 id。

- [ ] **Step 1: 测** `seen X`→非0；`run_done X ok`→runs/有文件且 `seen X`→0；`escalate Y low`→escalation/有文件
- [ ] **Step 2: FAIL**
- [ ] **Step 3:** log.sh 三函数，文件名含 id，repo_root 相对
- [ ] **Step 4: PASS**；**Commit** `feat: runs+escalation ledger`

---

### Task 4: loops/aone-triage.md 跑链 runbook

**Files:** Modify: `loops/aone-triage.md`

**Interfaces:** Consumes: scan.sh/log.sh/autonomy.md/技能。

- [ ] **Step 1:** 写：触发→scan.sh→plan.sh 出计划→supervised 等用户逐条授权(可行/允许/授权)→授权项逐条 技能 triage→autonomy 判→自动到预发 or escalate→run_done；done(到预发或入队)/仅人工(正式发布)
- [ ] **Step 2: 验** `grep -q plan.sh loops/aone-triage.md && grep -q 授权 loops/aone-triage.md && grep -q release loops/aone-triage.md`；**Commit** `feat: triage loop runbook`

---

### Task 5: 执行计划 + 授权门（plan.sh）

**Files:** Create: `bootstrap/plan.sh`; Test: `test/plan_test.sh`

**Interfaces:** Consumes: scan.sh items + log.sh seen。Produces: `plan.sh` 把入箱(去重后)出一份执行计划到 `runs/plan-<date>.md`：每条 = id/标题/拟动作/置信/自动or停/不可逆点；supervised 档下打印"待授权:逐条 可行?允许?授权?"并退码 2(待确认)，绝不直接写 Aone。

- [ ] **Step 1: 测** stub 2 items(1 高1 低)，跑 plan.sh：runs/plan-*.md 含两条、低置信标 escalate、退码 2；seen 的 id 不出现
- [ ] **Step 2: FAIL**
- [ ] **Step 3:** plan.sh：scan→过滤 seen→逐条排版动作+置信+auto/stop→supervised 退 2、unattended 退 0；不调写命令
- [ ] **Step 4: PASS**；**Commit** `feat: plan.sh execution-plan + supervised auth gate`

---

### Task 6: 收尾接线（CLAUDE.md 指 loop、README、定时说明）

**Files:** Modify: `CLAUDE.md`, `README.md`; Create: `bootstrap/cron.example`

- [ ] **Step 1:** CLAUDE.md 开局 2) 改为"跑 scan.sh→plan.sh→出计划等授权→按 loop 处理"；cron.example 给定时启 claude 跑 loop 样例（注明默认 supervised）
- [ ] **Step 2: 验** `bash bootstrap/verify.sh` 仍 exit0；**Commit** `feat: wire triage loop into bootstrap`

## 验收标准（怎么确认 plan2 达成）

端到端 demo，逐条可演示：
1. **不用喂单**：`scan.sh` 拉出当前 Aone 入箱真清单（≥1 条）。
2. **计划先行**：`plan.sh` 出执行计划，每条带 动作/置信/自动or停/不可逆点；supervised 下退码 2、**全程零 Aone 写**。
3. **去重**：已处理的 id 二次跑被跳过（`seen` 命中），不重复处理。
4. **授权后到预发**：授权一条 → 全链跑到预发/CR，runs/ 有该条审计；正式发布**未触发**。
5. **低置信入队**：一条查证冲突/缺源码 → 不发出，escalation/ 有该条+原因。
6. **硬门**：全程无任何正式发布动作。

全 6 条过 = plan2 事实达成、符合设计预期。映射：1=不用喂单；2/4=计划先行+授权才动；3/5/6=到预发即停、低置信入队、不漏不重、硬门守住。

## YAGNI
不做：多 app 交付、正式发布自动化、MCP 接入。先 triage→预发一条全绿。

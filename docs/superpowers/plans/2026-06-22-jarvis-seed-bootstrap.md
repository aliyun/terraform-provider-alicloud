# jarvis Seed Bootstrap Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让任何 Claude `git clone jarvis` + 注入凭证 + 一键装/验，即得到一个自包含、可干活的起点母版。

**Architecture:** `CLAUDE.md` 作自举入口（自动入上下文，`@import` 身份/清单/决策权）；`bootstrap/` 幂等装依赖、逐项验；技能 vendored 到 `skills/`。本计划只做 P0 自举 + P1 依赖骨架，不含 triage 跑链（plan #2）。

**Tech Stack:** bash, curl, git, a1/aliyun/cloudspec/gh CLI, Claude Code 技能体系。

## Global Constraints

- 主入口必须叫 `CLAUDE.md`（任意 IDENTITY.md 不会被自动读）
- `verify.sh`：每依赖一条独立 check，单独 PASS/FAIL，任一 FAIL 退非零，绝不聚合"全 ok"
- 每条 install 命令旁必带文档链接（脚本过时报"见 <doc>"）
- 凭证各验各的，不互相代替；CR 未合即可逆，唯一硬门=正式发布
- 不重写技能、不做 Web UI、先只支撑 triage 一条 loop

---

### Task 1: 自举入口 CLAUDE.md + 骨架目录

**Files:**
- Create: `CLAUDE.md`, `loops/.gitkeep`, `autonomy.md`, `runs/.gitkeep`, `escalation/.gitkeep`, `bootstrap/.gitkeep`, `skills/.gitkeep`

**Interfaces:**
- Produces: `CLAUDE.md` 含 `@import autonomy.md`、`@import loops/aone-triage.md`；后续任务填这些目标。

- [ ] **Step 1: 写 CLAUDE.md** —— 你是谁/对谁负责/开局动作三段，末尾 `@autonomy.md` `@loops/aone-triage.md`
- [ ] **Step 2: 建空目录占位** —— `loops runs escalation bootstrap skills` 各放 `.gitkeep`
- [ ] **Step 3: 验** —— `test -f CLAUDE.md && grep -q '@autonomy.md' CLAUDE.md` Expected: exit 0
- [ ] **Step 4: Commit** —— `git add -A && git commit -m "feat: bootstrap CLAUDE.md entry + skeleton"`

---

### Task 2: deps.lock + install.sh（CLI 安装，带文档链接）

**Files:**
- Create: `bootstrap/deps.lock`, `bootstrap/install.sh`

**Interfaces:**
- Consumes: 内部 CLI 安装源（见 spec）
- Produces: `install.sh` 幂等，已装则跳过，失败提示 `脚本可能过时，见 <doc>`

- [ ] **Step 1: deps.lock** —— a1/aliyun/cloudspec/gh + install cmd + doc url
- [ ] **Step 2: install.sh** —— `command -v` 判已装；a1 `curl ...aone-cli/install.sh|sh`、aliyun `aliyuncli.alicdn.com/install.sh`、cloudspec `acube...|sudo bash`，失败 echo 文档
- [ ] **Step 3: 验幂等** —— 二次跑全 "skip" Expected: exit 0
- [ ] **Step 4: Commit**

---

### Task 3: verify.sh 逐依赖独立 check

**Files:**
- Create: `bootstrap/verify.sh`

**Interfaces:**
- Produces: 输出 `PASS a1 / FAIL aliyun`，任一 FAIL 退非零

- [ ] **Step 1: 写测试** —— stub PATH 缺 a1，跑 verify 断言含 `FAIL a1` 且 `echo $?`≠0
- [ ] **Step 2: 跑测试 FAIL**
- [ ] **Step 3: verify.sh** —— 逐项 `cmd --version` 一行一 check，函数 `chk name cmd`
- [ ] **Step 4: 跑测试 PASS**
- [ ] **Step 5: Commit**

---

### Task 4: .env.example + 凭证各自验

**Files:** Create: `bootstrap/.env.example`; Modify: `bootstrap/verify.sh`

- [ ] **Step 1: .env.example** —— GH_TOKEN/a1/aliyun key 模板占位
- [ ] **Step 2: 加凭证 check** —— gh `gh auth status`、a1 调一次、aliyun `sts` 各独立 PASS/FAIL
- [ ] **Step 3: 验** —— 缺凭证→对应 FAIL 退非零
- [ ] **Step 4: Commit**

---

### Task 5: vendor aone-triage 进 skills/

**Files:** Create: `skills/aone-triage/**`; Modify: `verify.sh`

- [ ] **Step 1: 复制技能** —— 从 plugins cache 拷 aone-triage 全量入 `skills/aone-triage`
- [ ] **Step 2: verify 加技能 check** —— `test -f skills/aone-triage/SKILL.md`
- [ ] **Step 3: 验全绿** Expected: 全 PASS
- [ ] **Step 4: Commit**

---

### Task 6: README 补自举说明

**Files:** Modify: `README.md`

- [ ] **Step 1:** 写 clone→注 .env→`install.sh`→`verify.sh` 全绿→开局
- [ ] **Step 2: Commit**

## YAGNI
不含：triage 跑链、定时、MCP 装（plan #2）。

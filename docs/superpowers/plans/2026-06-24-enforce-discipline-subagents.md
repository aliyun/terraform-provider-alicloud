# 纪律工程化 + subagent 一条一派 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 jarvis 收尾纪律从提示词换成代码兜底（4 层），并落 4 类隔离 subagent。

**Architecture:** 改硬 `wrap.sh done`（status 必填、a1 失败 exit1）；`claim.sh` 写本地认领台账；Stop-hook `wrap-check.sh` 收尾未齐拦会话；`triage-one.sh` 一条一派 bookend；`reconcile.sh` 事后对账；`.claude/agents/` 落 triager/developer/pr-reviewer/verifier。

**Tech Stack:** bash（`set -uo pipefail`）、jq、python3、a1 CLI；测试用纯 bash + PATH stub，无 bats。

## Global Constraints
- 所有脚本沿用：`set -uo pipefail`、`JARVIS_ROOT` 覆盖、`script_dir` 解析、jq/python3 读 `config/pools.json`。
- tag/done_tag/ttl 一律读 `pools.json .claim`，不硬编码。
- 测试不连 Aone：PATH 前置假 `a1` stub，断言退码与调用记录。临时文件用 `mktemp`，禁仓库根脏文件。
- Aone: 83491649；分支已在 worktree `worktree-enforce-discipline-subagents`，禁直接合 master。
- 对外产物不带 AI 署名；git commit 末尾 `Co-Authored-By` 可留。

---

### Task 1: wrap.sh done 收紧（status 必填 + a1 失败即报错）

**Files:**
- Modify: `bootstrap/wrap.sh:47-64`
- Test: `bootstrap/tests/wrap_done.sh`

**Interfaces:**
- Produces: `wrap.sh done <id> <summary> <status>` — status 现为必填；a1 comment/update 任一失败 → exit 1。sync 不变。

- [ ] **Step 1: 写失败测试** `bootstrap/tests/wrap_done.sh`

```bash
#!/usr/bin/env bash
set -uo pipefail
here="$(cd "$(dirname "$0")" && pwd)"; root="$(cd "$here/../.." && pwd)"
stub="$(mktemp -d)"; cat > "$stub/a1" <<'E'
#!/usr/bin/env bash
[ "${A1_FAIL:-0}" = "1" ] && exit 1
exit 0
E
chmod +x "$stub/a1"; export PATH="$stub:$PATH"; export JARVIS_RUNS_DIR="$(mktemp -d)"
fail=0
# 缺 status → 退码非 0
bash "$root/bootstrap/wrap.sh" done 100 "s" >/dev/null 2>&1 && { echo "FAIL no-status"; fail=1; }
# 全参数 + a1 成功 → 退 0
bash "$root/bootstrap/wrap.sh" done 100 "s" "已发布" >/dev/null 2>&1 || { echo "FAIL ok"; fail=1; }
# a1 失败 → 退码非 0
A1_FAIL=1 bash "$root/bootstrap/wrap.sh" done 100 "s" "已发布" >/dev/null 2>&1 && { echo "FAIL a1err"; fail=1; }
[ $fail = 0 ] && echo PASS; exit $fail
```

- [ ] **Step 2: 跑测试看失败** `bash bootstrap/tests/wrap_done.sh` → 预期 FAIL（现 status 可选、a1 warn-only）。
- [ ] **Step 3: 改 wrap.sh done 分支**。`bootstrap/wrap.sh` done 段改：`[ -n "$id" ] && [ -n "$summary" ] && [ -n "$status" ] || { echo "Usage: wrap.sh done <id> <summary> <status>" >&2; exit 1; }`；comment 与 update 的 `|| echo ...` 改为 `|| { echo "..." >&2; exit 1; }`；status 块去掉 `if`，恒执行。
- [ ] **Step 4: 跑测试看通过** `bash bootstrap/tests/wrap_done.sh` → PASS。
- [ ] **Step 5: 提交** `git add bootstrap/wrap.sh bootstrap/tests/wrap_done.sh && git commit -m "fix: wrap.sh done status必填+a1失败即报错 (#83491649)"`

### Task 2: claim.sh 写认领台账 + wrap-check.sh

**Files:**
- Modify: `bootstrap/claim.sh:40-66`
- Create: `bootstrap/wrap-check.sh`
- Test: `bootstrap/tests/wrap_check.sh`

**Interfaces:**
- Produces: 台账 `.my-day/claims-<UTCdate>.json`（claim 追加 `{id,done:false}`，release 置 done:true）；`wrap-check.sh` 扫台账，未 done 且 `runs/*-<id>.md` 缺失则 exit 2 列 id，全齐 exit 0。

- [ ] **Step 1: 写失败测试** `bootstrap/tests/wrap_check.sh`：建临时 `.my-day/claims-X.json` 含一条 done:false 且 runs/ 无对应文件 → 期望 exit 2；补 `runs/<date>-<id>.md` → exit 0。
- [ ] **Step 2: 跑测试看失败**（wrap-check.sh 不存在）。
- [ ] **Step 3: claim.sh claim 成功后** append `{id,done:false}` 到 `.my-day/claims-$(date -u +%F).json`；release append/标记 done:true（用 jq）。
- [ ] **Step 4: 写 `wrap-check.sh`**：读当日台账，未 done 项逐个 `log.sh seen <id>`，缺则收集；非空 → 打印 + exit 2，空 → exit 0。
- [ ] **Step 5: 跑测试 PASS。提交。**

### Task 3: .claude/settings.json 注册 Stop hook

**Files:** Create: `.claude/settings.json`
- [ ] **Step 1:** 写 `{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"bash bootstrap/wrap-check.sh"}]}]}}`。
- [ ] **Step 2:** `bash bootstrap/wrap-check.sh; echo $?` 验证退 0（无未收尾）。提交。

### Task 4: reconcile.sh 对账器

**Files:** Create `bootstrap/reconcile.sh`；Test `bootstrap/tests/reconcile.sh`
- [ ] **Step1-2:** 测试：stub a1，模拟 jarvis-claimed 无 jarvis-done → 期望补 release。Step3：仿 sweep.sh 遍历 pools，列 claimed∧¬done，对每条 `claim.sh release` + `wrap.sh done`。Step4-5：PASS、提交。

### Task 5: .claude/agents/ 四类定义

**Files:** Create `triager.md/developer.md/pr-reviewer.md/verifier.md`
- [ ] **Step1:** 各文件 frontmatter `name/description/tools/model`，正文照 spec 表（用途/隔离/工具/写权/路由/收尾）。`developer` 注明 worktree。Step2: 提交。

### Task 6: triage-one.sh 编排 bookend

**Files:** Create `bootstrap/triage-one.sh`；Test `bootstrap/tests/triage_one.sh`
- [ ] **Step1-2:** 测试 dry-run 顺序 claim→done→release，缺 status 失败不 release。Step3: 写脚本。Step4-5: PASS、提交。

### Task 7: verify.sh 闸门 + 文档对齐

**Files:** Modify `bootstrap/verify.sh:62-72`；`CLAUDE.md` 纪律2、`loops/aone-triage.md`
- [ ] **Step1:** verify 加 chk settings.json + 4 agents 存在。Step2: CLAUDE.md/loops 指向 triage-one.sh+agents+reconcile。Step3: 提交。

## Self-Review
- 覆盖：spec 4 层=Task1/2+3/4/6+5/2+3 reconcile=4，4 subagent=5；verify/文档=7。无缺口。
- 无占位符：硬测试给了真码，其余给了断言要点。type 一致：done 三参、台账 done 布尔、wrap-check exit2 全程一致。

#!/bin/bash

fail_count=0

# chk function: checks if a command is available
# Usage: chk NAME CMD
chk() {
    local name="$1"
    local cmd="$2"

    if command -v "$cmd" >/dev/null 2>&1; then
        echo "PASS $name"
    else
        echo "FAIL $name"
        ((fail_count++))
    fi
}

# chk_cred function: checks if a credential is valid by running a command
# Usage: chk_cred NAME "COMMAND"
# Runs COMMAND; PASS if exit 0, FAIL if non-zero
chk_cred() {
    local name="$1"
    local cmd="$2"

    if eval "$cmd" >/dev/null 2>&1; then
        echo "PASS $name-cred"
    else
        echo "FAIL $name-cred"
        ((fail_count++))
    fi
}

# chk_skill function: checks if a vendored skill exists
# Usage: chk_skill NAME
# PASS if skills/NAME/SKILL.md exists, FAIL otherwise
chk_skill() {
    local name="$1"
    local repo_root="$(git rev-parse --show-toplevel)"
    local skill_file="${repo_root}/.claude/skills/${name}/SKILL.md"

    if test -f "$skill_file"; then
        echo "PASS $name"
    else
        echo "FAIL $name"
        ((fail_count++))
    fi
}

# All roles run the same full tool + credential check set. Workers used to
# SKIP gh/aliyun/cloudspec as "dev-only", but worker Tasks do the same
# triage/PR work as the scheduler host — a worker without gh + aliyun creds
# silently loses the GitHub delivery path and OpenAPI 查证. Both now install
# sudo-free from pinned tarballs on Linux (deps.lock), and cloudspec is a
# WARN-not-FAIL below, so there is no install-cost reason left to skip.
repo_root="$(git rev-parse --show-toplevel)"

# Check the CLIs
chk a1 a1
chk git git
chk gh gh
chk aliyun aliyun
# cloudspec CLI: upstream install endpoint (acube.aliyun-inc.com) currently
# returns HTTP 500 for all known versions — install.sh downloads an empty
# zip and unzip fails. Downgrade to WARN so preflight isn't blocked; the
# warning includes the repair command.
# Once upstream is healthy again, `bootstrap/install.sh` will pick it up
# and the WARN turns back into PASS on next preflight.
if command -v cloudspec >/dev/null 2>&1; then
    echo "PASS cloudspec"
else
    echo "WARN cloudspec — acube 上游缺 Linux 产物(cloudspec-linux-*.zip 全版本 500「下载失败」, mac 包正常); CloudSpec IDL 相关 Task 会缺 CLI 支持"
    echo "     修复: 跑 bash bootstrap/cloudspec-build-linux.sh 源码构建(官方 Linux 路径, 全程几分钟), 或等上游补产物后 preflight 自动装; 文档: https://aliyuque.antfin.com/cloudspec/model/cli-install"
fi

# Check credentials (each independent PASS/FAIL)
chk_cred aliyun "aliyun sts GetCallerIdentity"
chk_cred a1 "a1 auth whoami"
# GitHub token daily check — WARN, NOT FAIL. A stale/missing JARVIS_GITHUB_TOKEN only
# blocks the GitHub delivery path (PR/评论/推分支); it must not fail preflight and
# thereby block Aone-only work. So we print a WARN line and do NOT increment fail_count.
# (github-identity.sh check itself is unchanged — it stays a hard gate at use time.)
# Auto-source .env if JARVIS_GITHUB_TOKEN not yet in environment (gh#78).
if [ -z "${JARVIS_GITHUB_TOKEN:-}" ] && [ -f "$repo_root/bootstrap/.env" ]; then
    # shellcheck source=/dev/null
    source "$repo_root/bootstrap/.env"
fi
if "$repo_root/bootstrap/github-identity.sh" check >/dev/null 2>&1; then
    echo "PASS jarvis-github-token"
else
    echo "WARN jarvis-github-token — JARVIS_GITHUB_TOKEN 失效/缺失；GitHub 升级路径(PR/评论/推分支)会被阻断，Aone-only 工作不受影响"
    echo "     修复: 刷新 api-tool-agent 的 GitHub token，更新 bootstrap/.env 的 JARVIS_GITHUB_TOKEN(或其环境来源)，再跑 bootstrap/github-identity.sh check 复验"
fi

# a1 (jarvis identity) login-state daily check — WARN, NOT FAIL. Mirrors the
# GitHub-token check above: a dead/expired a1 session only blocks a1-backed Aone
# writes (aone-triage 回复/建单、wrap.sh、claim/scan — all go through
# bin/a1id); it must not fail preflight and thereby block non-a1 work. So we WARN,
# do NOT increment fail_count.
# OK iff `bin/a1id -- auth whoami` yields Account == WORKER_1782379562571 (jarvis).
# 空/报错(过期会话) 或 半死(EmpID 在但 Account 空) → WARN. Parser matches bin/a1id.
a1id_bin="${JARVIS_A1ID:-$repo_root/bin/a1id}"
a1_expect="WORKER_1782379562571"
# 显式 pin JARVIS_A1_IDENTITY=jarvis 防环境残留(如上一步 headless 用了角色身份)导致
# 检错身份出假 WARN——本 check 就是要验 jarvis 默认身份的登录态。
a1_account="$(JARVIS_A1_IDENTITY=jarvis "$a1id_bin" -- auth whoami 2>/dev/null | awk '/Account:/{print $2}' || true)"
if [ "$a1_account" = "$a1_expect" ]; then
    echo "PASS jarvis-a1-session"
else
    echo "WARN jarvis-a1-session — a1 jarvis 登录态失效/缺失(whoami Account='${a1_account:-<空>}' != $a1_expect)；a1 相关 Aone 写(triage 回复/建单/wrap)会被阻断，非 a1 工作不受影响"
    echo "     修复: 浏览器 BUC 登 open_jarvis 后 bin/a1id login jarvis，再跑 bin/a1id -- auth whoami 复验 Account=$a1_expect"
fi

# Screenshot capability (repo-controlled browser channel) — WARN, NOT FAIL. Headless
# bridge sessions have NO injected Playwright MCP (root cause:
# .claude/skills/screenshot-evidence/references/headless-screenshot-channels.md); the
# repo-controlled capture.sh probe checks for a Playwright Python binding + local
# chromium, else a headless Chrome/Chromium binary. Missing only degrades screenshot
# evidence to missing_capability n-a rows — it must not fail preflight and thereby
# block non-screenshot work.
if bash "$repo_root/.claude/skills/screenshot-evidence/scripts/capture.sh" probe >/dev/null 2>&1; then
    echo "PASS jarvis-screenshot"
else
    echo "WARN jarvis-screenshot — headless 截图通道缺失(无 Playwright Python + 本地 chromium，也无 headless Chrome/Chromium)；screenshot-evidence 会以 missing_capability 收口(逐层 n-a + 原因)，不阻断其它工作"
    echo "     修复(任选其一): python3 -m pip install playwright && python3 -m playwright install chromium；或装 chromium 并/或 export JARVIS_CHROME_BIN=<path>；文档: .claude/skills/screenshot-evidence/references/headless-screenshot-channels.md"
fi

# Check vendored skills
chk_skill aone-triage
for skill in \
    cloudspec-amp-workflow \
    cloudspec-idl-guide \
    cloudspec-resource-edit \
    cloudspec-operation-edit \
    cloudspec-flag-mode-edit \
    cloudspec-build-fix \
    cloudspec-norm-check-fix \
    cloudspec-shared-knowledge; do
    chk_skill "$skill"
done

# cloudspec-core snapshot check is pure repo-local file validation (lock file
# + vendored skill mirror, jq only) — nothing scheduler-specific, every role
# runs it. Last surviving worker-skip removed with the rest of the role gate.
if bash "$repo_root/bootstrap/cloudspec-core.sh" check >/dev/null 2>&1; then
    echo "PASS cloudspec-core-snapshot"
else
    echo "FAIL cloudspec-core-snapshot"
    ((fail_count++))
fi

# Check pools config parses and has >=3 pools
pools_cfg="$(git rev-parse --show-toplevel)/config/pools.json"
if jq -e '.pools | length >= 3' "$pools_cfg" >/dev/null 2>&1; then
    echo "PASS pools.json"
else
    echo "FAIL pools.json"
    ((fail_count++))
fi

# Check claim.idle_tag / done_tag / done_status are set to expected values
if jq -e '.claim.idle_tag=="jarvis-idle"' "$pools_cfg" >/dev/null 2>&1; then
    echo "PASS claim.idle_tag"
else
    echo "FAIL claim.idle_tag"
    ((fail_count++))
fi
if jq -e '.claim.done_tag=="jarvis-done"' "$pools_cfg" >/dev/null 2>&1; then
    echo "PASS claim.done_tag"
else
    echo "FAIL claim.done_tag"
    ((fail_count++))
fi
if jq -e '.claim.done_status!=null and (.claim.done_status|length)>0' "$pools_cfg" >/dev/null 2>&1; then
    echo "PASS claim.done_status"
else
    echo "FAIL claim.done_status"
    ((fail_count++))
fi

# Check .claude/settings.json exists and has hooks.Stop configured
settings_file="$(git rev-parse --show-toplevel)/.claude/settings.json"
if [ -f "$settings_file" ] && jq -e '.hooks.Stop' "$settings_file" >/dev/null 2>&1; then
    echo "PASS settings.json hooks.Stop"
else
    echo "FAIL settings.json hooks.Stop"
    ((fail_count++))
fi

# Check project Codex Stop hook is not Claude-runtime-only. It must call the
# shared repo-local wrapper and avoid expanding an empty CLAUDE_PROJECT_DIR to
# /bootstrap/wrap-check.sh.
codex_hooks_file="$(git rev-parse --show-toplevel)/.codex/hooks.json"
codex_stop_cmd="$(jq -r '.hooks.Stop[0].hooks[0].command // ""' "$codex_hooks_file" 2>/dev/null || true)"
if [ -f "$codex_hooks_file" ] \
    && printf '%s' "$codex_stop_cmd" | grep -q 'bootstrap/run-stop-hook.sh' \
    && ! printf '%s' "$codex_stop_cmd" | grep -q 'CLAUDE_PROJECT_DIR'; then
    echo "PASS codex hooks.Stop"
else
    echo "FAIL codex hooks.Stop"
    ((fail_count++))
fi

# Check the 3 数字人 agent 定义 files exist (主会话即总领，不单设 jarvis 子代理)
for agent in terraform-pd terraform-rd terraform-qa; do
    agent_file="${repo_root}/.claude/agents/${agent}.md"
    if [ -f "$agent_file" ]; then
        echo "PASS agent/${agent}"
    else
        echo "FAIL agent/${agent}"
        ((fail_count++))
    fi
done

# Exit with non-zero if any check failed
if [ $fail_count -gt 0 ]; then
    exit 1
else
    exit 0
fi

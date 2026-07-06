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

# Check the 4 CLIs
chk a1 a1
chk gh gh
chk git git
chk aliyun aliyun
chk cloudspec cloudspec

# Check credentials (each independent PASS/FAIL)
chk_cred aliyun "aliyun sts GetCallerIdentity"
chk_cred a1 "a1 auth whoami"

repo_root="$(git rev-parse --show-toplevel)"
# GitHub token daily check — WARN, NOT FAIL. A stale/missing JARVIS_GITHUB_TOKEN only
# blocks the GitHub escalation path (PR/评论/推分支); it must not fail preflight and
# thereby block Aone-only work. So we print a WARN line, do NOT increment fail_count,
# and drop an idempotent escalation note so the stale token stays visible/actionable.
# (github-identity.sh check itself is unchanged — it stays a hard gate at use time.)
if "$repo_root/bootstrap/github-identity.sh" check >/dev/null 2>&1; then
    echo "PASS jarvis-github-token"
else
    echo "WARN jarvis-github-token — JARVIS_GITHUB_TOKEN 失效/缺失；GitHub 升级路径(PR/评论/推分支)会被阻断，Aone-only 工作不受影响"
    echo "     修复: 刷新 api-tool-agent 的 GitHub token，更新 bootstrap/.env 的 JARVIS_GITHUB_TOKEN(或其环境来源)，再跑 bootstrap/github-identity.sh check 复验"
    esc_dir="${JARVIS_ESCALATION_DIR:-$repo_root/escalation}"
    esc_file="$esc_dir/github-token-invalid-$(date -u +%F).md"
    if [ ! -f "$esc_file" ]; then
        mkdir -p "$esc_dir"
        {
            echo "# GitHub token 失效/缺失 — $(date -u +%F)"
            echo ""
            echo "## 现象"
            echo "\`bootstrap/verify.sh\` 的 GitHub token 日检失败：\`bootstrap/github-identity.sh check\` 非零退出（JARVIS_GITHUB_TOKEN 过期 401 或缺失）。"
            echo ""
            echo "## 影响面"
            echo "仅阻断 GitHub 升级路径（terraform-provider-alicloud PR/评论/推分支，head=api-tool-agent:<branch>）。Aone-only 工作不受影响，故 verify 记 WARN 不硬失败。"
            echo ""
            echo "## 修复步骤"
            echo "1. 用 api-tool-agent 账号刷新 GitHub token（需 repo/workflow 权限）。"
            echo "2. 更新 \`bootstrap/.env\` 的 \`JARVIS_GITHUB_TOKEN\`（或其环境来源），确保 \`GH_TOKEN=\$JARVIS_GITHUB_TOKEN gh api user --jq .login\` 返回 \`api-tool-agent\`。"
            echo "3. 复验：\`bootstrap/github-identity.sh check\` 应打印 \`api-tool-agent\` 并退 0。"
        } > "$esc_file"
        echo "     已落 escalation 提示: $esc_file"
    fi
fi

# a1 (jarvis identity) login-state daily check — WARN, NOT FAIL. Mirrors the
# GitHub-token check above: a dead/expired a1 session only blocks a1-backed Aone
# writes (aone-triage 回复/建单、wrap.sh、claim/scan/reconcile — all go through
# bin/a1id); it must not fail preflight and thereby block non-a1 work. So we WARN,
# do NOT increment fail_count, and drop an idempotent escalation note.
# OK iff `bin/a1id -- auth whoami` yields Account == WORKER_1782379562571 (jarvis).
# 空/报错(过期会话) 或 半死(EmpID 在但 Account 空) → WARN. Parser matches bin/a1id.
a1id_bin="${JARVIS_A1ID:-$repo_root/bin/a1id}"
a1_expect="WORKER_1782379562571"
a1_account="$("$a1id_bin" -- auth whoami 2>/dev/null | awk '/Account:/{print $2}' || true)"
if [ "$a1_account" = "$a1_expect" ]; then
    echo "PASS jarvis-a1-session"
else
    echo "WARN jarvis-a1-session — a1 jarvis 登录态失效/缺失(whoami Account='${a1_account:-<空>}' != $a1_expect)；a1 相关 Aone 写(triage 回复/建单/wrap)会被阻断，非 a1 工作不受影响"
    echo "     修复: 浏览器 BUC 登 open_jarvis 后 bin/a1id login jarvis，再跑 bin/a1id -- auth whoami 复验 Account=$a1_expect"
    esc_dir="${JARVIS_ESCALATION_DIR:-$repo_root/escalation}"
    esc_file="$esc_dir/a1-session-expired-$(date -u +%F).md"
    if [ ! -f "$esc_file" ]; then
        mkdir -p "$esc_dir"
        {
            echo "# a1 jarvis 登录态失效/缺失 — $(date -u +%F)"
            echo ""
            echo "## 现象"
            echo "\`bootstrap/verify.sh\` 的 a1 登录态日检失败：\`bin/a1id -- auth whoami\` 的 Account 字段为空或非 jarvis 账号（期望 \`$a1_expect\`）。过期会话（报错/非零退出）与半死会话（EmpID 在、Account 空）均命中。"
            echo ""
            echo "## 影响面"
            echo "仅阻断 a1 相关 Aone 写路径（aone-triage 回复/建单、wrap.sh sync/done、claim/scan/reconcile 走 bin/a1id）。非 a1 工作不受影响，故 verify 记 WARN 不硬失败。"
            echo ""
            echo "## 修复步骤"
            echo "1. 浏览器登 BUC（https://buc.alibaba-inc.com/）为 open_jarvis 账号。"
            echo "2. 跑 \`bin/a1id login jarvis\`（走 BUC SSO，落盘 jarvis 身份）。"
            echo "3. 复验：\`bin/a1id -- auth whoami\` 应打印 \`Account: $a1_expect\` 并退 0。"
        } > "$esc_file"
        echo "     已落 escalation 提示: $esc_file"
    fi
fi

# Check vendored skills
chk_skill aone-triage

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

# Check the 3 agent definition files exist (主会话即总领，不单设 jarvis 子代理)
for agent in developer reviewer verifier; do
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

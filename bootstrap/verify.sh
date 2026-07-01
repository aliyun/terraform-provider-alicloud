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
if "$repo_root/bootstrap/github-identity.sh" check >/dev/null 2>&1; then
    echo "PASS jarvis-github-token"
else
    echo "FAIL jarvis-github-token"
    ((fail_count++))
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

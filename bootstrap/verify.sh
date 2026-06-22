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
chk_cred gh "gh auth status"
chk_cred aliyun "aliyun sts GetCallerIdentity"
chk_cred a1 "a1 auth whoami"

# Check vendored skills
chk_skill aone-triage

# Exit with non-zero if any check failed
if [ $fail_count -gt 0 ]; then
    exit 1
else
    exit 0
fi

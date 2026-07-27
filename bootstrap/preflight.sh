#!/usr/bin/env bash
# bootstrap/preflight.sh — 开局自检的日级闸门：install+verify 24h 跑一次即可。
# 上次全绿在 24h 内 → 直接放行;否则真跑 install→verify,全绿盖戳。
# 退码透传 verify.sh:0=全绿。强制重跑:JARVIS_PREFLIGHT_TTL=0,或 preflight.sh --force。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

# 保证 git 用仓库内的 hook(bootstrap/git-hooks/*);防止 skill 双份等漂移。
# git config 极廉价,每次跑都 idempotent 检查一次(TTL 命中跳过 install/verify 时也保这条)。
if git -C "$repo_root" rev-parse --git-dir >/dev/null 2>&1 && [ -d "$repo_root/bootstrap/git-hooks" ]; then
    if [ "$(git -C "$repo_root" config --local core.hooksPath 2>/dev/null)" != "bootstrap/git-hooks" ]; then
        git -C "$repo_root" config --local core.hooksPath bootstrap/git-hooks
        echo "preflight: set core.hooksPath=bootstrap/git-hooks"
    fi
fi

# AGENTS.md 从 CLAUDE.md 由 mirror sed 生成(AGENTS.md 已入库跟踪;preflight 兜底 +
# PostToolUse hook 实时 sync + pre-commit staged CLAUDE.md 自动重生成+add)。首次 checkout /
# codex 侧初始 / TTL 命中 skip 时也保它在 disk 上与 CLAUDE.md 一致,避免 codex 侧读到过期入口文档。
if [ -f "$repo_root/CLAUDE.md" ] && [ -f "$script_dir/skills-mirror-lib.sh" ]; then
    # shellcheck source=bootstrap/skills-mirror-lib.sh
    source "$script_dir/skills-mirror-lib.sh"
    mirror_sed_claude_to_codex < "$repo_root/CLAUDE.md" > "$repo_root/AGENTS.md"
fi

record_codex_stop_hook_root() {
    local codex_home state_dir roots_dir runner tmp_runner
    codex_home="${CODEX_HOME:-$HOME/.codex}"
    state_dir="$codex_home/jarvis"
    roots_dir="$state_dir/session-roots"
    runner="$state_dir/run-stop-hook.sh"

    mkdir -p "$roots_dir" || return 1
    printf '%s\n' "$repo_root" > "$state_dir/repo-root" || return 1
    if [ -n "${CODEX_THREAD_ID:-}" ]; then
        printf '%s\n' "$repo_root" > "$roots_dir/$CODEX_THREAD_ID" || return 1
    fi

    tmp_runner="$(mktemp "$state_dir/.run-stop-hook.XXXXXX")" || return 1
    if ! cat > "$tmp_runner" <<'EOF'
#!/usr/bin/env bash
set -uo pipefail
JARVIS_CODEX_HOOK_RUNNER_VERSION=2

codex_home="${CODEX_HOME:-$HOME/.codex}"
state_dir="$codex_home/jarvis"
repo=""

use_repo() {
    local candidate="$1"
    if [ -n "$candidate" ] && [ -f "$candidate/bootstrap/run-stop-hook.sh" ]; then
        repo="$candidate"
        return 0
    fi
    return 1
}

if [ -n "${JARVIS_ROOT:-}" ]; then
    use_repo "$JARVIS_ROOT" || true
fi

if [ -z "$repo" ]; then
    git_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
    use_repo "$git_root" || true
fi

if [ -z "$repo" ] && [ -n "${CODEX_THREAD_ID:-}" ]; then
    thread_root="$(cat "$state_dir/session-roots/$CODEX_THREAD_ID" 2>/dev/null || true)"
    use_repo "$thread_root" || true
fi

if [ -z "$repo" ]; then
    default_root="$(cat "$state_dir/repo-root" 2>/dev/null || true)"
    use_repo "$default_root" || true
fi

if [ -z "$repo" ]; then
    echo "stop-hook: cannot resolve repo root from cwd=$PWD" >&2
    exit 2
fi

if [ "${1:-}" = "prompt" ]; then
    exec bash "$repo/bootstrap/run-interactive-worker-hook.sh" codex UserPromptSubmit
fi
exec bash "$repo/bootstrap/run-stop-hook.sh" "$@"
EOF
    then
        rm -f "$tmp_runner"
        return 1
    fi
    chmod +x "$tmp_runner" || { rm -f "$tmp_runner"; return 1; }
    mv -f "$tmp_runner" "$runner" || { rm -f "$tmp_runner"; return 1; }
}

if ! record_codex_stop_hook_root; then
    echo "preflight: failed to install Codex hook runner" >&2
    exit 1
fi

ttl="${JARVIS_PREFLIGHT_TTL:-86400}"   # 24h
[ "${1:-}" = "--force" ] && ttl=0

# 84550781: 每次 preflight 都清一次 a1 孤儿锁（best-effort，无 a1 存活才动手）。
# 必须在 TTL 闸门之前——否则 24h 内 skip 时永远不清；脚本自带前置检查不会误伤活 a1。
bash "$script_dir/a1-locks-clean.sh" || true

if bash "$script_dir/cache.sh" fresh "preflight.ok" "$ttl"; then
    echo "preflight: skip (verified < $((ttl/3600))h ago)"
    exit 0
fi

bash "$script_dir/install.sh"
if bash "$script_dir/verify.sh"; then
    bash "$script_dir/cache.sh" get "preflight.ok" 0 -- date -u +%Y-%m-%dT%H:%M:%SZ >/dev/null
    echo "preflight: PASS (stamped)"
    exit 0
fi
echo "preflight: FAIL — 自检未全绿,见上方,先修再干活" >&2
exit 1

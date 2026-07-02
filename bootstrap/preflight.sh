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

ttl="${JARVIS_PREFLIGHT_TTL:-86400}"   # 24h
[ "${1:-}" = "--force" ] && ttl=0

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

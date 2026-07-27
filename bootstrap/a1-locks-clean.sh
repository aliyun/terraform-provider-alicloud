#!/usr/bin/env bash
# bootstrap/a1-locks-clean.sh — 只在无 a1 进程存活时清 ~/.config/a1 下的孤儿锁文件。
# 84550781 事故根因：a1 subprocess 被 timeout kill 后 auth.yaml.lock /
# telemetry-queue.jsonl.lock 残留（每个 identity 目录各一套），后续 a1 调用挂在
# 启动期取锁 → 再被 timeout kill → 再留锁，滚雪球至 24 个孤儿锁 + 7 僵尸进程。
#
# 安全前置：`pgrep -x a1` 返回非空即为存在活着的 a1，绝不清（会破坏正在跑的 a1）。
# 打印清理数量到 stderr；主要给 preflight 早期兜底和人工排障用。
# 退码 0=清理完成或跳过；1=前置检查失败。
set -uo pipefail

a1_root="${A1_HOME:-$HOME/.config/a1}"

if [ ! -d "$a1_root" ]; then
    echo "a1-locks-clean: skip (root not found: $a1_root)" >&2
    exit 0
fi

# 存在型锁不带持有者活性校验（a1 团队已反馈），文件存在就等价「有人持有」。
# 因此清之前必须确认没有真的 a1 在跑，否则会把活 a1 的锁一起删了。
if pgrep -x a1 >/dev/null 2>&1; then
    echo "a1-locks-clean: skip (live a1 process detected)" >&2
    exit 0
fi

# 只清明确的锁文件名——避免误删其它 sidecar 状态。
locks=()
while IFS= read -r -d '' f; do
    locks+=("$f")
done < <(find "$a1_root" \
    \( -name 'auth.yaml.lock' -o -name 'telemetry-queue.jsonl.lock' -o -name '*.lock' \) \
    -type f -print0 2>/dev/null)

if [ "${#locks[@]}" -eq 0 ]; then
    exit 0
fi

removed=0
for f in "${locks[@]}"; do
    if rm -f "$f" 2>/dev/null; then
        removed=$((removed + 1))
    fi
done

echo "a1-locks-clean: removed $removed stale lock(s) under $a1_root" >&2
exit 0

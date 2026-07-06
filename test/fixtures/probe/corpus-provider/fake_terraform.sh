#!/usr/bin/env bash
# fake_terraform.sh — hermetic terraform stub for probe-corpus.sh validate tests.
# 不联网、不装 provider。按 *.tf 内的 magic marker 决定 init/validate/fmt 成败:
#   CORPUS_FAIL_INIT     → init 退 1
#   CORPUS_FAIL_VALIDATE → validate 退 1
#   CORPUS_FAIL_FMT      → fmt -check 退 1
# 扫描范围:$PWD 下的 *.tf,以及任何以目录形式出现的位置参数。bash 3.2 兼容(无 mapfile)。
set -uo pipefail

sub="${1:-}"; shift 2>/dev/null || true

# _has_marker <marker> — 在待扫描 *.tf(当前目录 + 目录型参数)中命中 marker → 0
_has_marker() {
    local m="$1"; shift 2>/dev/null || true
    local f a
    shopt -s nullglob
    for f in ./*.tf; do [ -f "$f" ] && grep -q "$m" "$f" 2>/dev/null && return 0; done
    for a in "$@"; do
        [ -d "$a" ] || continue
        for f in "$a"/*.tf; do [ -f "$f" ] && grep -q "$m" "$f" 2>/dev/null && return 0; done
    done
    return 1
}

case "$sub" in
    version) echo "Terraform v9.9.9-corpus-stub"; exit 0 ;;
    init)
        _has_marker CORPUS_FAIL_INIT "$@" && { echo "stub: init failed (marker)"; exit 1; }
        echo "stub: init ok"; exit 0 ;;
    validate)
        _has_marker CORPUS_FAIL_VALIDATE "$@" && { echo "stub: validate error (marker)"; exit 1; }
        echo "stub: validate success"; exit 0 ;;
    fmt)
        is_check=0
        for a in "$@"; do [ "$a" = "-check" ] && is_check=1; done
        if [ "$is_check" -eq 1 ]; then
            _has_marker CORPUS_FAIL_FMT "$@" && { echo "stub: fmt drift (marker)"; exit 1; }
        fi
        exit 0 ;;
    *) echo "stub: ignoring '$sub'"; exit 0 ;;
esac

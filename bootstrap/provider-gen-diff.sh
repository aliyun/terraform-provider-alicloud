#!/usr/bin/env bash
# bootstrap/provider-gen-diff.sh
#
# Diff a Terraform alicloud provider resource's hand-modifications against the
# pure generator-v4 baseline.
#
# WHAT IT DOES
#   Regenerates the resource + data source + docs via generator-v4 (CloudSpec
#   pre env) into throwaway origin/master worktrees, then diffs each generated
#   file against the actual (hand-modified) file in your provider worktree.
#   The diff = what was hand-fixed on top of generator output. A checklist
#   summary then labels the 7 known generator-defect handfix categories
#   (see .claude/skills/terraform-provider-release/references/generator-post-gen-checklist.md).
#
#   If the generator cannot fetch AMP pre metadata, generation is skipped and
#   only the checklist annotation runs (degraded mode).
#
# USAGE
#   provider-gen-diff.sh <resource-name> [worktree-dir] [--refresh] [--no-color] [--checklist-only]
#   provider-gen-diff.sh vpc_route_target_group
#   provider-gen-diff.sh vpc_route_target_group /path/to/provider-worktree
#   provider-gen-diff.sh --popcode Vpc --resource RouteTargetGroup [worktree-dir]
#
#   <resource-name>  snake_case, first segment = product. e.g. vpc_route_target_group
#   [worktree-dir]   provider checkout holding the ACTUAL files. Defaults to
#                    `bootstrap/workspace.sh dir terraform_provider`.
#
# FLAGS
#   --refresh          regenerate the baseline even if cached
#   --no-color         plain diff (no ANSI color)
#   --checklist-only   skip generation/diff, run just the checklist annotation
#   --popcode X        explicit PopCode (namespace) instead of deriving
#   --resource Y       explicit resource type code (CamelCase) instead of deriving
#   --html[=PATH]      render diff+checklist to an HTML page and open it.
#                     PATH defaults to $cache_dir/diff.html. Implies --no-color.
#
# ENV
#   JARVIS_GEN_DIFF_CACHE       cache dir (default $HOME/.cache/jarvis-gen-diff)
#   JARVIS_GEN_DIFF_GEN_TIMEOUT per-type generation timeout in seconds (default 300)
#
# CACHING
#   The regenerated baseline is cached per (popcode,resource). Re-runs diff
#   against the cached baseline instantly unless --refresh. Pre metadata
#   evolves over time; use --refresh when pre has changed since the last run.
#
# NOTE ON PRE-DRIFT
#   The baseline reflects the CURRENT pre env. A PR generated from an older pre
#   may show extra diff hunks that are pre-drift, not handfixes. Cross-reference
#   with the checklist summary to tell handfixes from drift.

set -euo pipefail

prog="$(basename "$0")"

# Resolve jarvis repo root (this script lives in bootstrap/).
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="$(cd "$script_dir/.." && pwd)"

# ── helpers ─────────────────────────────────────────────────────────────────

die() {
    echo "$prog: $*" >&2
    if [ -n "${html_capture:-}" ]; then
        exec 1>&3 3>&- 2>/dev/null || true
        rm -f "$html_capture" 2>/dev/null || true
    fi
    exit 1
}

log() {
    if [ "$want_color" -eq 1 ]; then printf '\033[1;36m%s\033[0m\n' "$*"
    else printf '%s\n' "$*"; fi
}
warn() {
    if [ "$want_color" -eq 1 ]; then printf '\033[1;33m%s\033[0m\n' "$*" >&2
    else printf '%s\n' "$*" >&2; fi
}

# Capitalize first letter (vpc -> Vpc).
capitalize() {
    local s="$1"
    printf '%s%s' "$(printf '%s' "${s:0:1}" | tr '[:lower:]' '[:upper:]')" "${s:1}"
}

# snake_case -> CamelCase (route_target_group -> RouteTargetGroup).
to_camel() {
    awk -F_ '{
        out="";
        for (i=1;i<=NF;i++) {
            w=$i;
            out=out toupper(substr(w,1,1)) substr(w,2);
        }
        print out;
    }' <<<"$1"
}

# CamelCase -> snake_case (RealtimeCompute -> realtime_compute).
# macOS BSD sed has no \L, so insert '_' before each capital then tr lower.
to_snake() {
    printf '%s' "$1" | sed -E 's/([A-Z])/_\1/g' | tr '[:upper:]' '[:lower:]' | sed 's/^_//'
}

# product segment -> PopCode. Default = capitalize; override non-obvious ones.
popcode_for() {
    case "$1" in
        # Add non-obvious PopCode overrides here as discovered.
        *) capitalize "$1" ;;
    esac
}

# Resolve a single file by glob under a base dir. $1=base $2=glob-pattern.
# Prints the match, or empty + warning to stderr if 0/multiple. Excludes _test.go.
resolve_one() {
    local base="$1" pat="$2"
    local matches
    matches="$(find "$base" -path "$pat" -type f 2>/dev/null | grep -v '_test\.go$' || true)"
    local n
    n="$(printf '%s\n' "$matches" | grep -c . || true)"
    if [ "$n" -eq 1 ]; then
        printf '%s' "$matches"
    elif [ "$n" -eq 0 ]; then
        printf '' ; warn "  (no file matching $pat under $base)"
    else
        printf '%s\n' "$matches" | head -1
        warn "  (multiple files matching $pat under $base; using first)"
    fi
}

# ── arg parsing ─────────────────────────────────────────────────────────────

popcode=""; resource=""; worktree=""; resource_name=""
refresh=0; want_color=1; checklist_only=0; html_out=""

while [ $# -gt 0 ]; do
    case "$1" in
        --refresh) refresh=1; shift ;;
        --no-color) want_color=0; shift ;;
        --checklist-only) checklist_only=1; shift ;;
        --popcode) popcode="${2:-}"; shift 2 ;;
        --resource) resource="${2:-}"; shift 2 ;;
        --html) html_out=1; shift ;;
        --html=*) html_out="${1#--html=}"; shift ;;
        --help|-h) sed -n '2,/^$/p' "$0" | sed 's/^# \?//'; exit 0 ;;
        --*) die "unknown flag: $1" ;;
        *)
            if [ -z "$resource_name" ]; then resource_name="$1"
            elif [ -z "$worktree" ]; then worktree="$1"
            else die "unexpected positional arg: $1"; fi
            shift ;;
    esac
done

[ -n "$resource_name" ] || [ -n "$resource" ] || \
    die "usage: $prog <resource-name> [worktree-dir]  (or --popcode X --resource Y)"
[ -n "$resource" ] && [ -z "$popcode" ] && die "--resource requires --popcode"

# Resolve worktree (actual files).
if [ -z "$worktree" ]; then
    worktree="$("$jarvis_root/bootstrap/workspace.sh" dir terraform_provider)" \
        || die "cannot resolve terraform_provider workspace (bootstrap/workspace.sh)"
fi
[ -d "$worktree/alicloud" ] || die "worktree has no alicloud/ dir: $worktree"

# Resolve popcode/resource + snake_case segments.
if [ -z "$resource" ]; then
    rn="${resource_name#resource_alicloud_}"
    rn="${rn#data_source_alicloud_}"
    rn="${rn%.go}"
    product_lc="${rn%%_*}"                 # vpc
    res_lc="${rn#"${product_lc}"_}"         # route_target_group (strip leading "vpc_")
    resource="$(to_camel "$res_lc")"        # RouteTargetGroup (-r flag value)
    popcode="$(popcode_for "$product_lc")"  # Vpc (-n flag value)
else
    product_lc="$(to_snake "$popcode")"
    res_lc="$(to_snake "$resource")"
    # Strip the product prefix when the resource code carries it
    # (RealtimeComputeVariable -> variable; RouteTargetGroup stays route_target_group).
    res_lc="${res_lc#${product_lc}_}"
fi

# ── HTML render setup ─────────────────────────────────────────────────────────
# When --html is set: force --no-color (renderer needs plain text), resolve the
# default output path, and redirect stdout to a temp capture file so the entire
# diff+checklist output feeds the renderer. stderr (warn/die) still hits the
# terminal so errors stay visible.
render_py="$script_dir/provider-gen-diff-render.py"
if [ -n "$html_out" ]; then
    want_color=0
    if [ "$html_out" = "1" ]; then
        _cr="${JARVIS_GEN_DIFF_CACHE:-$HOME/.cache/jarvis-gen-diff}"
        html_out="$_cr/${popcode}_${resource}/diff.html"
    fi
    if ! command -v python3 >/dev/null 2>&1; then
        warn "python3 not found; --html disabled, will print text to terminal"
        html_out=""
    elif [ ! -f "$render_py" ]; then
        warn "renderer not found ($render_py); --html disabled"
        html_out=""
    fi
fi
if [ -n "$html_out" ]; then
    html_capture="$(mktemp -t gen-diff-capture)" || die "mktemp failed"
    exec 3>&1 1>"$html_capture"
fi

log "== provider-gen-diff =="
echo "  resource : $resource  (popcode=$popcode)"
echo "  worktree  : $worktree"
echo "  snake     : $product_lc / $res_lc"

# Resolve generator workspace.
gen="$("$jarvis_root/bootstrap/workspace.sh" dir terraform_generator_v4 2>/dev/null || true)"
if [ -z "$gen" ] || [ ! -f "$gen/cli.php" ]; then
    warn "generator (terraform_generator_v4) not found; falling back to checklist-only"
    checklist_only=1
else
    echo "  generator : $gen"
fi

# Actual (hand-modified) files in the worktree.
actual_r="$worktree/alicloud/resource_alicloud_${product_lc}_${res_lc}.go"
actual_d="$(resolve_one "$worktree" "*/alicloud/data_source_alicloud_${product_lc}_${res_lc}*.go")"
actual_doc_r="$worktree/website/docs/r/${product_lc}_${res_lc}.html.markdown"
actual_doc_d="$(resolve_one "$worktree" "*/website/docs/d/${product_lc}_${res_lc}*.html.markdown")"

# ── baseline cache ──────────────────────────────────────────────────────────

cache_root="${JARVIS_GEN_DIFF_CACHE:-$HOME/.cache/jarvis-gen-diff}"
cache_dir="$cache_root/${popcode}_${resource}"
baseline_r="$cache_dir/resource.go"
baseline_d="$cache_dir/datasource.go"
baseline_doc_r="$cache_dir/r.html.markdown"
baseline_doc_d="$cache_dir/d.html.markdown"

cache_hit=0
if [ "$checklist_only" -eq 0 ] && [ "$refresh" -eq 0 ] && [ -f "$baseline_r" ]; then
    cache_hit=1
    log "cached baseline found ($cache_dir) — diffing instantly (use --refresh to regenerate)"
fi

# ── generation ──────────────────────────────────────────────────────────────

generate_type() {
    # $1 = type (resource|datasource|document), $2 = baseline worktree path
    local t="$1" bt="$2"
    local timeout_s="${JARVIS_GEN_DIFF_GEN_TIMEOUT:-300}"
    log "  generating $t (timeout ${timeout_s}s)..."
    # best-effort: generator may exit non-zero on a provider.go map step but
    # still emit the target file; we check the file after, not the exit code.
    php "$gen/cli.php" -n "$popcode" -r "$resource" -e pre -i pre -t "$t" -p "$bt" \
        >"$cache_dir/gen-$t.log" 2>&1 < /dev/null &
    local pid=$!
    local waited=0
    while kill -0 "$pid" 2>/dev/null; do
        sleep 2; waited=$((waited+2))
        if [ "$waited" -ge "$timeout_s" ]; then
            warn "  $t generation exceeded ${timeout_s}s; killing"
            kill -TERM "$pid" 2>/dev/null || true; sleep 1
            kill -KILL "$pid" 2>/dev/null || true
            return 1
        fi
    done
    wait "$pid" 2>/dev/null || true
    return 0
}

regenerate_baseline() {
    [ -n "$gen" ] && [ -f "$gen/cli.php" ] || { warn "no generator; cannot regenerate"; return 1; }
    mkdir -p "$cache_dir"
    rm -f "$baseline_r" "$baseline_d" "$baseline_doc_r" "$baseline_doc_d"

    log "regenerating baseline from origin/master via generator-v4 (pre)..."
    # Fresh origin/master so the generator patches a clean tree.
    git -C "$worktree" fetch origin master --quiet 2>/dev/null || warn "  (fetch origin master warning, continuing)"

    # Three throwaway --detach worktrees so the three generation types (which
    # each patch provider.go) run in parallel without clobbering each other.
    local bt_r bt_d bt_doc
    bt_r="$(mktemp -d -t gen-diff-r)"   || die "mktemp failed"
    bt_d="$(mktemp -d -t gen-diff-d)"   || die "mktemp failed"
    bt_doc="$(mktemp -d -t gen-diff-doc)" || die "mktemp failed"
    # shellcheck disable=SC2064
    trap "git -C '$worktree' worktree remove --force '$bt_r' 2>/dev/null; git -C '$worktree' worktree remove --force '$bt_d' 2>/dev/null; git -C '$worktree' worktree remove --force '$bt_doc' 2>/dev/null; rm -rf '$bt_r' '$bt_d' '$bt_doc'" EXIT

    git -C "$worktree" worktree add --detach "$bt_r"   origin/master >/dev/null 2>&1 || die "worktree add (resource) failed"
    git -C "$worktree" worktree add --detach "$bt_d"   origin/master >/dev/null 2>&1 || die "worktree add (datasource) failed"
    git -C "$worktree" worktree add --detach "$bt_doc" origin/master >/dev/null 2>&1 || die "worktree add (document) failed"

    generate_type resource   "$bt_r"   &
    local pid_r=$!
    generate_type datasource "$bt_d"   &
    local pid_d=$!
    generate_type document   "$bt_doc" &
    local pid_doc=$!

    local fail=0
    wait "$pid_r"   || { warn "  resource generation failed/timeout"; fail=1; }
    wait "$pid_d"   || { warn "  datasource generation failed/timeout"; fail=1; }
    wait "$pid_doc" || { warn "  document generation failed/timeout"; fail=1; }

    # Copy target files out of each baseline worktree into the cache.
    # NOTE on doc placement: -t document produces the RESOURCE doc (r/) only;
    # the DATA SOURCE doc (d/) is produced by -t datasource. So gen_doc_d is
    # resolved from the datasource worktree (bt_d), not the document worktree.
    local gen_r="$bt_r/alicloud/resource_alicloud_${product_lc}_${res_lc}.go"
    local gen_d="$(resolve_one "$bt_d" "*/alicloud/data_source_alicloud_${product_lc}_${res_lc}*.go")"
    local gen_doc_r="$bt_doc/website/docs/r/${product_lc}_${res_lc}.html.markdown"
    local gen_doc_d="$(resolve_one "$bt_d" "*/website/docs/d/${product_lc}_${res_lc}*.html.markdown")"

    [ -f "$gen_r" ]     && cp "$gen_r"     "$baseline_r"     || warn "  baseline resource file not produced"
    [ -n "$gen_d" ] && [ -f "$gen_d" ]     && cp "$gen_d"     "$baseline_d"     || warn "  baseline datasource file not produced"
    [ -f "$gen_doc_r" ] && cp "$gen_doc_r" "$baseline_doc_r" || warn "  baseline r-doc not produced"
    [ -n "$gen_doc_d" ] && [ -f "$gen_doc_d" ] && cp "$gen_doc_d" "$baseline_doc_d" || warn "  baseline d-doc not produced"

    if [ ! -f "$baseline_r" ]; then
        warn "baseline regeneration did not produce the resource file; see $cache_dir/gen-*.log"
        warn "falling back to checklist-only"
        checklist_only=1
        return 1
    fi
    log "baseline regenerated -> $cache_dir"
}

if [ "$checklist_only" -eq 0 ]; then
    if [ "$cache_hit" -eq 0 ]; then
        regenerate_baseline || true
    fi
fi

# ── diff ─────────────────────────────────────────────────────────────────────

diff_files() {
    # $1 label, $2 baseline, $3 actual
    local label="$1" b="$2" a="$3"
    echo
    log "── $label ──"
    if [ ! -f "$b" ]; then echo "  (no baseline)"; return; fi
    if [ ! -f "$a" ]; then echo "  (no actual file in worktree: $a)"; return; fi
    local diff_args=(--no-index)
    [ "$want_color" -eq 1 ] && diff_args+=(--color=always) || diff_args+=(--color=never)
    # git diff --no-index exits 1 on differences; that is expected.
    git diff "${diff_args[@]}" -- "$b" "$a" || true
}

if [ "$checklist_only" -eq 0 ] && { [ "$cache_hit" -eq 1 ] || [ -f "$baseline_r" ]; }; then
    diff_files "resource: resource_alicloud_${product_lc}_${res_lc}.go" "$baseline_r" "$actual_r"
    diff_files "datasource: data_source_alicloud_${product_lc}_${res_lc}*.go" "$baseline_d" "$actual_d"
    diff_files "doc (r): ${product_lc}_${res_lc}.html.markdown" "$baseline_doc_r" "$actual_doc_r"
    diff_files "doc (d): ${product_lc}_${res_lc}*.html.markdown" "$baseline_doc_d" "$actual_doc_d"
fi

# ── checklist annotation ────────────────────────────────────────────────────
# Detects the 7 known generator-v4 defect handfix categories. When a baseline
# exists, compares actual vs baseline so the verdict reflects the handfix
# delta; otherwise reports presence/absence in the actual file only.

echo
log "== checklist (generator-v4 post-gen defect handfixes) =="

chk_grep() { # $1 file $2 pattern -> prints matching lines with numbers, or none
    [ -f "$1" ] || { echo "(file missing)"; return; }
    grep -nE "$2" "$1" 2>/dev/null || true
}
chk_count() { # $1 file $2 pattern -> count
    [ -f "$1" ] || { printf 0; return; }
    local c
    c="$(grep -cE "$2" "$1" 2>/dev/null || true)"
    printf '%s' "${c:-0}"
}

echo
echo "[Bug 1] Create-only attrs missing ForceNew (handfix = ForceNew added)"
echo "  actual  ForceNew: true count : $(chk_count "$actual_r" 'ForceNew:\s*true')"
if [ -f "$baseline_r" ]; then
    echo "  baseline ForceNew: true count: $(chk_count "$baseline_r" 'ForceNew:\s*true')"
fi
chk_grep "$actual_r" 'ForceNew:\s*true' | sed 's/^/    /' | head -8 || true

echo
echo "[Bug 2] Read not reverse-mapping (handfix = d.Set from nested/objectRaw)"
echo "  Read d.Set referencing objectRaw/nested builders:"
chk_grep "$actual_r" 'd\.Set\(' | grep -iE 'objectRaw|make\(\[\]|for .*range' | sed 's/^/    /' | head -8 || true

echo
echo "[Bug 3] Datasource nested schema missing fields (handfix = Elem/schema keys added)"
echo "  datasource Read mapping[] assignments (each key must exist in the nested Elem schema):"
chk_grep "$actual_d" 'mapping\["' | sed 's/^/    /' | head -10 || true

echo
echo "[Bug 4] Docs redundant NOTE (handfix = NOTE removed)"
for f in "$actual_doc_r" "$actual_doc_d"; do
    [ -f "$f" ] || continue
    local_n=$(basename "$f")
    actual_n=$(chk_count "$f" 'NOTE:.*only evaluated|NOTE:.*immutable|Modifying it in isolation|Changing it after creation has no effect')
    base_f=""
    case "$f" in
        */docs/r/*) base_f="$baseline_doc_r" ;;
        */docs/d/*) base_f="$baseline_doc_d" ;;
    esac
    printf '  %-32s NOTE lines: actual=%s' "$local_n" "$actual_n"
    if [ -f "$base_f" ]; then printf ' baseline=%s' "$(chk_count "$base_f" 'NOTE:.*only evaluated|NOTE:.*immutable|Modifying it in isolation|Changing it after creation has no effect')"; fi
    printf '\n'
done

echo
echo "[Bug 5] Docs example placeholder (handfix = real example written)"
for f in "$actual_doc_r" "$actual_doc_d"; do
    [ -f "$f" ] || continue
    local_n=$(basename "$f")
    actual_n=$(chk_count "$f" '没有资源测试用例|No resource test case|please pass the resource test case')
    printf '  %-32s placeholder lines: %s  %s\n' "$local_n" "$actual_n" "$([ "$actual_n" = 0 ] && echo '(handfixed: real example)' || echo '(placeholder still present)')"
done

echo
echo "[Bug 6] addDebug order after SetId (observation; no stable grep — review Create in the diff above)"

echo
echo "[Bug 7] Update trigger only on partial fields (handfix = HasChange branches)"
echo "  Update-func d.HasChange / update=true occurrences:"
chk_grep "$actual_r" 'd\.HasChange\(|update\s*=\s*true' | sed 's/^/    /' | head -8 || true

echo
log "done. baseline cache: $cache_dir"

# ── HTML render ───────────────────────────────────────────────────────────────
# Restore stdout (saved on fd3), render the captured text to HTML, open it.
if [ -n "${html_out:-}" ]; then
    exec 1>&3 3>&-
    mkdir -p "$(dirname "$html_out")" 2>/dev/null || true
    if python3 "$render_py" "$html_capture" "$html_out"; then
        if [ "$(uname)" = "Darwin" ]; then
            open "$html_out" 2>/dev/null || true
        fi
        log "HTML 渲染完成 → $html_out"
    else
        warn "HTML 渲染失败;原始文本见 $html_capture"
        cat "$html_capture" 2>/dev/null || true
    fi
    rm -f "$html_capture" 2>/dev/null || true
fi

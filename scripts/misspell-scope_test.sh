#!/bin/bash
# Regression test for the documentation spelling exemption used by
# scripts/local-ci-check.sh. It exercises the real helpers with a stand-in
# misspell binary: no network access and no provider build.
set -u

here="$(cd "$(dirname "$0")" && pwd)"
. "$here/lib/misspell-scope.sh"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

enum_doc="website/docs/d/ecs_image_pipeline_executions.html.markdown"
other_doc="website/docs/d/reserved_instances.html.markdown"

# The exemption is keyed on the real source position, so the stand-in binary
# reports findings on the exact line and column of the API enum literals, on an
# ordinary prose line of the same page, and at a wrong column of the enum line.
enum_positions="$(misspell_enum_positions)"
enum_line_first="$(printf '%s\n' "$enum_positions" | head -1)"
enum_line="${enum_line_first%%:*}"
enum_col_first="$(printf '%s\n' "$enum_positions" | head -1 | cut -d: -f2)"
enum_col_second="$(printf '%s\n' "$enum_positions" | sed -n 2p | cut -d: -f2)"
if [ -z "$enum_line" ] || [ -z "$enum_col_first" ] || [ -z "$enum_col_second" ]; then
  echo "FAIL: the API enum literal positions were not found in $enum_doc"
  exit 1
fi
if [ "$enum_line" -gt 1 ]; then
  prose_line=$((enum_line - 1))
else
  prose_line=$((enum_line + 1))
fi
# A column the enum literal does not occupy: ordinary prose on the same line.
wrong_col=$((enum_col_first + 1))
export FAKE_ENUM_LINE="$enum_line"
export FAKE_ENUM_COL_FIRST="$enum_col_first"
export FAKE_ENUM_COL_SECOND="$enum_col_second"
export FAKE_WRONG_COL="$wrong_col"
export FAKE_PROSE_LINE="$prose_line"

cat > "$tmpdir/misspell" <<'FAKE'
#!/bin/bash
for argument in "$@"; do last="$argument"; done
mode="${FAKE_MISSPELL_MODE:-clean}"
is_enum_page=0
case "$last" in
  *ecs_image_pipeline_executions*) is_enum_page=1 ;;
esac
case "$mode" in
  clean)
    exit 0
    ;;
  literals)
    # Both findings sit on the documented positions: the accepted exception.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_ENUM_LINE}:${FAKE_ENUM_COL_FIRST}: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:${FAKE_ENUM_LINE}:${FAKE_ENUM_COL_SECOND}: \"CANCELLED\" is a misspelling of \"CANCELED\""
      exit 2
    fi
    exit 0
    ;;
  same_line_wrong_column)
    # The same word on the enum line but at a column the API literal does not
    # occupy: ordinary prose, not exempt.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_ENUM_LINE}:${FAKE_WRONG_COL}: \"CANCELLING\" is a misspelling of \"CANCELING\""
      exit 2
    fi
    exit 0
    ;;
  prose)
    # Uppercase CANCELLING in ordinary prose on the same page: not exempt.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_PROSE_LINE}:1: \"CANCELLING\" is a misspelling of \"CANCELING\""
      exit 2
    fi
    exit 0
    ;;
  mixed)
    # The enum literals are exempt, the same-page prose is not.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_ENUM_LINE}:${FAKE_ENUM_COL_FIRST}: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:${FAKE_PROSE_LINE}:1: \"CANCELLED\" is a misspelling of \"CANCELED\""
      exit 2
    fi
    exit 0
    ;;
  otherpage)
    # The same literals on another page are not exempt.
    printf '%s\n' \
      "$other_doc:1:1: \"CANCELLING\" is a misspelling of \"CANCELING\""
    exit 2
    ;;
  unrelated)
    printf '%s\n' "website/docs/d/other.html.markdown:1:1: \"cancelling\" is a misspelling of \"canceling\""
    exit 2
    ;;
  empty)
    exit 2
    ;;
  usage)
    echo "misspell: flag provided but not defined"
    exit 1
    ;;
  crash)
    echo "misspell: internal error"
    exit 2
    ;;
esac
exit 0
FAKE
chmod +x "$tmpdir/misspell"
export PATH="$tmpdir:$PATH"

failures=0
expect() {
  local name="$1" expected="$2" mode="$3"
  shift 3
  local actual
  FAKE_MISSPELL_MODE="$mode" "$@" > "$tmpdir/out" 2>&1 && actual=pass || actual=fail
  if [ "$actual" != "$expected" ]; then
    echo "FAIL: $name (expected $expected, got $actual)"
    sed 's/^/    /' "$tmpdir/out"
    failures=$((failures + 1))
  else
    echo "ok: $name ($actual)"
  fi
}

# strict mode: one target directory or page
expect "strict clean" pass clean run_misspell_scoped ./website/docs
expect "strict api literals accepted" pass literals run_misspell_scoped "$enum_doc"
expect "strict same line wrong column rejected" fail same_line_wrong_column run_misspell_scoped "$enum_doc"
expect "strict prose in the same page rejected" fail prose run_misspell_scoped "$enum_doc"
expect "strict mixed findings rejected" fail mixed run_misspell_scoped "$enum_doc"
expect "strict literals on another page rejected" fail otherpage run_misspell_scoped "$other_doc"
expect "strict unrelated file rejected" fail unrelated run_misspell_scoped "$other_doc"
expect "strict empty tool failure rejected" fail empty run_misspell_scoped "$enum_doc"
expect "strict tool crash rejected" fail crash run_misspell_scoped "$enum_doc"
expect "strict usage error rejected" fail usage run_misspell_scoped "$enum_doc"

# changed-docs mode: newline separated document list
expect "changed docs multi-line list accepted" pass literals run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs same line wrong column rejected" fail same_line_wrong_column run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs prose rejected" fail prose run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs mixed findings rejected" fail mixed run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs empty tool failure rejected" fail empty run_misspell_changed_docs "$enum_doc"

if [ "$failures" -ne 0 ]; then
  echo "spelling scope test: $failures failure(s)"
  exit 1
fi
echo "spelling scope test: all cases passed"

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
# reports findings on the line of the API enum declaration and on an ordinary
# prose line of the same page.
enum_line="$(grep -n 'CANCELLING` and `CANCELLED`\.' "$enum_doc" | head -1 | cut -d: -f1)"
if [ -z "$enum_line" ]; then
  echo "FAIL: the API enum line was not found in $enum_doc"
  exit 1
fi
if [ "$enum_line" -gt 1 ]; then
  prose_line=$((enum_line - 1))
else
  prose_line=$((enum_line + 1))
fi
export FAKE_ENUM_LINE="$enum_line"
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
    # Both findings are on the API enum line: the documented exception.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_ENUM_LINE}:263: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:${FAKE_ENUM_LINE}:300: \"CANCELLED\" is a misspelling of \"CANCELED\""
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
    # The enum literal is exempt, the prose on the same page is not.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_ENUM_LINE}:263: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
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
expect "strict prose in the same page rejected" fail prose run_misspell_scoped "$enum_doc"
expect "strict mixed findings rejected" fail mixed run_misspell_scoped "$enum_doc"
expect "strict literals on another page rejected" fail otherpage run_misspell_scoped "$other_doc"
expect "strict unrelated file rejected" fail unrelated run_misspell_scoped "$other_doc"
expect "strict empty tool failure rejected" fail empty run_misspell_scoped "$enum_doc"
expect "strict tool crash rejected" fail crash run_misspell_scoped "$enum_doc"
expect "strict usage error rejected" fail usage run_misspell_scoped "$enum_doc"

# changed-docs mode: newline separated document list
expect "changed docs multi-line list accepted" pass literals run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs prose rejected" fail prose run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs mixed findings rejected" fail mixed run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs empty tool failure rejected" fail empty run_misspell_changed_docs "$enum_doc"

if [ "$failures" -ne 0 ]; then
  echo "spelling scope test: $failures failure(s)"
  exit 1
fi
echo "spelling scope test: all cases passed"

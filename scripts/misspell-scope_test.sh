#!/bin/bash
# Regression test for the documentation spelling exemption used by
# scripts/local-ci-check.sh. It exercises the real helpers with a stand-in
# misspell binary: no network access and no provider build.
set -u

here="$(cd "$(dirname "$0")" && pwd)"
. "$here/lib/misspell-scope.sh"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

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
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:76:263: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:76:280: \"CANCELLED\" is a misspelling of \"CANCELED\""
      exit 2
    fi
    exit 0
    ;;
  prose)
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:76:263: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:90:3: \"cancelling\" is a misspelling of \"canceling\""
      exit 2
    fi
    exit 0
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

enum_doc="website/docs/d/ecs_image_pipeline_executions.html.markdown"
other_doc="website/docs/d/reserved_instances.html.markdown"

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
expect "strict unrelated file rejected" fail unrelated run_misspell_scoped "$other_doc"
expect "strict empty tool failure rejected" fail empty run_misspell_scoped "$enum_doc"
expect "strict tool crash rejected" fail crash run_misspell_scoped "$enum_doc"
expect "strict usage error rejected" fail usage run_misspell_scoped "$enum_doc"

# changed-docs mode: newline separated document list
expect "changed docs multi-line list accepted" pass literals run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs prose rejected" fail prose run_misspell_changed_docs "$(printf '%s\n%s' "$enum_doc" "$other_doc")"
expect "changed docs empty tool failure rejected" fail empty run_misspell_changed_docs "$enum_doc"

if [ "$failures" -ne 0 ]; then
  echo "spelling scope test: $failures failure(s)"
  exit 1
fi
echo "spelling scope test: all cases passed"

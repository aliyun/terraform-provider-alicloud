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

if [ "$MISSPELL_ENUM_DOC" != "$enum_doc" ]; then
  echo "FAIL: the exemption page must be the fixed documented page, got $MISSPELL_ENUM_DOC"
  exit 1
fi

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
  duplicate_sections)
    # Another paragraph, an appendix or a fenced block of the same page repeats
    # the enum phrase at the shape of the status argument. Only the Argument
    # Reference definition may be exempt, so these findings stay reported.
    if [ "$is_enum_page" -eq 1 ]; then
      printf '%s\n' \
        "$last:${FAKE_DUPLICATE_LINE}:${FAKE_DUPLICATE_COL_FIRST}: \"CANCELLING\" is a misspelling of \"CANCELING\"" \
        "$last:${FAKE_DUPLICATE_LINE}:${FAKE_DUPLICATE_COL_SECOND}: \"CANCELLED\" is a misspelling of \"CANCELED\""
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
expect "changed docs tool crash rejected" fail crash run_misspell_changed_docs "$enum_doc"
expect "changed docs usage error rejected" fail usage run_misspell_changed_docs "$enum_doc"

# The page itself publishes an Example Usage fenced block before the argument
# list. The parser must still find exactly the two status literals of the
# Argument Reference argument.
argument_reference_line="$(grep -n '^## Argument Reference' "$enum_doc" | head -1 | cut -d: -f1)"
if [ -z "$argument_reference_line" ]; then
  echo "FAIL: the Argument Reference section was not found in $enum_doc"
  failures=$((failures + 1))
fi

# A duplicate of the status argument outside the Argument Reference section (an
# appendix), and one inside a fenced block of that section, must not be exempt.
duplicate_line='* `status` - (Optional, String) Valid values: `CANCELLING` and `CANCELLED`.'
duplicate_col_first="$(awk -v line="$duplicate_line" 'BEGIN { print index(line, "`CANCELLING`") }')"
duplicate_col_second="$(awk -v line="$duplicate_line" 'BEGIN { print index(line, "`CANCELLED`") }')"
export FAKE_DUPLICATE_COL_FIRST="$duplicate_col_first"
export FAKE_DUPLICATE_COL_SECOND="$duplicate_col_second"

paragraph_fixture="$tmpdir/paragraph/ecs_image_pipeline_executions.html.markdown"
mkdir -p "$(dirname "$paragraph_fixture")"
cp "$enum_doc" "$paragraph_fixture"
paragraph_line_number=$(( $(wc -l < "$paragraph_fixture") + 1 ))
printf '%s\n' "$duplicate_line" >> "$paragraph_fixture"

appendix_fixture="$tmpdir/appendix/ecs_image_pipeline_executions.html.markdown"
mkdir -p "$(dirname "$appendix_fixture")"
{
  cat "$enum_doc"
  printf '%s\n' '## Appendix' "$duplicate_line"
} > "$appendix_fixture"
appendix_line_number=$(( $(wc -l < "$enum_doc") + 2 ))

fence_fixture="$tmpdir/fence/ecs_image_pipeline_executions.html.markdown"
mkdir -p "$(dirname "$fence_fixture")"
awk -v heading="$argument_reference_line" -v line="$duplicate_line" '
  { print }
  NR == heading { print "```hcl"; print line; print "```" }
' "$enum_doc" > "$fence_fixture"
fence_line_number=$(( argument_reference_line + 2 ))

for fixture in "paragraph:$paragraph_fixture:$paragraph_line_number" "appendix:$appendix_fixture:$appendix_line_number" "fence:$fence_fixture:$fence_line_number"; do
  fixture_name="${fixture%%:*}"
  fixture_rest="${fixture#*:}"
  fixture_doc="${fixture_rest%%:*}"
  fixture_line="${fixture_rest##*:}"
  positions="$(misspell_enum_positions_from < "$fixture_doc")"
  position_count="$(printf '%s\n' "$positions" | grep -c .)"
  if [ "$position_count" -ne 2 ]; then
    echo "FAIL: the $fixture_name duplicate changed the accepted positions of $fixture_doc: $positions"
    failures=$((failures + 1))
  fi
  if printf '%s\n' "$positions" | grep -qx "${fixture_line}:${duplicate_col_first}"; then
    echo "FAIL: the $fixture_name duplicate was exempted"
    failures=$((failures + 1))
  fi
  export FAKE_DUPLICATE_LINE="$fixture_line"
  expect "strict $fixture_name copy is not exempt" fail duplicate_sections run_misspell_scoped "$fixture_doc"
  expect "changed docs $fixture_name copy is not exempt" fail duplicate_sections run_misspell_changed_docs "$fixture_doc"
done

if [ "$failures" -ne 0 ]; then
  echo "spelling scope test: $failures failure(s)"
  exit 1
fi
echo "spelling scope test: all cases passed"

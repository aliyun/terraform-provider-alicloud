#!/bin/bash
# Documentation spelling helpers shared by the strict and changed-docs checks.
# Sourced by scripts/local-ci-check.sh and exercised by scripts/misspell-scope_test.sh.

MISSPELL_ENUM_DOC="website/docs/d/ecs_image_pipeline_executions.html.markdown"

# DescribeImagePipelineExecutions publishes CANCELLING and CANCELLED as status
# literals, and misspell flags them. Only the argument line that declares that
# enum is accepted, and only in that one page: the same words in ordinary prose
# on the same page, and every finding anywhere else, stay reported.
MISSPELL_ENUM_LINE='^\* `status` - \(Optional, String\) The execution status to filter by\. .*`CANCELLING` and `CANCELLED`\.'

run_misspell_scoped() {
  local output status=0 remaining awk_status=0 enum_lines
  output=$(misspell -error -locale=US "$@" 2>&1) || status=$?
  if [ $status -eq 0 ]; then
    return 0
  fi
  # misspell -error exits 2 when it reported misspellings (its documented
  # -error contract). Any other status, an empty finding list, or a failing
  # filter is a tool failure and must not be reported as an accepted exception.
  if [ $status -ne 2 ]; then
    printf '%s\n' "$output"
    echo "misspell failed with exit status ${status}; the spelling check did not complete."
    return 1
  fi
  if [ -z "$output" ]; then
    echo "misspell exited with a finding status but printed nothing; the spelling check did not complete."
    return 1
  fi
  # Line numbers of the API enum declaration. The exemption is keyed on the
  # source position, not on the file alone, so an ordinary uppercase sentence
  # elsewhere on the same page is still reported.
  enum_lines=$(grep -nE "$MISSPELL_ENUM_LINE" "$MISSPELL_ENUM_DOC" 2>/dev/null | cut -d: -f1 | sort -n | paste -sd, -)
  remaining=$(printf '%s\n' "$output" | awk -v doc="$MISSPELL_ENUM_DOC" -v lines="$enum_lines" '
    BEGIN {
      count = split(lines, accepted, ",")
      for (i = 1; i <= count; i++) exempt[accepted[i]] = 1
    }
    {
      line = $0
      sub(/^\.\//, "", line)
      n = length(doc)
      if (substr(line, 1, n) == doc) {
        rest = substr(line, n + 1)
        if (match(rest, /^:[0-9]+:[0-9]+: "(CANCELLING|CANCELLED)" is a misspelling of "(CANCELING|CANCELED)"$/)) {
          split(substr(rest, 2), parts, ":")
          if (parts[1] in exempt) {
            next
          }
        }
      }
      print $0
    }') || awk_status=$?
  if [ $awk_status -ne 0 ]; then
    printf '%s\n' "$output"
    echo "the spelling exemption filter failed with exit status ${awk_status}; the spelling check did not complete."
    return 1
  fi
  if [ -z "$remaining" ]; then
    printf '%s\n' "$output"
    echo "Accepted the documented API status literals CANCELLING/CANCELLED on the API enum line."
    return 0
  fi
  printf '%s\n' "$remaining"
  return 1
}

# The changed-document list is newline separated, so it is read line by line
# instead of being expanded into a single `for` word list.
run_misspell_changed_docs() {
  local doc
  while IFS= read -r doc; do
    [ -n "$doc" ] || continue
    run_misspell_scoped "$doc" || return 1
  done <<< "$1"
  return 0
}

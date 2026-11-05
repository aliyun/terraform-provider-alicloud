#!/bin/bash
# Documentation spelling helpers shared by the strict and changed-docs checks.
# Sourced by scripts/local-ci-check.sh and exercised by scripts/misspell-scope_test.sh.

# The page that publishes the ECS status literals. It is the only page whose
# `status` argument may be exempted; the override exists so the regression test
# can exercise the same logic against a copy of the page.
: "${MISSPELL_ENUM_DOC:=website/docs/d/ecs_image_pipeline_executions.html.markdown}"

# DescribeImagePipelineExecutions publishes CANCELLING and CANCELLED as status
# literals, and misspell flags them. Only those literals in the `Valid values:`
# list of the `status` argument on that one page are accepted: the exemption is
# keyed on the exact line:column misspell reports, so the same word as ordinary
# prose on the same line, the same literals on another page, the same phrase in
# another paragraph of the same page, and every finding anywhere else stay
# reported.

# misspell_enum_positions prints one `line:column` pair per accepted occurrence.
# misspell reports the 0-based byte column of the word inside the line, so the
# column is the offset of the literal after its opening backtick. Only the top
# level `status` argument definition is read: an appendix, an example or another
# argument that happens to repeat the phrase does not grant an exemption.
misspell_enum_positions() {
  awk '
    /^\* `status` - / && /Valid values:/ {
      start = index($0, "Valid values:")
      if (start == 0) next
      segment = substr($0, start)
      stop = index(segment, ".")
      if (stop == 0) next
      segment = substr(segment, 1, stop - 1)
      base = start - 1
      consumed = 0
      while (match(segment, /`(CANCELLING|CANCELLED)`/)) {
        printf "%d:%d\n", NR, base + consumed + RSTART
        advance = RSTART + RLENGTH - 1
        consumed += advance
        segment = substr(segment, advance + 1)
      }
    }
  ' "$MISSPELL_ENUM_DOC"
}

run_misspell_scoped() {
  local output status=0 remaining awk_status=0 enum_positions
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
  # Source positions of the API enum declaration. The exemption is keyed on both
  # the line and the column, not on the file or the line alone, so an ordinary
  # uppercase sentence that happens to share the line is still reported.
  enum_positions=$(misspell_enum_positions | paste -sd' ' -)
  remaining=$(printf '%s\n' "$output" | awk -v doc="$MISSPELL_ENUM_DOC" -v positions="$enum_positions" '
    BEGIN {
      count = split(positions, accepted, " ")
      for (i = 1; i <= count; i++) {
        if (accepted[i] != "") exempt[accepted[i]] = 1
      }
    }
    {
      line = $0
      sub(/^\.\//, "", line)
      n = length(doc)
      if (substr(line, 1, n) == doc) {
        rest = substr(line, n + 1)
        if (match(rest, /^:[0-9]+:[0-9]+: "(CANCELLING|CANCELLED)" is a misspelling of "(CANCELING|CANCELED)"$/)) {
          split(substr(rest, 2), parts, ":")
          if ((parts[1] ":" parts[2]) in exempt) {
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
    echo "Accepted the documented API status literals CANCELLING/CANCELLED at the enum list positions."
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

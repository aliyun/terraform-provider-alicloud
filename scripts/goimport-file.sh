#!/bin/bash

set -euo pipefail

file="$1"

[[ -f "$file" ]] || exit 0

# Step 1: Remove empty lines inside import(...) block
awk '
/^import \($/ { in_import=1; print; next }
in_import && /^\)[[:space:]]*$/ { in_import=0; print; next }
in_import { if (NF == 0 || /^[[:space:]]*$/) next; else print; next }
{ print }
' "$file" > "$file.tmp" && mv "$file.tmp" "$file"

# Step 2: Run goimports (-local not supported)
goimports -w "$file"
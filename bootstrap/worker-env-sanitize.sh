#!/usr/bin/env bash
# Copy one dotenv-style runtime file for an ordinary worker while removing the
# operator-only control-plane admin credential.  Never print input lines: they
# may contain secrets.

set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "usage: $0 SOURCE DESTINATION" >&2
  exit 64
fi

source_file="$1"
destination="$2"
[ -f "$source_file" ] || {
  echo "worker-env-sanitize: source file is missing" >&2
  exit 66
}

mkdir -p "$(dirname "$destination")"
temporary="$(mktemp "${destination}.tmp.XXXXXX")"
trap 'rm -f "$temporary"' EXIT

# Runtime env files use one assignment per line. Drop every line mentioning the
# privileged variable (assignment, export, or comment) so it cannot enter the
# encrypted worker bundle. Other worker credentials remain unchanged.
sed '/JARVIS_CONTROL_PLANE_ADMIN_TOKEN/d' "$source_file" >"$temporary"
if grep -q 'JARVIS_CONTROL_PLANE_ADMIN_TOKEN' "$temporary"; then
  echo "worker-env-sanitize: admin credential marker survived filtering" >&2
  exit 65
fi

chmod 600 "$temporary"
mv "$temporary" "$destination"
trap - EXIT

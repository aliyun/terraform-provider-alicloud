#!/usr/bin/env bash
# Administrator-owned runner bootstrap; verify official archives before marking a tool complete.
set -euo pipefail
version="${1:?Go version required}"
case "$version" in
  1.25.8) checksum=ceb5e041bbc3893846bd1614d76cb4681c91dadee579426cf21a63f2d7e03be6 ;;
  1.26.10) checksum=ccfe7ae822d2c2f41f4fe7c956199967a4dac22d0ee0d2b2bfdb440aab3b20c5 ;;
  *) echo "No administrator seed for Go $version"; exit 0 ;;
esac
root="${RUNNER_TOOL_CACHE:-/home/runner/_work/_tool}"
target="$root/go/$version/x64"
if [[ -x "$target/bin/go" && -f "$target.complete" ]]; then
  exit 0
fi
temporary="$(mktemp -d)"
trap 'rm -rf "$temporary"' EXIT
curl --fail --show-error --connect-timeout 3 --max-time 25 --retry 4 --retry-delay 1 --retry-all-errors \
  "http://go-seed.terraform-ci-cache-maintenance.svc.cluster.local:8080/tools/go$version.linux-amd64.tar.gz" \
  -o "$temporary/go.tar.gz"
printf '%s  %s\n' "$checksum" "$temporary/go.tar.gz" | sha256sum -c -
mkdir -p "$target"
tar -xzf "$temporary/go.tar.gz" --strip-components=1 -C "$target"
[[ "$("$target/bin/go" version)" == "go version go$version linux/amd64" ]]
# @actions/tool-cache only discovers complete <tool>/<version>/<arch> entries.
touch "$target.complete"
echo "GO_TOOLCACHE_READY $version"

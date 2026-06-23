#!/bin/bash
# Test cache.sh: hit avoids re-run, bust forces re-run, TTL expiry re-runs, fresh probe.
set -uo pipefail
test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
cache="$proj_root/bootstrap/cache.sh"

tmp=$(mktemp -d); trap "rm -rf $tmp" EXIT
export JARVIS_CACHE_DIR="$tmp/cache"
counter="$tmp/n"; echo 0 > "$counter"
gen() { n=$(($(cat "$counter")+1)); echo "$n" > "$counter"; echo "val$n"; }
export -f gen; export counter

fail=0
a=$(bash "$cache" get k 60 -- bash -c gen)
b=$(bash "$cache" get k 60 -- bash -c gen)
[ "$a" = "val1" ] && [ "$b" = "val1" ] && echo "PASS hit" || { echo "FAIL hit: $a/$b"; fail=1; }

bash "$cache" bust k
c=$(bash "$cache" get k 60 -- bash -c gen)
[ "$c" = "val2" ] && echo "PASS bust" || { echo "FAIL bust: $c"; fail=1; }

bash "$cache" fresh k 60 && echo "PASS fresh" || { echo "FAIL fresh"; fail=1; }
d=$(bash "$cache" get k 0 -- bash -c gen)
[ "$d" = "val3" ] && echo "PASS ttl0" || { echo "FAIL ttl0: $d"; fail=1; }

# failing cmd must not poison cache
bash "$cache" bust k
bash "$cache" get k 60 -- false >/dev/null 2>&1
bash "$cache" fresh k 60 && { echo "FAIL no-poison"; fail=1; } || echo "PASS no-poison"

[ $fail -eq 0 ] && echo "✓ cache tests passed" || exit 1

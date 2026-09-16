#!/usr/bin/env bash
# Qoder interactive-session support: process matching, runtime-context
# resolution, and claim.sh's interactive gate for Qoder sessions.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
worker="$repo_root/bootstrap/jarvis-interactive-worker.py"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bin"

pass=0
fail=0
ok() {
  if [ "$1" = "0" ]; then
    pass=$((pass + 1))
    echo "PASS: $2"
  else
    fail=$((fail + 1))
    echo "FAIL: $2" >&2
  fi
}

# --- 1. _client_process_matches: Qoder host forms -------------------------
out="$(/usr/bin/python3 - "$worker" <<'PY'
import importlib.util, sys
spec = importlib.util.spec_from_file_location("worker", sys.argv[1])
mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)
m = mod._client_process_matches
cases = [
    ("/applications/qoder.app/contents/macos/qoder", "qoder", True),
    ("/users/x/.qoder/bin/qodercli --flag", "qoder", True),
    ("/users/x/.qoder/bin/qoder-computer-use mcp", "qoder", False),
    ("/users/x/.loongsuite-pilot/hooks/qoderwork-runtime-wrapper.mjs", "qoder", False),
    ("/applications/qoder.app/.../qoder-search.bundle.mjs mcp-bridge", "qoder", False),
    ("/usr/local/bin/claude --resume", "claude", True),
    ("/users/x/.codex/bin/codex run", "codex", True),
    ("/usr/bin/python3 worker.py", "qoder", False),
]
bad = [c for c in cases if bool(m(c[0], c[1])) != c[2]]
print("BAD" if bad else "GOOD")
PY
)"
[ "$out" = "GOOD" ]
ok $? "_client_process_matches: qoder basename forms accepted, helpers rejected"

# --- 2. _runtime_context: persisted qoder env resolves ---------------------
out="$(JARVIS_INTERACTIVE_CLIENT=qoder JARVIS_INTERACTIVE_SESSION_ID=sid-123 \
  /usr/bin/python3 - "$worker" <<'PY'
import importlib.util, sys
spec = importlib.util.spec_from_file_location("worker", sys.argv[1])
mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)
print("%s,%s" % mod._runtime_context())
PY
)"
[ "$out" = "qoder,sid-123" ]
ok $? "_runtime_context resolves persisted qoder env"

# --- 3. claim.sh interactive path when runtime-context yields qoder --------
cat > "$tmp/config-pools.json" <<'JSON'
{
  "claim": {
    "tag": "jarvis-claimed",
    "idle_tag": "jarvis-idle",
    "done_tag": "jarvis-done",
    "done_status": "已完成",
    "done_statuses": ["已完成"],
    "progress_status": "处理中",
    "start_statuses": ["待处理"]
  },
  "pools": []
}
JSON
mkdir -p "$tmp/config" && cp "$tmp/config-pools.json" "$tmp/config/pools.json"

cat > "$tmp/bin/a1" <<'STUB'
#!/usr/bin/env bash
echo "a1:$*" >> "$TEST_LOG"
if [ "$1 $2 $3" = "project workitem get" ]; then
  tags="$(cat "$A1_STATE" 2>/dev/null || true)"
  tags_ml="${tags//,/, }"
  printf '{"subject":"t","fields":[{"identifier":"tag","displayValue":"%s","value":"%s"},{"identifier":"workitemType","displayValue":"需求"},{"identifier":"status","displayValue":"处理中"}]}\n' "$tags_ml" "$tags_ml"
  exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
  args=("$@"); tag=""
  for ((i=0; i<${#args[@]}; i++)); do
    [ "${args[$i]}" = "--tag" ] && tag="${args[$((i+1))]}"
  done
  [ -n "$tag" ] && printf '%s' "$tag" > "$A1_STATE"
  exit 0
fi
exit 0
STUB
chmod +x "$tmp/bin/a1"

cat > "$tmp/runner.sh" <<'STUB'
#!/usr/bin/env bash
echo "worker:$*" >> "$TEST_LOG"
[ "$1" = "cli" ] && shift
case "${1:-}" in
  runtime-context) echo "qoder test-sid-1"; exit 0 ;;
  prepare-claim) printf '{"accepted":true,"proceed":true,"operationStatus":"SENDING"}\n'; exit 0 ;;
  operation-ack) exit 0 ;;
  operation-fail) exit 0 ;;
  *) exit 64 ;;
esac
STUB
chmod +x "$tmp/runner.sh"

run_claim() {
  (
    export JARVIS_ROOT="$tmp"
    export JARVIS_A1="$tmp/bin/a1"
    export JARVIS_INTERACTIVE_WORKER_RUNNER="$tmp/runner.sh"
    export JARVIS_CLAIM_SETTLE=0
    export TEST_LOG="$tmp/events.log"
    export A1_STATE="$tmp/a1.state"
    export HOME="$tmp"
    unset CLAUDE_CODE_SESSION_ID CODEX_THREAD_ID JARVIS_INTERACTIVE_CLIENT JARVIS_INTERACTIVE_SESSION_ID
    cd "$tmp"
    bash "$repo_root/bootstrap/claim.sh" claim 424242 1086837 >/dev/null 2>&1 || true
  )
}

: > "$tmp/events.log"
printf 'jarvis-claimed' > "$tmp/a1.state"
run_claim
grep -q "worker:cli prepare-claim" "$tmp/events.log"
ok $? "claim.sh runs interactive prepare-claim under qoder runtime-context"

# --- 4. negative: runtime-context failing must not enter interactive path --
cat > "$tmp/runner.sh" <<'STUB'
#!/usr/bin/env bash
echo "worker:$*" >> "$TEST_LOG"
[ "$1" = "cli" ] && shift
case "${1:-}" in
  runtime-context) exit 1 ;;
  prepare-claim) echo "MUST-NOT-CALL" >> "$TEST_LOG"; exit 0 ;;
  *) exit 64 ;;
esac
STUB
chmod +x "$tmp/runner.sh"
: > "$tmp/events.log"
printf 'jarvis-claimed' > "$tmp/a1.state"
run_claim
if grep -q "MUST-NOT-CALL\|worker:cli prepare-claim" "$tmp/events.log"; then exit 1; fi
ok $? "claim.sh skips interactive path when qoder context unresolvable"

echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]

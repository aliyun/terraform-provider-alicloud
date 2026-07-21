#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCRIPT="$ROOT/bootstrap/notify-dingtalk.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
mkdir -p "$TMP/config"

cat > "$TMP/streaming.py" <<'PY'
import argparse, json, os, sys
p = argparse.ArgumentParser()
p.add_argument("--to")
p.add_argument("--template-id")
p.add_argument("--no-stream", action="store_true")
p.add_argument("-m")
p.add_argument("--out-track-id")
a = p.parse_args()
print(json.dumps({"status": "created", "outTrackId": a.out_track_id}), flush=True)
if os.environ.get("FAKE_FAIL") == "1":
    sys.exit(1)
print(json.dumps({"status": "done", "outTrackId": a.out_track_id}), flush=True)
PY

base_env=(
  "JARVIS_ROOT=$TMP"
  "DINGTALK_APP_KEY=key"
  "DINGTALK_APP_SECRET=secret"
  "DINGTALK_TEMPLATE_ID=template"
  "DINGTALK_SKILL=$TMP/streaming.py"
)

out="$(env "${base_env[@]}" "$SCRIPT" --result-json --out-track-id stable-1 \
  484483 "title" "body")"
status="$(printf '%s\n' "$out" | tail -n 1 | jq -r .status)"
receipt="$(printf '%s\n' "$out" | tail -n 1 | jq -r .receipt)"
[ "$status" = "sent" ] && [ "$receipt" = "stable-1" ]

printf '484483\n' > "$TMP/config/dingtalk-optout.txt"
out="$(env "${base_env[@]}" "$SCRIPT" --result-json --out-track-id stable-2 \
  484483 "title" "body")"
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .status)" = "skipped" ]
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .reason)" = "optout" ]
rm "$TMP/config/dingtalk-optout.txt"

out="$(env JARVIS_ROOT="$TMP" DINGTALK_SKILL="$TMP/streaming.py" \
  "$SCRIPT" --result-json --out-track-id stable-3 484483 "title" "body")"
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .status)" = "failed" ]
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .reason)" = "missing_credentials" ]

out="$(env "${base_env[@]}" FAKE_FAIL=1 "$SCRIPT" --result-json \
  --out-track-id stable-4 484483 "title" "body")"
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .status)" = "failed" ]
[ "$(printf '%s\n' "$out" | tail -n 1 | jq -r .receipt)" = "stable-4" ]

# Legacy callers keep their old no-JSON stdout behavior and zero exit status.
legacy="$(env "${base_env[@]}" "$SCRIPT" 484483 "title" "body")"
[ -z "$legacy" ]

echo "notify_dingtalk_test: PASS"

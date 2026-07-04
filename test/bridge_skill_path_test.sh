#!/usr/bin/env bash
# test/bridge_skill_path_test.sh — skill_path() resolution for the DingTalk bridge.
#
# Guards the vendored-fallback fix: when the global ~/.claude copy of streaming.py is
# absent, skill_path() must fall back to the repo-vendored copy instead of returning a
# missing global path (which disabled all replies). Imports the bridge module directly;
# skips gracefully when the dingtalk_stream dependency is unavailable (non-bridge hosts).

set -uo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

python3 - "$repo_root" <<'PY'
import sys, os, tempfile
from pathlib import Path

repo_root = sys.argv[1]
sys.path.insert(0, os.path.join(repo_root, "bridge"))
try:
    import jarvis_dingtalk_bot as b
except Exception as e:  # dingtalk_stream not installed on this host → not the bridge machine
    print(f"SKIP: bridge module not importable here ({e})")
    sys.exit(0)

fails = 0
def ck(name, cond, extra=""):
    global fails
    if cond:
        print(f"PASS {name}")
    else:
        print(f"FAIL {name} {extra}")
        fails += 1

# 1. explicit DINGTALK_SKILL override wins
os.environ["DINGTALK_SKILL"] = "/tmp/override_streaming.py"
ck("override respected", str(b.skill_path()) == "/tmp/override_streaming.py", str(b.skill_path()))
os.environ.pop("DINGTALK_SKILL", None)

vendored = b.REPO_ROOT / ".claude" / "skills" / "dingtalk-ai-card" / "scripts" / "streaming.py"

# 2. global absent → fall back to the repo-vendored copy (and it exists)
b.DEFAULT_SKILL = Path("/definitely/nonexistent/streaming.py")
sp = b.skill_path()
ck("vendored fallback when global absent", sp == vendored, f"got {sp}")
ck("vendored fallback path exists", sp.exists(), f"missing {sp}")

# 3. global present → prefer the global copy over vendored
with tempfile.NamedTemporaryFile(suffix="_streaming.py", delete=False) as tf:
    tmp = Path(tf.name)
b.DEFAULT_SKILL = tmp
ck("prefer global when present", b.skill_path() == tmp, str(b.skill_path()))
tmp.unlink(missing_ok=True)

print("ALLPASS" if fails == 0 else f"FAILED ({fails})")
sys.exit(1 if fails else 0)
PY

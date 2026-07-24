#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

BRIDGE_DIR="$repo_root/bridge" python3 - <<'PY'
import inspect
import os
import sys

sys.path.insert(0, os.environ["BRIDGE_DIR"])
import jarvis_dingtalk_bot as bridge


def require(prompt, *terms):
    for term in terms:
        if term not in prompt:
            raise AssertionError(f"prompt missing {term!r}")


common_terms = (
    "Terraform 三层可视化证据契约",
    "visual_evidence_manifest",
    "OpenAPI",
    "CloudSpec/ACube",
    "Provider",
    "screenshot-evidence",
    "validate-manifest.py",
    "AONE_RESULT.reply_body",
    "严禁传 `--comment`",
    "禁止 claim/wrap/release/直接评论",
)

ticket = bridge._ticket_prompt(
    "12345678", "restore screenshot evidence", "tf_provider", "528766")
require(ticket, *common_terms)

ordinary = bridge._ticket_prompt(
    "12345678", "ordinary task", "api_toolkit", "2100304")
if "Terraform 三层可视化证据契约" in ordinary:
    raise AssertionError("non-Terraform ticket unexpectedly received Terraform evidence contract")

for role in ("terraform-pd", "terraform-rd", "terraform-qa"):
    persona = bridge._persona_prompt(
        "12345678", role, "continue", "note", 2, "snippet")
    require(persona, *common_terms)
    require(persona, "缺失或无效时先补起 terraform-pd")

wake_tail = bridge._task_result_instructions("12345678", True)
require(wake_tail, *common_terms)
wake_source = inspect.getsource(bridge.JarvisHandler._wake)
require(wake_source, "_task_result_instructions(aone_id, tf)")

print("aone_screenshot_bridge_prompt_test: PASS")
PY

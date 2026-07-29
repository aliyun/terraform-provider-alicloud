#!/usr/bin/env bash
# Repo-controlled browser screenshot capture for headless sessions.
#
# The interactive Playwright MCP (`mcp__playwright__*`) that this skill
# historically used is a per-user interactive-environment artifact and is NOT
# injected into headless bridge-launched sessions (see
# references/headless-screenshot-channels.md for the root cause). This wrapper
# exposes a stable, degradable channel that the agent can call from any cwd:
# Playwright Python binding first, then a headless Chrome/Chromium binary.
# When neither is available it exits 3 with a diagnosable missing_capability
# line — never a silent skip.
#
# Usage:
#   capture.sh probe
#       stdout: channel name, or "missing_capability: <reason>"
#       exit 0 = channel available; exit 3 = missing_capability
#   capture.sh capture <url> <out.png> \
#       [--wait N] [--full-page|--viewport] [--width W] [--height H]
#       stdout: channel name used
#       exit 0 = captured; exit 3 = missing_capability; exit 1 = capture_error
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/../../../.." && pwd)"
python_bin="${JARVIS_PYTHON_BIN:-python3}"

# Make THIS repo's bridge/ authoritative on the import path so a headless run
# started from any cwd (or under an inherited PYTHONPATH) still resolves the
# repo-controlled module rather than a sibling install.
export PYTHONPATH="$repo_root${PYTHONPATH:+:$PYTHONPATH}"
cd "$repo_root"

exec "$python_bin" -m bridge.jarvis_screenshot "$@"

#!/usr/bin/env bash
# Provision the isolated Python runtime used by bridge/run.sh.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BASE_PYTHON="${JARVIS_BRIDGE_BOOTSTRAP_PYTHON:-python3}"
VENV_DIR="${JARVIS_BRIDGE_VENV:-$REPO_ROOT/.venv/bridge}"

"$BASE_PYTHON" -m venv "$VENV_DIR"
"$VENV_DIR/bin/python" -m pip install --requirement "$REPO_ROOT/bridge/requirements.txt"
printf 'Bridge Python ready: %s\n' "$VENV_DIR/bin/python"

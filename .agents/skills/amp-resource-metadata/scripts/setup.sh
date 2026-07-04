#!/usr/bin/env bash
# Setup script for amp-resource-metadata skill
# Creates a virtual environment and installs dependencies.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Creating virtual environment in ${SCRIPT_DIR}/.venv ..."
python3 -m venv "${SCRIPT_DIR}/.venv"

echo "==> Activating virtual environment ..."
source "${SCRIPT_DIR}/.venv/bin/activate"

echo "==> Installing dependencies ..."
pip install --upgrade pip -q
pip install -r "${SCRIPT_DIR}/requirements.txt" -q

echo "==> Setup complete."
echo "    Activate with: source ${SCRIPT_DIR}/.venv/bin/activate"

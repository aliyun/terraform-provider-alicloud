#!/usr/bin/env bash
# Hermetic contract tests for the project-scoped AutoWonder MCP client config.
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

python3 - "$repo_root" <<'PY'
import json
import re
import sys
from pathlib import Path

root = Path(sys.argv[1])
endpoint = "https://auto-wonder.alibaba.net/api/mcp"
token_env = "AUTOWONDER_MCP_TOKEN"

claude = json.loads((root / ".mcp.json").read_text(encoding="utf-8"))
server = claude["mcpServers"]["autowonder"]
assert server == {
    "type": "http",
    "url": endpoint,
    "headers": {"Authorization": f"Bearer ${{{token_env}:-}}"},
}, server

settings = json.loads(
    (root / ".claude" / "settings.json").read_text(encoding="utf-8"))
assert "autowonder" in settings.get("enabledMcpjsonServers", []), settings
assert settings.get("enableAllProjectMcpServers") is not True, settings

codex = (root / ".codex" / "config.toml").read_text(encoding="utf-8")
assert re.search(r"(?m)^\[mcp_servers\.autowonder\]\s*$", codex), codex
assert re.search(rf'(?m)^url\s*=\s*"{re.escape(endpoint)}"\s*$', codex), codex
assert re.search(
    rf'(?m)^bearer_token_env_var\s*=\s*"{token_env}"\s*$', codex), codex
assert re.search(r"(?m)^enabled\s*=\s*true\s*$", codex), codex
assert "http_headers" not in codex, codex

tracked_configs = [
    root / ".mcp.json",
    root / ".codex" / "config.toml",
    root / ".claude" / "settings.json",
    root / "bootstrap" / ".env.example",
]
secret = re.compile(r"awmcp_[A-Za-z0-9_-]{43}")
for path in tracked_configs:
    assert not secret.search(path.read_text(encoding="utf-8")), path

print("autowonder_mcp_config_test: PASS")
PY

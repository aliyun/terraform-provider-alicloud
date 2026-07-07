#!/usr/bin/env python3
"""Smoke-test an Alibaba Cloud OpenAPI MCP connection end-to-end.

Drives a real MCP `initialize` + `tools/list` and prints the tool count or the
exact upstream error (e.g. 401 "Failed to exchange token"). Use it to verify a
server BEFORE wiring it into Claude, and to diagnose auth failures.

Usage:
  # Approach B (uvx proxy + aliyun CLI profile): spawns the proxy with the
  # profile and empty AK/SK (so the CLI profile wins the credential chain).
  python3 smoke-test.py <profile-name>

  # Also prints the "Discovered MCP server URL" for that profile when --debug
  # writes a log; pass --show-url to surface it.
  python3 smoke-test.py <profile-name> --show-url

Notes:
  - Approach A (mcp-remote/OAuth) has no static credential to inject here; verify
    it with a raw curl using a token from ~/.mcp-auth (see SKILL.md), or just
    connect it in Claude and check /mcp.
  - Exit code 0 = tools listed; 1 = error/no result.
"""
import json
import os
import subprocess
import sys
import threading
import time

def main() -> int:
    if len(sys.argv) < 2:
        print(__doc__)
        return 1
    profile = sys.argv[1]
    show_url = "--show-url" in sys.argv[2:]

    env = dict(os.environ)
    env["ALIBABA_CLOUD_PROFILE"] = profile
    # Neutralize any globally-exported AK/SK so the CLI profile is used.
    env["ALIBABA_CLOUD_ACCESS_KEY_ID"] = ""
    env["ALIBABA_CLOUD_ACCESS_KEY_SECRET"] = ""

    args = ["uvx", "alibabacloud.mcp-proxy@latest"]
    log = None
    if show_url:
        log = f"/tmp/aliyun-mcp-smoke-{profile}.log"
        open(log, "w").close()
        args += ["--debug", "--log-file", log]

    p = subprocess.Popen(
        args, stdin=subprocess.PIPE, stdout=subprocess.PIPE,
        stderr=subprocess.PIPE, text=True, env=env, bufsize=1,
    )
    for m in (
        {"jsonrpc": "2.0", "id": 1, "method": "initialize",
         "params": {"protocolVersion": "2024-11-05", "capabilities": {},
                    "clientInfo": {"name": "smoke", "version": "0"}}},
        {"jsonrpc": "2.0", "method": "notifications/initialized"},
        {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}},
    ):
        p.stdin.write(json.dumps(m) + "\n")
        p.stdin.flush()

    out: list[str] = []
    threading.Thread(
        target=lambda: [out.append(l.rstrip()) for l in p.stdout], daemon=True
    ).start()
    time.sleep(22)
    p.terminate()
    try:
        p.wait(5)
    except Exception:
        p.kill()

    if show_url and log:
        for line in open(log):
            if "Discovered MCP server URL:" in line:
                print("discovered:", line.split("Discovered MCP server URL:")[1].strip())

    rc = 1
    for line in out:
        try:
            o = json.loads(line)
        except Exception:
            continue
        if o.get("id") == 2:
            if "result" in o:
                tools = o["result"].get("tools", [])
                print(f"[{profile}] OK — {len(tools)} tools; "
                      f"sample: {[t['name'] for t in tools[:6]]}")
                rc = 0
            elif "error" in o:
                print(f"[{profile}] ERROR: "
                      f"{o['error'].get('message','').splitlines()[0]}")
    if rc:
        print(f"[{profile}] no successful tools/list result")
    return rc

if __name__ == "__main__":
    raise SystemExit(main())

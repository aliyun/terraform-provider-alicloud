#!/usr/bin/env bash
# bootstrap/serve.sh — tiny stdlib http server for the board.
# GET / → docs/board.html (未生成则自动 build) ; POST /refresh → run board-html.sh --refresh, 200 on success.
# Default port 8787; override: serve.sh [port]. Pure python3 http.server, no deps.
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${JARVIS_ROOT:-$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))}"
port="${1:-8787}"
# 未 build 则起服前先 build 一次（避免首屏 404）；失败也照常起服，可走 /refresh 重试。
[ -f "$root/docs/board.html" ] || { echo "serve: board.html 未生成，首次 build 中…" >&2; bash "$script_dir/board-html.sh" >/dev/null 2>&1 || true; }
echo "serve: http://localhost:$port  (Ctrl-C to stop)"
exec python3 - "$root" "$port" <<'PY'
import sys, os, subprocess
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
root, port = sys.argv[1], int(sys.argv[2])
board = os.path.join(root, "docs", "board.html")
board_html = os.path.join(root, "bootstrap", "board-html.sh")

class H(BaseHTTPRequestHandler):
    def log_message(self, *a): pass
    def _send(self, code, body=b"", ctype="text/plain; charset=utf-8"):
        self.send_response(code); self.send_header("Content-Type", ctype)
        self.send_header("Content-Length", str(len(body))); self.end_headers()
        if body: self.wfile.write(body)
    def do_GET(self):
        if self.path.split("?")[0] in ("/", "/board.html"):
            if not os.path.exists(board):
                subprocess.run(["bash", os.path.join(root, "bootstrap", "board-html.sh")],
                               cwd=root, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            try: b = open(board, "rb").read()
            except OSError: return self._send(503, b"board build failed; check bootstrap/scan.sh & pools.json")
            return self._send(200, b, "text/html; charset=utf-8")
        self._send(404, b"not found")
    def do_POST(self):
        if self.path != "/refresh": return self._send(404, b"not found")
        try:
            subprocess.run(["bash", board_html, "--refresh"], cwd=root, check=True,
                           stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            self._send(200, b"refreshed")
        except Exception as ex:
            self._send(500, ("refresh failed: %s" % ex).encode())

ThreadingHTTPServer(("0.0.0.0", port), H).serve_forever()
PY

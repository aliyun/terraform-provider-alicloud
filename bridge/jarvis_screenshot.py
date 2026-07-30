#!/usr/bin/env python3
"""Repo-controlled browser screenshot capability for headless sessions.

The interactive Playwright MCP (`mcp__playwright__*`) that the screenshot-evidence
skill historically assumed is a per-user interactive-environment artifact: it is
NOT part of the repo-controlled headless launch. ``jarvis_execution_runtime
.jarvis_cmd`` builds ``claude --settings <provider-settings.json> ...`` and the
provider settings files only carry model-provider env (`ANTHROPIC_*`) with no
declared `mcpServers`; there is no repo `.mcp.json` either. In headless sessions
that tool family is therefore absent and the evidence flow cannot capture
pages — which blocks Terraform three-layer verification.

This module provides a stable, detectable, degradable screenshot channel that
does not depend on the interactive MCP. It probes repo-controlled browser
channels in priority order (Playwright Python binding, then a headless
Chrome/Chromium binary) and, when none is available, surfaces a diagnosable
`missing_capability` result instead of failing mid-run or silently skipping the
screenshot step.

Exit codes (shared convention with preflight `missing_capability`):
    0  success / channel available
    3  missing_capability (no usable browser channel)
    1  capture_error (a channel was present but the page capture failed)

CLI:
    python3 -m bridge.jarvis_screenshot probe
        stdout: channel name, or ``missing_capability: <reason>``
    python3 -m bridge.jarvis_screenshot capture <url> <out.png> \
        [--wait N] [--full-page|--viewport] [--width W] [--height H] \
        [--text TARGET]
"""

from __future__ import annotations

import base64
import hashlib
import json
import math
import os
import shutil
import socket
import struct
import subprocess
import sys
import tempfile
import time
import urllib.request
from dataclasses import dataclass
from pathlib import Path
from typing import Callable, List, Optional, Sequence
from urllib.parse import urlparse

EXIT_OK = 0
EXIT_CAPTURE_ERROR = 1
EXIT_MISSING_CAPABILITY = 3

# Prefix that downstream code (skill n-a rows, autonomy escalate triggers) keys on.
MISSING_CAPABILITY_PREFIX = "missing_capability: "
DEFAULT_CHROME_STARTUP_TIMEOUT = 30.0
DEFAULT_CHROME_CAPTURE_ATTEMPTS = 3


class MissingCapability(Exception):
    """No usable browser channel is available on this host."""

    def __init__(self, reason: str, hint: str = ""):
        self.reason = reason
        self.hint = hint
        super().__init__(reason or "no screenshot channel")


class CaptureError(Exception):
    """A channel was available but the page capture failed."""


# ---------------------------------------------------------------------------
# Default capability probes (pure checks — no browser is launched to probe).
# ---------------------------------------------------------------------------


def _default_playwright_importable() -> bool:
    try:
        import playwright  # noqa: F401  — presence check only
        return True
    except Exception:  # noqa: BLE001 — any import failure means unavailable
        return False


def _default_playwright_browser_present() -> bool:
    """True iff a Playwright-managed chromium binary is installed locally."""
    cache = os.environ.get("PLAYWRIGHT_BROWSERS_PATH")
    roots: List[Path] = []
    if cache:
        roots.append(Path(cache))
    else:
        roots.append(Path.home() / "Library" / "Caches" / "ms-playwright")
        # Linux default cache location.
        roots.append(Path.home() / ".cache" / "ms-playwright")
    for root in roots:
        try:
            if root.is_dir():
                for entry in root.iterdir():
                    name = entry.name.lower()
                    if entry.is_dir() and ("chromium" in name or "chrome" in name):
                        return True
        except OSError:
            continue
    return False


_CHROME_BINARY_CANDIDATES = (
    "chromium",
    "chrome",
    "google-chrome",
    "google-chrome-stable",
    "chromium-browser",
)

_MACOS_CHROME_APPS = (
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/Applications/Chromium.app/Contents/MacOS/Chromium",
    "/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
)


def _default_chrome_binary() -> Optional[str]:
    """Find a headless-capable Chrome/Chromium binary, or None."""
    override = os.environ.get("JARVIS_CHROME_BIN", "").strip()
    if override and shutil.which(override):
        return override
    for candidate in _MACOS_CHROME_APPS:
        if os.path.isfile(candidate) and os.access(candidate, os.X_OK):
            return candidate
    for candidate in _CHROME_BINARY_CANDIDATES:
        found = shutil.which(candidate)
        if found:
            return found
    return None


# ---------------------------------------------------------------------------
# Channels
# ---------------------------------------------------------------------------


class ScreenshotChannel:
    """One browser-backed screenshot channel."""

    name: str = "channel"

    def available(self) -> bool:
        raise NotImplementedError

    def unavailability_reason(self) -> str:
        return "not available"

    def capture(self, url: str, out: str, *, wait_ms: int = 3000,
                full_page: bool = True, width: int = 1280,
                height: int = 2000,
                target_text: str = "") -> None:  # pragma: no cover - abstract
        raise NotImplementedError


@dataclass
class PlaywrightPythonChannel(ScreenshotChannel):
    """Playwright Python binding + locally installed chromium binary.

    Preferred over the Chrome binary because Playwright gives element-level
    screenshots and deterministic wait-for-render, which the SPA evidence
    pages (next.api, registry.terraform.io) need.
    """

    name: str = "playwright_python"
    import_checker: Callable[[], bool] = _default_playwright_importable
    browser_cache_checker: Callable[[], bool] = _default_playwright_browser_present

    def available(self) -> bool:
        return self.import_checker() and self.browser_cache_checker()

    def unavailability_reason(self) -> str:
        if not self.import_checker():
            return ("playwright python package not installed "
                    "(pip install playwright)")
        return ("playwright chromium browser not installed "
                "(run: playwright install chromium)")

    def capture(self, url: str, out: str, *, wait_ms: int = 3000,
                full_page: bool = True, width: int = 1280,
                height: int = 2000, target_text: str = "") -> None:
        try:
            from playwright.sync_api import sync_playwright
        except Exception as exc:  # noqa: BLE001
            raise CaptureError(
                "playwright import failed at capture time: %s" % exc) from exc
        try:
            with sync_playwright() as pw:
                browser = pw.chromium.launch()
                try:
                    page = browser.new_page()
                    page.set_viewport_size({"width": width, "height": height})
                    page.goto(url, wait_until="domcontentloaded", timeout=60000)
                    if wait_ms > 0:
                        page.wait_for_timeout(wait_ms)
                    if target_text:
                        target = page.locator(
                            "li, td, tr", has_text=target_text).first
                        target.wait_for(state="visible", timeout=10000)
                        target.screenshot(path=str(out))
                    else:
                        page.screenshot(path=str(out), full_page=full_page)
                finally:
                    browser.close()
        except CaptureError:
            raise
        except Exception as exc:  # noqa: BLE001
            raise CaptureError("playwright capture failed: %s" % exc) from exc


def _recv_exact(sock: socket.socket, size: int, buffered: bytearray) -> bytes:
    """Read exactly ``size`` bytes, consuming handshake leftovers first."""
    while len(buffered) < size:
        chunk = sock.recv(max(4096, size - len(buffered)))
        if not chunk:
            raise CaptureError("chrome devtools websocket closed")
        buffered.extend(chunk)
    data = bytes(buffered[:size])
    del buffered[:size]
    return data


class _WebSocket:
    """Minimal RFC6455 client sufficient for Chrome DevTools Protocol."""

    def __init__(self, websocket_url: str, timeout: float = 60.0):
        parsed = urlparse(websocket_url)
        if parsed.scheme != "ws" or not parsed.hostname or not parsed.port:
            raise CaptureError("invalid chrome devtools websocket URL")
        self.sock = socket.create_connection(
            (parsed.hostname, parsed.port), timeout=timeout)
        self.sock.settimeout(timeout)
        self.buffer = bytearray()
        path = parsed.path or "/"
        if parsed.query:
            path += "?" + parsed.query
        key = base64.b64encode(os.urandom(16)).decode("ascii")
        request = (
            "GET %s HTTP/1.1\r\n"
            "Host: %s:%s\r\n"
            "Upgrade: websocket\r\n"
            "Connection: Upgrade\r\n"
            "Sec-WebSocket-Key: %s\r\n"
            "Sec-WebSocket-Version: 13\r\n\r\n"
        ) % (path, parsed.hostname, parsed.port, key)
        self.sock.sendall(request.encode("ascii"))
        while b"\r\n\r\n" not in self.buffer:
            chunk = self.sock.recv(4096)
            if not chunk:
                raise CaptureError("chrome devtools websocket handshake closed")
            self.buffer.extend(chunk)
        headers, remainder = bytes(self.buffer).split(b"\r\n\r\n", 1)
        self.buffer[:] = remainder
        first = headers.split(b"\r\n", 1)[0]
        if b" 101 " not in first:
            raise CaptureError(
                "chrome devtools websocket handshake failed: %s"
                % first.decode("ascii", "replace"))
        expected = base64.b64encode(hashlib.sha1(
            (key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11").encode(
                "ascii")).digest()).decode("ascii")
        accept = ""
        for line in headers.split(b"\r\n")[1:]:
            if line.lower().startswith(b"sec-websocket-accept:"):
                accept = line.split(b":", 1)[1].strip().decode("ascii")
                break
        if accept != expected:
            raise CaptureError("chrome devtools websocket accept mismatch")

    def close(self) -> None:
        try:
            self._send_frame(b"", opcode=0x8)
        except (OSError, CaptureError):
            pass
        self.sock.close()

    def _send_frame(self, payload: bytes, opcode: int = 0x1) -> None:
        mask = os.urandom(4)
        length = len(payload)
        header = bytearray([0x80 | opcode])
        if length < 126:
            header.append(0x80 | length)
        elif length < 65536:
            header.append(0x80 | 126)
            header.extend(struct.pack("!H", length))
        else:
            header.append(0x80 | 127)
            header.extend(struct.pack("!Q", length))
        header.extend(mask)
        masked = bytes(value ^ mask[index % 4]
                       for index, value in enumerate(payload))
        self.sock.sendall(bytes(header) + masked)

    def send_json(self, payload: dict) -> None:
        self._send_frame(json.dumps(
            payload, separators=(",", ":")).encode("utf-8"))

    def recv_json(self) -> dict:
        fragments = bytearray()
        message_opcode = None
        while True:
            first, second = _recv_exact(self.sock, 2, self.buffer)
            final = bool(first & 0x80)
            opcode = first & 0x0F
            masked = bool(second & 0x80)
            length = second & 0x7F
            if length == 126:
                length = struct.unpack(
                    "!H", _recv_exact(self.sock, 2, self.buffer))[0]
            elif length == 127:
                length = struct.unpack(
                    "!Q", _recv_exact(self.sock, 8, self.buffer))[0]
            mask = (_recv_exact(self.sock, 4, self.buffer)
                    if masked else b"")
            payload = bytearray(_recv_exact(
                self.sock, length, self.buffer))
            if masked:
                for index in range(length):
                    payload[index] ^= mask[index % 4]
            if opcode == 0x8:
                raise CaptureError("chrome devtools websocket closed")
            if opcode == 0x9:
                self._send_frame(bytes(payload), opcode=0xA)
                continue
            if opcode == 0xA:
                continue
            if opcode in (0x1, 0x2):
                message_opcode = opcode
                fragments[:] = payload
            elif opcode == 0x0 and message_opcode is not None:
                fragments.extend(payload)
            else:
                continue
            if final:
                if message_opcode != 0x1:
                    raise CaptureError(
                        "unexpected binary chrome devtools response")
                try:
                    return json.loads(fragments.decode("utf-8"))
                except (UnicodeDecodeError, json.JSONDecodeError) as exc:
                    raise CaptureError(
                        "invalid chrome devtools response") from exc


class _ChromeDevTools:
    """Small synchronous Chrome DevTools Protocol client."""

    def __init__(self, websocket_url: str):
        self.websocket = _WebSocket(websocket_url)
        self.next_id = 1

    def close(self) -> None:
        self.websocket.close()

    def call(self, method: str, params: Optional[dict] = None) -> dict:
        command_id = self.next_id
        self.next_id += 1
        self.websocket.send_json({
            "id": command_id,
            "method": method,
            "params": params or {},
        })
        while True:
            message = self.websocket.recv_json()
            if message.get("id") != command_id:
                continue
            if message.get("error"):
                raise CaptureError(
                    "chrome devtools %s failed: %s"
                    % (method, message["error"].get("message", "unknown error")))
            return message.get("result") or {}


def _chrome_page_websocket(port: int, timeout: float = 30.0) -> str:
    deadline = time.monotonic() + timeout
    endpoint = "http://127.0.0.1:%s/json/list" % port
    target_create_attempted = False
    while time.monotonic() < deadline:
        try:
            with urllib.request.urlopen(endpoint, timeout=1) as response:
                pages = json.loads(response.read().decode("utf-8"))
            for page in pages:
                # Fresh Chrome profiles may still expose component-extension
                # background pages before the actual browser tab. CDP's
                # screenshot command must target a normal page.
                if page.get("type") != "page":
                    continue
                websocket_url = page.get("webSocketDebuggerUrl")
                if websocket_url:
                    return websocket_url
            # Chrome normally creates the command-line about:blank tab, but a
            # busy worker can expose DevToolsActivePort before that page is
            # visible in /json/list. Ask the browser for one target instead of
            # treating this short readiness race as a missing browser channel.
            if not target_create_attempted:
                target_create_attempted = True
                create = urllib.request.Request(
                    "http://127.0.0.1:%s/json/new?about%%3Ablank" % port,
                    method="PUT",
                )
                try:
                    with urllib.request.urlopen(
                            create, timeout=1) as response:
                        page = json.loads(response.read().decode("utf-8"))
                    websocket_url = page.get("webSocketDebuggerUrl")
                    if page.get("type") == "page" and websocket_url:
                        return websocket_url
                except (OSError, ValueError, json.JSONDecodeError):
                    pass
        except (OSError, ValueError, json.JSONDecodeError):
            pass
        time.sleep(0.05)
    raise CaptureError("chrome devtools page endpoint unavailable")


def _chrome_active_port(profile: Path, process: subprocess.Popen,
                        timeout: float = 30.0) -> int:
    active_port = profile / "DevToolsActivePort"
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if process.poll() is not None:
            raise CaptureError(
                "chrome exited before devtools became ready")
        try:
            first_line = active_port.read_text().splitlines()[0]
            return int(first_line)
        except (OSError, IndexError, ValueError):
            time.sleep(0.05)
    raise CaptureError("chrome devtools startup timed out")


def _target_clip(client: _ChromeDevTools, target_text: str) -> dict:
    expression = """(() => {
      const needle = %s;
      const element = Array.from(document.querySelectorAll('li, td, tr'))
        .find(node => (node.innerText || '').includes(needle));
      if (!element) return null;
      element.scrollIntoView({block: 'center', inline: 'nearest'});
      const rect = element.getBoundingClientRect();
      return {
        x: Math.max(0, rect.left + window.scrollX - 8),
        y: Math.max(0, rect.top + window.scrollY - 8),
        width: Math.max(1, rect.width + 16),
        height: Math.max(1, rect.height + 16)
      };
    })()""" % json.dumps(target_text)
    result = client.call("Runtime.evaluate", {
        "expression": expression,
        "returnByValue": True,
        "awaitPromise": True,
    })
    value = ((result.get("result") or {}).get("value"))
    if not isinstance(value, dict):
        raise CaptureError(
            "target text not found in rendered page: %s" % target_text)
    return {
        "x": float(value["x"]),
        "y": float(value["y"]),
        "width": float(value["width"]),
        "height": float(value["height"]),
        "scale": 1,
    }


def _chrome_cdp_capture(binary: str, url: str, out: str, *,
                        wait_ms: int, full_page: bool, width: int,
                        height: int, target_text: str = "") -> None:
    """Capture a page through CDP, including content beyond the viewport."""
    with tempfile.TemporaryDirectory(prefix="jarvis-chrome-") as profile_dir:
        profile = Path(profile_dir)
        command = [
            binary,
            "--headless=new",
            "--disable-gpu",
            "--no-sandbox",
            "--disable-extensions",
            "--hide-scrollbars",
            "--remote-debugging-port=0",
            "--user-data-dir=%s" % profile,
            "about:blank",
        ]
        process = subprocess.Popen(
            command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        client = None
        try:
            startup_timeout = _chrome_startup_timeout()
            port = _chrome_active_port(
                profile, process, timeout=startup_timeout)
            client = _ChromeDevTools(_chrome_page_websocket(
                port, timeout=startup_timeout))
            client.call("Page.enable")
            client.call("Runtime.enable")
            client.call("Emulation.setDeviceMetricsOverride", {
                "width": width,
                "height": height,
                "deviceScaleFactor": 1,
                "mobile": False,
            })
            client.call("Page.navigate", {"url": url})
            deadline = time.monotonic() + 60
            while True:
                ready = client.call("Runtime.evaluate", {
                    "expression": "document.readyState",
                    "returnByValue": True,
                })
                state = ((ready.get("result") or {}).get("value"))
                if state in ("interactive", "complete"):
                    break
                if time.monotonic() >= deadline:
                    raise CaptureError(
                        "chrome page navigation timed out")
                time.sleep(0.1)
            if wait_ms > 0:
                time.sleep(wait_ms / 1000.0)

            params = {
                "format": "png",
                "fromSurface": True,
                "captureBeyondViewport": bool(full_page or target_text),
            }
            if target_text:
                params["clip"] = _target_clip(client, target_text)
            elif full_page:
                metrics = client.call("Page.getLayoutMetrics")
                content = (metrics.get("cssContentSize")
                           or metrics.get("contentSize") or {})
                content_width = max(1, math.ceil(float(
                    content.get("width") or width)))
                content_height = max(1, math.ceil(float(
                    content.get("height") or height)))
                max_height = max(2000, int(os.environ.get(
                    "JARVIS_SCREENSHOT_MAX_HEIGHT", "50000")))
                if content_height > max_height:
                    raise CaptureError(
                        "rendered page height %s exceeds safety limit %s; "
                        "use --text for element capture"
                        % (content_height, max_height))
                params["clip"] = {
                    "x": 0,
                    "y": 0,
                    "width": content_width,
                    "height": content_height,
                    "scale": 1,
                }
            result = client.call("Page.captureScreenshot", params)
            try:
                image = base64.b64decode(result["data"], validate=True)
            except (KeyError, ValueError) as exc:
                raise CaptureError(
                    "chrome returned an invalid screenshot") from exc
            if not image.startswith(b"\x89PNG\r\n\x1a\n"):
                raise CaptureError("chrome screenshot is not a PNG")
            Path(out).write_bytes(image)
        except subprocess.TimeoutExpired as exc:
            raise CaptureError("chrome headless capture timed out") from exc
        finally:
            if client is not None:
                client.close()
            process.terminate()
            try:
                process.communicate(timeout=5)
            except subprocess.TimeoutExpired:
                process.kill()
                process.communicate()


def _chrome_startup_timeout() -> float:
    raw = os.environ.get(
        "JARVIS_SCREENSHOT_CHROME_STARTUP_TIMEOUT",
        str(DEFAULT_CHROME_STARTUP_TIMEOUT),
    )
    try:
        return max(1.0, float(raw))
    except ValueError:
        return DEFAULT_CHROME_STARTUP_TIMEOUT


def _chrome_capture_attempts() -> int:
    raw = os.environ.get(
        "JARVIS_SCREENSHOT_CHROME_ATTEMPTS",
        str(DEFAULT_CHROME_CAPTURE_ATTEMPTS),
    )
    try:
        return min(5, max(1, int(raw)))
    except ValueError:
        return DEFAULT_CHROME_CAPTURE_ATTEMPTS


def _retryable_chrome_capture_error(exc: CaptureError) -> bool:
    """Return whether a fresh Chrome profile may recover the failure."""
    message = str(exc)
    return not (
        message.startswith("target text not found in rendered page:")
        or message.startswith("rendered page height ")
    )


@dataclass
class ChromeBinaryChannel(ScreenshotChannel):
    """Headless Chrome/Chromium captured through the DevTools protocol.

    CDP is used instead of Chrome's ``--screenshot`` CLI because the latter
    only captures the configured viewport. CDP supports true full-page and
    target-text element capture without a Playwright dependency.
    """

    name: str = "chrome_binary"
    binary_finder: Callable[[], Optional[str]] = _default_chrome_binary

    def available(self) -> bool:
        return bool(self.binary_finder())

    def unavailability_reason(self) -> str:
        return ("no chrome/chromium binary found; set JARVIS_CHROME_BIN, "
                "install chromium, or `pip install playwright && "
                "playwright install chromium`")

    def capture(self, url: str, out: str, *, wait_ms: int = 3000,
                full_page: bool = True, width: int = 1280,
                height: int = 2000, target_text: str = "") -> None:
        binary = self.binary_finder()
        if not binary:
            raise CaptureError(
                "chrome binary vanished between probe and capture")
        attempts = _chrome_capture_attempts()
        for attempt in range(1, attempts + 1):
            try:
                _chrome_cdp_capture(
                    binary, url, out, wait_ms=wait_ms, full_page=full_page,
                    width=width, height=height, target_text=target_text)
                return
            except CaptureError as exc:
                if (not _retryable_chrome_capture_error(exc)
                        or attempt == attempts):
                    if attempt > 1:
                        raise CaptureError(
                            "%s after %s attempts" % (exc, attempt)) from exc
                    raise
                time.sleep(min(0.25 * (2 ** (attempt - 1)), 1.0))
            except Exception as exc:  # noqa: BLE001
                raise CaptureError(
                    "chrome devtools capture failed: %s" % exc) from exc


# ---------------------------------------------------------------------------
# Probe + capture dispatcher
# ---------------------------------------------------------------------------


def default_channels() -> List[ScreenshotChannel]:
    """Channels in priority order: Playwright Python first, Chrome binary next."""
    return [PlaywrightPythonChannel(), ChromeBinaryChannel()]


def probe(channels: Sequence[ScreenshotChannel]) -> ScreenshotChannel:
    """Return the first available channel or raise MissingCapability."""
    for channel in channels:
        try:
            if channel.available():
                return channel
        except Exception:  # noqa: BLE001 — a flaky probe must not hide others
            continue
    reasons = "; ".join(
        "%s: %s" % (c.name, c.unavailability_reason()) for c in channels)
    raise MissingCapability(
        reason=reasons or "no screenshot channels configured",
        hint=("install a repo-controlled browser channel: "
              "`pip install playwright && playwright install chromium`, "
              "or install chrome/chromium (override via JARVIS_CHROME_BIN)"))


def capture(url: str, out: str, *,
            channels: Optional[Sequence[ScreenshotChannel]] = None,
            wait_ms: int = 3000, full_page: bool = True,
            width: int = 1280, height: int = 2000,
            target_text: str = "") -> str:
    """Capture `url` to `out` PNG; return the channel name used."""
    channel = probe(channels if channels is not None else default_channels())
    Path(out).parent.mkdir(parents=True, exist_ok=True)
    channel.capture(url, out, wait_ms=wait_ms, full_page=full_page,
                    width=width, height=height, target_text=target_text)
    return channel.name


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------


def _parse_capture_args(argv: Sequence[str]) -> tuple:
    if len(argv) < 2:
        raise SystemExit(
            "usage: capture <url> <out.png> "
            "[--wait N] [--full-page|--viewport] [--width W] [--height H] "
            "[--text TARGET]")
    url, out = argv[0], argv[1]
    wait_ms = 3000
    full_page = True
    width, height = 1280, 2000
    target_text = ""
    rest = list(argv[2:])
    while rest:
        flag = rest.pop(0)
        if flag == "--wait":
            wait_ms = max(0, int(rest.pop(0)))
        elif flag == "--full-page":
            full_page = True
        elif flag == "--viewport":
            full_page = False
        elif flag == "--width":
            width = int(rest.pop(0))
        elif flag == "--height":
            height = int(rest.pop(0))
        elif flag == "--text":
            target_text = rest.pop(0)
        else:
            raise SystemExit("unknown option: %s" % flag)
    return url, out, wait_ms, full_page, width, height, target_text


def main(argv: Sequence[str]) -> int:
    if not argv:
        raise SystemExit(
            "usage: bridge.jarvis_screenshot probe|capture ...")
    command = argv[0]
    if command == "probe":
        try:
            channel = probe(default_channels())
            print(channel.name)
            return EXIT_OK
        except MissingCapability as exc:
            print(MISSING_CAPABILITY_PREFIX + exc.reason)
            if exc.hint:
                print("fix: " + exc.hint, file=sys.stderr)
            return EXIT_MISSING_CAPABILITY
    if command == "capture":
        (url, out, wait_ms, full_page, width, height,
         target_text) = _parse_capture_args(
            list(argv[1:]))
        try:
            name = capture(url, out, wait_ms=wait_ms, full_page=full_page,
                           width=width, height=height,
                           target_text=target_text)
            print(name)
            return EXIT_OK
        except MissingCapability as exc:
            print(MISSING_CAPABILITY_PREFIX + exc.reason)
            if exc.hint:
                print("fix: " + exc.hint, file=sys.stderr)
            return EXIT_MISSING_CAPABILITY
        except CaptureError as exc:
            print("capture_error: " + str(exc), file=sys.stderr)
            return EXIT_CAPTURE_ERROR
    raise SystemExit("unknown command: %s" % command)


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

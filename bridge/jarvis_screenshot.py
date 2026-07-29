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
        [--wait N] [--full-page|--viewport] [--width W] [--height H]
"""

from __future__ import annotations

import os
import shutil
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Callable, List, Optional, Sequence

EXIT_OK = 0
EXIT_CAPTURE_ERROR = 1
EXIT_MISSING_CAPABILITY = 3

# Prefix that downstream code (skill n-a rows, autonomy escalate triggers) keys on.
MISSING_CAPABILITY_PREFIX = "missing_capability: "


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
                height: int = 2000) -> None:  # pragma: no cover - abstract
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
                height: int = 2000) -> None:
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
                    page.screenshot(path=str(out), full_page=full_page)
                finally:
                    browser.close()
        except CaptureError:
            raise
        except Exception as exc:  # noqa: BLE001
            raise CaptureError("playwright capture failed: %s" % exc) from exc


@dataclass
class ChromeBinaryChannel(ScreenshotChannel):
    """Headless Chrome/Chromium via `--headless=new --screenshot`.

    Full-viewport only (no element-level); used where Playwright Python is not
    installed but a Chrome/Chromium binary exists. macOS app and Linux PATH
    candidates are probed; `JARVIS_CHROME_BIN` overrides.
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
                height: int = 2000) -> None:
        binary = self.binary_finder()
        if not binary:
            raise CaptureError(
                "chrome binary vanished between probe and capture")
        budget = max(1000, int(wait_ms))
        cmd = [
            binary, "--headless=new", "--disable-gpu", "--no-sandbox",
            "--hide-scrollbars",
            "--screenshot=%s" % out,
            "--window-size=%s,%s" % (width, height),
            "--virtual-time-budget=%s" % budget,
            url,
        ]
        try:
            result = subprocess.run(
                cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                timeout=120)
        except subprocess.TimeoutExpired as exc:
            raise CaptureError("chrome headless capture timed out") from exc
        if result.returncode != 0 or not Path(out).is_file():
            tail = (result.stderr or b"").decode("utf-8", "replace")[-200:]
            raise CaptureError(
                "chrome headless capture failed rc=%s %s"
                % (result.returncode, tail.strip()))
        # Chrome writes the file before the process exits on macOS even though
        # the message goes to stderr; the Path check above is authoritative.


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
            width: int = 1280, height: int = 2000) -> str:
    """Capture `url` to `out` PNG; return the channel name used."""
    channel = probe(channels if channels is not None else default_channels())
    Path(out).parent.mkdir(parents=True, exist_ok=True)
    channel.capture(url, out, wait_ms=wait_ms, full_page=full_page,
                    width=width, height=height)
    return channel.name


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------


def _parse_capture_args(argv: Sequence[str]) -> tuple:
    if len(argv) < 2:
        raise SystemExit(
            "usage: capture <url> <out.png> "
            "[--wait N] [--full-page|--viewport] [--width W] [--height H]")
    url, out = argv[0], argv[1]
    wait_ms = 3000
    full_page = True
    width, height = 1280, 2000
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
        else:
            raise SystemExit("unknown option: %s" % flag)
    return url, out, wait_ms, full_page, width, height


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
        url, out, wait_ms, full_page, width, height = _parse_capture_args(
            list(argv[1:]))
        try:
            name = capture(url, out, wait_ms=wait_ms, full_page=full_page,
                           width=width, height=height)
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

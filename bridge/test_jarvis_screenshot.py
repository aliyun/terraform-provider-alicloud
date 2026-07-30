#!/usr/bin/env python3
"""Hermetic regression tests for the repo-controlled headless screenshot channel.

These cover the `no Playwright MCP` scenario the work item exists for: when no
browser channel is available the probe must surface a diagnosable
`missing_capability` (exit 3) instead of failing mid-run, and when a channel is
available capture must produce a valid PNG file. A stub channel writes a real
1x1 PNG so the `valid screenshot file` acceptance criterion is exercised without
depending on an installed browser.
"""

from __future__ import annotations

import contextlib
import io
import os
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))         # repo convention
sys.path.insert(0, str(HERE.parent))  # make THIS repo's bridge/ authoritative

import bridge.jarvis_screenshot as js  # noqa: E402

# Minimal 1x1 transparent PNG (8-byte signature + IHDR + IDAT + IEND).
_PNG = (
    b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01"
    b"\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\rIDATx\x9cc\xf8\xcf"
    b"\xc0\x00\x00\x00\x03\x00\x01\xa7\xf1\x07\xd1\x00\x00\x00\x00IEND"
    b"\xaeB`\x82"
)


class StubChannel(js.ScreenshotChannel):
    """Test double that records captures and writes a real PNG."""

    def __init__(self, name, *, available=True, capture_ok=True):
        self.name = name
        self._available = available
        self._capture_ok = capture_ok
        self.captured = []

    def available(self):
        return self._available

    def unavailability_reason(self):
        return "stub %s disabled" % self.name

    def capture(self, url, out, **kwargs):
        self.captured.append((url, out, kwargs))
        if not self._capture_ok:
            raise js.CaptureError("stub capture failed for %s" % self.name)
        Path(out).parent.mkdir(parents=True, exist_ok=True)
        Path(out).write_bytes(_PNG)


class RaisingChannel(js.ScreenshotChannel):
    name = "raising"

    def available(self):
        raise RuntimeError("probe blew up")

    def unavailability_reason(self):
        return "probe raised"


class ProbeTests(unittest.TestCase):
    def test_missing_capability_when_no_channels_available(self):
        chans = [StubChannel("a", available=False),
                 StubChannel("b", available=False)]
        with self.assertRaises(js.MissingCapability) as ctx:
            js.probe(chans)
        self.assertIn("a:", ctx.exception.reason)
        self.assertIn("b:", ctx.exception.reason)
        self.assertTrue(ctx.exception.hint)

    def test_returns_first_available_channel(self):
        chans = [StubChannel("a", available=False),
                 StubChannel("b", available=True)]
        self.assertEqual(js.probe(chans).name, "b")

    def test_priority_playwright_over_chrome(self):
        chans = [StubChannel("playwright_python", available=True),
                 StubChannel("chrome_binary", available=True)]
        self.assertEqual(js.probe(chans).name, "playwright_python")

    def test_skips_unavailable_falls_to_next(self):
        chans = [StubChannel("playwright_python", available=False),
                 StubChannel("chrome_binary", available=True)]
        self.assertEqual(js.probe(chans).name, "chrome_binary")

    def test_skips_channel_whose_probe_raises(self):
        # A flaky probe must not hide later available channels.
        chans = [RaisingChannel(), StubChannel("chrome_binary", available=True)]
        self.assertEqual(js.probe(chans).name, "chrome_binary")

    def test_empty_channels_missing_capability(self):
        with self.assertRaises(js.MissingCapability):
            js.probe([])


class CaptureTests(unittest.TestCase):
    def setUp(self):
        self._dir = tempfile.mkdtemp(prefix="jarvis-shot-")

    def tearDown(self):
        import shutil
        shutil.rmtree(self._dir, ignore_errors=True)

    def path(self, *parts):
        return os.path.join(self._dir, *parts)

    def test_dispatches_to_probed_channel(self):
        stub = StubChannel("playwright_python", available=True)
        out = self.path("shot.png")
        name = js.capture("https://example.com", out, channels=[stub])
        self.assertEqual(name, "playwright_python")
        self.assertEqual(stub.captured[0][0], "https://example.com")
        self.assertEqual(Path(out).read_bytes(), _PNG)

    def test_raises_missing_capability_when_none(self):
        with self.assertRaises(js.MissingCapability):
            js.capture("https://example.com", self.path("shot.png"),
                       channels=[StubChannel("a", available=False)])

    def test_raises_capture_error_on_channel_failure(self):
        stub = StubChannel("chrome_binary", available=True, capture_ok=False)
        with self.assertRaises(js.CaptureError):
            js.capture("https://example.com", self.path("shot.png"),
                       channels=[stub])

    def test_writes_valid_png_file_without_browser(self):
        # Acceptance: a headless run with no Playwright MCP still produces a
        # valid screenshot FILE via the repo-controlled channel.
        stub = StubChannel("chrome_binary", available=True)
        out = self.path("nested", "shot.png")
        js.capture("https://example.com", out, channels=[stub], wait_ms=0)
        data = Path(out).read_bytes()
        self.assertEqual(data[:8], b"\x89PNG\r\n\x1a\n")

    def test_creates_parent_dirs(self):
        stub = StubChannel("chrome_binary", available=True)
        out = self.path("deep", "dir", "shot.png")
        js.capture("https://example.com", out, channels=[stub])
        self.assertTrue(Path(out).is_file())

    def test_forwards_target_text_to_channel(self):
        stub = StubChannel("chrome_binary", available=True)
        out = self.path("target.png")
        js.capture(
            "https://example.com", out, channels=[stub],
            target_text="TargetField")
        self.assertEqual(
            stub.captured[0][2]["target_text"], "TargetField")


class ChromeBinaryChannelTests(unittest.TestCase):
    def test_page_websocket_skips_extension_background_pages(self):
        payload = [
            {
                "type": "background_page",
                "webSocketDebuggerUrl": "ws://localhost:1/background",
            },
            {
                "type": "page",
                "webSocketDebuggerUrl": "ws://localhost:1/page",
            },
        ]
        response = mock.MagicMock()
        response.__enter__.return_value.read.return_value = (
            __import__("json").dumps(payload).encode())
        with mock.patch.object(
                js.urllib.request, "urlopen", return_value=response):
            self.assertEqual(
                js._chrome_page_websocket(9222, timeout=0.1),
                "ws://localhost:1/page")

    def test_page_websocket_creates_page_when_list_is_empty(self):
        list_response = mock.MagicMock()
        list_response.__enter__.return_value.read.return_value = b"[]"
        new_response = mock.MagicMock()
        new_response.__enter__.return_value.read.return_value = (
            b'{"type":"page","webSocketDebuggerUrl":"ws://localhost:1/new"}')
        with mock.patch.object(
                js.urllib.request, "urlopen",
                side_effect=[list_response, new_response]) as urlopen:
            self.assertEqual(
                js._chrome_page_websocket(9222, timeout=0.1),
                "ws://localhost:1/new")
        self.assertEqual(urlopen.call_count, 2)
        self.assertEqual(urlopen.call_args_list[1].args[0].method, "PUT")

    def test_uses_cdp_for_true_full_page_capture(self):
        channel = js.ChromeBinaryChannel(binary_finder=lambda: "/bin/chrome")
        with mock.patch.object(js, "_chrome_cdp_capture") as capture:
            channel.capture(
                "https://example.com", "/tmp/shot.png",
                wait_ms=25, full_page=True, width=900, height=700,
                target_text="FieldName")
        capture.assert_called_once_with(
            "/bin/chrome", "https://example.com", "/tmp/shot.png",
            wait_ms=25, full_page=True, width=900, height=700,
            target_text="FieldName")

    def test_missing_binary_still_fails_closed(self):
        channel = js.ChromeBinaryChannel(binary_finder=lambda: None)
        with self.assertRaises(js.CaptureError):
            channel.capture("https://example.com", "/tmp/shot.png")

    def test_retries_transient_capture_with_fresh_profile(self):
        channel = js.ChromeBinaryChannel(binary_finder=lambda: "/bin/chrome")
        with mock.patch.dict(
                os.environ, {"JARVIS_SCREENSHOT_CHROME_ATTEMPTS": "3"}), \
                mock.patch.object(
                    js, "_chrome_cdp_capture",
                    side_effect=[
                        js.CaptureError(
                            "chrome devtools page endpoint unavailable"),
                        None,
                    ]) as capture, \
                mock.patch.object(js.time, "sleep"):
            channel.capture("https://example.com", "/tmp/shot.png")
        self.assertEqual(capture.call_count, 2)

    def test_does_not_retry_target_text_miss(self):
        channel = js.ChromeBinaryChannel(binary_finder=lambda: "/bin/chrome")
        with mock.patch.dict(
                os.environ, {"JARVIS_SCREENSHOT_CHROME_ATTEMPTS": "3"}), \
                mock.patch.object(
                    js, "_chrome_cdp_capture",
                    side_effect=js.CaptureError(
                        "target text not found in rendered page: Field")
                ) as capture:
            with self.assertRaises(js.CaptureError):
                channel.capture(
                    "https://example.com", "/tmp/shot.png",
                    target_text="Field")
        capture.assert_called_once()


class CliTests(unittest.TestCase):
    """probe + capture exit codes (the no-Playwright-MCP regression surface)."""

    def setUp(self):
        self._saved = js.default_channels
        self._dir = tempfile.mkdtemp(prefix="jarvis-cli-")

    def tearDown(self):
        js.default_channels = self._saved
        import shutil
        shutil.rmtree(self._dir, ignore_errors=True)

    def _run(self, argv, channels):
        js.default_channels = lambda: list(channels)
        out, err = io.StringIO(), io.StringIO()
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            rc = js.main(list(argv))
        return rc, out.getvalue(), err.getvalue()

    def test_probe_available_exit0(self):
        rc, out, _ = self._run(["probe"],
                               [StubChannel("chrome_binary", available=True)])
        self.assertEqual(rc, js.EXIT_OK)
        self.assertEqual(out.strip(), "chrome_binary")

    def test_probe_missing_capability_exit3(self):
        # The `no Playwright MCP` regression: probe must exit 3 with a
        # diagnosable missing_capability line, not silently succeed.
        rc, out, _ = self._run(
            ["probe"],
            [StubChannel("playwright_python", available=False),
             StubChannel("chrome_binary", available=False)])
        self.assertEqual(rc, js.EXIT_MISSING_CAPABILITY)
        self.assertTrue(out.startswith(js.MISSING_CAPABILITY_PREFIX))
        self.assertIn("playwright_python", out)
        self.assertIn("chrome_binary", out)

    def test_capture_missing_capability_exit3(self):
        out = os.path.join(self._dir, "shot.png")
        rc, out_text, _ = self._run(
            ["capture", "https://example.com", out],
            [StubChannel("playwright_python", available=False)])
        self.assertEqual(rc, js.EXIT_MISSING_CAPABILITY)
        self.assertFalse(os.path.exists(out))
        self.assertIn(js.MISSING_CAPABILITY_PREFIX, out_text)

    def test_capture_success_exit0(self):
        out = os.path.join(self._dir, "shot.png")
        rc, out_text, _ = self._run(
            ["capture", "https://example.com", out, "--wait", "0"],
            [StubChannel("chrome_binary", available=True)])
        self.assertEqual(rc, js.EXIT_OK)
        self.assertTrue(Path(out).is_file())
        self.assertEqual(out_text.strip(), "chrome_binary")

    def test_capture_capture_error_exit1(self):
        out = os.path.join(self._dir, "shot.png")
        rc, _, err = self._run(
            ["capture", "https://example.com", out],
            [StubChannel("chrome_binary", available=True, capture_ok=False)])
        self.assertEqual(rc, js.EXIT_CAPTURE_ERROR)
        self.assertIn("capture_error", err)

    def test_unknown_command_exits(self):
        with self.assertRaises(SystemExit):
            js.main(["frobnicate"])

    def test_capture_unknown_option_exits(self):
        js.default_channels = lambda: [StubChannel("chrome_binary", available=True)]
        with self.assertRaises(SystemExit):
            js.main(["capture", "https://example.com",
                     os.path.join(self._dir, "x.png"), "--bogus"])

    def test_capture_target_text_is_forwarded(self):
        out = os.path.join(self._dir, "target.png")
        stub = StubChannel("chrome_binary", available=True)
        rc, _, _ = self._run(
            ["capture", "https://example.com", out,
             "--text", "TargetField"], [stub])
        self.assertEqual(rc, js.EXIT_OK)
        self.assertEqual(
            stub.captured[0][2]["target_text"], "TargetField")


class DefaultChannelsTests(unittest.TestCase):
    def test_priority_order(self):
        self.assertEqual([c.name for c in js.default_channels()],
                         ["playwright_python", "chrome_binary"])

    def test_probe_never_raises(self):
        # On a host with nothing installed this returns MissingCapability,
        # not an exception — that is the degrade contract.
        try:
            js.probe(js.default_channels())
        except js.MissingCapability:
            pass


if __name__ == "__main__":
    unittest.main()

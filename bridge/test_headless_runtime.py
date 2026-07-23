from __future__ import annotations

from types import SimpleNamespace
import unittest

from bridge.headless_runtime import (
    HeadlessRequest,
    HeadlessRuntime,
    Lane,
    SessionPolicy,
)


class HeadlessRuntimeTests(unittest.TestCase):
    def test_cross_host_resume_falls_back_to_fresh_with_same_session_id(self):
        attempts = []
        runtime = HeadlessRuntime(
            lambda attempt: attempts.append(attempt) or SimpleNamespace(
                text="ok", is_error=False, subtype="success"),
            transcript_exists=lambda _sid: False,
        )
        result = runtime.execute(HeadlessRequest(
            "prompt", "stable-sid", SessionPolicy.RESUME, Lane.TERRAFORM,
            max_retries=0,
        ))
        self.assertFalse(result.is_error)
        self.assertEqual(result.attempts, 1)
        self.assertEqual(
            [(item.request.session_id, item.prompt, item.resume)
             for item in attempts],
            [("stable-sid", "prompt", False)],
        )

    def test_transient_retry_resumes_only_with_text_and_local_transcript(self):
        attempts = []
        results = iter((
            SimpleNamespace(text="partial", is_error=True, subtype="network"),
            SimpleNamespace(text="done", is_error=False, subtype="success"),
        ))
        runtime = HeadlessRuntime(
            lambda attempt: attempts.append(attempt) or next(results),
            transcript_exists=lambda sid: sid == "stable-sid",
            sleeper=lambda _seconds: None,
        )
        result = runtime.execute(HeadlessRequest(
            "original", "stable-sid", max_retries=2,
            retry_backoff_seconds=0,
        ))
        self.assertEqual(result.attempts, 2)
        self.assertEqual(
            [(item.prompt, item.resume) for item in attempts],
            [
                ("original", False),
                (HeadlessRuntime.RESUME_PROMPT, True),
            ],
        )

    def test_no_text_retries_fresh_even_when_transcript_exists(self):
        attempts = []
        results = iter((
            SimpleNamespace(text="", is_error=True, subtype="network"),
            SimpleNamespace(text="done", is_error=False, subtype="success"),
        ))
        runtime = HeadlessRuntime(
            lambda attempt: attempts.append(attempt) or next(results),
            transcript_exists=lambda _sid: True,
            sleeper=lambda _seconds: None,
        )
        runtime.execute(HeadlessRequest(
            "original", "sid", max_retries=1, retry_backoff_seconds=0))
        self.assertEqual(
            [(item.prompt, item.resume) for item in attempts],
            [("original", False), ("original", False)],
        )

    def test_timeout_and_max_turns_fast_fail_without_inner_retry(self):
        for subtype in ("timeout", "error_max_turns"):
            with self.subTest(subtype=subtype):
                calls = []
                runtime = HeadlessRuntime(
                    lambda attempt: calls.append(attempt) or SimpleNamespace(
                        text="failed", is_error=True, subtype=subtype),
                    transcript_exists=lambda _sid: True,
                    sleeper=lambda _seconds: None,
                )
                result = runtime.execute(HeadlessRequest(
                    "prompt", "sid", max_retries=3,
                    retry_backoff_seconds=0,
                ))
                self.assertEqual(result.attempts, 1)
                self.assertEqual(len(calls), 1)

    def test_unsupported_model_is_rejected_before_execution(self):
        with self.assertRaisesRegex(ValueError, "model override"):
            HeadlessRequest("prompt", "sid", model="some-model")


if __name__ == "__main__":
    unittest.main()

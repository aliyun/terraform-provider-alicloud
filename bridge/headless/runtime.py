"""Business-neutral retry/resume orchestration for headless Jarvis sessions."""

from __future__ import annotations

import logging
import time
from typing import Any, Callable, Protocol

from .model import HeadlessAttempt, HeadlessRequest, HeadlessResult, SessionPolicy


class AttemptRunner(Protocol):
    def __call__(self, attempt: HeadlessAttempt) -> Any:
        """Return an object exposing text/is_error/subtype."""


class HeadlessRuntime:
    """Execute a stable session id through a bounded transient retry series."""

    TERMINAL_SUBTYPES = frozenset({"timeout", "error_max_turns"})
    RESUME_PROMPT = "上一次执行因瞬时错误中断，请从中断处继续完成本工单的 SOP。"

    def __init__(
        self,
        attempt_runner: AttemptRunner,
        *,
        transcript_exists: Callable[[str], bool],
        sleeper: Callable[[float], None] = time.sleep,
        logger: logging.Logger | None = None,
    ) -> None:
        if not callable(attempt_runner):
            raise TypeError("attempt_runner must be callable")
        if not callable(transcript_exists):
            raise TypeError("transcript_exists must be callable")
        if not callable(sleeper):
            raise TypeError("sleeper must be callable")
        self._attempt_runner = attempt_runner
        self._transcript_exists = transcript_exists
        self._sleep = sleeper
        self._log = logger or logging.getLogger(__name__)

    def execute(self, request: HeadlessRequest) -> HeadlessResult:
        if not isinstance(request, HeadlessRequest):
            raise TypeError("request must be a HeadlessRequest")
        prompt = request.prompt
        resume = request.session_policy is SessionPolicy.RESUME
        if resume and not self._transcript_exists(request.session_id):
            self._log.warning(
                "headless session %s requested resume without a local transcript; "
                "starting fresh",
                request.session_id,
            )
            resume = False
        retries = 0
        while True:
            raw = self._attempt_runner(
                HeadlessAttempt(request, prompt, resume, retries + 1))
            result = _normalize_result(raw, retries + 1)
            if not result.is_error:
                return result
            if result.subtype in self.TERMINAL_SUBTYPES:
                return result
            if retries >= request.max_retries:
                return result
            retries += 1
            # Output alone is not proof that Claude created a resumable
            # transcript: ``no_result`` often carries only one stderr line.
            if result.text and self._transcript_exists(request.session_id):
                resume = True
                prompt = self.RESUME_PROMPT
            else:
                resume = False
                prompt = request.prompt
            self._log.warning(
                "headless session %s transient error subtype=%s; retry %s/%s",
                request.session_id,
                result.subtype,
                retries,
                request.max_retries,
            )
            delay = min(request.retry_backoff_seconds * retries, 300)
            if delay:
                self._sleep(delay)


def _normalize_result(value: Any, attempts: int) -> HeadlessResult:
    try:
        text = value.text
        is_error = value.is_error
        subtype = value.subtype
    except AttributeError as exc:
        raise TypeError(
            "attempt_runner must return text/is_error/subtype attributes") from exc
    if not isinstance(is_error, bool):
        raise TypeError("attempt result is_error must be a bool")
    if not isinstance(subtype, str) or not subtype.strip():
        raise TypeError("attempt result subtype must be a nonblank string")
    if subtype == "no_result":
        is_error = True
    return HeadlessResult(
        text if isinstance(text, str) else "",
        is_error,
        subtype,
        attempts,
    )

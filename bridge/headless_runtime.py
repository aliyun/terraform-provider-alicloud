"""Validated contracts and retry orchestration for one-shot Jarvis sessions."""

from __future__ import annotations

from dataclasses import dataclass
from enum import Enum
import logging
import time
from typing import Any, Callable, Optional, Protocol

from bridge.jarvis_execution_runtime import (
    run_claude_buffered,
    session_file_exists,
)


class SessionPolicy(str, Enum):
    NEW = "NEW"
    RESUME = "RESUME"


class Lane(str, Enum):
    DEFAULT = "default"
    TERRAFORM = "terraform"


@dataclass(frozen=True)
class HeadlessRequest:
    prompt: str
    session_id: str
    session_policy: SessionPolicy = SessionPolicy.NEW
    lane: Lane = Lane.DEFAULT
    model: Optional[str] = None
    timeout_seconds: float = 43200
    max_retries: int = 2
    retry_backoff_seconds: float = 30
    on_spawn: Optional[Callable[[Any], None]] = None
    guarded: bool = False

    def __post_init__(self) -> None:
        if not isinstance(self.prompt, str) or not self.prompt.strip():
            raise ValueError("prompt must be a nonblank string")
        if not isinstance(self.session_id, str) or not self.session_id.strip():
            raise ValueError("session_id must be a nonblank string")
        if not isinstance(self.session_policy, SessionPolicy):
            raise ValueError("session_policy must be a SessionPolicy")
        if not isinstance(self.lane, Lane):
            raise ValueError("lane must be a Lane")
        if self.model is not None:
            if not isinstance(self.model, str) or not self.model.strip():
                raise ValueError("model must be a nonblank string when provided")
            raise ValueError("model override is not supported by the Jarvis adapter")
        if (isinstance(self.timeout_seconds, bool)
                or not isinstance(self.timeout_seconds, (int, float))
                or self.timeout_seconds <= 0):
            raise ValueError("timeout_seconds must be positive")
        if (isinstance(self.max_retries, bool)
                or not isinstance(self.max_retries, int)
                or self.max_retries < 0):
            raise ValueError("max_retries must be a non-negative integer")
        if (isinstance(self.retry_backoff_seconds, bool)
                or not isinstance(self.retry_backoff_seconds, (int, float))
                or self.retry_backoff_seconds < 0):
            raise ValueError("retry_backoff_seconds must be non-negative")
        if self.on_spawn is not None and not callable(self.on_spawn):
            raise TypeError("on_spawn must be callable")
        if not isinstance(self.guarded, bool):
            raise ValueError("guarded must be a bool")


@dataclass(frozen=True)
class HeadlessAttempt:
    request: HeadlessRequest
    prompt: str
    resume: bool
    number: int


@dataclass(frozen=True)
class HeadlessResult:
    text: str
    is_error: bool
    subtype: str
    attempts: int

    def __post_init__(self) -> None:
        if not isinstance(self.text, str):
            raise TypeError("text must be a string")
        if not isinstance(self.is_error, bool):
            raise TypeError("is_error must be a bool")
        if not isinstance(self.subtype, str) or not self.subtype.strip():
            raise ValueError("subtype must be a nonblank string")
        if (isinstance(self.attempts, bool)
                or not isinstance(self.attempts, int)
                or self.attempts <= 0):
            raise ValueError("attempts must be a positive integer")


class AttemptRunner(Protocol):
    def __call__(self, attempt: HeadlessAttempt) -> Any:
        """Return an object exposing text/is_error/subtype."""


class HeadlessRuntime:
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
                "starting fresh", request.session_id)
            resume = False
        retries = 0
        while True:
            result = _normalize_result(
                self._attempt_runner(
                    HeadlessAttempt(request, prompt, resume, retries + 1)),
                retries + 1,
            )
            if (not result.is_error
                    or result.subtype in self.TERMINAL_SUBTYPES
                    or retries >= request.max_retries):
                return result
            retries += 1
            if result.text and self._transcript_exists(request.session_id):
                resume = True
                prompt = self.RESUME_PROMPT
            else:
                resume = False
                prompt = request.prompt
            self._log.warning(
                "headless session %s transient error subtype=%s; retry %s/%s",
                request.session_id, result.subtype, retries, request.max_retries)
            delay = min(request.retry_backoff_seconds * retries, 300)
            if delay:
                self._sleep(delay)


def _normalize_result(value: Any, attempts: int) -> HeadlessResult:
    try:
        text, is_error, subtype = value.text, value.is_error, value.subtype
    except AttributeError as exc:
        raise TypeError(
            "attempt_runner must return text/is_error/subtype attributes") from exc
    if not isinstance(is_error, bool):
        raise TypeError("attempt result is_error must be a bool")
    if not isinstance(subtype, str) or not subtype.strip():
        raise TypeError("attempt result subtype must be a nonblank string")
    return HeadlessResult(
        text if isinstance(text, str) else "",
        True if subtype == "no_result" else is_error,
        subtype,
        attempts,
    )


def jarvis_transcript_exists(session_id: str) -> bool:
    return session_file_exists(session_id)


def run_jarvis_attempt(attempt: HeadlessAttempt):
    if attempt.request.model is not None:
        raise ValueError("Jarvis headless model override is not supported")
    return run_claude_buffered(
        attempt.prompt,
        attempt.request.session_id,
        attempt.resume,
        timeout=attempt.request.timeout_seconds,
        on_spawn=attempt.request.on_spawn,
        terraform=attempt.request.lane is Lane.TERRAFORM,
        guarded=attempt.request.guarded,
    )


__all__ = [
    "HeadlessAttempt", "HeadlessRequest", "HeadlessResult", "HeadlessRuntime",
    "Lane", "SessionPolicy", "jarvis_transcript_exists", "run_jarvis_attempt",
]

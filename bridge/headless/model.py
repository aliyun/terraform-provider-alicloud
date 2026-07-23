"""Immutable, business-neutral contracts for a headless Jarvis invocation."""

from __future__ import annotations

from dataclasses import dataclass
from enum import Enum
from typing import Any, Callable, Optional


class SessionPolicy(str, Enum):
    """How the first attempt should treat an existing session id."""

    NEW = "NEW"
    RESUME = "RESUME"


class Lane(str, Enum):
    """Credential/model routing lane selected for the whole retry series."""

    DEFAULT = "default"
    TERRAFORM = "terraform"


@dataclass(frozen=True)
class HeadlessRequest:
    """Everything needed to execute one bounded, retryable headless session."""

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
            # The current Jarvis command builder selects a model through its lane's
            # settings file.  Silently accepting a model override would lie about
            # the execution contract, so keep it fail-closed until the adapter has
            # an explicit model argument.
            raise ValueError("model override is not supported by the Jarvis adapter")
        if (
            isinstance(self.timeout_seconds, bool)
            or not isinstance(self.timeout_seconds, (int, float))
            or self.timeout_seconds <= 0
        ):
            raise ValueError("timeout_seconds must be positive")
        if (
            isinstance(self.max_retries, bool)
            or not isinstance(self.max_retries, int)
            or self.max_retries < 0
        ):
            raise ValueError("max_retries must be a non-negative integer")
        if (
            isinstance(self.retry_backoff_seconds, bool)
            or not isinstance(self.retry_backoff_seconds, (int, float))
            or self.retry_backoff_seconds < 0
        ):
            raise ValueError("retry_backoff_seconds must be non-negative")
        if self.on_spawn is not None and not callable(self.on_spawn):
            raise TypeError("on_spawn must be callable")
        if not isinstance(self.guarded, bool):
            raise ValueError("guarded must be a bool")


@dataclass(frozen=True)
class HeadlessAttempt:
    """One concrete attempt derived from a stable ``HeadlessRequest``."""

    request: HeadlessRequest
    prompt: str
    resume: bool
    number: int


@dataclass(frozen=True)
class HeadlessResult:
    """Normalized terminal result returned by ``HeadlessRuntime``."""

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
        if (
            isinstance(self.attempts, bool)
            or not isinstance(self.attempts, int)
            or self.attempts <= 0
        ):
            raise ValueError("attempts must be a positive integer")

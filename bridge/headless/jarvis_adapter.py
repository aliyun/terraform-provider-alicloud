"""Thin adapter from the generic HeadlessRuntime to Jarvis CLI execution."""

from __future__ import annotations

from .model import HeadlessAttempt, Lane


def _bot_module():
    try:
        import jarvis_dingtalk_bot as bot
    except ModuleNotFoundError:  # package import in tests/tools
        from bridge import jarvis_dingtalk_bot as bot
    return bot


def jarvis_transcript_exists(session_id: str) -> bool:
    """Lazily query the canonical transcript location used by Jarvis resume."""

    return bool(_bot_module()._session_file_exists(session_id))


def run_jarvis_attempt(attempt: HeadlessAttempt):
    """Run one attempt without constructing ``JarvisHandler`` or any executor."""

    if attempt.request.model is not None:
        raise ValueError("Jarvis headless model override is not supported")
    # Keep the legacy composition import lazy: a Scheduler with no due headless
    # slot must not load the DingTalk bridge, and this import never instantiates
    # JarvisHandler or starts PersistenceExecutor.
    bot = _bot_module()
    return bot.run_claude_buffered(
        attempt.prompt,
        attempt.request.session_id,
        attempt.resume,
        timeout=attempt.request.timeout_seconds,
        on_spawn=attempt.request.on_spawn,
        terraform=attempt.request.lane is Lane.TERRAFORM,
        guarded=attempt.request.guarded,
    )

"""Portable Task input shared by interactive and persistent Jarvis runtimes."""

from __future__ import annotations

import hashlib
import json
from typing import Any, Dict, Mapping, Optional

try:
    from bridge.task_policy import HEADLESS_POLICY_REVISION
except ModuleNotFoundError:  # bootstrap scripts append bridge/ as a top-level path.
    from task_policy import HEADLESS_POLICY_REVISION


PORTABLE_INPUT_CONTRACT = "PORTABLE_V1"
RUNTIME_INTERACTIVE = "INTERACTIVE"
RUNTIME_PERSISTENT = "PERSISTENT"
SUPPORTED_RUNTIMES = frozenset((RUNTIME_INTERACTIVE, RUNTIME_PERSISTENT))


class TaskInputContractError(ValueError):
    """The immutable Task input cannot be executed by the selected runtime."""

    retryable = False

    def __init__(self, code: str, message: str) -> None:
        self.code = str(code or "INVALID_PORTABLE_INPUT")
        super().__init__("%s: %s" % (self.code, str(message)))


def canonical_input_digest(payload: Mapping[str, Any]) -> str:
    if not isinstance(payload, Mapping):
        raise TaskInputContractError(
            "INVALID_INPUT_SHAPE", "Task input must be a JSON object")
    raw = json.dumps(
        dict(payload), ensure_ascii=False, sort_keys=True,
        separators=(",", ":"), default=str,
    ).encode("utf-8")
    return hashlib.sha256(raw).hexdigest()


def normalize_runtime(value: Any) -> str:
    runtime = str(value or "").strip().upper()
    if runtime not in SUPPORTED_RUNTIMES:
        raise TaskInputContractError(
            "UNSUPPORTED_TARGET_RUNTIME",
            "target runtime must be INTERACTIVE or PERSISTENT")
    return runtime


def fresh_replay_prompt(item_id: Any, project: Any, title: Any = "") -> str:
    """Return the deterministic minimum prompt needed for safe cross-runtime replay."""
    item = str(item_id or "").strip()
    project_id = str(project or "").strip()
    if not item or not project_id:
        raise TaskInputContractError(
            "INPUT_NOT_REHYDRATABLE",
            "canonical Aone itemId and project are required")
    title_text = str(title or "").strip()
    subject = "Aone 工单 #%s（项目 %s）" % (item, project_id)
    if title_text:
        subject += "，标题：%s" % title_text[:200]
    return (
        "继续处理%s。开始执行前，必须先读取当前 Aone 详情和最新评论，并读取控制面 "
        "Task、Session、Operation/receipt 状态；以最新人类输入和已确认回执为准。"
        "任何外部写操作前都要检查对应副作用是否已经完成，禁止重复评论、重复状态更新、"
        "重复认领或重复调用外部接口。随后按 aone-triage、当前仓库约束和正常 bookend "
        "继续处理；若状态冲突、回执不确定或输入仍不足，必须 fail-closed 并请求人工决策。"
    ) % subject


def portable_task_payload(
        *, item_id: Any, project: Any, kind: Any, prompt: Any,
        origin_runtime: Any, trigger: Any,
        policy_revision: Any = HEADLESS_POLICY_REVISION,
        rehydrated: bool = False, **extra: Any) -> Dict[str, Any]:
    """Build the one immutable input shape accepted by both runtimes."""
    runtime = normalize_runtime(origin_runtime)
    payload: Dict[str, Any] = {
        "itemId": str(item_id or "").strip(),
        "project": str(project or "").strip(),
        "kind": str(kind or "").strip(),
        "prompt": str(prompt or ""),
        "trigger": str(trigger or "").strip(),
        "policyRevision": str(policy_revision or "").strip(),
        "inputContract": PORTABLE_INPUT_CONTRACT,
        "inputContext": {
            "originRuntime": runtime,
            "rehydrated": bool(rehydrated),
        },
    }
    payload.update(extra)
    validate_portable_task_input(payload)
    return payload


def validate_portable_task_input(
        payload: Any, *, target_runtime: Optional[Any] = None) -> Dict[str, Any]:
    """Validate and return a copy of a PORTABLE_V1 input."""
    if not isinstance(payload, Mapping):
        raise TaskInputContractError(
            "INVALID_INPUT_SHAPE", "Task input must be a JSON object")
    result = dict(payload)
    if target_runtime is not None:
        normalize_runtime(target_runtime)
    missing = [
        key for key in ("itemId", "project", "prompt")
        if not str(result.get(key) or "").strip()
    ]
    if missing:
        raise TaskInputContractError(
            "INPUT_NOT_PORTABLE",
            "Task input is missing %s" % ", ".join(missing))
    contract = str(result.get("inputContract") or "").strip()
    if contract != PORTABLE_INPUT_CONTRACT:
        raise TaskInputContractError(
            "INPUT_CONTRACT_MISMATCH",
            "expected %s, got %s"
            % (PORTABLE_INPUT_CONTRACT, contract or "<missing>"))
    return result


def effective_input_from_timeline(
        timeline: Mapping[str, Any],
        selected_session_id: Optional[Any] = None) -> Dict[str, Any]:
    """Read the same frozen/effective input for which the server exposes a digest."""
    if not isinstance(timeline, Mapping):
        raise TaskInputContractError(
            "INVALID_TIMELINE", "task timeline must be an object")
    direct = timeline.get("effectiveInputPayload")
    if isinstance(direct, Mapping):
        return dict(direct)
    task = timeline.get("task")
    task = task if isinstance(task, Mapping) else {}
    session_id = (
        selected_session_id
        if selected_session_id is not None else task.get("currentSessionId"))
    sessions = timeline.get("sessions") or []
    selected = next((
        value for value in sessions
        if isinstance(value, Mapping)
        and session_id is not None
        and str(value.get("id")) == str(session_id)
    ), None)
    if isinstance(selected, Mapping):
        value = selected.get("inputPayload")
        if isinstance(value, str):
            try:
                value = json.loads(value)
            except (TypeError, ValueError) as exc:
                raise TaskInputContractError(
                    "INVALID_SOURCE_INPUT",
                    "frozen Session input is not valid JSON") from exc
        if isinstance(value, Mapping):
            return dict(value)
    value = task.get("payload")
    if isinstance(value, str):
        try:
            value = json.loads(value)
        except (TypeError, ValueError) as exc:
            raise TaskInputContractError(
                "INVALID_SOURCE_INPUT",
                "Task input is not valid JSON") from exc
    if isinstance(value, Mapping):
        return dict(value)
    raise TaskInputContractError(
        "SOURCE_INPUT_UNAVAILABLE",
        "timeline does not expose effectiveInputPayload, Session.inputPayload, or Task.payload")


def portable_replacement_for_redispatch(
        timeline: Mapping[str, Any], *, target_runtime: Any,
        selected_session_id: Optional[Any] = None) -> Optional[Dict[str, Any]]:
    """Return an audited replacement for legacy Aone input, or ``None`` if portable."""
    runtime = normalize_runtime(target_runtime)
    source = effective_input_from_timeline(timeline, selected_session_id)
    try:
        validate_portable_task_input(source, target_runtime=runtime)
        return None
    except TaskInputContractError:
        pass

    task = timeline.get("task")
    task = task if isinstance(task, Mapping) else {}
    source_ref = task.get("sourceRef")
    source_ref = source_ref if isinstance(source_ref, Mapping) else {}
    task_key = str(task.get("taskKey") or "")
    task_key_parts = task_key.split(":")
    canonical_project = (
        task_key_parts[1]
        if len(task_key_parts) == 3 and task_key_parts[0] == "aone" else "")
    canonical_item = (
        task_key_parts[2]
        if len(task_key_parts) == 3 and task_key_parts[0] == "aone" else "")
    item_id = str(
        source.get("itemId") or task.get("aoneId")
        or source_ref.get("aoneId") or canonical_item or "").strip()
    project = str(
        source.get("project") or source_ref.get("projectId")
        or canonical_project or "").strip()
    source_type = str(task.get("sourceType") or "").strip().upper()
    if source_type != "AONE" and not task_key.startswith("aone:"):
        raise TaskInputContractError(
            "INPUT_NOT_REHYDRATABLE",
            "only canonical AONE Tasks can be rebuilt without a business prompt")
    title = source_ref.get("title") or task.get("title") or ""
    extra = {
        key: value for key, value in source.items()
        if key not in {
            "itemId", "project", "kind", "prompt", "trigger",
            "policyRevision", "inputContract", "inputContext",
        }
    }
    replacement = portable_task_payload(
        item_id=item_id,
        project=project,
        kind=source.get("kind") or task.get("taskType") or "ticket",
        prompt=fresh_replay_prompt(item_id, project, title),
        origin_runtime=(
            source.get("inputContext", {}).get("originRuntime")
            if isinstance(source.get("inputContext"), Mapping) else
            RUNTIME_INTERACTIVE),
        trigger="REDISPATCH_REHYDRATE",
        rehydrated=True,
        **extra,
    )
    source_digest = str(timeline.get("effectiveInputDigest") or "").strip()
    calculated_source_digest = canonical_input_digest(source)
    if source_digest and source_digest != calculated_source_digest:
        raise TaskInputContractError(
            "SOURCE_INPUT_DIGEST_MISMATCH",
            "timeline payload does not match effectiveInputDigest")
    source_digest = source_digest or calculated_source_digest
    replacement_digest = canonical_input_digest(replacement)
    desired_revision = "portable:%s" % replacement_digest
    if desired_revision in {
            str(task.get("desiredRevision") or ""),
            str(task.get("processingRevision") or "")}:
        raise TaskInputContractError(
            "REPLACEMENT_REVISION_NOT_FRESH",
            "portable replacement is already the effective Task revision")
    return {
        "payload": replacement,
        "expectedSourceInputDigest": source_digest,
        "replacementDigest": replacement_digest,
        "desiredRevision": desired_revision,
    }

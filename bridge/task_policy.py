"""Versioned policy contract for durable Jarvis Task generations."""

from __future__ import annotations

import hashlib
import json
import re
from typing import Any, Mapping


HEADLESS_POLICY_REVISION = "terraform-rd-single-writer-v6"
STALE_TASK_POLICY_ERROR = "stale_task_policy_revision"


def _policy_label(policy_revision: str) -> str:
    match = re.search(r"(v\d+)$", str(policy_revision or ""))
    return match.group(1) if match else str(policy_revision or "unknown")[-24:]


def policy_desired_revision(
    source_revision: object,
    payload: Mapping[str, Any],
    *,
    policy_revision: str = HEADLESS_POLICY_REVISION,
) -> str:
    """Salt one business revision with policy and immutable input content.

    The source prefix remains human-readable while the canonical payload hash
    makes prompt/input changes produce a new desired generation even when the
    upstream Aone timestamp or comment cursor did not change.
    """
    canonical = json.dumps(
        dict(payload), ensure_ascii=False, sort_keys=True,
        separators=(",", ":"), default=str,
    ).encode("utf-8")
    digest = hashlib.sha256(canonical).hexdigest()[:16]
    source = re.sub(r"[^A-Za-z0-9._:@+-]+", "-", str(source_revision or "source"))
    suffix = "|policy:%s|input:%s" % (_policy_label(policy_revision), digest)
    return source[:max(1, 127 - len(suffix))] + suffix


class StaleTaskPolicyError(RuntimeError):
    """A frozen durable Task belongs to an older execution policy."""

    def __init__(self, frozen_revision: object, current_revision: object) -> None:
        frozen = str(frozen_revision or "<missing>").strip() or "<missing>"
        current = str(current_revision or "<missing>").strip() or "<missing>"
        super().__init__(
            "%s: frozen=%s current=%s"
            % (STALE_TASK_POLICY_ERROR, frozen, current))

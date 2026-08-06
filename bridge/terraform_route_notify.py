#!/usr/bin/env python3
"""Typed durable DingTalk notifications for Terraform route D.

The RD finalizer calls this only after the source work item's assignee and
per-type progress status have been synchronized.  Delivery and retry semantics
stay inside the shared DingTalk event ledger; callers never invoke
``notify-dingtalk.sh`` directly.
"""

from __future__ import annotations

import argparse
import json
import subprocess

from bridge.helpers import aone
from bridge.helpers import dingtalk


SOURCE_PROJECT = "1086837"
ASSIGN_SCRIPT = "bootstrap/aone-assign.sh"
ROUTE_SPECS = {
    "handwritten-urgent": {
        "staff_id": "521957",
        "label": "手写 resource 紧急",
    },
    "handwritten-normal": {
        "staff_id": "484483",
        "label": "手写 resource 非紧急",
    },
    "generated": {
        "staff_id": "429768",
        "label": "生成/发布（含 CloudSpec pre 收敛后）",
    },
}


def _recipient_owns_or_may_own(ticket: str, staff_id: str) -> tuple[bool, str]:
    """Would assigning ``ticket`` to ``staff_id`` be allowed by the owner policy?

    The DM recipient comes from the subtype table, never from the work item, so
    a DM sent while somebody else holds the ticket tells the wrong person to
    start work.  ``aone-assign.sh --check`` is the single source of that policy
    (API-team roster ∪ 云产品专属维护名单, with the same fail-closed reads), and
    its idempotent-no-op case is exactly "the recipient already owns it".

    Returns ``(allowed, detail)``.  Anything other than a clean rc=0 is treated
    as not-allowed: when ownership cannot be verified the assignee sync was
    refused too, so staying quiet keeps both halves of the route step honest.
    """
    script = aone.REPO_ROOT / ASSIGN_SCRIPT
    try:
        proc = subprocess.run(
            ["bash", str(script), "--check", str(ticket), str(staff_id)],
            capture_output=True, text=True, timeout=120)
    except Exception as exc:  # noqa: BLE001
        return False, "owner check failed to run: %s" % exc
    detail = (proc.stderr or proc.stdout or "").strip().splitlines()
    first = detail[0].strip() if detail else ""
    if proc.returncode == 0:
        return True, first
    return False, first or "owner check rc=%s" % proc.returncode


def _skipped_result(ticket: str, project: str, subtype: str, staff_id: str,
                    event_key: str, detail: str) -> dict:
    """A policy skip: nothing enqueued, nothing to compensate, nothing claimed."""
    return {
        "ticket": ticket,
        "project": project,
        "route": "d",
        "subtype": subtype,
        "staff_id": staff_id,
        "event_key": event_key,
        "ledger_id": "",
        "receipt": "",
        "state": "skipped_owner_protected",
        "skipped": True,
        "skip_reason": detail,
        "durable": False,
        "notification_complete": False,
    }


def _ledger_result(ticket: str, project: str, subtype: str, staff_id: str,
                   event_key: str) -> dict:
    ledger_id = aone._aone_event_ledger_id(ticket, event_key)
    ledger = dingtalk._dingtalk_event_load()
    bucket = ""
    record = {}
    for name in ("pending", "posted", "suppressed"):
        candidate = ledger.get(name, {}).get(ledger_id)
        if isinstance(candidate, dict):
            bucket, record = name, candidate
            break
    durable = bool(bucket)
    state = str(record.get("state") or bucket or "unpersisted")
    return {
        "ticket": ticket,
        "project": project,
        "route": "d",
        "subtype": subtype,
        "staff_id": staff_id,
        "event_key": event_key,
        "ledger_id": ledger_id,
        "receipt": str(record.get("receipt") or ""),
        "state": state,
        "skipped": False,
        "durable": durable,
        "notification_complete": bucket in {"posted", "suppressed"},
    }


def enqueue_d_route_notification(ticket: str, subtype: str,
                                 project: str = SOURCE_PROJECT) -> dict:
    """Persist and opportunistically deliver one idempotent route-D DM.

    A durable ``pending``/``failed``/``post_uncertain`` record is accepted and
    will be compensated by the shared ledger flush path.  An unpersisted result
    is explicitly returned as incomplete so the finalizer cannot claim that the
    notification succeeded.

    Nothing is enqueued when the owner policy would refuse assigning the ticket
    to this subtype's recipient — a branch-A 专属维护人 or any other human
    already holds it, so the DM would page the wrong person.  That comes back as
    ``state="skipped_owner_protected"`` with ``notification_complete`` false: a
    deliberate no-op the finalizer must report as "未发送", not retry.
    """
    ticket = str(ticket or "").strip()
    project = str(project or "").strip()
    subtype = str(subtype or "").strip()
    if not ticket.isdigit():
        raise ValueError("ticket must be numeric")
    if not project:
        raise ValueError("project is required")
    spec = ROUTE_SPECS.get(subtype)
    if spec is None:
        raise ValueError(
            "subtype must be one of: %s" % ", ".join(sorted(ROUTE_SPECS)))

    staff_id = spec["staff_id"]
    event_key = "terraform-route:d:%s:owner:%s" % (subtype, staff_id)
    allowed, detail = _recipient_owns_or_may_own(ticket, staff_id)
    if not allowed:
        return _skipped_result(
            ticket, project, subtype, staff_id, event_key, detail)
    dingtalk._dingtalk_event_enqueue(
        ticket,
        project,
        event_key,
        staff_id,
        "Terraform 路由已更新",
        "Terraform 工单 #%s 已按 %s 路由给你；Jarvis/TerraformRD 将同步主动开发，"
        "请关注后续验证与人工发布硬门。" % (ticket, spec["label"]),
    )
    return _ledger_result(ticket, project, subtype, staff_id, event_key)


def main(argv=None) -> int:
    parser = argparse.ArgumentParser(
        description="enqueue a durable Terraform route-D DingTalk notification")
    parser.add_argument("--ticket", required=True)
    parser.add_argument("--project", default=SOURCE_PROJECT)
    parser.add_argument("--subtype", required=True, choices=sorted(ROUTE_SPECS))
    args = parser.parse_args(argv)
    try:
        result = enqueue_d_route_notification(
            args.ticket, args.subtype, project=args.project)
    except ValueError as exc:
        print(json.dumps(
            {"status": "invalid", "error": str(exc)}, ensure_ascii=False))
        return 2
    print(json.dumps(result, ensure_ascii=False, sort_keys=True))
    if result.get("skipped"):
        # A policy skip is a successful decision, not a delivery failure: the
        # caller must not retry it. `notification_complete` stays false so the
        # finalizer still reports the DM as not sent.
        return 0
    return 0 if result["durable"] else 1


if __name__ == "__main__":
    raise SystemExit(main())

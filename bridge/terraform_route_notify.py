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

from bridge.helpers import aone
from bridge.helpers import dingtalk


SOURCE_PROJECT = "1086837"
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
    return 0 if result["durable"] else 1


if __name__ == "__main__":
    raise SystemExit(main())

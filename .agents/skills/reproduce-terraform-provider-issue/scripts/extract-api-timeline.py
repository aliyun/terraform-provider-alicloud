#!/usr/bin/env python3
"""Extract a caller-configured, allowlisted API timeline from Provider debug logs."""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any, Iterable


HEADER_RE = re.compile(
    r"^(?P<timestamp>\S+).*?: \*{15} (?P<api>[A-Za-z0-9_]+) Response \*{15}"
)
LOG_PREFIX_RE = re.compile(
    r"^.*provider\.terraform-provider-alicloud[^:]*: "
)

FORBIDDEN_FIELD_RE = re.compile(
    r"accesskey|secret|token|password|signature|cookie|authorization|userdata|privatekey",
    re.IGNORECASE,
)


def strip_log_prefix(line: str) -> str:
    return LOG_PREFIX_RE.sub("", line.rstrip("\n"), count=1)


def parse_json_payload(line: str) -> Any | None:
    payload = strip_log_prefix(line)
    payload = payload.split("Domain:", 1)[0].strip()
    try:
        value = json.loads(payload)
    except json.JSONDecodeError:
        return None
    if isinstance(value, str):
        try:
            value = json.loads(value)
        except json.JSONDecodeError:
            return None
    return value


def walk(value: Any) -> Iterable[tuple[str, Any]]:
    if isinstance(value, dict):
        for key, child in value.items():
            yield key, child
            yield from walk(child)
    elif isinstance(value, list):
        for child in value:
            yield from walk(child)


def unique_pairs(value: Any, allowed_keys: tuple[str, ...]) -> list[str]:
    seen: set[str] = set()
    result: list[str] = []
    for key, child in walk(value):
        if key not in allowed_keys or isinstance(child, (dict, list)):
            continue
        rendered = json.dumps(child, ensure_ascii=False, separators=(",", ":"))
        pair = f"{key}={rendered}"
        if pair not in seen:
            seen.add(pair)
            result.append(pair)
    return result


def allowlisted_request(value: Any, allowed_keys: tuple[str, ...]) -> dict[str, Any]:
    if not isinstance(value, dict):
        return {}
    return {key: value[key] for key in allowed_keys if key in value}


def observations(api: str, response: Any, observed_fields: tuple[str, ...]) -> list[str]:
    result: list[str] = []
    for specification in observed_fields:
        api_filter, separator, field = specification.partition(":")
        if not separator:
            field = api_filter
        elif api_filter != api:
            continue
        values = unique_pairs(response, (field,))
        result.extend(values if values else [f"{field}=missing"])
    return result


def extract_events(
    lines: list[str],
    request_fields: tuple[str, ...],
    target_fields: tuple[str, ...],
    observed_fields: tuple[str, ...],
) -> list[dict[str, Any]]:
    events: list[dict[str, Any]] = []
    for index, line in enumerate(lines):
        match = HEADER_RE.search(line)
        if not match:
            continue

        response: Any | None = None
        request: Any | None = None
        for candidate in lines[index + 1 : index + 14]:
            if HEADER_RE.search(candidate):
                break
            parsed = parse_json_payload(candidate)
            if parsed is None:
                continue
            if response is None and isinstance(parsed, dict) and "RequestId" in parsed:
                response = parsed
                continue
            if request is None and isinstance(parsed, dict):
                request = parsed

        response = response if isinstance(response, dict) else {}
        request_id = str(response.get("RequestId", "-"))
        event = {
            "timestamp": match.group("timestamp"),
            "api": match.group("api"),
            "request_id": request_id,
            "targets": unique_pairs(response, target_fields),
            "statuses": unique_pairs(response, ("Status",)),
            "request": allowlisted_request(request, request_fields),
            "observations": observations(match.group("api"), response, observed_fields),
        }
        events.append(event)
    return events


def markdown_escape(value: str) -> str:
    return value.replace("|", "\\|").replace("\n", " ")


def render_markdown(events: list[dict[str, Any]]) -> str:
    rows = [
        "| Timestamp | API | RequestId | Targets | Status | Allowlisted request | Observation |",
        "|---|---|---|---|---|---|---|",
    ]
    for event in events:
        request = json.dumps(
            event["request"], ensure_ascii=False, separators=(",", ":"), sort_keys=True
        )
        values = (
            event["timestamp"],
            event["api"],
            event["request_id"],
            ", ".join(event["targets"]) or "-",
            ", ".join(event["statuses"]) or "-",
            request if request != "{}" else "-",
            ", ".join(event["observations"]) or "-",
        )
        rows.append("| " + " | ".join(markdown_escape(str(value)) for value in values) + " |")
    return "\n".join(rows)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("log", type=Path, help="Raw Terraform/Provider debug log")
    parser.add_argument(
        "--format", choices=("markdown", "jsonl"), default="markdown"
    )
    parser.add_argument(
        "--request-field",
        action="append",
        default=[],
        metavar="NAME",
        help="Allowlisted top-level request field to emit; repeat for multiple fields",
    )
    parser.add_argument(
        "--target-field",
        action="append",
        default=[],
        metavar="NAME",
        help="Allowlisted response field to emit as a target; repeat for multiple fields",
    )
    parser.add_argument(
        "--observe-field",
        action="append",
        default=[],
        metavar="NAME",
        help="Allowlisted response field to report as value or missing; use API:NAME to scope it",
    )
    args = parser.parse_args()

    if not args.log.is_file():
        parser.error(f"log file does not exist: {args.log}")

    selected_fields = args.request_field + args.target_field + args.observe_field
    forbidden = sorted({field for field in selected_fields if FORBIDDEN_FIELD_RE.search(field)})
    if forbidden:
        parser.error("refusing sensitive field selection: " + ", ".join(forbidden))

    events = extract_events(
        args.log.read_text(encoding="utf-8", errors="replace").splitlines(),
        tuple(dict.fromkeys(args.request_field)),
        tuple(dict.fromkeys(args.target_field)),
        tuple(dict.fromkeys(args.observe_field)),
    )
    if not events:
        print("extract-api-timeline: no Provider response events found", file=sys.stderr)
        return 2

    if args.format == "jsonl":
        for event in events:
            print(json.dumps(event, ensure_ascii=False, sort_keys=True))
    else:
        print(render_markdown(events))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

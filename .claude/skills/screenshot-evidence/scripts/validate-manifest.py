#!/usr/bin/env python3
"""Validate the Terraform three-layer visual-evidence handoff manifest."""

from __future__ import annotations

import argparse
import sys
from pathlib import Path


REQUIRED_LAYERS = ("OpenAPI", "CloudSpec/ACube", "Provider")
REQUIRED_COLUMNS = ("layer", "result", "screenshot", "source", "note")
VALID_RESULTS = {"pass", "fail", "n-a"}


def _cells(line: str) -> list[str]:
    return [cell.strip() for cell in line.strip().strip("|").split("|")]


def _is_separator(cells: list[str]) -> bool:
    return bool(cells) and all(cell.replace("-", "").replace(":", "") == "" for cell in cells)


def validate(manifest: Path) -> list[str]:
    errors: list[str] = []
    if not manifest.is_file():
        return [f"manifest not found: {manifest}"]

    lines = manifest.read_text(encoding="utf-8").splitlines()
    header_index = -1
    header: list[str] = []
    for index, line in enumerate(lines):
        cells = _cells(line)
        if [cell.lower() for cell in cells] == list(REQUIRED_COLUMNS):
            header_index = index
            header = cells
            break
    if header_index < 0:
        return ["missing markdown table header: " + " | ".join(REQUIRED_COLUMNS)]

    rows: dict[str, dict[str, str]] = {}
    for line in lines[header_index + 1 :]:
        if not line.strip().startswith("|"):
            if rows:
                break
            continue
        cells = _cells(line)
        if _is_separator(cells):
            continue
        if len(cells) != len(header):
            errors.append(f"malformed table row: {line.strip()}")
            continue
        row = dict(zip(REQUIRED_COLUMNS, cells))
        layer = row["layer"]
        if layer not in REQUIRED_LAYERS:
            errors.append(f"unexpected layer: {layer or '<empty>'}")
            continue
        if layer in rows:
            errors.append(f"duplicate layer: {layer}")
            continue
        rows[layer] = row

    for layer in REQUIRED_LAYERS:
        row = rows.get(layer)
        if row is None:
            errors.append(f"missing layer: {layer}")
            continue
        result = row["result"].lower()
        if result not in VALID_RESULTS:
            errors.append(f"{layer}: result must be pass, fail, or n-a")
            continue
        if not row["source"] or row["source"].startswith("<"):
            errors.append(f"{layer}: source is required")
        if not row["note"] or row["note"].startswith("<"):
            errors.append(f"{layer}: note or N/A reason is required")
        screenshot = row["screenshot"]
        if result == "n-a":
            if screenshot.upper() != "N/A":
                errors.append(f"{layer}: n-a screenshot must be N/A")
            continue
        screenshot_path = Path(screenshot).expanduser()
        if not screenshot_path.is_absolute():
            errors.append(f"{layer}: screenshot path must be absolute")
        elif not screenshot_path.is_file():
            errors.append(f"{layer}: screenshot not found: {screenshot_path}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate a Terraform OpenAPI/CloudSpec/Provider evidence manifest.")
    parser.add_argument("manifest", type=Path)
    args = parser.parse_args()
    errors = validate(args.manifest)
    if errors:
        for error in errors:
            print(f"invalid visual evidence manifest: {error}", file=sys.stderr)
        return 1
    print(f"valid visual evidence manifest: {args.manifest}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

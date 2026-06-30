#!/usr/bin/env python3
"""Compare Cloudspec-generated Terraform files with handwritten provider files."""

from __future__ import annotations

import argparse
import difflib
import json
import re
import subprocess
import sys
from pathlib import Path
from typing import Any


def docs_name(resource: str) -> str:
    return resource.removeprefix("alicloud_")


def expected_paths(resource: str, include_provider: bool) -> list[str]:
    paths = [
        f"alicloud/resource_{resource}.go",
        f"alicloud/resource_{resource}_test.go",
        f"website/docs/r/{docs_name(resource)}.html.markdown",
    ]
    if include_provider:
        paths.append("alicloud/provider.go")
    return paths


def discover_generated_paths(generated_dir: Path, resource: str) -> set[str]:
    if not generated_dir.is_dir():
        return set()

    docs = docs_name(resource)
    discovered: set[str] = set()
    for path in generated_dir.rglob("*"):
        if not path.is_file():
            continue
        rel = path.relative_to(generated_dir).as_posix()
        if resource in rel or docs in rel or rel.startswith("alicloud/service_"):
            discovered.add(rel)
    return discovered


def read_generated(generated_dir: Path, rel_path: str) -> str | None:
    path = generated_dir / rel_path
    if not path.is_file():
        return None
    return path.read_text(encoding="utf-8", errors="replace")


def read_from_dir(root: Path, rel_path: str) -> str | None:
    path = root / rel_path
    if not path.is_file():
        return None
    return path.read_text(encoding="utf-8", errors="replace")


def read_from_git(repo: Path, ref: str, rel_path: str) -> str | None:
    proc = subprocess.run(
        ["git", "-C", str(repo), "show", f"{ref}:{rel_path}"],
        check=False,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    if proc.returncode != 0:
        return None
    return proc.stdout


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def walk_dicts(value: Any) -> list[dict[str, Any]]:
    found: list[dict[str, Any]] = []
    if isinstance(value, dict):
        found.append(value)
        for item in value.values():
            found.extend(walk_dicts(item))
    elif isinstance(value, list):
        for item in value:
            found.extend(walk_dicts(item))
    return found


def condition_leaves(condition: Any) -> list[dict[str, Any]]:
    if isinstance(condition, dict):
        if any(key.startswith("notExistCheck") for key in condition):
            return [condition]
        leaves: list[dict[str, Any]] = []
        for key in ("allOf", "anyOf"):
            value = condition.get(key)
            if isinstance(value, list):
                for item in value:
                    leaves.extend(condition_leaves(item))
        return leaves
    if isinstance(condition, list):
        leaves: list[dict[str, Any]] = []
        for item in condition:
            leaves.extend(condition_leaves(item))
        return leaves
    return []


def resource_not_exist_conditions(resource_type_json: Path | None) -> list[dict[str, Any]]:
    if not resource_type_json:
        return []

    raw = load_json(resource_type_json)
    conditions: list[dict[str, Any]] = []
    for item in walk_dicts(raw):
        condition = item.get("resourceNotExistCondition")
        conditions.extend(condition_leaves(condition))
    return conditions


def generated_go_text(generated_dir: Path) -> str:
    chunks: list[str] = []
    for path in sorted(generated_dir.rglob("*.go")):
        if path.is_file():
            chunks.append(path.read_text(encoding="utf-8", errors="replace"))
    return "\n".join(chunks)


def semantic_checks(generated_dir: Path, resource_type_json: Path | None) -> list[tuple[str, str]]:
    if not resource_type_json:
        return [("INFO", "resourceTypeCode JSON not provided; semantic checks skipped")]

    conditions = resource_not_exist_conditions(resource_type_json)
    if not conditions:
        return [("INFO", "no resourceNotExistCondition found in resourceTypeCode JSON")]

    go_text = generated_go_text(generated_dir)
    results: list[tuple[str, str]] = []
    for condition in conditions:
        target_type = str(condition.get("notExistCheckTargetValueType", ""))
        target_value = condition.get("notExistCheckTargetValue")
        prop = condition.get("notExistCheckProperty") or condition.get("notExistCheckPath") or "<unknown>"
        if target_type != "assertNotEqual" or target_value is None:
            results.append(("INFO", f"resourceNotExistCondition {prop} uses {target_type or '<unknown>'}; no local rule"))
            continue

        value = str(target_value)
        literal = re.escape(value)
        bad = re.search(rf"fmt\.Sprint\([^)]*\)\s*==\s*[\"`]{literal}[\"`]", go_text)
        good = re.search(rf"fmt\.Sprint\([^)]*\)\s*!=\s*[\"`]{literal}[\"`]", go_text)
        if bad:
            results.append(
                (
                    "WARN",
                    (
                        f"resourceNotExistCondition {prop} assertNotEqual {value!r}: generated Go "
                        "appears to use == before NotFound; expected != for not-exist"
                    ),
                )
            )
        elif good:
            results.append(
                (
                    "OK",
                    f"resourceNotExistCondition {prop} assertNotEqual {value!r}: generated Go uses !=",
                )
            )
        else:
            results.append(
                (
                    "WARN",
                    (
                        f"resourceNotExistCondition {prop} assertNotEqual {value!r}: no generated "
                        "fmt.Sprint comparison found"
                    ),
                )
            )
    return results


def summarize_presence(
    resource: str,
    paths: set[str],
    generated_dir: Path,
    read_handwritten,
    include_provider: bool,
) -> list[tuple[str, str, str]]:
    docs = docs_name(resource)
    service_paths = sorted(p for p in paths if p.startswith("alicloud/service_") and p.endswith(".go"))
    checks = [
        ("resource body", f"alicloud/resource_{resource}.go"),
        ("acceptance test", f"alicloud/resource_{resource}_test.go"),
        ("website doc", f"website/docs/r/{docs}.html.markdown"),
    ]
    if include_provider:
        checks.append(("provider registration", "alicloud/provider.go"))

    rows: list[tuple[str, str, str]] = []
    for label, rel_path in checks:
        generated = read_generated(generated_dir, rel_path) is not None
        handwritten = read_handwritten(rel_path) is not None
        if generated and handwritten:
            status = "generated+handwritten"
        elif generated:
            status = "generated only"
        elif handwritten:
            status = "handwritten only"
        else:
            status = "missing"
        rows.append((label, status, rel_path))

    if include_provider is False:
        generated = read_generated(generated_dir, "alicloud/provider.go") is not None
        handwritten = read_handwritten("alicloud/provider.go") is not None
        status = "comparison skipped"
        if generated and handwritten:
            status = "present; comparison skipped"
        elif generated:
            status = "generated present; comparison skipped"
        elif handwritten:
            status = "handwritten present; comparison skipped"
        rows.append(("provider registration", status, "alicloud/provider.go"))

    if service_paths:
        for rel_path in service_paths:
            generated = read_generated(generated_dir, rel_path) is not None
            handwritten = read_handwritten(rel_path) is not None
            if generated and handwritten:
                status = "generated+handwritten"
            elif generated:
                status = "generated only"
            elif handwritten:
                status = "handwritten only"
            else:
                status = "missing"
            rows.append(("service helper", status, rel_path))
    else:
        rows.append(("service helper", "missing", "alicloud/service_*.go"))

    return rows


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description=(
            "Compare Cloudspec generated Terraform provider files with a handwritten "
            "provider directory or git ref."
        )
    )
    parser.add_argument("--resource", required=True, help="Terraform resource type, e.g. alicloud_oss_bucket")
    parser.add_argument("--generated-dir", type=Path, help="Directory produced by cloudspec terraform")
    parser.add_argument(
        "--acube-dir",
        type=Path,
        help="Directory produced by tools/acube_terraform_generate.py; uses its generated/ and raw JSON files",
    )
    parser.add_argument(
        "--resource-type-json",
        type=Path,
        help="resourceTypeCode/get raw JSON for semantic checks. Defaults to <acube-dir>/resource_type_code_get.json.",
    )
    parser.add_argument("--handwritten-dir", type=Path, help="Directory containing handwritten provider files")
    parser.add_argument("--provider-repo", type=Path, help="Git repo containing handwritten provider files")
    parser.add_argument("--handwritten-ref", default="HEAD", help="Git ref to read from --provider-repo")
    parser.add_argument("--extra-path", action="append", default=[], help="Additional relative path to compare")
    parser.add_argument(
        "--include-provider",
        action="store_true",
        help="Also compare alicloud/provider.go. This can be noisy for full provider checkouts.",
    )
    return parser


def validate_args(args: argparse.Namespace, parser: argparse.ArgumentParser) -> None:
    if bool(args.handwritten_dir) == bool(args.provider_repo):
        parser.error("provide exactly one of --handwritten-dir or --provider-repo")
    if args.acube_dir:
        if not args.acube_dir.is_dir():
            parser.error(f"--acube-dir does not exist or is not a directory: {args.acube_dir}")
        if args.generated_dir is None:
            args.generated_dir = args.acube_dir / "generated"
        if args.resource_type_json is None and (args.acube_dir / "resource_type_code_get.json").is_file():
            args.resource_type_json = args.acube_dir / "resource_type_code_get.json"
    if args.generated_dir is None:
        parser.error("provide --generated-dir or --acube-dir")
    if not args.generated_dir.is_dir():
        parser.error(f"--generated-dir does not exist or is not a directory: {args.generated_dir}")
    if args.resource_type_json and not args.resource_type_json.is_file():
        parser.error(f"--resource-type-json does not exist or is not a file: {args.resource_type_json}")
    if args.handwritten_dir and not args.handwritten_dir.is_dir():
        parser.error(f"--handwritten-dir does not exist or is not a directory: {args.handwritten_dir}")
    if args.provider_repo and not (args.provider_repo / ".git").exists():
        parser.error(f"--provider-repo does not look like a git repository: {args.provider_repo}")


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    validate_args(args, parser)

    paths = set(expected_paths(args.resource, args.include_provider))
    paths.update(args.extra_path)
    paths.update(discover_generated_paths(args.generated_dir, args.resource))

    def read_handwritten(rel_path: str) -> str | None:
        if args.handwritten_dir:
            return read_from_dir(args.handwritten_dir, rel_path)
        return read_from_git(args.provider_repo, args.handwritten_ref, rel_path)

    only_generated: list[str] = []
    only_handwritten: list[str] = []
    changed: list[tuple[str, str]] = []
    identical: list[str] = []

    for rel_path in sorted(paths):
        generated = read_generated(args.generated_dir, rel_path)
        handwritten = read_handwritten(rel_path)

        if generated is None and handwritten is None:
            continue
        if generated is None:
            only_handwritten.append(rel_path)
            continue
        if handwritten is None:
            only_generated.append(rel_path)
            continue
        if generated == handwritten:
            identical.append(rel_path)
            continue

        diff = "".join(
            difflib.unified_diff(
                generated.splitlines(keepends=True),
                handwritten.splitlines(keepends=True),
                fromfile=f"generated/{rel_path}",
                tofile=f"handwritten/{rel_path}",
            )
        )
        changed.append((rel_path, diff))

    print(f"Resource: {args.resource}")
    print(f"Generated dir: {args.generated_dir}")
    if args.handwritten_dir:
        print(f"Handwritten dir: {args.handwritten_dir}")
    else:
        print(f"Provider repo: {args.provider_repo}")
        print(f"Handwritten ref: {args.handwritten_ref}")
    if args.acube_dir:
        print(f"Acube dir: {args.acube_dir}")
    if args.resource_type_json:
        print(f"resourceTypeCode JSON: {args.resource_type_json}")
    print()

    print("Structured summary:")
    for label, status, rel_path in summarize_presence(
        args.resource,
        paths,
        args.generated_dir,
        read_handwritten,
        args.include_provider,
    ):
        print(f"  {label}: {status} ({rel_path})")
    print()

    print("Semantic checks:")
    for level, message in semantic_checks(args.generated_dir, args.resource_type_json):
        print(f"  {level}: {message}")
    print()

    print("Only in generated:")
    if only_generated:
        for rel_path in only_generated:
            print(f"  {rel_path}")
    else:
        print("  (none)")
    print()

    print("Only in handwritten:")
    if only_handwritten:
        for rel_path in only_handwritten:
            print(f"  {rel_path}")
    else:
        print("  (none)")
    print()

    print("Identical:")
    if identical:
        for rel_path in identical:
            print(f"  {rel_path}")
    else:
        print("  (none)")
    print()

    if changed:
        for rel_path, diff in changed:
            print(f"Diff: {rel_path}")
            sys.stdout.write(diff)
            if not diff.endswith("\n"):
                print()
            print()
    else:
        print("Diff: (none)")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())

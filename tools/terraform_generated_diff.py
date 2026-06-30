#!/usr/bin/env python3
"""Compare Cloudspec-generated Terraform files with handwritten provider files."""

from __future__ import annotations

import argparse
import difflib
import subprocess
import sys
from pathlib import Path


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
        if resource in rel or docs in rel:
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


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description=(
            "Compare Cloudspec generated Terraform provider files with a handwritten "
            "provider directory or git ref."
        )
    )
    parser.add_argument("--resource", required=True, help="Terraform resource type, e.g. alicloud_oss_bucket")
    parser.add_argument("--generated-dir", required=True, type=Path, help="Directory produced by cloudspec terraform")
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
    if not args.generated_dir.is_dir():
        parser.error(f"--generated-dir does not exist or is not a directory: {args.generated_dir}")
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

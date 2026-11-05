#!/usr/bin/env python3
# Copyright (c) Alibaba, Inc.
# SPDX-License-Identifier: MPL-2.0
"""Generate CHANGELOG entries from git history since the last release.

The provider publishes a file-based CHANGELOG.md and a GitHub Release body
that share the same entry shape:

    - <summary>. ([#<pr>](https://github.com/aliyun/terraform-provider-alicloud/issues/<pr>))

Summaries are derived from squash-merge commit subjects (one commit per merged
PR on master), which carry the PR title followed by a trailing "(#NNNN)".

Usage:

    # Emit the generated markdown body for the version to stdout.
    python3 scripts/release/generate_changelog.py \
        --from v2.1.0 --version 2.2.0

    # Merge generated entries into the CHANGELOG.md "Unreleased" section in
    # place (existing manual entries without a PR reference are preserved).
    python3 scripts/release/generate_changelog.py \
        --from v2.1.0 --version 2.2.0 --update

    # Finalize a release: stamp the version header with today's date, prepend
    # a fresh empty "Unreleased" section, and keep generated entries.
    python3 scripts/release/generate_changelog.py \
        --from v2.1.0 --version 2.2.0 --release

Only the Python standard library is used so the script runs on any CI runner
without extra dependencies.
"""

from __future__ import annotations

import argparse
import datetime
import os
import re
import subprocess
import sys

PROVIDER_ISSUES_URL = "https://github.com/aliyun/terraform-provider-alicloud/issues"

# A trailing PR reference such as " (#11479)" at the end of a commit subject.
PR_REF_RE = re.compile(r"\s*\(#(\d+)\)\s*$")

# Match a bare "[GH-NNNN]" shorthand that changelog-links.sh rewrites later.
GH_SHORT_RE = re.compile(r"\[GH-(\d+)\]")

# Section ordering mirrors the existing CHANGELOG.md and the v2.x release body.
SECTION_ORDER = ["FEATURES", "ENHANCEMENTS", "BUG FIXES"]

# Title prefix -> section. Order matters: most specific first.
PREFIX_RULES = [
    ("New Resource:", "FEATURES"),
    ("New Data Source:", "FEATURES"),
    ("New Ephemeral Resource:", "FEATURES"),
    ("resource/alicloud_", "ENHANCEMENTS"),
    ("data-source/alicloud_", "ENHANCEMENTS"),
    ("docs:", "ENHANCEMENTS"),
    ("website:", "ENHANCEMENTS"),
    ("testcase:", "ENHANCEMENTS"),
]

# Prefixes that describe release tooling / chores and are not user-facing.
SKIP_PREFIXES = ("ci:", "chore:", "refactor:", "build:", "deps:", "release:")

# Words that reclassify an ENHANCEMENTS entry as a BUG FIX.
FIX_HINTS = (
    "fix",
    "fixes",
    "fixed",
    "bug",
    "bugfix",
    "panic",
    "crash",
    "regression",
    "correct",
    "repair",
    "restore",
    "tolerate",
    "prevent",
    "guard",
    "suppress",
    "normalize",
    "reconcile",
)


def run_git(args, cwd=None):
    """Run a git command and return stdout as text."""
    result = subprocess.run(
        ["git"] + args,
        check=True,
        capture_output=True,
        text=True,
        cwd=cwd,
    )
    return result.stdout


def resolve_from_ref(from_ref, cwd=None):
    """Resolve the --from argument to a full commit SHA.

    Accepts a tag name (e.g. v2.1.0), a short SHA, or a full SHA. Falls back to
    the most recent v* tag when nothing is given.
    """
    if not from_ref:
        tags = run_git(
            ["tag", "--list", "v*", "--sort=-creatordate"], cwd=cwd
        ).splitlines()
        if not tags:
            raise SystemExit(
                "generate-changelog: no previous release tag found; "
                "pass --from <commit> explicitly"
            )
        from_ref = tags[0]
    sha = run_git(["rev-parse", "--verify", f"{from_ref}^{{commit}}"], cwd=cwd).strip()
    return from_ref, sha


def parse_commits(from_ref, head_ref="HEAD", cwd=None):
    """Return a list of (pr_number_or_None, subject) tuples since from_ref."""
    log = run_git(
        ["log", "--no-merges", "--format=%s", f"{from_ref}..{head_ref}"],
        cwd=cwd,
    )
    entries = []
    for subject in log.splitlines():
        subject = subject.strip()
        if not subject:
            continue
        m = PR_REF_RE.search(subject)
        pr = int(m.group(1)) if m else None
        title = subject[: m.start()].strip() if m else subject
        # Expand [GH-NNNN] shorthand into an explicit PR reference so the
        # generated entry is already link-ready and changelog-links.sh is a
        # no-op for it.
        title = GH_SHORT_RE.sub(r"#\1", title)
        entries.append((pr, title))
    return entries


def categorize(title):
    """Map a PR title to a section name, or None when it should be skipped."""
    lower = title.lower()
    if any(title.startswith(p) for p in SKIP_PREFIXES):
        return None
    section = "ENHANCEMENTS"
    for prefix, sec in PREFIX_RULES:
        if title.startswith(prefix):
            section = sec
            break
    # Reclassify enhancement-flavored entries that are really bug fixes.
    if section == "ENHANCEMENTS" and title.lower().startswith("fix"):
        section = "BUG FIXES"
    if section == "ENHANCEMENTS" and any(
        hint in lower for hint in FIX_HINTS
    ):
        # Only promote to BUG FIXES when the hint appears in the leading
        # description, not inside a resource name.
        head = lower.split(":", 1)[1] if ":" in lower else lower
        if any(hint in head for hint in FIX_HINTS):
            section = "BUG FIXES"
    return section


def ensure_period(text):
    """Ensure a summary ends with a single period."""
    text = text.rstrip()
    if not text:
        return text
    if text.endswith(("...", ".", "!", "?")):
        return text
    return text + "."


def entry_line(pr, title):
    """Format one changelog entry line."""
    summary = ensure_period(title)
    if pr:
        return f"- {summary} ([#{pr}]({PROVIDER_ISSUES_URL}/{pr}))"
    return f"- {summary}"


def build_body(entries):
    """Build the markdown body (sections) from parsed entries."""
    buckets = {s: [] for s in SECTION_ORDER}
    seen_pr = set()
    for pr, title in entries:
        section = categorize(title)
        if section is None:
            continue
        if pr is not None and pr in seen_pr:
            continue
        if pr is not None:
            seen_pr.add(pr)
        buckets[section].append(entry_line(pr, title))

    parts = []
    for section in SECTION_ORDER:
        lines = buckets[section]
        if not lines:
            continue
        parts.append(f"{section}:\n")
        parts.append("\n".join(lines))
        parts.append("")  # blank line between sections
    return "\n".join(parts).rstrip() + "\n"


def read_changelog(path):
    with open(path, "r", encoding="utf-8") as fh:
        return fh.read()


def split_unreleased(content):
    """Split CHANGELOG.md into (pre, unreleased_body, post).

    The top-most "## <version> (Unreleased)" section is treated as the
    unreleased block. Returns (pre, body, post) where body excludes the
    header line; post starts at the next "## " header. When no Unreleased
    section exists, pre is the whole file and body/post are empty.
    """
    lines = content.splitlines(keepends=True)
    header_idx = None
    for i, line in enumerate(lines):
        if line.startswith("## ") and "(Unreleased)" in line:
            header_idx = i
            break
    if header_idx is None:
        return content, "", ""
    pre = "".join(lines[: header_idx + 1])
    rest = lines[header_idx + 1 :]
    end_idx = len(rest)
    for j, line in enumerate(rest):
        if line.startswith("## "):
            end_idx = j
            break
    body = "".join(rest[:end_idx])
    post = "".join(rest[end_idx:])
    return pre, body, post


def parse_existing_prs(body):
    """Extract PR numbers already referenced in an existing CHANGELOG body."""
    return set(int(m) for m in re.findall(r"\[#(\d+)\]", body))


def merge_entries(existing_body, generated_body, entries):
    """Merge generated entries into an existing Unreleased body.

    Preserves manual entries (no PR number) and appends new generated entries
    whose PR number is not already present, preserving the existing order.
    """
    existing_prs = parse_existing_prs(existing_body)
    # Rebuild the generated body with only new entries.
    new_entries = [(pr, title) for pr, title in entries if pr is None or pr not in existing_prs]
    addendum = build_body(new_entries)
    existing_body = existing_body.rstrip()
    if not addendum.strip():
        return existing_body + "\n"
    if not existing_body:
        return addendum
    return existing_body + "\n\n" + addendum


def next_unreleased_header(version):
    """Compute the next Unreleased header after releasing `version`.

    Increments the minor component (2.2.0 -> 2.3.0).
    """
    m = re.match(r"^(\d+)\.(\d+)\.(\d+)", version)
    if not m:
        return "## Unreleased\n"
    major, minor, patch = int(m.group(1)), int(m.group(2)), int(m.group(3))
    return f"## {major}.{minor + 1}.0 (Unreleased)\n"


def release_date():
    """Return today's date in the CHANGELOG.md style: 'Month D, YYYY'."""
    return datetime.date.today().strftime("%B %-d, %Y")


def cmd_emit(args):
    from_ref, _ = resolve_from_ref(args.from_ref, cwd=args.repo)
    entries = parse_commits(from_ref, head_ref=args.head, cwd=args.repo)
    body = build_body(entries)
    sys.stdout.write(body)


def cmd_update(args):
    from_ref, _ = resolve_from_ref(args.from_ref, cwd=args.repo)
    entries = parse_commits(from_ref, head_ref=args.head, cwd=args.repo)
    path = os.path.join(args.repo or ".", "CHANGELOG.md")
    content = read_changelog(path)
    pre, body, post = split_unreleased(content)
    if not body:
        # No Unreleased section yet: create one for the requested version.
        pre = f"## {args.version} (Unreleased)\n\n"
        body = ""
        post = content if content else "\n"
    merged = merge_entries(body, build_body(entries), entries)
    new_content = pre.rstrip("\n") + "\n\n" + merged.rstrip("\n") + "\n" + post
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(new_content)
    sys.stderr.write(f"generate-changelog: updated {path}\n")


def cmd_release(args):
    from_ref, _ = resolve_from_ref(args.from_ref, cwd=args.repo)
    entries = parse_commits(from_ref, head_ref=args.head, cwd=args.repo)
    path = os.path.join(args.repo or ".", "CHANGELOG.md")
    content = read_changelog(path)
    pre, body, post = split_unreleased(content)
    generated = build_body(entries)
    if body.strip():
        final_body = merge_entries(body, generated, entries)
    else:
        final_body = generated
    date = release_date()
    # Stamp the version header and prepend a fresh Unreleased section.
    stamped_header = re.sub(
        r"## .*?\(Unreleased\)\s*",
        f"## {args.version} ({date})\n",
        pre,
        count=1,
    )
    new_unreleased = next_unreleased_header(args.version)
    new_content = (
        new_unreleased
        + "\n"
        + stamped_header.rstrip("\n")
        + "\n\n"
        + final_body.rstrip("\n")
        + "\n"
        + post
    )
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(new_content)
    sys.stderr.write(
        f"generate-changelog: finalized {args.version} ({date}) in {path}\n"
    )


def main(argv=None):
    parser = argparse.ArgumentParser(
        description="Generate CHANGELOG entries from git history."
    )
    parser.add_argument(
        "--from",
        dest="from_ref",
        default=None,
        help="previous release tag or commit (default: latest v* tag)",
    )
    parser.add_argument(
        "--version",
        required=True,
        help="target release version, e.g. 2.2.0",
    )
    parser.add_argument(
        "--head",
        default="HEAD",
        help="head ref to generate entries from (default: HEAD)",
    )
    parser.add_argument(
        "--repo",
        default=None,
        help="repository root (default: current directory)",
    )
    group = parser.add_mutually_exclusive_group()
    group.add_argument(
        "--update",
        action="store_true",
        help="merge generated entries into CHANGELOG.md Unreleased section",
    )
    group.add_argument(
        "--release",
        action="store_true",
        help="finalize the version: stamp date, prepend a new Unreleased section",
    )
    args = parser.parse_args(argv)

    if not re.match(r"^\d+\.\d+\.\d+(-beta\d+|-rc\d+)?$", args.version):
        raise SystemExit(
            f"generate-changelog: invalid version '{args.version}' "
            "(expected N.N.N or N.N.N-betaN / N.N.N-rcN)"
        )

    if args.update:
        cmd_update(args)
    elif args.release:
        cmd_release(args)
    else:
        cmd_emit(args)


if __name__ == "__main__":
    main()

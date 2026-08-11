#!/usr/bin/env python3
"""Fail-closed detection of CloudSpec ``operations/`` changes in a Git repo.

The task manager captures :class:`RepositoryBaseline` once at the beginning of
an execution epoch and inspects the same repository before allowing a write or
finalizing the task.  Inspection compares the baseline commit with the current
working tree, so committed, staged, and unstaged changes are covered together;
untracked files are added separately.  All Git filename output is NUL-delimited
to preserve spaces, newlines, and Unicode paths.

This module is deliberately local-only.  It neither mutates the repository nor
calls a control-plane service.
"""

from __future__ import annotations

import os
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path, PurePosixPath
from typing import Iterable, Optional, Sequence, Union


_OBJECT_ID_RE = re.compile(r"^[0-9a-fA-F]{40}(?:[0-9a-fA-F]{24})?$")
_CHANGE_STATUSES = frozenset({"A", "M", "D", "R", "C", "T"})
_GIT_TIMEOUT_SECONDS = 5.0


class OperationDiffGuardError(RuntimeError):
    """The repository state cannot be proven safe."""


@dataclass(frozen=True)
class RepositoryBaseline:
    """Immutable Git identity captured at one task epoch's start."""

    repo_root: str
    head: str
    branch: str
    components: tuple[str, ...]


@dataclass(frozen=True)
class OperationChange:
    """One Git name-status record involving a CloudSpec operation path."""

    status: str
    old_path: Optional[str]
    new_path: Optional[str]

    @property
    def paths(self) -> tuple[str, ...]:
        return tuple(
            dict.fromkeys(
                path for path in (self.old_path, self.new_path) if path is not None
            )
        )


@dataclass(frozen=True)
class RepositoryInspection:
    """Structured, non-throwing result returned to the task manager."""

    blocked: bool
    reason: str
    paths: tuple[str, ...]
    current_head: str
    current_branch: str
    changes: tuple[OperationChange, ...] = ()
    components: tuple[str, ...] = ()
    repo_root: str = ""
    baseline_sha: str = ""

    @property
    def has_operation_changes(self) -> bool:
        return bool(self.paths)


@dataclass(frozen=True)
class _RepositoryContext:
    root: Path
    head: str
    branch: str


def _decode(value: bytes) -> str:
    return os.fsdecode(value)


def _display_stderr(value: bytes) -> str:
    text = _decode(value).replace("\x00", "\\0").strip()
    return text[-800:] if text else "unknown Git error"


def _git(root: Path, args: Sequence[str]) -> bytes:
    try:
        completed = subprocess.run(
            ["git", "-C", os.fspath(root), *args],
            stdin=subprocess.DEVNULL,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            check=False,
            timeout=_GIT_TIMEOUT_SECONDS,
        )
    except subprocess.TimeoutExpired as exc:
        raise OperationDiffGuardError(
            "git %s timed out after %.1f seconds"
            % (args[0] if args else "command", _GIT_TIMEOUT_SECONDS)) from exc
    except (OSError, ValueError) as exc:
        raise OperationDiffGuardError("cannot execute Git: %s" % exc) from exc
    if completed.returncode != 0:
        raise OperationDiffGuardError(
            "git %s failed: %s" % (args[0] if args else "command",
                                    _display_stderr(completed.stderr)))
    return completed.stdout


def _single_line(raw: bytes, label: str) -> str:
    value = _decode(raw).strip()
    if not value or "\x00" in value or "\n" in value or "\r" in value:
        raise OperationDiffGuardError("Git returned an invalid %s" % label)
    return value


def _repository_context(repo_root: Union[str, os.PathLike[str]]) -> _RepositoryContext:
    supplied = Path(repo_root).expanduser().resolve()
    if not supplied.is_dir():
        raise OperationDiffGuardError("repository root is not a directory")

    discovered = Path(_single_line(
        _git(supplied, ["rev-parse", "--show-toplevel"]),
        "repository root",
    )).resolve()
    if discovered != supplied:
        raise OperationDiffGuardError(
            "repository root mismatch: expected %s, Git resolved %s"
            % (supplied, discovered))

    head = _single_line(
        _git(discovered, ["rev-parse", "--verify", "HEAD"]), "HEAD")
    if _OBJECT_ID_RE.fullmatch(head) is None:
        raise OperationDiffGuardError("Git returned an invalid HEAD object id")

    try:
        branch = _single_line(
            _git(discovered, ["symbolic-ref", "--quiet", "--short", "HEAD"]),
            "branch",
        )
    except OperationDiffGuardError as exc:
        raise OperationDiffGuardError(
            "repository is detached or its current branch cannot be resolved") from exc
    return _RepositoryContext(discovered, head.lower(), branch)


def _split_nul(raw: bytes, label: str) -> list[str]:
    if raw and not raw.endswith(b"\x00"):
        raise OperationDiffGuardError("Git returned unterminated %s data" % label)
    values = raw.split(b"\x00")
    if values and values[-1] == b"":
        values.pop()
    return [_decode(value) for value in values]


def _component_roots(paths: Iterable[str]) -> tuple[str, ...]:
    components: set[str] = set()
    for value in paths:
        path = PurePosixPath(value)
        if path.name != "main.cspec":
            continue
        parent = path.parent.as_posix()
        components.add("" if parent == "." else parent)
    return tuple(sorted(components))


def _current_components(context: _RepositoryContext) -> tuple[str, ...]:
    paths = _split_nul(
        _git(context.root, ["ls-files", "--cached", "--others",
                            "--exclude-standard", "-z"]),
        "current filename",
    )
    # A staged addition removed again from the working tree must not turn an
    # arbitrary sibling operations/ directory into a live component.
    existing_main_files = (
        path for path in paths
        if PurePosixPath(path).name == "main.cspec"
        and (context.root / Path(path)).is_file()
    )
    return _component_roots(existing_main_files)


def _baseline_components(context: _RepositoryContext,
                         baseline_sha: str) -> tuple[str, ...]:
    paths = _split_nul(
        _git(context.root, ["ls-tree", "-r", "--name-only", "-z",
                            baseline_sha, "--"]),
        "baseline filename",
    )
    return _component_roots(paths)


def discover_cloudspec_components(
        repo_root: Union[str, os.PathLike[str]]) -> tuple[str, ...]:
    """Return current repo-relative component roots containing ``main.cspec``.

    Only tracked or non-ignored untracked files in the supplied Git repository
    participate.  The empty string denotes a component rooted at the repository
    root.  Invalid or ambiguous repository state raises
    :class:`OperationDiffGuardError`, allowing callers to fail closed.
    """
    return _current_components(_repository_context(repo_root))


def _normalize_repo_path(context: _RepositoryContext,
                         value: Union[str, os.PathLike[str]]) -> str:
    raw = os.fspath(value)
    if not raw or "\x00" in raw:
        raise OperationDiffGuardError("operation path is empty or contains NUL")
    candidate = Path(raw)
    if candidate.is_absolute():
        try:
            relative = candidate.resolve(strict=False).relative_to(context.root)
        except ValueError as exc:
            raise OperationDiffGuardError(
                "operation path is outside the repository") from exc
    else:
        relative = candidate
    posix = PurePosixPath(relative.as_posix())
    if posix.is_absolute() or any(part in {"", ".", ".."} for part in posix.parts):
        raise OperationDiffGuardError("operation path is not repo-relative")
    return posix.as_posix()


def _path_is_operation(path: str, components: Iterable[str]) -> bool:
    for component in components:
        prefix = "%soperations/" % (component + "/" if component else "")
        if path.startswith(prefix) and len(path) > len(prefix):
            return True
    return False


def path_is_cloudspec_operation(
        repo_root: Union[str, os.PathLike[str]],
        path: Union[str, os.PathLike[str]],
) -> bool:
    """Whether a prospective path is under a current CloudSpec component.

    The function is intended for write-time interception.  Invalid repositories
    and paths raise :class:`OperationDiffGuardError` rather than returning a
    permissive false result.
    """
    context = _repository_context(repo_root)
    relative = _normalize_repo_path(context, path)
    # Include HEAD as well as the working tree.  Otherwise deleting or staging
    # the deletion of main.cspec immediately before a write would make a real
    # component disappear from prospective-path classification.
    components = tuple(sorted(set(_current_components(context)) |
                              set(_baseline_components(context, context.head))))
    return _path_is_operation(relative, components)


def capture_repository_baseline(
        repo_root: Union[str, os.PathLike[str]]) -> RepositoryBaseline:
    """Capture the immutable repository identity for one task epoch."""
    context = _repository_context(repo_root)
    return RepositoryBaseline(
        repo_root=os.fspath(context.root),
        head=context.head,
        branch=context.branch,
        components=_baseline_components(context, context.head),
    )


def _parse_name_status(raw: bytes) -> tuple[OperationChange, ...]:
    fields = _split_nul(raw, "name-status")
    changes: list[OperationChange] = []
    index = 0
    while index < len(fields):
        token = fields[index]
        index += 1
        if not token:
            raise OperationDiffGuardError("Git returned an empty change status")
        status = token[0]
        if status not in _CHANGE_STATUSES:
            raise OperationDiffGuardError(
                "unsupported Git change status %r" % token)
        if status in {"R", "C"}:
            if index + 1 >= len(fields):
                raise OperationDiffGuardError(
                    "Git returned a truncated %s record" % status)
            old_path, new_path = fields[index], fields[index + 1]
            index += 2
        else:
            if index >= len(fields):
                raise OperationDiffGuardError(
                    "Git returned a truncated %s record" % status)
            path = fields[index]
            index += 1
            old_path = path if status == "D" else None
            new_path = None if status == "D" else path
        if not (old_path or new_path):
            raise OperationDiffGuardError("Git returned an empty change path")
        changes.append(OperationChange(status, old_path, new_path))
    return tuple(changes)


def _validated_baseline(context: _RepositoryContext, baseline_sha: str) -> str:
    if not isinstance(baseline_sha, str) or _OBJECT_ID_RE.fullmatch(baseline_sha) is None:
        raise OperationDiffGuardError("baseline SHA is missing or malformed")
    baseline = baseline_sha.lower()
    resolved = _single_line(
        _git(context.root, ["rev-parse", "--verify", baseline + "^{commit}"]),
        "baseline commit",
    ).lower()
    if resolved != baseline:
        raise OperationDiffGuardError("baseline SHA does not resolve exactly")
    _git(context.root, ["merge-base", "--is-ancestor", baseline, context.head])
    return baseline


def _blocked(reason: str, *, repo_root: str = "", baseline_sha: str = "",
             current_head: str = "", current_branch: str = "") -> RepositoryInspection:
    return RepositoryInspection(
        blocked=True,
        reason=reason,
        paths=(),
        current_head=current_head,
        current_branch=current_branch,
        repo_root=repo_root,
        baseline_sha=baseline_sha,
    )


def _safe_path_text(value: object) -> str:
    try:
        return os.fspath(value)  # type: ignore[arg-type]
    except (TypeError, ValueError, OSError):
        return ""


def inspect_repository(
        repo_root: Union[str, os.PathLike[str]],
        baseline_sha: Union[str, RepositoryBaseline],
        *,
        expected_branch: Optional[str] = None,
) -> RepositoryInspection:
    """Inspect operation changes since ``baseline_sha`` without leaking errors.

    ``baseline_sha`` may be the plain object id persisted by a manager, or the
    :class:`RepositoryBaseline` returned by :func:`capture_repository_baseline`.
    Passing the object also binds inspection to the original repo and branch.
    With a plain SHA, managers should pass their persisted ``expected_branch``.

    Every repository, baseline, branch, conflict, parsing, and subprocess error
    is returned as ``blocked=True`` with a reason; this function does not expose
    an operational exception to its caller.
    """
    context: Optional[_RepositoryContext] = None
    raw_baseline = baseline_sha.head if isinstance(
        baseline_sha, RepositoryBaseline) else baseline_sha
    baseline_text = raw_baseline if isinstance(raw_baseline, str) else ""
    try:
        context = _repository_context(repo_root)
        if isinstance(baseline_sha, RepositoryBaseline):
            if Path(baseline_sha.repo_root).resolve() != context.root:
                raise OperationDiffGuardError(
                    "baseline repository does not match the inspected repository")
            if expected_branch is not None and expected_branch != baseline_sha.branch:
                raise OperationDiffGuardError(
                    "caller branch disagrees with the captured baseline")
            expected_branch = baseline_sha.branch
        if expected_branch is not None:
            if not isinstance(expected_branch, str) or not expected_branch.strip():
                raise OperationDiffGuardError("expected branch is malformed")
            if context.branch != expected_branch:
                raise OperationDiffGuardError(
                    "branch changed since baseline: expected %s, found %s"
                    % (expected_branch, context.branch))

        baseline = _validated_baseline(context, baseline_text)
        if _git(context.root, ["ls-files", "--unmerged", "-z"]):
            raise OperationDiffGuardError(
                "repository has unresolved merge conflicts")

        baseline_components = _baseline_components(context, baseline)
        current_components = _current_components(context)
        component_roots = tuple(sorted(set(baseline_components) |
                                       set(current_components)))
        changes = list(_parse_name_status(_git(context.root, [
            "diff", "--name-status", "-z", "--find-renames",
            "--find-copies", "--find-copies-harder", baseline, "--",
        ])))
        untracked = _split_nul(
            _git(context.root, ["ls-files", "--others", "--exclude-standard", "-z"]),
            "untracked filename",
        )
        changes.extend(OperationChange("A", None, path) for path in untracked)

        relevant = tuple(
            change for change in changes
            if any(_path_is_operation(path, component_roots)
                   for path in change.paths)
        )
        paths = tuple(dict.fromkeys(
            path for change in relevant for path in change.paths
        ))
        return RepositoryInspection(
            blocked=False,
            reason="",
            paths=paths,
            current_head=context.head,
            current_branch=context.branch,
            changes=relevant,
            components=component_roots,
            repo_root=os.fspath(context.root),
            baseline_sha=baseline,
        )
    except Exception as exc:  # The public inspection boundary is fail-closed.
        reason = str(exc).strip() or exc.__class__.__name__
        return _blocked(
            reason,
            repo_root=(os.fspath(context.root) if context is not None
                       else _safe_path_text(repo_root)),
            baseline_sha=baseline_text,
            current_head=context.head if context is not None else "",
            current_branch=context.branch if context is not None else "",
        )


__all__ = [
    "OperationChange",
    "OperationDiffGuardError",
    "RepositoryBaseline",
    "RepositoryInspection",
    "capture_repository_baseline",
    "discover_cloudspec_components",
    "inspect_repository",
    "path_is_cloudspec_operation",
]

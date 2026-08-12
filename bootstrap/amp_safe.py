#!/usr/bin/env python3
"""Strict, non-interactive safety wrapper for the AMP CLI.

The wrapper accepts a deliberately small AMP-compatible argv surface, rebuilds
the command from parsed values, and never invokes a shell.  Mutating commands
must first be authorized by the current Jarvis task fence.  Real daily/pre
publishes automatically run one successful dry-run; online publishing remains
a human-only production gate.
"""

from __future__ import annotations

import json
import os
import pwd
import re
import socket
import struct
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Callable, Mapping, Optional, Sequence


_FEATURE_BRANCH_RE = re.compile(
    r"^feature/[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?"
    r"(?:/[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?)*$")
_VALUE_LIMIT = 1024
_FIXED_FLAGS = ("-o", "json", "--no-interactive")
_TESTING_ENV = "JARVIS_AMP_SAFE_TESTING"
_AMP_BIN_ENV = "JARVIS_AMP_BIN"
_HEADLESS_BROKER_ENV = "JARVIS_HEADLESS_AMP_BROKER"
_UNSAFE_ENV_NAMES = frozenset({
    "BASH_ENV", "CDPATH", "ENV", "GIT_CONFIG", "GIT_CONFIG_GLOBAL",
    "GIT_CONFIG_SYSTEM", "GIT_DIR", "GIT_EXEC_PATH", "GIT_SSH",
    "GIT_SSH_COMMAND", "GIT_WORK_TREE", "LD_LIBRARY_PATH", "LD_PRELOAD",
    "PYTHONHOME", "PYTHONPATH", "RUBYLIB",
})
_UNSAFE_ENV_PREFIXES = ("DYLD_", "GIT_CONFIG_KEY_", "GIT_CONFIG_VALUE_")


class AmpSafeError(RuntimeError):
    """The requested AMP invocation is outside the safe allowlist."""


@dataclass(frozen=True)
class AmpInvocation:
    """Canonical AMP argv and the task-fence action it requires."""

    argv: tuple[str, ...]
    authorize_action: Optional[str] = None
    publish_environment: Optional[str] = None
    dry_run: bool = False


Authorizer = Callable[[str, Path], int]
BranchResolver = Callable[[Path], str]


def _value(value: object, label: str) -> str:
    text = str(value)
    if (not text or len(text) > _VALUE_LIMIT or "\x00" in text
            or "\n" in text or "\r" in text):
        raise AmpSafeError("%s is empty or malformed" % label)
    return text


def _feature_branch(value: object) -> str:
    branch = _value(value, "branch")
    if (_FEATURE_BRANCH_RE.fullmatch(branch) is None
            or ".." in branch or "//" in branch):
        raise AmpSafeError("branch must match feature/<short-action>")
    return branch


def _parse_options(
        tokens: Sequence[str], *,
        value_flags: Sequence[str] = (),
        boolean_flags: Sequence[str] = (),
        positional_count: int = 0,
) -> tuple[dict[str, str], set[str], tuple[str, ...]]:
    allowed_values = frozenset(value_flags)
    allowed_booleans = frozenset(boolean_flags)
    values: dict[str, str] = {}
    booleans: set[str] = set()
    positional: list[str] = []
    index = 0
    while index < len(tokens):
        token = _value(tokens[index], "argument")
        if token.startswith("-"):
            name, separator, inline = token.partition("=")
            if name in allowed_booleans:
                if separator or name in booleans:
                    raise AmpSafeError("duplicate or valued flag is not allowed: %s" % name)
                booleans.add(name)
                index += 1
                continue
            if name not in allowed_values:
                raise AmpSafeError("unknown or unsafe AMP flag: %s" % name)
            if name in values:
                raise AmpSafeError("duplicate AMP flag: %s" % name)
            if separator:
                option_value = _value(inline, name)
                index += 1
            else:
                if index + 1 >= len(tokens):
                    raise AmpSafeError("missing value for AMP flag: %s" % name)
                option_value = _value(tokens[index + 1], name)
                if option_value.startswith("-"):
                    raise AmpSafeError("missing value for AMP flag: %s" % name)
                index += 2
            values[name] = option_value
            continue
        positional.append(token)
        index += 1
    if len(positional) != positional_count:
        raise AmpSafeError(
            "expected %d positional argument(s), found %d"
            % (positional_count, len(positional)))
    return values, booleans, tuple(positional)


def _canonical_options(
        values: Mapping[str, str], booleans: set[str],
        value_order: Sequence[str], boolean_order: Sequence[str] = (),
) -> tuple[str, ...]:
    argv: list[str] = []
    for name in value_order:
        if name in values:
            argv.extend((name, values[name]))
    for name in boolean_order:
        if name in booleans:
            argv.append(name)
    return tuple(argv)


_SCOPE_FLAGS = ("--pop-code", "--pop-version", "--project-id")


def _parse_init(tokens: Sequence[str]) -> AmpInvocation:
    order = (*_SCOPE_FLAGS, "--branch", "--api-name")
    values, booleans, _ = _parse_options(tokens, value_flags=order)
    if "--branch" in values:
        values["--branch"] = _feature_branch(values["--branch"])
    if "--project-id" in values and (
            "--pop-code" in values or "--pop-version" in values):
        raise AmpSafeError(
            "init accepts project-id or pop-code/pop-version, not both")
    if ("--pop-code" in values) != ("--pop-version" in values):
        raise AmpSafeError("init requires pop-code and pop-version together")
    return AmpInvocation(
        ("init", *_canonical_options(values, booleans, order)),
        authorize_action="init",
    )


def _parse_branch(tokens: Sequence[str]) -> AmpInvocation:
    if not tokens:
        raise AmpSafeError("branch requires list, get, create, or switch")
    command = tokens[0]
    rest = tokens[1:]
    if command == "list":
        values, booleans, _ = _parse_options(rest, value_flags=_SCOPE_FLAGS)
        return AmpInvocation((
            "branch", "list", *_canonical_options(
                values, booleans, _SCOPE_FLAGS)))
    if command == "get":
        order = (*_SCOPE_FLAGS, "--branch", "--status")
        values, booleans, _ = _parse_options(rest, value_flags=order)
        return AmpInvocation((
            "branch", "get", *_canonical_options(values, booleans, order)))
    if command == "create":
        order = (*_SCOPE_FLAGS, "--branch", "--description")
        values, booleans, _ = _parse_options(rest, value_flags=order)
        if "--branch" not in values:
            raise AmpSafeError("branch create requires --branch")
        values["--branch"] = _feature_branch(values["--branch"])
        return AmpInvocation(
            ("branch", "create", *_canonical_options(
                values, booleans, order)),
            authorize_action="branch-create",
        )
    if command == "switch":
        values, booleans, positional = _parse_options(
            rest, value_flags=_SCOPE_FLAGS, positional_count=1)
        branch = _feature_branch(positional[0])
        return AmpInvocation(
            ("branch", "switch", branch, *_canonical_options(
                values, booleans, _SCOPE_FLAGS)),
            authorize_action="branch-switch",
        )
    raise AmpSafeError("branch subcommand is not allowed: %s" % command)


def _parse_context(tokens: Sequence[str]) -> AmpInvocation:
    if tuple(tokens[:2]) != ("set", "branch"):
        raise AmpSafeError("only 'context set branch <feature/...>' is allowed")
    values, booleans, positional = _parse_options(
        tokens[2:], value_flags=_SCOPE_FLAGS, positional_count=1)
    branch = _feature_branch(positional[0])
    return AmpInvocation(
        ("context", "set", "branch", branch, *_canonical_options(
            values, booleans, _SCOPE_FLAGS)),
        authorize_action="context-set-branch",
    )


def _parse_api(tokens: Sequence[str]) -> AmpInvocation:
    if not tokens:
        raise AmpSafeError("api requires list or get")
    command = tokens[0]
    rest = tokens[1:]
    if command == "list":
        values, booleans, _ = _parse_options(rest, value_flags=_SCOPE_FLAGS)
        return AmpInvocation((
            "api", "list", *_canonical_options(values, booleans, _SCOPE_FLAGS)))
    if command == "get":
        order = (*_SCOPE_FLAGS, "--api-name")
        values, booleans, _ = _parse_options(rest, value_flags=order)
        return AmpInvocation((
            "api", "get", *_canonical_options(values, booleans, order)))
    raise AmpSafeError("api write or unknown subcommand is not allowed: %s" % command)


def _parse_publish(tokens: Sequence[str]) -> AmpInvocation:
    if not tokens:
        raise AmpSafeError("publish requires daily, pre, online, or status")
    environment = tokens[0]
    rest = tokens[1:]
    if environment == "status":
        values, booleans, _ = _parse_options(
            rest, value_flags=("--publish-id",))
        return AmpInvocation(
            ("publish", "status", *_canonical_options(
                values, booleans, ("--publish-id",))),
            authorize_action="publish",
            publish_environment="status",
        )
    if environment not in {"daily", "pre", "online"}:
        raise AmpSafeError("publish environment is not allowed: %s" % environment)
    values, booleans, _ = _parse_options(
        rest, boolean_flags=("--dry-run",))
    return AmpInvocation(
        ("publish", environment,
         *_canonical_options(values, booleans, (), ("--dry-run",))),
        authorize_action="publish",
        publish_environment=environment,
        dry_run="--dry-run" in booleans,
    )


def parse_amp_argv(argv: Sequence[str]) -> AmpInvocation:
    """Parse and canonicalize the complete safe AMP argv surface."""
    tokens = tuple(_value(value, "argument") for value in argv)
    if not tokens:
        raise AmpSafeError("an AMP command is required")
    command = tokens[0]
    rest = tokens[1:]
    if command == "--version":
        if rest:
            raise AmpSafeError("--version accepts no arguments")
        return AmpInvocation(("--version",))
    if command in {"doctor", "whoami"}:
        if rest:
            raise AmpSafeError("%s accepts no wrapper arguments" % command)
        return AmpInvocation((command,))
    if command == "init":
        return _parse_init(rest)
    if command == "branch":
        return _parse_branch(rest)
    if command == "context":
        return _parse_context(rest)
    if command == "api":
        return _parse_api(rest)
    if command == "publish":
        return _parse_publish(rest)
    raise AmpSafeError("AMP command is not allowed: %s" % command)


def _amp_binary(environ: Mapping[str, str]) -> str:
    override = environ.get(_AMP_BIN_ENV, "").strip()
    testing = environ.get(_TESTING_ENV) == "1"
    if override:
        if not testing:
            raise AmpSafeError(
                "%s is accepted only when %s=1" % (_AMP_BIN_ENV, _TESTING_ENV))
        candidate = Path(override)
        if not candidate.is_absolute():
            raise AmpSafeError("test AMP binary override must be an absolute path")
        resolved = os.fspath(candidate.resolve())
    else:
        # AMP's installer places the binary here on every Jarvis host.  Never
        # consult Agent-controlled PATH: otherwise a repository-local fake
        # executable could inherit the wrapper's trusted authority.
        try:
            account_home = Path(pwd.getpwuid(os.getuid()).pw_dir).resolve()
        except (KeyError, OSError, RuntimeError) as exc:
            raise AmpSafeError("cannot resolve the local account home: %s" % exc) from exc
        resolved = os.fspath((account_home / ".local" / "bin" / "amp").resolve())
    if not resolved or not Path(resolved).is_file() or not os.access(resolved, os.X_OK):
        raise AmpSafeError("amp executable was not found or is not executable")
    return resolved


def _account_home() -> Path:
    try:
        return Path(pwd.getpwuid(os.getuid()).pw_dir).resolve()
    except (KeyError, OSError, RuntimeError) as exc:
        raise AmpSafeError("cannot resolve the local account home: %s" % exc) from exc


def _worker_environment(environ: Mapping[str, str]) -> dict[str, str]:
    """Preserve Worker routing while removing interpreter/process injection."""
    result = dict(environ)
    for name in tuple(result):
        if name in _UNSAFE_ENV_NAMES or name.startswith(_UNSAFE_ENV_PREFIXES):
            result.pop(name, None)
    result.pop(_TESTING_ENV, None)
    result.pop(_AMP_BIN_ENV, None)
    result["HOME"] = os.fspath(_account_home())
    result["PATH"] = "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
    return result


def _child_environment(environ: Mapping[str, str]) -> dict[str, str]:
    """Build the minimal environment visible to the credentialed AMP child."""
    result: dict[str, str] = {
        "HOME": os.fspath(_account_home()),
        "PATH": "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
    }
    for name in ("AMP_BUC_TOKEN", "LANG", "LC_ALL", "LC_CTYPE", "TZ"):
        value = environ.get(name)
        if value:
            result[name] = value
    # Test stubs need only their result channel; production can never opt in
    # because the binary override itself is rejected without this exact flag.
    if environ.get(_TESTING_ENV) == "1":
        for name, value in environ.items():
            if name.startswith("AMP_SAFE_"):
                result[name] = value
    return result


def _authorize(action: str, repo: Path) -> int:
    broker_path = os.environ.get(_HEADLESS_BROKER_ENV, "").strip()
    if broker_path:
        try:
            broker = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
            broker.settimeout(10)
            broker.connect(broker_path)
            if sys.platform == "darwin":
                peer_pid = int.from_bytes(
                    broker.getsockopt(0, 2, 4), sys.byteorder)
            else:
                peercred = getattr(socket, "SO_PEERCRED", None)
                if peercred is None:
                    raise OSError("peer credentials unavailable")
                peer_pid = struct.unpack(
                    "3i", broker.getsockopt(
                        socket.SOL_SOCKET, peercred, 12))[0]
            # A model-controlled child can bind its own socket and override the
            # environment.  The real broker lives in the bridge executor, which
            # must be an ancestor of this wrapper process.
            ancestor = os.getppid()
            peer_is_ancestor = False
            for _depth in range(64):
                if ancestor <= 1:
                    break
                if ancestor == peer_pid:
                    peer_is_ancestor = True
                    break
                completed = subprocess.run(
                    ["/bin/ps", "-o", "ppid=", "-p", str(ancestor)],
                    capture_output=True, text=True, timeout=2, check=False)
                ancestor = int(completed.stdout.strip() or 0)
            if not peer_is_ancestor:
                raise OSError("broker peer is not the headless executor")
            broker.sendall(json.dumps({
                "action": action,
                "repo": os.fspath(repo.resolve(strict=True)),
            }, ensure_ascii=True, sort_keys=True,
                separators=(",", ":")).encode("utf-8"))
            broker.shutdown(socket.SHUT_WR)
            response = broker.recv(32)
        except (OSError, RuntimeError, ValueError) as exc:
            raise AmpSafeError(
                "headless AMP authorization broker unavailable: %s"
                % type(exc).__name__) from exc
        finally:
            try:
                broker.close()
            except (NameError, OSError):
                pass
        if response != b"OK\n":
            print("amp-safe: headless Task fence denied AMP mutation",
                  file=sys.stderr)
            return 2
        # The broker proves the executor-owned remote lease only. Continue
        # through the local authorizer as a second, independent gate for the
        # pre-established CloudSpec baseline and operations/ diff scan.
    worker = Path(__file__).resolve().with_name("jarvis-interactive-worker.py")
    try:
        completed = subprocess.run(
            [sys.executable, "-I", os.fspath(worker), "amp-authorize",
             "--action", action, "--repo", os.fspath(repo)],
            cwd=os.fspath(repo),
            env=_worker_environment(os.environ),
            stdin=subprocess.DEVNULL,
            check=False,
        )
    except OSError as exc:
        raise AmpSafeError("cannot invoke AMP task-fence authorization: %s" % exc) from exc
    return int(completed.returncode)


def _execute_amp(binary: str, argv: Sequence[str], cwd: Path,
                 environ: Mapping[str, str]) -> int:
    command = [binary, *argv, *_FIXED_FLAGS]
    try:
        completed = subprocess.run(
            command,
            cwd=os.fspath(cwd),
            env=dict(environ),
            stdin=subprocess.DEVNULL,
            check=False,
        )
    except OSError as exc:
        raise AmpSafeError("cannot invoke amp: %s" % exc) from exc
    return int(completed.returncode)


def _git_feature_branch(cwd: Path) -> str:
    git = Path("/usr/bin/git")
    if not git.is_file() or not os.access(git, os.X_OK):
        raise AmpSafeError("trusted /usr/bin/git is unavailable")
    try:
        completed = subprocess.run(
            [os.fspath(git), "-C", os.fspath(cwd), "symbolic-ref",
             "--quiet", "--short", "HEAD"],
            cwd=os.fspath(cwd),
            env=_child_environment(os.environ),
            stdin=subprocess.DEVNULL,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=False,
        )
    except OSError as exc:
        raise AmpSafeError("cannot resolve Git branch: %s" % exc) from exc
    if completed.returncode != 0:
        raise AmpSafeError("publish requires an attached Git feature branch")
    return _feature_branch(completed.stdout.strip())


def run(
        argv: Sequence[str], *,
        cwd: Optional[Path] = None,
        environ: Optional[Mapping[str, str]] = None,
        authorizer: Optional[Authorizer] = None,
        branch_resolver: Optional[BranchResolver] = None,
) -> int:
    """Authorize and execute one canonical AMP invocation."""
    invocation = parse_amp_argv(argv)
    working_directory = (Path.cwd() if cwd is None else Path(cwd)).resolve()
    if not working_directory.is_dir():
        raise AmpSafeError("working directory does not exist")
    environment = os.environ if environ is None else environ
    binary = _amp_binary(environment)
    child_environment = _child_environment(environment)

    if invocation.publish_environment == "online" and not invocation.dry_run:
        raise AmpSafeError(
            "online is a production release; real publication requires a human")

    amp_argv = invocation.argv
    if invocation.publish_environment in {"daily", "pre", "online"}:
        resolve_branch = _git_feature_branch if branch_resolver is None else branch_resolver
        branch = _feature_branch(resolve_branch(working_directory))
        amp_argv = (
            invocation.argv[0], invocation.argv[1], "--branch", branch,
            *invocation.argv[2:],
        )

    if invocation.authorize_action is not None:
        authorize = _authorize if authorizer is None else authorizer
        authorization_code = int(authorize(
            invocation.authorize_action, working_directory))
        if authorization_code != 0:
            return authorization_code

    if (invocation.publish_environment in {"daily", "pre"}
            and not invocation.dry_run):
        dry_run_argv = (*amp_argv, "--dry-run")
        dry_run_code = _execute_amp(
            binary, dry_run_argv, working_directory, child_environment)
        if dry_run_code != 0:
            return dry_run_code
        # A dry-run is executable code and may change the checkout. Recheck the
        # task lease and operation diff immediately before the real publish.
        assert invocation.authorize_action is not None
        authorization_code = int(authorize(
            invocation.authorize_action, working_directory))
        if authorization_code != 0:
            return authorization_code
    return _execute_amp(
        binary, amp_argv, working_directory, child_environment)


def main(argv: Optional[Sequence[str]] = None) -> int:
    try:
        return run(sys.argv[1:] if argv is None else argv)
    except AmpSafeError as exc:
        print("amp-safe: %s" % exc, file=sys.stderr)
        return 2


__all__ = [
    "AmpInvocation",
    "AmpSafeError",
    "main",
    "parse_amp_argv",
    "run",
]


if __name__ == "__main__":
    raise SystemExit(main())

#!/usr/bin/env python3
"""Fail-closed safety gate for destructive a1 CR-flow commands.

``app pipeline exit-cr`` removes every CR attached to a publish-flow instance,
so Jarvis must never execute it.  The only supported recovery operation is the
targeted ``app cr quit <cr-id> --pipeline-id <id>`` form.  Before forwarding
that command this module proves that the CR is bound to the currently claimed
Aone work item and to the repository/branch of the calling worktree.

The module is used in two places:

* ``bin/a1id`` routes every a1 invocation through :func:`run_guarded`.
* the interactive-worker PreToolUse fence uses
  :func:`pretool_a1_block_reason` to reject direct ``a1`` bypasses.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import shlex
import subprocess
import sys
from pathlib import Path
from typing import Any, Iterable, Mapping, Optional, Sequence
from urllib.parse import urlsplit


BLOCKED_PIPELINE_QUIT_MESSAGE = (
    "'app pipeline exit-cr' and 'app pipeline quit' are disabled because they "
    "remove every CR in the pipeline instance; use the ownership-checked targeted form "
    "'bin/a1id -- app cr quit <cr-id> --pipeline-id <id>'"
)
DIRECT_QUIT_MESSAGE = (
    "a1 safety: direct 'a1 app cr quit' bypasses ownership verification; use "
    "'bin/a1id -- app cr quit <cr-id> --pipeline-id <id>'"
)

_GLOBAL_NO_VALUE = {
    "--debug", "--no-update-check", "-q", "--quiet", "--verbose",
}
_GLOBAL_WITH_VALUE = {"--config", "-f", "--format"}
_SHELL_OPERATORS = {
    ";", "&&", "||", "|", "&", "<", ">", "<<", ">>", "(", ")", "{", "}",
}


class GuardError(RuntimeError):
    """A command cannot be proven safe and must not reach the real a1."""


def _without_global_flags(argv: Sequence[str]) -> tuple[list[str], list[str]]:
    """Return semantic argv and normalized harmless output/debug flags.

    Cobra accepts inherited global flags before or between command words.  The
    safety classifier therefore removes their space and equals forms before it
    identifies a command path.  ``--config`` is intentionally not returned for
    a targeted quit: identity/config substitution would make the metadata read
    and mutation use different authority contexts.
    """
    semantic: list[str] = []
    forwarded: list[str] = []
    index = 0
    while index < len(argv):
        token = str(argv[index])
        if token in _GLOBAL_NO_VALUE:
            forwarded.append(token)
            index += 1
            continue
        if (any(token.startswith(flag + "=") for flag in (
                "--debug", "--no-update-check", "--quiet", "--verbose"))
                or token.startswith("-q=")):
            forwarded.append(token)
            index += 1
            continue
        if token in _GLOBAL_WITH_VALUE:
            if index + 1 >= len(argv):
                # Keep malformed argv semantic so a destructive command cannot
                # hide behind a missing global-flag value.
                semantic.append(token)
                index += 1
                continue
            value = str(argv[index + 1])
            if token == "--config":
                forwarded.extend([token, value])
            else:
                forwarded.extend([token, value])
            index += 2
            continue
        matched = False
        for flag in ("--config", "--format"):
            prefix = flag + "="
            if token.startswith(prefix):
                forwarded.append(token)
                matched = True
                break
        if matched:
            index += 1
            continue
        if token.startswith("-f=") or (token.startswith("-f") and len(token) > 2):
            forwarded.append(token)
            index += 1
            continue
        semantic.append(token)
        index += 1
    return semantic, forwarded


def _command_kind(argv: Sequence[str]) -> str:
    semantic, _forwarded = _without_global_flags(argv)
    if semantic[:3] == ["app", "pipeline", "exit-cr"]:
        return "pipeline-bulk-quit"
    if semantic[:3] == ["app", "pipeline", "quit"]:
        return "pipeline-bulk-quit"
    if semantic[:3] == ["app", "cr", "quit"]:
        return "cr-quit"
    return "other"


def _a1id_payload(argv: Sequence[str]) -> Sequence[str]:
    if not argv:
        return ()
    if argv[0] == "--":
        return argv[1:]
    if argv[0] == "as" and len(argv) >= 2:
        remaining = list(argv[2:])
        if remaining[:1] == ["--"]:
            remaining = remaining[1:]
        return remaining
    return ()


def _shell_tokens(command: str) -> list[str]:
    lexer = shlex.shlex(command, posix=True, punctuation_chars=";&|<>(){}")
    lexer.whitespace_split = True
    lexer.commenters = ""
    return list(lexer)


def _command_slices(tokens: Sequence[str]) -> Iterable[Sequence[str]]:
    start = 0
    for index, token in enumerate(tokens):
        if token in _SHELL_OPERATORS or re.fullmatch(r"[;&|<>]+", token):
            if start < index:
                yield tokens[start:index]
            start = index + 1
    if start < len(tokens):
        yield tokens[start:]


def _basename(token: str) -> str:
    return token.rsplit("/", 1)[-1]


def _execution_index(invocation: Sequence[str]) -> Optional[int]:
    """Locate the executable token without treating data arguments as commands."""
    index = 0
    # A shell compound-list can leave its control word at the start of the
    # same simple-command slice (for example ``if a1 ...`` or
    # ``then a1 ...``).  Negation is likewise a shell prefix rather than the
    # executable.  Skip only leading control words; arguments later in the
    # command remain data and cannot trigger the classifier.
    while index < len(invocation) and invocation[index] in {
            "if", "then", "elif", "else", "while", "until", "do", "!",
    }:
        index += 1
    while index < len(invocation) and re.fullmatch(
            r"[A-Za-z_][A-Za-z0-9_]*=.*", invocation[index]):
        index += 1
    if index >= len(invocation):
        return None
    if _basename(invocation[index]) == "env":
        index += 1
        while index < len(invocation):
            token = invocation[index]
            if re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*=.*", token):
                index += 1
                continue
            if token in {"-i", "--ignore-environment", "-0", "--null"}:
                index += 1
                continue
            if token in {"-u", "--unset", "-C", "--chdir", "-S", "--split-string"}:
                index += 2
                continue
            if token.startswith("--unset=") or token.startswith("--chdir="):
                index += 1
                continue
            break
    if index < len(invocation) and _basename(invocation[index]) == "command":
        index += 1
        # command -v/-V is an inspection, not execution.
        if index < len(invocation) and invocation[index] in {"-v", "-V"}:
            return None
        while index < len(invocation) and invocation[index] == "-p":
            index += 1
    if index < len(invocation) and _basename(invocation[index]) == "exec":
        index += 1
        while index < len(invocation):
            if invocation[index] in {"-c", "-l"}:
                index += 1
                continue
            if invocation[index] == "-a":
                index += 2
                continue
            break
    # ``time`` and ``nohup`` execute the following command.  Treat their
    # option prefixes as wrappers so direct a1 cannot hide behind them.
    while index < len(invocation) and _basename(invocation[index]) in {
            "time", "nohup",
    }:
        wrapper = _basename(invocation[index])
        index += 1
        while index < len(invocation) and invocation[index].startswith("-"):
            option = invocation[index]
            index += 1
            if wrapper == "time" and option in {
                    "-f", "-o", "--format", "--output",
            }:
                index += 1
        while index < len(invocation) and re.fullmatch(
                r"[A-Za-z_][A-Za-z0-9_]*=.*", invocation[index]):
            index += 1
    # Wrappers may be nested in either order (``time env a1`` or
    # ``nohup command a1``).  Re-run the prefix parser on the remaining tail
    # instead of granting one particular ordering a bypass.
    if index < len(invocation) and _basename(invocation[index]) in {
            "env", "command", "exec", "time", "nohup",
    }:
        nested = _execution_index(invocation[index:])
        return None if nested is None else index + nested
    return index if index < len(invocation) else None


def _pretool_reason_from_command(command: str, depth: int = 0) -> Optional[str]:
    if depth > 2:
        return None
    try:
        tokens = _shell_tokens(command)
    except ValueError:
        # The Worker fence handles malformed shell input independently.  This
        # classifier does not grant a permit, so returning None is safe here.
        return None
    for invocation in _command_slices(tokens):
        exec_index = _execution_index(invocation)
        if exec_index is None:
            continue
        executable = _basename(invocation[exec_index])
        exec_args = invocation[exec_index + 1:]
        if executable == "a1":
            kind = _command_kind(exec_args)
            if kind == "pipeline-bulk-quit":
                return BLOCKED_PIPELINE_QUIT_MESSAGE
            if kind == "cr-quit":
                return DIRECT_QUIT_MESSAGE
        elif executable == "a1id":
            kind = _command_kind(_a1id_payload(exec_args))
            if kind == "pipeline-bulk-quit":
                return BLOCKED_PIPELINE_QUIT_MESSAGE
        elif executable in {"bash", "sh", "zsh"}:
            # ``bash bin/a1id -- ...`` is a common explicit-wrapper form.
            script_index = exec_index + 1
            while (script_index < len(invocation)
                   and invocation[script_index].startswith("-")
                   and invocation[script_index] not in {"-c", "-lc", "-fc"}):
                script_index += 1
            if (script_index < len(invocation)
                    and _basename(invocation[script_index]) == "a1id"):
                kind = _command_kind(_a1id_payload(invocation[script_index + 1:]))
                if kind == "pipeline-bulk-quit":
                    return BLOCKED_PIPELINE_QUIT_MESSAGE

        # Cover common ``bash -c/-lc/-fc 'a1 ...'`` wrappers.  shlex
        # deliberately leaves the quoted script as one token.
        if executable in {"bash", "sh", "zsh"}:
            for index in range(exec_index + 1, len(invocation) - 1):
                if invocation[index] in {"-c", "-lc", "-fc"}:
                    nested = _pretool_reason_from_command(
                        str(invocation[index + 1]), depth + 1)
                    if nested:
                        return nested
    return None


def pretool_a1_block_reason(event: Mapping[str, Any]) -> Optional[str]:
    """Return a user-facing block reason for unsafe Bash a1 invocations."""
    tool_name = str(event.get("tool_name") or "").strip().lower()
    if not (tool_name == "bash" or tool_name == "exec_command"
            or tool_name.endswith("__exec_command")
            or tool_name.endswith(".exec_command")
            or tool_name.endswith(":exec_command")):
        return None
    tool_input = event.get("tool_input")
    if not isinstance(tool_input, Mapping):
        return None
    command = tool_input.get("command")
    if command is None:
        command = tool_input.get("cmd")
    if not isinstance(command, str) or not command.strip():
        return None
    return _pretool_reason_from_command(command)


def _wrapper_quit_argv_in_command(command: str, depth: int = 0) -> Optional[list[str]]:
    if depth > 2:
        return None
    try:
        tokens = _shell_tokens(command)
    except ValueError:
        return None
    for invocation in _command_slices(tokens):
        exec_index = _execution_index(invocation)
        if exec_index is None:
            continue
        executable = _basename(invocation[exec_index])
        args = invocation[exec_index + 1:]
        if executable == "a1id" and _command_kind(_a1id_payload(args)) == "cr-quit":
            return list(_a1id_payload(args))
        if executable in {"bash", "sh", "zsh"}:
            script_index = exec_index + 1
            while (script_index < len(invocation)
                   and invocation[script_index].startswith("-")
                   and invocation[script_index] not in {"-c", "-lc", "-fc"}):
                script_index += 1
            if (script_index < len(invocation)
                    and _basename(invocation[script_index]) == "a1id"
                    and _command_kind(_a1id_payload(
                        invocation[script_index + 1:])) == "cr-quit"):
                return list(_a1id_payload(invocation[script_index + 1:]))
            for index in range(exec_index + 1, len(invocation) - 1):
                if invocation[index] in {"-c", "-lc", "-fc"}:
                    nested = _wrapper_quit_argv_in_command(
                        invocation[index + 1], depth + 1)
                    if nested is not None:
                        return nested
    return None


def pretool_has_targeted_a1id_quit(event: Mapping[str, Any]) -> bool:
    """Whether a shell tool asks bin/a1id to perform a targeted CR quit."""
    tool_input = event.get("tool_input")
    if not isinstance(tool_input, Mapping):
        return False
    command = tool_input.get("command")
    if command is None:
        command = tool_input.get("cmd")
    return (isinstance(command, str)
            and _wrapper_quit_argv_in_command(command) is not None)


def pretool_targeted_a1id_quit(event: Mapping[str, Any]) -> Optional[tuple[str, str]]:
    """Return the positive ``(cr_id, pipeline_id)`` requested via bin/a1id."""
    tool_input = event.get("tool_input")
    if not isinstance(tool_input, Mapping):
        return None
    command = tool_input.get("command")
    if command is None:
        command = tool_input.get("cmd")
    if not isinstance(command, str):
        return None
    argv = _wrapper_quit_argv_in_command(command)
    if argv is None:
        return None
    try:
        cr_id, pipeline_id, _app, _globals = _parse_targeted_quit(argv)
    except GuardError:
        return None
    return cr_id, pipeline_id


def _worker_authority(phase: str, cr_id: str, pipeline_id: str,
                      token: str = "") -> tuple[dict[str, str], str]:
    """Read current Aone ownership through the trusted Worker manager.

    The manager reuses its locked host-incarnation, local-state and lease-proof
    checks.  Reading the JSON state file directly here would let stale local
    lineage authorize a mutation after the server fence was lost.
    """
    manager = Path(__file__).resolve().with_name("jarvis-interactive-worker.py")
    try:
        command = [
            sys.executable, "-I", str(manager), "current-authority",
            phase, cr_id, pipeline_id,
        ]
        if token:
            command.append(token)
        result = subprocess.run(
            command,
            capture_output=True, text=True, timeout=10, check=False)
    except (OSError, subprocess.SubprocessError) as exc:
        raise GuardError("cannot verify the current Worker authority") from exc
    if result.returncode != 0:
        detail = result.stderr.strip().splitlines()[-1] if result.stderr.strip() else "unavailable"
        raise GuardError("current Worker authority is invalid: %s" % detail)
    try:
        parsed = json.loads(result.stdout)
    except (TypeError, ValueError) as exc:
        raise GuardError("current Worker authority returned invalid JSON") from exc
    if not isinstance(parsed, Mapping):
        raise GuardError("current Worker authority returned no context")
    aone_id = str(parsed.get("aoneId") or "").strip()
    if not aone_id.isdigit() or int(aone_id) <= 0:
        raise GuardError("current Worker authority has no positive Aone id")
    permit_token = str(parsed.get("permitToken") or "").strip()
    if phase == "begin" and not permit_token:
        raise GuardError("current Worker authority returned no CR-quit permit")
    assignment_epoch = str(parsed.get("assignmentEpoch") or "").strip()
    worker_digest = str(parsed.get("workerKeyDigest") or "").strip()
    if not assignment_epoch or not re.fullmatch(r"[0-9a-f]{16}", worker_digest):
        raise GuardError("current Worker authority returned an incomplete fingerprint")
    return {
        "aoneId": aone_id,
        "assignmentEpoch": assignment_epoch,
        "workerKeyDigest": worker_digest,
    }, permit_token


def _git_output(*args: str) -> str:
    try:
        result = subprocess.run(
            ["git", *args], capture_output=True, text=True, timeout=5,
            check=False)
    except (OSError, subprocess.SubprocessError) as exc:
        raise GuardError("cannot inspect the current git worktree") from exc
    if result.returncode != 0:
        raise GuardError("targeted CR quit must run inside the CR git worktree")
    return result.stdout.strip()


def _normalize_repo_url(raw: str) -> str:
    value = str(raw or "").strip().split(" ", 1)[0]
    if not value:
        return ""
    if re.match(r"^[^/@]+@[^/:]+:", value):
        user_host, path = value.split(":", 1)
        host = user_host.rsplit("@", 1)[-1]
    else:
        parsed = urlsplit(value)
        if parsed.scheme and parsed.netloc:
            host = parsed.hostname or ""
            path = parsed.path
        else:
            return value.removesuffix(".git").rstrip("/")
    return (host.lower() + "/" + path.lstrip("/").removesuffix(".git").rstrip("/"))


def _current_repo_branch() -> tuple[str, str]:
    branch = _git_output("branch", "--show-current")
    if not branch or branch in {"main", "master"}:
        raise GuardError("targeted CR quit requires a non-default worktree branch")
    origin = _normalize_repo_url(_git_output("remote", "get-url", "origin"))
    if not origin:
        raise GuardError("current worktree has no verifiable origin remote")
    return origin, branch


def _parse_targeted_quit(argv: Sequence[str]) -> tuple[str, str, Optional[str], list[str]]:
    semantic, global_flags = _without_global_flags(argv)
    if semantic[:3] != ["app", "cr", "quit"]:
        raise GuardError("internal error: command is not a targeted CR quit")
    if any(flag == "--config" or flag.startswith("--config=") for flag in global_flags):
        raise GuardError("targeted CR quit does not allow --config identity substitution")

    cr_ids: list[str] = []
    pipeline_ids: list[str] = []
    apps: list[str] = []
    remaining = semantic[3:]
    index = 0
    while index < len(remaining):
        token = remaining[index]
        if token in {"--pipeline-id", "--app"}:
            if index + 1 >= len(remaining):
                raise GuardError("targeted CR quit flag %s requires a value" % token)
            target = pipeline_ids if token == "--pipeline-id" else apps
            target.append(remaining[index + 1])
            index += 2
            continue
        if token.startswith("--pipeline-id="):
            pipeline_ids.append(token.split("=", 1)[1])
        elif token.startswith("--app="):
            apps.append(token.split("=", 1)[1])
        elif token.startswith("-"):
            raise GuardError("unsupported targeted CR quit flag: %s" % token)
        else:
            cr_ids.append(token)
        index += 1
    if (len(cr_ids) != 1 or not cr_ids[0].isdigit()
            or int(cr_ids[0]) <= 0):
        raise GuardError("targeted CR quit requires exactly one positive numeric <cr-id>")
    if (len(pipeline_ids) != 1 or not pipeline_ids[0].isdigit()
            or int(pipeline_ids[0]) <= 0):
        raise GuardError("targeted CR quit requires exactly one positive numeric --pipeline-id")
    if len(apps) > 1 or (apps and not re.fullmatch(r"[A-Za-z0-9_.-]+", apps[0])):
        raise GuardError("targeted CR quit accepts at most one simple --app value")
    # --config was excluded above; the remaining global flags are output/debug
    # controls and are safe to preserve on the final targeted mutation.
    return cr_ids[0], pipeline_ids[0], (apps[0] if apps else None), global_flags


def _read_cr(a1_bin: str, cr_id: str) -> Mapping[str, Any]:
    try:
        result = subprocess.run(
            [a1_bin, "app", "cr", "get", cr_id, "--workitems", "--format", "json"],
            capture_output=True, text=True, timeout=30, check=False)
    except (OSError, subprocess.SubprocessError) as exc:
        raise GuardError("cannot read CR metadata for ownership verification") from exc
    if result.returncode != 0:
        raise GuardError("a1 app cr get failed; targeted quit is blocked")
    try:
        parsed = json.loads(result.stdout)
    except (TypeError, ValueError) as exc:
        raise GuardError("a1 app cr get returned invalid JSON; targeted quit is blocked") from exc
    if not isinstance(parsed, Mapping):
        raise GuardError("a1 app cr get returned no CR object; targeted quit is blocked")
    return parsed


def _read_a1_json(a1_bin: str, argv: Sequence[str], label: str) -> Mapping[str, Any]:
    try:
        result = subprocess.run(
            [a1_bin, *argv], capture_output=True, text=True, timeout=30,
            check=False)
    except (OSError, subprocess.SubprocessError) as exc:
        raise GuardError("cannot read %s for ownership verification" % label) from exc
    if result.returncode != 0:
        raise GuardError("a1 %s failed; targeted quit is blocked" % label)
    try:
        parsed = json.loads(result.stdout)
    except (TypeError, ValueError) as exc:
        raise GuardError("a1 %s returned invalid JSON; targeted quit is blocked" % label) from exc
    if not isinstance(parsed, Mapping):
        raise GuardError("a1 %s returned no object; targeted quit is blocked" % label)
    return parsed


def _item_repo_branch(item: Mapping[str, Any]) -> tuple[str, str]:
    content = item.get("crItemContent")
    if not isinstance(content, Mapping):
        raw = item.get("content")
        try:
            content = json.loads(raw) if isinstance(raw, str) else {}
        except (TypeError, ValueError):
            content = {}
    branch = str(content.get("branchName") or "").strip()
    repo = _normalize_repo_url(str(
        content.get("trunkUrl") or content.get("branchUrl") or ""))
    return repo, branch


def _verify_cr_ownership(cr: Mapping[str, Any], *, cr_id: str,
                         requested_app: Optional[str], claimed_aone_id: str,
                         origin: str, branch: str) -> str:
    actual_cr_id = str(cr.get("crId") or cr.get("id") or "").strip()
    if actual_cr_id != cr_id:
        raise GuardError("CR metadata id does not match the requested CR")
    workitems = cr.get("workItemIds")
    if (not isinstance(workitems, list)
            or claimed_aone_id not in {str(value) for value in workitems}):
        raise GuardError(
            "CR is not bound to the currently claimed Aone %s" % claimed_aone_id)

    items = cr.get("crItems")
    matched = False
    if isinstance(items, list):
        for item in items:
            if not isinstance(item, Mapping):
                continue
            item_repo, item_branch = _item_repo_branch(item)
            if item_repo == origin and item_branch == branch:
                matched = True
                break
    if not matched:
        raise GuardError(
            "CR repository/branch does not match the current worktree (%s)" % branch)

    app_id = str(cr.get("appId") or "").strip()
    app_name = str(cr.get("appName") or "").strip()
    if not app_id.isdigit() or int(app_id) <= 0:
        raise GuardError("CR metadata has no positive numeric appId")
    if requested_app and requested_app not in {app_id, app_name}:
        raise GuardError("requested --app does not match the CR application")
    return app_id


def _positive_id(value: Any, label: str) -> str:
    text = str(value or "").strip()
    if not text.isdigit() or int(text) <= 0:
        raise GuardError("%s has no positive numeric id" % label)
    return text


def _verify_pipeline_membership(a1_bin: str, *, app_id: str,
                                pipeline_id: str, cr_id: str) -> None:
    """Prove the target CR is in the latest instance of the requested flow."""
    status = _read_a1_json(a1_bin, [
        "app", "pipeline", "status", "--app", app_id,
        "--pipeline-id", pipeline_id, "--format", "json",
    ], "app pipeline status")
    if _positive_id(status.get("appId"), "pipeline status appId") != app_id:
        raise GuardError("pipeline status appId does not match the CR")
    if _positive_id(status.get("pipelineId"), "pipeline status pipelineId") != pipeline_id:
        raise GuardError("pipeline status pipelineId does not match the request")
    instance_id = _positive_id(
        status.get("pipelineInstanceId"), "pipeline status pipelineInstanceId")

    branches = _read_a1_json(a1_bin, [
        "app", "pipeline", "branch", "--app", app_id,
        "--instance-id", instance_id, "--format", "json",
    ], "app pipeline branch")
    if _positive_id(branches.get("appId"), "pipeline branch appId") != app_id:
        raise GuardError("pipeline branch appId does not match the CR")
    if _positive_id(branches.get("pipelineId"), "pipeline branch pipelineId") != pipeline_id:
        raise GuardError("pipeline branch pipelineId does not match the request")
    if (_positive_id(branches.get("pipelineInstanceId"),
                     "pipeline branch pipelineInstanceId") != instance_id):
        raise GuardError("pipeline instance changed during membership verification")
    change_requests = branches.get("changeRequests")
    ids = {
        str(item.get("crId") or item.get("id") or "").strip()
        for item in change_requests
        if isinstance(item, Mapping)
    } if isinstance(change_requests, list) else set()
    if cr_id not in ids:
        raise GuardError("target CR is not attached to the latest requested pipeline instance")


def run_guarded(a1_bin: str, argv: Sequence[str]) -> None:
    kind = _command_kind(argv)
    if kind == "pipeline-bulk-quit":
        raise GuardError(BLOCKED_PIPELINE_QUIT_MESSAGE)
    if kind != "cr-quit":
        os.execvpe(a1_bin, [a1_bin, *argv], os.environ.copy())

    cr_id, pipeline_id, requested_app, global_flags = _parse_targeted_quit(argv)
    authority, permit_token = _worker_authority(
        "begin", cr_id, pipeline_id)
    claimed_aone_id = authority["aoneId"]
    origin, branch = _current_repo_branch()
    cr = _read_cr(a1_bin, cr_id)
    app_id = _verify_cr_ownership(
        cr, cr_id=cr_id, requested_app=requested_app,
        claimed_aone_id=claimed_aone_id, origin=origin, branch=branch)
    _verify_pipeline_membership(
        a1_bin, app_id=app_id, pipeline_id=pipeline_id, cr_id=cr_id)
    # Close the metadata/network TOCTOU window: the exact task/session/fence
    # authority must still be current immediately before the mutation.
    confirmed_authority, _unused = _worker_authority(
        "confirm", cr_id, pipeline_id, permit_token)
    if confirmed_authority != authority:
        raise GuardError("Worker task/session/fence authority changed during CR quit verification")
    canonical = [
        a1_bin, "app", "cr", "quit", cr_id,
        "--pipeline-id", pipeline_id, "--app", app_id,
        *global_flags,
    ]
    os.execvpe(a1_bin, canonical, os.environ.copy())


def _parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--a1-bin")
    parser.add_argument("--check-a1id-argv", action="store_true")
    parser.add_argument("a1_args", nargs=argparse.REMAINDER)
    return parser


def main(argv: Optional[Sequence[str]] = None) -> int:
    args = _parser().parse_args(argv)
    raw_args = list(args.a1_args)
    if args.check_a1id_argv:
        if _command_kind(_a1id_payload(raw_args)) == "pipeline-bulk-quit":
            print("a1 safety: %s" % BLOCKED_PIPELINE_QUIT_MESSAGE, file=sys.stderr)
            return 2
        return 0
    a1_args = raw_args
    if a1_args[:1] == ["--"]:
        a1_args = a1_args[1:]
    if not args.a1_bin:
        print("a1 safety: --a1-bin is required", file=sys.stderr)
        return 64
    if not a1_args:
        print("a1 safety: missing a1 command", file=sys.stderr)
        return 64
    try:
        run_guarded(args.a1_bin, a1_args)
    except GuardError as exc:
        print("a1 safety: %s" % exc, file=sys.stderr)
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

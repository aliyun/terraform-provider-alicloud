#!/usr/bin/env python3
"""Fail-closed safety gate for destructive a1 CR-flow commands.

Jarvis must never execute ``app pipeline exit-cr``, ``app pipeline quit`` or
``app cr quit``.  All three mutate CR/branch membership in a publish flow;
conflicts and withdrawal requests must stop and be handed to a human instead.

The module is used in two places:

* ``bin/a1id`` routes every a1 invocation through :func:`run_guarded`.
* the interactive-worker PreToolUse fence uses
  :func:`pretool_a1_block_reason` to reject direct ``a1`` bypasses.
"""

from __future__ import annotations

import argparse
import os
import re
import shlex
import sys
from typing import Any, Iterable, Mapping, Optional, Sequence


BLOCKED_CR_EXIT_MESSAGE = (
    "'app pipeline exit-cr', 'app pipeline quit', and 'app cr quit' are "
    "permanently disabled for Jarvis; stop and ask a human to handle the CR/branch"
)
DIRECT_QUIT_MESSAGE = (
    "a1 safety: 'app cr quit' is permanently disabled for Jarvis; stop and "
    "ask a human to handle the CR/branch"
)
ACUBE_BUILD_TASK_MESSAGE = (
    "Acube createBuildTaskV2 is permanently disabled for Jarvis; "
    "Terraform work must continue on the source Aone Task"
)
TERRAFORM_SOURCE_AONE_WRITE_MESSAGE = (
    "Terraform 528766 source Task model is Aone read-only; "
    "the bridge executor owns source-ticket bookend and downstream workitems "
    "must not be created or related"
)

_GLOBAL_NO_VALUE = {
    "--debug", "--no-update-check", "-q", "--quiet", "--verbose",
}
_GLOBAL_WITH_VALUE = {"--config", "-f", "--format"}
_SHELL_OPERATORS = {
    ";", "&&", "||", "|", "&", "<", ">", "<<", ">>", "(", ")", "{", "}",
}
_SIMPLE_ASSIGNMENT = re.compile(r"([A-Za-z_][A-Za-z0-9_]*)=(.*)", re.DOTALL)
_SIMPLE_VAR_REF = re.compile(
    r"\$(?:\{([A-Za-z_][A-Za-z0-9_]*)\}|([A-Za-z_][A-Za-z0-9_]*))")
_HTTP_CLIENTS = frozenset({"curl", "wget", "http", "https"})
_SCRIPT_CLIENT = re.compile(
    r"^(?:python(?:\d+(?:\.\d+)*)?|pypy(?:\d+)?|node(?:js)?|ruby|perl|php)$")
_SCRIPT_NETWORK_HINT = re.compile(
    r"(?:requests?\s*\.|httpx\s*\.|aiohttp\s*\.|urllib(?:3|\.request)|"
    r"http\.client|fetch\s*\(|axios\s*[\.(]|https?\s*\.\s*request|"
    r"superagent\s*[\.(]|child_process|subprocess|curl\b|wget\b)",
    re.IGNORECASE,
)


class GuardError(RuntimeError):
    """A command cannot be proven safe and must not reach the real a1."""


def _without_global_flags(argv: Sequence[str]) -> tuple[list[str], list[str]]:
    """Return semantic argv and normalized harmless output/debug flags.

    Cobra accepts inherited global flags before or between command words.  The
    safety classifier therefore removes their space and equals forms before it
    identifies a command path regardless of harmless global flag placement.
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
        return "blocked-quit"
    if semantic[:3] == ["app", "pipeline", "quit"]:
        return "blocked-quit"
    if semantic[:3] == ["app", "cr", "quit"]:
        return "blocked-quit"
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


def _aone_workitem_write(argv: Sequence[str]) -> bool:
    semantic, _forwarded = _without_global_flags(argv)
    if semantic[:2] != ["project", "workitem"] or len(semantic) < 3:
        return False
    operation = semantic[2]
    if operation in {"get", "list", "activity"}:
        return False
    if operation == "comment" and semantic[3:4] == ["list"]:
        return False
    if operation in {"relation", "attachment", "field"}:
        return tuple(semantic[3:4]) not in {
            ("list",), ("download",), ("options",),
        }
    return True


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


def _expand_simple_vars(value: str, variables: Mapping[str, str]) -> str:
    """Expand bounded ``$name``/``${name}`` references without executing shell.

    This intentionally supports only literal same-command assignments.  It
    never evaluates command substitutions, arithmetic, parameter operators, or
    files, keeping the PreToolUse classifier deterministic and side-effect free.
    """
    result = str(value)
    for _unused in range(8):
        expanded = _SIMPLE_VAR_REF.sub(
            lambda match: variables.get(
                match.group(1) or match.group(2), match.group(0)),
            result,
        )
        if expanded == result:
            break
        result = expanded
    return result


def _remember_simple_assignments(
    invocation: Sequence[str],
    variables: dict[str, str],
) -> None:
    """Remember literal assignment prefixes/standalone assignment commands."""
    for token in invocation:
        match = _SIMPLE_ASSIGNMENT.fullmatch(str(token))
        if match is None:
            break
        value = match.group(2)
        # Complex shell evaluation is deliberately not approximated.  Leaving
        # it unresolved is safer than inventing a value that can cross-match an
        # unrelated later command.
        if "$(" in value or "`" in value:
            continue
        variables[match.group(1)] = _expand_simple_vars(value, variables)


def _is_acube_build_invocation(
    invocation: Sequence[str],
    variables: Mapping[str, str],
) -> bool:
    """Recognize an actual network/script client targeting createBuildTaskV2."""
    exec_index = _execution_index(invocation)
    if exec_index is None:
        return False
    executable = _basename(
        _expand_simple_vars(str(invocation[exec_index]), variables))
    expanded_args = [
        _expand_simple_vars(str(token), variables)
        for token in invocation[exec_index + 1:]
    ]
    joined = " ".join(expanded_args)
    if "createBuildTaskV2" not in joined:
        return False
    if executable in _HTTP_CLIENTS:
        return True
    # Inline script clients are powerful enough to make the request without a
    # curl/wget child.  Require a network/process API hint in that same
    # invocation so ``python -c 'print("createBuildTaskV2")'`` remains an
    # ordinary read-only audit.
    return bool(
        _SCRIPT_CLIENT.fullmatch(executable)
        and _SCRIPT_NETWORK_HINT.search(joined)
    )


def _pretool_reason_from_command(command: str, depth: int = 0) -> Optional[str]:
    if depth > 2:
        return None
    try:
        tokens = _shell_tokens(command)
    except ValueError:
        # The Worker fence handles malformed shell input independently.  This
        # classifier does not grant a permit, so returning None is safe here.
        return None
    invocations = list(_command_slices(tokens))
    # Track only literal assignments across this compound command.  The target
    # marker must resolve inside the same actual client invocation; a later
    # ``printf createBuildTaskV2`` cannot taint an unrelated earlier curl.
    variables: dict[str, str] = {}
    for invocation in invocations:
        _remember_simple_assignments(invocation, variables)
        if _is_acube_build_invocation(invocation, variables):
            return ACUBE_BUILD_TASK_MESSAGE
        exec_index = _execution_index(invocation)
        if exec_index is None:
            continue
        executable = _basename(invocation[exec_index])
        exec_args = invocation[exec_index + 1:]
        if executable == "a1":
            kind = _command_kind(exec_args)
            if kind == "blocked-quit":
                return (DIRECT_QUIT_MESSAGE if
                        _without_global_flags(exec_args)[0][:3] == ["app", "cr", "quit"]
                        else BLOCKED_CR_EXIT_MESSAGE)
        elif executable == "a1id":
            kind = _command_kind(_a1id_payload(exec_args))
            if kind == "blocked-quit":
                return BLOCKED_CR_EXIT_MESSAGE
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
                if kind == "blocked-quit":
                    return BLOCKED_CR_EXIT_MESSAGE

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


def _pretool_aone_write_from_command(command: str, depth: int = 0) -> bool:
    if depth > 2:
        return False
    try:
        tokens = _shell_tokens(command)
    except ValueError:
        return False
    for invocation in _command_slices(tokens):
        exec_index = _execution_index(invocation)
        if exec_index is None:
            continue
        executable = _basename(invocation[exec_index])
        args = invocation[exec_index + 1:]
        if executable == "a1" and _aone_workitem_write(args):
            return True
        if executable == "a1id" and _aone_workitem_write(_a1id_payload(args)):
            return True
        if executable in {"bash", "sh", "zsh"}:
            script_index = exec_index + 1
            while (script_index < len(invocation)
                   and invocation[script_index].startswith("-")
                   and invocation[script_index] not in {"-c", "-lc", "-fc"}):
                script_index += 1
            if (script_index < len(invocation)
                    and _basename(invocation[script_index]) == "a1id"
                    and _aone_workitem_write(
                        _a1id_payload(invocation[script_index + 1:]))):
                return True
            for index in range(exec_index + 1, len(invocation) - 1):
                if (invocation[index] in {"-c", "-lc", "-fc"}
                        and _pretool_aone_write_from_command(
                            str(invocation[index + 1]), depth + 1)):
                    return True
    return False


def pretool_aone_write_block_reason(
    event: Mapping[str, Any],
) -> Optional[str]:
    """Return a block reason for direct a1/a1id workitem mutations."""
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
    if (isinstance(command, str)
            and _pretool_aone_write_from_command(command)):
        return TERRAFORM_SOURCE_AONE_WRITE_MESSAGE
    return None


def run_guarded(a1_bin: str, argv: Sequence[str]) -> None:
    kind = _command_kind(argv)
    if kind == "blocked-quit":
        raise GuardError(BLOCKED_CR_EXIT_MESSAGE)
    os.execvpe(a1_bin, [a1_bin, *argv], os.environ.copy())


def _parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--a1-bin")
    parser.add_argument("--check-a1id-argv", action="store_true")
    parser.add_argument("--check-aone-write-argv", action="store_true")
    parser.add_argument("--check-pretool-command")
    parser.add_argument("a1_args", nargs=argparse.REMAINDER)
    return parser


def main(argv: Optional[Sequence[str]] = None) -> int:
    args = _parser().parse_args(argv)
    raw_args = list(args.a1_args)
    if args.check_pretool_command is not None:
        reason = _pretool_reason_from_command(args.check_pretool_command)
        if reason:
            print("a1 safety: %s" % reason, file=sys.stderr)
            return 2
        return 0
    if args.check_aone_write_argv:
        raw_args = raw_args[1:] if raw_args[:1] == ["--"] else raw_args
        return 2 if _aone_workitem_write(raw_args) else 0
    if args.check_a1id_argv:
        if _command_kind(_a1id_payload(raw_args)) == "blocked-quit":
            print("a1 safety: %s" % BLOCKED_CR_EXIT_MESSAGE, file=sys.stderr)
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

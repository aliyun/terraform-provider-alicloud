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


def run_guarded(a1_bin: str, argv: Sequence[str]) -> None:
    kind = _command_kind(argv)
    if kind == "blocked-quit":
        raise GuardError(BLOCKED_CR_EXIT_MESSAGE)
    os.execvpe(a1_bin, [a1_bin, *argv], os.environ.copy())


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

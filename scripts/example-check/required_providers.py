#!/usr/bin/env python3
"""HCL-aware provider requirement handling for local example checks.

``scripts/local-example-check.sh`` extracts each documentation example into
``<example_dir>/main.tf`` and used to unconditionally write a second
``terraform.tf`` declaring ``required_providers``. When the example already
declares ``terraform.required_providers`` the two files then form one module
with two required provider configurations and configuration loading fails
with ``Duplicate required providers configuration``.

This module parses the example HCL properly (plain strings, comments and
heredocs are never blocks) and decides:

* no ``required_providers`` block found          -> write ``terraform.tf``
  with the alicloud requirement (previous behaviour);
* one block that already declares alicloud       -> keep the example's own
  declaration (source, version and provider alias untouched);
* one block declaring other providers only       -> complete that very block
  in place with the alicloud source, no second block;
* two or more real blocks                        -> leave everything as-is so
  terraform still reports the duplicate configuration error.

Only the temporary working copy under ``--example-dir`` is modified; the
documentation source itself is never rewritten.
"""

import argparse
import re
import sys
from pathlib import Path

ALICLOUD_SOURCE = "registry.terraform.io/aliyun/alicloud"

# Exactly the file local-example-check.sh used to write unconditionally.
INJECTED_TERRAFORM_TF = """terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
    }
  }
}

"""

IDENT_RE = re.compile(r"[A-Za-z_][A-Za-z0-9_-]*")
HEREDOC_RE = re.compile(r"<<-?([A-Za-z_][A-Za-z0-9_-]*)")
# Structural tokens scanned on the masked text: identifiers, (emptied)
# strings, braces and the attribute assignment sign.
STRUCT_RE = re.compile(r'[A-Za-z_][A-Za-z0-9_-]*|"[^"\n]*"|[{}]|=')


def _scan_quoted(text, start):
    """Return the index of the closing quote of the string at ``start``.

    Interpolations ``${ ... }`` and template directives ``%{ ... }`` may
    contain nested quotes, braces and newlines, so they are tracked with
    a depth counter instead of ending the string early: the expression
    inside an open interpolation may legally span several lines, and
    object literals or for-expressions inside it contribute their own
    balanced ``{ ... }`` pairs to that depth. Escaped introducers
    (``$${`` and ``%%{``) are literal text and open nothing. Outside an
    interpolation a raw newline ends the scan (HCL quoted strings cannot
    span lines; heredocs are handled elsewhere).
    """
    i = start + 1
    n = len(text)
    interp = 0
    while i < n:
        ch = text[i]
        if ch == "\\":
            i += 2
            continue
        if ch == '"':
            if interp:
                i = _scan_quoted(text, i) + 1
                continue
            return i
        if ch == "\n" and not interp:
            return i
        if ch in "$%":
            triple = text[i:i + 3]
            if triple in ("$${", "%%{"):
                # Escaped template introducer: literal "${" or "%{".
                i += 3
                continue
            if text[i:i + 2] in ("${", "%{"):
                # An interpolation or template directive opens an
                # expression context that may legally span lines.
                interp += 1
                i += 2
                continue
        if ch == "{" and interp:
            # A bare object literal or for-expression inside the open
            # interpolation: its braces must balance, otherwise the
            # closing brace of the nested object would end the
            # interpolation early and a later newline would truncate the
            # scan before the string's real closing quote.
            interp += 1
            i += 1
            continue
        if ch == "}" and interp:
            interp -= 1
        i += 1
    return n - 1


def _blank(chars, begin, end):
    """Replace every non-newline character in ``chars[begin:end]``."""
    for k in range(begin, min(end, len(chars))):
        if chars[k] != "\n":
            chars[k] = " "


def _heredoc_end(text, tag, marker_end):
    """Return the end offset of a heredoc whose tag ends at ``marker_end``.

    The heredoc body ends after the first line that contains only the
    opening tag; an EOF without a closing line ends it at end of text.
    """
    n = len(text)
    nl = text.find("\n", marker_end)
    if nl < 0:
        return n
    k = nl + 1
    while k < n:
        line_end = text.find("\n", k)
        if line_end < 0:
            line_end = n
        if text[k:line_end].strip() == tag:
            return min(line_end + 1, n)
        k = line_end + 1
    return n


def mask_hcl(text):
    """Return ``text`` with comments, string bodies and heredocs blanked.

    Newlines are preserved so the masked text keeps the original line
    structure; positions in the masked text map 1:1 onto the original.
    """
    chars = list(text)
    i = 0
    n = len(text)
    while i < n:
        ch = text[i]
        if ch == "#":
            j = text.find("\n", i)
            _blank(chars, i, n if j < 0 else j)
            i = n if j < 0 else j
        elif ch == "/" and text[i:i + 2] == "//":
            j = text.find("\n", i)
            _blank(chars, i, n if j < 0 else j)
            i = n if j < 0 else j
        elif ch == "/" and text[i:i + 2] == "/*":
            j = text.find("*/", i + 2)
            _blank(chars, i, n if j < 0 else j + 2)
            i = n if j < 0 else j + 2
        elif ch == "<" and text[i:i + 2] == "<<":
            m = HEREDOC_RE.match(text, i)
            if not m:
                i += 1
                continue
            end = _heredoc_end(text, m.group(1), m.end())
            _blank(chars, i, end)
            i = end
        elif ch == '"':
            j = _scan_quoted(text, i)
            _blank(chars, i + 1, j)
            i = j + 1
        else:
            i += 1
    return "".join(chars)


def _unquote(token):
    if token.startswith('"') and token.endswith('"') and len(token) >= 2:
        return token[1:-1]
    return token


def find_required_providers_blocks(text):
    """Return the ``required_providers`` blocks of a (masked) HCL document.

    Each block is a dict with ``open_pos``/``close_pos`` brace offsets in
    the original text, the ``ident_pos`` of the ``required_providers``
    keyword (for indentation) and the provider names declared by its
    ``name = { ... }`` entries. Only blocks nested directly inside a
    top-level ``terraform`` block count; attribute forms
    (``required_providers = { ... }``) are not blocks.
    """
    masked = mask_hcl(text)
    blocks = []
    stack = []  # open frames: {'name': block name or None, 'record': block dict or None}
    prev_ident = None
    prev_ident_pos = -1
    prev_name_token = None  # identifier or quoted string usable as entry name
    expect_block = False  # last tokens look like a block header, not an object value

    def stack_names():
        return [frame["name"] for frame in stack]

    for m in STRUCT_RE.finditer(masked):
        tok = m.group(0)
        if tok == "{":
            name = prev_ident if expect_block else None
            record = None
            if name == "required_providers" and stack_names() == ["terraform"]:
                record = {
                    "open_pos": m.start(),
                    "close_pos": None,
                    "ident_pos": prev_ident_pos,
                    "entries": [],
                }
                blocks.append(record)
            stack.append({"name": name, "record": record})
            expect_block = False
        elif tok == "}":
            if stack:
                frame = stack.pop()
                if frame["record"] is not None:
                    frame["record"]["close_pos"] = m.start()
            expect_block = False
        elif tok == "=":
            if stack_names() == ["terraform", "required_providers"] and prev_name_token:
                blocks[-1]["entries"].append(_unquote(prev_name_token))
            expect_block = False
        elif tok.startswith('"'):
            # Quoted block label or entry key; keeps a pending block header.
            # Read the slice from the original text: mask_hcl blanks string
            # bodies, so the masked token would only be quotes and spaces.
            prev_name_token = text[m.start():m.end()]
        else:
            prev_ident = tok
            prev_ident_pos = m.start()
            prev_name_token = tok
            expect_block = True
    return blocks


def _block_indent(text, block):
    """Leading whitespace of the line holding the block keyword.

    One-liner style (keyword sharing its line with other code) has no
    indentation to inherit.
    """
    line_start = text.rfind("\n", 0, block["ident_pos"]) + 1
    prefix = text[line_start:block["ident_pos"]]
    return "" if prefix.strip() else prefix


def complete_block(text, block):
    """Insert the alicloud requirement into an existing block, in place."""
    prefix = _block_indent(text, block)
    indent = prefix + "  "
    entry = (
        "\n" + indent + "alicloud = {\n"
        + indent + "  source = \"" + ALICLOUD_SOURCE + "\"\n"
        + indent + "}"
    )
    open_pos = block["open_pos"]
    close_pos = block["close_pos"]
    if close_pos is not None and "\n" not in text[open_pos + 1:close_pos]:
        # Single-line block: it can hold only one argument, so adding the
        # alicloud entry requires the multiline form. The existing body
        # and the closing brace each move onto their own line, otherwise
        # terraform reports "Missing newline after argument".
        entry += "\n"
        return (
            text[:open_pos + 1] + entry
            + text[open_pos + 1:close_pos]
            + "\n" + prefix + "}"
            + text[close_pos + 1:]
        )
    tail = text[open_pos + 1:]
    if tail.split("\n", 1)[0].strip():
        # Code shares the opening-brace line with the first entry; keep it
        # on its own line, otherwise HCL sees two arguments separated only
        # by whitespace.
        entry += "\n"
        if close_pos is None:
            return text[:open_pos + 1] + entry + tail
        close_line_start = text.rfind("\n", 0, close_pos) + 1
        if not text[close_line_start:close_pos].strip():
            return text[:open_pos + 1] + entry + tail
        # While an entry shares the opening-brace line HCL parses the block
        # in single-line mode and tolerates the entry's closing brace on
        # the block's closing line ("} }"). Pushing the entry to its own
        # line exits that mode, so the block's closing brace must move
        # onto its own line as well, otherwise terraform reports "Missing
        # newline after argument".
        return (
            text[:open_pos + 1] + entry
            + text[open_pos + 1:close_pos].rstrip(" \t")
            + "\n" + prefix + "}"
            + text[close_pos + 1:]
        )
    return text[:open_pos + 1] + entry + tail


def prepare_example(example_dir, main_tf="main.tf", terraform_tf="terraform.tf"):
    """Apply the provider requirement decision to ``example_dir``.

    Returns one of ``injected`` (terraform.tf written), ``declared``
    (example declaration kept, nothing written), ``completed`` (alicloud
    added to the example's own block) or ``duplicate`` (two or more real
    blocks; nothing is changed so terraform reports the error).
    """
    example_dir = Path(example_dir)
    main_path = example_dir / main_tf
    tf_path = example_dir / terraform_tf
    text = main_path.read_text(encoding="utf-8")

    blocks = find_required_providers_blocks(text)
    if not blocks:
        tf_path.write_text(INJECTED_TERRAFORM_TF, encoding="utf-8")
        return "injected"
    if len(blocks) > 1:
        return "duplicate"
    if "alicloud" in blocks[0]["entries"]:
        return "declared"
    main_path.write_text(complete_block(text, blocks[0]), encoding="utf-8")
    return "completed"


def main(argv=None):
    parser = argparse.ArgumentParser(
        description="Prepare the alicloud provider requirement for one example."
    )
    parser.add_argument(
        "--example-dir",
        required=True,
        help="temporary directory holding main.tf (and terraform.tf)",
    )
    parser.add_argument("--main-tf", default="main.tf")
    parser.add_argument("--terraform-tf", default="terraform.tf")
    args = parser.parse_args(argv)
    status = prepare_example(args.example_dir, args.main_tf, args.terraform_tf)
    print(status)
    return 0


if __name__ == "__main__":
    sys.exit(main())

#!/usr/bin/env bash
# aone-comment-format.sh — normalize Jarvis Aone comments into readable plaintext.
#
# Aone's API stores Markdown-ish text, but the UI rendering can vary. Keep the
# output readable even when rendered as plain text.

set -euo pipefail

if [ "$#" -gt 0 ]; then
    printf '%s' "$*"
else
    cat
fi | python3 -c '
import re
import sys

text = sys.stdin.read().replace("\r\n", "\n").replace("\r", "\n")


def convert_inline_numbered(line: str) -> list[str]:
    matches = list(re.finditer(r"(?<!\d)(\d{1,2})[）)]", line))
    if len(matches) < 2:
        return [line]

    prefix = line[: matches[0].start()].rstrip("；; ")
    out: list[str] = []
    if prefix:
        out.append(prefix)

    for idx, match in enumerate(matches):
        start = match.end()
        end = matches[idx + 1].start() if idx + 1 < len(matches) else len(line)
        item = line[start:end].strip()
        item = item.lstrip("；;、,， ")
        item = item.rstrip("；; ")
        out.append(f"{match.group(1)}. {item}")

    return out


def normalize_line(line: str) -> list[str]:
    stripped = line.strip()
    if not stripped:
        return [""]

    heading = re.match(r"^#{1,6}\s+(.+?)\s*$", stripped)
    if heading:
        return [f"【{heading.group(1)}】"]

    bullet = re.match(r"^[-*]\s+(.+?)\s*$", stripped)
    if bullet:
        return [f"• {bullet.group(1)}"]

    if re.match(r"^```", stripped):
        return []

    line = re.sub(r"`([^`\n]+)`", r"\1", line)
    return convert_inline_numbered(line.rstrip())


lines: list[str] = []
for raw in text.split("\n"):
    lines.extend(normalize_line(raw))

cleaned: list[str] = []
blank_count = 0
for line in lines:
    if line.startswith("代码：") and cleaned and cleaned[-1] != "":
        cleaned.append("")
        blank_count = 1

    if line == "":
        blank_count += 1
        if blank_count <= 2:
            cleaned.append(line)
        continue

    blank_count = 0
    cleaned.append(line.rstrip())

while cleaned and cleaned[0] == "":
    cleaned.pop(0)
while cleaned and cleaned[-1] == "":
    cleaned.pop()

sys.stdout.write("\n".join(cleaned))
'

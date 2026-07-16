#!/usr/bin/env bash
# aone-comment-format.sh — normalize Jarvis Aone comments into readable Markdown.
#
# Aone renders Markdown headings, emphasis, inline code, and fenced code blocks.
# Normalize list syntax that is easy to collapse in the UI, while preserving
# Markdown features that the UI renders correctly.

set -euo pipefail

if [ "$#" -gt 0 ]; then
    printf '%s' "$*"
else
    cat
fi | python3 -c '
import os
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
        out.append(f"{match.group(1)}、{item}")

    return out


def normalize_line(line: str) -> list[str]:
    stripped = line.strip()
    if not stripped:
        return [""]

    bullet = re.match(r"^[-*]\s+(.+?)\s*$", stripped)
    if bullet:
        return [f"· {bullet.group(1)}"]

    ordered = re.match(r"^(\d{1,2})[.)）]\s*(.+?)\s*$", stripped)
    if ordered:
        return [f"{ordered.group(1)}、{ordered.group(2)}"]

    return convert_inline_numbered(line.rstrip())


lines: list[tuple[str, bool]] = []
in_fence = False
for raw in text.split("\n"):
    stripped = raw.strip()
    if re.match(r"^```", stripped):
        lines.append((raw.rstrip(), True))
        in_fence = not in_fence
        continue
    if in_fence:
        lines.append((raw.rstrip(), True))
        continue

    lines.extend((line, False) for line in normalize_line(raw))

cleaned: list[str] = []
blank_count = 0
def is_list_item(line: str) -> bool:
    return bool(re.match(r"^\d{1,2}、", line) or line.startswith("· "))


for line, raw_markdown in lines:
    if raw_markdown:
        cleaned.append(line)
        blank_count = 0
        continue

    if line.startswith("代码：") and cleaned and cleaned[-1] != "":
        cleaned.append("")
        blank_count = 1

    if is_list_item(line):
        if cleaned and cleaned[-1] != "":
            cleaned.append("")
        cleaned.append(line)
        cleaned.append("")
        blank_count = 1
        continue

    if line == "":
        blank_count += 1
        if blank_count <= 2:
            cleaned.append(line)
        continue

    blank_count = 0
    cleaned.append(line.rstrip())

# Bare-URL → clickable markdown link gate.
# Aone does NOT autolink bare URLs in comments (renders them as dead text);
# the only reliable clickable form is markdown `[text](url)` (see aone-triage SKILL
# §4). Auto-wrap any bare http(s) URL that is not already inside a markdown link,
# inline code span, or fenced code block — so the rule is enforced at the write
# path (wrap.sh → aone-comment-format.sh), not left to agent memory mid-flow.
# Set JARVIS_COMMENT_URL_AUTOLINK=0 to disable.
def _wrap_bare_urls_line(line: str) -> str:
    if os.environ.get("JARVIS_COMMENT_URL_AUTOLINK", "1") == "0":
        return line
    url_re = re.compile(r"https?://[^\s)\]<>"
                        r"　-〿一-鿿＀-￯]+")

    def stash(store: list[str], m: re.Match) -> str:
        store.append(m.group(0))
        return "\x00%d\x00" % (len(store) - 1)

    codes: list[str] = []
    links: list[str] = []
    work = re.sub(r"`[^`]*`", lambda m: stash(codes, m), line)
    work = re.sub(r"\[[^\]]*\]\([^)]+\)", lambda m: stash(links, m), work)

    def wrap(m: re.Match) -> str:
        url = m.group(0).rstrip(".,;:!?。，；：！？）】、》")
        return "[%s](%s)" % (url, url)

    work = url_re.sub(wrap, work)

    def restore(store: list[str], s: str) -> str:
        for i, val in enumerate(store):
            s = s.replace("\x00%d\x00" % i, val)
        return s

    work = restore(links, work)
    work = restore(codes, work)
    return work

in_fence_final = False
final: list[str] = []
for line in cleaned:
    if re.match(r"^```", line.strip()):
        in_fence_final = not in_fence_final
        final.append(line)
        continue
    final.append(line if in_fence_final else _wrap_bare_urls_line(line))
cleaned = final

while cleaned and cleaned[0] == "":
    cleaned.pop(0)
while cleaned and cleaned[-1] == "":
    cleaned.pop()

sys.stdout.write("\n".join(cleaned))
'

#!/usr/bin/env python3
"""Render a Markdown evidence report as self-contained, WAF-safe HTML."""

from __future__ import annotations

import argparse
import html
import json
import re
from pathlib import Path


STYLE = """
:root { color-scheme: light; --ink:#172033; --line:#dce2ea; }
* { box-sizing: border-box; }
body { margin:0; background:#eef2f7; color:var(--ink); font:15px/1.68 -apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif; }
main { width:min(1180px,calc(100% - 32px)); margin:28px auto; padding:40px 48px 64px; background:#fff; border:1px solid var(--line); border-radius:14px; box-shadow:0 10px 30px rgba(29,45,72,.08); }
h1 { margin:0 0 24px; font-size:30px; line-height:1.3; }
h2 { margin:38px 0 14px; padding-bottom:8px; border-bottom:2px solid #e8edf4; font-size:22px; }
h3 { margin-top:28px; font-size:18px; }
p,li { max-width:1000px; }
code { padding:2px 5px; border-radius:4px; background:#f0f3f7; color:#b4235a; font-family:"SFMono-Regular",Consolas,monospace; }
pre { overflow:auto; padding:16px; border:1px solid var(--line); border-radius:8px; background:#111827; color:#e5edf8; }
pre code { padding:0; background:transparent; color:inherit; }
table { width:100%; margin:16px 0 26px; border-collapse:collapse; font-size:13px; }
th,td { padding:9px 10px; border:1px solid var(--line); text-align:left; vertical-align:top; overflow-wrap:anywhere; }
th { background:#edf4ff; color:#123967; }
tr:nth-child(even) td { background:#fafbfd; }
blockquote { margin:18px 0; padding:12px 16px; border-left:4px solid #e59b16; background:#fff8e8; color:#654400; }
@media (max-width:720px) { main { width:100%; margin:0; padding:24px 16px 48px; border:0; border-radius:0; } h1{font-size:24px;} table{display:block;overflow-x:auto;} }
""".strip()

FORBIDDEN_HTML_RE = re.compile(
    r"<\s*(?:script|iframe|object|embed|link|meta|base|form|input|button|style)\b"
    r"|\bon[a-z]+\s*="
    r"|(?:href|src)\s*=\s*['\"]?\s*javascript:",
    flags=re.IGNORECASE,
)
FORBIDDEN_RENDERED_HTML_RE = re.compile(
    r"<\s*(?:script|iframe|object|embed|base|form|input|button)\b"
    r"|\bon[a-z]+\s*="
    r"|(?:href|src)\s*=\s*['\"]?\s*javascript:",
    flags=re.IGNORECASE,
)
LANGUAGE_TAG_RE = re.compile(r"[A-Za-z]{2,8}(?:-[A-Za-z0-9]{1,8})*")


def decode_html_entities(value: str) -> str:
    decoded = value
    for _ in range(3):
        candidate = html.unescape(decoded)
        if candidate == decoded:
            break
        decoded = candidate
    return decoded


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("input", type=Path)
    parser.add_argument("output", type=Path)
    parser.add_argument(
        "--lang",
        default="zh-CN",
        help="document language tag (default: zh-CN)",
    )
    args = parser.parse_args()
    if not LANGUAGE_TAG_RE.fullmatch(args.lang):
        raise SystemExit(
            "render-report-html: --lang must be a simple BCP-47 language tag"
        )
    language = html.escape(args.lang, quote=True)

    try:
        import markdown
    except ImportError as error:
        raise SystemExit(
            "render-report-html: Python package 'markdown' is required"
        ) from error

    source = args.input.read_text(encoding="utf-8")
    decoded_source = decode_html_entities(source)
    if re.search(r"data:image/[^;]+;base64,", decoded_source, flags=re.IGNORECASE):
        raise SystemExit("render-report-html: base64 images are forbidden by report WAF")
    if FORBIDDEN_HTML_RE.search(decoded_source):
        raise SystemExit("render-report-html: executable HTML is forbidden in reports")

    title_match = re.search(r"^#\s+(.+)$", source, flags=re.MULTILINE)
    title = title_match.group(1).strip() if title_match else args.input.stem
    body = markdown.markdown(
        source,
        extensions=("tables", "fenced_code", "toc"),
        output_format="html5",
    )
    document = f"""<!doctype html>
<html lang="{language}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{html.escape(title)}</title>
<style>{STYLE}</style>
</head>
<body><main>{body}</main></body>
</html>
"""
    decoded_document = decode_html_entities(document)
    if "data:image" in decoded_document.lower():
        raise SystemExit("render-report-html: rendered HTML contains a base64 image")
    if FORBIDDEN_RENDERED_HTML_RE.search(decoded_document):
        raise SystemExit("render-report-html: rendered executable HTML is forbidden")

    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(document, encoding="utf-8")
    print(
        json.dumps(
            {
                "output": str(args.output),
                "bytes": len(document.encode("utf-8")),
                "base64_images": 0,
                "title": title,
                "lang": args.lang,
            },
            ensure_ascii=False,
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

#!/usr/bin/env python3
"""Standard visual evidence report generator for screenshot-evidence.

Reads a layers JSON (stdin or --layers-file) and emits an HTML report where
each screenshot is wrapped in `<a href="<signed-url>" target="_blank">` so
reviewers click to open the full-resolution image in a new tab. The pre-agent
WAF forbids JS / `<script>` / `<button>` / `on*=` (rgv587_flag:sm), so a
plain anchor is the only allowed zoom path — no lightbox.

Layout uses `<figure>` (max-width 1200px) instead of a narrow table cell, so
screenshot details stay legible without clicking; the anchor is the
"view original" affordance.

Input JSON (array of layer objects):
    [
      {"name": "OpenAPI", "result": "pass",
       "screenshot_url": "https://...", "source_url": "https://...",
       "source_label": "DescribeDBInstances 返回参数", "note": "..."},
      {"name": "CloudSpec/ACube", "result": "n-a", "note": "pure datasource"},
      {"name": "Provider", "result": "pass",
       "screenshot_url": "https://...", "source_url": "https://...",
       "source_label": "PR 10086 diff", "note": "..."}
    ]

Usage:
    gen-report.py --title "可视化查证报告 — Aone #XXX" \
                  --layers-file layers.json [--summary "..."] > report.html
    cat layers.json | gen-report.py --title "..." > report.html
"""
import argparse
import html as html_mod
import json
import sys
from pathlib import Path

_BADGE = {"pass": "✅ pass", "n-a": "— n/a", "fail": "❌ fail"}


def _esc(value):
    return html_mod.escape(str(value if value is not None else ""), quote=True)


def _layer_block(layer):
    name = layer.get("name", "")
    result = layer.get("result", "n-a")
    shot_url = layer.get("screenshot_url", "") or ""
    source_url = layer.get("source_url", "") or ""
    source_label = layer.get("source_label", "") or source_url or "N/A"
    note = layer.get("note", "")
    badge = _BADGE.get(result, result)
    if shot_url:
        img = (
            f'<a href="{_esc(shot_url)}" target="_blank" rel="noopener">'
            f'<img src="{_esc(shot_url)}" alt="{_esc(name)}" '
            f'style="max-width:100%;border:1px solid #ddd;border-radius:4px"></a>'
            f'<div style="font-size:.8em;color:#6b7280;margin-top:4px">'
            f'点击图片在新标签打开原图</div>')
    else:
        img = '<em>N/A</em>'
    src = (f'<a href="{_esc(source_url)}">{_esc(source_label)}</a>'
           if source_url else _esc(source_label))
    return f"""
<figure style="margin:16px 0;padding:12px;border:1px solid #e5e7eb;border-radius:6px;background:#fff">
  <figcaption><strong>{_esc(name)}</strong> <span class="badge">{badge}</span> · {src}</figcaption>
  <div style="margin-top:8px">{img}</div>
  <p style="margin:8px 0 0;font-size:.92em;color:#374151">{_esc(note)}</p>
</figure>"""


def render(title, layers, summary=""):
    blocks = "".join(_layer_block(l) for l in layers)
    summary_html = (
        f'<p style="font-size:.85em;color:#6b7280">{_esc(summary)}</p>'
        if summary else "")
    return f"""<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{_esc(title)}</title>
<style>
  body{{font-family:-apple-system,"PingFang SC","Microsoft YaHei",sans-serif;max-width:1200px;margin:24px auto;padding:0 16px;color:#222;line-height:1.6}}
  h1{{font-size:1.35em;border-bottom:2px solid #2563eb;padding-bottom:8px}}
  .badge{{display:inline-block;padding:2px 8px;border-radius:10px;font-size:.8em;background:#eef2ff;color:#3730a3}}
  a{{color:#2563eb}}
  figure img{{cursor:zoom-in}}
</style></head><body>
<h1>{_esc(title)}</h1>
{summary_html}
{blocks}
</body></html>"""


def main(argv):
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("--title", required=True)
    p.add_argument("--layers-file", help="JSON file of layers (default: stdin)")
    p.add_argument("--summary", default="")
    p.add_argument("--out", help="output path (default: stdout)")
    args = p.parse_args(argv)
    if args.layers_file:
        layers = json.loads(Path(args.layers_file).read_text(encoding="utf-8"))
    else:
        layers = json.loads(sys.stdin.read())
    report = render(args.title, layers, args.summary)
    if args.out:
        Path(args.out).write_text(report, encoding="utf-8")
        print("wrote %s (%d bytes)" % (args.out, len(report)))
    else:
        sys.stdout.write(report)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

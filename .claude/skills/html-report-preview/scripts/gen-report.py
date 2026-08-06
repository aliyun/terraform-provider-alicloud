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

    # From an evidence-manifest.md: uploads each local screenshot via
    # upload-screenshots.sh and embeds the returned signed URL. Fails closed
    # (exit 3) if any manifest screenshot has no signed URL, so a screenshot-less
    # "visual evidence" report can never be produced by accident.
    gen-report.py --title "..." --manifest evidence-manifest.md \
                  --aone-id 12345678 > report.html
    # Or reuse an earlier upload:
    gen-report.py --title "..." --manifest evidence-manifest.md \
                  --image-urls image-urls.txt > report.html
"""
import argparse
import html as html_mod
import json
import os
import subprocess
import sys
from pathlib import Path

_BADGE = {"pass": "✅ pass", "n-a": "— n/a", "fail": "❌ fail"}
_IMAGE_SUFFIXES = (".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp")


class ManifestError(RuntimeError):
    """Manifest could not be turned into a report with every screenshot embedded."""


def _split_row(line):
    return [cell.strip() for cell in line.strip().strip("|").split("|")]


def parse_manifest(text):
    """Extract layer rows from an evidence-manifest.md pipe table.

    Returns a list of dicts with the manifest's own column names, so a row's
    `screenshot` value is still the LOCAL path at this point.
    """
    header = None
    rows = []
    for line in text.splitlines():
        stripped = line.strip()
        if not stripped.startswith("|"):
            continue
        cells = _split_row(stripped)
        if set("".join(cells)) <= set("-: "):
            continue
        if header is None:
            header = [cell.lower() for cell in cells]
            continue
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise ManifestError("no pipe table found in manifest")
    for required in ("layer", "result", "screenshot"):
        if required not in header:
            raise ManifestError("manifest table has no '%s' column" % required)
    if not rows:
        raise ManifestError("manifest table has no layer rows")
    return rows


def parse_image_urls(text):
    """Parse `name|signed_url` lines emitted by upload-screenshots.sh."""
    mapping = {}
    for line in text.splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "|" not in line:
            continue
        name, _, url = line.partition("|")
        name, url = name.strip(), url.strip()
        if name and url:
            mapping[name] = url
    return mapping


def _local_screenshot(value):
    """Local screenshot path in a manifest cell, or "" when absent/already a URL."""
    value = (value or "").strip()
    if not value or value.lower() in {"n/a", "na", "-", "—"}:
        return ""
    if "://" in value:
        return ""
    if not value.lower().endswith(_IMAGE_SUFFIXES):
        return ""
    return value


def default_uploader():
    """Sibling skill's uploader, resolved without hardcoding the skills root.

    This file lives at <skills>/html-report-preview/scripts/gen-report.py, so the
    uploader is found the same way in every mirrored skills tree.
    """
    return (Path(__file__).resolve().parents[2]
            / "screenshot-evidence" / "scripts" / "upload-screenshots.sh")


def upload_screenshots(aone_id, directory, uploader):
    """Run the repo-owned uploader and return its name->signed_url mapping."""
    script = Path(uploader) if uploader else default_uploader()
    if not script.is_file():
        raise ManifestError("uploader not found: %s" % script)
    try:
        done = subprocess.run(
            ["bash", str(script), str(aone_id), str(directory)],
            capture_output=True, text=True, check=False)
    except OSError as exc:
        raise ManifestError("uploader could not run: %s" % exc.strerror)
    if done.returncode != 0:
        # stderr carries only the uploader's fixed sanitized codes.
        raise ManifestError("uploader failed (exit %d): %s"
                            % (done.returncode, done.stderr.strip()))
    return parse_image_urls(done.stdout)


def manifest_layers(rows, image_urls):
    """Attach signed URLs to manifest rows, failing closed on any gap.

    A screenshot named in the manifest but missing from image_urls would render
    as an image-free "visual evidence" report, so it is an error, not a warning.
    """
    layers = []
    missing = []
    for row in rows:
        local = _local_screenshot(row.get("screenshot"))
        shot_url = ""
        if local:
            stem = os.path.splitext(os.path.basename(local))[0]
            shot_url = image_urls.get(stem, "")
            if not shot_url:
                missing.append(stem)
        layers.append({
            "name": row.get("layer", ""),
            "result": row.get("result", "n-a"),
            "screenshot_url": shot_url,
            "source_url": row.get("source", ""),
            "source_label": row.get("source", "") or "N/A",
            "note": row.get("note", ""),
        })
    if missing:
        raise ManifestError(
            "no signed URL for screenshot(s): %s; run upload-screenshots.sh and "
            "pass --image-urls, or supply --aone-id to upload here"
            % ", ".join(sorted(set(missing))))
    return layers


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
    p.add_argument("--manifest",
                   help="evidence-manifest.md; screenshots become signed URLs")
    p.add_argument("--image-urls",
                   help="name|signed_url file from upload-screenshots.sh; "
                        "omit to upload the manifest's screenshots here")
    p.add_argument("--aone-id", help="Aone id, required to upload from --manifest")
    p.add_argument("--uploader",
                   help="path to upload-screenshots.sh (default: sibling skill)")
    p.add_argument("--summary", default="")
    p.add_argument("--out", help="output path (default: stdout)")
    args = p.parse_args(argv)
    if args.manifest and args.layers_file:
        p.error("--manifest and --layers-file are mutually exclusive")
    if args.manifest:
        try:
            rows = parse_manifest(Path(args.manifest).read_text(encoding="utf-8"))
            if args.image_urls:
                urls = parse_image_urls(
                    Path(args.image_urls).read_text(encoding="utf-8"))
            elif any(_local_screenshot(r.get("screenshot")) for r in rows):
                if not args.aone_id:
                    raise ManifestError(
                        "--aone-id is required to upload screenshots; or pass "
                        "--image-urls from a previous upload-screenshots.sh run")
                shots = {os.path.dirname(_local_screenshot(r.get("screenshot")))
                         for r in rows if _local_screenshot(r.get("screenshot"))}
                if len(shots) != 1:
                    raise ManifestError(
                        "screenshots span %d directories; upload them yourself "
                        "and pass --image-urls" % len(shots))
                urls = upload_screenshots(args.aone_id, shots.pop(), args.uploader)
            else:
                urls = {}
            layers = manifest_layers(rows, urls)
        except (OSError, ManifestError) as exc:
            sys.stderr.write("gen-report: %s\n" % exc)
            return 3
    elif args.layers_file:
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

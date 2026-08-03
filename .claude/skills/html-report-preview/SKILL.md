---
name: html-report-preview
description: Upload HTML report artifacts to AutomationAgent online preview storage and return shareable review links. Use when Jarvis/Claude Code needs to upload local .html/.htm files, .zip files containing HTML, directories of HTML reports, or Aone work item attachments for online review; when asked to comment preview links back to Aone; or when configuring/troubleshooting JARVIS_HTML_REPORT_TOKEN, /api/reports/aone uploads, or /reports/aone preview links.
---

# HTML Report Preview

## Overview

Use Jarvis's repository-owned helper to publish HTML reports to AutomationAgent and return a browser-viewable link. The durable implementation is `bootstrap/html-report-preview.sh`; do not reimplement upload or Aone attachment parsing in chat.

## Quick Start

Run commands from the Jarvis repository root or a Jarvis worktree:

```bash
bash bootstrap/html-report-preview.sh upload <aone-id> <html|zip|dir>
bash bootstrap/html-report-preview.sh from-aone <aone-id>
```

Useful options:

- `--comment`: post the generated preview links back to a non-Terraform Aone work item.
- `--format jsonl`: emit machine-readable one-line JSON records.
- `--attachment-id <id>`: upload one specific Aone attachment.
- `--all`: upload all HTML/ZIP attachments from the Aone work item.
- `--base-url <url>`: override `JARVIS_HTML_REPORT_BASE_URL`.

Default target is preprod: `https://pre-agent.aliyun-inc.com`.

### Terraform single-writer boundary

For a Terraform main-processing run, PD/QA never upload. The RD finalizer uploads exactly once
without `--comment`, captures the returned preview URL, and returns it in the only aggregate reply.
When bridge executor owns the bookend, that link must go into `AONE_RESULT.reply_body`; the model run
must not call `wrap.sh` or post any Aone comment. `--comment` is allowed only for non-Terraform
workflows or an explicitly non-executor-managed context.

## Token Handling

Uploads require `Authorization: Bearer $JARVIS_HTML_REPORT_TOKEN`. The helper also loads missing variables from `bootstrap/.env`, which is gitignored.

Never commit the plaintext token to tracked files, skill files, tests, or Aone comments. If the token is missing, the helper must return exit code `3`; JSONL mode emits `success:false,status:blocked,code:missing_token` before any curl or Aone call. Ask for it to be configured in the environment or gitignored `bootstrap/.env`; do not fall back to personal browser cookies or a personal BUC session for this server-token upload path.

## Image Handling (WAF constraint)

**Do NOT embed screenshots as base64 `data:` URIs in the uploaded HTML.** The preprod WAF (Anti-Bot `rgv587`) runs content inspection on the multipart upload and blocks any HTML whose body carries base64 image data — PNG and JPEG alike, regardless of size. Symptom: `POST /api/reports/aone/<id>` returns HTTP 200 whose body is a `waf_block*.html` / `punish` page instead of the `{"success":true,...,"viewUrl":...}` JSON. This is not a token, endpoint, or method problem (an empty-body POST reaches the app and returns a clean 415).

Note: the WAF classification gate requires an `X-Request-Context` header — the helper sends it automatically on every upload (default `rctx_a3f90b7e2d41c8f6`, override via env `JARVIS_HTML_REPORT_WAF_HEADER`). It does **not** exempt base64-image payloads: even with the header, inline `data:` images are still blocked — keep images external per the workaround below.

Before the first network request, `html-report-preview.sh` validates every HTML member in the
input file, directory, or ZIP. Image references are allowed only when all of these conditions hold:

- every `<img src>` and every candidate in `<img srcset>` is an absolute `https://` URL;
- every candidate in `<picture><source srcset>` is an absolute `https://` URL;
- relative paths, protocol-relative URLs (`//...`), `http:`, `file:`, and `data:` are rejected.

The whole batch fails closed with `invalid_image_reference`; no report in that batch is uploaded.
`<source srcset>` outside `<picture>` is not an image source and is outside this check.

**Workaround — use AutomationAgent's private image upload API and signed GET URLs:**

Never make report screenshots public-read. Use the repository-owned screenshot helper; it loads the
same server token through `bootstrap/runtime-config.sh` and uploads each PNG/JPG as multipart:

```bash
bash .claude/skills/screenshot-evidence/scripts/upload-screenshots.sh \
  <aone-id> <screenshot-dir> > image-urls.txt
```

The helper calls `POST /api/reports/aone/<aone-id>/images` with multipart field `file` and emits
`name|signed_url`. In the report HTML use the returned HTTPS URL as
`<img src="<signed-url>">` (no base64, no relative path), then upload the HTML normally.

Image storage is a service-side boundary. AutomationAgent alone accesses the private
`jarvis-upload-files` bucket. Before any storage operation it must verify through STS that the caller
account is exactly `1983056807138283`, discover the bucket region/endpoint, write private objects,
and return only time-limited signed GET URLs. Identity mismatch, unavailable identity/region,
upload failure, or signing failure must fail closed.

Validation before sharing the report:

- Server acceptance tests must prove an unsigned object URL is not publicly readable (normally 403).
- Signed URL must load with GET, for example `curl -fL "<signed-url>" -o /tmp/shot.jpg`. Do not use HEAD as the proof; OSS signatures are method-sensitive and a GET-signed URL can return 403 for HEAD.
- The uploaded preview page must load the image in a browser (`naturalWidth > 0` / image `onload`). Record the signed URL expiry time in notes when the report is intended to stay reviewable for months.

Credential boundary: HTML and screenshot clients rely only on `JARVIS_HTML_REPORT_TOKEN`. Never
configure, decrypt, log, or pass storage credentials to an agent or local client; do not fall back
to personal AK/SK, browser cookies, a different bucket, or direct local storage CLI calls.

Keep screenshots reasonably compressed (e.g. `sips -Z 900 -s format jpeg`). Text/tables/CSS in the HTML are fine; only base64 image blobs trip the gate.

## Executable HTML and Viewer boundary

Upload only static HTML. Any inline script, including a local-only copy button, can be rejected by
the pre-agent WAF with `rgv587_flag:sm`. Do not add `<script>`, `<button>`, `on*=`,
`javascript:` or other executable markup, and **禁止规避** or disguise the WAF signature.

Viewer-side HCL copy is owned by the AutomationAgent app/viewer. Until the platform implements it,
record `viewer_copy=platform_blocked`; do not emulate the feature inside uploaded HTML.

## Standard report template

Use `scripts/gen-report.py` to build the report HTML so screenshots are legible and clickable
instead of each PD/finalizer hand-rolling a narrow-table layout. It wraps each `<img>` in
`<a href="<signed-url>" target="_blank">` (the only zoom path the WAF allows — no JS lightbox)
and uses a `<figure>` layout at `max-width:1200px` so details stay readable without clicking.

```bash
# layers.json: [{name, result, screenshot_url, source_url, source_label, note}, ...]
python3 .claude/skills/html-report-preview/scripts/gen-report.py \
  --title "可视化查证报告 — Aone #<id>" \
  --layers-file layers.json [--summary "..."] > report.html
```

`screenshot_url` must be the signed URL returned by `upload-screenshots.sh`. Do NOT hand-write
`<img style="max-width:100%">` inside a narrow table `<td>` — that produced 36px-tall unreadable
screenshots (SPA unrendered + td shrink + no click-to-zoom). The template floors layout at
viewport width and anchors the original.

## Workflows

For a non-Terraform local report:

```bash
bash bootstrap/html-report-preview.sh upload 83873535 /path/to/report.html --comment
```

For a ZIP or directory, pass the ZIP or directory path. The helper extracts/uploads each `.html` or `.htm` file and prints one preview link per file.

For an Aone work item attachment:

```bash
bash bootstrap/html-report-preview.sh from-aone 83873535 --comment
```

Without `--attachment-id` or `--all`, `from-aone` selects the newest `.html`, `.htm`, or `.zip` attachment. Aone operations must keep using `bin/a1id` through the helper's default `JARVIS_A1_BIN`.

For a Terraform finalizer, omit `--comment` and hand the returned URL to the orchestrator:

```bash
bash bootstrap/html-report-preview.sh upload <aone-id> <report.html>
```

## Link Rules

The upload API is server-token protected:

```text
POST /api/reports/aone/<aone-id>
```

The returned review link must be the public readonly route:

```text
/reports/aone/<aone-id>/<report-id>/view
```

The helper rejects any mismatched `/buc/*` or other route. In JSONL mode a successful upload emits
`success:true,status:uploaded` together with `reportId`, `viewUrl`, and the absolute `url`.
Rejected, malformed, transport-failed, or wrong-route server responses are never echoed. JSONL mode
emits a fixed `success:false,status:failed,code:<sanitized-code>` record so a reflected bearer token
or backend detail cannot leak through stderr/stdout.

If a viewer sees `{"error":"buc_required" ...}`, the shared link is probably a stale `/buc/*` URL or the deployed service is returning the wrong view route. Re-upload after the preprod fix or manually switch to the `/reports/aone/.../view` route only when the report id is known.

## Verification

After upload, verify the returned preview URL without credentials:

```bash
curl -fsSI '<preview-url>'
```

Expect HTTP 200. For content checks, use `curl -fsSL '<preview-url>' | head` and confirm the expected report title or marker text is present.

---
name: html-report-preview
description: Upload HTML report artifacts to AutomationAgent online preview storage and return shareable review links. Use when Jarvis/Codex needs to upload local .html/.htm files, .zip files containing HTML, directories of HTML reports, or Aone work item attachments for online review; when asked to comment preview links back to Aone; or when configuring/troubleshooting JARVIS_HTML_REPORT_TOKEN, /api/reports/aone uploads, or /reports/aone preview links.
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

**Workaround — host images externally with private OSS objects and signed GET URLs:**

Never make report screenshots public-read. Upload them to a private bucket or as private objects, then reference time-limited signed GET URLs from the HTML.

1. Upload each screenshot to OSS with private ACL:
   ```bash
   aliyun oss cp shot.jpg oss://<bucket>/<path>/shot.jpg --acl private -e oss-<region>.aliyuncs.com -f
   ```
2. Generate a signed GET URL and use that URL in the report HTML. For a half-year expiry, use `15768000` seconds:
   ```bash
   aliyun oss sign oss://<bucket>/<path>/shot.jpg --timeout 15768000 -e oss-<region>.aliyuncs.com
   ```
3. In the report HTML use `<img src="<signed-url>">` (HTTPS absolute URL, no base64, no relative path).
4. Upload the HTML — with only URLs in the body it passes the WAF, and the preview renders the images from OSS.

Validation before sharing the report:

- Unsigned direct object URL must not be publicly readable; `curl -fsS "https://<bucket>.oss-<region>.aliyuncs.com/<path>/shot.jpg"` should fail, normally with 403.
- Signed URL must load with GET, for example `curl -fL "<signed-url>" -o /tmp/shot.jpg`. Do not use HEAD as the proof; OSS signatures are method-sensitive and a GET-signed URL can return 403 for HEAD.
- The uploaded preview page must load the image in a browser (`naturalWidth > 0` / image `onload`). Record the signed URL expiry time in notes when the report is intended to stay reviewable for months.

Credential boundary: HTML report upload only relies on `JARVIS_HTML_REPORT_TOKEN`; OSS upload/signing needs a Jarvis/team-approved private bucket and least-privilege credentials. Do not use arbitrary personal AKSK or random buckets. If no approved private OSS signing capability exists, escalate instead of claiming image evidence is fixed.

Keep screenshots reasonably compressed (e.g. `sips -Z 900 -s format jpeg`). Text/tables/CSS in the HTML are fine; only base64 image blobs trip the gate.

## Executable HTML and Viewer boundary

Upload only static HTML. Any inline script, including a local-only copy button, can be rejected by
the pre-agent WAF with `rgv587_flag:sm`. Do not add `<script>`, `<button>`, `on*=`,
`javascript:` or other executable markup, and **禁止规避** or disguise the WAF signature.

Viewer-side HCL copy is owned by the AutomationAgent app/viewer. Until the platform implements it,
record `viewer_copy=platform_blocked`; do not emulate the feature inside uploaded HTML.

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

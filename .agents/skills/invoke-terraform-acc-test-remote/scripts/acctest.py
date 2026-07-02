#!/usr/bin/env python3
"""
Terraform AccTest CLI - ACube AccTest API client

Usage:
  acctest.py upload     --namespace NS --resource RES --dir PATH [--test-case NAME] [--base-url URL]
  acctest.py status     --task-id ID [--log-limit N] [--log-since TS] [--base-url URL]
  acctest.py logs       --task-id ID [--base-url URL] [--download-dir DIR]
  acctest.py download   --url URL [--download-dir DIR]
  acctest.py cancel     --task-id ID [--base-url URL]
  acctest.py upload-run --namespace NS --resource RES --dir PATH [--test-case NAME] [--base-url URL] [--poll-interval SEC] [--download-dir DIR]

Commands:
  upload      Upload local code zip and submit AccTest task
  status      Query task status (response includes recent log tail + FC lifecycle fields)
  logs        Get presigned download URLs for run.log / tf-debug.log, optionally download
  download    Download a single file from a presigned OSS URL
  cancel      Cancel a running task (calls FC StopAsyncTask and marks ACube task FAILED)
  upload-run  Full workflow: upload local code -> poll until done -> download logs

Note: The legacy GitLab-branch submission path (submit / run) has been removed —
use upload-run to package local code instead.
"""

import argparse
import io
import json
import os
import sys
import tempfile
import time
import urllib.request
import urllib.error
import urllib.parse
import zipfile

DEFAULT_BASE_URL = "https://acube.aliyun-inc.com"
DEFAULT_POLL_INTERVAL = 60
DEFAULT_DOWNLOAD_DIR = "./acctest_logs"


def api_request(method, url, data=None):
    """Make HTTP request and return parsed JSON response."""
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    body = json.dumps(data).encode("utf-8") if data else None
    req = urllib.request.Request(url, data=body, headers=headers, method=method)

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            raw = resp.read().decode("utf-8")
            return json.loads(raw)
    except urllib.error.HTTPError as e:
        body_text = e.read().decode("utf-8", errors="replace")
        print(f"HTTP {e.code}: {body_text}", file=sys.stderr)
        sys.exit(1)
    except urllib.error.URLError as e:
        print(f"Connection error: {e.reason}", file=sys.stderr)
        sys.exit(1)


def download_file(url, dest_path):
    """Download a file from URL to local path."""
    try:
        urllib.request.urlretrieve(url, dest_path)
        size = os.path.getsize(dest_path)
        return size
    except Exception as e:
        print(f"Download failed: {e}", file=sys.stderr)
        return -1


# ── Commands ──────────────────────────────────────────────────────────

def cmd_status(args):
    """Query task status."""
    url = f"{args.base_url}/api/v1/terraform_acc_test/{args.task_id}"
    resp = api_request("GET", url)

    if resp.get("code") != "SUCCESS":
        print(json.dumps(resp, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    data = resp["data"]
    print(json.dumps(data, indent=2, ensure_ascii=False))
    return data


def cmd_logs(args):
    """Get presigned download URLs, optionally download files."""
    url = f"{args.base_url}/api/v1/terraform_acc_test/{args.task_id}/logs"
    resp = api_request("GET", url)

    if resp.get("code") != "SUCCESS":
        print(json.dumps(resp, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    data = resp["data"]

    # If --download-dir is set, download files locally
    download_dir = getattr(args, "download_dir", None)
    if download_dir:
        os.makedirs(download_dir, exist_ok=True)
        downloaded = {}
        for key, dl_url in data.items():
            if not dl_url:
                continue
            filename = "run.log" if key == "runLog" else "tf-debug.log"
            dest = os.path.join(download_dir, f"{args.task_id}_{filename}")
            size = download_file(dl_url, dest)
            if size >= 0:
                downloaded[key] = {"path": os.path.abspath(dest), "size": size}
                print(f"[AccTest] Downloaded {filename}: {dest} ({size} bytes)", file=sys.stderr)
        result = {"taskId": args.task_id, "logs": downloaded}
        print(json.dumps(result, indent=2, ensure_ascii=False))
        return result
    else:
        print(json.dumps(data, indent=2, ensure_ascii=False))
        return data


def cmd_download(args):
    """Download a file from a presigned OSS URL."""
    dl_url = args.url
    download_dir = args.download_dir
    os.makedirs(download_dir, exist_ok=True)

    # Infer filename from URL path
    url_path = urllib.parse.urlparse(dl_url).path
    filename = os.path.basename(url_path) or "downloaded_file"
    dest = os.path.join(download_dir, filename)

    print(f"[AccTest] Downloading: {filename}", file=sys.stderr)
    size = download_file(dl_url, dest)
    if size < 0:
        print(f"[AccTest] Download failed", file=sys.stderr)
        sys.exit(1)

    result = {"path": os.path.abspath(dest), "size": size, "filename": filename}
    print(json.dumps(result, indent=2, ensure_ascii=False))
    print(f"[AccTest] Saved: {dest} ({size} bytes)", file=sys.stderr)
    return result


def cmd_cancel(args):
    """Cancel a running AccTest task.

    Calls FC StopAsyncTask to terminate the FC instance and marks the ACube task as FAILED.
    Returns the task's snapshot after cancellation.
    Idempotency: cancelling a terminal task returns an error from the server.
    """
    url = f"{args.base_url}/api/v1/terraform_acc_test/{args.task_id}/cancel"
    resp = api_request("POST", url)

    if resp.get("code") != "SUCCESS":
        print(json.dumps(resp, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    data = resp.get("data") or {}
    print(json.dumps(data, indent=2, ensure_ascii=False))
    print(f"[AccTest] Cancelled task {args.task_id}: status={data.get('status')}, fcStatus={data.get('fcStatus')}",
          file=sys.stderr)
    return data


# ── Zip packaging (whitelist-based, minimal for `go test ./alicloud`) ──
# Remote runner executes `TF_ACC=1 go test ./alicloud`; anything outside of
# these paths is irrelevant. Whitelist avoids pulling in built binaries,
# node_modules, acctest_logs, doc artifacts, etc. that can bloat the zip
# from ~12 MB to 500+ MB.
INCLUDED_TOP_LEVEL_FILES = frozenset({"go.mod", "go.sum"})
INCLUDED_DIRS = frozenset({"alicloud"})

# Excluded at any depth inside the included dirs
EXCLUDED_SUBDIRS = frozenset({"vendor", "node_modules", "testdata_bak"})

# Build/test artifact extensions — never needed for remote runs
EXCLUDED_EXTENSIONS = frozenset({".test", ".prof", ".coverprofile", ".out"})

# Files larger than this inside included dirs are almost certainly build
# artifacts (compiled binaries, dumped coverage, etc.), skip them.
LARGE_FILE_BYTES = 10 * 1024 * 1024  # 10 MB

# Required string in go.mod for a directory to be considered the alicloud Terraform provider.
PROVIDER_MODULE_MARKER = "terraform-provider-alicloud"


def validate_provider_dir(dir_path):
    """Hard-fail if dir_path is not a terraform-provider-alicloud working tree.

    Catches the common mistake of running from an unrelated directory (acube, home,
    a different Go project) — otherwise we'd zip the wrong tree and waste minutes
    of upload + FC startup before seeing a "0 tests found" failure.

    Checks (any failure exits 1 with a clear hint):
      1. dir exists
      2. <dir>/go.mod exists
      3. <dir>/go.mod references "terraform-provider-alicloud"
      4. <dir>/alicloud/ is a directory
    Soft warning (no exit):
      - <dir>/alicloud/ contains at least one *_test.go file
    """
    abs_dir = os.path.abspath(dir_path)
    hint = (f"[AccTest] Hint: pass --dir pointing to your local terraform-provider-alicloud root,\n"
            f"[AccTest]       or cd into that directory first.")

    if not os.path.isdir(abs_dir):
        print(f"[AccTest] Error: directory not found: {abs_dir}", file=sys.stderr)
        print(hint, file=sys.stderr)
        sys.exit(1)

    go_mod = os.path.join(abs_dir, "go.mod")
    if not os.path.isfile(go_mod):
        print(f"[AccTest] Error: not a terraform-provider-alicloud directory — "
              f"go.mod not found in {abs_dir}.",
              file=sys.stderr)
        print(hint, file=sys.stderr)
        sys.exit(1)

    try:
        with open(go_mod, "r", encoding="utf-8") as f:
            go_mod_content = f.read()
    except Exception as e:
        print(f"[AccTest] Error: could not read {go_mod}: {e}", file=sys.stderr)
        sys.exit(1)

    if PROVIDER_MODULE_MARKER not in go_mod_content:
        print(f"[AccTest] Error: {go_mod} does not reference '{PROVIDER_MODULE_MARKER}'.",
              file=sys.stderr)
        print(f"[AccTest] This looks like a Go module but not the alicloud Terraform provider.",
              file=sys.stderr)
        print(hint, file=sys.stderr)
        sys.exit(1)

    alicloud_dir = os.path.join(abs_dir, "alicloud")
    if not os.path.isdir(alicloud_dir):
        print(f"[AccTest] Error: {alicloud_dir} not found.",
              file=sys.stderr)
        print(f"[AccTest] Expected layout: <provider_dir>/alicloud/*_test.go", file=sys.stderr)
        print(hint, file=sys.stderr)
        sys.exit(1)

    # Soft warning if no test files
    try:
        has_test = any(f.endswith("_test.go") for f in os.listdir(alicloud_dir))
    except OSError:
        has_test = False
    if not has_test:
        print(f"[AccTest] Warning: no *_test.go files found under {alicloud_dir}", file=sys.stderr)

    print(f"[AccTest] Verified provider directory: {abs_dir}", file=sys.stderr)


def create_code_zip(dir_path):
    """Create a zip containing only files needed to run `go test ./alicloud`.

    Whitelist approach: package go.mod, go.sum, and the alicloud/ tree
    (with hidden dirs / vendor / node_modules / large binaries filtered out).
    Caller should have already run validate_provider_dir() to ensure this is
    the right directory.
    """
    dir_path = os.path.abspath(dir_path)
    if not os.path.isdir(dir_path):
        print(f"Error: directory not found: {dir_path}", file=sys.stderr)
        sys.exit(1)

    buf = io.BytesIO()
    base_name = os.path.basename(dir_path)
    file_count = 0
    skipped_large = []

    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
        # 1. Whitelisted top-level files
        for name in INCLUDED_TOP_LEVEL_FILES:
            src = os.path.join(dir_path, name)
            if os.path.isfile(src):
                zf.write(src, os.path.join(base_name, name))
                file_count += 1
            else:
                print(f"[AccTest] Warning: missing top-level file: {name}",
                      file=sys.stderr)

        # 2. Whitelisted directories (recursive)
        for included in INCLUDED_DIRS:
            dir_full = os.path.join(dir_path, included)
            if not os.path.isdir(dir_full):
                print(f"[AccTest] Warning: missing directory: {included}",
                      file=sys.stderr)
                continue
            for root, dirs, files in os.walk(dir_full):
                dirs[:] = [
                    d for d in dirs
                    if not d.startswith(".") and d not in EXCLUDED_SUBDIRS
                ]
                for f in files:
                    if f.startswith("."):
                        continue
                    if any(f.endswith(ext) for ext in EXCLUDED_EXTENSIONS):
                        continue
                    full_path = os.path.join(root, f)
                    try:
                        size = os.path.getsize(full_path)
                    except OSError:
                        continue
                    if size > LARGE_FILE_BYTES:
                        skipped_large.append(
                            (os.path.relpath(full_path, dir_path), size)
                        )
                        continue
                    arc_name = os.path.join(
                        base_name, os.path.relpath(full_path, dir_path)
                    )
                    zf.write(full_path, arc_name)
                    file_count += 1

    if skipped_large:
        print(f"[AccTest] Skipped {len(skipped_large)} file(s) "
              f">{LARGE_FILE_BYTES // (1024*1024)}MB (likely build artifacts):",
              file=sys.stderr)
        for path, size in skipped_large[:5]:
            print(f"[AccTest]   {path} ({size // (1024*1024)} MB)",
                  file=sys.stderr)
        if len(skipped_large) > 5:
            print(f"[AccTest]   ... and {len(skipped_large) - 5} more",
                  file=sys.stderr)

    print(f"[AccTest] Packaged {file_count} file(s)", file=sys.stderr)
    buf.seek(0)
    return buf


def multipart_upload(url, namespace, resource, zip_buf, test_case=None):
    """Send multipart/form-data POST with zip file.

    test_case (optional): appended as query parameter (per backend convention).
    """
    if test_case:
        sep = "&" if "?" in url else "?"
        url = f"{url}{sep}testCaseName={urllib.parse.quote(test_case)}"

    boundary = "----AccTestBoundary" + str(int(time.time()))
    body = io.BytesIO()

    def write_field(name, value):
        body.write(f"--{boundary}\r\n".encode())
        body.write(f'Content-Disposition: form-data; name="{name}"\r\n\r\n'.encode())
        body.write(f"{value}\r\n".encode())

    write_field("namespace", namespace)
    write_field("resourceTypeCode", resource)

    # File part
    body.write(f"--{boundary}\r\n".encode())
    body.write(f'Content-Disposition: form-data; name="file"; filename="code.zip"\r\n'.encode())
    body.write(b"Content-Type: application/zip\r\n\r\n")
    body.write(zip_buf.read())
    body.write(b"\r\n")
    body.write(f"--{boundary}--\r\n".encode())

    data = body.getvalue()
    headers = {
        "Content-Type": f"multipart/form-data; boundary={boundary}",
        "Accept": "application/json",
    }
    req = urllib.request.Request(url, data=data, headers=headers, method="POST")

    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            raw = resp.read().decode("utf-8")
            return json.loads(raw)
    except urllib.error.HTTPError as e:
        body_text = e.read().decode("utf-8", errors="replace")
        print(f"HTTP {e.code}: {body_text}", file=sys.stderr)
        sys.exit(1)
    except urllib.error.URLError as e:
        print(f"Connection error: {e.reason}", file=sys.stderr)
        sys.exit(1)


def cmd_upload(args):
    """Upload local code zip and submit AccTest task."""
    dir_path = args.dir
    validate_provider_dir(dir_path)
    print(f"[AccTest] Zipping directory: {dir_path}", file=sys.stderr)
    zip_buf = create_code_zip(dir_path)
    zip_size = zip_buf.getbuffer().nbytes
    print(f"[AccTest] Zip created: {zip_size} bytes", file=sys.stderr)

    url = f"{args.base_url}/api/v1/terraform_acc_test/upload"
    print(f"[AccTest] Uploading to: {url}", file=sys.stderr)
    resp = multipart_upload(url, args.namespace, args.resource, zip_buf,
                            test_case=getattr(args, "test_case", None))

    if resp.get("code") != "SUCCESS":
        print(json.dumps(resp, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    task_id = resp["data"]
    result = {"taskId": task_id, "status": "submitted", "zipSize": zip_size}
    print(json.dumps(result, ensure_ascii=False))
    return task_id


def cmd_upload_run(args):
    """Full workflow: upload local code -> poll until done -> download logs."""
    # 1. Upload
    dir_path = args.dir
    validate_provider_dir(dir_path)
    print(f"[AccTest] Zipping directory: {dir_path}", file=sys.stderr)
    zip_buf = create_code_zip(dir_path)
    zip_size = zip_buf.getbuffer().nbytes
    print(f"[AccTest] Zip created: {zip_size} bytes", file=sys.stderr)

    url = f"{args.base_url}/api/v1/terraform_acc_test/upload"
    print(f"[AccTest] Uploading to: {url}", file=sys.stderr)
    resp = multipart_upload(url, args.namespace, args.resource, zip_buf,
                            test_case=getattr(args, "test_case", None))

    if resp.get("code") != "SUCCESS":
        print(json.dumps(resp, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    task_id = resp["data"]
    print(f"[AccTest] Task submitted: taskId={task_id}", file=sys.stderr)

    # 2. Poll with live log streaming
    # Each poll fetches recent log entries via logSinceTimestamp cursor; prints new lines as they appear.
    # Tracks staleness (time since last new log entry); emits an explicit STALL marker after 5 minutes
    # so the calling agent can decide whether to ask the user about cancelling. The CLI itself never
    # auto-cancels.
    base_status_url = f"{args.base_url}/api/v1/terraform_acc_test/{task_id}"
    poll_count = 0
    # FC task mode supports up to 24h execution time
    max_polls = (24 * 3600) // args.poll_interval
    last_log_ts = ""        # cursor for incremental log fetch
    last_new_log_at = time.time()
    stall_warning_seconds = 5 * 60   # warn agent after 5 minutes of silence
    stall_warned = False

    status = "unknown"
    while poll_count < max_polls:
        time.sleep(args.poll_interval)
        poll_count += 1

        # Build status URL with log tail + incremental cursor
        params = {"logLimit": "200"}
        if last_log_ts:
            params["logSinceTimestamp"] = last_log_ts
        query = "&".join(f"{k}={urllib.parse.quote(v)}" for k, v in params.items())
        poll_url = f"{base_status_url}?{query}"

        status_resp = api_request("GET", poll_url)
        if status_resp.get("code") != "SUCCESS":
            print(f"[AccTest] Poll error: {json.dumps(status_resp)}", file=sys.stderr)
            continue

        data = status_resp["data"]
        status = data.get("status", "unknown")
        fc_status = data.get("fcStatus")
        elapsed = poll_count * args.poll_interval

        # Print any new log entries
        new_logs = data.get("recentLogs") or []
        for log in new_logs:
            ts = log.get("timestamp", "")
            level = log.get("level", "")
            stage = log.get("stage", "")
            message = log.get("message", "")
            print(f"  {ts} [{level}] {stage}: {message}", file=sys.stderr)

        if new_logs:
            last_log_ts = data.get("logLastTimestamp") or last_log_ts
            last_new_log_at = time.time()
            stall_warned = False

        # Concise heartbeat line so absent-log polls are still visible
        print(f"[AccTest] [{elapsed}s] status={status} fcStatus={fc_status} "
              f"totalLogs={data.get('totalLogs')} progress={data.get('logProgress')}",
              file=sys.stderr)

        # Stall detection — emit a clear marker for the agent to act on
        silent_for = time.time() - last_new_log_at
        if silent_for >= stall_warning_seconds and not stall_warned:
            print(
                f"[AccTest] STALL_DETECTED taskId={task_id} silent_seconds={int(silent_for)} "
                f"status={status} fcStatus={fc_status} -- "
                f"no new log entries for {int(silent_for // 60)}+ minutes. "
                f"Agent: ask the user whether to cancel (use 'acctest.py cancel --task-id {task_id}') "
                f"or keep waiting. Do NOT cancel without user confirmation.",
                file=sys.stderr,
            )
            stall_warned = True

        if status in ("success", "failed"):
            break
    else:
        print("[AccTest] Timeout: task did not complete within 24 hours", file=sys.stderr)
        result = {"taskId": task_id, "status": "timeout"}
        print(json.dumps(result, ensure_ascii=False))
        sys.exit(1)

    # 3. Get log URLs
    logs_url = f"{args.base_url}/api/v1/terraform_acc_test/{task_id}/logs"
    logs_resp = api_request("GET", logs_url)
    log_urls = logs_resp.get("data", {}) if logs_resp.get("code") == "SUCCESS" else {}

    # 4. Download logs
    download_dir = args.download_dir
    os.makedirs(download_dir, exist_ok=True)

    downloaded = {}
    for key, dl_url in log_urls.items():
        if not dl_url:
            continue
        filename = "run.log" if key == "runLog" else "tf-debug.log"
        dest = os.path.join(download_dir, f"{task_id}_{filename}")
        size = download_file(dl_url, dest)
        if size >= 0:
            downloaded[key] = {"path": dest, "size": size}
            print(f"[AccTest] Downloaded {filename}: {dest} ({size} bytes)", file=sys.stderr)

    # 5. Output result
    result = {
        "taskId": task_id,
        "status": status,
        "logs": downloaded,
    }
    print(json.dumps(result, indent=2, ensure_ascii=False))

    return 0 if status == "success" else 1


# ── CLI Parser ────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(
        description="Terraform AccTest CLI - ACube AccTest API client",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL,
                        help=f"ACube API base URL (default: {DEFAULT_BASE_URL})")

    sub = parser.add_subparsers(dest="command", help="Available commands")

    # status
    p_status = sub.add_parser("status", help="Query task status (incl. recent log tail)")
    p_status.add_argument("--task-id", required=True, help="Task ID")

    # logs
    p_logs = sub.add_parser("logs", help="Get log download URLs (optionally download)")
    p_logs.add_argument("--task-id", required=True, help="Task ID")
    p_logs.add_argument("--download-dir", default=None,
                        help="If set, download log files to this directory")

    # download (single file from presigned URL)
    p_dl = sub.add_parser("download", help="Download a file from presigned OSS URL")
    p_dl.add_argument("--url", required=True, help="Presigned OSS download URL")
    p_dl.add_argument("--download-dir", default=DEFAULT_DOWNLOAD_DIR,
                      help=f"Directory to save file (default: {DEFAULT_DOWNLOAD_DIR})")

    # cancel
    p_cancel = sub.add_parser("cancel", help="Cancel a running AccTest task (FC StopAsyncTask + DB FAILED)")
    p_cancel.add_argument("--task-id", required=True, help="Task ID")

    # upload
    p_upload = sub.add_parser("upload", help="Upload local code zip and submit AccTest task")
    p_upload.add_argument("--namespace", required=True, help="Product namespace (e.g. VPC, ECS)")
    p_upload.add_argument("--resource", required=True, help="Resource type code (e.g. VSwitch, Instance)")
    p_upload.add_argument("--dir", required=True, help="Path to terraform-provider-alicloud directory")
    p_upload.add_argument("--test-case", default=None,
                          help="Optional: run a single test case by exact name "
                               "(e.g. TestAccAliCloudVPCVSwitch_basic). If unset, runs all matching cases.")

    # upload-run (all-in-one from local code)
    p_upload_run = sub.add_parser("upload-run", help="Full workflow: upload local code -> poll -> download logs")
    p_upload_run.add_argument("--namespace", required=True, help="Product namespace (e.g. VPC, ECS)")
    p_upload_run.add_argument("--resource", required=True, help="Resource type code (e.g. VSwitch, Instance)")
    p_upload_run.add_argument("--dir", required=True, help="Path to terraform-provider-alicloud directory")
    p_upload_run.add_argument("--test-case", default=None,
                               help="Optional: run a single test case by exact name")
    p_upload_run.add_argument("--poll-interval", type=int, default=DEFAULT_POLL_INTERVAL,
                               help=f"Poll interval in seconds (default: {DEFAULT_POLL_INTERVAL})")
    p_upload_run.add_argument("--download-dir", default=DEFAULT_DOWNLOAD_DIR,
                               help=f"Directory to save downloaded logs (default: {DEFAULT_DOWNLOAD_DIR})")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    if args.command == "upload":
        cmd_upload(args)
    elif args.command == "status":
        cmd_status(args)
    elif args.command == "logs":
        cmd_logs(args)
    elif args.command == "download":
        cmd_download(args)
    elif args.command == "cancel":
        cmd_cancel(args)
    elif args.command == "upload-run":
        exit_code = cmd_upload_run(args)
        sys.exit(exit_code)


if __name__ == "__main__":
    main()

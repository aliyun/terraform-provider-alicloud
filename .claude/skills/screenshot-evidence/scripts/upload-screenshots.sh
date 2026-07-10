#!/usr/bin/env bash
# upload-screenshots.sh — upload PNG screenshots to OSS and generate signed URLs.
#
# Usage: upload-screenshots.sh <aone-id> <screenshot-dir>
# Output: name|signed_url lines to stdout (one per file)
#
# Requires: aliyun CLI authenticated with OSS access.

set -euo pipefail

die() { echo "upload-screenshots: $*" >&2; exit 1; }

aone_id="${1:?usage: upload-screenshots.sh <aone-id> <screenshot-dir>}"
screenshot_dir="${2:?usage: upload-screenshots.sh <aone-id> <screenshot-dir>}"

[ -d "$screenshot_dir" ] || die "directory not found: $screenshot_dir"

BUCKET="oss://jarvis-report-images"
ENDPOINT="oss-cn-hangzhou.aliyuncs.com"
PREFIX="reports/${aone_id}"
TIMEOUT=15768000  # 6 months

count=0
for f in "$screenshot_dir"/*.png "$screenshot_dir"/*.jpg "$screenshot_dir"/*.jpeg; do
    [ -f "$f" ] || continue
    basename_f="$(basename "$f")"
    oss_path="${BUCKET}/${PREFIX}/${basename_f}"

    # Upload (private ACL)
    aliyun oss cp "$f" "$oss_path" --acl private -e "$ENDPOINT" -f >/dev/null 2>&1 \
        || die "upload failed: $basename_f"

    # Generate signed URL
    signed_url="$(aliyun oss sign "$oss_path" --timeout "$TIMEOUT" -e "$ENDPOINT" 2>/dev/null | head -1)"
    [ -n "$signed_url" ] || die "sign failed: $basename_f"

    # Output: name (without extension)|signed_url
    name="${basename_f%.*}"
    echo "${name}|${signed_url}"
    count=$((count + 1))
done

[ "$count" -gt 0 ] || die "no PNG/JPG files found in $screenshot_dir"
echo "# uploaded $count files" >&2

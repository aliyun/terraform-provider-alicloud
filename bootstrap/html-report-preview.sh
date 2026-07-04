#!/usr/bin/env bash
# bootstrap/html-report-preview.sh - upload HTML reports to AutomationAgent preview storage.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"

load_env_defaults() {
    local env_file="$jarvis_root/bootstrap/.env"
    [ -f "$env_file" ] || return 0

    local line key value
    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in
            ""|\#*) continue ;;
            *=*) ;;
            *) continue ;;
        esac
        key="${line%%=*}"
        value="${line#*=}"
        [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || continue
        if [ -z "${!key+x}" ]; then
            export "$key=$value"
        fi
    done < "$env_file"
}

load_env_defaults

DEFAULT_BASE_URL="${JARVIS_HTML_REPORT_BASE_URL:-https://pre-agent.aliyun-inc.com}"
CURL_BIN="${JARVIS_CURL_BIN:-curl}"
A1_BIN="${JARVIS_A1_BIN:-$jarvis_root/bin/a1id}"

usage() {
    cat <<'EOF'
Usage:
  bootstrap/html-report-preview.sh upload <aone-id> <html|zip|dir> [--base-url URL] [--comment] [--format markdown|jsonl]
  bootstrap/html-report-preview.sh from-aone <aone-id> [--attachment-id ID|--all] [--base-url URL] [--comment] [--format markdown|jsonl]

Defaults:
  --base-url defaults to $JARVIS_HTML_REPORT_BASE_URL or https://pre-agent.aliyun-inc.com
  Set $JARVIS_HTML_REPORT_TOKEN to send Authorization: Bearer <token>.
  from-aone uploads the newest .html/.htm/.zip attachment unless --attachment-id or --all is set.
EOF
}

die() {
    echo "html-report-preview: $*" >&2
    exit 1
}

is_safe_aone_id() {
    [[ "${1:-}" =~ ^[A-Za-z0-9_-]+$ ]]
}

is_html_path() {
    local lower
    lower="$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')"
    [[ "$lower" == *.html || "$lower" == *.htm ]]
}

is_zip_path() {
    local lower
    lower="$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')"
    [[ "$lower" == *.zip ]]
}

join_url() {
    local base="$1"
    local path="$2"
    if [[ "$path" =~ ^https?:// ]]; then
        printf '%s\n' "$path"
    else
        printf '%s/%s\n' "${base%/}" "${path#/}"
    fi
}

a1_run() {
    "$A1_BIN" -- "$@"
}

parse_common_flags() {
    base_url="$DEFAULT_BASE_URL"
    comment=0
    output_format="markdown"
    attachment_id=""
    all_attachments=0
    positional=()

    while [ "$#" -gt 0 ]; do
        case "$1" in
            --base-url)
                [ $# -ge 2 ] || die "--base-url requires a value"
                base_url="$2"
                shift 2
                ;;
            --comment)
                comment=1
                shift
                ;;
            --format)
                [ $# -ge 2 ] || die "--format requires a value"
                output_format="$2"
                case "$output_format" in
                    markdown|jsonl) ;;
                    *) die "--format must be markdown or jsonl" ;;
                esac
                shift 2
                ;;
            --attachment-id)
                [ $# -ge 2 ] || die "--attachment-id requires a value"
                attachment_id="$2"
                shift 2
                ;;
            --all)
                all_attachments=1
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            --)
                shift
                while [ "$#" -gt 0 ]; do positional+=("$1"); shift; done
                ;;
            -*)
                die "unknown flag: $1"
                ;;
            *)
                positional+=("$1")
                shift
                ;;
        esac
    done
}

extract_zip_htmls() {
    local zip_path="$1"
    local out_dir="$2"

    python3 - "$zip_path" "$out_dir" <<'PY'
import os
import posixpath
import re
import sys
import zipfile

zip_path, out_dir = sys.argv[1], sys.argv[2]
count = 0
with zipfile.ZipFile(zip_path) as zf:
    for info in zf.infolist():
        if info.is_dir():
            continue
        name = info.filename
        lower = name.lower()
        if not (lower.endswith(".html") or lower.endswith(".htm")):
            continue
        base = posixpath.basename(name)
        if not base:
            continue
        count += 1
        safe_base = re.sub(r"[^A-Za-z0-9._-]+", "_", base)
        target = os.path.join(out_dir, f"{count:03d}-{safe_base}")
        with zf.open(info) as src, open(target, "wb") as dst:
            dst.write(src.read())
        print(f"{target}\t{base}")
if count == 0:
    sys.exit(3)
PY
}

collect_html_inputs() {
    local input="$1"
    local tmp_dir="$2"

    [ -e "$input" ] || die "input not found: $input"
    if [ -d "$input" ]; then
        local found=0
        while IFS= read -r -d '' file; do
            printf '%s\t%s\n' "$file" "$(basename "$file")"
            found=1
        done < <(find "$input" -type f \( -iname '*.html' -o -iname '*.htm' \) -print0 | sort -z)
        [ "$found" -eq 1 ] || die "directory has no HTML files: $input"
    elif is_html_path "$input"; then
        printf '%s\t%s\n' "$input" "$(basename "$input")"
    elif is_zip_path "$input"; then
        extract_zip_htmls "$input" "$tmp_dir" || die "zip has no HTML files: $input"
    else
        die "only .html, .htm, .zip, or a directory of HTML files are supported: $input"
    fi
}

upload_one() {
    local aone_id="$1"
    local file="$2"
    local label="$3"
    local endpoint="${base_url%/}/api/reports/aone/$aone_id"
    local response
    local curl_args=(-fsS -X POST)

    if [ -n "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
        curl_args+=(-H "Authorization: Bearer $JARVIS_HTML_REPORT_TOKEN")
    fi

    response="$("$CURL_BIN" "${curl_args[@]}" -F "file=@$file;type=text/html" "$endpoint")" \
        || die "upload failed for $label"
    if ! jq -e '.success == true and .data.viewUrl != null' >/dev/null <<<"$response"; then
        die "upload failed for $label: $response"
    fi

    local view_url report_id object_key size absolute_url
    view_url="$(jq -r '.data.viewUrl' <<<"$response")"
    report_id="$(jq -r '.data.reportId // ""' <<<"$response")"
    object_key="$(jq -r '.data.objectKey // ""' <<<"$response")"
    size="$(jq -r '.data.size // ""' <<<"$response")"
    absolute_url="$(join_url "$base_url" "$view_url")"

    if [ "$output_format" = "jsonl" ]; then
        jq -cn \
            --arg label "$label" \
            --arg url "$absolute_url" \
            --arg viewUrl "$view_url" \
            --arg reportId "$report_id" \
            --arg objectKey "$object_key" \
            --arg size "$size" \
            '{label:$label,url:$url,viewUrl:$viewUrl,reportId:$reportId,objectKey:$objectKey,size:$size}'
    else
        printf -- '- %s: %s\n' "$label" "$absolute_url"
    fi

    links+=("$label|$absolute_url")
}

upload_input() {
    local aone_id="$1"
    local input="$2"
    local tmp_dir="$3"
    local count=0

    while IFS=$'\t' read -r file label; do
        [ -n "$file" ] || continue
        upload_one "$aone_id" "$file" "$label"
        count=$((count + 1))
    done < <(collect_html_inputs "$input" "$tmp_dir")

    [ "$count" -gt 0 ] || die "no HTML files uploaded"
}

select_attachments() {
    local attachments_json="$1"

    if [ -n "$attachment_id" ]; then
        jq -c --arg id "$attachment_id" '
          map(select((.id|tostring) == $id)) |
          if length == 0 then error("attachment not found: " + $id) else .[] end
        ' <<<"$attachments_json"
    elif [ "$all_attachments" -eq 1 ]; then
        jq -c '
          map(select((.name // "" | ascii_downcase | test("\\.(html?|zip)$")))) |
          sort_by(.created // "") |
          .[]
        ' <<<"$attachments_json"
    else
        jq -c '
          map(select((.name // "" | ascii_downcase | test("\\.(html?|zip)$")))) |
          sort_by(.created // "") |
          if length == 0 then error("no html or zip attachment found") else .[-1] end
        ' <<<"$attachments_json"
    fi
}

download_and_upload_attachments() {
    local aone_id="$1"
    local tmp_dir="$2"
    local attachments_json attachment_json attachment_name attachment_id_value out_path

    attachments_json="$(a1_run project workitem attachment list "$aone_id" -f json)"
    while IFS= read -r attachment_json; do
        [ -n "$attachment_json" ] || continue
        attachment_id_value="$(jq -r '.id' <<<"$attachment_json")"
        attachment_name="$(jq -r '.name // ("attachment-" + (.id|tostring))' <<<"$attachment_json")"
        out_path="$tmp_dir/$attachment_id_value-$(basename "$attachment_name")"
        a1_run project workitem attachment download "$aone_id" "$attachment_id_value" -o "$out_path" >/dev/null
        upload_input "$aone_id" "$out_path" "$tmp_dir"
    done < <(select_attachments "$attachments_json")
}

comment_links() {
    local aone_id="$1"
    [ "${#links[@]}" -gt 0 ] || return 0

    local message
    message="已上传 HTML 报告，可在线评阅："
    local item label url
    for item in "${links[@]}"; do
        label="${item%%|*}"
        url="${item#*|}"
        message="${message}"$'\n'"- ${label}: ${url}"
    done
    a1_run project workitem comment create "$aone_id" -m "$message"
}

cmd="${1:-}"
[ -n "$cmd" ] || { usage; exit 1; }
shift || true

case "$cmd" in
    upload)
        parse_common_flags "$@"
        [ "${#positional[@]}" -eq 2 ] || die "upload requires <aone-id> <html|zip|dir>"
        aone_id="${positional[0]}"
        input="${positional[1]}"
        ;;
    from-aone)
        parse_common_flags "$@"
        [ "${#positional[@]}" -eq 1 ] || die "from-aone requires <aone-id>"
        aone_id="${positional[0]}"
        input=""
        ;;
    -h|--help)
        usage
        exit 0
        ;;
    *)
        usage >&2
        exit 1
        ;;
esac

is_safe_aone_id "$aone_id" || die "Aone ID format is invalid: $aone_id"

tmp_dir="$(mktemp -d)"
links=()
trap 'rm -rf "$tmp_dir"' EXIT

case "$cmd" in
    upload)
        upload_input "$aone_id" "$input" "$tmp_dir"
        ;;
    from-aone)
        download_and_upload_attachments "$aone_id" "$tmp_dir"
        ;;
esac

if [ "$comment" -eq 1 ]; then
    comment_links "$aone_id"
fi

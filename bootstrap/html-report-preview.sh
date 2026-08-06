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
  $JARVIS_HTML_REPORT_TOKEN is required for every upload.
  from-aone uploads the newest .html/.htm/.zip attachment unless --attachment-id or --all is set.
EOF
}

die() {
    echo "html-report-preview: $*" >&2
    exit 1
}

blocked_missing_token() {
    local message
    message="JARVIS_HTML_REPORT_TOKEN is required; no curl or Aone request was sent"
    if [ "$output_format" = "jsonl" ]; then
        jq -cn \
            --arg message "$message" \
            '{success:false,status:"blocked",code:"missing_token",message:$message}'
    else
        printf 'html-report-preview: blocked (missing_token): %s\n' "$message"
    fi
    exit 3
}

upload_failed() {
    local code="$1"
    if [ "$output_format" = "jsonl" ]; then
        jq -cn \
            --arg code "$code" \
            '{success:false,status:"failed",code:$code,message:"HTML report upload failed"}'
    else
        printf 'html-report-preview: upload failed (%s)\n' "$code" >&2
    fi
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

# Only HTML members are ever uploaded, so images sitting beside the report are
# silently dropped and the report renders with broken/absent screenshots. Refuse
# the input instead: screenshots must become signed URLs via upload-screenshots.sh.
reject_dropped_image_files() {
    local found="$1"

    [ -n "$found" ] || return 0
    echo "html-report-preview: images are dropped, not uploaded with the report: $found" >&2
    echo "html-report-preview: run .claude/skills/screenshot-evidence/scripts/upload-screenshots.sh first, then reference the signed https URLs" >&2
    upload_failed "image_files_dropped"
}

zip_image_members() {
    python3 - "$1" <<'PY'
import posixpath
import sys
import zipfile

SUFFIXES = (".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp")
names = []
try:
    with zipfile.ZipFile(sys.argv[1]) as zf:
        for info in zf.infolist():
            if info.is_dir():
                continue
            if info.filename.lower().endswith(SUFFIXES):
                base = posixpath.basename(info.filename)
                if base:
                    names.append(base)
except (OSError, zipfile.BadZipFile):
    raise SystemExit(0)
print(" ".join(sorted(names)[:5]))
PY
}

collect_html_inputs() {
    local input="$1"
    local tmp_dir="$2"

    [ -e "$input" ] || die "input not found: $input"
    if [ -d "$input" ]; then
        reject_dropped_image_files "$(find "$input" -type f \
            \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' \
               -o -iname '*.webp' -o -iname '*.gif' -o -iname '*.bmp' \) \
            -exec basename {} \; 2>/dev/null | sort | head -5 | tr '\n' ' ')"
    fi
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
        reject_dropped_image_files "$(zip_image_members "$input")"
        extract_zip_htmls "$input" "$tmp_dir" || die "zip has no HTML files: $input"
    else
        die "only .html, .htm, .zip, or a directory of HTML files are supported: $input"
    fi
}

validate_html_image_references() {
    local file="$1"

    python3 - "$file" <<'PY'
import re
import sys
from html.parser import HTMLParser
from urllib.parse import urlsplit

IMAGE_SUFFIX_RE = re.compile(r"\.(?:png|jpe?g|webp|gif|bmp)$", re.IGNORECASE)
# A screenshot filename mentioned in running text, e.g. "截图：openapi_create.png".
# Tokens containing "://" are URLs, not the prose-only pattern we are hunting.
PROSE_IMAGE_RE = re.compile(r"[\w./-]+\.(?:png|jpe?g|webp|gif|bmp)\b", re.IGNORECASE)
# Text inside these elements is sample code/config, not a screenshot caption.
LITERAL_TAGS = {"pre", "code", "script", "style", "textarea"}
# The documented degradation marker from the screenshot-evidence skill: a report
# that explains why screenshots are absent is legitimately image-free.
DEGRADED_MARKERS = ("missing_capability",)


def valid_https_url(value):
    if not value or value != value.strip():
        return False
    if any(char.isspace() or ord(char) < 0x20 for char in value):
        return False
    if "\\" in value:
        return False
    try:
        parsed = urlsplit(value)
        # Accessing port also rejects malformed values such as ":not-a-port".
        parsed.port
    except ValueError:
        return False
    return (
        parsed.scheme == "https"
        and bool(parsed.hostname)
        and parsed.username is None
        and parsed.password is None
    )


def srcset_urls(value):
    if not value or value != value.strip():
        raise ValueError("empty srcset")
    urls = []
    for candidate in value.split(","):
        fields = candidate.strip().split()
        if not fields or len(fields) > 2:
            raise ValueError("invalid srcset candidate")
        if len(fields) == 2 and not re.fullmatch(
            r"(?:[1-9][0-9]*w|(?:[1-9][0-9]*(?:\.[0-9]+)?|0\.[0-9]*[1-9][0-9]*)x)",
            fields[1],
        ):
            raise ValueError("invalid srcset descriptor")
        urls.append(fields[0])
    return urls


class ImageReferenceValidator(HTMLParser):
    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.picture_depth = 0
        self.invalid = False
        # Screenshot-presence tracking: a "visual evidence" report that names its
        # screenshots only in prose renders with nothing to look at, and the
        # reference checks above cannot catch it because there is no ref to check.
        self.img_count = 0
        self.embedded_image_url = False
        self.literal_depth = 0
        self.prose = []

    def validate_url(self, value):
        if not valid_https_url(value or ""):
            self.invalid = True

    def validate_srcset(self, value):
        try:
            urls = srcset_urls(value or "")
        except ValueError:
            self.invalid = True
            return
        if not all(valid_https_url(url) for url in urls):
            self.invalid = True

    def note_embedded_url(self, value):
        for candidate in re.split(r"[\s,]+", value or ""):
            candidate = candidate.strip()
            if candidate.lower().startswith("https://") and IMAGE_SUFFIX_RE.search(
                urlsplit(candidate).path
            ):
                self.embedded_image_url = True

    def inspect(self, tag, attrs):
        tag = tag.lower()
        if tag == "img":
            self.img_count += 1
        for _, value in attrs:
            self.note_embedded_url(value)
        if tag == "img":
            for name, value in attrs:
                name = name.lower()
                if name == "src":
                    self.validate_url(value)
                elif name == "srcset":
                    self.validate_srcset(value)
        elif tag == "source" and self.picture_depth > 0:
            for name, value in attrs:
                if name.lower() == "srcset":
                    self.validate_srcset(value)

    def handle_starttag(self, tag, attrs):
        self.inspect(tag, attrs)
        tag = tag.lower()
        if tag == "picture":
            self.picture_depth += 1
        if tag in LITERAL_TAGS:
            self.literal_depth += 1

    def handle_startendtag(self, tag, attrs):
        self.inspect(tag, attrs)

    def handle_endtag(self, tag):
        tag = tag.lower()
        if tag == "picture" and self.picture_depth > 0:
            self.picture_depth -= 1
        if tag in LITERAL_TAGS and self.literal_depth > 0:
            self.literal_depth -= 1

    def handle_data(self, data):
        if self.literal_depth == 0:
            self.prose.append(data)

    def prose_only_screenshot(self):
        """Report names screenshot files in prose yet embeds no image at all."""
        if self.img_count or self.embedded_image_url:
            return False
        text = "".join(self.prose)
        if any(marker in text for marker in DEGRADED_MARKERS):
            return False
        return any(
            "://" not in match.group(0) for match in PROSE_IMAGE_RE.finditer(text)
        )


try:
    with open(sys.argv[1], "r", encoding="utf-8") as handle:
        document = handle.read()
    validator = ImageReferenceValidator()
    validator.feed(document)
    validator.close()
except (OSError, UnicodeError, ValueError):
    raise SystemExit(1)

if validator.invalid:
    raise SystemExit(1)
raise SystemExit(4 if validator.prose_only_screenshot() else 0)
PY
}

upload_one() {
    local aone_id="$1"
    local file="$2"
    local label="$3"
    local endpoint="${base_url%/}/api/reports/aone/$aone_id"
    local response
    local curl_args=(-fsS -X POST)

    curl_args+=(-H "Authorization: Bearer $JARVIS_HTML_REPORT_TOKEN")

    # WAF classification gate requires X-Request-Context header (pre-agent rgv587 Anti-Bot)
    local waf_header="${JARVIS_HTML_REPORT_WAF_HEADER:-rctx_a3f90b7e2d41c8f6}"
    curl_args+=(-H "X-Request-Context: $waf_header")

    if ! response="$("$CURL_BIN" "${curl_args[@]}" \
        -F "file=@$file;type=text/html" "$endpoint" 2>/dev/null)"; then
        upload_failed "transport_error"
    fi
    if ! jq -e '.success == true and .data.viewUrl != null' \
        >/dev/null 2>&1 <<<"$response"; then
        upload_failed "upload_rejected"
    fi

    local view_url report_id object_key size absolute_url
    view_url="$(jq -r '.data.viewUrl' <<<"$response")"
    report_id="$(jq -r '.data.reportId // ""' <<<"$response")"
    object_key="$(jq -r '.data.objectKey // ""' <<<"$response")"
    size="$(jq -r '.data.size // ""' <<<"$response")"
    absolute_url="$(join_url "$base_url" "$view_url")"
    if [ -z "$report_id" ] ||
       [ "$view_url" != "/reports/aone/$aone_id/$report_id/view" ]; then
        upload_failed "invalid_view_route"
    fi

    if [ "$output_format" = "jsonl" ]; then
        jq -cn \
            --arg label "$label" \
            --arg url "$absolute_url" \
            --arg viewUrl "$view_url" \
            --arg reportId "$report_id" \
            --arg objectKey "$object_key" \
            --arg size "$size" \
            '{success:true,status:"uploaded",label:$label,url:$url,viewUrl:$viewUrl,reportId:$reportId,objectKey:$objectKey,size:$size}'
    else
        printf -- '- %s: %s\n' "$label" "$absolute_url"
    fi

    links+=("$label|$absolute_url")
}

upload_input() {
    local aone_id="$1"
    local input="$2"
    local tmp_dir="$3"
    local inputs_file="$tmp_dir/html-inputs.tsv"
    local count=0

    collect_html_inputs "$input" "$tmp_dir" >"$inputs_file"

    # Validate the entire batch before the first upload so a ZIP/directory cannot
    # partially publish valid entries before a later unsafe report is rejected.
    while IFS=$'\t' read -r file label; do
        [ -n "$file" ] || continue
        local check_rc=0
        validate_html_image_references "$file" || check_rc=$?
        case "$check_rc" in
            0) ;;
            4) upload_failed "screenshot_prose_without_img" ;;
            *) upload_failed "invalid_image_reference" ;;
        esac
    done <"$inputs_file"

    while IFS=$'\t' read -r file label; do
        [ -n "$file" ] || continue
        upload_one "$aone_id" "$file" "$label"
        count=$((count + 1))
    done <"$inputs_file"

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

    # Aone 评论区按 markdown 渲染:可点击链接唯一可靠格式 = [text](url);
    # 裸 URL(独行/行内)不 autolink,<a href>/<url> 被当 HTML tag 剥掉——
    # 先例:84307546 评论 124870464 四格式对照实测(aone-triage SKILL §4 quirk)。
    local message
    message="已上传 HTML 报告，可在线评阅："
    local item label url
    for item in "${links[@]}"; do
        label="${item%%|*}"
        url="${item#*|}"
        message="${message}"$'\n'"- [${label}](${url})"
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
[ -n "${JARVIS_HTML_REPORT_TOKEN:-}" ] || blocked_missing_token

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

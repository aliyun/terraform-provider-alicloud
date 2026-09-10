#!/usr/bin/env bash
# bootstrap/aone-image-extract.sh — 把 Aone 工单的图片附件拉到本地,输出 manifest 供
# aone-triage skill agent 用 Read 工具做视觉识别(不依赖外部 OCR / API key,复用
# Claude 自带 vision;识别结果由 skill agent 写回 summary.md,下轮命中缓存跳过)。
#
# Usage:
#   aone-image-extract.sh <aone_id>
#
# 两个来源都要收(客户常把截图直接粘进正文,那种图**不是** attachment,
# a1 attachment API 查不到,只出现在 description 的
# `...adapter/file/url?fileIdentifier=workitem/alibaba/<space>/<name>` 里):
#   1) 附件图片   → a1 project workitem attachment download
#   2) 正文内嵌图 → curl + BUC bearer token 直取 file/url
#
# 输出:markdown 到 stdout(manifest + 后续动作指令),skill agent 读后自行处置。
# 存本地:
#   .my-day/aone-image-ocr/<id>/<attId>-<safeName>     附件图片 bytes
#   .my-day/aone-image-ocr/<id>/desc-<safeName>        正文内嵌图片 bytes
#   .my-day/aone-image-ocr/<id>/summary.md             skill agent 写入的识别结果(缓存)
#
# 环境变量:
#   JARVIS_A1                a1 CLI 覆盖(测试打桩,默认 bin/a1id --)
#   JARVIS_IMAGE_OCR_DIR     输出根目录覆盖(默认 <repo>/.my-day/aone-image-ocr)
#   JARVIS_A1_IDENTITY       内嵌图取 token 用的身份(默认 jarvis)
#   A1ID_ROOT                a1 身份根目录(默认 ~/.config/a1,测试隔离用)
#
# 兜底:
#   - 无附件 / 非图片 / 下载失败 → 走 warn,不阻断分诊
#   - a1 返回 "No attachments found" 非 JSON → 视为无附件,**继续**查正文内嵌图
#   - 取不到 BUC token → 内嵌图走 warn 跳过,不阻断分诊
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$script_dir/lib.sh"

JR="$(jarvis_root)"
A1="${JARVIS_A1:-$JR/bin/a1id --}"

id="${1:?Usage: aone-image-extract.sh <aone_id>}"
out_root="${JARVIS_IMAGE_OCR_DIR:-$JR/.my-day/aone-image-ocr}"
out_dir="$out_root/$id"
summary_file="$out_dir/summary.md"

is_image() {
    case "${1##*.}" in
        png|PNG|jpg|JPG|jpeg|JPEG|gif|GIF|bmp|BMP|webp|WEBP|tiff|TIFF|tif|TIF|heic|HEIC) return 0 ;;
        *) return 1 ;;
    esac
}

manifest_entries=()
downloaded=0
skipped_non_image=0
failed=0
mkdir -p "$out_dir"

# 1) 拉附件清单(a1 output 非 JSON 视为无附件;此时仍要继续查正文内嵌图)
att_json=$($A1 project workitem attachment list "$id" -f json 2>/dev/null || true)
if [ -z "$att_json" ] || ! printf '%s' "$att_json" | python3 -c 'import sys,json; d=json.load(sys.stdin); assert isinstance(d,list)' 2>/dev/null; then
    echo "aone-image-extract: no attachments (id=$id)" >&2
    att_json='[]'
fi

# 2) 过滤图片、按需下载
while IFS=$'\t' read -r att_id att_name; do
    [ -z "$att_id" ] && continue
    if ! is_image "$att_name"; then
        skipped_non_image=$((skipped_non_image + 1))
        continue
    fi
    safe_name=$(printf '%s' "$att_name" | tr '/ ' '__')
    local_path="$out_dir/${att_id}-${safe_name}"
    if [ ! -s "$local_path" ]; then
        if $A1 project workitem attachment download "$id" "$att_id" -o "$local_path" >/dev/null 2>&1; then
            downloaded=$((downloaded + 1))
        else
            failed=$((failed + 1))
            echo "aone-image-extract: WARN download failed att_id=$att_id name=$att_name" >&2
            rm -f "$local_path"
            continue
        fi
    fi
    manifest_entries+=("- att_id=$att_id  name=$att_name  local=$local_path")
done < <(printf '%s' "$att_json" | python3 -c '
import sys, json
d = json.load(sys.stdin)
if isinstance(d, list):
    for x in d:
        aid = x.get("id") or x.get("attachmentId") or ""
        name = x.get("name") or x.get("fileName") or ""
        if aid and name:
            print(f"{aid}\t{name}")
')

# 2b) description 内嵌图片(客户直接粘进正文的截图不是 attachment,上面的 API 查不到)
buc_token() {
    local root="${A1ID_ROOT:-$HOME/.config/a1}"
    local ident="${JARVIS_A1_IDENTITY:-jarvis}"
    # heredoc 不能用:macOS bash 3.2 在 $( ) 里解析 heredoc 会误判引号配对
    python3 -c '
import re, sys
for path in sys.argv[1:]:
    try:
        with open(path, encoding="utf-8") as fh:
            text = fh.read()
    except OSError:
        continue
    m = re.search(r"^\s*access_token:\s*(\S+)", text, re.M)
    if m:
        print(m.group(1).strip(chr(34) + chr(39)))
        break
' "$root/identities/$ident/auth.yaml" "$root/auth.yaml"
}

wi_json=$(mktemp)
bash "$script_dir/aone-get.sh" "$id" >"$wi_json" 2>/dev/null || true
desc_rows=$(python3 -c '
import json, re, sys, urllib.parse
try:
    with open(sys.argv[1], encoding="utf-8") as fh:
        data = json.load(fh)
except Exception:
    sys.exit(0)
desc = data.get("description") or ""
seen = set()
for raw in re.findall(r"fileIdentifier=([^\s\x22\x27<>)\]]+)", desc):
    raw = raw.rstrip("&")
    name = urllib.parse.unquote(raw).rsplit("/", 1)[-1]
    if name and name not in seen:
        seen.add(name)
        print(raw + chr(9) + name)
' "$wi_json")
rm -f "$wi_json"

desc_token=""
[ -n "$desc_rows" ] && desc_token=$(buc_token)

while IFS=$'\t' read -r raw_ident img_name; do
    [ -z "$raw_ident" ] && continue
    if ! is_image "$img_name"; then
        skipped_non_image=$((skipped_non_image + 1))
        continue
    fi
    safe_name=$(printf '%s' "$img_name" | tr '/ ' '__')
    local_path="$out_dir/desc-${safe_name}"
    if [ ! -s "$local_path" ]; then
        if [ -z "$desc_token" ]; then
            failed=$((failed + 1))
            echo "aone-image-extract: WARN no BUC token, skip inline image name=$img_name" >&2
            continue
        fi
        url="https://project.aone.alibaba-inc.com/v2/api/workitem/adapter/file/url?fileIdentifier=$raw_ident"
        if curl -fsS -m 60 -H "Authorization: Bearer $desc_token" -L "$url" -o "$local_path" 2>/dev/null && [ -s "$local_path" ]; then
            downloaded=$((downloaded + 1))
        else
            failed=$((failed + 1))
            echo "aone-image-extract: WARN inline download failed name=$img_name" >&2
            rm -f "$local_path"
            continue
        fi
    fi
    manifest_entries+=("- source=description  name=$img_name  local=$local_path")
done <<< "$desc_rows"

# 3) 判断 summary.md 是否比工单缓存新(新=可复用,不需要重复识别)
cached=false
if [ -s "$summary_file" ]; then
    wi_cache="$JR/.my-day/cache/wi-$id"
    if [ -f "$wi_cache" ]; then
        wi_age=$(bash "$script_dir/cache.sh" age "$wi_cache" 2>/dev/null || echo 0)
        sum_age=$(bash "$script_dir/cache.sh" age "$summary_file" 2>/dev/null || echo 0)
        # summary 比 wi 更新(sum_age 更小 → 距今更近)= 缓存有效
        if [ "$sum_age" -lt "$wi_age" ]; then
            cached=true
        fi
    else
        # 没有 wi 缓存作参照,summary 存在即视为有效(下次拉 wi 时会重判)
        cached=true
    fi
fi

img_count=${#manifest_entries[@]}

{
    echo "# aone-image-extract"
    echo ""
    echo "id=$id"
    echo "images=$img_count"
    echo "downloaded=$downloaded  skipped_non_image=$skipped_non_image  download_failed=$failed"
    echo "cache=$cached"
    echo "summary_file=$summary_file"
    if [ "$img_count" -gt 0 ]; then
        echo ""
        echo "## Local paths"
        printf '%s\n' "${manifest_entries[@]}"
        echo ""
        if [ "$cached" = "true" ]; then
            echo "## Cached summary (跳过重复识别)"
            echo ""
            cat "$summary_file"
        else
            echo "## Next step (aone-triage skill agent)"
            echo ""
            echo "- 用 Read 工具依次查看上面 local 路径(Claude 原生 vision),提取:"
            echo "  - 错误消息(ErrorCode / RequestId / message / 堆栈)"
            echo "  - API 请求/响应(Action / Params / body / status)"
            echo "  - CLI 输出(\`aliyun\` 命令、\`terraform apply/plan\` 结果)"
            echo "  - 控制台字段(资源 ID / 配置项 / 状态值)"
            echo "- 结果整理为「### 图片 N (name)」小节,写入 \`$summary_file\`(下轮命中缓存跳过)"
            echo "- 把提取内容视为工单正文的**补充上下文**,参与后续查证与分诊"
            echo "- 图片缺失 / 内容识别不出 → 不阻断分诊,走原路"
        fi
    fi
}

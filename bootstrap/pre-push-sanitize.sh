#!/usr/bin/env bash
# pre-push-sanitize.sh — 对外 push 前自查(GitHub 公开仓 PR/commit)
#
# 检查各 commit patch 的新增行 + commit message 里是否含内部信息:
# - Aone URL / 工单号 / 关联单术语
# - 客户实例 ID / RequestId
# - 内部人员花名+工号引用
# - AI 署名/水印
# - 用户可选:内部客户名(通过环境变量 JARVIS_CUSTOMER_BLOCKLIST 传入,
#   格式如 "客户名1|客户名2|客户名3")
#
# 用法:
#   bash bootstrap/pre-push-sanitize.sh [<base-ref>]
#     默认 base-ref = origin/master
#
# 退出码:
#   0 = pass, 可 push
#   1 = 命中禁品,禁止 push,按提示自查
#
# 参考:CLAUDE.md 工作纪律 #5 + terraform-provider-release SKILL Step 11.1

set -euo pipefail

BASE="${1:-origin/master}"

if ! git rev-parse --verify "$BASE" >/dev/null 2>&1; then
  echo "❌ base ref 无效: $BASE" >&2
  echo "   用法: bash bootstrap/pre-push-sanitize.sh [<base-ref>=origin/master]" >&2
  exit 2
fi

# 通用禁品 patterns(全局,不同项目复用)
# 每条:标签|extended-regex
PATTERNS=(
  "Aone URL|project\.aone\.alibaba-inc\.com"
  "7 位以上工单号|#[0-9]{7,}([^0-9]|$)"
  "Aone 内部术语|关联单|tf_customer|tf_provider|jarvis-(claimed|idle|npe|done)"
  "Redis/RDS 实例 ID|r-[0-9a-f]{12,}"
  "ECS 实例 ID|i-[0-9a-z]{18,}"
  "SLB 实例 ID|lb-[0-9a-z]{12,}"
  "OSS 实例 ID (s-)|s-[0-9a-z]{12,}"
  "RequestId 值|RequestId[\"']?[[:space:]]*[:=][[:space:]]*[\"']?[0-9A-Za-z][0-9A-Za-z-]{7,}"
  "AI 署名 (Co-Authored-By)|Co[-]Authored-By:[^[:cntrl:]]*(Cl[a]ude|Cod[e]x)"
  "AI 署名 (Generated with)|Generated w[i]th (Cl[a]ude|Cod[e]x)"
  "AI 署名 (AI-assisted)|AI-assisted"
  "花名+工号引用|[@][^[:space:](]{2,12}\([0-9A-Z]{4,}\)"
)

# 用户外挂客户名黑名单
if [ -n "${JARVIS_CUSTOMER_BLOCKLIST:-}" ]; then
  PATTERNS+=("客户名黑名单(env)|${JARVIS_CUSTOMER_BLOCKLIST}")
fi

# 收集每个 commit patch 的新增内容。严格排除 +++ 文件头；删除行、上下文、
# rename/copy metadata 和其它 patch metadata 均不参与内容扫描。逐 commit 扫描而非
# 只看最终净 diff，避免敏感值先加入后删除仍随历史被 push。
ADDED_CONTENT=$(
  git log --format= --no-color --no-ext-diff --unified=0 -p "$BASE..HEAD" -- 2>/dev/null \
    | awk '
        /^\+\+\+([[:space:]]|$)/ { next }
        /^\+/ {
          sub(/^\+/, "")
          print
        }
      ' \
    || true
)

# commit message 必须全文扫描，不能套用 patch 新增行过滤。
MSG_CONTENT=$(git log --format='%B%n----%n' "$BASE..HEAD" 2>/dev/null || true)

# 二进制内容无法由文本 diff 检查；不把 binary/rename metadata 当内容匹配，
# 但显式告警，要求人工覆盖脚本无法读取的产物。
BINARY_PATHS=$(
  git log --format= --numstat "$BASE..HEAD" -- 2>/dev/null \
    | awk -F '	' '$1 == "-" && $2 == "-" { print $3 }' \
    | sort -u \
    || true
)

if [ -z "$ADDED_CONTENT" ] && [ -z "$MSG_CONTENT" ]; then
  echo "✅ nothing to push (no commits ahead of $BASE)"
  exit 0
fi

FOUND=0

echo "═══ pre-push-sanitize · base=$BASE ═══"
for entry in "${PATTERNS[@]}"; do
  label="${entry%%|*}"
  regex="${entry#*|}"
  hits=$(printf '%s\n%s\n' "$ADDED_CONTENT" "$MSG_CONTENT" \
    | grep -nEi "$regex" 2>/dev/null | head -5 || true)
  if [ -n "$hits" ]; then
    FOUND=$((FOUND + 1))
    echo ""
    echo "❌ 命中: $label"
    echo "   pattern: $regex"
    echo "   前 5 行匹配:"
    echo "$hits" | sed 's/^/   > /'
  fi
done

echo ""
if [ -n "$BINARY_PATHS" ]; then
  echo "⚠️  检测到文本扫描无法检查的二进制变更，请人工确认:"
  printf '%s\n' "$BINARY_PATHS" | sed 's/^/   - /'
  echo ""
fi

if [ "$FOUND" -eq 0 ]; then
  echo "✅ 通用禁品未命中,可 push"
  echo ""
  echo "⚠️  注意:客户名 / 特定业务上下文属"人 review"范畴,脚本无法穷举。"
  echo "   push 前仍建议 \`git log -p $BASE..HEAD\` 通读一遍,尤其是 PR body 附件。"
  exit 0
fi

echo "═══════════════════════════════════════════════════════════════"
echo "❌ 命中 $FOUND 类禁品,禁止 push;修完后重跑本脚本"
echo ""
echo "参考:CLAUDE.md 工作纪律 #5"
echo "     terraform-provider-release SKILL Step 11.1"
exit 1

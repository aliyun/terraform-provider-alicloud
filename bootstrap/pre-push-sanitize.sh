#!/usr/bin/env bash
# pre-push-sanitize.sh — 对外 push 前自查(GitHub 公开仓 PR/commit)
#
# 检查 diff + commit message 里是否含内部信息:
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
  "7 位以上工单号|#[0-9]{7,}[^0-9]"
  "Aone 内部术语|关联单|tf_customer|tf_provider|jarvis-(claimed|idle|npe|done)"
  "Redis/RDS 实例 ID|r-[0-9a-f]{12,}"
  "ECS 实例 ID|i-[0-9a-z]{18,}"
  "SLB 实例 ID|lb-[0-9a-z]{12,}"
  "OSS 实例 ID (s-)|s-[0-9a-z]{12,}"
  "RequestId 字面量|RequestId"
  "AI 署名 (Co-Authored-By)|Co-Authored-By: Claude"
  "AI 署名 (Generated with)|Generated with Claude"
  "AI 署名 (AI-assisted)|AI-assisted"
  "花名+工号引用|@[^\s(]{2,6}\([0-9A-Z]{4,}\)"
)

# 用户外挂客户名黑名单
if [ -n "${JARVIS_CUSTOMER_BLOCKLIST:-}" ]; then
  PATTERNS+=("客户名黑名单(env)|${JARVIS_CUSTOMER_BLOCKLIST}")
fi

# 收集 diff + 每个 commit message 的完整文本
DIFF_CONTENT=$(git log -p "$BASE..HEAD" 2>/dev/null || true)
MSG_CONTENT=$(git log --format='%B%n----%n' "$BASE..HEAD" 2>/dev/null || true)

if [ -z "$DIFF_CONTENT" ] && [ -z "$MSG_CONTENT" ]; then
  echo "✅ nothing to push (no commits ahead of $BASE)"
  exit 0
fi

FOUND=0

echo "═══ pre-push-sanitize · base=$BASE ═══"
for entry in "${PATTERNS[@]}"; do
  label="${entry%%|*}"
  regex="${entry#*|}"
  hits=$(printf '%s\n%s\n' "$DIFF_CONTENT" "$MSG_CONTENT" \
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

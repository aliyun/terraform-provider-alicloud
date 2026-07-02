#!/usr/bin/env bash
# test/aone_comment_format_test.sh — unit tests for Aone-safe comment formatting.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/../bootstrap" && pwd)"
FORMAT="$BOOTSTRAP_DIR/aone-comment-format.sh"

pass_count=0
fail_count=0

assert_contains() {
    local desc="$1"
    local haystack="$2"
    local needle="$3"
    if [[ "$haystack" == *"$needle"* ]]; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc"
        printf '  expected to find: %q\n' "$needle"
        printf '  actual: %q\n' "$haystack"
        fail_count=$((fail_count + 1))
    fi
}

assert_not_contains() {
    local desc="$1"
    local haystack="$2"
    local needle="$3"
    if [[ "$haystack" != *"$needle"* ]]; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc"
        printf '  did not expect: %q\n' "$needle"
        printf '  actual: %q\n' "$haystack"
        fail_count=$((fail_count + 1))
    fi
}

echo "Test 1: inline Chinese numbered problems become separate list lines"
inline="结论：生成器已修复；剩余问题：1）测试用例缺失；2）业务语义缺失；3）文档说明缺失。"
formatted="$(printf '%s' "$inline" | bash "$FORMAT")"
assert_contains "keeps conclusion text" "$formatted" "结论：生成器已修复"
assert_contains "splits first problem with Aone-safe marker" "$formatted" $'\n\n1、测试用例缺失\n\n'
assert_contains "splits second problem with Aone-safe marker" "$formatted" $'\n\n2、业务语义缺失\n\n'
assert_contains "splits third problem with Aone-safe marker" "$formatted" $'\n\n3、文档说明缺失。'
assert_not_contains "does not keep compact 1） marker" "$formatted" "1）测试用例缺失"
assert_not_contains "does not emit Markdown ordered list marker" "$formatted" $'\n1. 测试用例缺失'

echo "Test 2: markdown headings, emphasis, inline code, and fenced code blocks are preserved"
markdown="$(cat <<'EOF'
## 结论

**重点：标题和加粗可以保留 Markdown。**

提示：不要把 `##` 和 `**` 去掉。

- 第一条问题
- 第二条问题

```go
currentStatus := object["Status"]
if currentStatus != "Accepted" {
    return NotFound
}
```
EOF
)"
formatted_markdown="$(printf '%s' "$markdown" | bash "$FORMAT")"
assert_contains "heading preserved" "$formatted_markdown" "## 结论"
assert_contains "bold emphasis preserved" "$formatted_markdown" "**重点：标题和加粗可以保留 Markdown。**"
assert_contains "inline code preserved" "$formatted_markdown" '`##`'
assert_contains "first bullet converted with hard paragraph break" "$formatted_markdown" $'\n\n· 第一条问题\n\n'
assert_contains "second bullet converted with hard paragraph break" "$formatted_markdown" $'\n\n· 第二条问题\n\n'
assert_contains "fenced code block preserved" "$formatted_markdown" $'```go\ncurrentStatus := object["Status"]\nif currentStatus != "Accepted" {\n    return NotFound\n}\n```'
assert_not_contains "does not convert headings to plaintext markers" "$formatted_markdown" "【结论】"
assert_not_contains "does not emit UI-collapsible bullet glyph" "$formatted_markdown" "• 第一条问题"
assert_not_contains "does not keep Markdown bullet marker" "$formatted_markdown" $'\n- 第一条问题'

echo "Test 3: mentions and footer stay readable"
with_footer=$'@谜拟 请看下\n代码：jarvis @ codex/aone-comment-format (abc123)'
formatted_footer="$(printf '%s' "$with_footer" | bash "$FORMAT")"
assert_contains "mention preserved" "$formatted_footer" "@谜拟"
assert_contains "footer separated" "$formatted_footer" $'\n\n代码：jarvis @ codex/aone-comment-format (abc123)'

echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi

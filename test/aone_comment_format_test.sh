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

echo "Test 4: bare URLs become clickable markdown links (autolink gate)"
bare=$'PR: https://github.com/aliyun/terraform-provider-alicloud/pull/9977 已提交\n见 https://example.com/path'
formatted_bare="$(printf '%s' "$bare" | bash "$FORMAT")"
assert_contains "bare GitHub URL wrapped" "$formatted_bare" '[https://github.com/aliyun/terraform-provider-alicloud/pull/9977](https://github.com/aliyun/terraform-provider-alicloud/pull/9977)'
assert_contains "bare plain URL wrapped" "$formatted_bare" '[https://example.com/path](https://example.com/path)'
assert_not_contains "no leftover unwrapped GitHub URL" "$formatted_bare" 'PR: https://github.com'

# existing markdown link must NOT be double-wrapped
link=$'已有 [文本](https://keep.com) 不重包'
formatted_link="$(printf '%s' "$link" | bash "$FORMAT")"
assert_contains "existing link display text intact" "$formatted_link" '[文本](https://keep.com)'
assert_not_contains "existing link not double-wrapped" "$formatted_link" '[https://keep.com](https://keep.com)'

# URLs inside inline code + fenced code must NOT be touched
code=$'代码 `https://incode.com` 不动\n```\nhttps://in-fence.com 不动\n```'
formatted_code="$(printf '%s' "$code" | bash "$FORMAT")"
assert_not_contains "inline-code URL not wrapped" "$formatted_code" '[https://incode.com]'
assert_not_contains "fenced-code URL not wrapped" "$formatted_code" '[https://in-fence.com]'

# URL stops at CJK punctuation (not absorbed into the link target)
cjk=$'尾句 https://end.com。另一条 https://a.com，还有 https://b.com；'
formatted_cjk="$(printf '%s' "$cjk" | bash "$FORMAT")"
assert_contains "URL stops before CJK full stop" "$formatted_cjk" '[https://end.com](https://end.com)。'
assert_contains "URL stops before CJK comma" "$formatted_cjk" '[https://a.com](https://a.com)，'
assert_not_contains "CJK text not absorbed into URL" "$formatted_cjk" '[https://end.com。另一条]'

# kill switch disables autolink
off="$(printf '%s' '裸 https://off.com' | JARVIS_COMMENT_URL_AUTOLINK=0 bash "$FORMAT")"
assert_not_contains "kill switch leaves URL bare" "$off" '[https://off.com]'
assert_contains "kill switch preserves raw URL text" "$off" 'https://off.com'

echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi

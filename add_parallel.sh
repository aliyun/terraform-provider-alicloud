#!/bin/bash

# 脚本：给所有测试函数添加 t.Parallel()
# 用法：./add_parallel.sh

FILE="alicloud/resource_alicloud_cs_kubernetes_node_pool_test.go"

# 检查文件是否存在
if [ ! -f "$FILE" ]; then
    echo "错误: 文件 $FILE 不存在"
    exit 1
fi

echo "正在处理文件: $FILE"

# 创建临时文件
TMPFILE=$(mktemp)
COUNT=0

# 使用 awk 处理文件
awk '
BEGIN {
    in_test_func = 0
    has_parallel = 0
    func_line = ""
}

# 匹配测试函数定义
/^func Test.*\(t \*testing\.T\) \{/ {
    # 如果之前有测试函数没有处理完，先输出
    if (in_test_func && !has_parallel) {
        print func_line
        print "	t.Parallel()"
        COUNT++
    } else if (in_test_func) {
        print func_line
    }
    
    # 开始新的测试函数
    in_test_func = 1
    has_parallel = 0
    func_line = $0
    next
}

# 检查是否已经有 t.Parallel()
/^\s*t\.Parallel\(\)/ {
    if (in_test_func) {
        has_parallel = 1
    }
    print $0
    next
}

# 如果遇到下一个函数定义或文件结束
/^func / {
    if (in_test_func && !has_parallel) {
        print func_line
        print "	t.Parallel()"
        COUNT++
        in_test_func = 0
        has_parallel = 0
    } else if (in_test_func) {
        print func_line
        in_test_func = 0
        has_parallel = 0
    }
    print $0
    next
}

# 其他行
{
    if (in_test_func && func_line != "") {
        # 如果这是函数体的第一行，检查是否需要添加 t.Parallel()
        if (!has_parallel && $0 !~ /^\s*\/\// && $0 !~ /^\s*$/) {
            print func_line
            print "	t.Parallel()"
            COUNT++
            in_test_func = 0
            has_parallel = 0
            func_line = ""
        }
    }
    print $0
}

END {
    # 处理最后一个测试函数
    if (in_test_func && !has_parallel) {
        print func_line
        print "	t.Parallel()"
        COUNT++
    } else if (in_test_func) {
        print func_line
    }
}
' "$FILE" > "$TMPFILE"

# 替换原文件
mv "$TMPFILE" "$FILE"

echo "完成: 已检查并确保所有测试函数都有 t.Parallel()"

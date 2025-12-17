#!/bin/bash

# 简单脚本：给所有测试函数添加 t.Parallel()
# 用法：./add_parallel_simple.sh

FILE="alicloud/resource_alicloud_cs_kubernetes_node_pool_test.go"

if [ ! -f "$FILE" ]; then
    echo "错误: 文件 $FILE 不存在"
    exit 1
fi

echo "正在处理文件: $FILE"

# 使用 sed 和 awk 组合处理
# 1. 找到所有测试函数
# 2. 检查下一行是否有 t.Parallel()
# 3. 如果没有，则添加

python3 << 'PYTHON_SCRIPT'
import re
import sys

file_path = "alicloud/resource_alicloud_cs_kubernetes_node_pool_test.go"

with open(file_path, 'r') as f:
    lines = f.readlines()

output = []
i = 0
count = 0

while i < len(lines):
    line = lines[i]
    output.append(line)
    
    # 检查是否是测试函数定义
    if re.match(r'^func Test.*\(t \*testing\.T\) \{', line):
        # 检查下一行是否是 t.Parallel()
        if i + 1 < len(lines):
            next_line = lines[i + 1]
            if not re.match(r'^\s*t\.Parallel\(\)', next_line):
                # 添加 t.Parallel()
                output.append('\tt.Parallel()\n')
                count += 1
                i += 1
                continue
    
    i += 1

with open(file_path, 'w') as f:
    f.writelines(output)

print(f"完成: 已为 {count} 个测试函数添加 t.Parallel()")
PYTHON_SCRIPT



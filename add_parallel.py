#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
脚本：给所有测试函数添加 t.Parallel()
用法：python3 add_parallel.py
"""

import re
import sys

FILE = "alicloud/resource_alicloud_cs_kubernetes_node_pool_test.go"

def add_parallel_to_tests():
    try:
        with open(FILE, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except FileNotFoundError:
        print(f"错误: 文件 {FILE} 不存在")
        sys.exit(1)
    
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
    
    # 写回文件
    with open(FILE, 'w', encoding='utf-8') as f:
        f.writelines(output)
    
    print(f"完成: 已为 {count} 个测试函数添加 t.Parallel()")

if __name__ == '__main__':
    add_parallel_to_tests()



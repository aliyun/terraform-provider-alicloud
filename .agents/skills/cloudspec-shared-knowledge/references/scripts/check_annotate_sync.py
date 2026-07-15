#!/usr/bin/env python3
"""
验证 rawdocs/annotate 下文档中定义的 annotate 是否都在
docs/corpora/common/annotates 中有完整的文档定义。

核心匹配策略：
  只通过 markdown 标题（### AX.X name annotate）精确匹配，
  避免因正文中偶然提到关键词而误判为已同步。

用法：
    python scripts/check_annotate_sync.py

输出：
    - 列出每个 rawdocs 源文件中提取到的 annotate 名称
    - 标记哪些 annotate 在 corpora 中找不到
    - 最终汇总所有缺失的 annotate
"""

import os
import re
import sys

RAWDOCS_DIR = os.path.join(os.path.dirname(__file__), "..", "rawdocs", "annotate")
CORPORA_DIR = os.path.join(os.path.dirname(__file__), "..", "docs", "corpora", "common", "annotates")


def clean_markdown(content):
    """清理 font 标签和 ** 加粗标记，方便统一匹配。"""
    cleaned = re.sub(r'<font[^>]*>', '', content)
    cleaned = re.sub(r'</font>', '', cleaned)
    cleaned = cleaned.replace('**', '')
    return cleaned


def extract_annotate_names(filepath):
    """
    从 markdown 文件中提取所有通过标题定义的 annotate 名称。
    匹配格式示例：
        ### A1.1 length annotate
        ### A4.2 horizontalCodeMapping annotate
        ### A8.5 actionTrail annotate
        ## A6.1 http 注解
    提取出的名称为：length, horizontalCodeMapping, actionTrail, http 等。
    """
    with open(filepath, "r", encoding="utf-8") as file:
        content = file.read()

    cleaned = clean_markdown(content)

    # 匹配 ## 或 ### 开头，后跟 AX.X 编号、annotate 名称、annotate/注解 关键词
    pattern = r'#{2,}\s+A[\d.]+\s+(\S+)\s+(?:[Aa]nnotate|注解)'
    matches = re.findall(pattern, cleaned)

    annotate_names = []
    for match in matches:
        name = match.strip().strip('*').strip()
        if name:
            annotate_names.append(name)

    return annotate_names


def load_corpora_annotates(corpora_dir):
    """
    读取 corpora 目录下所有 .md 文件，通过标题精确提取每个文件中
    定义的 annotate 名称。

    返回：
      - annotate_to_file: dict，annotate名称(小写) -> 所在文件名
      - file_list: list，corpora 目录下的文件名列表
    """
    annotate_to_file = {}
    file_list = []

    if not os.path.isdir(corpora_dir):
        print(f"❌ corpora 目录不存在: {corpora_dir}")
        return annotate_to_file, file_list

    for filename in sorted(os.listdir(corpora_dir)):
        if filename.endswith(".md"):
            file_list.append(filename)
            filepath = os.path.join(corpora_dir, filename)
            names = extract_annotate_names(filepath)
            for name in names:
                annotate_to_file[name.lower()] = filename

    return annotate_to_file, file_list


def find_annotate_in_corpora(annotate_name, corpora_annotate_map):
    """
    检查某个 annotate 名称是否在 corpora 文件中有完整的文档定义。
    只通过标题匹配，避免因为正文中偶然提到关键词而误判为已同步。
    """
    name_lower = annotate_name.lower()
    if name_lower in corpora_annotate_map:
        return True, corpora_annotate_map[name_lower]
    return False, None


def main():
    rawdocs_dir = os.path.abspath(RAWDOCS_DIR)
    corpora_dir = os.path.abspath(CORPORA_DIR)

    print("=" * 70)
    print("Annotate 同步检查工具")
    print("=" * 70)
    print(f"📂 rawdocs 目录: {rawdocs_dir}")
    print(f"📂 corpora 目录: {corpora_dir}")
    print()

    if not os.path.isdir(rawdocs_dir):
        print(f"❌ rawdocs 目录不存在: {rawdocs_dir}")
        sys.exit(1)

    # 加载 corpora 中已定义的 annotate（通过标题精确匹配）
    corpora_annotate_map, corpora_file_list = load_corpora_annotates(corpora_dir)
    if not corpora_file_list:
        print("❌ corpora 目录下没有找到任何 .md 文件")
        sys.exit(1)

    print(f"📄 corpora 文件列表:")
    for filename in corpora_file_list:
        print(f"   - {filename}")
    print()

    print(f"📋 corpora 中已定义的 annotate 数量: {len(corpora_annotate_map)}")
    for name, filename in sorted(corpora_annotate_map.items(), key=lambda x: x[1]):
        print(f"   - {name} ({filename})")
    print()

    # 遍历 rawdocs 文件，提取 annotate 并检查
    all_missing = []
    total_annotates = 0

    rawdoc_files = sorted([
        f for f in os.listdir(rawdocs_dir) if f.endswith(".md")
    ])

    for rawdoc_file in rawdoc_files:
        filepath = os.path.join(rawdocs_dir, rawdoc_file)
        annotate_names = extract_annotate_names(filepath)

        if not annotate_names:
            continue

        print(f"📄 {rawdoc_file}")
        print(f"   提取到 {len(annotate_names)} 个 annotate:")

        for name in annotate_names:
            total_annotates += 1
            found, found_in = find_annotate_in_corpora(name, corpora_annotate_map)
            if found:
                print(f"   ✅ {name}  (found in {found_in})")
            else:
                print(f"   ❌ {name}  — 未在 corpora 中找到!")
                all_missing.append((rawdoc_file, name))

        print()

    # 汇总
    print("=" * 70)
    print("📊 汇总")
    print("=" * 70)
    print(f"总计检查 annotate 数量: {total_annotates}")
    print(f"已同步: {total_annotates - len(all_missing)}")
    print(f"缺失:   {len(all_missing)}")
    print()

    if all_missing:
        print("🚨 以下 annotate 在 docs/corpora/common/annotates 中缺失:")
        print("-" * 70)
        for source_file, annotate_name in all_missing:
            print(f"  [{source_file}]  →  {annotate_name}")
        print("-" * 70)
        print()
        print("请检查是否需要将这些 annotate 同步到 corpora 目录中。")
        sys.exit(1)
    else:
        print("✅ 所有 annotate 均已同步，没有遗漏！")
        sys.exit(0)

if __name__ == "__main__":
    main()

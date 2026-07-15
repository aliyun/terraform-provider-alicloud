# CloudSpec CodeFix 使用说明

源自 Aliyun CLI cspec plugin 的 idl_codefix 模块，用于自动修复 CloudSpec 规范违规问题。

## 概述

`aliyun cspec codefix` 可自动修复多种规范问题，支持按规则 ID、组件类型、组件名称进行精确或批量修复。

## 基本命令

```bash
# 预览修复效果（不实际修改）
aliyun cspec codefix -r <规则ID> -t <类型> -c <组件名> --dry-run

# 实际修复
aliyun cspec codefix -r <规则ID> -t <类型> -c <组件名>

# 批量修复（不指定 -c 时）
aliyun cspec codefix -r <规则ID> -t <类型> --dry-run
```

## 参数说明

| 参数 | 说明 |
|------|------|
| `-r, --rule` | 修复规则 ID（必需），如 E-AT-0006 |
| `-t, --type` | 组件类型（必需）：`operation`、`resource`、`service` |
| `-c, --component` | 指定要修复的组件名称（可选） |
| `-n, --namespace` | 指定命名空间（可选，会从 CloudSpec 推断） |
| `--dry-run` | 预览模式，只显示修复效果，不实际修改 |
| `--verbose` | 详细输出模式 |

## 使用建议

- **语法类错误**：优先尝试 `aliyun cspec codefix`，可精确修复已知规则
- **其他规范问题**：若 codefix 无对应规则，需手动分析修复
- **修复后验证**：执行 `aliyun cspec build` 和 `aliyun cspec check --name <组件名>` 确认

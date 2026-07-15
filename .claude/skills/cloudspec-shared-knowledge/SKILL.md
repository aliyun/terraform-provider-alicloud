---
name: cloudspec-shared-knowledge
description: |
  CloudSpec IDL 完整知识库与共享参考文档库。不直接调用，作为其他 cloudspec-* skill 的共享文档库。包含注解规范、基础语法、设计指南、测试指南和 CLI 参考。
  Triggers: "cspec语法", "注解规范", "基础类型", "保留字", "命名空间", "CloudSpec知识库", "复合类型", "设计指南", "CloudSpec文档".
---

# CloudSpec 共享知识库

本 skill 是 CloudSpec 的完整参考文档库，供其他 cloudspec-* skill 按需引用。

> **本 skill 不需要直接调用**，它会被以下 skill 在需要时自动引用：
> - cloudspec-operation-edit
> - cloudspec-resource-edit
> - cloudspec-test-fix
> - cloudspec-resource-test
> - cloudspec-norm-check-fix
> - cloudspec-build-fix

## 文档索引

### 注解文档（references/docs/corpora/common/annotates/）

| 文件 | 内容 |
|------|------|
| operation-annotate.md | Operation 相关注解完整规范（@http、@backendConfigurationHttp、@operationInfo、@paginated 等） |
| resource-annotate.md | Resource 相关注解完整规范（@arn、@resourceBaseInfo、@operationMapping、@rac 等） |
| enterprise-annotate.md | 企业级能力注解（@terraform、@rmc、@tagService、@ros、@config、@sdk 等） |
| iam-annotate.md | IAM/RAM 权限注解（@ram、@defineAction、@defineCondition、@requiredPermission、@conditions 等） |
| assignment-constraint-annotate.md | 传值约束注解（@required、@readonly、@clientOptional、@clientProhibited 等） |
| value-constraint-annotate.md | 值约束注解（@default、@length、@range、@regexPattern、@enums 等） |
| document-annotate.md | 文档注解（@document、@deprecated、@sensitive 等） |
| test-annotate.md | 测试注解（@testConfig、@before、@after、@step 等） |
| event-annotate.md | 事件类注解（@eventConfig、@gatewayTransportType、@events） |
| auto-generate-annotate.md | 自动化生成注解（@autoGenerateOperations、@autoGenerateResource 等） |
| constraint-annotate.md | 约束器注解（@skipConstraint） |
| error-annotate.md | 错误注解 |
| service-types.md | Service 注解详解 |

### 基础语法（references/docs/corpora/common/basic-grammer/）

| 文件 | 内容 |
|------|------|
| quick-start.md | CloudSpec IDL 快速入门 |
| basic-types.md | 基础类型详解 |
| composite-types.md | 复合类型说明 |
| model-composition.md | 模型组合 |
| namespace.md | 命名空间管理 |
| reserved-words.md | 保留字列表 |
| cloudspec-idl-grammer.md | CloudSpec IDL 语法指南 |

### 设计指南

| 文件 | 内容 |
|------|------|
| references/docs/corpora/common/how-to-design-operation.md | Operation 设计指南（设计原则、配置规范、最佳实践、检查清单） |
| assets/roa-vs-rpc.md | RPC/ROA 风格对照（skills 仓库自有文件） |
| ../cloudspec-idl-guide/references/local-development-workflow.md | 本地环境、build/check/test、本地转 yaml 和资源测试生成前检查 |
| ../cloudspec-idl-guide/references/common-examples.md | ROA body/backend、operationMapping body 路径、资源测试和 operation 测试高频示例 |

### 测试指南（references/docs/corpora/resource-test/）

| 文件 | 内容 |
|------|------|
| test-guide/resource-test-quick-start.md | 资源测试快速开始 |
| resource-test-faq/ | 资源测试常见问题（网络错误、重试策略、资源不存在条件、异步、差异、枚举、依赖等） |

### 插件内资源测试指南

| 文件 | 内容 |
|------|------|
| ../cloudspec-resource-test/references/cloudspec-agent.md | cloudspec-agent 安装、MCP 配置提醒和 cover 命令 |
| ../cloudspec-resource-test/references/resource-test-syntax.md | `$test`、init/modifies/destroy、operation input、依赖和函数语法 |
| ../cloudspec-resource-test/references/verification-commands.md | build、test run、check、本地转 yaml 的验证命令 |
| ../cloudspec-resource-test/references/diagnostic-report.md | 资源测试失败时面向产品/后端的诊断报告模板 |

### CLI 指南（references/docs/corpora/cli-guide/）

| 文件 | 内容 |
|------|------|
| cloudspec-cli-quick-start.md | Aliyun CLI cspec plugin 完整使用指南 |

## 实测注解与类型陷阱速查

> 以下是经过 4 轮实测验证的高频易错点，`aliyun cspec build` 会直接报错。详细案例见 `cloudspec-build-fix/references/build-errors.md`。

### 不存在的类型名

| 错误 | 正确 | 说明 |
|------|------|------|
| `integer` | `int32` 或 `int64` | CloudSpec 没有 `integer` |
| `long` | `int64` | 没有 `long` |
| `number` | `float` 或 `double` | 没有 `number` |
| `list` | `array` | 没有 `list` 关键字 |

### 不存在的注解（报 `自定义 annotate xxx 未定义`）

| 错误注解 | 正确替代 |
|----------|----------|
| `@backend({...})` | `@backendConfigurationHttp({...})` |
| `@idempotent` | 移除 |
| `@pattern("regex")` | `@regexPattern("regex")` |
| `@maxLength(n)` | `@length({ min: 0, max: n })` |
| `@min(n)` / `@max(n)` | `@range({ min: n, max: m })` |

### 语法易错

| 错误 | 正确 | 说明 |
|------|------|------|
| `@enums([{...}{...}])` | `@enums([{...}, {...}])` | 数组项之间必须有逗号 |

## 文档来源

本 skill 的 `references/` 目录是 [cloudspec-docs-for-ai](git@gitlab.alibaba-inc.com:cloudspec/cloudspec-docs-for-ai.git) 仓库的 **git submodule**，无需手动同步。

更新文档：

```bash
cd cloudspec-shared-knowledge/references
git pull origin master
cd ../..
git add cloudspec-shared-knowledge/references
git commit -m "chore: update shared-knowledge submodule"
```

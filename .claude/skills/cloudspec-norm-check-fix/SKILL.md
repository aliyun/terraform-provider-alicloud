---
name: cloudspec-norm-check-fix
description: |
  CloudSpec 规范检查与自动修复 skill。对组件运行 aliyun cspec check，发现规范问题后尽力自动修复；若存在冲突、语义不确定或需要用户确认的问题，必须报告结论和原因。支持 aliyun cspec codefix 按规则编号修复（如 M-RL-0044、E-AT-0006）。
  Triggers: "规范检查", "check", "aliyun cspec check", "norm check", "修复规范", "codefix", "检查规范性并修复", "规范修复", "fix norm", "规则编号", "M-RL-", "E-AT-", "注解缺失", "命名规范", "name 是空的", "有规范问题的话顺手修掉".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec 规范检查与自修复

给定 CloudSpec 项目及待检查的组件（资源或操作），**必须运行 `aliyun cspec check` 进行规范检查，并尽力自动修复发现的问题**。如果规范问题与现有语义、兼容性或用户意图冲突，不要为了通过 check 擅自改动契约；无法安全修完时必须给出结论、原因、剩余问题和需用户确认项。

> **前置**：若尚未了解整体编写流程，可先调用 **cloudspec-idl-guide**。

## Step1 识别 API 风格（前置必做）

CloudSpec 项目分为 **RPC** 和 **ROA（RESTful）** 两种风格，修复规则因风格不同而有差异。**在检查前必须先识别风格**。

**优先使用 CLI 命令获取项目信息**：

```bash
aliyun cspec baseinfo
```

该命令输出 JSON 格式的项目基本信息，包含 `apiStyle`（rpc/roa）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`（资源列表）、`operations`（操作列表）等字段。从 `apiStyle` 字段即可判断 API 风格。

仅当 CLI 不可用时，再手动检查 `main.cspec` 或操作文件，详见 [roa-vs-rpc.md](references/roa-vs-rpc.md)。

> **需要查阅注解规范？** 修复规范问题时如需了解某个注解的完整属性，可按需阅读共享知识库 `cloudspec-shared-knowledge` 中的对应文档：
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/operation-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/resource-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/enterprise-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/iam-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/constraint-annotate.md`（@skipConstraint）

## Step2 运行 aliyun cspec check

> **⚠️ Inner API 豁免**：运行 `aliyun cspec baseinfo`，若 `isInnerApi` 为 `true`，则该项目为 inner API，**不需要做规范检查**，可跳过本 skill 的所有步骤。

在项目根目录执行：

```bash
aliyun cspec check --name <组件名称>
```

该命令会输出规范校验结果。若包含 `ERROR`，则表示存在需修复的规范问题。

## Step3 解析错误并修复

根据 `aliyun cspec check` 的输出逐项修复。check 的报错信息通常已经明确指出了问题所在和期望值，应优先信任报错内容，按需查阅文档。

### 修复边界

- 能确定不改变 API 语义、兼容性和用户意图的问题，直接修复。
- 涉及字段删除、重命名、类型变化、必填性变化、枚举收窄、路径/方法变化、错误码语义变化、资源 identifyDefinition 或 operationMapping 语义变化时，视为可能的非兼容变更。
- 非兼容变更以原定义为准；不要为了规范检查通过擅自改原有契约，必须先向用户说明影响并请求确认。
- 当 check 规则之间冲突、同项目参考写法冲突、缺少业务语义、或修复会破坏兼容时，停止该项自动修复并记录为待确认问题。

### 优先尝试 codefix

若错误信息中包含**规范 ID/规则编号**（如 `M-RL-0044`、`E-AT-0006`），可优先尝试：

```bash
aliyun cspec codefix -r <规则编号> -t operation|resource|service -c <组件名> --dry-run   # 预览
aliyun cspec codefix -r <规则编号> -t operation|resource|service -c <组件名>           # 执行修复
```

codefix 不适用或无法覆盖的问题，再手动修复。常见问题与修复方式参考 [review-rules.md](references/review-rules.md)，基础语法参考 [quick-start.md](references/quick-start.md)：

| 问题类型 | 示例 | 修复方式 |
|----------|------|----------|
| @document name 为空 | 操作/资源的 @document 中 `name: ""` | **先看 check 报错是否直接给出期望值**；若无，再参考：1) 同项目同类型操作的 name 命名规律（如"创建/查询/修改/删除 + 资源中文名"）2) 操作名称本身的语义（如 `CreateSession` → `"创建会话资源"`）。**禁止随意填写无意义的占位描述**（如 "Schema of Response"、"操作描述"）。 |
| 命名不规范 | 未使用 PascalCase | 统一为 PascalCase |
| 缺少必要注解 | 缺 @visibility、@ram 等 | 参考同项目其他文件补充 |
| 错误定义不完整 | 缺 httpCode、errorCode | 补充完整错误定义 |

## Step4 迭代直至通过

修复后**必须**重新运行完整的验证流程：

```bash
# 1. 先确保编译通过
aliyun cspec build

# 2. 再运行规范检查
aliyun cspec check --name <组件名称>
```

若仍有 ERROR 则继续修复。**最多迭代 10 次**，每次修复后都必须重新运行上述两个命令验证。

> **重要**：`aliyun cspec check` 可能一次报告多个错误，必须逐一处理所有 ERROR。能安全修复的必须修复；无法安全修复的，必须明确留下结论、原因和需用户确认项，不能假装已全部通过。特别是涉及多个组件的场景（如同时修复 CreateConfigRule 和 ActiveAggregateConfigRules），需要对每个组件分别运行 `aliyun cspec check --name <组件名>` 确认状态。

## Step5 未全部修复时的报告格式

如果最终仍有 ERROR，或部分修复需要用户决策，输出必须包含：

- **结论**：已修复 / 部分修复 / 未修复。
- **已修复项**：规则编号、组件名、改动摘要。
- **剩余问题**：规则编号、组件名、check 原因。
- **不能自动修复的原因**：冲突、缺少业务语义、可能非兼容、需要用户确认等。
- **建议下一步**：保留原定义、请求用户确认具体语义，或由业务 owner 决策。

## 注意事项

- **本 skill 会主动修改代码**，以修复规范问题为目标。
- 修复时参考同项目其他文件的写法，保持风格一致。
- 区分 RPC/ROA 风格，不同风格的规范要求不同。
- 修复过程中不可引入编译错误，每次修改后需运行 `aliyun cspec build` 验证。
- 规范修复不得覆盖用户明确需求；非兼容变更以原定义为准，除非用户明确确认。

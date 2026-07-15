---
name: cloudspec-build-fix
description: |
  CloudSpec 编译错误分析与自动修复 skill。运行 aliyun cspec build 分析错误类型（语法、引用、拼写、缺少import、保留字、类型不匹配）并自动修复。支持 aliyun cspec codefix 基于规则的自动修复。
  Triggers: "编译报错", "build 失败", "修复编译错误", "aliyun cspec build", "compilation error", "build error", "语法错误", "syntax error", "fix build", "codefix", "import缺失", "引用错误", "保留字", "类型不匹配".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec 编译修复

给定一个CloudSpec项目，**必须严格按照以下步骤运行编译并修复错误**。

> **前置**：若尚未了解整体编写流程，可先调用 **cloudspec-idl-guide**。

## Step1 了解项目结构

CloudSpec IDL语法和项目结构查看[quick-start.md](references/quick-start.md)，完整的IDL语法规则（所有合法的类型定义、注解语法、值类型、保留字等）详见[idl-grammar-reference.md](references/idl-grammar-reference.md)。

### 识别API风格（前置必做）

CloudSpec项目分为**RPC**和**ROA（RESTful）**两种风格，编译错误的修复方式因风格不同而有差异。**在修复前必须先识别风格**。

**优先使用 CLI 命令获取项目信息**：

```bash
aliyun cspec baseinfo
```

该命令输出 JSON 格式的项目基本信息，包含 `apiStyle`（rpc/roa）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`（资源列表）、`operations`（操作列表）等字段。从 `apiStyle` 字段即可判断 API 风格。

仅当 CLI 不可用时，再手动检查 `main.cspec` 或操作文件，详见[roa-vs-rpc.md](references/roa-vs-rpc.md)。

> **需要查阅注解规范？** 修复编译错误时如需了解某个注解的完整属性，可按需阅读共享知识库 `cloudspec-shared-knowledge` 中的对应文档：
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/operation-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/resource-annotate.md`
> - `../cloudspec-shared-knowledge/references/docs/corpora/common/annotates/enterprise-annotate.md`

标准项目结构：

```
project/
├── main.cspec          # 服务定义 + import
├── resources/          # 资源定义
├── operations/         # 操作定义
└── tests/              # 测试用例
```

> **最小完整项目示例**：如需了解一个可编译通过的最小项目骨架（包含 main.cspec、resource、operation 及所有必填注解的完整引用关系），参见[minimal-project-example.md](references/minimal-project-example.md)。注意：示例仅供参考骨架结构，**实际修复时必须以用户项目中已有文件的风格为准**。

确认当前目录为CloudSpec项目根目录（包含`main.cspec`文件），如果不是，需要先定位到正确目录。

## Step2 运行编译

在CloudSpec项目根目录下执行：

```bash
aliyun cspec build
```

- 如果编译成功（无错误输出），告知用户编译通过，直接结束。
- 如果编译失败，进入Step3分析并修复错误。

## Step3 分析错误并修复

### 3.1 错误分类

编译错误按以下类别分类处理，常见错误及修复方式详见[build-errors.md](references/build-errors.md)。

| 类别 | 说明 | 示例 |
|------|------|------|
| **语法错误** | IDL语法不符合规范，完整语法规则见[idl-grammar-reference.md](references/idl-grammar-reference.md) | 缩进错误、括号不匹配、注解格式错误、保留字未转义、类型关键字误用 |
| **引用错误** | 引用的类型/操作/错误码未定义 | operation引用了不存在的struct或error |
| **命名空间错误** | namespace相关问题 | import的namespace未定义、namespace不一致 |
| **资源生命周期错误** | 资源与操作的关联问题 | 资源引用了不存在的操作、属性引用路径错误 |
| **类型错误** | 类型不匹配或不合法 | 测试用例中属性类型与资源定义不一致 |
| **约束错误** | constraint校验不通过 | 资源关系、映射规则不满足约束条件 |

### 3.2 修复策略

按以下优先级修复：

1. **语法错误优先**：语法错误会阻塞后续检查，必须首先修复。
2. **引用错误次之**：找到未定义的引用，确认是拼写错误还是缺少定义，对应修复。
3. **逻辑错误最后**：类型不匹配、映射关系错误等需要理解业务语义后修复。

### 3.3 修复原则

- **最小改动原则**：每次只修复一类错误，避免引入新问题。
- **保留原始意图**：修复时保持原有代码的设计意图，不随意改变业务语义。
- **不修改operations目录**：除非用户明确要求，否则不修改`./operations`目录下已有的API元数据文件。
- **不添加import**：除非确认缺少必要的import语句，否则不添加新的import。

### 3.4 尝试CLI自动修复

如果错误信息中包含规则编号（如`M-RL-0044`、`E-AT-0006`），优先尝试使用CLI自动修复：

```bash
# 预览修复效果（不实际修改）
aliyun cspec codefix -r <规则编号> -t <类型> --dry-run

# 确认无误后执行修复
aliyun cspec codefix -r <规则编号> -t <类型>
```

参数说明：
- `-r`：规则编号，如`M-RL-0044`、`E-AT-0006`
- `-t`：组件类型，`operation`、`resource`或`service`
- `-c`：可选，指定组件名称（不指定则修复所有该类型组件）
- `--dry-run`：预览模式，只显示不修改

如果是资源级别的问题，也可以尝试：

```bash
aliyun cspec fix resource -n <资源名称>
```

如果`main.cspec`的import索引与实际文件不一致：

```bash
aliyun cspec fix index
```

### 3.5 手动逐项修复

CLI自动修复不适用或未覆盖的错误，手动修复：

1. **定位**：根据错误信息中的文件路径和行号，读取对应文件。
2. **分析**：确认错误根因，参考[build-errors.md](references/build-errors.md)中的解决方案。语法类错误可对照[idl-grammar-reference.md](references/idl-grammar-reference.md)确认合法写法。
3. **修复**：执行最小改动修复。
4. **记录**：记录修复内容（文件、行号、修改前后的代码）。

## Step4 重新编译验证

修复完成后，重新运行编译：

```bash
aliyun cspec build
```

- 如果编译成功，进入Step5。
- 如果仍有错误，回到Step3继续修复。
- 如果同一个错误反复出现无法解决，记录问题并告知用户需要人工介入。

最多循环**10次**。如果10次后仍有编译错误，停止修复，在Step5中报告未解决的问题。

## Step5 输出修复报告

```
## 编译修复报告

### 编译结果
- 最终状态：成功 / 失败
- 编译轮次：N
- 修复问题数：N

### 修复记录
1. [文件路径:行号] 错误描述
   - 根因：...
   - 修复：...

### 未解决问题（如有）
1. [文件路径:行号] 错误描述
   - 原因：...
   - 建议：...
```

## 注意事项

- 绝不放弃。以坚持不懈为荣，以轻易放弃为耻。只要编译不通过，就要持续分析和修复，直到编译成功或达到最大尝试次数。
- 禁止修改`./operations`目录下已有的API元数据文件，除非用户明确要求。
- 如果编译错误明确指出缺少 import，**必须恢复或添加对应的 import**。优先使用 `aliyun cspec fix index` 自动修复，或手动添加通配符导入（`import "./operations/*"`）。
- 修复时充分参考同项目中其他文件的写法，保持风格一致。
- 如果错误涉及资源属性与API映射，需要同时检查资源定义和对应的操作定义来判断根因。
- 每轮修复后必须重新编译验证，不能假设修复一定正确。

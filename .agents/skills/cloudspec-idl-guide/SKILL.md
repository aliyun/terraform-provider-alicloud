---
name: cloudspec-idl-guide
description: |
  ⚠️ 仅限当前目录已有 main.cspec 文件时使用。用于在已 clone 的 cspec 项目目录内编辑 .cspec 文件。
  🚫 禁止在以下场景使用本 skill（必须改用 cloudspec-amp-workflow）：用户消息包含"初始化"、"镇元"、"namespace:"、"namespace"、"发布"、"publish"、"clone"、"拉代码"、"拉到本地"、"推送远程"、"新建API"、"新增API"、"API接口定义"、"接口定义"、"创建操作"中的任何一个。
  Triggers: "修改cspec文件", "编辑cspec", "编辑资源", "修改cspec", "curl转IDL", "edit resource", "cspec文件编辑", "修改cspec属性", "cspec check", "cspec build", "修改资源", "编辑操作", "修改操作", ".cspec", "cspec文件".
allowed-tools: Bash, Read, WebSearch, Write, Grep, MultiEdit, Edit, WebFetch, Glob
---

# CloudSpec IDL 编写指南

> **🚨 强制前置拦截（MUST DO FIRST）**：加载本 skill 后，在做任何其他事情之前，必须先执行 `ls main.cspec 2>/dev/null`。
> - 如果 `main.cspec` **不存在**，说明当前目录不是 cspec 项目。**立即停止本 skill 的所有后续步骤**，转交 `cloudspec-amp-workflow` skill 完成项目初始化（amp bootstrap → branch create → git clone）。**严禁在空目录手动创建 main.cspec 或任何 .cspec 文件**。
> - 如果用户消息包含"初始化"、"镇元"、"namespace:"、"发布"、"clone"、"拉代码"等 amp 生命周期关键词，更应转交 `cloudspec-amp-workflow`。
> - 只有当 `main.cspec` 存在时，才继续下面的流程。

**重要**：只要涉及**阅读、编辑、修改或删除**任意 .cspec 文件，都必须先查阅本 skill，了解流程与调用顺序，再根据任务类型调用相应专项 skill。与 Aliyun CLI cspec plugin 的 AI 工作流一致。

本文档是 CloudSpec IDL（.cspec 文件）编写的**总览与入口**，指导何时调用哪些专项 skill 以及执行顺序。

## 一、项目结构与基本概念

CloudSpec 模型简介、基本概念、类型、Resource/Operation 定义等详见 [quick-start.md](references/quick-start.md)。本地环境、build/check/test/yaml 统一流程见 [local-development-workflow.md](references/local-development-workflow.md)，高频示例见 [common-examples.md](references/common-examples.md)。

标准 CloudSpec 项目结构：

```
project/
├── main.cspec          # 服务定义、namespace、import
├── resources/          # 资源定义（.cspec）
├── operations/         # 操作定义（.cspec）
└── tests/              # 资源测试用例
```

- **Namespace**：`alicloud.{Product}.{PopCode}.v{PopVersion}`，所有组件在其下
- **Service**：main.cspec 中的 service 定义，包含 resources、operations
- **Resource**：资源元数据，含 identifyDefinition、properties、CRUD 操作映射
- **Operation**：OpenAPI 定义，含 input、output、errors
- **Annotate**：注解（@http、@document、@backendConfigurationHttp 等）用于补充元信息

## 二、编写流程（何时调用哪个 Skill）

### 流程总览

```
用户任务
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│ 0. 分支确认（必做）                                              │
│    检查当前 git 分支 + amp 分支，确保不在 master/main 上直接编辑  │
└─────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. aliyun cspec baseinfo（必做）                                   │
│    识别 apiStyle（RPC/ROA）、isInnerApi、namespace、资源/操作列表 │
└─────────────────────────────────────────────────────────────────┘
    │
    ├── 创建/编辑 Operation  ──►  cloudspec-operation-edit
    │
    ├── 编辑 Resource        ──►  cloudspec-resource-edit
    │
    ├── 编译报错             ──►  cloudspec-build-fix
    │
    ├── 规范检查并修复（只对本次操作的增量内容）       ──►  cloudspec-norm-check-fix
    │
    ├── 资源测试语法/生成/运行 ──►  cloudspec-resource-test
    │
    ├── 完成后验证           ──►  aliyun cspec build && 资源变更必跑 aliyun cspec check --name <ResourceName>
    │
    └── 变更持久化（验证通过后）──►  提示用户 commit + push 或转 cloudspec-amp-workflow 发布
```

### Step -1：项目存在性检查（最先执行）

> **为什么**：如果当前目录不是 cspec 项目（没有 `main.cspec`），任何编辑操作都无法进行。此时应转交 `cloudspec-amp-workflow` 先完成初始化（amp init → branch create → git clone），而不是尝试手动创建项目结构。

执行以下检查：

```bash
ls main.cspec 2>/dev/null
```

**根据结果决定下一步**：

| main.cspec 存在？ | 用户意图涉及初始化/namespace/发布？ | 处理 |
|---|---|---|
| ✅ 存在 | — | 继续 Step 0 |
| ❌ 不存在 | 是（提到"初始化"、"镇元项目"、"namespace:"、"发布"等） | **转交 `cloudspec-amp-workflow` skill**，由其完成 bootstrap → branch → clone → 回到本 skill |
| ❌ 不存在 | 否（用户只是问 cspec 语法问题等） | 提示用户：当前目录不是 cspec 项目，需要先通过 `cloudspec-amp-workflow` 初始化或 `cd` 到已有的 cspec 项目目录 |

> **⚠️ 严禁在空目录手动创建 main.cspec 骨架**。cspec 项目结构由 amp 后端维护，必须通过 `amp init` + `git clone` 获取。

### Step -0.5：本地环境检查

编辑、生成测试或运行验证前，先按 [local-development-workflow.md](references/local-development-workflow.md) 检查 Aliyun CLI、cspec 插件、用户 profile、项目 baseinfo 和必要的测试工具。环境不完整时先修环境，不进入编辑或测试生成。

### Step 0：分支确认（编辑前必做）

> **为什么**：所有 cspec 编辑必须在 feature 分支进行，直接改 master/main 无法通过 `amp publish` 发布。

执行以下检查：

```bash
git branch --show-current 2>/dev/null
cat .amp/context.yaml 2>/dev/null
```

**根据结果决定下一步**：

| 当前分支 | .amp/context.yaml 有 branch | 处理 |
|---------|---------------------------|------|
| `feature/*` 或非 master/main | 有 | ✅ 确认后继续 Step 1 |
| `feature/*` 或非 master/main | 无 | 提醒用户运行 `amp context set branch <当前分支名>` |
| `master` / `main` | — | ⚠️ **必须先建分支**，见下方 |

**在 master/main 上时**，用 `AskUserQuestion` 提示用户：

```
检测到当前在 master 分支，直接编辑无法通过 amp publish 发布。请选择：
- 选项 A: 新建 feature 分支（推荐）—— 请提供分支名，如 feature/add-xxx
- 选项 B: 我确认只做本地实验，不需要发布
- 选项 C: 我已有分支，告诉你分支名切过去
```

- 选 A → 转交 `cloudspec-amp-workflow` step 3.3 创建分支 + 切换，完成后回来继续 Step 1。
- 选 B → 跳过分支检查，继续 Step 1（记录用户已确认）。
- 选 C → `git checkout <branch>` + `amp context set branch <branch>`，继续 Step 1。

### 各 Skill 调用时机

| 用户意图 | 调用 Skill | 说明 |
|----------|------------|------|
| 新建 API、添加参数、改操作注解 | cloudspec-operation-edit | 先运行 `aliyun cspec baseinfo` 获取风格 |
| 添加/修改/删除资源属性 | cloudspec-resource-edit | 先运行 `aliyun cspec baseinfo` 获取风格 |
| 编译报错、build 失败 | cloudspec-build-fix | 定位并修复语法、引用、拼写等 |
| 规范检查、check 有 ERROR | cloudspec-norm-check-fix | 运行 aliyun cspec check 并自动修复 |
| 资源测试语法、生成资源测试、提升覆盖率、运行资源测试 | cloudspec-resource-test | 安装/使用 cloudspec-agent，生成覆盖率用例，运行资源测试并跑 build/test/check |
| 迁移/订正资源测试 | cloudspec-test-migrate / test-fix | 迁移旧用例或修复已有失败用例 |

### 核心原则

1. **先获取项目信息**：任何编辑前，优先运行 `aliyun cspec baseinfo` 获取项目的 API 风格、namespace、资源列表、操作列表等基本信息。该命令输出 JSON 格式，包含 `apiStyle`（RPC/ROA）、`namespace`、`popCode`、`apiVersion`、`isInnerApi`、`resources`、`operations` 等字段。
2. **先风格，后编辑**：从 `aliyun cspec baseinfo` 的 `apiStyle` 字段判断风格（`RPC` 或 `ROA`），再按对应规范编辑。
3. **编辑后验证**：修改 .cspec 后必须运行 `aliyun cspec build`。资源变更必须运行 aliyun cspec check --name <ResourceName>（实际命令：`aliyun cspec check --name <ResourceName>`），除非 `aliyun cspec baseinfo` 显示 `isInnerApi: true`。非资源组件按任务要求运行 `aliyun cspec check --name <组件名>`。
4. **参考同项目**：新建或修改时，参照同项目已有操作的注解风格、命名、配置模式。
5. **典型例子优先复用**：ROA body 参数、ROA 后端配置、operationMapping body 路径、资源测试和 operation 测试示例见 [common-examples.md](references/common-examples.md)。

## 三、基本语法速查

### 文件头

```cspec
$version: 1
namespace: alicloud.{Product}.{PopCode}.v{PopVersion}
```

### 注解语法

- 使用冒号：`key: value`，不用等号
- 注解与目标之间无空行
- 2 空格缩进

### 命名规范

- 资源名、操作名、结构体、字段：PascalCase
- List 操作用复数：`ListInstances`

## 四、验证命令

| 命令 | 用途 |
|------|------|
| `aliyun cspec build` | 编译项目，检查语法、引用、类型 |
| `aliyun cspec check --name <组件名>` | 单组件规范检查（注解、命名、映射等） |
| `aliyun cspec test run -n <MainTestName>` | 运行指定入口资源测试 |
| `aliyun cspec yaml --help` | 先确认本机 CLI 的 yaml 转换参数，再执行本地转 yaml |
| `aliyun cspec codefix` | 自动修复部分规范问题（配合 -r/-t/-c 参数） |

编辑完成后，先 `aliyun cspec build`。如果改动了 `resources/*.cspec`，必须对每个涉及资源运行 `aliyun cspec check --name <ResourceName>`；非资源组件按任务要求运行 `aliyun cspec check`。更多 CLI 说明见 [cli-commands.md](references/cli-commands.md)、[codefix-usage.md](references/codefix-usage.md)。

> **⚠️ Inner API 豁免**：运行 `aliyun cspec baseinfo`，若 `isInnerApi` 为 `true`，则该项目为 inner API，**不需要运行 `aliyun cspec check`**，只需确保 `aliyun cspec build` 通过即可。

> **共享知识库**：所有注解的完整规范、基础语法、设计指南等详细文档位于 `cloudspec-shared-knowledge` skill 中，各专项 skill 会按需引用。完整文档索引见 `../cloudspec-shared-knowledge/SKILL.md`。

## 五、变更持久化（验证通过后必做）

当 `aliyun cspec build` 和 `aliyun cspec check` 全部通过后，**必须提示用户持久化变更**，避免本地修改丢失。

使用 `AskUserQuestion` 向用户确认下一步：

```
✅ build + check 已全部通过。请选择下一步操作：
- 选项 A: 提交并推送到远端（推荐）—— 将变更 commit + push 到当前 feature 分支，确保工作落盘保存
- 选项 B: 继续编辑 —— 暂不提交，稍后手动处理
- 选项 C: 直接发布 —— 跳转 cloudspec-amp-workflow 执行 amp publish（publish 会自动同步源码）
```

**各选项处理**：

- **选 A**（推荐）→ 执行以下步骤：
  1. `git add` 本次新增/修改的 .cspec 文件（**只添加具体文件，不用 `git add .`**）
  2. `git commit -m "feat: <简要描述本次变更>"`（提交信息由 agent 根据实际变更自动生成）
  3. `git push origin <当前分支名>`
  4. 告知用户：「变更已推送到远端分支 `<branch>`，可随时通过 cloudspec-amp-workflow 执行发布。」

- **选 B** → 提醒用户：「请记得在结束工作前手动 commit + push，避免本地修改丢失。」

- **选 C** → 转交 `cloudspec-amp-workflow` skill step 5（发布前验证清单）。

> **注意**：此步骤是**提示用户确认后再执行**，不是自动替用户 push。agent 在用户明确选择"提交并推送"后才执行 git 操作。

## 六、专项 Skill 说明

各 skill 的详细步骤与参考文档见各自 SKILL.md：

- **cloudspec-operation-edit**：创建/编辑 Operation，含注解与参数；ROA body/backend 高频示例也见 [common-examples.md](references/common-examples.md)
- **cloudspec-resource-edit**：编辑 Resource 属性与注解
- **cloudspec-build-fix**：分析编译错误并修复
- **cloudspec-norm-check-fix**：运行 check 并自动修复规范问题
- **cloudspec-resource-test**：资源测试语法、生成、运行、失败分析和验证入口
- **cloudspec-test-migrate / test-fix**：迁移旧测试或修复已有失败测试

## 七、跨域衔接：cloudspec-amp-workflow

如果用户的任务**不是**纯 .cspec 编辑，而是涉及以下场景，本 skill 不处理，**请转交 `cloudspec-amp-workflow` skill**：

- 创建 / 切换 / 删除 amp 分支（`amp branch *`）
- 拉 cspec 源码仓库到本地（`git clone -b <branch> git@gitlab.alibaba-inc.com:cloudspec-model/<repo>.git`）
- amp 工作区初始化 / 体检 / BUC 登录配置（`amp init` / `amp doctor` / `amp login`；`amp config set-ak` 仅低频可选）
- 发布到 daily / pre（`amp publish *`）

典型完整流程：用户先在 `cloudspec-amp-workflow` 里建分支 + clone 仓库，进入 cspec 工作区后回到本 skill 做 .cspec 编辑，编辑完成 + `aliyun cspec build` & `check` 全绿，再回到 `cloudspec-amp-workflow` 做 publish。

---

**总结**：本 skill 作为 CloudSpec IDL 编写的入口，负责串联 cspec 域内流程与调用顺序。具体编写步骤由上述专项 skill 完成。涉及 amp 分支 / 仓库 clone / 发布请走 `cloudspec-amp-workflow`。

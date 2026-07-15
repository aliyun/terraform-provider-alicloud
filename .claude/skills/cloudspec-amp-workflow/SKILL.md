---
name: cloudspec-amp-workflow
description: |
  镇元/amp 项目的唯一入口技能。无论用户是要初始化项目、新建API、创建操作、拉代码、还是发布，都必须先经过本 skill。本 skill 负责完整生命周期：amp bootstrap → 分支管理 → git clone cspec 仓库 → API/接口定义编辑 → 发布 daily/pre。同时支持 policy/domain/error-code/gateway 管理。
  🔴 路由规则：只要用户消息包含以下任意关键词，必须使用本 skill：初始化、镇元、namespace、发布、publish、clone、拉代码、拉到本地、推送远程、新建API、新增API、API接口定义、接口定义、创建操作、IDL、IDL定义。
  Triggers: "amp", "amp branch", "amp publish", "建分支", "创建分支", "切分支", "分支管理", "发布到 daily", "发布到 pre", "publish daily", "publish pre", "集成发布", "amp 流程", "amp 初始化", "amp init", "amp doctor", "amp login", "工作区初始化", "镇元", "镇元项目", "镇元 cli", "初始化项目", "初始化镇元", "namespace:", "namespace", "amp workflow", "克隆 cspec", "克隆仓库", "clone 仓库", "git clone", "cspec 仓库", "cspec 源码", "拉代码", "拉源码", "拉到本地", "cloudspec-model", "发布预发", "推送远程", "新建API", "create API", "新建API并发布", "初始化并新建API", "API接口定义", "接口定义", "API定义", "新增API", "创建操作", "IDL", "IDL定义", "CloudSpec IDL", "amp policy", "amp domain", "amp error-code", "amp gateway", "策略管理", "域名管理", "错误码管理", "网关".
allowed-tools: Bash, Read, Write, Edit, AskUserQuestion
---

# amp-workflow · amp CLI 全生命周期编排

> **Skill 定位**：把 amp CLI 的「前置配置 → 分支生命周期 → 发布远端」串成一次对话内可完成的流程，同时支持 policy / domain / error-code / gateway 管理。中间「API / 资源设计 / 测试」明确交接给 `cloudspec-*` skill，本 skill 不越界。

## 1. 范围与边界

| 归本 skill | **不归**（明确交接） |
|---|---|
| `amp doctor / config / login / init`（bootstrap 体检与代跳） | `aliyun cspec build / check`（cspec 工具链） |
| `amp branch list / get / create / update / delete / switch` + `amp context set branch` | `cloudspec-operation-edit` / `cloudspec-resource-edit`（API/资源设计） |
| `amp publish daily / pre` + `--dry-run` + 发布前 `branch get` & `api list` 验证 | `cloudspec-test-fix / cloudspec-test-migrate`（测试） |
| `amp policy *`（策略管理） | |
| `amp domain *`（域名管理） | |
| `amp error-code *`（错误码管理） | |
| `amp gateway *`（网关元数据查询） | |
| 失败排查（按 amp 手册第 10 章错误码） | 流控 → `cloudspec-amp-rate-limit` skill |

### 硬约束（违反即拒绝执行）

1. **`amp publish prod` 永不调** —— prod 段处于白屏阶段，本 skill 仅支持 `daily` / `pre`。
2. **禁止在 master/main 分支上做任何编辑操作** —— 所有 cspec 编辑、`amp api` 写操作、publish 必须在 feature 分支进行。检测到当前在 master/main 时，**必须先引导用户创建或切换到 feature 分支**（step 3.3），再继续后续流程。
3. **publish 必走 `--dry-run` 预演** —— dry-run 通过才允许真发；dry-run 失败转 step 7 排查。
4. **所有命令默认带 `-o json`** —— 便于解析 `nextActions` / `error` 字段。
5. **危险操作显式 `--yes`**（branch delete / api delete 等），并先 `--dry-run`。
6. **环境名全小写**（`daily / pre / online`）—— amp 手册明确要求，错大小写会触发 `INVALID_INPUT`。
7. **不写认证信息到 git 或日志** —— amp 常规认证只依赖 BUC 登录；BUC token 不读、不打印、不让用户在对话里粘贴。AK/SK 仅在用户明确要求兼容低频旧链路时作为可选配置，并且只能由用户在本地安全终端处理。

### 设计取向

- **task-oriented**（与 amp CLI 同款）：按"用户要完成的事"分组，不按后端 Action 透出。
- **可脚本化**：所有 Bash 调用统一 `-o json --no-interactive`，删除类追加 `--yes`。
- **代跳但不代决策**：能用文档默认值的（endpoint / openapi-version）skill 自动 set；需要 BUC 登录或工作区参数时，引导用户完成登录或确认参数。

---

## 2. bootstrap 自动代跳

> 这是本 skill 与"普通 routing skill"的核心差异：**用户调用一次本 skill，bootstrap 必须让 `amp doctor` 的常规检查就绪**，再继续后续动作。下面 8 步幂等执行，已就绪的步直接跳过；AK/SK 相关检查只作为低频可选项记录。

### 2.1 amp 二进制体检与版本更新

```bash
amp --version 2>/dev/null || echo "AMP_NOT_INSTALLED"
```

- 输出形如 `amp <version>` → 先保持到最新版本，再继续 step 2.2：

  ```bash
  amp upgrade --check -o json
  amp upgrade --yes --no-interactive
  amp --version
  ```

  如果 `amp upgrade` 因网络、权限或下载源失败，说明失败原因并暂停；不要在旧版本上继续执行分支、clone 或 publish。

- 输出 `AMP_NOT_INSTALLED` 或非零退出 → **不要替用户装**，引导用户在终端跑（来自 amp 手册 3.2），装完立即升级到最新：

  ```bash
  export AMP_DOWNLOAD_URL=https://amp-cli.oss-cn-beijing.aliyuncs.com/0.1.0  # bootstrap 入口，随后必须 upgrade 到最新
  curl -fsSL $AMP_DOWNLOAD_URL/install.sh | sh
  amp upgrade --yes --no-interactive
  amp --version
  ```

  装完让用户回到对话再触发本 skill。

### 2.2 amp doctor 总体检

```bash
amp doctor -o json
```

按 JSON 输出里 `data.checks` 的失败项分别走 step 2.3-2.7。**所有常规缺失项一次代跳完毕后，回到 step 2.8 复查。** 如果只剩 `AKSK_MISSING`，按 step 2.6 记录为可选项，不阻塞常规 bootstrap。

### 2.3 endpoint 缺失

```bash
amp config set endpoint https://ampv2inner-share.aliyuncs.com
```

（手册 4.2 默认值）

### 2.4 openapi-version 缺失

```bash
amp config set openapi-version 2026-04-20
```

（手册 4.2 默认值，业务版本另在 step 2.7 工作区里写）

### 2.5 BUC 登录缺失

按优先级查找：

1. 环境变量 `AMP_BUC_TOKEN` 已存在 —— 直接使用，不打印、不复述。
2. `amp whoami` 输出 token 健康 —— 跳过。
3. 都没有 —— 主动执行浏览器登录，并等待用户完成：

   ```bash
   amp login
   ```

   如果当前环境无法打开浏览器或监听本地端口，提示用户在自己的安全终端运行 `amp login` 后回来继续。不要要求用户在对话里粘贴 BUC token；只有用户已经通过安全方式把 token 放进 `AMP_BUC_TOKEN` 时才使用它。

### 2.6 AK/SK 可选配置（低频）

常规 amp 分支、clone、publish、policy/domain/error-code/gateway 和 flow 操作只要求 BUC 登录，不要求 AK/SK。不要把 `AKSK_MISSING` 当作常规 bootstrap 阻塞项。

仅在以下情况处理 AK/SK：

1. 用户明确要求配置 AK/SK；或
2. 某个低频旧链路命令在 BUC 登录健康后仍明确返回 `AKSK_MISSING`，且用户确认要走该旧链路。

处理规则：

- 可以用 `amp config ak-status` 查看是否已配置，但不要读取或打印任何值。
- 禁止让用户在对话里明文提供 AK/SK。
- 如必须配置，让用户在本地安全终端执行：

  ```bash
  amp config set-ak --access-key-id <你的AK> --access-key-secret <你的SK>
  ```

- AK/SK 未配置不影响继续后续常规流程。

### 2.7 工作区上下文缺失

> **注意**：此步在当前 CWD 执行 `amp init`，主要目的是获取 `project_id` 和 `SshUrl`（供 step 4.1 使用）。**最终的 amp 工作区应在 clone 出来的 cspec 目录内**（step 4.4 会再次 init）。

```bash
ls .amp/context.yaml 2>/dev/null
```

- 文件存在 → 读 `pop_code / version / project_id`，跳过本步。
- 不存在 → `AskUserQuestion` 问用户：

  ```
  Q: 当前目录还不是 amp 工作区。请提供：
  - pop-code（产品 POP Code，例如 ecs / vpc）
  - version（业务产品版本，例如 2014-05-26）
  - 默认分支名（可选，不填则后续手动 amp branch create）
  - 默认 api-name（可选）
  ```

  拿到后执行：

  ```bash
  amp init \
    --pop-code <popCode> \
    --version <version> \
    [--branch <branch>] \
    [--api-name <apiName>]
  ```

  会同时生成 `.amp/context.yaml`（建议入 git）和 `.amp/local.yaml`（建议加 `.gitignore`）。

### 2.8 doctor 复查

```bash
amp doctor -o json
```

JSON 里除可选 `AKSK_MISSING` 外，所有常规 `checks[*].status == "ok"` 才算 bootstrap 完成。任意常规项还红 → 回到对应 step 2.3-2.7，**不要继续 step 3**。如果仅剩 AK/SK 相关检查失败，记录为低频可选项后继续。

---

## 3. 分支管理

> amp 分支是后端逻辑分支（不是 git 分支），用于隔离一组未发布的 API 改动。

### 3.1 列出分支

```bash
amp branch list -o json
```

返回里看 `data.branches[*]` 的 `name / state / description`。

### 3.2 查询单个分支

```bash
amp branch get --branch <name> -o json
```

或上下文已 `amp context set branch <name>`，省略 `--branch`。

### 3.3 创建分支

```bash
amp branch create \
  --branch <name> \
  --description "<msg>" \
  -o json
```

**命名约定**（amp 手册 5.7）：`feature/<short-action>`，例如 `feature/create-user`。

### 3.4 切分支（本地上下文）

```bash
amp context set branch <name>
amp branch switch <name>     # 同步本地默认分支
```

`amp branch switch` 仅本地切换，不调后端。

### 3.5 更新分支元信息

```bash
amp branch update \
  --branch <name> \
  --description "<msg>" \
  [--member-emp-id-list WB01073675,117319] \
  -o json
```

### 3.6 关闭/删除分支

```bash
# 先预演
amp branch delete --branch <name> --dry-run -o json

# 确认无误再真删
amp branch delete --branch <name> --yes -o json
```

---

## 4. 克隆 cspec 源码仓库 + 进入 cspec 工作流

> amp 后端在 `amp init` 时已返回完整 `SshUrl`，本节利用该 URL 直接 clone，无需用户手动拼仓库名。
>
> **硬约束**：amp branch 与 git branch **同名**，但**只能用 `amp branch create` 创建**。直接 `git checkout -b` 创建的 git 分支后端不识别，无法 publish。本 skill 永远先 `amp branch create`，再 `git clone -b`。

### 4.1 获取仓库地址（自动，从 amp init 获取）

`amp init --debug` 调用后端 `GetProjectByPopCodeAndVersion`，响应里包含：

```json
{
  "SshUrl": "git@gitlab.alibaba-inc.com:cloudspec-model/ECS_pop_Ecs_2014-05-26.git",
  "ProjectName": "ECS::pop::Ecs::2014-05-26",
  "ProjectId": 2928938
}
```

**skill 执行步骤**：

```bash
# step 2.7 已经执行过 amp init，这里复用 --debug 输出解析 SshUrl
amp init --pop-code <popCode> --pop-version <version> --debug --no-interactive -o json 2>&1
```

从 debug 输出中提取 `response_body` 里的 `SshUrl` 字段，即为完整 git clone 地址。

**用户只需提供**：
- `pop-code`（产品 POP Code，如 ecs、polardb、OpenAPIExplorer）
- `pop-version`（业务版本，如 2014-05-26）
- 分支名（如 feature/add-tag）

这三个信息在 step 2.7 工作区初始化时已经问过，**不需要额外再问**。

**验证过的示例**：

| pop-code | pop-version | 后端返回 SshUrl |
|---|---|---|
| `ecs` | `2014-05-26` | `git@...cloudspec-model/ECS_pop_Ecs_2014-05-26.git` |
| `polardb` | `2017-08-01` | `git@...cloudspec-model/PolarDB_pop_polardb_2017-08-01.git` |
| `OpenAPIExplorer` | `2024-11-30` | `git@...cloudspec-model/OpenAPIExplorer_pop_OpenAPIExplorer_2024-11-30.git` |

### 4.2 Fallback：用户直接提供仓库信息

当 `amp init` 返回空（后端找不到项目）时，按以下优先级 fallback：

**A. 用户提供 namespace**

namespace 格式：`alicloud.{Product}.{popCode}.v{YYYYMMDD}`

解析规则：
1. 按 `.` 分割得 4 段：`[alicloud, Product, popCode, vDate]`
2. 去掉 `alicloud`
3. `Product` 保持原样（大小写敏感）
4. `popCode` 保持原样
5. `vDate` 去 `v` 前缀，格式化为 `YYYY-MM-DD`
6. 中间连接符为网关类型：`pop`（默认）或 `none`
7. 拼接：`{Product}_{pop|none}_{popCode}_{YYYY-MM-DD}`

示例：`alicloud.AmpV6.amp-2.v20140526` → `AmpV6_pop_amp-2_2014-05-26`

**B. 用户直接提供完整仓库名或 git URL**

直接用，不做解析。

### 4.3 选定本地路径 + git clone

```
Q: clone 到哪里？
- 选项 A: 当前目录 $PWD/<repoName>（推荐）
- 选项 B: 我指定路径
- 选项 C: 已经 clone 过了，告诉你绝对路径，跳过 clone
```

选 C 时跳过 clone，直接进 step 4.4。

```bash
git clone -b <branchName> <SshUrl> <localPath>
```

`<branchName>` = step 3.3 刚 `amp branch create` 出来的同名分支。

**失败诊断**：

| 失败信号 | 原因 | 处理 |
|---|---|---|
| `Permission denied (publickey)` | SSH key 没加到 GitLab | 引导用户去 `https://code.alibaba-inc.com/profile/keys` 加公钥，**不要改 ~/.ssh** |
| `remote: ERROR: ... not found` | 仓库不存在 / 后端 SshUrl 有误 | 让用户复核 pop-code，或走 step 4.2 fallback 手动提供 URL |
| `error: Remote branch <name> not found` | amp branch 还没同步到 git 远端 | `amp branch get --branch <name>` 确认存在 → 等 30s 重试；仍失败则开 issue |

### 4.4 在 cspec 目录初始化 amp 工作区

> **关键**：`amp init` 在**当前工作目录**创建 `.amp/context.yaml`（存 project_id / pop_code / version）。step 2.7 的 init 可能在任意目录执行过，但 **cspec 编辑和 publish 必须在 clone 出来的元数据仓库目录内操作**，所以 **必须在 cspec 目录内再执行一次 `amp init`**。

```bash
cd <localPath>
ls main.cspec resources/ operations/ 2>/dev/null   # 嗅探 cspec 项目结构

# 分支守卫：确认不在 master/main 上
git branch --show-current
# 如果输出 master 或 main → 必须先切到 feature 分支，不允许继续

# 在 cspec 项目目录初始化 amp 工作区（幂等，已存在则更新）
amp init --pop-code <popCode> --pop-version <version> --no-interactive -o json

# 设置当前分支
amp context set branch <branchName>
```

> **⚠️ 分支守卫**：如果 `git branch --show-current` 返回 `master` 或 `main`，**必须停下来**，用 `AskUserQuestion` 让用户提供 feature 分支名，先 `amp branch create` + `git checkout` 再继续。禁止在 master/main 上做任何编辑或发布操作。

初始化后目录结构：

```
<localPath>/
├── main.cspec
├── resources/
├── operations/
├── tests/
└── .amp/                    ← amp init 创建
    ├── context.yaml         ← project_id + pop_code + version（建议入 git）
    └── local.yaml           ← 当前分支等本地状态（建议加 .gitignore）
```

**注意**：`.amp/local.yaml` 包含本地分支状态，不应入 git。建议在 cspec 仓库的 `.gitignore` 中加入 `.amp/local.yaml`。

确认 cspec 项目结构存在后，**本 skill 暂停**，把控制权交出去：

```
仓库 clone + 分支切换完成 ──┐
                            │
                            ▼
               cloudspec-idl-guide  （编辑 .cspec）
                            │
                            ├── cloudspec-operation-edit  / cloudspec-resource-edit
                            ├── cloudspec-build-fix / cloudspec-norm-check-fix
                            └── cloudspec-test-fix / cloudspec-test-migrate
                            │
                            ▼
               aliyun cspec build && aliyun cspec check 全绿
                            │
                            ▼
               回到本 skill step 5 准备发布
```

**衔接规则**：

1. 本 skill 输出一句明确指引：「仓库已 clone 到 `<localPath>`，分支 `<branchName>` 就绪。请在该目录做 .cspec 编辑——按 cloudspec-idl-guide 路由。完成后回到 cloudspec-amp-workflow 做发布。」
2. **元数据生产统一走 cspec 文件**，不使用 `amp api create / update / delete`。生产路径：编辑 `.cspec` → `aliyun cspec build` → `amp publish`。`amp api list / get` 仅用于发布前只读验证。
3. 不要**自动**替用户 `git add / commit / push`——cspec 工具链 + amp publish 会自己处理源码同步。但在 build + check 全绿后，**应提示用户是否要 commit + push 到远端分支以持久化变更**（详见 cloudspec-idl-guide 第五节），用户明确确认后方可执行 git 操作。

---

## 5. 发布前验证清单

> publish 是写类操作 + 跨环境影响，**任何一步失败都不能跳过**。

### 5.1 分支状态

```bash
amp branch get -o json
```

确认 `state` 字段为 `active`（或文档定义的可发布态），`apis` 数组非空。

### 5.2 列出本次要发的 API

```bash
amp api list -o json
```

人工/agent 二次确认这是预期的发布清单——尤其是有 delete / revert 的情况。

### 5.3 dry-run 预演

```bash
amp publish daily --dry-run -o json
```

- 成功（`success: true`）→ step 6 真发。
- 失败 → step 7 排查；**不要忽略 dry-run 失败直接真发**。

---

## 6. 发布执行

### 6.1 daily 发布（默认目标）

```bash
amp publish daily -o json
```

发布完读 JSON 信封的 `nextActions[]` 数组——通常会列出"去哪验证 / 是否需要继续 pre / 文档/SDK 链接"。

### 6.2 pre 发布（按需）

```bash
amp publish pre --dry-run -o json   # 必须先 dry-run
amp publish pre -o json
```

### 6.3 prod ❌ 永不调

本 skill 不支持 `amp publish prod`。如果用户明说要 publish 到 prod，回复："prod 段当前白屏，本 skill 不执行；请通过镇元工作台或 SOP 流程发布。"

---

## 7. 失败排查

amp CLI 错误码统一格式（手册第 10 章）：

```
Error: <可读说明>
Code: <错误码>
Suggestion: <建议>
```

下面列 6 条本 skill 高频遇到的，**严格按 Suggestion 走，不要循环重试**：

| Code | 触发场景 | 修复 |
|---|---|---|
| `ENDPOINT_MISSING` | 没配 endpoint | `amp config set endpoint https://ampv2inner-share.aliyuncs.com` |
| `AUTH_TOKEN_MISSING` | 没登录 / token 过期 | `amp login`；不要要求用户在对话里粘贴 BUC token |
| `AKSK_MISSING` | 低频旧链路要求 AK/SK | 常规流程忽略；仅在用户明确确认旧链路时，让用户本地执行 `amp config set-ak ...` |
| `PROJECT_ID_MISSING` | 工作区缺 project_id | `amp init --pop-code <code> --version <version>` |
| `BRANCH_MISSING` | 命令需要分支但上下文没有 | `amp context set branch <name>` 或命令加 `--branch` |
| `INVALID_INPUT`（env 大小写）| 环境名传成 `Daily` / `PRE` | 改全小写：`daily` / `pre` / `online` |

### `BACKEND_ERROR` 处理协议

amp 手册第 10 章原话："不要循环重试，多数情况是参数/权限/状态问题，重试会放大故障。" 本 skill 严格遵守：

1. 用 `-o json` 重跑一次，记录 `requestId` + 后端 `Code`。
2. 加 `--debug` 看请求/响应细节。
3. 把 `requestId + Code + 命令（脱敏）+ 时间` 报给用户/后端团队。
4. **最多人工再跑一次确认偶发**，仍失败就停下来开 issue。

---

## 8. AI 脚本约定

写 Bash 调用时遵守：

```bash
amp <cmd> [args] \
  -o json \              # 让输出可解析
  --no-interactive \     # 避免阻塞在确认提示上
  [--yes]                # 删除/危险操作显式确认
```

- **永远不读 `~/.amp/tokens/`** —— 该目录文件已 0600，agent 无需也不应直接读。
- **不缓存 amp 状态** —— 每次都用 `amp doctor` / `amp branch get` 拉最新。
- **调用顺序遵循 step 2 → step 3 → 交接 → step 5 → step 6**，不要乱序。

---

## 9. 典型对话剧本

### 剧本 A：用户提供 pop-code + version + 分支名（最常见）

> 用户："帮我初始化 ecs 2014-05-26 的本地工作空间，分支 feature/add-user-tag，改完发 daily。"

skill 内部流程：

1. step 2 bootstrap 自动代跳（doctor → 补 endpoint/openapi-version/BUC 登录/workspace；AK/SK 仅低频可选）。
2. step 2.7 `amp init --pop-code ecs --pop-version 2014-05-26 --debug -o json`
   → 解析 debug 获得 `SshUrl = git@...cloudspec-model/ECS_pop_Ecs_2014-05-26.git`。
3. step 3.3 `amp branch create --branch feature/add-user-tag --description "..."`。
4. step 3.4 `amp context set branch feature/add-user-tag`。
5. step 4.3 AskUserQuestion 问 clone 路径 → `git clone -b feature/add-user-tag <SshUrl> <localPath>`。
6. step 4.4 输出："仓库 clone 到 `<localPath>`，分支 `feature/add-user-tag` 就绪。请按 cloudspec-idl-guide 编辑 .cspec，完成后回来 publish。"
7. **暂停**——等用户回来。
8. 用户："cspec 改完测试也过了，发 daily。"
9. step 5.1-5.3 验证 + dry-run。
10. step 6.1 `amp publish daily`。
11. 输出 `nextActions[]`，结束。

**用户只需说一句话**，提供 3 个关键信息：`pop-code` + `version` + `分支名`。
全程 6-8 个 amp 命令 + 1 次 git clone；bootstrap 缺 BUC 登录时等待一次登录，clone 问一次路径，**仅此两次交互**。

### 剧本 B：用户提供 namespace

> 用户："namespace 是 alicloud.PolarDB.polardb.v20170801，帮我建分支 feature/add-query 并 clone。"

skill 流程：
1. 从 namespace 解析出 `pop-code=polardb`、`version=2017-08-01`。
2. 走 step 2 bootstrap → `amp init --pop-code polardb --pop-version 2017-08-01 --debug` 拿 SshUrl。
3. 后续同剧本 A 步骤 3-11。

### 剧本 C：用户已有本地工作区，只需建分支

> 用户："我在 ~/code/ECS_pop_Ecs_2014-05-26 目录已经有仓库了，帮我建个新分支 feature/fix-param。"

skill 流程：
1. step 2 bootstrap（已有 .amp/context.yaml → 跳过 init）。
2. step 3.3 `amp branch create --branch feature/fix-param`。
3. `cd ~/code/ECS_pop_Ecs_2014-05-26 && git fetch && git checkout feature/fix-param`。
4. 交接给 cloudspec-idl-guide。

---

### 用户输入示例汇总

| 用户说的话 | skill 需要的信息 | 需要额外问的 |
|---|---|---|
| "初始化 ecs 2014-05-26，分支 feature/add-tag" | pop-code + version + branch ✅ | 仅问 clone 路径 |
| "namespace: alicloud.DRDS.polardbx.v20200202，建分支 feature/fix" | namespace + branch ✅ | 仅问 clone 路径 |
| "帮我 clone polardb 2017-08-01 的仓库" | pop-code + version ✅ | 问分支名 + clone 路径 |
| "我要在 ~/code/Ecs 目录开发，分支 feature/x" | 本地路径 + branch ✅ | 无需额外问 |
| "git@gitlab...cloudspec-model/Foo.git 分支 bar" | 完整 URL + branch ✅ | 仅问 clone 路径 |

---
name: cloudspec-amp-workflow
description: |
  镇元/amp 项目的唯一入口技能。无论用户是要初始化项目、新建API、创建操作、拉代码、还是发布，都必须先经过本 skill。本 skill 负责完整生命周期：amp bootstrap → 分支管理 → git clone cspec 仓库 → API/接口定义编辑 → 发布 daily/pre。同时支持 policy/domain/error-code/gateway 管理。
  🔴 路由规则：只要用户消息包含以下任意关键词，必须使用本 skill：初始化、镇元、namespace、发布、publish、clone、拉代码、拉到本地、推送远程、新建API、新增API、API接口定义、接口定义、创建操作、IDL、IDL定义。
  Triggers: "amp", "amp branch", "amp publish", "建分支", "创建分支", "切分支", "分支管理", "发布到 daily", "发布到 pre", "publish daily", "publish pre", "集成发布", "amp 流程", "amp 初始化", "amp init", "amp doctor", "amp whoami", "工作区初始化", "镇元", "镇元项目", "镇元 cli", "初始化项目", "初始化镇元", "namespace:", "namespace", "amp workflow", "克隆 cspec", "克隆仓库", "clone 仓库", "git clone", "cspec 仓库", "cspec 源码", "拉代码", "拉源码", "拉到本地", "cloudspec-model", "发布预发", "推送远程", "新建API", "create API", "新建API并发布", "初始化并新建API", "API接口定义", "接口定义", "API定义", "新增API", "创建操作", "IDL", "IDL定义", "CloudSpec IDL", "amp policy", "amp domain", "amp error-code", "amp gateway", "策略管理", "域名管理", "错误码管理", "网关".
allowed-tools: Bash, Read, Write, Edit, AskUserQuestion
---

# amp-workflow · amp CLI 全生命周期编排

> **Skill 定位**：把 amp CLI 的「前置配置 → 分支生命周期 → 发布远端」串成一次对话内可完成的流程，同时支持 policy / domain / error-code / gateway 管理。中间「API / 资源设计 / 测试」明确交接给 `cloudspec-*` skill，本 skill 不越界。

## 1. 范围与边界

| 归本 skill | **不归**（明确交接） |
|---|---|
| `amp doctor / whoami / config / init`（bootstrap 体检与非认证配置） | `aliyun cspec build / check`（cspec 工具链） |
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
7. **AMP 经可信 wrapper 复用 Jarvis Code PAT，skill 只验证、不接触明文** —— 凭据固定来自真实 OS 账号 home 下 `~/.config/a1/identities/jarvis/auth.yaml` 的 `platforms.code`，由 wrapper 内部读取校验。必须通过 step 2.5 的 `authType=private_token`、`authenticated=true` 及目标项目只读后端验证；不执行或指导 `login`、`logout`，不回退 BUC/个人/TerraformRD，不要求另配 AMP 凭据。不读、不打印、不写入 git/日志、不索取 token 或 AK/SK；认证异常立即停止并报告源凭据运行环境问题。
8. **Jarvis 访问 Code 私库只用 Jarvis private token** —— 只读 API 走
   `bin/a1id -- repo ...`；clone/push 走本 skill 的
   `scripts/jarvis-code-git.sh`。禁止使用 TerraformRD/个人身份、SSH key、`SshUrl`，也禁止把
   token 拼进 URL、remote、命令参数或日志。private token 不可用时 fail-closed，不回退 SSH。

### Jarvis 可信执行覆盖（Agent 必须遵守）

本文后续的 `amp ...` 命令是人工终端的原生 CLI 参考。Jarvis/
Claude/Codex **禁止直接执行 raw amp**，必须使用仓库内可信 wrapper：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <absolute-workspace-or-model-repo> \
  <allowed-amp-subcommand-and-arguments>
```

- `--repo-root` 必须是绝对路径，每一次调用都显式传入；不得假设上一个
  Bash 工具调用里的 `cd` 会在下一次保留。
- clone 前的 `init` / `branch create` 使用该任务独占的 bootstrap workspace；
  clone 后的 `init` / `context` / `branch switch` / `publish` 使用含
  `main.cspec` 的 model Git 仓库根。
- Jarvis 执行 `branch create` 时必须显式传
  `--project-id <context.project_id>`；不允许省略 scope，也不允许用
  `--pop-code/--pop-version` 代替 project-id 执行远程建分支。
- 任何远程写操作必须由 wrapper 将实际 AMP `project_id` + feature branch
  与任务 fence 绑定；禁止使用含 `/` 的混合 `pop-code/version` 字符串，
  也不得用可歧义 pop-code 覆盖已确定的 project-id。
- `publish` 额外要求任务开始前已建立的 Git baseline；任何
  `operations/` diff 仍永久阻断，daily/pre 仍执行 dry-run 后二次授权。
- wrapper 返回的 `reason=<code>` 是可恢复性判定真源。错误 repo/scope
  被拒绝后应修正参数重试，不应要求需求方“再评论一次”。

### 设计取向

- **task-oriented**（与 amp CLI 同款）：按"用户要完成的事"分组，不按后端 Action 透出。
- **可脚本化**：所有 Bash 调用统一 `-o json --no-interactive`，删除类追加 `--yes`。
- **代跳但不代决策**：能用文档默认值的（endpoint / openapi-version）skill 自动 set；工作区参数不足时向用户确认。AMP 由可信 wrapper 复用运行环境已有的 Jarvis Code PAT；此复用已经用户授权，不需要用户为 AMP 手动配置第二份凭据。skill 仅验证本地认证状态和目标服务权限，异常交运行环境维护。

---

## 2. bootstrap 自动代跳

> 这是本 skill 与"普通 routing skill"的核心差异：**用户调用一次本 skill，bootstrap 必须让 `amp doctor` 的常规检查就绪**，再继续后续动作。下面各步幂等执行，已就绪的非认证配置直接跳过；认证始终按 step 2.5 验证 wrapper 的 Jarvis PAT 复用及目标项目后端访问，不另配或复制 AMP 凭据。`--version` 无需读取 token；请求签名按 step 2.6 保留既有 OAuth 配置并验证，不自动配置 AK/SK。

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

按 JSON 输出里 `data.checks` 的失败项分别走 step 2.3-2.7；无论 doctor 是否通过，都必须执行 step 2.5。只补齐非认证配置，认证失败立即报告运行环境问题；不得执行输出中建议的登录、登出或 BUC 回退动作。**配置就绪后，回到 step 2.8 复查。** `AKSK_MISSING` 按 step 2.6 排查既有 OAuth 请求签名来源，不能仅因 PAT 本地认证成功而忽略。

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

### 2.5 验证 Jarvis PAT 复用

wrapper 在内部复用同仓 `scripts/jarvis-code-git.sh` 的 Code 凭据读取校验逻辑，固定从真实 OS
账号 home 下 `~/.config/a1/identities/jarvis/auth.yaml` 的 `platforms.code` 获取 PAT。
仅通过 AMP 子进程环境 `AMP_PRIVATE_TOKEN` 传递：模型不读取明文，token 不进入 argv/日志，
不复制写入 `~/.amp` 凭据。每次调用使用自动清理的临时 `AMP_HOME`，其配置只保留原当前 profile
的 endpoint、OpenAPI 版本、HTTP 超时，以及 `credentials` 的严格非 secret 白名单
`source`、`oauth_profile`、`oauth_site`（仅字符串，`source` 仅 `local` / `oauth`）；
固定 `auth.type=private_token` 并关闭 HTTP debug。不复制其它 profile、token、AK/SK、
`auth.token_file` 或任何 secret 文件，不修改真实配置。缺失的签名 metadata 不自动补齐，
由原生 AMP 报错。AMP 子进程的 `HOME` 仍是真实 OS 账号 home，让原生 OAuth 按既有
profile/site 获取请求签名凭据；这不代表回退 BUC 身份认证。除无需读 token 的 `--version`
外，所有 AMP 调用的身份认证都使用 Jarvis PAT；不依赖 AMP default profile 已缓存的
BUC/私有 token 或 ambient token，也不接受它们作回退。

通过可信 wrapper 执行只读状态检查：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <absolute-workspace-or-model-repo> whoami
```

wrapper 自动附加 `-o json --no-interactive`；调用 wrapper 的 `whoami` 时不再重复传这些参数。

仅在命令成功、JSON 可解析，且返回的认证状态同时满足以下两项时，本地状态检查通过：

- `authType` 严格等于字符串 `private_token`；
- `authenticated` 严格等于布尔值 `true`。

字段缺失、类型不符、其他认证类型（包括 BUC）、未认证或命令失败均视为运行环境认证异常，立即停止后续分支、clone、publish 等操作。仅报告脱敏后的错误码和状态，不输出原始凭据。

**`whoami` 的本地成功不等于后端认证/权限有效**。还须对已验证的目标项目执行只读查询：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <absolute-workspace-or-model-repo> branch list --project-id <verified-project-id>
```

缺少已验证的 `project-id` 时，先按 step 2.7 定位项目，再回到此处；不得猜测 ID 或用其它项目替代。
只有目标项目查询成功（包括后端业务结果成功）才确认 AMP 后端有效，允许继续分支、clone、publish。

源凭据缺失、非 `private_token`/内容错误、文件权限不安全或 HTTP 401/403 均 fail-closed；
由运行环境维护源 `auth.yaml` 及对应服务权限后重新验证，不自动重置凭据。
不检查 `AMP_BUC_TOKEN` 是否存在来代替状态验证，也不以 doctor 通过代替上述检查。
skill 不执行或指导 `amp login`、`amp logout`，不打开浏览器，不回退 BUC/个人/TerraformRD；
此同源复用已经用户授权，不再要求用户手动为 AMP 配置第二份凭据。

### 2.6 POP 请求签名与既有 OAuth 配置

`private_token` 是身份认证方式，不等于免除 POP 请求签名；即使 `whoami` 返回
`authenticated=true`，普通 `api get` / `branch list` 等后端调用仍可能需要 AK/SK 签名。
不能断言常规流程不需要 AK，也不能把实际请求的 `AKSK_MISSING` 当作可忽略的低频旧链路提示。

原当前 profile 已有 `credentials.source=oauth` 时，临时 profile 保留既有
`oauth_profile` / `oauth_site`，由原生 AMP 通过该 OAuth 来源获取签名凭据，无需另配 PAT。
先确认这三个非 secret metadata 未在投影中遗漏；不读取签名 secret 文件，不自动配置、
打印或索取 AK/SK，不修改真实配置，不触发登录或 BUC 回退。metadata / 签名来源缺失时
沿用原生错误，停止并报告运行环境问题；签名就绪不能替代 step 2.5 的 PAT 与后端权限验证。

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

重新执行 step 2.5 的 wrapper `whoami`，确认 `authType=private_token` 且 `authenticated=true`，并对已验证的目标 `project-id` 执行只读 `branch list` 确认后端成功；所有常规 `checks[*].status == "ok"` 才算 bootstrap 完成。AMP 与 Code 使用同源 Jarvis PAT，但此处通过不替代 step 2.9 的 Code 服务权限检查。认证或签名异常立即停止并报告源 `auth.yaml`、既有 OAuth 签名来源或服务权限运行环境问题，不执行 doctor 建议的登录动作。常规项还红 → 回到对应 step 2.3-2.7，**不要继续 step 3**；`AKSK_MISSING` 按 step 2.6 处理，不能标为可选后继续。

### 2.9 Jarvis Code private token 体检

Jarvis/数字人环境在 clone 前仍须验证 Code 服务权限；AMP 与 Code 复用同一份
`~/.config/a1/identities/jarvis/auth.yaml` 的 `platforms.code` PAT，不是两套独立凭据。
两端服务权限仍分别验证：AMP 后端成功不能替代 Code 仓库访问成功，反之亦然。

```bash
bash <skill-dir>/scripts/jarvis-code-git.sh check
bin/a1id -- repo view <group/repo> -f json
```

缺凭据、凭据失效或仓库访问失败时，停止 clone/push 并报告源 `auth.yaml` 或 Code 仓库权限运行环境问题。skill 不执行或指导 Code 登录、登出或另配凭据，不读取或索取 token，不得借 TerraformRD、个人身份或 SSH 绕过。

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
  --project-id <context.project_id> \
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

## 4. 用 Jarvis private token 克隆 cspec 源码 + 进入 cspec 工作流

> amp 后端在 `amp init` 时返回项目 ID、项目名和 `SshUrl`。Jarvis 只用这些信息确定
> `<group/repo>`，**不得执行或保存 `SshUrl`**；实际 Git transport 固定为 Jarvis Code
> private-token HTTPS。
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

从 debug 输出提取项目 ID，并把 `SshUrl`/`ProjectName` 仅解析为 `<group/repo>`。随后用 Jarvis
身份 point-read 仓库和目标分支：

```bash
bin/a1id -- repo view <group/repo> -f json
bin/a1id -- repo branch list --repo <group/repo> \
  --keyword <branchName> --per-page 100 -f json
```

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

只解析出 `<group/repo>`；丢弃 URL 中的 scheme、host 与认证信息。Jarvis 禁止直接使用用户给的
URL clone，更禁止接受带 token 的 URL。

### 4.3 选定本地路径 + private-token HTTPS clone

```
Q: clone 到哪里？
- 选项 A: 当前目录 $PWD/<repoName>（推荐）
- 选项 B: 我指定路径
- 选项 C: 已经 clone 过了，告诉你绝对路径，跳过 clone
```

选 C 时跳过 clone，直接进 step 4.4。

```bash
bash <skill-dir>/scripts/jarvis-code-git.sh clone \
  <group/repo> <branchName> <absolute-localPath>
```

`<branchName>` = step 3.3 刚 `amp branch create` 出来的同名分支。

helper 从 `~/.config/a1/identities/jarvis/auth.yaml` 读取 `platforms.code` private token，并仅通过
`GIT_ASKPASS` 交给 Git；remote 保持无凭证 HTTPS URL。禁止自行读取/打印 token，禁止临时改
全局 `credential.helper`，禁止把 token 写进 `.git/config`。

**失败诊断**：

| 失败信号 | 原因 | 处理 |
|---|---|---|
| `jarvis Code auth ... missing` | 共享源 `auth.yaml` 无 Jarvis Code private token | 停止并报告源凭据运行环境问题；不另配 AMP 凭据、不执行或指导登录，不回退 SSH |
| `HTTP 401/403` | 共享 PAT 失效或无 Code 仓库权限 | 停止 clone/push，报告源 `auth.yaml` 或 Code repo ACL 运行环境问题；AMP 成功不代表 Code 可用，不自动重置或切身份 |
| `Permission denied (publickey)` | 错误地进入了 SSH 路径 | 停止并改用 `jarvis-code-git.sh`；**不要配置或借用 SSH key** |
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
3. 不要**自动**替用户 `git add / commit / push`——cspec 工具链 + amp publish 会自己处理源码同步。但在 build + check 全绿且已获授权后，push 必须继续使用 Jarvis private-token helper：

   ```bash
   bash <skill-dir>/scripts/jarvis-code-git.sh push \
     <model-repo-dir> <group/repo> HEAD refs/heads/<feature-branch>
   ```

   helper 永久拒绝 master/main；不得改用 SSH 或把 token 持久化到 remote。

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

下面列出本 skill 高频遇到的错误，**按下表处理，不要循环重试**。`Suggestion` / `nextActions` 中的登录、登出、BUC 回退、另配或自动重置凭据建议不执行；认证故障统一按 step 2.5 / 2.9 报告源 `auth.yaml` 或对应服务权限的运行环境维护问题：

| Code | 触发场景 | 修复 |
|---|---|---|
| `ENDPOINT_MISSING` | 没配 endpoint | `amp config set endpoint https://ampv2inner-share.aliyuncs.com` |
| `AUTH_TOKEN_MISSING` / wrapper 凭据校验失败 | 源 `auth.yaml` 缺失、非 `private_token`/内容错误或权限不安全 | fail-closed；交运行环境维护 Jarvis `platforms.code` 源凭据，不另配 AMP 凭据；修复后重跑 step 2.5 |
| `HTTP 401/403` | 共享 PAT 失效或无 AMP 目标项目权限 | 停止并报告源 `auth.yaml` 或 AMP 服务权限问题；不回退 BUC/个人/TerraformRD、不自动重置；修复后重跑目标项目只读验证 |
| `AKSK_MISSING` | POP 请求签名来源缺失，PAT 本地认证成功也可能发生 | 按 step 2.6 检查既有 OAuth 非 secret metadata 是否完整投影；停止并报告签名来源问题，不自动配置或打印 AK，不回退 BUC |
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

1. step 2 bootstrap 体检（doctor → 补 endpoint/openapi-version → wrapper 复用 Jarvis PAT，保留既有 OAuth 签名配置，whoami 验证 `authType=private_token` 且 `authenticated=true` → workspace → 目标项目只读 branch list 成功）。认证或请求签名异常立即报告源凭据、既有 OAuth 签名来源或服务权限运行环境问题，不进入后续动作。
2. step 2.7 `amp init --pop-code ecs --pop-version 2014-05-26 --debug -o json`
   → 解析 debug 获得 `SshUrl = git@...cloudspec-model/ECS_pop_Ecs_2014-05-26.git`。
3. step 3.3 `amp branch create --project-id <context.project_id> --branch feature/add-user-tag --description "..."`。
4. step 3.4 `amp context set branch feature/add-user-tag`。
5. step 4.3 AskUserQuestion 问 clone 路径 → `git clone -b feature/add-user-tag <SshUrl> <localPath>`。
6. step 4.4 输出："仓库 clone 到 `<localPath>`，分支 `feature/add-user-tag` 就绪。请按 cloudspec-idl-guide 编辑 .cspec，完成后回来 publish。"
7. **暂停**——等用户回来。
8. 用户："cspec 改完测试也过了，发 daily。"
9. step 5.1-5.3 验证 + dry-run。
10. step 6.1 `amp publish daily`。
11. 输出 `nextActions[]`，结束。

**用户只需说一句话**，提供 3 个关键信息：`pop-code` + `version` + `分支名`。
bootstrap 验证 wrapper 复用的 Jarvis PAT 与 AMP 目标项目权限；clone 前另验同源 PAT 的 Code 仓库权限，并在路径未确定时确认 clone 路径。无需第二份 AMP 凭据；全程不执行或指导登录，认证或权限异常时停止并报告源 `auth.yaml` 或服务权限运行环境问题。

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

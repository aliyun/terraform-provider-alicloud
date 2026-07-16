# Jarvis 接入自动化服务台工作项池与交付上下文设计

> Aone: [84215393](https://project.aone.alibaba-inc.com/v2/project/2100304/req/84215393)
>
> 状态：需求范围与方向已由用户确认；本文固化实现边界与验收口径。

## 1. 背景

Jarvis 已能扫描并处理 Terraform、MCP Server、Cloudspec 和内部工具类工作项，但自动化服务台目前只有一个不完整的后端 workspace 登记：没有独立 Aone 池、关联仓库、交付 runbook 和分诊路由。因此，定时扫描不会接收自动化服务台工作项，直接给出工单时也缺少稳定的研发与发布上下文。

本次目标是让 Jarvis 能接收、分诊、研发并交付自动化服务台需求，同时保持两条安全边界：

- 定时扫描只处理明确指派给 `open-jarvis` 的工作项；
- 正式发布继续是人工门，Jarvis 最多自动推进到 CR/MR 和预发验收。

## 2. 已锁定范围

### 2.1 纳入范围

- Aone 项目：自动化服务台，project ID `1091779`。
- 默认扫描负责人：`open-jarvis`，账号 `WORKER_1782379562571`。
- 主应用：`aliyun-automation-platform`，app ID `172823`，repo ID `2156624`。
- 交付流水线：预发 `66`、正式 `67`。
- 关联仓库：
  - 后端 `aliyun-automation-platform/aliyun-automation-platform`
  - 前端 `aliyun-api/iac-service`
  - 运行时 `opensource-tools/iac-service-runtime`
  - 集成测试 `aliyun-automation-platform/automation-function-test`
  - 公开 API `cloudspec-model/IaCService`
  - 内部 API `cloudspec-model/IaCService-inner`
- SLS 排障入口：生产 `systemlog-prod`、预发 `systemlog-pre`、POP 报警日志，并复用 `sls-log-query-aliyun-automation-platform` skill。

### 2.2 明确排除

- 不接入 Agent 门户、AgentRuntime、PlayGround、FC sandbox、WebSocket/STS 编排或 `aliyun-automation-agent` 等 Agent 链路。
- 不固化历史 feature 分支、临时 flowId 或单个历史需求。
- 不自动合并 CR/MR，不自动触发正式流水线。

当一个平台需求确实需要修改 Agent 仓库时，应拆成独立关联工作项并切换到现有 Agent 交付链路，不能把两套应用的发布绑定成一次交付。

## 3. 配置模型

### 3.1 独立工作线与 Aone 池

在 `config/pools.json` 中新增独立工作线 `automation_platform`，展示名为“Automation Platform”；新增同名池：

```json
{
  "project": 1091779,
  "name": "自动化服务台",
  "line": "automation_platform",
  "assignee": "WORKER_1782379562571",
  "dev": true
}
```

定时扫描继续复用现有 per-pool assignee 机制。`scan.sh` 会为该池的 req/bug/task 请求附加 `assignedTo=WORKER_1782379562571`，因此不会自动消费团队其他成员的工作项。用户直接给出 Aone URL 或 ID 时走单工单入口，不经过 scan 过滤，仍可读取和处理任意负责人工作项。

状态映射使用 2026-07-13 从项目 `1091779` 实时读取的枚举：

| 工作项类型 | 认领后状态 | 完成状态 |
|---|---|---|
| 产品类需求 | 开发中 | 已发布 |
| 功能缺陷 | Open | Fixed |
| 线上问题 | Open | Fixed |
| 任务 | 处理中 | 已完成 |

终态过滤覆盖上述完成态、取消态以及 Bug 的 Closed/Won'tfix/Worksforme/Duplicate/Invalid/External/ByDesign，防止定时扫描反复派发已经结束的工作项。

### 3.2 应用与仓库 workspace

扩充现有 `automation_platform` workspace：补充 pool、project、pipelines 和 delivery reference。其余五个仓库各自登记为独立 workspace，使用本机实际目录名作为 `repo`，共同绑定 `automation_platform` 池。

不新增私有绝对路径；本地目录仍由 `bootstrap/workspace.sh`、`JARVIS_WORKSPACE_ROOT` 和 gitignored 的 `workspaces.local.json` 解析。仓库登记只保存 git URL、默认分支、职责和可验证的 build/test 操作。

API workspace 与后端 workspace 分离。涉及请求、响应、POP Action 或资源定义变化时，交付 runbook 要求同步检查公开 API；仅内部接口变化时检查内部 API。API 定义修改后执行 `cloudspec build`，避免 Java 后端已支持但契约、文档和前端类型未更新。

## 4. 分诊与交付链路

### 4.1 路由

`aone-triage` 增加 `1091779 automation_platform` 的显式路由，并新增 `delivery-aliyun-automation-platform.md`。路由优先级放在泛化的 MCP/Agent 关键词之前，防止包含“自动化”或“IaC”字样的平台工作项误入 Agent/Cloudspec 链路。

`loops/adhoc-intake.md` 的候选池和 workspace 说明同步改为包含 `automation_platform`，并删除已经不再是入箱池的 `cloudspec` 候选描述。

### 4.2 Delivery reference

交付 reference 固化以下流程：

1. 读取工作项全文并判断后端、前端、runtime、集成测试、公开 API、内部 API 中哪些仓库受影响。
2. 按每个受影响仓库各自的主干创建 worktree；禁止直改 master。
3. 缺陷必须补回归测试；跨 API 契约变更必须同步 schema 仓库。
4. 创建/关联 CR 或 MR，并立即用 `wrap.sh sync` 回贴 Aone。
5. 运行仓库级 build/test，再提交预发流水线 `66`。
6. 预发成功后结合业务验证和 SLS 日志确认结果；等待用户明确验收。
7. 只有收到用户正式发布授权后，才允许提交流水线 `67`。
8. CR/MR 合入且发布闭环后再 finish 和清理 worktree；否则 done + release，保留人工合并门。

SLS 排障不复制凭据或查询脚本，只引用现有 skill 和三个 logstore。生产问题优先查 `systemlog-prod`，预发验证查 `systemlog-pre`；结构化字段查询无结果时，按已有经验回退到短关键词全文搜索。

## 5. 运行时数据流

### 5.1 定时扫描

```text
bridge ScanScheduler
  → bootstrap/scan.sh 动态遍历 config/pools.json
  → project 1091779 + assignedTo=open-jarvis
  → 标准 scan item(pool=automation_platform, pool_project=1091779)
  → 通用 DispatchPool
  → aone-triage 路由到平台 delivery reference
  → workspace 解析 / worktree / CR·MR / 预发
  → 人工验收与正式发布门
```

`scan.sh`、bridge 派发、revisit、claim、wrap 和 reconcile 都已经动态消费池配置，本次不修改 bridge 业务代码。部署环境若显式设置了 `JARVIS_DISPATCH_POOLS`，运维侧需要将 `automation_platform` 加入白名单；未设置时现有代码默认放行所有配置池。

### 5.2 直接工单入口

```text
用户提供 Aone URL/ID
  → aone-get 读取 space=1091779
  → aone-triage 显式路由
  → 平台 delivery reference
```

该入口不依赖工作项当前负责人，满足人工临时委派、协查和历史问题处理场景。

## 6. 看板与可观测性

`board.sh` 已经动态保留任意 pool，但 `board-html.sh` 当前写死旧四池，新增池会在前端默认过滤中消失。本次将池列表改为从 `config/pools.json` 动态生成，并为未知工作线提供稳定的颜色 fallback。这样以后新增池不再需要手工修改 HTML 常量。

看板只展示 Aone 状态，不改变扫描或派发策略。SLS 继续用于交付后的运行结果验证，两者职责分离。

## 7. 错误处理与安全边界

- Aone 项目无权限或单 category 查询 403：沿用 scan 的按 category/池跳过逻辑，不让一个池拖垮整轮扫描，并在日志中保留可诊断信息。
- workspace 未登记或本地仓库无法解析：停止开发并按 `missing_capability` 升级，不臆造路径。
- CR/MR、构建或预发失败：回贴 Aone，保留 claim/bookend 审计，不进入正式发布。
- API 契约与后端实现不一致：验收失败，不以“Java 通用结构可以承载”为理由跳过 API 定义更新。
- 正式发布、合并和跨应用 Agent 改动：始终需要人工授权或拆单。

## 8. 测试与验收

### 8.1 配置测试

- `test/pools_config_test.sh`
  - 池数从 4 更新为 5，工作线从 2 更新为 3；
  - 断言 project `1091779`、独立 line、`open-jarvis` assignee；
  - 断言项目真实的 per-type progress/done 状态；
  - 同时修正基线中 `api_toolkit.done_status["需求"]` 与实际 `"产品类需求"` 键不一致的问题。
- `test/workspaces_config_test.sh`
  - 断言主应用、六个仓库、pipeline、delivery、默认分支和关键 ops。

### 8.2 扫描与看板测试

- 在 `test/scan_test.sh` 增加真实池回归：project `1091779` 的查询必须带 `assignedTo=WORKER_1782379562571`，输出必须标记正确 pool/project。
- 为 `board-html.sh` 增加动态池回归，确保 `automation_platform` 出现在过滤器、统计和默认可见列表中。
- bridge 已有“未知新池默认可派”通用测试，不新增特定池分支。

### 8.3 Skill 与镜像测试

- 在 canonical `.claude/skills/aone-triage` 中修改路由并新增 delivery reference。
- 使用 `bootstrap/mirror.sh to-codex` 同步 `.agents/skills`，不手工维护两份不同内容。
- 更新 `test/aone_triage_templates_sync_test.sh` 的 reference 清单，并运行 `bootstrap/mirror.sh check`。

### 8.4 完成标准

- 目标 JSON、shell 和镜像测试全部通过；
- 通过 stub scan 证明只扫描指派给 `open-jarvis` 的自动化服务台工作项；
- 通过直接 ID 的路由断言证明不受 assignee 限制；
- 看板能展示新池；
- 交付配置和 workspace 不接入 Agent 链路，也不固化历史分支或临时 flowId；
- 提交只存在于功能 worktree 分支，通过 CR/MR 等待人工合入。

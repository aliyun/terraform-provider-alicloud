# 交付链路：aliyun-automation-platform

> 适用于自动化服务台 / IaCService 产品研发与交付。任何 tracked 文件改动都在独立 worktree 和功能分支完成，经 CR/MR 评审合入；禁止直接 push 或 merge master。

## 路由与边界

- Aone project：`1091779`；池 key：`automation_platform`。
- 主应用：`aliyun-automation-platform`；app：`172823`；后端 repo ID：`2156624`。
- 固定流水线坐标：`prestage 66`，`prod 67`。
- 定时扫描只接 `assignedTo=WORKER_1782379562571`（open-jarvis）的工作项；**直接 URL/ID 处理不受 assignedTo 限制**，按工作项的 `space=1091779` 进入本链路。
- 输入已经是项目 `1091779` 的工作项时直接复用，不为开发重复创建 Aone；CR/MR、验证结果和交付状态回填原工作项。

本链路只覆盖自动化服务台产品的后端、前端、IaC runtime、集成测试、公开 API 和内部 API。明确排除 Agent portal（Agent 门户）、AgentRuntime、PlayGround、FC sandbox、WebSocket/STS 编排以及 `aliyun-automation-agent`；这些仍走 `delivery-aliyun-automation-agent.md`，不能因“自动化”字样混入本池。

只记录上述稳定坐标；不要把历史 feature branch、临时 flow ID 或单次流水线实例 ID 沉淀为长期配置。

## Repository selection

| 变更面 | `WORKSPACE_KEY` | 仓库职责 | 必做验证 |
|---|---|---|---|
| 后端服务、任务编排、持久化、POP 实现 | `automation_platform` | `aliyun-automation-platform/aliyun-automation-platform` | 读取 `config/workspaces.json` 中注册的 build/test ops |
| 前端页面与交互 | `automation_platform_frontend` | `aliyun-api/iac-service` | 注册的 build/test/lint ops |
| IaC 执行运行时 | `automation_platform_runtime` | `opensource-tools/iac-service-runtime` | 注册的 Go build/test ops |
| 端到端与集成验证 | `automation_platform_function_test` | `aliyun-automation-platform/automation-function-test` | 注册的 Maven test op |
| 公开 POP API 契约 | `automation_platform_api` | `cloudspec-model/IaCService_pop_IaCService_2021-08-06` | `cloudspec build` |
| 内部 API 契约 | `automation_platform_api_inner` | `cloudspec-model/IaCService-inner_pop_IaCService-inner_2021-09-01` | `cloudspec build` |

涉及 request/response、POP Action、错误码或资源定义时，先判断公开还是内部 API，选择对应 API workspace，并运行 `cloudspec build`。若契约和实现都要改变，API 仓与后端仓分别开分支、验证和 CR/MR，不把跨仓改动塞进一个工作区。

## Worktree、验证与 CR/MR

由配置解析仓库，禁止猜本机绝对路径：

```bash
WORKSPACE_KEY=automation_platform_api  # 按上表替换
WORKSPACE_DIR="$(bash bootstrap/workspace.sh dir "$WORKSPACE_KEY")"
cd "$WORKSPACE_DIR"
```

先在默认分支执行 `git pull --ff-only`，再创建 `worktree-<slug>` 功能分支和 worktree。按 `config/workspaces.json` 的 `ops` 运行该仓 build/test/lint；接口变更同时补契约生成与相关集成测试。提交不得带 AI 署名或内部工单/客户信息。

创建或关联 CR 时复用原工作项 ID：

```bash
bin/a1id -- app link 172823
bin/a1id -- app cr create "<3-5 行变更摘要>" --branch <branch-suffix> --workitem-ids "$WORKITEM_ID"
```

CR/MR 一经创建，立即用 `bootstrap/wrap.sh sync` 把内部评审链接和当前验证结果回填原工作项。未合入 master 时只能 release claim，不能 finish 工作项。

## Prestage 与 production

在 app `172823` 下提交预发：

```bash
bin/a1id -- app link 172823
bin/a1id -- app cr submit "$CR_ID" --pipeline-id 66
bin/a1id -- app pipeline status --pipeline-id 66
```

预发部署成功后，结合业务用例和下方 SLS 上下文验证。到这里必须停下等待验收；不能因为 pipeline 绿色就自动进入正式。

**正式发布必须取得明确的人工批准**。只有仓库主人明确确认预发验收通过后才允许：

```bash
bin/a1id -- app cr submit "$CR_ID" --pipeline-id 67
bin/a1id -- app pipeline status --pipeline-id 67
```

正式仍处于人工卡点、CR/MR 未合并或流水线未整体成功时，不得把 Aone 标为完成，也不得清理仍需排障的 worktree。

## 发布冲突安全

- 冲突优先在当前功能分支合入兄弟分支、解冲突并重新 submit；**永久禁止** `app pipeline exit-cr` 和等价的 `app pipeline quit`，因为两者会退出最新实例中的全部 CR。`bin/a1id` 与 PreToolUse 会在真实 a1 前硬拒绝。
- 确需撤回当前 CR 时，只允许已 claim 对应 Aone 的根 Worker 在 CR worktree 内执行 `bin/a1id -- app cr quit <cr-id> --pipeline-id <id>`。执行层核对 CR 工作项、origin、分支、指定流水线最新实例成员，并在网络查证前后复验同一 task/session/fence；无法证明归属即 fail closed。直接 `a1 app cr quit` 同样被拒绝。
- 该护栏针对 Jarvis 正常工具入口的误操作和 wrapper 绕过；它不构成同一 UID 恶意本地代码的密码学隔离，因为后者本就能读取本机凭据并直调 a1。

## SLS 诊断与验收

必须使用现有 `sls-log-query-aliyun-automation-platform` skill，不在本文件复制凭证或临时查询脚本：

| 场景 | SLS 上下文 |
|---|---|
| 预发部署、业务回归 | `systemlog-pre` |
| 正式环境、线上回归 | `systemlog-prod` |
| POP 请求与报警 | `aliyun-automation-platform-1252907582134651-pop-aliyun-cn` |

按 workitem 里的时间窗、Action、RequestId/traceId 和关键业务字段缩小范围。预发结论至少包含功能验证结果和对应日志证据；线上问题先查 prod/POP 上下文，再决定落后端、API 契约、runtime 或前端仓。

## 完成条件

- 选仓与变更面一致，相关仓库的注册 ops 全部通过。
- API 契约变更已在对应 IaCService API 仓运行 `cloudspec build`，并验证调用侧兼容性。
- CR/MR 链接、预发结果和 SLS 证据已经回填原工作项。
- production 仅在明确人工批准后执行，且流水线与业务验收均成功。
- 没有引入任何 Agent 链路仓库或范围。

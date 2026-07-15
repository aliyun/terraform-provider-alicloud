---
name: mcp-server-playground-acceptance
description: >-
  Use when Alibaba Cloud MCP Server, AutomationAgent/AgentAutomation, IaCService, cloudspec/OpenAPI MCP Server,
  RunIaC, API query tools, or related service changes need end-to-end testing, realistic user-flow
  simulation, or manual experience validation through pre-agent or agent Playground. Also use the
  TerraformAgent Playground when validating MCP Server Core with a managed MCP token, token-level
  safety policy, RunIaC HITL, InitializeApiMcpServerConnection, or GetLatestMcpToken.
---

# MCP Server Playground 验收

真实验收 = 新 Playground 会话 + 实际 MCP 工具调用 + 可回填证据。本地单测、部署成功、口头解释都只是前置条件。

## 配合技能

- 有 Aone URL/ID 时,先用 `aone-triage` 读单、claim、回填和 release。
- 通过浏览器打开 Playground 实测;不要用本地脚本或后端日志替代 UI 链路。
- 需要上传 HTML 验收报告时,先用仓库内已有的 `html-report-preview` skill;报告包含截图时必须遵循其 Image Handling 约束。
- 涉及 cloudspec / OpenAPI MCP Server / RunIaC / API 查询工具时,默认参考 `aone-triage/references/delivery-cloudspec.md` 的 app 与流水线坐标。
- 涉及 AutomationAgent 时,默认参考 `aone-triage/references/delivery-aliyun-automation-agent.md` 的 app 与流水线坐标。
- 涉及托管 MCP token、Token 级安全策略、RunIaC HITL、`InitializeApiMcpServerConnection` 或 `GetLatestMcpToken` 时,必须读取并执行 [TerraformAgent Playground 验收](references/terraform-agent-playground.md)。

## 目标环境

| 环境 | URL | 使用条件 |
|---|---|---|
| 通用预发 | `https://pre-agent.aliyun-inc.com/playground` | 不依赖托管 token 或 Token 级策略的常规 MCP 验收 |
| TerraformAgent 预发 | `https://pre-agent.aliyun-inc.com/terraform-agent/playground` | 托管 MCP token、Token 级策略、RunIaC HITL、`InitializeApiMcpServerConnection`、`GetLatestMcpToken` |
| 通用线上 | `https://agent.aliyun-inc.com/playground` | 仅用户明确要求线上验证时使用 |

线上验证不能顺手做;必须有用户明确指令。不要用通用 Playground 的 settings 页面代替 TerraformAgent 链路;它不提供同一条托管 MCP token 初始化和 Token 级策略链路。RunIaC 按本次验收的实际测试场景选择 action。

## 预发门禁

1. 先读 Aone/CR 评论,确认本轮要验证的 CR、环境和需求口径。
2. 按变更面选择门禁:
   - AgentAutomation / TerraformAgent 页面、托管 token 或策略初始化变更:app `283346`,预发 pipeline `66`。
   - Cloudspec / MCP Server Core / RunIaC / HITL 执行变更:app `260634`,预发 pipeline `420`。
   - 同时跨 AgentAutomation 与 Cloudspec:两边都查,任一侧未就绪都不能开始验收。
3. 分别查所需应用的最新预发实例:
   ```bash
   bin/a1id -- app pipeline status --app 283346 --pipeline-id 66 --format json
   bin/a1id -- app pipeline status --app 260634 --pipeline-id 420 --format json
   ```
4. 如果最新实例仍是 `RUNNING` / `SCHEDULING` / `PENDING` / `CREATED`,继续等:
   ```bash
   bin/a1id -- app pipeline status --app 283346 --pipeline-id 66 --wait-until-settled --format json
   bin/a1id -- app pipeline status --app 260634 --pipeline-id 420 --wait-until-settled --format json
   ```
5. 只有“预发部署”阶段为 `SUCCESS` 才能开始 Playground 验收。后续人工验证阶段处于 `WAITING` 不阻塞;如果最新实例失败或出现更新的未完成实例,不要验收,先报告阻塞。

## Playground 流程

1. 打开目标 URL,确认当前页面已登录且连接到目标 MCP,例如 `mcp_core_playground`。
2. 新建会话,不要复用旧 session。
3. Prompt 必须要求 Agent “实际调用 MCP 工具”,并禁止只解释、猜测、本地模拟。
4. 观察工具调用记录。没有实际工具调用就判定未验收,重新提示或记录失败。
5. 记录 Playground session id、工具名、关键入参、返回的任务 id 或 processID。
6. 对异步任务继续要求调用查询工具到终态;不能停在“任务已提交”。

## RunIaC 专项规则

- Prompt 中明确本次要使用的 action,以实际测试场景为准。
- 验收 Token 级 HITL 时,只走 TerraformAgent Playground,并执行 [TerraformAgent Playground 验收](references/terraform-agent-playground.md) 的最小矩阵。
- 涉及真实资源创建、修改或删除时,先确认用户授权、影响范围和清理方式。
- `RunIaC` 返回 `processID` 后,继续调用 `GetTask` 直到终态。
- 收集 `status`、`nextAction`、`summary`、`message` 原文或关键片段。
- 按需求定义最小测试矩阵;涉及 Terraform 表达式或聚合逻辑时,覆盖 direct / for / for_each / 异步查询等真实用户路径。

## 证据合同

报告必须包含:

| 字段 | 要求 |
|---|---|
| 环境 | 预发或线上 URL |
| 部署门禁 | app、pipeline、instance、预发部署阶段状态 |
| Playground | session id、新会话说明 |
| Token 策略 | TerraformAgent 路径、三态选择、切换不同显式策略值时是否使用 fresh token;不得记录 token、AK、UID 或 ARN |
| 初始化 API | 单独列出 `InitializeApiMcpServerConnection`、`GetCallerIdentity`;不要计入 MCP 工具列表 |
| 初始化错误 | 记录 `code`、`initializationConsumed`、`canRetry` 和据此选择的复用或刷新动作 |
| 工具调用 | 实际 MCP 工具列表;RunIaC HITL 验收至少包含 `RunIaC`、`GetTask` |
| RunIaC | action、processID、审批状态迁移、GetTask 终态 |
| 结果 | status、nextAction、summary、message 关键片段 |
| 结论 | PASS / FAIL / BLOCKED,并说明原因 |
| Aone | 已回填的工作项或待回填草稿 |

如需上传带截图的 HTML 报告到 AutomationAgent:

1. 先加载 `html-report-preview` skill,按其 Image Handling 章节处理图片。
2. 禁止在 HTML 中使用相对图片路径或 base64 `data:image`;截图必须上传为 OSS 私有对象,再用半年期签名 GET URL 让 HTML 引用。严禁 public-read bucket/object。
3. 凭证边界:HTML 报告上传到 AutomationAgent 只依赖 `JARVIS_HTML_REPORT_TOKEN`;截图上传和签名 OSS URL 才依赖对象存储能力。不要使用任意个人 AKSK 或随便找一个 bucket;只能使用 Jarvis/团队认可的私有 bucket、最小权限上传/签名凭证,或已有可安全限时访问的图片 URL。缺合规 OSS 签名能力时不要声称图片已修复,应升级人工处理。
4. 调 `bootstrap/html-report-preview.sh upload` 前,检查 HTML 里所有 `<img src>` 都是 HTTP(S) 绝对 URL。
5. 分享前,确认未签名的 OSS 直链不可公开读取,签名 URL 用 GET 返回 200,并用浏览器确认报告页里每张图已实际加载。

## Aone 追评模板

```text
本轮已通过 <pre-agent/agent> Playground 真实验收 MCP Server,非本地模拟。

部署门禁:
- app: <appId>
- pipeline: <pipelineId>
- instance: <instanceId>
- 预发部署阶段: <SUCCESS/FAIL/BLOCKED>

Playground 证据:
- URL: <url>
- session: <sessionId>
- 初始化 API: <InitializeApiMcpServerConnection/GetCallerIdentity 或不适用>
- 实际 MCP 工具: <tool names>

验收项:
| 用例 | processID | GetTask 终态 | summary/message 观察 | 结论 |
|---|---|---|---|---|
| <case> | <processID> | <status>/<nextAction> | <summary> | <PASS/FAIL> |

结论:
<一句话说明是否通过;失败时写清缺失属性、重复计数或丢失的诊断类型。>
```

## 常见错误

- 只看 CR 或 pipeline 成功就回“已验证”。
- 在最新预发实例还在跑时打开旧环境验证。
- Playground 会话里 Agent 只解释 Terraform,没有实际调用 MCP。
- RunIaC 拿到 processID 后不调用 GetTask。
- 把 MCP endpoint path id 当成 bearer token;两者不是同一个值。
- 用通用 Playground/settings 替代 TerraformAgent Playground 验证托管 token 或 Token 级 HITL。
- 在同一个 token 上切换“开启/关闭”,忽略 `UpdateBearerTokenSafetyPolicy` 的 create-only 语义。
- 遇到 `TOKEN_FETCH_RETRYABLE/false/true` 仍刷新重做,或遇到 `TOKEN_POLICY_FAILED/true/false` 仍在原 initialization 重试。
- 把初始化 API `InitializeApiMcpServerConnection`、`GetCallerIdentity` 混入实际 MCP 工具清单。
- 将 token、AK、UID、RAM userArn 或完整审批链接写入日志、截图、HTML 或 Aone。
- 只测单一路径,漏测 direct / for / for_each / 异步查询等需求相关路径。
- 把线上 URL 当默认目标。
- Aone 回填没有 session id / processID / summary / message。
- HTML 报告截图用相对路径或 base64,导致 AutomationAgent 预览页图片 404 或被 WAF 拦截。

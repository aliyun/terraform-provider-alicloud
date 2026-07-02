---
name: mcp-server-playground-acceptance
description: >-
  Use when Alibaba Cloud MCP Server, AutomationAgent, IaCService, cloudspec/OpenAPI MCP Server,
  RunIaC, API query tools, or related service changes need end-to-end testing, realistic user-flow
  simulation, or manual experience validation through pre-agent or agent Playground.
---

# MCP Server Playground 验收

真实验收 = 新 Playground 会话 + 实际 MCP 工具调用 + 可回填证据。本地单测、部署成功、口头解释都只是前置条件。

## 配合技能

- 有 Aone URL/ID 时,先用 `aone-triage` 读单、claim、回填和 release。
- 通过浏览器打开 Playground 实测;不要用本地脚本或后端日志替代 UI 链路。
- 涉及 cloudspec / OpenAPI MCP Server / RunIaC / API 查询工具时,默认参考 `aone-triage/references/delivery-cloudspec.md` 的 app 与流水线坐标。
- 涉及 AutomationAgent 时,默认参考 `aone-triage/references/delivery-aliyun-automation-agent.md` 的 app 与流水线坐标。

## 目标环境

| 环境 | URL | 使用条件 |
|---|---|---|
| 预发 | `https://pre-agent.aliyun-inc.com/playground` | 默认目标;所有预发验收先测这里 |
| 线上 | `https://agent.aliyun-inc.com/playground` | 仅用户明确要求线上验证时使用 |

线上验证不能顺手做;必须有用户明确指令。RunIaC 按本次验收的实际测试场景选择 action。

## 预发门禁

1. 先读 Aone/CR 评论,确认本轮要验证的 CR、环境和需求口径。
2. cloudspec 默认坐标:app `260634`,预发 pipeline `420`。
3. 查最新预发实例:
   ```bash
   bin/a1id -- app pipeline status --app 260634 --pipeline-id 420 --format json
   ```
4. 如果最新实例仍是 `RUNNING` / `SCHEDULING` / `PENDING` / `CREATED`,继续等:
   ```bash
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
| 工具调用 | 实际工具列表,至少包含被验收 MCP 工具 |
| RunIaC | action、processID、GetTask 终态 |
| 结果 | status、nextAction、summary、message 关键片段 |
| 结论 | PASS / FAIL / BLOCKED,并说明原因 |
| Aone | 已回填的工作项或待回填草稿 |

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
- 实际调用工具: <tool names>

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
- 只测单一路径,漏测 direct / for / for_each / 异步查询等需求相关路径。
- 把线上 URL 当默认目标。
- Aone 回填没有 session id / processID / summary / message。

# Acube createBuildTaskV2 工作流（历史只读）

本页仅用于审计旧路由留下的证据，不再提供任何可执行工作流。

## D/E/G 禁令

D/E/G 必须在源工单上下文由 Jarvis/TerraformRD 主动开发。严禁从本页或其它入口：

- 请求 Acube 创建构建任务；
- 调用 `createBuildTaskV2`；
- 创建、复用或承载 528766；
- 新建或补写 relation；
- 改派源单或历史关联单；
- claim、wrap、release、finish 或其它 bookend。

D-generated（含 E 的 pre Meta 收敛且 QA 通过后）应在源单同步 owner=临钧（429768），通过
`bridge.terraform_route_notify` 的 typed durable DM 通知，然后直接完成 Provider
dev、CI、remote ACC 与 PR。D 手写和 G 也按 active routing skill 在源单开发。

## 唯一允许的只读用途

唯一用途是历史 relation/PR 防重取证：可以读取 relation、taskId、aoneId 和既有 Provider
PR，用于审计来源及避免重复代码工作。
这些历史记录不是 carrier、完成信号、阻塞门或 observe/wait 理由，不得据此恢复旧路由。

## I/H 边界

I 的 2169561 文档质量主腿、必要的 Provider docs 528766 兜底腿，以及 H 的 528766 路径仍按
当前 aone-triage active skill 执行。本历史页不是 I/H 的执行入口，也不得据此构造任何动作。

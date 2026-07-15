# cap: github-identity 管 push 身份但不管 commit 作者 → CLA 失败

> **[归档]** 两项补丁均已落地:`github-identity.sh` 增 `commit` 子命令(作者 `api-tool-agent <cloudspec_bot@alibaba-inc.com>`,`JARVIS_GIT_AUTHOR_NAME/EMAIL` 可覆盖)+ `push` 前 tip 作者 WARN 兜底;CLA commit 作者硬门已写入 provider-resource-dev / terraform-pr-review skill。归档保留正文作历史缺口记录。

- **缺口类型**: 能力/流程缺陷（missing_capability，工具覆盖不全）
- **阻塞任务**: jarvis 开发的 terraform-provider PR（工单 83884678 → PR 9924）CI `license/cla` 失败
- **置信度**: high_conf（已实测复现 + 修复验证）

## 现象

PR 9924 全部检查通过后，唯 `license/cla` 卡在 "Contributor License Agreement is not signed yet"。
根因：developer 子代理在 provider 仓本地 `git commit`，commit 作者被记为本地默认伪身份
`jarvis <jarvis@jarvis.local>`。CLA-assistant 按 **commit 作者邮箱**核验，而该邮箱未签 CLA。
已合并的 api-tool-agent PR（9913–9919）commit 作者均为 `api-tool-agent <cloudspec_bot@alibaba-inc.com>`（CLA 已签）。

## 为什么会漏

`bootstrap/github-identity.sh` 只封装了 **push 身份**（GIT_ASKPASS + token）与 `gh` 身份校验，
**不涉及 commit 作者**。子代理用裸 `git commit`，作者取本地 `user.name/user.email`（或伪默认），
于是 CLA 必挂。这是每个 jarvis 自开发 provider PR 的通用陷阱，非本单个例。

## 修复（本 PR）

1. `bootstrap/github-identity.sh`:
   - 新增 `commit` 子命令：以 `api-tool-agent <cloudspec_bot@alibaba-inc.com>` 作者+提交者身份 `git commit`（作者身份可用 `JARVIS_GIT_AUTHOR_NAME/EMAIL` 覆盖）。
   - `push` 前对本地 ref tip commit 作者邮箱做校验，不符预期时**显式 WARN** 并给出 `commit --amend` 修法（backstop，不阻断）。
2. skill 文档补该硬门：`provider-resource-dev`（步骤 8 PR）、`terraform-pr-review`（dev 路径）、`aone-triage/references/probe-ticket-routing.md`（PR 硬门段）——提交一律走 `github-identity.sh commit`，或确保作者=api-tool-agent。

## 复发防护

- 主防线：子代理提交走 `github-identity.sh commit`（作者天然正确）。
- 兜底：`push` 时作者不符 → WARN 提示 re-author，避免静默 CLA 失败。

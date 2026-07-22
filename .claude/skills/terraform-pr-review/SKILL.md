---
name: terraform-pr-review
description: >-
  Use when reviewing a terraform-provider-alicloud (upstream aliyun) GitHub pull request — a
  github.com/aliyun/terraform-provider-alicloud/pull/N link, "评审/看下这个 PR", or asked to check a
  TF provider PR. Default is READ-ONLY review: read diff → two-layer verify (OpenAPI + source) →
  report → comment only after authorization. Also covers the dev path when the review finds a fix
  is needed: fork branch → develop → push → CR, never push master / never direct release. ad-hoc
  PR with no Aone ticket defaults to tf_provider pool 528766 (customer-sourced: tf_customer 1086837),
  confirm before creating. Not for Aone tickets — that is aone-triage.
---

# Terraform Provider PR 评审

> **默认只读评审**——只出报告 + 草拟评论,`gh pr comment` 仅在用户授权后发;不碰分支、不推代码。
> **红线:任何代码改动一律新分支 + CR/MR 评审,master 只接已评审合入。评审本身只读,不改任何文件。**

## 前置
- 先确认 `gh` 命令存在;Jarvis 代表执行 GitHub 写操作(评论/创建 PR/推分支)前必须 `bootstrap/github-identity.sh check`;写操作用 `bootstrap/github-identity.sh gh ...`,推分支用 `bootstrap/github-identity.sh push ...`,登录名必须是 `api-tool-agent`,禁止依赖本机 ambient `gh auth`、个人账号或 ambient git 凭据。
- 评审看 alicloud(=upstream aliyun)PR;只读查证可读 upstream。改动若有,head/push 落 `api-tool-agent:<branch>`(见 `config/workspaces.json` `terraform_provider`)。
- 写操作(评论/建单)先授权;低置信不发评论,将 Task `SUSPENDED` 并发布人工决策事件。

## 1. 读 PR
```
gh pr view <url>            # 标题/作者/状态/关联 issue/labels
gh pr diff <url>            # 全 diff(大可重定向到 scratchpad)
gh pr view <url> --json files -q '.files[].path'   # 改了哪些文件
```
取改了哪些资源 .go、新增/改了哪些 schema 字段。

## 2. 进 workspace(给查证上下文)
读 `config/workspaces.json` 取 `workspaces.terraform_provider`(`ops.build`/`ops.vet` 等);本地路径用 `bootstrap/workspace.sh dir terraform_provider` 解析(base 不存绝对路径)。cd 进去取上下文;仓不在则 `scripts/sync-provider.sh` 同步/clone(**有库 fetch + `reset --hard` 强对齐 upstream——主目录是只读查证镜像,改动一律走 worktree**)。

## 3. 双层查证(顺序固定,不凭记忆)
1. **OpenAPI 全集**:解析 product+action → `AlibabaCloud ListApis`/`GetApiDefinition` 核字段名/类型/枚举/action 是否存在。JMESPath 用单引号,反引号会失败:`parameters[?name=='X'].schema.properties|[0]|keys(@)`。
2. **映射**:`curl "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x"` —— 仅判 TF↔Cloudspec 映射,**不代表实现**。
3. **实现以源码为准**:`scripts/sync-provider.sh` 同步,在 path(go fork)grep 资源 .go,核 schema/`Importer`/Create 下发参数。单复数陷阱:`*_instances` 多半是数据源。同族资源对照(如 wafv3_defense_rule)判惯例 vs 新坑。
4. **文档兜底**:GitHub raw markdown。
- OpenAPI 与源码冲突 / 缺映射 / 命不中 → low_conf,不发评论,escalate。
- **错误码语义查证**(PR 涉及 retry 白名单 / `IsExpectedErrors` / 错误码补丁等) → 读 `.claude/skills/aone-triage/references/aliyun-error-code-lookup.md`(跨 skill 复用,给定 product+code 出 HTTP/中英 message/官方 retry 建议/相邻错误码)。

## 4. 看点
- 字段名/类型对齐 OpenAPI;required/optional/computed 合理;set vs list 正确。
- **ID 组装/拆解**:`SetId` 用的字段与 Read/Delete/Import 拆 `parts[]` 对得上(alicloud 高发坑)。
- **ForceNew vs Update**:不可变字段标 ForceNew,doc 与 schema 一致。
- 错误处理:NotFound→`d.SetId("")`、retry/NeedRetry、Delete 幂等;死代码/恒假 HasChange。
- import:`ImportStateVerify` 开,computed-only 字段不漏存致 import 后 diff。
- 回归用例覆盖 create→update→清空→reimport;无破坏性改动。
- 文档(r/d markdown)字段、Import 节、示例与 schema 一致。

## 5. 写报告
落 `runs/<UTCdate>-pr-<n>.md`(UTC,如 2026-06-23):结论 → 逐项(字段/import/用例)+ 证据(源码行/OpenAPI/grep) → 风险/建议。`go build`/`vet` 跑了就附,没跑注明局限。

## 6. 出口(写操作,先授权)
- 评论:草稿过目 → `bootstrap/github-identity.sh gh pr comment <url> --body "..."`(结论→逐项+证据→建议)。**仅授权后**。
- 无 Aone 工单:ad-hoc PR 默认链 tf_provider 池 **528766**;客户来源用 tf_customer **1086837**。落池前先反问。

## 开发路径(独立、需另行授权)
评审命中需改代码 → 切 worktree 分支开发,绝不在主目录/master 改;推送和 PR 创建必须走 `api-tool-agent` 身份:
```
prov="$(bootstrap/workspace.sh dir terraform_provider)"
git -C "$prov" worktree add -b <branch> <worktree_dir> origin/master
# 提交走 identity 助手,确保 commit 作者 = api-tool-agent(CLA 硬门,见下)
bootstrap/github-identity.sh commit -m "resource/...: ..."
bootstrap/github-identity.sh check
bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD <branch>
bootstrap/github-identity.sh gh pr create --repo aliyun/terraform-provider-alicloud --head api-tool-agent:<branch> ...
```
改完通过 `bootstrap/github-identity.sh push` 推到 `api-tool-agent:<branch>` → 走评审/CR,不直发正式。

### PR 修复/验证补充
- GitHub 写操作若在 provider worktree shell 里缺 `JARVIS_GITHUB_TOKEN`,不要回退个人账号;回到 Jarvis 已认证 shell 调 `bootstrap/github-identity.sh push <owner/repo> <local-ref> <remote-ref>`。
- 上游 master 前进后,CI 若报 `jitterbit/get-changed-files` / `head commit ... is not ahead of the base commit`,先 `fetch` upstream,确认 PR commit 后 rebase 到最新 `origin/master`/`alicloud/master`,保持单提交再 force update `api-tool-agent:<branch>`。**force update 后 synchronize 事件仍必挂这两检查**(action 对 synchronize 比对 event.before→after,历史被改写 before 就不是 after 祖先)——rebase+push 后必须 `gh pr close` + `gh pr reopen` 触发 `reopened` 事件(比对 base..head)才能转绿;且 master 在 rebase 与 CI 起跑之间再前移也会挂(base.sha=master 当前 tip 必须是 head 祖先),操作要一气呵成,挂了就再追平重来(先例:PR 9977 连清三轮假红)。
- 推送前先做单提交门禁: `git rev-list --count <base>..HEAD` 必须是 `1`;若 GitHub CI `Pull Request Max Commits` 报 `commitNum>1`,不要叠加修复提交,应 squash 成一个提交后再 force-with-lease 更新 `api-tool-agent:<branch>`。
  - 对**自有 fork `api-tool-agent:<branch>`（PR-head）的 force-push 是 `autonomy.md` 预授权的例行动作（`fork_push`）**：headless 下**直接执行,不 SUSPEND、不 escalate、不等工单评论放行**（授权来自策略,工单评论不作破坏性操作授权源）。前提:目标是自有 fork PR-head(非上游 `aliyun/…`、非任何 master)、内容已过 ACC 验收。force-push 上游 / master = `release_prod` 人工硬门,不在此列。
- **commit 作者硬门(CLA)**:`license/cla` 报 "Contributor License Agreement is not signed yet" 多因 **commit 作者邮箱**不是 CLA-signed 的 `cloudspec_bot@alibaba-inc.com`(裸 `git commit` 会落本地伪身份如 `jarvis@jarvis.local`)。提交走 `bootstrap/github-identity.sh commit`;已错则 `bootstrap/github-identity.sh commit --amend --no-edit` 重署名后 force-push。CLA 校验的是作者,不是 push token / PR opener。
- PR 评论要求“可用 Example”时,先在本地用 PR provider 包/override 验证 `terraform init/validate`。示例必须含 `required_providers`,跨账号资源用 aliased providers;AK/SK 只通过 `sensitive` 变量或环境变量传入,禁止写真实值。
- 跨账号 AccTest 不只看 `TF_ACC`:先隔离 ambient `ALICLOUD_ACCESS_KEY`/`ALICLOUD_SECRET_KEY`,再显式检查 `ALICLOUD_ACCESS_KEY_1/2` 解析到的账号是否符合预期。测试前置清理只用于清历史脏关系,不能替代 provider Delete;若 CLI/API 能清理关系,资源 Delete 也应实现同等删除并校验幂等。
- CI 失败诊断必须按失败 check 的 job id 拉日志:先 `gh pr checks --json name,link,state,bucket,workflow`,从失败项 URL 拿 run/job,再用 `gh run view <run_id> --job <job_id> --log`;同一个 workflow 里的其它 job 日志不能替代失败 job。
- **`TestingCoverageRate` 红时不许简单定性"存量必挂/非阻塞"**:该检查按被改资源全量 schema 三层核验(must-set 100% / ignore 数组合法 / modify 覆盖)。评审/修复都按 provider-resource-dev SKILL 步骤 7.5「TestingCoverageRate CI 门」执行——能真实补的补(含 create-only 无 import step 的属性覆盖 test 解法),明确补不了的逐属性在关联工单说明原因,禁占位值凑覆盖。本地复现:`go run scripts/testing/testing_coverage_rate_check.go -resource=alicloud_<name>`。
- Terraform Provider PR 的 Jarvis 内部研发动作要单独落内部 Aone:PR 评审/CI/AccTest/skill 沉淀默认进 `tf_provider`(528766)或对应 Jarvis/API 工具内部池;禁止同步到 tf_customer 客户池,除非用户明确要求同步客户主单关键节点。
- **Reviewer 在公开 PR 追问 RequestId/请求链路等内部诊断细节时走双通道**:公开评论只贴脱敏对照表(合成测试数据的入参 vs 返回值),注明完整证据已 shared internally;RequestId 全链路只发内部通道——Aone 工单评论 + 钉钉私信 reviewer 本人(`config/contacts.json` 的 `github` 字段可把 GitHub login 反查花名/工号)。RequestId/内部单号/内部系统名一律不入公开评论(CLAUDE.md 工作纪律 #5)。证据不必重跑测试:历史远程 ACC task 日志随时可 `acctest.py logs --task-id <n> --download-dir <dir>` 重下载,`tf-debug.log` 含每次调用的完整请求/响应/RequestId 对照(先例:PR 9993 CAS 掩码争议,Create 明文入参后紧接 Get 即返回掩码值,链路证据直接终结"是否服务端行为"的争论)。
- **共享测试环境的 List 类响应可能带出账号内其它历史数据**(他人邮箱/公司名等,即便已脱敏):公开引用前逐条剔除,只引用本轮自建 fixture 的数据;对无法实证的字段(如账号门槛挡住的场景)如实写"未验证",不跟着既有结论笼统外推——reviewer 要的是证据,范围窄但每条有据优于全而无据。

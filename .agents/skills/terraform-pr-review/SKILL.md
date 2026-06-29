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
- `gh auth status` 验登录;无 `gh` 或未登录则停下提示用户,不绕道。
- 评审看 alicloud(=upstream aliyun)PR;改动若有,落 origin=ChenHanZhang fork(见 `config/workspaces.json` `terraform_provider`)。
- 写操作(评论/建单)先授权;低置信不发评论,写 `escalation/`。

## 1. 读 PR
```
gh pr view <url>            # 标题/作者/状态/关联 issue/labels
gh pr diff <url>            # 全 diff(大可重定向到 scratchpad)
gh pr view <url> --json files -q '.files[].path'   # 改了哪些文件
```
取改了哪些资源 .go、新增/改了哪些 schema 字段。

## 2. 进 workspace(给查证上下文)
读 `config/workspaces.json` 取 `workspaces.terraform_provider`:`path`(go fork,查证/开发同一仓)+ `ops.build`/`ops.vet`。cd 到 path 取上下文;不在则 `scripts/sync-provider.sh`(从 workspaces.json 读 path,fetch-only)同步/clone。

## 3. 双层查证(顺序固定,不凭记忆)
1. **OpenAPI 全集**:解析 product+action → `AlibabaCloud ListApis`/`GetApiDefinition` 核字段名/类型/枚举/action 是否存在。JMESPath 用单引号,反引号会失败:`parameters[?name=='X'].schema.properties|[0]|keys(@)`。
2. **映射**:`curl "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x"` —— 仅判 TF↔Cloudspec 映射,**不代表实现**。
3. **实现以源码为准**:`scripts/sync-provider.sh` 同步,在 path(go fork)grep 资源 .go,核 schema/`Importer`/Create 下发参数。单复数陷阱:`*_instances` 多半是数据源。同族资源对照(如 wafv3_defense_rule)判惯例 vs 新坑。
4. **文档兜底**:GitHub raw markdown。
- OpenAPI 与源码冲突 / 缺映射 / 命不中 → low_conf,不发评论,escalate。

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
- 评论:草稿过目 → `gh pr comment <url> --body "..."`(结论→逐项+证据→建议)。**仅授权后**。
- 无 Aone 工单:ad-hoc PR 默认链 tf_provider 池 **528766**;客户来源用 tf_customer **1086837**。落池前先反问。

## 开发路径(独立、需另行授权)
评审命中需改代码 → 切 origin fork 分支开发,绝不在主目录/master 改:
```
git -C <workspaces.terraform_provider.path> checkout -b <branch> origin/master
```
改完 push origin fork → 走评审/CR,不直发正式。

# GitHub PR 评审入口

> 触发:github.com/.../pull/N 链接 或「评审/看下这个 PR」。**默认只读评审**——只出报告 + 草拟评论,`gh pr comment` 仅在用户授权后发;不碰分支、不推代码。开发是另一条、需另行授权的链路(末尾)。

> **红线:任何代码改动一律新分支 + CR/MR 评审,master 只接已评审合入。评审本身只读,不改任何文件。**

## 前置
- `gh auth status` 验登录;无 `gh` 或未登录则停下提示用户,不绕道。
- 评审看 alicloud(=upstream aliyun)PR;改动若有,落 origin=ChenHanZhang fork(见 config/workspaces.json `terraform_provider`)。

## 1. 读 PR
```
gh pr view <url>            # 标题/作者/状态/关联 issue
gh pr diff <url>            # 全 diff
```
取改了哪些资源 .go、新增/改了哪些 schema 字段。

## 2. 进 workspace(给查证上下文)
```
bootstrap/workspace.sh resolve terraform_provider   # → 仓库路径
bootstrap/workspace.sh ops terraform-provider-alicloud build|vet   # go build/vet 命令
```
cd 到该 path 取上下文;`ensure terraform-provider-alicloud` 不在则按提示 clone。

## 3. 双层查证(复用 SKILL.md「查证」核,顺序固定不凭记忆)
1. OpenAPI 全集:`AlibabaCloud ListApis`/`GetApiDefinition` 核 product+action 与字段。
2. 映射:`getTerraformResourceSpec?terraformResourceType=alicloud_x` 判 TF↔Cloudspec 映射(不代表实现)。
3. 实现以源码为准:在 workspace grep 资源 .go,核 schema/`Importer`/Create 下发参数;单复数陷阱。
4. 文档兜底:GitHub raw markdown。
PR 看点:字段名/类型对齐 OpenAPI、import 完整、回归用例、无破坏性改动。

## 4. 写报告
报告落 `runs/<UTCdate>-pr-<n>.md`(UTC,如 2026-06-23):结论 → 逐项(字段/import/用例)+ 证据(源码行/OpenAPI/grep) → 风险/建议。`go build`/`vet` 结果附上。

## 5. 出口(写操作,先授权)
- 评论:草稿过目 → `gh pr comment <url> --body "..."`(结论→逐项+证据→建议)。**仅用户授权后**。
- low_conf(OpenAPI 与源码冲突 / 缺映射 / 命不中)→ 不发评论,写 `escalation/`,通知用户决策。
- 无 Aone 工单:ad-hoc PR 默认链 tf_provider 池 **528766**;客户来源用 tf_customer **1086837**。落池前先反问确认,不默认建。

## 开发路径(独立、需另行授权)
评审命中需改代码 → 切 origin fork 分支开发,绝不在主目录/master 改:
```
git -C $(bootstrap/workspace.sh resolve terraform_provider) checkout -b <branch> origin/master
```
改完 push origin fork → 走评审,不直发。

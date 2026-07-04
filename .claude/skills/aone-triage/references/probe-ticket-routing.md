# probe 探测工单诊断与路由(aone-triage reference)

> 触发条件:aone-triage 读单后,发现工单在 **tf_provider 池(528766)** 且带标签 `jarvis-probe` 或标题前缀 `[probe]`。
> 这类单是 jarvis 自己的 tf-customer-probe 探测产出(见 `escalation/cap-tf-customer-probe.md` / `escalation/cap-probe-fix-flywheel.md`),
> 本 reference 定义修复飞轮第③段(修复)与第④段(验证)的路由与状态机。

## 一、识别

- **标签 `jarvis-probe`** 或**标题前缀 `[probe]`**(528766 池)即 probe 工单。
- 正文含两节:
  - **「溯源」节**:tier-0 资源(`alicloud_xxx`)或 tier-1 场景(`terraform_playground/<product>/<id>/`)+ verdict 路径(`runs/probe/<日期>-*.json`)。
  - **「修复建议」节**:探测判定给出的处置方向。
- 路由判定**读工单正文**(尤其「修复建议」+「溯源」),不臆断。

## 二、四分类路由

| 类型 | 判定特征(读正文) | 处理路径 |
|------|------------------|----------|
| **provider 代码修** | 修复建议指向 schema / ValidateFunc / CRUD 逻辑 / 重试 / import 等**代码** | `provider-resource-dev` skill:fork 分支 → 改 + 单测 → `invoke-terraform-acc-test-remote` 验收 → GitHub PR。**`github-identity` 硬门**:提交走 `bootstrap/github-identity.sh commit`(commit 作者必须 `api-tool-agent <cloudspec_bot@alibaba-inc.com>`,否则 `license/cla` 挂);PR/推分支前 `bootstrap/github-identity.sh check`,token 账号必须 `api-tool-agent`,PR head `api-tool-agent:<branch>` |
| **TF 文档修** | 仅 `website/docs` 变更(doc gap / 字段说明 / 废弃标注) | 同 PR 路径 **docs-only**;**可多单打包一个 PR**,PR 正文里逐单列 Aone 链接 |
| **上游协作** | 问题在 **OpenAPI 文档 / 云产品侧**,不在 provider | **不改 provider**;按 `config/pools.json` upstream(`submit_only`)渠道或工单评论**明确转交**,标注「上游依赖」后 release(等待时不空占) |
| **需实验定性** | 处置建议含「tier-1 变体实测」/ 需实证才能定类 | 先 `bootstrap/probe.sh run <变体场景>`(按工单建变体场景),拿实证后**归入前三类**;实证(verdict 摘录)写回工单评论 |

## 三、复验与关单状态机（收尾必做）

修复 PR **合并后**,按溯源节复验:

1. **tier-0 溯源** → `bootstrap/probe.sh tier0 <resource>` 复扫,**该 finding 项应消失**。
2. **tier-1 溯源** → `bootstrap/probe.sh run <scenario>` 复跑,**应绿**(无对应 finding)。
3. master 已修但**未发布** → 工单评论标「**已修复待发布**」+ `bootstrap/claim.sh release <id> 528766`。
4. **发布后**复验绿 → `bootstrap/wrap.sh done` + `bootstrap/claim.sh finish <id> 528766`(打 `jarvis-done`)。

**复验证据(verdict 路径 / 输出摘录)必须贴回工单**——溯源可追、闭环可查。

## 四、回灌（关单前直落 playground + 工单报备）

- 场景语料库**外置在 jarvis 仓外**(`terraform_playground/`,按云产品维度两级归档),回灌**无需 worktree/MR**:
  关单前 jarvis **直接落** regression 场景到 `terraform_playground/<product>/regression-<aone-id>/`
  (`main.tf` 可从修复 PR 的验证配置改造,`source_docs` 换成 Aone 工单链接,`checks.md` 记「修复前症状/修复后期望」)。
- **落后必在对应工单评论报备场景路径**(`terraform_playground/<product>/regression-<aone-id>/`)供仓库主人查验。
- 原「`escalation/scenario-drafts` + 周批 MR 入 `probes/`」流程**已废弃**(飞轮第⑥段)。

## 五、纪律重申

- **凡动 probe 单同样 bookend**:`claim` → `wrap`(sync/done) → `release`/`finish`,与普通工单一致。
- **PR 禁 AI 署名水印**(Co-Authored-By / 🤖 Generated 等),发出前剥掉。
- **upstream PR merge 与 release_prod 是人工门**:等待时 release 工单不空占,jarvis 可自动 remind / 催 maintainer。
- probe 单的**发现方与修复方解耦**:探测会话不 claim 工单,修复由本 triage 流程认领——避免既当运动员又当裁判员。

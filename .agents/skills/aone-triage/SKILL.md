---
name: aone-triage
description: >-
  Master triage skill for ANY Aone work item. Trigger on any project.aone.alibaba-inc.com URL,
  工单ID, 「看下/处理/回复这个工作项」, or a standalone "does X 支持 Y" question (bare
  next.api.aliyun.com / product / resource / attribute question — no ticket needed). Universal
  core flow: read (aone-get) → classify (工单类型 + 所属池 + 领域) → 查证 → draft reply →
  user auth → write (a1) + bookend (claim/wrap/release/finish)。Domain-specific routing 通过
  references 加载 on-demand:
  · references/tf-customer-request-routing.md  — Terraform tf_customer 池(1086837)客户单:接入
    alicloud_xxx 新资源 / 缺属性/值/行为 / 进度催办 / 类比 PolarDB 之类 / bug
  · references/delivery-aliyun-automation-agent.md — 自家应用 Agent门户/AgentRuntime/
    aliyun-automation-agent/PlayGround 交付(app 283346)
  · references/delivery-cloudspec.md — 自家应用 cloudspec / OpenAPI MCP Server 交付(app 260634)
  · references/aliyun-error-code-lookup.md — 阿里云错误码官方定义查证(跨 skill 复用,给定 product+code 出 HTTP/message/retry 建议)
  · references/templates.md — 回复/需求骨架、机读 JSON
  · references/probe-ticket-routing.md — jarvis-probe 探测工单(标签 jarvis-probe/标题 [probe],528766 池):四分类路由/复验关单/场景回灌
  NOT for: terraform-provider-alicloud GitHub PR 评审(用 terraform-pr-review)/ 资源从零开发
  (用 provider-resource-dev)/ 特定客户单接入进度催办不属 tf_customer(视场景自定)。
---

# Aone 工单主诊断与处理

> **红线(交付改动)**:任何代码改动新分支 + CR/MR 评审,master 只接已评审合入。禁 push master,禁零 diff 空 CR 直发正式。

## 前置 · a1 CLI

全流程走 `a1`(封装 `bootstrap/a1id -- <args>`,默认 jarvis 身份;个人身份 `chenyi/guozai/linjun` 仅当仓库主人当面授权本轮才可 `a1id as <id> -- ...`)。
先 `bin/a1id -- auth whoami` 验登录;`command not found` 或认证错误 → **征得用户同意后**装:
```
curl -fsSL https://git.cn-hangzhou.oss-cdn.aliyun-inc.com/aone-cli/install.sh | sh
a1 skill install a1@0.28.0
```
未同意则停下提示手动装,别绕 MCP。

## 入口分类

| 输入形态 | 处理路径 |
|---|---|
| 有 workitemId(URL/ID + 「看下/处理/回复」) | 走"工单全流程"(下方) |
| 只有 API/产品/资源(next.api 链接 / "alicloud 支持 X 吗" / "把 X 接入 Terraform") | 直接走"查证"输出"已支持 / 缺哪些属性",**无需建单** |
| GitHub PR 链接 | 切换到 **terraform-pr-review** skill(不走本 skill) |
| 两者都有 | 工单流程内嵌查证 |

## 工单全流程(通用骨架)

### 1. 读单 + 归类

```bash
bash bootstrap/aone-get.sh <id>                    # 3h 缓存
bin/a1id -- project workitem comment list <id>
bin/a1id -- project workitem activity <id>         # 可选,看流转
```

从返回 JSON 抽:`workitemType` / `status` / `assignedTo` / `priority` / `space`(= 所属池)/ `涉及云产品(140097)` / `工单ID(104264)` / `description` **全文**(尤其末段,常藏真实诉求)/ `creator` / `计划截止日期(80)`。

**归类 = 决定加载哪本 reference**:

| space(池) | 领域 | Reference |
|---|---|---|
| **1086837** tf_customer | Terraform 客户需求(接入/属性/催办等) | `references/tf-customer-request-routing.md` |
| **528766** tf_provider | Terraform Provider 内部研发(通常由客户主单派生的关联单) | 无独立 reference;跟 tf_customer 主单同域 |
| **2124589** mcp_server | 自家应用交付(Agent门户/AgentRuntime/aliyun-automation-agent/PlayGround) | `references/delivery-aliyun-automation-agent.md` |
| cwd 在 cloudspec repo 或诉求涉 cloudspec / OpenAPI MCP Server | 自家应用交付(cloudspec) | `references/delivery-cloudspec.md` |
| **2165097** upstream.cloudspec_gap | 上游 Cloudspec 需求池(谜拟) | 提单 only,submit_only |
| 其它 | 无 domain reference → 走本文件通用流程 | — |

**判断规则**:先看 `space` 命中 池,再看 `涉及云产品` / 标题 / cwd 辅助定位。有 domain reference 就**加载并跟随**它的决策树;无 reference 走本文件下方通用查证。

### 2. 按类型分诊(通用)

| workitemType | 通用动作 |
|---|---|
| **需求/咨询**(req / 需求问题) | 查证 → 回复(授权) → 若真缺口:tf_customer 走 [[tf-customer-request-routing]] 决策;其它域 → `upstream.cloudspec_gap`(2165097)双向关联 + 源单改「待上游排期」 |
| **缺陷**(bug) | 复现要点 + 源码定位 → 回复/指派;确认 spec 缺口才转需求 |
| **任务**(task) | 直接执行或拆解,产出+回执,无需查证/转需求 |
| **自己交付**(改我们自家应用) | 走 `references/delivery-<app>.md`:建需求→建 CR→worktree 开发→预发→**等用户验证反馈**→正式→关单清 worktree |

### 3. 查证(领域无关的通用套路)

**分层顺序,不凭记忆**:

1. **OpenAPI 层**:`aliyun <product> <Action> --help` 拿官方 meta(next.api 网页是 SPA,curl 拿不到 JSON);或 `AlibabaCloud ListApis` / `GetApiDefinition` 若 MCP 可用。JMESPath 用单引号,反引号会失败:`parameters[?name=='X'].schema.properties|[0]|keys(@)`。
2. **Terraform 映射层**(仅当涉 Terraform 资源):`curl "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x"` —— 仅判"TF 资源 ↔ Cloudspec 是否建映射",不代表实现。
3. **provider 源码层**(仅当涉 Terraform 资源):先 `bash .Codex/skills/aone-triage/scripts/sync-provider.sh` 同步,再在 `$(bash bootstrap/workspace.sh dir terraform_provider)/alicloud/` grep 资源 .go,核对 schema / `Importer` / Create 实际下发参数。单复数陷阱:`*_instances` 多半是数据源。
4. **文档兜底**:GitHub raw markdown、`aliyun help <product>` 命令列表。

**Terraform-specific 领域**详细 branch(专属维护名单 / 类比 API 原生 vs Provider 适配 / 镇元覆盖度 / 生成器 vs 手写)全在 `references/tf-customer-request-routing.md`。tf_customer 域必读,其它域按需借鉴。

### 4. 回复草稿(结构固定,先给用户过目)

```
结论(一句话对齐真实诉求)
├─ 逐问 + 证据(引 provider 源码行 / OpenAPI meta 引用 / RequestId)
└─ 建议行动(转谁 / 谁 @ / 状态怎么改 / 用户侧要做什么)
```

**@ 语法** = `@花名(工号)`。团队常用工号见 memory `team-roster-tf-alicloud`(涵盖专属维护名单 11 人 + 通用路由 4 人)。
**不带 AI 署名**(AGENTS.md 工作纪律 #7):对外产物剥掉「🤖 Generated with Codex」等。

### 5. 写操作(全部先授权 — supervised 默认模式)

| 动作 | 命令 |
|---|---|
| 回复评论 | 走 wrap.sh done(见 bookend;多行用 `--summary-stdin`/`--summary-file`),别单独 `a1 comment create`(会与 wrap 里的重复,a1 无 delete) |
| 转需求(Cloudspec 缺口) | `bin/a1id -- project workitem create --project 2165097 --category req --assignee 479782 -m "<7 字段模板+机读 JSON,见 references/templates.md>"` |
| 建关联单(自家团队接手) | tf_customer 域走 `references/tf-customer-request-routing.md` 分工表 |
| 双向关联 | `bin/a1id -- project workitem relation add <A> relate:<B>` **调两次**(A→B, B→A) |
| 状态更新 | `bin/a1id -- project workitem update <id> --status "<value>"` |
| 更新详情(description) | `bin/a1id -- project workitem update <id> --body-file <path>`(单行小改可 `--body "<text>"`)。**何时必须**:重审/复核推翻了 description 里的根因或方案、方案实施与描述已相左、验收证伪原描述——评论只是过程审计追加在尾部,新读者第一眼看的是详情,详情停在已否决结论=持续误导接手者。重写时开头加一行 `> ⚠️ 本 description 于 <date> 重写:<被否决的旧结论一句话>,演进见评论区`,保住审计链。**边界**:仅限我方创建/维护的工单(tf_provider 关联单/研发单/probe 单);客户主单 description 是客户原声,禁改 |
| 字段必填缺失 | `bin/a1id -- project workitem field options <field> --project <id>` 查枚举补 `--cfs` |
| GitHub PR/评论/推分支(Jarvis 身份) | 必须先 `bootstrap/github-identity.sh check`;`gh` 走 `bootstrap/github-identity.sh gh ...`;推分支 `bootstrap/github-identity.sh push`;账号必须 `api-tool-agent`;PR head 必须 `api-tool-agent:<branch>` |

## Bookend(动工必走)

任何"要写工单"的场景都必须走完整 bookend(AGENTS.md 工作纪律 #5)。纯只读查证可免 claim。

```bash
# 1. claim(打 jarvis-claimed 标签,冻结 prefix 到 .my-day/claim-prefix-<id>.txt)
bash bootstrap/claim.sh claim <id> <pool-project>

# 2. wrap.sh done —— 一次发完整评论 + run_done + (可选)改状态
bash bootstrap/wrap.sh done <id> --summary-stdin <status|--no-status> <<'EOF'
<完整回复>
EOF

# 3. release / finish 二选一
bash bootstrap/claim.sh release <id> <pool-project>   # 本轮释放,等对方接手 → jarvis-idle
bash bootstrap/claim.sh finish  <id> <pool-project>   # 真闭环 → jarvis-done + status=已发布待需求排期
```

**wrap.sh 参数陷阱**(memory `wrap-done-single-comment`):
- 单行可继续用位置参数: `bash bootstrap/wrap.sh done <id> "<完整回复>" <status|--no-status>`
- 多行正文用 `--summary-stdin` heredoc 或先写文件再 `--summary-file <path>`;不要把换行写成字面量 `\n`
- 不支持 `--status` 命名参数;status 仍放在最后一个位置参数
- 用之前**先起草完整评论内容**,一次发完(先手动 `a1 comment create` 再 wrap.sh done 会重复,a1 无 comment delete)

**release vs finish**:默认 release(路由 ≠ 真闭环,需下游响应);仅当查证发现"其实已支持 + 只是客户版本旧"这类无缺口场景走 finish。

**MR/CR 未合并禁 finish**:当 MR/CR 已提交但未合并(PR state ≠ merged / CR 未合入 master)时,**禁止调 `claim.sh finish`**。正确路径:
- `wrap.sh done <id> --no-status` —— 发评论记录进展,不改 Aone 状态
- `claim.sh release <id> <project>` —— 释放为 jarvis-idle,等人工合并验收
- Aone 评论说明「MR 已提交待合并验收,链接: <PR_URL>」
- `claim.sh finish` 内置了硬闸门(退码 2),即使遗漏也会拦截

**关联单 claim 规则**:指派给过载(484483)的关联单,jarvis 直接 claim 跟进解决,bookend 同时处理客户主单与关联单(研发细节 wrap 关联单,客户主单只 wrap 关键节点,收尾两边各自 done+release);指派其他人(谜拟/新山/临钧等)或 Cloudspec(2165097)池的关联单不 claim,建单 + @ 即可,不 touch 标签。

## 自己交付(改自家应用)

自家应用交付走对应 reference 的完整链路(需求→CR→worktree→预发→**等用户验证反馈**→正式→关单→清 worktree)。app IDs / 预发正式流水线 ID / 常见坑见:

- `references/delivery-aliyun-automation-agent.md`(app 283346,预发 66/正式 67)
- `references/delivery-cloudspec.md`(app 260634,预发 420/正式 67)

**红线**:worktree 上开发,分支只走 CR/PR/MR;master 只接已评审合入;正式发布(release_prod)永远人工确认后触发。

## 无头模式挂起（headless suspend）

当 Jarvis 由 bridge/Tata 后台委派（无终端交互）且遇到必须等待人类确认/决策的点时：

1. 先在 Aone 工单上评论你的问题，使用 `@花名(工号)` 格式 @需要回答的人
2. 在**本轮最终回复的末尾**单独一行输出挂起哨兵：
   ```
   [[SUSPEND:{"aone_id":"<工单ID>","wait_for":"<花名>","reason":"<一句话说明等什么>"}]]
   ```
3. 输出哨兵后**立即停止工作**——系统会自动挂起你的会话、释放进程资源
4. 对方在 Aone 工单评论回复后，系统会用 `--resume` 唤醒你继续处理，上下文完整保留

**触发挂起的典型场景**：
- 需确认走生成器链路还是手写代码
- 需要使用非 jarvis 身份（chenyi/guozai 等），须在工单上获得本人授权
- 查证结果矛盾，需要人类判断
- 方案有多个可行路径，需要决策

**注意**：只有在无头模式下才使用 `[[SUSPEND:...]]`。交互模式（终端）直接在终端提问即可。

## 反模式

- ❌ 读单只看标题不读 description 末段"限制/差异/仍需" —— 常藏真实诉求(memory `read-description-last-paragraph`)
- ❌ Terraform 相关工单不加载 `references/tf-customer-request-routing.md` —— 会漏专属维护名单直接被路由到过载/谜拟
- ❌ 用 next.api 网页 curl 拿 API meta —— SPA 拿不到 JSON,用 `aliyun <product> <Action> --help` 或 MCP `ListApis/GetApiDefinition`
- ❌ 写操作跳过用户授权(supervised 模式) —— 每一步 comment/create/relation/status 都先给草稿等 yes
- ❌ 用 wrap.sh done 之前先手动 `a1 comment create` —— 重复评论且 a1 无 delete
- ❌ 多行 wrap 评论写成 `"第一行\n第二行"` —— 字面量 `\n` 会被拦截;用 heredoc 的 `--summary-stdin` 或 `--summary-file`
- ❌ 关联单不双向 —— `relation add` 必须调两次
- ❌ jarvis 自行 push master / merge PR / release_prod —— 永久停止项(autonomy.md `stop`)
- ❌ 对外产物带 AI 署名 —— AGENTS.md 工作纪律 #7,发出前剥掉
- ❌ 使用非默认 a1 身份(chenyi/guozai/linjun)未经仓库主人当面授权 —— 红线
- ❌ 推翻性结论只发评论、不改我方工单的 description —— 研发单详情停留在已否决的根因/方案,后续接手者被第一屏误导;重审/方案演进必须同步重写详情(写操作表「更新详情」行,`--body-file`;客户主单原声禁改)。案例:83998772 方案 A→E→R 两次演进,详情滞后在 A

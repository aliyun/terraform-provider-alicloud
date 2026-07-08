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
| **2165097** upstream.cloudspec_gap | 镇元 agent 池:上游 Cloudspec 需求 + tf_customer 分支 E agent 关联单(谜拟做人类兜底 owner) | submit_only(建单 assignee=`WORKER_1783326253279` + body 必带 `## 机读信息` JSON,agent 自动接单,不 claim;详见 `references/templates.md` 硬契约) |
| 其它 | 无 domain reference → 走本文件通用流程 | — |

**判断规则**:先看 `space` 命中 池,再看 `涉及云产品` / 标题 / cwd 辅助定位。有 domain reference 就**加载并跟随**它的决策树;无 reference 走本文件下方通用查证。

**528766 特例 —— Terraform 自动审核流程单**(平台发布流水线自动建单,标题形如 `[Terraform X发布自动审核流程] 产品 [P] 资源 [R]`,评论区是流水线各闸门的自动结果),**按标题子类型分流,勿一律当复核**:

| 标题子类型 | 正确处理 |
|---|---|
| **[Terraform 资源发布自动审核流程]** | **调 `terraform-provider-release` skill 跑完整 SOP**——需求差距分析(AMP 元数据 vs provider 代码)+ 远程 ACC 实测 + 出 PR。平台流水线的「源码生成/打包上传容器」只是构建产物,**不等于代码已进 provider 仓、更不等于 ACC 验证过**;jarvis 仍须按 SOP 补 ACC + PR(PR merge 是人工门)。**禁止只复核告警就 release**——那是漏跑发布流程,不算处理达标。 |
| **[Terraform 文档发布自动审核流程]** | **复核确认闸门**(不跑资源 SOP):用镇元 `GetResourceType` 核验文档告警落在 provider 公开 schema 还是镇元元数据侧——落 provider 侧 → `provider-fix-documentation` 补文档+PR;落镇元侧/无缺口 → 复核结论 + 路由发布流程或上游 owner。无资源开发,不跑 ACC。 |

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
3. **provider 源码层**(仅当涉 Terraform 资源):先 `bash .claude/skills/aone-triage/scripts/sync-provider.sh` 同步,再在 `$(bash bootstrap/workspace.sh dir terraform_provider)/alicloud/` grep 资源 .go,核对 schema / `Importer` / Create 实际下发参数。单复数陷阱:`*_instances` 多半是数据源。
4. **文档兜底**:GitHub raw markdown、`aliyun help <product>` 命令列表。

**「文档 vs 代码一致性」分诊起点**(涉 endpoint/schema 键名/枚举值/属性行为类问题):**先对比"文档承诺"vs"代码实际行为"两侧**,不要只看代码或只看文档,更不要只看客户报错。**两侧不一致时,修哪一侧不是一定的**——代码追上文档 / 文档追上代码 / 两侧都改(通常最全面)/ 引导客户 workaround 都是可行方案,取决于兼容性、云产品长期路线图、客户改动成本、历史 tf 广泛度等,分诊者的**职责是把可行方案(常见 3-4 种)完整提交给决策者**(见下方"转单/建关联单 body 内容原则"),不预设"哪一侧才是权威"的结论。

**Terraform-specific 领域**详细 branch(专属维护名单 / 类比 API 原生 vs Provider 适配 / 镇元覆盖度 / 生成器 vs 手写)全在 `references/tf-customer-request-routing.md`。tf_customer 域必读,其它域按需借鉴。

### 4. 回复草稿(结构固定,先给用户过目)

```
结论(一句话对齐真实诉求)
├─ 逐问 + 证据(引 provider 源码行 / OpenAPI meta 引用 / RequestId)
└─ 建议行动(转谁 / 谁 @ / 状态怎么改 / 用户侧要做什么)
```

**@ 语法** = `@花名(工号)`。团队常用工号见 memory `team-roster-tf-alicloud`(涵盖专属维护名单 11 人 + 通用路由 4 人)。
**不带 AI 署名**(CLAUDE.md 工作纪律 #7):对外产物剥掉「🤖 Generated with Claude Code」等。

### 5. 写操作(全部先授权 — supervised 默认模式)

| 动作 | 命令 |
|---|---|
| 回复评论 | 走 wrap.sh done(见 bookend;多行用 `--summary-stdin`/`--summary-file`),别单独 `a1 comment create`(会与 wrap 里的重复,a1 无 delete) |
| 转需求(Cloudspec 缺口) | `bin/a1id -- project workitem create --project 2165097 --category req --assignee WORKER_1783326253279 --body-file <path>`;body **必须严格按** `references/templates.md` 的「Cloudspec 关联单 · 镇元 agent 接单硬契约」骨架(`## 背景` / `## 需求` / `## 机读信息` + ```json 代码块 + 7 字段全);缺 marker/字段/JSON 语法错 = agent 无法接单 = 单沉底(不指派谜拟 479782;她做人类兜底 owner 挂在源客户主单上) |
| 建关联单(自家团队接手) | tf_customer 域走 `references/tf-customer-request-routing.md` 分工表 |
| 双向关联 | `bin/a1id -- project workitem relation add <A> relate:<B>` **调两次**(A→B, B→A) |
| 状态更新 | `bin/a1id -- project workitem update <id> --status "<value>"` |
| 更新详情(description) | `bin/a1id -- project workitem update <id> --body-file <path>`(单行小改可 `--body "<text>"`)。**何时必须**:重审/复核推翻了 description 里的根因或方案、方案实施与描述已相左、验收证伪原描述——评论只是过程审计追加在尾部,新读者第一眼看的是详情,详情停在已否决结论=持续误导接手者。重写时开头加一行 `> ⚠️ 本 description 于 <date> 重写:<被否决的旧结论一句话>,演进见评论区`,保住审计链。**边界**:仅限我方创建/维护的工单(tf_provider 关联单/研发单/probe 单);客户主单 description 是客户原声,禁改 |
| 字段必填缺失 | `bin/a1id -- project workitem field options <field> --project <id>` 查枚举补 `--cfs` |
| GitHub PR/评论/推分支(Jarvis 身份) | 必须先 `bootstrap/github-identity.sh check`;`gh` 走 `bootstrap/github-identity.sh gh ...`;推分支 `bootstrap/github-identity.sh push`;账号必须 `api-tool-agent`;PR head 必须 `api-tool-agent:<branch>` |
| **钉钉私信**(所有实质动作补充通知) | `bash bootstrap/notify-dingtalk.sh <staffId> "<title>" "<body>"`;jarvis 做实质动作就私信相关方——转单/补建关联单/分支 F 上游缺口→承接方或提单人;模板 D/F/E→提单人+承接方;仅"观察等待<30 天"不发。**agent 承接方一律用人类 owner 顶替**:凡承接方是镇元 agent(`WORKER_1783326253279`)——`WORKER_` 前缀无 IM 通道,`notify-dingtalk.sh` 传该 id 会 400/静默——分支 E 一律改私信谜拟(479782,agent 归她维护);紧急双单再加私信新山(521957)。缺凭据/`JARVIS_NOTIFY_DINGTALK=0`/opt-out 均静默降级不阻断。详见 `references/tf-customer-request-routing.md` §"钉钉私信 · 通用调用姿势" |

### 转单/建关联单 body 内容原则

给他人写关联单 body(尤其转出的分支 D-新山 / E-镇元 agent / G-新山 / H-夏节 等**非过载承接**)时,body 是承接方唯一的**决策依据**,jarvis 给的信息完备度决定承接方能否直接拍板还是回来反复问(**分支 E 镇元 agent 走机读契约,body 完备度直接决定 agent 能否接单,详见 references/templates.md**)。规则:

- **多方案覆盖**:分析出根因后,不要只写"我倾向"的单一方案;至少列 3-4 种可行方案(**代码修复 / 只改文档 / 代码+文档双改 / 客户 workaround** 是常见 4 种),给对比表(改动量/兼容性/长期语义/优缺点),jarvis 明说自己倾向哪个 + 依赖对方哪个信息拍板(如"取决于产品长期路线图")。
- **根因定位含完整代码引用**:引用具体文件:行号(如 `alicloud/provider.go:2273-2286`)+ 完整代码片段(不是省略号或概括)+ 一句话解释每处做什么。承接方一眼能顺着行号跳读,不用自己 grep。
- **推理链完整**:根因是"A 导致 B 导致 C"时,把每一环写清,不跳步。让承接方能验证 jarvis 的定位对不对,不是被动接受结论。
- **影响面 + acc test 建议**:影响范围(仅某云产品 / 全局)、客户 tf 兼容性、需要覆盖的 acc test case 列表,承接方能据此评估回归风险。
- **body 可 update 不必只追加 comment**:分析扩展后走 `a1 update --body-file` 覆盖(见上表"更新详情"),让新读者第一眼看到完整版;evolution 走评论区 audit trail。**追加 comment** 适合 diff-式补充(如"补充方案 Z-C/Z-D"),**update body** 适合完整重写。两者可并用:先追加 comment 让承接方收到通知,再 update body 让详情整洁。

## Bookend(动工必走)

任何"要写工单"的场景都必须走完整 bookend(CLAUDE.md 工作纪律 #5)。纯只读查证可免 claim。

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

**关联单 claim 规则**:指派给过载(484483)的关联单,jarvis 直接 claim 跟进解决,bookend 同时处理客户主单与关联单(研发细节 wrap 关联单,客户主单只 wrap 关键节点,收尾两边各自 done+release);指派其他人(新山/临钧等)或 Cloudspec(2165097)池的关联单(镇元 agent 接单)不 claim,建单 + @ 即可,不 touch 标签。

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
- ❌ Terraform 相关工单不加载 `references/tf-customer-request-routing.md` —— 会漏专属维护名单直接被路由到过载/新山/镇元 agent
- ❌ 建 2165097 池 Cloudspec 关联单 body 缺 `## 机读信息` + JSON 段 —— 镇元 agent 靠机读契约驱动 spec/映射/覆盖度动作,marker/字段/JSON 缺任何一项都接不了单,单沉底(2165097 池又不在 jarvis 视检范围,会长期烂在里面);正确姿势严格按 `references/templates.md` 「Cloudspec 关联单 · 镇元 agent 接单硬契约」骨架
- ❌ 分支 E 关联单 assignee 写谜拟 479782 —— 谜拟已不解单(2026-07-08 切换到镇元 agent),关联单硬指派 `WORKER_1783326253279`;写 479782 = 落回人手不再走 agent 自动化(谜拟保留在**源客户主单** assignee 上做人类兜底 owner)
- ❌ 用 next.api 网页 curl 拿 API meta —— SPA 拿不到 JSON,用 `aliyun <product> <Action> --help` 或 MCP `ListApis/GetApiDefinition`
- ❌ 写操作跳过用户授权(supervised 模式) —— 每一步 comment/create/relation/status 都先给草稿等 yes
- ❌ 用 wrap.sh done 之前先手动 `a1 comment create` —— 重复评论且 a1 无 delete
- ❌ 多行 wrap 评论写成 `"第一行\n第二行"` —— 字面量 `\n` 会被拦截;用 heredoc 的 `--summary-stdin` 或 `--summary-file`
- ❌ 关联单不双向 —— `relation add` 必须调两次
- ❌ jarvis 自行 push master / merge PR / release_prod —— 永久停止项(autonomy.md `stop`)
- ❌ 对外产物带 AI 署名 —— CLAUDE.md 工作纪律 #7,发出前剥掉
- ❌ 使用非默认 a1 身份(chenyi/guozai/linjun)未经仓库主人当面授权 —— 红线
- ❌ 推翻性结论只发评论、不改我方工单的 description —— 研发单详情停留在已否决的根因/方案,后续接手者被第一屏误导;重审/方案演进必须同步重写详情(写操作表「更新详情」行,`--body-file`;客户主单原声禁改)。案例:83998772 方案 A→E→R 两次演进,详情滞后在 A

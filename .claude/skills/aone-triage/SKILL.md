---
name: aone-triage
description: >-
  Use when handling any Aone work item, project.aone.alibaba-inc.com URL, workitem ID,
  「看下/处理/回复这个工作项」 request, or an Alibaba Cloud support question involving
  next.api.aliyun.com, a product, resource, API, attribute, Terraform, Cloudspec, MCP Server,
  自动化服务台, IaCService, or project 1091779. Not for terraform-provider-alicloud GitHub
  PR review or zero-to-one provider resource development.
---

# Aone 工单主诊断与处理

> **红线(交付改动)**:任何代码改动新分支 + CR/MR 评审,master 只接已评审合入。禁 push master,禁零 diff 空 CR 直发正式。

## 前置 · a1 CLI

全流程走 `a1`。非 Terraform 工单默认 `bin/a1id -- <args>`（jarvis）；Terraform 工单在同一
headless run 内按 `terraform-pd/rd/qa` subagent 分工，PD/QA 只返回结构化结果、禁止外写，
开发阶段 RD 也不发工单进展，本次主处理 run 最后由 terraform-rd finalizer 聚合回复一次。
后续重要生命周期事件由 bridge 继续以 TerraformRD 身份幂等更新；其中开放 Terraform idle 单
满 8 天无实质进展时固定执行 Aone 评论区 @ + 钉钉私信，两个通道独立补偿，同一
anchor/owner epoch 各自成功一次后静默。其它无变化和重复事件静默。外写使用
`bin/a1id as terraform-rd -- ...` 或 `JARVIS_A1_IDENTITY=terraform-rd`。RD 未登录即阻断，
不回退 jarvis。旧
`pd/qa/terraform-pd/terraform-qa` 只是一版兼容别名到 RD，会告警且不读取旧 auth。个人身份
`chenyi/guozai/linjun/shanye` 仅仓库主人当面授权本轮才可临时使用。
**PD/QA 不外写**；terraform-rd finalizer 是 downstream single-writer，负责路由动作与下游
副作用；executor 只负责原主单 bookend（claim、唯一回复、outcome status/tag、release/finish）。
Terraform 线先 `bin/a1id ready terraform-rd`，其它线先 `bin/a1id -- auth whoami` 验登录；
`command not found` 或认证错误 → **征得用户同意后**装:
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
bash bootstrap/aone-image-extract.sh <id>          # 附件截图→本地,skill 自识别
```

从返回 JSON 抽:`workitemType` / `status` / `assignedTo` / `priority` / `space`(= 所属池)/ `涉及云产品(140097)` / `工单ID(104264)` / `description` **全文**(尤其末段,常藏真实诉求)/ `creator` / `计划截止日期(80)`。

**图像内容提取**:`aone-image-extract.sh` 输出 manifest,含图片本地路径(`.my-day/aone-image-ocr/<id>/`)。若 `images>0` 且 `cache=false`,**逐张 `Read` 查看**(Claude 原生 vision),提取错误消息 / API 请求-响应 / CLI 输出 / 控制台字段,整理成结构化文本作为**工单正文的补充上下文**参与后续查证与分诊,并把整理结果写入脚本给出的 `summary_file`(下轮命中缓存直接回显,跳过重复识别);`cache=true` 时脚本已回显 summary,直接用。**图缺失/识别不出不阻断**,退回纯文字流程。客户常粘控制台报错 / next.api 响应 / `aliyun` CLI 截图,漏识别 = 分诊时抽不到真实诉求,易误判 `jarvis-npe` 或路由到错的人。

**归类 = 决定加载哪本 reference**:

| space(池) | 领域 | Reference |
|---|---|---|
| **1086837** tf_customer | Terraform 客户需求(接入/属性/催办等) | `references/tf-customer-request-routing.md` |
| **528766** tf_provider | Terraform Provider 内部研发(通常由客户主单派生的关联单) | 无独立 reference;跟 tf_customer 主单同域 |
| **2124589** mcp_server | 自家应用交付(Agent门户/AgentRuntime/aliyun-automation-agent/PlayGround) | `references/delivery-aliyun-automation-agent.md` |
| **1091779** automation_platform | 自动化服务台 / IaCService 产品研发与交付(不含 Agent 链路) | `references/delivery-aliyun-automation-platform.md` |
| cwd 在 cloudspec repo 或诉求涉 cloudspec / OpenAPI MCP Server | 自家应用交付(cloudspec) | `references/delivery-cloudspec.md` |
| CloudSpec 文档文本 metadata | **分支 I** | 仅限 description、字段解释、NOTE、枚举文案且不改变字段集合/类型/约束/CRUD；创建或复用 2169561 指派念依，公开 Provider docs 同错时独立补 528766 紧急兜底腿 |
| CloudSpec 资源/schema/结构 metadata | **分支 E → 源单 Provider dev** | 原主单用 CloudSpec skills + AMP 修到 pre Meta 收敛；pre 验收后在同一源单上下文继续 Provider dev/CI/远程 ACC/PR，不走 2165097/Acube/528766 |
| 其它 | 无 domain reference → 走本文件通用流程 | — |

**判断规则**:先看 `space` 命中 池,再看 `涉及云产品` / 标题 / cwd 辅助定位。有 domain reference 就**加载并跟随**它的决策树;无 reference 走本文件下方通用查证。

`automation_platform` 的定时扫描仅消费 `assignedTo=WORKER_1782379562571`；用户直接给 Aone URL/ID 时按 `space=1091779` 路由，不以当前负责人作为拒绝处理条件。

**528766 特例 —— Terraform 自动审核流程单**(平台发布流水线自动建单,标题形如 `[Terraform X发布自动审核流程] 产品 [P] 资源 [R]`,评论区是流水线各闸门的自动结果),**按标题子类型分流,勿一律当复核**:

| 标题子类型 | 正确处理 |
|---|---|
| **[Terraform 资源发布自动审核流程]** | **调 `terraform-provider-release` skill 跑完整 SOP**——需求差距分析(AMP 元数据 vs provider 代码)+ 远程 ACC 实测 + 出 PR。平台流水线的「源码生成/打包上传容器」只是构建产物,**不等于代码已进 provider 仓、更不等于 ACC 验证过**;jarvis 仍须按 SOP 补 ACC + PR(PR merge 是人工门)。**禁止只复核告警就 release**——那是漏跑发布流程,不算处理达标。 |
| **[Terraform 文档发布自动审核流程]** | **复核确认闸门**(不跑资源 ACC):对比 OpenAPI、CloudSpec 文档源与 Provider 公开文档。text-only CloudSpec 文档 metadata 走 I→2169561；公开 Provider docs 同错时独立补 528766 紧急兜底。CloudSpec 源正确、仅 Provider 本地生成/展示偏差才走 D；涉及结构变化则走 E，pre 验收后源单继续 Provider 开发。 |

### 2. 按类型分诊(通用)

| workitemType | 通用动作 |
|---|---|
| **需求/咨询**(req / 需求问题) | 查证 → 回复(授权) → 若真缺口按 [[tf-customer-request-routing]] 区分 I 文档、E 结构、D Provider 与 F OpenAPI |
| **缺陷**(bug) | 复现要点 + 源码定位 → 回复/指派;确认 spec 缺口才转需求 |
| **任务**(task) | 直接执行或拆解,产出+回执,无需查证/转需求 |
| **自己交付**(改我们自家应用) | 走 `references/delivery-<app>.md`:建需求→建 CR→worktree 开发→预发→**等用户验证反馈**→正式→关单清 worktree |

### 2.5 Canned 缺参前置门

领域 reference 判定命中 canned 时，先起草补料/澄清内容，由当前流程的唯一写者回复后
`release/idle` **等待补料**。所需输入未齐前不得进入正式路由、开发、建单、改派或其它写操作。
允许继续不改变外部状态的**只读安全查证**，用于核实已有事实和缩小补料范围，但
**不得据此形成正式路由结论**。Terraform PD 应返回 `status: blocked`、
`requested_external_actions: []` 与
`next=terraform-rd-finalizer/finalize`；RD 不得把缺参任务转成开发。具体 canned 清单和所需材料
见 `references/tf-customer-request-routing.md`。

### 3. 查证(领域无关的通用套路)

**分层顺序,不凭记忆**:

1. **OpenAPI 层**:`aliyun <product> <Action> --help` 拿官方 meta(next.api 网页是 SPA,curl 拿不到 JSON);或 `AlibabaCloud ListApis` / `GetApiDefinition` 若 MCP 可用。JMESPath 用单引号,反引号会失败:`parameters[?name=='X'].schema.properties|[0]|keys(@)`。
2. **Terraform 映射层**(仅当涉 Terraform 资源):`curl "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x"` —— 仅判"TF 资源 ↔ Cloudspec 是否建映射",不代表实现。
3. **provider 源码层**(仅当涉 Terraform 资源):先 `bash .claude/skills/aone-triage/scripts/sync-provider.sh` 同步,再在 `$(bash bootstrap/workspace.sh dir terraform_provider)/alicloud/` grep 资源 .go,核对 schema / `Importer` / Create 实际下发参数。单复数陷阱:`*_instances` 多半是数据源。
4. **文档兜底**:GitHub raw markdown、`aliyun help <product>` 命令列表。

**「文档 vs 代码一致性」分诊起点**(涉 endpoint/schema 键名/枚举值/属性行为类问题):**先对比"文档承诺"vs"代码实际行为"两侧**,不要只看代码或只看文档。Terraform 资源文档必须增加 OpenAPI、CloudSpec 文档源、Provider 本地生成/展示三侧比较：严格 text-only 的 CloudSpec 文档 metadata 走 I；CloudSpec 源正确、Provider 本地偏差走 D；涉及字段集合/类型/约束/CRUD 则走 E，pre 验收后在源单上下文继续 Provider 开发。

**Terraform-specific 领域**详细 branch(专属维护名单 / 类比 API 原生 vs Provider 适配 / 镇元覆盖度 / 生成器 vs 手写)全在 `references/tf-customer-request-routing.md`。tf_customer 域必读,其它域按需借鉴。

5. **可视化截图取证**(查证完成后追加):调用 `.claude/skills/screenshot-evidence` skill，把关键页面组装成 HTML 可视化报告并上传 pre-agent 预览。**Terraform 工单的 OpenAPI、CloudSpec/ACube 映射、Provider 源码三层查证必须各有截图证据；确实不适用或页面不可达时，在 manifest 中逐层写明 N/A 原因，禁止静默省略。**为保持单写者约束，PD 只在 `.my-day/screenshots/<aone-id>/` 生成本地截图和 `evidence-manifest.md`，不上传 OSS、不写 Aone；最终 RD finalizer 运行 `screenshot-evidence` 的 manifest 校验器后统一上传一次，调用 `html-report-preview.sh upload` 时不得传 `--comment`，并把报告链接写入本轮唯一聚合回复。executor 托管时链接进入 `AONE_RESULT.reply_body`，run 内不得调用 `wrap.sh`；非 Terraform 流程仍可由当前处理者端到端生成、上传和回贴。**评论内 URL/图片渲染规则见下方 §4「Aone 评论渲染 quirk」**；签名图片受 img src query 剥离影响，只能在 pre-agent 在线报告里展示，评论区不直接内嵌。

### 4. 回复草稿(结构固定,先给用户过目)

```
结论(一句话对齐真实诉求)
├─ 逐问 + 证据(引 provider 源码行 / OpenAPI meta 引用 / RequestId)
└─ 建议行动(转谁 / 谁 @ / 状态怎么改 / 用户侧要做什么)
```

**@ 语法** = `@花名(工号)`。以下直接评论规则只适用于非 Terraform；Terraform 的 @内容也必须并入最终 RD 唯一聚合回复。团队常用工号见 `references/team-roster.md`(专属维护名单 + 通用路由)。**同名歧义陷阱**:`a1 comment create -m "@花名"` 的自动解析对**目录里同名的人会挑错工号**(实例:工单 84043785 负责人「刘源」是 `WB01269865`,只写 `@刘源` 被解析成同名同事 `WB01437449`,通知发错人)。**凡要 @ 的人可能同名(尤其外包 WB 工号、常见姓名),一律显式写全 `@花名(工号)`**——a1 对已带工号的形式原样保留、不再重解析;工号从工单 `assignedTo`/评论作者/roster 取,别靠裸名字赌解析。
**不带 AI 署名**(CLAUDE.md 工作纪律 #5):对外产物剥掉「🤖 Generated with Claude Code」等。

**Aone 评论渲染 quirk(所有评论 URL 都受此规则约束,不只截图链接)**:

- **可点击链接的唯一可靠格式 = markdown `[text](url)`**——评论区按 markdown 渲染(先例:84307546 评论 124870464 四格式对照实测,仅 markdown 链接可点);
- 裸 URL **不 autolink**——独占一行+前后空行、或行内紧贴文字,统统渲染为不可点的死文本;
- `<a href>` 锚标签与 `<url>` 尖括号包 URL 也不行——`<...>` 被当 HTML tag 剥掉/转义,不渲染为链接;
- 评论正文仍应主动写 markdown 链接；`wrap.sh` 的 `aone-comment-format.sh` 与 bridge 幂等事件发布器各自在自己的统一出站边界自动兜底裸 URL，直调 `a1 comment create` 不经过这两道闸门;
- 详情区(description)同为 markdown 渲染,`[text](url)` 同样适用(评论与详情口径一致)。

### 5. 写操作(全部先授权 — supervised 默认模式)

| 动作 | 命令 |
|---|---|
| 回复评论 | 非 Terraform 走 wrap.sh done；Terraform 源工单由 executor 在最终 RD 聚合后 done 一次。D/E/G/pure datasource 禁止任何 528766 承载动作；I 的 Provider docs 紧急兜底腿与 H 的合法 528766 实际 claim 时由 RD finalizer 另做一次聚合 bookend。后续重要事件只走 bridge RD-only event publisher。多行用 `--summary-stdin`/`--summary-file`，别先单独回复 |
| CloudSpec 文档文本 metadata（I） | terraform-rd finalizer 单写创建或复用 2169561 并指派念依；Provider 公开 docs 同错时按分池防重独立补 528766 紧急兜底；executor 只做原主单 bookend |
| CloudSpec 结构 metadata（E） | PD 返回 `requested_external_actions: []` 与 `next=terraform-rd/dev`；RD 修 CloudSpec 到 pre Meta 收敛，QA pre 验收后回 RD 在同一源单上下文继续 Provider dev/CI/PR，再由 QA 远程 ACC；不走 Acube/528766 |
| pure datasource | 紧急源单指派新山、非紧急源单指派过载；Jarvis/TerraformRD 在源单直接开发，严禁创建或复用 528766 |
| D/G source-only | D 手写紧急/非紧急/生成发布分别同步源单给新山/过载/临钧并主动开发，D finalizer 通过 ledger 幂等 enqueue owner DM；G 源单给新山并主动开发、不发 route DM；D/E/G 严禁 528766 承载 |
| 建关联单(自家团队接手) | tf_customer 域走 `references/tf-customer-request-routing.md` 分工表 |
| 双向关联 | `bin/a1id -- project workitem relation add <A> relate:<B>` **单次即自动双向**(重复调第二次返回 400 已存在) |
| 状态更新 | `bin/a1id -- project workitem update <id> --status "<value>"` |
| 更新详情(description) | `bin/a1id -- project workitem update <id> --body-file <path>`(单行小改可 `--body "<text>"`)。**何时必须**:重审/复核推翻了 description 里的根因或方案、方案实施与描述已相左、验收证伪原描述——评论只是过程审计追加在尾部,新读者第一眼看的是详情,详情停在已否决结论=持续误导接手者。重写时开头加一行 `> ⚠️ 本 description 于 <date> 重写:<被否决的旧结论一句话>,演进见评论区`,保住审计链。**边界**:仅限我方创建/维护的工单(tf_provider 关联单/研发单/probe 单);客户主单 description 是客户原声,禁改 |
| 字段必填缺失 | `bin/a1id -- project workitem field options <field> --project <id>` 查枚举补 `--cfs` |
| GitHub PR/评论/推分支(Jarvis 身份) | 必须先 `bootstrap/github-identity.sh check`;`gh` 走 `bootstrap/github-identity.sh gh ...`;推分支 `bootstrap/github-identity.sh push`;账号必须 `api-tool-agent`;PR head 必须 `api-tool-agent:<branch>` |
| **钉钉私信**(所有实质动作补充通知) | 非 Terraform 可按路由规则调用。Terraform 主处理 run 仅 D route 由 RD finalizer 在 owner/status 同步后调用 `python3 -m bridge.terraform_route_notify` 类型化 enqueue；模型禁止裸调 notify 脚本。其它阶段通知禁用；满 8 天重访仍走 bridge 双通道 ledger |

### 纯 datasource source-only 契约

只在诉求仅涉及 `data.alicloud_xxx` 的查询、过滤、分页、输出字段或 Read，且不含 resource
变更时命中。resource+datasource 混合诉求、G Provider 全局改造、手写 resource D 均不属于 pure datasource。
该契约优先于旧 G/urgent-D，但不得外溢：

- 紧急源单 assignee=新山（521957），非紧急源单 assignee=过载（484483）；均由
  Jarvis/TerraformRD 在源单直接开发。
- 严禁为 pure datasource create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766。
- 历史 relation 只读保留，不删、不迁、不关、不改派；不是开发、完成或 blocker 门，
  允许引用已有 PR 防重复。
- RD route phase 只幂等同步源单 assignee + per-type progress_status；
  bridge executor 独占源单 claim/唯一回复/tag/release/finish。
- CI pending/fail 或 QA fail 均回 RD 修复，不得标为 blocked；只有真实 `missing_capability`、
  `retry exhausted`、明确外部依赖或人工决策才可 blocked/SUSPENDED。
  open PR + QA pass 时源单 release，不 finish。
- D/E/G 同样严禁 528766 承载；I/H/pure datasource/A/F 保持原边界。

### D/G source-only + D route DM 契约

- D 手写紧急源单给新山（521957），非紧急给过载（484483），生成/发布腿（含 E pre 收敛后）
  给临钧（429768）；G 源单给新山（521957）。均由 Jarvis/TerraformRD 在源单上下文主动开发。
- D/E/G 严禁 create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766；
  历史 relation 只读保留且不构成开发、完成或 blocker 门。
- D finalizer 先幂等同步源单 assignee + per-type progress_status，再运行
  `python3 -m bridge.terraform_route_notify --ticket <id> --subtype
  <handwritten-urgent|handwritten-normal|generated>`，最后交 AONE_RESULT。禁止模型裸调
  `notify-dingtalk.sh`。event key 固定 `terraform-route:d:<subtype>:owner:<staffId>`；
  durable pending 不阻断开发，posted/suppressed 不重发，post_uncertain 保持同 receipt；
  ledger 无法持久化不得宣称通知完成。G 不发新增 route DM。
- 不得因 owner/status、通知或历史 relation 观察等待；build/test/CI/QA fail 回 RD 修复。
  open PR + QA pass 时源单 release 不 finish，正式发布仍为人工硬门。
- I 仍 2169561→念依，public docs 同错时独立 528766→过载；H 仍 528766→夏节；A/F 与 pure
  datasource 保持原边界。

### 转单/建关联单 body 内容原则

给他人写关联单 body(例如 I-念依 / H-夏节等承接)时,body 是承接方唯一的
**决策依据**,jarvis 给的信息完备度决定承接方能否直接拍板还是回来反复问。D/E/G 不建关联单；
I/H 按各自专用契约。其它关联单规则:

- **多方案覆盖**:分析出根因后,不要只写"我倾向"的单一方案;至少列 3-4 种可行方案(**代码修复 / 只改文档 / 代码+文档双改 / 客户 workaround** 是常见 4 种),给对比表(改动量/兼容性/长期语义/优缺点),jarvis 明说自己倾向哪个 + 依赖对方哪个信息拍板(如"取决于产品长期路线图")。
- **根因定位含完整代码引用**:引用具体文件:行号(如 `alicloud/provider.go:2273-2286`)+ 完整代码片段(不是省略号或概括)+ 一句话解释每处做什么。承接方一眼能顺着行号跳读,不用自己 grep。
- **推理链完整**:根因是"A 导致 B 导致 C"时,把每一环写清,不跳步。让承接方能验证 jarvis 的定位对不对,不是被动接受结论。
- **影响面 + acc test 建议**:影响范围(仅某云产品 / 全局)、客户 tf 兼容性、需要覆盖的 acc test case 列表,承接方能据此评估回归风险。
- **body 可 update 不必只追加 comment**:非 Terraform 分析扩展后可走 `a1 update --body-file` 覆盖(见上表"更新详情"),让新读者第一眼看到完整版;evolution 走评论区 audit trail。**追加 comment** 适合 diff-式补充(如"补充方案 Z-C/Z-D"),**update body** 适合完整重写。Terraform 不执行这类阶段性 comment/update 组合，路由与说明统一由最终 RD 审查后一次落地。

## Bookend(动工必走)

任何"要写工单"的场景都必须走完整 bookend(CLAUDE.md 工作纪律 #3)。纯只读查证可免 claim。
非 Terraform 使用第一组通用骨架。Terraform 工单先在同一 run 内完成 PD→RD→QA 结构化协作，
PD/QA 不写外部系统；最后由 RD finalizer 汇总全部证据、MR/CR 链接、路由动作与下一步。
控制面 executor 托管时，源工单禁止模型直接 claim/wrap/release/评论，只返回
`AONE_RESULT` 供 executor 收口。D/E/G/pure datasource 严禁对 528766 执行承载动作；I 的
Provider docs 紧急兜底腿与 H 仍按各自合法边界处理，RD finalizer 只对本 run 实际 claim 的
I/H 研发单执行一次关联单 bookend。D route DM 是唯一主处理通知例外：finalizer 必须先同步
源单 owner/status，再经类型化 ledger enqueue，最后交 AONE_RESULT；G 不发 route DM。后续
重要事件由 bridge 统一幂等发布，不复用阶段评论通道。

```bash
# 非 Terraform
bash bootstrap/claim.sh claim <id> <pool-project>
bash bootstrap/wrap.sh done <id> --summary-stdin <status|--no-status> <<'EOF'
<完整回复>
EOF
bash bootstrap/claim.sh release <id> <pool-project>   # 本轮释放,等对方接手 → jarvis-idle
bash bootstrap/claim.sh finish  <id> <pool-project>   # 真闭环 → jarvis-done + status = pools.json 里该池 × workitemType 的 done_status

# Terraform 源工单（executor 托管）：finalizer 只返回一条聚合结果，不直接执行下列命令
# [[AONE_RESULT:{"outcome":"done|idle|suspend","reply_body":"<完整聚合回复>",...}]]

# I Provider docs 紧急兜底腿或 H 的合法 528766：仅 RD finalizer 执行，每条命令最多一次
# D/E/G/pure datasource 严禁进入本骨架。
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh claim <related-id> 528766
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/wrap.sh done <related-id> --summary-stdin --no-status <<'EOF'
<PD + RD + QA + 路由动作 + 下一步的完整聚合回复>
EOF
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh release <related-id> 528766
# PR 未合并、blocked 或待人工确认时只 release，不得 finish。

# 非 executor 托管的独立 finalizer 才对源工单使用显式 claim/wrap/release|finish。
```

**claim 退码 3 = 工单缺必填字段，不是 lost race**：Aone 对存量老单执行任意 update 时可能返回
`【<字段>】(<fieldId>)不能为空`。先运行 `bash bootstrap/aone-fields.sh missing <id>`；若字段定义的
options 为空，脚本会继续查询 field options API，并返回合法候选。由 agent 根据工单语义明确选值后，
运行 `bash bootstrap/aone-fields.sh fill <id> <fieldId>=<value>`，再重试 claim。禁止自动选择或盲填分类。
`triage-one.sh` 收到退码 3 只输出候选并 escalate，不会 wrap/release。

**finish 时的 status 由 pools.json per-池 × per-workitemType 决定**,不是全局默认。看 `config/pools.json` 里 `pools.<pool>.done_status` 是**对象**(按 workitemType displayValue 分派),例:
- tf_customer(1086837) · 需求问题 → `已发布待需求方验收`
- tf_provider(528766) · **产品类需求** → `待发布`(**不是** `已发布`——workflow 不允许从 `已选择` 直跳 `已发布`,发版是人工门,jarvis 到 `待发布` 停手)
- tf_provider · 功能缺陷/线上问题 → `Fixed`
- tf_provider · 任务 → `已完成`
- **`.claim.done_status` 是全局兜底**(`已发布待需求排期`),只用于池未映射该 workitemType 的边缘情况;主流 workitem 都走 per-池 per-category。

**遇到 `status 'X' was rejected` 的处理路径**:
1. 先查该池该 type 的合法 status enum:`bin/a1id -- project workitem field options status --project <id> --type <workitemType>`
2. 判定被拒是 **enum 不存在** 还是 **workflow 转换不合法**(状态转移图不允许从当前态直跳目标态)
3. 更新 `config/pools.json` 里的 `done_status[<workitemType>]` 到合法且**能从 claim 后进行中态到达**的终态
4. 已 tag `jarvis-done` 但 status 卡在中间态时,手动一步:`bin/a1id -- project workitem update <id> --status "<正确终态>"`

**wrap.sh 参数陷阱**:
- 非 Terraform 单行可继续用位置参数: `bash bootstrap/wrap.sh done <id> "<完整回复>" <status|--no-status>`；Terraform 使用上方显式 RD identity 版本
- 多行正文用 `--summary-stdin` heredoc 或先写文件再 `--summary-file <path>`;不要把换行写成字面量 `\n`
- status 可用位置参数(最后一个),也可用 `--status <值>` / `--status=<值>` / `--no-status` 命名参数(任意位置);flag 与位置 status 互斥,二给其一
- 非 Terraform 用之前**先起草完整评论内容**,一次发完(先手动 `a1 comment create` 再 wrap.sh done 会重复,a1 无 comment delete)；Terraform 禁止手动 comment

**Terraform 单写者附加门**:
- PD 的路由动作和 QA 的缺陷/验收结论都只是 `requested_external_actions` 提案，由最终 RD 审查。
- PD 必须把三层截图和 `evidence-manifest.md` 路径放入 `visual_evidence_manifest`；finalizer 用 `screenshot-evidence/scripts/validate-manifest.py` 校验三层齐全后统一上传报告，链接只进入本轮唯一聚合回复（executor 托管时为 `AONE_RESULT.reply_body`）。缺层、缺 manifest 或上传失败必须明确标为 blocked/missing_capability，禁止无说明跳过。
- QA fail 在内部退回 RD 修复并重跑，不对外同步失败过程。
- normal / close / escalate 的每个主处理 headless run 都最多一条 RD 聚合回复。
- MR/CR 链接只写入最终聚合；开 MR/CR 后不立即 sync。
- 后续重要事件（revisit 新结论、满 8 天无实质进展的 Aone @ + 钉钉私信、PR merged/closed/merged+npe、CI 修复达上限、派发
  retries-exhausted/timeout/max-turns/stale-orphan）允许 RD 更新一次；每次轮询无变化、
  CI pending/单次 retry/new head、普通 reviewer comment、内部交接和同 key 重复检测静默。
- 后续事件只能走 bridge 独立 pending/posted ledger + 固定摘要 marker 的发布器；semantic
  source 经 SHA-256 短摘要后才落盘/出现在 marker。正文统一清理裸/括号 PD/RD/QA
  分诊/开发/验收头与 handoff，结构化脱敏 RequestId、资源 ID、Authorization、
  AccessKey/token/secret/password/user/用户名/钉钉密钥并限长；revisit 模型 summary 只接受
  240 字内单行纯文本，命中内部/敏感/JSON/URL/多行/超长则固定降级。create 返回 comment id
  即 posted；rc=0 无 id 且 marker 未可见则进入 `post_uncertain`，后续只查 marker、绝不再次
  create；明确写失败才保留 pending 重试。禁止各调度器自行 `comment create`。

**release vs finish**:默认 release(路由 ≠ 真闭环,需下游响应);仅当查证发现"其实已支持 + 只是客户版本旧"这类无缺口场景走 finish。

**收尾蒸馏钩子(涉及 terraform 云产品的工单必挂)**:工单涉及某个 terraform 产品(客户单/内部研发单/probe 单皆算),在 `wrap.sh done` 之后、`claim.sh release/finish` 之前,按 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md` 契约把本单学到的产品级事实蒸馏进 `<playground>/<product>/KNOWLEDGE.md`(触发点②aone-triage bookend 收尾——这是评审阻断项,客户单场合的蒸馏钩子必须挂在主流程,不能只挂 probe 侧)。收录判据:可执行 / 跨场景复用 / 非文档已明示;条目格式 `- [YYYY-MM-DD][来源: 工单URL/verdict路径/PR URL] <可执行的产品级事实>`。playground 路径解析走 `bootstrap/workspace.sh dir tf_playground` 或 env `JARVIS_TF_PLAYGROUND`。

**MR/CR 未合并禁 finish**:当 MR/CR 已提交但未合并(PR state ≠ merged / CR 未合入 master)时,**禁止调 `claim.sh finish`**。正确路径:
- `JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/wrap.sh done <id> "<final aggregate>" --no-status`(多行用 `--summary-stdin --no-status`)—— Terraform 由 RD finalizer 一次写完完整结果,不改 Aone 状态
- `JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh release <id> <project>` —— 释放为 jarvis-idle,等人工合并验收
- 最终聚合说明「MR 已提交待合并验收,链接: <PR_URL>」
- `JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh finish` 内置了硬闸门(退码 2),即使遗漏也会拦截

**关联单 claim 规则**:D/E/G/pure datasource 的历史 relation 均仅供只读防重，不 claim、
不改派、不关单。I 的公开 docs 紧急兜底腿由内部链跟进；2169561 念依单不 claim。H 的
528766→夏节保持原边界。最终 RD 审查动作并把结果合入主单唯一聚合回复。

## 自己交付(改自家应用)

自家应用交付走对应 reference 的完整链路(需求→CR→worktree→预发→**等用户验证反馈**→正式→关单→清 worktree)。app IDs / 预发正式流水线 ID / 常见坑见:

- `references/delivery-aliyun-automation-agent.md`(app 283346,预发 66/正式 67)
- `references/delivery-aliyun-automation-platform.md`(project 1091779,app 172823,预发 66/正式 67)
- `references/delivery-cloudspec.md`(app 260634,预发 420/正式 67)

**红线**:worktree 上开发,分支只走 CR/PR/MR;master 只接已评审合入;正式发布(release_prod)永远人工确认后触发。

## 无头模式挂起（headless suspend）

当 Jarvis 由 bridge/Tata 后台委派（无终端交互）且遇到必须等待人类确认/决策的点时：

1. 非 Terraform 先在 Aone 工单上评论问题，使用 `@花名(工号)` 格式 @需要回答的人。Terraform 不先发独立评论；由最终 RD 把问题和 @对象并入本 run 的唯一聚合 `done`
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

- ❌ 读单只看标题不读 description 末段"限制/差异/仍需" —— 常藏真实诉求
- ❌ Terraform 相关工单不加载 `references/tf-customer-request-routing.md` —— 会漏 I/E/D 精确分流
- ❌ 命中 canned 且材料未齐仍进入正式路由/开发；只读安全查证只能补事实，不构成路由授权
- ❌ pure datasource 复用历史 528766 当承载单，或对其 create/reassign/relation/claim/
  wrap/release/finish；source-only 只允许源单 owner/status 同步与源单直接开发
- ❌ 把 text-only 文档 metadata 送进 E，或只建 528766 而漏建 I 的 2169561 念依主腿
- ❌ 把结构 metadata 送进 I；字段集合、类型、约束或 CRUD 变化必须走 E
- ❌ E pre 未验收就开始 Provider 生成/开发；正确路径是 pre QA 通过后在源单上下文继续 Provider dev/CI/PR/远程 ACC
- ❌ `amp publish pre` 成功就 `finish` 或宣称正式发布 —— prod/online、master/main 与正式发布仍是人工硬门
- ❌ 用 next.api 网页 curl 拿 API meta —— SPA 拿不到 JSON,用 `aliyun <product> <Action> --help` 或 MCP `ListApis/GetApiDefinition`
- ❌ 写操作跳过用户授权(supervised 模式) —— 每一步 comment/create/relation/status 都先给草稿等 yes
- ❌ 非 Terraform 用 wrap.sh done 之前先手动 `a1 comment create`；Terraform 在 PD/RD/QA 阶段手动 comment —— 都会制造重复评论且 a1 无 delete
- ❌ 把“内部交接零评论”误解成“Terraform 全生命周期只准一条评论” —— 正确规则是：主处理 run
  只由 RD 聚合一次；重要生命周期事件由 RD-only 幂等发布器更新；无变化 poll、CI pending/单次
  retry/new head、普通 reviewer comment、瞬时错误恢复、内部交接与同 key 重复事件静默
- ❌ Revisit/PrWatch/dispatch failure 各自直接 comment/sync —— 会绕过 pending/posted ledger、
  摘要 marker、`post_uncertain` 防重与统一 sanitizer；必须提交稳定 semantic source 给 publisher
- ❌ 多行 wrap 评论写成 `"第一行\n第二行"` —— 字面量 `\n` 会被拦截;用 heredoc 的 `--summary-stdin` 或 `--summary-file`
- ❌ 给已关联的两单重复 `relation add` —— 单次已自动双向,第二次 400 已存在
- ❌ jarvis 自行 push master / merge PR / release_prod —— 永久停止项(autonomy.md `stop`)
- ❌ 对外产物带 AI 署名 —— CLAUDE.md 工作纪律 #5,发出前剥掉
- ❌ 使用非默认 a1 身份(chenyi/guozai/linjun/shanye)未经仓库主人当面授权 —— 红线
- ❌ 推翻性结论只发评论、不改我方工单的 description —— 研发单详情停留在已否决的根因/方案,后续接手者被第一屏误导;重审/方案演进必须同步重写详情(写操作表「更新详情」行,`--body-file`;客户主单原声禁改)。案例:83998772 方案 A→E→R 两次演进,详情滞后在 A
- ❌ Aone 评论里贴裸 URL(独行或行内)、`<url>` 尖括号、`<a href>` 锚标签 —— 都不渲染为可点击链接(评论区不 autolink 纯文本,HTML tag 被剥);**唯一可点 = markdown `[text](url)`**(§4「Aone 评论渲染 quirk」,先例:84307546 评论 124870464 四格式对照)。`wrap.sh done` 与手工 `a1 comment create` 的正文都要用 markdown 链接格式

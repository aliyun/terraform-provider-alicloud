# TF 客户需求单诊断与路由(aone-triage reference)

> 触发条件：aone-triage 读单后，发现工单在 **tf_customer 池(1086837)**，或标题/
> 涉及云产品包含 `alicloud_xxx`、Terraform、CloudSpec 资源定义或资源文档源头问题。
> 目标是先定位缺口层，再决定由 Provider 路线处理、等待上游 API，还是由 open-jarvis
> 在原主单内完成 CloudSpec pre 自闭环。

## Terraform 单写者执行覆盖层

本文是 **terraform-pd 的只读诊断与路由知识**。PD 只返回证据、路由判断和结构化提案；
不得直接 comment、wrap、notify、建单、关联、改派或改状态。主处理 run 只由最终
terraform-rd finalizer 聚合一次：

```bash
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/wrap.sh done \
  <id> --summary-stdin <status|--no-status>
```

- 查证、CloudSpec/Provider 开发、QA、MR/CR、PR、阻塞与下一步都并入最终聚合。
- 后续重要事件走 bridge RD-only publisher；revisit gate 新结论、再次阻塞、PR 终态、
  CI 修复耗尽等才允许更新。
- publisher 只接收稳定 semantic source，ledger/marker 保存短摘要；`post_uncertain`
  只查 marker，不再次 create。
- revisit 模型文本不满足短单行/无敏感信息约束时使用固定安全降级。
- 无变化、CI pending/单次 retry/new head、普通 reviewer comment、内部交接与重复
  event key 必须静默。

## 核心决策树

```text
读单，抽取真实诉求、产品、资源、字段、API、优先级、DDL 与验收目标
│
├─ Provider 全局改造（region 白名单、公共 endpoint、provider.go、SDK 等）
│  └─ 分支 G：普通 Provider 研发路由
│
├─ 产品命中专属维护名单
│  └─ 分支 A：直接指派专属维护人，不查 CloudSpec，不建共享研发单
│
├─ 上游 OpenAPI 本身缺能力
│  └─ 分支 F：@提单人，状态待上游排期，不建关联单
│
├─ 纯 datasource 问题
│  └─ 分支 D：跳过 CloudSpec；按紧急度走普通 Provider 研发路由
│
└─ 资源或资源文档问题
   │
   ├─ CloudSpec 文档源错误
   │  └─ 在 CloudSpec OK 判定前短路到分支 E：原主单自闭环
   │     （不得被 schema、properties 或 CoverageScore 全绿覆盖）
   │
   └─ 非文档问题，或已证明 CloudSpec 文档源正确
      │
      ├─ CloudSpec OK（四条件全满足）
      │  └─ 缺口在 Provider
      │     ├─ 生成器产出 → 分支 D-临钧（Acube 任务自动建研发单）
      │     └─ 手写实现问题，或 Provider 本地文档生成/展示偏差
      │        → 分支 D（普通 Provider 研发路由）
      │
      └─ CloudSpec NOT OK（任一条件不满足）
         └─ 分支 E：CloudSpec 原主单自闭环
            ├─ 资源结构：字段/枚举/约束/CRUD/映射
            ├─ 文档源头：资源描述/字段解释/枚举文案
            └─ 必要时继续 Provider 从 pre 生成/手改、PR、CI、ACC
```

以上都无法定位且职责边界模糊时才走分支 H（NPE 分诊兜底），不能把“需要查证”当 NPE。

## CloudSpec OK 判定

资源或资源文档问题必须满足全部四项：

1. pre 与 online 均能查到对应资源，且 released list 命中；
2. 当前资源 schema properties 覆盖客户真实诉求；
3. Acube V2 `CoverageDetail.CoverageScore == 1.0`。
4. **文档源正确性**：涉及文档诉求时，CloudSpec 资源/属性/操作的 description、枚举文案与
   OpenAPI 长期语义一致；非文档诉求记为 N/A。

覆盖度为 1.0 但缺客户要的字段仍是 NOT OK。手写资源也不豁免，因为 CloudSpec 是长期生成、
文档和迁移合同。CloudSpec 文档源错误必须**在 CloudSpec OK 判定前短路**到分支 E，
不得被 schema、properties 或 CoverageScore 全绿覆盖。只有 CloudSpec 文档源正确、差异仅发生在
Provider 本地文档生成/展示结果时，才允许进入分支 D。纯 datasource 不进入此判定。

## CloudSpec 原主单自闭环

资源定义、metadata 与资源文档源头问题走同一路径，不再按“资源结构/文档文本”拆到不同
Aone 承接方，也不创建 Provider 文档兜底 Aone。原主单是唯一 Aone 真源。

### PD 返回契约

```yaml
requested_external_actions: []
next=terraform-rd/dev
```

PD evidence 至少包含 OpenAPI、初始 pre Meta、Provider 源码/文档、目标资源与验收标准；
不得提案 `create_related`、`relation`、`assign` 或个人身份动作。

### RD 执行契约

1. 加载 `terraform-provider-release/references/cloudspec-pre-resource-loop.md`，运行
   `bash bootstrap/cloudspec-core.sh doctor`。
2. 调用 `cloudspec-amp-workflow` 创建/切换 **task 专属 feature 分支**，使用
   **AMP 返回的 SSH URL** clone 对应 cloudspec-model。禁止在 master/main 编辑。
3. 模型目录已有 `main.cspec` 后调用 `cloudspec-idl-guide`；按变更使用
   `cloudspec-resource-edit`，必要时 `cloudspec-operation-edit`。
4. 运行 `aliyun cspec build`；失败才调用 `cloudspec-build-fix`。对每个变更资源运行
   `aliyun cspec check --name <ResourceName>`，用 `cloudspec-norm-check-fix` 收敛增量。
5. 提交并推送 CloudSpec feature 分支，执行 `amp publish pre --dry-run`；通过后执行
   `amp publish pre`，轮询 pre Meta 直到属性、CRUD 与文档收敛。
6. 需要 Provider 时必须从收敛后的 pre 重新生成/修改，完成 diff、PR、CI 与远程 ACC。
7. CloudSpec 分支、commit、MR/CR、build/check、pre、Provider PR/CI/ACC 与 blocker
   全部交 finalizer 聚合到原主单。

AMP 登录、SSH、模型仓权限、CLI 或 pre 发布能力失败时返回 `missing_capability` / `blocked`，
把缺口写入原主单；不得改派外部承接人、切个人身份或另建 Aone 绕过。`amp publish prod`、
prod/online、master/main merge/push 与正式发布始终是人工硬门。pre 完成后只能
`release/idle`，不得 finish，也不得宣称正式发布。

## 团队分工速查

普通 Provider 分支见 [team-roster.md](./team-roster.md)。CloudSpec 原主单自闭环没有新的
assignee；保持当前主单承接关系，由内部 PD→RD→QA→RD-finalizer 链处理。

## Step 1 — 读单与前置分诊

必须抽取：

- description 全文，尤其末段“限制/差异/仍需/不支持类似 X”；
- `alicloud_<product>_<resource>`、目标 API/字段/行为；
- assignedTo、priority、计划截止日期、涉及云产品、creator、space、workitemType；
- 附件中的完整错误、Terraform HCL、API 请求/响应与期望。

缺陷类型（功能缺陷/线上问题/性能瓶颈）在普通 Provider 分支按紧急处理；CloudSpec
自闭环不因紧急而另建并行单，必要的 Provider 紧急补丁仍在同一原主单内完成。

### Canned 缺参前置门

以下任一 canned 命中时，先生成补料/澄清 `reply_fragment`，由 RD finalizer 在本轮唯一回复中
发出，然后 `release/idle` **等待补料**。在该 canned 所需信息补齐前，不得进入正式路由、
开发、建单、改派或其它写操作。可以继续不改变外部状态的**只读安全查证**，用于确认已知事实
或把补料问题问得更精确，但不得据此形成正式路由结论，也不得以“先查了一部分”为由启动开发。

1. 信息不完整：索要完整 HCL、完整错误与期望。
2. OpenAPI 错误含 RequestId：核实 API error code 和参数，必要时走上游 API 分支。
3. 询问“TF 是否支持某功能”：先确认上游是否有 API/action/parameter。
4. 控制台改动导致 plan diff：说明 Terraform 全生命周期管理和 import/apply 选择。
5. Provider source 路径问题：解释历史 source 与当前兼容关系。
6. 本地无法复现：索要 Terraform debug 日志，公开回复不泄露敏感值。
7. endpoint 超时：按网络/endpoint 排查，不误判为 schema 缺口。
8. 能力边界咨询：Terraform 只能编排上游 API 已支持的能力。

发 canned 前检查同期重复单；已有承接单时避免重复打搅客户。PD 返回 `status: blocked`、
`requested_external_actions: []` 与 `next=terraform-rd-finalizer/finalize`，RD 不得把该返回
解释为已获开发授权。补料到齐后的新 run 才重新执行完整决策树。

## Step 2 — 定位缺口

加载 [镇元查证与路由分支](../../provider-resource-dev/references/zhenyuan-verification.md)：

1. upstream PR / commit 前扫；
2. OpenAPI 全集；
3. CloudSpec pre/online 资源与 properties；
4. Acube 映射/覆盖度；
5. Provider schema、CRUD、Importer、ID 与文档。

文档问题必须比较 OpenAPI、CloudSpec 资源文档与 Provider docs 三侧，并先判定 CloudSpec
文档源正确性。CloudSpec 源错误直接进入 E；只有 CloudSpec 源正确而 Provider 本地生成物/
展示有偏差才进入 D。只改 Provider markdown 会被后续生成覆盖，不能算源头闭环。

## Step 3 — 执行动作

### 写前 Gate

执行任何外部动作前 point-read 最新评论、状态、assignee 与关联关系。若已有同题 merged/open
PR 或团队成员给出命中根因的修复证据，复用现有结果，避免重复建单。

状态非 New、已有 @ 或只有讨论不等于根因已闭环；必须核对 PR/修复是否覆盖真实诉求。

### Existing-related 状态机

PR/源码前扫只能避免“同题代码已经存在”的重复工作，不能替代 Aone relation 状态机。完成
写前 Gate 后，必须横向检查源单的 relation、目标池、关联单状态、最后实质进展、assignee
activity 与已有 Acube task，再决定是否创建或触发：

| 当前事实 | 动作 |
|---|---|
| 新工单，或前次路由错误 | 分支 A 只同步源单；D/G/H 仅在正确目标池没有关联单时创建。错误的历史关联单不迁移、不关闭，但也不能阻止补齐正确目标池 |
| 路由正确，但 **D/G/H 目标池缺失** | 补建一次 528766 关联单；源单 assignee/status 已正确时保持不动 |
| 路由正确、目标池 relation 齐全，但源单 assignee 或 per-type status 漏同步/漂移 | **relation 齐全但源单映射字段漂移**：只幂等修复 assignee/status；不 create、不触发 Acube、不重复阶段回复。映射字段一致后才进入观察等待 |
| 路由正确且目标池关联齐全，距上次实质进展不足 8 天 | **观察等待**：不评论、不改状态、不改派、不 create、不触发 Acube |
| 路由正确且目标池关联齐全，**距上次实质进展 ≥8 天** | 由 bridge 的稳定 epoch 走固定 Aone @ + 钉钉双通道催办；不启动新的 PD/RD/QA run |
| 承接方已有结论，但仍缺客户或云产品材料 | **追料/补料**：finalizer 唯一回复提出材料清单并 `release/idle`，不新建关联单 |
| 关联单或承接方已有可复核终结论 | **终局收敛**：按“非 PR 终局结果表”更新源单；不新建关联单、不改到共享兜底 assignee |

约束与优先级：

- **分支 A 不要求关联单**；assignee 已是专属维护人即表示路由齐全，绝不能按“缺关联单”补建
  528766。
- D/G/H 的正确目标池是 528766；D-临钧的 528766 单由 Acube 创建。relation 中已经有正确
  目标池，或评论/activity 已有尚未回填 relation 的同一 Acube `taskId/aoneId` 时，禁止再次
  create，尤其**禁止重复触发 `createBuildTaskV2`**，也禁止重复阶段回复。已有正确 relation
  时仍须对照 roster 与 `progress_status[workitemType]` 检查源单；字段漂移只修差异字段，
  映射字段一致时也禁止重复改派或重复回复。
- **缺关联单优先于 8 天催办**：先补齐一次正确目标池并产生新的 assignee/activity epoch，
  不能一边漏建研发单一边只催进度。
- 终结论优先于追料；只有结论本身仍依赖客户材料时才进入追料。无法确认 relation 是否属于
  当前诉求时，宁可观察或 blocked，请求人工核对，不能猜测后重复建单。
- “实质进展”只包括技术/验证结论、明确 blocker 与下一步、排期承诺，以及 PR/MR/commit/
  版本证据；claim/release、handoff、纯 @、pending/收到等 canned 不计入 8 天 anchor。

### D/G/H 创建与源单同步契约

D/G/H 需要 528766 关联单时，finalizer 必须先通过上述状态机确认目标池确实缺失，再执行一次
create + relation。创建参数保持原客户单语义：

- `--project 528766`。
- **category 跟随源单 workitemType**：需求问题/产品类需求 → `req`，功能缺陷/线上问题/
  性能瓶颈 → `bug`，任务 → `task`；不得为方便统一写成 task。
- **缺陷类型强制紧急**；非缺陷复制源单 priority。
- 研发单截止日期为 **源单 DDL - 2 天**，并且**不得早于 tomorrow**；源单
  **无 DDL 时 today + 3 天**。源单自身 priority 与 DDL 不改。
- 528766 的 Task/Req 必带 `计划开始日期`、`计划截止日期`、`实际工时=0`；Bug 除三项外
  还必须带 `Terraform需求类型`，默认值和合法枚举以 Aone field options 为准。
- create 成功后只调用一次 `relation add <源单> relate:<研发单>`；relation 自动双向，禁止
  反向再调一次。

新建、Acube 返回关联单，或 Existing-related 状态机发现历史映射字段漂移时，必须
**同步源单 assignee** 到 roster 承接人；已有 relation 时只修差异字段，不重复 create、
Acube 或阶段回复：

| 分支 | 源单 assignee |
|---|---|
| D-手写/纯 datasource，紧急 | 新山（521957） |
| D-手写/纯 datasource，非紧急 | 过载（484483） |
| D-临钧生成器 | 临钧（429768） |
| G Provider 全局改造 | 新山（521957） |
| H NPE 兜底 | 夏节（401498） |

源单进度状态不得硬编码为同一个中文值；从 `config/pools.json` 的
`.pools.tf_customer.progress_status[workitemType]` 解析并由 finalizer 与 assignee 同步写入：

- 需求问题 → `问题解决中`
- 功能缺陷/线上问题 → `Open`
- 性能瓶颈 → `Open`
- 任务 → `处理中`

未映射的 workitemType 必须 blocked 并查合法枚举，不能取 progress_status 的第一个值兜底。
PD 只在 `requested_external_actions` 提案；executor 托管时由 executor 执行结构化动作并把结果
交 RD finalizer 聚合，角色不得自行中途 comment/wrap。独立 finalizer 才使用 terraform-rd
身份执行动作与本轮唯一回复。

### 分支 A — 专属维护名单

只更新源单 assignee + 按 `config/pools.json` 当前 workitemType 对应的 progress_status，
由 finalizer 聚合 @负责人；不建 tf_provider 关联单。专属名单见 team-roster。

### 分支 D/G/H — 普通 Provider 研发路由

这些分支保留既有 tf_provider(528766) 纪律：

- D：CloudSpec 四条件全满足且 Provider 手写实现有问题、纯 datasource，或已证明
  **CloudSpec 文档源正确，Provider 本地文档生成/展示偏差**；
- G：Provider 全局改造；
- H：确实无法定位的 NPE 兜底，并打 `jarvis-npe`。

普通研发单严格使用上方 category/priority/DDL/CFS 契约；指派过载的单由内部链 claim 开发，
指派其他人的单只由 finalizer/executor 建单、关联并聚合通知。创建后同步源单 assignee 与
per-type progress_status。CloudSpec 源头问题和文档源头问题**不适用本段**，不得借普通研发单
规避原主单自闭环。

### H 分支标签合并保护

H 除创建 528766 关联单、同步 assignee=夏节和 per-type progress_status 外，还要增加
`jarvis-npe`。标签更新是覆盖式 API，必须 point-read 源单 `.fields[].tag.value` 的数字 ID
列表，解析 `jarvis-npe` 的合法 option ID 后 **merge existing tag IDs** 再写回：

- **保留 `jarvis-idle`**、`jarvis-claimed` 等 Jarvis 生命周期 tag（以操作当刻实际存在者为准）；
- **保留全部业务 tag**，包括跨项目或名称不在当前 options 白名单的 tag；
- 只添加缺失的 `jarvis-npe` ID，已存在时保持幂等；
- **禁止裸名称覆盖**，不得用 `update --tag "jarvis-idle,jarvis-npe"` 之类名称列表覆盖现有值。

无法解析现有 ID 或 `jarvis-npe` option ID 时返回 blocked，不得用裸名称重试。若标签动作安排
在 `claim.sh release` 后，必须重新 point-read，因为 release 会改变 Jarvis 生命周期 tag。

### 分支 E — CloudSpec 原主单自闭环

不执行 create/relation/assign，不发阶段评论或钉钉。按本文件“CloudSpec 原主单自闭环”
完整执行并进入 QA：

```yaml
requested_external_actions: []
next:
  role: terraform-rd
  action: dev
```

成功到 pre 且静态/远程验证通过后，finalizer 汇总并 `release/idle`；等待 prod/online、
主干合并或正式发布期间保持开放。能力失败时 finalizer 报 `missing_capability` / `blocked`
并同样 release，不得 finish。

### 分支 D-临钧 — 生成器产出

当 CloudSpec 已对齐且 Provider 资源是生成器产出时，执行入口固定为
[acube-createBuildTaskV2-workflow.md](./acube-createBuildTaskV2-workflow.md)。接口内部创建
528766 研发单、指派临钧并启动生成/PR；Jarvis 不同时手动建单。

执行边界：

1. PD 只读确认“生成器产出”并提案，不调用 Acube。
2. RD finalizer 在写前 Gate 和 Existing-related 状态机再次确认：没有正确的 528766 relation，
   也没有已有 `taskId/aoneId`，才允许调用一次 `createBuildTaskV2`。
3. 调用后使用 `queryAoneByTaskId` 轮询，最多 **60 秒**；拿到 aoneId 后只做一次 relation，
   同步源单 assignee=临钧以及 per-type progress_status，并把 taskId/aoneId 放入唯一聚合回复。
4. 60 秒内拿不到 aoneId 时返回 blocked/suspend，保留 taskId 供下轮查询；禁止回退手动 create，
   因为任务可能已成功而 Aone 结果尚未可见。
5. **executor 托管**的 run 由 executor 执行上述结构化副作用，RD finalizer 返回
   `AONE_RESULT` 聚合内容且不得自行 `wrap.sh`；独立 finalizer 才用 terraform-rd 身份执行
   relation/源单同步与唯一 `wrap.sh done`。

### 分支 F — 上游 API 缺口

只有 OpenAPI 本身缺 action/parameter/合法值且 Provider/CloudSpec 无法实现时才命中：

- 不建关联单；
- @提单人协助转云产品 API 团队；
- 状态改“待上游排期”；
- 由 finalizer 一次聚合。

### 非 PR 终局结果表

关联单、承接方评论或云产品结论已经结束当前处理路径时，finalizer 必须先复核证据，再按
`config/pools.json` 与 tf_customer 当前合法枚举收敛；不得沿用旧的“一律当已发布”口径：

| 已验证结果 | 源单状态与 bookend |
|---|---|
| 能力已经存在，只是客户版本旧、配置方式错误或查证后确认无需开发 | `方案功能已存在`；无待验收事项时可 finish |
| 云产品/API 明确拒接，或正式结论为不支持且无后续排期 | `已拒绝`；finalizer 聚合拒接证据后 finish |
| 能力已经正式发布，等待需求方验证 | `已发布待需求方验收`；`release/idle` 等客户验收，不提前 finish |
| 客户验收通过 | `验收通过`；无其它开放关联或硬门时 finish |
| 已至少一次明确追料且约定期限届满仍未回应 | **先查当前合法枚举**，并执行“客户未响应三重门”：只有 field option 合法、`.claim.done_statuses` 已登记、`.pools.tf_customer.exclude_status` 已登记三项同时满足，才可把 `客户未响应` 作为终态并 finish；否则只可选择两份配置均已登记、语义与证据匹配的合法终态，无法选择则返回 blocked |
| 功能缺陷/线上问题已有可复核修复终态 | 使用 tf_customer `done_status` 的 `Fixed` |
| 任务已完成且交付物可复核 | 使用 tf_customer `done_status` 的 `已完成` |

PR 路径单独使用当前配置：需求类 PR merged 先写
`.pools.tf_customer.pr_merged_status` 的 `已合入主线` 并 `release/idle` 等发布/验收，不得恢复
旧的“merged 即已发布待验收”；缺陷类根据合法缺陷状态和真实发布门选择 `Fixed` 或继续等待。

所有终局/待验收动作都**保持最后处理人**为源单 assignee；不得在关单阶段改到过载或其它共享
兜底人。最后处理人从给出终结论的关联单 assignee、最后实质评论作者或 PR/发布 owner 取证。
executor 托管时由 executor 执行状态动作并交 RD finalizer 生成唯一回复；独立 finalizer 才
直接 bookend。存在仍开放的正确关联单、未满足的发布硬门或验收失败时不得 finish。
**客户未响应三重门**缺一项时，不得先写 `客户未响应` 再调用 finish；否则 `claim.sh finish`
可能按 workitemType 的 done_status 覆盖状态，scan 也可能再次派发。

## Step 4 — 回复骨架

所有最终回复遵循“结论 → 查证证据 → 已执行动作 → 验证 → 未决硬门/下一步”。

### 模板 A — 上游 API 缺口

```markdown
### 结论
客户诉求 <...> 的缺口位于 <Product>::<Action>；当前 API 无对应能力，Terraform
Provider / CloudSpec 无法凭空实现。

### 证据
- OpenAPI：<字段/枚举/action>
- Provider：<透传或无替代路径的源码行>

### 下一步
@<提单人>(<工号>) 请协助转云产品 API 团队；状态已改待上游排期。
```

### 模板 B — 普通 Provider 路由

```markdown
### 结论
CloudSpec 四条件已对齐；如为文档问题，已证明源头正确，缺口仅位于 Provider
<手写/生成器/datasource/本地文档生成或展示/全局> 层。

### 证据
- CloudSpec：资源、properties、CoverageScore
- Provider：文件、函数、行号
- 研发单/Acube task：<链接>

### 下一步
<承接与验证安排>
```

### 模板 C — CloudSpec 原主单自闭环

```markdown
### 结论
<资源定义/metadata/资源文档源头> 已在当前主单内完成 CloudSpec pre 修复。

### 证据
- AMP feature 分支与 MR/CR：<链接>
- build/check：<结果>
- pre dry-run/发布与 Meta 收敛：<结果>
- Provider PR/CI/ACC：<结果或 N/A>

### 未决硬门
<prod/online、master/main merge/push、正式发布或 blocker>

### 下一步
本轮 release/idle；pre 不代表正式发布，不执行 finish。
```

### 模板 D — 8 天无实质进展

bridge 以稳定 epoch 直接走 RD-only Aone + 钉钉双通道 ledger；主处理角色不执行。
同一 anchor/owner epoch 只发一次，无变化静默。

### 模板 E — 终局/待验收收敛

```markdown
### 结论
<能力已存在 / 已拒接 / 已正式发布待验收 / 已修复 / 客户未响应>；依据为 <关联单/版本/评论>。

### 证据
- 关联单或承接方终结论：<链接、状态、作者、时间>
- PR/版本/API：<适用时填写；非 PR 结论明确写 N/A>

### 下一步
源单状态更新为 <非 PR 终局结果表或当前 merged config 的值>；assignee 保持最后处理人。
<release/idle 等验收，或在无开放硬门时 finish>。
```

### 模板 F — 追料/补料

```markdown
### 当前结论
承接方已确认 <根因/限制>，继续处理还缺以下材料。

### 待补材料与期限
- <完整 HCL / 错误 / 日志 / 云产品结论>
- 建议在 <YYYY-MM-DD> 前补齐；到期仍无回复时先查 status field options、
  `.claim.done_statuses` 与 `.pools.tf_customer.exclude_status`。只有三重门全绿才使用
  `客户未响应` 并 finish；否则选择已配置且证据匹配的合法终态，无法选择就返回 blocked。

### 下一步
assignee 保持最后处理人，本轮 release/idle；材料到齐后重新进入决策树，不新建重复关联单。
```

## Bookend 与状态

- 控制面 executor 托管时，模型只返回 `AONE_RESULT`，不调用 claim/wrap/release。
- 独立 finalizer 才显式使用 `JARVIS_A1_IDENTITY=terraform-rd` 做一次 done。
- 普通 tf_provider 研发单的状态按 `config/pools.json`。
- CloudSpec 原主单自闭环到 pre 后一律 release/idle；prod/online、主干合并、正式发布与必要
  验收未完成前不得 finish。
- PR/MR/CR 未合并时不得 finish；登记 pr-watch 等后续事件。

## 反模式

### 分诊

- ❌ 只看标题，不读 description 末段与附件。
- ❌ 只看 CoverageScore，不检查客户要的 properties。
- ❌ 命中 canned 且信息未补齐，仍进入正式路由、建单或开发。
- ❌ 把只读安全查证结果当成已补齐输入并形成正式路由结论。
- ❌ 纯 datasource 去查 CloudSpec 覆盖度。
- ❌ Provider 全局改造进入资源 schema 判定。
- ❌ 未扫 upstream PR 就按本地旧 workspace 重复建单。

### CloudSpec 自闭环

- ❌ CloudSpec NOT OK 后创建新 Aone、relation、assign 或改派外部承接人。
- ❌ 因 schema/properties/CoverageScore 全绿而忽略 CloudSpec 文档源错误并进入 D。
- ❌ 资源文档源头问题只开 Provider 文档补丁/兜底单，不改 CloudSpec。
- ❌ AMP 登录、SSH、模型仓权限失败后切个人身份或退回旧承接路径。
- ❌ 在 master/main 直接编辑、merge 或 push。
- ❌ build/check 未绿就发 pre，或 dry-run 失败后继续真发。
- ❌ pre Meta 未收敛就运行生成器，或复用发布前缓存/生成物。
- ❌ 把 `amp publish pre` 写成 prod/online 正式完成。
- ❌ pre 成功就 finish；正确动作是 release/idle，等待人工硬门。

### 普通 Provider 路由

- ❌ 上游 API 缺口仍建 Provider 研发单。
- ❌ 专属维护产品污染共享研发池。
- ❌ Acube 已自动建单后再手动创建重复单。
- ❌ CloudSpec 文档源头问题用 528766 文档兜底 Aone 代替原主单自闭环。

### CLI 与出站

- ❌ 多行正文用字面量 `\n`；应使用 `--summary-file` / `--summary-stdin`。
- ❌ relation add 调两次；Aone 单次即双向。
- ❌ 评论贴裸 URL；必须使用 markdown link。
- ❌ 公开产物带客户、Aone、实例、RequestId、内部人员或 AI 署名。
- ❌ 主处理阶段直接 comment/notify；必须由最终 RD 单写者聚合。

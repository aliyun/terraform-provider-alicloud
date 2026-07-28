# persona-collab — Terraform 单写者与内部三角色

> Terraform 对外只有 TerraformRD 一个数字人。terraform-pd / terraform-rd / terraform-qa
> 仍作为同一 headless run 内的 subagent 分工，但过程不写入 Aone；主处理 run 最后由 RD
> 汇总回复一次，后续重要生命周期事件仍由 RD 幂等更新。

## 一、角色与身份

| internal_role | 内部职责 | 是否可对外回复 |
|---|---|---|
| `terraform-pd` | 分诊、三层查证、本地截图 manifest、路由提案、需求分析 | 否 |
| `terraform-rd` | 开发、PR/CR、CI 修复；最终聚合；重要事件更新 | finalizer / event publisher |
| `terraform-qa` | 普通 Provider 远程 AccTest；分支 E 的 CloudSpec pre-only 验证；回归与缺陷证据 | 否 |

唯一公共身份：

| public_identity | BUC worker | 范围 |
|---|---|---|
| `terraform-rd` | `WORKER_1783582458263` | Terraform Aone 写入、最终回复、钉钉、MR/CR |

硬约束：

1. PD/QA 不使用任何公开身份，不做 Aone、钉钉、GitHub、MR/CR 写入。
2. RD 开发阶段不发 Aone 进展；主处理 run 只有最终 RD finalizer 可回复。后续重要事件只走
   bridge 的 RD-only 幂等发布器，不由 PD/QA 或普通 headless 阶段直接评论。
3. TerraformRD 未登录即阻断，不回退 jarvis。
4. 旧 PD/QA 身份只保留入站识别兼容，不能成为新出站作者。
5. GitHub 仍须先过 `bootstrap/github-identity.sh check`，账号为 `api-tool-agent`。

旧 TerraformPD/TerraformQA 的账号、权限或存量会话由平台管理员在代码外回收。代码侧不登记、
不读取旧凭据；旧 worker id 仅用于迁移期作者、@mention 与 tracker 识别。

## 二、同一 run 的内部协议

每个 subagent 只把结果返回给编排层，不把 Aone 评论区当内部消息总线。统一字段：

```yaml
internal_role: terraform-pd | terraform-rd | terraform-qa
status: done | pass | fail | blocked | low_conf | missing_capability
summary: 一句话结论
evidence:
  - 可复核证据
visual_evidence_manifest: <Terraform PD 生成的本地三层截图 manifest；其它角色可为空>
requested_external_actions:
  - 由最终 RD 审查和执行的动作提案
next:
  role: terraform-rd | terraform-qa | terraform-rd-finalizer
  action: dev | fix | acc_verify | cloudspec_pre_verify | pre_handoff | finalize
reply_fragment: 可直接纳入最终回复的片段
```

`requested_external_actions` 对 PD/QA 永远是提案，不表示已执行。编排层必须把前序角色的完整返回
传给后序角色，不能只传一句摘要。

## 三、新工单执行链

一个 Terraform 工单只启动一个 headless run：

1. 编排层去重、探测 TerraformRD 登录；控制面 executor 在模型进程外持有 lease 并托管 claim。
2. Task 起 PD：调用 `aone-triage` 与 `screenshot-evidence`，完成三层查证、本地截图 manifest
   和路由提案；不上传、不回贴。
3. Task 起 RD：根据 PD 分支开发或 no-op；普通 Provider 路径 worktree 隔离，必要时创建 PR/CR，
   PR CI 必须全绿；分支 E 只做到 CloudSpec pre Meta 收敛。
4. Task 起 QA：普通 Provider 路径走远程 AccTest；分支 E 走 `cloudspec_pre_verify`。
5. QA fail → 把缺陷草稿与证据内部退回 RD 修复 → 重跑 QA。
6. QA pass、blocked、low_conf 或达到循环上限 → Task 起 RD finalizer。
7. finalizer 用 `screenshot-evidence/scripts/validate-manifest.py` 校验 PD manifest，统一上传一次
   可视化报告（不传 `--comment`），汇总全部返回和报告链接，执行允许的外部动作，把完整正文
   交编排层写入 `AONE_RESULT.reply_body`；executor 随后回复一次并 release/finish。

无开发需求时，RD 可返回 no-op，QA 对支持性结论或复现证据做独立校验后进入 finalizer。不得为了
形式跳过 PD/QA，也不得让 PD/QA 代替 RD 发声。

### G / 紧急普通 D 的双 owner 契约

G Provider 全局改造，以及紧急普通 D（纯 datasource；或 CloudSpec 结构 OK + 手写 Provider）
必须由内部链自助研发：

1. PD 只读提出双 owner 路由：源客户主单 assignee 保持新山（521957），528766 研发关联单
   assignee 固定过载（484483）；不得提案“交给新山等待”。
2. TerraformRD control plane point-read 同题单和 claim。healthy existing claim 不抢占，
   当前重复 run 不并发写入；同题旧单优先复用，禁止重复 create。
3. 无健康 claim 时，同一 terraform-rd Task 先进入 route-finalizer phase，fail-closed claim
   实际研发单；claim 成功后才幂等改派过载、补缺失 relation 并切 dev。claim 失败立即停止，
   禁止 relation/update/wrap。RD 随后执行 TDD、Provider 修改与 CI，QA 验收。
   relation/assignee/status 只表示路由物化；没有 PR/CI/QA 完成信号就必须继续 RD。
4. build/test/CI 或 QA fail 只走 RD ↔ QA 修复重验，绝不转交新山。`missing_capability` 或
   retries exhausted 才进入 blocked/SUSPENDED，保持源单新山、研发单过载，不 finish。
5. 源工单仍由 executor bookend；本 run 实际 claim 的 528766 则由 RD finalizer 独立
   claim/bookend。两张工单各自最多一次聚合 bookend；PR 未合并只 release。只有已有 PR
   待人工合并、明确外部依赖或人工决策，才允许 observe/release 或 blocked。
6. D-临钧/A/F/H/非紧急 D 边界不变；I→念依、E→CloudSpec pre→D-临钧保持原流程。

CloudSpec 路由分为两个互斥分支：

- **分支 I — CloudSpec 文档文本 metadata**：resource/property/operation description、字段解释、
  NOTE 与枚举文案，且不改变字段集合、类型、约束或 CRUD。PD 只提出创建或复用项目 2169561、
  指派念依（373108）的动作；公开 Provider docs 同时错误时，再提出独立 528766 紧急兜底腿。
  两腿分池防重，一个池已有 relation 不能抑制另一个池的缺失补建。RD/QA 都 no-op 验证，
  finalizer `single-writer` 执行动作。
- **分支 E — CloudSpec 结构 metadata 原主单自闭环**：字段集合、类型、约束、CRUD、
  operationMapping 或生命周期缺口。PD 返回 `requested_external_actions: []`、
  `next=terraform-rd/dev`；RD 调用 `terraform-provider-release` 的
  `cloudspec-pre-resource-loop.md` 与 CloudSpec skills，在 AMP task 专属 feature 分支完成
  build/check/publish pre 与 pre Meta 收敛，不创建 2165097。
- E 的 QA 使用 `verification_mode: cloudspec_pre`，只验 build/check/pre Meta 收敛，不运行远程
  AccTest、不要求 Provider PR/CI/ACC；pass 后返回 `pre_handoff`。finalizer 执行
  **E → D-临钧**：已有正确 relation/taskId/aoneId 时只查询/复用，否则通过 Acube
  `createBuildTaskV2` 自动创建或复用 528766 并指派临钧（429768）。
- pre 未收敛不得触发 Acube；不得由 E 直接执行 Provider PR/CI/ACC，也不得在 E 完成后直接
  release/idle。只允许分支 E 进入此转换，不得泛化到 A/F/G/H/I、纯 datasource 或纯手写
  Provider-only bug。普通分支 D 保持 PR CI + 远程 ACC。
- 权限、AMP 登录、SSH、pre 或 Acube 能力失败返回 `missing_capability` / `blocked`，不得回退
  个人身份。prod/online、master/main merge/push 与正式 release 仍是人工硬门。

## 四、主处理 run 的单次最终回复

finalizer 回复正文至少包含：

1. 总结论；
2. PD 查证与路由结论；
3. 三层可视化查证报告链接；
4. RD 改动、PR/CR 与 CI；
5. QA pass/fail/blocked 及证据；
6. 已执行的外部动作；
7. 未决项和下一步。

MR/CR 链接只在这条最终聚合回复中同步。Terraform 例外：开 MR/CR 后不立即做中途 Aone 回填。

控制面 Task 下，**源工单禁止模型直接 claim/wrap/release/评论**，**源工单仍由 executor
bookend**。**源工单禁令不约束按既有契约由内部链承接的 528766**：G/紧急普通 D，以及
**非紧急 D 与 I 的 Provider docs 紧急兜底腿**仍按各自既有边界处理；实际 claim 的
**528766 由 RD finalizer claim/bookend**，并把同一份 RD/QA 聚合写入研发单一次。G/紧急
普通 D hard gate 只新增双 owner、先 claim 后 dev 与不可观察语义，不得对源工单重复认领或
回复。模型必须在末尾返回：

```text
[[AONE_RESULT:{"outcome":"done|idle|suspend","reply_body":"<含报告链接的唯一完整回复>",...}]]
```

executor 使用 terraform-rd 身份对源工单单次落账。RD finalizer 对本 run 按既有契约实际
claim 的内部 528766 执行 `claim → wrap.sh done → release/finish`，每条命令至多一次；PR
未合并、blocked 或待人工确认时只 release，不得 finish。非 executor 托管的独立 finalizer
才对源工单执行相同 bookend。
禁止阶段回复、中途 `wrap.sh sync`、钉钉进展通知或新公开接力标记。此约束不等于“工单全生命周期
只能一条评论”。

回复后：

- 真闭环且满足 MR/CR 合并与状态门 → `claim.sh finish`；
- PR/CR 未合并、blocked、待人工确认 → `claim.sh release`；
- 有 PR → 登记 `pr-watch.sh add`。

## 五、后续重要事件与防刷屏

后续所有 Aone 更新仍只允许 TerraformRD。bridge 使用两份互不耦合的通道台账：
`.my-day/bridge/aone-event-ledger.json` 与
`.my-day/bridge/dingtalk-event-ledger.json`。Aone 成功不抑制私信失败重试，私信成功也不
抑制 Aone 失败重试：

1. 调度器只提交稳定 semantic source；ledger/marker 只保存其 SHA-256 24 位摘要，不保存
   semantic source、session、ticket 等可读细节；
2. 评论 marker 固定为 `[[JARVIS-EVENT:v1:<24hex>]]`；
3. 发前查询近期评论，marker 已存在即本地补 posted、不重发；
4. create 响应含 comment id 即作为可持久化成功证据记 posted；若 rc=0 但无 id 且 marker
   尚未可见，记 `post_uncertain`，后续只查 marker、绝不再次 create；
5. 明确写失败保留 pending/failed 重试；PrWatch 条目删除不丢事件；
6. 所有正文统一过出站 sanitizer：删除裸/方括号/全角括号 PD/RD/QA
   分诊/开发/验收标记与 handoff/sentinel，结构化脱敏 RequestId、实例/资源 ID、
   Authorization、AccessKey/token/secret/password/user/用户名及钉钉密钥并限制长度；
7. revisit 的模型 summary 只接受最长 240 字的单行纯文本；命中内部协议、敏感键、URL、
   RequestId/资源 ID、JSON/多行/超长时统一降级为「状态发生变化，详情见内部记录。」；
8. dispatch 终态只发布固定 RD 摘要（kind/subtype/尝试次数/release 结果/下一步），原始 tail
   将 Task 置为 `SUSPENDED` 并发布人工决策事件；
9. 钉钉通过 `notify-dingtalk.sh --result-json --out-track-id` 返回
   sent/skipped/failed 与稳定 receipt；opt-out 记 suppressed，凭据/网络/API 失败保留 pending
   退避。`WORKER_*` 承接人按 `config/contacts.json.agent_fallbacks` 改发真人；
10. Terraform 重访扫描所有开放 idle 单：实质进展或 assignee 变化形成新 epoch；自 anchor
    起满 8×24h（Asia/Shanghai）后固定发一次「Aone 评论区 @ + 钉钉私信」。催办模板、
    纯 @、canned、claim/release、PD/RD/QA/Jarvis/系统评论均不重置 anchor。

semantic source 使用业务事实，不使用 LLM 自由措辞；revisit 的 `semantic_id` 只允许最长 96
字符的小写 `[a-z0-9._:-]` slug。对外与落盘统一取 source 的短摘要：

- revisit：`revisit:<gate>:<transition>:<semantic-id>`；
- stale reminder：`revisit:stale:<ticket>:<anchor-kind>:<anchor-id>:<anchor-time>:<owner>`；
- PR：`pr:<url>:merged:<mergedAt>`、`pr:<url>:closed`、
  `pr:<url>:merged-npe:<mergedAt>`、`pr:<url>:ci-exhausted:<head>:<max>`；
- terminal dispatch：`dispatch:<kind>:<session-id>:<subtype>`。

需要更新一次：

- revisit gate 首次解锁形成新结论、再次阻塞、blocker 语义变化；
- PR merged、closed-unmerged、merged + `jarvis-npe`；
- CI 自动修复达到上限；
- headless retries exhausted、timeout、max-turns、stale orphan。
- 开放 Terraform idle 单满 8 天无实质进展（Aone @ + 钉钉私信双通道）。

必须静默：

- 每次定时检查无变化；
- CI pending、单次 retry、修复 push 后 new head；
- 普通 reviewer comment（仅在 GitHub 内回应）；
- PD/RD/QA 内部交接；
- 已 posted 或远端已有同 marker 的重复事件；
- 瞬时错误后正常 resume / retry 成功。

## 六、关单与升级

关单仍是人工门。收到关单请求后：

1. PD 只读核验无待接入资源、缺属性、缺陷、未决澄清、未合并 MR/CR 或未闭环关联单。
2. 必要时 RD/QA 内部复核。
3. finalizer 在唯一回复中给出证据并 @提单真人确认；提单方是数字人时升级辰羿(320687)和
   过载(484483)。
4. release，不由数字人代 finish。

达到 `JARVIS_PERSONA_MAX_ROUNDS` 时不再循环；finalizer 在唯一回复中说明超限原因并 @过载，
然后 release。close / normal / escalate 三种路径在同一主处理 run 内都不得产生第二条 RD
聚合回复；后续新的重要生命周期事件仍按 §五幂等更新。

## 七、人类入口

人类只需公开 `@TerraformRD` 或 `@WORKER_1783582458263`，bridge 默认从内部 PD 开始。

迁移期仍识别：

- 旧 `@TerraformPD` / `@WORKER_1783582374386` → 从内部 PD 开始；
- 旧 `@TerraformQA` / `@WORKER_1783582593461` → 从内部 QA 开始；
- 裸角色名没有 `@` 不触发。

响应作者始终是 TerraformRD。

## 八、旧公开接力仅入站兼容

旧评论可能带：

```text
[[PERSONA-HANDOFF:{"from":"terraform-pd","to":"terraform-rd","ticket":84129378,"action":"dev","round":1,"note":"补字段"}]]
```

bridge 保留解析、action 白名单、`from/to` 校验、last-one-wins、轮次和 ledger 逻辑，但只把它当
迁移期入站触发：

1. 消费旧 `to/action/note`；
2. 在同一个新 headless run 内从对应 internal_role 继续剩余链；
3. 最后由 RD finalizer 聚合回复一次；
4. 新流程不得再生成该格式。

统一 RD 作者下，历史 `from == to` 仍判 `self_addressed`；数字作者的旧格式缺合法 `from` 仍判
`bad_from`。这些仅是历史评论防重规则，不是新协作协议。

## 九、PersonaScheduler 迁移安全阀

PersonaScheduler 继续扫 `jarvis-idle` 与 tracker，用于消费历史接力和人类 @：

| 环境变量 | 默认值 | 含义 |
|---|---:|---|
| `JARVIS_PERSONA_WATCH` | `0` | 显式设 1 才启用 |
| `JARVIS_PERSONA_INTERVAL` | `600` | 轮询秒数 |
| `JARVIS_PERSONA_MAX_ROUNDS` | `6` | 迁移补位上限 |
| `JARVIS_PERSONA_TRACKER_SCAN` | `1` | 扫 RD/jarvis tracker |
| `JARVIS_PERSONA_LEGACY_TRACKER_SCAN` | `1` | 扫旧 PD/QA tracker |
| `JARVIS_PERSONA_NICKS` | 空 | 显示名兼容映射 |

台账 `.my-day/bridge/persona-ledger.json` 保留 `last_seen`、`processed`、`dispatch_count`、
`escalated`。只有 worker 真正启动才提交 ledger；pool 拒收不占坑，旧评论过期、终态、同工单
在飞实例照旧跳过。

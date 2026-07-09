# cap-tf-customer-probe

## 背景

Jarvis 入箱工单大头是 terraform-provider-alicloud 问题（缺资源/缺属性/行为不符/bug）。现状全是**被动**：
等客户踩坑提单，Jarvis 才 triage 修。仓库主人 2026-07-02 指令：让 Jarvis **像真实客户一样**——参考阿里云
Terraform 官方文档，用 terraform 创建真实资源，主动把潜在问题探出来；发现的问题提单到 terraform-alicloud
池（528766）指派 jarvis（WORKER_1782379562571），并**按危害定优先级**。

本能力把 Jarvis 从「被动修单」升级为「主动体检」：自己发现 → 自己建单 → 自己修 → 修完场景回归复跑。

## 能力定义与闭环

**合成客户探测（synthetic customer probing）**,两条腿(2026-07-03 重定义):

- **tier-0 静态三方一致性扫描**(以资源为单位):核对 **TF 文档 ↔ OpenAPI 文档 ↔ provider 源码**。机械部分做本地
  文档↔源码 diff(五类 gap);OpenAPI 一侧留 `judgment_queue` 交 skill 层双层查证。**只测已接入 TF 的面**。
- **tier-1 真实 apply 生命周期探测**(以场景为单位):以不同 persona 参考 `website/docs` 写 HCL,跑
  init→validate→plan→apply→re-plan→step2→import→destroy,用行为差异(validate/plan/apply 失败、永久 diff、
  意外重建、import 断链、destroy 残留)暴露 provider 潜在 bug。默认开启,region 默认 focus=eu-central-1。

```
probe 发现(tier-0 doc gap / tier-1 生命周期 bug) → 建单(tf_provider 528766) → aone-triage 认领 → 修复 → 发布 → 复跑回归 → 闭环
```

## 五层防线中的定位

| 层 | 防线 | 状态 |
|----|------|------|
| ① | 上游 GitHub issue 挖掘（从社区已报问题反查） | 未建（未来） |
| **②** | **合成客户探测 = 本能力** | **P0 本轮** |
| ③ | 真实工单回灌回归（regression-<aone-id> 场景） | 规则已立（scenario-authoring.md），直落外置 playground + 工单报备 |
| ④ | 发布前 RC 门禁（发版前全场景过一遍） | P2 |
| ⑤ | cloudspec/OpenAPI 覆盖矩阵驱动生成（探从未被示例覆盖的属性） | P3 |

本能力是②，并为③④⑤留接口（场景语料库结构、verdict schema、tier 分层、config 开关都可复用）。

## 架构组件

| 组件 | 路径 | 职责 |
|------|------|------|
| 配置 | `config/probe.json` | provider/tf 版本、regions(focus/matrix)、tiers(tier1.enabled/prepaid_guard)、limits、ticket、paths |
| runner | `bootstrap/probe.sh` | doctor / list / **tier0** / run / sweep;分层执行 + findings/env 分流 + verdict 落盘 |
| 场景语料库(独立 git 数据仓) | `tf_playground/<product>/<id>/` | 云产品维度两级布局;scenario.yaml(无 tier 键) + main.tf + checks.md (+ step2/);根解析 env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` > `workspace.sh dir tf_playground` > 默认 `<jarvis 父目录>/terraform_playground` |
| tier-0 fixture | `test/fixtures/probe/` | 手造迷你 doc + go,供解析器 hermetic 单测(五类 gap) |
| 技能 | `.claude/skills/tf-customer-probe/` | 全流程 + severity/ticket/authoring references |
| 循环 | `loops/tf-probe.md` | 触发→预检→tier0→tier1→分流→清理→Done runbook |
| 审计 | `runs/probe/<YYYYMMDD>-<HHMMSS>-<id>.{json,md}` + `<YYYYMMDD>-<HHMMSS>-tier0.{json,md}` | verdict 落盘(本地缓存,真源在 Aone) |
| draft 队列 | `escalation/probe-drafts/` | 工单草稿(未跟踪文件 = 待审信号) |

## tier 分层与当前开关(2026-07-03 重定义)

| tier | 动作 | 单位 | 当前 |
|------|------|------|------|
| tier-0 | 静态三方一致性扫描(TF 文档 ↔ OpenAPI 文档 ↔ 源码,只测已接入面) | 资源 | **默认可跑**(需本地 provider 仓) |
| tier-1 | 真实 apply 全生命周期(init→…→destroy) | 场景 | **默认开**(`tiers.tier1.enabled=true`) |

- tier-1 关掉(`enabled=false`)→ **封顶 plan-only**(init/validate/plan,不 apply,记 `tier1_disabled_plan_only`),非「降级 tier-0」。
- **成本白名单(tier1_allowlist)已撤销**——测试账号付费不设限;换 **prepaid 销毁性守门**:命中 PrePaid/Subscription 默认阻断(`prepaid_block`)。
- region:默认 `regions.focus`(eu-central-1,重点方向),`--region` 可切,`regions.matrix` 为未来矩阵候选。

## 危害分级 → 优先级映射

S1 紧急 / S2 高 / S3 中 / S4 低（详见 skill `references/severity-rubric.md`）→ Aone 优先级 紧急/高/中/低。
具体枚举值在首次真实建单（mode=file，2026-07-05 已毕业）用 a1 查证项目 528766 字段后固化。

## 护栏

- **销毁性(替代成本门)**：**prepaid 守门**——apply 前扫 plan 的 `*charge_type`/`*payment_type`,命中
  PrePaid/Subscription 默认阻断(包年包月多数无法 API 销毁,破坏零残留);场景 `allow_prepaid:true` 或
  `prepaid_guard=false` 豁免。强制 `destroy`(trap EXIT 兜底);`sweep` 残留核查;**残留即停并升级**。
- **工单**：mode=file 已毕业(2026-07-05,直接建 Aone 单);draft 冷启动保留为可回退开关
  (`ticket.mode=draft` 只写 `escalation/probe-drafts/`、不碰 Aone);建单前去重(a1 标签 + GitHub issue 只读);日上限
  `daily_new_tickets`(默认 100);统一 `jarvis-probe` 标签。
- **身份/账号**：a1 一律 jarvis(`bin/a1id`);probe 会话不 claim 工单;**只用环境注入的测试 AK/SK,绝不用生产账号**。
- **凭证**：AK/SK 绝不落日志 / verdict / draft / 工单;doctor 只报 set/unset。
- **tier-0 范围红线**:只核对已接入 TF 的面;未接入的资源/参数不报 gap(需求非 bug,走 tf_customer 需求路径)。
- **红线**：绝不碰生产存量资源;state 隔离在 `.my-day/probe/<ts>-<id>/`;只 destroy 本 run 自建 state 内的资源。

## 路线图

- **P0（本 MR）**：骨架 + 5 场景 + draft 模式 + tier-0 静态扫描(文档↔源码机械 diff)+ tier-1 默认真实 apply。
- **P1**：**tier-0 OpenAPI 侧机械化**——✅ **首发线 T0-mech 已落地**(2026-07-05,F3):`bootstrap/probe-meta.sh` 拉 amp
  `get_api_definition.py` 元数据(+`cache.sh` 7d 缓存),`probe.sh tier0` 抽 (product,version,action) 三元组 + 解析
  `StringInSlice`/`IntBetween`/`Default`,机械 diff 出六类 `api_gap_*`(deprecated_action/enum_superset/required/type/range/default),
  `judgment_queue` 从「人海判定」收窄为「机械层拿不准的疑点(prose/映射不上/OSS 无 action)」+ 抑制表/容差表/coverage_note 精度护栏;
  未竟:接 cloudspec/镇元 resourceType spec 做更深属性 diff、`--all` 全量轮换的规模化实跑(AMP 白名单凭证 2026-07-06 已配置,可全量实跑)。
  **源码 schema 嵌套深层解析**(当前只顶层,深挖 Elem 内层字段);terraform 二进制入 `bootstrap/install.sh` 依赖;
  工单回灌机制落地;`sweep` 接 aliyun CLI 按标签扫真实孤儿资源;自动建单**已毕业**(2026-07-05 主人拍板,采纳率
  7/8=87.5% 提前毕业,`ticket.mode=file`;draft 保留为可回退开关);a1 建单命令与优先级枚举固化。
- **P2**：cron/bridge 定时接入(**调度与修复闭环统一由 `cap-probe-fix-flywheel` F2 承接**);场景库批量扩容
  (website docs 全量 tier-0 覆盖);发布前 RC 门禁(接 terraform-changelog 发版流程,发版前全资源 tier-0 + 全场景
  tier-1 过一遍);~~upgrader persona~~ **已落地(2026-07-08)——`provider_version_from` 键 + upgrader dance**;
  `probe-corpus.sh` 生成 persona 变体(migrator/refactorer/ds-checker,同资源不同角度)。
- **P3**：cloudspec/OpenAPI 覆盖矩阵驱动属性组合生成(优先探从未被示例覆盖的属性);真实架构级组合场景;
  度量看板(发现数/采纳率/发现→修复周期,接 board.sh);**scale/throttling persona**(API 限流/大规模并发);
  **provider 新版本三件套**(发布检测 + config/playground pin 批量 bump + 触发全量轮);**`_quarantine` 自动出队**
  (再校验 + 修好归位);**acc-test / pr-review / release 场合蒸馏钩子**(把 KNOWLEDGE 契约挂到 provider-resource-review
  / terraform-pr-review / terraform-provider-release 的收尾流程);**KNOWLEDGE → 产品级 skill 毕业**
  (某产品 ≥15 条 + 消费 ≥3 次 → 起草 `.claude/skills/<product>-*`)。

## 置信度

- **high**：机械组件本轮 hermetic 测试覆盖(`test/probe_test.sh` 154 断言):新 config schema、场景键齐全且无 tier 键、
  pin 版本、list、run `--dry`、`tier1.enabled=false → plan-only` 封顶、region 解析优先级、prepaid 守门(阻断/放行/豁免)、
  **tier-0 解析器五类 gap 全抓 + 嵌套字段不误报**(fixture)、doctor 缺 terraform、sweep 残留。tier-0 亦对真实 provider 仓
  试扫验证(alicloud_vpc 等)。
- **待验**：tier-1 真实执行路径(apply/import/destroy)待真实测试账号凭证环境端到端验证;tier-0 OpenAPI 机械 diff(T0-mech)
  的**六类 api_gap_* 全被 hermetic 单测覆盖**(fixture 元数据桩,206 断言),真实标定亦对 07-03 首轮 8 资源逐条对照(复现
  ClassicLink deprecated、RAM/SG 零误报、OSS 正确降级 queue);AMP 白名单凭证 2026-07-06 已配置,`tier0 --all` 全量实拉元数据可端到端跑(机械面全量点亮)。

## 决策记录

### 2026-07-02（P0 初版）
1. **范围**：P0 全套骨架(config + runner + 5 场景 + skill + loop + cap + test)。
2. **云资源**：初版 tier-1 骨架完整但默认关。
3. **不写 Aone**：工单走 draft;建设单目标池 api_toolkit(2100304)。

### 2026-07-03（仓库主人分层重定义指令）
1. **tier-0 重定义**：从「init/validate/plan」改为**静态三方一致性扫描**(TF 文档 ↔ OpenAPI 文档 ↔ 源码),
   机械只做本地 文档↔源码 diff,OpenAPI 侧留 `judgment_queue`;**范围红线=只测已接入 TF 的面**。
2. **tier-1 重定义**：**必须真实 apply,默认开启**(不再「骨架就绪默认关」,授权已给);region 多选、eu-central-1 重点。
3. **成本门撤销**：测试账号费用不设限,删 `tier1_allowlist`;保留 **prepaid 销毁性守门**(PrePaid/Subscription 阻断)。
   建设单与 MR 仍待主人追认(目标池 api_toolkit 2100304)。

### 2026-07-03（场景语料库外置 terraform_playground）
仓库主人指令:tier-1 场景库从 jarvis 仓 `probes/scenarios/` **迁出到 `~/workspace/terraform_playground/`**,
按**云产品维度**一级归档、场景两级布局 `<product>/<id>/`。动机与设计:
1. **terraform 专家 skill 语料底座**:场景库作为未来「Terraform 专家 skill」的示例与验证语料,按云产品组织更利检索/扩容。
2. **回灌免 MR 闭环提速**:语料库在 jarvis 仓外,regression 场景由 jarvis **直接落** `<playground>/<product>/regression-<aone-id>/`
   (无需 worktree/MR),但须在对应工单评论**报备场景路径**供主人查验(取代原 `escalation/scenario-drafts` + 周批 MR)。
3. **provider 仓统一迁 `terraform_space` 归档**:本机工作区路径覆盖走 gitignored `config/workspaces.local.json`
   (已入 `bootstrap/master-allowlist` 长期免 worktree),base 配置不写绝对路径。
4. **repo 侧适配**(本 commit):`probe.sh` 场景根解析改 env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` >
   默认 `<jarvis 父目录>/terraform_playground`;目录布局升两级(list 加 PRODUCT 列、跨 product 同 id 报错、doctor 查场景根);
   `git rm` 仓内 `probes/`;测试 fixtures 迁两级布局保持 hermetic(不依赖真实 playground)。

### 2026-07-05（F3 首发线 T0-mech —— tier-0 OpenAPI 侧机械化,Aone 83929676）

规格:tier-0 三方一致性的 OpenAPI 侧从「verifier 人海判定」升级为「机械 diff 预筛 + verifier 只判疑点」。落地(developer 子代理):
1. **元数据层** `bootstrap/probe-meta.sh`:薄封装 `amp-resource-metadata/get_api_definition.py`,fetch/cached-fetch(cache.sh 7d)/clear/available;
   probe.sh **不直接调 python**;网络/venv/凭证不可用 → 干净降级(退非零 + 明确提示)。
2. **机械 diff** `probe.sh tier0`:(product,version,action) 三元组抽取(`RpcPost("P","V",action)` 同行可得,OSS SDK 风格抽不到进 queue)+
   源码 `StringInSlice`/`IntBetween`/`Default` 解析(解析不动标 unknown 进 queue,不猜)+ 六类 `api_gap_*` 机械 diff。**精度护栏**:
   snake→Camel 精确映射(映射不上 queue)、抑制表 `suppress_params`+容差表 `type_tolerance`(命中入 `suppressed[]` 可审计)、
   废弃双轨对不报、TF 更严记 `coverage_notes[]`。**CLI**:`--no-mech`(降级现行为)/`--all`(全资源清单)/`--limit N`/`--rotate N`(LRU 轮换)。
3. **真实标定(验收核心)**:对 07-03 首轮 8 资源标定——机械层复现 `alicloud_vpc` ClassicLink `api_gap_deprecated_action` ×2
   (= 单 83881282);`alicloud_ram_role`/`alicloud_security_group` 三方一致 → 零新增误报;`alicloud_oss_bucket` SDK 风格 → `no_triple` queue
   (enum 超集需元数据可得,符合规格);`dns_hostname_status` convert 改名 → `unmapped_params` queue(机械层不猜)。hermetic 单测 206 断言全绿。
   **约束**:本机无 AMP 白名单凭证,全量实拉元数据(`tier0 --all`)待有凭证环境端到端验证(→ **2026-07-06 已配置凭证,约束解除**,见下);机械管道由 fixture 元数据桩完整覆盖。

### 2026-07-06（主人指令四条,详见 `cap-probe-fix-flywheel.md` 决策记录）

与本能力相关三条:①**成本门收窄为「值敏感」**——prepaid 守门不再对所有 PrePaid/Subscription 一刀切阻断,收窄为只挡真正值敏感
(不可逆计费 / 无法 API 销毁);**订阅类资源若无对应 data source(无 ds 可读回校验)则放行 `apply`**(等价 `allow_prepaid`)。
②**场景语料库 git 化**——`terraform_playground` 升级为 git 仓 **`terraflow/tf_playground`**,**直推 master + 工单报备**模型
(取代原仓外裸目录无 MR 方案,同时天然收敛多机语料分叉)。③**AMP 白名单凭证已配置**——T0-mech 机械面全量点亮,
`tier0 --all` 全量实拉元数据可端到端跑,解除上条 07-05「待有凭证环境」约束。(第四条 `bridge/run.sh` 单一入口属 bridge 侧,不涉本能力。)

### 2026-07-08（探测链 v2 —— 归档自动处理 + persona 扩容 + 知识蒸馏契约 + board 适配,锚单 Aone 84080465）

四块一并落地(worktree `probe-chain-v2` → MR):

1. **归档自动处理**——新子命令 `probe.sh archive [--dry]`(幂等):draft frontmatter `status=filed/rejected*` 收入
   `escalation/probe-drafts/archived/`;`runs/probe/` 顶层 verdict 超 `limits.audit_retention_days`(默认 60)搬
   `runs/probe/archive/<YYYYMM>/`(排 `ledger.jsonl`/`*-summary.md`);`.my-day/probe/<ts-sid>/` 过期 + tfstate 空 → rm
   (排 `.plugin-cache/`/`manual-*`/索引文件);plugin-cache 陌生版本只报体积不删;pending drafts + `_quarantine/` +
   `origin: generated` 未校订三类待办清单一次打印。新增 config `paths.drafts_archived`、`limits.audit_retention_days`、
   `limits.workdir_retention_days`;**所有新 config 键必有代码内默认值**(config 分裂防御)。**verdict 同日覆盖修复**:
   审计副本文件名加 `HHMMSS`(`<YYYYMMDD>-<HHMMSS>-tier0.json`/`<YYYYMMDD>-<HHMMSS>-<sid>.json`),同日多轮不再互相覆盖;
   board.sh 侧 findings 周窗口按 (code,resource,attribute) 去重相应适配。台账 `runs/probe/ledger.jsonl` append-only:
   tier0/tier1 finalize 追加 + skill Step D.5 建单追加 + archive 追加(dry=`archive_dry`);LRU 索引
   `.my-day/probe/t1-last-run.json`(仅推进到 plan 完成及以后才更新,`probe.sh list` 行尾列 LAST_RUN 追加)。
2. **persona 扩容 5 类 + 新键**——`migrator`/`upgrader`/`refactorer`/`drifter`/`ds-checker`(+`ci-runner` taxonomy 骨架);
   scenario.yaml 新键(全部可选叠加):`steps: step2,step3` CSV(泛化 update_step)+ `step<N>_expect: no_changes|changed|fail`
   逐步声明期望 → `refactor_replace`(no_changes 却出 diff;delete+create S1);`expect_fail: validate|plan|apply`
   +`expect_error_contains` 四态(`expected`/`expected_but_error_mismatch` S3/`late_validation` S3/`expected_fail_missed` S2);
   `provider_version_from: <old-pin>` upgrader 版本 dance → `upgrade_diff`(S2;delete+create S1);
   `drift_cli: aliyun <product> <Action> ...` drifter 五重护栏(tokenize+元字符黑名单/受限占位符/凭证显式映射
   `ALIBABA_CLOUD_ACCESS_KEY_ID/SECRET/REGION_ID`/action 白名单锚在 jarvis config `tiers.tier1.drift_action_allow`
   /默认 `drift_enabled: false` 关)→ `drift_undetected`(S2;Claude 复核后静默错配可升 S1)。判定抽成纯函数
   (`_expect_fail_verdict` 源可单测,先例 `_prepaid_should_block`)。
3. **知识蒸馏契约**——新 reference `.claude/skills/tf-customer-probe/references/knowledge-distillation.md`
   (**跨 skill 单点**,先例 `provider-resource-dev/references/zhenyuan-verification.md`)。数据在
   `tf_playground/<product>/KNOWLEDGE.md`(五节:命名/参数 quirk/生命周期/API 行为/报错→原因→解法),条目格式
   `- [YYYY-MM-DD][来源: 链接/路径] <可执行的产品级事实>`。**三个触发点**:①probe 轮 Step E 收尾;
   ②**aone-triage bookend 收尾**(客户单场合的蒸馏钩子挂在 aone-triage 主流程——评审阻断项,不能只挂 probe 侧);
   ③provider-resource-dev 完成开发后。收录判据=可执行/跨场景复用/非文档已明示;消费约定=dev/review/probe 涉及
   产品 X 先读 `<playground>/<X>/KNOWLEDGE.md`(存在即读);对外流转按 CLAUDE.md #5 禁品清单 sanitize;毕业标准=
   某产品 ≥15 条 + 消费 ≥3 次 → 起草产品级 skill(cap P3)。
4. **board 适配 + bridge 修复**——board.sh 加 archived/ drafts + verdict archive/<YYYYMM>/ 支持 + findings 周窗口
   (code,resource,attribute) 去重;bridge `_probe_prompt` 删「mode=draft 人审」硬编码改按 SKILL Step C/D 与
   `config/probe.json ticket.mode` 执行,尾部追加 `probe.sh archive` + 按 knowledge-distillation 契约蒸馏两步;
   `_DailyScheduler._run_once` bool 契约(仅 queue_full 返 False 当日不 mark,5min tick 重试)+ probe- 前缀 final
   文本落 `runs/probe/<rid>-summary.md`。

## 关联

- **修复侧闭环设计**：见 `escalation/cap-probe-fix-flywheel.md`(探测→建单→修复→验证→发布→回灌 六段飞轮 + F0–F4 阶段计划)。
- probe 发现的 provider 问题单：落 tf_provider 池 528766，指派 WORKER_1782379562571，标签 jarvis-probe。
- 相关文件：`config/probe.json`、`bootstrap/probe.sh`、外置 `terraform_playground/`(场景库)、
  `.claude/skills/tf-customer-probe/`、`loops/tf-probe.md`、`test/probe_test.sh`。
- 相关技能：aone-triage（接单修复）、provider-resource-dev（资源开发）、terraform-pr-review（PR 评审）、
  invoke-terraform-acc-test-remote（远程 AccTest）、terraform-changelog（发版，未来接 RC 门禁）。
- 工作纪律：CLAUDE.md #4 工作区登记、#6 身份纪律；autonomy.md（probe 产出的建单/CR 受策略约束）。

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
| ③ | 真实工单回灌回归（regression-<aone-id> 场景） | 规则已立（probes/README.md），P1 落地 |
| ④ | 发布前 RC 门禁（发版前全场景过一遍） | P2 |
| ⑤ | cloudspec/OpenAPI 覆盖矩阵驱动生成（探从未被示例覆盖的属性） | P3 |

本能力是②，并为③④⑤留接口（场景语料库结构、verdict schema、tier 分层、config 开关都可复用）。

## 架构组件

| 组件 | 路径 | 职责 |
|------|------|------|
| 配置 | `config/probe.json` | provider/tf 版本、regions(focus/matrix)、tiers(tier1.enabled/prepaid_guard)、limits、ticket、paths |
| runner | `bootstrap/probe.sh` | doctor / list / **tier0** / run / sweep;分层执行 + findings/env 分流 + verdict 落盘 |
| 场景语料库 | `probes/scenarios/<id>/` | scenario.yaml(无 tier 键) + main.tf + checks.md (+ step2/) |
| tier-0 fixture | `test/fixtures/probe/` | 手造迷你 doc + go,供解析器 hermetic 单测(五类 gap) |
| 技能 | `.claude/skills/tf-customer-probe/` | 全流程 + severity/ticket/authoring references |
| 循环 | `loops/tf-probe.md` | 触发→预检→tier0→tier1→分流→清理→Done runbook |
| 审计 | `runs/probe/<日期>-<id>.{json,md}` + `<日期>-tier0.{json,md}` | verdict 落盘(本地缓存,真源在 Aone) |
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
具体枚举值在首次真实建单（mode=file 毕业后）用 a1 查证项目 528766 字段后固化。

## 护栏

- **销毁性(替代成本门)**：**prepaid 守门**——apply 前扫 plan 的 `*charge_type`/`*payment_type`,命中
  PrePaid/Subscription 默认阻断(包年包月多数无法 API 销毁,破坏零残留);场景 `allow_prepaid:true` 或
  `prepaid_guard=false` 豁免。强制 `destroy`(trap EXIT 兜底);`sweep` 残留核查;**残留即停并升级**。
- **工单**：draft 冷启动(当前不写 Aone);建单前去重(a1 标签 + GitHub issue 只读);日上限
  `daily_new_tickets`(默认 100);统一 `jarvis-probe` 标签。
- **身份/账号**：a1 一律 jarvis(`bin/a1id`);probe 会话不 claim 工单;**只用环境注入的测试 AK/SK,绝不用生产账号**。
- **凭证**：AK/SK 绝不落日志 / verdict / draft / 工单;doctor 只报 set/unset。
- **tier-0 范围红线**:只核对已接入 TF 的面;未接入的资源/参数不报 gap(需求非 bug,走 tf_customer 需求路径)。
- **红线**：绝不碰生产存量资源;state 隔离在 `.my-day/probe/<ts>-<id>/`;只 destroy 本 run 自建 state 内的资源。

## 路线图

- **P0（本 MR）**：骨架 + 5 场景 + draft 模式 + tier-0 静态扫描(文档↔源码机械 diff)+ tier-1 默认真实 apply。
- **P1**：**tier-0 OpenAPI 侧机械化**(接 cloudspec/镇元 spec 自动 diff,让 `judgment_queue` 从人判走向机判);
  **源码 schema 嵌套深层解析**(当前只顶层,深挖 Elem 内层字段);terraform 二进制入 `bootstrap/install.sh` 依赖;
  工单回灌机制落地;`sweep` 接 aliyun CLI 按标签扫真实孤儿资源;自动建单毕业条件(累计 ≥10 draft 且采纳率 ≥90% 后
  `ticket.mode` 切 `file`);a1 建单命令与优先级枚举固化。
- **P2**：cron/bridge 定时接入;场景库批量扩容(website docs 全量 tier-0 覆盖);发布前 RC 门禁
  (接 terraform-changelog 发版流程,发版前全资源 tier-0 + 全场景 tier-1 过一遍);upgrader persona(版本升级 state 兼容探测)。
- **P3**：cloudspec/OpenAPI 覆盖矩阵驱动属性组合生成(优先探从未被示例覆盖的属性);真实架构级组合场景;
  度量看板(发现数/采纳率/发现→修复周期,接 board.sh)。

## 置信度

- **high**：机械组件本轮 hermetic 测试覆盖(`test/probe_test.sh` 154 断言):新 config schema、场景键齐全且无 tier 键、
  pin 版本、list、run `--dry`、`tier1.enabled=false → plan-only` 封顶、region 解析优先级、prepaid 守门(阻断/放行/豁免)、
  **tier-0 解析器五类 gap 全抓 + 嵌套字段不误报**(fixture)、doctor 缺 terraform、sweep 残留。tier-0 亦对真实 provider 仓
  试扫验证(alicloud_vpc 等)。
- **待验**：tier-1 真实执行路径(apply/import/destroy)待真实测试账号凭证环境端到端验证;tier-0 OpenAPI 侧判定当前靠 skill 层人判(P1 机械化)。

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

## 关联

- probe 发现的 provider 问题单：落 tf_provider 池 528766，指派 WORKER_1782379562571，标签 jarvis-probe。
- 相关文件：`config/probe.json`、`bootstrap/probe.sh`、`probes/`、`.claude/skills/tf-customer-probe/`、
  `loops/tf-probe.md`、`test/probe_test.sh`。
- 相关技能：aone-triage（接单修复）、provider-resource-dev（资源开发）、terraform-pr-review（PR 评审）、
  invoke-terraform-acc-test-remote（远程 AccTest）、terraform-changelog（发版，未来接 RC 门禁）。
- 工作纪律：CLAUDE.md #4 工作区登记、#6 身份纪律；autonomy.md（probe 产出的建单/CR 受策略约束）。

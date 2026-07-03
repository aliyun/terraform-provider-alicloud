# cap-tf-customer-probe

## 背景

Jarvis 入箱工单大头是 terraform-provider-alicloud 问题（缺资源/缺属性/行为不符/bug）。现状全是**被动**：
等客户踩坑提单，Jarvis 才 triage 修。仓库主人 2026-07-02 指令：让 Jarvis **像真实客户一样**——参考阿里云
Terraform 官方文档，用 terraform 创建真实资源，主动把潜在问题探出来；发现的问题提单到 terraform-alicloud
池（528766）指派 jarvis（WORKER_1782379562571），并**按危害定优先级**。

本能力把 Jarvis 从「被动修单」升级为「主动体检」：自己发现 → 自己建单 → 自己修 → 修完场景回归复跑。

## 能力定义与闭环

**合成客户探测（synthetic customer probing）**：以不同 persona（新手/组合/更新/导入）为视角，参考
`website/docs` 官方示例写 HCL，跑 terraform 生命周期（init→validate→plan→apply→re-plan→import→destroy），
用行为差异（validate/plan/apply 失败、永久 diff、意外重建、import 断链、destroy 残留）暴露 provider 潜在 bug。

```
probe 发现 → 建单(tf_provider 528766) → aone-triage 认领 → 修复 → 发布 → 场景回归复跑 → 闭环
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
| 配置 | `config/probe.json` | provider/tf 版本、tier 开关、tier1_allowlist、limits、ticket、paths |
| runner | `bootstrap/probe.sh` | doctor / list / run / sweep;分层执行 + findings/env 分流 + verdict 落盘 |
| 场景语料库 | `probes/scenarios/<id>/` | scenario.yaml + main.tf + checks.md (+ step2/) |
| 技能 | `.claude/skills/tf-customer-probe/` | 单场景全流程 + severity/ticket/authoring references |
| 循环 | `loops/tf-probe.md` | 触发→预检→选场景→执行→分流→清理→Done runbook |
| 审计 | `runs/probe/<日期>-<id>.{json,md}` | verdict 落盘（本地缓存,真源在 Aone） |
| draft 队列 | `escalation/probe-drafts/` | P0 工单草稿（未跟踪文件 = 待审信号） |

## tier 分层与当前开关

| tier | 动作 | 成本 | 当前 |
|------|------|------|------|
| tier-0 | init / validate / plan | 零资源创建 | **默认开** |
| tier-1 | apply / re-plan / import / destroy（免费 allowlist 资源） | 免费资源 | **骨架就绪，默认关**（`tiers.tier1_enabled=false`） |
| tier-2 | 付费资源 | 计费 | **永不**（P0 硬封顶，绝不放行） |

有效 tier = min(场景声明, `--tier` 请求, 配置允许最高)。`tier1_enabled=false` 时一律降级 tier-0 并记 `tier_downgraded`。

## 危害分级 → 优先级映射

S1 紧急 / S2 高 / S3 中 / S4 低（详见 skill `references/severity-rubric.md`）→ Aone 优先级 紧急/高/中/低。
具体枚举值在首次真实建单（mode=file 毕业后）用 a1 查证项目 528766 字段后固化。

## 护栏

- **成本**：只 `cost: free` 资源；tier-1 双门（config 开关 + allowlist 硬门）；强制 `destroy`（trap EXIT 兜底）；
  `sweep` 残留核查；**残留即停并升级**。
- **工单**：draft 冷启动（P0 不写 Aone）；建单前去重（a1 标签 + GitHub issue 只读）；日上限
  `daily_new_tickets`（默认 3）；统一 `jarvis-probe` 标签。
- **身份**：a1 一律 jarvis（`bin/a1id`）；probe 会话不 claim 工单。
- **凭证**：AK/SK 绝不落日志 / verdict / draft / 工单；doctor 只报 set/unset。
- **红线**：绝不碰生产存量资源；state 隔离在 `.my-day/probe/<ts>-<id>/`；只 destroy 本 run 自建 state 内的资源。

## 路线图

- **P0（本 MR）**：骨架 + 5 免费场景 + draft 模式 + tier-0 默认（tier-1 骨架就绪默认关）。
- **P1**：主人授权后开 tier-1；terraform 二进制入 `bootstrap/install.sh` 依赖；工单回灌机制落地；
  `sweep` 接 aliyun CLI 按标签扫真实孤儿资源；自动建单毕业条件（累计 ≥10 draft 且采纳率 ≥90% 后
  `ticket.mode` 切 `file`）；a1 建单命令与优先级枚举固化。
- **P2**：cron/bridge 定时接入；场景库批量扩容（website docs 全量生成 tier-0 语料）；发布前 RC 门禁
  （接 terraform-changelog 发版流程，发版前全场景过一遍）；upgrader persona（provider 版本升级 state 兼容探测）。
- **P3**：cloudspec/OpenAPI 覆盖矩阵驱动属性组合生成（优先探从未被示例覆盖的属性）；真实架构级组合场景；
  度量看板（发现数/采纳率/发现→修复周期，接 board.sh）。

## 置信度

- **high**：全部组件本轮 hermetic 测试覆盖（`test/probe_test.sh`：config 键、场景键齐全、tier-1 resources ⊆
  allowlist、pin 版本、list、`--dry` 降级、allowlist 拒绝、doctor 缺 terraform、sweep 残留）。
- **待验**：tier-1 真实执行路径（apply/import/destroy）待真实凭证 + `tier1_enabled=true` 环境验证；本机当前
  为 tier-0/dry 覆盖。

## 决策记录（2026-07-02）

三项默认值（供仓库主人追认）：
1. **范围**：P0 全套骨架（config + runner + 5 场景 + skill + loop + cap + test），不只做最小 PoC。
2. **云资源**：只到 tier-0（免 apply）为止；tier-1 骨架完整但默认关，等主人授权再开。
3. **本轮不写 Aone**：工单走 draft 落 `escalation/probe-drafts/`；本能力的**建设单**与 MR 待主人追认，
   建设单目标池 api_toolkit（2100304）。

## 关联

- probe 发现的 provider 问题单：落 tf_provider 池 528766，指派 WORKER_1782379562571，标签 jarvis-probe。
- 相关文件：`config/probe.json`、`bootstrap/probe.sh`、`probes/`、`.claude/skills/tf-customer-probe/`、
  `loops/tf-probe.md`、`test/probe_test.sh`。
- 相关技能：aone-triage（接单修复）、provider-resource-dev（资源开发）、terraform-pr-review（PR 评审）、
  invoke-terraform-acc-test-remote（远程 AccTest）、terraform-changelog（发版，未来接 RC 门禁）。
- 工作纪律：CLAUDE.md #4 工作区登记、#6 身份纪律；autonomy.md（probe 产出的建单/CR 受策略约束）。

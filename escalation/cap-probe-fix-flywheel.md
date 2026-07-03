# cap-probe-fix-flywheel

## 背景

项目终态 = 数字员工自治运转。探测能力（`cap-tf-customer-probe`）P0 已建成并实战验证：2026-07-03 首轮
产出 **7 张 probe 工单**（83881282 / 83881291 / 83881297 / 83881301 / 83881314 / 83882090 / 83882094,
均指派 jarvis、打 `jarvis-probe` 标签）+ **1 份 verify 否决样本**（判定层挡下的伪 finding)。

「探测出来」只是半程。本 cap 设计**「发现的问题被 jarvis 自己修完并回归验证」**的完整飞轮——从发现到发布再到
回灌复跑,形成自闭环,让整个项目在三个人工硬门之外自动运转。

## 探测侧完成度判定（如实）

**已建成（实战验证）**：

- tier-0 机械扫描 + OpenAPI 判定管道（文档↔源码 diff + `judgment_queue`）
- tier-1 真实 apply 全生命周期（init→…→destroy）
- findings / env_issues 严格分流
- draft → 人审 → 建单管道（`escalation/probe-drafts/`）
- 去重台账（`jarvis-probe` 标签 + draft filed 回写）
- 护栏（prepaid 销毁性守门 / 强制 destroy / sweep 残留核查 / a1id 身份 / 测试账号边界）
- 审计与 **154 项 hermetic 测试**

**未完成（自动运转三缺口）**：

1. **调度**——目前手动触发,需 bridge 定时化。
2. **规模**——5 场景 / 8 资源 vs 全 provider 资源面,需语料生成器 + cloudspec 覆盖矩阵。
3. **判定成本**——OpenAPI 侧靠 verifier 子代理人海,需 cloudspec 机械 diff 预筛降本。

## 飞轮六段架构

| 段 | 触发 | 执行机件 | 现状 | 缺口 | 人工门 |
|----|------|----------|------|------|--------|
| ① 发现 Probe | 定时 + 新版本发布 | `tf-customer-probe` skill + `probe.sh`(tier-0/tier-1) | 已建 | 缺 bridge cron 接入 | 无 |
| ② 立项 File | findings → 分级 → 去重 → 建单 528766 指派 jarvis | draft 管道 + 建单模板/标签/上限(100) | draft 管道已建、模板/标签/上限齐 | `mode=draft→file` 毕业(条件:累计 ≥10 draft 且人审采纳率 ≥90%;2026-07-03 首轮 7/8=87.5%,下一轮达标即翻) | 毕业前人审 draft(临时门) |
| ③ 修复 Fix | bridge ScanScheduler 扫池(probe 单指派 jarvis 被 `scan.sh` 自然扫到) → headless dispatch → aone-triage 认领 → **按 `probe-ticket-routing` 路由** | (a) provider 代码修 → `provider-resource-dev` → fork+UT → `invoke-terraform-acc-test-remote` 验收 → GitHub PR(`github-identity` 硬门,`api-tool-agent`);(b) TF 文档修 → 同 PR 路径 docs-only;(c) 上游协作 → cloudspec_gap 等 `submit_only` 转发;(d) 需实验定性 → 先跑 tier-1 变体场景再归入 a/b/c | 机件**全部已存在** | 只缺**路由规范**(本 commit 落地) | upstream PR merge(maintainer) |
| ④ 验证 Verify-fix | PR 合并 → master 复验 | tier-0 重扫该项应消失 / tier-1 场景复跑应绿 →「已修未发布」→ 发布后复跑绿 → `claim.sh finish`(jarvis-done) | 靠工单溯源字段映射回场景/资源,无需新 runner 子命令 | 复验编排规范(routing reference 内定义)+ 状态机落 tag/评论 | 无 |
| ⑤ 发布 Release | changelog 聚合 → 发版 | `terraform-changelog` skill(已有) | 已有 | 发布前 RC 门禁 = 全场景语料过一遍(P2) | release_prod(autonomy.md 永久停止项) |
| ⑥ 回灌 Regress | 每张修复完成的 probe 单 + 每张真实客户单 | `regression-<aone-id>` 场景**直落**外置 `terraform_playground/<product>/regression-<aone-id>/`(仓外,无需 worktree/MR)+ 工单评论报备场景路径 | 规则已立(2026-07-03 外置简化) | 收尾清单挂钩 | 工单评论报备(仓库主人查验) |

## 三个永久人工硬门（其余全自动）

1. **jarvis 仓 MR 合并**（仓库主人）。
2. **terraform-provider-alicloud upstream PR merge**（团队 maintainer;jarvis 可自动 remind / 催）。
3. **release_prod**（正式发布,autonomy.md 永久停止项）。

外加 **draft 毕业前的人审 draft**（临时门,`mode` 毕业翻 `file` 后撤）。

## 兜底与健康度

- **escalate 面不变**：`low_conf` / `verify_fail` / `redline` / `missing_capability` + probe 新增 `destroy_fail` / 残留。
- **僵死收敛**：`reconcile.sh` / watchdog / coord 孤儿收养。
- **飞轮度量**（P3 接 `board.sh`）：每轮发现数、draft 采纳率、建单→修复→发布周期、回归通过率、场景覆盖率。

## 阶段计划

- **F0（本 commit）**：飞轮 cap 设计 + aone-triage probe 单路由 reference——已有的 7 张单下一轮 triage 即可正确消化。
- **F1（首证）**：选 **83881291**（dns_hostname_status,最干净的代码修）人工触发走完整链一遍:
  `provider-resource-dev` 修复 → acc test → PR → merge 后 master 复验 → 回灌 regression 场景 → finish。
  产出飞轮首个端到端样本与卡点清单。
- **F2（自动化）**：bridge ScanScheduler 增 probe 定时轮次（日常 tier-0 批扫 + tier-1 轮换 + provider 新版本事件触发）;
  `mode=file` 毕业翻开关;复验步骤进 triage 收尾清单。
- **F3（规模化）**：场景语料生成器（website docs 全量 → tier-0 语料）、cloudspec 覆盖矩阵驱动生成与 OpenAPI 机械
  diff 预筛、发布前 RC 门禁、度量看板。
- **F4（目标态）**：无人值守运转,人只守三硬门与 escalation 队列。

## 决策记录

- **2026-07-03**：仓库主人提问「探测部分是否建成?如何让 jarvis 自闭环解决探测出的问题,让整个项目自动运转」并指示继续设计
  → 本 cap（飞轮总设计）+ F0 落地（aone-triage probe 单路由 reference)。
- **待主人决策**：
  1. `mode=file` 翻开关时点（建议下一轮采纳率 ≥90% 后)。
  2. bridge probe cron 频率（建议每日一轮 tier-0 批扫 + 3 场景 tier-1 轮换)。
  3. F1 首证启动授权。

## 关联

- 探测能力本体：`cap-tf-customer-probe.md`。
- 能力建设单：https://project.aone.alibaba-inc.com/v2/project/2100304/req/83879813
- MR：https://code.alibaba-inc.com/terraflow/jarvis/codereview/28364380
- 首轮 7 张 probe 单：83881282 / 83881291 / 83881297 / 83881301 / 83881314 / 83882090 / 83882094。
- F0 路由 reference：`.claude/skills/aone-triage/references/probe-ticket-routing.md`。

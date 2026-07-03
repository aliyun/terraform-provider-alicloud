# tf-probe 合成客户探测循环

> 主动探测 terraform-provider-alicloud 潜在问题的 runbook。与 aone-triage 形成闭环：
> probe 产出的工单被 triage loop 扫到后按正常流程认领修复；修复发布后对应场景/资源复跑即回归验证。
> 能力全景/路线图见 `escalation/cap-tf-customer-probe.md`；全流程技能见 `.claude/skills/tf-customer-probe`。

## 分层(2026-07-03 重定义)

- **tier-0 = 静态三方一致性扫描**(TF 文档 ↔ OpenAPI 文档 ↔ provider 源码),**以资源为单位**,不跑 terraform。
- **tier-1 = 真实 apply 全生命周期探测**(默认开启),**以场景为单位**,region 默认 focus=eu-central-1(重点方向)。
- **执行顺序建议:先 tier-0 扫本轮涉及资源(快、零成本)→ 再 tier-1 跑场景生命周期。**

---

## 一、触发

| 方式 | 说明 |
|------|------|
| 手动触发 | 用户在会话中执行 `/tf-probe` 或直接发指令「跑一轮场景探测 / 三方一致性扫描」 |
| provider 新版 | 新版本发布后应尽快全量跑一轮(版本升级易引入 state 不兼容 / 新永久 diff / 文档漂移) |
| cron / bridge | P2 才接入,本轮**不接** |

---

## 二、预检（doctor）

```bash
bootstrap/probe.sh doctor
```

- 全绿退 0 才继续。
- 缺 **terraform** → tier-1 不可跑;env 问题,**不硬闯**,按 loops/self-improve.md 记缺口 / escalation。
- 缺 **本地 provider 仓**(website/docs + alicloud) → tier-0 静态扫描不可跑(WARN)。
- 缺凭证 → tier-1 只能 plan-only(plan 亦需凭证;无凭证 plan 记 `no_creds`,不是 finding)。

---

## 三、tier-0 静态扫描(先跑)

```bash
bootstrap/probe.sh tier0                 # 无参 = 全部场景 resources 并集(去重)
bootstrap/probe.sh tier0 alicloud_vpc    # 指定资源
bootstrap/probe.sh tier0 --dry           # 只列将扫资源 + 文档/源码文件存在性
```

- 机械产出 `findings`(五类 doc↔source gap)+ 每资源一条 `judgment_queue`(OpenAPI 侧待查证 + source 实际 API action)。
- **OpenAPI 判定接棒(Claude/verifier)**:对 `judgment_queue` 逐资源走双层查证惯例,核验 API 参数/枚举/行为 vs TF 文档;
  产 `doc_api_gap`。**范围红线:只核对已接入 TF 的面,未接入的资源/参数不报 gap(需求非 bug)**。
- 落盘 `runs/probe/<日期>-tier0.json` + 人读 md。退出码:0 无 findings / 1 有 findings / 2 runner 错误。

---

## 四、tier-1 场景生命周期探测

### 4.1 场景选择
```bash
bootstrap/probe.sh list
```
- 默认**最久未跑优先**:看 `runs/probe/` 各场景最新日期,挑最旧的;单轮 ≤ `config.limits.max_scenarios_per_run`(默认 3)。

### 4.2 执行
```bash
bootstrap/probe.sh run <id>                     # region 默认 config.regions.focus(eu-central-1)
bootstrap/probe.sh run <id> --region cn-hangzhou # 切 region(重点 eu-central-1,但保留多 region)
bootstrap/probe.sh run <id> --dry               # 只看步骤计划(region 解析、prepaid 守门)
```
- region 优先级:`--region` > scenario.yaml `region:` > config `regions.focus` > 环境 `ALICLOUD_REGION`。
- **prepaid 销毁性守门**:apply 前扫 plan,命中 PrePaid/Subscription 计费类型默认阻断(`prepaid_block`);场景 `allow_prepaid: true` 豁免。
- `tiers.tier1.enabled=false` → **封顶 plan-only**(不 apply,env_issue `tier1_disabled_plan_only`)。
- 读产出 `verdict.json`(工作目录 + `runs/probe/<日期>-<id>.json` + md)。退出码:0/1/2/**3 清理失败(最高优先级人工介入)**。
- **findings 判定是 Claude 的职责**:逐 finding 对照 evidence、`references/severity-rubric.md`、CHANGELOG Unreleased 段(已修掉的标「已修复未发布」不建单)。

---

## 五、产物分流

| 产物 | 去向 |
|------|------|
| `findings`(provider 疑似 bug,tier-0 或 tier-1) | 去重后 → draft(当前 mode)到 `escalation/probe-drafts/`;毕业后 → adhoc-intake 建单(tf_provider 528766,见 skill) |
| `env_issues`(凭证/网络/`prepaid_block`/`tier1_disabled_plan_only`) | **不建单**;缺 terraform / 缺 provider 仓等能力类走 loops/self-improve.md 或 escalation |

- 去重:a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词;GitHub 上游 open issues 只读检索;重复则追加 evidence 不新建。
- **probe 会话本身不 claim 工单**:产出的新单由 aone-triage loop 后续接管。

---

## 六、清理核查（sweep，残留即停）

```bash
bootstrap/probe.sh sweep
```

- 无残留退 0。
- **有残留退 1 → 立即停并升级**:`.my-day/probe/*/terraform.tfstate` 有非空 `resources`,按提示手动 `terraform destroy`
  或按 `managed_by=jarvis-probe` 标签用 aliyun CLI 清理(prepaid 守门已尽量把不可销毁资源挡在 apply 前)。

---

## 七、Done — 本轮结束标准

- tier-0:本轮涉及资源都扫过,`runs/probe/<日期>-tier0.json` 落盘,judgment_queue 已交判定。
- tier-1:每个跑过的场景都有 `verdict.json`;`sweep` 零残留。
- 每条 provider finding 都已去重并落 draft(或标注「已修复未发布」/「上游已报」跳过)。
- 会话汇报:`tier0 资源数/findings/judgment 数` + `tier1 场景数/findings/draft 数/env 数`,逐 draft 附路径与建议优先级。

---

## 八、与 aone-triage 的闭环

```
tf-probe 发现问题(tier-0 doc gap / tier-1 生命周期 bug) → draft/建单(tf_provider 528766, tag jarvis-probe)
        ↓
aone-triage loop 扫到该单 → claim → provider-resource-dev/修复 → PR → 发布
        ↓
发布后:tf-probe 复跑对应资源(tier0)/场景(tier1) → 无 finding = 回归通过(闭环)
```

真实客户工单也应回灌为 `regression-<aone-id>` 场景(probes/README.md),成为发版前回归项。

### 修复侧衔接

probe 单指派 jarvis,会被 `scan.sh` 自然扫到进 triage;aone-triage 按 **四分类路由**(provider 代码修 / TF 文档修 /
上游协作 / 需实验定性)修复,PR 合并后按溯源**复验关单**,并回灌 regression 场景——完整飞轮见
`escalation/cap-probe-fix-flywheel.md`,路由与复验状态机见 `.claude/skills/aone-triage/references/probe-ticket-routing.md`。

---

## 九、工具链速查

| 工具 | 作用 |
|------|------|
| `bootstrap/probe.sh doctor` | 环境预检(terraform/jq/凭证/config/本地 provider 仓) |
| `bootstrap/probe.sh list` | 列全部场景(id/persona/resources/detect) |
| `bootstrap/probe.sh tier0 [alicloud_xxx ...] [--dry]` | tier-0 静态三方一致性扫描 |
| `bootstrap/probe.sh run <id> [--region r] [--dry] [--keep]` | tier-1 真实 apply 生命周期探测 |
| `bootstrap/probe.sh sweep` | 扫残留 state,残留退 1 |
| `config/probe.json` | regions / tiers(tier1.enabled, prepaid_guard) / limits / ticket / paths |
| `probes/scenarios/` | tier-1 场景语料库 |
| `.claude/skills/tf-customer-probe` | 全流程技能 + references |
| `escalation/cap-tf-customer-probe.md` | 能力路线图 |

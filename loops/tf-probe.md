# tf-probe 合成客户探测循环

> 主动探测 terraform-provider-alicloud 潜在问题的 runbook。与 aone-triage 形成闭环：
> probe 产出的工单被 triage loop 扫到后按正常流程认领修复；修复发布后对应场景复跑即回归验证。
> 能力全景/路线图见 `escalation/cap-tf-customer-probe.md`；单场景全流程技能见 `.claude/skills/tf-customer-probe`。

---

## 一、触发

| 方式 | 说明 |
|------|------|
| 手动触发 | 用户在会话中执行 `/tf-probe` 或直接发指令「跑一轮场景探测」 |
| provider 新版 | 新版本发布后应尽快全量跑一轮（版本升级易引入 state 不兼容 / 新永久 diff） |
| cron / bridge | P2 才接入，本轮**不接** |

---

## 二、预检（doctor）

```bash
bootstrap/probe.sh doctor
```

- 全绿退 0 才继续。
- 缺 terraform → env 问题，**不硬闯**：按 loops/self-improve.md 记缺口 / escalation，本机只能跑 `--dry`。
- 缺凭证 → 只能 tier-0（plan 亦需凭证；无凭证时 plan 记 `no_creds` env_issue，不是 finding）。

---

## 三、场景选择

```bash
bootstrap/probe.sh list
```

- 默认**最久未跑优先**：看 `runs/probe/` 里各场景最新日期，挑最旧的。
- 单轮 ≤ `config.limits.max_scenarios_per_run`（默认 3）。
- provider 刚发新版 → 建议全量过一遍。

---

## 四、执行（调 tf-customer-probe 技能）

对每个选中的场景：

```bash
bootstrap/probe.sh run <id>          # 自然 tier,被配置封顶(tier1_enabled=false → 降 tier-0)
bootstrap/probe.sh run <id> --dry    # 只看步骤计划,不需 terraform
```

- 读产出 `verdict.json`（工作目录 + `runs/probe/<日期>-<id>.json` + 人读 `.md`）。
- 退出码：0 无 findings / 1 有 findings / 2 env 阻断 / **3 清理失败（最高优先级人工介入）**。
- **findings 判定是 Claude 的职责**（不是脚本）：逐 finding 对照 evidence 日志、`references/severity-rubric.md`、
  provider 仓 CHANGELOG Unreleased 段（已在 master 修掉的标「已修复未发布」不建单）。

---

## 五、产物分流

| 产物 | 去向 |
|------|------|
| `findings`（provider 疑似 bug） | 去重后 → draft（P0）到 `escalation/probe-drafts/`；毕业后 → adhoc-intake 建单（tf_provider 528766，见 skill） |
| `env_issues`（凭证/网络/allowlist/降级） | **不建单**；缺 terraform / 缺能力类走 loops/self-improve.md 或 escalation |

- 去重：a1 检索 528766 池 `jarvis-probe` 标签 + 标题关键词；GitHub 上游 open issues 只读检索；重复则追加 evidence 不新建。
- **probe 会话本身不 claim 工单**：产出的新单由 aone-triage loop 后续接管。

---

## 六、清理核查（sweep，残留即停）

```bash
bootstrap/probe.sh sweep
```

- 无残留退 0。
- **有残留退 1 → 立即停并升级**：`.my-day/probe/*/terraform.tfstate` 里有非空 `resources`，可能有计费资源没删干净，
  按提示手动 `terraform destroy` 或按 `managed_by=jarvis-probe` 标签用 aliyun CLI 清理。

---

## 七、Done — 本轮结束标准

- 每个跑过的场景都有一份 `verdict.json`（runs/probe/ 落盘）。
- `sweep` 零残留。
- 每条 provider finding 都已去重并落 draft（或标注「已修复未发布」/「上游已报」跳过）。
- 会话汇报：`场景数 N / 发现数 F / draft 数 D / env 问题数 E`，逐 draft 附路径与建议优先级。

---

## 八、与 aone-triage 的闭环

```
tf-probe 发现问题 → draft/建单(tf_provider 528766, tag jarvis-probe)
        ↓
aone-triage loop 扫到该单 → claim → provider-resource-dev/修复 → PR → 发布
        ↓
发布后:tf-probe 复跑对应场景 → verdict 无 finding = 回归通过(闭环)
```

真实客户工单也应回灌为 `regression-<aone-id>` 场景（probes/README.md），成为发版前回归项。

---

## 九、工具链速查

| 工具 | 作用 |
|------|------|
| `bootstrap/probe.sh doctor` | 环境预检（terraform/jq/凭证 set-unset/config 可解析） |
| `bootstrap/probe.sh list` | 列全部场景（id/tier/persona/resources/detect） |
| `bootstrap/probe.sh run <id> [--tier N] [--dry] [--keep]` | 跑单场景，落 verdict + 审计 |
| `bootstrap/probe.sh sweep` | 扫残留 state，残留退 1 |
| `config/probe.json` | tier 开关 / allowlist / limits / ticket / paths |
| `probes/scenarios/` | 场景语料库 |
| `.claude/skills/tf-customer-probe` | 单场景全流程技能 + references |
| `escalation/cap-tf-customer-probe.md` | 能力路线图 |

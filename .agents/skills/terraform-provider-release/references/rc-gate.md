# 发版前强化门禁（可选）— RC gate

> **Additive 指引,不改本 skill 既有步骤。** 这是 jarvis 侧 F3 探测线(`bootstrap/rc-gate.sh`)提供的
> 一条**发版前全量过闸**能力,供本 SOP 在提 PR / 合并前当可选闸门用。是否插入、插在哪,由跑者按发版风险自定;
> 门禁本身不阻塞 SOP 既有流程。

## 它做什么

把 tf-probe 的 tier-0（三方一致性全量机械层）+ tier-1（全场景真实生命周期）串成一次过闸,给红/黄/绿判定与退码:

```bash
bootstrap/rc-gate.sh <provider-dir>            # 完整模式:tier-1 真实 apply(顺跑,并发 1)
bootstrap/rc-gate.sh <provider-dir> --quick    # 快扫:tier-1 plan 为止(零创建),不 apply
```

- `<provider-dir>` = 本仓 worktree 里的 terraform-provider-alicloud 检出（需含 `website/docs/r` + `alicloud`）。
- 报告落 `runs/rc-gate/<date>-report.md`,顶部 `## VERDICT: <色> (exit N)`。

## 判定与退码（怎么用在本 SOP 里）

| 判定 | 退码 | 对本 SOP 的含义 |
|------|------|------|
| 🟢 GREEN | 0 | 过闸,照常进入远程 ACC / PR / 评审 / 人工合并环节。 |
| 🟡 YELLOW | 0 | 放行但**知情**:报告标注的黄项（doc_gap、机械层降级、queue 激增、quick 未 apply 等）非硬阻断,建议合并前尽量清;`--quick` 跑出的黄务必用完整模式复跑一遍再发。 |
| 🔴 RED | 1 | **禁发**:tier-0 `api_gap` S3+ / 场景 fail / destroy 残留。回到「修复方案」环节修红项,复跑门禁转绿/黄再提交。 |
| ⚪ CANNOT_CERTIFY | 2 | 门禁不完整（tier-0 跑不起来,provider 仓/probe-meta 环境缺）。修环境后重跑,**勿把不可判当绿放行**。 |

## 衔接位置（建议,可选）

```
本 SOP: 需求澄清 → gap 分析 → (生成/改码) → [可选: rc-gate.sh 过闸] → 远程 ACC → PR → 人工合并(release_prod 硬门)
```

- 门禁**不替代**本 SOP 的远程 ACC 验证——两者互补:ACC 深验单资源正确性,RC 门禁广验全量一致性 + 全场景生命周期。
- 门禁**不触碰** `release_prod`;合并仍是人工硬门(见本 skill「Jarvis 自治边界」)。
- 完整读法 / 触发时机 / env 调参见 jarvis 仓 `loops/tf-probe.md`「四点五、RC 门禁」节;实现见 `bootstrap/rc-gate.sh`。

# severity-rubric —— 危害分级与优先级映射

probe 发现的 finding 先落 detector 的 `severity_hint`（机判，见下表），再由 Claude 按本 rubric 复核、
升降级，最后映射到 Aone 优先级写进 draft/工单。

## 危害分级（S1 最重）

| 级别 | 定义 | 典型 finding |
|------|------|--------------|
| **S1 紧急** | 静默错配 / 数据或安全风险 / 清理失败 | 静默错配（`apply` 成功但云上实际 ≠ 声明）、意外替换重建（`unexpected_replace`）、数据丢失风险、安全缺省错误（如默认放通 0.0.0.0/0 未告警）、`destroy_fail` / `state_residue` |
| **S2 高** | 合法配置跑不通 / 收敛不了 | 合法配置 `apply_fail`、更新不生效、`perpetual_diff`、`import_diff`（import 断链） |
| **S3 中** | 文档示例级别的校验/映射错误 | 官方文档示例 `validate_fail` / `plan_fail`、schema 误报、报错信息严重误导 |
| **S4 低** | 体验问题 | 文档缺漏、字段说明不清、提示不友好、废弃字段未标注 |

## Aone 优先级映射

| severity | Aone 优先级 |
|----------|-------------|
| S1 | 紧急 |
| S2 | 高 |
| S3 | 中 |
| S4 | 低 |

> **枚举值待固化**：Aone 项目 528766 的优先级字段实际枚举值（及其 id），在**首次真实建单**（`ticket.mode=file`
> 毕业后）时用 `a1` 查证项目字段后固化到本表；draft 阶段先用中文标签「紧急/高/中/低」。

## detector 默认 severity_hint（runner 机判，Claude 可改）

### tier-1（真实 apply 生命周期）

| finding code | 默认 | 语义 |
|--------------|------|------|
| `validate_fail` | S3 | 官方示例组合 validate 不过 |
| `plan_fail` | S2 | 合法配置 plan 失败 |
| `plan_crash` | S1 | plan 触发 provider panic |
| `apply_fail` | S2 | 合法配置 apply 失败 |
| `perpetual_diff` | S2 | apply/更新后立即 plan 仍有 diff |
| `unexpected_replace` | S1 | apply 后 plan 出现 delete+create（意外重建） |
| `import_diff` | S2 | import 失败或 import 后 plan 非空 |
| `destroy_fail` | S1 | destroy 失败 |
| `state_residue` | S1 | destroy 后 state 仍残留资源 |

### tier-0（静态三方一致性扫描）

机械 diff(文档↔源码)产出:

| finding code | 默认 | 语义 |
|--------------|------|------|
| `doc_gap_phantom` | S3 | 文档有此参数但源码 schema 无(文档幻影参数) |
| `doc_gap_undocumented` | S3 | 源码有此参数但文档未记(未文档化参数) |
| `doc_gap_flag_mismatch` | S3 | Required/Optional 标注文档≠源码 |
| `doc_gap_forcenew` | S2 | ForceNew 标注文档≠源码(客户会被意外重建;仅查活跃非废弃字段) |
| `doc_gap_deprecated` | S4 | 源码已 Deprecated 但文档仍作正常参数列出(未标注废弃) |

OpenAPI 侧判定(skill 层查 `judgment_queue`,非机械)产出:

| finding code | 默认 | 语义 |
|--------------|------|------|
| `doc_api_gap` | S2–S4 | provider 实际调用的 API 的参数/枚举/行为与 TF 文档不一致,severity 由判定按影响面给。**只核对已接入面**;未接入 TF 的资源/参数不报(需求非 bug) |

## 升降级判则（Claude 复核时应用）

- **安全类资源升一级**：安全组/RAM/密钥/网络 ACL 等，其 `perpetual_diff` / 更新不生效影响安全边界 → S2 升 S1。
- **静默错配一律 S1**：只要“apply 报成功但云端实际配置与声明不符”，无论原始 code 是什么，升 S1（客户最难自查）。
- **仅文档示例问题降级**：若 `validate_fail` 根因是官方文档示例本身写错（而非 provider schema bug）→ 降 S4 文档类，
  建单走文档修正而非代码修复。
- **已修复未发布不建单**：finding 若对应 provider 仓 CHANGELOG `Unreleased` 段已合入的修复 → 不建单，
  在 draft 里标注「已在 master 修复，待 vX.Y.Z 发布」，转为发版后回归复跑项。
- **废弃字段的 flag/forcenew 差异不报**：tier-0 机械 diff 已把 `flag_mismatch`/`forcenew` 限定在活跃(非废弃)字段;
  废弃字段的标注差异是低价值噪声,只留 `doc_gap_deprecated`。
- **env_issue 永不升级为 finding**：鉴权/网络/`prepaid_block`/`tier1_disabled_plan_only` 都是环境噪声,不进危害分级、不建单。

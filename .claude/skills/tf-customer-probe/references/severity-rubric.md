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

> **枚举值待固化**：Aone 项目 528766 的优先级字段实际枚举值（及其 id），在**首次真实建单**（`ticket.mode=file`）时
> 用 `a1` 查证项目字段后固化到本表；`mode=draft` 时先用中文标签「紧急/高/中/低」。

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

OpenAPI 侧**机械三方 diff**(T0-mech,`probe.sh tier0` 元数据 diff 预筛)产出——精度命门,拿不准一律进 `judgment_queue` 不硬报:

| finding code | 默认 | 语义 |
|--------------|------|------|
| `api_gap_deprecated_action` | S3 | 源码在调的 (product,version,action) 在 OpenAPI 标 `deprecated:true`(如 ClassicLink 退役),TF 侧无提示 |
| `api_gap_enum_superset` | S3 | TF `ValidateFunc` 枚举 ⊋ API 枚举——客户端放行 API 必拒的值(如 OSS `storage_class` 含 `Standard`) |
| `api_gap_required` | S3 | TF 标 `Optional` 但 API `docRequired=true`(且无 `Default`/`Computed` 兜底、非废弃双轨) |
| `api_gap_type` | S3 | 类型硬冲突(经 `type_tolerance` 容差表过滤后仍冲突,如 TF list vs API string) |
| `api_gap_range` | S4 | TF `IntBetween(min,max)` 越过 API 数值范围(TF 放行 API 拒绝的范围) |
| `api_gap_default` | S4 | TF 显式 `Default` 与 API 默认值冲突 |

> 方向安全不报:TF 枚举 ⊊ API 枚举 / TF 范围更严 = TF 比 API 严,记 `coverage_notes[]` 而非 finding。
> 抑制护栏:`config.tier0_mech.suppress_params`(ClientToken/RegionId 等)与 `type_tolerance`(API int↔TF string 建模惯例)命中记 `suppressed[]`(可审计),不报 finding。废弃双轨对沿用 deprecated 检测(参与双轨的参数 required/名称差异不报)。

OpenAPI 侧**剩余非机械项**(机械层拿不准 → `judgment_queue` 带 reason,交 skill 层/verifier 双层查证)产出:

| finding code | 默认 | 语义 |
|--------------|------|------|
| `doc_api_gap` | S2–S4 | queue 里 `prose_review`(长度/字符集/基数等纯 prose 约束、行为一致性)/`unmapped_params`(snake→Camel 映射不上,如 convert 改名)/`enum_unparsed`(枚举非字面 slice)/`no_triple`(OSS SDK 风格抽不到 action)经人工核实后确为不一致。severity 由判定按影响面给。**只核对已接入面**;未接入 TF 的资源/参数不报(需求非 bug) |

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

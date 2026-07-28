# Templates & verified facts

## 同步 provider 源码(查证前)
`bash .Codex/skills/aone-triage/scripts/sync-provider.sh` —— 无库 clone,有库 fetch + `reset --hard` 强制对齐 upstream(**主目录会被重置,开发一律走 worktree**);repo 路径走 `bootstrap/workspace.sh dir terraform_provider`(本机覆盖 `workspaces.local.json` / `JARVIS_WORKSPACE_ROOT`)。

## 缺陷骨架
复现要点 → `$(bash bootstrap/workspace.sh dir terraform_provider)` 源码定位(资源.go + 行号)→ 根因 → 修复/绕过 → **补/改一个会因该 bug 失败的用例锁定回归(无可测则在 CR 说明为何)** → 仅 spec 缺口才转需求。

## 任务骨架
拆解步骤 → 直接执行 → 产出+回执;无需查证/转需求。

## Comment skeleton (回写工单)
```
【经查证回复 - <provider 源码 + OpenAPI 双向确认>】

核心结论:<一句话差异定位>

一、<问题1>:<是否支持> + 证据(provider 源码行 / grep 命中数 / 现有可配置项)
二、<问题2>:import / 数据源 vs 资源(单复数)说明 + 命令示例
三、解决方案:1) 短期绕过(OpenAPI/CLI 或 null_resource)2) 正解(给 provider+spec 提 PR)
```

## Canned 补料等待骨架

命中 canned 且正式判断所需材料不完整时，先由唯一写者回复并**等待补料**：

```markdown
### 当前结论
现有材料不足以判断 <API / CloudSpec / Provider> 的责任层，当前不进入正式路由或开发。

### 请补充
- <完整 HCL / 资源和字段>
- <完整错误、API error code；RequestId 仅在内部证据中脱敏保存>
- <期望行为与当前行为>
- <必要的 Terraform debug 日志或最小复现>

### 已完成的安全检查
<只读安全查证事实；明确这些事实不构成正式路由结论>

### 下一步
收到材料后重新执行完整决策树；本轮 release/idle 等待补料。
```

## Terraform 文档分层判定

| 证据结论 | 路由 |
|---|---|
| **pure datasource**：只涉及 `data.alicloud_xxx` 查询、过滤、分页、输出字段或 Read，不含 resource 变更 | **source-only**：紧急源单新山、非紧急源单过载，由 TerraformRD 在源单直接开发；历史 relation 只读，严禁任何 528766 承载动作 |
| **CloudSpec 文档文本 metadata**：只改 resource/property/operation description、字段解释、NOTE 或枚举文案，不改变字段集合、类型、约束或 CRUD | **分支 I**，创建或复用 2169561 并指派念依（373108）；Provider 公开 docs 同时错误时独立补 528766 紧急兜底腿 |
| **CloudSpec 文档源正确，Provider 本地文档生成/展示偏差** | **分支 D**，仅处理 Provider 本地生成、发布或展示链路 |
| **CloudSpec 结构 metadata**：新增/删除字段，或改变类型、约束、枚举集合、CRUD/operation/映射 | **分支 E**，在原主单修到 pre Meta 收敛，再强制 E → D-临钧 |
| 尚未取得足以判断文档源和结构边界的证据 | 使用上方 canned 骨架等待补料，不得猜测 I/D/E |

## 分支 I · 文档文本 metadata 路由骨架

I 的 text-only 文档主腿固定走 `upstream.cloudspec_docs_quality`；PD/QA 不外写，
terraform-rd finalizer 作为 downstream `single-writer` 执行，executor 只负责原主单
bookend，不解析或重放路由动作：

```markdown
### CloudSpec 文档文本 metadata

- 边界：<resource/property/operation description、字段解释、NOTE、枚举文案>
- 结构不变证据：<字段集合、类型、约束、CRUD 均未变化>
- 2169561：<created/reused、念依（373108）、relation>
- Provider 公开 docs：<正确 / 同时错误>
- 独立 528766 紧急兜底腿：<N/A / created/reused、过载（484483）、relation>
- 分池防重：<2169561 point-read>；<528766 point-read>
```

一个池已有 relation 不能抑制另一个池的缺失补建；528766 只能临时兜底公开 docs，不能替代
2169561 文档源主腿。

## 分支 E · CloudSpec 结构 metadata 原主单自闭环骨架

E 仅处理字段集合、类型、约束、CRUD/operation/映射。open-jarvis 在当前原主单修 CloudSpec
到 pre Meta 收敛，然后交 D-临钧；不得直接做 Provider PR/CI/ACC：

```markdown
### CloudSpec 结构 metadata 原主单自闭环

- 缺口：<字段集合/类型/约束/CRUD/operation/映射 + 预期语义>
- OpenAPI 证据：<Product::Action、字段、类型、枚举>
- 初始 pre Meta：<资源是否存在、属性/CRUD/文档差异>
- AMP/Git：<task 专属 feature 分支、AMP 返回的 SSH URL、cloudspec-model commit/MR>
- 编辑链：cloudspec-amp-workflow → cloudspec-idl-guide →
  cloudspec-resource-edit / cloudspec-operation-edit → cloudspec-build-fix →
  cloudspec-norm-check-fix
- 验证：`aliyun cspec build`、资源级 `aliyun cspec check`、`amp publish pre --dry-run`
  与 `amp publish pre`、pre Meta 收敛结果
- E → D-临钧：<已有 relation/taskId/aoneId 复用，或 createBuildTaskV2 的 taskId/aoneId>
- Provider：<由 D-临钧生成器链路后续处理；E 不直接执行 PR/CI/ACC>
- 当前门：<pre 未收敛，不触发 Acube / handoff 完成 / missing_capability / blocked>
- 下一步：<handoff 回执后 release/idle；不得 finish；prod/online 与主干仍为人工硬门>
```

权限、AMP 登录、SSH 或仓库访问失败时，保留已取得的只读证据并返回
`missing_capability` / `blocked`；不得换个人身份、改派旧 2165097 路径或另建 Aone
规避能力缺口。pre 未收敛不得触发 Acube；不得在 E 完成后直接 release/idle。

## Requirement skeleton (Terraform 生成器问题/API 工具团队)
默认池: `api_toolkit` / project `2100304`;产品字段优先选 Terraform;标题聚焦生成器行为,不要写成客户资源诉求。
```
background: Terraform 资源 <terraform_resource> 研发中发现 Acube/terraform-generator-v4 生成行为与 Cloudspec resourceTypeCode 配置不一致,影响资源生成验收。
requirement: 修复 terraform-generator-v4 对 <Cloudspec字段/条件> 的生成逻辑。现状:<实际生成片段/行为>;期望:<应生成片段/行为>。
evidence: 1. resourceTypeCode/get: <接口/文件路径>,关键配置=<字段和值>. 2. createLocalBuildTask: <任务/文件路径>,生成代码=<文件:行>. 3. 关联资源/工单:<source>。
acceptance: 1. 给定相同 resourceTypeCode,Acube 重新生成代码符合 Cloudspec 语义。2. 补生成器单测覆盖该条件。3. 关联 Terraform 资源可通过 `tools/terraform_generated_diff.py --acube-dir ...` 语义检查,不再报 WARN。
deadline: <date>
source: <源 Aone/PR/CR url>

```json
{"background":"…","requirement":"…","evidence":"…","acceptance":"…","deadline":"YYYY-MM-DD","source":"…"}
```
```

## Verified ROS facts (查证基准 provider 1.281.0)
- `alicloud_ros_stack_group` → `ROS::StackGroup`;无 `rd_folder_ids`,变更部署支持。
- `alicloud_ros_stack_instance` → `ROS::StackInstance`;支持 import `<group>:<account>:<region>`;Create 仅发 AccountIds/RegionIds,无 DeploymentTargets。
- `alicloud_ros_stack_instances`(复数)= 数据源,不可 import。
- OpenAPI `CreateStackInstances`/`UpdateStackInstances` `DeploymentTargets` 子字段 = `RdFolderIds`,`AccountIds`。
- 优先级 id: 94紧急 95高 96中 97低;关系 20001关联 20013阻塞 20014被阻塞。

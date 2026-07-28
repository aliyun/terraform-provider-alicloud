# Templates & verified facts

## 同步 provider 源码(查证前)
`bash .claude/skills/aone-triage/scripts/sync-provider.sh` —— 无库 clone,有库 fetch + `reset --hard` 强制对齐 upstream(**主目录会被重置,开发一律走 worktree**);repo 路径走 `bootstrap/workspace.sh dir terraform_provider`(本机覆盖 `workspaces.local.json` / `JARVIS_WORKSPACE_ROOT`)。

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
| **CloudSpec 文档源错误**：resource/property/operation description 或枚举文案与 OpenAPI 长期语义不一致 | **分支 E**，在 CloudSpec OK 判定前进入原主单自闭环；schema/coverage 全绿不能覆盖该结论 |
| **CloudSpec 文档源正确，Provider 本地文档生成/展示偏差** | **分支 D**，仅处理 Provider 本地生成、发布或展示链路 |
| 尚未取得足以判断文档源的证据 | 使用上方 canned 骨架等待补料，不得猜测 E/D |

## CloudSpec 原主单自闭环记录骨架

CloudSpec 资源定义、metadata 或资源文档源头缺口不再创建任何镇元侧/文档兜底关联单。由
open-jarvis 在当前原主单内完成，并把下列证据交最终 RD 聚合：

```markdown
### CloudSpec 原主单自闭环

- 缺口：<资源/属性/操作/文档位置 + 预期语义>
- OpenAPI 证据：<Product::Action、字段、类型、枚举>
- 初始 pre Meta：<资源是否存在、属性/CRUD/文档差异>
- AMP/Git：<task 专属 feature 分支、AMP 返回的 SSH URL、cloudspec-model commit/MR>
- 编辑链：cloudspec-amp-workflow → cloudspec-idl-guide →
  cloudspec-resource-edit / cloudspec-operation-edit → cloudspec-build-fix →
  cloudspec-norm-check-fix
- 验证：`aliyun cspec build`、资源级 `aliyun cspec check`、`amp publish pre --dry-run`
  与 `amp publish pre`、pre Meta 收敛结果
- Provider：<是否需从 pre 重新生成、diff、PR、CI、远程 ACC>
- 当前门：<pre 已完成 / missing_capability / blocked / 待 prod/online 或主干人工动作>
- 下一步：<release/idle；不得 finish，除非正式发布与其它硬门另行完成>
```

权限、AMP 登录、SSH 或仓库访问失败时，保留已取得的只读证据并返回
`missing_capability` / `blocked`；不得换个人身份、改派外部承接人或另建 Aone 规避能力缺口。

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

# Templates & verified facts

## 同步 provider 源码(查证前)
`bash .Codex/skills/aone-triage/scripts/sync-provider.sh` —— 无库 clone,有库只 fetch 不重置(保护本地开发);repo 路径取 config/workspaces.json 的 terraform_provider.path（改这里即可换路径）。

## 缺陷骨架
复现要点 → `~/go/src/github.com/chenhanzhang/terraform-provider-alicloud` 源码定位(资源.go + 行号)→ 根因 → 修复/绕过 → **补/改一个会因该 bug 失败的用例锁定回归(无可测则在 CR 说明为何)** → 仅 spec 缺口才转需求。

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

## Requirement skeleton (Cloudspec 关联单 · 镇元 agent 接单硬契约)

**这是硬契约,不是可选参考**。落 **Terraform镇元对接(2165097)** 池的关联单(包括 tf_customer 分支 E 谜拟单、其它域上游 Cloudspec 缺口单)交由 **镇元 agent(`WORKER_1783326253279`)** 自动接手,agent 从 description 里解析结构化字段驱动后续处理(spec 补齐 / 映射建立 / 生成器触发)。**字段缺失 / marker 缺失 / JSON 语法错 = agent 无法接单 = 单沉底**——jarvis 自己不能替 agent 干活,漏契约就是把单丢黑洞。

**契约来源**:谜拟本人钉钉文档 https://alidocs.dingtalk.com/i/nodes/YQBnd5ExVEjea40qC2vjQyxPJyeZqMmz

### 字段清单(全部必填,无可选)

| 字段 | 放什么 | 示例 |
|---|---|---|
| `background` | 原文【背景】段落:客户诉求 + 现状(为什么现在做不到) | "OSS 已支持通过 OpenAPI 配置桶清单规则,但 Terraform Provider 中不存在 alicloud_oss_bucket_inventory 资源,客户无法通过 Terraform 管理 OSS 桶清单规则。" |
| `requirement` | 原文【需求】段落整体:要什么资源/属性/接口,一条一条列 | "新增 alicloud_oss_bucket_inventory 资源,支持桶清单规则的完整 CRUD 操作。" |
| `documentUrl` | 官方 API 文档链接(api.aliyun.com/document/) | "https://api.aliyun.com/document/Oss/2019-05-17/PutBucketInventory" |
| `mappingCheckUrl` | acube getTerraformResourceSpec 映射查询接口 | "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_oss_bucket_inventory" |
| `acceptance` | 验收标准:spec 补齐 / 映射建立 / 查询返回有效 / 偏差检测可见 | "Cloudspec 资源 spec 补齐,alicloud_oss_bucket_inventory ↔ Cloudspec 映射建立,getTerraformResourceSpec 查询返回有效映射。" |
| `deadline` | 交付时间(YYYY-MM-DD;无明确 DDL 传空串,别省字段) | "2026-06-30" 或 "" |
| `source` | 需求来源(源工单号 + 类别 或 咨询上下文) | "工单 82845574（Terraform - 客户问题）" |

### description 完整正文骨架(必须严格照抄结构)

````markdown
## 背景
<原文【背景】段落 · 与 JSON background 字段一致>

## 需求
<原文【需求】段落 · 与 JSON requirement 字段一致>

## 机读信息
```json
{
  "background": "<与上方【背景】段落一致>",
  "requirement": "<与上方【需求】段落一致>",
  "documentUrl": "https://api.aliyun.com/document/<Product>/<ver>/<Action>",
  "mappingCheckUrl": "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x",
  "acceptance": "<spec补齐 + 映射建立 + getTerraformResourceSpec 有效 + 偏差检测可见>",
  "deadline": "YYYY-MM-DD 或 空串",
  "source": "工单 <ID>（Terraform - 客户问题）"
}
```
````

### 硬性规则

1. **`## 机读信息` marker 段名不可省不可改**——agent 靠这个 marker 定位 JSON 块;写成 "机读信息"、"机读 JSON"、`### 机读信息` 都不认。
2. **JSON 必须在 ` ```json ` 代码块内**——纯文本 JSON 或用 ` ``` ` (无语言标记) 也可能不认。
3. **7 字段全部必填,顺序建议同上**;`deadline` 无明确日期时传空串 `""`,别删字段。
4. **JSON 字段与上方【背景】【需求】段落必须逐字一致**——两处不一致 agent 有可能取错。
5. **需求导向,不罗列现有属性清单**——写"缺什么/对应哪个API",不写"现有 x/y/z 属性"。
6. **`source` 必须能反查源客户工单**——形如 `工单 <ID>（Terraform - 客户问题）` 或 `用户咨询 <API url>`;别只写"客户需求"。

### 参考样本

样本工单 [83120127](https://project.aone.alibaba-inc.com/v2/project/2165097/req/83120127) 是标准姿势(谜拟本人给的示范单),需要抄骨架时对着它抄。

### 反模式速查

- ❌ description 只有【背景】【需求】两段,没有 `## 机读信息` + JSON → agent 完全不接单,单沉底
- ❌ 只放 JSON、没有【背景】【需求】自然语言段落 → 人类看不懂,评审无从下手
- ❌ JSON 字段名拼错(如 `docUrl` / `mappingUrl` / `dueDate`)→ agent 解析错,把 default 值当输入
- ❌ 引用 <mcreference>...</mcreference> 之类的模型标记留在 background/requirement 里 → agent 语义污染,建议先剥掉
- ❌ `documentUrl` 用 help.aliyun.com/document_detail/xxx.html → 应用 api.aliyun.com/document/<Product>/<ver>/<Action>(agent 后续拉 OpenAPI spec 用)
- ❌ `mappingCheckUrl` 里 `terraformResourceType` 参数值和 requirement 里资源名不一致 → agent 查错映射
- ❌ 主观倾向/裁决建议塞进 requirement → 只放客户诉求本身,方案由 agent 或后手研发决定

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

## Verified ROS facts (本次确认, provider 1.281.0)
- `alicloud_ros_stack_group` → `ROS::StackGroup`;无 `rd_folder_ids`,变更部署支持。
- `alicloud_ros_stack_instance` → `ROS::StackInstance`;支持 import `<group>:<account>:<region>`;Create 仅发 AccountIds/RegionIds,无 DeploymentTargets。
- `alicloud_ros_stack_instances`(复数)= 数据源,不可 import。
- OpenAPI `CreateStackInstances`/`UpdateStackInstances` `DeploymentTargets` 子字段 = `RdFolderIds`,`AccountIds`。
- 优先级 id: 94紧急 95高 96中 97低;关系 20001关联 20013阻塞 20014被阻塞。

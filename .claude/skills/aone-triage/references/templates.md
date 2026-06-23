# Templates & verified facts

## 同步 provider 源码(查证前)
`bash .claude/skills/aone-triage/scripts/sync-provider.sh` —— 无库 clone,有库只 fetch 不重置(保护本地开发);repo=`~/go/src/github.com/chenhanzhang/terraform-provider-alicloud`,可用 `TERRAFLOW_ALICLOUD` 覆盖。

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

## Requirement skeleton (Cloudspec 需求池标准格式)
7 字段正文 + 末尾机读 JSON。标准来源:钉钉文档 https://alidocs.dingtalk.com/i/nodes/YQBnd5ExVEjea40qC2vjQyxPJyeZqMmz
```
background: 缺口一句话 + 来源诉求
requirement: 1. 资源:<Namespace::Type>(alicloud_x) 缺<attr> 对应<CreateX> 预期透传+import. 2. …
documentUrl: https://api.aliyun.com/document/<Product>/<ver>/<Action>
mappingCheckUrl: https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x
acceptance: <spec补齐+映射建立+getTerraformResourceSpec有效+偏差检测可见>
deadline: <date>
source: <工单 url 或 用户咨询>

```json
{"background":"…","requirement":"…","documentUrl":"…","mappingCheckUrl":"…","acceptance":"…","deadline":"YYYY-MM-DD","source":"…"}
```
```
保持需求导向:写"缺什么/对应哪个API",不罗列现有属性清单。两层正文 + JSON 字段必须一致。

## Verified ROS facts (本次确认, provider 1.281.0)
- `alicloud_ros_stack_group` → `ROS::StackGroup`;无 `rd_folder_ids`,变更部署支持。
- `alicloud_ros_stack_instance` → `ROS::StackInstance`;支持 import `<group>:<account>:<region>`;Create 仅发 AccountIds/RegionIds,无 DeploymentTargets。
- `alicloud_ros_stack_instances`(复数)= 数据源,不可 import。
- OpenAPI `CreateStackInstances`/`UpdateStackInstances` `DeploymentTargets` 子字段 = `RdFolderIds`,`AccountIds`。
- 优先级 id: 94紧急 95高 96中 97低;关系 20001关联 20013阻塞 20014被阻塞。

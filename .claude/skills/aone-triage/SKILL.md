---
name: aone-triage
description: >-
  Triage and resolve Aone work items of any source (小蜜/customer ticket, 缺陷/bug, 任务/task,
  需求/req) and answer standalone "does Terraform/Cloudspec support X" questions. Use whenever the
  user shares an Aone work item URL (project.aone.alibaba-inc.com), a 工单ID, says "看下/处理/回复
  这个工作项", pastes a next.api.aliyun.com product/API link, or asks whether an alicloud provider
  resource / Cloudspec resource supports some attribute or whether a product is integrated into
  Terraform. Covers: read ticket → triage → verify (OpenAPI + Cloudspec mapping + provider SOURCE) →
  reply → escalate gap to Cloudspec 需求池 → link → set status. 小蜜 is just one source. Trigger even
  with no ticket — a bare API/product link is enough. Also covers the deliver-it-yourself path: 建需求
  → 建变更/CR → worktree 开发 → 预发 → 正式 via a1 CLI, when the fix lands in one of our own apps —
  trigger when 提需求给自己 / 给 Agent门户 / AgentRuntime / aliyun-automation-agent / PlayGround 提需求,
  or cwd is the aliyun-automation-agent repo (routes to project 2124589 + app 283346,见 config/pools.json).
---

# Aone 工单 / Terraform 能力查证

> **红线(交付改动):任何代码改动一律新分支 + CR/MR 评审,master 只接已评审合入。禁止直接 push master,禁止从 master 拉零 diff 空 CR 直发正式——没有可评审 diff 就不上线。**

## 前置:依赖 a1 CLI
全流程(工单读写、建需求、CR、发布)都走 `a1`。参考 <https://open.aone.alibaba-inc.com/skill/a1>。
先 `a1 auth whoami` 验登录;命令报 `command not found` 或认证错误时,**征得用户同意后**装:
```
curl -fsSL https://git.cn-hangzhou.oss-cdn.aliyun-inc.com/aone-cli/install.sh | sh
a1 skill install a1@0.28.0
```
未同意则停下提示手动安装,不要绕过用 MCP。

Repeatable support-engineering loop. Two ways in: a work item (any source/type) or a bare
"does Terraform/Cloudspec support X" question. Side-effectful steps (comment, create, assign,
status) need a clear user yes first — confirm before each.

## 入口分诊
- **有 workitemId**(贴工单链接/ID)→ 工单全流程(下方)
- **只有 API/产品**(next.api.aliyun.com 链接、"alicloud 支持 X 吗"、"把 X 接入 Terraform")→ 直接走「查证」,输出"已支持 / 缺哪些属性",无需建工单
- 两者都有 → 工单流程内嵌查证

## 选池 + 是否开发的闸门
建需求前查 `references/routing.md` 选池;命不中先反问。建完按内容分叉:
- 分析/调研/咨询/统计 → 建需求即停,不拉分支
- 改代码/接资源/修 bug → 仅命中开发行才进链路(link app→CR→worktree)

## 工单全流程
`a1 project workitem get <id>`(URL 末尾数字,`-f json` 取字段)→ 概要表(id/标题/类型/状态/指派/工单ID)+ 诉求 → 按类型分诊:
1. **需求/咨询** → 查证 → 回复(授权)→ 真缺口转 Cloudspec 池 + 双向关联 + 源工单改"待上游排期"
2. **缺陷** → 复现要点 + 源码定位 → 回复/指派,确认 spec 缺口才转需求
3. **任务** → 直接执行或拆解
4. **自己交付**(改我们自家应用,不转上游)→ 走对应应用的交付链路:建需求→建 CR→worktree 开发→预发→**等用户验证反馈**→正式→关单清 worktree。先按目标应用选链路文件,IDs/坑见 reference:
   - **Agent门户 / AgentRuntime / aliyun-automation-agent / PlayGround**(或 cwd 在该 repo),给自己提需求 → `references/delivery-aliyun-automation-agent.md`(默认项目 2124589 + app 283346)
   - 其余应用追加同名 `delivery-<app>.md`

## 查证(两层,顺序固定,不凭记忆)
1. **OpenAPI 全集**:解析 product+action(next.api 链接或描述)→ `AlibabaCloud ListApis` / `GetApiDefinition`。JMESPath 用单引号,反引号会失败:`parameters[?name=='X'].schema.properties|[0]|keys(@)`。
2. **映射**:`curl "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x"` —— **仅判断 TF 资源 ↔ Cloudspec 资源是否已建映射,不代表实现**。
3. **实现以源码为准**:先 `scripts/sync-provider.sh` 同步,再在 `~/terraflow/providers/alicloud` grep 资源 .go,核对 schema 字段 / `Importer` / Create 实际下发参数。单复数陷阱:`*_instances` 多半是数据源。
4. **文档兜底**:GitHub raw markdown。

## 回复/转需求/关联(写操作,先授权)
- 回复:草稿过目 → `a1 project workitem comment create <id> -m "..."`。结构=结论→逐问+证据→方案。
- 转需求:Cloudspec 池 `2165097`,指派谜拟 `479782`(见 config/pools.json) → `a1 project workitem create --project 2165097 --category req --assignee 479782`。描述用 7 字段标准格式(background/requirement/documentUrl/mappingCheckUrl/acceptance/deadline/source)+ 末尾机读 ```json,字段一致;不堆现状。模板见 templates.md。
- 关联:`a1 project workitem relation add <id> relate:<target>` 调两次(A→B、B→A)双向。
- 源工单改状态"待上游排期":`a1 project workitem update <id> --status "待上游排期"`(仅真缺口已转池时)。
- 创建报必填字段缺失 → `a1 project workitem field options <field> --project <id>` 查枚举补 `--cfs`(同自家应用建需求那套)。

骨架与 ROS 事实见 `references/templates.md`。

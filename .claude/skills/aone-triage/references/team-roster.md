# 团队分工速查(TF Provider 路由)

> 本 reference 是 tf_customer 路由的团队分工单点维护(P3.a 从 tf-customer-request-routing 抽出),便于跨 skill/loop 复用。

### 专属维护名单(直接指派该负责人,不走镇元/生成器)

这批云产品的 provider 代码由专人维护,**不接镇元**;不用查覆盖度,直接指派 + @。

| 云产品 | 花名 | 工号 |
|---|---|---|
| 容器服务 Kubernetes (ACK) | 若即 | 377376 |
| 日志服务 SLS | 豁朗 | 269032 |
| 消息服务 MNS | 曼红 | 38570 |
| OSS | 凡修 | 71145 |
| 弹性伸缩 ESS | 扶柳 | WB530580 |
| 表格存储 OTS | 景哲 | 263417 |
| E-MapReduce (EMR) | 鱼戏 | 373227 |
| RDS | 柴天生 | WB01586841 |
| PolarDB | 米汐 | 527630 |
| MSE | 棠溪 | 401341 |
| ClickHouse | 逸颉 | 439859 |

### 通用路由角色(其他云产品走此表)

| 场景 | 花名/名称 | 工号 |
|---|---|---|
| **Provider 侧全局改造**(非单一产品/资源:region 白名单/框架 utility/公共 endpoint/provider.go 基础/SDK bump) | 新山 | 521957 |
| **镇元资源文档修改**(仅描述/字段解释/枚举值文案层,不涉资源本身新增字段/结构)——**分支 I 关联单落 CloudSpec 文档质量问题(2169561) 池**指派念依修镇元源头。TF provider docs 从镇元自动生成,provider PR 只是紧急兜底,不修镇元源头下次发版会覆盖回旧值。**文档改造分支通常与 528766 过载单双建**(过载紧急合 provider PR,念依修镇元源头) | 念依(陈旖旎) | 373108 |
| **与镇元相关且镇元 NOT OK · 关联单**(资源未定义 / 属性不满足诉求 / 覆盖度<100%)——镇元侧根因主责,**无论紧急与否都建单**;关联单落 2165097 池指派**镇元 agent 自动接单**(2026-07-08 谜拟本人切换,不再手动解单)。**body 必须严格按 [templates.md 硬契约](./templates.md) 骨架**(`## 背景` / `## 需求` / `## 机读信息` + ```json 代码块 + 7 字段全),缺 marker/字段/JSON 语法错 = agent 无法接单 = 单沉底。**注意 agent 只接"资源本身需变更"**——纯文档描述修改(枚举值文案/字段说明)走念依 · 2169561 池,不投 agent | 镇元 agent(无花名,agent 身份) | WORKER_1783326253279 |
| 同上 · 源客户主单 assignee(客户可见,agent 断电/复杂决策时人类兜底 owner) | 谜拟 | 479782 |
| 同上且**紧急**(优先级=紧急 或 距 DDL<14 天 或 缺陷类型覆写):在镇元 agent 单之外**再建一张关联单并行**(紧急兜底 provider 侧,agent+新山两张都建) | 新山 | 521957 |
| 与镇元不相关 + 资源代码由生成器产出(修复=acube 重跑生成器,管道不变) | 临钧 | 429768 |
| **与镇元不相关**(纯 datasource / 镇元 OK 但 provider 侧问题,手写代码)且**紧急** | 新山 | 521957 |
| 与镇元不相关(同上)且**不紧急**(默认兜底) | 过载 | 484483 |
| **NPE 兜底**(以上所有分支均未匹配 / 跨多产品无单一负责人 / 分诊模糊超出团队职责)+ 打标签 `jarvis-npe` | 夏节 | 401498 |

**镇元 agent 特点**(2026-07-08 上线):
- 工号 `WORKER_1783326253279`,是 agent/机器人身份(非自然人),**无钉钉 IM 通道**(notify-dingtalk.sh 传 `WORKER_` 前缀会 400/静默);私信必须发谜拟(479782)/新山(521957) 等真人
- Bridge 视角仍算"人工介入"(见 config/contacts.json 说明),不会触发 jarvis-idle 自我重派
- **只识机读 JSON**,不看自然语言 body 补充说明——契约见 `references/templates.md` "Requirement skeleton (Cloudspec 关联单 · 镇元 agent 接单硬契约)"
- 只接分支 E 场景(与镇元相关且镇元 NOT OK);其它场景不能强指派

**与镇元不相关的问题**(2026-07-06 定义,镇元 agent 不接这两类):
1. **纯 datasource 问题**:诉求只涉 `data.alicloud_xxx`(查询/过滤/输出字段),不涉资源 schema/生命周期——datasource 是 provider 侧对查询 API 的只读封装,镇元只管资源 schema,**跳过镇元查证**;resource+datasource 混合不算"纯"
2. **镇元侧无问题、provider 侧存在问题**:镇元 OK 三条件全满足,缺口在 provider 实现(bug/适配缺失/文档行为不符)

**镇元 OK 三条件**(全满足才算 OK,任一不满足即视为 NOT OK=与镇元相关):
1. **API 在镇元有对应资源**:资源已在镇元定义并发布(get 返回 data 且 released list 命中)
2. **当前资源属性满足客户诉求**:比对客户抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖**(缺字段=NOT OK,即便覆盖度分再高也不算 OK)
3. **测试覆盖度 100%**:acube V2 `CoverageDetail.CoverageScore == 1.0`

详见 Step 2 分支 D 前的判定说明。

### 上游 API 缺口(纯上游产品团队问题)

- 不走上述任一路由,**@提单人**(工单 creator)+ status=待上游排期
- 由提单人协助转对应云产品的 API/OpenAPI 团队评估;jarvis 不代跨团队协调

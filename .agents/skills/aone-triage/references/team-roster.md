# 团队分工速查(TF Provider 路由)

> 本 reference 是 tf_customer 路由的团队分工单点维护,便于跨 skill/loop 复用。

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
| **CloudSpec 文档文本 metadata（I）**：resource/property/operation description、字段解释、NOTE、枚举文案，且不改变字段集合/类型/约束/CRUD | 念依（2169561 submit_only） | 373108 |
| **CloudSpec 结构 metadata（E）**：资源未定义、字段集合/类型/约束/CRUD/映射不满足诉求；原主单修到 pre Meta 收敛后强制 **E → D-临钧** | open-jarvis → 临钧 | 原主单内部执行 → 429768 |
| 与镇元不相关 + 资源代码由生成器产出(修复=acube 重跑生成器,管道不变) | 临钧 | 429768 |
| **与镇元不相关**(纯 datasource / 镇元 OK 但 provider 侧问题；文档场景仅限 CloudSpec 源正确、Provider 本地生成/展示偏差)且**紧急** | 新山 | 521957 |
| 与镇元不相关(同上)且**不紧急**(默认兜底) | 过载 | 484483 |
| **NPE 兜底**(以上所有分支均未匹配 / 跨多产品无单一负责人 / 分诊模糊超出团队职责)+ 打标签 `jarvis-npe` | 夏节 | 401498 |

**I/E 路由契约**:
- I 创建或复用 2169561 并指派念依；Provider 公开 docs 同时错误时独立补 528766 紧急兜底腿。
  两池分池防重，一个池已有 relation 不能压掉另一个池的缺失补建。
- E 的 PD 返回 `requested_external_actions: []`、`next=terraform-rd/dev`；RD 使用
  `cloudspec-amp-workflow` 与 IDL/resource/operation/build/norm skills修到 pre Meta 收敛。
- E pre 未收敛不得触发 Acube；收敛后必须 E → D-临钧，经 `createBuildTaskV2` 创建或复用
  528766 并指派临钧（429768）。E 不直接做 Provider PR/CI/ACC，也不直接 release/idle。
- PD/QA 不外写；terraform-rd finalizer 是 downstream single-writer，负责 I/E 路由动作与
  下游副作用；executor 只负责原主单 bookend。AMP 登录、SSH 或仓库
  权限失败返回 `missing_capability` / `blocked`；不得回退 2165097。
- prod/online、master/main merge/push 与正式发布仍是人工硬门。

**与镇元不相关的问题**(不进入 CloudSpec 原主单自闭环的两类):
1. **纯 datasource 问题**:诉求只涉 `data.alicloud_xxx`(查询/过滤/输出字段),不涉资源 schema/生命周期——datasource 是 provider 侧对查询 API 的只读封装,镇元只管资源 schema,**跳过镇元查证**;resource+datasource 混合不算"纯"
2. **镇元侧无问题、provider 侧存在问题**:CloudSpec 结构 OK 三条件全满足,缺口在 provider 实现；
   文档场景仅包含 CloudSpec 文档源正确、Provider 本地文档生成/展示偏差

**CloudSpec 结构 OK 三条件**(全满足才算 OK,任一不满足即视为结构 NOT OK):
1. **API 在镇元有对应资源**:资源已在镇元定义并发布(get 返回 data 且 released list 命中)
2. **当前资源属性满足客户诉求**:比对客户抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖**(缺字段=NOT OK,即便覆盖度分再高也不算 OK)
3. **测试覆盖度 100%**:acube V2 `CoverageDetail.CoverageScore == 1.0`

text-only 文档 metadata 先按 I 判定，不混入结构 OK；文档与结构同时变化时按 E 处理。

详见 Step 2 分支 D 前的判定说明。

### 上游 API 缺口(纯上游产品团队问题)

- 不走上述任一路由,**@提单人**(工单 creator)+ status=待上游排期
- 由提单人协助转对应云产品的 API/OpenAPI 团队评估;jarvis 不代跨团队协调

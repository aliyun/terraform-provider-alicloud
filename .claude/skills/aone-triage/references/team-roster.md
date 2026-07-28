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
| **CloudSpec 资源定义或文档源头缺口**(资源未定义 / 属性不满足诉求 / 覆盖度<100% / CloudSpec 描述、字段解释、枚举文案错误)——文档源错误在 OK 判定前短路；不改派、不建关联单；open-jarvis 在原主单内修 CloudSpec feature 分支、发布 pre、验证 Meta，并按需继续 Provider 生成/PR/ACC | open-jarvis（内部执行） | 原主单保持当前承接关系 |
| 与镇元不相关 + 资源代码由生成器产出(修复=acube 重跑生成器,管道不变) | 临钧 | 429768 |
| **与镇元不相关**(纯 datasource / 镇元 OK 但 provider 侧问题；文档场景仅限 CloudSpec 源正确、Provider 本地生成/展示偏差)且**紧急** | 新山 | 521957 |
| 与镇元不相关(同上)且**不紧急**(默认兜底) | 过载 | 484483 |
| **NPE 兜底**(以上所有分支均未匹配 / 跨多产品无单一负责人 / 分诊模糊超出团队职责)+ 打标签 `jarvis-npe` | 夏节 | 401498 |

**CloudSpec 原主单自闭环契约**:
- PD 判定 CloudSpec 缺口后返回 `requested_external_actions: []`、`next=terraform-rd/dev`，不得提出建单、关联、改派或个人身份动作。
- RD 使用 `cloudspec-amp-workflow` 取得 task 专属 feature 分支与模型仓 SSH URL，再调用 IDL/resource/operation/build/norm skills；build/check 与 pre 证据统一交 finalizer。
- AMP 登录、SSH 或仓库权限失败返回 `missing_capability` / `blocked` 并写原主单阻塞；不得改派外部承接人。
- `amp publish pre` 完成后保持 `release/idle`；prod/online、master/main merge/push 与正式发布仍是人工硬门。

**与镇元不相关的问题**(不进入 CloudSpec 原主单自闭环的两类):
1. **纯 datasource 问题**:诉求只涉 `data.alicloud_xxx`(查询/过滤/输出字段),不涉资源 schema/生命周期——datasource 是 provider 侧对查询 API 的只读封装,镇元只管资源 schema,**跳过镇元查证**;resource+datasource 混合不算"纯"
2. **镇元侧无问题、provider 侧存在问题**:镇元 OK 四条件全满足,缺口在 provider 实现；
   文档场景仅包含 CloudSpec 文档源正确、Provider 本地文档生成/展示偏差

**镇元 OK 四条件**(全满足才算 OK,任一不满足即视为 NOT OK=与镇元相关):
1. **API 在镇元有对应资源**:资源已在镇元定义并发布(get 返回 data 且 released list 命中)
2. **当前资源属性满足客户诉求**:比对客户抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖**(缺字段=NOT OK,即便覆盖度分再高也不算 OK)
3. **测试覆盖度 100%**:acube V2 `CoverageDetail.CoverageScore == 1.0`
4. **文档源正确性**:涉及文档诉求时，CloudSpec resource/property/operation description 与枚举
   文案必须符合 OpenAPI 长期语义；源错误不受前三项全绿影响，直接进入分支 E

详见 Step 2 分支 D 前的判定说明。

### 上游 API 缺口(纯上游产品团队问题)

- 不走上述任一路由,**@提单人**(工单 creator)+ status=待上游排期
- 由提单人协助转对应云产品的 API/OpenAPI 团队评估;jarvis 不代跨团队协调

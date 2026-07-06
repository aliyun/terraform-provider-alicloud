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

| 场景 | 花名 | 工号 |
|---|---|---|
| **Provider 侧全局改造**(非单一产品/资源:region 白名单/框架 utility/公共 endpoint/provider.go 基础/SDK bump) | 新山 | 521957 |
| 镇元 NOT OK(资源未定义 / 属性不满足诉求 / 覆盖度<100%),非紧急且距 DDL ≥14 天 | 谜拟 | 479782 |
| 镇元 NOT OK(同上),紧急 或 距 DDL <14 天 | 新山 | 521957 |
| 镇元 OK + provider 代码由生成器产出 | 临钧 | 429768 |
| 镇元 OK + provider 代码手写(默认兜底) | 过载 | 484483 |
| **NPE 兜底**(以上所有分支均未匹配 / 跨多产品无单一负责人 / 分诊模糊超出团队职责)+ 打标签 `jarvis-npe` | 夏节 | 401498 |

**镇元 OK 三条件**(全满足才算 OK,任一不满足即视为 NOT OK):
1. **API 在镇元有对应资源**:资源已在镇元定义并发布(get 返回 data 且 released list 命中)
2. **当前资源属性满足客户诉求**:比对客户抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖**(缺字段=NOT OK,即便覆盖度分再高也不算 OK)
3. **测试覆盖度 100%**:acube V2 `CoverageDetail.CoverageScore == 1.0`

详见 Step 2 分支 D 前的判定说明。

### 上游 API 缺口(纯上游产品团队问题)

- 不走上述任一路由,**@提单人**(工单 creator)+ status=待上游排期
- 由提单人协助转对应云产品的 API/OpenAPI 团队评估;jarvis 不代跨团队协调

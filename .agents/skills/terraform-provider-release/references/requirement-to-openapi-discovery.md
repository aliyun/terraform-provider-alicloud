# 需求 → OpenAPI 圈定 → CloudSpec 定义闭环

> A/B 路是**普通分支 D / 常规 release** 的需求建模与设计 review，不因涉及 CloudSpec 就自动
> 归入分支 E。triage 已确认字段集合、类型、约束、CRUD 等结构 metadata 缺陷时才进入
> **分支 E**：修到 pre Meta 收敛并经 QA `cloudspec_pre` 验证后，回 RD 在同一源单继续
> Provider 生成、PR CI 与远程 ACC；禁止调用 Acube 或操作 528766。

一个完整的 Terraform 需求,往往**从找到正确的 OpenAPI 开始**:把 OpenAPI 能力透出到 Terraform。本 reference 覆盖三件事:

- **A 路**:用户只给了需求描述(没有现成资源)——从需求出发圈定 OpenAPI → CloudSpec 定义 → 发布 pre → 生成器产出 resource/datasource/测试;
- **B 路**:用户已提供资源(schema/cspec/文档)——review 其设计是否正确,给出修复清单;
- **C 规**:资源文档必须与 API 文档对应,描述不清就从 OpenAPI 文档补。

---

## A 路:从需求描述到 pre 发布

### A.1 需求理解

从工单/用户描述中抽取三类信号:

| 信号 | 例子 | 用途 |
|---|---|---|
| **产品线索** | 「NAS 的日志转储」→ 产品=NAS | 定位 popCode |
| **能力动词** | 开启/创建/配置/查询/删除/转储 | 推断需要哪些 CRUDL API |
| **对象名词** | 日志转储、生命周期策略、访问点 | 推断资源名(如 LogAnalysis) |

写下一句话需求重述(「用户要在 Terraform 里管理 <产品> 的 <对象> 的 <生命周期>」),回填 Aone 留痕,后续所有判定以此为锚。

### A.2 OpenAPI 检索(三通道,按序)

1. **OpenAPI Explorer 站内搜**:`https://api.aliyun.com/product/<Product>` 列该产品全部 API;按对象名词关键字(中英文都试)过滤,如 `LogAnalysis` / `日志`。
2. **amp-resource-metadata skill**:`get_runtime_api.py` / `get_api_definition.py` 拉产品 API 全集与单 API 定义(入参/出参/错误码/描述)。
3. **aliyun CLI meta**:`aliyun <product> <ApiName> --help` 快速核对单个 API 的官方 meta。

圈定结果按 CRUDL 归位,并记录每个 API 的 version:

```
Create: CreateLogAnalysis (2017-06-26)
Read:   DescribeLogAnalysis (List 型,无单查)
Update: (无 → 全属性 ForceNew)
Delete: DeleteLogAnalysis
List:   DescribeLogAnalysis
```

### A.3 完备性判定(硬门)

| 缺口 | 判定 |
|---|---|
| 无 Create 或无 Delete | 资源生命周期不完整,Terraform 无法管理 → 记录后走人工会审/推动 API 侧补齐,**不得硬做** |
| 无 Update | 可接受:全属性 `ForceNew`,在需求重述里明示 |
| 只有 List 型 Read(无按主键单查) | 可接受但要预判规范门(M-RL-0019/0043/0047 类),Read 走 List+客户端过滤;提前准备加白理由 |
| API 存在但非 OpenAPI(内部 API) | 记录到 Aone(主流程 Step 5.5),该属性/操作不纳入本次发布 |

### A.4 CloudSpec 落定义并发布 pre

前置:确认 API 全部就位后才动 CloudSpec。

1. **定位产品空间**:`amp init --pop-code <Product> --pop-version <version>`(或 set-amp-workspace 把 popCode+version 换 projectId)。
2. **创建分支**:`amp branch create` 建 `feature/aone-<工单id>-<resource>`(amp 分支与 git 分支同名,**必须先 amp branch create 再 clone**,直接 `git checkout -b` 的分支后端不识别、无法 publish)。
3. **clone 元数据仓**,在仓目录内再 `amp init` 一次(cspec 编辑和 publish 必须在 clone 出来的目录内)。
4. **落资源定义**:
   - 存量 CRUDL API 已在 cspec `operations/` → 走 [新资源从存量 OpenAPI 推断](cloudspec-new-resource-infer.md)(`cloudspec fix resource` 反推);
   - operations 也缺 → 先用 `cloudspec-operation-edit` 从 API meta 建 operation,再用 `cloudspec-resource-edit` 建 resource(identifyDefinition 主键、operationMapping、属性注解逐项落)。
5. **build + check**:`aliyun cspec build` → `aliyun cspec check --name <Resource>`;规范问题按 `cloudspec-norm-check-fix` 修,确属 API 现实无法满足的(如 List 型 Read 三件套)走加白流程并留痕。
6. **发布 pre**:`amp publish pre --dry-run` 通过后 `amp publish pre`;等 pre 元数据收敛(可检索)。
7. **回主流程**:回到 SKILL.md Step 1.5 重新对齐(此时应得 `PRE_CLOUDSPEC_ALIGNED`),继续 Step 2 → Step 6,生成器**必须消费 pre** 产出 resource / datasource / 测试用例。

---

## B 路:用户已提供资源时的设计 review

用户给了现成的资源物料(Terraform schema 草案 / cspec 片段 / 资源文档 / 竞品对照)时,**不能直接采信**,必须 review 后再进主流程。

### B.1 反查 API

从用户物料里抽出资源名与属性清单,按 A.2 三通道反查这些属性背后的真实 API:每个属性都要能回答「它来自哪个 API 的哪个字段」。回答不了的属性 → 标记 `UNMAPPED`,是 review 的第一批发现。

### B.2 用 CloudSpec 工具对照设计

在 feature 分支上把用户设计落成(或比对现有)cspec 定义,逐项核:

| 检查项 | 工具/依据 |
|---|---|
| 主键定义(identifyDefinition.uniqueKeyFields)是否成立 | API 单查/删除入参是否真用这个键 |
| CRUDL operationMapping 是否与 A.2 圈定的 API 集一致 | 漏 update?错把 List 当 Get? |
| 属性类型/required/readOnly/enum 与 API 参数表是否一致 | `get_api_definition.py` 拉官方参数表逐项比 |
| 命名/规范 | `aliyun cspec check`,HCL 保留字冲突检查(SKILL.md Step 1.5) |
| 生命周期语义 | 无 Update 的属性是否都标了 ForceNew;异步操作是否有终态可查 |

### B.3 输出 review 结论

固定格式回填 Aone(或答复用户):

```
资源设计 Review:<resource>
- 正确项:<n> 项(列关键 3-5 条)
- 需修复:<m> 项,逐条:
  1. <属性/映射> — 问题:<与 API 差异> — 依据:<API 文档 URL/字段> — 修法:<改 cspec 哪里>
- UNMAPPED 属性:<清单>(找不到支撑 API,需用户确认来源或砍掉)
```

需修复项走 A.4 的 CloudSpec 闭环落修 → publish pre → 回主流程生成；这是普通分支 D/release
设计路径。若输入已被 triage 判为分支 E，则修到 pre Meta 收敛并经 QA 验证后，回 RD
在同一源单继续 Provider 主流程。
**用户设计正确也要留痕**(「review 通过,无修改」),不留痕视为未 review。

---

## C 规:资源文档与 API 文档对齐

资源/数据源文档(website/docs)的每个 attribute 描述,必须与其来源 API 的官方文档对应:

- **来源**:`get_api_definition.py`(amp-resource-metadata)或 `api.aliyun.com/api/<Product>/<Version>/<Action>` 页面的参数 description(中英文)。
- **规则**:
  1. 文档描述与 API 参数 description **语义一致**,不得自造语义;
  2. 资源文档描述**不够清晰**(如只有「The name of the resource」这类空话)→ 从 OpenAPI 文档的参数描述取内容补充(取值范围、格式、默认值、约束条件);
  3. 枚举值、取值范围、单位必须与 API 文档一致并写全;
  4. API 文档本身缺失/含糊 → 不编造;标记该属性文档待补,必要时在 Aone 记录推动 API 文档补齐。
- **核对时点**:Step 6 生成文档后(6.1 后处理清单)与 Step 11 提 PR 前各过一遍;PR 里文档改动要能经得起 reviewer 拿 API 文档对照。

# scenario-authoring —— 怎么写一个 probe 场景

写新场景时按此规范。核心信念：**像真实客户一样参考文档**——宁可保守照抄官方示例，不要凭记忆编字段。

> 场景库里全是 **tier-1 生命周期场景**(真实 apply),`scenario.yaml` **无 `tier` 键**。
> tier-0 是**资源级**静态扫描(`probe.sh tier0`),不写成场景。

## 场景库位置与目录结构(独立数据仓 + 两级布局)

- 场景语料库是**独立 git 数据仓** `tf_playground`(gitlab `terraflow/tf_playground`),与 jarvis 仓分离,
  按**云产品维度**两级归档。**写入模型 = 直推 master + 工单报备**(数据仓,非代码,不走 MR;多机 `git pull` 同步):

  ```
  terraform_playground/
    <product>/            ← 一级 = 云产品短名(小写),e.g. vpc / oss / ram
      <id>/               ← 二级 = 场景(= 目录名,scenario.yaml 的 id 须与之一致且全局唯一)
        scenario.yaml     ← 扁平元数据(probe.sh 用 grep/sed 解析)
        main.tf           ← 场景 HCL(以官方文档示例为骨架)
        checks.md         ← 期望行为 + 文档链接 + 探测点
        step2/main.tf     ← 可选:update 场景覆盖层
  ```

- **场景 id 跨 product 全局唯一**:`probe.sh run <id>` 按 id 跨 product 目录检索;同 id 出现在多个
  product 目录会**明确报错**(退 2)。新产品新建一级目录即可。
- **路径解析优先级**(probe.sh `probe_playground_dir`):env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` >
  **`bootstrap/workspace.sh dir tf_playground`**(数据仓已登记 `config/workspaces.json`,多机 clone 到
  `${JARVIS_WORKSPACE_ROOT}/tf_playground` 后**零配置可用**)> 默认 `<jarvis 根目录的父目录>/terraform_playground`。
  本机若目录名为 `terraform_playground`,用 env 或 `workspaces.local.json` 覆盖;base 配置不写绝对路径。

## 场景来源优先级

1. **provider 仓 `website/docs` 示例**（最高）：`bootstrap/workspace.sh dir terraform_provider` 拿本地仓，
   逐资源读 `website/docs/r/<name>.html.markdown` 的 Example Usage + Argument Reference，确认目标版本下
   资源名存在、必填参数、字段是否废弃/改名、块是内联还是拆分资源。
2. **阿里云官方 TF 教程 / 最佳实践**：registry 上的 how-to、云产品 IaC 文档。
3. **terraform-alicloud registry modules**（alibaba/* 官方 module）的组合方式。
4. **真实工单回灌**：aone-triage 处理过的客户问题 → `regression-<aone-id>` 场景(回灌流程见下「工单回灌」)。

## persona 定义

| persona | 客户画像 | 探测侧重 |
|---------|----------|----------|
| `beginner` | 新手，照抄单资源/最小组合示例 | 文档示例能否直接跑通（validate/plan/apply） |
| `composer` | 进阶，把多资源/多内联块拼成真实用法 | 组合下的 schema 冲突、JSON 归一化、永久 diff |
| `updater` | 已上线后改配置重跑（step2 覆盖层） | 更新是否生效、改/删字段后是否永久 diff、误重建 |
| `importer` | 用 `terraform import` 纳管存量资源 | import 是否还原完整 state（import 断链） |

## 批量生成(probe-corpus.sh)—— 从 website docs 机械造场景

手写场景之外,`bootstrap/probe-corpus.sh` 把 provider 仓 `website/docs/r/<name>.html.markdown` 的
**Example Usage** 抽出来机械改造成场景三件套,直落 playground `<product>/<id>/`。playground 是独立数据仓
`tf_playground`(直推 master + 工单报备,不走 MR),入库后 `git add/commit/push` 数据仓并在工单报备路径。

```bash
# 单资源(--force 覆盖已存在)
bootstrap/probe-corpus.sh gen alicloud_vpc
# 批量:全量 website 资源 diff 掉 playground 已有,免费族(config corpus.free_prefixes)优先,生成 N 个
bootstrap/probe-corpus.sh gen --batch 30
# 质量门:terraform init+validate+fmt -check;失败场景移 <playground>/_quarantine/<id>/ + QUARANTINE_REASON.txt
bootstrap/probe-corpus.sh validate --all          # 或 validate <id>...
```

- **机械改造项**:pin provider 版本(`config .provider.version`)、注入 `variable "run_id"`、剥离示例里的
  `terraform/provider/backend` 块(heredoc 安全)、引用到的额外 provider(如 `random`)补进
  `required_providers`(hashicorp/ 命名空间)、可命名 `name`/`*_name` 的**纯字面量或 `var.name`** → `"probe-${var.run_id}"`
  (**已带 `${...}` 插值的名称不动**,保留其原生唯一化)、示例既有 `tags` 块注入 `managed_by`。
- **产品目录**:源码 `RpcPost` 三元组 product 小写 > docs `subcategory` 净化 > `misc`。
- **产物是骨架**:`scenario.yaml`/`checks.md` 标注「生成骨架, 待人工校订」,`generated_by: jarvis-probe-corpus` +
  **`origin: generated`**;字段/探测点须**人工核对官方文档后收敛**,别直接当成品发探测。
- **限制(机械层)**:内联单行 `tags = {...}` 不注入(只认多行 opener);`domain_name`/`host_name` 等语义非资源名的
  `*_name` 也会被注入 → 若破坏 validate,会在质量门被 quarantine(可人工挪回改)。

## generated 来源纪律(`origin: generated`)—— 防机器配置污染 528766 池

- 所有生成场景 `scenario.yaml` 带 **`origin: generated`**。机器抽文档 + 机械改造的配置**本身可能有缺陷**
  (示例多块交叉引用、参数缺省、机械注入破坏格式),其 `apply_fail`/`validate_fail` **未必是 provider bug**。
- **判定纪律(skill 判定层必守)**:`origin: generated` 且**未经人工校订**的场景跑出 `apply_fail` 时,
  **先归「场景质量疑点」入人工校订队列,不得直接当 provider bug 建 528766 内部研发单**。须人工确认
  「配置无误、确系 provider 行为」后才升级为 bug。手写/回灌(`regression-*`)场景无此限,按常规分诊。
- 人工校订通过后,把 `origin: generated` 改为 `origin: curated`(或删键),该场景才进入常规 bug 分诊面。

## apply 门(成本安全门)—— `apply:` 键(值敏感)

- `scenario.yaml` 可选键 **`apply:`(默认 `true`)**。写 `apply: false` → `probe.sh run` **止步 plan**
  (init/validate/plan 后不 apply,verdict 记 `env_issue apply_disabled_by_scenario`)。
- 生成器按 `config .tier1_risk_denylist` **值敏感**判定 —— **按量付费(PostPaid/PayAsYouGo)大件一律直接 apply**
  (无按量大件资源名清单)。仅命中**订阅语义**才收窄:
  - `charge_value_fields`(`instance_charge_type`/`payment_type`/`charge_type`/`pricing_cycle` 等)的**字面量值**
    ∈ `subscription_values`(`prepaid`/`subscription`/…,大小写不敏感);**或**
  - 独立 `period` 字段取**订阅时长**(纯整数 `1..period_subscription_max`,默认 12)——借此把订阅 `period=1`
    与秒级 metric `period=900`/`"60"` 区分开(后者放行);`retention_period`/`period_unit` 等非独立 period 不误伤。
- **命中订阅门后按有无对应 data source 分岔**:
  - **有 ds**(`website/docs/d/` 存在对应 data source 文档)→ `cost: paid` + **`apply: false`** + 生成 `ds-<id>` 只读变体,
    按**订阅类资源规范**(见下)引存量、勿真实创建;
  - **无 ds**(取不到 data source 文档,无从读回校验)→ **`apply: true` + `allow_prepaid: true`** 放行真实创建探全生命周期
    (锁死会漏探),交 runner 运行时 prepaid_guard 兜底;订阅/包年包月 destroy 可能失败,届时 `destroy_fail`(S1)兜底上报待人工清理。
- 与 **runner 侧 prepaid 守门**互补:`apply:` 是**生成期静态值敏感门**;prepaid 守门是 **plan 期运行时值敏感兜底**
  (apply 前扫 plan 命中 PrePaid/Subscription 才阻断)。手写场景确需真跑订阅生命周期时,用 `allow_prepaid: true`(且自证可销毁)。

## 订阅类资源规范 —— 有 data source 引存量,无则放行真跑

- 订阅/包年包月资源**难以 API 干净销毁**。命中订阅门后按有无对应 data source 分岔(见上「apply 门」):
  **有 ds** → `apply: false` + 生成 `ds-<id>` 只读变体引存量、不真实创建;**无 ds** → `apply: true` + `allow_prepaid: true`
  放行真跑(无从 ds 读回校验,锁死漏探),`checks.md` 注明规范与 destroy 风险。
- `ds-<id>` 只读变体:若 provider `website/docs/d/<name|names|namees>` 存在对应 data source 文档,
  抽其 `data "alicloud_…" {}` 块 + `output`,落**同 product 目录** `ds-<id>/`(纯 data,天然 `apply` 安全,探测
  data source 读路径)。取不到 data source 文档的,只在主场景 `checks.md` 留规范注记、不产 ds- 变体。

## 资源选择:prepaid 守门(不是成本白名单)

- 本机是**测试账号**,付费不设限(成本白名单成本门已撤销)。
- **真正的门是「可销毁性」**:包年包月/订阅(PrePaid/Subscription)资源多数无法 API 销毁,会破坏「零残留」纪律。
  runner apply 前扫 plan 的 `*charge_type`/`*payment_type` 字段,命中 PrePaid/Subscription 默认阻断。
- **写场景时**:优先按量付费(PostPaid)或本身不计费的资源;若资源默认包年包月,显式在 HCL 里设按量付费字段
  (如 `instance_charge_type = "PostPaid"` / `payment_type = "PayAsYouGo"`)。
- **确需探测 prepaid 生命周期**:`scenario.yaml` 声明 `allow_prepaid: true` 显式豁免(慎用,须能确认可销毁)。

## 命名 / 标签 / pin 纪律

- 每个 `main.tf` 头部固定：
  ```hcl
  terraform {
    required_version = ">= 1.5.0"
    required_providers {
      alicloud = { source = "aliyun/alicloud", version = "1.284.0" }
    }
  }
  variable "run_id" { type = string }
  ```
- provider block **不写显式 region**（沿用 `ALICLOUD_REGION` env，贴近客户用法;runner 注入解析后的 region）。
- 可命名资源名称含 `${var.run_id}`；支持 tags 的资源统一加 `managed_by = "jarvis-probe"`（外加 `run_id` 标签）。
- OSS 桶名等有格式约束的（小写、全局唯一、长度）直接用 `var.run_id`（runner 传入 `jvp-<id>-<短哈希>`，天然合规）。

## region 声明

- 默认走 config `regions.focus`（eu-central-1,重点探测方向)——**多数场景不写 `region:`**。
- 只有当场景**必须**在某 region(资源可用性/配额)时才在 `scenario.yaml` 写 `region: <r>`。
- 命令行 `probe.sh run <id> --region <r>` 临时覆盖(优先级最高)。`regions.matrix` 是未来矩阵探测的候选集。

## update / import 声明法

- **update 场景**：`scenario.yaml` 设 `update_step: true`，并放 `step2/main.tf`——**完整配置**（含 terraform 块 +
  variable + 资源），runner 会用它覆盖工作目录的 `main.tf` 再 apply。改动尽量原地更新（避免改 ForceNew 字段，
  否则测的是重建不是更新）。
- **import 场景**：`import_check: true` + 成对声明 `import_address`（如 `alicloud_vpc.main`）与 `import_id_output`
  （提供真实 id 的 output 名）。runner 会 `state rm` → `import`（id 取自该 output）→ `plan`。

## tier-0 扫描范围红线(资源级,不写成场景)

- tier-0 是 `probe.sh tier0 [alicloud_xxx ...]`(无参 = 场景 resources 并集),做文档↔源码机械 diff。
- **只测已接入 TF 的面**:某云产品资源/参数**没接入 provider**,不在 tier-0 检查范围——那是「需求」不是「bug」,
  真要提需求走 tf_customer 需求路径,不当 gap 报。

## 工单回灌(regression-<aone-id>)—— 直落数据仓 + 直推 + 工单报备

每一个被 aone-triage 处理过的**真实客户 TF 问题**,修复合入后回灌为一个 `regression-<aone-id>` 场景,用作发版前回归项。

- **落点**:场景库是独立数据仓 `tf_playground`,回灌**不走 jarvis worktree/MR**——jarvis
  **直接落** `<playground>/<product>/regression-<aone-id>/`(选对应云产品一级目录);`scenario.yaml` 的
  `source_docs` 换成 Aone 工单链接,`main.tf` 为该工单最小复现配置,`checks.md` 记「修复前症状 / 修复后期望」;
  然后在数据仓 `git add/commit/push` 直推 master(数据仓写入模型)。
- **报备**:直推后**必须在对应工单评论里报备场景路径**(`tf_playground/<product>/regression-<aone-id>/`),
  供仓库主人查验 + 提醒他机 `git pull`。
- **自检**:`bootstrap/probe.sh list` 能看到该 regression 场景、`run <id> --dry` 退 0。

## 写完自检

- `bootstrap/probe.sh list` 能看到新场景(含 PRODUCT 列)。
- `bootstrap/probe.sh run <id> --dry` 退 0 且步骤计划符合预期(region 解析正确、prepaid 守门说明合理)。
- `bash test/probe_test.sh` 全绿（校验 yaml 键齐全、id 跨 product 唯一、无 tier 键、pin 版本、prepaid 守门、region 与场景根解析优先级等）。

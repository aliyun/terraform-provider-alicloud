# CloudSpec pre 资源定义检查与修复闭环

本 reference 只服务 `terraform-provider-release` 的 pre 资源合同校验，并严格区分路由来源：

- **分支 I — CloudSpec 文档文本 metadata**：resource/property/operation description、字段解释、
  NOTE 与枚举文案，且不改变字段集合、类型、约束或 CRUD。I 不进入本修复闭环；由 finalizer
  创建或复用 `upstream.cloudspec_docs_quality`（项目 2169561，念依 373108，
  `submit_only`）。若公开 Provider docs 也错误，另保留独立 528766 紧急兜底腿，按池分别防重；
  一个池已有 relation 不能抑制另一个池的缺失补建。
- **分支 E — CloudSpec 结构 metadata 原主单自闭环**：只处理字段集合、类型、约束、CRUD、
  operationMapping、生命周期等结构合同。原主单用 CloudSpec skills + AMP 完成
  “修 IDL → build/check → publish pre → pre Meta 收敛”，不创建 2165097。
- **普通分支 D / 常规 release**：若 CloudSpec pre 本来已对齐，仍按 Provider release SOP
  继续生成、PR CI 与 ACC；不得把 E 的停点泛化到 A/F/G/H/I、纯 datasource 或纯手写
  Provider-only bug。

E 的第一段硬门是 pre Meta 收敛并通过 QA `verification_mode: cloudspec_pre`。QA pass 后返回
RD，在同一源单上下文从 pre 继续 Provider 生成/开发、PR CI，再由 QA 运行远程 ACC。E 与
D/G 一样禁止 create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766，
也不得调用 Acube `createBuildTaskV2`。历史 relation 只读，不是开发、完成、阻塞或 observe
门。CloudSpec prod/online、master/main 合并和正式 release 始终是人工硬门。

## 0. 机器人能力入口

机器人使用仓库内锁定的 `cloudspec-core` 技能快照，不依赖操作者个人目录：

```bash
bash bootstrap/cloudspec-core.sh doctor
```

依次加载并遵守：

1. `cloudspec-amp-workflow`：bootstrap、feature 分支、clone、`publish pre`。
2. `cloudspec-idl-guide`：确认 `main.cspec` 与 IDL 编辑路由。
3. `cloudspec-resource-edit`，必要时配合 `cloudspec-operation-edit` / `cloudspec-flag-mode-edit`。
4. `cloudspec-build-fix`：只在 build 失败时修编译问题。
5. `cloudspec-norm-check-fix`：修复本次资源增量的规范问题。
6. `cloudspec-resource-infer`：**仅当** cspec 里已有 CRUDL operations 但 resources/ 缺定义时使用（见 [新资源推断](cloudspec-new-resource-infer.md)）。

插件来源和纳入范围锁在 `config/cloudspec-core.lock.json`。Marketplace hooks、MCP 和 telemetry 没有 vendoring；禁止为使用这些技能临时切到个人身份。

## 1. 接单初期的 pre 一致性闸门

从工单提取 `product_code`、`resource_code`、需求属性/行为及验收标准。随后只取 pre 资源 Meta：

```bash
python3 .Codex/skills/amp-resource-metadata/scripts/get_resource_type.py \
  --service-code <product_code> \
  --resource-code <resource_code> \
  --env pre
```

逐项比较：

- 属性是否存在，名称、类型、required/optional、默认值、枚举和约束是否一致；
- Create/Read/Update/Delete 及 data source 的 List 映射是否覆盖需求；
- identifyDefinition、不可变/可更新语义、返回字段与工单验收是否一致；
- pre 中不存在资源时，按“明确缺口”处理，不得用 online 冒充 pre。

判定只允许三种：

| 结论 | 后续动作 |
|---|---|
| `PRE_CLOUDSPEC_ALIGNED` | 先跑 [HCL 保留字硬门](#hcl-保留字硬门)。普通分支 D 通过后继续 provider release SOP；分支 E 通过后视为 pre Meta 已收敛，跳到第 5 节 QA pre 验证，pass 后回 RD 继续 Provider |
| `PRE_CLOUDSPEC_GAP`（已有资源需扩展/订正） | 需求清楚且能从工单唯一确定目标定义，进入第 3 节修复闭环 |
| `PRE_CLOUDSPEC_MISSING`（全新资源） | pre 完全没这资源。查 cspec `operations/`：<br>• CRUDL 5 件套齐 → 走 [新资源从存量 OpenAPI 推断](cloudspec-new-resource-infer.md)<br>• 缺 API → API 侧尚未建，走第 2 节人工会审 |
| `PRE_REQUIREMENT_AMBIGUOUS` | 停止生成/编码，执行第 2 节人工会审通知 |

不能因为 online 定义”看起来正确”就跳过 pre；禁止回退 online、缓存 Meta 或历史生成物。

### HCL 保留字硬门

Terraform HCL 保留字**绝不能**作为顶层 schema 字段名——generator 直译，用户写 `provider = “openai”` 会被解析成 meta-argument，Terraform 找不到那个 provider 插件就在 Step 0 挂。

保留字清单（Terraform SDK v1/v2 通用）：

```
provider  data  resource  variable  output  module  locals  terraform
count     for_each  depends_on  lifecycle  connection  dynamic  self
```

pre 对齐时检查每个属性名（PascalCase 转 snake_case 后比对）：

```bash
python3 .Codex/skills/amp-resource-metadata/scripts/get_resource_type.py \
  --service-code <product_code> --resource-code <resource_code> --env pre \
  | python3 -c 'import json,sys,re; d=json.load(sys.stdin);
props=d.get(“Data”,{}).get(“Properties”,{})
BAD={“provider”,”data”,”resource”,”variable”,”output”,”module”,”locals”,”terraform”,
     “count”,”for_each”,”depends_on”,”lifecycle”,”connection”,”dynamic”,”self”}
for k in props:
    snake=re.sub(r”([A-Z])”,r”_\1”,k).lower().lstrip(“_”)
    if snake in BAD: print(f”CLASH: {k} -> {snake}”)
'
```

命中即 `PRE_CLOUDSPEC_GAP`，进入第 3 节，在 cspec 里改名（保持 `@backendName` 不动，API 契约不变）：

```cspec
# 冲突前
@backendName(“provider”)
Provider: string

# 改为
@backendName(“provider”)
ModelProvider: string
```

同步改所有引用：
- `resources/<R>.cspec` 里 `@operationMapping` 块的 `resourceProperty: “$.Provider”` → `$.ModelProvider`
- `tests/*_test.cspec` 里所有 `Provider:` 字面量引用

## 2. 需求不明确时的人工会审

以下任一情况都算无法判断：工单没有可验收的字段/行为描述；描述相互矛盾；无法唯一确定类型、约束或 CRUD 语义；pre 与工单差异可能是需求变化也可能是定义错误。

在 Aone 评论中同时 @ 四人，缺一不可：

- @辰羿(320687)
- @临钧(429768)
- @过载(484483)
- @原根(265607)

阻塞内容至少包含：资源标识、工单原文摘要、pre Meta 证据、冲突/缺失点、需要人类拍板的单一问题。开发阶段不得单独评论或私信；由最终 RD 将问题与 @对象并入原主单唯一聚合回复，随后 `claim.sh release` / `jarvis-idle` 并按 `loops/persona-collab.md` 写 `[[SUSPEND]]`；不得继续 Terraform 生成、provider 编码或 CloudSpec 发布。

## 3. 测试期 CloudSpec 定义修复

先用失败请求、生成结果和 pre Meta 证明根因位于资源合同，而不是生成器实现：

- Meta 缺属性/类型/约束错误、CRUD operationMapping 错、生命周期语义错 → CloudSpec 定义问题；
- Meta 正确但生成代码丢字段、模板渲染错误、Terraform 名称映射错误 → generator 问题，转 `provider-resource-dev`。

定义问题按下列顺序执行：

1. 用 `cloudspec-amp-workflow` 完成 amp doctor、创建/切换 `feature/*` 分支并 clone 对应 cspec 模型。禁止在 master/main 编辑。
2. 确认目录已有 `main.cspec`，加载 `cloudspec-idl-guide` 后用 `cloudspec-resource-edit` 修改资源；涉及操作或 flag 模式时再加载对应专项技能。
3. 完成本轮 IDL 编辑后只运行一次 `aliyun cspec build`；失败时批量修复后再重跑。
4. 对每个变更资源逐个、前台、串行运行 `aliyun cspec check --name <ResourceName>`；同一模型
   目录禁止后台或多 Agent 并行 check。失败时批量修复，再重跑本轮 build/check。
5. 提交并推送 CloudSpec feature 分支，把 commit、MR/CR、build/check 输出摘要交 finalizer 聚合回原主单。

语义不确定时不得借 codefix 猜答案，回到第 2 节会审。AMP 登录、SSH、模型仓权限或 pre
能力失败时返回 `missing_capability` / `blocked`，不得回退个人身份、外部承接人或另建 Aone。

## 4. 发布 pre 与收敛验证

全部资源 check 通过后，发布继续由 `cloudspec-amp-workflow` 执行，且先预演后真发：

```bash
amp branch get -o json
amp publish pre --dry-run -o json --no-interactive
amp publish pre -o json --no-interactive
```

`dry-run` 失败即停止。真发成功后循环读取 `--env pre` 的资源 Meta，直到变更字段和 CRUD 映射与本次 IDL 一致；记录发布结果、轮询时间和最终 Meta 摘要。超时或内容不一致时挂起并升级，不能进入生成器。

`amp publish prod` / prod/online 正式发布以及 master/main merge/push 是人工硬门，本 SOP 永不执行。
pre 成功只表示预发模型已收敛，不得 finish，也不得宣称正式发布。分支 E 必须继续第 5 节
QA pre 验证，并在 pass 后继续 Provider/CI/ACC。

## 5. 分支 E 的 QA pre 验证与源单 Provider 续跑

QA 先使用 `verification_mode: cloudspec_pre` 独立核验 build/check/pre Meta 收敛；此阶段
不运行远程 AccTest。QA pass 后返回 `next=terraform-rd/dev`，由 RD 在同一源单上下文执行：

1. 使用已收敛的 pre Meta 重新生成或修正 Provider resource/data source/tests/docs。
2. 跑本地定向检查并提交 Provider PR，等待远程 PR CI 全绿。
3. CI 绿后交 QA 运行远程 ACC；失败回 RD 修复并重验。
4. 把 CloudSpec 分支/commit、build/check、pre publish、最终 pre Meta、Provider PR/CI 和 ACC
   证据交 finalizer 聚合到源单。

E 禁止调用 Acube `createBuildTaskV2`，也禁止任何 528766 承载动作。历史 relation 不得作为
等待交接的理由。open PR + QA pass 时 release 源单、不 finish。

普通分支 D 或常规 release 不受此 E 专用停点影响：其 CloudSpec pre 已对齐时，按主 SOP 显式
选择 pre 生成，继续 Provider PR CI 和远程 ACC。

## 6. 原主单证据清单

每个闭环至少由 finalizer 在原主单留下：

- 初始结论：`PRE_CLOUDSPEC_ALIGNED` / `PRE_CLOUDSPEC_GAP` / `PRE_REQUIREMENT_AMBIGUOUS`；
- 初始 pre Meta 的关键字段与 CRUD/List 对照；
- CloudSpec 分支与 commit（若修复）；
- `aliyun cspec build`、资源级 check 结果；
- `amp publish pre --dry-run` 与正式 pre 发布结果；
- pre Meta 收敛证据；
- QA `verification_mode: cloudspec_pre` 的 build/check/pre Meta 收敛结论；
- Provider 生成参数、PR/CI 与远程 ACC 结果；
- CloudSpec MR/CR 与仍待人工执行的正式发布工作；
- 仍待人工执行的 prod/online、master/main 与正式发布硬门。

E 的证据清单必须包含实际 Provider PR/CI/ACC 结果。普通分支 D 的 release 证据同样按主
SOP 记录生成参数、Provider PR/CI 与 ACC。

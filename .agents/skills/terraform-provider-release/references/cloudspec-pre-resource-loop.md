# CloudSpec pre 资源定义检查与修复闭环

本 reference 只服务 `terraform-provider-release`：接单初期校验 pre 资源定义，或 Terraform 测试证明 CloudSpec IDL 定义错误时，完成“修 IDL → 发布 pre → 从 pre 重新生成 → ACC 复验”。这是 **CloudSpec 原主单自闭环**：资源定义、metadata 与资源文档源头问题均以当前原主单为唯一 Aone 真源，不创建镇元侧、文档质量或 Provider 文档兜底关联单；CloudSpec 分支、MR/CR、pre、Provider PR、ACC 与 blocker 全部交最终 RD 聚合回原主单。CloudSpec 正式发布不在本闭环内。

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
| `PRE_CLOUDSPEC_ALIGNED` | 先跑 [HCL 保留字硬门](#hcl-保留字硬门)。通过 → 回填对齐证据，继续 provider release SOP。命中 → 转入第 3 节改名 |
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
3. 运行 `aliyun cspec build`。失败时用 `cloudspec-build-fix` 定位并修复，直至通过。
4. 对每个变更资源运行 `aliyun cspec check --name <ResourceName>`；用 `cloudspec-norm-check-fix` 修复本次增量，直至通过。
5. 提交并推送 CloudSpec feature 分支，把 commit、MR/CR、build/check 输出摘要交 finalizer 聚合回原主单。

语义不确定时不得借 codefix 猜答案，回到第 2 节会审。AMP 登录、SSH、模型仓权限或 pre
能力失败时返回 `missing_capability` / `blocked`，不得回退个人身份、外部承接人或另建 Aone。

## 4. 发布 pre 与收敛验证

发布必须由 `cloudspec-amp-workflow` 执行，且先预演后真发：

```bash
amp branch get -o json
amp publish pre --dry-run -o json --no-interactive
amp publish pre -o json --no-interactive
```

`dry-run` 失败即停止。真发成功后循环读取 `--env pre` 的资源 Meta，直到变更字段和 CRUD 映射与本次 IDL 一致；记录发布结果、轮询时间和最终 Meta 摘要。超时或内容不一致时挂起并升级，不能进入生成器。

`amp publish prod` / prod/online 正式发布以及 master/main merge/push 是人工硬门，本 SOP 永不执行。
pre 成功只表示预发模型已收敛：主单必须 `release/idle`，不得 finish，也不得宣称正式发布。

## 5. 强制从 pre 重新生成

pre Meta 收敛后才允许运行 Terraform generator。先查看当前 generator 的 CLI/help 或任务参数，找到其真实的 pre 环境选择器；在 Aone 证据中写明 `GENERATOR_META_ENV=pre`，并附实际命令/任务参数与生成时间。

硬门：

- 生成器必须显式选择 pre，不能依赖默认环境；
- 若当前 generator 没有可验证的 pre 选择器，返回 `missing_capability` 并停止；
- 禁止回退 online；禁止复用发布前的生成目录、缓存 Meta 或旧代码；
- 生成后检查 diff，确认修复字段/映射确实来自已收敛的 pre 定义。

随后运行 SOP 的静态检查和全部远程 ACC。仍失败时重新分类：定义仍错则回第 3 节；Meta 已正确但产物错则转 generator 问题。每次修复最多重试 3 轮，超过上限升级人工。

## 6. 原主单证据清单

每个闭环至少由 finalizer 在原主单留下：

- 初始结论：`PRE_CLOUDSPEC_ALIGNED` / `PRE_CLOUDSPEC_GAP` / `PRE_REQUIREMENT_AMBIGUOUS`；
- 初始 pre Meta 的关键字段与 CRUD/List 对照；
- CloudSpec 分支与 commit（若修复）；
- `aliyun cspec build`、资源级 check 结果；
- `amp publish pre --dry-run` 与正式 pre 发布结果；
- pre Meta 收敛证据；
- `GENERATOR_META_ENV=pre` 及实际生成参数；
- 生成 diff 摘要与全部 ACC 结果。
- MR/CR、Provider PR/CI 与仍待人工执行的 prod/online、master/main、正式发布硬门。

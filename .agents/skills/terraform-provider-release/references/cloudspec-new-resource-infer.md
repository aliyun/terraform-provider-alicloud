# 全新资源：从存量 OpenAPI 推断 CloudSpec 资源

> 本文是**普通分支 D / 常规 new-resource release** 的建模路径，不因修改 CloudSpec 就自动
> 归入分支 E。只有 triage 已确认字段集合、类型、约束、CRUD 等结构 metadata 缺陷时才是
> **分支 E**：修到 pre Meta 收敛后强制 **E → D-临钧**，由 finalizer 通过 Acube
> `createBuildTaskV2` 创建或复用 528766；不得由 E 直接执行 Provider PR/CI/ACC。

本 reference 只服务 `terraform-provider-release` 的 **"接单时 pre 资源不存在，但 CRUDL 存量 OpenAPI 已在 cspec 里"** 场景。用 `cloudspec-resource-infer` 让 CLI 自动推断资源属性和 CRUD 映射，人工 review 后 build+publish pre，再走 Terraform generator。

## 何时进入本流程

Step 1.5 拉 pre Meta 得到 `PRE_CLOUDSPEC_GAP` 且工单目标是**新建资源**（provider 侧 `alicloud/resource_alicloud_<product>_<resource>.go` 也不存在），先查 CloudSpec 项目：

```bash
# CloudSpec 项目根目录
CSPEC_DIR=$(bash bootstrap/cloudspec-core.sh workspace <PopCode>_<PopVersion>)

ls "$CSPEC_DIR/operations/" | grep -iE "^(Create|Get|Update|Delete|List)<Resource>"
ls "$CSPEC_DIR/resources/"  | grep -iE "^<Resource>\.cspec$"
```

- `operations/*` 有完整 CRUDL 且 `resources/<Resource>.cspec` 缺失 → 本流程
- `operations/*` 也缺 → 属需求侧新 API 待镇元建，回 Step 1.5 挂人工会审
- `resources/<Resource>.cspec` 已存在 → 走 `cloudspec-pre-resource-loop.md` 的 §3 修复闭环，不走本文

## 步骤

### 1. 加载 resource-infer 技能

```bash
bash bootstrap/cloudspec-core.sh doctor        # 通盘检查
# 若 lock 未包含 cloudspec-resource-infer，插件来源在
# https://code.alibaba-inc.com/agenthub-inner-skills/cloudspec-core
# 补进 config/cloudspec-core.lock.json 后 doctor 重跑
```

### 2. 用 CLI 发现 CRUDL 候选并留痕

只做发现，不改文件：

```bash
cd "$CSPEC_DIR"
ls operations/ | grep -iE "^(Create|Get|Update|Delete|List)<Resource>"
grep -h "^operation " operations/{Create,Get,Update,Delete,List}<Resource>.cspec
```

回填 Aone：候选资源名、主键（从 Get/Delete path 参数）、5 个 CRUDL 操作名。

### 3. 写最小资源骨架

`resources/<Resource>.cspec` 最小可编译骨架（详见 `cloudspec-resource-infer` skill 的 Step 3 模板）：

```cspec
$version: 1
namespace: alicloud.<Product>.<PopCode>.v<PopVersion>

@document({ name: "…", nameEn: "…", zh: "…", en: "…" })
@resourceBaseInfo({
  paidType: "free"
  getRegionIdByEndpoint: true
  classification: "normal"
  deliveryScope: "region"
  availableSites: ["china"]
  hozComponentList: ["RAM", "ACTIONTRAIL", "RG"]
  statusDefine: []
})
@terraform({ enable: true })
resource <Resource> {
  identifyDefinition: {
    @readonly
    @document({ name: "资源一级ID" zh: "…" en: "…" })
    <PrimaryKey>: string
  }
  create: Create<Resource>
  get:    Get<Resource>
  update: Update<Resource>
  delete: Delete<Resource>
  list:   List<Resources>
}
```

`cloudspec build` 通过才进第 4 步。

### 4. 跑推断

```bash
cloudspec fix resource -n <Resource>
cloudspec build
```

build 挂了走 `cloudspec-build-fix`。

### 5. 人工 review CLI 产物（关键）

CLI 是规则推导，**必须**读 `resources/<Resource>.cspec` 逐项检查：

| 检查项 | 怎么看 |
|---|---|
| Create-only 字段 | 从 Create 输入 body 找出所有 `@required`；Get 输出没有对应可读回的字段一律 `@rac({operatePrivateType: ["create"]})` + `@readonly` |
| 只读服务端生成字段 | Get/Update/List 输出里有但 Create/Update 输入里没有 → `@readonly @clientProhibited` |
| 相同语义字段的名字打架 | e.g. Get 里 `boundServices[*].serviceId` 是 Create 输入 `serviceIds` 的回读——CLI 不会自动识别，需要在 Get 侧属性上补 `@resourceProperty("ServiceIds")` 或在生成产物里手工在 Read 里反向映射（见 `generator-post-gen-checklist.md` bug 2） |
| RequestId / totalCount / pageSize / 包装 `data` | 不该进属性；输出层字段上标 `@notResourceProperty`，或对属性做手工清除 |
| List 有数组根节点 | Meta 里若 List 用 `rootMapping: $.data.items[*]`，检查生成的 `@nested` 是否正确 |
| **HCL 保留字冲突** | `<Resource>` 或任一属性名映射到 `provider` / `data` / `resource` / `variable` / `output` / `module` / `locals` / `terraform` / `count` / `for_each` / `depends_on` / `lifecycle` / `connection` / `dynamic` → **必须改名**，见下节 |
| 主键类型 | `identifyDefinition` 的类型和 Get 输入 path 参数、Delete path 参数一致（一般 string） |

有问题的字段：优先在 operation 上加 `@notResourceProperty` / `@resourceProperty("Name")` / `@nested` / `@skipMapping` 让 CLI 重生成；不重生成则手工订正映射块。

### 6. HCL 保留字冲突（硬门）

Terraform HCL 保留字**绝不能**作为 schema 属性名——generator 直译，用户写 `provider = "openai"` 会被解析成 meta-argument，然后找不到名叫 openai 的 provider 插件，Step 0 就挂。

保留字清单（黑名单，2026 年 Terraform SDK v1/v2 通用）：

```
provider  data  resource  variable  output  module  locals  terraform
count     for_each  depends_on  lifecycle  connection  dynamic  self
```

在 `resources/<Resource>.cspec` 里 grep 一下顶层属性名（PascalCase 转 snake_case 后跟保留字比对）：

```bash
awk '/^resource /,/^}$/' resources/<Resource>.cspec | grep -oE '^\s+[A-Z][a-zA-Z0-9]+:' \
  | tr -d ' :' | sed 's/\([A-Z]\)/_\l\1/g;s/^_//'
```

命中保留字 → 在 cspec 里改名（保持 `@backendName` 不动，API 契约不变），再 `cloudspec build`：

```cspec
# 原
@backendName("provider")
Provider: string

# 改
@backendName("provider")
ModelProvider: string
```

同时改 `@operationMapping` 的 `resourceProperty: "$.Provider"` → `$.ModelProvider`，以及 `tests/*_test.cspec` 里所有 `Provider:` 引用。

### 7. Build → commit → push → publish pre

```bash
cd "$CSPEC_DIR"
git switch -c infer-<resource>     # 或复用工单绑定的 feature 分支
git add resources/<Resource>.cspec operations/*.cspec tests/*_test.cspec
git commit -m "feat(<Resource>): 新增 <Resource> 资源定义"
git push origin HEAD

# 用 cloudspec-amp-workflow 发预发
amp login                          # 已登录跳过
amp init --pop-code <PopCode> --pop-version <PopVersion> --branch <cspec_branch> --yes
amp publish pre --release-mode single --dry-run --yes
amp publish pre --release-mode single --yes
amp publish status --publish-id <ProcessId> --output json --yes    # 轮询到 status=FINISH
```

轮询到 `overall status: FINISH` 且四大阶段（构建/准入/发布/准出）全 PASS/SKIP 才进 Terraform generator。

### 8. 回填 Aone

- 资源名 + CRUDL 操作名 + 主键
- CLI 推断结果差异项（review 修了什么）
- 保留字冲突是否命中（若命中，附改名前后对照）
- `cloudspec build` / feature 分支 commit / `amp publish pre` ProcessId + 最终 status
- 收敛后的 pre Meta 摘要（`get_resource_type.py --env pre` 输出）

## 与 `cloudspec-pre-resource-loop.md` 的关系

- 本文 = **无中生有**：resources/ 里没定义，用 CLI 从 operations/ 反推。
- pre-resource-loop = **修已有**：resources/ 已发布到 pre，测试证明定义错，改属性/映射/生命周期后 republish。
- 本文的普通分支 D/new-resource release 在 `amp publish pre` 收敛后回到 SKILL.md Step 6
  从 pre 生成。
- `cloudspec-pre-resource-loop.md` 若由分支 E 进入，则以 pre Meta 收敛为 E 的停点，转
  E → D-临钧 / `createBuildTaskV2`，不由 E 自己继续 Provider。

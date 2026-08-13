---
name: amp-doc-backend
description: 通过 Jarvis 受控 AMP wrapper 查询和治理镇元 API、struct、resource 文档。用户提到 AMP 文档、文档治理、API/Struct/资源文档、文档审核、SaveDraftDocMeta 或 RecommendResourceDocs 时必须使用；API/struct 只保存草稿并申请审核，resource 只走 recommend-resource 并在操作后验证 online 状态。
compatibility: 需要 Jarvis 仓库 bootstrap/amp_safe.py、已登录 AMP CLI、task 专属 CloudSpec feature 分支及有效 .amp/context.yaml。
---

# AMP 文档后端治理

这个 skill 只通过 Jarvis 的安全 wrapper 操作 AMP。不要直接调用底层 CLI，也不要安装、升级或替换本机 CLI；wrapper 报缺失或版本不兼容时，返回 `missing_capability`。

先解析两个绝对路径：

- `<jarvis-root>`：当前 Jarvis 仓库根目录。
- `<model-root>`：AMP 创建并 clone 的 task 专属 CloudSpec 模型仓，当前 Git 分支和 `.amp/context.yaml` 的 `branch` 都必须是同一个 `feature/*` 分支。

每条命令使用以下固定入口；`--repo-root` 防止从错误目录取得 project/branch 身份：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> context show
```

所有 wrapper 调用自动补 `-o json --no-interactive`，命令示例不要重复传这两个参数。写操作还会绑定当前 Aone task fence、模型仓基线、project、feature branch 与文档目标；授权失败就停止，不绕过 wrapper。

## 类型与不可跨越的边界

先明确 `type=api|struct|resource`，不能根据名称猜类型。

| 类型 | 唯一路径 | 成功语义 |
| --- | --- | --- |
| `api` | `get` → 白名单编辑 docmeta/apimeta → `create` 草稿 → `list-approver` → `submit-audit` → `get-audit-url` | 仅表示草稿已提交审核，绝不表示正式发布 |
| `struct` | 与 API 相同 | 仅表示草稿已提交审核，绝不表示正式发布 |
| `resource` | 只执行 `recommend-resource`，随后 `get --type resource --env online` 验证 | 只有 online 查询证明新结果已生效时才可宣称发布；仅出现审核单/审核链接时必须写“待审核” |

API/struct 没有任何允许的直发命令。resource 不允许 `create` 或 `submit-audit`，也不存在另一个 resource publish 命令。结构字段集合、类型、约束、CRUD、operationMapping 或生命周期发生变化时退出本 skill，转 CloudSpec 分支 E。

## 只读命令

只读命令不申请 task-fence 写授权，但仍必须走 wrapper：

```bash
# 查询文档；type 可为 api、struct 或 resource。
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc get \
  --type <api|struct|resource> --doc-key <name> \
  --language ZH_CN --env online --project-id <project-id>

# 查询可用审核人。
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc list-approver --project-id <project-id>

# 查询审核链接。
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc get-audit-url \
  --type <api|struct|resource> --doc-key <name> \
  --language ZH_CN --project-id <project-id>

# 只查询审核人角色；不带 reason，因此是只读。
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc approver-role --project-id <project-id>
```

只查询/预览时到这里为止，不保存草稿、不提交审核、不推荐资源。

## API/Struct：草稿与审核申请

### 1. 读取两份元数据

执行 `doc get`，从成功响应的 `data.docMetaJsonString` 解码 docmeta，从 `data.metaJsonString` 解码 apimeta。它们是 JSON 字符串，先解码一层；不要把响应外层对象传给 `create`。

若任一元数据缺失，只有在能从当前权威 API/Struct 定义构造时才补齐。定义也缺失时停止，不能凭空编造参数结构。JSON 形状和初始化边界见 [references/metadata.md](references/metadata.md)。

### 2. 只改文档字段

复制原始 apimeta，只允许更新 `description`、`title`、`example`、`enumValueTitles`。保留参数/
属性集合、数组顺序、`required`、`type`、`format`、`default`、约束、`$ref`、响应码、扩展字段和原枚举值。
枚举值集合的新增、删除或替换属于结构 metadata，必须退出 I 并转分支 E。

- API：遍历 `paths.*.*.parameters[*].schema` 与 `responses.*.schema.properties`，沿 `properties`/`items` 递归。
- Struct：遍历 `components.schemas.*` 及其递归 `properties`/`items`。
- docmeta 只补可读标题、摘要、说明和与实际文档 diff 对齐的 `paramExtraInfo`；无可靠来源时不填 branch、commit、cloudType、权限、SDK 或业务数据。
- 保留已有高质量内容，只修用户指出的问题或补缺。提交前做递归 diff，任何非白名单变化都停止。

### 3. 保存草稿

wrapper 只接受不超过 256 KiB 的 JSON object，并拒绝重复键、异常深度和非标准 JSON。把两份 JSON 从受控临时文件读成两个独立参数：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc create \
  --type <api|struct> --doc-key <name> \
  --doc-meta "$(< /absolute/path/docmeta.json)" \
  --meta "$(< /absolute/path/apimeta.json)" \
  --language ZH_CN --project-id <project-id>
```

`create` 成功只表示草稿保存成功。

### 4. 选择审核人并提交

先调用 `list-approver`，从返回列表中选择第一个非空、合法的 `empId`，不得伪造。空列表时停止；如果源单明确要求申请审核人角色，可执行以下写操作：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc approver-role \
  --reason <reason> --project-id <project-id>
```

有审核人后提交 JSON 数组，再查询审核链接：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc submit-audit \
  --type <api|struct> --doc-key <name> \
  --auditor-emp-ids '["<emp-id>"]' \
  --language ZH_CN --project-id <project-id>

/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc get-audit-url \
  --type <api|struct> --doc-key <name> \
  --language ZH_CN --project-id <project-id>
```

最终结论写明：文档类型、名称、白名单 diff、草稿结果、审核人、提交结果与审核 URL，并明确“已进入审核，尚未正式发布”。

## Resource：推荐并做 online 后验验证

resource 唯一写命令是 `recommend-resource`。必须指定单个 `--resource-name`、显式 `--env online` 和 project；不允许省略资源名去批量处理整个服务。

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc recommend-resource \
  --resource-name <resource-name> --env online \
  --project-id <project-id>
```

若已从只读 `list-approver` 得到合法审核人，可额外传 `--auditor-emp-id <emp-id>`。命令由服务端生成资源文档并发起后端流程；返回 0 本身不是“已发布”的充分证据。

成功后必须立即查询 online 文档：

```bash
/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py \
  --repo-root <model-root> doc get \
  --type resource --doc-key <resource-name> \
  --language ZH_CN --env online --project-id <project-id>
```

根据响应做后验判定：

- online 文档内容/版本/生效时间能证明本次推荐结果已生效：可以宣称资源文档已发布，并附验证证据。
- 推荐响应或后续查询只显示审核状态、审核 ID、审核 URL，或 online 仍是旧版本：如实写“已提交，待审核”，可再调用只读 `get-audit-url`。
- online 查询失败、响应缺关键字段或无法关联本次变更：写“发布状态未验证”，不能宣称发布成功。

## 失败与安全处理

- `PROJECT_ID_MISSING`：用受控 `context show` 查上下文；不要使用 skill 自带或临时伪造的 `.amp/context.yaml`。
- wrapper 拒绝 project、branch、document 或 receipt：停止并修正当前 task 专属工作区，不复制授权、不切个人身份。
- JSON 解析失败、重复键、超限或白名单 diff 失败：不执行 `create`。
- 审核人为空：不执行 `submit-audit`；只有源单授权治理时才带 reason 申请角色。
- `recommend-resource` 失败：检查 project、资源名、online 环境与权限；不要降级成手工 resource 草稿或杜撰发布命令。
- 正式 CloudSpec prod/online publish、主干合并和 Provider release 的既有人工硬门不因本 skill 改变。

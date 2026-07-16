# Terraformer 资源开发 Skill v0.1 设计

> Aone: [84375416](https://project.aone.alibaba-inc.com/v2/project/2100304/req/84375416)
>
> 状态：初版范围与关键技术边界已由用户确认；Skill v0.1 已按本文实现。

## 1. 背景

Jarvis 已有 `provider-resource-dev`，能够处理 Terraform Provider 资源的生成、手工修复、回归测试与 PR 交付，但尚无面向 Terraformer 的专用开发流程。Terraformer 的核心任务不是实现资源 CRUD，而是：

1. 枚举当前账号下可导出的远端资源；
2. 为每个远端资源构造与 Terraform Provider 完全一致的 Import ID；
3. 调用已安装的 Terraform Provider 读取资源状态；
4. 将状态转换成可用的 Terraform state 和 HCL。

因此，Terraformer 接入不能照搬 Data Source 或 Provider Resource 的开发流程。当子资源 List 必须提供父上下文时，Data Source 可以要求用户填写父资源 ID，Terraformer 面向全量导出则需要自行发现父资源，再枚举子资源。这个“父资源 List”只是一类资源发现模式；复合 Import ID 本身并不意味着必须遍历父资源。

本设计新增 `terraformer-resource-dev` Skill，复用 `provider-resource-dev` 的源码取证、工作区治理和验证纪律，但建立 Terraformer 自己的资源发现、Import ID、分页与生成结果验收规则。

## 2. 目标与边界

### 2.1 目标

- 支持一个新 Alicloud Terraformer 资源从分析到可导出的完整接入。
- 支持现有资源的发现不全、Import ID 错误、分页缺陷、endpoint/client 缺失和生成结果异常修复。
- 从 Terraform Provider 源码和 Import 文档确定资源身份，不凭名称猜测 Import ID。
- 将 `InitResources` 的实现选择收敛为可审查的四类模式。
- 对复合 ID、父子枚举和逐父分页给出明确实现门禁。
- 用目标包测试、构建、CLI 注册检查及可选的真实只读导出验证交付结果。
- 让 Jarvis 能按既有 Aone、worktree、数字人和 QA 流程编排 Terraformer 开发任务。

### 2.2 明确排除

- 不生产、推导或维护 Terraformer 资源关联关系规则。
- 关联关系由统一位置生产；Skill 只读取其明确产物，并将适用内容补入本次改动。
- 不修改 Terraform Provider 以迁就 Terraformer；发现 Provider 契约问题时应拆分或转交 Provider 流程。
- 不自动合并分支，不执行正式发布。
- 不把 Data Source 的“父 ID 必填”设计直接复制成 Terraformer 的用户输入要求。
- 不要求每个资源都走父资源 List；只有 API 的子资源枚举确实需要父上下文时才走该分支。

## 3. Skill 入口与文件布局

Skill 名称：`terraformer-resource-dev`

建议 description：

```text
Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer, including unsupported resources, incomplete discovery, incorrect import IDs, pagination defects, endpoint failures, or invalid generated Terraform state or HCL.
```

Canonical 文件布局：

```text
.claude/skills/terraformer-resource-dev/
├── SKILL.md
└── references/
    └── alicloud-resource-development.md
```

通过 `bootstrap/mirror.sh to-codex` 同步到：

```text
.agents/skills/terraformer-resource-dev/
├── SKILL.md
└── references/
    └── alicloud-resource-development.md
```

初版不包含可选的 `agents/openai.yaml`；Jarvis 依靠 `SKILL.md` frontmatter 发现和触发该 Skill，后续只有在需要 Codex UI 展示元数据时才单独增加。

同时新增规则测试：

```text
test/terraformer_resource_dev_skill_rules_test.sh
```

职责划分：

- `SKILL.md`：触发条件、总流程、分支决策、硬门和 Jarvis 编排。
- `references/alicloud-resource-development.md`：源码定位、四类 `InitResources` 模板、复合 ID、分页、文件矩阵和验证细节。
- 规则测试：锁定关键语义并校验 `.claude` / `.agents` 镜像一致。

主 `SKILL.md` 保持短而可执行，技术细节集中在 reference，避免将某一种资源发现模式写成整个 Skill 的核心定义。

## 4. 总体开发流程

### 4.1 入口分类

收到需求后先判断：

- **新增资源**：Terraformer 尚未注册该资源；需要完成发现实现、注册、必要 client/endpoint 和测试。
- **修复资源**：已注册但导出不全、ID 错误、分页错误、读取失败或生成结果无效；只修改与根因相关的最小文件，并补锁定回归。

### 4.2 八步流程

1. 解析 Terraform 资源名、产品和 Terraformer workspace。
2. 按真源顺序查 Provider 契约、API 和同产品实现。
3. 确认 Provider 的 Import ID 结构及每一段来源。
4. 选择一种 `InitResources` 资源发现模式。
5. 实现资源文件，并按条件补注册、client、endpoint，或读取统一关系产物并按既有方式消费。
6. 补资源级测试，重点锁定 ID、分页和错误传播。
7. 执行静态测试、目标包测试、构建与 CLI 注册检查。
8. 在具备账号和现存资源时执行只读真实导出，检查 state/HCL，并运行 Terraform 校验。

### 4.3 Terraformer 运行链路

Skill 需要保持各层职责清晰：

```text
resource.InitResources
  → 枚举远端对象并返回 Provider 可识别的 Import ID
  → ProviderWrapper.Refresh
  → 调用已安装的 Alicloud Provider Import/Read 填充状态
  → ConvertTFstate
  → 生成 Terraform state 与 HCL
```

资源接入的自定义代码主要位于第一层。`InitResources` 不应复制 Provider Read 逻辑，也不应把完整 API 对象直接转换为 Terraform schema；Provider 已负责将 Import ID 读取成规范状态。生成结果异常时则沿链路判断问题属于发现、ID、Provider Read 还是 state/HCL 转换，避免在错误层打补丁。

## 5. 事实来源与取证顺序

分析资源时按以下顺序读取，前一层用于确定契约，后一层用于补充 API 调用细节：

1. Terraform Provider Resource 源码；
2. Provider Import 文档；
3. Provider Data Source 源码；
4. Provider service/client 实现；
5. Terraformer 同产品、同 API 形态的资源；
6. OpenAPI 元数据或官方 API 文档；
7. 真实账号中的只读调用结果。

各层用途：

- Resource 的 `d.SetId(...)`、`ParseResourceId(...)` 和 Import 文档共同定义最终 Import ID。
- Data Source 用于参考 List API、过滤字段和分页方式；其必填输入不自动成为 Terraformer 的输入契约。
- Provider service/client 用于确认 endpoint、API version、请求与响应结构以及错误处理。
- Terraformer 同类资源只提供本仓代码形态参考，不能覆盖 Provider 的 Import ID 契约。

如果这些来源互相冲突，必须记录冲突并以 Provider 实际 Import/Read 行为为准；不得用“看起来像”完成接入。

## 6. `InitResources` 的四类发现模式

`InitResources` 的职责是返回能够被 Provider Import/Read 的资源 ID 集合。每个资源必须先归入以下一种模式，再写代码：

| 模式 | 适用条件 | 实现要点 |
|---|---|---|
| A. 直接全量 List + 单字段 ID | List API 不需要父上下文，响应项提供 Provider 的单个 ID 字段 | 完整处理分页；直接使用该字段构造 Import ID |
| B. 单次 List 返回多段 ID 全部片段 | 同一个响应项已经包含多段 Import ID 所需全部字段 | 按 Provider 顺序和分隔符组合；不额外调用父 List |
| C. 父子遍历 | 子资源 List 必须提供父 ID，且 Terraformer 需要全量发现 | 先列父资源，再逐父列子资源；每个父资源独立重置分页 |
| D. 无法完整枚举 | API 只支持精确查询、缺少父资源枚举或受权限/区域限制 | 使用明确的 scoped/filter 能力；否则报告能力边界，不伪造“全量支持” |

### 6.1 模式 C：父子遍历

仅当子资源 API 确实要求父上下文时使用：

1. 调用父资源的 List API，取得父资源 ID；
2. 对每个父资源创建独立的子 List 请求；
3. 为每个父资源重新初始化 page number、page size 或 next token；
4. 从子响应取得子 ID；
5. 按 Provider 契约组合最终 Import ID；
6. 一处父资源失败时返回带父 ID 上下文的错误，不能静默漏资源。

Data Source 在同一场景中可以把父 ID 声明为 Required，因为用户主动查询一个作用域；Terraformer 的全量导出没有这个前提，必须自己完成父资源发现。该规则只约束模式 C，不影响可直接 List 的模式 A、B。

### 6.2 模式 D：能力边界

如果父资源也无法枚举，或当前 API 只能根据已知 ID 精确查询，Skill 不应猜测父 ID。可以采用 Terraformer 已有的 filter/scoped 入口，但必须在结果、文档或诊断结论中说明限制。若仓库当前没有可承载该约束的入口，则停止实现并报告缺口。

## 7. 复合 Import ID

### 7.1 真源

复合 ID 的段数、顺序和分隔符只能来自：

- Provider Resource 的 `d.SetId(...)`；
- Provider Read/Update/Delete 中的 `ParseResourceId(...)`；
- Provider Import 文档与 Import acceptance test。

名称相似的 Data Source、API 主键或其他 Terraformer 资源不能替代这些证据。

### 7.2 实现规则

- 在遍历阶段以结构化字段携带 ID 片段，例如父 ID、子 ID、附件 ID。
- 只在生成叶子资源 ID 时按 Provider 约定连接一次。
- 任何必需片段为空都应返回可诊断错误，不能生成带空段的 ID。
- 不重复 encode、trim 或更换分隔符。
- 兼容旧代码时允许解析已有字符串，但新实现优先结构化传递，避免逐层字符串拼接。
- 单元测试至少覆盖正常复合 ID、缺段、段顺序和特殊字符边界。

父 ID 的来源取决于发现模式：可能与子 ID 同时出现在一次 List 响应中，也可能来自父资源 List。不能因为 Import ID 是多段式就自动推导出必须调用父 List。

## 8. 分页与错误处理

分页是资源发现的正确性门禁：

- 每个父资源拥有独立分页状态，禁止把上一个父资源的 page number 或 next token 带入下一个父资源。
- 请求中的 page size 与终止判断使用同一变量，不能一边请求固定值、一边用另一个常量判断。
- 优先使用 API 返回的 next token、total count、is truncated 等明确信号。
- 若只能根据返回条数判断结束，条件应与实际请求 page size 对齐。
- 空页、最后一页和恰好整页都要有测试。
- API 错误原样向上返回并增加资源、父 ID、页码或 token 上下文。
- 不把权限错误、endpoint 错误或单父资源失败转换成“无资源”。

## 9. 关联关系处理

关联关系不是本 Skill 的分析与生产职责。执行时只做三件事：

1. 定位统一生产位置输出的关系产物；
2. 读取其中与本资源明确匹配的声明；
3. 按既有消费方式补入本次改动并验证格式。

产物不存在、资源未声明或语义不清时，只记录缺口，不从 Provider schema、Data Source 参数或 API 字段反向推导关系。除非关联关系本身是明确验收项，这不阻塞资源的核心 discovery 和 Import ID 接入。

## 10. 文件变更矩阵

| 文件 | 何时修改 |
|---|---|
| `providers/alicloud/resource_alicloud_<name>.go` | 每个新增资源必需；修复时通常也是主文件 |
| `providers/alicloud/alicloud_provider.go` | 新资源需要加入 `SupportedResourceByProduct` 时 |
| Alicloud client/service 文件 | 当前产品缺 client、API version 或公共调用能力时 |
| endpoint 配置 | Provider/SDK 的默认 endpoint 无法满足 Terraformer 调用时 |
| 资源级 `_test.go` | 新增和修复都应修改，锁定 ID、分页或错误行为 |
| 统一关系产物的消费位置 | 统一产物已明确声明该资源关系时 |

默认不修改：

- Terraform Provider 仓库；
- Terraformer `cmd` 或 module 入口；
- README 和全局文档；
- 与本资源无关的公共重构。

参考同仓实现时优先选择相同发现模式，而不是只按产品名或文件名相似度选择。例如，附件类复合 ID 可参考 `resource_alicloud_ecs_disk_attachment.go`，需要作用域遍历的实现可参考具有父子 List 结构的资源，但最终 ID 仍以 Provider 为准。

## 11. 测试与验证

### 11.1 资源级测试

新增资源至少覆盖：

- Import ID 构造与段顺序；
- 分页结束条件；
- 空结果；
- API 错误传播。

模式 C 额外覆盖：

- 多父资源均被遍历；
- 分页状态按父资源重置；
- 一个父资源包含多页子资源；
- 父 List 或子 List 错误包含上下文。

修复任务必须增加一个未修复时失败、修复后通过的锁定回归；确实无法自动化时，需写明原因和替代验证。

### 11.2 仓库级静态验证

依次执行：

1. 目标 Go 文件 `gofmt` 检查；
2. `go test ./providers/alicloud`；
3. 将 Terraformer binary 构建到 `/tmp`，避免污染仓库；
4. 使用 CLI 列表或等价入口确认资源已注册。

当前 Terraformer checkout 的 `go test ./...` 存在与本资源无关的既有失败，因此不能把全仓测试单独作为本次成败判断。仍应运行或记录全仓基线；目标 package 新增失败必须修复，既有失败则以基线差异说明。

### 11.3 真实只读验证

有可用账号和现存资源时：

1. 使用目标产品和资源执行只读导出；
2. 检查发现数量、复合 ID、state 和 HCL；
3. 对生成目录执行 `terraform init`、`terraform validate`；
4. 执行 `terraform plan -refresh-only`，确认没有由错误 ID 或字段转换造成的读取失败。

没有账号、权限或现存资源时，不创建云资源来伪造验证；交付结论明确标记为“仅静态验证”，并列出未完成的真实验收项。

## 12. Jarvis 编排

- Terraformer 路径通过 `bootstrap/workspace.sh dir terraformer` 解析，不硬编码本机目录。
- 需求进入 `tf_customer` 路由时，沿用 Aone 分诊与 bookend；纯 adhoc 任务按 `loops/adhoc-intake.md` 建单。
- 资源分析与编码交 `terraform-rd`，验收交 `terraform-qa`；主 Jarvis 负责编排与收口。
- 修改前从 Terraformer 主干创建独立 worktree，保留用户已有 dirty changes。
- 开发结论同步 Aone；MR/PR 或 CR 创建后立即 `wrap.sh sync`。
- 合并与正式发布仍由人工决定。

## 13. Skill 自身规则测试

`test/terraformer_resource_dev_skill_rules_test.sh` 至少锁定：

- description 能触发新增资源和既有资源修复；
- reference 明确列出四类 `InitResources` 模式；
- Import ID 必须来自 Provider 契约；
- Data Source 的父 ID 必填只用于解释模式 C，不能成为全局规则；
- Skill 不生产或推导关联关系；
- 目标包测试与全仓既有失败的处理方式明确；
- `.claude` 与 `.agents` 内容一致；

增加三组前向场景用于人工或自动评估：

1. 直接全量 List 且使用单字段 ID 的新资源，应选择模式 A，不引入父资源遍历；同一响应返回多段 ID 全部片段时应选择模式 B；
2. 子 List 需要 workspace/project ID 的资源，应选择模式 C，并逐父重置分页；
3. 已有资源的复合 ID 顺序错误，应从 Provider 取证，只修改 ID 与回归测试所需文件。

## 14. 完成标准

初版 Skill 实现完成需同时满足：

- canonical Skill、technical reference 和 Codex 镜像齐全；
- 规则测试与镜像检查通过；
- 三个前向场景的决策符合本设计；
- 文档没有把父资源 List、Data Source 父 ID 或复合 Import ID 绑定成错误的全局规则；
- 文档明确关联关系只消费统一产物；
- Jarvis 的 Aone、workspace、worktree、数字人和 QA 流程均有可执行入口。

设计评审通过后已单独产出实现计划，并据此创建 Skill v0.1。

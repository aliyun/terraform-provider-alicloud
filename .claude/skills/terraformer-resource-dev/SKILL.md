---
name: terraformer-resource-dev
description: 用于开发、诊断或修复 Terraformer 中的阿里云资源，包括资源未支持、资源发现不完整、Import ID 错误或多段式、父级作用域枚举、分页缺陷、端点故障，以及生成的 Terraform 状态或 HCL 无效等场景。
---

# Terraformer 资源开发

## 核心模型

把 Terraformer 资源视为“资源发现适配器”，不要把它实现成第二套 Terraform Provider 资源：

```text
InitResources
  -> 枚举远端对象，生成与 Provider 兼容的 Import ID
  -> ProviderWrapper.Refresh 用该 ID 构造先验状态，并调用 Provider ReadResource
     （实现中还保留 ImportResourceState 回退路径）
  -> ConvertTFstate 生成 Terraform 状态和 HCL
```

按正确层次诊断问题：资源发现、Import ID、Provider 读取或状态/HCL 转换。禁止把 Provider CRUD 逻辑复制进 `InitResources`。

## 每次任务的起始动作

1. 用 `bash bootstrap/workspace.sh dir terraformer` 解析 Terraformer 仓库，用 `bash bootstrap/workspace.sh dir terraform_provider` 解析 Provider 证据仓库。任一命令失败或返回的路径不是现有目录时，停止并按 missing_capability 升级；禁止猜测或直接使用不存在的路径。克隆或同步前，先用 `bash bootstrap/workspace.sh config <key>` 读取登记的仓库、远端和默认分支。
2. 保留 Terraformer 检出目录中已有的脏文件。先拉取登记的默认分支，再为受 Git 跟踪的文件修改创建独立工作树。
3. 输入含 Aone URL 或 ID 时，仓库查证前先调用 [aone-triage](../aone-triage/SKILL.md)。没有工作项但需要修改受 Git 跟踪的文件时，先按 [loops/adhoc-intake.md](../../../loops/adhoc-intake.md) 创建或复用工作项；纯只读查证可使用仓库已定义的豁免。
4. 开工执行 `bash bootstrap/claim.sh claim <id> <pool-project>`。创建 CR/MR 后，用 `bash bootstrap/wrap.sh sync <id> "<包含 [CR](url) 的进展>"` 回填链接。未合并收尾时，依次执行 `bash bootstrap/wrap.sh done <id> "<总结>" --no-status` 和 `bash bootstrap/claim.sh release <id> <pool-project>`。
5. 将实现工作交给 `terraform-rd`，将验收验证交给 `terraform-qa`。
6. 选择 API 或写代码前，完整阅读 [references/alicloud-resource-development.md](references/alicloud-resource-development.md)。
7. 先判断任务是新资源接入还是现有资源修复。修复任务只修改已证明与根因相关的文件，并增加回归测试。

## 证据优先级

按以下顺序查证，并记录决定性证据：

1. Terraform Provider 资源源码。
2. Provider 导入文档和 Import 验收测试。
3. Provider Data Source 源码，仅用于参考 List API、过滤和分页行为。
4. Provider 服务/客户端实现。
5. Terraformer 中采用相同资源发现模式的资源。
6. OpenAPI 元数据或官方 API 文档。
7. 有凭据且已有资源时的只读 API/导出结果。

Import ID 由 Provider 的 `d.SetId(...)`、`ParseResourceId(...)`、Import 文档和 Import 测试共同定义。禁止根据名称或 Data Source 参数猜测。

## 选择一种资源发现模式

每个资源只选择一个主要 `InitResources` 模式：

- **A. 直接全量 List + 单字段 Import ID：** List API 不需要父级作用域，每条记录直接提供 Provider 所需的单个 ID 字段。
- **B. 单次 List 返回多段 Import ID 的全部片段：** 一条响应记录已包含 Provider 多段式 Import ID 所需的全部片段。
- **C. 父子遍历：** 子资源 List API 强制要求父资源 ID，因此先枚举父资源，再逐父枚举子资源；每个父资源都重新初始化分页状态。
- **D. 无法完整枚举：** 使用已有的显式作用域/过滤器输入；若无法表达缺少的作用域，则报告能力边界。

多段式 Import ID 本身不等于模式 C。Data Source 可以要求父资源 ID，因为调用方会提供查询作用域；只有子资源 List API 本身要求父级作用域时，Terraformer 才需要发现父资源。

## 只修改适用文件

- 新资源增加 `providers/alicloud/resource_alicloud_<name>.go`；只有根因属于资源自身时才修复该文件。
- 仅在缺少 `SupportedResourceByProduct` 注册或全局资源分类时修改 `providers/alicloud/alicloud_provider.go`。
- 仅在现有产品客户端无法调用目标 API 时增加客户端、服务层或端点支持。
- 增加资源级测试，锁定 Import ID 构造、分页、空结果和错误传播。
- Terraformer 任务中不修改 Terraform Provider；发现 Provider 契约缺陷时，拆分到 `provider-resource-dev` 流程。
- 禁止生产或推导资源关联关系。只读取统一关系产物，并且仅消费其中明确匹配当前资源的声明。

## 验证门禁

先执行目标检查，再执行广泛检查：

1. 确认 `gofmt` 没有报告目标文件。
2. 运行资源回归测试和 `go test ./providers/alicloud`。
3. 将二进制构建到 `/tmp/terraformer`，保持仓库干净。
4. 通过 Terraformer CLI 的注册路径确认资源可见。
5. 运行或记录 `go test ./...`；与基线比较，不隐藏既有无关失败。
6. 有账号和现有资源时，执行只读导出，检查状态文件和 HCL，并运行 `terraform validate` 和 `terraform plan -refresh-only`。

无法执行真实验证时，明确报告“仅完成静态验证”，并列出缺少的验收证据。除非用户明确授权，禁止为了验证 Terraformer 资源发现而创建云资源。

## 交付

创建 CR/MR 后保留工作树，立即关联 Aone，禁止自行合并或发布。汇报时说明所选资源发现模式、Import ID 证据、修改文件、执行过的测试、既有基线失败，以及真实验证缺口。

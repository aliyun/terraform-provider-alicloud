# Terraformer Resource Development Skill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a first usable `terraformer-resource-dev` Skill to Jarvis and update the open design CR with the implementation.

**Architecture:** Keep `.claude/skills/terraformer-resource-dev` as the canonical source and generate the `.agents/skills` mirror with `bootstrap/mirror.sh`. Keep the entrypoint concise and move Terraformer-specific discovery, Import ID, pagination, file-selection, and validation details into one reference. Enforce the agreed boundaries with a shell contract test.

**Tech Stack:** Markdown Skill files, Bash contract tests, Jarvis mirror tooling, Python Skill validator.

## Global Constraints

- Final Skill layout contains only `SKILL.md` and `references/alicloud-resource-development.md`; do not add `agents/openai.yaml`.
- `.claude/skills/terraformer-resource-dev` is canonical; `.agents/skills/terraformer-resource-dev` is produced with `bootstrap/mirror.sh to-codex`.
- Treat parent-resource listing as only discovery pattern C; a multipart Import ID alone does not imply parent traversal.
- Derive Import ID segment count, order, and delimiter from Terraform Provider `d.SetId`, `ParseResourceId`, Import docs, and Import tests.
- Data Source may require a parent ID, but Terraformer must discover it only when the selected child List API requires parent scope.
- Do not produce, infer, or maintain resource relationships; only read and consume the unified relationship artifact when it explicitly contains the resource.
- Do not modify `/Users/shanye/programs/terraformer`; it is read-only evidence for this Skill implementation.
- Preserve the existing CR worktree and update CR 28627638; do not merge or delete the worktree.

---

### Task 1: Add the Skill contract and implementation

**Files:**
- Create: `test/terraformer_resource_dev_skill_rules_test.sh`
- Create: `.claude/skills/terraformer-resource-dev/SKILL.md`
- Create: `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Create via mirror: `.agents/skills/terraformer-resource-dev/SKILL.md`
- Create via mirror: `.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Modify: `docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md`

**Interfaces:**
- Consumes: `bootstrap/mirror.sh`, `bootstrap/skills-mirror-lib.sh`, workspace keys `terraformer` and `terraform_provider`, existing `provider-resource-dev` governance conventions.
- Produces: a discoverable `terraformer-resource-dev` Skill and a deterministic contract test named `test/terraformer_resource_dev_skill_rules_test.sh`.

- [ ] **Step 1: Write the failing contract test**

Create `test/terraformer_resource_dev_skill_rules_test.sh` with this content:

```bash
#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

for rel in \
  "SKILL.md" \
  "references/alicloud-resource-development.md"; do
  claude_file="$repo_root/.claude/skills/terraformer-resource-dev/$rel"
  codex_file="$repo_root/.agents/skills/terraformer-resource-dev/$rel"
  test -f "$claude_file"
  test -f "$codex_file"

  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$codex_file" > "$expected"
  diff -u "$expected" "$claude_file"
done

expected_layout=$'SKILL.md\nreferences/alicloud-resource-development.md'
for skill_root in \
  "$repo_root/.claude/skills/terraformer-resource-dev" \
  "$repo_root/.agents/skills/terraformer-resource-dev"; do
  actual_layout="$(find "$skill_root" -type f | sed "s#^$skill_root/##" | LC_ALL=C sort)"
  if [[ "$actual_layout" != "$expected_layout" ]]; then
    echo "terraformer_resource_dev_skill_rules_test: unexpected layout in $skill_root" >&2
    diff -u <(printf '%s\n' "$expected_layout") <(printf '%s\n' "$actual_layout") >&2 || true
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/terraformer-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/SKILL.md"; do
  for term in \
    "description: 用于开发、诊断或修复 Terraformer 中的阿里云资源" \
    "# Terraformer 资源开发" \
    "## 核心模型" \
    "## 每次任务的起始动作" \
    "## 证据优先级" \
    "## 选择一种资源发现模式" \
    "## 只修改适用文件" \
    "## 验证门禁" \
    "## 交付" \
    "bootstrap/workspace.sh dir terraformer" \
    "停止并按 missing_capability 升级" \
    "aone-triage" \
    "loops/adhoc-intake.md" \
    "bootstrap/claim.sh claim" \
    "bootstrap/wrap.sh sync <id>" \
    "bootstrap/wrap.sh done" \
    "bootstrap/claim.sh release" \
    "references/alicloud-resource-development.md" \
    "terraform-rd" \
    "terraform-qa" \
    "InitResources" \
    "禁止生产或推导资源关联关系"; do
    if ! grep -Fq -- "$term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
  for old_term in \
    "description: Use when developing" \
    "# Terraformer Resource Development" \
    "## Core model" \
    "## Start every task" \
    "## Evidence hierarchy" \
    "## Evidence order" \
    "## Select one discovery pattern" \
    "## Choose one discovery pattern" \
    "## Change only applicable files" \
    "## Change only the applicable files" \
    "## Validation gates" \
    "## Delivery" \
    "prior state" \
    "fallback" \
    "schema flatten" \
    "state/HCL" \
    "Terraformer checkout" \
    "tracked files" \
    "scope/filter" \
    "child List API" \
    "client/service" \
    "endpoint"; do
    if grep -Fq -- "$old_term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: unexpected English prose '$old_term' in $skill" >&2
      exit 1
    fi
  done
done

for reference in \
  "$repo_root/.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md"; do
  for term in \
    "# Alicloud Terraformer 资源开发" \
    "## 目录" \
    "## 1. 运行时架构" \
    "## 2. 证据真源检查清单" \
    "## 3. InitResources 资源发现模式" \
    "A. 直接全量 List + 单字段 Import ID" \
    "B. 单次 List 返回多段 Import ID 的全部片段" \
    "C. 父子遍历" \
    "D. 无法完整枚举" \
    "## 4. 多段式 Import ID" \
    "## 5. 分页与错误处理" \
    "## 6. 文件选择" \
    "## 7. 测试与验证" \
    "## 8. 常见错误" \
    'd.SetId(...)' \
    'ParseResourceId(...)' \
    "多段式 Import ID 本身并不意味着必须遍历父资源" \
    "Data Source 可以要求父资源 ID" \
    "每个父资源都必须重置分页状态" \
    "使用 token 分页时，只要返回的 next token 为空就终止，不受当前页数量影响" \
    "使用页码分页时" \
    "禁止从 Provider schema、Data Source 参数或 API 字段名生产或推导关联关系" \
    "不阻塞核心的资源发现与 Import ID 支持" \
    "go test ./providers/alicloud" \
    "go test ./..." \
    "/tmp/terraformer"; do
    if ! grep -Fq -- "$term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $reference" >&2
      exit 1
    fi
  done
  for old_term in \
    "# Alicloud Terraformer resource development" \
    "## Contents" \
    "## 1. Runtime architecture" \
    "## 2. Source-of-truth checklist" \
    "## 3. InitResources discovery patterns" \
    "### A. Direct full List with a single-field Import ID" \
    "### B. One List returns every multipart-ID segment" \
    "### C. Parent-child traversal" \
    "### D. Complete enumeration is unavailable" \
    "## 4. Multipart Import IDs" \
    "## 5. Pagination and errors" \
    "## 6. File selection" \
    "## 7. Tests and validation" \
    "## 8. Common mistakes" \
    "prior state" \
    "fallback" \
    "schema flatten" \
    "import round trip" \
    "child List API" \
    "attachment" \
    "page size" \
    "page number" \
    "action" \
    "endpoint" \
    "decode" \
    "retry helper" \
    "consumer" \
    "producer" \
    "connection map" \
    "drift"; do
    if grep -Fq -- "$old_term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: unexpected English prose '$old_term' in $reference" >&2
      exit 1
    fi
  done
done

evaluation_report="$repo_root/docs/superpowers/reports/2026-07-16-terraformer-resource-dev-forward-evaluation.md"
test -f "$evaluation_report"
for term in \
  "Scenario A — PASS" \
  "Scenario B — PASS" \
  "Scenario C — PASS"; do
  if ! grep -Fq -- "$term" "$evaluation_report"; then
    echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $evaluation_report" >&2
    exit 1
  fi
done

echo "terraformer_resource_dev_skill_rules_test: PASS"
```

- [ ] **Step 2: Run the contract test and verify RED**

Run:

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
```

Expected: non-zero because `.claude/skills/terraformer-resource-dev/SKILL.md` does not exist. Confirm the failure is caused by the missing Skill, not shell syntax.

- [ ] **Step 3: Run the official initializer in scratch space**

Run:

```bash
scratch_root="$(mktemp -d /tmp/terraformer-skill-scaffold.XXXXXX)"
python3 /Users/shanye/.codex/skills/.system/skill-creator/scripts/init_skill.py \
  terraformer-resource-dev \
  --path "$scratch_root" \
  --resources references \
  --interface display_name="Terraformer Resource Development" \
  --interface short_description="Develop and repair Alibaba Cloud Terraformer resources" \
  --interface 'default_prompt=Use $terraformer-resource-dev to implement or diagnose an Alibaba Cloud Terraformer resource.'
```

Expected: initializer succeeds in `/tmp`. Do not copy its optional `agents/openai.yaml` into Jarvis; the user explicitly selected the minimal repository layout.

- [ ] **Step 4: Write the canonical Skill entrypoint**

创建 `.claude/skills/terraformer-resource-dev/SKILL.md`，所有说明性内容使用中文，专有标识、命令和文件路径保持原样：

```markdown
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
```

- [ ] **Step 5: Write the Chinese technical reference**

创建 `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md`。所有说明性内容使用中文，并保留契约测试依赖的标题和精确规则句。

```markdown
# Alicloud Terraformer 资源开发

## 目录

1. 运行时架构
2. 证据真源检查清单
3. InitResources 资源发现模式
4. 多段式 Import ID
5. 分页与错误处理
6. 文件选择
7. 测试与验证
8. 常见错误

## 1. 运行时架构

`Generator.InitResources()` 加载 Alicloud 客户端，调用一个或多个只读 API，把发现的对象转换为 `terraformutils.Resource`，再追加到 `g.Resources`。`ProviderWrapper.Refresh` 通常用该 ID 构造先验状态，并调用已安装 Provider 的 `ReadResource`；实现中还保留 `ImportResourceState` 回退路径。`ConvertTFstate` 把 Provider 返回的状态转换成 Terraform 状态和 HCL。

让 `InitResources` 只负责资源发现和生成与 Provider 兼容的 ID。禁止在其中重复实现 Provider 的 Create、Update、Delete、schema 扁平化或漂移逻辑。

## 2. 证据真源检查清单

按以下顺序阅读源码和文档：

1. Provider 资源：查找 `d.SetId(...)`、所有 `ParseResourceId(...)`、Importer 和 Read 查询参数。
2. Import 文档/测试：确认片段顺序、分隔符和导入往返验证。
3. Provider Data Source：只复用 List API、过滤条件、响应路径和分页语义。
4. Provider 服务/客户端：确认产品端点、API 版本、RPC/ROA 类型、可重试错误和响应归一化逻辑。
5. Terraformer 同模式资源：复用仓库代码惯例，不复用未经证明的身份语义。
6. OpenAPI：Provider 代码间接或由生成器生成时，用它核对请求与响应字段。
7. 真实只读调用：只有存在凭据和已有资源时才执行。

证据冲突时，以 Provider 导入/读取行为定义 ID 契约；记录冲突，禁止猜测。

## 3. InitResources 资源发现模式

### A. 直接全量 List + 单字段 Import ID

适用于一个 List API 无需父资源 ID 就能枚举全部资源，并且每条记录直接提供 Provider 所需的单个 ID 字段。优先使用 API 的显式结束信号完成分页；没有更强信号时才使用短页判断。

### B. 单次 List 返回多段 Import ID 的全部片段

适用于一条响应记录已经包含 Provider Import ID 所需的全部片段。严格保持 Provider 定义的片段顺序和分隔符。禁止仅因为 ID 是多段式就额外调用父资源 List。

### C. 父子遍历

仅当子资源 List API 要求父级作用域，并且 Terraformer 需要枚举整个账号或地域时使用：

1. 完整分页列出所有父资源。
2. 为每个父资源创建新的子资源请求。
3. 每个父资源都必须重置分页状态。
4. 完整列出每一页子资源。
5. 仅在叶子资源处按 Provider 契约拼接一次父、子 ID 片段。
6. 错误中带上父资源 ID 和页码/token 上下文；禁止静默跳过任一父资源。

Data Source 可以要求父资源 ID，因为 Terraform 调用方会主动提供查询作用域。Terraformer 的全量导出不能把这个 Data Source 输入直接转嫁给用户；只有在本模式下才自行发现父资源。

以下代码仅展示循环结构，不是可直接复制的 SDK 调用：

```go
for _, parentID := range parentIDs {
    nextToken := ""

    for {
        children, returnedNextToken, err := listChildren(parentID, nextToken, pageSize)
        if err != nil {
            return nil, fmt.Errorf("list children for parent %s: %w", parentID, err)
        }
        for _, child := range children {
            importID, err := buildProviderImportID(parentID, child.ID)
            if err != nil {
                return nil, err
            }
            ids = append(ids, importID)
        }
        if returnedNextToken == "" {
            break
        }
        nextToken = returnedNextToken
    }
}
```

以上示例只使用 token。使用 token 分页时，只要返回的 next token 为空就终止，不受当前页数量影响。使用页码分页时，递增页码，并按 API 返回的总数/页码元数据结束；没有更强信号时才使用短页判断。除非 API 明确定义两者同时存在，否则禁止混用 token 和页码契约。

### D. 无法完整枚举

适用于服务只提供精确查询、父资源无法枚举，或权限不足以执行账号级发现。已有 Terraformer 作用域/过滤器机制能够表达缺少的输入时，复用该机制；否则停止并报告限制，禁止宣称已支持完整枚举，也禁止猜测 ID。

## 4. 多段式 Import ID

片段数量、顺序和分隔符只能由 Provider Resource 的 `d.SetId(...)`、`ParseResourceId(...)`、Import 文档和 Import 测试确定。

多段式 Import ID 本身并不意味着必须遍历父资源。全部片段可能已由同一个 List 响应返回（模式 B），也可能需要先发现父资源才能得到前置片段（模式 C）。

实现规则：

- 遍历期间把父资源、子资源、挂载关系或账号片段保存在独立变量中。
- 拼接前校验每个必需片段。
- 创建叶子 `terraformutils.Resource` ID 时只拼接一次。
- 没有 Provider 证据时，禁止修剪、编码、重排或更换分隔符。
- 测试正常 ID、缺失片段、顺序、分隔符和特殊字符边界。

## 5. 分页与错误处理

- 优先使用 `NextToken`、`TotalCount`、`IsTruncated` 或同类显式信号。
- 使用返回条目数量判断时，必须与请求实际发送的每页数量比较。
- 每个父资源都重新初始化页码/token；分页状态放在父资源循环内部。
- 覆盖空首页、短末页、恰好满页和多页结果。
- 包装错误时带上操作、资源类型、父资源 ID、页码或 token。
- 权限、端点、解码和单父资源失败都必须返回错误，禁止转换成空结果。
- 复用仓库现有重试辅助函数和产品客户端惯例，禁止再实现一套重试框架。

## 6. 文件选择

| 文件 | 何时修改 |
|---|---|
| `providers/alicloud/resource_alicloud_<name>.go` | 新资源必加；修复任务仅在根因属于资源自身时修改 |
| `providers/alicloud/alicloud_provider.go` | 缺少注册或全局资源分类 |
| 产品客户端/服务层文件 | 现有客户端无法调用目标 API |
| 端点配置 | 已证明当前端点解析不足 |
| 资源 `_test.go` | 锁定 ID、分页、空结果和错误处理 |
| 统一关系消费端 | 共享产物明确声明了当前资源 |

禁止从 Provider schema、Data Source 参数或 API 字段名生产或推导关联关系。统一生产端负责定义关系语义。

统一产物没有匹配声明时，保持关系消费端不变并记录缺口。除非关联关系本身是明确验收项，否则该缺失不阻塞核心的资源发现与 Import ID 支持。

除非仓库证据证明资源无法工作，否则不要修改 `cmd`、模块入口、README、Provider 源码或无关共享代码。

## 7. 测试与验证

修复任务使用 TDD：先复现当前失败并增加最小回归测试，再实现修复。

静态门禁：

```bash
RESOURCE_FILE=providers/alicloud/resource_alicloud_example.go
gofmt -l "$RESOURCE_FILE"
go test ./providers/alicloud
go build -o /tmp/terraformer .
```

通过 CLI 的支持资源列表或等价代码路径确认注册。运行或记录 `go test ./...`；当前仓库存在既有无关失败，因此用已捕获基线比较广泛测试结果，同时要求目标包测试通过。

可以执行真实只读验证时：

1. 只导出目标产品和资源。
2. 将发现数量和 ID 与 API 响应比较。
3. 检查生成的状态文件和 HCL。
4. 在生成目录运行 `terraform init` 和 `terraform validate`。
5. 运行 `terraform plan -refresh-only`，排查读取/导入漂移。

缺少凭据或现有资源时，报告“仅完成静态验证”，并列出未验证的真实步骤。

## 8. 常见错误

| 错误 | 正确做法 |
|---|---|
| 把每个多段式 ID 都当成父子发现 | 同一响应已返回全部片段时选择模式 B |
| 复制 Data Source 的必填父资源参数 | 只有子资源 List API 要求父作用域时才枚举父资源 |
| 在父资源循环外初始化页码 | 每个父资源都重新初始化分页状态 |
| 请求使用一种每页数量，结束判断却使用另一个值 | 请求与结束判断使用同一个每页数量变量 |
| 根据 API 主键猜 Import ID | 阅读 Provider 的 `d.SetId(...)`、`ParseResourceId(...)` 和 Import 证据 |
| 根据字段名直接编辑关联映射 | 只读取统一关系产物中的明确声明 |
| 把 `go test ./...` 的基线失败当成成功或新回归 | 报告相对基线的变化，并要求目标测试通过 |
| 为了方便验证而创建云资源 | 使用已有资源，或先取得明确授权 |
```

- [ ] **Step 6: Remove optional metadata from the approved design**

Update `docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md` so section 3 shows only `SKILL.md` and `references/alicloud-resource-development.md` in both mirrors. Remove the `agents/openai.yaml` responsibility and rule-test requirement. In the completion criteria, replace `agent metadata` with `technical reference`. Add this sentence after the layout:

```markdown
初版不包含可选的 `agents/openai.yaml`；Jarvis 依靠 `SKILL.md` frontmatter 发现和触发该 Skill，后续只有在需要 Codex UI 展示元数据时才单独增加。
```

- [ ] **Step 7: Generate the Codex mirror**

Run:

```bash
bash bootstrap/mirror.sh to-codex \
  .claude/skills/terraformer-resource-dev/SKILL.md \
  .claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md
```

Expected: both `.agents/skills/terraformer-resource-dev` files are created with the repository's Claude-to-Codex token transformation.

- [ ] **Step 8: Run the contract test and verify GREEN**

Run:

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
```

Expected: `terraformer_resource_dev_skill_rules_test: PASS`.

- [ ] **Step 9: Validate the Skill and mirror**

Run:

```bash
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  .claude/skills/terraformer-resource-dev
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  .agents/skills/terraformer-resource-dev
bash bootstrap/mirror.sh check
```

Expected: both Skill validators succeed and mirror check exits zero.

- [ ] **Step 10: Run Jarvis regression tests**

Run these separately:

```bash
bash test/aone_comment_format_test.sh
bash test/wrap_done_test.sh
bash test/provider_resource_dev_skill_sync_test.sh
```

Expected: every command exits zero.

- [ ] **Step 11: Commit the implementation**

```bash
git add \
  docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md \
  docs/superpowers/plans/2026-07-16-terraformer-resource-dev.md \
  test/terraformer_resource_dev_skill_rules_test.sh \
  .claude/skills/terraformer-resource-dev \
  .agents/skills/terraformer-resource-dev
git commit -m "feat: add terraformer resource development skill"
```

### Task 2: Forward-test and deliver the Skill

**Files:**
- Modify only if forward tests expose a concrete gap: `.claude/skills/terraformer-resource-dev/SKILL.md`
- Modify only if forward tests expose a concrete gap: `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Regenerate after any canonical edit: `.agents/skills/terraformer-resource-dev/**`
- Create: `docs/superpowers/reports/2026-07-16-terraformer-resource-dev-forward-evaluation.md`

**Interfaces:**
- Consumes: Task 1 Skill, `/Users/shanye/programs/terraformer` as read-only evidence, three evaluation prompts below.
- Produces: forward-test evidence, any minimal refinement, and an updated CR 28627638.

- [ ] **Step 1: Run three fresh-context application scenarios with the Skill**

Use fresh `terraform-rd` evaluation agents. Each agent must read the Skill first and must not modify Terraformer:

```text
Scenario A: A new resource has a global List API whose items contain every segment of a two-part Provider Import ID. Explain the InitResources approach and files to inspect.

Scenario B: A child List API requires workspace_id; Provider Import ID is workspace_id:member_id. Explain discovery, pagination, and ID construction.

Scenario C: Provider schema implies a parent relation, but the unified relationship artifact has no declaration. Explain what the Terraformer change should do.
```

Expected:

- Scenario A selects pattern B and does not add parent traversal.
- Scenario B selects pattern C, discovers parents, resets child pagination per parent, and uses Provider evidence for `workspace_id:member_id`.
- Scenario C refuses to infer a relationship and records/consumes only unified producer output.

- [ ] **Step 2: Refine only demonstrated gaps**

If an evaluator violates an expected behavior, add the smallest explicit sentence or example that closes that gap, regenerate the mirror, and rerun the failed scenario. Do not add hypothetical material.

- [ ] **Step 3: Run final verification**

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py .claude/skills/terraformer-resource-dev
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py .agents/skills/terraformer-resource-dev
bash bootstrap/mirror.sh check
bash test/aone_comment_format_test.sh
bash test/wrap_done_test.sh
git diff --check origin/master...HEAD
```

Expected: every command exits zero.

- [ ] **Step 4: Commit refinements when present**

```bash
git add .claude/skills/terraformer-resource-dev .agents/skills/terraformer-resource-dev
git commit -m "docs: tighten terraformer resource development guidance"
```

Skip this commit when forward tests require no file changes.

- [ ] **Step 5: Push and update the existing CR**

```bash
git push origin worktree-terraformer-resource-dev-design
```

Verify CR 28627638 contains the design, plan, Skill, reference, mirror, and rule test; then sync the CR link and verification summary to Aone 84375416 and release the claim without changing status.

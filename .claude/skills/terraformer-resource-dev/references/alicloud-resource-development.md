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

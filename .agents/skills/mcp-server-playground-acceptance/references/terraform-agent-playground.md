# TerraformAgent Playground 验收

使用 TerraformAgent 独立 Playground 验收 MCP Server Core 的托管 token 链路和 RunIaC Token 级 HITL。不要用通用 Playground 或 settings 页面替代这条链路。

## 初始化托管 token

1. 打开 `https://pre-agent.aliyun-inc.com/terraform-agent/playground`,在 TerraformAgent 页面生成初始化步骤并打开页面给出的 Workbench。为“不设置/开启/关闭”等不同策略值分别创建全新初始化;同一“开启”临时 Agent session 可以连续执行 reject 与 approve 用例。
2. 只在 Workbench 调用 `InitializeApiMcpServerConnection`;把代理地址选为“杭州Region”,把 AK 账号选为“Terraform测试”。
3. 仍在 Workbench 使用同一个“Terraform测试”AK 调用 `GetCallerIdentity`,取得账号 UID 或完整 RAM userArn。不要在 Workbench 选择或传递 `enableHitl`。
4. 返回 TerraformAgent 页面,输入上一步得到的 UID 或完整 RAM userArn,再选择 Token 级 HITL 三态:
   - “不设置”:省略 `enableHitl`,保留服务端默认行为。
   - “开启”:传 `enableHitl=true`。
   - “关闭”:传 `enableHitl=false`。
5. 在 TerraformAgent 页面点击创建临时 Agent session。让该页面只提交身份输入和 `enableHitl`;让后端通过 `GetLatestMcpToken` 取得托管 token,再调用 `UpdateBearerTokenSafetyPolicy`。禁止让 token 返回浏览器或出现在页面、日志和错误详情里。
6. 等待临时 Agent session 创建成功,再从 TerraformAgent 页面进入实际 MCP 验收。Workbench 到第 3 步即结束,不要把它当作临时 Agent session。

## 隔离策略用例

- 将显式“开启/关闭”视为 create-only 设置。每个 token 只首次写入一个显式策略值;切换到不同策略值时创建 fresh initialization/token,不要覆写同一个 token 的策略。
- 同一个 `enableHitl=true` session/token 可以连续执行 reject 与 approve 用例,无需重新写入策略或更换 token。
- 仅当取 token 尚未消费初始化且服务端明确允许重试,才复用同一个 initialization;按“处理初始化错误”的稳定字段执行,不要一律刷新。
- 不要把 MCP endpoint URL 的 path id 当成 bearer token。path id 只标识 endpoint,不能用于调用 Token 安全策略接口。

## 验收 RunIaC HITL

使用内置 `terraform_data` 构造非零 diff,避免创建、修改或删除云资源。要求 Agent 实际调用 `RunIaC` 和 `GetTask`,不要接受文字模拟。

至少执行以下矩阵:

| Token 策略 | action | 预期 |
|---|---|---|
| 关闭 | plan | 不触发审批,返回计划结果 |
| 关闭 | apply | 不进入 `ApprovalPending`,最终 `Succeeded` |
| 开启 | plan | 不触发审批,返回计划结果 |
| 开启 | apply + reject | 先进入 `ApprovalPending`;在 `skills.aliyun.com/hitl` 拒绝后,持续调用 `GetTask` 到 `ApprovalRejected` |
| 开启 | apply + approve | 先进入 `ApprovalPending`;在 `skills.aliyun.com/hitl` 批准后,持续调用 `GetTask` 到 `Succeeded` |

每次 `RunIaC` 返回 `processID` 后,持续调用 `GetTask` 到终态。记录 `status`、`nextAction`、`summary`、`message` 的必要脱敏片段和审批状态迁移。若“开启”用例的 apply 直接成功或“关闭”用例进入审批,判定 FAIL。

## 处理初始化错误

优先用稳定字段判定是否可重试:

- 记录 `code` 作为错误分类。
- 记录 `initializationConsumed` 判断本次初始化额度或状态是否已经消耗。
- 记录 `canRetry` 决定原页面重试还是刷新并创建新会话。
- 不要把可变 message、RequestId 或敏感值作为自动化稳定断言。

按以下固定规则处理,不要只凭提示文案或 HTTP 状态猜测:

| 响应 | 动作 |
|---|---|
| `code=TOKEN_FETCH_RETRYABLE`、`initializationConsumed=false`、`canRetry=true` | 保留页面并复用同一个 initialization 重试;不要刷新或重新执行 Workbench 初始化 |
| `code=TOKEN_POLICY_FAILED`、`initializationConsumed=true`、`canRetry=false` | 刷新 TerraformAgent 页面,重新走完整初始化并取得 fresh token |
| 未知 502、网络错误或非 JSON 响应 | 刷新 TerraformAgent 页面,重新走完整初始化;不要假设旧 initialization 可用 |
| 其他结构化错误 | 以 `initializationConsumed` 和 `canRetry` 组合决定动作;只有 `false/true` 才复用,其余刷新重来 |

需要切换“开启/关闭”或在策略写入已消费后重设 HITL 时,始终使用 fresh initialization/token,避免违反 create-only 语义。

## 证据与脱敏

记录以下证据:

- TerraformAgent 预发 URL 和新 session id。
- AgentAutomation 与 Cloudspec 所需预发门禁结果。
- 三态选择以及切换不同显式策略值时使用 fresh token 的说明。
- 初始化 API `InitializeApiMcpServerConnection`、`GetCallerIdentity`;将它们与 MCP 工具清单分开记录。
- 实际 MCP 工具名 `RunIaC`、`GetTask`,以及 RunIaC action、processID、审批状态迁移和 GetTask 终态。
- 稳定错误字段 `code`、`initializationConsumed`、`canRetry`。

禁止记录或回填以下内容:

- 托管 MCP token、Bearer header、Authorization 内容。
- AK、账号 UID、完整 RAM userArn。
- 可还原敏感信息的完整审批链接或初始化请求/响应正文。

截图、HTML 报告和 Aone 追评都执行同一脱敏规则。若证据中出现上述内容,先遮盖或删除,再分享或回填。

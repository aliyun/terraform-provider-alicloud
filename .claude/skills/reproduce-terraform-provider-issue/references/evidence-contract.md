# 证据与脱敏契约

## 证据分层

保持以下层次相互独立，使产品团队无需反向推导 Terraform 输出即可判断问题归属：

| 层次 | 必需证据 | 证明内容 |
|---|---|---|
| Terraform 输入 | Provider 版本和触发字段 | 用户要求 Provider 下发的内容 |
| 创建 API | API 名、白名单请求字段、RequestId、创建资源 ID | 服务端接受请求并创建了资源 |
| 读取 API | API 名、RequestId、返回或缺失字段 | 服务端在创建后能够回读的内容 |
| 关联 API | RequestId 和关系字段 | 缺失数据是否保存在子资源或关联资源上 |
| Provider 状态 | Read/refresh 后的 state 字段 | Provider 如何映射 API 响应 |
| Plan | 详细退出码、动作、替换路径 | Terraform 生命周期后果 |

## 时间线行契约

每一行必须包含：

```text
时间戳 | API | RequestId | 目标 ID | 状态 | 白名单请求字段 | 观察结果
```

保留重复的轮询调用，用于证明缺失字段是暂时还是持续存在，并向产品团队提供多个 RequestId。

## 允许保留的请求字段

解析器不内置产品或资源字段。每次复现必须根据当前 API 契约，通过 `--request-field`、`--target-field` 和 `--observe-field` 显式选择诊断必要字段。观察字段优先写成 `API:字段`，仅在该 API 响应中判断字段是否缺失。例如：

- Region/Zone；
- 资源 ID 和名称；
- 资源类型、协议、存储和加密选项；
- VPC/VSwitch 及相关网络拓扑 ID；
- 访问组和访问规则语义；
- 复现使用的 CIDR。

未显式传入白名单时，解析器不输出请求参数、目标 ID 或观察字段。禁止将认证、授权、签名、Session、Cookie、Token、密码、UserData、私钥或任意业务载荷字段加入白名单；解析器会额外拒绝常见敏感字段名。

## 敏感产物

即使控制台输出看似已脱敏，也默认以下产物包含密钥：

- 保存的 Terraform plan 二进制文件；
- 原始 `TF_LOG` 或 Provider 调试日志；
- Shell 环境变量导出；
- Provider 配置值；
- 签名 URL 和 Authorization Header。

原始产物仅允许在本地用于提取白名单证据。生成脱敏时间线和报告后立即删除。禁止把原始日志上传到 AutomationAgent 或作为 Aone 附件。

## 报告必备内容

报告必须包含：

1. 现场状态和保留或清理警告；
2. 环境信息及相对用户输入的偏差；
3. 完整实例清单；
4. 创建请求和响应证据；
5. 首次读取响应中存在争议的字段；
6. 包含轮询调用的完整 apply API 时间线；
7. 创建后的直接 API 验证；
8. 完整 refresh API 时间线；
9. State 差异和替换路径；
10. 已知的 Provider 代码和 Schema 位置；
11. 需要产品团队确认的问题；
12. 后续排查使用的只读命令；
13. `template/main.tf` 的唯一完整 `hcl` fence；
14. 使用同一 HCL 的在线 Terraform 预填链接；
15. 交付包、fmt/init/validate、在线预览和平台状态。

## 报告交付包契约

默认交付 `REPORT.md`、`REPORT.html`、`template/main.tf`、`template/README.md`，init 生成
lockfile 时保留 `template/.terraform.lock.hcl`。禁止交付 `.terraform/`、plan、state/backup、
原始 TF_LOG、crash log、tfvars、凭据、可执行 HTML、base64 或 data URI。

安全扫描覆盖包内所有 UTF-8 文本，包括 `.txt`。JSON/HCL/YAML/Shell 中
`accessKeyId`、`accessKeySecret`、`api_token` 等字段只在带有实际值时判为凭据；仅列举字段名
的文档允许通过。HTML 检查先解码 entity，`java&#x73;cript:`、编码后的 `on*=` 等仍须拒绝。

`template/main.tf` 是唯一 HCL 源，必须以 UTF-8 换行结束。`REPORT.md` 只能有一个完整
`hcl` fence；`REPORT.html` 只能有一个 `language-hcl` code block；两者解码后都必须与
`main.tf` 字节相同。`profile` 变量固定为：

```hcl
variable "profile" {
  type     = string
  default  = null
  nullable = true
}
```

凭据只通过 `TF_VAR_profile` 或本地 CLI Profile 传入。

在线链接必须使用固定 query、同一 13 位时间戳两次，并令 `params` 精确等同 Java
`URLEncoder.encode(mainTf, UTF_8).replace("+", "%20")`。不得双重编码。链接只预填代码，
不执行 Terraform。

## 验证和平台状态

用 `scripts/validate-report-package.py` 统一执行安全扫描、字节/URL/profile 校验，以及
`terraform fmt -check -recursive`、隔离 `init -backend=false -input=false -no-color`、
`validate -no-color`。退出码 `0` 为通过，`2` 为确定性错误，`3` 为外部或平台阻塞。
Terraform 子进程只接收 PATH、locale、代理/证书、临时目录等白名单环境以及新建的隔离
HOME/TF_DATA_DIR；禁止继承 `TF_CLI_ARGS*`、`TF_VAR_*` 和任何云凭据。

在线预览记录必须包含 `success=true`、`status=uploaded`、`reportId` 和只读
`/reports/aone/.../view` 路由；要求在线交付时，再用匿名 GET 验证 HTTP 200、`text/html`、
标题、完整 HCL 和至少一个非空用户 marker。绝对 `url` 的 scheme/host/port 必须等于
`--preview-origin`，原始 path 必须逐字等于 `viewUrl`，且禁止 userinfo、query、fragment、
百分号编码和重定向。GET 禁止携带 Authorization/Cookie。HTTP preview origin 只允许本机
loopback 测试。

报告 HTML 必须为零脚本静态文档。Viewer 复制能力归 AutomationAgent 平台，当前状态记录为
`platform_blocked`；不得用 `<script>`、`<button>`、`on*=` 或其他上传内容绕过。

## 结论约束

- 只有响应内容和 RequestId 共同证明时，才能写“API 返回空值或缺失字段”。
- 只有检查 Read/refresh 后的 state 时，才能写“Provider 清空了状态”。
- 只有记录 plan 动作和替换路径后，才能写“Terraform 将替换资源”。
- 不得把“API 契约歧义”和“Provider 映射缺陷”合并为同一个结论。
- 如果使用隔离网络替换了无法访问的用户网络 ID，必须在报告靠前位置说明。

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

内置解析器仅允许诊断必要字段，例如：

- Region/Zone；
- 资源 ID 和名称；
- 资源类型、协议、存储和加密选项；
- VPC/VSwitch 及相关网络拓扑 ID；
- 访问组和访问规则语义；
- 复现使用的 CIDR。

禁止将认证、授权、签名、Session、Cookie、Token、密码、UserData、私钥或任意业务载荷字段加入白名单。

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
12. 后续排查使用的只读命令。

## 结论约束

- 只有响应内容和 RequestId 共同证明时，才能写“API 返回空值或缺失字段”。
- 只有检查 Read/refresh 后的 state 时，才能写“Provider 清空了状态”。
- 只有记录 plan 动作和替换路径后，才能写“Terraform 将替换资源”。
- 不得把“API 契约歧义”和“Provider 映射缺陷”合并为同一个结论。
- 如果使用隔离网络替换了无法访问的用户网络 ID，必须在报告靠前位置说明。

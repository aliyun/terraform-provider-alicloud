# checks: ram-role-policy

persona: composer —— 客户建一个可被 ECS 扮演的角色 + 一条自定义 OSS 只读策略，并把策略附加到角色。

## 期望行为（客户视角）

1. `validate` / `plan` 通过；三资源组合是 r/ram_role + r/ram_policy + r/ram_role_policy_attachment 的直接组合。
2. `apply` 成功；`apply` 后立即 `plan` 应为空 diff——尤其 `assume_role_policy_document` / `policy_document` 的 JSON 归一化（provider 回读云端 JSON 后与本地字符串比对）是永久 diff 高发区，重点盯 `perpetual_diff`（S2）。
3. `destroy`：`force = true` 应先解除附加再删除角色/策略，`state list` 清空。

## 文档依据（1.284.0）

- ram_role：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/ram_role
  - 新字段 `role_name`（`name` 自 v1.252.0 废弃）、`assume_role_policy_document`（`document` 自 v1.252.0 废弃）、`tags`（Map，v1.252.0）、`force`。
- ram_policy：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/ram_policy
  - 新字段 `policy_name`（`name` 自 v1.114.0 废弃）、`policy_document`（`document` 自 v1.114.0 废弃）、`tags`（Map，v1.246.0）、`force`；导出 `type`。
- ram_role_policy_attachment：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/ram_role_policy_attachment
  - 必填 `policy_name` / `policy_type`（Custom|System）/ `role_name`。

## 探测点（潜在 provider 问题）

- **JSON 归一化永久 diff**：本地 heredoc JSON 缩进/键序与云端返回不一致时，provider 若未做语义归一化会永久 diff。
- **文档 drift（已发现，值得回灌为 finding）**：`ram_role_policy_attachment` 的官方示例仍在 `alicloud_ram_role` 上用**已废弃**的 `name` / `document`，与 `ram_role` 独立文档的新字段 `role_name` / `assume_role_policy_document` 冲突——照抄 attachment 示例的客户会踩废弃字段告警。属 S4 文档类，见最终汇报疑点清单。
- 本场景 attachment 有意引用**非废弃**属性（`policy_name` / `type` / `role_name`）以对照。

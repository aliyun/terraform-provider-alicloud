# checks: net-vpc-basic

persona: beginner —— 新手照抄官方文档示例，把 VPC / VSwitch / 安全组 / 安全组规则拼成一套最典型的组网。

## 期望行为（客户视角）

1. `validate` 通过：四个资源 + `data.alicloud_zones` 的组合是官方文档示例的直接组合，应无语法/schema 报错。
2. `plan`（有凭证时）产出 4 个待建资源，`data.alicloud_zones` 正常读取到至少一个可用区。
3. `apply` 成功创建；`apply` 后立即 `plan` 应为空 diff（无 `perpetual_diff`）。
4. 安全组规则 `nic_type = "intranet"` 对 VPC 安全组是唯一合法值——若 provider 对 VPC 安全组的 `internet` 未报清晰错误或默认值处理不当，属可疑（`plan_fail` / `validate_fail`）。
5. `destroy` 能干净删除，`state list` 清空（无 `destroy_fail` / `state_residue`）。

## 文档依据（1.284.0）

- vpc：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/vpc —— 字段 `vpc_name` / `cidr_block`（`name` 已废弃）。
- vswitch：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/vswitch —— 示例用 `data.alicloud_zones.foo.zones.0.id` 取 `zone_id`；`vpc_id` / `cidr_block` 必填。
- security_group：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/security_group —— 字段 `security_group_name`（`name` 自 v1.239.0 废弃）。
- security_group_rule：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/security_group_rule —— 示例 `type` / `ip_protocol` / `nic_type=intranet` / `policy` / `port_range` / `priority` / `security_group_id` / `cidr_ip`。
- zones 数据源：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/zones —— `available_resource_creation = "VSwitch"`。

## 探测点

- VSwitch `cidr_block`（172.16.0.0/24）必须是 VPC `cidr_block`（172.16.0.0/16）子集——provider 对越界子网是否给清晰报错。
- `data.alicloud_zones` 在无凭证时于 `plan` 阶段即需鉴权（读 API）——鉴权失败归 `env_issue: auth_error`，不算 finding。

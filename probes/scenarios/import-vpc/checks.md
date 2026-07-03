# checks: import-vpc

persona: importer —— 客户对已存在的 VPC 用 `terraform import` 纳管，验证 import 是否还原完整状态。

## 期望行为（客户视角）

1. `apply` 建出 VPC，`output vpc_id` 为真实 VPC id。
2. runner 执行 `terraform state rm alicloud_vpc.main` 后 `terraform import alicloud_vpc.main <vpc_id>`：import **应成功**（无 `import_diff` 的失败态）。
3. import 后立即 `plan`：**应为空 diff**——若非空说明 import 未把云端属性完整读回本地 state（import 断链），即 `import_diff`（S2）。
4. `destroy` 干净，`state list` 清空。

## 文档依据（1.284.0）

- vpc：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/vpc
  - 文末 Import 段：`terraform import alicloud_vpc.example <id>`，import id 即 VPC id。

## 探测点（潜在 provider 问题）

- **import 后永久 diff**：Read 未回填某些可选/计算属性（如 `tags`、`description`、`route_table_id`）→ import 后 `plan` 非空。这是 import 断链最常见形态。
- **import 直接失败**：资源不支持 import 或 id 格式解析错误 → import 步骤非零退出。
- 断链定位：对照 import 后 `plan` 的 diff 字段，即“哪些属性 Read 没读回”，是 provider Read 函数补全的直接线索。

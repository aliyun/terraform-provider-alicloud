# checks: oss-bucket-config

persona: composer —— 进阶客户在一个 `alicloud_oss_bucket` 上内联叠加版本控制、生命周期、标签。

## 期望行为（客户视角）

1. `validate` / `plan` 通过：`versioning` / `lifecycle_rule` / `tags` 均为 1.284.0 主资源支持的内联块。
2. `apply` 成功；`apply` 后立即 `plan` 应为空 diff。
3. `destroy`：`force_destroy = true` 应允许删除（空桶本无对象，但开启版本控制后删除标记也需清理）。

## 文档依据（1.284.0）

- oss_bucket：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/oss_bucket
  - 主资源仍支持内联 `versioning { status }`、`lifecycle_rule { ... }`、`tags`、`force_destroy`、`storage_class`。
  - 内联 `acl` 自 1.220.0 起废弃，改用独立资源 `alicloud_oss_bucket_acl`（本场景不设 `acl`，桶默认 private）。

## 探测点（潜在 provider 问题）

- **生命周期永久 diff**：`expiration { expired_object_delete_marker = true }` 与 `noncurrent_version_expiration { days }` 组合，是 OSS provider 历史上易出永久 diff 的写法——`apply` 后立即 `plan` 若非空即 `perpetual_diff`（S2）。
- **内联 vs 拆分资源二义性**：1.284.0 文档在头部强 NOTE：若同时用内联属性和独立子资源（如 `alicloud_oss_bucket_versioning`）管理同一项，必须加 `lifecycle { ignore_changes }`，否则每次 apply 相互覆盖产生永久 diff。本场景只用内联块规避；**拆分资源冲突值得单开一个 regression 场景专门探测**（见 scenario-authoring.md 场景来源④）。
- **版本控制桶的 destroy**：开启 versioning 后，`force_destroy` 是否真能清掉所有版本 + 删除标记，值得盯 `destroy_fail` / `state_residue`。

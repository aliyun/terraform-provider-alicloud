# checks: vpc-tags-update

persona: updater —— 客户先建 VPC，再改 `vpc_name` / `description` / `tags`（不动 `cidr_block`）后重跑，模拟“改配置再 apply”。

## 期望行为（客户视角）

1. step1 `apply` 成功，`apply` 后立即 `plan` 空 diff。
2. step2（覆盖层）`apply`：`vpc_name` / `description` / `tags` **原地更新**，`cidr_block` 未变故不应触发 ForceNew 重建。
3. step2 `apply` 后立即 `plan` 空 diff——**若仍有 diff 说明“更新不生效/写回不完整”**（`perpetual_diff`，抓的就是这个）。
4. `destroy` 干净，`state list` 清空。

## 文档依据（1.284.0）

- vpc：https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/vpc
  - `vpc_name` / `description` / `tags` 均可原地更新；`cidr_block` 变更才会 ForceNew。

## 探测点（潜在 provider 问题）

- **更新不生效**：改标签（尤其**删除**某个 key，如 step2 移除 `phase=initial` 改为 `phase=updated` 并新增 `extra`）后，provider 若只做增量 add 不做 remove，会导致云端残留旧标签 → 下次 `plan` 永久 diff。
- **意外重建**：`cidr_block` 未变，若 provider 误判某字段为 ForceNew 触发 delete+create → `unexpected_replace`（S1）。
- **name 改动**：`vpc_name` 改动应为 in-place 更新调用（ModifyVpcAttribute），不应重建。

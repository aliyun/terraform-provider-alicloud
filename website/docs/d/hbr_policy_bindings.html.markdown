---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_policy_bindings"
sidebar_current: "docs-alicloud-datasource-hbr-policy-bindings"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Policy Bindings to the user.
---

# alicloud_hbr_policy_bindings

This data source provides the Hbr Policy Bindings of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.221.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_policy_bindings" "ids" {
  policy_id   = "example-policy-id"
  source_type = "OSS"
}

output "hbr_policy_binding_id_1" {
  value = data.alicloud_hbr_policy_bindings.ids.policy_bindings.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Policy Binding IDs. The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
* `policy_id` - (Optional) The ID of the backup policy. If specified, only the policy bindings under this policy are listed.
* `source_type` - (Optional) The data source type of the policy binding. Valid values: `UDM_ECS`, `NAS`, `OSS`, `File`, `ECS_FILE`, `OTS`.
  - `UDM_ECS` - indicates the ECS instance backup.
  - `OSS` - indicates an OSS backup.
  - `NAS` - indicates an Alibaba Cloud NAS Backup.
  - `ECS_FILE` - indicates that the ECS file is backed up.
  - `File` - indicates a local File backup.
  - `OTS` - indicates the Tablestore backup.
* `name_regex` - (Optional) A regex string to filter results by the policy binding description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Policy Binding IDs.
* `policy_bindings` - A list of Hbr Policy Bindings. Each element contains the following attributes:
  * `id` - The ID of the Policy Binding. The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
  * `policy_id` - The ID of the backup policy.
  * `source_type` - The data source type of the policy binding.
  * `data_source_id` - The ID of the data source bound to the policy.
  * `disabled` - Whether the policy is effective for the data source. `true`: paused; `false`: not paused.
  * `exclude` - The file type that does not need to be backed up. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `include` - The file types to be backed up. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `source` - The OSS prefix to be backed up. Valid only when `source_type` is `OSS`.
  * `speed_limit` - The backup flow control. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `policy_binding_description` - The description of the policy binding.
  * `cross_account_type` - The cross-account type. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
  * `cross_account_role_name` - The name of the cross-account authorization role. Valid only when `cross_account_type` is `CROSS_ACCOUNT`.
  * `cross_account_user_id` - The ID of the actual account to which the data source belongs. Valid only when `cross_account_type` is `CROSS_ACCOUNT`.
  * `create_time` - The creation time of the policy binding.
  * `advanced_options` - The advanced backup options of the policy binding.
    * `oss_detail` - The advanced options for OSS backup.
      * `ignore_archive_object` - Whether archived objects are not prompted in task statistics and failed file lists.
      * `inventory_cleanup_policy` - Whether to delete the inventory file after the backup.
      * `inventory_id` - The name of the OSS inventory.
    * `udm_detail` - The advanced options for ECS backup.
      * `destination_kms_key_id` - The custom KMS key ID of encrypted copy.
      * `disk_id_list` - The list of backup disks.
      * `exclude_disk_id_list` - The list of cloud disk IDs that are not backed up.

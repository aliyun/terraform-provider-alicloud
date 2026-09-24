---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_policy_bindings"
sidebar_current: "docs-alicloud-datasource-hbr-policy-bindings"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Policy Bindings to the user.
---

# alicloud\_hbr\_policy\_bindings

This data source provides the Hbr Policy Bindings of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_policy_bindings" "default" {
  policy_id   = "po-000xxxxxxxxxxxx56y"
  source_type = "OSS"
}

output "first_policy_binding_id" {
  value = data.alicloud_hbr_policy_bindings.default.bindings.0.policy_binding_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Policy Binding IDs. The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
* `policy_id` - (Optional) The policy ID.
* `source_type` - (Optional) The type of the data source. Valid values: `UDM_ECS`, `OSS`, `NAS`, `ECS_FILE`, `File`, `OTS`.
* `data_source_ids` - (Optional) A list of data source IDs. **NOTE:** The API matches `data_source_ids` only when `source_type` is set at the same time; when `source_type` is not set, the matching is performed on the client side.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Policy Binding IDs.
* `bindings` - A list of Hbr Policy Bindings. Each element contains the following attributes:
  * `id` - The ID of the Policy Binding. The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
  * `policy_binding_id` - The ID of the policy binding.
  * `policy_id` - The policy ID.
  * `source_type` - The type of the data source.
  * `data_source_id` - The data source ID. When `source_type` is `OSS`, this value is the bucket name.
  * `disabled` - Whether the policy is effective for the data source. `true`: paused. `false`: not paused.
  * `source` - The backup source. When `source_type` is `OSS`, it specifies a prefix to be backed up; if not specified, the entire root directory of the bucket is backed up. When `source_type` is `ECS_FILE` or `File`, it specifies the file directory to be backed up; if not specified, all directories are backed up.
  * `include` - The file types to be backed up. All files of these types are backed up. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `exclude` - The file types that do not need to be backed up. All files of these types are not backed up. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `speed_limit` - The backup flow control. The format is `{start}{end}{bandwidth}`. Valid only when `source_type` is `ECS_FILE` or `File`.
  * `policy_binding_description` - The description of the policy binding.
  * `create_time` - The time when the policy binding was created. This value is a UNIX timestamp. Unit: seconds.
  * `cross_account_type` - The type of the account to which the data source belongs. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
  * `cross_account_user_id` - The ID of the actual account to which the data source belongs. Valid only when `cross_account_type` is `CROSS_ACCOUNT`.
  * `cross_account_role_name` - The name of the cross-account authorization role of the data source. Valid only when `cross_account_type` is `CROSS_ACCOUNT`.
  * `created_by_tag` - Whether the resource is automatically associated through a backup policy resource tag.
  * `hit_tags` - The hit tag rules.
    * `key` - The tag key.
    * `value` - The tag value.
    * `operator` - The tag matching rule.
  * `advanced_options` - The backup advanced options.
    * `oss_detail` - The OSS backup advanced options.
      * `ignore_archive_object` - Whether archived objects are not prompted in task statistics and failed file lists.
      * `inventory_cleanup_policy` - Whether to delete the inventory file after the backup. Valid only when the OSS inventory is used. Valid values: `NO_CLEANUP`, `DELETE_CURRENT`, `DELETE_CURRENT_AND_PREVIOUS`.
      * `inventory_id` - The name of the OSS inventory.
    * `udm_detail` - The ECS instance backup advanced options.
      * `disk_id_list` - The list of disks to be protected. If it is empty, all disks are protected.
      * `exclude_disk_id_list` - The list of disk IDs that do not need to be protected. This parameter is ignored when `disk_id_list` is not empty.
      * `destination_kms_key_id` - The ID of the custom KMS key in the destination region. If this field is not empty and cross-region replication is enabled, this key is used for encrypted cross-region replication.
      * `app_consistent` - Whether to create application-consistent snapshots. Valid only when all disks are ESSD.
      * `snapshot_group` - Whether to create a snapshot consistency group. Valid only when all disks are ESSD.
      * `ram_role_name` - The RAM role name required to create application-consistent snapshots. Valid only when `app_consistent` is `true`.
      * `pre_script_path` - The path of the freeze script executed before creating application-consistent snapshots. Valid only when `app_consistent` is `true`.
      * `post_script_path` - The path of the thaw script executed after creating application-consistent snapshots. Valid only when `app_consistent` is `true`.
      * `enable_fs_freeze` - Whether to use the Linux FsFreeze mechanism to keep the file system read-only consistent before creating application-consistent snapshots. Valid only when `app_consistent` is `true`. Default value: `true`.
      * `timeout_in_seconds` - The IO freeze timeout period. Unit: seconds. Valid only when `app_consistent` is `true`. Default value: `30`.
      * `enable_writers` - Whether to create application-consistent snapshots. `true`: creates application-consistent snapshots. `false`: creates file-system-consistent snapshots. Valid only when `app_consistent` is `true`. Default value: `true`.

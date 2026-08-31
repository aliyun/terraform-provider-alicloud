---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_policy_binding"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Policy Binding resource.
---

# alicloud_hbr_policy_binding

Provides a Hybrid Backup Recovery (HBR) Policy Binding resource.

A policy binding relationship consists of a data source, a policy, and binding options.

For information about Hybrid Backup Recovery (HBR) Policy Binding and how to use it, see [What is Policy Binding](https://www.alibabacloud.com/help/en/cloud-backup/developer-reference/api-hbr-2017-09-08-createpolicybindings).

-> **NOTE:** Available since v1.221.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_policy_binding&exampleId=75afb166-ee6a-a9ce-4a6b-11e459b0d21555d1a658&activeTab=example&spm=docs.r.hbr_policy_binding.0.75afb166ee&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_hbr_vault" "defaultyk84Hc" {
  vault_type = "STANDARD"
  vault_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = "example-value-${random_integer.default.result}"
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.defaultyk84Hc.id
  }
  policy_description = "policy example"
}

resource "alicloud_oss_bucket" "defaultKtt2XY" {
  storage_class = "Standard"
  bucket        = "example-value-${random_integer.default.result}"
}

resource "alicloud_hbr_policy_binding" "default" {
  source_type                = "OSS"
  disabled                   = "false"
  policy_id                  = alicloud_hbr_policy.defaultoqWvHQ.id
  data_source_id             = alicloud_oss_bucket.defaultKtt2XY.bucket
  policy_binding_description = "policy binding example (update)"
  source                     = "prefix-example-update/"
}
```


ECS Instance Backup With App-Consistent Snapshot Group

This example migrates an `alicloud_hbr_server_backup_plan` configuration (deprecated since v1.249.0) to `alicloud_hbr_policy_binding` using `alicloud_hbr_policy` + `advanced_options.udm_detail` with `app_consistent` and `snapshot_group`.

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_hbr_vault" "default" {
  vault_type = "STANDARD"
  vault_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_hbr_policy" "default" {
  policy_name = "example-value-${random_integer.default.result}"
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.default.id
  }
  policy_description = "policy example"
}

resource "alicloud_instance" "default" {
  instance_name = "example-value-${random_integer.default.result}"
  instance_type = "ecs.g7.large"
  image_id      = "aliyun_2_1903_x64_7h_cor_4.0.40_alibase"
  system_disk {
    category = "cloud_essd"
    size     = "40"
  }
}

resource "alicloud_hbr_policy_binding" "default" {
  source_type    = "UDM_ECS"
  policy_id      = alicloud_hbr_policy.default.id
  data_source_id = alicloud_instance.default.id
  disabled       = "false"
  advanced_options {
    udm_detail {
      app_consistent     = true
      snapshot_group     = true
      ram_role_name      = "AliyunECSBackupRole"
      pre_script_path    = "/opt/prescript.sh"
      post_script_path   = "/opt/postscript.sh"
      enable_fs_freeze   = true
      timeout_in_seconds = 60
      enable_writers     = true
    }
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_hbr_policy_binding&spm=docs.r.hbr_policy_binding.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `advanced_options` - (Optional, ForceNew, Computed, Set) Backup Advanced Options See [`advanced_options`](#advanced_options) below.
* `cross_account_role_name` - (Optional, ForceNew, Available since v1.230.0) Valid only when CrossAccountType = CROSS_ACCOUNT, indicating the name of the cross-account authorization role of the data source, and the management account uses this role to access the data source.
* `cross_account_type` - (Optional, ForceNew, Computed, Available since v1.230.0) Cross-account type, supported
* `cross_account_user_id` - (Optional, ForceNew, Int, Available since v1.230.0) Valid only when CrossAccountType = CROSS_ACCOUNT, indicating the ID of the actual account to which the data source belongs.
* `data_source_id` - (Optional, ForceNew, Computed) The data source ID.
* `disabled` - (Optional) Whether the policy is effective for the data source.
  - true: Pause
  - false: not paused
* `exclude` - (Optional) This parameter is required only when the value of SourceType is ECS_FILE or File. Indicates a file type that does not need to be backed up. All files of this type are not backed up. A maximum of 255 characters is supported.
* `include` - (Optional) This parameter is required only when the value of SourceType is ECS_FILE or File. Indicates the file types to be backed up, and all files of these types are backed up. A maximum of 255 characters is supported.
* `policy_binding_description` - (Optional) Resource Description
* `policy_id` - (Optional, ForceNew, Computed) The policy ID.
* `source` - (Optional) When SourceType is OSS, a prefix is specified to be backed up. If it is not specified, the entire root directory of the Bucket is backed up.
* `source_type` - (Optional, ForceNew, Computed, Available since v1.260.1) Data source type, value range:
  - `UDM_ECS`: indicates the ECS instance backup.
  - `OSS`: indicates an OSS backup.
  - `NAS`: indicates an Alibaba Cloud NAS Backup. When you bind a file system to a policy, Cloud Backup automatically creates a mount point for the file system. If you no longer need the mount point, delete it manually.
  - `ECS_FILE`: indicates that the ECS file is backed up.
  - `File`: indicates a local File backup.
  - `OTS`: indicates the Tablestore backup.
* `speed_limit` - (Optional) This parameter is required only when the value of SourceType is ECS_FILE or File. Indicates backup flow control. The format is {start}{end}{bandwidth}. Multiple flow control configurations use partitioning, and no overlap in configuration time is allowed. start: start hour. end: end of hour. bandwidth: limit rate, in KB/s.

### `advanced_options`

The advanced_options supports the following:
* `oss_detail` - (Optional, Set, Available since v1.273.0) OSS Backup Advanced options See [`oss_detail`](#advanced_options-oss_detail) below.
* `udm_detail` - (Optional, ForceNew, Computed, Set) ECS Backup Advanced options See [`udm_detail`](#advanced_options-udm_detail) below.

### `advanced_options-oss_detail`

The advanced_options-oss_detail supports the following:
* `ignore_archive_object` - (Optional, Available since v1.273.0) Archived objects are not prompted in task statistics and failed file lists
* `inventory_cleanup_policy` - (Optional, Available since v1.273.0) Whether to delete the inventory file after the backup. Valid only when using the OSS inventory. Supported: NO_CLEANUP: Do not delete. DELETE_CURRENT: Deletes the current file. DELETE_CURRENT_AND_PREVIOUS: Deletes all files.
* `inventory_id` - (Optional, Available since v1.273.0) The name of the OSS inventory. If the value is not empty, the OSS inventory will be used for performance tuning. We recommend that you use a list to improve incremental performance when backing up more than 0.1 billion OSS objects. OSS charges the storage fee for the list file separately. It takes time to generate the OSS inventory file. The backup may fail before the OSS inventory file is generated. You can wait for the next cycle.

### `advanced_options-udm_detail`

The advanced_options-udm_detail supports the following:
* `app_consistent` - (Optional) Whether to enable application-consistent backup. When enabled, the system uses a snapshot group together with pre/post scripts to guarantee application data consistency. Only supported when all cloud disk types of the instance are ESSD.
* `destination_kms_key_id` - (Optional) Custom KMS key ID of encrypted copy
* `disk_id_list` - (Optional, List) The list of backup disks. If it is empty, all disks are backed up.
* `enable_fs_freeze` - (Optional) Whether to enable file system freeze before taking a snapshot.
* `enable_writers` - (Optional) Whether to enable VSS writers.
* `exclude_disk_id_list` - (Optional, List) List of cloud disk IDs that are not backed up
* `post_script_path` - (Optional) The path of the post-backup script, executed after the snapshot is taken. Required when `app_consistent` is `true`.
* `pre_script_path` - (Optional) The path of the pre-backup script, executed before the snapshot is taken. Required when `app_consistent` is `true`.
* `ram_role_name` - (Optional) The RAM role name used by ECS to run the pre/post scripts. Required when `app_consistent` is `true`.
* `snapshot_group` - (Optional) Whether to use a snapshot group. Valid when `app_consistent` is `true`.
* `timeout_in_seconds` - (Optional) The timeout in seconds for the pre/post script execution.

-> **NOTE:** `app_consistent`, `snapshot_group`, `ram_role_name`, `pre_script_path`, `post_script_path`, `enable_fs_freeze`, `timeout_in_seconds` and `enable_writers` are only supported when `source_type` is `UDM_ECS`. When `app_consistent` is set to `true`, `ram_role_name`, `pre_script_path` and `post_script_path` are required by the [CreatePolicyBindings API](https://www.alibabacloud.com/help/en/cloud-backup/developer-reference/api-hbr-2017-09-08-createpolicybindings). This set of options supersedes the deprecated `alicloud_hbr_server_backup_plan` `detail` block.

-> **NOTE:** Application-consistent backup relies on [ECS snapshot consistency groups](https://help.aliyun.com/zh/ecs/user-guide/create-a-snapshot-consistency-group), which have the following runtime constraints: every disk attached to the instance must be an ESSD; multi-attach (`Multi-Attach`) shared disks are not supported; the disks in a consistency group must be in the same zone (the group can span multiple instances); a single consistency group contains at most 128 disks and at most 256 TiB of total capacity; and snapshots produced by a consistency group are retained permanently by default and are not subject to the backup policy retention period. When [Cloud Backup](https://help.aliyun.com/zh/cloud-backup/user-guide/back-up-ecs-instances) runs a whole-instance backup, it creates a consistency group if the instance supports application-consistent backup; otherwise it falls back to per-disk crash-consistent snapshots.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy Binding.
* `delete` - (Defaults to 5 mins) Used when delete the Policy Binding.
* `update` - (Defaults to 5 mins) Used when update the Policy Binding.

## Import

Hybrid Backup Recovery (HBR) Policy Binding can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_policy_binding.example <policy_id>:<source_type>:<data_source_id>
```
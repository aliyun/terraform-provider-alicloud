---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_policy_binding"
description: |-
  Provides a Alicloud HBR Policy Binding resource.
---

# alicloud_hbr_policy_binding

Provides a HBR Policy Binding resource.

For information about HBR Policy Binding and how to use it, see [What is Policy Binding](https://www.alibabacloud.com/help/en/cloud-backup/developer-reference/api-hbr-2017-09-08-createpolicybindings).

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

## Argument Reference

The following arguments are supported:
* `advanced_options` - (Optional, ForceNew, List) Backup Advanced Options See [`advanced_options`](#advanced_options) below.
* `cross_account_role_name` - (Optional, ForceNew, Available since v1.230.0) Valid only when CrossAccountType = CROSS_ACCOUNT, indicating the name of the cross-account authorization role of the data source, and the management account uses this role to access the data source.
* `cross_account_type` - (Optional, ForceNew, Available since v1.230.0) Cross-account type, supported
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
* `source_type` - (Optional, ForceNew, Computed) Data source type, value range:
  - `UDM_ECS`: indicates the ECS instance backup.
  - `OSS`: indicates an OSS backup.
  - `NAS`: indicates an Alibaba Cloud NAS Backup. When you bind a file system to a policy, Cloud Backup automatically creates a mount point for the file system. If you no longer need the mount point, delete it manually.
  - `ECS_FILE`: indicates that the ECS file is backed up.
  - `File`: indicates a local File backup.
* `speed_limit` - (Optional) This parameter is required only when the value of SourceType is ECS_FILE or File. Indicates backup flow control. The format is {start}{end}{bandwidth}. Multiple flow control configurations use partitioning, and no overlap in configuration time is allowed. start: start hour. end: end of hour. bandwidth: limit rate, in KB/s.

### `advanced_options`

The advanced_options supports the following:
* `udm_detail` - (Optional, ForceNew, List) ECS Backup Advanced options. See [`udm_detail`](#advanced_options-udm_detail) below.

### `advanced_options-udm_detail`

The advanced_options-udm_detail supports the following:
* `destination_kms_key_id` - (Optional) Custom KMS key ID of encrypted copy.
* `disk_id_list` - (Optional, List) The list of backup disks. If it is empty, all disks are backed up.
* `exclude_disk_id_list` - (Optional, List) List of cloud disk IDs that are not backed up.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<policy_id>:<source_type>:<data_source_id>`.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy Binding.
* `delete` - (Defaults to 5 mins) Used when delete the Policy Binding.
* `update` - (Defaults to 5 mins) Used when update the Policy Binding.

## Import

HBR Policy Binding can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_policy_binding.example <policy_id>:<source_type>:<data_source_id>
```
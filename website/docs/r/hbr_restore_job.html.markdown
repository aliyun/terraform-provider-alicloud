---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_restore_job"
sidebar_current: "docs-alicloud-resource-hbr-restore-job"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Restore Job resource.
---

# alicloud\_hbr\_restore\_job

Provides a Hybrid Backup Recovery (HBR) Restore Job resource.

For information about Hybrid Backup Recovery (HBR) Restore Job and how to use it, see [What is Restore Job](https://www.alibabacloud.com/help/doc-detail/186575.htm).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_restore_job&exampleId=68841786-d5c7-0c64-5689-df09846ab9dfe76c6065&activeTab=example&spm=docs.r.hbr_restore_job.0.68841786d5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_hbr_ecs_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_oss_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_nas_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_snapshots" "ecs_snapshots" {
  source_type = "ECS_FILE"
  vault_id    = data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
  instance_id = data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
}

data "alicloud_hbr_snapshots" "oss_snapshots" {
  source_type = "OSS"
  vault_id    = data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
  bucket      = data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
  source_type    = "NAS"
  vault_id       = data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
  file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
  create_time    = data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
}

resource "alicloud_hbr_restore_job" "nasJob" {
  snapshot_hash         = data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash
  vault_id              = data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
  source_type           = "NAS"
  restore_type          = "NAS"
  snapshot_id           = data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id
  target_file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
  target_create_time    = data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
  target_path           = "/"
  options               = <<EOF
    {"includes":[], "excludes":[]}
  EOF
}

resource "alicloud_hbr_restore_job" "ossJob" {
  snapshot_hash = data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_hash
  vault_id      = data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
  source_type   = "OSS"
  restore_type  = "OSS"
  snapshot_id   = data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_id
  target_bucket = data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
  target_prefix = ""
  options       = <<EOF
    {"includes":[], "excludes":[]}
  EOF
}

resource "alicloud_hbr_restore_job" "ecsJob" {
  snapshot_hash      = data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash
  vault_id           = data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
  source_type        = "ECS_FILE"
  restore_type       = "ECS_FILE"
  snapshot_id        = data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id
  target_instance_id = data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
  target_path        = "/"
}
```

-> **NOTE:** This resource can only be created, cannot be modified or deleted. Therefore, any modification of the resource attribute will not affect exist resource.

## Argument Reference

The following arguments are supported:

* `restore_job_id` - (Optional, Computed, ForceNew) Restore Job ID. It's the unique key of this resource, if you want to set this argument by yourself, you must specify a unique keyword that never appears.
* `vault_id` - (Required, ForceNew) The ID of backup vault.
* `source_type` - (Required, ForceNew) The type of data source. Valid values: `ECS_FILE`, `NAS`, `OSS`,`OTS_TABLE`,`UDM_ECS`.
* `restore_type` - (Required, ForceNew) The type of recovery destination. Valid values: `ECS_FILE`, `NAS`, `OSS`,`OTS_TABLE`,`UDM_ECS_ROLLBACK`. **Note**: Currently, there is a one-to-one correspondence between the data source type with the recovery destination type.
* `snapshot_id` - (Required, ForceNew) The ID of Snapshot.
* `snapshot_hash` - (Required, ForceNew) The hashcode of Snapshot.
* `options` - (Optional, ForceNew) Recovery options. **NOTE:** Required while source_type equals `OSS` or `NAS`, invalid while source_type equals `ECS_FILE`. It's a json string with format:`"{"includes":[],"excludes":[]}",`. Recovery options. When restores OTS_TABLE and real target time is the rangEnd time of the snapshot, it should be a string with format: `{"UI_TargetTime":1650032529018}`.
* `exclude` - (Optional) The exclude path. **NOTE:** Invalid while source_type equals `OSS` or `NAS`. It's a json string with format:`["/excludePath]`, up to 255 characters. **WARNING:** If this value filled in incorrectly, the task may not start correctly, so please check the parameters before executing the plan.
* `include` - (Optional) The include path. **NOTE:** Invalid while source_type equals `OSS` or `NAS`. It's a json string with format:`["/includePath"]`, Up to 255 characters. **WARNING:** The field is required while source_type equals `OTS_TABLE` which means source table name. If this value filled in incorrectly, the task may not start correctly, so please check the parameters before executing the plan. 
* `target_bucket` - (Optional, ForceNew) The target name of OSS bucket. **NOTE:** Required while source_type equals `OSS`,
* `target_prefix` - (Optional, ForceNew) The target prefix of the OSS object. **WARNING:** Required while source_type equals `OSS`. If this value filled in incorrectly, the task may not start correctly, so please check the parameters before executing the plan.
* `target_file_system_id` - (Optional, ForceNew) The ID of destination File System. **NOTE:** Required while source_type equals `NAS`
* `target_create_time` - (Optional, ForceNew) The creation time of destination File System. **NOTE:** While source_type equals `NAS`, this parameter must be set. **Note:** The time format of the API adopts the ISO 8601 format, such as `2021-07-09T15:45:30CST` or `2021-07-09T07:45:30Z`.
* `target_path` - (Optional, ForceNew) The target file path of (ECS) instance. **WARNING:** Required while source_type equals `NAS` or `ECS_FILE`, If this value filled in incorrectly, the task may not start correctly, so please check the parameters before executing the plan.
* `target_instance_id` - (Optional, ForceNew)  The target ID of ECS instance. **NOTE:** Required while source_type equals `ECS_FILE`
* `target_client_id` - (Optional, ForceNew) The target client ID.
* `target_data_source_id` - (Optional, ForceNew) The target data source ID.
* `target_time` - (Optional, Available in v1.164.0) The time when data is restored to the Table store instance. This value is a UNIX timestamp. Unit: seconds. **WARNING:** Required while source_type equals `OTS_TABLE`. **Note:** The time when data is restored to the Tablestore instance. It should be 0 if restores data at the End time of the snapshot.
* `udm_detail` - (Optional, Available in v1.164.0) The full machine backup details.
* `target_instance_name` - (Optional, Available in v1.164.0) The name of the Table store instance to which you want to restore data.**WARNING:** Required while source_type equals `OTS_TABLE`.
* `target_table_name` - (Optional, Available in v1.164.0) The name of the table that stores the restored data. **WARNING:** Required while source_type equals `OTS_TABLE`.
* `ots_detail` - (Optional, Computed, Available in v1.186.0) The details about the Tablestore instance. See the following `Block ots_detail`.
* `cross_account_type` - (Optional, ForceNew, Computed, Available in v1.189.0+) The type of the cross account backup. Valid values: `SELF_ACCOUNT`, `CROSS_ACCOUNT`.
* `cross_account_user_id` - (Optional, ForceNew, Available in v1.189.0+) The original account ID of the cross account backup managed by the current account.
* `cross_account_role_name` - (Optional, ForceNew, Available in v1.189.0+) The role name created in the original account RAM backup by the cross account managed by the current account.

#### Block ots_detail

The ots_detail supports the following:
* `overwrite_existing` - (Optional, ForceNew, Computed) Whether to overwrite the existing table storage recovery task. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Restore Job. The value formats as `<restore_job_id>:<restore_type>`.
* `status` - The Restore Job Status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Restore Job.

## Import

Hybrid Backup Recovery (HBR) Restore Job can be imported using the id. Format to `<restore_job_id>:<restore_type>`, e.g.

```shell
$ terraform import alicloud_hbr_restore_job.example your_restore_job_id:your_restore_type
```

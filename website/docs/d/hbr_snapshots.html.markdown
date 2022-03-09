---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_snapshots"
sidebar_current: "docs-alicloud-datasource-hbr-snapshots"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Snapshots to the user.
---

# alicloud\_hbr\_snapshots

This data source provides the Hbr Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

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
  source_type           = "OSS"
  vault_id              = data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
  bucket                = data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
  complete_time         = "2021-07-20T14:17:15CST,2021-07-24T14:17:15CST"
  complete_time_checker = "BETWEEN"
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
  source_type           = "NAS"
  vault_id              = data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
  file_system_id        = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
  create_time           = data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
  complete_time         = "2021-08-23T14:17:15CST"
  complete_time_checker = "GREATER_THAN_OR_EQUAL"
}

output "hbr_snapshot_id_1" {
  value = data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `vault_id` - (Required, ForceNew) The ID of Vault.
* `source_type` - (Required, ForceNew) Data source type, valid values: `ECS_FILE`, `OSS`, `NAS`.
* `status` - (Optional, ForceNew) The status of snapshot, valid values: `COMPLETE`, `PARTIAL_COMPLETE`.
* `instance_id` - (Optional, ForceNew) The ID of ECS instance. The ecs backup client must have been installed on the host. While source_type equals `ECS_FILE`, this parameter must be set.
* `bucket` - (Optional, ForceNew) The bucket name of OSS. While source_type equals `OSS`, this parameter must be set.
* `file_system_id` - (Optional, ForceNew) The File System ID of Nas. While source_type equals `NAS`, this parameter must be set.
* `create_time` - (Optional, ForceNew) File system creation timestamp of Nas. While source_type equals `NAS`, this parameter must be set. **Note** The time format of the API adopts the ISO 8601 format, such as `2021-07-09T15:45:30CST` or `2021-07-09T07:45:30Z`.
* `complete_time` - (Optional, ForceNew) Timestamp of Snapshot completion. Note The time format of the API adopts the ISO 8601 format, such as 2021-07-09T15:45:30CST or 2021-07-09T07:45:30Z. **Note**: While `complete_time_checker` equals `BETWEEN`, this field should be formatted such as `"2021-08-20T14:17:15CST,2021-08-26T14:17:15CST"`, The first part of this string is the start time, the second part is the end time, and the two parts should be separated by commas.
* `complete_time_checker` - (Optional, ForceNew) Complete time filter operator. Optional values: `MATCH_TERM`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL`, `LESS_THAN`, `LESS_THAN_OR_EQUAL`, `BETWEEN`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `snapshots` - A list of Hbr Snapshots. Each element contains the following attributes:
	* `id` - The ID of the Snapshot.
	* `snapshot_id` - The ID of the Snapshot.
	* `snapshot_hash` - The hashcode of Snapshot.
	* `job_id` - The job ID of backup task.
	* `source_type` - Data source type, optional values: `ECS_FILE`, `OSS`, `NAS`.
	* `backup_type` - Backup type. Possible values: `COMPLETE` (full backup).
	* `actual_bytes` - The actual data volume of the snapshot. Unit byte.
	* `actual_items` - The actual number of items in the snapshot. (Currently only file backup is available).
	* `bytes_done` - The incremental amount of backup data. Unit byte.
	* `bytes_total` - The total amount of data sources. Unit byte.
	* `created_time` - Snapshot creation time. UNIX time in seconds.
	* `start_time` - The start time of the snapshot. UNIX time in seconds.
	* `updated_time` - The update time of snapshot. UNIX time in seconds.
	* `complete_time` - The time when the snapshot completed. UNIX time in seconds.
	* `instance_id` - (ECS_FILE) The ID of ECS instance.
	* `client_id` - (ECS_FILE) The ID of ECS backup client.
	* `bucket` - (OSS) The name of OSS bucket.
	* `file_system_id` - (NAS) The ID of NAS File system.
	* `create_time` - (NAS) File System Creation Time of Nas. Unix Time Seconds.
	* `path` - (ECS_FILE, NAS) Backup Path.
	* `prefix` - (OSS) Backup file prefix.
	* `items_done` - The number of backup items. (Currently only file backup is available).
	* `items_total` - The total number of data source items. (Currently only file backup is available).
	* `parent_snapshot_hash` - The hashcode of parent backup snapshot.
	* `retention` - The number of days to keep.
	* `status` - The status of snapshot execution. Possible values: `COMPLETE`, `PARTIAL_COMPLETE`, `FAILED`.

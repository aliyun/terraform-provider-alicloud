---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_restore_jobs"
sidebar_current: "docs-alicloud-datasource-hbr-restore-jobs"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Restore Jobs to the user.
---

# alicloud\_hbr\_restore\_jobs

This data source provides the Hbr Restore Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ecs_backup_plans" "default" {
  name_regex = "plan-name"
}

data "alicloud_hbr_restore_jobs" "default" {
  restore_type       = "ECS_FILE"
  vault_id           = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  target_instance_id = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
}
```

## Argument Reference

The following arguments are supported:

* `restore_type` - (Required, ForceNew) The Recovery Destination Types. Valid values: `ECS_FILE`, `NAS`, `OSS`,`OTS_TABLE`, `UDM_ECS_ROLLBACK`.
* `restore_id` - (Optional, ForceNew) The list of restore job IDs.
* `vault_id` - (Optional, ForceNew) The list of backup vault IDs.
* `source_type` - (Optional, ForceNew) The list of data source types. Valid values: `ECS_FILE`, `NAS`, `OSS`, `OTS_TABLE`,`UDM_ECS_ROLLBACK`.
* `status` - (Optional, ForceNew) The status of restore job. Valid values: `CANCELED`, `CANCELING`, `COMPLETE`, `CREATED`, `EXPIRED`, `FAILED`, `PARTIAL_COMPLETE`, `QUEUED`, `RUNNING`.
* `target_bucket` - (Optional, ForceNew) The name of target OSS bucket.
* `target_file_system_id` - (Optional, ForceNew) Valid while source_type equals `NAS`. The list of destination File System IDs.
* `target_instance_id` - (Optional, ForceNew) The ID of target ECS instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of restore job IDs.
* `jobs` - A list of Hbr restore jobs. Each element contains the following attributes:
	* `id` - The ID of the restore job.
	* `vault_id` - The ID of backup vault.
	* `restore_job_id` - The ID of restore job.
	* `snapshot_id` - The ID of Snapshot.
	* `snapshot_hash` - The hashcode of Snapshot.
	* `restore_type` - The type of recovery destination. Valid Values: `ECS_FILE`, `OSS`, `NAS`.
	* `create_time` - The creation time of restore job.
	* `start_time` - The start time of restore job. Unix Time in Seconds.
	* `updated_time` - The update Time of restore job. Unix Time in Seconds.
	* `complete_time` - The completion time of restore Job.
	* `target_bucket` - The name of target ofo OSS bucket.
	* `target_create_time` - The creation time of destination file system.
	* `target_file_system_id` - The ID of destination file system.
	* `target_instance_id` - The ID of target ECS instance.
	* `target_path` - The target file path of ECS instance.
	* `target_prefix` - The file prefix of target OSS object.
	* `status` - The status of restore job.
	* `actual_bytes` - The actual size of Snapshot.
	* `actual_items` - The actual number of files.
	* `bytes_done` - The size of restore job recovered.
	* `bytes_total` - The total size of restore job recovered.
	* `error_message` - The error message of recovery task execution.
	* `expire_time` - The expiration time of restore job. Unix Time in seconds.
	* `items_done` - The number of items restore job recovered.
	* `items_total` - The total number of items restore job recovered.
	* `options` - Recovery Options.
	* `progress` - The recovery progress.

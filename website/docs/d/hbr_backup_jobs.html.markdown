---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_backup_jobs"
sidebar_current: "docs-alicloud-datasource-hbr-backup-jobs"
description: |-
  Provides a list of Hbr Backup Jobs to the user.
---

# alicloud\_hbr\_backup\_jobs

This data source provides the Hbr Backup Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ecs_backup_plans" "default" {
  name_regex = "plan-name"
}

data "alicloud_hbr_backup_jobs" "default" {
  source_type = "ECS_FILE"
  filter {
    key      = "VaultId"
    operator = "IN"
    values   = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  }
  filter {
    key      = "InstanceId"
    operator = "IN"
    values   = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
  }
  filter {
    key      = "CompleteTime"
    operator = "BETWEEN"
    values   = ["2021-08-23T14:17:15CST", "2021-08-24T14:17:15CST"]
  }
}

data "alicloud_hbr_backup_jobs" "example" {
  source_type = "ECS_FILE"
  status      = "COMPLETE"
  filter {
    key      = "VaultId"
    operator = "IN"
    values   = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  }
  filter {
    key      = "InstanceId"
    operator = "IN"
    values   = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
  }
  filter {
    key      = "CompleteTime"
    operator = "LESS_THAN"
    values   = ["2021-10-20T20:20:20CST"]
  }
}

output "alicloud_hbr_backup_jobs_default_1" {
  value = data.alicloud_hbr_backup_jobs.default.jobs.0.id
}

output "alicloud_hbr_backup_jobs_example_1" {
  value = data.alicloud_hbr_backup_jobs.example.jobs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Backup Job IDs.
* `source_type` - (Required, ForceNew) The type of data source. Valid values: `ECS_FILE`, `NAS`, `OSS`, `UDM_ECS`, `UDM_ECS_DISK`.
* `status` - (Optional, ForceNew) The status of backup job. Valid values: `COMPLETE`, `PARTIAL_COMPLETE`, `FAILED`, `UNAVAILABLE`.
* `sort_direction` - (Optional, ForceNew) The sort direction, sort results by ascending or descending order based on the value jobs id. Valid values: `ASCEND`, `DESCEND`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


### Block filter

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:
* `key`      - (Required) The key of the field to filter. Valid values: `PlanId`, `VaultId`, `InstanceId`, `Bucket`, `FileSystemId`, `CompleteTime`.
* `operator` - (Required) The operator of the field to filter. Valid values: `EQUAL`, `NOT_EQUAL`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL`, `LESS_THAN`, `LESS_THAN_OR_EQUAL`, `BETWEEN`, `IN`.
* `values`   - (Required) Set of values that are accepted for the given field.

-> **NOTE:** Numeric types such as `CompleteTime` do not support `IN` operations for the time being.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Hbr Backup Jobs. Each element contains the following attributes:
	* `id` - The ID of the backup job.
	* `backup_job_id` - The ID of the backup job.
	* `source_type` - The type of data source. Valid Values: `ECS_FILE`, `OSS`, `NAS`, `UDM_DISK`.
	* `back_job_name` - The name of backup job.
	* `plan_id` - The ID of a backup plan.
	* `vault_id` - The ID of backup vault.
	* `actual_bytes` - The actual data volume of the backup task (After deduplication) . Unit byte.
	* `actual_items` - The actual number of items in the backup task. (Currently only file backup is available).
	* `backup_type` - Backup type. Valid values: `COMPLETE`(full backup).
	* `bucket` - The name of target OSS bucket.
	* `bytes_done` - The amount of backup data (Incremental). Unit byte.
	* `bytes_total` - The total amount of data sources. Unit byte.
	* `create_time` - The creation time of backup job. UNIX time seconds.
	* `start_time` - The scheduled backup start time. UNIX time seconds.
	* `complete_time` -  The completion time of backup job. UNIX time seconds.
	* `updated_time` - The update time of backup job. UNIX time seconds.
	* `file_system_id` - The ID of destination file system.
	* `nas_create_time` - File system creation time. UNIX time in seconds.
	* `instance_id` - The ID of target ECS instance.
	* `items_done` - The number of items restore job recovered.
	* `items_total` - The total number of items restore job recovered.
	* `paths` - List of backup path. e.g. `["/home", "/var"]`.
	* `prefix` - The prefix of Oss bucket files.
	* `include` - Include path. String of Json list. Up to 255 characters. e.g. `"[\"/var\"]"`
	* `exclude` - Exclude path. String of Json list. Up to 255 characters. e.g. `"[\"/home/work\"]"`
	* `status` - The status of restore job. Valid values: `COMPLETE` , `PARTIAL_COMPLETE`, `FAILED`.
	* `progress` - Backup progress. The value is 100%*100.
	* `error_message` - Error message.
	

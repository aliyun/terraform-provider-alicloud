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
    key       = "VaultId"
    operator  = "IN"
    values    = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  }
  filter {
    key       = "InstanceId"
    operator  = "IN"
    values    = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
  }
  filter {
    key       = "CompleteTime"
    operator  = "BETWEEN"
    values    = ["2019-08-23T14:17:15CST", "2020-08-23T14:17:15CST"]
  }
}

data "alicloud_hbr_backup_jobs" "example" {
  source_type = "ECS_FILE"
  status      = "COMPLETE"
  filter {
    key       = "VaultId"
    operator  = "IN"
    values    = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  }
  filter {
    key       = "InstanceId"
    operator  = "IN"
    values    = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
  }
  filter {
    key       = "CompleteTime"
    operator  = "LESS_THAN"
    values    = ["2021-08-20T20:20:20CST"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Backup Job IDs.
* `source_type` - (Required, ForceNew) The type of data source. Valid values: `ECS_FILE`, `NAS`, `OSS`.
* `status` - (Optional, ForceNew) The status of backup job. Valid values: `COMPLETE`, `PARTIAL_COMPLETE`, `FAILED`.
* `sort_direction` - (Optional, ForceNew) The sort direction.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


### Block filter

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:
* `key`      - (Required) The key of the field to filter. Valid values: `PlanId`, `VaultId`, `InstanceId`, `Bucket`, `FileSystemId`, `CompleteTime`.
* `operator` - (Required) The operator of the field to filter. Valid values: `MATCH_TERM`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL`, `LESS_THAN`, `LESS_THAN_OR_EQUAL`, `BETWEEN`.
* `values`   - (Required) Set of values that are accepted for the given field.


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Hbr Backup Jobs. Each element contains the following attributes:
	* `id` - The ID of the backup job.
	* `source_type` - The type of data source. Valid Values: `ECS_FILE`, `OSS`, `NAS`.
	* `backup_job_id` - The ID of backup job.
	* `back_job_name` - The name of backup job.
	* `plan_id` - The IF of a backup plan.
	* `vault_id` - The ID of backup vault.
	* `actual_bytes` - The actual size of backup job.
	* `actual_items` - The actual number of files.
	* `backup_type` - Backup type. Valid values: `COMPLETE`(full backup).
	* `bucket` - The name of target ofo OSS bucket.
	* `bytes_done` - The size of backup job recovered.
	* `bytes_total` - The total size of backup job recovered.
	* `create_time` - The creation time of backup job. UNIX time seconds.
	* `start_time` - The scheduled backup start time. UNIX time seconds.
	* `complete_time` -  The completion time of backup job. UNIX time seconds.
	* `updated_time` - The update time of backup job. UNIX time seconds.
	* `file_system_id` - The ID of destination file system.
	* `nas_create_time` - File system creation time. UNIX time in seconds.
	* `instance_id` - The ID of target ECS instance.
	* `items_done` - The number of items restore job recovered.
	* `items_total` - The total number of items restore job recovered.
	* `paths` - Backup path. e.g. `["/home", "/var"]`
	* `prefix` - The prefix of Oss bucket files.
	* `include` - Include path. String of Json list. Up to 255 characters. e.g. `"[\"/var\"]"`
	* `exclude` - Exclude path. String of Json list. Up to 255 characters. e.g. `"[\"/home/work\"]"`
	* `status` - The status of restore job. Valid values: `COMPLETE` , `PARTIAL_COMPLETE`, `FAILED`.
	
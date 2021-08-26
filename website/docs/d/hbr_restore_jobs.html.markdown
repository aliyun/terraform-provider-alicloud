---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_restore_jobs"
sidebar_current: "docs-alicloud-datasource-hbr-restore-jobs"
description: |-
  Provides a list of Hbr Restore Jobs to the user.
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
  restore_type =       "ECS_FILE"
  vault_id =           [data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id]
  target_instance_id = [data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, ForceNew) The filters.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `restore_id` - (Optional, ForceNew) The restore id.
* `restore_type` - (Required, ForceNew) The Recovery Destination Types. Valid Values: ECS_FILE OSS NAS. Valid values: `ECS_FILE`, `NAS`, `OSS`.
* `source_type` - (Optional, ForceNew) The Type of Data Source. Valid Values: ECS_FILE OSS NAS. Valid values: `ECS_FILE`, `NAS`, `OSS`.
* `status` - (Optional, ForceNew) The Restore Job Status.
* `target_bucket` - (Optional, ForceNew) The Target ofo OSS Bucket Name.
* `target_file_system_id` - (Optional, ForceNew) The Destination File System ID.
* `target_instance_id` - (Optional, ForceNew) Objective to ECS Instance Id.
* `vault_id` - (Optional, ForceNew) The ID of Vault.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Restore Job IDs.
* `jobs` - A list of Hbr Restore Jobs. Each element contains the following attributes:
	* `id` - The ID of the Restore Job.
	* `vault_id` - The ID of Vault.
	* `restore_job_id` - Restore Job ID.
	* `snapshot_id` - The ID of Snapshot.
	* `snapshot_hash` - Restore Snapshot of HashCode.
	* `restore_type` - The Recovery Destination Types. Valid Values: `ECS_FILE`, `OSS`, `NAS`.
	* `source_type` - The Type of Data Source. Valid Values: `ECS_FILE`, `OSS`, `NAS`.
	* `create_time` - The Restore Job Creation Time.
	* `start_time` - Restoring the Start Time. Unix Time in Seconds.
	* `updated_time` - Update Time.
	* `complete_time` - Restore Completion Time.
	* `target_bucket` - The Target ofo OSS Bucket Name.
	* `target_client_id` - The ID of Target Client.
	* `target_create_time` - The Destination File System Creation Time.
	* `target_data_source_id` - The Destination ID.
	* `target_file_system_id` - The Destination File System ID.
	* `target_instance_id` - Objective to ECS Instance Id.
	* `target_path` - The Target of (ECS) Instance Changes the ECS File Path.
	* `target_prefix` - The Target of the OSS Object Prefix.
	* `status` - The Restore Job Status.
	* `actual_bytes` - The Actual Size of Snapshot.
	* `actual_items` - The Actual Number of Files.
	* `bytes_done` - Recovery Is Successful, Size.
	* `bytes_total` - The Restored Total.
	* `error_message` - The Recovery Task Execution Error Message.
	* `expire_time` - Restore the Expiration Time. Unix Time in Seconds.
	* `items_done` - Log of Files Successfully Recovered the Number.
	* `items_total` - File the Total Number.
	* `options` - Recovery Options.
	* `parent_id` - The Parent Node.
	* `progress` - The Recovery Progress 100% * 100.
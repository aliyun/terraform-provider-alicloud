---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_nas_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-nas_backup_plans"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) NasBackupPlans to the user.
---

# alicloud\_hbr\_nas\_backup\_plans

This data source provides the Hbr NasBackupPlans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_nas_backup_plans" "ids" {
  name_regex = "^my-NasBackupPlan"
}

output "hbr_nas_backup_plan_id" {
  value = data.alicloud_hbr_nas_backup_plans.ids.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of NasBackupPlan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by NasBackupPlan name.
* `file_system_id` - (Optional, ForceNew) The Nas fileSystem instance ID of the EcsBackupPlan used.
* `vault_id` - (Optional, ForceNew) The backup vault ID of the NasBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of NasBackupPlan names.
* `plans` - A list of Hbr NasBackupPlans. Each element contains the following attributes:
    * `id` - The ID of Nas backup plan.
    * `vault_id` - The ID of backup vault.
    * `nas_backup_plan_id` - The ID of Nas backup plan.
    * `nas_backup_plan_name` - The name of the resource.
    * `retention` - Backup retention days, the minimum is 1.
    * `schedule` - Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
        * `startTime` Backup start time, UNIX time seconds.
        * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
    * `file_system_id` - The File System ID of Nas.
    * `create_time` - File System Creation Time. **Note** The time format of the API adopts the ISO 8601 format, such as `2021-07-09T15:45:30CST`.
    * `path` - List of backup path. Up to 65536 Characters. e.g.`["/home", "/var"]`
    * `backup_type` - Backup type. Valid values: `COMPLETE`.
    * `disabled` - Whether to be suspended. Valid values: `true`, `false`.
    * `options` - Options field of Nas backup plan.
    * `created_time` - The creation time of the backup plan. UNIX time in seconds.
    * `updated_time` - The update time of the backup plan. UNIX time in seconds.


---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ots_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-ots_backup_plans"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) OtsBackupPlans to the user.
---

# alicloud\_hbr\_ots\_backup\_plans

This data source provides the Hbr OtsBackupPlans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ots_backup_plans" "ids" {
  name_regex = "^my-otsBackupPlan"
}
output "hbr_ots_backup_plan_id" {
  value = data.alicloud_hbr_ots_backup_plans.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of OtsBackupPlan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by OtsBackupPlan name.
* `plan_id` - (Optional, ForceNew) The ID of the backup plan.
* `plan_name` - (Optional, ForceNew) The ID of the backup plan.
* `vault_id` - (Optional) The ID of backup vault the OtsBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of OtsBackupPlan names.
* `plans` - A list of Hbr OtsBackupPlans. Each element contains the following attributes:
  * `vault_id` - The ID of backup vault.
  * `backup_type` - The Backup type. Valid values: `COMPLETE`.
  * `source_type` - The type of the data source.
  * `disabled` - Whether to be suspended. Valid values: `true`, `false`.
  * `retention` - The Backup retention days, the minimum is 1.
  * `created_time` - The creation time of the backup plan. UNIX time in seconds.
  * `id` - The ID of ots backup plan.
  * `ots_backup_plan_id` - The ID of ots backup plan.
  * `ots_backup_plan_name` - The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
  * `schedule` - The Backup strategy. Optional format: I|{startTime}|{interval}. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
  * `updated_time` - The update time of the backup plan. UNIX time in seconds.
  *ots_detail - The details about the Tablestore instance.



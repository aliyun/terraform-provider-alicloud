---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_oss_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-oss_backup_plans"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) OssBackupPlans to the user.
---

# alicloud\_hbr\_oss\_backup\_plans

This data source provides the Hbr OssBackupPlans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_oss_backup_plans" "ids" {
  name_regex = "^my-OssBackupPlan"
}
output "hbr_oss_backup_plan_id" {
  value = data.alicloud_hbr_oss_backup_plans.ids.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of OssBackupPlan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by OssBackupPlan name.
* `bucket` - (Required, ForceNew) The name of OSS bucket.
* `vault_id` - (Optional) The ID of backup vault the OssBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of OssBackupPlan names.
* `plans` - A list of Hbr OssBackupPlans. Each element contains the following attributes:
    * `id` - The ID of Oss backup plan.
    * `oss_backup_plan_id` - The ID of Oss backup plan.
    * `oss_backup_plan_name` - The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
    * `vault_id` - The ID of backup vault.
    * `prefix` - Backup prefix.
    * `bucket` - (Required, ForceNew) The name of OSS bucket.
    * `retention` - (Required) Backup retention days, the minimum is 1.
    * `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval}. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed yet, the next backup task will not be triggered.
        * `startTime` Backup start time, UNIX time seconds.
        * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
    * `backup_type` - (Optional, Computed, ForceNew) Backup type. Valid values: `COMPLETE`.
    * `disabled` - Whether to be suspended. Valid values: `true`, `false`.
    * `created_time` - The creation time of the backup plan. UNIX time in seconds.
    * `updated_time` - The update time of the backup plan. UNIX time in seconds.

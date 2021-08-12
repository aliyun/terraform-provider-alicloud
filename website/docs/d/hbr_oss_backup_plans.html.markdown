---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_oss_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-oss_backup_plans"
description: |-
  Provides a list of Hbr OssBackupPlans to the user.
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
* `bucket` - (Required, ForceNew) The OSS Bucket Name.
* `vault_id` - (Optional) The Vault ID of the OssBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of OssBackupPlan names.
* `plans` - A list of Hbr OssBackupPlans. Each element contains the following attributes:
    * `oss_backup_plan_name` - (Required) The Configuration Page of a Backup Plan Name. 1-64 Characters, requiring a Single Warehouse under Each of the Data Source Type Drop-down List of the Configuration Page of a Backup Plan Name Is Unique.
    * `vault_id` - (Required, ForceNew) Vault ID.
    * `bucket` - (Required, ForceNew) The OSS Bucket Name.
    * `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
    * `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval} * startTime Backup start time, UNIX time, in seconds. * interval ISO8601 time interval. E.g: ** PT1H, one hour apart. ** P1D, one day apart. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed, the next backup task will not be triggered.
    * `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.


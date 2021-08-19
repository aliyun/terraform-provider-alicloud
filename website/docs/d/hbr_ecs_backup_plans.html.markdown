---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-ecs_backup_plans"
description: |-
  Provides a list of Hbr EcsBackupPlans to the user.
---

# alicloud\_hbr\_ecs\_backup\_plans

This data source provides the Hbr EcsBackupPlans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ecs_backup_plans" "ids" {
  name_regex = "^my-EcsBackupPlan"
}

output "hbr_ecs_backup_plan_id" {
  value = data.alicloud_hbr_ecs_backup_plans.ids.plans.0.id
}           
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of EcsBackupPlan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by EcsBackupPlan name.
* `instance_id` - (Optional) The ECS instance ID of the EcsBackupPlan used.
* `vault_id` - (Optional) The Vault ID of the EcsBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of EcsBackupPlan names.
* `plans` - A list of Hbr EcsBackupPlans. Each element contains the following attributes:
	* `ecs_backup_plan_name` - (Required) The Configuration Page of a Backup Plan Name. 1-64 Characters, requiring a Single Warehouse under Each of the Data Source Type Drop-down List of the Configuration Page of a Backup Plan Name Is Unique.
	* `vault_id` - (Required, ForceNew) Vault ID.
	* `instance_id` - (Required, ForceNew) The ECS Instance Id. Must Have Installed the Client.
	* `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
	* `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval} * startTime Backup start time, UNIX time, in seconds. * interval ISO8601 time interval. E.g: ** PT1H, one hour apart. ** P1D, one day apart. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed, the next backup task will not be triggered.
	* `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.
	* `options` - (Optional) Windows System with Application Consistency Using VSS. eg: {`UseVSS`:false}.
	* `speed_limit` - (Optional) flow control. The format is: {start}|{end}|{bandwidth} * start starting hour * end end hour * bandwidth limit rate, in KiB ** Use | to separate multiple flow control configurations; ** Multiple flow control configurations are not allowed to have overlapping times.
	* `path` - (Optional) Backup Path. e.g. `["/home", "/var"]`
	* `exclude` - (Optional) Exclude Path. String of Json List, most 255 Characters. e.g. `"[\"/home/work\"]"`
	* `include` - (Optional) Include Path. String of Json List, most 255 Characters. e.g. `"[\"/var\"]"`



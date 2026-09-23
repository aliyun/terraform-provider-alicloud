---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-ecs_backup_plans"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) EcsBackupPlans to the user.
---

# alicloud\_hbr\_ecs\_backup\_plans

This data source provides the Hbr EcsBackupPlans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ecs_backup_plans" "ids" {
  name_regex = "plan-name"
}

output "hbr_ecs_backup_plan_id" {
  value = data.alicloud_hbr_ecs_backup_plans.ids.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of EcsBackupPlan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by EcsBackupPlan name.
* `instance_id` - (Optional, ForceNew) The ECS instance ID of the EcsBackupPlan used.
* `vault_id` - (Optional, ForceNew) The Vault ID of the EcsBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ecs backup plan names.
* `plans` - A list of Hbr Ecs backup plans. Each element contains the following attributes:
    * `id` - The ID of ecs backup plan.
    * `ecs_backup_plan_id` - The ID of ecs backup plan.
    * `ecs_backup_plan_name` - The name of the backup plan.
    * `vault_id` - The ID of Backup vault.
    * `instance_id` - The ID of ECS instance.
    * `retention` - Backup retention days, the minimum is 1.
    * `schedule` - Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
        * `startTime` Backup start time, UNIX time seconds.
        * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart.
    * `backup_type` - Backup type. Valid values: `COMPLETE`.
    * `source_type` - The type of backup source.
    * `options` - Windows operating system with application consistency using VSS. eg: {`UseVSS`:false}.
    * `speed_limit` - Flow control. The format is: {start}|{end}|{bandwidth}. Use `|` to separate multiple flow control configurations, multiple flow control configurations not allowed to have overlapping times.
    * `path` - Backup path. e.g. `["/home", "/var"]`
    * `exclude` - Exclude path. String of Json list. Up to 255 characters. e.g. `"[\"/home/work\"]"`
    * `include` - Include path. String of Json list. Up to 255 characters. e.g. `"[\"/var\"]"`
    * `created_time` - The creation time of the backup plan. UNIX time in seconds.
    * `updated_time` - The update time of the backup plan. UNIX time in seconds.
    * `disabled` - Whether to be suspended. Valid values: `true`, `false`.

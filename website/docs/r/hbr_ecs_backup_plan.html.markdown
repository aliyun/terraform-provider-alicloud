---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-ecs-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Ecs Backup Plan resource.
---

# alicloud\_hbr\_ecs\_backup\_plan

Provides a HBR Ecs Backup Plan resource.

For information about HBR Ecs Backup Plan and how to use it, see [What is Ecs Backup Plan](https://www.alibabacloud.com/help/doc-detail/186574.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbr_ecs_backup_plan" "example" {
  ecs_backup_plan_name = "example_value"
  instance_id          = "i-bp1567rc0oxxxxxxxxxx"
  vault_id             = "v-0003gxoksflhxxxxxxxx"
  retention            = "1"
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  speed_limit          = "0:24:5120"
  path                 = ["/home", "/var"]
  exclude              = <<EOF
  ["/home/exclude"]
  EOF
  include              = <<EOF
  ["/home/include"]
  EOF
}
```

## Argument Reference

The following arguments are support:

* `ecs_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `vault_id` - (Required, ForceNew) The ID of Backup vault.
* `instance_id` - (Required, ForceNew) The ID of ECS instance. The ecs backup client must have been installed on the host.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval}. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `1D` means one day apart. 
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `backup_type` - (Optional, Computed, ForceNew) Backup type. Valid values: `COMPLETE`.
* `options` - (Optional) Windows operating system with application consistency using VSS. eg: {`UseVSS`:false}.
* `speed_limit` - (Optional) Flow control. The format is: {start}|{end}|{bandwidth}. Use `|` to separate multiple flow control configurations, multiple flow control configurations not allowed to have overlapping times.
    * `start` starting hour 
    * `end` end hour 
    * `bandwidth` limit rate, in KiB
* `path` - (Optional) Backup path. e.g. `["/home", "/var"]`
* `exclude` - (Optional) Exclude path. String of Json list, up to 255 characters. e.g. `"[\"/home/work\"]"`
* `include` - (Optional) Include path. String of Json list, up to 255 characters. e.g. `"[\"/var\"]"`


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ecs Backup Plan.

## Import

HBR Ecs Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_ecs_backup_plan.example <id>
```
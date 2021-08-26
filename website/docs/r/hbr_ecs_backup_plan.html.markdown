---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-ecs-backup-plan"
description: |-
  Provides a Alicloud HBR Ecs Backup Plan resource.
---

# alicloud\_hbr\_ecs\_backup\_plan

Provides a HBR Ecs Backup Plan resource.

For information about HBR Ecs Backup Plan and how to use it, see [What is Ecs Backup Plan](https://www.alibabacloud.com/help/doc-detail/186568.htm).

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
  speed_limit          = "I|1602673264|PT2H"
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

The following arguments are supported:

* `ecs_backup_plan_name` - (Required) The Configuration Page of a Backup Plan Name. 1-64 Characters, requiring a Single Warehouse under Each of the Data Source Type Drop-down List of the Configuration Page of a Backup Plan Name Is Unique.
* `vault_id` - (Required, ForceNew) Vault ID.
* `instance_id` - (Required, ForceNew) The ECS Instance Id. Must Have Installed the Client.
* `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
* `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval} * startTime Backup start time, UNIX time, in seconds. * interval ISO8601 time interval. E.g: ** PT1H, one hour apart. ** P1D, one day apart. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed, the next backup task will not be triggered.
* `disabled` - (Optional) Whether to Disable the Backup Task. Valid Values: true, false.
* `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.
* `options` - (Optional) Windows System with Application Consistency Using VSS. eg: {`UseVSS`:false}.
* `speed_limit` - (Optional) flow control. The format is: {start}|{end}|{bandwidth} * start starting hour * end end hour * bandwidth limit rate, in KiB ** Use | to separate multiple flow control configurations; ** Multiple flow control configurations are not allowed to have overlapping times.
* `path` - (Optional) Backup Path. e.g. `["/home", "/var"]`
* `exclude` - (Optional) Exclude Path. String of Json List, most 255 Characters. e.g. `"[\"/home/work\"]"`
* `include` - (Optional) Include Path. String of Json List, most 255 Characters. e.g. `"[\"/var\"]"`


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ecs Backup Plan.

## Import

HBR Ecs Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_ecs_backup_plan.example <id>
```
---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_nas_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-nas-backup-plan"
description: |-
  Provides a Alicloud HBR Nas Backup Plan resource.
---

# alicloud\_hbr\_nas\_backup\_plan

Provides a HBR Nas Backup Plan resource.

For information about HBR Nas Backup Plan and how to use it, see [What is Nas Backup Plan](https://www.alibabacloud.com/help/doc-detail/132248.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbr_nas_backup_plan" "example" {
  nas_backup_plan_name = "example_value"
  vault_id             = "v-0003gxoksflhu46w185s"
  file_system_id       = "031cf4964f"
  create_time          = "1603163444"
  schedule             = "I|1602673264|PT2H"
  retention            = "1"
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

* `nas_backup_plan_name` - (Required) The name of the resource.
* `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
* `schedule` - (Required) The Backup Policy. Formats: I | {Range Specified by the StartTime }|{ Interval}\n* The Time Range Specified by the StartTime Backup Start Time in Unix Time Seconds.\n* Interval ISO8601 Time Intervals. For Example:\n**PT1H Interval for an Hour.\n**P1D Interval Day.\nMeaning from {Range Specified by the Starttime} Every {Interval} of the Time Where We Took Backups Once a Task. Does Not Compensate the Has Elapsed Time the Backup Task. If the Last Backup Has Not Been Completed without Triggering the next Backup.
* `disabled` - (Optional) Whether to Disable the Backup Task. Valid Values: true, false.
* `file_system_id` - (Optional, ForceNew) The File System ID.
* `create_time` - (Optional, ForceNew) File System Creation Time. **Note** The time format of the API adopts the ISO 8601 format, such as `2021-07-09T15:45:30CST` or `2021-07-09T07:45:30Z`.
* `include` - (Optional) The include path. String of Json List, most 255 Characters. e.g. `"[\"/home/work\"]"`
* `exclude` - (Optional) The exclude path. String of Json List, most 255 Characters. e.g. `"[\"/var\"]"`
* `path` - (Optional) Backup Path. Up to 65536 Characters. e.g.`["/home", "/var"]`
* `speed_limit` - (Optional) flow control. The format is: {start}|{end}|{bandwidth} * start starting hour * end end hour * bandwidth limit rate, in KiB ** Use | to separate multiple flow control configurations; ** Multiple flow control configurations are not allowed to have overlapping times.
* `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.
* `options` - (Optional) Options. NAS Backup Plan Does Not Support Yet.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nas Backup Plan.

## Import

HBR Nas Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_nas_backup_plan.example <id>
```
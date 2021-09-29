---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_nas_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-nas-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Nas Backup Plan resource.
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

* `nas_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `vault_id` - (Required, ForceNew) The ID of Backup vault.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `schedule` - (Required) Backup strategy. Optional format: I|{startTime}|{interval}. It means to execute a backup task every {interval} starting from {startTime}. The backup task for the elapsed time will not be compensated. If the last backup task is not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `1D` means one day apart.
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `backup_type` - (Required, ForceNew) Backup type. Valid values: `COMPLETE`.
* `file_system_id` - (Required, ForceNew) The File System ID of Nas.
* `create_time` - (Required, ForceNew) File System Creation Time. **Note** The time format of the API adopts the ISO 8601 format, such as `2021-07-09T15:45:30CST` or `2021-07-09T07:45:30Z`.
* `include` - (Optional) The include path. String of Json list, up to 255 characters. e.g. `"[\"/home/work\"]"`
* `exclude` - (Optional) The exclude path. String of Json list, up to 255 characters. e.g. `"[\"/var\"]"`
* `path` - (Optional) Backup path. Up to 65536 characters. e.g.`["/home", "/var"]`
* `speed_limit` - (Optional) Flow control. The format is: {start}|{end}|{bandwidth}. Use `|` to separate multiple flow control configurations, multiple flow control configurations not allowed to have overlapping times.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nas Backup Plan.

## Import

HBR Nas Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_nas_backup_plan.example <id>
```
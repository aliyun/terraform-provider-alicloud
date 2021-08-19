---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_nas_backup_plans"
sidebar_current: "docs-alicloud-datasource-hbr-nas_backup_plans"
description: |-
  Provides a list of Hbr NasBackupPlans to the user.
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
* `file_system_id` - (Optional) The Nas fileSystem instance ID of the EcsBackupPlan used.
* `vault_id` - (Optional) The Vault ID of the EcsBackupPlan used.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of NasBackupPlan names.
* `plans` - A list of Hbr NasBackupPlans. Each element contains the following attributes:
	* `nas_backup_plan_name` - (Required) The name of the resource.
	* `retention` - (Required) Backup Retention Period, the Minimum Value of 1.
	* `schedule` - (Required) The Backup Policy. Formats: I | {Range Specified by the Starttime }|{ Interval}\n* The Time Range Specified by the Starttime Backup Start Time in Unix Time Seconds.\n* Interval ISO8601 Time Intervals. For Example:\n**PT1H Interval for an Hour.\n**P1D Interval Day.\nMeaning from {Range Specified by the Starttime} Every {Interval} of the Time Where We Took Backups Once a Task. Does Not Compensate the Has Elapsed Time the Backup Task. If the Last Backup Has Not Been Completed without Triggering the next Backup.
	* `file_system_id` - (Optional, ForceNew) The File System ID.
	* `create_time` - (Optional, ForceNew) File System Creation Time. Unix Time Seconds.
	* `include` - (Optional) The include path. String of Json List, most 255 Characters. e.g. `"[\"/home/work\"]"`
	* `exclude` - (Optional) The exclude path. String of Json List, most 255 Characters. e.g. `"[\"/var\"]"`
	* `path` - (Optional) Backup Path. Up to 65536 Characters. e.g.`["/home", "/var"]`
	* `speed_limit` - (Optional) flow control. The format is: {start}|{end}|{bandwidth} * start starting hour * end end hour * bandwidth limit rate, in KiB ** Use | to separate multiple flow control configurations; ** Multiple flow control configurations are not allowed to have overlapping times.
	* `backup_type` - (Optional, Computed, ForceNew) Backup Type. Valid Values: * Complete. Valid values: `COMPLETE`.
	* `options` - (Optional) Options. NAS Backup Plan Does Not Support Yet.
	


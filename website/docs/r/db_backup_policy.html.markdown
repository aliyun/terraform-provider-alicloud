---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_backup_policy"
sidebar_current: "docs-alicloud-resource-db-backup-policy"
description: |-
  Provides an RDS backup policy resource.
---

# alicloud\_db\_backup\_policy

Provides an RDS instance backup policy resource and used to configure instance backup policy.

~> **NOTE:** Each DB instance has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```
resource "alicloud_db_backup_policy" "default" {
	instance_id = "rm-2eps..."
	backup_period = ["Monday", "Wednesday"]
	backup_time = "02:00Z-03:00Z"
	retention_period = 7
	log_backup = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance that can run database.
* `backup_period` - (Optional) DB Instance backup period. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to ["Tuesday", "Thursday", "Saturday"].
* `backup_time` - (Optional) DB instance backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `retention_period` - (Optional) Instance backup retention days. Valid values: [7-730]. Default to 7.
* `log_backup` - (Optional) Whether to backup instance log. Default to true.
* `log_retention_period` - (Optional) Instance log backup retention days. Valid values: [7-730]. Default to 7. It can be larger than 'retention_period'.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'instance_id'.
* `instance_id` - The Id of DB instance.
* `backup_period` - DB Instance backup period.
* `backup_time` - DB instance backup time.
* `retention_period` - Instance backup retention days.
* `log_backup` - Whether to backup instance log.
* `log_retention_period` - Instance log backup retention days.

## Import

RDS backup policy can be imported using the id or instance id, e.g.

```
$ terraform import alicloud_db_backup_policy.example "rm-12345678"
```

---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_backup_policy"
sidebar_current: "docs-alicloud-resource-db-backup-policy"
description: |-
  Provides an RDS backup policy resource.
---

# alicloud\_db\_backup\_policy

Provides an RDS instance backup policy resource and used to configure instance backup policy.

-> **NOTE:** Each DB instance has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbbackuppolicybasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_backup_policy" "policy" {
  instance_id = alicloud_db_instance.instance.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `backup_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'preferred_backup_period' instead.
* `preferred_backup_period` - (Optional, available in 1.69.0+) DB Instance backup period. Please set at least two days to ensure backing up at least twice a week. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].
* `backup_time` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'preferred_backup_time' instead.
* `preferred_backup_time` - (Optional, available in 1.69.0+) DB instance backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `retention_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'backup_retention_period' instead.
* `backup_retention_period` - (Optional, available in 1.69.0+) Instance backup retention days. Valid values: [7-730]. Default to 7. But mysql local disk is unlimited.
* `log_backup` - (Deprecated) It has been deprecated from version 1.68.0, and use field 'enable_backup_log' instead.
* `enable_backup_log` - (Optional, available in 1.68.0+) Whether to backup instance log. Valid values are `true`, `false`, Default to `true`. Note: The 'Basic Edition' category Rds instance does not support setting log backup. [What is Basic Edition](https://www.alibabacloud.com/help/doc-detail/48980.htm).
* `log_retention_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'log_backup_retention_period' instead.
* `log_backup_retention_period` - (Optional, available in 1.69.0+) Instance log backup retention days. Valid when the `enable_backup_log` is `1`. Valid values: [7-730]. Default to 7. It cannot be larger than `backup_retention_period`.
* `local_log_retention_hours` - (Optional, available in 1.69.0+) Instance log backup local retention hours. Valid when the `enable_backup_log` is `true`. Valid values: [0-7*24].
* `local_log_retention_space` - (Optional, available in 1.69.0+) Instance log backup local retention space. Valid when the `enable_backup_log` is `true`. Valid values: [0-50].
* `high_space_usage_protection` - (Optional, available in 1.69.0+) Instance high space usage protection policy. Valid when the `enable_backup_log` is `true`. Valid values are `Enable`, `Disable`.
* `log_backup_frequency` - (Optional, available in 1.69.0+) Instance log backup frequency. Valid when the instance engine is `SQLServer`. Valid values are `LogInterval`.
* `compress_type` - (Optional, available in 1.69.0+) The compress type of instance policy. Valid values are `1`, `4`, `8`.
* `archive_backup_retention_period` - (Optional, available in 1.69.0+) Instance archive backup retention days. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values: [30-1095], and `archive_backup_retention_period` must larger than `backup_retention_period` 730.
* `archive_backup_keep_count` - (Optional, available in 1.69.0+) Instance archive backup keep count. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. When `archive_backup_keep_policy` is `ByMonth` Valid values: [1-31]. When `archive_backup_keep_policy` is `ByWeek` Valid values: [1-7].
* `archive_backup_keep_policy` - (Optional, available in 1.69.0+) Instance archive backup keep policy. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values are `ByMonth`, `ByWeek`, `KeepAll`.
* `released_keep_policy` - (Optional, available in 1.147.0+) The policy based on which ApsaraDB RDS retains archived backup files if the instance is released. Default value: None. Valid values:
  * **None**: No archived backup files are retained.
  * **Lastest**: Only the most recent archived backup file is retained.
  * **All**: All archived backup files are retained.

-> **NOTE:** Currently, the SQLServer instance does not support to modify `log_backup_retention_period`.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'instance_id'.

## Import

RDS backup policy can be imported using the id or instance id, e.g.

```
$ terraform import alicloud_db_backup_policy.example "rm-12345678"
```

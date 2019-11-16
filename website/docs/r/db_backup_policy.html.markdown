---
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
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "alicloud_db_backup_policy" "policy" {
  instance_id = "${alicloud_db_instance.instance.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `backup_period` - (Optional) DB Instance backup period. Please set at least two days to ensure backing up at least twice a week. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].
* `backup_time` - (Optional) DB instance backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `retention_period` - (Optional) Instance backup retention days. Valid values: [7-730]. Default to 7.
* `log_backup` - (Optional) Whether to backup instance log. Note: The 'Basic Edition' category Rds instance does not support setting log backup. [What is Basic Edition](https://www.alibabacloud.com/help/doc-detail/48980.htm).

-> **NOTE:** Do not support SQLServer's close operation.
* `log_retention_period` - (Optional) Instance log backup retention days. Valid when the `log_backup` is `true`. Valid values: [7-730]. Default to 7. It cannot be larger than `retention_period`.

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

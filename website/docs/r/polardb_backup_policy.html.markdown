---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_backup_policy"
sidebar_current: "docs-alicloud-resource-polardb-backup-policy"
description: |-
  Provides an PolarDB backup policy resource.
---

# alicloud\_polardb\_backup\_policy

Provides an PolarDB cluster backup policy resource and used to configure cluster backup policy.

-> **NOTE:** Available in v1.66.0+. Each DB cluster has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```
variable "name" {
  default = "polardbClusterconfig"
}

variable "creation" {
  default = "PolarDB"
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

resource "alicloud_polardb_cluster" "default" {
  db_type              = "MySQL"
  db_version           = "8.0"
  db_node_class        = "polar.mysql.x4.large"
  pay_type             = "PostPaid"
  description          = "${var.name}"
  vswitch_id           = "vsw-t4nq4tr8wcuj7397rbws2"
}


resource "alicloud_polardb_backup_policy" "policy" {
  cluster_id    = "${alicloud_polardb_cluster.default.id}"
  backup_period = "Tuesday,Wednesday"
  backup_time   = "10:00Z-11:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `backup_period` - (Optional) PolarDB Cluster backup period. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]. Default to ["Tuesday", "Thursday", "Saturday"].
* `backup_time` - (Optional) PolarDB Cluster backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'cluster_id'.
* `cluster_id` - The Id of PolarDB Cluster.
* `backup_period` - PolarDB Cluster backup period.
* `backup_time` - PolarDB Cluster backup time.
* `retention_period` - Cluster backup retention days.

## Import

PolarDB backup policy can be imported using the id or cluster id, e.g.

```
$ terraform import alicloud_polardb_backup_policy.example "rm-12345678"
```

---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_backup"
sidebar_current: "docs-alicloud-resource-polardb-backup"
description: |-
  Provides a PolarDB backup resource.
---

# alicloud\_polardb\_backup

Provides a PolarDB cluster backup resource and used to backup PolarDB cluster manually. For more information about PolarDB backup, see [PolarDB backup data](https://www.alibabacloud.com/help/doc-detail/72672.htm).

-> **NOTE:** Available in v1.97.0+.

-> **NOTE:** Each PolarDB cluster can have up to 3 manually created backups. If there are already 3 manually created backups, you need to delete the backup before you can manually create a backup.

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
  db_node_class        = "polar.mysql.x4.medium"
  pay_type             = "PostPaid"
  description          = "${var.name}"
  vswitch_id           = "${alicloud_vswitch.default.id}"
}


resource "alicloud_polardb_backup" "polardb_backup" {
  db_cluster_id    = "${alicloud_polardb_cluster.default.id}"
}
```

 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup resource ID. Composed of cluster ID, backup start time, backup ID and backup end time with format `<cluster_id>:<start_time>:<backup_id>:<end_time>`.
* `backup_id` - Id of polarDB cluster backup.
* `backup_status` - Status of polarDB cluster backup.
* `backup_start_time` - Start time of polarDB cluster backup.
* `backup_end_time` - End time of polarDB cluster backup.

## Import

PolarDB backup can be imported using the id , e.g.

```
$ terraform import alicloud_polardb_backup.example "pc-123456:2000-01-01T01-00Z:4567890:2000-01-01T01-20Z"
```

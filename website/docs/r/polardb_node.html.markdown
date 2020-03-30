---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_node"
sidebar_current: "docs-alicloud-resource-polardb-node"
description: |-
  Provides a PolarDB node resource.
---

# alicloud\_polardb\_node

Provides a PolarDB node resource.

-> **NOTE:** Available in v1.77.0+. 

## Example Usage

```
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbnodepgsql"
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

resource "alicloud_polardb_cluster" "cluster" {
  db_type               = "PostgreSQL"
  db_version            = "11"
  db_node_class         = "polar.pg.x4.medium"
  pay_type              = "PostPaid"
  vswitch_id            = "${alicloud_vswitch.default.id}"
  description           = "${var.name}"
}

resource "alicloud_polardb_node" "default" {
  db_cluster_id         = "${alicloud_polardb_cluster.cluster.id}"
  db_node_class         = "polar.pg.x4.medium"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which node belongs.
* `db_node_class` - (Required) The db_node_class of cluster node.
* `modify_type` - (Optional) Use as `db_node_class` change class , define upgrade or downgrade.  Valid values are `Upgrade`, `Downgrade`.

## Attributes Reference

The following attributes are exported:

* `id` - The current node resource ID. Composed of instance ID and node ID with format `<instance_id>:<node_id>`.

## Import

PolarDB node can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_node.example "pc-12345:pi-67890"
```
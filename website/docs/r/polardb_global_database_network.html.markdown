---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_global_database_network"
sidebar_current: "docs-alicloud-resource-polardb-global-database-network"
description: |-
  Provides a Alicloud PolarDB Global Database Network resource.
---

# alicloud\_polardb\_global\_database\_network

Provides a PolarDB Global Database Network resource.

For information about PolarDB Global Database Network and how to use it, see [What is Global Database Network](https://www.alibabacloud.com/help/en/polardb/api-polardb-2017-08-01-createglobaldatabasenetwork).

-> **NOTE:** Available since v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
  pay_type   = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_global_database_network" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  description   = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the primary cluster.
* `status` - (Computed) The status of the Global Database Network.
* `description` - (Optional, Computed) The description of the Global Database Network.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Global Database Network.
* `status` - The status of the Global Database Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the PolarDB Global Database Network.
* `update` - (Defaults to 3 mins) Used when update the PolarDB Global Database Network.
* `delete` - (Defaults to 10 mins) Used when delete the PolarDB Global Database Network.

## Import

PolarDB Global Database Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_global_database_network.example <id>
```
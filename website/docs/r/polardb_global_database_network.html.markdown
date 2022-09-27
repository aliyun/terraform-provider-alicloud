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

For information about PolarDB Global Database Network and how to use it, see [What is Global Database Network](https://www.alibabacloud.com/help/en/polardb-for-mysql/latest/createglobaldatabasenetwork).

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_polardb_node_classes" "default" {
  zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "8.0"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  vswitch_id    = data.alicloud_vswitches.default.ids.0
  description   = "example_value"
}

resource "alicloud_polardb_global_database_network" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  description   = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the primary cluster.
* `description` - (Optional, Computed) The description of the Global Database Network.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Global Database Network.
* `status` - The status of the Global Database Network.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the PolarDB Global Database Network.
* `update` - (Defaults to 3 mins) Used when update the PolarDB Global Database Network.
* `delete` - (Defaults to 10 mins) Used when delete the PolarDB Global Database Network.

## Import

PolarDB Global Database Network can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_global_database_network.example <id>
```
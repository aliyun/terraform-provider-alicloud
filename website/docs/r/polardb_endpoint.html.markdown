---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-endpoint"
description: |-
  Provides a PolarDB instance endpoint resource.
---

# alicloud\_polardb\_endpoint

Provides a PolarDB endpoint resource to manage custom endpoint of PolarDB cluster.

-> **NOTE:** Available in v1.80.0+. Only used to manage PolarDB MySQL custom cluster endpoint.

## Example Usage

```terraform
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbconnectionbasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_endpoint" "endpoint" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  endpoint_type    = "Custom"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `endpoint_type` - (Required, ForceNew) Type of endpoint. Valid value: `Custom`. Currently supported only `Custom`.
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. Default to `ReadOnly`.
* `nodes` - (Optional) Node id list for endpoint configuration. At least 2 nodes if specified, or if the cluster has more than 3 nodes, read-only endpoint is allowed to mount only one node. Default is all nodes.
* `auto_add_new_nodes` - (Optional) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. Default to `Disable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.

## Import

PolarDB endpoint can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_endpoint.example pc-abc123456:pe-abc123456
```
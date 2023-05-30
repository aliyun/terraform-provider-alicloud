---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoint_address"
sidebar_current: "docs-alicloud-resource-poalrdb-endpoint"
description: |-
  Provides a PolarDB instance endpoint resource.
---

# alicloud\_polardb\_endpoint_address

Provides a PolarDB endpoint address resource to allocate an Internet endpoint address string for PolarDB instance.

-> **NOTE:** Available in v1.68.0+. Each PolarDB instance will allocate a intranet connection string automatically and its prefix is Cluster ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
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

data "alicloud_polardb_endpoints" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
}

resource "alicloud_polardb_endpoint_address" "default" {
  db_cluster_id     = alicloud_polardb_cluster.default.id
  db_endpoint_id    = data.alicloud_polardb_endpoints.default.endpoints[0].db_endpoint_id
  connection_prefix = "polardbexample"
  net_type          = "Public"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `db_endpoint_id` - (Required, ForceNew) The Id of endpoint that can run database.
* `connection_prefix` - (Optional) Prefix of the specified endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter.
* `net_type` - (Optional, ForceNew) Internet connection net type. Valid value: `Public`. Default to `Public`. Currently supported only `Public`.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `port` - Connection cluster or endpoint port.
* `connection_string` - Connection cluster or endpoint string.
* `ip_address` - The ip address of connection string.

## Import

PolarDB endpoint address can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_endpoint_address.example pc-abc123456:pe-abc123456
```

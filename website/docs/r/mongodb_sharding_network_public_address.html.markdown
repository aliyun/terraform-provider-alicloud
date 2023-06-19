---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_network_public_address"
sidebar_current: "docs-alicloud-resource-mongodb-sharding-network-public-address"
description: |-
  Provides a Alicloud MongoDB Sharding Network Public Address resource.
---

# alicloud_mongodb_sharding_network_public_address

Provides a MongoDB Sharding Network Public Address resource.

For information about MongoDB Sharding Network Public Address and how to use it, see [What is Sharding Network Public Address](https://www.alibabacloud.com/help/doc-detail/67602.html).

-> **NOTE:** Available since v1.149.0.

-> **NOTE:** This operation supports sharded cluster instances only.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_mongodb_zones" "default" {}
locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = local.zone_id
  vswitch_id     = alicloud_vswitch.default.id
  engine_version = "4.2"
  name           = var.name
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = "10"
  }
  shard_list {
    node_class        = "dds.shard.standard"
    node_storage      = "20"
    readonly_replicas = "1"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}

resource "alicloud_mongodb_sharding_network_public_address" "example" {
  db_instance_id = alicloud_mongodb_sharding_instance.default.id
  node_id        = alicloud_mongodb_sharding_instance.default.mongo_list.0.node_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the instance.
* `node_id` - (Required, ForceNew) The ID of the `mongos`, `shard`, or `Configserver` node in the sharded cluster instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Sharding Network Public Address. The value formats as `<db_instance_id>:<node_id>`.
* `network_address` - An array that consists of the endpoints of ApsaraDB for MongoDB instances.
  * `network_type` - The network type.
  * `network_address` - The endpoint of the instance.
  * `node_type` - The type of the node.
  * `port` - The port number.
  * `role` - The role of the node.
  * `vpc_id` - The ID of the VPC.
  * `expired_time` - The remaining duration of the classic network address. Unit: `seconds`.
  * `ip_address` - The IP address of the instance.
  * `vswitch_id` - The vSwitch ID of the VPC.
  * `node_id` - The ID of the `mongos`, `shard`, or `Configserver` node in the sharded cluster instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the MongoDB Sharding Network Public Address.
* `delete` - (Defaults to 5 mins) Used when terminating the MongoDB Sharding Network Public Address.

## Import

MongoDB Sharding Network Public Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_sharding_network_public_address.example <db_instance_id>:<node_id>
```
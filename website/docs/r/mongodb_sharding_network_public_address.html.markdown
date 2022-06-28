---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_network_public_address"
sidebar_current: "docs-alicloud-resource-mongodb-sharding-network-public-address"
description: |-
  Provides a Alicloud MongoDB Sharding Network Public Address resource.
---

# alicloud\_mongodb\_sharding\_network\_public\_address

Provides a MongoDB Sharding Network Public Address resource.

For information about MongoDB Sharding Network Public Address and how to use it, see [What is Sharding Network Public Address](https://www.alibabacloud.com/help/doc-detail/67602.html).

-> **NOTE:** Available in v1.149.0+.

-> **NOTE:** This operation supports sharded cluster instances only.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name = "subnet-for-local-test"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  engine_version = "3.4"
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the MongoDB Sharding Network Public Address.
* `delete` - (Defaults to 5 mins) Used when terminating the MongoDB Sharding Network Public Address.

## Import

MongoDB Sharding Network Public Address can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_sharding_network_public_address.example <db_instance_id>:<node_id>
```
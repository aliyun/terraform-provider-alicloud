---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_network_private_address"
sidebar_current: "docs-alicloud-resource-mongodb-sharding-network-private-address"
description: |-
  Provides a Alicloud MongoDB Sharding Network Private Address resource.
---

# alicloud_mongodb_sharding_network_private_address

Provides a MongoDB Sharding Network Private Address resource.

For information about MongoDB Sharding Network Private Address and how to use it, see [What is Sharding Network Private Address](https://www.alibabacloud.com/help/en/doc-detail/141403.html).

-> **NOTE:** Available since v1.157.0.

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

resource "alicloud_mongodb_sharding_network_private_address" "default" {
  db_instance_id   = alicloud_mongodb_sharding_instance.default.id
  node_id          = alicloud_mongodb_sharding_instance.default.shard_list.0.node_id
  zone_id          = alicloud_mongodb_sharding_instance.default.zone_id
  account_name     = "example"
  account_password = "Example_123"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Optional) The name of the account. 
  - The name must be 4 to 16 characters in length and can contain lowercase letters, digits, and underscores (_). It must start with a lowercase letter.
  - You need to set the account name and password only when you apply for an endpoint for a shard or Configserver node for the first time. In this case, the account name and password are used for all shard and Configserver nodes.
  - The permissions of this account are fixed to read-only.
* `account_password` - (Optional, Sensitive) Account password. 
  - The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `!#$%^&*()_+-=`.
  - The password must be 8 to 32 characters in length.
* `db_instance_id` - (Required) The db instance id.
* `zone_id` - (Required) The zone ID of the instance.
* `node_id` - (Required) The ID of the Shard node or the ConfigServer node.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Sharding Network Private Address. The value formats as `<db_instance_id>:<node_id>`.
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

## Import

MongoDB Sharding Network Private Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_sharding_network_private_address.example <db_instance_id>:<node_id>
```
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_sharding_network_private_address&exampleId=fefc11d2-8aa2-b0ef-9865-e393b93452836da95b7d&activeTab=example&spm=docs.r.mongodb_sharding_network_private_address.0.fefc11d28a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_mongodb_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
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

* `db_instance_id` - (Required, ForceNew) The ID of the sharded cluster instance.
* `node_id` - (Required, ForceNew) The ID of the Shard node or ConfigServer node.
* `zone_id` - (Required, ForceNew) The zone ID of the instance.
* `account_name` - (Optional) The username of the account.
  - The name must be 4 to 16 characters in length and can contain lowercase letters, digits, and underscores (_). It must start with a lowercase letter.
  - You need to set the account name and password only when you apply for an endpoint for a shard or ConfigServer node for the first time. In this case, the account name and password are used for all shard and ConfigServer nodes.
  - The permissions of this account are fixed to read-only.
* `account_password` - (Optional, Sensitive) The password for the account.
  - The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `!#$%^&*()_+-=`.
  - The password must be 8 to 32 characters in length.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sharding Network Private Address. It formats as `<db_instance_id>:<node_id>`.
* `network_address` - The endpoints of ApsaraDB for MongoDB instances.
  * `node_id` - The ID of the `Shard`, or `ConfigServer` node.
  * `node_type` - The type of the node.
  * `role` - The role of the node.
  * `vpc_id` - The ID of the VPC.
  * `vswitch_id` - The ID of the vSwitch in the VPC.
  * `network_type` - The network type of the instance.
  * `network_address` - The connection string of the instance.
  * `ip_address` - The IP address of the instance.
  * `port` - The port that is used to connect to the instance.
  * `expired_time` - The remaining duration of the classic network endpoint.

## Import

MongoDB Sharding Network Private Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_sharding_network_private_address.example <db_instance_id>:<node_id>
```

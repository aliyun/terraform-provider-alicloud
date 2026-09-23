---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_network_public_addresses"
sidebar_current: "docs-alicloud-datasource-mongodb-sharding-network-public-addresses"
description: |-
  Provides a list of Mongodb Sharding Network Public Addresses to the user.
---

# alicloud\_mongodb\_sharding\_network\_public\_addresses

This data source provides the Mongodb Sharding Network Public Addresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_sharding_network_public_addresses" "example" {
  db_instance_id = "example_value"
  node_id        = "example_value"
  role           = "Primary"
}
output "mongodb_sharding_network_public_address_db_instance_id_1" {
  value = data.alicloud_mongodb_sharding_network_public_addresses.example.addresses.0.db_instance_id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The db instance id.
* `node_id` - (Optional, ForceNew) The ID of the `mongos`, `shard`, or `Configserver` node in the sharded cluster instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `role` - (Optional, ForceNew) The role of the node. Valid values: `Primary` or `Secondary`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `addresses` - A list of Mongodb Sharding Network Addresses. Each element contains the following attributes:
  * `db_instance_id` - The db instance id.
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
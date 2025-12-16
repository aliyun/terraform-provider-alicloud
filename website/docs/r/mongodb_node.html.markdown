---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_node"
description: |-
  Provides a Alicloud Mongodb Node resource.
---

# alicloud_mongodb_node

Provides a Mongodb Node resource.

The sub-resources of the ShardingInstance, including the cs, shards, and mongos nodes.

For information about Mongodb Node and how to use it, see [What is Node](https://next.api.alibabacloud.com/document/Dds/2015-12-01/CreateNode).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "ipv4_cidr" {
  default = "10.0.0.0/24"
}

resource "alicloud_vpc" "default" {
  description = "tf-example"
  vpc_name    = "tf-vpc-shanghai-b"
  cidr_block  = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = var.zone_id
  cidr_block   = var.ipv4_cidr
  vswitch_name = "tf-shanghai-B"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  engine_version   = "4.2"
  vswitch_id       = alicloud_vswitch.default.id
  zone_id          = var.zone_id
  name             = var.name
  storage_type     = "cloud_auto"
  provisioned_iops = 60
  config_server_list {
    node_class   = "mdb.shard.2x.xlarge.d"
    node_storage = 40
  }
  mongo_list {
    node_class = "mdb.shard.2x.xlarge.d"
  }
  mongo_list {
    node_class = "mdb.shard.2x.xlarge.d"
  }
  shard_list {
    node_class   = "mdb.shard.2x.xlarge.d"
    node_storage = 40
  }
  shard_list {
    node_class   = "mdb.shard.2x.xlarge.d"
    node_storage = 40
  }
  lifecycle {
    ignore_changes = [shard_list]
  }
}

resource "alicloud_mongodb_node" "default" {
  account_password  = "q1w2e3r4!"
  auto_pay          = "true"
  node_class        = "mdb.shard.4x.large.d"
  shard_direct      = "false"
  business_info     = "1234"
  node_storage      = "40"
  readonly_replicas = "0"
  db_instance_id    = alicloud_mongodb_sharding_instance.default.id
  node_type         = "shard"
  account_name      = "root"
}
```

## Argument Reference

The following arguments are supported:
* `account_name` - (Optional) Account name, value description:
  - Begins with a lowercase letter.
  - Consists of lowercase letters, numbers, or underscores (_).
  - 4~16 characters in length.

-> **NOTE:** - apsaradb for MongoDB does not support using keywords as accounts.
  - The permissions of the account are fixed to read-only permissions.
  - When applying for a direct connection address of a Shard node for the first time, you need to set an account and password.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `account_password` - (Optional) Account password, value description:
  - Consists of at least three of uppercase letters, lowercase letters, numbers, and special characters.
  - Oh-! @#$%^& *()_+-= is a special character.
  - Length is 8~32 characters.

-> **NOTE:**  apsaradb for MongoDB does not support resetting the account and password of the Shard node.


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `auto_pay` - (Optional) Whether to pay automatically. Value description:
  - `true` (default): Pay automatically, please ensure that the account has sufficient balance.
  - `false`: paid manually.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `business_info` - (Optional) Additional parameters, business information.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `db_instance_id` - (Required, ForceNew) The ID of the sharded cluster instance.
* `effective_time` - (Optional) Effective time of configuration change. Value description:
  - `Immediately` (default): takes effect Immediately.
  - `MaintainTime`: takes effect during the O & M period of the instance.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `from_app` - (Optional) Request Source, value description:
  - `OpenApi`: The request source is OpenApi.
  - `mongo_buy`: The request source is the console.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `node_class` - (Required) The specifications of the Shard node or Mongos node. For more information, see [Instance Specifications](~~ 57141 ~~).
* `node_storage` - (Optional, Int) The disk space of the Node. Unit: GB.

Value range: `10` to `2000`, with a step size of 10GB.

-> **NOTE:**  When the node type is `Shard`, you need to configure this parameter.

* `node_type` - (Required, ForceNew) Node type, value description:
  - `shard`:Shard node.
  - `mongos`: the Mongos node.
* `order_type` - (Optional) Order type, value description:
  - `UPGRADE`: UPGRADE.
  - `DOWNGRADE`: downgrading.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `readonly_replicas` - (Optional, Computed, Int) The number of read-only nodes in the Shard.

Value range: `0` to `5` (integer). Default value: **0 * *.

-> **NOTE:**  This parameter is currently only supported by China Station.

* `shard_direct` - (Optional) Whether to apply for the direct connection address of the Shard node. Value description:
  - `true`: Apply for the direct connection address of the Shard node.
  - `false`:(default) Do not apply for the direct connection address of the Shard node.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `switch_time` - (Optional) The execution time of the change configuration, in the format of  yyyy-MM-dd T  HH:mm:ss Z(UTC time).

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<node_id>`.
* `node_id` - The first ID of the resource
* `status` - Running status of node in sharded cluster

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 30 mins) Used when create the Node.
* `delete` - (Defaults to 5 mins) Used when delete the Node.
* `update` - (Defaults to 40 mins) Used when update the Node.

## Import

Mongodb Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_node.example <db_instance_id>:<node_id>
```
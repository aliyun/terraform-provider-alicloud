---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_node"
sidebar_current: "docs-alicloud-resource-mongodb-sharding-node"
description: |-
  Provides a Alicloud MongoDB Sharding Node resource.
---

# alicloud\_mongodb\_sharding\_node

Provides a MongoDB Sharding Node resource.

For information about MongoDB Sharding Node and how to use it, see [What is Sharding Node](https://www.alibabacloud.com/help/doc-detail/61911.htm).

-> **NOTE:** Available in v1.149.0+.

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
resource "alicloud_mongodb_sharding_node" "example" {
  db_instance_id = alicloud_mongodb_sharding_instance.default.id
  node_class     = "dds.shard.mid"
  node_storage   = "10"
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the sharded cluster instance.
* `node_class` - (Required) The instance type of the shard node. For more information, see [Instance types](https://www.alibabacloud.com/help/doc-detail/57141.htm).
* `node_storage` - (Required) The disk space of the Node. The value range is `10` to `2000` GB and the step size is `10` GB.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Sharding Node. The value formats as `<db_instance_id>:<node_id>`.
* `node_id` - The ID of the node.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Sharding Node.
* `delete` - (Defaults to 20 mins) Used when delete the Sharding Node.
* `update` - (Defaults to 20 mins) Used when update the Sharding Node.

## Import

MongoDB Sharding Node can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_sharding_node.example <db_instance_id>:<node_id>
```
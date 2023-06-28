---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_resource_group"
sidebar_current: "docs-alicloud-resource-adb-resource-group"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) Resource Group resource.
---

# alicloud_adb_resource_group

Provides a Adb Resource Group resource.

For information about Adb Resource Group and how to use it, see [What is Adb Resource Group](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/api-doc-adb-2019-03-15-api-doc-createdbresourcegroup).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_adb_zones" "default" {
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_adb_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster" "default" {
  compute_resource    = "48Core192GBNEW"
  db_cluster_category = "MixedStorage"
  db_cluster_version  = "3.0"
  db_node_class       = "E32"
  db_node_count       = 1
  db_node_storage     = 100
  description         = var.name
  elastic_io_resource = 1
  maintain_time       = "04:00Z-05:00Z"
  mode                = "flexible"
  payment_type        = "PayAsYouGo"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_ips        = ["10.168.1.12", "10.168.1.11"]
  vpc_id              = alicloud_vpc.default.id
  vswitch_id          = alicloud_vswitch.default.id
  zone_id             = data.alicloud_adb_zones.default.zones[0].id
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_adb_resource_group" "default" {
  group_name    = "TF_EXAMPLE"
  group_type    = "batch"
  node_num      = 1
  db_cluster_id = alicloud_adb_db_cluster.default.id
}
```

## Argument Reference

The following arguments are supported:
* `db_cluster_id` - (Required, ForceNew) DB cluster id.
* `group_name` - (Required, ForceNew) The name of the resource pool. The group name must be 2 to 30 characters in length, and can contain upper case letters, digits, and underscore(_).
* `group_type` - (Optional, ForceNew) Query type, value description:
  * **etl**: Batch query mode.
  * **interactive**: interactive Query mode.
  * **default_type**: the default query mode.
* `node_num` - (Optional) The number of nodes. The default number of nodes is 0. The number of nodes must be less than or equal to the number of nodes whose resource name is USER_DEFAULT.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<db_cluster_id>:<group_name>`.
* `create_time` - Creation time.
* `update_time` - Update time.
* `user` - Binding User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Group.
* `delete` - (Defaults to 1 mins) Used when delete the Resource Group.
* `update` - (Defaults to 5 mins) Used when update the Resource Group.

## Import

Adb Resource Group can be imported using the id, e.g.

```shell
$terraform import alicloud_adb_resource_group.example <db_cluster_id>:<group_name>
```
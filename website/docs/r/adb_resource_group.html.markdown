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

For information about Adb Resource Group and how to use it, see [What is Adb Resource Group](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/create-db-resource-group).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_adb_resource_group" "default" {
  group_name    = "TESTOPENAPI"
  group_type    = "batch"
  node_num      = 0
  db_cluster_id = "am-bp1a16357gty69185"
}
```

## Argument Reference

The following arguments are supported:
* `db_cluster_id` - (Required,ForceNew) DB cluster id.
* `group_name` - (Required,ForceNew) The name of the resource pool. The group name must be 2 to 30 characters in length, and can contain upper case letters, digits, and underscore(_).
* `group_type` - (ForceNew,Optional) Query type, value description:
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Group.
* `delete` - (Defaults to 1 mins) Used when delete the Resource Group.
* `update` - (Defaults to 5 mins) Used when update the Resource Group.

## Import

Adb Resource Group can be imported using the id, e.g.

```shell
$terraform import alicloud_adb_resource_group.example <db_cluster_id>:<group_name>
```
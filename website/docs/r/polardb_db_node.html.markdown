---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_db_node"
description: |-
  Provides a Alicloud PolarDB DB Node resource.
---

# alicloud_polardb_db_node

Provides a PolarDB DB Node resource.

For information about PolarDB DB Node and how to use it, see [CreateDBNodes](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateDBNodes).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_node {
    target_class = "polar.mysql.x4.large"
  }
}
```

## Argument Reference

The following arguments are supported:
* `auto_use_coupon` - (Optional) Whether to automatically use coupons. Valid values:
  - `true` (default): use coupons.
  - `false`: do not use coupons.
* `cloud_provider` - (Optional) The cloud service provider to which the instance belongs.
* `db_cluster_id` - (Required, ForceNew) The ID of the cluster.
* `db_node` - (Required, ForceNew, List) Add node information. See [`db_node`](#db_node) below.
* `db_node_ids` - (Computed) The DBNodeId list created.
* `db_node_target_class` - (Optional) The target specifications of all nodes, see [Compute Node Specifications](https://www.alibabacloud.com/help/en/polardb/polardb-for-mysql/user-guide/specifications-of-compute-nodes).
* `db_node_type` - (Optional) Node type. Value range:
  - RO
  - STANDBY
  - DLNode
* `endpoint_bind_list` - (Optional, ForceNew) The address ID of the cluster connection address to which the new node needs to join. If you need to add multiple connection addresses at the same time, separate multiple connection address IDs with commas (,).
* `from_time_service` - (Optional) Whether to execute the restart node task immediately or on schedule. Value range:
  - `false` (default): scheduled execution
  - `true`: immediate execution
* `imci_switch` - (Optional, ForceNew) Whether to enable column storage index. Value range:
  - `ON`
  - `OFF`
* `modify_type` - (Optional) Change type, the range of values is as follows:
  - `Upgrade`: Upgrade specification
  - `Downgrade`: Downgrade specification
* `planned_end_time` - (Optional) The latest time to start the target scheduled task. The format is `YYYY-MM-DDThh:mm:ssZ` (UTC).
* `planned_flashing_off_time` - (Optional) Planned flashing-off time. The time uses ISO 8601 format with UTC timezone (format: YYYY-MM-DDTHH:mm:ssZ).
* `planned_start_time` - (Optional) Start execution timing (that is, the earliest time when the new node task is added within the target time period). The format is `YYYY-MM-DDThh:mm:ssZ` (UTC).
* `promotion_code` - (Optional) Coupon number. If not specified, the default coupon is used.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `restart_planned_end_time` - (Optional) The planned end time for the RestartDBNode operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC). The latest time should be 30 minutes or more later than the start time.
* `restart_planned_start_time` - (Optional) The planned start time for the RestartDBNode operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC).
* `scc_mode` - (Optional) Cluster strong consistent SCC. Valid values: `ON`, `OFF`.
* `sub_category` - (Optional, Computed) Cluster subseries. The range of values is as follows:
  - `normal_exclusive`: exclusive specification
  - `Normal_General`: General specification

### `db_node`

The db_node supports the following:
* `target_class` - (Optional, ForceNew) The specifications of the new node must be the same as those of the original node. For more information about the node specifications, see [Compute Node Specifications](https://www.alibabacloud.com/help/en/polardb/polardb-for-mysql/user-guide/specifications-of-compute-nodes).
* `zone_id` - (Optional, ForceNew) The zone of the new node must be consistent with the original node. You can call the DescribeRegions operation to query the zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the DB Node.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the DB Node.
* `delete` - (Defaults to 5 mins) Used when delete the DB Node.
* `update` - (Defaults to 5 mins) Used when update the DB Node.

## Import

PolarDB DB Node can be imported using the composite id `<db_cluster_id>:<db_node_id>`, e.g.

```shell
$ terraform import alicloud_polardb_db_node.example <db_cluster_id>:<db_node_id>
```

---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_db_node"
description: |-
  Provides a Alicloud Polardb Db Node resource.
---

# alicloud_polardb_db_node

Provides a Polardb Db Node resource.

Adds a read-only, standby or DLNode to an existing PolarDB cluster. The resource manages exactly one node: `db_node` accepts a single element, and the created node ID is exposed as `id` and `db_node_ids`.

For information about Polardb Db Node and how to use it, see [What is Db Node](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateDBNodes).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** Destroying this resource removes the node from the cluster. The cluster and its other nodes are not affected.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type            = "MySQL"
  db_version         = "8.0"
  pay_type           = "PostPaid"
  db_node_class      = "polar.mysql.x4.medium"
  vswitch_id         = alicloud_vswitch.default.id
  encrypt_new_tables = "OFF"
}

resource "alicloud_polardb_db_node" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_node_type  = "RO"

  db_node = [{
    target_class = "polar.mysql.x4.medium"
    zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  }]

  # The cluster is created in the same apply; its order can still be processing.
  timeouts {
    create = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:
* `auto_use_coupon` - (Optional) Whether to use coupons automatically. Valid values: `true` and `false`. Defaults to `true`.

  -> **NOTE:** This parameter is not returned by the API; Terraform keeps the configured value in state. It is sent with `CreateDBNodes` and `ModifyDBNodeClass`, so a change alone does not trigger an API call.

* `cloud_provider` - (Optional) The cloud service provider of the node.

  -> **NOTE:** This parameter is not returned by the API; Terraform keeps the configured value in state. It is sent with `CreateDBNodes` and `ModifyDBNodeClass`, so a change alone does not trigger an API call.

* `db_cluster_id` - (Required, ForceNew) The ID of the cluster.
* `db_node` - (Required, ForceNew, List) Add node information. See [`db_node`](#db_node) below.

  -> **NOTE:** `db_node` describes the node to create. It is a list of objects, so both `target_class` and `zone_id` must be specified in every element, and only one element can be set because the resource manages a single node. Use `db_node_target_class` and `modify_type` to change the class of the existing node; changing `db_node` replaces the resource.

* `db_node_target_class` - (Optional) The target class of the node for the `ModifyDBNodeClass` operation. For more information about the node specifications, see [Compute Node Specifications](~~ 102542 ~~).

  -> **NOTE:** This parameter is only used by `ModifyDBNodeClass`, which requires both `db_node_target_class` and `modify_type`. Set them together to resize the node after creation; setting them at creation time does not resize the node. The operation applies the target class to all nodes of the cluster, not only the node managed by this resource.

* `db_node_type` - (Optional) The node type. Valid values: `RO`, `STANDBY` and `DLNode`.

  -> **NOTE:** This parameter is not returned by the API; Terraform keeps the configured value in state. It is sent with `CreateDBNodes` and `ModifyDBNodeClass`, so a change alone does not trigger an API call.

* `endpoint_bind_list` - (Optional, ForceNew) The address ID of the cluster connection address to which the new node needs to join. If you need to add multiple connection addresses at the same time, separate multiple connection address IDs with commas (,).

  -> **NOTE:**  * You can call the [DescribeDBClusterEndpoints](~~ 98205 ~~) interface to view the details of the cluster connection address, including the connection address ID.
  * You can enter the default cluster connection address and the address ID of the custom cluster connection address.
  * If this parameter is left blank, the new node is added to all cluster connection addresses that enable the **New node auto-add** function (that is, the value of 'AutoAddNewNodes' is 'Enable').

* `from_time_service` - (Optional) Whether to execute the restart node task immediately or on schedule. Valid values:
  - `false` (default): scheduled execution
  - `true`: immediate execution

  -> **NOTE:** Changing this parameter triggers `RestartDBNode` together with `restart_planned_start_time` and `restart_planned_end_time`.

* `imci_switch` - (Optional, Computed, ForceNew) Whether to enable column storage index. Valid values:
  - `ON`: enable the column storage index
  - `OFF`: disable the column storage index

  -> **NOTE:** This parameter is only evaluated at creation. Changing it replaces the resource.

* `modify_type` - (Optional) Change type, the range of values is as follows:
  - `Upgrade`: Upgrade specification
  - `Downgrade`: Downgrade specification

  -> **NOTE:** This parameter is only used by `ModifyDBNodeClass`, which requires both `db_node_target_class` and `modify_type`. Set them together to resize the node after creation; setting them at creation time does not resize the node.

* `planned_end_time` - (Optional) The latest time to start the target scheduled task. The format is' YYYY-MM-DDThh:mm:ssZ'(UTC).

  -> **NOTE:**  * The latest time should be 30 minutes or more later than the start time.
  * This parameter is sent with `CreateDBNodes` and `ModifyDBNodeClass`; changing it alone does not trigger an API call.
  * If plannedstarttime' is set but this parameter is not specified, the latest time to execute the target task is 'start time +30 minutes' by default'. For example, when the' plannedstarttime' is set to' 2021-01-14 T09:00:00Z' and this parameter is left blank, the target task will start execution at the latest '2021-01-14 T09:30:00Z'.

* `planned_flashing_off_time` - (Optional) Planned flashing-off time. The time uses ISO 8601 format with UTC timezone (format: YYYY-MM-DDTHH:mm:ssZ).

  -> **NOTE:** This parameter is only used by `ModifyDBNodeClass`; changing it alone does not trigger an API call.

* `planned_start_time` - (Optional) Start execution timing (that is, the earliest time when the new node task is added within the target time period). The format is' YYYY-MM-DDThh:mm:ssZ'(UTC).

  -> **NOTE:**  * The value range of the start time is any time point in the next 24 hours. For example, the current time is' 2021-01-14 T09:00:00Z', and the starting time range allowed here is' 2021-01-14 T09:00:00Z' ~ '2021-01-15 T09:00:00Z'.
  * This parameter is sent with `CreateDBNodes` and `ModifyDBNodeClass`; changing it alone does not trigger an API call.
  * If this parameter is left blank, the new node task will be executed immediately by default.

* `promotion_code` - (Optional) The coupon code. If it is left empty, the default coupon is used.

  -> **NOTE:** This parameter is not returned by the API; Terraform keeps the configured value in state. It is sent with `CreateDBNodes` and `ModifyDBNodeClass`, so a change alone does not trigger an API call.

* `resource_group_id` - (Optional, Computed) The ID of the resource group of the cluster.

  -> **NOTE:** This parameter is only sent with `CreateDBNodes` and is reported by the cluster; changing it after creation does not trigger an API call.
* `restart_planned_end_time` - (Optional) The planned end time for the `RestartDBNode` operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC). The latest time should be 30 minutes or more later than the start time.

  -> **NOTE:** Changing this parameter triggers `RestartDBNode`; leave it empty to restart immediately.

* `restart_planned_start_time` - (Optional) The planned start time for the `RestartDBNode` operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC).

  -> **NOTE:** Changing this parameter triggers `RestartDBNode`; leave it empty to restart immediately.

* `scc_mode` - (Optional) Cluster strong consistent SCC. Valid values: `ON` and `OFF`.

  -> **NOTE:** Changing this parameter triggers `ModifyDBNodeSccMode`.

* `sub_category` - (Optional, Computed) The sub-series of the cluster. Valid values:
  - `normal_exclusive`: exclusive specification
  - `normal_general`: general specification

  -> **NOTE:** This parameter is only used by `ModifyDBNodeClass`, so it must be set together with `db_node_target_class` and `modify_type`; changing it alone does not trigger an API call. It is required when the class change switches between exclusive and general specifications.

### `db_node`

The db_node supports the following:
* `target_class` - (Required, ForceNew) The specifications of the new node must be the same as those of the original node. For more information about the node specifications, see [Compute Node Specifications](~~ 102542 ~~).

  -> **NOTE:**  * DBNode.N.ZoneId and DBNode.N.TargetClass, fill in at least one. N is an integer starting from 1, and the maximum value is 16-the current number of nodes.
  * This parameter is required when the cluster is PolarDB O engine or PolarDB PostgreSQL engine. This parameter is not required when the cluster is PolarDB MySQL engine.

* `zone_id` - (Required, ForceNew) The zone of the new node must be consistent with the original node. You can call the [DescribeRegions](~~ 98041 ~~) operation to query the zone ID.

  -> **NOTE:**  * DBNode.N.ZoneId and DBNode.N.TargetClass, fill in at least one. N is an integer starting from 1, and the maximum value is 16-the current number of nodes.
  * This parameter is required when the cluster is PolarDB O engine or PolarDB PostgreSQL engine. This parameter is not required when the cluster is PolarDB MySQL engine.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `db_node_ids` - The ID list of the created node. The resource tracks a single node, so the list contains exactly one element and equals `id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Db Node.
* `delete` - (Defaults to 5 mins) Used when delete the Db Node.
* `update` - (Defaults to 5 mins) Used when update the Db Node.

## Import

Polardb Db Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_db_node.example <db_node_id>
```

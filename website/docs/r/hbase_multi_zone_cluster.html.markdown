---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_multi_zone_cluster"
sidebar_current: "docs-alicloud-resource-hbase-multi-zone-cluster"
description: |-
  Provides a Alicloud HBase Multi Zone Cluster resource.
---

# alicloud\_hbase\_multi\_zone\_cluster

Provides a HBase Multi Zone Cluster resource.

For information about HBase Multi Zone Cluster and how to use it, see [What is Multi Zone Cluster](https://www.alibabacloud.com/help/en/apsaradb-for-hbase/latest/createmultizonecluster).

-> **NOTE:** Available in v1.180.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbase_multi_zone_cluster" "example" {
  arbiter_vswitch_id     = "example_value"
  arbiter_zone_id        = "example_value"
  arch_version           = "example_value"
  core_disk_size         = 1
  core_disk_type         = "example_value"
  core_instance_type     = "example_value"
  core_node_count        = 1
  engine                 = "example_value"
  engine_version         = "example_value"
  log_disk_size          = 1
  log_disk_type          = "example_value"
  log_instance_type      = "example_value"
  log_node_count         = 1
  master_instance_type   = "example_value"
  multi_zone_combination = "example_value"
  payment_type           = "example_value"
  primary_vswitch_id     = "example_value"
  primary_zone_id        = "example_value"
  standby_vswitch_id     = "example_value"
  standby_zone_id        = "example_value"
  vpc_id                 = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `arbiter_vswitch_id` - (Required) The ID of the vSwitch that is specified for the arbiter. The vSwitch must be deployed in the zone specified by the ArbiterZoneId parameter.
* `arbiter_zone_id` - (Required, ForceNew) The ID of the zone where the arbiter is deployed.
* `arch_version` - (Required) The version of the deployment architecture. This parameter only takes effect if you specify hbaseue for Engine. Valid values: `2.0`.
* `auto_renew_period` - (Optional) The duration for which the ApsaraDB for HBase cluster is automatically renewed after the cluster expires. Unit: month.
* `cluster_name` - (Optional, ForceNew) The name of the ApsaraDB for HBase cluster. The specified name must meet the following requirements:
  - The name must be 2 to 128 characters in length.
  - The name must start with a letter.
  - The name can contain digits, periods (.), hyphens (-), and underscores (_).
* `core_disk_size` - (Required) The disk size of each core node. The valid values range from `400` to `64000`. Unit: GB. Step size: `40`.
* `core_disk_type` - (Required, ForceNew) The disk type of the core nodes. Valid values:
  - `cloud_efficiency`: efficient cloud disk.
  - `cloud_ssd`: SSD disk.
  - `local_hdd_pro`: Throughput-intensive local disk.
  - `local_ssd_pro`:I/O intensive local disk.
* `core_instance_type` - (Required) The instance type of the core nodes.
* `core_node_count` - (Required) The number of core nodes. The valid values range from 2 to 20. Step size: 2.
* `engine` - (Required, ForceNew) The type of engine. This parameter takes effect only on ApsaraDB for HBase Performance-enhanced Edition. Set the value to hbaseue.
* `engine_version` - (Required) The version of the engine. Valid values: `2.0`.
* `immediate_delete_flag` - (Optional) The immediate delete flag. Valid values: `false`, `true`.
* `log_disk_size` - (Required) The disk size of each log node. The valid values range from 400 to 64000. Unit: GB. Step size: 40.
* `log_disk_type` - (Required, ForceNew) The disk type of the log nodes. Valid values:
  - `cloud_efficiency`: efficient cloud disk.
  - `cloud_ssd`:SSD disk.
  - `local_hdd_pro`: Throughput-intensive local disk.
  - `local_ssd_pro`:I/O intensive local disk.
* `log_instance_type` - (Required) The instance type of the log nodes.
* `log_node_count` - (Required) The number of log nodes. The valid values range from 4 to 400. The specified value must be a multiple of 4.
* `master_instance_type` - (Required) The instance type of the master node.
* `multi_zone_combination` - (Required, ForceNew) The combination of zones.
* `payment_type` - (Required, ForceNew) The billing method of the ApsaraDB for HBase cluster. Valid values: `PayAsYouGo`, `Subscription`.
* `period` - (Optional) The duration of the subscription. Valid values:
  - If PeriodUnit is set to year, the valid values of the Period parameter are `1`, `2`, and `3`.
  - If PeriodUnit is set to month, the valid values of the Period parameter are integers that range from `1` to `9`.
* `period_unit` - (Optional) The unit of the subscription duration. Valid values: `year`, `month`.
* `primary_core_node_count` - (Optional) The primary core node count.
* `primary_vswitch_id` - (Required) The ID of the vSwitch that is specified for the primary instance. The vSwitch must be deployed in the zone specified by the PrimaryZoneId parameter.
* `primary_zone_id` - (Required, ForceNew) The ID of the zone where the primary instance is deployed.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group. You can query available resource groups in the Resource Management console. If you do not specify this parameter, the default resource group is used.
* `security_ip_list` - (Optional) The IP addresses or CIDR blocks that you want to add to the whitelist of the ApsaraDB for HBase cluster. If you specify multiple IP addresses or CIDR blocks, separate them with commas (,).
* `standby_core_node_count` - (Optional) The standby core node count.
* `standby_vswitch_id` - (Required) The ID of the vSwitch that is specified for the secondary instance. The vSwitch must be deployed in the zone specified by the StandbyZoneId parameter.
* `standby_zone_id` - (Required, ForceNew) The ID of the zone where the secondary instance is deployed.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC). The VPC must be deployed in the region specified by the RegionId parameter.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Multi Zone Cluster.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Multi Zone Cluster.
* `delete` - (Defaults to 1 mins) Used when delete the Multi Zone Cluster.
* `update` - (Defaults to 1 mins) Used when update the Multi Zone Cluster.

## Import

HBase Multi Zone Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_hbase_multi_zone_cluster.example <id>
```
---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_zonal_db_cluster"
sidebar_current: "docs-alicloud-resource-polardb-zonal-db-cluster"
description: |-
  Provides a PolarDB zonal cluster resource.
---

# alicloud_polardb_zonal_db_cluster

Provides an PolarDB zonal cluster resource. An PolarDB zonal cluster is an isolated database
environment in the cloud. An PolarDB zonal cluster can contain multiple user-created
databases.

-> **NOTE:** Available since v1.261.0.

## Example Usage

Create a PolarDB MySQL zonal cluster

```terraform
variable "db_cluster_nodes_configs" {
  description = "The advanced configuration for all nodes in the cluster except for the RW node, including db_node_class, hot_replica_mode, and imci_switch properties."
  type = map(object({
    db_node_class    = string
    db_node_role     = optional(string, null)
    hot_replica_mode = optional(string, null)
    imci_switch      = optional(string, null)
  }))
  default = {
    db_node_1 = {
      db_node_class = "polar.mysql.x4.medium.c"
      db_node_role  = "Writer"
    }
    db_node_2 = {
      db_node_class = "polar.mysql.x4.medium.c"
      db_node_role  = "Reader"
    }
  }
}

resource "alicloud_ens_network" "default" {
  network_name = "terraform-example"

  description   = "LoadBalancerNetworkDescription_test"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "tr-Istanbul-1"
}

resource "alicloud_ens_vswitch" "default" {
  description  = "LoadBalancerVSwitchDescription_test"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "terraform-example"

  ens_region_id = "tr-Istanbul-1"
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_polardb_zonal_db_cluster" "default" {
  db_node_class = "polar.mysql.x4.medium.c"
  description   = "terraform-example"
  ens_region_id = "tr-Istanbul-1"
  vpc_id        = alicloud_ens_network.default.id
  vswitch_id    = alicloud_ens_vswitch.default.id
  db_cluster_nodes_configs = {
    for node, config in var.db_cluster_nodes_configs : node => jsonencode({ for k, v in config : k => v if v != null })
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Optional, Required, ForceNew) Database type. Value options: MySQL, Oracle, PostgreSQL.
* `db_version` - (Optional, Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `DBVersion`.
* `db_minor_version` - (Optional, ForceNew) Database minor version. Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `DBMinorVersion`. This parameter takes effect only when `db_type` is MySQL and `db_version` is 8.0.
* `db_node_class` - (Required) The db_node_class of cluster node.Only effective when the cluster is created for the first time. After the cluster is started, the cluster specification is maintained through the node class.
* `ens_region_id` - (Required, ForceNew) The Zone to launch the DB cluster.
* `pay_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when pay_type is `PrePaid`. Valid value:1, 2, 3, 6, 12, 24, 36, Default to 1.
* `used_time` - (Optional) The duration that you will buy DB cluster (in month). It is valid when pay_type is `PrePaid`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to `PostPaid`. Once effect, it will not be modified that means running `terraform apply` will not affect the resource.
* `vswitch_id` - (Required, ForceNew) The ENS virtual switch ID to launch DB instances in one VPC.
* `description` - (Optional, Computed) The description of cluster.
* `sub_category` - The category of the cluster. Valid values are `Exclusive`, `General`, Default to `Exclusive`.
* `creation_category` - (Optional, ForceNew) The edition of the PolarDB service. Valid values are `SENormal`.
* `vpc_id` - (Required, ForceNew) The id of the ENS VPC.
* `storage_type` - (Optional, ForceNew) The storage type of the cluster. Valid values are `ESSDPL1`, `ESSDPL0`.
* `storage_space` - (Optional) Storage space charged by space (monthly package). Unit: GB.
-> **NOTE:**  Valid values for PolarDB for MySQL Standard Edition: 20 to 32000. It is valid when pay_type are `PrePaid` ,`PostPaid`.
* `storage_pay_type` - The billing method of the storage. Valid values `Prepaid`.
* `upgrade_type` - Version upgrade type. Valid values are PROXY, DB, ALL. PROXY means upgrading the proxy version, DB means upgrading the db version, ALL means upgrading both db and proxy versions simultaneously.
* `planned_start_time` - The earliest time to start executing a scheduled (i.e. within the target time period) kernel version upgrade task. The format is YYYY-MM-DDThh: mm: ssZ (UTC).
-> **NOTE:** The starting time range is any time point within the next 24 hours. For example, the current time is 2021-01-14T09:00:00Z, and the allowed start time range for filling in here is 2021-01-14T09:00:00Z~2021-01-15T09:00:00Z. If this parameter is left blank, the kernel version upgrade task will be executed immediately by default.
-> **NOTE:** The latest time must be 30 minutes or more later than the start time. If PlannedStartTime is set but this parameter is not specified, the latest time to execute the target task defaults to the start time+30 minutes. For example, when the PlannedStartTime is set to 2021-01-14T09:00:00Z and this parameter is left blank, the target task will start executing at the latest on 2021-01-14T09:30:00Z.
* `target_minor_version` - (Optional) The Version Code of the target version, whose parameter values can be obtained from the [DescribeDBClusterVersionZonal](https://www.alibabacloud.com/help/en/polardb/api-polardb-2017-08-01-describedbclusterversionzonal) interface.
* `db_cluster_nodes_configs` - (Optional, Required) Map of node needs to be created after DB cluster was launched.
* `cluster_version` - (Optional, ForceNew, Computed) current DB Cluster revision Version.
* `storage_pay_type` - (Computed) The billing method of the storage. Valid values `Prepaid`.

## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB zonal cluster ID.
* `create_time` - PolarDB zonal cluster creation time.  
* `cluster_latest_version` - PolarDB zonal cluster latest version.
* `db_cluster_nodes_attributes` - Cache of the relationship between node key and node ID for PolarDB zonal Cluster
* `db_cluster_nodes_ids` - Cache of node ID for PolarDB zonal Cluster
* `region_id`  - PolarDB zonal cluster region

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the PolarDB zonal cluster (until it reaches the initial `Running` status).
* `update` - (Defaults to 38 mins) Used when updating the PolarDB zonal cluster (until it reaches the initial `Running` status).
* `delete` - (Defaults to 5 mins) Used when terminating the PolarDB zonal cluster.

## Import

PolarDB zonal cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_zonal_db_cluster.example pc-abc12345678
```

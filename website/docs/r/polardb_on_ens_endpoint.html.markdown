---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_on_ens_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-on-ens-endpoint"
description: |-
  Provides a PolarDB ON ENS instance endpoint resource.
---

# alicloud_polardb_on_ens_endpoint

Provides a PolarDB ON ENS endpoint resource to manage custom endpoint of PolarDB cluster.

-> **NOTE:** Available since v1.80.0.
-> **NOTE:** The primary endpoint and the default cluster endpoint can not be created or deleted manually.

## Example Usage

```terraform
variable "db_cluster_nodes_configs" {
  description = "The advanced configuration for all nodes in the cluster except for the RW node, including db_node_class, hot_replica_mode, and imci_switch properties."
  type        = map(object({
    db_node_class           = string
    db_node_role            = optional(string,null)
    hot_replica_mode        = optional(string,null)
    imci_switch             = optional(string,null)
  }))
  default     = {}
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

resource "alicloud_polardb_on_ens_cluster" "default" {
  db_node_class = "polar.mysql.x4.medium.c"
  description   = "terraform-example"
  ens_region_id = "tr-Istanbul-1"
  vpc_id = alicloud_ens_network.default.id
  vswitch_id    = alicloud_ens_vswitch.default.id
  db_cluster_nodes_configs = {
    for node, config in var.db_cluster_nodes_configs : node => jsonencode({for k, v in config : k => v if v != null})
  }
}

resource "alicloud_polardb_on_ens_endpoint" "default" {
  db_cluster_id           = alicloud_polardb_on_ens_cluster.default.id
  db_cluster_nodes_ids    = alicloud_polardb_on_ens_cluster.default.db_cluster_nodes_ids
  endpoint_config         = {"MasterAcceptReads":"on"}
  nodes_key               = ["node_reader_1"]
  read_write_mode         = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `endpoint_type` - (Optional, ForceNew) Type of the endpoint. Valid values are `Custom`, `Cluster`, `Primary`, default to `Custom`. However when creating a new endpoint, it also only can be `Custom`. 
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. When creating a new custom endpoint, default to `ReadOnly`.
* `nodes` - (Optional) Node id list for endpoint configuration. At least 2 nodes if specified, or if the cluster has more than 3 nodes, read-only endpoint is allowed to mount only one node. Default is all nodes.
* `auto_add_new_nodes` - (Optional) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. When creating a new custom endpoint, default to `Disable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).
* `net_type` - (Optional) The network type of the endpoint address.
* `db_endpoint_description` - (Optional) The name of the endpoint.
* `port` - (Optional) Port of the specified endpoint. Valid values: 3000 to 5999.
* `vpc_id` - (Optional) The ID of ENS VPC where to use the DB.  
* `vswitch_id` - (Optional) The ID of ENS virtual switch where to use the DB.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `db_endpoint_id` - The ID of the cluster endpoint.

## Import

PolarDB ON ENS endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_on_ens_endpoint.example pc-abc123456:pe-abc123456
```

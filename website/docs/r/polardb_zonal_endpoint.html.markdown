---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_zonal_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-on-ens-endpoint"
description: |-
  Provides a PolarDB Zonal instance endpoint resource.
---

# alicloud_polardb_zonal_endpoint

Provides a PolarDB Zonal endpoint resource to manage custom endpoint of PolarDB cluster.

-> **NOTE:** Available since v1.262.0.
-> **NOTE:** The primary endpoint and the default cluster endpoint can not be created or deleted manually.

## Example Usage

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

resource "alicloud_polardb_zonal_endpoint" "default" {
  db_cluster_id        = alicloud_polardb_zonal_db_cluster.default.id
  db_cluster_nodes_ids = alicloud_polardb_zonal_db_cluster.default.db_cluster_nodes_ids
  endpoint_config      = {}
  nodes_key            = ["db_node_1", "db_node_2"]
  read_write_mode      = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `endpoint_type` - (Optional, Computed, ForceNew) Type of the endpoint. Valid values are `Custom`, `Cluster`, `Primary`, default to `Custom`. However when creating a new endpoint, it also only can be `Custom`. 
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. When creating a new custom endpoint, default to `ReadOnly`.
* `nodes` - (Computed) Node id list for endpoint configuration.
* `db_cluster_nodes_ids` - (Required) referenced from the db_cluster_nodes_ids attribute of alicloud_polardb_zonal_db_cluster.. 
* `nodes_key` - (Optional) The list of backend nodes for the endpoint, with the attribute values derived from the map key of db_cluster_nodes_ids.
* `auto_add_new_nodes` - (Optional, Computed) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. When creating a new custom endpoint, default to `Enable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).
* `connection_prefix` - (Computed) Prefix of the specified endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter.
* `net_type` - (Optional, Computed, ForceNew) The network type of the endpoint address.
* `db_endpoint_description` - (Optional) The name of the endpoint.
* `port` - (Computed) Port of the specified endpoint. Valid values: 3000 to 5999.
* `vpc_id` - (Optional, ForceNew, Computed) The ID of ENS VPC where to use the DB.  
* `vswitch_id` - (Optional, ForceNew, Computed) The ID of ENS virtual switch where to use the DB.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `db_endpoint_id` - The ID of the cluster endpoint.

## Import

PolarDB Zonal endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_zonal_endpoint.example pc-abc123456:pe-abc123456
```

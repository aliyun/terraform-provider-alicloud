---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_on_ens_database"
sidebar_current: "docs-alicloud-resource-polardb-on-ens-database"
description: |-
  Provides a PolarDB ON ENS database resource.
---

# alicloud_polardb_on_ens_database

Provides a PolarDB ON ENS database resource. A database deployed in a PolarDB ON ENS cluster. A PolarDB ON ENS cluster can own multiple databases.

-> **NOTE:** Available since v1.66.0.

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

resource "alicloud_polardb_on_ens_database" "default" {
  db_cluster_id         = alicloud_polardb_on_ens_cluster.default.id
  db_description        = "terraform-example"
  db_name               = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `db_name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letterand have no more than 64 characters.
* `character_set_name` - (Optional, ForceNew) Character set. The value range is limited to the following: [ utf8, gbk, latin1, utf8mb4, Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ], default is "utf8" \(`utf8mb4` only supports versions 5.5 and 5.6\).
* `db_description` - (Optional) Database description. It must start with a Chinese character or English letter, cannot start with "http://" or "https://". It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length must be 2-256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of cluster ID and database name with format `<cluster_id>:<name>`.

## Import

PolarDB database can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_on_ens_database.example "pc-12345:tf_database"
```

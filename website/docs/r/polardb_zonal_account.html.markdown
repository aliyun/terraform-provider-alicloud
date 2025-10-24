---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_zonal_account"
sidebar_current: "docs-alicloud-resource-polardb-zonal-account"
description: |-
  Provides a PolarDB Zonal account resource.
---

# alicloud_polardb_zonal_account

Provides a PolarDB Zonal account resource and used to manage databases.

-> **NOTE:** Available since v1.262.0. 

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

resource "alicloud_polardb_zonal_account" "default" {
  db_cluster_id    = alicloud_polardb_zonal_db_cluster.default.id
  account_name     = "terraform_example"
  account_password = "Example1234"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `account_password` - (Optional, Sensitive, Computed) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `account_description` - (Optional) Account description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_type` - (Optional, ForceNew, Computed) Account type, Valid values are `Normal`, `Super`, Default to `Normal`.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

PolarDB Zonal account can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_zonal_account.example "pc-12345:tf_account"
```

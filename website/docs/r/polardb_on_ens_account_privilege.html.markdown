---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_on_ens_account_privilege"
sidebar_current: "docs-alicloud-resource-polardb-on-ens-account-privilege"
description: |-
  Provides a PolarDB ON ENS account privilege resource.
---

# alicloud_polardb_on_ens_account_privilege

Provides a PolarDB ON ENS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

-> **NOTE:** Available in v1.67.0+.

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
  ens_region_id = "vn-hanoi-3"
}

resource "alicloud_ens_vswitch" "default" {
  description  = "LoadBalancerVSwitchDescription_test"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "terraform-example"

  ens_region_id = "vn-hanoi-3"
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_polardb_on_ens_cluster" "default" {
  db_node_class = "polar.mysql.x4.medium.c"
  description   = "terraform-example"
  ens_region_id = "vn-hanoi-3"
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

resource "alicloud_polardb_on_ens_account" "default" {
  db_cluster_id          = alicloud_polardb_on_ens_cluster.default.id
  account_name           = "terraform-example"
  account_password       = "Example1234"
}

resource "alicloud_polardb_on_ens_account_privilege" "default" {
  db_cluster_id         = alicloud_polardb_on_ens_cluster.default.id
  account_name          = alicloud_polardb_on_ens_account.default.account_name
  account_privilege     = "ReadWrite"
  db_names              = [alicloud_polardb_on_ens_database.default.db_name]
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `account_privilege` - (Optional, ForceNew) The privilege of one account access database. Valid values: ["ReadOnly", "ReadWrite"], ["DMLOnly", "DDLOnly"] added since version v1.101.0. Default to "ReadOnly".
* `db_names` - (Required) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<db_cluster_id>:<account_name>:<account_privilege>`.

## Import

PolarDB ON ENS account privilege can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_on_ens_account_privilege.example "pc-12345:tf_account:ReadOnly"
```

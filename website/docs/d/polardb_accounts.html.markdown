---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_accounts"
sidebar_current: "docs-alicloud-datasource-polardb-accounts"
description: |-
    Provides a collection of PolarDB endpoints according to the specified filters.
---

# alicloud\_polardb\_accounts

The `alicloud_polardb_accounts` data source provides a collection of PolarDB cluster database account available in Alibaba Cloud account.
Filters support regular expression for the account name, searches by clusterId.

-> **NOTE:** Available since v1.70.0+.

## Example Usage

```terraform

data "alicloud_polardb_node_classes" "this" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PrePaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.this.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_count = "2"
  db_node_class = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
  vswitch_id    = alicloud_vswitch.default.id
}

data "alicloud_polardb_clusters" "polardb_clusters_ds" {
  description_regex = alicloud_polardb_cluster.cluster.description
  status            = "Running"
}

resource "alicloud_polardb_account" "account" {
  db_cluster_id       = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
  account_name        = "tfnormal_01"
  account_password    = "Test12345"
  account_description = "tf_account_description"
  account_type        = "Normal"
}

data "alicloud_polardb_accounts" "default" {
  db_cluster_id = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
  name_regex    = alicloud_polardb_account.account.account_name
}

output "account" {
  value = data.alicloud_polardb_accounts.default.accounts[0].account_name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required) The polarDB cluster ID. 
* `name_regex` - (Optional) A regex string to filter results by account name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - Account name of the cluster.
* `accounts` - A list of PolarDB cluster accounts. Each element contains the following attributes:
  * `account_description` - Account description.
  * `account_lock_state` - Account lock state, Valid values are `Lock`, `UnLock`.
  * `account_name` - Account name.
  * `account_status` - Cluster address type.`Cluster`: the default address of the Cluster.`Primary`: Primary address.`Custom`: Custom cluster addresses.
  * `account_type` - Account type, Valid values are `Normal`, `Super`.
  * `database_privileges` - A list of database privilege. Each element contains the following attributes.
      * `account_privilege` - Account privilege of database
      * `db_name` - The account owned database name 

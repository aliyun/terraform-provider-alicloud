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

-> **NOTE:** Available in v1.70.0+.

## Example Usage

```terraform
data "alicloud_polardb_clusters" "polardb_clusters_ds" {
  description_regex = "pc-\\w+"
  status            = "Running"
}

data "alicloud_polardb_accounts" "default" {
  db_cluster_id = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
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

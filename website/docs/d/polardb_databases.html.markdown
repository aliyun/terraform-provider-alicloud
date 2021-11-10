---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_databases"
sidebar_current: "docs-alicloud-datasource-polardb-databases"
description: |-
    Provides a collection of PolarDB endpoints according to the specified filters.
---

# alicloud\_polardb\_databases

The `alicloud_polardb_databases` data source provides a collection of PolarDB cluster database available in Alibaba Cloud account.
Filters support regular expression for the database name, searches by clusterId.

-> **NOTE:** Available in v1.70.0+.

## Example Usage

```terraform
data "alicloud_polardb_clusters" "polardb_clusters_ds" {
  description_regex = "pc-\\w+"
  status            = "Running"
}

data "alicloud_polardb_databases" "default" {
  db_cluster_id = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
}

output "database" {
  value = data.alicloud_polardb_databases.default.databases[0].db_name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required) The polarDB cluster ID. 
* `name_regex` - (Optional) A regex string to filter results by database name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - database name of the cluster.
* `databases` - A list of PolarDB cluster databases. Each element contains the following attributes:
  * `character_set_name` - The character set name of database.
  * `db_description` - Database description.
  * `db_name` - Database name.
  * `db_status` - The status of database.
  * `engine` - The engine of database.
  * `accounts` - A list of accounts of database. Each element contains the following attributes.
      * `account_name` - Account name.
      * `account_status` - Account status.
      * `privilege_status` - The privilege status of account.

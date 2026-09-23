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

-> **NOTE:** Available since v1.70.0+.

## Example Usage

```terraform
data "alicloud_polardb_node_classes" "this" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
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

resource "alicloud_polardb_database" "default" {
  db_cluster_id  = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
  db_name        = "tfaccountpri_${data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id}"
  db_description = "from terraform"
}

data "alicloud_polardb_databases" "default" {
  db_cluster_id = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
  name_regex    = alicloud_polardb_database.default.db_name
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

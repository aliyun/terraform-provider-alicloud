---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_global_database_networks"
sidebar_current: "docs-alicloud-datasource-polardb-global-database-networks"
description: |-
  Provides a list of PolarDB Global Database Networks to the user.
---

# alicloud\_polardb\_global\_database\_networks

This data source provides the PolarDB Global Database Networks of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.181.0+.

## Example Usage

Basic Usage

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

resource "alicloud_polardb_global_database_network" "default" {
  db_cluster_id = alicloud_polardb_cluster.cluster.id
  description   = alicloud_polardb_cluster.cluster.id
}
data "alicloud_polardb_global_database_networks" "ids" {
  ids = [alicloud_polardb_global_database_network.default.id]
}
output "polardb_global_database_network_id_1" {
  value = data.alicloud_polardb_global_database_networks.ids.networks.0.id
}

data "alicloud_polardb_global_database_networks" "description" {
  description = alicloud_polardb_cluster.cluster.id
}
output "polardb_global_database_network_id_2" {
  value = data.alicloud_polardb_global_database_networks.description.networks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Global Database Network IDs.
* `gdn_id` - (Optional, ForceNew)  The ID of the Global Database Network.
* `db_cluster_id` - (Optional, ForceNew) The ID of the cluster.
* `description` - (Optional, ForceNew, Computed) The description of the Global Database Network.
* `status` - (Optional, ForceNew) The status of the Global Database Network. Valid values:
	- `creating`: The Global Database Network is being created.
	- `active`: The Global Database Network is running.
	- `deleting`: The Global Database Network is being deleted.
	- `locked`: The Global Database Network is locked. If the Global Database Network is locked, you cannot perform operations on clusters in the Global Database Network.
	- `removing_member`: The secondary cluster is being removed from the Global Database Network.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `networks` - A list of PolarDB Global Database Networks. Each element contains the following attributes:
	* `id` - The ID of the Global Database Network.
	* `gdn_id` - The ID of the Global Database Network.
	* `description` - The description of the Global Database Network.
	* `db_type` - The type of the database engine. Only MySQL is supported.
	* `db_version` - The version number of the database engine. Only the 8.0 version is supported.
	* `create_time` - The time when the Global Database Network was created. The time is in the YYYY-MM-DDThh:mm:ssZ format. The time is displayed in UTC.
	* `status` - The status of the Global Database Network.
	* `db_clusters` - The details of each cluster in the Global Database Network.
		* `db_cluster_id` - The ID of the PolarDB cluster.
		* `role` - The role of the cluster.
		* `region_id` - The region ID of the cluster.
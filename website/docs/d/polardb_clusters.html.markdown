---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_clusters"
sidebar_current: "docs-alicloud-datasource-polardb-clusters"
description: |-
    Provides a collection of PolarDB clusters according to the specified filters.
---

# alicloud\_polardb\_clusters

The `alicloud_polardb_clusters` data source provides a collection of PolarDB clusters available in Alibaba Cloud account.
Filters support regular expression for the cluster description, searches by tags, and other filters which are listed below.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

```terraform
data "alicloud_polardb_clusters" "polardb_clusters_ds" {
  description_regex = "pc-\\w+"
  status            = "Running"
}

output "first_polardb_cluster_id" {
  value = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to filter results by cluster description.
* `ids` - (Optional) A list of PolarDB cluster IDs. 
* `status` - (Optional) status of the cluster.
* `db_type` - (Optional) Database type. Options are `MySQL`, `Oracle` and `PostgreSQL`. If no value is specified, all types are returned.
* `tags` - (Optional, Available in v1.68.0+) A mapping of tags to assign to the resource.
      - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
      - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of RDS cluster IDs. 
* `descriptions` - A list of RDS cluster descriptions. 
* `clusters` - A list of PolarDB clusters. Each element contains the following attributes:
  * `id` - The ID of the PolarDB cluster.
  * `description` - The description of the PolarDB cluster.
  * `charge_type` - Billing method. Value options: `PostPaid` for Pay-As-You-Go and `PrePaid` for subscription.
  * `network_type` - The DBClusterNetworkType of the PolarDB cluster.
  * `region_id` - Region ID the cluster belongs to.
  * `zone_id` - The ZoneId of the PolarDB cluster.
  * `expire_time` - Expiration time. Pay-As-You-Go clusters never expire.
  * `expired` - The expired of the PolarDB cluster.
  * `status` - Status of the cluster.
  * `engine` - Database type. Options are `MySQL`, `Oracle` and `PostgreSQL`. If no value is specified, all types are returned.
  * `db_type` - `Primary` for primary cluster, `ReadOnly` for read-only cluster, `Guard` for disaster recovery cluster, and `Temp` for temporary cluster.
  * `db_version` - The DBVersion of the PolarDB cluster.
  * `lock_mode` - The LockMode of the PolarDB cluster.
  * `delete_lock` - The DeleteLock of the PolarDB cluster.
  * `create_time` - The CreateTime of the PolarDB cluster.
  * `vpc_id` - ID of the VPC the cluster belongs to.
  * `db_node_number` - The DBNodeNumber of the PolarDB cluster.
  * `db_node_class` - The DBNodeClass of the PolarDB cluster.
  * `storage_used` - The StorageUsed of the PolarDB cluster.
  * `db_nodes` - The DBNodes of the PolarDB cluster.
    * `db_node_class` - The db_node_class of the db_nodes.
    * `max_iops` - The max_iops of the db_nodes.
    * `region_id` - The region_id of the db_nodes.
    * `db_node_role` - The db_node_role of the db_nodes.
    * `max_connections` - The max_connections of the db_nodes.
    * `zone_id` - The zone_id of the db_nodes.
    * `db_node_status` - The db_node_status of the db_nodes.
    * `db_node_id` - The db_node_id of the db_nodes.
    * `create_time` - The create_time of the db_nodes.

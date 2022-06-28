---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_clusters"
sidebar_current: "docs-alicloud-datasource-adb-db-clusters"
description: |-
  Provides a list of Adb DBClusters to the user.
---

# alicloud\_adb\_db\_clusters

This data source provides the Adb DBClusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_adb_db_clusters" "example" {
  description_regex = "example"
}

output "first_adb_db_cluster_id" {
  value = data.alicloud_adb_db_clusters.example.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of DBCluster.
* `description_regex` - (Optional, ForceNew) A regex string to filter results by DBCluster description.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of DBCluster IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The status of the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `descriptions` - A list of DBCluster descriptions.
* `clusters` - A list of Adb Db Clusters. Each element contains the following attributes:
	* `auto_renew_period` - Auto-renewal period of an cluster, in the unit of the month.
	* `commodity_code` - The name of the service.
	* `compute_resource` - The specifications of computing resources in elastic mode. The increase of resources can speed up queries. AnalyticDB for MySQL automatically scales computing resources. For more information, see [Specifications](https://www.alibabacloud.com/help/en/doc-detail/144851.htm).
	* `connection_string` - The endpoint of the cluster.
	* `create_time` - The CreateTime of the ADB cluster.
	* `db_cluster_category` - The db cluster category.
	* `db_cluster_id` - The db cluster id.
	* `db_cluster_network_type` - The db cluster network type.
	* `db_cluster_type` - The db cluster type.
	* `db_cluster_version` - The db cluster version.
	* `db_node_class` - The db node class.
	* `db_node_count` - The db node count.
	* `db_node_storage` - The db node storage.
	* `description` - The description of DBCluster.
	* `disk_type` - The type of the disk.
	* `dts_job_id` - The ID of the data synchronization task in Data Transmission Service (DTS). This parameter is valid only for analytic instances.
	* `elastic_io_resource` - The elastic io resource.
	* `engine` - The engine of the database.
	* `engine_version` - The engine version of the database..
	* `executor_count` - The number of nodes. The node resources are used for data computing in elastic mode.
	* `expire_time` - The time when the cluster expires.
	* `expired` - Indicates whether the cluster has expired.
	* `id` - The ID of the DBCluster.
	* `lock_mode` - The lock mode of the cluster.
	* `lock_reason` - The reason why the cluster is locked.
	* `maintain_time` - The maintenance window of the cluster.
	* `payment_type` - The payment type of the resource.
	* `charge_type` - The payment type of the resource.
	* `port` - The port that is used to access the cluster.
	* `rds_instance_id` - The ID of the ApsaraDB RDS instance from which data is synchronized to the cluster. This parameter is valid only for analytic instances.
	* `renewal_status` - The status of renewal.
	* `resource_group_id` - The ID of the resource group.
	* `security_ips` - List of IP addresses allowed to access all databases of an cluster.
	* `status` - The status of the resource.
	* `storage_resource` - The specifications of storage resources in elastic mode. The resources are used for data read and write operations. The increase of resources can improve the read and write performance of your cluster. For more information, see [Specifications](https://www.alibabacloud.com/help/en/doc-detail/144851.htm).
	* `tags` - The tag of the resource.
		* `tag_key` - The key of the tags.
		* `tag_value` - The value of the tags.
	* `vpc_cloud_instance_id` - The vpc cloud instance id.
	* `vpc_id` - The vpc id.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The zone ID  of the resource.
	* `mode` - The lock mode of the cluster.	

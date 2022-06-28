---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_clusters"
sidebar_current: "docs-alicloud-datasource-adb-clusters"
description: |-
    Provides a collection of ADB clusters according to the specified filters.
---

# alicloud\_adb\_clusters

The `alicloud_adb_clusters` data source provides a collection of ADB clusters available in Alibaba Cloud account.
Filters support regular expression for the cluster description, searches by tags, and other filters which are listed below.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.121.0`. Please use new datasource [alicloud_adb_db_clusters](https://www.terraform.io/docs/providers/alicloud/d/adb_db_clusters).

-> **NOTE:** Available in v1.71.0+.

## Example Usage

```
data "alicloud_adb_clusters" "adb_clusters_ds" {
  description_regex = "am-\\w+"
  status     = "Running"
}

output "first_adb_cluster_id" {
  value = data.alicloud_adb_clusters.adb_clusters_ds.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to filter results by cluster description.
* `ids` - (Optional) A list of ADB cluster IDs. 
* `status` - (Optional, ForceNew, Available in v1.102.0+) The status of the cluster. Valid values: `Preparing`, `Creating`, `Restoring`, `Running`, `Deleting`, `ClassChanging`, `NetAddressCreating`, `NetAddressDeleting`. For more information, see [Cluster status](https://www.alibabacloud.com/help/doc-detail/143075.htm).
* `tags` - (Optional, Available in v1.68.0+) A mapping of tags to assign to the resource.
      - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
      - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ADB cluster IDs. 
* `descriptions` - A list of ADB cluster descriptions. 
* `clusters` - A list of ADB clusters. Each element contains the following attributes:
  * `id` - The ID of the ADB cluster.
  * `description` - The description of the ADB cluster.
  * `charge_type` - Billing method. Value options: `PostPaid` for Pay-As-You-Go and `PrePaid` for subscription.
  * `network_type` - The DBClusterNetworkType of the ADB cluster.
  * `region_id` - Region ID the cluster belongs to.
  * `zone_id` - The ZoneId of the ADB cluster.
  * `expire_time` - Expiration time. Pay-As-You-Go clusters never expire.
  * `expired` - The expired of the ADB cluster.
  * `status` - Status of the cluster.
  * `lock_mode` - The LockMode of the ADB cluster.
  * `create_time` - The CreateTime of the ADB cluster.
  * `vpc_id` - ID of the VPC the cluster belongs to.
  * `db_node_count` - The DBNodeCount of the ADB cluster.
  * `db_node_class` - The DBNodeClass of the ADB cluster.
  * `db_node_storage` - The DBNodeStorage of the ADB cluster.

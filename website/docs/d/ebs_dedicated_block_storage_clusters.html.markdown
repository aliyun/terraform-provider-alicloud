---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_dedicated_block_storage_clusters"
sidebar_current: "docs-alicloud-datasource-ebs-dedicated-block-storage-clusters"
description: |-
  Provides a list of Ebs Dedicated Block Storage Cluster owned by an Alibaba Cloud account.
---

# alicloud_ebs_dedicated_block_storage_clusters

This data source provides Ebs Dedicated Block Storage Cluster available to the user.

-> **NOTE:** Available in 1.196.0+

## Example Usage

```terraform
data "alicloud_ebs_dedicated_block_storage_clusters" "default" {
  ids        = ["example_id"]
  name_regex = alicloud_ebs_dedicated_block_storage_cluster.default.name
}

output "alicloud_ebs_dedicated_block_storage_cluster_example_id" {
  value = data.alicloud_ebs_dedicated_block_storage_clusters.default.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Dedicated Block Storage Cluster IDs.
* `dedicated_block_storage_cluster_names` - (Optional, ForceNew) The name of the Dedicated Block Storage Cluster. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Dedicated Block Storage Cluster IDs.
* `names` - A list of name of Dedicated Block Storage Clusters.
* `clusters` - A list of Dedicated Block Storage Cluster Entries. Each element contains the following attributes:
  * `available_capacity` - The available capacity of the dedicated block storage cluster. Unit: GiB.
  * `category` - The type of cloud disk that can be created by a dedicated block storage cluster.
  * `create_time` - The creation time of the resource
  * `dedicated_block_storage_cluster_id` - The first ID of the resource
  * `dedicated_block_storage_cluster_name` - The name of the resource
  * `delivery_capacity` - Capacity to be delivered in GB.
  * `description` - The description of the dedicated block storage cluster.
  * `expired_time` - The expiration time of the dedicated block storage cluster, in the Unix timestamp format, in seconds.
  * `performance_level` - Cloud disk performance level, possible values:-PL0.-PL1.-PL2.-PL3.> Only valid in SupportedCategory = cloud_essd.
  * `resource_group_id` - The ID of the resource group
  * `status` - The status of the resource
  * `supported_category` - This parameter is not supported.
  * `total_capacity` - The total capacity of the dedicated block storage cluster. Unit: GiB.
  * `type` - The dedicated block storage cluster performance type. Possible values:-Standard: Basic type. This type of dedicated block storage cluster can create an ESSD PL0 cloud disk.-Premium: performance type. This type of dedicated block storage cluster can create an ESSD PL1 cloud disk.
  * `used_capacity` - The used (created disk) capacity of the current cluster, in GB
  * `zone_id` - The zone ID  of the resource

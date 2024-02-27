---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_dedicated_block_storage_cluster"
sidebar_current: "docs-alicloud-resource-ebs-dedicated-block-storage-cluster"
description: |-
  Provides a Alicloud Ebs Dedicated Block Storage Cluster resource.
---

# alicloud_ebs_dedicated_block_storage_cluster

Provides a Ebs Dedicated Block Storage Cluster resource.

For information about Ebs Dedicated Block Storage Cluster and how to use it, see [What is Dedicated Block Storage Cluster](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ebs-2021-07-30-creatededicatedblockstoragecluster).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ebs_dedicated_block_storage_cluster" "default" {
  type                                 = "Premium"
  zone_id                              = "cn-heyuan-b"
  dedicated_block_storage_cluster_name = "dedicated_block_storage_cluster_name"
  total_capacity                       = 61440
  region_id                            = "cn-heyuan"
}
```

## Argument Reference

The following arguments are supported:
* `dedicated_block_storage_cluster_name` - (Required) The name of the resource
* `description` - (Computed,Optional) The description of the dedicated block storage cluster.
* `total_capacity` - (Required,ForceNew) The total capacity of the dedicated block storage cluster. Unit: GiB.
* `type` - (Required,ForceNew) The dedicated block storage cluster performance type. Possible values:-Standard: Basic type. This type of dedicated block storage cluster can create an ESSD PL0 cloud disk.-Premium: performance type. This type of dedicated block storage cluster can create an ESSD PL1 cloud disk.
* `zone_id` - (Required,ForceNew) The zone ID  of the resource



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `available_capacity` - The available capacity of the dedicated block storage cluster. Unit: GiB.
* `category` - The type of cloud disk that can be created by a dedicated block storage cluster.
* `create_time` - The creation time of the resource
* `dedicated_block_storage_cluster_id` - The first ID of the resource
* `delivery_capacity` - Capacity to be delivered in GB.
* `description` - The description of the dedicated block storage cluster.
* `expired_time` - The expiration time of the dedicated block storage cluster, in the Unix timestamp format, in seconds.
* `performance_level` - Cloud disk performance level, possible values:-PL0.-PL1.-PL2.-PL3.> Only valid in SupportedCategory = cloud_essd.
* `status` - The status of the resource
* `supported_category` - This parameter is not supported.
* `used_capacity` - The used (created disk) capacity of the current cluster, in GB
* `resource_group_id` - The ID of the resource group

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Dedicated Block Storage Cluster.
* `update` - (Defaults to 5 mins) Used when update the Dedicated Block Storage Cluster.
* `delete` - (Defaults to 5 mins) Used when update the Dedicated Block Storage Cluster.

## Import

Ebs Dedicated Block Storage Cluster can be imported using the id, e.g.

```shell
$terraform import alicloud_disk_dedicated_block_storage_cluster.example <id>
```
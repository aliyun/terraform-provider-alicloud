---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_dedicated_host_clusters"
sidebar_current: "docs-alicloud-datasource-ecs-dedicated-host-clusters"
description: |-
  Provides a list of Ecs Dedicated Host Clusters to the user.
---

# alicloud\_ecs\_dedicated\_host\_clusters

This data source provides the Ecs Dedicated Host Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_dedicated_host_clusters" "ids" {
  ids = ["example_id"]
}
output "ecs_dedicated_host_cluster_id_1" {
  value = data.alicloud_ecs_dedicated_host_clusters.ids.clusters.0.id
}

data "alicloud_ecs_dedicated_host_clusters" "nameRegex" {
  name_regex = "^my-DedicatedHostCluster"
}
output "ecs_dedicated_host_cluster_id_2" {
  value = data.alicloud_ecs_dedicated_host_clusters.nameRegex.clusters.0.id
}

data "alicloud_ecs_dedicated_host_clusters" "zoneId" {
  zone_id = "example_value"
}
output "ecs_dedicated_host_cluster_id_3" {
  value = data.alicloud_ecs_dedicated_host_clusters.zoneId.clusters.0.id
}

data "alicloud_ecs_dedicated_host_clusters" "clusterName" {
  dedicated_host_cluster_name = "example_value"
}
output "ecs_dedicated_host_cluster_id_4" {
  value = data.alicloud_ecs_dedicated_host_clusters.clusterName.clusters.0.id
}

data "alicloud_ecs_dedicated_host_clusters" "clusterIds" {
  dedicated_host_cluster_ids = ["example_id"]
}
output "ecs_dedicated_host_cluster_id_5" {
  value = data.alicloud_ecs_dedicated_host_clusters.clusterIds.clusters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `dedicated_host_cluster_ids` - (Optional, ForceNew) The IDs of dedicated host clusters.
* `dedicated_host_cluster_name` - (Optional, ForceNew) The name of the dedicated host cluster.
* `ids` - (Optional, ForceNew, Computed)  A list of Dedicated Host Cluster IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Dedicated Host Cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `zone_id` - (Optional, ForceNew) The zone ID of the dedicated host cluster.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dedicated Host Cluster names.
* `clusters` - A list of Ecs Dedicated Host Clusters. Each element contains the following attributes:
  * `dedicated_host_cluster_id` - The ID of the dedicated host cluster.
  * `dedicated_host_cluster_name` - The name of the dedicated host cluster.
  * `description` - The description of the dedicated host cluster.
  * `id` - The ID of the Dedicated Host Cluster.
  * `resource_group_id` - The ID of the resource group to which the dedicated host cluster belongs.
  * `zone_id` - The zone ID of the dedicated host cluster.
  * `tags` - A mapping of tags to assign to the resource.
  * `dedicated_host_ids` - The IDs of dedicated hosts in the dedicated host cluster.
  * `dedicated_host_cluster_capacity` - The capacity of the dedicated host cluster.
    * `available_memory` - The available memory size. Unit: `GiB`.
    * `available_vcpus` - The number of available vCPUs.
    * `total_memory` - The total memory size. Unit: `GiB`.
    * `total_vcpus` - The total number of vCPUs.
    * `local_storage_capacities` - The local storage.
      * `available_disk` - The available capacity of the local disk. Unit: `GiB`.
      * `data_disk_category` - The category of the data disk. Valid values:`cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`, `cloud_essd`.
      * `total_disk` - The total capacity of the local disk. Unit: `GiB`.
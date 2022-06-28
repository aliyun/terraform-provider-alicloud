---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_hpc_clusters"
sidebar_current: "docs-alicloud-datasource-ecs-hpc-clusters"
description: |-
  Provides a list of Ecs Hpc Clusters to the user.
---

# alicloud\_ecs\_hpc\_clusters

This data source provides the Ecs Hpc Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_hpc_clusters" "example" {
  ids        = ["hpc-bp1i09xxxxxxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_hpc_cluster_id" {
  value = data.alicloud_ecs_hpc_clusters.example.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Hpc Cluster IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Hpc Cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Hpc Cluster names.
* `clusters` - A list of Ecs Hpc Clusters. Each element contains the following attributes:
	* `description` - The description of ECS Hpc Cluster.
	* `hpc_cluster_id` - The ID of the Hpc Cluster.
	* `id` - The ID of the Hpc Cluster.
	* `name` - The name of ECS Hpc Cluster.

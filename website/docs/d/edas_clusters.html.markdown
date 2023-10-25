---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_clusters"
sidebar_current: "docs-alicloud-datasource-edas-clusters"
description: |-
    Provides a list of EDAS clusters available to the user.
---

# alicloud\_edas\_clusters

This data source provides a list of EDAS clusters in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.82.0+

## Example Usage

```
data "alicloud_edas_clusters" "clusters" {
  logical_region_id = "cn-shenzhen:xxx"
  ids   =   ["addfs-dfsasd"]
  output_file = "clusters.txt"
}

output "first_cluster_name" {
  value = data.alicloud_alikafka_consumer_groups.clusters.clusters.0.cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `logical_region_id` - (Required) ID of the namespace in EDAS.
* `ids` - (Optional) An ids string to filter results by the cluster id. 
* `name_regex` - (Optional) A regex string to filter results by the cluster name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of cluster IDs.
* `names` - A list of cluster names.
* `clusters` - A list of clusters.
  * `cluster_id` - The ID of the cluster that you want to create the application.
  * `cluster_name` - The name of the cluster.
  * `cluster_type` - The type of the cluster, Valid values: 1: Swarm cluster. 2: ECS cluster. 3: Kubernetes cluster.
  * `create_time` - Cluster's creation time.
  * `update_time` - The time when the cluster was last updated.
  * `cpu` - The total number of CPUs in the cluster.
  * `cpu_used` - The number of used CPUs in the cluster.
  * `mem` - The total amount of memory in the cluser. Unit: MB.
  * `mem_used` - The amount of used memory in the cluser. Unit: MB.
  * `network_mode` - The network type of the cluster. Valid values: 1: classic network. 2: VPC.
  * `node_num` - The number of the Elastic Compute Service (ECS) instances that are deployed to the cluster.
  * `vpc_id` - The ID of the Virtual Private Cloud (VPC) for the cluster.
  * `region_id` - The ID of the namespace the application belongs to.


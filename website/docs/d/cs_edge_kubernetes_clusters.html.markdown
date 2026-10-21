---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_edge_kubernetes_clusters"
sidebar_current: "docs-alicloud-datasource-cs-edge-kubernetes-clusters"
description: |-
  Provides a list of Container Service Edge Kubernetes Clusters to be used by the alicloud_cs_edge_kubernetes_clusters resource.
---

# alicloud_cs_edge_kubernetes_clusters

-> **DEPRECATION NOTICE:** This data source has been deprecated since v1.276.0 and will be removed in a future release. Please use `alicloud_cs_clusters` instead.

This data source provides a list Container Service Edge Kubernetes Clusters on Alibaba Cloud.

-> **NOTE:** Available since v1.103.0.

## Example Usage

```terraform
# Declare the data source
data "alicloud_cs_edge_kubernetes_clusters" "k8s_clusters" {
  name_regex  = "my-first-k8s"
  output_file = "my-first-k8s-json"
}

output "output" {
  value = data.alicloud_cs_edge_kubernetes_clusters.k8s_clusters.clusters
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) Cluster IDs to filter.
* `name_regex` - (Optional) A regex string to filter results by cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Boolean, false by default, only `id` and `name` are exported. Set to true if more details are needed, e.g., `vpc_id`, `security_group_id`, `connections`. See full list in attributes.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Kubernetes clusters' ids.
* `names` - A list of matched Kubernetes clusters' names.
* `clusters` - A list of matched Kubernetes clusters.

### `clusters`

The clusters supports the following:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_nodes` - List of cluster worker nodes.
* `connections` - A list of kubernetes cluster connection information.

### `clusters.worker_nodes`

The worker_nodes supports the following:

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.

### `clusters.connections`

The connections supports the following:

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.

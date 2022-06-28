---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_edge_kubernetes_clusters"
sidebar_current: "docs-alicloud-datasource-cs-edge-kubernetes-clusters"
description: |-
  Provides a list of Container Service Edge Kubernetes Clusters to be used by the alicloud_cs_edge_kubernetes_clusters resource.
---

# alicloud\_cs\_edge\_kubernetes\_clusters

This data source provides a list Container Service Edge Kubernetes Clusters on Alibaba Cloud.

-> **NOTE:** Available in v1.103.0+

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
* `enabled_details` - (Optional) Boolean, false by default, only `id` and `name` are exported. Set to true if more details are needed, e.g., `master_disk_category`, `slb_internet_enabled`, `connections`. See full list in attributes.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Kubernetes clusters' ids.
* `names` - A list of matched Kubernetes clusters' names.
* `clusters` - A list of matched Kubernetes clusters. Each element contains the following attributes:
  * `id` - The ID of the container cluster.
  * `name` - The name of the container cluster.
  * `availability_zone` - The ID of availability zone.
  * `key_name` - The keypair of ssh login cluster node, you have to create it first.
  * `worker_numbers` - The ECS instance node number in the current container cluster.
  * `vswitch_ids` - The ID of VSwitches where the current cluster is located.
  * `vpc_id` - The ID of VPC where the current cluster is located.
  * `security_group_id` - The ID of security group where the current cluster worker node is located.
  * `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
  * `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
  * `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
  * `log_config` - A list of one element containing information about the associated log store. It contains the following attributes:
    * `type` - Type of collecting logs.
    * `project` - Log Service project name. 

### Block Nodes

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.

### Block Connections

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.

---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_clusters"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-clusters"
description: |-
  Provides a list of Container Service Kubernetes Clusters to be used by the alicloud_cs_kubernetes_clusters resource.
---

# alicloud\_cs\_kubernetes\_clusters

This data source provides a list Container Service Kubernetes Clusters on Alibaba Cloud.

-> **NOTE:** Available in v1.34.0+.

-> **NOTE:** From version 1.177.0+, We supported batch export of clusters' kube config information by `kube_config_file_prefix`.

## Example Usage

```
# Declare the data source
data "alicloud_cs_kubernetes_clusters" "k8s_clusters" {
  name_regex              = "my-first-k8s"
  output_file             = "my-first-k8s-json"
  kube_config_file_prefix = "~/.kube/k8s"
}

output "output" {
  value = "${data.alicloud_cs_kubernetes_clusters.k8s_clusters.clusters}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) Cluster IDs to filter.
* `name_regex` - (Optional) A regex string to filter results by cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enabled_details` - (Optional) Boolean, false by default, only `id` and `name` are exported. Set to true if more details are needed, e.g., `master_disk_category`, `slb_internet_enabled`, `connections`. See full list in attributes.
* `kube_config_file_prefix` - (Optional, Available in 1.177.0+) The path prefix of kube config. You could store kube config in a specified directory by specifying this field, like `~/.kube/k8s`, then it will be named with `~/.kube/k8s-clusterID-kubeconfig`. If you don't specify this field, it will be stored in the current directory and named with `clusterID-kubeconfig`.

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
  * `slb_internet_enabled` - Whether internet load balancer for API Server is created
  * `security_group_id` - The ID of security group where the current cluster worker node is located.
  * `image_id` - The ID of node image.
  * `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
  * `master_instance_types` - The instance type of master node.
  * `worker_instance_types` - The instance type of worker node.
  * `master_disk_category` - The system disk category of master node.
  * `master_disk_size` - The system disk size of master node.
  * `worker_disk_category` - The system disk category of worker node.
  * `worker_disk_size` - The system disk size of worker node.
  * `worker_data_disk_category` - The data disk size of worker node.
  * `worker_data_disk_size` - The data disk category of worker node.
  * `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
  * `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
  * `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
  * `node_cidr_mask` - The network mask used on pods for each node.
  * `log_config` - A list of one element containing information about the associated log store. It contains the following attributes:
    * `type` - Type of collecting logs.
    * `project` - Log Service project name.
  * `resource_group_id` - (Optional, ForceNew, Available in 1.101.0+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.

### Block Nodes

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.
* `role` - (Deprecated from version 1.9.4)

### Block Connections

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.
* `service_domain` - Service Access Domain.

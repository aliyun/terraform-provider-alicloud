---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_serverless_kubernetes_clusters"
sidebar_current: "docs-alicloud-datasource-cs-serverless-kubernetes-clusters"
description: |-
  Provides a list of Container Service Serverless Kubernetes Clusters to be used by the alicloud_cs_serverless_kubernetes_clusters resource.
---

# alicloud\_cs\_serverless\_kubernetes\_clusters

This data source provides a list Container Service Serverless Kubernetes Clusters on Alibaba Cloud.

-> **NOTE:** Available in 1.58.0+

-> **NOTE:** From version 1.177.0+, We supported batch export of clusters' kube config information by `kube_config_file_prefix`.

## Example Usage

```
# Declare the data source
data "alicloud_cs_serverless_kubernetes_clusters" "k8s_clusters" {
  name_regex              = "my-first-k8s"
  output_file             = "my-first-k8s-json"
  kube_config_file_prefix = "~/.kube/serverless"
}

output "output" {
  value = "${data.alicloud_cs_serverless_kubernetes_clusters.k8s_clusters.clusters}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) Cluster IDs to filter.
* `name_regex` - (Optional) A regex string to filter results by cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enabled_details` - (Optional) Boolean, false by default, only `id` and `name` are exported. Set to true if more details are needed, e.g.,  `deletion_protection`, `connections`. See full list in attributes.
* `kube_config_file_prefix` - (Optional, Available in 1.177.0+) The path prefix of kube config. You could store kube config in a specified directory by specifying this field, like `~/.kube/serverless`, then it will be named with `~/.kube/serverless-clusterID-kubeconfig`. If you don't specify this field, it will be stored in the current directory and named with `clusterID-kubeconfig`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Kubernetes clusters' ids.
* `names` - A list of matched Kubernetes clusters' names.
* `clusters` - A list of matched Kubernetes clusters. Each element contains the following attributes:
  * `id` - The ID of the container cluster.
  * `name` - The name of the container cluster.
  * `vswitch_id` - The ID of VSwitch where the current cluster is located.
  * `vpc_id` - The ID of VPC where the current cluster is located.
  * `security_group_id` - The ID of security group where the current cluster  is located.
  * `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
  * `deletion_protection` - Whether the cluster support delete protection.  
  * `connections` - Map of serverless cluster connection information. It contains several attributes to `Block Connections`.
  * `resource_group_id` - (Optional, ForceNew, Available in 1.101.0+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
  
### Block Connections

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.

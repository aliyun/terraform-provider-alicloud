---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_clusters"
sidebar_current: "docs-alicloud-datasource-ehpc-clusters"
description: |-
  Provides a list of Ehpc Clusters to the user.
---

# alicloud\_ehpc\_clusters

This data source provides the Ehpc Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ehpc_clusters" "ids" {
  ids = ["example_id"]
}
output "ehpc_cluster_id_1" {
  value = data.alicloud_ehpc_clusters.ids.clusters.0.id
}

data "alicloud_ehpc_clusters" "nameRegex" {
  name_regex = "^my-Cluster"
}
output "ehpc_cluster_id_2" {
  value = data.alicloud_ehpc_clusters.nameRegex.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Cluster IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values:
  * `uninit`: The cluster is not initialized.
  * `creating`: The cluster is being created.
  * `init`: The cluster is being initialized.
  * `running`: The cluster is running.
  * `exception`: The cluster encounters an exception.
  * `releasing`: The cluster is being released.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Cluster names.
* `clusters` - A list of Ehpc Clusters. Each element contains the following attributes:
  * `account_type` - The server type of the account.
		* `client_version` - The version number of the client used by the cluster.
  * `cluster_id` - The id of E-HPC Cluster.
  * `cluster_name` - The name of E-HPC cluster.
  * `compute_count` - The number of compute nodes in the cluster.
  * `compute_instance_type` - Cluster compute node specifications.
  * `create_time` - The creation time of the resource.
  * `deploy_mode` - Cluster deployment mode. Possible values:
    - Standard: separate deployment of account nodes, scheduling nodes, login nodes, and compute nodes.
    - Advanced:HA mode deployment.
    - Simple: the account node and the scheduling node are deployed on one node, and the login node and the compute node are deployed separately.
    - Tiny: account nodes, scheduling nodes, and login nodes are deployed on one node, and compute nodes are deployed separately.
    - OneBox: account node, scheduling node, login node and compute node are deployed on one node.
  * `description` - The description of E-HPC cluster.
  * `ha_enable` - Whether to turn on high availability. > If high availability is enabled, each control role in the cluster will use two primary and secondary instances.
  * `id` - The ID of the Cluster.
  * `image_id` - The ID of the Image.
  * `image_owner_alias` - The type of the image.
  * `login_count` - The number of cluster login nodes. Only configuration 1 is supported.
  * `login_instance_type` - Cluster login node specifications.
  * `manager_instance_type` - The instance type of manager nodes.
  * `os_tag` - The image tag of the operating system.
  * `remote_directory` - Mount the remote directory of the shared storage.
  * `scc_cluster_id` - The SccCluster ID used by the cluster. If the cluster is not an SCC model, it is empty.
  * `scheduler_type` - Dispatch server type.
  * `security_group_id` - The ID of the security group.
  * `status` - The status of the resource.
  * `volume_id` - The ID of the NAS instance. Currently, you cannot automatically create an Alibaba Cloud NAS instance.
  * `volume_mountpoint` - The mount target of the file system. Mount targets cannot be automatically created for NAS file systems.
  * `volume_protocol` - The type of the protocol that is used by the file system.
  * `volume_type` - The type of the network shared storage. Valid value: NAS.
  * `vpc_id` - The ID of the VPC network.
  * `vswitch_id` - The vswitch id.
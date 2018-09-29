---
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-kubernetes"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster.
---

# alicloud\_cs\_kubernetes

This resource will help you to manager a Kubernetes Cluster. The cluster is same as container service created by web console.

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** If there is no specified `vswitch_ids`, the resource will create a new VPC and VSwitch while creating kubernetes cluster.

-> **NOTE:** Each kubernetes cluster contains 3 master nodes and those number cannot be changed at now.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.

-> **NOTE:** From version 1.9.4, the provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** From version 1.16.0, the provider supports Multiple Availability Zones Kubernetes Cluster. To create a cluster of this kind,
you must specify three items in `vswitch_ids`, `master_instance_types` and `worker_instance_types`.

## Example Usage

Basic Usage

```
data "alicloud_zones" "default" {
  "available_resource_creation"= "VSwitch"
}

resource "alicloud_cs_kubernetes" "main" {
  name_prefix = "my-first-k8s"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  new_nat_gateway = true
  master_instance_types = ["ecs.n4.small"]
  worker_instance_types = ["ecs.n4.small"]
  worker_numbers = [3]
  password = "Test12345"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `availability_zone` - (Force new resource) The Zone where new kubernetes cluster will be located. If it is not be specified, the value will be vswitch's zone.
* `vswitch_id` - (Deprecated from version 1.16.0)(Force new resource) The vswitch where new kubernetes cluster will be located. If it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `vswitch_ids` - (Force new resource) The vswitch where new kubernetes cluster will be located. For SingleAZ Cluster, if it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified. For MultiAZ Cluster, you must create three vswitches firstly, specify them here.
* `new_nat_gateway` - (Force new resource) Whether to create a new nat gateway while creating kubernetes cluster. Default to true.
* `master_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of master node.
* `master_instance_types` - (Required, Force new resource) The instance type of master node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of worker node.
* `worker_instance_types` - (Required, Force new resource) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_number` - The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `password` - (Required, Force new resource) The password of ssh login cluster node. You have to specify one of `password` and `key_name` fields.
* `key_name` - (Required, Force new resource) The keypair of ssh login cluster node, you have to create it first.
* `cluster_network_type` - (Required, Force new resource) The network that cluster uses, use `flannel` or `terway`.
* `pod_cidr` - (Required, Force new resource) The CIDR block for the pod network. It will be allocated automatically when `vswitch_ids` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
Maximum number of hosts allowed in the cluster: 256. Refer to [Plan Kubernetes CIDR blocks under VPC](https://www.alibabacloud.com/help/doc-detail/64530.htm).
* `service_cidr` - (Required, Force new resource) The CIDR block for the service network.  It will be allocated automatically when `vswitch_id` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional, Force new resource) The network mask used on pods for each node, ranging from `24` to `28`.
Larger this number is, less pods can be allocated on each node. Default value is `24`, means you can create 256 pods on echo node.
* `log_config` - (Optional, Force new resource) A list of one element containing information about the associated log store. It contains the following attributes:
  * `type` - Type of collecting logs, only `SLS` are supported currently.
  * `project` - Log Service project name, if none specified, it will help you create one SLS project.
* `enable_ssh` - (Force new resource) Whether to allow to SSH login kubernetes. Default to false.
* `master_disk_category` - (Force new resource) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Force new resource) The system disk size of master node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_disk_category` - (Force new resource) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Force new resource) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_data_disk_category` - (Force new resource) The data disk size of worker node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_data_disk_size` - (Force new resource) The data disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `install_cloud_monitor` - (Force new resource) Whether to install cloud monitor for the kubernetes' node.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `key_name` - The keypair of ssh login cluster node, you have to create it first.
* `worker_number` - (Deprecated from version 1.16.0) The ECS instance node number in the current container cluster.
* `worker_numbers` - The ECS instance node number in the current container cluster.
* `vswitch_id` - (Deprecated from version 1.16.0) The ID of VSwitch where the current cluster is located.
* `vswitch_ids` - The ID of VSwitches where the current cluster is located.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_id` - (Deprecated from version 1.9.2).
* `slb_internet` - The ID of public load balancer where the current cluster master node is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `image_id` - The ID of node image.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `master_instance_type` - (Deprecated from version 1.16.0) The instance type of master node.
* `master_instance_types` - The instance type of master node.
* `worker_instance_type` - (Deprecated from version 1.16.0)The instance type of worker node.
* `worker_instance_types` - The instance type of worker node.
* `master_disk_category` - The system disk category of master node.
* `master_disk_size` - The system disk size of master node.
* `worker_disk_category` - The system disk category of worker node.
* `worker_disk_size` - The system disk size of worker node.
* `worker_data_disk_category` - The data disk size of worker node.
* `worker_data_disk_size` - The data disk category of worker node.
* `nodes` - (Deprecated from version 1.9.4) It has been deprecated from provider version 1.9.4. New field `master_nodes` and `worker_nodes` replace it.
* `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
* `node_cidr_mask` - The network mask used on pods for each node.
* `log_config` - A list of one element containing information about the associated log store. It contains the following attributes:
  * `type` - Type of collecting logs.
  * `project` - Log Service project name.

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

## Import

Kubernetes cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_kubernetes.main ce4273f9156874b46bb
```
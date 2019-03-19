---
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_managed_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-managed-kubernetes"
description: |-
  Provides a Alicloud resource to manage container managed kubernetes cluster.
---

# alicloud\_cs\_managed\_kubernetes

This resource will help you to manager a Managed Kubernetes Cluster. The cluster is same as container service created by web console.

-> **NOTE:** Managed Kubernetes cluster only supports single availability zone currently. Arguments `vswitch_ids`, `worker_numbers`, `worker_instance_types` only accept one item.

-> **NOTE:** Managed Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** If there is no specified `vswitch_ids`, the resource will create a new VPC and VSwitch while creating managed kubernetes cluster.

-> **NOTE:** Creating managed kubernetes cluster need to install several packages and it will cost about 10 minutes. Please be patient.

-> **NOTE:** The provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** If you want to manage managed Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

## Example Usage

Basic Usage

```
variable "name" {
	default = "my-first-k8s"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_numbers = [2]
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `availability_zone` - (Force new resource) The Zone where new kubernetes cluster will be located. If it is not be specified, the value will be vswitch's zone.
* `vswitch_ids` - (Force new resource) The vswitch where new kubernetes cluster will be located. Specify one vswitch's id, if it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `new_nat_gateway` - (Force new resource) Whether to create a new nat gateway while creating kubernetes cluster. Default to true.
* `password` - (Required, Force new resource) The password of ssh login cluster node. You have to specify one of `password` and `key_name` fields.
* `key_name` - (Required, Force new resource) The keypair of ssh login cluster node, you have to create it first.
* `pod_cidr` - (Required, Force new resource) The CIDR block for the pod network. It will be allocated automatically when `vswitch_ids` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
Maximum number of hosts allowed in the cluster: 256. Refer to [Plan Kubernetes CIDR blocks under VPC](https://www.alibabacloud.com/help/doc-detail/64530.htm).
* `service_cidr` - (Required, Force new resource) The CIDR block for the service network.  It will be allocated automatically when `vswitch_id` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `slb_internet_enabled` - (Force new resource) Whether to create internet load balancer for API Server. Default to false.
* `install_cloud_monitor` - (Force new resource) Whether to install cloud monitor for the kubernetes' node.
* `worker_disk_size` - (Force new resource) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_disk_category` - (Force new resource) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_data_disk_size` - (Force new resource) The data disk size of worker node. Its valid value range [20~32768] in GB. When `worker_data_disk_category` is presented, it defaults to 40.
* `worker_data_disk_category` - (Force new resource) The data disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`, if not set, data disk will not be created.
* `worker_numbers` - The worker node number of the kubernetes cluster. Default to [3]. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_instance_types` - (Required, Force new resource) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
You can get the available kubetnetes master node instance types by [datasource instance_types](https://www.terraform.io/docs/providers/alicloud/d/instance_types.html#kubernetes_node_role)
* `worker_instance_charge_type` - (Optional, Force new resource) Worker payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `worker_period_unit` - (Optional) Worker payment period unit. `Month` or `Week`, defaults to `Month`.
* `worker_period` - (Optional) Worker payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `worker_auto_renew` - (Optional) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Optional) Worker payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `cluster_network_type` - (Optional, Force new resource) The network that cluster uses, use `flannel` or `terway`.
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
* `worker_numbers` - The ECS instance node number in the current container cluster.
* `vswitch_ids` - The ID of VSwitches where the current cluster is located.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `image_id` - The ID of node image.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_instance_types` - The instance type of worker node.
* `worker_disk_size` - The system disk size of worker node.
* `worker_disk_category` - The system disk category of worker node.
* `worker_data_disk_size` - The data disk category of worker node.
* `worker_data_disk_category` - The data disk size of worker node.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.

### Block Nodes

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.

## Import

Managed Kubernetes cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_managed_kubernetes.main ce4273f9156874b46bb
```
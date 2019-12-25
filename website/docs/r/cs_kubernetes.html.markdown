---
subcategory: "Container Service (CS)"
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

-> **NOTE:** From version 1.20.0, the provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** You need to activate several other products and confirm Authorization Policy used by Container Service before using this resource.
Please refer to the `Authorization management` and `Cluster management` sections in the [Document Center](https://www.alibabacloud.com/help/doc-detail/86488.htm).

-> **NOTE:** From version 1.50.0, when `force_update` is set to `false`, updates to the following arguments will be ignored: `vswitch_ids`, `master_instance_types`, `worker_instance_types`, `worker_numbers`, `password`, `key_name`, `user_ca`, `pod_cidr`, `service_cidr`, `cluster_network_type`, `node_cidr_mask`, `log_config`, `enable_ssh`, `master_disk_size`, `master_disk_category`, `worker_disk_size`, `worker_disk_category`, `worker_data_disk_category`, `master_instance_charge_type`, `worker_instance_charge_type`, `install_cloud_monitor`, `is_outdated`.


## Example Usage

Single AZ Kubernetes Cluster

```
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cs_kubernetes" "main" {
  name_prefix           = "my-first-k8s"
  availability_zone     = data.alicloud_zones.default.zones.0.id
  new_nat_gateway       = true
  master_instance_types = ["ecs.n4.small"]
  worker_instance_types = ["ecs.n4.small"]
  worker_numbers        = [3]
  password              = "Yourpassword1234"
  pod_cidr              = "192.168.1.0/16"
  service_cidr          = "192.168.2.0/24"
  enable_ssh            = true
  install_cloud_monitor = true
}
```

Three AZ Kubernetes Cluster

```
variable "name" {
  default = "my-first-3az-k8s"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "instance_types_1_master" {
  availability_zone    = data.alicloud_zones.main.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Master"
}
data "alicloud_instance_types" "instance_types_2_master" {
  availability_zone    = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 1) % length(data.alicloud_zones.main.zones)], "id")}"
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Master"
}
data "alicloud_instance_types" "instance_types_3_master" {
  availability_zone    = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 2) % length(data.alicloud_zones.main.zones)], "id")}"
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Master"
}

data "alicloud_instance_types" "instance_types_1_worker" {
  availability_zone    = data.alicloud_zones.main.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
data "alicloud_instance_types" "instance_types_2_worker" {
  availability_zone    = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 1) % length(data.alicloud_zones.main.zones)], "id")}"
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
data "alicloud_instance_types" "instance_types_3_worker" {
  availability_zone    = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 2) % length(data.alicloud_zones.main.zones)], "id")}"
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
resource "alicloud_vpc" "foo" {
  name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "vsw1" {
  name              = var.name
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.main.zones.0.id
}

resource "alicloud_vswitch" "vsw2" {
  name              = var.name
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "10.1.2.0/24"
  availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 1) % length(data.alicloud_zones.main.zones)], "id")}"
}

resource "alicloud_vswitch" "vsw3" {
  name              = var.name
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "10.1.3.0/24"
  availability_zone = "${lookup(data.alicloud_zones.main.zones[(length(data.alicloud_zones.main.zones) - 2) % length(data.alicloud_zones.main.zones)], "id")}"
}

resource "alicloud_nat_gateway" "nat_gateway" {
  name          = var.name
  vpc_id        = alicloud_vpc.foo.id
  specification = "Small"
}

resource "alicloud_snat_entry" "snat_entry_1" {
  snat_table_id     = alicloud_nat_gateway.nat_gateway.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vsw1.id
  snat_ip           = alicloud_eip.eip.ip_address
}

resource "alicloud_snat_entry" "snat_entry_2" {
  snat_table_id     = alicloud_nat_gateway.nat_gateway.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vsw2.id
  snat_ip           = alicloud_eip.eip.ip_address
}

resource "alicloud_snat_entry" "snat_entry_3" {
  snat_table_id     = alicloud_nat_gateway.nat_gateway.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vsw3.id
  snat_ip           = alicloud_eip.eip.ip_address
}

resource "alicloud_eip" "eip" {
  name      = var.name
  bandwidth = "100"
}

resource "alicloud_eip_association" "eip_asso" {
  allocation_id = alicloud_eip.eip.id
  instance_id   = alicloud_nat_gateway.nat_gateway.id
}

resource "alicloud_cs_kubernetes" "k8s" {
  name                      = var.name
  vswitch_ids               = [alicloud_vswitch.vsw1.id, alicloud_vswitch.vsw2.id, alicloud_vswitch.vsw3.id]
  new_nat_gateway           = true
  master_instance_types     = [data.alicloud_instance_types.instance_types_1_master.instance_types.0.id, data.alicloud_instance_types.instance_types_2_master.instance_types.0.id, data.alicloud_instance_types.instance_types_3_master.instance_types.0.id]
  worker_instance_types     = [data.alicloud_instance_types.instance_types_1_worker.instance_types.0.id, data.alicloud_instance_types.instance_types_2_worker.instance_types.0.id, data.alicloud_instance_types.instance_types_3_worker.instance_types.0.id]
  worker_numbers            = [1, 2, 3]
  master_disk_category      = "cloud_ssd"
  worker_disk_size          = 50
  worker_data_disk_category = "cloud_ssd"
  worker_data_disk_size     = 50
  password                  = "Yourpassword1234"
  pod_cidr                  = "192.168.1.0/16"
  service_cidr              = "192.168.2.0/24"
  enable_ssh                = true
  slb_internet_enabled      = true
  node_cidr_mask            = 25
  install_cloud_monitor     = true
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Optional) The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `force_update` - (Optional, Available in 1.50.0+) Whether to force the update of kubernetes cluster arguments. Default to false.
* `availability_zone` - (Optional, ForceNew) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.
* `vswitch_id` - (Deprecated from version 1.16.0)(Force new resource) The vswitch where new kubernetes cluster will be located. If it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `vswitch_ids` - (Required, ForceNew) The vswitch where new kubernetes cluster will be located. Specify one or more vswitch's id. It must be in the zone which `availability_zone` specified.
* `new_nat_gateway` - (Optional, ForceNew) Whether to create a new nat gateway while creating kubernetes cluster. Default to true.
* `master_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of master node.
* `master_instance_types` - (Required, ForceNew) The instance type of master node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
You can get the available kubetnetes master node instance types by [datasource instance_types](https://www.terraform.io/docs/providers/alicloud/d/instance_types.html#kubernetes_node_role)
* `worker_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of worker node.
* `worker_instance_types` - (Required, ForceNew) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
You can get the available kubetnetes master node instance types by [datasource instance_types](https://www.terraform.io/docs/providers/alicloud/d/instance_types.html#kubernetes_node_role)
* `worker_number` - (Required) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `password` - (Optional, ForceNew, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Optional, ForceNew) The keypair of ssh login cluster node, you have to create it first.
* `kms_encrypted_password` - (Optional, ForceNew, Available in 1.57.1+) An KMS encrypts password used to a cs kubernetes. It is conflicted with `password` and `key_name`.
* `kms_encryption_context` - (Optional, ForceNew, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `user_ca` - (Optional, ForceNew) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `cluster_network_type` - (Optional, ForceNew) The network that cluster uses, use `flannel` or `terway`.
* `pod_cidr` - (Optional, ForceNew) The CIDR block for the pod network. It will be allocated automatically when `vswitch_ids` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
Maximum number of hosts allowed in the cluster: 256. Refer to [Plan Kubernetes CIDR blocks under VPC](https://www.alibabacloud.com/help/doc-detail/64530.htm).
* `service_cidr` - (Optional, ForceNew) The CIDR block for the service network. 
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `master_instance_charge_type` - (Optional, ForceNew) Master payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `master_period_unit` - (Optional) Master payment period unit. `Month` or `Week`, defaults to `Month`.
* `master_period` - (Optional) Master payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `master_auto_renew` - (Optional) Enable master payment auto-renew, defaults to false.
* `master_auto_renew_period` - (Optional) Master payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `worker_instance_charge_type` - (Optional, Force new resource) Worker payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `worker_period_unit` - (Optional) Worker payment period unit. `Month` or `Week`, defaults to `Month`.
* `worker_period` - (Optional) Worker payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `worker_auto_renew` - (Optional) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Optional) Worker payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `node_cidr_mask` - (Optional, Force new resource) The network mask used on pods for each node, ranging from `24` to `28`.
Larger this number is, less pods can be allocated on each node. Default value is `24`, means you can allocate 256 pods on each node.
* `log_config` - (Optional, ForceNew) A list of one element containing information about the associated log store. It contains the following attributes:
  * `type` - Type of collecting logs, only `SLS` are supported currently.
  * `project` - Log Service project name, cluster logs will output to this project.
* `enable_ssh` - (Optional, ForceNew) Whether to allow to SSH login kubernetes. Default to false.
* `slb_internet_enabled` - (Optional, ForceNew) Whether to create internet load balancer for API Server. Default to true.
* `master_disk_category` - (Optional, ForceNew) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Optional, ForceNew) The system disk size of master node. Its valid value range [20~500] in GB. Default to 20.
* `worker_disk_category` - (Optional, ForceNew) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional, ForceNew) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_data_disk_size` - (Optional, ForceNew) The data disk size of worker node. Its valid value range [20~32768] in GB. When `worker_data_disk_category` is presented, it defaults to 40.
* `worker_data_disk_category` - (Optional, ForceNew) The data disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`, if not set, data disk will not be created.
* `install_cloud_monitor` - (Optional, ForceNew) Whether to install cloud monitor for the kubernetes' node.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

### Timeouts

-> **NOTE:** Available in 1.58.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.

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

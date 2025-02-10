---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_managed_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-managed-kubernetes"
description: |-
  Provides a Alicloud resource to manage container managed kubernetes cluster.
---

# alicloud_cs_managed_kubernetes

This resource will help you to manage a ManagedKubernetes Cluster in Alibaba Cloud Kubernetes Service. 

-> **NOTE:** Available since v1.26.0.

-> **NOTE:** It is recommended to create a cluster with zero worker nodes, and then use a node pool to manage the cluster nodes. 

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.

-> **NOTE:** From version 1.9.4, the provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** From version 1.20.0, the provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** You need to activate several other products and confirm Authorization Policy used by Container Service before using this resource.
Please refer to the `Authorization management` and `Cluster management` sections in the [Document Center](https://www.alibabacloud.com/help/doc-detail/86488.htm).

-> **NOTE:** From version 1.72.0, Some parameters have been removed from resource,You can check them below and re-import the cluster if necessary.

-> **NOTE:** From version 1.120.0, Support for cluster migration from Standard cluster to professional.

-> **NOTE:** From version 1.177.0, `runtime`,`enable_ssh`,`rds_instances`,`exclude_autoscaler_nodes`,`worker_number`,`worker_instance_types`,`password`,`key_name`,`kms_encrypted_password`,`kms_encryption_context`,`worker_instance_charge_type`,`worker_period`,`worker_period_unit`,`worker_auto_renew`,`worker_auto_renew_period`,`worker_disk_category`,`worker_disk_size`,`worker_data_disks`,`node_name_mode`,`node_port_range`,`os_type`,`platform`,`image_id`,`cpu_policy`,`user_data`,`taints`,`worker_disk_performance_level`,`worker_disk_snapshot_policy_id`,`install_cloud_monitor` are deprecated.
We Suggest you using resource **`alicloud_cs_kubernetes_node_pool`** to manage your cluster worker nodes.

-> **NOTE:** From version 1.212.0, `runtime`,`enable_ssh`,`rds_instances`,`exclude_autoscaler_nodes`,`worker_number`,`worker_instance_types`,`password`,`key_name`,`kms_encrypted_password`,`kms_encryption_context`,`worker_instance_charge_type`,`worker_period`,`worker_period_unit`,`worker_auto_renew`,`worker_auto_renew_period`,`worker_disk_category`,`worker_disk_size`,`worker_data_disks`,`node_name_mode`,`node_port_range`,`os_type`,`platform`,`image_id`,`cpu_policy`,`user_data`,`taints`,`worker_disk_performance_level`,`worker_disk_snapshot_policy_id`,`install_cloud_monitor`,`kube_config`,`availability_zone` are removed.
Please use resource **`alicloud_cs_kubernetes_node_pool`** to manage your cluster worker nodes.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_managed_kubernetes&exampleId=27aadcc9-05f0-2e9f-38f6-9c130b175087b7374a3b&activeTab=example&spm=docs.r.cs_managed_kubernetes.0.27aadcc905&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

# leave it to empty would create a new one
variable "vpc_id" {
  description = "Existing vpc id used to create several vswitches and other resources."
  default     = ""
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

# leave it to empty then terraform will create several vswitches
variable "vswitch_ids" {
  description = "List of existing vswitch id."
  type        = list(string)
  default     = []
}


variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16"]
}

# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}

# options: ipvs|iptables
variable "proxy_mode" {
  description = "Proxy mode is option of kube-proxy."
  default     = "ipvs"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "192.168.0.0/16"
}

variable "terway_vswitch_ids" {
  description = "List of existing vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "terway_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'terway_vswitch_cidrs' is not specified."
  type        = list(string)
  default     = ["10.4.0.0/16", "10.5.0.0/16"]
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

# If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count      = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "terway_vswitches" {
  count      = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.terway_vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name         = var.name
  cluster_spec = "ack.pro.small"
  # version can not be defined in variables.tf.
  # version            = "1.26.3-aliyun.1"
  vswitch_ids     = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)) : length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids = length(var.terway_vswitch_ids) > 0 ? split(",", join(",", var.terway_vswitch_ids)) : length(var.terway_vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  new_nat_gateway = true
  node_cidr_mask  = var.node_cidr_mask
  proxy_mode      = var.proxy_mode
  service_cidr    = var.service_cidr

  addons {
    name = "terway-eniip"
  }
  addons {
    name = "csi-plugin"
  }
  addons {
    name = "csi-provisioner"
  }
  addons {
    name = "logtail-ds"
    config = jsonencode({
      IngressDashboardEnabled = "true"
    })
  }
  addons {
    name = "nginx-ingress-controller"
    config = jsonencode({
      IngressSlbNetworkType = "internet"
    })
    # to disable install nginx-ingress-controller automatically
    # disabled = true
  }
  addons {
    name = "arms-prometheus"
  }
  addons {
    name = "ack-node-problem-detector"
    config = jsonencode({
      # sls_project_name = "your-sls-project"
    })
  }
}
```

## Argument Reference

The following arguments are supported:

*Global params*

* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `worker_vswitch_ids` - (Optional, Deprecated since v1.241.0) The vswitches used by control plane. Modification after creation will not take effect. Please use `vswitch_ids` to managed control plane vswtiches, which supports modifying control plane vswtiches.
* `vswitch_ids` - (Optional, Available since v1.241.0) The vSwitches of the control plane.
-> **NOTE:** Please take of note before updating the `vswitch_ids`:
  * This parameter overwrites the existing configuration. You must specify all vSwitches of the control plane. 
  * The control plane restarts during the change process. Exercise caution when you perform this operation. 
  * Ensure that all security groups of the cluster, including the security groups of the control plane, all node pools, and container network, are allowed to access the CIDR blocks of the new vSwitches. This ensures that the nodes and containers can connect to the API server. 
  * If the new vSwitches of the control plane are configured with an ACL, ensure that the ACL allows communication between the new vSwitches and CIDR blocks such as those of the cluster nodes and the container network.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `timezone` - (Optional, ForceNew, Available since v1.103.2) When you create a cluster, set the time zones for the Master and Worker nodes. You can only change the managed node time zone if you create a cluster. Once the cluster is created, you can only change the time zone of the Worker node.
* `resource_group_id` - (Optional, Available since v1.101.0) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `version` - (Optional, Available since 1.70.1) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK. Do not specify if cluster auto upgrade is enabled, see [cluster_auto_upgrade](#operation_policy-cluster_auto_upgrade) for more information. 
* `security_group_id` - (Optional, ForceNew, Available since v1.91.0) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, ForceNew, Available since v1.91.0) Enable to create advanced security group. default: false. Only works for **Create** Operation. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm). 
* `proxy_mode` - (Optional, ForceNew) Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `cluster_domain` - (Optional, ForceNew, Available since v1.103.2) Cluster local domain name, Default to `cluster.local`. A domain name consists of one or more sections separated by a decimal point (.), each of which is up to 63 characters long, and can be lowercase, numerals, and underscores (-), and must be lowercase or numerals at the beginning and end.
* `custom_san` - (Optional, Available since v1.103.2) Customize the certificate SAN, multiple IP or domain names are separated by English commas (,).
-> **NOTE:** Make sure you have specified all certificate SANs before updating. Updating this field will lead APIServer to restart.
* `user_ca` - (Optional) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `deletion_protection` - (Optional, Available since v1.103.2)  Whether to enable cluster deletion protection.
* `enable_rrsa` - (Optional, Available since v1.171.0) Whether to enable cluster to support RRSA for kubernetes version 1.22.3+. Default to `false`. Once the RRSA function is turned on, it is not allowed to turn off. If your cluster has enabled this function, please manually modify your tf file and add the rrsa configuration to the file, learn more [RAM Roles for Service Accounts](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/use-rrsa-to-enforce-access-control).
* `service_account_issuer` - (Optional, ForceNew, Available since v1.92.0) The issuer of the Service Account token for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm), corresponds to the `iss` field in the token payload. Set this to `"https://kubernetes.default.svc"` to enable the Token Volume Projection feature (requires specifying `api_audiences` as well). From cluster version 1.22, Service Account Token Volume Projection will be enabled by default.
* `api_audiences` - (Optional, ForceNew, Available since v1.92.0) A list of API audiences for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm). Set this to `["https://kubernetes.default.svc"]` if you want to enable the Token Volume Projection feature (requires specifying `service_account_issuer` as well. From cluster version 1.22, Service Account Token Volume Projection will be enabled by default.
* `tags` - (Optional, Available since v1.97.0) Default nil, A map of tags assigned to the kubernetes cluster and work nodes. See [`tags`](#tags) below.
* `cluster_spec` - (Optional, Available since v1.101.0) The cluster specifications of kubernetes cluster,which can be empty. Valid values:
  * ack.standard : Standard managed clusters.
  * ack.pro.small : Professional managed clusters.
* `encryption_provider_key` - (Optional, ForceNew, Available since v1.103.2) The disk encryption key.
* `maintenance_window` - (Optional, Available since v1.109.1) The cluster maintenance window，effective only in the professional managed cluster. Managed node pool will use it. See [`maintenance_window`](#maintenance_window) below.
* `operation_policy` - (Optional, Available since v1.232.0) The cluster automatic operation policy. See [`operation_policy`](#operation_policy) below.
* `load_balancer_spec` - (Optional, Deprecated since v1.232.0) The cluster api server load balancer instance specification. For more information on how to select a LB instance specification, see [SLB instance overview](https://help.aliyun.com/document_detail/85931.html). Only works for **Create** Operation. The spec will not take effect because the charge of the load balancer has been changed to PayByCLCU.
* `control_plane_log_ttl` - (Optional, Available since v1.141.0) Control plane log retention duration (unit: day). Default `30`. If control plane logs are to be collected, `control_plane_log_ttl` and `control_plane_log_components` must be specified.
* `control_plane_log_components` - (Optional, Available since v1.141.0) List of target components for which logs need to be collected. Supports `apiserver`, `kcm`, `scheduler`, `ccm` and `controlplane-events`.
* `control_plane_log_project` - (Optional, Available since v1.141.0) Control plane log project. If this field is not set, a log service project named k8s-log-{ClusterID} will be automatically created.
* `retain_resources` - (Optional, Available since v1.141.0) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.
* `delete_options` - (Optional, Available since v1.223.2) Delete options, only work for deleting resource. Make sure you have run `terraform apply` to make the configuration applied. See [`delete_options`](#delete_options) below.
* `addons` - (Optional, Available since v1.88.0) The addon you want to install in cluster. See [`addons`](#addons) below. Only works for **Create** Operation, use [resource cs_kubernetes_addon](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_addon) to manage addons if cluster is created.

*Network params*

* `pod_cidr` - (Optional, ForceNew) - [Flannel Specific] The CIDR block for the pod network when using Flannel.
* `pod_vswitch_ids` - (Optional) - [Terway Specific] The vswitches for the pod network when using Terway. It is recommended that `pod_vswitch_ids` is not belong to `vswitch_ids` but must be in same availability zones. Only works for **Create** Operation. 
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice. Only works for **Create** Operation. 
* `service_cidr` - (Optional, ForceNew) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional, ForceNew) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.
* `ip_stack` - (Optional, ForceNew, Available since v1.243.0) The IP address family that the cluster network uses. Valid values:
  * `ipv4`: IPv4 stack.
  * `dual`: IPv4/IPv6 dual stack. IPv4 addresses are used for communication between worker nodes and the control plane. The VPC used by the cluster must support IPv4 and IPv6 dual-stack. This feature is only supported for Kubernetes version 1.22 and later, and you must select `Terway` as CNI network plugin. If you use the shared ENI mode of `Terway`, the ECS instance type must support IPv6 addresses and the number of IPv4 addresses supported by the ECS instance type must be the same as the number of IPv6 addresses, for more information about ECS instance types, see [Overview of instance families](https://www.alibabacloud.com/help/zh/ecs/user-guide/overview-of-instance-families#concept-sx4-lxv-tdb). Dual stack is not supported if you want to use [Elastic Remote Direct Memory Access (eRDMA)](https://www.alibabacloud.com/help/zh/ack/ack-managed-and-ack-dedicated/user-guide/use-erdma-in-ack-clusters) in the cluster.

-> **NOTE:** If you want to use `Terway` as CNI network plugin, You need to specify the `pod_vswitch_ids` field and addons with `terway-eniip`.
If you want to use `Flannel` as CNI network plugin, You need to specify the `pod_cidr` field and addons with `flannel`.

*Computed params*

* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

*Removed params*

* `worker_instance_type` - (Removed from version 1.16.0) The instance type of worker node.
* `force_update` - (Removed) Whether to force the update of kubernetes cluster arguments. Default to false.
* `log_config` - (Removed) A list of one element containing information about the associated log store. See [`log_config`](#log_config) below.
* `cluster_network_type` - (Removed) The network that cluster uses, use `flannel` or `terway`.
* `worker_data_disk_category` - (Removed) The data disk category of worker, use `worker_data_disks` to instead it.
* `worker_data_disk_size` - (Removed) The data disk size of worker, use `worker_data_disks` to instead it.
* `worker_numbers` - (Removed) The number of workers, use `worker_number` to instead it.
* `runtime` - (Removed since v1.212.0) The runtime of containers. If you select another container runtime, see [Comparison of Docker, containerd, and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm). See [`runtime`](#runtime) below.
* `enable_ssh` - (Removed since v1.212.0) Enable login to the node through SSH. Default to `false`.
* `rds_instances` - (Removed since v1.212.0) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `install_cloud_monitor` - (Removed since v1.212.0) Install cloud monitor agent on ECS. Default is `true` in previous version. From provider version 1.208.0, the default value is `false`.
* `exclude_autoscaler_nodes` - (Removed since v1.212.0) Exclude autoscaler nodes from `worker_nodes`. Default to `false`.
* `kube_config` - (Removed since v1.212.0) The path of kube config, like `~/.kube/config`. You can set some file paths to save kube_config information, but this way is cumbersome. Since version 1.105.0, we've written it to tf state file. About its use，see export attribute certificate_authority. From version 1.187.0, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kube_config.
* `availability_zone` - (Removed since v1.212.0) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.
* `worker_number` - (Removed since v1.212.0) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us. From version 1.109.1, It is not necessary in the professional managed cluster, but it is necessary in other types of clusters.
* `worker_instance_types` - (Removed since v1.212.0) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster. From version 1.109.1, It is not necessary in the professional managed cluster, but it is necessary in other types of clusters.
* `password` - (Removed since v1.212.0) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. From ersion 1.109.1, It is not necessary in the professional managed cluster.
* `key_name` - (Removed since v1.212.0) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. From ersion 1.109.1, It is not necessary in the professional managed cluster.
* `kms_encrypted_password` - (Removed since v1.212.0) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Removed since v1.212.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `worker_instance_charge_type` - (Removed since v1.212.0) Worker payment type, its valid value is either or `PostPaid` or `PrePaid`. Defaults to `PostPaid`. If value is `PrePaid`, the files `worker_period`, `worker_period_unit`, `worker_auto_renew` and `worker_auto_renew_period` are required, default is `PostPaid`.
* `worker_period` - (Removed since v1.212.0) Worker payment period. The unit is `Month`. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `worker_period_unit` - (Removed since v1.212.0) Worker payment period unit, the valid value is `Month`.
* `worker_auto_renew` - (Removed since v1.212.0) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Removed since v1.212.0) Worker payment auto-renew period, it can be one of {1, 2, 3, 6, 12}.
* `worker_disk_category` - (Removed since v1.212.0) The system disk category of worker node. Its valid value are `cloud`, `cloud_ssd`, `cloud_essd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Removed since v1.212.0) The system disk size of worker node. Its valid value range [40~500] in GB.
* `worker_data_disks` - (Removed since v1.212.0) The data disk configurations of worker nodes, such as the disk type and disk size.  See [`worker_data_disks`](#worker_data_disks) below.
* `node_name_mode` - (Removed since v1.212.0) Each node name consists of a prefix, an IP substring, and a suffix, the input format is `customized,<prefix>,IPSubStringLen,<suffix>`. For example "customized,aliyun.com-,5,-test", if the node IP address is 192.168.59.176, the prefix is aliyun.com-, IP substring length is 5, and the suffix is -test, the node name will be aliyun.com-59176-test.
* `node_port_range`- (Removed since v1.212.0) The service port range of nodes, valid values: `30000` to `65535`. Default to `30000-32767`.
* `os_type` - (Removed since v1.212.0) The operating system of the nodes that run pods, its valid value is either `Linux` or `Windows`. Default to `Linux`.
* `platform` - (Removed since v1.212.0) The architecture of the nodes that run pods, its valid value is either `CentOS` or `AliyunLinux`. Default to `CentOS`.
* `image_id` - (Removed since v1.212.0) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `cpu_policy` - (Removed since v1.212.0) Kubelet cpu policy. For Kubernetes 1.12.6 and later, its valid value is either `static` or `none`. Default to `none`.
* `user_data` - (Removed since v1.212.0) Custom data that can execute on nodes. For more information, see [Prepare user data](https://www.alibabacloud.com/help/doc-detail/49121.htm).
* `taints` - (Removed since v1.212.0) Taints ensure pods are not scheduled onto inappropriate nodes. One or more taints are applied to a node; this marks that the node should not accept any pods that do not tolerate the taints. For more information, see [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/). See [`taints`](#taints) below.
* `worker_disk_performance_level` - (Removed since v1.212.0) Worker node system disk performance level, when `worker_disk_category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `worker_disk_snapshot_policy_id` - (Removed since v1.212.0) Worker node system disk auto snapshot policy.

### `maintenance_window`

The following arguments are supported in the `maintenance_window` configuration block:

* `enable` - (Optional) Whether to open the maintenance window. The following parameters take effect only `enable = true`.
* `maintenance_time` - (Optional) Initial maintenance time, RFC3339 format. For example: "2024-10-15T12:31:00.000+08:00".
* `duration` - (Optional) The maintenance time, values range from 1 to 24,unit is hour. For example: "3h".
* `weekly_period` - (Optional) Maintenance cycle, you can set the values from Monday to Sunday, separated by commas when the values are multiple. The default is Thursday.

for example:
```
  maintenance_window {
    enable            = true
    maintenance_time  = "2024-10-15T12:31:00.000+08:00"
    duration          = "3h"
    weekly_period     = "Monday,Friday"
  }
```

### `operation_policy`
* `cluster_auto_upgrade` - (Optional) Automatic cluster upgrade policy. See [`cluster_auto_upgrade`](#operation_policy-cluster_auto_upgrade) below.

### `operation_policy-cluster_auto_upgrade`
Automatic cluster upgrade policy. If `enabled` is set to `true`, ACK will automatically upgrade cluster depends on the `channel` value. The `version` field may show diffs if set in config, please remove the field or ignore it.  

* `enabled` - (Optional) Whether to enable automatic cluster upgrade.
* `channel` - (Optional) The automatic cluster upgrade channel. Valid values: `patch`, `stable`, `rapid`.

for example:
```
  operation_policy {
    cluster_auto_upgrade {
      enabled = true
      channel = "stable"
    }
  }
```

### `addons`

The following arguments are supported in the `addons` configuration block:

* `name` - (Optional) This parameter specifies the name of the component.
* `config` - (Optional) If this parameter is left empty, no configurations are required. For more config information, see [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata).
* `version` - (Optional) It specifies the version of the component.
* `disabled` - (Optional) It specifies whether to disable automatic installation. 

It is a new field since 1.75.0. You can specific network plugin, log component,ingress component and so on.

You can get more information about addons on ACK web console. When you create a ACK cluster. You can get openapi-spec before creating the cluster on submission page.

`logtail-ds` - You can specify `IngressDashboardEnabled` and `sls_project_name` in config. If you switch on `IngressDashboardEnabled` and `sls_project_name`,then logtail-ds would use `sls_project_name` as default log store.

`nginx-ingress-controller` - You can specific `IngressSlbNetworkType` in config. Options: internet|intranet.

The `main.tf`:

```
resource "alicloud_cs_managed_kubernetes" "k8s" {
  # ... other configuration ...

  dynamic "addons" {
    for_each = var.cluster_addons
    content {
      name     = lookup(addons.value, "name", var.cluster_addons)
      config   = lookup(addons.value, "config", var.cluster_addons)
      version  = lookup(addons.value, "version", var.cluster_addons)
      disabled = lookup(addons.value, "disabled", var.cluster_addons)
    }
  }
}
```

The `varibales.tf`:

```
# Network-flannel is required, Conflicts With Network-terway
variable "cluster_addons" {
  description = "Addon components in kubernetes cluster"

  type = list(object({
    name      = string
    config    = string
  }))

  default = [
    {
      "name"     = "flannel",
      "config"   = "",
    }
  ]
}

# Network-terway is required, Conflicts With Network-flannel
variable "cluster_addons" {
  type = list(object({
    name      = string
    config    = string
  }))

  default = [
    {
      "name"     = "terway-eniip",
      "config"   = "",
    }
  ]
}

# Storage-csi is required, Conflicts With Storage-flexvolume
variable "cluster_addons" {
  type = list(object({
    name      = string
    config    = string
  }))

  default = [
    {
      "name"     = "csi-plugin",
      "config"   = "",
    },
    {
      "name"     = "csi-provisioner",
      "config"   = "",
    }
  ]
}

# Storage-flexvolume is required, Conflicts With Storage-csi
variable "cluster_addons" {
  type = list(object({
    name      = string
    config    = string
  }))
  default = [
    {
      "name"     = "flexvolume",
      "config"   = "",
    }
  ]
}

# Log, Optional
variable "cluster_addons" {
  type = list(object({
    name      = string
    config    = string
  }))
  default = [
    {
      "name"     = "logtail-ds",
      "config"   = "{\"IngressDashboardEnabled\":\"true\",\"sls_project_name\":\"your-sls-project-name\"}",
    }
  ]
}

# Ingress,Optional
variable "cluster_addons" {
  type = list(object({
    name      = string
    config    = string
  }))

  default = [
    {
      "name"     = "nginx-ingress-controller",
      "config"   = "{\"IngressSlbNetworkType\":\"internet\"}",
    }
  ]
}

# Ingress-Disable, Optional
variable "cluster_addons" {
  type = list(object({
      name      = string
      config    = string
      disabled  = bool
  }))

  default = [
    {
      "name"     = "nginx-ingress-controller",
      "config"   = "",
      "disabled": true,
    }
  ]

# Prometheus, Optional.
variable "cluster_addons" {
  type = list(object({
      name      = string
      config    = string
  }))

  default = [
    {
      "name"     = "arms-prometheus",
      "config"   = "",
    }
  ]
}

# Event Center, Optional.
variable "cluster_addons" {
  type = list(object({
      name      = string
      config    = string
  }))
  default = [
    {
      "name"     = "ack-node-problem-detector",
      "config"   = "{\"sls_project_name\":\"\"}",
    }
  ]
}
# ACK default alert, Optional.
variable "cluster_addons" {
  type = list(object({
      name      = string
      config    = string
  }))
  default = [
    {
      "name"     = "alicloud-monitor-controller",
      "config"   = "{\"group_contact_ids\":\"[159]\"}",
    }
  ]
}
```

### `worker_data_disks`

The following arguments are supported in the `worker_data_disks` configuration block:

* `category` - (Optional) The type of the data disks. Valid values: `cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`. Default to `cloud_efficiency`.
* `size` - (Optional) The size of a data disk, at least 40. Unit: GiB.
* `encrypted` - (Optional) Specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
* `performance_level` - (Optional, Available since v1.120.0) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `auto_snapshot_policy_id` - (Optional, Available since v1.120.0) Worker node data disk auto snapshot policy.
* `kms_key_id` - (Optional) The ID of the Key Management Service (KMS) key to use for data disk N.
* `device` - (Optional) The mount point of data disk N.
* `name` - (Optional) The name of data disk N. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (.), underscores (_), and hyphens (-).
* `snapshot_id` - (Optional) The ID of the snapshot to be used to create data disk N. Valid values of N: 1 to 16. When DataDisk.N.SnapshotId is specified, DataDisk.N.Size is ignored. The data disk is created based on the size of the specified snapshot. Use snapshots that were created on or after July 15, 2013. Otherwise, an error is returned and your request is rejected.

### `taints`

The following arguments are supported in the `taints` configuration block:

* `key` - (Optional) The taint key.
* `value` - (Optional) The taint value.
* `effect` - (Optional) The taint effect.

The following example is the definition of taints block:

```
resource "alicloud_cs_managed_kubernetes" "k8s" {
  # ... other configuration ...

  #  defining two taints
  taints {
    key = "key-a"
    value = "value-a"
    effect = "NoSchedule"
  }
  taints {
    key = "key-b"
    value = "value-b"
    effect = "NoSchedule"
  }
}
```

### `log_config`

The following arguments are supported in the `log_config` configuration block:

* `type` - (Required) Type of collecting logs, only `SLS` are supported currently.
* `project` - (Optional) Log Service project name, cluster logs will output to this project.

### `runtime`

* `name` - (Optional) The name of the runtime. Supported runtimes can be queried by data source [alicloud_cs_kubernetes_version](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_version).
* `version` - (Optional) The version of the runtime.

The following example is the definition of runtime block:

```
  runtime = {
    name = "containerd"
    version = "1.6.28"
  }
```

### tags

The following example is the definition of tags block. The type of this field is map:

```
  # for example, define three tags

  tags = {
    "key1" = "value1"
    "key2" = "value2"
    "name" = "tf"
  }
```

### `delete_options`

The following arguments are supported in the `delete_options` configuration block:
* `delete_mode` - (Optional) The deletion mode of the cluster. Different resources may have different default behavior, see `resource_type` for details. Valid values:
  - `delete`: delete resources created by the cluster.
  - `retain`: retain resources created by the cluster.
* `resource_type` - (Optional) The type of resources that are created by cluster. Valid values:
  - `SLB`: SLB resources created by the Nginx Ingress Service, default behavior is to delete, option to retain is available.  
  - `ALB`: ALB resources created by the ALB Ingress Controller, default behavior is to retain, option to delete is available. 
  - `SLS_Data`: SLS Project used by the cluster logging feature, default behavior is to retain, option to delete is available. 
  - `SLS_ControlPlane`: SLS Project used for the managed cluster control plane logs, default behavior is to retain, option to delete is available.

```
  ...
  // Specify delete_options as below when deleting cluster
  // delete SLB resources created by the Nginx Ingress Service
  delete_options {
    delete_mode = "delete"
    resource_type = "SLB"
  }
  // delete ALB resources created by the ALB Ingress Controller
  delete_options {
    delete_mode = "delete"
    resource_type = "ALB"
  }
  // delete SLS Project used by the cluster logging feature
  delete_options {
    delete_mode = "delete"
    resource_type = "SLS_Data"
  }
  // delete SLS Project used for the managed cluster control plane logs
  delete_options {
    delete_mode = "delete"
    resource_type = "SLS_ControlPlane"
  }
```
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `slb_internet` - The public ip of load balancer.
* `slb_id` - The ID of APIServer load balancer.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_nodes` - (Removed from version 1.212.0) List of cluster worker nodes.
  * `id` - ID of the node.
  * `name` - Node name.
  * `private_ip` - The private IP address of node.
* `connections` - Map of kubernetes cluster connection information.
  * `api_server_internet` - API Server Internet endpoint.
  * `api_server_intranet` - API Server Intranet endpoint.
  * `master_public_ip` - Master node SSH IP address.
  * `service_domain` - Service Access Domain.
* `worker_ram_role_name` - The RamRole Name attached to worker node.
* `certificate_authority` - (Available since v1.105.0) Nested attribute containing certificate authority data for your cluster.
  * `cluster_cert` - The base64 encoded cluster certificate data required to communicate with your cluster. Add this to the certificate-authority-data section of the kubeconfig file for your cluster.
  * `client_cert` - The base64 encoded client certificate data required to communicate with your cluster. Add this to the client-certificate-data section of the kubeconfig file for your cluster.
  * `client_key` - The base64 encoded client key data required to communicate with your cluster. Add this to the client-key-data section of the kubeconfig file for your cluster.
* `rrsa_metadata` - (Optional, Available since v1.185.0) Nested attribute containing RRSA related data for your cluster.
  * `enabled` - Whether the RRSA feature has been enabled.
  * `rrsa_oidc_issuer_url` - The issuer URL of RRSA OIDC Token.
  * `ram_oidc_provider_name` - The name of OIDC Provider that was registered in RAM.
  * `ram_oidc_provider_arn` -  The arn of OIDC provider that was registered in RAM.

## Timeouts

-> **NOTE:** Available since v1.58.0.
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

## Import

Kubernetes managed cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`.

```shell
$ terraform import alicloud_cs_managed_kubernetes.main cluster_id
```

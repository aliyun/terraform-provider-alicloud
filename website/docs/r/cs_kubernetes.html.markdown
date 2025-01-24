---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-kubernetes"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster.
---

# alicloud_cs_kubernetes

This resource will help you to manage a Kubernetes Cluster in Alibaba Cloud Kubernetes Service, see [What is kubernetes](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/create-an-ask-cluster-1).

-> **NOTE:** From August 21, 2024, Container Service for Kubernetes (ACK) discontinues the creation of ACK dedicated clusters, see [Product announcement](https://www.alibabacloud.com/help/en/ack/product-overview/product-announcement-announcement-on-stopping-new-ack-dedicated-cluster) for more details.

-> **NOTE:** Available since v1.9.0.

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Each kubernetes cluster contains 3 master nodes and those number cannot be changed at now.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.

-> **NOTE:** From version 1.9.4, the provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** From version 1.16.0, the provider supports Multiple Availability Zones Kubernetes Cluster. To create a cluster of this kind, you must specify 3 or 5 items in `master_vswitch_ids` and `master_instance_types`.

-> **NOTE:** From version 1.20.0, the provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** You need to activate several other products and confirm Authorization Policy used by Container Service before using this resource.
Please refer to the `Authorization management` and `Cluster management` sections in the [Document Center](https://www.alibabacloud.com/help/doc-detail/86488.htm).

-> **NOTE:** From version 1.75.0, Some parameters have been removed from resource,You can check them below and re-import the cluster if necessary.

-> **NOTE:** From version 1.101.0+, We supported the `professional managed clusters(ack-pro)`, You can create a pro cluster by setting the the value of `cluster_spec`.

-> **NOTE:** From version 1.177.0+, `exclude_autoscaler_nodes`,`worker_number`,`worker_vswitch_ids`,`worker_instance_types`,`worker_instance_charge_type`,`worker_period`,`worker_period_unit`,`worker_auto_renew`,`worker_auto_renew_period`,`worker_disk_category`,`worker_disk_size`,`worker_data_disks`,`node_port_range`,`cpu_policy`,`user_data`,`taints`,`worker_disk_performance_level`,`worker_disk_snapshot_policy_id` are deprecated.
We Suggest you using resource **`alicloud_cs_kubernetes_node_pool`** to manage your cluster worker nodes.

-> **NOTE:** From version 1.212.0, `exclude_autoscaler_nodes`,`worker_number`,`worker_vswitch_ids`,`worker_instance_types`,`worker_instance_charge_type`,`worker_period`,`worker_period_unit`,`worker_auto_renew`,`worker_auto_renew_period`,`worker_disk_category`,`worker_disk_size`,`worker_data_disks`,`node_port_range`,`cpu_policy`,`user_data`,`taints`,`worker_disk_performance_level`,`worker_disk_snapshot_policy_id`,`kube_config`,`availability_zone` are removed.
Please use resource **`alicloud_cs_kubernetes_node_pool`** to manage your cluster worker nodes.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_kubernetes&exampleId=8ecfabbf-301e-0a64-48fa-e099644a8db27cad3ce9&activeTab=example&spm=docs.r.cs_kubernetes.0.8ecfabbf30&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-kubernetes-example"
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
  default     = ["10.1.0.0/16", "10.2.0.0/16", "10.3.0.0/16"]
}

variable "terway_vswitch_ids" {
  description = "List of existing vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "terway_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'terway_vswitch_cidrs' is not specified."
  type        = list(string)
  default     = ["10.4.0.0/16", "10.5.0.0/16", "10.6.0.0/16"]
}

variable "cluster_addons" {
  type = list(object({
    name   = string
    config = map(string)
  }))

  default = [
    # If use terway network, must specify addons with `terway-eniip` and param `pod_vswitch_ids`
    {
      "name"   = "terway-eniip",
      "config" = {},
    },
    {
      "name"   = "csi-plugin",
      "config" = {},
    },
    {
      "name"   = "csi-provisioner",
      "config" = {},
    },
    {
      "name" = "logtail-ds",
      "config" = {
        "IngressDashboardEnabled" = "true",
      }
    },
    {
      "name" = "nginx-ingress-controller",
      "config" = {
        "IngressSlbNetworkType" = "internet"
      }
    },
    {
      "name"   = "arms-prometheus",
      "config" = {},
    },
    {
      "name" = "ack-node-problem-detector",
      "config" = {
        "sls_project_name" = ""
      },
    }
  ]
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
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index < length(data.alicloud_enhanced_nat_available_zones.enhanced.zones) ? count.index : 0].zone_id
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "terway_vswitches" {
  count      = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.terway_vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index < length(data.alicloud_enhanced_nat_available_zones.enhanced.zones) ? count.index : 0].zone_id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_instance_types" "cloud_essd" {
  count                = 3
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index < length(data.alicloud_enhanced_nat_available_zones.enhanced.zones) ? count.index : 0].zone_id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
}

resource "alicloud_cs_kubernetes" "default" {
  master_vswitch_ids    = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)) : length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids       = length(var.terway_vswitch_ids) > 0 ? split(",", join(",", var.terway_vswitch_ids)) : length(var.terway_vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  master_instance_types = [data.alicloud_instance_types.cloud_essd.0.instance_types.0.id, data.alicloud_instance_types.cloud_essd.1.instance_types.0.id, data.alicloud_instance_types.cloud_essd.2.instance_types.0.id]
  master_disk_category  = "cloud_essd"
  password              = "Yourpassword1234"
  service_cidr          = "172.18.0.0/16"
  install_cloud_monitor = "true"
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  deletion_protection   = "false"
  timezone              = "Asia/Shanghai"
  os_type               = "Linux"
  platform              = "AliyunLinux3"
  cluster_domain        = "cluster.local"
  proxy_mode            = "ipvs"
  custom_san            = "www.terraform.io"
  new_nat_gateway       = "true"
  dynamic "addons" {
    for_each = var.cluster_addons
    content {
      name   = lookup(addons.value, "name", var.cluster_addons)
      config = jsonencode(lookup(addons.value, "config", var.cluster_addons))
    }
  }
}
```

## Argument Reference

*Global params*

* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `name_prefix` - (Optional, Deprecated) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `timezone` - (Optional, ForceNew, Available since v1.103.2) When you create a cluster, set the time zones for the Master and Worker nodes. You can only change the managed node time zone if you create a cluster. Once the cluster is created, you can only change the time zone of the Worker node.
* `resource_group_id` - (Optional, Available since v1.101.0) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `version` - (Optional, Available since v1.70.1) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `runtime` - (Optional, Available since v1.103.2) The runtime of containers. If you select another container runtime, see [How do I select between Docker and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm?spm=a2c63.p38356.b99.440.22563866AJkBgI). See [`runtime`](#runtime) below.
* `enable_ssh` - (Optional, ForceNew) Enable login to the node through SSH. Default to `false`.
* `rds_instances` - (Optional, Available since v1.103.2) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `security_group_id` - (Optional, ForceNew, Available since v1.91.0) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, ForceNew, Available since v1.91.0) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm). Only works for **Create** Operation. 
* `proxy_mode` - (Optional, ForceNew) Proxy mode is option of kube-proxy. options: iptables | ipvs. default: ipvs.
* `image_id` - (Optional, ForceNew) Custom Image support. Must be based on AliyunLinux or AliyunLinux3.
* `cluster_domain` - (Optional, ForceNew, Available since v1.103.2) Cluster local domain name, Default to `cluster.local`. A domain name consists of one or more sections separated by a decimal point (.), each of which is up to 63 characters long, and can be lowercase, numerals, and underscores (-), and must be lowercase or numerals at the beginning and end.
* `custom_san` - (Optional, ForceNew, Available since v1.103.2) Customize the certificate SAN, multiple IP or domain names are separated by English commas (,).
* `user_ca` - (Optional) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `deletion_protection` - (Optional, Available since v1.103.2)  Whether to enable cluster deletion protection.
* `install_cloud_monitor` - (Optional, ForceNew) Install cloud monitor agent on ECS. Default to `true`.
* `service_account_issuer` - (Optional, ForceNew, Available since v1.92.0) The issuer of the Service Account token for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm), corresponds to the `iss` field in the token payload. Set this to `"https://kubernetes.default.svc"` to enable the Token Volume Projection feature (requires specifying `api_audiences` as well). From cluster version 1.22+, Service Account Token Volume Projection will be enabled by default.
* `api_audiences` - (Optional, ForceNew, Available since v1.92.0) A list of API audiences for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm). Set this to `["https://kubernetes.default.svc"]` if you want to enable the Token Volume Projection feature requires specifying `service_account_issuer` as well. From cluster version 1.22+, Service Account Token Volume Projection will be enabled by default.
* `tags` - (Optional, Available since v1.97.0) Default nil, A map of tags assigned to the kubernetes cluster and work nodes.
* `load_balancer_spec` - (Optional, Deprecated since v1.232.0) The cluster api server load balancer instance specification. For more information on how to select a LB instance specification, see [SLB instance overview](https://help.aliyun.com/document_detail/85931.html). Only works for **Create** Operation. The spec will not take effect because the charge of the load balancer has been changed to PayByCLCU. 
* `retain_resources` - (Optional, Available since v1.141.0) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.
* `delete_options` - (Optional, Available since v1.223.2) Delete options, only work for deleting resource. Make sure you have run `terraform apply` to make the configuration applied. See [`delete_options`](#delete_options) below.
* `password` - (Optional, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Optional, ForceNew) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `os_type` - (Optional, ForceNew, Available since v1.103.2) The operating system of the nodes that run pods, its valid value is either `Linux` or `Windows`. Default to `Linux`.
* `platform` - (Optional, ForceNew, Available since v1.103.2) The architecture of the nodes that run pods, its valid value `AliyunLinux`, `AliyunLinux3`. Default to `AliyunLinux3`.
* `node_name_mode` - (Optional, ForceNew, Available since v1.88.0) Each node name consists of a prefix, an IP substring, and a suffix, the input format is `customized,<prefix>,IPSubStringLen,<suffix>`. For example "customized,aliyun.com-,5,-test", if the node IP address is 192.168.59.176, the prefix is aliyun.com-, IP substring length is 5, and the suffix is -test, the node name will be aliyun.com-59176-test.
* `addons` - (Optional, Available since v1.88.0) The addon you want to install in cluster. See [`addons`](#addons) below. Only works for **Create** Operation, use [resource cs_kubernetes_addon](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_addon) to manage addons if cluster is created.

*Network params*

* `pod_cidr` - (Optional, ForceNew) [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `pod_vswitch_ids` - (Optional) - [Terway Specific] The vswitches for the pod network when using Terway. It is recommended that `pod_vswitch_ids` is not belong to `worker_vswitch_ids` and `master_vswitch_ids` but must be in same availability zones. Only works for **Create** Operation. 
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice. Your cluster nodes and applications will have public network access. If there is a NAT gateway in the selected VPC, ACK will use this gateway by default; if there is no NAT gateway in the selected VPC, ACK will create a new NAT gateway for you and automatically configure SNAT rules. Only works for **Create** Operation. 
* `service_cidr` - (Optional, ForceNew) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional, ForceNew) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true. Only works for **Create** Operation. 

-> **NOTE:** If you want to use `Terway` as CNI network plugin, You need to specify the `pod_vswitch_ids` field and addons with `terway-eniip`.
If you want to use `Flannel` as CNI network plugin, You need to specify the `pod_cidr` field and addons with `flannel`.

*Master params*

* `master_vswitch_ids` - (Required, ForceNew) The vswitches used by master, you can specific 3 or 5 vswitches because of the amount of masters. Detailed below.
* `master_instance_types` - (Required, ForceNew) The instance type of master node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `master_instance_charge_type` - (Optional, ForceNew) Master payment type. or `PostPaid` or `PrePaid`, defaults to `PostPaid`. If value is `PrePaid`, the files `master_period`, `master_period_unit`, `master_auto_renew` and `master_auto_renew_period` are required.
* `master_period` - (Optional, ForceNew) Master payment period.Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `master_period_unit` - (Optional, ForceNew) Master payment period unit, the valid value is `Month`.
* `master_auto_renew` - (Optional, ForceNew) Enable master payment auto-renew, defaults to false.
* `master_auto_renew_period` - (Optional, ForceNew) Master payment auto-renew period, it can be one of {1, 2, 3, 6, 12}.
* `master_disk_category` - (Optional, ForceNew) The system disk category of master node. Its valid value are `cloud_ssd`, `cloud_essd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Optional, ForceNew) The system disk size of master node. Its valid value range [20~500] in GB. Default to 20.
* `master_disk_performance_level` - (Optional, ForceNew, Available since v1.120.0) Master node system disk performance level. When `master_disk_category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `master_disk_snapshot_policy_id` - (Optional, ForceNew, Available since v1.120.0) Master node system disk auto snapshot policy.

*Computed params*

* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

*Removed params*

* `master_instance_type` - (Removed) The instance type of master node.
* `worker_instance_type` - (Removed) The instance type of worker node.
* `vswitch_id` - (Removed) The vswitch where new kubernetes cluster will be located. If it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `vswitch_ids` - (Removed) The vswitch where new kubernetes cluster will be located. Specify one or more vswitch's id. It must be in the zone which `availability_zone` specified.
* `force_update` - (Removed) Whether to force the update of kubernetes cluster arguments. Default to false.
* `log_config` - (Removed) A list of one element containing information about the associated log store. See [`log_config`](#log_config) below.
* `cluster_network_type` - (Removed) The network that cluster uses, use `flannel` or `terway`.
* `worker_data_disk_category` - (Removed) The data disk category of worker, use `worker_data_disks` to instead it.
* `worker_data_disk_size` - (Removed) The data disk size of worker, use `worker_data_disks` to instead it.
* `worker_numbers` - (Removed) The number of workers, use `worker_number` to instead it.
* `nodes` - (Removed) The master nodes, use `master_nodes` to instead it.
* `exclude_autoscaler_nodes` - (Removed since v1.212.0) Exclude autoscaler nodes from `worker_nodes`. Default to `false`.
* `kube_config` - (Removed since v1.212.0) The path of kube config, like `~/.kube/config`. You can set some file paths to save kube_config information, but this way is cumbersome. Since version 1.105.0, we've written it to tf state file. About its use，see export attribute certificate_authority. From version 1.187.0+, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kube_config.
* `availability_zone` - (Removed since v1.212.0) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.
* `worker_number` - (Removed since v1.212.0) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswitch_ids` - (Removed since v1.212.0) The vswitches used by workers.
* `worker_instance_types` - (Removed since v1.212.0) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_instance_charge_type` - (Removed since v1.212.0) Worker payment type, its valid value is either or `PostPaid` or `PrePaid`. Defaults to `PostPaid`. If value is `PrePaid`, the files `worker_period`, `worker_period_unit`, `worker_auto_renew` and `worker_auto_renew_period` are required, default is `PostPaid`.
* `worker_period` - (Removed since v1.212.0) Worker payment period. The unit is `Month`. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `worker_period_unit` - (Removed since v1.212.0) Worker payment period unit, the valid value is `Month`.
* `worker_auto_renew` - (Removed since v1.212.0) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Removed since v1.212.0) Worker payment auto-renew period, it can be one of {1, 2, 3, 6, 12}.
* `worker_disk_category` - (Removed since v1.212.0) The system disk category of worker node. Its valid value are `cloud`, `cloud_ssd`, `cloud_essd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Removed since v1.212.0) The system disk size of worker node. Its valid value range [40~500] in GB.
* `worker_data_disks` - (Removed since v1.212.0) The data disk configurations of worker nodes, such as the disk type and disk size. See [`worker_data_disks`](#worker_data_disks) below.
* `node_port_range`- (Removed since v1.212.0) The service port range of nodes, valid values: `30000` to `65535`. Default to `30000-32767`.
* `cpu_policy` - (Removed since v1.212.0) Kubelet cpu policy. For Kubernetes 1.12.6 and later, its valid value is either `static` or `none`. Default to `none`.
* `user_data` - (Removed since v1.212.0) Custom data that can execute on nodes. For more information, see [Prepare user data](https://www.alibabacloud.com/help/doc-detail/49121.htm).
* `taints` - (Removed since v1.212.0) Taints ensure pods are not scheduled onto inappropriate nodes. One or more taints are applied to a node; this marks that the node should not accept any pods that do not tolerate the taints. For more information, see [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/). See [`taints`](#taints) below.
* `worker_disk_performance_level` - (Removed since v1.212.0) Worker node system disk performance level, when `worker_disk_category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `worker_disk_snapshot_policy_id` - (Removed since v1.212.0) Worker node system disk auto snapshot policy.

### `taints`

The taints supports the following:

* `key` - (Optional) The key of a taint.
* `value` - (Optional) The key of a taint.
* `effect` - (Optional) The scheduling policy. Valid values: NoSchedule | NoExecute | PreferNoSchedule. Default value: NoSchedule.

### `addons`

The addons supports the following:

* `name` - (Optional) Name of the ACK add-on. The name must match one of the names returned by [DescribeAddons](https://help.aliyun.com/document_detail/171524.html).
* `config` - (Optional) The ACK add-on configurations. For more config information, see [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata).
* `version` - (Optional) The version of the component.
* `disabled` - (Optional) Disables the automatic installation of a component. Default is `false`.

The following example is the definition of addons block, The type of this field is list:

```
# install nginx ingress, conflict with SLB ingress
addons {
  name = "nginx-ingress-controller"
  # use internet
  config = "{\"IngressSlbNetworkType\":\"internet",\"IngressSlbSpec\":\"slb.s2.small\"}"
  # if use intranet, detail below.
  # config = "{\"IngressSlbNetworkType\":\"intranet",\"IngressSlbSpec\":\"slb.s2.small\"}"
}
```

### `worker_data_disks`

The worker_data_disks supports the following:

* `category` - (Optional) The type of the data disks. Valid values: `cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`. Default to `cloud_efficiency`.
* `size` - (Optional) The size of a data disk, Its valid value range [40~32768] in GB. Unit: GiB.
* `encrypted` - (Optional) Specifies whether to encrypt data disks. Valid values: true and false.
* `performance_level` - (Optional, Available since v1.120.0) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `auto_snapshot_policy_id` - (Optional, Available since v1.120.0) Worker node data disk auto snapshot policy.
* `snapshot_id` - (Optional) The id of snapshot.
* `kms_key_id` - (Optional) The id of the kms key.
* `name` - (Optional) The name of the data disks.
* `device` - (Optional) The device of the data disks.

### `log_config`

The log_config supports the following:

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
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.
* `connections` - (Map) Map of kubernetes cluster connection information.
  * `api_server_internet` - API Server Internet endpoint.
  * `api_server_intranet` - API Server Intranet endpoint.
  * `master_public_ip` - Master node SSH IP address.
  * `service_domain` - Service Access Domain.
* `certificate_authority` - (Map, Available since v1.105.0) Nested attribute containing certificate authority data for your cluster.
  * `cluster_cert` - The base64 encoded cluster certificate data required to communicate with your cluster. Add this to the certificate-authority-data section of the kubeconfig file for your cluster.
  * `client_cert` - The base64 encoded client certificate data required to communicate with your cluster. Add this to the client-certificate-data section of the kubeconfig file for your cluster.
  * `client_key` - The base64 encoded client key data required to communicate with your cluster. Add this to the client-key-data section of the kubeconfig file for your cluster.
* `slb_id` - The ID of APIServer load balancer.
* `master_nodes` - (Optional) The master nodes. See [`master_nodes`](#master_nodes) below.
* `worker_nodes` - (Removed since v1.212.0) List of cluster worker nodes. See [`worker_nodes`](#worker_nodes) below.

### `master_nodes`

The master_nodes supports the following:

* `id` - (Optional) The id of a node.
* `name` - (Optional) The name of a node.
* `private_ip` - (Optional) The private ip of a node.

### `worker_nodes`

The worker_nodes supports the following:

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.

## Timeouts

-> **NOTE:** Available since v1.58.0.
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

## Import

Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`.

```shell
$ terraform import alicloud_cs_kubernetes.main cluster-id
```

---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_edge_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-edge-kubernetes"
description: |-
  Provides a Alicloud resource to manage container edge kubernetes cluster.
---

# alicloud\_cs\_edge\_kubernetes

This resource will help you to manage a Edge Kubernetes Cluster in Alibaba Cloud Kubernetes Service. 

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.

-> **NOTE:** The provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** The provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** Available in v1.103.0+.

## Example Usage

```
# If vpc_id is not specified, a new one will be created
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count      = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = element(var.availability_zone, count.index)
}

resource "alicloud_cs_edge_kubernetes" "k8s" {
  name                  = var.cluster_name
  count                 = var.k8s_number
  worker_vswitch_ids    = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)): length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  worker_instance_types = var.worker_instance_types
  worker_number         = var.worker_number
  node_cidr_mask        = var.node_cidr_mask
  install_cloud_monitor = var.install_cloud_monitor
  proxy_mode            = var.proxy_mode
  password              = var.password
  service_cidr          = var.service_cidr
  pod_cidr              = var.pod_cidr
  # version can not be defined in variables.tf.
  version               = "1.20.11-aliyunedge.1"

  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name   = lookup(addons.value, "name", var.cluster_addons)
        config = lookup(addons.value, "config", var.cluster_addons)
      }
  }
  slb_internet_enabled         = var.slb_enabled
  is_enterprise_security_group = var.enterprise_sg
}

```

## Argument Reference

The following arguments are supported:

### Global params

* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `security_group_id` - (Optional) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).
* `rds_instance` - (Optional, Available in 1.103.2+) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `resource_group_id` - (Optional, ForceNew, Available in 1.103.2+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `deletion_protection` - (Optional, Available in 1.103.2+)  Whether to enable cluster deletion protection.
* `force_update` - (Optional, ForceNew, Available in 1.113.0+) Default false, when you want to change `vpc_id`, you have to set this field to true, then the cluster will be recreated.
* `tags` - (Optional, Available in 1.120.0+) Default nil, A map of tags assigned to the kubernetes cluster and work node.
* `retain_resources` - (Optional, Available in 1.141.0+) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.

### Network params

* `pod_cidr` - (Required) [Flannel Specific] The CIDR block for the pod network when using Flannel.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

->NOTE: If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

### Worker params

* `password` - (**Required**, Sensitive) The password of ssh login cluster node. You have to specify one of `password`, `key_name` `kms_encrypted_password` fields.
* `key_name` - (**Required**) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `worker_number` - (**Required**) The cloud worker node number of the edge kubernetes cluster. Default to 1. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswtich_ids` - (**Required**) The vswitches used by workers.
* `worker_instance_types` - (**Required**, ForceNew) The instance types of worker node, you can set multiple types to avoid NoStock of a certain type
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_efficiency`, `cloud_ssd` and `cloud_essd` and . Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.
* `worker_data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size.
  * `category` - The type of the data disks. Valid values: `cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`. Default to `cloud_efficiency`.
  * `size` - The size of a data disk, at least 40. Unit: GiB.
  * `encrypted` - Specifies whether to encrypt data disks. Valid values: true and false. Default is `false`.
  * `performance_level` - (Optional, Available in 1.120.0+) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
  * `auto_snapshot_policy_id` - (Optional, Available in 1.120.0+) Worker node data disk auto snapshot policy.
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. default: `true`.
* `proxy_mode` - Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `worker_disk_performance_level` - (Optional, Available in 1.120.0+) Worker node system disk performance level, when `worker_disk_category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `worker_disk_snapshot_policy_id` - (Optional, Available in 1.120.0+) Worker node system disk auto snapshot policy.

### Addons

You can specific network plugin,log component,ingress component and so on.You can get more information about addons on [ACK](https://cs.console.aliyun.com/) web console. When you create a ACK cluster. You can get openapi-spec before creating the cluster on submission page.

To create an edge cluster, you only need to install additional logging components. addons name:

* `logtail-ds-docker` - Logtail addon.
* `alibaba-log-controller` - In edge cluster, you should set the `IngressDashboardEnabled` to false in config. Detail below.
* `alicloud-monitor-controller` - Install the cloud monitoring plug-in on the node to view monitoring information for the created ECS instance in the cloud monitoring console. Detail below.

The `main.tf`:

```
resource "alicloud_cs_edg_kubernetes" "k8s" {
  # ... other configuration ...

  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name          = lookup(addons.value, "name", var.cluster_addons)
        config        = lookup(addons.value, "config", var.cluster_addons)
      }
  }
}
```

The `varibales.tf`:

```
# flannel network component, Required.
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

# logtail and monitor component, optional
variable "cluster_addons" {
  type = list(object({
      name      = string
      config    = string
  }))

  default = [
    {
      "name"     = "logtail-ds-docker",
      "config"   = "",
    }
    {
      "name": "alibaba-log-controller",
      "config": "{\"IngressDashboardEnabled\":\"false\"}"
    }
    {
      "name": "alicloud-monitor-controller",
      "config"   = "",
    }
  ]
}
```

### Computed params (No need to configure)

You can set some file paths to save kube_config information, but this way is cumbersome. Since version 1.105.0, we've written it to tf state file. About its use，see export attribute certificate_authority.

* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_nodes` - List of cluster worker nodes.
  * `id` - ID of the node.
  * `name` - Node name.
  * `private_ip` - The private IP address of node.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.
* `certificate_authority` - (Available in 1.105.0+) Nested attribute containing certificate authority data for your cluster.
  * `cluster_cert` - The base64 encoded cluster certificate data required to communicate with your cluster. Add this to the certificate-authority-data section of the kubeconfig file for your cluster.
  * `client_cert` - The base64 encoded client certificate data required to communicate with your cluster. Add this to the client-certificate-data section of the kubeconfig file for your cluster.
  * `client_key` - The base64 encoded client key data required to communicate with your cluster. Add this to the client-key-data section of the kubeconfig file for your cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster. 

## Import

Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`

```
  $ terraform import alicloud_cs_edge_kubernetes.main cluster-id
```

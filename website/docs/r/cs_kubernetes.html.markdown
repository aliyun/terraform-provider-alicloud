---
subcategory: "Container Service (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-kubernetes"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster.
---

# alicloud\_cs\_kubernetes

This resource will help you to manage a Kubernetes Cluster in Alibaba Cloud Kubernetes Service. 

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
## Example Usage
```$xslt

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count             = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id            = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block        = element(var.vswitch_cidrs, count.index)
  availability_zone = element(var.availability_zone, count.index)
}

resource "alicloud_cs_kubernetes" "k8s" {
  count                 = 1
  master_vswitch_ids    = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)): length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  worker_vswitch_ids    = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)): length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  master_instance_types = var.master_instance_types
  worker_instance_types = var.worker_instance_types
  worker_number         = var.worker_number
  node_cidr_mask        = var.node_cidr_mask
  enable_ssh            = var.enable_ssh
  install_cloud_monitor = var.install_cloud_monitor
  cpu_policy            = var.cpu_policy
  proxy_mode            = var.proxy_mode
  password              = var.password
  pod_cidr              = var.pod_cidr
  service_cidr          = var.service_cidr
  # version can not be defined in variables.tf. Options: 1.16.6-aliyun.1|1.14.8-aliyun.1
  version               = "1.16.6-aliyun.1"
  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name                    = lookup(addons.value, "name", var.cluster_addons)
        config                  = lookup(addons.value, "config", var.cluster_addons)
      }
  }
}
```

## Argument Reference

The following arguments are supported:

#### Global params
* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional, Available since 1.70.1) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Required) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Required, Available in 1.57.1+) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `user_ca` - (Optional, ForceNew) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `enable_ssh` - (Optional) Enable login to the node through SSH. default: false 
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. default: true 
* `cpu_policy` - kubelet cpu policy. options: static|none. default: none.
* `proxy_mode` - Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `image_id` - Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `user_data` - (Optional, Available in 1.81.0+) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `exclude_autoscaler_nodes` - (Optional, Available in 1.88.0+) Exclude autoscaler nodes from `worker_nodes`. default: false 
* `node_name_mode` - (Optional, Available in 1.88.0+) Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test. 
* `worker_data_disks` - (Optional, Available in 1.91.0+) The data disk configurations of worker nodes, such as the disk type and disk size. 
  - category: the type of the data disks. Valid values:
      + cloud: basic disks.
      + cloud_efficiency: ultra disks.
      + cloud_ssd: SSDs.
  - size: the size of a data disk. Unit: GiB.
  - encrypted: specifies whether to encrypt data disks. Valid values: true and false.
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, Available in 1.91.0+) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).
* `service_account_issuer` - (Optional, ForceNew, Available in 1.92.0+) The issuer of the Service Account token for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm), corresponds to the `iss` field in the token payload. Set this to `"kubernetes.default.svc"` to enable the Token Volume Projection feature (requires specifying `api_audiences` as well).
* `api_audiences` - (Optional, ForceNew, Available in 1.92.0+) A list of API audiences for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm). Set this to `["kubernetes.default.svc"]` if you want to enable the Token Volume Projection feature (requires specifying `service_account_issuer` as well.

#### Addons 
It is a new field since 1.75.0. You can specific network plugin,log component,ingress component and so on.     
 
```$xslt
  main.tf
   
  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name                    = lookup(addons.value, "name", var.cluster_addons)
        config                  = lookup(addons.value, "config", var.cluster_addons)
      }
  }
```
```$xslt
    varibales.tf 
    
    // Network-flannel 
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
       
    
    // Network-terway 
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
    
    // Storage-csi
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
    
    // Storage-flexvolume
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
    
    // Log
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
    
    // Ingress
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
    
    // Ingress-Disable
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
    } 
  
```
* `logtail-ds` - You can specific `IngressDashboardEnabled` and `sls_project_name` in config. If you switch on `IngressDashboardEnabled` and `sls_project_name`,then logtail-ds would use `sls_project_name` as default log store.
* `nginx-ingress-controller` - You can specific `IngressSlbNetworkType` in config. Options: internet|intranet.     
You can get more information about addons on ACK web console. When you create a ACK cluster. You can get openapi-spec before creating the cluster on submission page. 



#### Network
* `pod_cidr` - (Required) [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `pod_vswitch_ids` - (Required) [Terway Specific] The vswitches for the pod network when using Terway.Be careful the `pod_vswitch_ids` can not equal to `worker_vswtich_ids` or `master_vswtich_ids` but must be in same availability zones.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

If you want to use `Terway` as CNI network plugin, You need to specific the `pod_vswitch_ids` field and addons with `terway-eniip`.    
If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

#### Master params
* `master_vswtich_ids` - (Required) The vswitches used by master, you can specific 3 or 5 vswitches because of the amount of masters. You can also specific 
* `master_instance_types` - (Required) The instance type of master node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `master_instance_charge_type` - (Optional) Master payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `master_period_unit` - (Optional) Master payment period unit. `Month` or `Week`, defaults to `Month`.
* `master_period` - (Optional) Master payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `master_auto_renew` - (Optional) Enable master payment auto-renew, defaults to false.
* `master_auto_renew_period` - (Optional) Master payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `master_disk_category` - (Optional) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Optional) The system disk size of master node. Its valid value range [20~500] in GB. Default to 20.

#### Worker params 
* `worker_number` - (Required) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswtich_ids` - (Required) The vswitches used by workers. 
* `worker_instance_types` - (Required, ForceNew) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_instance_charge_type` - (Optional, Force new resource) Worker payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `worker_period_unit` - (Optional) Worker payment period unit. `Month` or `Week`, defaults to `Month`.
* `worker_period` - (Optional) Worker payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `worker_auto_renew` - (Optional) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Optional) Worker payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.

#### Computed params (No need to configure) 
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `availability_zone` - (Optional) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.

#### Removed params (Never Supported)
* `master_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of master node.
* `worker_instance_type` - (Deprecated from version 1.16.0)(Required, Force new resource) The instance type of worker node.
* `vswitch_id` - (Deprecated from version 1.16.0)(Force new resource) The vswitch where new kubernetes cluster will be located. If it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `vswitch_ids` - (Required, ForceNew) The vswitch where new kubernetes cluster will be located. Specify one or more vswitch's id. It must be in the zone which `availability_zone` specified.
* `force_update` - (Optional, Available in 1.50.0+) Whether to force the update of kubernetes cluster arguments. Default to false.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `log_config` - (Optional, ForceNew) A list of one element containing information about the associated log store. It contains the following attributes:
  * `type` - Type of collecting logs, only `SLS` are supported currently.
  * `project` - Log Service project name, cluster logs will output to this project.
* `cluster_network_type` - (Optional) The network that cluster uses, use `flannel` or `terway`.
  
  
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
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.

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

Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`

```
$ terraform import alicloud_cs_kubernetes.main cluster-id
```

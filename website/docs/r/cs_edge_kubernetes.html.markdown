---
subcategory: "Container Service (CS)"
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

-> **NOTE:** From version 1.9.4, the provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** From version 1.20.0, the provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

## Example Usage
```$xslt
// If vpc_id is not specified, a new one will be created
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


resource "alicloud_cs_edge_kubernetes" "k8s" {
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
  # version can not be defined in variables.tf. Options: 1.14.8-aliyunedge.1|1.12.6-aliyunedge.2
  version               = "1.12.6-aliyunedge.2"
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
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. default: true 
* `proxy_mode` - Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `user_data` - (Optional, Available in 1.81.0+) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `worker_data_disks` - (Optional, Available in 1.91.0+) The data disk configurations of worker nodes, such as the disk type and disk size. 
  - category: the type of the data disks. Valid values:
      + cloud: basic disks.
      + cloud_efficiency: ultra disks.
      + cloud_ssd: SSDs.
  - size: the size of a data disk. Unit: GiB.
  - encrypted: specifies whether to encrypt data disks. Valid values: true and false.
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, Available in 1.91.0+) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).

#### Addons 
It is a new field since 1.75.0. You can specific network plugin,log component,ingress component and so on.     

-> **NOTE:** If you want to upgrade provider to 1.90.0+, you need to pay attention to the disabled value. If the value is `""`, you need to modify it to `"false"`, and then run `terraform apply` to make it effect. After that, you can modify the provider version to upgrade smoothly. Otherwise, there will throw an error: `Error: a bool is required`.
 
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
       
    
    // Log
    variable "cluster_addons" {
        type = list(object({
            name      = string
            config    = string
        }))
    
        default = [
            {
                "name"     = "logtail-ds-docker",
                "config"   = "{\"IngressDashboardEnabled\":\"false\"}",
            }
            {
                "name"     = "alibaba-log-controller",
                "config"   = "",
            }
        ]
    } 
 
  
```
* `logtail-ds-docker` - You can specific `IngressDashboardEnabled` and `sls_project_name` in config. If you switch on `IngressDashboardEnabled` and `sls_project_name`,then logtail-ds-docker would use `sls_project_name` as default log store. 

You can get more information about addons on ACK web console. When you create a ACK cluster. You can get openapi-spec before creating the cluster on submission page. 



#### Network
* `pod_cidr` - (Required) [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

#### Worker params 
* `worker_number` - (Required) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswtich_ids` - (Required) The vswitches used by workers. 
* `worker_instance_types` - (Required, ForceNew) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.

#### Computed params (No need to configure) 
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`

  
### Timeouts
-> **NOTE:** Available in 1.58.0+.
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
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
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.

### Block Nodes
* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.

### Block Connections
* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `service_domain` - Service Access Domain.

## Import

Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`

```
$ terraform import alicloud_cs_edge_kubernetes.main cluster-id
```

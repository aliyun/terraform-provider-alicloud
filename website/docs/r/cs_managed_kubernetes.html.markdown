---
subcategory: "Container Service (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_managed_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-managed-kubernetes"
description: |-
  Provides a Alicloud resource to manage container managed kubernetes cluster.
---

# alicloud\_cs\_managed\_kubernetes

This resource will help you to manage a ManagedKubernetes Cluster in Alibaba Cloud Kubernetes Service. 

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

## Example Usage

```
# If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count             = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id            = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block        = element(var.vswitch_cidrs, count.index)
  availability_zone = element(var.availability_zone, count.index)
}


# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "terway_vswitches" {
  count             = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cirds)
  vpc_id            = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block        = element(var.terway_vswitch_cirds, count.index)
  availability_zone = element(var.availability_zone, count.index)
}

resource "alicloud_cs_managed_kubernetes" "k8s" {
  count                 = var.k8s_number
  # version can not be defined in variables.tf. Options: 1.18.8-aliyun.1|1.16.9-aliyun.1
  version               = "1.16.9-aliyun.1"
  worker_vswitch_ids    = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)): length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids       = length(var.terway_vswitch_ids) > 0 ? split(",", join(",", var.terway_vswitch_ids)): length(var.terway_vswitch_cirds) < 1 ? [] : split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  worker_instance_types = var.worker_instance_types
  worker_number         = var.worker_number
  node_cidr_mask        = var.node_cidr_mask
  enable_ssh            = var.enable_ssh
  install_cloud_monitor = var.install_cloud_monitor
  cpu_policy            = var.cpu_policy
  proxy_mode            = var.proxy_mode
  password              = var.password
  service_cidr          = var.service_cidr

  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name                    = lookup(addons.value, "name", var.cluster_addons)
        config                  = lookup(addons.value, "config", var.cluster_addons)
      }
  }
  runtime = {
    name    = "docker"
    version = "19.03.5"
  }
}

```

## Argument Reference

The following arguments are supported:

### Global params

* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `timezone` - (Optional, ForceNew, Available in 1.103.2+) When you create a cluster, set the time zones for the Master and Woker nodes. You can only change the managed node time zone if you create a cluster. Once the cluster is created, you can only change the time zone of the Worker node.
* `resource_group_id` - (Optional, ForceNew, Available in 1.101.0+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `version` - (Optional, Available since 1.70.1) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `runtime` - (Optional, Available in 1.103.2+) The runtime of containers. Default to `docker`. If you select another container runtime, see [How do I select between Docker and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm?spm=a2c63.p38356.b99.440.22563866AJkBgI). Detailed below.
* `enable_ssh` - (Optional) Enable login to the node through SSH. Default to `false`.
* `rds_instance` - (Optional, Available in 1.103.2+) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, Available in 1.91.0+) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).
* `tags` - (Optional, Available in 1.97.0+) Default nil, A map of tags assigned to the kubernetes cluster .
* `proxy_mode` - (Optional) Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `image_id` - (Optional) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `cluster_domain` - (Optional, ForceNew, Available in 1.103.2+) Cluster local domain name, Default to `cluster.local`. A domain name consists of one or more sections separated by a decimal point (.), each of which is up to 63 characters long, and can be lowercase, numerals, and underscores (-), and must be lowercase or numerals at the beginning and end.
* `custom_san` - (Optional, ForceNew, Available in 1.103.2+) Customize the certificate SAN, multiple IP or domain names are separated by English commas (,).
* `user_ca` - (Optional, ForceNew) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `deletion_protection` - (Optional, Available in 1.103.2+)  Whether to enable cluster deletion protection.
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. Default to `true`.
* `exclude_autoscaler_nodes` - (Optional, Available in 1.88.0+) Exclude autoscaler nodes from `worker_nodes`. Default to `false`.
* `service_account_issuer` - (Optional, ForceNew, Available in 1.92.0+) The issuer of the Service Account token for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm), corresponds to the `iss` field in the token payload. Set this to `"kubernetes.default.svc"` to enable the Token Volume Projection feature (requires specifying `api_audiences` as well).
* `api_audiences` - (Optional, ForceNew, Available in 1.92.0+) A list of API audiences for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm). Set this to `["kubernetes.default.svc"]` if you want to enable the Token Volume Projection feature (requires specifying `service_account_issuer` as well.
* `tags` - (Optional, Available in 1.97.0+) Default nil, A map of tags assigned to the kubernetes cluster . Detailed below.
* `encryption_provider_key` - (Optional, ForceNew, Available in 1.103.2+) The disk encryption key.

##### runtime

The following example is the definition of runtime block:

```
  runtime = {
    name = "docker"
    version = "19.03.5"
  }
```

##### tags

The following example is the definition of tags block. The type of this field is map:

```
  # for example, define three tags

  tags = {
    "key1" = "value1"
    "key2" = "value2"
    "name" = "tf"
  }
```

### Network

* `pod_cidr` - (**Required**) - [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `pod_vswitch_ids` - (**Required**) - [Terway Specific] The vswitches for the pod network when using Terway.Be careful the `pod_vswitch_ids` can not equal to `worker_vswtich_ids` or `master_vswtich_ids` but must be in same availability zones.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

-> **NOTE:** If you want to use `Terway` as CNI network plugin, You need to specific the `pod_vswitch_ids` field and addons with `terway-eniip`.
If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

### Worker params

* `worker_number` - (**Required**) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswtich_ids` - (**Required**) The vswitches used by workers.
* `worker_instance_types` - (**Required**, ForceNew) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `password` - (**Required**, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (**Required**) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encrypted_password` - (**Required**, Available in 1.57.1+) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

* `worker_instance_charge_type` - (Optional, Force new resource) Worker payment type, its valid value is either or `PostPaid` or `PrePaid`. Defaults to `PostPaid`. If value is `PrePaid`, the files `worker_period`, `worker_period_unit`, `worker_auto_renew` and `worker_auto_renew_period` are required.
* `worker_period` - (Optional) Worker payment period. The unit is `Month`. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `worker_period_unit` - (Optional) Worker payment period unit, the valid value is `Month`.
* `worker_auto_renew` - (Optional) Enable worker payment auto-renew, defaults to false.
* `worker_auto_renew_period` - (Optional) Worker payment auto-renew period,, it can be one of {1, 2, 3, 6, 12}.
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud`, `cloud_ssd`, `cloud_essd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.
* `worker_data_disks` - (Optional, Available in 1.91.0+) The data disk configurations of worker nodes, such as the disk type and disk size.
  * `category`: the type of the data disks. Valid values:
    * cloud: basic disks.
    * cloud_efficiency: ultra disks.
    * cloud_ssd: SSDs.
    * cloud_essd: essd.
  * `size`: the size of a data disk. Unit: GiB.
  * `encrypted`: specifies whether to encrypt data disks. Valid values: true and false.
* `node_name_mode` - (Optional, Available in 1.88.0+) Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be `aliyun.com00055test`.
* `node_port_range`- (Optional, ForceNew, Available in 1.103.2+) The service port range of nodes, valid values: `30000` to `65535`. Default to `30000-32767`.
* `os_type` - (Optional, ForceNew, Available in 1.103.2+) The operating system of the nodes that run pods, its valid value is either `Linux` or `Windows`. Default to `Linux`.
* `platform` - (Optional, ForceNew, Available in 1.103.2+) The architecture of the nodes that run pods, its valid value is either `CentOS` or `AliyunLinux`. Default to `CentOS`.
* `cpu_policy` - (Optional) Kubelet cpu policy. For Kubernetes 1.12.6 and later, its valid value is either `static` or `none`. Default to `none`.
* `user_data` - (Optional, Available in 1.81.0+) Custom data that can execute on nodes. For more information, see [Prepare user data](https://www.alibabacloud.com/help/doc-detail/49121.htm).
* `taints` - (Optional, Available in 1.103.2+) Taints ensure pods are not scheduled onto inappropriate nodes. One or more taints are applied to a node; this marks that the node should not accept any pods that do not tolerate the taints. For more information, see [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/). Detailed below.

##### worker_vswtich_ids

The following example is the definition of `worker_vswtich_ids` block.

```
  # the ID can be the same, At least one.

  worker_vswtich_ids = ["vsw-id1", "vsw-id1", "vsw-id2"]
```

##### taints

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

### Addons

It is a new field since 1.75.0. You can specific network plugin,log component,ingress component and so on.

You can get more information about addons on ACK web console. When you create a ACK cluster. You can get openapi-spec before creating the cluster on submission page.

* `logtail-ds` - You can specific `IngressDashboardEnabled` and `sls_project_name` in config. If you switch on `IngressDashboardEnabled` and `sls_project_name`,then logtail-ds would use `sls_project_name` as default log store.
* `nginx-ingress-controller` - You can specific `IngressSlbNetworkType` in config. Options: internet|intranet.

The `main.tf`:

```terraform
resource "alicloud_cs_managed_kubernetes" "k8s" {
  # ... other configuration ...

  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name          = lookup(addons.value, "name", var.cluster_addons)
        config        = lookup(addons.value, "config", var.cluster_addons)
        disabled      = lookup(addons.value, "disabled", var.cluster_addons)
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
}
```

### Computed params (No need to configure)

* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `availability_zone` - (Optional) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.

### Removed params (Never Supported)

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
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.

### Block Nodes

The following arguments are supported in the `worker_nodes` configuration block:

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.
* `role` - (Deprecated from version 1.9.4)

### Block Connections

The following arguments are supported in the `connections` configuration block:

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.
* `service_domain` - Service Access Domain.

## Import

Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`.

```
  $ terraform import alicloud_cs_managed_kubernetes.main cluster-id
```

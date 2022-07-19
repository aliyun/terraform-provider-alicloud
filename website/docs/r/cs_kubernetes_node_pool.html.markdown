---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_node_pool"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-node-pool"
description: |-
  Provides a Alicloud resource to manage container kubernetes node pool.
---

# alicloud\_cs\_kubernetes\_node\_pool

This resource will help you to manage node pool in Kubernetes Cluster. 

-> **NOTE:** Available in 1.97.0+.

-> **NOTE:** From version 1.109.1, support managed node pools, but only for the professional managed clusters.

-> **NOTE:** From version 1.109.1, support remove node pool nodes.

-> **NOTE:** From version 1.111.0, support auto scaling node pool. For more information on how to use auto scaling node pools, see [Use Terraform to create an elastic node pool](https://help.aliyun.com/document_detail/197717.htm). With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.

-> **NOTE:** ACK adds a new RamRole (AliyunCSManagedAutoScalerRole) for the permission control of the node pool with auto-scaling enabled. If you are using a node pool with auto scaling, please click [AliyunCSManagedAutoScalerRole](https://ram.console.aliyun.com/role/authorization?request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedAutoScalerRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedAutoScalerRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization. 

-> **NOTE:** ACK adds a new RamRole（AliyunCSManagedNlcRole） for the permission control of the management node pool. If you use the management node pool, please click [AliyunCSManagedNlcRole](https://ram.console.aliyun.com/role/authorization?spm=5176.2020520152.0.0.387f16ddEOZxMv&request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedNlcRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedNlcRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization.

-> **NOTE:** From version 1.123.1, supports the creation of a node pool of spot instance.

-> **NOTE:** It is recommended to create a cluster with zero worker nodes, and then use a node pool to manage the cluster nodes. 

-> **NOTE:** From version 1.127.0, support for adding existing nodes to the node pool. In order to distinguish automatically created nodes, it is recommended that existing nodes be placed separately in a node pool for management. 

-> **NOTE:** From version 1.149.0, support for specifying deploymentSet for node pools. 

-> **NOTE:** From version 1.158.0, Support for specifying the desired size of nodes for the node pool, for more information, visit [Modify the expected number of nodes in a node pool](https://www.alibabacloud.com/help/en/doc-detail/160490.html#title-mpp-3jj-oo3)

-> **NOTE:** From version 1.166.0, Support configuring system disk encryption.

-> **NOTE:** From version 1.177.0+, Support `kms_encryption_context`, `rds_instances`, `system_disk_snapshot_policy_id` and `cpu_policy`, add spot strategy `SpotAsPriceGo` and `NoSpot`.

## Example Usage

The managed cluster configuration,

```terraform
variable "name" {
  default = "tf-test"
}
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.1.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
}
```

Create a node pool.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alicloud_vswitch.default.id]
  instance_types = [data.alicloud_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be 0
  desired_size = 1
}
```

The parameter `node_count` are deprecated from version 1.158.0，but it can still works. If you want to use the new parameter `desired_size` instead, you can update it as follows. for more information of `desired_size`, visit [Modify the expected number of nodes in a node pool](https://www.alibabacloud.com/help/en/doc-detail/160490.html#title-mpp-3jj-oo3). 

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alicloud_vswitch.default.id]
  instance_types = [data.alicloud_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # comment out node_count and specify a new field desired_size
  # node_count = 1

  desired_size = 1
}
```

Create a managed node pool. If you need to enable maintenance window, you need to set the maintenance window in `alicloud_cs_managed_kubernetes`.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40

  # only key_name is supported in the management node pool
  key_name = alicloud_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be zero
  desired_size = 1

  # management node pool configuration.
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }

}
```

Enable automatic scaling for the node pool. `scaling_config` is required.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # automatic scaling node pool configuration.
  # With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
  scaling_config {
    min_size = 1
    max_size = 10
  }

}
```

Enable automatic scaling for managed node pool.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name
  # management node pool configuration.
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }
  # enable auto-scaling
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
  # Rely on auto-scaling configuration, please create auto-scaling configuration through alicloud_cs_autoscaling_config first.
  depends_on = [alicloud_cs_autoscaling_config.default]
}
```

Create a `PrePaid` node pool.
```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name
  # use PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # open cloud monitor
  install_cloud_monitor = true

  # enable auto-scaling
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
}
```

Create a node pool with spot instance.
```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alicloud_vswitch.default.id]
  instance_types = [data.alicloud_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be 0
  desired_size = 1

  # spot config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.default.instance_types.0.id
    # Different instance types have different price caps
    price_limit = "0.70"
  }
}
```

Use Spot instances to create a node pool with auto-scaling enabled 
```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # spot price config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.default.instance_types.0.id
    price_limit   = "0.70"
  }
}
```

Create a node pool with platform as Windows 
```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "windows-np"
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"
  desired_size         = 1

  // if the instance platform is windows, the password is requered.
  password = "Hello1234"
  platform = "Windows"
  image_id = "${window_image_id}"
}
```

Add an existing node to the node pool

In order to distinguish automatically created nodes, it is recommended that existing nodes be placed separately in a node pool for management. 

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "existing-node"
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"

  # add existing node to nodepool
  instances = ["instance_id_01", "instance_id_02", "instance_id_03"]
  # default is false
  format_disk = false
  # default is true
  keep_instance_name = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `name` - (Required) The name of node pool.
* `vswitch_ids` - (Required) The vswitches used by node pool workers.
* `instance_types` (Required) The instance type of worker node.
* `password` - (Optional, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Optional) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. Only `key_name` is supported in the management node pool.
* `kms_encrypted_password` - (Optional, Available in 1.177.0) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, Available in 1.177.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `node_count` (Optional, Deprecated) The worker node number of the node pool. From version 1.111.0, `node_count` is not required.
* `desired_size` (Optional, Available in 1.158.0+) The desired size of nodes of the node pool. From version 1.158.0, `desired_size` is not required.
* `system_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd`, `cloud_efficiency` and `cloud_essd`. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) The system disk category of worker node. Its valid value range [40~500] in GB. Default to `120`.
* `system_disk_performance_level` - (Optional) The performance of system disk, only valid for ESSD disk. You have to specify one of `PL0` `PL1` `PL2` `PL3` fields.
* `data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size. 
  * `category` - The type of the data disks. Valid values:`cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`.
  * `size` - The size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
  * `encrypted` - Specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
  * `performance_level` - (Optional, Available in 1.120.0+) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
  * `kms_key_id` - (Optional, Available in 1.97.0+) The kms key id used to encrypt the data disk. It takes effect when `encrypted` is true.
* `security_group_id` - (Optional, Deprecated) The security group id for worker node. Field `security_group_id` has been deprecated from provider version 1.145.0. New field `security_group_ids` instead.
* `platform` - (Optional, Deprecated, Available in 1.127.0+) The platform. One of `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`. If you select `Windows` or `WindowsCore`, the `passord` is required. Field `platform` has been deprecated from provider version 1.145.0. New field `image_type` instead.
* `image_id` - (Optional) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `node_name_mode` - (Optional) Each node name consists of a prefix, an IP substring, and a suffix. For example "customized,aliyun.com,5,test", if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `tags` - (Optional) A Map of tags to assign to the resource. It will be applied for ECS instances finally.
* `labels` - (Optional) A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument.
  * `key` - The label key.
  * `value` - The label value.
* `taints` - (Optional) A List of Kubernetes taints to assign to the nodes.
* `management` - (Optional, Available in 1.109.1+) Managed node pool configuration. When using a managed node pool, the node key must use `key_name`. Detailed below.
* `scaling_policy` - (Optional, Available in 1.127.0+) The scaling mode. Valid values: `release`, `recycle`, default is `release`. Standard mode(release): Create and release ECS instances based on requests.Swift mode(recycle): Create, stop, and restart ECS instances based on needs. New ECS instances are only created when no stopped ECS instance is avalible. This mode further accelerates the scaling process. Apart from ECS instances that use local storage, when an ECS instance is stopped, you are only chatged for storage space.
* `scaling_config` - (Optional, Available in 1.111.0+) Auto scaling node pool configuration. For more details, see `scaling_config`. With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
* `instance_charge_type`- (Optional, Available in 1.119.0+) Node payment type. Valid values: `PostPaid`, `PrePaid`, default is `PostPaid`. If value is `PrePaid`, the arguments `period`, `period_unit`, `auto_renew` and `auto_renew_period` are required.
* `period`- (Optional, Available in 1.119.0+) Node payment period. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `period_unit`- (Optional, Available in 1.119.0+) Node payment period unit, valid value: `Month`. Default is `Month`.
* `auto_renew`- (Optional, Available in 1.119.0+) Enable Node payment auto-renew, default is `false`.
* `auto_renew_period`- (Optional, Available in 1.119.0+) Node payment auto-renew period, one of `1`, `2`, `3`,`6`, `12`.
* `install_cloud_monitor`- (Optional, Available in 1.119.0+) Install the cloud monitoring plug-in on the node, and you can view the monitoring information of the instance through the cloud monitoring console. Default is `true`.
* `unschedulable`- (Optional, Available in 1.119.0+) Set the newly added node as unschedulable. If you want to open the scheduling option, you can open it in the node list of the console. If you are using an auto-scaling node pool, the setting will not take effect. Default is `false`.
* `resource_group_id` - (Optional, ForceNew, Available in 1.123.1+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `internet_charge_type` - (Optional, Available in 1.123.1+) The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`. Conflict with `eip_internet_charge_type`, EIP and public network IP can only choose one. 
* `internet_max_bandwidth_out` - (Optional, Available in 1.123.1+) The maximum outbound bandwidth for the public network. Unit: Mbit/s. Valid values: 0 to 100.
* `spot_strategy` - (Optional, Available in 1.123.1+) The preemption policy for the pay-as-you-go instance. This parameter takes effect only when `instance_charge_type` is set to `PostPaid`. Valid value `SpotWithPriceLimit`,`SpotAsPriceGo` and `NoSpot`.
* `spot_price_limit` - (Optional, Available in 1.123.1+) The maximum hourly price of the instance. This parameter takes effect only when `spot_strategy` is set to `SpotWithPriceLimit`. You could enable multiple spot instances by setting this field repeatedly.
  * `instance_type` - (Optional, Available in 1.123.1+) Spot instance type.
  * `price_limit` - (Optional, Available in 1.123.1+) The maximum hourly price of the spot instance. A maximum of three decimal places are allowed.
* `instances` - (Optional, Available in 1.127.0+) The instance list. Add existing nodes under the same cluster VPC to the node pool. 
* `keep_instance_name` - (Optional, Available in 1.127.0+) Add an existing instance to the node pool, whether to keep the original instance name. It is recommended to set to `true`.
* `format_disk` - (Optional, Available in 1.127.0+) After you select this check box, if data disks have been attached to the specified ECS instances and the file system of the last data disk is uninitialized, the system automatically formats the last data disk to ext4 and mounts the data disk to /var/lib/docker and /var/lib/kubelet. The original data on the disk will be cleared. Make sure that you back up data in advance. If no data disk is mounted on the ECS instance, no new data disk will be purchased. Default is `false`.
* `security_group_ids` - (Optional, Available in 1.145.0+) Multiple security groups can be configured for a node pool. If both `security_group_ids` and `security_group_id` are configured, `security_group_ids` takes effect. This field cannot be modified. 
* `image_type` - (Optional, Available in 1.145.0+) The image type, instead of `platform`. This field cannot be modified. One of `AliyunLinux`, `AliyunLinux3`, `AliyunLinux3Arm64`, `AliyunLinuxUEFI`, `CentOS`, `Windows`,`WindowsCore`,`AliyunLinux Qboot`,`ContainerOS`. If you select `Windows` or `WindowsCore`, the `passord` is required.
* `runtime_name` - (Optional, ForceNew, Available in 1.145.0+) The runtime name of containers. If not set, the cluster runtime will be used as the node pool runtime. If you select another container runtime, see [Comparison of Docker, containerd, and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm).
* `runtime_version` - (Optional, ForceNew, Available in 1.145.0+) The runtime version of containers. If not set, the cluster runtime will be used as the node pool runtime.
* `deployment_set_id` - (Optional, ForceNew, Available in 1.149.0+) The deployment set of node pool. Specify the deploymentSet to ensure that the nodes in the node pool can be distributed on different physical machines.
* `system_disk_encrypted` - (Optional, Available in 1.166.0+) Whether to enable system disk encryption.
* `system_disk_kms_key` - (Optional, Available in 1.166.0+) The kms key id used to encrypt the system disk. It takes effect when system_disk_encrypted is true.
* `system_disk_encrypt_algorithm` - (Optional, Available in 1.166.0+) The encryption Algorithm for Encrypting System Disk. It takes effect when system_disk_encrypted is true. Valid values `aes-256` and `sm4-128`.
* `cis_enabled` - (Optional, Available in 1.173.0+) Whether enable worker node to support cis security reinforcement, its valid value `true` or `false`. Default to `false` and apply to `image_type/platform=AliyunLinux`, see [CIS Reinforcement](https://help.aliyun.com/document_detail/223744.html).
* `soc_enabled` - (Optional, Available in 1.173.0+) Whether enable worker node to support soc security reinforcement, its valid value `true` or `false`. Default to `false` and apply to `image_type/platform=AliyunLinux`, see [SOC Reinforcement](https://help.aliyun.com/document_detail/196148.html).  
  -> **NOTE:** It is forbidden to set both `cis_enabled` and `soc_enabled` to `true`at the same time.
* `rds_instances` - (Optional, Available in 1.177.0+) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `system_disk_snapshot_policy_id` - (Optional, Available in 1.177.0+) The system disk snapshot policy id.
* `cpu_policy` - (Optional, Available in 1.177.0+) Kubelet cpu policy. For Kubernetes 1.12.6 and later, its valid value is either `static` or `none`. Default to `none` and modification is not supported.

#### tags

The tags example：
```
tags = {
  "key-a" = "value-a"
  "key-b" = "value-b"
  "env"   = "prod"
}
```

#### management

The following arguments are supported in the `management` configuration block:

* `auto_repair` - (Optional) Whether automatic repair, Default to `false`.
* `auto_upgrade`- (Optional) Whether auto upgrade, Default to `false`.
* `surge` - (Optional) Number of additional nodes. You have to specify one of surge, surge_percentage.
* `surge_percentage` - (Optional) Proportion of additional nodes. You have to specify one of surge, surge_percentage.
* `max_unavailable` - (Required) Max number of unavailable nodes. Default to `1`.

#### scaling_config

The following arguments are supported in the `scaling_config` configuration block:

* `min_size` - (Required, Available in 1.111.0+) Min number of instances in a auto scaling group, its valid value range [0~1000].
* `max_size` - (Required, Available in 1.111.0+) Max number of instances in a auto scaling group, its valid value range [0~1000]. `max_size` has to be greater than `min_size`.
* `type` - (Optional, Available in 1.111.0+) Instance classification, not required. Vaild value: `cpu`, `gpu`, `gpushare` and `spot`. Default: `cpu`. The actual instance type is determined by `instance_types`.
* `is_bond_eip` - (Optional, Available in 1.111.0+) Whether to bind EIP for an instance. Default: `false`.
* `eip_internet_charge_type` - (Optional, Available in 1.111.0+) EIP billing type. `PayByBandwidth`: Charged at fixed bandwidth. `PayByTraffic`: Billed as used traffic. Default: `PayByBandwidth`. Conflict with `internet_charge_type`, EIP and public network IP can only choose one. 
* `eip_bandwidth` - (Optional, Available in 1.111.0+) Peak EIP bandwidth. Its valid value range [1~500] in Mbps. Default to `5`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the node pool, format cluster_id:nodepool_id.
* `cluster_id` - The cluster id.
* `name` - The name of the nodepool.
* `vswitch_ids` - The vswitches used by node pool workers.
* `vpc_id` - The VPC of the nodes in the node pool.
* `image_id` - The image used by node pool workers.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `scaling_group_id` - (Available in 1.105.0+) Id of the Scaling Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating node-pool in the kubernetes cluster (until it reaches the initial `active` status). 
* `update` - (Defaults to 60 mins) Used when activating the node-pool in the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when deleting node-pool in kubernetes cluster. 

## Import

Cluster nodepool can be imported using the id, e.g. Then complete the nodepool.tf accords to the result of `terraform plan`.

```
  $ terraform import alicloud_cs_kubernetes_node_pool.custom_nodepool cluster_id:nodepool_id
```

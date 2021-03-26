---
subcategory: "Container Service for Kubernetes (CSK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_node_pool"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-node-pool"
description: |-
  Provides a Alicloud resource to manage container kubernetes node pool.
---

# alicloud\_cs\_kubernetes\_node\_pool

This resource will help you to manager node pool in Kubernetes Cluster. 

-> **NOTE:** Available in 1.97.0+.

-> **NOTE:** From version 1.109.1, support managed node pools, but only for the professional managed clusters.

-> **NOTE:** From version 1.109.1, support remove node pool nodes.

-> **NOTE:** From version 1.111.0, support auto scaling node pool. For more information on how to use auto scaling node pools, see [Use Terraform to create an elastic node pool](https://help.aliyun.com/document_detail/197717.htm).

-> **NOTE:** ACK adds a new RamRole (AliyunCSManagedAutoScalerRole) for the permission control of the auto-scaling node pool. If you use the auto-scaling node pool, please click [AliyunCSManagedAutoScalerRole](https://ram.console.aliyun.com/role/authorization?request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedAutoScalerRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedAutoScalerRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization. 

-> **NOTE:** ACK adds a new RamRole（AliyunCSManagedNlcRole） for the permission control of the management node pool. If you use the management node pool, please click [AliyunCSManagedNlcRole](https://ram.console.aliyun.com/role/authorization?spm=5176.2020520152.0.0.387f16ddEOZxMv&request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedNlcRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedNlcRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization. 
## Example Usage

The managed cluster configuration,

```terraform
variable "name" {
  default = "tf-test"
}
data "alicloud_zones" default {
  available_resource_creation  = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone            = data.alicloud_zones.default.zones.0.id
  cpu_core_count               = 2
  memory_size                  = 4
  kubernetes_node_role         = "Worker"
}
resource "alicloud_vpc" "default" {
  name                         = var.name
  cidr_block                   = "10.1.0.0/21"
}
resource "alicloud_vswitch" "default" {
  vswitch_name                 = var.name
  vpc_id                       = alicloud_vpc.default.id
  cidr_block                   = "10.1.1.0/24"
  availability_zone            = "data.alicloud_zones.default.zones.0.id
}
resource "alicloud_key_pair" "default" {
  key_name                     = var.name
}
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}
```

Custom node pool in kubernetes cluster.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                         = var.name
  cluster_id                   = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids                  = [alicloud_vswitch.default.id]
  instance_types               = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category         = "cloud_efficiency"
  system_disk_size             = 40
  key_name                     = alicloud_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be 0
  node_count                   = 1
}
```

Management node pool in kubernetes cluster. If you need to enable maintenance window, you need to set the maintenance window in `alicloud_cs_managed_kubernetes`.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                         = var.name
  cluster_id                   = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids                  = [alicloud_vswitch.default.id]
  instance_types               = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category         = "cloud_efficiency"
  system_disk_size             = 40
  
  # only key_name is supported in the management node pool
  key_name                     = alicloud_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be zero
  node_count                   = 1

  # management node pool configuration.
  management {
    auto_repair      = true
    auto_upgrade     = true
    surge            = 1
    max_unavailable  = 1
  }

}
```

Automatic scaling node pool in kubernetes cluster. `node_count` is not required when using an automatic scaling node pool.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                         = var.name
  cluster_id                   = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids                  = [alicloud_vswitch.default.id]
  instance_types               = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category         = "cloud_efficiency"
  system_disk_size             = 40
  key_name                     = alicloud_key_pair.default.key_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size      = 1
    max_size      = 10
  }

}
```

Enables auto-scaling of the managed node pool in kubernetes cluster.

```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                         = var.name
  cluster_id                   = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids                  = [alicloud_vswitch.default.id]
  instance_types               = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category         = "cloud_efficiency"
  system_disk_size             = 40
  key_name                     = alicloud_key_pair.default.key_name
  # management node pool configuration.
  management {
    auto_repair      = true
    auto_upgrade     = true
    surge            = 1
    max_unavailable  = 1
  }
  # enable auto-scaling
  scaling_config {
    min_size         = 1
    max_size         = 10
    type             = "cpu"
  }
}
```
Create a `PrePaid` node pool.
```terraform
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                         = var.name
  cluster_id                   = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids                  = [alicloud_vswitch.default.id]
  instance_types               = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category         = "cloud_efficiency"
  system_disk_size             = 40
  key_name                     = alicloud_key_pair.default.key_name
  # use PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # open cloud monitor
  install_cloud_monitor        = true
  
  # enable auto-scaling
  scaling_config {
    min_size         = 1
    max_size         = 10
    type             = "cpu"
  }
}
```


## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `name` - (Required) The name of node pool.
* `vswitch_ids` - (Required) The vswitches used by node pool workers.
* `instance_types` (Required) The instance type of worker node.
* `node_count` (Optional) The worker node number of the node pool. From version 1.111.0, `node_count` is not required.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Required) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. Only `key_name` is supported in the management node pool.
* `kms_encrypted_password` - (Required) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `system_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) The system disk category of worker node. Its valid value range [40~500] in GB. Default to `120`.
* `data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size. 
  * `category` - the type of the data disks. Valid values:`cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`.
  * `size` - the size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
  * `encrypted` - specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
  * `performance_level` - (Optional, Available in 1.119.1+) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `security_group_id` - (Optional) The system disk size of worker node. 
* `image_id` - (Optional) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `node_name_mode` - (Optional) Each node name consists of a prefix, an IP substring, and a suffix. For example "customized,aliyun.com,5,test", if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `tags` - (Optional) A map of tags to assign to the resource. It will be applied for ECS instances finally.
* `labels` - (Optional) A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument.
  * `key` - The label key.
  * `value` - The label value.
* `taints` - (Optional) A List of Kubernetes taints to assign to the nodes.
* `management` - (Optional, Available in 1.109.1+) Managed node pool configuration. When using a managed node pool, the node key must use `key_name`. Detailed below.
* `scaling_config` - (Optional, Available in 1.111.0+) Auto scaling node pool configuration. For more details, see `scaling_config`.
* `instance_charge_type`- (Optional, Available in 1.119.0+) Node payment type. Valid values: `PostPaid`, `PrePaid`, default is `PostPaid`. If value is `PrePaid`, the arguments `period`, `period_unit`, `auto_renew` and `auto_renew_period` are required.
* `period`- (Optional, Available in 1.119.0+) Node payment period. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `period_unit`- (Optional, Available in 1.119.0+) Node payment period unit, valid value: `Month`. Default is `Month`.
* `auto_renew`- (Optional, Available in 1.119.0+) Enable Node payment auto-renew, default is `false`.
* `auto_renew_period`- (Optional, Available in 1.119.0+) Node payment auto-renew period, one of `1`, `2`, `3`,`6`, `12`.
* `install_cloud_monitor`- (Optional, Available in 1.119.0+) Install the cloud monitoring plug-in on the node, and you can view the monitoring information of the instance through the cloud monitoring console. Default is `true`.
* `unschedulable`- (Optional, Available in 1.119.0+) Set the newly added node as unschedulable. If you want to open the scheduling option, you can open it in the node list of the console. If you are using an auto-scaling node pool, the setting will not take effect. Default is `false`.

#### tags
The tags example：
```
tags {
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
* `eip_internet_charge_type` - (Optional, Available in 1.111.0+) EIP billing type. `PayByBandwidth`: Charged at fixed bandwidth. `PayByTraffic`: Billed as used traffic. Default: `PayByBandwidth`.
* `eip_bandwidth` - (Optional, Available in 1.111.0+) Peak EIP bandwidth. Its valid value range [1~500] in Mbps. Default to `5`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the node pool, format cluster_id:nodepool_id.
* `cluster_id` - The cluster id.
* `name` - The name of the nodepool.
* `vswitch_ids` - The vswitches used by node pool workers.
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
  $ terraform import alicloud_cs_node_pool.custom_nodepool cluster_id:nodepool_id
```
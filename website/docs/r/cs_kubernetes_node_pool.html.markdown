---
subcategory: "Container Service (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_node_pool"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-node-pool"
description: |-
  Provides a Alicloud resource to manage container kubernetes node pool.
---

# alicloud\_cs\_kubernetes\_node\_pool

This resource will help you to manager node pool in Kubernetes Cluster. 

-> **NOTE:** Available in 1.97.0+.

-> **NOTE:** From version 1.109.0, support managed node pools, but only for the professional managed clusters.

-> **NOTE:** From version 1.109.0, support remove node pool nodes.

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
  name                         = var.name
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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `name` - (Required) The name of node pool.
* `vswitch_ids` - (Required) The vswitches used by node pool workers.
* `instance_types` (Required) The instance type of worker node.
* `node_count` (Required) The worker node number of the node pool.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Required) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. Only `key_name` is supported in the management node pool.
* `kms_encrypted_password` - (Required) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `system_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) The system disk category of worker node. Its valid value range [40~500] in GB. Default to `120`.
* `worker_data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size. 
  * category: the type of the data disks. Valid values:`cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`.
  * size: the size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
  * encrypted: specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
* `security_group_id` - (Optional) The system disk size of worker node. 
* `image_id` - (Optional) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `node_name_mode` - (Optional) Each node name consists of a prefix, an IP substring, and a suffix. For example "customized,aliyun.com,5,test", if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `tags` - (Optional) A List of tags to assign to the resource. It will be applied for ECS instances finally.
  * key: It can be up to 64 characters in length. It cannot begin with "aliyun", "http://", or "https://". It cannot be a null string.
  * value: It can be up to 128 characters in length. It cannot begin with "aliyun", "http://", or "https://" It can be a null string.
* `labels` - (Optional) A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument.
* `taints` - (Optional) A List of Kubernetes taints to assign to the nodes.
* `management` - (Optional, Available in 1.109.0+) Managed node pool configuration. When using a managed node pool, the node key must use `key_name`. Detailed below.

#### management

The following arguments are supported in the `management` configuration block:

* `auto_repair` - (Optional) Whether automatic repair, Default to `false`.
* `auto_upgrade`- (Optional) Whether auto upgrade, Default to `false`.
* `surge` - (Optional) Number of additional nodes. You have to specify one of surge, surge_percentage.
* `surge_percentage` - (Optional) Proportion of additional nodes. You have to specify one of surge, surge_percentage.
* `max_unavailable` - (Required) Max number of unavailable nodes. Default to `1`.

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
* `scaling_group_id` - (Available in 1.105.0+) Id of the Scaling Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating node-pool in the kubernetes cluster (until it reaches the initial `active` status). 
* `update` - (Defaults to 60 mins) Used when activating the node-pool in the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when deleting node-pool in kubernetes cluster. 

---
subcategory: "Container Service (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_swarm"
sidebar_current: "docs-alicloud-resource-cs-swarm"
description: |-
  Provides a Alicloud container swarm cluster resource.
---

# alicloud\_cs\_swarm

-> **DEPRECATED:** This resource manages swarm cluster, which is being deprecated and will be replaced by Kubernetes cluster.

This resource will help you to manager a Swarm Cluster.

-> **NOTE:** Swarm cluster only supports VPC network and you can specify a VPC network by filed `vswitch_id`.

## Example Usage

Basic Usage

```
resource "alicloud_cs_swarm" "my_cluster" {
  password      = "Yourpassword1234"
  instance_type = "ecs.n4.small"
  name          = "ClusterFromAlicloud"
  node_number   = 2
  disk_category = "cloud_efficiency"
  disk_size     = 20
  cidr_block    = "172.18.0.0/24"
  image_id      = "${var.image_id}"
  vswitch_id    = "${var.vswitch_id}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - The container cluster's name. It is the only in one Alicloud account.
* `name_prefix` - The container cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to 'Terraform-Creation'.
* `size` - Field 'size' has been deprecated from provider version 1.9.1. New field 'node_number' replaces it.
* `node_number` - The ECS node number of the container cluster. Its value choices are 1~50, and default to 1.
* `cidr_block` - (Required, ForceNew) The CIDR block for the Container. It can not be same as the CIDR used by the VPC.
  Valid value:
    - 192.168.0.0/16
    - 172.19-30.0.0/16
    - 10.0.0.0/16

  System reserved private network address: 172.16/17/18/31.0.0/16.
  Maximum number of hosts allowed in the cluster: 256.

* `image_id` - (ForceNew) The image ID of ECS instance node used. Default to System automate allocated.
* `instance_type` - (Required, ForceNew) The type of ECS instance node.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `password` - (Required, ForceNew, Sensitive) The password of ECS instance node.
* `disk_category` - (ForceNew) The data disk category of ECS instance node. Its valid value are `cloud`, `cloud_ssd`, `cloud_essd`, `ephemeral_essd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `disk_size` - (ForceNew) The data disk size of ECS instance node. Its valid value is 20~32768 GB. Default to 20.
* `vswitch_id` - (Required, ForceNew) The password of ECS instance node. If it is not specified, the container cluster's network mode will be `Classic`.
* `release_eip` - Whether to release EIP after creating swarm cluster successfully. Default to false.
* `need_slb`- (ForceNew) Whether to create the default simple routing Server Load Balancer instance for the cluster. The default value is true.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `size` - It has been deprecated from provider version 1.9.1. New field 'node_number' replaces it.
* `node_number` - The node number.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `vswitch_id` - The ID of VSwitch where the current cluster is located.
* `slb_id` - The ID of load balancer where the current cluster worker node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `agent_version` - The nodes agent version.
* `instance_type` - The instance type of nodes.
* `disk_category` - The data disk category of nodes.
* `disk_size` - The data disk size of nodes.
* `nodes` - List of cluster nodes. It contains several attributes to `Block Nodes`.

### Block Nodes

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.
* `eip` - The Elastic IP address of node.
* `status` - The node current status. It is different with instance status.

## Import

Swarm cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_swarm.foo cf123456789
```

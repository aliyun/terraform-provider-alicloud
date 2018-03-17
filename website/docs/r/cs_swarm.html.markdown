---
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_swarm"
sidebar_current: "docs-alicloud-resource-cs-swarm"
description: |-
  Provides a Alicloud container swarm cluster resource.
---

# alicloud\_container\_cluster

This resource will help you to manager a Swarm Cluster.

~> **NOTE:** Swarm cluster only supports VPC network and you can specify a VPC network by filed `vswitch_id`.

## Example Usage

Basic Usage

```
resource "alicloud_cs_swarm" "my_cluster" {
  password = "Test12345"
  instance_type = "ecs.n4.small"
  name = "ClusterFromAlicloud"
  size = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.18.0.0/24"
  image_id = "${var.image_id}"
  vswitch_id = "${var.vswitch_id}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - The container cluster's name. It is the only in one Alicloud account.
* `name_prefix` - The container cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to 'Terraform-Creation'.
* `size` - The ECS node number of the container cluster. Its value choices are 1~50, and default to 1.
* `cidr_block` - (Required, Force new resource) The CIDR block for the Container. It can not be same as the CIDR used by the VPC.
  Valid value:
    - 192.168.0.0/16
    - 172.19-30.0.0/16
    - 10.0.0.0/16

  System reserved private network address: 172.16/17/18/31.0.0/16.
  Maximum number of hosts allowed in the cluster: 256.

* `image_id` - (Force new resource) The image ID of ECS instance node used. Default to System automate allocated.
* `instance_type` - (Required, Force new resource) The type of ECS instance node.
* `password` - (Required, Force new resource) The password of ECS instance node.
* `disk_category` - (Force new resource) The data disk category of ECS instance node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `disk_size` - (Force new resource) The data disk size of ECS instance node. Its valid value is 20~32768 GB. Default to 20.
* `vswitch_id` - (Required, Force new resource) The password of ECS instance node. If it is not specified, the container cluster's network mode will be `Classic`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `size` The ECS instance node number in the current container cluster.
* `vpc_id` - The ID of VPC that current cluster launched.
* `vswitch_id` - The ID of VSwitch that current cluster launched.

## Import

Swarm cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_swarm.foo cf123456789
```
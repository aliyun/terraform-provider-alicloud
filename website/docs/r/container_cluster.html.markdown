---
layout: "alicloud"
page_title: "Alicloud: alicloud_container_cluster"
sidebar_current: "docs-alicloud-resource-container-cluster"
description: |-
  Provides a Alicloud container cluster resource.
---

# alicloud\_container\_cluster

Provides a container cluster resource.

## Example Usage

Basic Usage

```
resource "alicloud_container_cluster" "my_cluster" {
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

* `name` - (Force new resource) The container cluster's name. It is the only in one Alicloud account.
* `name_prefix` - (Force new resource) The container cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name.
* `size` - The ECS node number of the container cluster. Its value choices are 1~20, and default to 1.
* `cidr_block` - (Required, Force new resource) The CIDR block for the Container. Its valid value are `192.168.X.0/24` or `172.18.X.0/24` ~ `172.31.X.0/24`. And it cannot be equal to vswitch's cidr_block and sub cidr block.
* `image_id` - (Force new resource) The image ID of ECS instance node used. Default to System automate allocated.
* `instance_type` - (Required, Force new resource) The type of ECS instance node.
* `password` - (Required, Force new resource) The password of ECS instance node.
* `disk_category` - (Force new resource) The data disk category of ECS instance node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `disk_size` - (Force new resource) The data disk size of ECS instance node. Its valid value is 20~32768 GB. Default to 20.
* `vswitch_id` - (Force new resource) The password of ECS instance node. If it is not specified, the container cluster's network mode will be `Classic`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `size` The ECS instance node number in the current container cluster.
* `vpc_id` - The ID of VPC that current cluster launched.
* `vswitch_id` - The ID of VSwitch that current cluster launched.

## Import

Container cluster can be imported using the id, e.g.

```
$ terraform import alicloud_container_cluster.foo cf123456789
```
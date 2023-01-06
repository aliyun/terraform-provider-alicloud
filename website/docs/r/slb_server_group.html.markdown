---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_group"
sidebar_current: "docs-alicloud-resource-slb-server-group"
description: |-
  Provides a Load Banlancer Virtual Backend Server Group resource.
---

# alicloud\_slb\_server\_group

A virtual server group contains several ECS instances. The virtual server group can help you to define multiple listening dimension,
and to meet the personalized requirements of domain name and URL forwarding.

-> **NOTE:** One ECS instance can be added into multiple virtual server groups.

-> **NOTE:** One virtual server group can be attached with multiple listeners in one load balancer.

-> **NOTE:** One Classic and Internet load balancer, its virtual server group can add Classic and VPC ECS instances.

-> **NOTE:** One Classic and Intranet load balancer, its virtual server group can only add Classic ECS instances.

-> **NOTE:** One VPC load balancer, its virtual server group can only add the same VPC ECS instances.

For information about server group and how to use it, see [Configure a server group](https://www.alibabacloud.com/help/en/doc-detail/35215.html).


## Example Usage

```terraform
variable "slb_server_group_name" {
  default = "forSlbServerGroup"
}

data "alicloud_zones" "server_group" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "server_group" {
  vpc_name   = var.slb_server_group_name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "server_group" {
  vpc_id       = alicloud_vpc.server_group.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.server_group.zones[0].id
  vswitch_name = var.slb_server_group_name
}


resource "alicloud_slb_load_balancer" "server_group" {
  load_balancer_name   = var.slb_server_group_name
  vswitch_id           = alicloud_vswitch.server_group.id
  instance_charge_type = "PayByCLCU"
}

resource "alicloud_slb_server_group" "server_group" {
  load_balancer_id = alicloud_slb_load_balancer.server_group.id
  name             = var.slb_server_group_name
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new virtual server group.
* `name` - (Optional) Name of the virtual server group. Our plugin provides a default name: "tf-server-group".
* `servers` - A list of ECS instances to be added. **NOTE:** Field 'servers' has been deprecated from provider version 1.163.0 and it will be removed in the future version. Please use the new resource 'alicloud_slb_server_group_server_attachment'. At most 20 ECS instances can be supported in one resource. It contains three sub-fields as `Block server` follows.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Block servers

The servers mapping supports the following:

* `server_ids` - (Required) A list backend server ID (ECS instance ID).
* `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.
* `type` - (Optional, Available in 1.51.0+) Type of the backend server. Valid value ecs, eni. Default to eni.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual server group.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new virtual server group.
* `name` - The name of the virtual server group.
* `servers` - A list of ECS instances that have be added.

## Import

Load balancer backend server group can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_server_group.example abc123456
```

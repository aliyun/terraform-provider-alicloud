---
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

## Example Usage

```
# Create a new load balancer and virtual server group

resource "alicloud_instance" "instance" {
  instance_name = "for-slb-server"
  count = 3
  ...
}

resource "alicloud_slb" "instance" {
  name = "new-slb"
  vswitch_id = "<one vswitch id>"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  servers = [
    {
      server_ids = ["${alicloud_instance.instance.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new virtual server group.
* `name` - (Optional) Name of the virtual server group. Our plugin provides a default name: "tf-server-group".
* `servers` - A list of ECS instances to be added. At most 20 ECS instances can be supported in one resource. It contains three sub-fields as `Block server` follows.

## Block servers

The servers mapping supports the following:

* `server_ids` - (Required) A list backend server ID (ECS instance ID).
* `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual server group.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new virtual server group.
* `name` - The name of the virtual server group.
* `servers` - A list of ECS instances that have be added.

## Import

Load balancer backend server group can be imported using the id, e.g.

```
$ terraform import alicloud_slb_server_group.example abc123456
```

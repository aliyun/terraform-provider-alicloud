---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_group_server_attachment"
sidebar_current: "docs-alicloud-resource-slb-server-group-server-attachment"
description: |-
  Provides a Load Banlancer Virtual Backend Server Group Server Attachment resource.
---

# alicloud\_slb\_server\_group\_server\_attachment

-> **NOTE:** Available in v1.163.0+.

For information about server group server attachment and how to use it, see [Configure a server group server attachment](https://www.alibabacloud.com/help/en/doc-detail/35218.html).

## Example Usage

```
variable "name" {
  default = "slbservergroupvpc"
}

variable "num" {
  default = 5
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  count                      = var.num
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = data.alicloud_vswitches.default.vswitches.0.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = var.name
}

resource "alicloud_slb_server_group_server_attachment" "default" {
  count           = var.num
  server_group_id = alicloud_slb_server_group.default.id
  server_id       = alicloud_instance.default[count.index].id
  port            = 8080
  weight          = 0
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  backend_port     = "80"
  frontend_port    = "80"
  protocol         = "tcp"
  bandwidth        = 10
  scheduler        = "rr"
  server_group_id  = alicloud_slb_server_group.default.id
}
```

## Argument Reference

The following arguments are supported:

* `server_group_id` - (Required, ForceNew) The ID of the server group.
* `server_id` - (Required, ForceNew) The ID of the backend server. You can specify the ID of an Elastic Compute Service (ECS) instance or an elastic network interface (ENI).
* `port` - (Required, ForceNew) The port that is used by the backend server. Valid values: `1` to `65535`.
* `weight` - (Optional, ForceNew, Computed) The weight of the backend server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no requests are forwarded to the backend server.
* `type` - (Optional, ForceNew, Computed) The type of backend server. Valid values: `ecs`, `eni`.
* `description` - (Optional, ForceNew, Computed) The description of the backend server.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual server group server attachment. The value formats as `<server_group_id>:<server_id>:<port>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.


## Import

Load balancer backend server group server attachment can be imported using the id, e.g.

```
$ terraform import alicloud_slb_server_group_server_attachment.example <server_group_id>:<server_id>:<port>
```

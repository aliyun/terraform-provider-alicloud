---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group_server_attachment"
description: |-
  Provides a Alicloud NLB Server Group Server Attachment resource.
---

# alicloud_nlb_server_group_server_attachment

Provides a NLB Server Group Server Attachment resource. Network Server Load Balancer.

For information about NLB Server Group Server Attachment and how to use it, see [What is Server Group Server Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/addserverstoservergroup-nlb).

-> **NOTE:** Available since v1.192.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultlHBOhp" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultV4HL2d" {
  vpc_id       = alicloud_vpc.defaultlHBOhp.id
  cidr_block   = "10.2.0.0/16"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

}

resource "alicloud_nlb_server_group" "defaultg9h9VW" {
  address_ip_version = "Ipv4"
  scheduler          = "Wrr"
  health_check {
  }
  server_group_type = "Instance"
  vpc_id            = alicloud_vpc.defaultlHBOhp.id
  protocol          = "TCP"
  server_group_name = var.name

}

resource "alicloud_ecs_instance" "defaultNzHh7X" {
  system_disk {
    performance_level = "PL1"
    size              = "100"
    category          = "cloud_essd"
  }
  spot_strategy = "NoSpot"
  image_id      = "aliyun_2_1903_x64_20G_alibase_20200324.vhd"
  vpc_attributes {
    vpc_id     = alicloud_vpc.defaultlHBOhp.id
    vswitch_id = alicloud_vswitch.defaultV4HL2d.id
  }
  internet_charge_type = "PayByTraffic"
  instance_name        = var.name

  internet_max_bandwidth_out = "2"
  description                = "MyFirstEcsInstance"
  instance_type              = "ecs.e-c1m1.large"
  zone_id                    = "cn-beijing-j"
}


resource "alicloud_nlb_server_group_server_attachment" "default" {
  server_type     = "Ecs"
  server_group_id = alicloud_nlb_server_group.defaultg9h9VW.id
  server_id       = alicloud_ecs_instance.defaultNzHh7X.id
  port            = "80"
  description     = "ertwgs"
  server_ip       = alicloud_ecs_instance.defaultNzHh7X.private_ip_address[0]
  weight          = "80"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the servers. The description must be 2 to 256 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (_), and hyphens (-).
* `port` - (Optional, ForceNew, Computed) The port used by the backend server. Valid values: 1 to 65535.
* `server_group_id` - (Required, ForceNew) The ID of the server group.
* `server_id` - (Required, ForceNew) The ID of the server.
  - If the server group type is Instance, set the ServerId parameter to the ID of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. These backend servers are specified by Ecs, Eni, or Eci. 
  - If the server group type is Ip, set the ServerId parameter to an IP address.
* `server_ip` - (Optional, ForceNew) The IP address of the server. If the server group type is Ip, set the ServerId parameter to an IP address.
* `server_type` - (Required, ForceNew) The type of the backend server. Valid values: `Ecs`, `Eni`, `Eci`, `Ip`.
* `weight` - (Optional) The weight of the backend server. Valid values: 0 to 100. Default value: 100. If the weight of a backend server is set to 0, no requests are forwarded to the backend server.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<server_group_id>:<server_id>:<server_type>:<port>`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group Server Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group Server Attachment.
* `update` - (Defaults to 5 mins) Used when update the Server Group Server Attachment.

## Import

NLB Server Group Server Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_server_group_server_attachment.example <server_group_id>:<server_id>:<server_type>:<port>
```
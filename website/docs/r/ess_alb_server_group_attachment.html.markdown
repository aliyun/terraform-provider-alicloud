---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_alb_server_group_attachment"
sidebar_current: "docs-alicloud-resource-ess-alb-server-group-attachment"
description: |-
  Provides a ESS Attachment resource to attach or remove alb server group.
---

# alicloud_ess_alb_server_group_attachment

Attaches/Detaches alb server group to a specified scaling group.

For information about alb server group attachment, see [AttachAlbServerGroups](https://www.alibabacloud.com/help/en/doc-detail/266800.html).

-> **NOTE:** If scaling group's network type is `VPC`, the alb server groups must be in the same `VPC`.

-> **NOTE:** Alb server group attachment is defined uniquely by `scaling_group_id`, `alb_server_group_id`, `port`.

-> **NOTE:** Resource `alicloud_ess_alb_server_group_attachment` don't support modification.

-> **NOTE:** Available since v1.158.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_security_group" "default" {
  name   = local.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = "0"
  max_size           = "2"
  scaling_group_name = local.name
  default_cooldown   = 200
  removal_policies   = ["OldestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
  enable            = true
}

resource "alicloud_alb_server_group" "default" {
  server_group_name = local.name
  vpc_id            = alicloud_vpc.default.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-example"
    sticky_session_type    = "Server"
  }
}

resource "alicloud_ess_alb_server_group_attachment" "default" {
  scaling_group_id    = alicloud_ess_scaling_configuration.default.scaling_group_id
  alb_server_group_id = alicloud_alb_server_group.default.id
  port                = 9000
  weight              = 50
  force_attach        = true
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group.
* `alb_server_group_id` - (Required, ForceNew) ID of Alb Server Group.
* `port` - (Required, ForceNew) - The port will be used for Alb Server Group backend server.
* `weight` - (Required, ForceNew) The weight of an ECS instance attached to the Alb Server Group.
* `force_attach` - (Optional) If instances of scaling group are attached/removed from slb backend server when attach/detach alb
  server group from scaling group. Default to false.

## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS alb server group attachment resource ID，in the follwing format: scaling_group_id:
  alb_server_group_id:port.

## Import

ESS alb server groups can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_alb_server_group_attachment.example asg-xxx:sgp-xxx:5000 
```

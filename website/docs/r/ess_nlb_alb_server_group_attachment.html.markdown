---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_nlb_alb_server_group_attachment"
sidebar_current: "docs-alicloud_ess_nlb_alb_server_group_attachment"
description: |-
  Provides a ESS Attachment resource to attach or remove nlb or alb server group.
---

# alicloud_ess_nlb_alb_server_group_attachment

Attaches/Detaches server groups to a specified scaling group.

For information about server group attachment, see [AttachServerGroups](https://www.alibabacloud.com/help/en/auto-scaling/developer-reference/api-attachservergroups?spm=a2c63.p38356.0.0.2bf515dd9wXWwC).

-> **NOTE:** If scaling group's network type is `VPC`, the alb server groups must be in the same `VPC`.

-> **NOTE:** server group attachment is defined uniquely by `scaling_group_id`, `server_group_id`, `type`, `port`.

-> **NOTE:** Resource `alicloud_ess_nlb_alb_server_group_attachment` don't support modification.

-> **NOTE:** Available since v1.228.0.

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

resource "alicloud_ess_nlb_alb_server_group_attachment" "default" {
  scaling_group_id = alicloud_ess_scaling_configuration.default.scaling_group_id
  server_groups {
    server_group_id = alicloud_alb_server_group.default.id
    type            = "ALB"
    port            = "210"
    weight          = "60"
  }
  force_attach = true
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group.
* `server_groups` - (Required) A list of vserver groups attached on scaling group. See [`server_groups`](#server_groups) below.
* `force_attach` - (Optional) If instances of scaling group are attached/removed from server when attach/detach server group from scaling group. Default to false.

### `server_groups`

the server_group supports the following:

* `server_group_id` - (Required) The ID of server group.
* `type` - (Required) The type of server group. Valid values: NLB, ALB.
* `weight` - (Required) The weight of an ECS instance or elastic container instance as a backend server of server group. Valid values: 0 to 100.
* `port` - (Required) The port used by ECS instances or elastic container instances after being added as backend servers to server group. Valid values: 1 to 65535.



## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS server group attachment resource ID.

## Import

ESS server group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_nlb_alb_server_group_attachment.example abc123456
```

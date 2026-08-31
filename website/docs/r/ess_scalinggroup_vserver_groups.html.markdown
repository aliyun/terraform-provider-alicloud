---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scalinggroup_vserver_groups"
sidebar_current: "docs-alicloud-resource-ess_scalinggroup_vserver_groups"
description: |-
  Provides a ESS Attachment resource to attach or remove vserver groups.
---

# alicloud_ess_scalinggroup_vserver_groups

Attaches/Detaches vserver groups to a specified scaling group.

-> **NOTE:** The load balancer of which vserver groups belongs to must be in `active` status.

-> **NOTE:** If scaling group's network type is `VPC`, the vserver groups must be in the same `VPC`.
 
-> **NOTE:** A scaling group can have at most 5 vserver groups attached by default.

-> **NOTE:** Vserver groups and the default group of loadbalancer share the same backend server quota.

-> **NOTE:** When attach vserver groups to scaling group, existing ECS instances will be added to vserver groups; Instead, ECS instances will be removed from vserver group when detach.

-> **NOTE:** Detach action will be executed before attach action.

-> **NOTE:** Vserver group is defined uniquely by `loadbalancer_id`, `vserver_group_id`, `port`.

-> **NOTE:** Modifing `weight` attribute means detach vserver group first and then, attach with new weight parameter.

-> **NOTE:** Available since v1.53.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_scalinggroup_vserver_groups&exampleId=9ae7a110-84b0-e4da-64b4-6132a7bb2a1227867466&activeTab=example&spm=docs.r.ess_scalinggroup_vserver_groups.0.9ae7a11084&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_slb_load_balancer" "default" {
  count              = 2
  load_balancer_name = format("terraform-example%d", count.index + 1)
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_server_group" "default1" {
  count            = "2"
  load_balancer_id = alicloud_slb_load_balancer.default.0.id
  name             = local.name
}

resource "alicloud_slb_server_group" "default2" {
  count            = "2"
  load_balancer_id = alicloud_slb_load_balancer.default.1.id
  name             = local.name
}

resource "alicloud_slb_listener" "default" {
  count             = 2
  load_balancer_id  = alicloud_slb_load_balancer.default[count.index].id
  backend_port      = "22"
  frontend_port     = "22"
  protocol          = "tcp"
  bandwidth         = "10"
  health_check_type = "tcp"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = "2"
  max_size           = "2"
  scaling_group_name = local.name
  default_cooldown   = 200
  removal_policies   = ["OldestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
  loadbalancer_ids   = alicloud_slb_listener.default.*.load_balancer_id
}

resource "alicloud_ess_scalinggroup_vserver_groups" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  vserver_groups {
    loadbalancer_id = alicloud_slb_load_balancer.default.0.id
    vserver_attributes {
      vserver_group_id = alicloud_slb_server_group.default1.0.id
      port             = "100"
      weight           = "60"
    }
    vserver_attributes {
      vserver_group_id = alicloud_slb_server_group.default1.1.id
      port             = "110"
      weight           = "60"
    }
  }
  vserver_groups {
    loadbalancer_id = alicloud_slb_load_balancer.default.1.id
    vserver_attributes {
      vserver_group_id = alicloud_slb_server_group.default2.0.id
      port             = "200"
      weight           = "60"
    }
    vserver_attributes {
      vserver_group_id = alicloud_slb_server_group.default2.1.id
      port             = "210"
      weight           = "60"
    }
  }
  force = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ess_scalinggroup_vserver_groups&spm=docs.r.ess_scalinggroup_vserver_groups.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group.
* `vserver_groups` - (Required) A list of vserver groups attached on scaling group. See [`vserver_groups`](#vserver_groups) below.
* `force` - (Optional, Available in 1.64.0+) If instances of scaling group are attached/removed from slb backend server when attach/detach vserver group from scaling group. Default to true.

### `vserver_groups`

the vserver_group supports the following:

* `loadbalancer_id` - (Required) Loadbalancer server ID of VServer Group.
* `vserver_attributes` - (Required) A list of VServer Group attributes. See [`vserver_attributes`](#vserver_groups-vserver_attributes) below.

### `vserver_groups-vserver_attributes`

* `vserver_group_id` - (Required) ID of VServer Group.
* `port` - (Required) - The port will be used for VServer Group backend server.
* `weight` - (Required) The weight of an ECS instance attached to the VServer Group.

## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS vserver groups attachment resource ID.

## Import

ESS vserver groups can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_vserver_groups.example abc123456
```

---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_backend_server"
sidebar_current: "docs-alicloud-resource-slb-backend-server"
description: |-
  Provides an Application Load Balancer Backend Server resource.
---

# alicloud\_slb\_backend\_server

Add a group of backend servers (ECS or ENI instance) to the Server Load Balancer or remove them from it.

-> **NOTE:** Available in 1.53.0+

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_backend_server&exampleId=6a67a5d7-eaa1-8639-5767-c28204091ffd76966ef1&activeTab=example&spm=docs.r.slb_backend_server.0.6a67a5d7ea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform

// alicloud_slb_backend_server
variable "slb_backend_server_name" {
  default = "slbbackendservertest"
}

data "alicloud_zones" "backend_server" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "backend_server" {
  availability_zone = data.alicloud_zones.backend_server.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "backend_server" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "backend_server" {
  vpc_name   = var.slb_backend_server_name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "backend_server" {
  vpc_id       = alicloud_vpc.backend_server.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.backend_server.zones[0].id
  vswitch_name = var.slb_backend_server_name
}

resource "alicloud_security_group" "backend_server" {
  name   = var.slb_backend_server_name
  vpc_id = alicloud_vpc.backend_server.id
}

resource "alicloud_instance" "backend_server" {
  image_id                   = data.alicloud_images.backend_server.images[0].id
  instance_type              = data.alicloud_instance_types.backend_server.instance_types[0].id
  instance_name              = var.slb_backend_server_name
  count                      = "2"
  security_groups            = alicloud_security_group.backend_server.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.backend_server.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.backend_server.id
}

resource "alicloud_slb_load_balancer" "backend_server" {
  load_balancer_name   = var.slb_backend_server_name
  vswitch_id           = alicloud_vswitch.backend_server.id
  instance_charge_type = "PayByCLCU"
}

resource "alicloud_slb_backend_server" "backend_server" {
  load_balancer_id = alicloud_slb_load_balancer.backend_server.id

  backend_servers {
    server_id = alicloud_instance.backend_server[0].id
    weight    = 100
  }

  backend_servers {
    server_id = alicloud_instance.backend_server[1].id
    weight    = 100
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the load balancer.
* `backend_servers` - (Optional) A list of instances to added backend server in the SLB. It contains three sub-fields as `Block server` follows.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Block servers

The servers mapping supports the following:

* `server_id` - (Required) A list backend server ID (ECS instance ID).
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. 
* `type` - (Optional) Type of the backend server. Valid value `ecs`, `eni`, `eci`. Default to `ecs`. **NOTE:** From 1.170.0+, The `eci` is valid. 
* `server_ip` - (Optional, Computed, Available in 1.93.0+) ServerIp of the backend server. This parameter can be specified when the type is `eni`. `ecs` type currently does not support adding `server_ip` parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource and the value same as load balancer id.

## Import

Load balancer backend server can be imported using the load balancer id.

```shell
$ terraform import alicloud_slb_backend_server.example <load_balancer_id>
```

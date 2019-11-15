---
subcategory: "Server Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_backend_server"
sidebar_current: "docs-alicloud-resource-slb-backend-server"
description: |-
  Provides an Application Load Balancer Backend Server resource.
---

# alicloud\_slb\_backend\_server

Add a group of backend servers (ECS instance) to the Server Load Balancer or remove them from it.

-> **NOTE:** Available in 1.53.0+

## Example Usage

```
variable "name" {
  default = "slbbackendservertest"
}
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name       = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_backend_server" "default" {
  	load_balancer_id = "${alicloud_slb.default.id}"
  	
	backend_servers {
      server_id = "${alicloud_instance.default.0.id}"
      weight     = 100
    }

    backend_servers {
      server_id = "${alicloud_instance.default.1.id}"
      weight     = 100
    }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the load balancer.
* `backend_servers` - (Required) A list of instances to added backend server in the SLB. It contains three sub-fields as `Block server` follows.

## Block servers

The servers mapping supports the following:

* `server_id` - (Required) A list backend server ID (ECS instance ID).
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. 
* `type` - (Optional) Type of the backend server. Valid value ecs, eni. Default to eni.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource and the value same as load balancer id.

## Import

Load balancer backend server can be imported using the load balancer id.

```
$ terraform import alicloud_slb_backend_server.example lb-abc123456
```

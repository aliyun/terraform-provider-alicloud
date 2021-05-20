---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_master_slave_server_group"
sidebar_current: "docs-alicloud-resource-slb-master-slave-server-group"
description: |-
  Provides a Load Banlancer Master Slave Server Group resource.
---

# alicloud\_slb\_master\_slave\_server\_group

A master slave server group contains two ECS instances. The master slave server group can help you to define multiple listening dimension.

-> **NOTE:** One ECS instance can be added into multiple master slave server groups.

-> **NOTE:** One master slave server group can only add two ECS instances, which are master server and slave server.

-> **NOTE:** One master slave server group can be attached with tcp/udp listeners in one load balancer.

-> **NOTE:** One Classic and Internet load balancer, its master slave server group can add Classic and VPC ECS instances.

-> **NOTE:** One Classic and Intranet load balancer, its master slave server group can only add Classic ECS instances.

-> **NOTE:** One VPC load balancer, its master slave server group can only add the same VPC ECS instances.

-> **NOTE:** Available in 1.54.0+

## Example Usage

```
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount        = 2
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccSlbMasterSlaveServerGroupVpc"
}

variable "number" {
  default = "1"
}

resource "alicloud_vpc" "main" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id            = alicloud_vpc.main.id
  cidr_block        = "172.16.0.0/16"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = alicloud_vpc.main.id
}

resource "alicloud_instance" "instance" {
  image_id                   = data.alicloud_images.image.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  count                      = "2"
  security_groups            = [alicloud_security_group.group.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.main.id
}

resource "alicloud_slb_load_balancer" "instance" {
  load_balancer_name  = var.name
  vswitch_id    = alicloud_vswitch.main.id
  load_balancer_spec = "slb.s2.small"
}

resource "alicloud_network_interface" "default" {
  count           = var.number
  name            = var.name
  vswitch_id      = alicloud_vswitch.main.id
  security_groups = [alicloud_security_group.group.id]
}

resource "alicloud_network_interface_attachment" "default" {
  count                = var.number
  instance_id          = alicloud_instance.instance[0].id
  network_interface_id = element(alicloud_network_interface.default.*.id, count.index)
}

resource "alicloud_slb_master_slave_server_group" "group" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
  name             = var.name

  servers {
    server_id   = alicloud_instance.instance[0].id
    port        = 100
    weight      = 100
    server_type = "Master"
  }

  servers {
    server_id   = alicloud_instance.instance[1].id
    port        = 100
    weight      = 100
    server_type = "Slave"
  }
}

resource "alicloud_slb_listener" "tcp" {
  load_balancer_id             = alicloud_slb_load_balancer.instance.id
  master_slave_server_group_id = alicloud_slb_master_slave_server_group.group.id
  frontend_port                = "22"
  protocol                     = "tcp"
  bandwidth                    = "10"
  health_check_type            = "tcp"
  persistence_timeout          = 3600
  healthy_threshold            = 8
  unhealthy_threshold          = 8
  health_check_timeout         = 8
  health_check_interval        = 5
  health_check_http_code       = "http_2xx"
  health_check_connect_port    = 20
  health_check_uri             = "/console"
  established_timeout          = 600
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new master slave server group.
* `name` - (Required, ForceNew) Name of the master slave server group. 
* `servers` - (Optional, ForceNew) A list of ECS instances to be added. Only two ECS instances can be supported in one resource. It contains six sub-fields as `Block server` follows.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Block servers

The servers mapping supports the following:

* `server_ids` - (Required) A list backend server ID (ECS instance ID).
* `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.
* `type` - (Optional, Available in 1.51.0+) Type of the backend server. Valid value ecs, eni. Default to eni.
* `server_type` - (Optional) The server type of the backend server. Valid value Master, Slave.
* `is_backup` - (Removed from v1.63.0) Determine if the server is executing. Valid value 0, 1. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the master slave server group.

## Import

Load balancer master slave server group can be imported using the id, e.g.

```
$ terraform import alicloud_slb_master_slave_server_group.example abc123456
```

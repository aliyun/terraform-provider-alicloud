---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_master_slave_server_group"
sidebar_current: "docs-alicloud-resource-slb-master-slave-server-group"
description: |-
  Provides a Load Banlancer Master Slave Server Group resource.
---

# alicloud_slb_master_slave_server_group

A master slave server group contains two ECS instances. The master slave server group can help you to define multiple listening dimension.

-> **NOTE:** One ECS instance can be added into multiple master slave server groups.

-> **NOTE:** One master slave server group can only add two ECS instances, which are master server and slave server.

-> **NOTE:** One master slave server group can be attached with tcp/udp listeners in one load balancer.

-> **NOTE:** One Classic and Internet load balancer, its master slave server group can add Classic and VPC ECS instances.

-> **NOTE:** One Classic and Intranet load balancer, its master slave server group can only add Classic ECS instances.

-> **NOTE:** One VPC load balancer, its master slave server group can only add the same VPC ECS instances.

-> **NOTE:** Available since v1.54.0+

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_master_slave_server_group&exampleId=801c7454-2bd4-26f6-8a51-b5e8e13657270a27203e&activeTab=example&spm=docs.r.slb_master_slave_server_group.0.801c74542b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "ms_server_group" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "ms_server_group" {
  availability_zone    = data.alicloud_zones.ms_server_group.zones[0].id
  cpu_core_count       = 2
  memory_size          = 8
  instance_type_family = "ecs.g6"
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "slb_master_slave_server_group" {
  default = "forSlbMasterSlaveServerGroup"
}

resource "alicloud_vpc" "main" {
  vpc_name   = var.slb_master_slave_server_group
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.ms_server_group.zones[0].id
  vswitch_name = var.slb_master_slave_server_group
}

resource "alicloud_security_group" "group" {
  security_group_name = var.slb_master_slave_server_group
  vpc_id              = alicloud_vpc.main.id
}

resource "alicloud_instance" "ms_server_group" {
  image_id                   = data.alicloud_images.image.images[0].id
  instance_type              = data.alicloud_instance_types.ms_server_group.instance_types[0].id
  instance_name              = var.slb_master_slave_server_group
  count                      = 2
  security_groups            = [alicloud_security_group.group.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.ms_server_group.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.main.id
}

resource "alicloud_slb_load_balancer" "ms_server_group" {
  load_balancer_name = var.slb_master_slave_server_group
  vswitch_id         = alicloud_vswitch.main.id
  load_balancer_spec = "slb.s2.small"
}

resource "alicloud_ecs_network_interface" "ms_server_group" {
  network_interface_name = var.slb_master_slave_server_group
  vswitch_id             = alicloud_vswitch.main.id
  security_group_ids     = [alicloud_security_group.group.id]
}

resource "alicloud_ecs_network_interface_attachment" "ms_server_group" {
  instance_id          = alicloud_instance.ms_server_group[0].id
  network_interface_id = alicloud_ecs_network_interface.ms_server_group.id
}

resource "alicloud_slb_master_slave_server_group" "group" {
  load_balancer_id = alicloud_slb_load_balancer.ms_server_group.id
  name             = var.slb_master_slave_server_group

  servers {
    server_id   = alicloud_instance.ms_server_group[0].id
    port        = 100
    weight      = 100
    server_type = "Master"
  }

  servers {
    server_id   = alicloud_instance.ms_server_group[1].id
    port        = 100
    weight      = 100
    server_type = "Slave"
  }
}

resource "alicloud_slb_listener" "tcp" {
  load_balancer_id             = alicloud_slb_load_balancer.ms_server_group.id
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_slb_master_slave_server_group&spm=docs.r.slb_master_slave_server_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new master slave server group.
* `name` - (Required, ForceNew) Name of the master slave server group. 
* `servers` - (Optional, ForceNew) A list of ECS instances to be added. Only two ECS instances can be supported in one resource. See [`servers`](#servers) below.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

### `servers`

The servers mapping supports the following:

* `server_id` - (Required) A list backend server ID (ECS instance ID).
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

```shell
$ terraform import alicloud_slb_master_slave_server_group.example abc123456
```

---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_master_slave_server_groups"
sidebar_current: "docs-alicloud-datasource-slb-master-slave-server-groups"
description: |-
    Provides a list of master slave server groups related to a server load balancer to the user.
---

# alicloud\_slb\_master\_slave\_server\_groups

This data source provides the master slave server groups related to a server load balancer.

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
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
  cidr_block = "172.16.0.0/16"
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
  load_balancer_name = var.name
  vswitch_id    = alicloud_vswitch.main.id
  load_balancer_spec = "slb.s2.small"
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

data "alicloud_slb_master_slave_server_groups" "sample_ds" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
}

output "first_slb_server_group_id" {
  value = data.alicloud_slb_master_slave_server_groups.sample_ds.groups[0].id
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - ID of the SLB.
* `ids` - (Optional) A list of master slave server group IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by master slave server group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB master slave server groups IDs.
* `names` - A list of SLB master slave server groups names.
* `groups` - A list of SLB master slave server groups. Each element contains the following attributes:
  * `id` - master slave server group ID.
  * `name` - master slave server group name.
  * `servers` - ECS instances associated to the group. Each element contains the following attributes:
    * `instance_id` - ID of the attached ECS instance.
    * `weight` - Weight associated to the ECS instance.
    * `port` - The port used by the master slave server group.
    * `server_type` - The server type of the attached ECS instance.
    * `is_backup` - (Removed from v1.63.0) Determine if the server is executing.


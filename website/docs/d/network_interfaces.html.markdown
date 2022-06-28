---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interfaces"
sidebar_current: "docs-alicloud-datasource-network-interfaces"
description: |-
  Provides a data source to get a list of elastic network interfaces according to the specified filters.
---

# alicloud\_network_interfaces

-> **DEPRECATED:** This datasource has been renamed to [alicloud_ecs_network_interfaces](https://www.terraform.io/docs/providers/alicloud/d/ecs_network_interfaces) from version 1.123.1.

Use this data source to get a list of elastic network interfaces according to the specified filters in an Alibaba Cloud account.

For information about elastic network interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html)

## Example Usage

```
variable "name" {
  default = "networkInterfacesName"
}

resource "alicloud_vpc" "vpc" {
  vpc_name = "${var.name}"
  cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
  vswitch_name = "${var.name}"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "group" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "interface" {
  name = "${var.name}%d"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
  security_groups = [
    "${alicloud_security_group.group.id}"]
  description = "Basic test"
  private_ip = "192.168.0.2"
  tags = {
    TF-VER = "0.11.3"
  }
}

resource "alicloud_instance" "instance" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  security_groups = [
    "${alicloud_security_group.group.id}"]
  instance_type = "ecs.e3.xlarge"
  system_disk_category = "cloud_efficiency"
  image_id = "centos_7_04_64_20G_alibase_201701015.vhd"
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
  internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface_attachment" "attachment" {
  instance_id = "${alicloud_instance.instance.id}"
  network_interface_id = "${alicloud_network_interface.interface.id}"
}

data "alicloud_network_interfaces" "default" {
  ids = [
    "${alicloud_network_interface_attachment.attachment.network_interface_id}"]
  name_regex = "${var.name}"
  tags = {
    TF-VER = "0.11.3"
  }
  vpc_id = "${alicloud_vpc.vpc.id}"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
  private_ip = "192.168.0.2"
  security_group_id = "${alicloud_security_group.group.id}"
  type = "Secondary"
  instance_id = "${alicloud_instance.instance.id}"
}

output "eni0_name" {
    value = "${data.alicloud_network_interfaces.default.interfaces.0.name}"
}
```

##  Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of ENI IDs.
* `name_regex` - (Optional) A regex string to filter results by ENI name.
* `vpc_id` - (Optional) The VPC ID linked to ENIs.
* `vswitch_id` - (Optional) The VSwitch ID linked to ENIs.
* `private_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID linked to ENIs.
* `name` - (Optional) The name of the ENIs.
* `type` - (Optional) The type of ENIs, Only support for "Primary" or "Secondary".
* `instance_id` - (Optional) The ECS instance ID that the ENI is attached to.
* `tags` - (Optional) A map of tags assigned to ENIs.
* `output_file` - (Optional) The name of output file that saves the filter results.
* `resource_group_id` - (Optional, ForceNew, Available in 1.57.0+) The Id of resource group which the network interface belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `interfaces` - A list of ENIs. Each element contains the following attributes:
    * `id` - ID of the ENI.
    * `status` - Current status of the ENI.
    * `vpc_id` - ID of the VPC that the ENI belongs to.
    * `vswitch_id` - ID of the VSwitch that the ENI is linked to.
    * `zone_id` - ID of the availability zone that the ENI belongs to.
    * `public_ip` - Public IP of the ENI.
    * `private_ip` - Primary private IP of the ENI.
    * `private_ips` - A list of secondary private IP address that is assigned to the ENI.
    * `mac` - MAC address of the ENI.
    * `security_groups` - A list of security group that the ENI belongs to.
    * `name` - Name of the ENI.
    * `description` - Description of the ENI.
    * `instance_id` - ID of the instance that the ENI is attached to.
    * `creation_time` - Creation time of the ENI.
    * `tags` - A map of tags assigned to the ENI.
    * `resource_group_id` - The Id of resource group.

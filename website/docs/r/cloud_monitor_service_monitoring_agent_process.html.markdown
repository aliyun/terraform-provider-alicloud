---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_monitoring_agent_process"
description: |-
  Provides a Alicloud Cloud Monitor Service Monitoring Agent Process resource.
---

# alicloud_cloud_monitor_service_monitoring_agent_process

Provides a Cloud Monitor Service Monitoring Agent Process resource. 

For information about Cloud Monitor Service Monitoring Agent Process and how to use it, see [What is Monitoring Agent Process](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2019-01-01-createmonitoragentprocess).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vswitch.default.vpc_id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}

resource "alicloud_cloud_monitor_service_monitoring_agent_process" "default" {
  instance_id  = alicloud_instance.default.id
  process_name = var.name
  process_user = "root"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `process_name` - (Required, ForceNew) The name of the process.
* `process_user` - (Optional, ForceNew) The user who launches the process.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitoring Agent Process. It formats as `<instance_id>:<process_id>`.
* `process_id` - The ID of the process.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Monitoring Agent Process.
* `delete` - (Defaults to 5 mins) Used when delete the Monitoring Agent Process.

## Import

Cloud Monitor Service Monitoring Agent Process can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_monitoring_agent_process.example <instance_id>:<process_id>
```

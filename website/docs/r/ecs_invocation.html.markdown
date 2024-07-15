---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_invocation"
sidebar_current: "docs-alicloud-resource-ecs-invocation"
description: |-
  Provides a Alicloud ECS Invocation resource.
---

# alicloud_ecs_invocation

Provides a ECS Invocation resource.

For information about ECS Invocation and how to use it, see [What is Invocation](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/invokecommand#t9958.html).

-> **NOTE:** Available since v1.168.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
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
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_instance" "default" {
  vswitch_id                 = alicloud_vswitch.default.id
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = ["${alicloud_security_group.default.id}"]
  instance_name              = var.name
}

resource "alicloud_ecs_command" "default" {
  name             = var.name
  command_content  = "ZWNobyBoZWxsbyx7e25hbWV9fQ=="
  description      = "For Terraform Test"
  type             = "RunShellScript"
  working_dir      = "/root"
  enable_parameter = true
}

resource "alicloud_ecs_invocation" "default" {
  command_id  = alicloud_ecs_command.default.id
  instance_id = [alicloud_instance.default.id]
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Required, ForceNew) The ID of the command.
* `instance_id` - (Required, ForceNew) The list of instances to execute the command. You can specify up to 50 instance IDs.
* `repeat_mode` - (Optional, ForceNew, Computed) Specifies how to run the command. Valid values: `Once`, `Period`, `NextRebootOnly`, `EveryReboot`. Default value: When `timed` is set to false and Frequency is not specified, the default value of `repeat_mode` is `Once`. When `Timed` is set to true and Frequency is specified, `period` is used as the value of RepeatMode regardless of whether `repeat_mode` is specified.
* `timed` - (Optional, ForceNew, Computed) Specifies whether to periodically run the command. Default value: `false`.
* `frequency` - (Optional, ForceNew) The schedule on which the recurring execution of the command takes place. Take note of the following items:
  * The interval between two consecutive executions must be 10 seconds or longer. The minimum interval cannot be less than the timeout period of the execution.
  * When you set Timed to true, you must specify Frequency.
  * The value of the Frequency parameter is a cron expression. For more information, see [Cron expression](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/cron-expression).
* `parameters` - (Optional, ForceNew) The key-value pairs of custom parameters to be passed in when the custom parameter feature is enabled.  Number of custom parameters: 0 to 10.
* `username` - (Optional, ForceNew, Computed) The username that is used to run the command on the ECS instance. 
  * For Linux instances, the root username is used. 
  * For Windows instances, the System username is used.
  * You can also specify other usernames that already exist in the ECS instance to run the command. It is more secure to run Cloud Assistant commands as a regular user. For more information, see [Configure a regular user to run Cloud Assistant commands](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/run-cloud-assistant-commands-as-a-regular-user).
* `windows_password_name` - (Optional, ForceNew) The name of the password used to run the command on a Windows instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Invocation.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the ECS Invocation.
* `delete` - (Defaults to 1 mins) Used when stop the ECS Invocation.

## Import

ECS Invocation can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_invocation.example <id>
```
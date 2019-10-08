---
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_schedule"
sidebar_current: "docs-alicloud-resource-ess-schedule"
description: |-
  Provides a ESS schedule resource.
---

# alicloud\_ess\_scheduled\_task

Provides a ESS schedule resource.

## Example Usage

```
variable "name" {
  default = "essscheduleconfig"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
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
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  vswitch_ids        = ["${alicloud_vswitch.default.id}"]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = "${alicloud_ess_scaling_group.default.id}"
  image_id          = "${data.alicloud_images.default.images.0.id}"
  instance_type     = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_group_id = "${alicloud_security_group.default.id}"
  force_delete      = "true"
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "alicloud_ess_scheduled_task" "default" {
  scheduled_action    = "${alicloud_ess_scaling_rule.default.ari}"
  launch_time         = "2019-05-21T11:37Z"
  scheduled_task_name = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Required) Operations performed when the scheduled task is triggered. Fill in the unique identifier of the scaling rule.
* `launch_time` - (Required) Operations performed when the scheduled task is triggered. Fill in the unique identifier of the scaling rule.
* `scheduled_task_name` - (Optional) Display name of the scheduled task, which must be 2-40 characters (English or Chinese) long.
* `description` - (Optional) Description of the scheduled task, which is 2-200 characters (English or Chinese) long.
* `launch_expiration_time` - (Optional) Time period within which the failed scheduled task is retried. The default value is 600s. Value range: [0, 21600]
* `recurrence_type` - (Optional) Type of the scheduled task to be repeated. RecurrenceType, RecurrenceValue and RecurrenceEndTime must be specified. Optional values:
    - Daily: Recurrence interval by day for a scheduled task.
    - Weekly: Recurrence interval by week for a scheduled task.
    - Monthly: Recurrence interval by month for a scheduled task.
* `recurrence_value` - (Optional) Value of the scheduled task to be repeated. RecurrenceType, RecurrenceValue and RecurrenceEndTime must be specified.
    - Daily: Only one value in the range [1,31] can be filled.
    - Weekly: Multiple values can be filled. The values of Sunday to Saturday are 0 to 6 in sequence. Multiple values shall be separated by a comma “,”.
    - Monthly: In the format of A-B. The value range of A and B is 1 to 31, and the B value must be greater than the A value.
* `recurrence_end_time` - (Optional) End time of the scheduled task to be repeated. The date format follows the ISO8601 standard and uses UTC time. It is in the format of YYYY-MM-DDThh:mmZ. A time point 90 days after creation or modification cannot be entered. RecurrenceType, RecurrenceValue and RecurrenceEndTime must be specified.                                  
* `task_enabled` - (Optional) Whether to enable the scheduled task. The default value is true.
                                  
                                 
## Attributes Reference

The following attributes are exported:

* `id` - The schedule task ID.

## Import

ESS schedule task can be imported using the id, e.g.

```
$ terraform import alicloud_ess_scheduled_task.example abc123456
```
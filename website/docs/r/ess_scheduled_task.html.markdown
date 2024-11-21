---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scheduled_task"
sidebar_current: "docs-alicloud-resource-ess-schedule"
description: |-
  Provides a ESS schedule resource.
---

# alicloud_ess_scheduled_task

Provides a ESS schedule resource.

For information about ess schedule task, see [Scheduled Tasks](https://www.alibabacloud.com/help/en/auto-scaling/latest/createscheduledtask).

-> **NOTE:** Available since v1.60.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_scheduled_task&exampleId=ab8ce419-74cd-8c67-296b-cbfb8df30213f9b4c48f&activeTab=example&spm=docs.r.ess_scheduled_task.0.ab8ce41974&intl_lang=EN_US" target="_blank">
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

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
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
  zone_id      = data.alicloud_zones.default.zones[0].id
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

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}-${random_integer.default.result}"
  vswitch_ids        = [alicloud_vswitch.default.id]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = "true"
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "alicloud_ess_scheduled_task" "default" {
  scheduled_action    = alicloud_ess_scaling_rule.default.ari
  launch_time         = formatdate("YYYY-MM-DD'T'hh:mm'Z'", timeadd(timestamp(), "24h"))
  scheduled_task_name = "${var.name}-${random_integer.default.result}"

  # for test
  lifecycle {
    ignore_changes = [launch_time]
  }
}
```

## Module Support

You can use to the existing [autoscaling-rule module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling-rule/alicloud) 
to create scheduled task, different type rules and alarm task one-click.

## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Optional) The operation to be performed when a scheduled task is triggered. Enter the unique identifier of a scaling rule. **NOTE:** You cannot specify `scheduled_action` and `scaling_group_id` at the same time.
* `scheduled_task_name` - (Optional) Display name of the scheduled task, which must be 2-40 characters (English or Chinese) long.
* `scaling_group_id` - (Optional, Available in 1.90.0+) The ID of the scaling group where the number of instances is modified when the scheduled task is triggered. After the `ScalingGroupId` parameter is specified, the scaling method of the scheduled task is to specify the number of instances in a scaling group. You must specify at least one of the following parameters: `MinValue`, `MaxValue`, and `DesiredCapacity`. **NOTE:** You cannot specify `scheduled_action` and `scaling_group_id` at the same time.
* `min_value` - (Optional, Available in 1.90.0+) The minimum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
* `max_value` - (Optional, Available in 1.90.0+) The maximum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
* `desired_capacity` - (Optional, Available in 1.90.0+) The expected number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group. **NOTE:** You must specify the `DesiredCapacity` parameter when you create the scaling group.
* `description` - (Optional) Description of the scheduled task, which is 2-200 characters (English or Chinese) long.
* `launch_time` - (Optional) The time at which the scheduled task is triggered. Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mmZ format. 
The time must be in UTC. You cannot enter a time point later than 90 days from the date of scheduled task creation. 
If the `recurrence_type` parameter is specified, the task is executed repeatedly at the time specified by LaunchTime. 
Otherwise, the task is only executed once at the date and time specified by LaunchTime.
* `launch_expiration_time` - (Optional) The time period during which a failed scheduled task is retried. Unit: seconds. Valid values: 0 to 21600. Default value: 600
* `recurrence_type` - (Optional) Specifies the recurrence type of the scheduled task. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time. Valid values:
    - Daily: The scheduled task is executed once every specified number of days.
    - Weekly: The scheduled task is executed on each specified day of a week.
    - Monthly: The scheduled task is executed on each specified day of a month.
    - Cron: (Available in 1.60.0+) The scheduled task is executed based on the specified cron expression.
* `recurrence_value` - (Optional) Specifies how often a scheduled task recurs. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time. The valid value depends on `recurrence_type`
    - Daily: You can enter one value. Valid values: 1 to 31.
    - Weekly: You can enter multiple values and separate them with commas (,). For example, the values 0 to 6 correspond to the days of the week in sequence from Sunday to Saturday.
    - Monthly: You can enter two values in A-B format. Valid values of A and B: 1 to 31. The value of B must be greater than or equal to the value of A.
    - Cron: You can enter a cron expression which is written in UTC and consists of five fields: minute, hour, day of month (date), month, and day of week. The expression can contain wildcard characters including commas (,), question marks (?), hyphens (-), asterisks (*), number signs (#), forward slashes (/), and the L and W letters.
* `recurrence_end_time` - (Optional) Specifies the end time after which the scheduled task is no longer repeated. Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mm:ssZ format. 
The time must be in UTC. You cannot enter a time point later than 365 days from the date of scheduled task creation. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time.                                
* `task_enabled` - (Optional) Specifies whether to start the scheduled task. Default to true.
                                  
                                 
## Attributes Reference

The following attributes are exported:

* `id` - The schedule task ID.

## Import

ESS schedule task can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_scheduled_task.example abc123456
```


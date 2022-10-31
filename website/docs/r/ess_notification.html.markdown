---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_notification"
sidebar_current: "docs-alicloud-resource-ess-notification"
description: |-
  Provides a ESS notification resource.
---

# alicloud\_ess\_notification

Provides a ESS notification resource. More about Ess notification, see [Autoscaling Notification](https://www.alibabacloud.com/help/doc-detail/71114.htm).

-> **NOTE:** Available in 1.55.0+

## Example Usage
```
variable "name" {
  default = "tf-testAccEssNotification-%d"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_mns_queue" "default" {
  name = var.name
}

resource "alicloud_ess_notification" "default" {
  scaling_group_id   = alicloud_ess_scaling_group.default.id
  notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS", "AUTOSCALING:SCALE_OUT_ERROR"]
  notification_arn   = "acs:ess:${data.alicloud_regions.default.regions[0].id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the Auto Scaling group.
* `notification_arn` - (Required, ForceNew) The Alibaba Cloud Resource Name (ARN) of the notification object, The value must be in `acs:ess:{region}:{account-id}:{resource-relative-id}` format.
    * region: the region ID of the scaling group. For more information, see `Regions and zones`
    * account-id: the ID of your account.
    * resource-relative-id: the notification method. Valid values : `cloudmonitor`, MNS queue: `queue/{queuename}`, Replace the queuename with the specific MNS queue name, MNS topic: `topic/{topicname}`, Replace the topicname with the specific MNS topic name.
* `notification_types` - (Required) The notification types of Auto Scaling events and resource changes. Supported notification types: 'AUTOSCALING:SCALE_OUT_SUCCESS', 'AUTOSCALING:SCALE_IN_SUCCESS', 'AUTOSCALING:SCALE_OUT_ERROR', 'AUTOSCALING:SCALE_IN_ERROR', 'AUTOSCALING:SCALE_REJECT', 'AUTOSCALING:SCALE_OUT_START', 'AUTOSCALING:SCALE_IN_START', 'AUTOSCALING:SCHEDULE_TASK_EXPIRING'.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of notification resource, which is composed of 'scaling_group_id' and 'notification_arn' in the format of '<scaling_group_id>:<notification_arn>'.

## Import

Ess notification can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_notification.example 'scaling_group_id:notification_arn'
```

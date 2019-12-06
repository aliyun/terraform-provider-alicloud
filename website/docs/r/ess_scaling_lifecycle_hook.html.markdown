---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_lifecycle_hook"
sidebar_current: "docs-alicloud-resource-ess-lifecycle-hook"
description: |-
  Provides a ESS lifecycle hook resource.
---

# alicloud\_ess\_lifecycle\_hook

Provides a ESS lifecycle hook resource. More about Ess lifecycle hook, see [LifecycleHook](https://www.alibabacloud.com/help/doc-detail/73839.htm).

## Example Usage
```
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name       = "testAccEssScalingGroup_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = "${alicloud_vpc.foo.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  vpc_id            = "${alicloud_vpc.foo.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ess_scaling_group" "foo" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "testAccEssScaling_group"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = ["${alicloud_vswitch.foo.id}", "${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_lifecycle_hook" "foo" {
  scaling_group_id      = "${alicloud_ess_scaling_group.foo.id}"
  name                  = "testAccEssLifecycle_hook"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}
```

## Module Support

You can use to the existing [autoscaling module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling/alicloud) 
to create a lifecycle hook, scaling group and configuration directly.

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the Auto Scaling group to which you want to assign the lifecycle hook.
* `name` - (Optional, ForceNew) The name of the lifecycle hook, which must contain 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain number, underscores `_`, hypens `-`, and decimal point `.`. If this parameter value is not specified, the default value is lifecycle hook id.
* `lifecycle_transition` - (Required) Type of Scaling activity attached to lifecycle hook. Supported value: SCALE_OUT, SCALE_IN.
* `heartbeat_timeout` - (Optional) Defines the amount of time, in seconds, that can elapse before the lifecycle hook times out. When the lifecycle hook times out, Auto Scaling performs the action defined in the default_result parameter. Default value: 600.
* `default_result` - (Optional) Defines the action the Auto Scaling group should take when the lifecycle hook timeout elapses. Applicable value: CONTINUE, ABANDON, default value: CONTINUE.
* `notification_arn` - (Optional) The Arn of notification target.
* `notification_metadata` - (Optional) Additional information that you want to include when Auto Scaling sends a message to the notification target.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of lifecycle hook.
* `scaling_group_id` - The scalingGroupId to which lifecycle belongs.
* `name` - The name of lifecycle hook.
* `default_result` - The action the Auto Scaling group should take when the lifecycle hook timeout elapses.
* `heartbeat_timeout` - The amount of time that can elapse before the lifecycle hook time out.
* `lifecycle_transition` - Type of Scaling activity attached to lifecycle hook.
* `notification_metadata` - Additional information that will be sent to notification target.
* `notification_arn` - The arn of notification target.

## Import

Ess lifecycle hook can be imported using the id, e.g.

```
$ terraform import alicloud_ess_lifecycle_hook.example ash-l12345
```

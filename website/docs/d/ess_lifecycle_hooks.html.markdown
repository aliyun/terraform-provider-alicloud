---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_lifecycle_hooks"
sidebar_current: "docs-alicloud_ess_lifecycle_hooks"
description: |-
    Provides a list of lifecycle hooks available to the user.
---

# alicloud_ess_lifecycle_hooks

This data source provides available lifecycle hook resources. 

-> **NOTE:** Available in 1.72.0+

## Example Usage

```
data "alicloud_ess_lifecycle_hooks" "ds" {
  scaling_group_id = "scaling_group_id"
  name_regex       = "lifecyclehook_name"
}

output "first_lifecycle_hook" {
  value = "${data.alicloud_ess_lifecycle_hooks.ds.hooks.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) Scaling group id the lifecycle hooks belong to.
* `name_regex` - (Optional) A regex string to filter resulting lifecycle hook by name.
* `ids` - (Optional) A list of lifecycle hook IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of lifecycle hook ids.
* `names` - A list of lifecycle hook names.
* `hooks` - A list of lifecycle hooks. Each element contains the following attributes:
  * `id` - ID of the lifecycle hook.
  * `scaling_group_id` - ID of the scaling group.
  * `name` - Name of the lifecycle hook.
  * `default_result` - Defines the action the Auto Scaling group should take when the lifecycle hook timeout elapses. 
  * `heartbeat_timeout` - Defines the amount of time, in seconds, that can elapse before the lifecycle hook times out. When the lifecycle hook times out, Auto Scaling performs the action defined in the default_result parameter.
  * `lifecycle_transition` - Type of Scaling activity attached to lifecycle hook.
  * `notification_arn` - The Arn of notification target.
  * `notification_metadata` - Additional information that you want to include when Auto Scaling sends a message to the notification target.

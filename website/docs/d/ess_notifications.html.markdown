---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_notifications"
sidebar_current: "docs-alicloud_ess_notifications"
description: |-
    Provides a list of notifications available to the user.
---

# alicloud_ess_notifications

This data source provides available notification resources. 

-> **NOTE:** Available in 1.72.0+

## Example Usage

```
data "alicloud_ess_notifications" "ds" {
  scaling_group_id = "scaling_group_id"
}

output "first_notification" {
  value = "${data.alicloud_ess_notifications.ds.notifications.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) Scaling group id the notifications belong to.
* `ids` - (Optional)A list of notification ids.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of notification ids.
* `notifications` - A list of notifications. Each element contains the following attributes:
  * `id` - ID of the notification.
  * `scaling_group_id` - ID of the scaling group.
  * `notification_arn` - The Alibaba Cloud Resource Name (ARN) for the notification object. 
  * `notification_types` - The notification types of Auto Scaling events and resource changes.

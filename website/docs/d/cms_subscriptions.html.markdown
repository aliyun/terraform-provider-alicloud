---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_subscriptions"
sidebar_current: "docs-alicloud-datasource-cms-subscriptions"
description: |-
  Provides a list of Cloud Monitor Service Subscriptions to the user.
---

# alicloud_cms_subscriptions

This data source provides the Cloud Monitor Service Subscriptions of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_subscriptions" "default" {
  subscription_name = "tf_example"
  ids              = ["<id>"]
}

output "first_subscription_id" {
  value = data.alicloud_cms_subscriptions.default.subscriptions.0.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of subscription IDs.
* `subscription_name` - (Optional, ForceNew) The name of the subscription.
* `enable` - (Optional, ForceNew) Whether the subscription is enabled.
* `workspace` - (Optional, ForceNew) The workspace of the subscription.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by subscription name.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of subscription IDs.
* `subscriptions` - A list of subscriptions. Each element contains the following attributes:
  * `id` - The ID of the subscription.
  * `subscription_id` - The ID of the subscription.
  * `subscription_name` - The name of the subscription.
  * `description` - The description of the subscription.
  * `notify_strategy_id` - The ID of the notify strategy.
  * `workspace` - The workspace of the subscription.
  * `user_id` - The user ID.
  * `create_time` - The creation time of the subscription.
  * `update_time` - The update time of the subscription.
  * `enable` - Whether the subscription is enabled.
  * `region_id` - The region ID of the subscription.
  * `subscription_type` - The type of the subscription.
  * `filter_setting` - The filter setting of the subscription.
  * `pushing_setting` - The pushing setting of the subscription.

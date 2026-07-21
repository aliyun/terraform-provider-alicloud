---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_subscriptions"
sidebar_current: "docs-alicloud-datasource-message-service-subscriptions"
description: |-
  Provides a list of Message Notification Service Subscriptions to the user.
---

# alicloud\_message\_service\_subscriptions

This data source provides the Message Notification Service Subscriptions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.188.0+.

## Example Usage

```terraform
data "alicloud_message_service_subscriptions" "ids" {
  ids        = ["example_id"]
  topic_name = "tf-example"
}

output "subscription_id_1" {
  value = data.alicloud_message_service_subscriptions.ids.subscriptions.0.id
}

data "alicloud_message_service_subscriptions" "name" {
  topic_name = "tf-example"
}

output "subscription_id_2" {
  value = data.alicloud_message_service_subscriptions.name.subscriptions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Subscription IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Subscription name.
* `topic_name` - (Required, ForceNew) The name of the topic.
* `subscription_name` - (Optional, ForceNew) The name of the subscription.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Subscription names.
* `subscriptions` - A list of Subscriptions. Each element contains the following attributes:
  * `id` - The id of the Subscription.
  * `topic_name` - The name of the topic.
  * `subscription_name` - The name of the subscription.
  * `endpoint` - The endpoint to which the messages are pushed.
  * `filter_tag` - The tag that is used to filter messages. Only the messages that are attached with the specified tag can be pushed.
  * `notify_content_format` - The content format of the messages that are pushed to the endpoint.
  * `notify_strategy` - The retry policy that is applied if an error occurs when MNS pushes messages to the endpoint.
  * `topic_owner` - The account ID of the topic owner.
  * `subscription_url` - The url of the subscription.
  * `last_modify_time` - The time when the subscription was last modified. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.
  * `create_time` - The time when the subscription was created. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.

---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic_subscriptions"
sidebar_current: "docs-alicloud-datasource-mns-topic-subscriptions"
description: |-
  Provides a list of mns topic subscriptions available to the user.
---

# alicloud\_mns\_topic_subscriptions

This data source provides a list of MNS topic subscriptions in an Alibaba Cloud account according to the specified parameters.

-> **DEPRECATED:**  This datasource has been deprecated from version `1.188.0`. Please use new datasource [message_service_subscriptions](https://www.terraform.io/docs/providers/alicloud/d/message_service_subscriptions).

## Example Usage

```terraform
data "alicloud_mns_topic_subscriptions" "subscriptions" {
  topic_name  = "topic_name"
  name_prefix = "tf-"
}

output "first_topic_subscription_id" {
  value = "${data.alicloud_mns_topic_subscriptions.subscriptions.subscriptions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required) Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `name_prefix` - (Optional) A string to filter resulting subscriptions of the topic by their name prefixs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of subscription names.
* `subscriptions` - A list of subscriptions. Each element contains the following attributes:
  * `id` - The ID of the topic subscription. The value is set to `name`.
  * `name` - The name of the subscription.
  * `topic_name`- The topic which The subscription belongs to was named with the name.
  * `notify_strategy` - The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails.
  * `notify_content_format` - The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users.
  * `endpoint` - Describe the terminal address of the message received in this subscription.
  * `filter_tag`- A string to filter resulting messages of the topic by their message tag.

---
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic_subscriptions"
sidebar_current: "docs-alicloud-datasource-mns-topic-subscriptions"
description: |-
    Provides a list of mns topic subscriptions available to the user.
---

# alicloud\_msn\_topic_subscriptions

This data source provides a list of MNS topic subscriptions  in an Alibaba Cloud account according to the specified parameters.

## Example Usage

```
data "alicloud_mns_topic_subscriptions" "subscriptions" {
  topic_name="topic_name"
  name_prefix = "tf-"
}

output "first_topic_subscription_id" {
  value = "${data.alicloud_mns_topic_subscriptions.subscriptions.subscriptions.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `topic_name`  -  (Required) Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `name_prefix` - (Optional) A  string to filter resulting subscriptions of the topic by their name prefixs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `subscriptions` - A list of users. Each element contains the following attributes:
   * `name` - Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
   * `notify_strategy` - The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. the attribute has two value EXPONENTIAL_DECAY_RETR or BACKOFF_RETRY.
   * `notify_content_format` - he NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. the attribute has two value SIMPLIFIED or XML.
   * `endpoint` -  "email format: mail:directmail:XXX@YYY.com ,   queue format: http(s)://AccountId.mns.regionId.aliyuncs.com/, http format: http(s)://www.xxx.com/xx"
   * `filter_tag`-  A  string to filter resulting messages of the topic by their message tag.

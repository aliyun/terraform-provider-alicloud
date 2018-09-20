--
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic_subscription"
sidebar_current: "docs-alicloud-resource-mns-topic_subscription"
description: |-
  Provides a Alicloud MNS Topic  Subscription resource.
---

# alicloud\_mns\_topic\_subscription

Provides a MNS topic subscription resource.

~> **NOTE:** Terraform will auto build a mns topic subscription  while it uses `alicloud_mns_topic_subscription` to build a mns topic subscription resource.

## Example Usage

Basic Usage

```
resource "alicloud_mns_topic" "topic"{
    name="${var.name}"
    maximum_message_size=${var.maximum_message_size}
    logging_enabled=${var.logging_enabled} 
}

resource "alicloud_mns_topic_subscription" "subscription"{
    topic_name="${alicloud_mns_topic.topic.name}"
    name="${var.subscription_name}"
    filter_tag="${var.fitler_flag}"
    endpoint="${var.endpoint}"
    notify_strategy="${var.notify_strategy}"
    notify_content_format="${var.notify_content_format}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `notify_strategy` - (Required)The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. the attribute has two value EXPONENTIAL_DECAY_RETR or BACKOFF_RETRY.
* `notify_content_format` - (Required, Forces new resource)he NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. the attribute has two value SIMPLIFIED or XML.
* `endpoint` - (Required, Forces new resource) "email format: mail:directmail:XXX@YYY.com ,   queue format: http(s)://AccountId.mns.regionId.aliyuncs.com/, http format: http(s)://www.xxx.com/xx"
* `filter_tag`-(Optional, Forces new resource)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the topic subscription.
* `name` - The name of the topic subscription.

## Import

the import will be supported in the next version e.g.
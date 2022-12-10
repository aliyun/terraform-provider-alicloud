---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_subscription"
sidebar_current: "docs-alicloud-resource-message-service-subscription"
description: |-
  Provides a Alicloud Message Notification Service Subscription resource.
---

# alicloud\_message\_service\_subscription

Provides a Message Notification Service Subscription resource.

For information about Message Notification Service Subscription and how to use it, see [What is Subscription](https://www.alibabacloud.com/help/en/message-service/latest/subscribe-1).

-> **NOTE:** Available in v1.188.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_message_service_topic" "default" {
  topic_name       = "tf-example-value"
  max_message_size = 12357
  logging_enabled  = true
}

resource "alicloud_message_service_subscription" "default" {
  topic_name            = alicloud_message_service_topic.default.topic_name
  subscription_name     = "tf-example-value"
  endpoint              = "http://www.test.com/test"
  push_type             = "http"
  filter_tag            = "tf-test"
  notify_content_format = "XML"
  notify_strategy       = "BACKOFF_RETRY"
}
```

## Argument Reference

The following arguments are supported:

* `topic_name`- (Required, ForceNew) The topic which The subscription belongs to was named with the name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `subscription_name` - (Required, ForceNew) Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `endpoint` - (Required, ForceNew) The endpoint has three format. Available values format:
  - `HTTP Format`: http://xxx.com/xxx
  - `Queue Format`: acs:mns:{REGION}:{AccountID}:queues/{QueueName}
  - `Email Format`: mail:directmail:{MailAddress}
* `push_type` - (Required, ForceNew) The Push type of Subscription. The Valid values: `http`, `queue`, `mpush`, `alisms` and `email`.
* `filter_tag` - (Optional, ForceNew) The tag that is used to filter messages. Only the messages that have the same tag can be pushed. A tag is a string that can be up to 16 characters in length. By default, no tag is specified to filter messages.
* `notify_content_format` - (Optional, Computed, ForceNew) The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. Valid values: `XML`, `JSON` and `SIMPLIFIED`. Default value: `XML`.
* `notify_strategy` - (Optional, Computed) The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. Default value: `BACKOFF_RETRY`. Valid values:
  - `BACKOFF_RETRY`: retries with a fixed backoff interval.
  - `EXPONENTIAL_DECAY_RETRY`: retries with exponential backoff.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Subscription. The value formats as `<topic_name>:<subscription_name>`.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Subscription.
* `update` - (Defaults to 3 mins) Used when update the Subscription.
* `delete` - (Defaults to 3 mins) Used when delete the Subscription.

## Import

Message Notification Service Subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_subscription.example <topic_name>:<subscription_name>
```

---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_subscription"
description: |-
  Provides a Alicloud Message Service Subscription resource.
---

# alicloud_message_service_subscription

Provides a Message Service Subscription resource. 

For information about Message Service Subscription and how to use it, see [What is Subscription](https://www.alibabacloud.com/help/en/message-service/latest/subscribe-1).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_message_service_topic" "defaultTopic" {
  max_message_size = 65536
  topic_name       = var.name
}


resource "alicloud_message_service_subscription" "default" {
  push_type             = "http"
  endpoint              = "http://example.com"
  notify_strategy       = "BACKOFF_RETRY"
  notify_content_format = "SIMPLIFIED"
  subscription_name     = var.name
  filter_tag            = "important"
  topic_name            = alicloud_message_service_topic.defaultTopic.topic_name
}
```

## Argument Reference

The following arguments are supported:
* `endpoint` - (Required, ForceNew) The endpoint has three format. Available values format:
  - `HTTP Format`: http://example.com/example
  - `Queue Format`: acs:mns:{REGION}:{AccountID}:queues/{QueueName}
  - `Email Format`: mail:directmail:{MailAddress}.
* `filter_tag` - (Optional, ForceNew) The tag that is used to filter messages. Only the messages that have the same tag can be pushed. A tag is a string that can be up to 16 characters in length. By default, no tag is specified to filter messages.
* `notify_content_format` - (Optional, ForceNew) The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. Valid values: `XML`, `JSON` and `SIMPLIFIED`. Default value: `XML`.
* `notify_strategy` - (Optional) The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. Default value: `BACKOFF_RETRY`. Valid values:
  - `BACKOFF_RETRY`: retries with a fixed backoff interval.
  - `EXPONENTIAL_DECAY_RETRY`: retries with exponential backoff.
* `push_type` - (Required) The Push type of Subscription. The Valid values: `http`, `queue`, `mpush`, `alisms` and `email`.
* `subscription_name` - (Required, ForceNew) Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `topic_name` - (Required, ForceNew) The topic which The subscription belongs to was named with the name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<topic_name>:<subscription_name>`.
* `create_time` - Represents when the subscription was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Subscription.
* `delete` - (Defaults to 5 mins) Used when delete the Subscription.
* `update` - (Defaults to 5 mins) Used when update the Subscription.

## Import

Message Service Subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_subscription.example <topic_name>:<subscription_name>
```
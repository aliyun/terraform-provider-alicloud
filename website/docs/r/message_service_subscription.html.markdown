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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_subscription&exampleId=0c87558b-ba9d-70e9-726a-e079c2e192a76e7dbaf5&activeTab=example&spm=docs.r.message_service_subscription.0.0c87558bba&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_message_service_topic" "default" {
  topic_name       = var.name
  max_message_size = 16888
  enable_logging   = true
}

resource "alicloud_message_service_subscription" "default" {
  topic_name            = alicloud_message_service_topic.default.topic_name
  subscription_name     = var.name
  endpoint              = "http://example.com"
  push_type             = "http"
  filter_tag            = var.name
  notify_content_format = "XML"
  notify_strategy       = "BACKOFF_RETRY"
}
```

## Argument Reference

The following arguments are supported:
* `dlq_policy` - (Optional, Set, Available since v1.244.0) The dead-letter queue policy. See [`dlq_policy`](#dlq_policy) below.
* `topic_name`- (Required, ForceNew) The topic which The subscription belongs to was named with the name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `subscription_name` - (Required, ForceNew) Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `endpoint` - (Required, ForceNew) The endpoint has three format. Available values format:
  - `HTTP Format`: http://xxx.com/xxx
  - `Queue Format`: acs:mns:{REGION}:{AccountID}:queues/{QueueName}
  - `Email Format`: mail:directmail:{MailAddress}
* `push_type` - (Required, ForceNew) The Push type of Subscription. The Valid values: `http`, `queue`, `mpush`, `alisms` and `email`.
* `filter_tag` - (Optional, ForceNew) The tag that is used to filter messages. Only the messages that have the same tag can be pushed. A tag is a string that can be up to 16 characters in length. By default, no tag is specified to filter messages.
* `notify_content_format` - (Optional, Computed, ForceNew) The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. Valid values: `XML`, `JSON` and `SIMPLIFIED`. Default value: `XML`.
* `notify_strategy` - (Optional) The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. Default value: `BACKOFF_RETRY`. Valid values:
  - `BACKOFF_RETRY`: retries with a fixed backoff interval.
  - `EXPONENTIAL_DECAY_RETRY`: retries with exponential backoff.

### `dlq_policy`

The dlq_policy supports the following:
* `dead_letter_target_queue` - (Optional) The queue to which dead-letter messages are delivered.
* `enabled` - (Optional, Bool) Specifies whether to enable the dead-letter message delivery. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Subscription. The value formats as `<topic_name>:<subscription_name>`.
* `create_time` - (Available since v1.244.0) The time when the subscription was created.

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

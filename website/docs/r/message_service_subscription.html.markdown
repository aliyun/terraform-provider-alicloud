---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_subscription"
description: |-
  Provides a Alicloud Message Service Subscription resource.
---

# alicloud_message_service_subscription

Provides a Message Service Subscription resource.



For information about Message Service Subscription and how to use it, see [What is Subscription](https://www.alibabacloud.com/help/en/mns/developer-reference/api-mns-open-2022-01-19-subscribe).

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_message_service_subscription&spm=docs.r.message_service_subscription.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `dlq_policy` - (Optional, Set, Available since v1.244.0) The dead-letter queue policy. See [`dlq_policy`](#dlq_policy) below.
* `dm_attributes` - (Optional, ForceNew, Set, Available since v1.284.0) The email push attributes. This parameter is required when `push_type` is set to `dm`. See [`dm_attributes`](#dm_attributes) below.
* `dysms_attributes` - (Optional, ForceNew, Set, Available since v1.284.0) The SMS push attributes. This parameter is required when `push_type` is set to `dysms`. See [`dysms_attributes`](#dysms_attributes) below.
* `tenant_rate_limit_policy` - (Optional, Set, Available since v1.284.0) The rate limit policy. See [`tenant_rate_limit_policy`](#tenant_rate_limit_policy) below.
* `topic_name`- (Required, ForceNew) The topic which The subscription belongs to was named with the name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `subscription_name` - (Required, ForceNew) Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 255 characters.
* `endpoint` - (Required, ForceNew) The endpoint of the subscription. The format varies with `push_type`. Available values format:
  - `HTTP Format`: An HTTP URL that starts with http:// or https://.
  - `Queue Format`: A queue name.
  - `Dm Format`: `smq-ep:dm:<account_id>:__dynamic`, where `<account_id>` is the ID of your Alibaba Cloud account.
  - `Dysms Format`: `smq-ep:dysms:<account_id>:<phone_number>`.
  - `MPush Format`: An AppKey.
  - `Sms Format`: A mobile number
  - `Email Format`: An email address.
* `sts_role_arn` - (Optional, ForceNew, Available since v1.259.0) The ARN of the RAM role assumed by the service. The format is `acs:ram::<account_id>:role/<role_name>`. This parameter is required when `push_type` is set to `dm` or `dysms`.
* `push_type` - (Required) The Push type of Subscription. Valid values: `http`, `queue`, `dm`, `dysms`, `fc` and `eventbus`. The values `mpush`, `alisms` and `email` are deprecated and are retained only for compatibility with existing subscriptions.
* `filter_tag` - (Optional, ForceNew) The tag that is used to filter messages. Only the messages that have the same tag can be pushed. A tag is a string that can be up to 16 characters in length. By default, no tag is specified to filter messages.
* `notify_content_format` - (Optional, Computed, ForceNew) The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. Valid values: `XML`, `JSON` and `SIMPLIFIED`. Default value: `XML`.
* `notify_strategy` - (Optional) The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. Default value: `BACKOFF_RETRY`. Valid values:
  - `BACKOFF_RETRY`: retries with a fixed backoff interval.
  - `EXPONENTIAL_DECAY_RETRY`: retries with exponential backoff.

### `dlq_policy`

The dlq_policy supports the following:
* `dead_letter_target_queue` - (Optional) The queue to which dead-letter messages are delivered.
* `enabled` - (Optional, Bool) Specifies whether to enable the dead-letter message delivery. Valid values: `true`, `false`.

### `dm_attributes`

The dm_attributes supports the following:
* `account_name` - (Optional, ForceNew) The sender address of the email push. The value must be a sender address that has been configured in Direct Mail.
* `subject` - (Optional, ForceNew) The subject of the pushed email.

### `dysms_attributes`

The dysms_attributes supports the following:
* `sign_name` - (Optional, ForceNew) The signature name of the SMS push. The value must be a signature that has been configured in Short Message Service.
* `template_code` - (Optional, ForceNew) The template code of the SMS push. You can obtain the value from the Short Message Service console.

### `tenant_rate_limit_policy`

The tenant_rate_limit_policy supports the following:
* `enabled` - (Optional, Bool) Specifies whether to enable the rate limit policy. Valid values: `true`, `false`.
* `max_receives_per_second` - (Optional, Int) The maximum number of messages that can be pushed or consumed per second.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Subscription. The value formats as `<topic_name>:<subscription_name>`.
* `create_time` - (Available since v1.244.0) The time when the subscription was created.
* `last_modify_time` - (Available since v1.284.0) The time when the subscription was last modified.
* `topic_owner` - (Available since v1.284.0) The ID of the Alibaba Cloud account that owns the subscribed topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Subscription.
* `delete` - (Defaults to 5 mins) Used when delete the Subscription.
* `update` - (Defaults to 5 mins) Used when update the Subscription.

## Import

Message Service Subscription can be imported using the id, which consists of topic_name and subscription_name, e.g.

```shell
$ terraform import alicloud_message_service_subscription.example <topic_name>:<subscription_name>
```

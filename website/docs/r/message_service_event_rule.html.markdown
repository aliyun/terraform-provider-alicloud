---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_event_rule"
description: |-
  Provides a Alicloud Message Service Event Rule resource.
---

# alicloud_message_service_event_rule

Provides a Message Service Event Rule resource.



For information about Message Service Event Rule and how to use it, see [What is Event Rule](https://next.api.alibabacloud.com/document/Mns-open/2022-01-19/CreateEventRule).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_event_rule&exampleId=cc3e6632-2b29-7c57-b8d9-af824dba9a789f446f0f&activeTab=example&spm=docs.r.message_service_event_rule.0.cc3e66322b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "queue_name" {
  default = "tf-exampe-topic2queue"
}

variable "rule_name" {
  default = "tf-exampe-topic-1"
}

variable "topic_name" {
  default = "tf-exampe-topic2queue"
}

resource "alicloud_message_service_topic" "CreateTopic" {
  max_message_size = "65536"
  topic_name       = var.topic_name
  logging_enabled  = false
}

resource "alicloud_message_service_queue" "CreateQueue" {
  delay_seconds            = "2"
  polling_wait_seconds     = "2"
  message_retention_period = "566"
  maximum_message_size     = "1123"
  visibility_timeout       = "30"
  queue_name               = var.queue_name
  logging_enabled          = false
}

resource "alicloud_message_service_subscription" "CreateSub" {
  push_type             = "queue"
  notify_strategy       = "BACKOFF_RETRY"
  notify_content_format = "SIMPLIFIED"
  subscription_name     = "RDK-example-sub"
  filter_tag            = "important"
  topic_name            = alicloud_message_service_topic.CreateTopic.topic_name
  endpoint              = format("acs:mns:cn-hangzhou:1511928242963727:/queues/%s", alicloud_message_service_queue.CreateQueue.id)
}

resource "alicloud_message_service_event_rule" "default" {
  event_types = [
    "ObjectCreated:PutObject"
  ]
  match_rules = [
    [
      {
        suffix      = ""
        match_state = "true"
        name        = "acs:oss:cn-hangzhou:1511928242963727:accccx"
        prefix      = ""
      }
    ]
  ]
  endpoint {
    endpoint_value = alicloud_message_service_subscription.CreateSub.topic_name
    endpoint_type  = "topic"
  }

  rule_name = var.rule_name
}
```

## Argument Reference

The following arguments are supported:
* `delivery_mode` - (Optional, ForceNew) -DIRECT: directly delivers to a single queue (1:1) without creating a Topic;
  - BROADCAST: BROADCAST to all subscription queues (1:N). You need to create a Topic;
* `endpoint` - (Optional, ForceNew, List) Message Receiving Terminal Endpoint Object. See [`endpoint`](#endpoint) below.
* `event_types` - (Required, ForceNew, List) Event Type List
* `match_rules` - (Optional, ForceNew, List) Matching rules, or relationships between multiple rules. See [`match_rules`](#match_rules) below.
* `rule_name` - (Required, ForceNew) The event notification rule name.

### `endpoint`

The endpoint supports the following:
* `endpoint_type` - (Optional, ForceNew) Message receiving terminal endpoint type
* `endpoint_value` - (Optional, ForceNew) Message Receiving Terminal Endpoint

### `match_rules`

Event Matching Rules-Atomic Objects.
-> **Note:** Full Name Matching Rule: If this item is filled in, other items cannot be filled in.
Prefix match and suffix match: either of these two items can be filled in. If both items are filled in, the front and suffix match.

* `match_state` - (Optional, ForceNew)  Match state. valid values: `true`, `false`.
* `match_name` - (Optional, ForceNew) Full name matching rule.
* `prefix` - (Optional, ForceNew) Prefix matching rule.
* `suffix` - (Optional, ForceNew) Suffix matching rule.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Event Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Event Rule.

## Import

Message Service Event Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_event_rule.example <id>
```
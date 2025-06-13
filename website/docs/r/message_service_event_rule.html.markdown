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

-> **NOTE:** Available since v1.252.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "queue_name" {
  default = "tf-example-topic2queue"
}

variable "rule_name" {
  default = "exampleRule-topic-1"
}

variable "topic_name" {
  default = "tf-example-topic2queue"
}

resource "alicloud_message_service_topic" "CreateTopic" {
  maximum_message_size = "65536"
  name                 = var.topic_name
  logging_enabled      = false
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
  topic_name            = alicloud_message_service_topic.CreateTopic.name
  endpoint              = alicloud_message_service_queue.CreateQueue.id
}


resource "alicloud_message_service_event_rule" "default" {
  endpoints {
    endpoint_type  = "topic"
    endpoint_value = alicloud_message_service_queue.CreateQueue.id
  }
  rule_name   = var.rule_name
  event_types = ["ObjectCreated:PutObject"]
  match_rules "items" {
    match_state = true
    name        = "aaac"
  }
}
```

## Argument Reference

The following arguments are supported:
* `endpoints` - (Required, ForceNew, List) Message receiving terminal list See [`endpoints`](#endpoints) below.
* `event_types` - (Required, ForceNew, List) Event Type List
* `match_rules` - (Required, ForceNew, List) Matching rules, or relationships between multiple rules
* `rule_name` - (Required, ForceNew) The event notification rule name.

### `endpoints`

The endpoints supports the following:
* `endpoint_type` - (Required, ForceNew) Receiving terminal type
* `endpoint_value` - (Required, ForceNew) Receiving terminal actual value

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
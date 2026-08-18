---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_streaming"
sidebar_current: "docs-alicloud-resource-event-bridge-event-streaming"
description: |-
  Provides a Alicloud Event Bridge Event Streaming resource.
---

# alicloud_event_bridge_event_streaming

Provides a Event Bridge Event Streaming resource.

An event streaming allows you to stream events from an event source to an event target, with optional event filtering and transformation.

For information about Event Bridge Event Streaming and how to use it, see [What is Event Streaming](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createeventstreaming).

-> **NOTE:** Available since v1.290.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_message_service_queue" "source" {
  queue_name = "${var.name}-source"
}

resource "alicloud_message_service_queue" "sink" {
  queue_name = "${var.name}-sink"
}

resource "alicloud_event_bridge_event_streaming" "default" {
  event_streaming_name = var.name
  description          = "terraform-example-event-streaming"
  filter_pattern       = "{}"
  source = jsonencode({
    SourceMNSParameters = {
      RegionId       = "cn-hangzhou"
      QueueName      = alicloud_message_service_queue.source.queue_name
      IsBase64Decode = true
    }
  })
  sink = jsonencode({
    SinkMNSParameters = {
      QueueName = {
        Value = alicloud_message_service_queue.sink.queue_name
        Form  = "CONSTANT"
      }
      Body = {
        Value = "$.data"
        Form  = "JSONPATH"
      }
      IsBase64Encode = {
        Value = "true"
        Form  = "CONSTANT"
      }
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `event_streaming_name` - (Required, ForceNew) The name of the event streaming.
* `source` - (Required, Json) The event provider, which is also known as the event source. You must and can specify only one event source. The value is a JSON string, for example `{"SourceMNSParameters":{"QueueName":"example","RegionId":"cn-hangzhou","IsBase64Decode":true}}`.
* `filter_pattern` - (Required, Json) The rule that is used to filter events. If you leave this parameter empty, all events are matched. The value is a JSON string, for example `{}`.
* `sink` - (Required, Json) The event target. You must and can specify only one event target. The value is a JSON string, for example `{"SinkMNSParameters":{"QueueName":{"Value":"example","Form":"CONSTANT"}}}`.
* `run_options` - (Optional, Json) The parameters that are configured for the runtime environment. The value is a JSON string, for example `{"ErrorsTolerance":"ALL","MaximumTasks":1}`.
* `transforms` - (Optional, Json) The rules that are used to transform events. The value is a JSON array string, for example `[{"Arn":"acs:fc:cn-hangzhou:ACCOUNT_ID:functions/example"}]`.
* `description` - (Optional) The description of the event streaming.
* `status` - (Optional) The expected status of the event streaming. Valid values: `RUNNING` and `PAUSED`. If you set this parameter to `RUNNING`, the event streaming is started. If you set this parameter to `PAUSED`, the event streaming is paused. If you do not set this parameter, the event streaming stays in the `READY` status after it is created.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Streaming. It formats as `<event_streaming_name>`.
* `status` - The current status of the event streaming. Valid values: `READY`, `STARTING`, `STARTING_FAILED`, `RUNNING`, `RUNNING_FAILED` and `PAUSED`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Event Streaming.
* `delete` - (Defaults to 10 mins) Used when delete the Event Streaming.
* `update` - (Defaults to 10 mins) Used when update the Event Streaming.

## Import

Event Bridge Event Streaming can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_streaming.example <event_streaming_name>
```

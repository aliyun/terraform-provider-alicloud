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

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_event_streaming&exampleId=74ffd385-6f9d-9bd2-11a6-9342a95ab1cd9bbbb5ef&activeTab=example&spm=docs.r.event_bridge_event_streaming.0.74ffd3856f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_event_bridge_event_streaming&spm=docs.r.event_bridge_event_streaming.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `event_streaming_name` - (Required, ForceNew) The name of the event streaming.
* `source` - (Required, String) The event provider, which is also known as the event source. You must and can specify only one event source. The value is a JSON string, for example `{"SourceMNSParameters":{"QueueName":"example","RegionId":"cn-hangzhou","IsBase64Decode":true}}`.
* `filter_pattern` - (Required, String) The rule that is used to filter events. If you leave this parameter empty, all events are matched. The value is a JSON string, for example `{}`.
* `sink` - (Required, String) The event target. You must and can specify only one event target. The value is a JSON string, for example `{"SinkMNSParameters":{"QueueName":{"Value":"example","Form":"CONSTANT"}}}`.
* `run_options` - (Optional, String) The parameters that are configured for the runtime environment. The value is a JSON string, for example `{"ErrorsTolerance":"ALL","MaximumTasks":1}`.
* `transforms` - (Optional, String) The rules that are used to transform events. The value is a JSON array string, for example `[{"Arn":"acs:fc:cn-hangzhou:ACCOUNT_ID:functions/example"}]`.
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

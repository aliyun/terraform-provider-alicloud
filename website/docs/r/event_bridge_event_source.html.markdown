---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_source"
sidebar_current: "docs-alicloud-resource-event-bridge-event-source"
description: |-
  Provides a Alicloud Event Bridge Event Source resource.
---

# alicloud_event_bridge_event_source

Provides a Event Bridge Event Source resource.

For information about Event Bridge Event Source and how to use it, see [What is Event Source](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createeventsource).

-> **NOTE:** Available since v1.130.0.

-> **NOTE:** Deprecated since v1.269.0.

-> **DEPRECATED:** This resource has been deprecated from version `1.269.0`. Please use new resource [alicloud_event_bridge_event_source_v2](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/event_bridge_event_source_v2).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_mns_queue" "default" {
  name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_event_bridge_event_bus" "default" {
  event_bus_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_event_bridge_event_source" "default" {
  event_bus_name         = alicloud_event_bridge_event_bus.default.event_bus_name
  event_source_name      = "${var.name}-${random_integer.default.result}"
  description            = var.name
  linked_external_source = true
  external_source_type   = "MNS"
  external_source_config = {
    QueueName = alicloud_mns_queue.default.name
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_event_bridge_event_source&spm=docs.r.event_bridge_event_source.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `event_bus_name` - (Required, ForceNew) The name of the event bus to which the event source is attached.
* `event_source_name` - (Required, ForceNew) The name of the event source.
* `external_source_type` - (Optional) The type of the external data source. Valid values: `RabbitMQ`, `RocketMQ` and `MNS`.
* `external_source_config`- (Optional, Map) The configuration of the external data source.
  When `external_source_type` is `RabbitMQ`, The following attributes are supported:
  `RegionId` - The region ID of RabbitMQ.
  `InstanceId` - The instance ID of RabbitMQ.
  `VirtualHostName` - The virtual host name of RabbitMQ.
  `QueueName` - The queue name of RabbitMQ.
  When `external_source_type` is `RabbitMQ`, The following attributes are supported:
  `RegionId` - The region ID of RabbitMQ.
  `InstanceId` - The instance ID of RabbitMQ.
  `Topic` - The topic of RabbitMQ.
  `Offset` -  The offset of RabbitMQ, valid values: `CONSUME_FROM_FIRST_OFFSET`, `CONSUME_FROM_LAST_OFFSET` and `CONSUME_FROM_TIMESTAMP`.
  `GroupID` - The group ID of consumer.
  When `external_source_type` is `MNS`, The following attributes are supported:
  `QueueName` - The queue name of MNS.
* `description` - (Optional) The description of the event source.
* `linked_external_source` - (Optional) Specifies whether to connect to an external data source. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Source. Value as `event_source_name`.

## Import

Event Bridge Event Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_source.example <event_source_name>
```

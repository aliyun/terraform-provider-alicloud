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

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_event_source&exampleId=cc6d99f0-df8d-c48f-df08-84324242fee34a3aabf6&activeTab=example&spm=docs.r.event_bridge_event_source.0.cc6d99f0df&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = var.name
}
resource "alicloud_mns_queue" "example" {
  name = "${var.name}-${random_integer.default.result}"
}
resource "alicloud_event_bridge_event_source" "example" {
  event_bus_name         = alicloud_event_bridge_event_bus.example.event_bus_name
  event_source_name      = var.name
  description            = var.name
  linked_external_source = true
  external_source_type   = "MNS"
  external_source_config = {
    QueueName = alicloud_mns_queue.example.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_name` - (Required, ForceNew) The name of event bus.
* `description` - (Optional) The detail describe of event source.
* `event_source_name` - (Required, ForceNew) The code name of event source.
* `linked_external_source` - (Optional, Computed) Whether to connect to an external data source. Default value: `false`
* `external_source_type` - (Optional) The type of external data source. Valid value : `RabbitMQ`, `RocketMQ` and `MNS`. **NOTE:** Only When `linked_external_source` is `true`, This field is valid.
* `external_source_config`- (Optional, Map) The config of external source.
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

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Source. Value as `event_source_name`.

## Import

Event Bridge Event Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_source.example <event_source_name>
```

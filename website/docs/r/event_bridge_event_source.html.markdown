---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_source"
sidebar_current: "docs-alicloud-resource-event-bridge-event-source"
description: |-
  Provides a Alicloud Event Bridge Event Source resource.
---

# alicloud\_event\_bridge\_event\_source

Provides a Event Bridge Event Source resource.

For information about Event Bridge Event Source and how to use it, see [What is Event Source](https://www.alibabacloud.com/help/doc-detail/188425.htm).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_event_source" "example" {
  event_bus_name         = "bus_name"
  event_source_name      = "tftest"
  description            = "tf-test"
  linked_external_source = true
  external_source_type   = "MNS"
  external_source_config = {
    QueueName = "mns_queuqe_name"
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

```
$ terraform import alicloud_event_bridge_event_source.example <event_source_name>
```

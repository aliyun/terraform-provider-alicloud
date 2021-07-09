---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_bus"
sidebar_current: "docs-alicloud-resource-event-bridge-event-bus"
description: |-
  Provides a Alicloud Event Bridge Event Bus resource.
---

# alicloud\_event\_bridge\_event\_bus

Provides a Event Bridge Event Bus resource.

For information about Event Bridge Event Bus and how to use it, see [What is Event Bus](https://www.alibabacloud.com/help/en/doc-detail/163897.htm).

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = "tf-testacc1234"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of event bus. 
* `event_bus_name` - (Required, ForceNew) The name of event bus. No more than 127 characters, starting with a letter or number, and individuals must include numbers, numbers, and dashes (-). The default is a reserved keyword and cannot be used as the name of the event. It cannot start with a string beginning with `eventbridge-reserved-`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Bus. Value as `event_bus_name`.

## Import

Event Bridge Event Bus can be imported using the id, e.g.

```
$ terraform import alicloud_event_bridge_event_bus.example <event_bus_name>
```
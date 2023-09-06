---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_bus"
sidebar_current: "docs-alicloud-resource-event-bridge-event-bus"
description: |-
  Provides a Alicloud Event Bridge Event Bus resource.
---

# alicloud_event_bridge_event_bus

Provides a Event Bridge Event Bus resource.

For information about Event Bridge Event Bus and how to use it, see [What is Event Bus](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createeventbus).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of event bus.
* `event_bus_name` - (Required, ForceNew) The name of event bus. The length is limited to 2 ~ 127 characters, which can be composed of letters, numbers or hyphens (-)

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Bus. Its value is same as `event_bus_name`.

## Import

Event Bridge Event Bus can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_bus.example <event_bus_name>
```

---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_buses"
sidebar_current: "docs-alicloud-datasource-event-bridge-event-buses"
description: |-
  Provides a list of Event Bridge Event Buses to the user.
---

# alicloud\_event\_bridge\_event\_buses

This data source provides the Event Bridge Event Buses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_event_bridge_event_buses" "example" {
  ids        = ["tf-testacc1234"]
  name_regex = "the_resource_name"
}

output "first_event_bridge_event_bus_id" {
  value = data.alicloud_event_bridge_event_buses.example.buses.0.id
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_type` - (Optional, ForceNew) The event bus type.
* `ids` - (Optional, ForceNew, Computed)  A list of Event Bus IDs.
* `name_prefix` - (Optional, ForceNew) The name prefix.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Event Bus name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Event Bus names.
* `buses` - A list of Event Bridge Event Buses. Each element contains the following attributes:
	* `description` - The description of event bus.
	* `event_bus_name` - The name of event bus.
	* `id` - The ID of the Event Bus.
	* `create_time` - The create time of Event Bus.
---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_sources"
sidebar_current: "docs-alicloud-datasource-event-bridge-event-sources"
description: |-
  Provides a list of Event Bridge Event Sources to the user.
---

# alicloud\_event\_bridge\_event\_sources

This data source provides the Event Bridge Event Sources of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_event_bridge_event_sources" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_event_bridge_event_source_id" {
  value = data.alicloud_event_bridge_event_sources.example.sources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Event Source IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Event Source name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Event Source names.
* `sources` - A list of Event Bridge Event Sources. Each element contains the following attributes:
	* `description` - The detail describe of event source.
	* `event_source_name` - The code name of event source.
	* `external_source_config` - The config of external data source.
	* `external_source_type` - The type of external data source.
	* `id` - The ID of the Event Source.
	* `linked_external_source` - Whether to connect to an external data source.
	* `event_bus_name` - The Event Source and Event Source Bound to the Event Bus Name.

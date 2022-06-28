---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_rules"
sidebar_current: "docs-alicloud-datasource-event-bridge-rules"
description: |-
  Provides a list of Event Bridge Rules to the user.
---

# alicloud\_event\_bridge\_rules

This data source provides the Event Bridge Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_event_bridge_rules" "example" {
  event_bus_name = "example_value"
  ids            = ["example_value"]
  name_regex     = "the_resource_name"
}

output "first_event_bridge_rule_id" {
  value = data.alicloud_event_bridge_rules.example.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_name` - (Required, ForceNew) The name of event bus.
* `ids` - (Optional, ForceNew, Computed)  A list of Rule IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `rule_name_prefix` - (Optional, ForceNew) The rule name prefix.
* `status` - (Optional, ForceNew) Rule status, either Enable or Disable. Valid values: `DISABLE`, `ENABLE`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of Event Bridge Rules. Each element contains the following attributes:
	* `description` - The description of rule.
	* `event_bus_name` - The name of event bus.
	* `filter_pattern` - The pattern to match interested events.
	* `id` - The ID of the Rule.
	* `rule_name` - The name of rule.
	* `status` - Rule status, either Enable or Disable.
	* `targets` - The target for rule.
		* `endpoint` - The endpoint.
		* `target_id` - The id of target.
		* `type` - The type of target.

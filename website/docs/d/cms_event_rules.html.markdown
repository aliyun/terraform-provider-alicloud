---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_event_rules"
sidebar_current: "docs-alicloud-datasource-cms-event-rules"
description: |-
  Provides a list of Cms Event Rules to the user.
---

# alicloud\_cms\_event\_rules

This data source provides the Cms Event Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.182.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_event_rules" "ids" {
  ids = ["example_id"]
}
output "cms_event_rule_id_1" {
  value = data.alicloud_cms_event_rules.ids.rules.0.id
}

data "alicloud_cms_event_rules" "nameRegex" {
  name_regex = "^my-EventRule"
}
output "cms_event_rule_id_2" {
  value = data.alicloud_cms_event_rules.nameRegex.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Event Rule IDs. Its element value is same as Event Rule Name.
* `name_prefix` - (Optional, ForceNew) The name prefix.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Event Rule name.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `DISABLED`, `ENABLED`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Event Rule names.
* `rules` - A list of Cms Event Rules. Each element contains the following attributes:
  	* `id` - The ID of the Event Rule. Its value is same as Event Rule Name.
  	* `event_rule_name` - The name of the event rule.	
	* `description` - The description of the rule.
	* `event_type` - The type of event.
	* `group_id` - The ID of the application Group.
	* `silence_time` - The mute period during which new alerts are not sent even if the trigger conditions are met.
	* `status` - The status of the resource.
	* `event_pattern` - Event mode, used to describe the trigger conditions for this event.
		* `event_type_list` - The list of event types.
		* `level_list` - The list of event levels.
		* `name_list` - The list of event names.
		* `product` - The type of the cloud service.
		* `sql_filter` - The SQL condition that is used to filter events.
		* `keyword_filter` - The filter keyword.
			* `key_words` - The keywords that are used to match events.
			* `relation` - The relationship between multiple keywords in a condition.
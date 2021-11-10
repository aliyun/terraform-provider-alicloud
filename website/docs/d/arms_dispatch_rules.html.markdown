---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_dispatch_rules"
sidebar_current: "docs-alicloud-datasource-arms-dispatch-rules"
description: |-
  Provides a list of Arms Dispatch Rules to the user.
---

# alicloud_arms_dispatch_rules

This data source provides the Arms Dispatch Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_dispatch_rules" "ids" {}
output "arms_dispatch_rule_id_1" {
  value = data.alicloud_arms_dispatch_rules.ids.rules.0.id
}

data "alicloud_arms_dispatch_rules" "nameRegex" {
  name_regex = "^my-DispatchRule"
}
output "arms_dispatch_rule_id_2" {
  value = data.alicloud_arms_dispatch_rules.nameRegex.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `dispatch_rule_name` - (Optional, ForceNew) The name of the dispatch rule.
* `ids` - (Optional, ForceNew, Computed)  A list of dispatch rule id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Dispatch Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dispatch Rule names.
* `rules` - A list of Arms Dispatch Rules. Each element contains the following attributes:
	* `dispatch_rule_id` - Dispatch rule ID.
	* `dispatch_rule_name` - The name of the dispatch rule.
	* `id` - The ID of the Dispatch Rule.
	* `status` - The resource status of Alert Dispatch Rule.
	* `group_rules` - Sets the event group.
		* `group_wait_time` - The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
		* `group_interval` - The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
		* `grouping_fields` - The fields that are used to group events. Events with the same field content are assigned to a group. Alerts with the same specified grouping field are sent to the handler in separate notifications.
		* `repeat_interval` - The silence period of repeated alerts. All alerts are repeatedly sent at specified intervals until the alerts are cleared. The minimum value is 61. Default to 600.

	* `label_match_expression_grid` - Sets the dispatch rule.
		* `label_match_expression_groups` - Sets the dispatch rule.
			* `label_match_expressions` - Sets the dispatch rule.
				* `key` - The key of the tag of the dispatch rule.
				* `value` - The value of the tag.
				* `operator` - The operator used in the dispatch rule. 
	
	* `notify_rules` - Sets the notification rule. 
		* `notify_objects` - Sets the notification object.
			* `notify_object_id` - The ID of the contact or contact group.
			* `name` - The name of the contact or contact group.
			* `notify_type` - The type of the alert contact.
			
  		* `notify_channels` - The notification method.

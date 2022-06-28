---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_groups"
sidebar_current: "docs-alicloud-datasource-cms-monitor-groups"
description: |-
  Provides a list of Cms Monitor Groups to the user.
---

# alicloud\_cms\_monitor\_groups

This data source provides the Cms Monitor Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_monitor_groups" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_cms_monitor_group_id" {
  value = data.alicloud_cms_monitor_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dynamic_tag_rule_id` - (Optional, ForceNew) The ID of the tag rule.
* `ids` - (Optional, ForceNew, Computed)  A list of Monitor Group IDs.
* `include_template_history` - (Optional, ForceNew) The include template history.
* `keyword` - (Optional, ForceNew) The keyword to be matched.
* `monitor_group_name` - (Optional, ForceNew) The name of the application group.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Monitor Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `select_contact_groups` - (Optional, ForceNew) The select contact groups.
* `type` - (Optional, ForceNew) The type of the application group. Valid values: `custom`, `ehpc_cluster`, `kubernetes`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Monitor Group names.
* `groups` - A list of Cms Monitor Groups. Each element contains the following attributes:
	* `bind_url` - The URL of the Kubernetes cluster from which the application group is synchronized.
	* `contact_groups` - The list of  alert groups that receive alert notifications for the application group.
	* `dynamic_tag_rule_id` - The ID of the tag rule.
	* `gmt_create` - The time when the application group was created.
	* `gmt_modified` - The time when the application group was modified.
	* `group_id` - The ID of the application group.
	* `id` - The ID of the Monitor Group.
	* `monitor_group_name` - The name of the application group.
	* `service_id` - The ID of the Alibaba Cloud service.
	* `tags` - A map of tags assigned to the Cms Monitor Group.
		* `tag_key` - The key of the tag attached to the application group.
		* `tag_value` - The value of the tag attached to the application group.
	* `template_ids` - The alert templates applied to the application group.
	* `type` - The type of the application group.

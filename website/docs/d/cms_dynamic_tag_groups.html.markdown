---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dynamic_tag_groups"
sidebar_current: "docs-alicloud-datasource-cms-dynamic-tag-groups"
description: |-
  Provides a list of Cms Dynamic Tag Groups to the user.
---

# alicloud\_cms\_dynamic\_tag\_groups

This data source provides the Cms Dynamic Tag Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}
resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
  describe                 = "example_value"
  enable_subscribed        = true
}
resource "alicloud_cms_dynamic_tag_group" "default" {
  contact_group_list = [alicloud_cms_alarm_contact_group.default.id]
  tag_key            = "your_tag_key"
  match_express {
    tag_value                = "your_tag_value"
    tag_value_match_function = "all"
  }
}
data "alicloud_cms_dynamic_tag_groups" "ids" {
  ids = [alicloud_cms_dynamic_tag_group.default.id]
}
output "cms_dynamic_tag_group_id_1" {
  value = data.alicloud_cms_dynamic_tag_groups.ids.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Dynamic Tag Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `RUNNING`, `FINISH`.
* `tag_key` - (Optional, ForceNew) The tag key of the tag.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of Cms Dynamic Tag Groups. Each element contains the following attributes:
	* `dynamic_tag_rule_id` - The ID of the tag rule.
	* `id` - The ID of the Dynamic Tag Group.
	* `match_express_filter_relation` - The relationship between conditional expressions. Valid values: `and`, `or`.
	* `match_express` - The label generates a matching expression that applies the grouping. See the following `Block match_express`.
		* `tag_value` - The tag value. The Tag value must be used in conjunction with the tag value matching method TagValueMatchFunction.
		* `tag_value_match_function` - Matching method of tag value. Valid values: `all`, `startWith`,`endWith`,`contains`,`notContains`,`equals`.
	* `status` -  The status of the resource. Valid values: `RUNNING`, `FINISH`.
	* `tag_key` - The tag key of the tag.
	

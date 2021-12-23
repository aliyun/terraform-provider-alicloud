---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dynamic_tag_group"
sidebar_current: "docs-alicloud-resource-cms-dynamic-tag-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Dynamic Tag Group resource.
---

# alicloud\_cms\_dynamic\_tag\_group

Provides a Cloud Monitor Service Dynamic Tag Group resource.

For information about Cloud Monitor Service Dynamic Tag Group and how to use it, see [What is Dynamic Tag Group](https://www.alibabacloud.com/help/doc-detail/150123.html).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = "example_value"
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

```

## Argument Reference

The following arguments are supported:

* `contact_group_list` - (Required, ForceNew) Alarm contact group. The value range of N is 1~100. The alarm notification of the application group is sent to the alarm contact in the alarm contact group.
* `enable_install_agent` - (Optional) The enable install agent.
* `enable_subscribe_event` - (Optional) The enable subscribe event.
* `match_express` - (Optional, ForceNew) The label generates a matching expression that applies the grouping. See the following `Block match_express`.
* `match_express_filter_relation` - (Optional, ForceNew) The relationship between conditional expressions. Valid values: `and`, `or`.
* `tag_key` - (Required, ForceNew) The tag key of the tag.
* `template_id_list` - (Optional, ForceNew) Alarm template ID list.

#### Block match_express

The match_express supports the following: 

* `tag_value` - (Optional) The tag value. The Tag value must be used in conjunction with the tag value matching method TagValueMatchFunction.
* `tag_value_match_function` - (Optional) Matching method of tag value. Valid values: `all`, `startWith`,`endWith`,`contains`,`notContains`,`equals`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dynamic Tag Group.
* `status` - The status of the resource. Valid values: `RUNNING`, `FINISH`.

## Import

Cloud Monitor Service Dynamic Tag Group can be imported using the id, e.g.

```
$ terraform import alicloud_cms_dynamic_tag_group.example <id>
```

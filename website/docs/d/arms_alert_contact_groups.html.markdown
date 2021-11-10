---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_alert_contact_groups"
sidebar_current: "docs-alicloud-datasource-arms-alert-contact-groups"
description: |-
  Provides a list of Arms Alert Contact Groups to the user.
---

# alicloud\_arms\_alert\_contact\_groups

This data source provides the Arms Alert Contact Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_arms_alert_contact_groups" "nameRegex" {
  name_regex = "^my-AlertContactGroup"
}
output "arms_alert_contact_group_id" {
  value = data.alicloud_arms_alert_contact_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `alert_contact_group_name` - (Optional, ForceNew) The name of the resource.
* `contact_id` - (Optional, ForceNew) The contact id.
* `contact_name` - (Optional, ForceNew) The contact name.
* `ids` - (Optional, ForceNew, Computed)  A list of Alert Contact Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Alert Contact Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Alert Contact Group names.
* `groups` - A list of Arms Alert Contact Groups. Each element contains the following attributes:
	* `alert_contact_group_id` - The first ID of the resource.
	* `alert_contact_group_name` - The name of the resource.
	* `contact_ids` - contact ids.
	* `create_time` - The creation time of the resource.
	* `id` - The ID of the Alert Contact Group.

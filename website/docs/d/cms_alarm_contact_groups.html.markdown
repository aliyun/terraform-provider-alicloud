---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm_contact_groups"
sidebar_current: "docs-alicloud-datasource-cms-alarm-contact-groups"
description: |-
  Provides a list of CMS Groups to the user.
---

# alicloud\_cms\_contact\_groups

This data source provides the CMS Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.101.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alarm_contact_groups" "example" {
  name_regex = "tf-testacc"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew)  A list of Alarm Contact Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Alarm Contact Group name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of CMS Alarm Contact Group names.
* `groups` - A list of CMS groups. Each element contains the following attributes:
	* `id` - The ID of the CMS.
	* `alarm_contact_group_name` - The name of Alarm Contact Group.
	* `contacts` - The alarm contacts in the alarm group.
	* `describe` - The description of the Alarm Group.
	* `enable_subscribed` - Indicates whether the alarm group subscribes to weekly reports. 

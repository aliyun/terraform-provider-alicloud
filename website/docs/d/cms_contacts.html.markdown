---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_contacts"
sidebar_current: "docs-alicloud-data-source-cms-contacts"
description: |-
  Provides a list of Cloud Monitor Service (CMS) Contacts.
---

# alicloud_cms_contacts

This data source provides a list of Cloud Monitor Service (CMS) Contacts in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_contacts" "default" {
  contact_name = "tf-example"
}

output "first_contact_id" {
  value = data.alicloud_cms_contacts.default.contacts.0.contact_id
}
```

## Argument Reference

The following arguments are supported:

* `contact_name` - (Optional) The name of the alarm contact used to filter results.
* `phone` - (Optional) The phone number of the alarm contact used to filter results.
* `email` - (Optional) The email address of the alarm contact used to filter results.
* `query_ungrouped_contacts` - (Optional) Whether to query ungrouped contacts.
* `name_regex` - (Optional) A regex string to filter results by contact name.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of contact IDs.
* `contacts` - A list of CMS Contacts. Each element contains the following attributes:
  * `contact_id` - The ID of the alarm contact.
  * `contact_name` - The name of the alarm contact.
  * `email` - The email address of the alarm contact.
  * `phone` - The phone number of the alarm contact.
  * `lang` - The language type of the alarm.
  * `workspace` - The workspace ID.
  * `im_user_ids` - The DingTalk and other communication tools user IDs.

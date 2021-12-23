---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_contacts"
sidebar_current: "docs-alicloud-datasource-msc-sub-contacts"
description: |-
    Provides a list of Message Center Contacts to the user.
---

# alicloud\_msc\_sub\_contacts

This data source provides the Message Center Contacts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_msc_sub_contacts" "ids" {}
output "msc_sub_contact_id_1" {
  value = data.alicloud_msc_sub_contacts.ids.contacts.0.id
}

data "alicloud_msc_sub_contacts" "nameRegex" {
  name_regex = "^my-Contact"
}
output "msc_sub_contact_id_2" {
  value = data.alicloud_msc_sub_contacts.nameRegex.contacts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Contact IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Contact name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Contact names.
* `contacts` - A list of Msc Sub Contacts. Each element contains the following attributes:
    * `account_uid` - UID.
    * `contact_id` - The first ID of the resource.
    * `contact_name` - The User's Contact Name. **Note:** The name must be 2 to 12 characters in length, and can contain uppercase and lowercase letters.
    * `email` - The User's Contact Email Address.
    * `id` - The ID of the Contact.
    * `is_account` - Indicates Whether the BGP Group Is the Account Itself.
    * `is_obsolete` - Whether They Have Expired Or Not.
    * `is_verified_email` - Email Validation for.
    * `is_verified_mobile` - If the Phone Verification.
    * `last_email_verification_time_stamp` - Last Verification Email Transmission Time.
    * `last_mobile_verification_time_stamp` - The Pieces of Authentication SMS Sending Time.
    * `mobile` - The User's Telephone.
    * `position` - The User's Position. Valid values: `CEO`, `Technical Director`, `Maintenance Director`, `Project Director`,`Finance Director` and `Other`.

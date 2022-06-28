---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_contact"
sidebar_current: "docs-alicloud-resource-msc-sub-contact"
description: |-
  Provides a Alicloud Message Center Contact resource.
---

# alicloud\_msc\_sub\_contact

Provides a Msc Sub Contact resource.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_msc_sub_contact" "default" {
  contact_name = example_value
  position     = "CEO"
  email        = "123@163.com"
  mobile       = "153xxxxx906"
}
```

## Argument Reference

The following arguments are supported:

* `contact_name` - (Required) The User's Contact Name. **Note:** The name must be 2 to 12 characters in length.
* `email` - (Required) The User's Contact Email Address.
* `mobile` - (Required) The User's Telephone.
* `position` - (Required, ForceNew) The User's Position. Valid values: `CEO`, `Technical Director`, `Maintenance Director`, `Project Director`,`Finance Director` and `Other`.

-> **NOTE:** When the user creates a contact, the user should use `alicloud_msc_sub_contact_verification_message` to receive the verification message and confirm it.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Contact.

## Import

Msc Sub Contact can be imported using the id, e.g.

```
$ terraform import alicloud_msc_sub_contact.example <id>
```

---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_contact_verification_message"
sidebar_current: "docs-alicloud-datasource-msc-sub-contact-verification-message"
description: |-
    Provide a data source to send the verification message to the user.
---

# alicloud\_msc\_sub\_contact\_verification\_message


-> **NOTE:** Available in v1.156.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_msc_sub_contact" "default" {
  contact_name = "example_value"
  position     = "CEO"
  email        = "123@163.com"
  mobile       = "153xxxxx906"
}

data "alicloud_msc_sub_contact_verification_message" "default" {
  contact_id = alicloud_msc_sub_contact.default.id
  type       = 1
}
```

## Argument Reference

The following arguments are supported:

* `contact_id` - (Required, ForceNew)  The ID of the Contact.
* `type` - (Required, ForceNew) How a user receives verification messages. Valid values : `1`, `2`.
  * `1`: Send a verification message through the user's mobile.
  * `2`: Send a verification message through the user's mail.

  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The sending status of the message. Valid values : `Success`, `Failed`.


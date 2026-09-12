---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_contact"
sidebar_current: "docs-alicloud-resource-cms-contact"
description: |-
  Provides a Cloud Monitor Service (CMS) Contact resource.
---

# alicloud_cms_contact

Provides a Cloud Monitor Service (CMS) Contact resource.

For information about CMS Contact and how to use it, see [What is CMS Contact](https://www.alibabacloud.com/help/en/cloudmonitor/).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_contact" "default" {
  contact_name = "tf-example"
  email        = "terraform@test.com"
  phone        = "86-13800000000"
  lang         = "zh_CN"
  im_user_ids = {
    key1 = "user1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `contact_name` - (Required) The name of the alarm contact.
* `email` - (Optional) The email address of the alarm contact. After you add or modify an email address, the recipient receives an email that contains an activation link. The system adds the recipient to the list of alarm contacts only after the recipient activates the email address.
* `phone` - (Optional) The phone number of the alarm contact. After you add or modify a phone number, the recipient receives a message that contains an activation link.
* `lang` - (Optional) The language type of the alarm. Valid values: `zh_CN`, `en_US`.
* `im_user_ids` - (Optional) The DingTalk and other communication tools user IDs. It is a map of string.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alarm contact. It is the same as `contact_id`.
* `contact_id` - The ID of the alarm contact.
* `region_id` - The region ID of the resource.
* `workspace` - The workspace ID.

## Import

CMS Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_contact.example abc12345
```

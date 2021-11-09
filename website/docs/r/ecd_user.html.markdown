---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_user"
sidebar_current: "docs-alicloud-resource-ecd-user"
description: |-
  Provides a Alicloud Elastic Desktop Service(EDS) User resource.
---

# alicloud\_ecd\_user

Provides a Elastic Desktop Service(EDS) User resource.

For information about Elastic Desktop Service(EDS) User and how to use it, see [What is User](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_user" "example" {
  email       = "your_email"
  end_user_id = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `email` - (Required, ForceNew) The email of the user email.
* `end_user_id` - (Required, ForceNew) The Username. The custom setting is composed of lowercase letters, numbers and underscores, and the length is 3~24 characters.
* `password` - (Optional, ForceNew) The password of the user password.
* `phone` - (Optional, ForceNew) The phone of the mobile phone number.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Unlocked`, `Locked`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User. Its value is same as `end_user_id`.

## Import

ECD User can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_user.example <end_user_id>
```

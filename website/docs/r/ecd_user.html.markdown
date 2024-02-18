---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_user"
sidebar_current: "docs-alicloud-resource-ecd-user"
description: |-
  Provides a Alicloud Elastic Desktop Service (ECD) User resource.
---

# alicloud_ecd_user

Provides a Elastic Desktop Service (ECD) User resource.

For information about Elastic Desktop Service (ECD) User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-eds-user-2021-03-08-createusers-desktop).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_ecd_user" "default" {
  end_user_id = "terraform_example123"
  email       = "tf.example@abc.com"
  phone       = "18888888888"
  password    = "Example_123"
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

```shell
$ terraform import alicloud_ecd_user.example <end_user_id>
```

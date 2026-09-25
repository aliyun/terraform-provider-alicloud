---
subcategory: "Hologres (Hologram)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hologram_user"
description: |-
  Provides a Hologres (Hologram) User resource.
---

# alicloud_hologram_user

Provides a Hologres (Hologram) User resource.

For information about Hologres (Hologram) User and how to use it, see [What is Hologram User](https://www.alibabacloud.com/help/zh/hologres/developer-reference/api-hologram-2022-06-01-createuser).

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_hologram_instance" "default" {
  instance_type = "starter"
  instance_name = var.name
}

resource "alicloud_hologram_user" "default" {
  instance_id = alicloud_hologram_instance.default.id
  user_name   = "tf_example_user"
  super_user  = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Hologres instance.
* `user_name` - (Required, ForceNew) The name of the instance user. It must start with a letter and can contain lowercase letters, digits, and underscores (_).
* `super_user` - (Optional, ForceNew) Specifies whether the user is a super user.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID of Hologram User. The value is formulated as `<instance_id>:<user_name>`.

## Import

Hologram User can be imported using the id, e.g.

```shell
$ terraform import alicloud_hologram_user.example <instance_id>:<user_name>
```

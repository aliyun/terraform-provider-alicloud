---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_user_group_attachment"
description: |-
  Provides a Alicloud RAM User Group Attachment resource.
---

# alicloud_ram_user_group_attachment

Provides a RAM User Group Attachment resource.



For information about RAM User Group Attachment and how to use it, see [What is User Group Attachment](https://next.api.alibabacloud.com/document/Ram/2015-05-01/AddUserToGroup).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_user" "default" {
  name         = "terraform-example-${random_integer.default.result}"
  display_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_ram_group" "default" {
  name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_ram_user_group_attachment" "default" {
  group_name = alicloud_ram_group.default.id
  user_name  = alicloud_ram_user.default.name
}
```

## Argument Reference

The following arguments are supported:
* `group_name` - (Required, ForceNew) The user group name.
* `user_name` - (Required, ForceNew) The user name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<group_name>:<user_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the User Group Attachment.

## Import

RAM User Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_user_group_attachment.example <group_name>:<user_name>
```
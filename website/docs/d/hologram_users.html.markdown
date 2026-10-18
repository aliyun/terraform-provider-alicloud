---
subcategory: "Hologres (Hologram)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hologram_users"
sidebar_current: "docs-alicloud-datasource-hologram-users"
description: |-
  Provides a list of Hologres (Hologram) User owned by an Alibaba Cloud account.
---

# alicloud_hologram_users

This data source provides a list of Hologres (Hologram) Users owned by an Alibaba Cloud account.

For information about Hologres (Hologram) User, see [What is Hologram User](https://www.alibabacloud.com/help/zh/hologres/developer-reference/api-hologram-2022-06-01-createuser).

-> **NOTE:** Available since v1.295.0.

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

data "alicloud_hologram_users" "default" {
  instance_id = alicloud_hologram_instance.default.id
}

output "first_user_name" {
  value = data.alicloud_hologram_users.default.users.0.user_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the Hologres instance.
* `ids` - (Optional) A list of User IDs. The value is formulated as `<instance_id>:<user_name>`.
* `name_regex` - (Optional) A regex string to filter results by user name.
* `output_file` - (Optional) File name where to save the data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of User IDs.
* `names` - A list of user names.
* `users` - A list of Hologram Users. Each element contains the following attributes:
  * `id` - The ID of the user, `<instance_id>:<user_name>`.
  * `user_name` - The name of the instance user.
  * `super_user` - Whether the user is a super user. Always null: ListUsers does not return this field.

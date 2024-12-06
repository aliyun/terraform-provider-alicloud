---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_project_member"
description: |-
  Provides a Alicloud Data Works Project Member resource.
---

# alicloud_data_works_project_member

Provides a Data Works Project Member resource.



For information about Data Works Project Member and how to use it, see [What is Project Member](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
variable "admin_code" {
  default = "role_project_admin"
}

variable "name" {
  default = "tf_example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

resource "alicloud_ram_user" "defaultKCTrB2" {
  display_name = "${var.name}${random_integer.randint.id}"
  name         = "${var.name}${random_integer.randint.id}"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "defaultQeRfvU" {
  status                  = "Available"
  description             = "tf_desc"
  project_name            = "${var.name}${random_integer.randint.id}"
  pai_task_enabled        = "false"
  display_name            = "tf_new_api_display"
  dev_role_disabled       = "true"
  dev_environment_enabled = "false"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_data_works_project_member" "default" {
  user_id    = alicloud_ram_user.defaultKCTrB2.id
  project_id = alicloud_data_works_project.defaultCoMnk8.id
  roles {
    code = var.admin_code
  }
}
```

## Argument Reference

The following arguments are supported:
* `project_id` - (Required, ForceNew, Int) Project ID
* `roles` - (Optional, List) List of roles owned by members. See [`roles`](#roles) below.
* `user_id` - (Required, ForceNew) The user ID of the member.

### `roles`

The roles supports the following:
* `code` - (Optional) Project Role Code.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_id>:<user_id>`.
* `roles` - List of roles owned by members.
  * `name` - project role name
  * `type` - project role type
* `status` - The status of the user in project

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project Member.
* `delete` - (Defaults to 5 mins) Used when delete the Project Member.
* `update` - (Defaults to 5 mins) Used when update the Project Member.

## Import

Data Works Project Member can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_project_member.example <project_id>:<user_id>
```
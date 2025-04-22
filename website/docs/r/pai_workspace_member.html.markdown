---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_member"
description: |-
  Provides a Alicloud PAI Workspace Member resource.
---

# alicloud_pai_workspace_member

Provides a PAI Workspace Member resource.



For information about PAI Workspace Member and how to use it, see [What is Member](https://www.alibabacloud.com/help/en/pai/developer-reference/api-aiworkspace-2021-02-04-createmember).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_pai_workspace_workspace" "Workspace" {
  description    = "156"
  display_name   = var.name
  workspace_name = "${var.name}_${random_integer.default.result}"
  env_types      = ["prod"]
}

resource "alicloud_ram_user" "default" {
  name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_pai_workspace_member" "default" {
  user_id      = alicloud_ram_user.default.id
  workspace_id = alicloud_pai_workspace_workspace.Workspace.id
  roles        = ["PAI.AlgoDeveloper", "PAI.AlgoOperator", "PAI.LabelManager"]
}
```

## Argument Reference

The following arguments are supported:
* `roles` - (Required, List) The list of roles. see [how to use it](https://www.alibabacloud.com/help/en/pai/developer-reference/api-aiworkspace-2021-02-04-createmember).
* `user_id` - (Required, ForceNew) The ID of the User.
* `workspace_id` - (Required, ForceNew) The ID of the Workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<workspace_id>:<member_id>`.
* `create_time` - The time when the workspace is created, in UTC. The time follows the ISO 8601 standard.
* `member_id` - The member ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Member.
* `delete` - (Defaults to 5 mins) Used when delete the Member.
* `update` - (Defaults to 5 mins) Used when update the Member.

## Import

PAI Workspace Member can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_member.example <workspace_id>:<member_id>
```

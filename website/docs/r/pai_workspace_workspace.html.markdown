---
subcategory: "P A I Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_workspace"
description: |-
  Provides a Alicloud P A I Workspace Workspace resource.
---

# alicloud_pai_workspace_workspace

Provides a P A I Workspace Workspace resource.



For information about P A I Workspace Workspace and how to use it, see [What is Workspace](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.233.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_pai_workspace_workspace" "default" {
  description    = var.name
  workspace_name = var.name
  display_name   = var.name
  env_types      = ["prod"]
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Required) Workspace description, no more than 80 characters.
* `display_name` - (Optional) It is recommended that you name the workspace based on the business attribute to identify the purpose of the workspace. If not configured, the default value is the workspace name.
* `env_types` - (Required, ForceNew, List) Environments contained in the workspace:
  - Simple mode only production environment (prod).
  - Standard mode includes development environment (dev) and production environment (prod).
* `workspace_name` - (Required, ForceNew) The workspace name. The format is as follows:
  - 3 to 23 characters in length and can contain letters, underscores, or numbers.
  - Must start with a large or small letter.
  - Unique in the current region.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The UTC time when the workspace is created. The time format is ISO8601.
* `status` - Workspace state, possible values:

  ENABLED: normal.

  INITIALIZING: INITIALIZING.

  FAILURE: failed.

  DISABLED: manually DISABLED.

  FROZEN: Arrears are FROZEN.

  UPDATING: The project is being updated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Workspace.
* `delete` - (Defaults to 5 mins) Used when delete the Workspace.
* `update` - (Defaults to 5 mins) Used when update the Workspace.

## Import

P A I Workspace Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_workspace.example <id>
```
---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_default_workspace"
description: |-
  Provides a Alicloud PAI Workspace Default Workspace resource.
---

# alicloud_pai_workspace_default_workspace

Provides a PAI Workspace Default Workspace resource.

Default Workspace Resources.

For information about PAI Workspace Default Workspace and how to use it, see [What is Default Workspace](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_pai_workspace_default_workspace" "default" {
  description = "defaultWorkspace"
  env_types   = ["prod"]
}
```

### Deleting `alicloud_pai_workspace_default_workspace` or removing it from your configuration

Terraform cannot destroy resource `alicloud_pai_workspace_default_workspace`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `description` - (Required, ForceNew) Description of resource description
* `env_types` - (Required, ForceNew, List) Resource description for environment type list
* `workspace_id` - (Optional, Computed) Resource Description of the workspace ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.
* `status` - Resource description of workspace state

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Default Workspace.
* `update` - (Defaults to 5 mins) Used when update the Default Workspace.

## Import

PAI Workspace Default Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_default_workspace.example 
```
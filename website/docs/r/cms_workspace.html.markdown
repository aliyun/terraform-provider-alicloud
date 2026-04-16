---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_workspace"
description: |-
  Provides a Alicloud Cms Workspace resource.
---

# alicloud_cms_workspace

Provides a Cms Workspace resource.



For information about Cms Workspace and how to use it, see [What is Workspace](https://next.api.alibabacloud.com/document/Cms/2024-03-30/PutWorkspace).

-> **NOTE:** Available since v1.276.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the workspace.
* `display_name` - (Optional) The dispalyName of the workspace.
* `sls_project` - (Required, ForceNew) The project bind to workspace.
* `workspace_name` - (Required, ForceNew) The name of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the workspace.
* `region_id` - The region of the workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Workspace.
* `delete` - (Defaults to 5 mins) Used when delete the Workspace.
* `update` - (Defaults to 5 mins) Used when update the Workspace.

## Import

Cms Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_workspace.example <workspace_name>
```

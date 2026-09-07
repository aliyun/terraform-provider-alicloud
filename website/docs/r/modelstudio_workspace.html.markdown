---
subcategory: "Model Studio"
layout: "alicloud"
page_title: "Alicloud: alicloud_modelstudio_workspace"
description: |-
  Provides a Alicloud Model Studio workspace resource.
---

# alicloud_modelstudio_workspace

Provides a Model Studio workspace resource.

For information about Model Studio workspace and how to use it, see [Create a workspace](https://www.alibabacloud.com/help/en/model-studio/developer-reference/create-workspace).

-> **NOTE:** Available since v1.244.0.

## Example Usage

```terraform
resource "alicloud_modelstudio_workspace" "default" {
  workspace_name = "tf-example"
  service_site   = "global"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_name` - (Required) The name of the business workspace. The name must be 1 to 30 characters in length.
* `service_site` - (Optional, ForceNew, Computed) The service deployment scope. Services that rely on GPU resources (such as inference, training, and deployment) are deployed in this scope.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the workspace.
* `workspace_id` - The ID of the workspace.
* `create_time` - The creation time of the workspace.
* `region_id` - The region ID of the workspace.
* `api_host` - The access endpoint of the workspace, which is a dedicated isolated access point for the workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the workspace.
* `update` - (Defaults to 5 mins) Used when updating the workspace.
* `delete` - (Defaults to 5 mins) Used when deleting the workspace.

## Import

Model Studio workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_modelstudio_workspace.example <workspace-id>
```

---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_default_workspace"
description: |-
  Provides a Alicloud PAI Workspace default workspace resource.
---

# alicloud_pai_workspace_default_workspace

Provides a PAI Workspace default workspace resource. The default workspace is a singleton per region per account: it is created once, cannot be deleted, and its `description` and `env_types` are immutable after creation (changing either forces replacement).

-> **NOTE:** Available since v1.295.0.

-> **NOTE:** Because the default workspace is a singleton that cannot be deleted, forcing a replacement (by changing `description` or `env_types`) leaves the existing cloud resource in place; Terraform will keep reporting the drift until the configuration matches the cloud again.

## Example Usage

```terraform
resource "alicloud_pai_workspace_default_workspace" "default" {
  description = "terraform-default-workspace"
  env_types   = ["prod"]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, ForceNew) The description of the default workspace, up to 80 characters. Immutable after creation.
* `env_types` - (Required, ForceNew) The environments contained in the workspace. Simple mode only contains `prod`; standard mode contains `dev` and `prod`. Immutable after creation.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `workspace_id`.
* `workspace_id` - The ID of the default workspace.
* `workspace_name` - The name of the workspace, unique within the region.
* `display_name` - The display name of the workspace.
* `creator` - The creator of the workspace.
* `status` - The status of the workspace. Valid values: `INITIALIZING`, `ENABLED`, `FAILURE`, `DISABLE`, `FROZEN`, `UPDATING`.
* `region_id` - The region ID of the workspace.
* `gmt_create_time` - The UTC creation time of the workspace, in ISO 8601 format.
* `gmt_modified_time` - The UTC modification time of the workspace, in ISO 8601 format.
* `owner_user_id` - The user ID of the workspace owner.
* `owner_user_name` - The user name of the workspace owner.
* `owner_user_kp` - The user KP of the workspace owner.

## Import

PAI Workspace default workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_default_workspace.default <workspace_id>
```

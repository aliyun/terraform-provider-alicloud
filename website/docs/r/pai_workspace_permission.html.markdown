---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_permission"
sidebar_current: "docs-alicloud-resource-pai-workspace-permission"
description: |-
  Provides a PAI Workspace Permission resource.
---

# alicloud_pai_workspace_permission

Provides a PAI Workspace Permission resource. PAI Workspace permissions are
read-only catalog entries surfaced by the `GetPermission` and `ListPermissions`
APIs of the **AIWorkSpace 2021-02-04** product. A permission point is scoped to
a workspace and identified by a `permission_code` (e.g. `PaiDLC:ListJobs`).

This resource is **read-only**: it has no create, update or delete lifecycle
against the cloud. Applying it imports an existing permission point into
Terraform state; destroying it removes the resource from state only, leaving
the catalog entry untouched.

For information about PAI Workspace Permission and how to use it, see
[GetPermission](https://www.alibabacloud.com/help/en/pai/user-guide/developer-reference/api-aiworkspace-2021-02-04-getpermission).

-> **NOTE:** Available since v2.0.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_workspace_workspace" "example" {
  workspace_name = "terraform-example"
  display_name   = "Terraform Example"
  description    = "Managed by Terraform"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_permission" "example" {
  workspace_id    = alicloud_pai_workspace_workspace.example.id
  permission_code = "PaiDLC:ListJobs"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The ID of the PAI Workspace. How to get
  the workspace ID, see ListWorkspaces.
* `permission_code` - (Required, ForceNew) The name of the permission point,
  unique within the region, e.g. `PaiDLC:ListJobs`. Permission codes contain a
  colon and are split on the first colon of the resource id. How to get the
  permission code, see ListPermissions.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<workspace_id>:<permission_code>`.
  Because permission codes contain a colon, the id is split on the first colon
  only.
* `permission_rule_accessibilities` - The Accessibility value of each rule of
  the permission point, in catalog order. Valid values: `PUBLIC`, `PRIVATE`,
  `ANY`. Index-corresponds with `permission_rule_entity_access_types`.
* `permission_rule_entity_access_types` - The EntityAccessType value of each
  rule of the permission point, in catalog order. Valid values: `CREATOR`,
  `ANY`. Index-corresponds with `permission_rule_accessibilities`.

-> **NOTE:** The permission rules are flattened into two parallel list
  attributes rather than a nested `permission_rules` block because this resource
  is served over protocol v5 by a muxed provider server; the nested-attribute
  representation is rejected by the v5 schema conversion. Pair rules across the
  two lists by index.

## Import

PAI Workspace Permission can be imported using the formatted id
`<workspace_id>:<permission_code>`, e.g.

```shell
$ terraform import alicloud_pai_workspace_permission.example <workspace_id>:PaiDLC:ListJobs
```
